package process

import (
	"context"
	"fmt"
	"time"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

// DeleteRoomWaiter is a process that waits until TTL to delete a room.
// NOTE: THIS IS A PURE PROCESS; AND IS NOT ADVISED TO BE TAGGED IN A MESSAGE
type DeleteRoomWaiter struct {
	TTL time.Duration `json:"ttl"`
	DeleteRoom
	AsyncProcess
}

func NewDeleteRoomWaiter(ctx context.Context, roomid types.RoomID, ttl time.Duration) *DeleteRoomWaiter {
	return &DeleteRoomWaiter{
		AsyncProcess: ManualAsyncProcessInitialisation(context.WithTimeout(ctx, ttl)),
		DeleteRoom:   DeleteRoom{RoomID: roomid},
		TTL:          ttl,
	}
}

func (p *DeleteRoomWaiter) Process(ctx context.Context, processor interceptor.CanProcess, _ interceptor.State) error {
	timer := time.NewTimer(p.TTL)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			if err := p.DeleteRoom.Process(ctx, processor, nil); err != nil {
				return fmt.Errorf("error while processing DeleteRoomWaiter process; err: %s", err.Error())
			}
			return nil
		case <-ctx.Done():
			return fmt.Errorf("context cancelled before process completion")
		}
	}
}
