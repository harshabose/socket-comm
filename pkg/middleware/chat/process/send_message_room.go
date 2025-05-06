package process

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/harshabose/socket-comm/internal/util"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type SendMessageRoom struct {
	msgFactory func() message.Message
	roomid     types.RoomID
	interval   time.Duration
	err        error
	mux        sync.RWMutex
	done       chan struct{}
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewSendMessageRoom(ctx context.Context, cancel context.CancelFunc, msgFactory func() message.Message, roomid types.RoomID, interval time.Duration) *SendMessageRoom {
	return &SendMessageRoom{
		ctx:        ctx,
		cancel:     cancel,
		msgFactory: msgFactory,
		roomid:     roomid,
		interval:   interval,
		done:       make(chan struct{}),
	}
}

func (p *SendMessageRoom) Process(r interfaces.CanGetRoom, s interfaces.State) error {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := p.process(r, s); err != nil {
				fmt.Println("error while processing SendMessageRoom process; err: ", err.Error())
			}
		case <-p.ctx.Done():
			return nil
		}
	}
}

func (p *SendMessageRoom) ProcessBackground(r interfaces.CanGetRoom, s interfaces.State) interfaces.CanProcessBackground {
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

func (p *SendMessageRoom) Wait() error {
	<-p.done
	p.mux.RLock()
	defer p.mux.RUnlock()

	return p.err
}

func (p *SendMessageRoom) Stop() {
	p.cancel()
}

func (p *SendMessageRoom) process(r interfaces.CanGetRoom, _ interfaces.State) error {
	room, err := r.GetRoom(p.roomid)
	if err != nil {
		return err
	}

	participants := room.GetParticipants()

	w, ok := room.(interfaces.CanWriteRoomMessage)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	merr := util.NewMultiError()

	for _, participant := range participants {
		if err := w.WriteRoomMessage(p.roomid, p.msgFactory(), "", participant); err != nil {
			merr.Add(err)
		}
	}

	return merr.ErrorOrNil()
}
