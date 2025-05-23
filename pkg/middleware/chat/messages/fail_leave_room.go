package messages

import (
	"context"
	"fmt"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

const FailLeaveRoomProtocol message.Protocol = "room:fail_leave_room"

type FailLeaveRoom struct {
	interceptor.BaseMessage
	RoomID types.RoomID `json:"room_id"`
	Error  string       `json:"error"`
}

func NewFailLeaveRoomMessage(id types.RoomID, err error) (*FailLeaveRoom, error) {
	msg := &FailLeaveRoom{
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

func NewFailLeaveRoomMessageFactory(id types.RoomID, err error) func() (message.Message, error) {
	return func() (message.Message, error) {
		return NewFailLeaveRoomMessage(id, err)
	}
}

func (m *FailLeaveRoom) GetProtocol() message.Protocol {
	return FailLeaveRoomProtocol
}

func (m *FailLeaveRoom) ReadProcess(_ context.Context, _i interceptor.Interceptor, _ interceptor.Connection) error {
	_, ok := _i.(*chat.ClientInterceptor)
	if !ok {
		return interceptor.ErrInvalidInterceptor
	}

	fmt.Println("failed to leave room:", m.Error)

	// NOTE: INTENTIONALLY EMPTY
	return nil
}
