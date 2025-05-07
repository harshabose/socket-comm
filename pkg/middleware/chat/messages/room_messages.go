package messages

import (
	"fmt"
	"time"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

var (
	CreateRoomProtocol       message.Protocol = "room:create_room"
	DeleteRoomProtocol       message.Protocol = "room:delete_room"
	ForwardMessageProtocol   message.Protocol = "room:forward_message"
	ForwardedMessageProtocol message.Protocol = "room:forwarded_message"
	JoinRoomProtocol         message.Protocol = "room:join_room"
	LeaveRoomProtocol        message.Protocol = "room:leave_room"
)

type CreateRoom struct {
	interceptor.BaseMessage
	RoomID  types.RoomID     `json:"room_id"`
	Allowed []types.ClientID `json:"allowed"`
	TTL     time.Duration    `json:"ttl"`
}

func (m *CreateRoom) GetProtocol() message.Protocol {
	return CreateRoomProtocol
}

func (m *CreateRoom) ReadProcess(_i interceptor.Interceptor, connection interceptor.Connection) error {
	i, ok := _i.(*chat.ServerInterceptor)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	s, err := i.GetState(connection)
	if err != nil {
		return err
	}

	return i.Rooms.Process(m, s)
}

func (m *CreateRoom) Process(p interfaces.Processor, s *state.State) error {
	r, ok := p.(interfaces.CanCreateRoom)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	room, err := r.CreateRoom(m.RoomID, m.Allowed, m.TTL)
	if err != nil {
		return err
	}

	return room.Add(m.RoomID, s)
}

type DeleteRoom struct {
	interceptor.BaseMessage
	RoomID types.RoomID `json:"room_id"`
}

func (m *DeleteRoom) GetProtocol() message.Protocol {
	return DeleteRoomProtocol
}

func (m *DeleteRoom) ReadProcess(_i interceptor.Interceptor, _ interceptor.Connection) error {
	i, ok := _i.(*chat.ServerInterceptor)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	return i.Rooms.Process(m, nil)
}

func (m *DeleteRoom) Process(p interfaces.Processor, _ *state.State) error {
	r, ok := p.(interfaces.CanDeleteRoom)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	return r.DeleteRoom(m.RoomID)
}

type ToForwardMessage struct {
	interceptor.BaseMessage
	RoomID types.RoomID     `json:"room_id"`
	To     []types.ClientID `json:"to"`
}

func (m *ToForwardMessage) GetProtocol() message.Protocol {
	return ForwardMessageProtocol
}

func (m *ToForwardMessage) ReadProcess(_i interceptor.Interceptor, connection interceptor.Connection) error {
	i, ok := _i.(*chat.ServerInterceptor)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}
	s, err := i.GetState(connection)
	if err != nil {
		return err
	}

	return i.Rooms.Process(m, s)
}

func (m *ToForwardMessage) Process(p interfaces.Processor, s *state.State) error {
	msg, err := newForwardedMessage(m)
	if err != nil {
		return err
	}

	clientID, err := s.GetClientID()
	if err != nil {
		return err
	}

	if types.ClientID(m.CurrentHeader.Sender) != clientID {
		return fmt.Errorf("error while read processing 'ToForwardMessage'; From and ClientID did not match")
	}

	w, ok := p.(interfaces.CanWriteRoomMessage)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	if err := w.WriteRoomMessage(m.RoomID, msg, clientID, m.To...); err != nil {
		return err
	}

	return nil
}

type ForwardedMessage struct {
	interceptor.BaseMessage
}

func (m *ForwardedMessage) GetProtocol() message.Protocol {
	return ForwardedMessageProtocol
}

func newForwardedMessage(forward *ToForwardMessage) (*ForwardedMessage, error) {
	msg := &ForwardedMessage{}
	bmsg, err := interceptor.NewBaseMessage(forward.GetNextProtocol(), forward.NextPayload, msg)
	if err != nil {
		return nil, err
	}

	msg.BaseMessage = bmsg

	return msg, nil
}

type JoinRoom struct {
	interceptor.BaseMessage
	RoomID       types.RoomID  `json:"room_id"`
	WaitDuration time.Duration `json:"wait_duration"`
}

func (m *JoinRoom) GetProtocol() message.Protocol {
	return JoinRoomProtocol
}

func (m *JoinRoom) ReadProcess(_i interceptor.Interceptor, connection interceptor.Connection) error {
	i, ok := _i.(*chat.ServerInterceptor)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	s, err := i.GetState(connection)
	if err != nil {
		return err
	}

	return i.Rooms.Process(m, s)
}

func (m *JoinRoom) Process(p interfaces.Processor, s *state.State) error {
	a, ok := p.(interfaces.CanAdd)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	timer := time.NewTimer(m.WaitDuration)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			return fmt.Errorf("error while read processing 'JoinRoom' msg; err: %s", errors.ErrContextCancelled)
		default:
			err := a.Add(m.RoomID, s)
			if err == nil {
				return nil
			}
			fmt.Println(fmt.Errorf("error while read processing 'JoinRoom' msg; err: %s. retrying", err.Error()))
		}
	}
}

type LeaveRoom struct {
	interceptor.BaseMessage
	RoomID types.RoomID
}

func (m *LeaveRoom) GetProtocol() message.Protocol {
	return LeaveRoomProtocol
}

func (m *LeaveRoom) ReadProcess(_i interceptor.Interceptor, connection interceptor.Connection) error {
	i, ok := _i.(*chat.ServerInterceptor)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	s, err := i.GetState(connection)
	if err != nil {
		return err
	}

	return i.Rooms.Process(m, s)
}

func (m *LeaveRoom) Process(p interfaces.Processor, s *state.State) error {
	r, ok := p.(interfaces.CanRemove)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	return r.Remove(m.RoomID, s)
}
