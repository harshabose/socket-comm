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

func NewCreateRoom(roomID types.RoomID, allowed []types.ClientID, ttl time.Duration) *CreateRoom {
	return &CreateRoom{
		RoomID:  roomID,
		Allowed: allowed,
		TTL:     ttl,
	}
}

// Process requires room processor to be passed in.
func (p *CreateRoom) Process(ctx context.Context, processor interfaces.Processor, _ *state.State) error {
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

		return nil
	}
}
