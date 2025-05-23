package messages

import (
	"context"
	"fmt"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

const FailStartHealthStreamingProtocol message.Protocol = "chat:fail_get_health_snapshot"

type FailStartHealthStreaming struct {
	interceptor.BaseMessage
	RoomID types.RoomID `json:"room_id"`
	Error  string       `json:"error"`
}

func NewFailStartHealthStreamingMessage(id types.RoomID, err error) (*FailStartHealthStreaming, error) {
	msg := &FailStartHealthStreaming{
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

func NewFailStartHealthStreamingMessageFactory(id types.RoomID, err error) func() (message.Message, error) {
	return func() (message.Message, error) {
		return NewFailStartHealthStreamingMessage(id, err)
	}
}

func (m *FailStartHealthStreaming) GetProtocol() message.Protocol {
	return FailStartHealthStreamingProtocol
}

func (m *FailStartHealthStreaming) ReadProcess(_ context.Context, _i interceptor.Interceptor, _ interceptor.Connection) error {
	_, ok := _i.(*chat.ClientInterceptor)
	if !ok {
		return interceptor.ErrInvalidInterceptor
	}

	fmt.Println("failed to start health streaming:", m.Error)

	// NOTE: INTENTIONALLY EMPTY
	return nil
}
