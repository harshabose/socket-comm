package process

import (
	"context"
	"fmt"

	"github.com/harshabose/socket-comm/internal/util"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type SendMessageAllRoom struct {
	msgFactory func() (message.Message, error)
	Roomid     types.RoomID `json:"roomid"`
	AsyncProcess
}

func NewSendMessageAllRoom(roomid types.RoomID, msgFactory func() (message.Message, error)) *SendMessageAllRoom {
	return &SendMessageAllRoom{
		msgFactory: msgFactory,
		Roomid:     roomid,
	}
}

func (p *SendMessageAllRoom) Process(ctx context.Context, processor interfaces.Processor, _ *state.State) error {
	select {
	case <-ctx.Done():
		return errors.ErrContextCancelled
	default:
		r, ok := processor.(interfaces.CanGetRoom)
		if !ok {
			return errors.ErrInterfaceMisMatch
		}

		room, err := r.GetRoom(p.Roomid)
		if err != nil {
			return err
		}

		participants := room.GetParticipants()
		merr := util.NewMultiError()

		for _, participant := range participants {
			if err := NewSendMessageRoom(p.Roomid, participant, p.msgFactory).Process(ctx, processor, nil); err != nil {
				fmt.Println("error while sending message to room; err: ", err.Error())
			}
		}

		return merr.ErrorOrNil()
	}
}
