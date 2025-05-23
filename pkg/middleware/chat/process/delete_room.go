package process

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type DeleteRoom struct {
	RoomID types.RoomID `json:"room_id"`
	AsyncProcess
}

// Process needs room processor to be passed in.
func (p *DeleteRoom) Process(ctx context.Context, processor interceptor.CanProcess, _ interceptor.State) error {
	select {
	case <-ctx.Done():
		return interceptor.ErrContextCancelled
	default:
		d, ok := processor.(interfaces.CanDeleteRoom)
		if !ok {
			return interceptor.ErrInterfaceMisMatch
		}

		return d.DeleteRoom(p.RoomID)
	}
}
