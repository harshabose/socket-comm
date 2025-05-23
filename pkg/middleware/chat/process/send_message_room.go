package process

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type SendMessageRoom struct {
	RoomID         types.RoomID         `json:"room_id"`
	ClientID       interceptor.ClientID `json:"client_id"`
	messageFactory func() (message.Message, error)
	AsyncProcess
}

func NewSendMessageBetweenParticipantsInRoom(roomID types.RoomID, clientID interceptor.ClientID, messageFactory func() (message.Message, error)) *SendMessageRoom {
	return &SendMessageRoom{
		RoomID:         roomID,
		ClientID:       clientID,
		messageFactory: messageFactory,
	}
}

func (p *SendMessageRoom) Process(ctx context.Context, processor interceptor.CanProcess, _ interceptor.State) error {
	select {
	case <-ctx.Done():
		return interceptor.ErrContextCancelled
	default:
		w, ok := processor.(interfaces.CanWriteRoomMessage)
		if !ok {
			return interceptor.ErrInterfaceMisMatch
		}

		msg, err := p.messageFactory()
		if err != nil {
			return err
		}

		return w.WriteRoomMessage(p.RoomID, msg, p.ClientID)
	}
}
