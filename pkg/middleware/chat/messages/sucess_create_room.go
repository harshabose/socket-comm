package messages

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

var SuccessCreateRoomProtocol message.Protocol = "room:success_create_room"

// SuccessCreateRoom is the message sent by the server to the client, which requested to create the room
// when creation of the room is successfully.
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
		return errors.ErrInvalidInterceptor
	}

	// TODO: AS OF NOW, THIS IS A EMPTY PROCESS.
	return nil
}
