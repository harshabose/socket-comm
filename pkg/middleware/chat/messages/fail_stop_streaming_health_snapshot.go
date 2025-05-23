package messages

import (
	"context"
	"fmt"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

const FailStopHealthStreamingProtocol message.Protocol = "chat:fail_stop_streaming_health_snapshot"

type FailStopHealthStreaming struct {
	interceptor.BaseMessage
	RoomID types.RoomID `json:"room_id"`
	Error  string       `json:"error"`
}

func NewFailStopHealthStreamingMessage(id types.RoomID, err error) (*FailStopHealthStreaming, error) {
	msg := &FailStopHealthStreaming{
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

func NewFailStopHealthStreamingMessageFactory(id types.RoomID, err error) func() (message.Message, error) {
	return func() (message.Message, error) {
		return NewFailStopHealthStreamingMessage(id, err)
	}
}

func (m *FailStopHealthStreaming) GetProtocol() message.Protocol {
	return FailStopHealthStreamingProtocol
}

func (m *FailStopHealthStreaming) ReadProcess(_ context.Context, _i interceptor.Interceptor, _ interceptor.Connection) error {
	_, ok := _i.(*chat.ClientInterceptor)
	if !ok {
		return interceptor.ErrInvalidInterceptor
	}

	fmt.Println("failed to stop health streaming:", m.Error)

	// NOTE: INTENTIONALLY EMPTY
	return nil
}
