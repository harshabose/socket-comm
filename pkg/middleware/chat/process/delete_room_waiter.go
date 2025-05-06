package process

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type DeleteRoomWaiter struct {
	deleter interfaces.CanDeleteRoom
	ttl     time.Duration
	roomid  types.RoomID
	err     error
	ctx     context.Context
	cancel  context.CancelFunc
	mux     sync.RWMutex
	done    chan struct{}
}

func NewDeleteRoomWaiter(ctx context.Context, cancel context.CancelFunc, deleter interfaces.CanDeleteRoom, roomid types.RoomID, ttl time.Duration) *DeleteRoomWaiter {
	return &DeleteRoomWaiter{
		ctx:     ctx,
		cancel:  cancel,
		roomid:  roomid,
		deleter: deleter,
		err:     nil,
		ttl:     ttl,
		done:    make(chan struct{}),
	}
}

func (p *DeleteRoomWaiter) Process(r interfaces.CanGetRoom, _ interfaces.State) error {
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

func (p *DeleteRoomWaiter) ProcessBackground(r interfaces.CanGetRoom, s interfaces.State) interfaces.CanProcessBackground {
	go func() {
		if err := p.Process(r, s); err != nil {
			p.mux.Lock()
			p.err = err
			p.mux.Unlock()
			p.done <- struct{}{}

			fmt.Println(p.err)
		}
	}()

	return p
}

func (p *DeleteRoomWaiter) Wait() error {
	<-p.done
	p.mux.RLock()
	defer p.mux.RUnlock()

	return p.err
}

func (p *DeleteRoomWaiter) Stop() {
	p.cancel()
}

func (p *DeleteRoomWaiter) process(r interfaces.CanGetRoom) error {
	room, err := r.GetRoom(p.roomid)
	if err != nil {
		return fmt.Errorf("error while processing DelteRoomWaiter process; err: %s", err.Error())
	}

	if err := room.Close(); err != nil {
		return err
	}

	if err := p.deleter.DeleteRoom(p.roomid); err != nil {
		return err
	}

	return nil
}
