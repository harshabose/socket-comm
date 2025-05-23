package messages

import (
	"context"
	"time"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/process"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

const SuccessStartHealthStreamingProtocol message.Protocol = "chat:success_start_health_streaming"

type SuccessStartHealthStreaming struct {
	interceptor.BaseMessage
	process.CreateHealthRoom
}

func NewSuccessStartHealthStreamingMessage(id types.RoomID, allowed []interceptor.ClientID, ttl time.Duration) (*SuccessStartHealthStreaming, error) {
	msg := &SuccessStartHealthStreaming{
		CreateHealthRoom: process.NewCreateHealthRoom(id, allowed, ttl),
	}
	bmsg, err := interceptor.NewBaseMessage(message.NoneProtocol, nil, msg)
	if err != nil {
		return nil, err
	}

	msg.BaseMessage = bmsg
	return msg, nil
}

func NewSuccessStartHealthStreamingMessageFactory(id types.RoomID, allowed []interceptor.ClientID, ttl time.Duration) func() (message.Message, error) {
	return func() (message.Message, error) {
		return NewSuccessStartHealthStreamingMessage(id, allowed, ttl)
	}
}

func (m *SuccessStartHealthStreaming) GetProtocol() message.Protocol {
	return SuccessStartHealthStreamingProtocol
}

func (m *SuccessStartHealthStreaming) ReadProcess(ctx context.Context, _i interceptor.Interceptor, _ interceptor.Connection) error {
	i, ok := _i.(*chat.ClientInterceptor)
	if !ok {
		return interceptor.ErrInvalidInterceptor
	}

	if err := i.Health.Process(ctx, m, nil); err != nil {
		return err
	}

	return nil
	// NOTE: NO SUCCESS TRAIL MESSAGE
}
