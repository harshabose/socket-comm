package process

import (
	"context"
	"fmt"
	"time"

	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type DeleteHealthWaiter struct {
	TTL    time.Duration `json:"ttl"`
	RoomID types.RoomID  `json:"room_id"`
	AsyncProcess
}

func NewDeleteHealthWaiter(ctx context.Context, roomid types.RoomID, ttl time.Duration) *DeleteHealthWaiter {
	return &DeleteHealthWaiter{
		AsyncProcess: ManualAsyncProcessInitialisation(context.WithTimeout(ctx, ttl)),
		TTL:          ttl,
		RoomID:       roomid,
	}
}

func (p *DeleteHealthWaiter) Process(ctx context.Context, processor interfaces.Processor, _ *state.State) error {
	d, ok := processor.(interfaces.CanDeleteHealth)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	timer := time.NewTimer(p.TTL)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return errors.ErrContextCancelled
		case <-timer.C:
			if err := p.process(d); err != nil {
				return fmt.Errorf("error while processing DeleteHealthWaiter process; err: %s", err.Error())
			}
			return nil
		}
	}
}

func (p *DeleteHealthWaiter) process(d interfaces.CanDeleteHealth) error {
	if err := d.DeleteHealth(p.RoomID); err != nil {
		return err
	}

	return nil
}
