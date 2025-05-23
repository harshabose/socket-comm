package messages

import (
	"context"
	"fmt"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

const FailDeleteRoomProtocol message.Protocol = "room:fail_delete_room"

type FailDeleteRoom struct {
	interceptor.BaseMessage
	RoomID types.RoomID `json:"room_id"`
	Error  string       `json:"error"`
}

func NewFailDeleteRoomMessage(id types.RoomID, err error) (*FailDeleteRoom, error) {
	msg := &FailDeleteRoom{
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

func NewFailDeleteRoomMessageFactory(id types.RoomID, err error) func() (message.Message, error) {
	return func() (message.Message, error) {
		return NewFailDeleteRoomMessage(id, err)
	}
}

func (m *FailDeleteRoom) GetProtocol() message.Protocol {
	return FailDeleteRoomProtocol
}

func (m *FailDeleteRoom) ReadProcess(_ context.Context, _i interceptor.Interceptor, _ interceptor.Connection) error {
	_, ok := _i.(*chat.ClientInterceptor)
	if !ok {
		return interceptor.ErrInvalidInterceptor
	}

	fmt.Println("failed to create room:", m.Error)

	// NOTE: INTENTIONALLY EMPTY
	return nil
}
