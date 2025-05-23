package process

import (
	"context"

	"github.com/harshabose/socket-comm/internal/util"
	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/health"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type UpdateHealthSnapshot struct {
	Snapshot health.Snapshot `json:"snapshot"`
	AsyncProcess
}

func NewUpdateHealthSnapshot(roomID types.RoomID, snapGetter interfaces.CanGetHealthSnapshot) (*UpdateHealthSnapshot, error) {
	snapshot, err := snapGetter.GetHealthSnapshot(roomID)
	if err != nil {
		return nil, err
	}
	return &UpdateHealthSnapshot{
		Snapshot: snapshot,
	}, nil
}

func (p *UpdateHealthSnapshot) Process(ctx context.Context, processor interceptor.CanProcess, _ interceptor.State) error {
	select {
	case <-ctx.Done():
		return interceptor.ErrContextCancelled
	default:
		c, ok := processor.(interfaces.CanGetHealth)
		if !ok {
			return interceptor.ErrInterfaceMisMatch
		}

		h, err := c.GetHealth(p.Snapshot.Roomid)
		if err != nil {
			return err
		}

		merr := util.NewMultiError()
		for client, stat := range p.Snapshot.Participants {
			merr.Add(h.Update(p.Snapshot.Roomid, client, stat))
		}

		return merr.ErrorOrNil()
	}
}
