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
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type SendMessageRoom struct {
	msgFactory func() (message.Message, error)
	roomid     types.RoomID
	interval   time.Duration
	err        error
	mux        sync.RWMutex
	done       chan struct{}
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewSendMessageRoom(ctx context.Context, cancel context.CancelFunc, msgFactory func() (message.Message, error), roomid types.RoomID, interval time.Duration) *SendMessageRoom {
	return &SendMessageRoom{
		ctx:        ctx,
		cancel:     cancel,
		msgFactory: msgFactory,
		roomid:     roomid,
		interval:   interval,
		done:       make(chan struct{}),
	}
}

func (p *SendMessageRoom) Process(r interfaces.Processor, s *state.State) error {
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

func (p *SendMessageRoom) ProcessBackground(_r interfaces.Processor, s *state.State) interfaces.CanBeProcessedBackground {
	go func() {
		if err := p.Process(_r, s); err != nil {
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

func (p *SendMessageRoom) process(_r interfaces.Processor, _ *state.State) error {
	r, ok := _r.(interfaces.CanGetRoom)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	room, err := r.GetRoom(p.roomid)
	if err != nil {
		return err
	}

	participants := room.GetParticipants()

	merr := util.NewMultiError()

	for _, participant := range participants {
		msg, err := p.msgFactory()
		if err != nil {
			merr.Add(err)
			continue
		}

		if err := room.WriteRoomMessage(p.roomid, msg, "", participant); err != nil {
			merr.Add(err)
		}
	}

	return merr.ErrorOrNil()
}
