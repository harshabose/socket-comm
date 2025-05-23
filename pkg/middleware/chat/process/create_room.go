package process

import (
	"context"
	"time"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

// CreateRoom is a process that creates a room as per requested by the client and adds the client to it.
type CreateRoom struct {
	RoomID  types.RoomID           `json:"room_id"`
	Allowed []interceptor.ClientID `json:"allowed"`
	TTL     time.Duration          `json:"ttl"`
	AsyncProcess
}

func NewCreateRoom(roomID types.RoomID, allowed []interceptor.ClientID, ttl time.Duration) *CreateRoom {
	return &CreateRoom{
		RoomID:  roomID,
		Allowed: allowed,
		TTL:     ttl,
	}
}

// Process requires room processor to be passed in.
func (p *CreateRoom) Process(ctx context.Context, processor interceptor.CanProcess, _ interceptor.State) error {
	select {
	case <-ctx.Done():
		return interceptor.ErrContextCancelled
	default:
		r, ok := processor.(interfaces.CanCreateRoom)
		if !ok {
			return interceptor.ErrInterfaceMisMatch
		}

		_, err := r.CreateRoom(p.RoomID, p.Allowed, p.TTL)
		if err != nil {
			return err
		}

		return nil
	}
}
