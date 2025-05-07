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

var JoinRoomProtocol message.Protocol = "room:join_room"

type JoinRoom struct {
	interceptor.BaseMessage
	RoomID       types.RoomID  `json:"room_id"`
	JoinDeadline time.Duration `json:"join_deadline"`
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

	timer := time.NewTimer(m.JoinDeadline)
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

	// TODO: AFTER SUCCESS, SEND ROOM CURRENT STATE TO THE CLIENT
	// TODO: THEN SEND SuccessJoinRoom MESSAGE TO THE CLIENT
}
