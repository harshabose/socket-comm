package process

import (
	"context"
	"fmt"
	"time"

	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type SendMessageStreamRoom struct {
	msgFactory func() (message.Message, error)
	roomid     types.RoomID
	interval   time.Duration
	AsyncProcess
}

func NewSendMessageStreamRoom(ctx context.Context, msgFactory func() (message.Message, error), roomid types.RoomID, interval time.Duration, duration time.Duration) *SendMessageStreamRoom {
	return &SendMessageStreamRoom{
		AsyncProcess: ManualAsyncProcessInitialisation(context.WithTimeout(ctx, duration)),
		msgFactory:   msgFactory,
		roomid:       roomid,
		interval:     interval,
	}
}

func (p *SendMessageStreamRoom) Process(ctx context.Context, r interfaces.Processor, _ *state.State) error {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := NewSendMessageAllRoom(p.roomid, p.msgFactory).Process(ctx, r, nil); err != nil {
				fmt.Println("error while sending message to room; err: ", err.Error())
			}
		case <-ctx.Done():
			return nil
		}
	}
}
