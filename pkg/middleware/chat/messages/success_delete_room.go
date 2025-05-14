package messages

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

var SuccessDeleteRoomProtocol message.Protocol = "room:success_delete_room"

type SuccessDeleteRoom struct {
	interceptor.BaseMessage
	RoomID types.RoomID `json:"room_id"`
}

func NewSuccessDeleteRoomMessage(id types.RoomID) (*SuccessDeleteRoom, error) {
	msg := &SuccessDeleteRoom{
		RoomID: id,
	}

	bmsg, err := interceptor.NewBaseMessage(message.NoneProtocol, nil, msg)
	if err != nil {
		return nil, err
	}

	msg.BaseMessage = bmsg
	return msg, nil
}

func NewSuccessDeleteRoomMessageFactory(id types.RoomID) func() (message.Message, error) {
	return func() (message.Message, error) {
		return NewSuccessDeleteRoomMessage(id)
	}
}

func (m *SuccessDeleteRoom) GetProtocol() message.Protocol {
	return SuccessDeleteRoomProtocol
}

func (m *SuccessDeleteRoom) ReadProcess(_ context.Context, _i interceptor.Interceptor, _ interceptor.Connection) error {
	_, ok := _i.(*chat.ClientInterceptor)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	// NOTE: INTENTIONALLY EMPTY

	return nil
}
