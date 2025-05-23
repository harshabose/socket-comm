package process

import (
	"context"
	"time"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type CreateHealthRoom struct {
	RoomID  types.RoomID           `json:"room_id"`
	Allowed []interceptor.ClientID `json:"allowed"`
	TTL     time.Duration          `json:"ttl"`
	AsyncProcess
}

func NewCreateHealthRoom(id types.RoomID, allowed []interceptor.ClientID, ttl time.Duration) CreateHealthRoom {
	return CreateHealthRoom{
		RoomID:  id,
		Allowed: allowed,
		TTL:     ttl,
	}
}

func (p *CreateHealthRoom) Process(ctx context.Context, processor interceptor.CanProcess, _ interceptor.State) error {
	select {
	case <-ctx.Done():
		return interceptor.ErrContextCancelled
	default:
		r, ok := processor.(interfaces.CanCreateHealth)
		if !ok {
			return interceptor.ErrInterfaceMisMatch
		}

		_, err := r.CreateHealth(p.RoomID, p.Allowed, p.TTL)
		if err != nil {
			return err
		}

		return nil
	}
}
