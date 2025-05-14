package process

import (
	"context"
	"time"

	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type CreateHealthRoom struct {
	RoomID  types.RoomID     `json:"room_id"`
	Allowed []types.ClientID `json:"allowed"`
	TTL     time.Duration    `json:"ttl"`
	AsyncProcess
}

func NewCreateHealthRoom(id types.RoomID, allowed []types.ClientID, ttl time.Duration) CreateHealthRoom {
	return CreateHealthRoom{
		RoomID:  id,
		Allowed: allowed,
		TTL:     ttl,
	}
}

func (p *CreateHealthRoom) Process(ctx context.Context, processor interfaces.Processor, _ *state.State) error {
	select {
	case <-ctx.Done():
		return errors.ErrContextCancelled
	default:
		r, ok := processor.(interfaces.CanCreateHealth)
		if !ok {
			return errors.ErrInterfaceMisMatch
		}

		_, err := r.CreateHealth(p.RoomID, p.Allowed, p.TTL)
		if err != nil {
			return err
		}

		return nil
	}
}
