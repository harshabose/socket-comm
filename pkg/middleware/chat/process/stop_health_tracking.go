package process

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type StopHealthTracking struct {
	RoomID types.RoomID
	AsyncProcess
}

func NewUnMarkRoomForHealthTracking(roomID types.RoomID) *StopHealthTracking {
	return &StopHealthTracking{
		RoomID: roomID,
	}
}

func (p *StopHealthTracking) Process(ctx context.Context, processor interceptor.CanProcess, _ interceptor.State) error {
	select {
	case <-ctx.Done():
		return interceptor.ErrContextCancelled
	default:
		u, ok := processor.(interfaces.CanStopHealthTracking)
		if !ok {
			return interceptor.ErrInterfaceMisMatch
		}

		if err := u.StopHealthTracking(p.RoomID); err != nil {
			return err
		}

		return nil
	}
}
