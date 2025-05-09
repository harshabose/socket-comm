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

func NewIdentWaiter(options ...IdentWaiterOptions) *IdentWaiter {
	p := &IdentWaiter{
		duration: 500 * time.Millisecond,
	}

	for _, option := range options {
		option(p)
	}

	return p
}

type IdentWaiter struct {
	duration time.Duration
	AsyncProcess
}

func (p *IdentWaiter) Process(ctx context.Context, _ interfaces.Processor, s *state.State) error {
	ticker := time.NewTicker(p.duration)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return errors.ErrContextCancelled
		case <-ticker.C:
			if err := p.process(s); err == nil {
				return nil
			}
			fmt.Println("waiting for ident...")
		}
	}
}

func (p *IdentWaiter) process(s *state.State) error {
	_, err := s.GetClientID()
	return err
}
