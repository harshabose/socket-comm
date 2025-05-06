package process

import (
	"context"
	"fmt"
	"time"

	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
)

type WaitUntilIdent struct {
	ctx      context.Context
	duration time.Duration
}

type WaitUntilIdentOption func(*WaitUntilIdent)

func WithTickerDuration(duration time.Duration) WaitUntilIdentOption {
	return WaitUntilIdentOption(func(ident *WaitUntilIdent) {
		ident.duration = duration
	})
}

func NewWaitUntilIdentComplete(ctx context.Context, options ...WaitUntilIdentOption) *WaitUntilIdent {
	i := &WaitUntilIdent{
		ctx:      ctx,
		duration: 500 * time.Millisecond,
	}

	for _, option := range options {
		option(i)
	}

	return i
}

func (p WaitUntilIdent) Process(_ interfaces.CanGetRoom, s interfaces.State) error {
	ticker := time.NewTicker(p.duration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if p.process(s) {
				return nil
			}
		case <-p.ctx.Done():
			return fmt.Errorf("error while processing WaitUntilIdent process; err: %s", errors.ErrContextCancelled.Error())
		}
	}
}

func (p WaitUntilIdent) process(s interfaces.State) bool {
	_, err := s.GetClientID()
	return err == nil
}
