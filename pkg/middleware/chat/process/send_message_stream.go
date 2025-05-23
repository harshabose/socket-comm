package process

import (
	"context"
	"time"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
)

type SendMessageStream struct {
	Interval time.Duration `json:"interval"`
	*SendMessage
	AsyncProcess
}

func NewSendMessageStream(factory func() (message.Message, error), interval time.Duration) *SendMessageStream {
	return &SendMessageStream{
		Interval:    interval,
		SendMessage: NewSendMessage(factory),
	}
}

func (p *SendMessageStream) Process(ctx context.Context, _ interceptor.CanProcess, s interceptor.State) error {
	ticker := time.NewTicker(p.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := p.SendMessage.Process(ctx, nil, s); err != nil {
				return err
			}
		case <-ctx.Done():
			return interceptor.ErrContextCancelled
		}
	}
}
