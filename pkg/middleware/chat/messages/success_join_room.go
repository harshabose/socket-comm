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

var SuccessJoinRoomProtocol message.Protocol = "room:success_join_room"

// SuccessJoinRoom is the message sent by the server to the clients (including the requested client and roommates)
// when the client joins a room successfully.
type SuccessJoinRoom struct {
	interceptor.BaseMessage
	RoomID   types.RoomID   `json:"room_id"`
	ClientID types.ClientID `json:"client_id"`
}

func (m *SuccessJoinRoom) GetProtocol() message.Protocol {
	return SuccessJoinRoomProtocol
}

func (m *SuccessJoinRoom) ReadProcess(_i interceptor.Interceptor, _ interceptor.Connection) error {
	i, ok := _i.(*chat.ClientInterceptor)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	return i.Health.Process(m, nil)
}

func (m *SuccessJoinRoom) Process(p interfaces.Processor, _ *state.State) error {
	a, ok := p.(interfaces.CanAddHealth)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}
	// NOTE: MIGHT FAIL IF THE ROOM CREATION MESSAGE IS NOT RECEIVED BY THE CLIENT BEFORE THIS MESSAGE IS SENT.
	return a.Add(m.RoomID, m.ClientID)
}
