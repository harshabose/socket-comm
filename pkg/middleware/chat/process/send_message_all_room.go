package process

import (
	"context"
	"fmt"

	"github.com/harshabose/socket-comm/internal/util"
	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

// SendMessageToAllParticipantsInRoom is a process that sends a message to all participants of a room.
// NOTE: THIS IS A PURE PROCESS; AND IS NOT ADVISED TO BE TAGGED IN A MESSAGE
type SendMessageToAllParticipantsInRoom struct {
	msgFactory func() (message.Message, error)
	Roomid     types.RoomID `json:"roomid"`
	AsyncProcess
}

func NewSendMessageToAllParticipantsInRoom(roomid types.RoomID, msgFactory func() (message.Message, error)) *SendMessageToAllParticipantsInRoom {
	return &SendMessageToAllParticipantsInRoom{
		msgFactory: msgFactory,
		Roomid:     roomid,
	}
}

// Process needs Room processor to be passed in.
func (p *SendMessageToAllParticipantsInRoom) Process(ctx context.Context, processor interceptor.CanProcess, _ interceptor.State) error {
	select {
	case <-ctx.Done():
		return interceptor.ErrContextCancelled
	default:
		r, ok := processor.(interfaces.CanGetRoom)
		if !ok {
			return interceptor.ErrInterfaceMisMatch
		}

		room, err := r.GetRoom(p.Roomid)
		if err != nil {
			return err
		}

		participants := room.GetParticipants()
		merr := util.NewMultiError()

		for _, participant := range participants {
			if err := NewSendMessageBetweenParticipantsInRoom(p.Roomid, participant, p.msgFactory).Process(ctx, processor, nil); err != nil {
				fmt.Println("error while sending message to room; err: ", err.Error())
			}
		}

		return merr.ErrorOrNil()
	}
}
