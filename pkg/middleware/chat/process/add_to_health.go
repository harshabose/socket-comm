package process

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type AddToHealth struct {
	RoomID      types.RoomID   `json:"room_id"`
	Participant types.ClientID `json:"participant"`
	AsyncProcess
}

func NewAddToHealth(roomID types.RoomID, participant types.ClientID) interfaces.CanBeProcessed {
	return &AddToHealth{
		RoomID:      roomID,
		Participant: participant,
	}
}

func (p *AddToHealth) Process(ctx context.Context, processor interfaces.Processor, _ *state.State) error {
	select {
	case <-ctx.Done():
		return errors.ErrContextCancelled
	default:
		a, ok := processor.(interfaces.CanAddHealth)
		if !ok {
			return errors.ErrInterfaceMisMatch
		}

		return a.Add(p.RoomID, p.Participant)
	}
}
