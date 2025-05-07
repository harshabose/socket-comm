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

var DeleteRoomProtocol message.Protocol = "room:delete_room"

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
