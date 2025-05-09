package process

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type RemoveFromRoom struct {
	RoomID types.RoomID `json:"room_id"`
}

func (p *RemoveFromRoom) Process(ctx context.Context, processor interfaces.Processor, s *state.State) error {
	r, ok := processor.(interfaces.CanRemove)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	return r.Remove(p.RoomID, s)
}
