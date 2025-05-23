package process

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
)

type SendMessage struct {
	factory func() (message.Message, error)
	AsyncProcess
}

func NewSendMessage(factory func() (message.Message, error)) *SendMessage {
	return &SendMessage{
		factory: factory,
	}
}

func (p *SendMessage) Process(ctx context.Context, _ interceptor.CanProcess, s interceptor.State) error {
	select {
	case <-ctx.Done():
		return interceptor.ErrContextCancelled
	default:
		msg, err := p.factory()
		if err != nil {
			return err
		}
		return s.Write(ctx, msg)
	}
}
