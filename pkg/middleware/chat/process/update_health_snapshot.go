package process

import (
	"context"

	"github.com/harshabose/socket-comm/internal/util"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/health"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
)

type SnapshotGetter interface {
	GetSnapshot() (health.Snapshot, error)
}

type UpdateHealthSnapshot struct {
	health.Snapshot
	AsyncProcess
}

func NewUpdateHealthSnapshot(snapGetter SnapshotGetter) (*UpdateHealthSnapshot, error) {
	snapshot, err := snapGetter.GetSnapshot()
	if err != nil {
		return nil, err
	}
	return &UpdateHealthSnapshot{
		Snapshot: snapshot,
	}, nil
}

func (p *UpdateHealthSnapshot) Process(ctx context.Context, processor interfaces.Processor, s *state.State) error {
	select {
	case <-ctx.Done():
		return errors.ErrContextCancelled
	default:
		c, ok := processor.(interfaces.CanGetHealth)
		if !ok {
			return errors.ErrInterfaceMisMatch
		}

		h, err := c.GetHealth(p.Roomid)
		if err != nil {
			return err
		}

		// NOTE: I COULD JUST DO THE FOLLOWING:
		// h.Snapshot = p.Snapshot
		// BUT MAYBE SOME MISMATCH WITH TTL OR ROOMID OR SOMETHING CAN CORRUPT THIS DATA.

		merr := util.NewMultiError()
		for client, stat := range p.Participants {
			merr.Add(h.Update(p.Roomid, client, stat))
		}

		return merr.ErrorOrNil()
	}
}
