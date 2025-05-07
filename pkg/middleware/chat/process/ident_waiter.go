package process

import (
	"context"
	"fmt"
	"time"

	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
)

type IdentWaiterOptions func(*IdentWaiter)

func NewIdentWaiter(ctx context.Context, options ...IdentWaiterOptions) *IdentWaiter {
	p := &IdentWaiter{
		ctx:      ctx,
		duration: 500 * time.Millisecond,
	}

	for _, option := range options {
		option(p)
	}

	return p
}

type IdentWaiter struct {
	ctx      context.Context
	duration time.Duration
}

func (p *IdentWaiter) Process(_ interfaces.Processor, s *state.State) error {
	ticker := time.NewTicker(p.duration)
	defer ticker.Stop()

	for {
		select {
		case <-p.ctx.Done():
			return errors.ErrContextCancelled
		case <-ticker.C:
			if err := p.process(nil, s); err == nil {
				return nil
			}
			fmt.Println("waiting for ident...")
		}
	}
}

func (p *IdentWaiter) process(_ interfaces.Processor, s *state.State) error {
	_, err := s.GetClientID()
	return err
}
