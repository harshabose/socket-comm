package messages

import (
	"context"
	"fmt"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

const FailCreateRoomProtocol message.Protocol = "room:fail_create_room"

type FailCreateRoom struct {
	interceptor.BaseMessage
	RoomID types.RoomID `json:"room_id"`
	Error  string       `json:"error"`
}

func NewFailCreateRoomMessage(id types.RoomID, err error) (*FailCreateRoom, error) {
	msg := &FailCreateRoom{
		RoomID: id,
		Error:  err.Error(),
	}

	bmsg, err := interceptor.NewBaseMessage(message.NoneProtocol, nil, msg)
	if err != nil {
		return nil, err
	}

	msg.BaseMessage = bmsg
	return msg, nil
}

func NewFailCreateRoomMessageFactory(id types.RoomID, err error) func() (message.Message, error) {
	return func() (message.Message, error) {
		return NewFailCreateRoomMessage(id, err)
	}
}

func (m *FailCreateRoom) GetProtocol() message.Protocol {
	return FailCreateRoomProtocol
}

func (m *FailCreateRoom) ReadProcess(_ context.Context, _i interceptor.Interceptor, _ interceptor.Connection) error {
	_, ok := _i.(*chat.ClientInterceptor)
	if !ok {
		return interceptor.ErrInvalidInterceptor
	}

	fmt.Println("failed to create room:", m.Error)

	// NOTE: INTENTIONALLY EMPTY
	return nil
}
