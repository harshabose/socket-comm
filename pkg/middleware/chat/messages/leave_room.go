package messages

import (
	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

var LeaveRoomProtocol message.Protocol = "room:leave_room"

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
