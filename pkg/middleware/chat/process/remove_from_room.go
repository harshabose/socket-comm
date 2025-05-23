package process

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type RemoveFromRoom struct {
	RoomID types.RoomID `json:"room_id"`
	AsyncProcess
}

func (p *RemoveFromRoom) Process(_ context.Context, processor interceptor.CanProcess, s interceptor.State) error {
	r, ok := processor.(interfaces.CanRemove)
	if !ok {
		return interceptor.ErrInterfaceMisMatch
	}

	return r.Remove(p.RoomID, s)
}
