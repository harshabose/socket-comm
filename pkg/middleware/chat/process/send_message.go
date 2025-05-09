package process

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
)

type SendMessage struct {
	factory func() (message.Message, error)
	AsyncProcess
}

func NewSendMessage(factory func() (message.Message, error)) interfaces.CanBeProcessed {
	return &SendMessage{
		factory: factory,
	}
}

func (p *SendMessage) Process(ctx context.Context, _ interfaces.Processor, s *state.State) error {
	select {
	case <-ctx.Done():
		return errors.ErrContextCancelled
	default:
		msg, err := p.factory()
		if err != nil {
			return err
		}
		return s.Write(msg)
	}
}
