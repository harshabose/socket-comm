package process

import (
	"context"
	"fmt"
	"time"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

// SendMessageStreamRoomToAllParticipants is a process that sends a message to all participants of a room.
// This wraps NewSendMessageToAllParticipantsInRoom and sends the message at every given interval.
// NOTE: THIS IS A PURE PROCESS; AND IS NOT ADVISED TO BE TAGGED IN A MESSAGE
type SendMessageStreamRoomToAllParticipants struct {
	msgFactory func() (message.Message, error)
	roomid     types.RoomID
	interval   time.Duration
	AsyncProcess
}

func NewSendMessageStreamToAllParticipants(ctx context.Context, msgFactory func() (message.Message, error), roomid types.RoomID, interval time.Duration, duration time.Duration) *SendMessageStreamRoomToAllParticipants {
	return &SendMessageStreamRoomToAllParticipants{
		AsyncProcess: ManualAsyncProcessInitialisation(context.WithTimeout(ctx, duration)),
		msgFactory:   msgFactory,
		roomid:       roomid,
		interval:     interval,
	}
}

func (p *SendMessageStreamRoomToAllParticipants) SetInterval(interval time.Duration) {
	p.interval = interval
}

// Process needs Room processor to be passed in.
func (p *SendMessageStreamRoomToAllParticipants) Process(ctx context.Context, r interceptor.CanProcess, _ interceptor.State) error {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := NewSendMessageToAllParticipantsInRoom(p.roomid, p.msgFactory).Process(ctx, r, nil); err != nil {
				fmt.Println("error while sending message to room; err: ", err.Error())
			}
		case <-ctx.Done():
			return nil
		}
	}
}
