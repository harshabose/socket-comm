package room

import (
	"context"
	"fmt"
	"time"

	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type DeleteRoomWaiter struct {
	canDelete interfaces.CanDeleteRoom
	ttl       time.Duration
	roomid    types.RoomID
	ctx       context.Context
}

func NewDeleteRoomWaiter(ctx context.Context, manager interfaces.CanDeleteRoom, roomid types.RoomID, ttl time.Duration) *DeleteRoomWaiter {
	return &DeleteRoomWaiter{
		ctx:       ctx,
		roomid:    roomid,
		canDelete: manager,
		ttl:       ttl,
	}
}

func (p DeleteRoomWaiter) Process(r interfaces.CanGetRoom, _ interfaces.State) error {
	timer := time.NewTimer(p.ttl)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			if err := p.process(r); err != nil {
				return fmt.Errorf("error while processing DeleteRoomWaiter process; err: %s", err.Error())
			}
		case <-p.ctx.Done():
			return fmt.Errorf("context cancelled before process completion")
		}
	}
}

func (p DeleteRoomWaiter) process(r interfaces.CanGetRoom) error {
	room, err := r.GetRoom(p.roomid)
	if err != nil {
		return fmt.Errorf("error while processing DelteRoomWaiter process; err: %s", err.Error())
	}

	if err := room.Close(); err != nil {
		return err
	}

	if err := p.canDelete.DeleteRoom(p.roomid); err != nil {
		return err
	}

	return nil
}
