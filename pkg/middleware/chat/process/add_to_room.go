package process

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

// AddToRoom is a process that adds a state (client) to a room.
type AddToRoom struct {
	RoomID types.RoomID
	AsyncProcess
}

func NewAddToRoom(roomID types.RoomID) interceptor.CanBeProcessed {
	return &AddToRoom{
		RoomID: roomID,
	}
}

// Process needs room processor
func (p *AddToRoom) Process(ctx context.Context, processor interceptor.CanProcess, s interceptor.State) error {
	select {
	case <-ctx.Done():
		return interceptor.ErrContextCancelled
	default:
		r, ok := processor.(interfaces.CanAdd)
		if !ok {
			return interceptor.ErrInterfaceMisMatch
		}

		return r.Add(p.RoomID, s)
	}
}
