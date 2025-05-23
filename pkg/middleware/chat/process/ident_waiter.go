package process

import (
	"context"
	"fmt"
	"time"

	"github.com/harshabose/socket-comm/pkg/interceptor"
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

func (p *IdentWaiter) Process(ctx context.Context, _ interceptor.CanProcess, s interceptor.State) error {
	ticker := time.NewTicker(p.duration)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return interceptor.ErrContextCancelled
		case <-ticker.C:
			if err := p.process(s); err == nil {
				return nil
			}
			fmt.Println("waiting for ident...")
		}
	}
}

func (p *IdentWaiter) process(s interceptor.State) error {
	_, err := s.GetClientID()
	return err
}
