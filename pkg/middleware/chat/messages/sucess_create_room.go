package messages

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

var SuccessCreateRoomProtocol message.Protocol = "room:success_create_room"

// SuccessCreateRoom is the message sent by the server to the client after successful creation of the requested room.
// This marks the end of the CreateRoom topic.
type SuccessCreateRoom struct {
	interceptor.BaseMessage
	RoomID types.RoomID `json:"room_id"`
}

func NewSuccessCreateRoomMessage(id types.RoomID) (*SuccessCreateRoom, error) {
	msg := &SuccessCreateRoom{
		RoomID: id,
	}

	bmsg, err := interceptor.NewBaseMessage(message.NoneProtocol, nil, msg)
	if err != nil {
		return nil, err
	}

	msg.BaseMessage = bmsg
	return msg, nil
}

func NewSuccessCreateRoomMessageFactory(id types.RoomID) func() (message.Message, error) {
	return func() (message.Message, error) {
		return NewSuccessCreateRoomMessage(id)
	}
}

func (m *SuccessCreateRoom) GetProtocol() message.Protocol {
	return SuccessCreateRoomProtocol
}

func (m *SuccessCreateRoom) ReadProcess(_ context.Context, _i interceptor.Interceptor, _ interceptor.Connection) error {
	_, ok := _i.(*chat.ClientInterceptor)
	if !ok {
		return interceptor.ErrInvalidInterceptor
	}

	// NOTE: INTENTIONALLY EMPTY
	return nil
}
