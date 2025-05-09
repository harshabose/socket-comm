package process

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

// CreateHealth is a process that creates health for a room
type CreateHealth struct {
	RoomID  types.RoomID     `json:"room_id"`
	Allowed []types.ClientID `json:"allowed"`
	AsyncProcess
}

func NewCreateHealth(roomID types.RoomID, allowed []types.ClientID) interfaces.CanBeProcessed {
	return &CreateHealth{
		RoomID:  roomID,
		Allowed: allowed,
	}
}

func (p *CreateHealth) Process(ctx context.Context, processor interfaces.Processor, _ *state.State) error {
	select {
	case <-ctx.Done():
		return errors.ErrContextCancelled
	default:
		c, ok := processor.(interfaces.CanCreateHealth)
		if !ok {
			return errors.ErrInterfaceMisMatch
		}

		_, err := c.CreateHealth(p.RoomID, p.Allowed)
		if err != nil {
			return err
		}

		return nil
	}
}
