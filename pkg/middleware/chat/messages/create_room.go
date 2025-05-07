package messages

import (
	"time"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

var CreateRoomProtocol message.Protocol = "room:create_room"

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
