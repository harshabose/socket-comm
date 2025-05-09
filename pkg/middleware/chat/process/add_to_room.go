package process

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

// AddToRoom is a process that adds a state (client) to a room.
type AddToRoom struct {
	RoomID types.RoomID
	AsyncProcess
}

func NewAddToRoom(roomID types.RoomID) interfaces.CanBeProcessed {
	return &AddToRoom{
		RoomID: roomID,
	}
}

func (p *AddToRoom) Process(ctx context.Context, processor interfaces.Processor, s *state.State) error {
	select {
	case <-ctx.Done():
		return errors.ErrContextCancelled
	default:
		r, ok := processor.(interfaces.CanAdd)
		if !ok {
			return errors.ErrInterfaceMisMatch
		}

		return r.Add(p.RoomID, s)
	}
}
