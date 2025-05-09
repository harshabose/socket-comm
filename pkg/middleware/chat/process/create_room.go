package process

import (
	"context"
	"time"

	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

// CreateRoom is a process that creates a room as per requested by the client and adds the client to it.
type CreateRoom struct {
	RoomID  types.RoomID     `json:"room_id"`
	Allowed []types.ClientID `json:"allowed"`
	TTL     time.Duration    `json:"ttl"`
	AsyncProcess
}

func (p *CreateRoom) Process(ctx context.Context, processor interfaces.Processor, s *state.State) error {
	// NOTE: CTX HERE MIGHT BE SHORT-LIVED
	select {
	case <-ctx.Done():
		return errors.ErrContextCancelled
	default:
		r, ok := processor.(interfaces.CanCreateRoom)
		if !ok {
			return errors.ErrInterfaceMisMatch
		}

		_, err := r.CreateRoom(p.RoomID, p.Allowed, p.TTL)
		if err != nil {
			return err
		}

		// TODO: USING SHORT-LIVED CTX IS NOT APPROPRIATE HERE; FOR NOW USING NIL
		_ = NewDeleteRoomWaiter(ctx, p.RoomID, p.TTL).ProcessBackground(nil, processor, s)

		return nil
	}
}
