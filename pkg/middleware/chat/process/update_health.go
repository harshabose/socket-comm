package process

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/health"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

// UpdateHealthStat process updates the health processor with the given health.Stat
type UpdateHealthStat struct {
	RoomID types.RoomID `json:"room_id"`
	health.Stat
	AsyncProcess
}

func (p *UpdateHealthStat) Process(ctx context.Context, processor interceptor.CanProcess, s interceptor.State) error {
	select {
	case <-ctx.Done():
		return interceptor.ErrContextCancelled
	default:
		id, err := s.GetClientID()
		if err != nil {
			return err
		}

		u, ok := processor.(interfaces.CanUpdate)
		if !ok {
			return interceptor.ErrInterfaceMisMatch
		}

		if err := u.Update(p.RoomID, id, &p.Stat); err != nil {
			return err
		}
		return nil
	}
}
