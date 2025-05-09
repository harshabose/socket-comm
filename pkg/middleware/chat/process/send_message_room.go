package process

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type SendMessageRoom struct {
	RoomID         types.RoomID   `json:"room_id"`
	ClientID       types.ClientID `json:"client_id"`
	messageFactory func() (message.Message, error)
	AsyncProcess
}

func NewSendMessageRoom(roomID types.RoomID, clientID types.ClientID, messageFactory func() (message.Message, error)) *SendMessageRoom {
	return &SendMessageRoom{
		RoomID:         roomID,
		ClientID:       clientID,
		messageFactory: messageFactory,
	}
}

func (p *SendMessageRoom) Process(ctx context.Context, processor interfaces.Processor, _ *state.State) error {
	select {
	case <-ctx.Done():
		return errors.ErrContextCancelled
	default:
		w, ok := processor.(interfaces.CanWriteRoomMessage)
		if !ok {
			return errors.ErrInterfaceMisMatch
		}

		msg, err := p.messageFactory()
		if err != nil {
			return err
		}

		return w.WriteRoomMessage(p.RoomID, msg, p.ClientID)
	}
}
