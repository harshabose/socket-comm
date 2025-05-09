package process

import (
	"context"
	"time"

	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/health"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type UpdateHealth struct {
	RoomID   types.RoomID  `json:"room_id"`
	Validity time.Duration `json:"validity"`
	health.Stat
	AsyncProcess
}

func (p *UpdateHealth) Process(ctx context.Context, processor interfaces.Processor, s *state.State) error {
	u, ok := processor.(interfaces.CanUpdate)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	select {
	case <-ctx.Done():
		return errors.ErrContextCancelled
	default:
		id, err := s.GetClientID()
		if err != nil {
			return err
		}

		if err := u.Update(p.RoomID, id, &p.Stat); err != nil {
			return err
		}
		return nil
	}
}
