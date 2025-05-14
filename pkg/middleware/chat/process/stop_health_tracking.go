package process

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type UnMarkRoomForHealthTracking struct {
	RoomID types.RoomID
	AsyncProcess
}

func NewUnMarkRoomForHealthTracking(roomID types.RoomID) *UnMarkRoomForHealthTracking {
	return &UnMarkRoomForHealthTracking{
		RoomID: roomID,
	}
}

func (p *UnMarkRoomForHealthTracking) Process(ctx context.Context, processor interfaces.Processor, _ *state.State) error {
	select {
	case <-ctx.Done():
		return errors.ErrContextCancelled
	default:
		u, ok := processor.(interfaces.CanStopHealthTracking)
		if !ok {
			return errors.ErrInterfaceMisMatch
		}

		if err := u.StopHealthTracking(p.RoomID); err != nil {
			return err
		}

		return nil
	}
}
