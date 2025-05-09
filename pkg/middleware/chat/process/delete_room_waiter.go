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

// DeleteRoomWaiter is a process that waits until TTL to delete a room.
// NOTE: THIS IS A PURE PROCESS; AND IS NOT ADVISED TO BE TAGGED IN A MESSAGE
type DeleteRoomWaiter struct {
	TTL    time.Duration `json:"ttl"`
	RoomID types.RoomID  `json:"room_id"`
	AsyncProcess
}

func NewDeleteRoomWaiter(ctx context.Context, roomid types.RoomID, ttl time.Duration) *DeleteRoomWaiter {
	return &DeleteRoomWaiter{
		AsyncProcess: ManualAsyncProcessInitialisation(context.WithTimeout(ctx, ttl)),
		RoomID:       roomid,
		TTL:          ttl,
	}
}

func (p *DeleteRoomWaiter) Process(ctx context.Context, processor interfaces.Processor, _ *state.State) error {
	d, ok := processor.(interfaces.CanDeleteRoom)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	timer := time.NewTimer(p.TTL)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			if err := p.process(d); err != nil {
				return fmt.Errorf("error while processing DeleteRoomWaiter process; err: %s", err.Error())
			}
			return nil
		case <-ctx.Done():
			return fmt.Errorf("context cancelled before process completion")
		}
	}
}

func (p *DeleteRoomWaiter) process(d interfaces.CanDeleteRoom) error {
	if err := d.DeleteRoom(p.RoomID); err != nil {
		return err
	}

	return nil
}
