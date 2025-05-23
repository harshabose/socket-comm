package messages

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

const SuccessLeaveRoomProtocol message.Protocol = "room:success_leave_room"

type SuccessLeaveRoom struct {
	interceptor.BaseMessage
	RoomID   types.RoomID         `json:"room_id"`
	ClientID interceptor.ClientID `json:"client_id"`
}

func NewSuccessLeaveRoomMessage(id types.RoomID, clientID interceptor.ClientID) (*SuccessLeaveRoom, error) {
	msg := &SuccessLeaveRoom{
		RoomID:   id,
		ClientID: clientID,
	}

	bmsg, err := interceptor.NewBaseMessage(message.NoneProtocol, nil, msg)
	if err != nil {
		return nil, err
	}

	msg.BaseMessage = bmsg
	return msg, nil
}

func NewSuccessLeaveRoomMessageFactory(id types.RoomID, clientID interceptor.ClientID) func() (message.Message, error) {
	return func() (message.Message, error) {
		return NewSuccessLeaveRoomMessage(id, clientID)
	}
}

func (m *SuccessLeaveRoom) GetProtocol() message.Protocol {
	return SuccessLeaveRoomProtocol
}

func (m *SuccessLeaveRoom) ReadProcess(_ context.Context, _i interceptor.Interceptor, _ interceptor.Connection) error {
	_, ok := _i.(*chat.ClientInterceptor)
	if !ok {
		return interceptor.ErrInterfaceMisMatch
	}

	// NOTE: INTENTIONALLY EMPTY
	return nil
}
