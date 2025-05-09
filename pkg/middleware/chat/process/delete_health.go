package process

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type DeleteHealth struct {
	RoomID types.RoomID `json:"room_id"`
}

func (p *DeleteHealth) Process(ctx context.Context, processor interfaces.Processor, _ *state.State) error {
	select {
	case <-ctx.Done():
		return errors.ErrContextCancelled
	default:
		d, ok := processor.(interfaces.CanDeleteHealth)
		if !ok {
			return errors.ErrInterfaceMisMatch
		}

		return d.DeleteHealth(p.RoomID)
	}
}
