package messages

import (
	"context"
	"fmt"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

const FailJoinRoomProtocol message.Protocol = "room:fail_join_room"

type FailJoinRoom struct {
	interceptor.BaseMessage
	RoomID types.RoomID `json:"room_id"`
	Error  string       `json:"error"`
}

func NewFailJoinRoomMessage(id types.RoomID, err error) (*FailJoinRoom, error) {
	msg := &FailJoinRoom{
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

func NewFailJoinRoomMessageFactory(id types.RoomID, err error) func() (message.Message, error) {
	return func() (message.Message, error) {
		return NewFailJoinRoomMessage(id, err)
	}
}

func (m *FailJoinRoom) GetProtocol() message.Protocol {
	return FailJoinRoomProtocol
}

func (m *FailJoinRoom) ReadProcess(_ context.Context, _i interceptor.Interceptor, _ interceptor.Connection) error {
	_, ok := _i.(*chat.ClientInterceptor)
	if !ok {
		return interceptor.ErrInvalidInterceptor
	}

	fmt.Println("failed to join room:", m.Error)

	// NOTE: INTENTIONALLY EMPTY
	return nil
}
