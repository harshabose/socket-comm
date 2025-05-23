package process

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type DeleteHealthRoom struct {
	RoomID types.RoomID `json:"room_id"`
	AsyncProcess
}

func NewDeleteHealthRoom(roomID types.RoomID) DeleteHealthRoom {
	return DeleteHealthRoom{
		RoomID: roomID,
	}
}

func (p *DeleteHealthRoom) Process(ctx context.Context, processor interceptor.CanProcess, _ interceptor.State) error {
	select {
	case <-ctx.Done():
		return interceptor.ErrContextCancelled
	default:
		d, ok := processor.(interfaces.CanDeleteHealth)
		if !ok {
			return interceptor.ErrInterfaceMisMatch
		}

		return d.DeleteHealth(p.RoomID)
	}
}
