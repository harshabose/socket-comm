package keyexchange

import (
	"context"
	"fmt"
	"time"

	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/encryptionerr"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/types"
)

type SessionStateTargetWaiter struct {
	target types.SessionState
	ctx    context.Context
}

func NewSessionStateTargetWaiter(ctx context.Context, target types.SessionState) SessionStateTargetWaiter {
	return SessionStateTargetWaiter{
		target: target,
		ctx:    ctx,
	}
}

func (w SessionStateTargetWaiter) Process(protocol interfaces.Protocol, _ interfaces.State) error {
	p, ok := protocol.(interfaces.CanGetSessionState)
	if !ok {
		return encryptionerr.ErrInvalidMessageType
	}

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		// TODO: Implement BLOCKing (ideally using sync.Cond, or simple ticker)
		select {
		case <-ticker.C:
			if p.GetState() == w.target {
				return nil
			}
		case <-w.ctx.Done():
			return fmt.Errorf("timeout waiting for state %v: %w", w.target, w.ctx.Err())
		}
	}
}

// TODO: WRITE MORE KEY EXCHANGE PROCESSES HERE
