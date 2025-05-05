package room

import (
	"fmt"
	"time"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
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
	s, ok := _i.(interfaces.CanGetState)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	i, ok := _i.(interfaces.CanCreateRoom)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	ss, err := s.GetState(connection)
	if err != nil {
		return fmt.Errorf("error while read processing 'CreateRoom' msg; err: %s", err.Error())
	}

	if !ss.GetConfig().IsServer {
		return fmt.Errorf("error while read processing 'CreateRoom' msg; err: %s", errors.ErrMessageForServerOnly)
	}

	room, err := i.CreateRoom(m.RoomID, m.Allowed, m.TTL)
	if err != nil {
		return fmt.Errorf("error while read processing 'CreateRoom' msg; err: %s", err.Error())
	}

	if err := room.Add(m.RoomID, ss); err != nil {
		return fmt.Errorf("error while read processing 'CreateRoom' msg; err: %s", err.Error())
	}

	return nil
}

type DeleteRoom struct {
	interceptor.BaseMessage
	RoomID types.RoomID `json:"room_id"`
}

func (m *DeleteRoom) GetProtocol() message.Protocol {
	return DeleteRoomProtocol
}

func (m *DeleteRoom) ReadProcess(_i interceptor.Interceptor, connection interceptor.Connection) error {
	s, ok := _i.(interfaces.CanGetState)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	i, ok := _i.(interfaces.CanDeleteRoom)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	ss, err := s.GetState(connection)
	if err != nil {
		return fmt.Errorf("error while read processing 'DeleteRoom' msg; err: %s", err.Error())
	}

	if !ss.GetConfig().IsServer {
		return fmt.Errorf("error while read processing 'DeleteRoom' msg; err: %s", errors.ErrMessageForServerOnly)
	}

	if err := i.DeleteRoom(m.RoomID); err != nil {
		return fmt.Errorf("error while read processing 'DeleteRoom' msg; err: %s", errors.ErrMessageForServerOnly)
	}

	return nil
}

type ToForwardMessage struct {
	interceptor.BaseMessage
	RoomID types.RoomID   `json:"room_id"`
	From   types.ClientID `json:"from"`
	To     types.ClientID `json:"to"`
}

func (m *ToForwardMessage) GetProtocol() message.Protocol {
	return ForwardMessageProtocol
}

func (m *ToForwardMessage) ReadProcess(_i interceptor.Interceptor, connection interceptor.Connection) error {
	s, ok := _i.(interfaces.CanGetState)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	msg, err := newForwardedMessage(m)
	if err != nil {
		return fmt.Errorf("error while read processing 'ToForwardMessage'; err: %s", err.Error())
	}

	ss, err := s.GetState(connection)
	if err != nil {
		return fmt.Errorf("error while read processing 'ToForwardMessage' msg; err: %s", err.Error())
	}

	if !ss.GetConfig().IsServer {
		return fmt.Errorf("error while read processing 'ToForwardMessage' msg; err: %s", errors.ErrMessageForServerOnly)
	}

	w, ok := _i.(interfaces.CanWriteRoomMessage)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	if err := w.WriteMessage(m.RoomID, msg, m.From, m.To); err != nil {
		return fmt.Errorf("error while read processing 'ToForwardMessage' msg; err: %s", err.Error())
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
	RoomID types.RoomID
}

func (m *JoinRoom) GetProtocol() message.Protocol {
	return JoinRoomProtocol
}

func (m *JoinRoom) ReadProcess(_i interceptor.Interceptor, connection interceptor.Connection) error {
	s, ok := _i.(interfaces.CanGetState)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	ss, err := s.GetState(connection)
	if err != nil {
		return fmt.Errorf("error while read processing 'JoinRoom' msg; err: %s", err.Error())
	}

	a, ok := _i.(interfaces.CanAdd)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	if err := a.Add(m.RoomID, ss); err != nil {
		return fmt.Errorf("error while read processing 'JoinRoom' msg; err: %s", err.Error())
	}

	return nil
}

type LeaveRoom struct {
	interceptor.BaseMessage
	RoomID types.RoomID
}

func (m *LeaveRoom) GetProtocol() message.Protocol {
	return LeaveRoomProtocol
}

func (m *LeaveRoom) ReadProcess(_i interceptor.Interceptor, connection interceptor.Connection) error {
	s, ok := _i.(interfaces.CanGetState)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	ss, err := s.GetState(connection)
	if err != nil {
		return fmt.Errorf("error while read processing 'JoinRoom' msg; err: %s", err.Error())
	}

	a, ok := _i.(interfaces.CanRemove)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	if err := a.Remove(m.RoomID, ss); err != nil {
		return fmt.Errorf("error while read processing 'JoinRoom' msg; err: %s", err.Error())
	}

	return nil
}
