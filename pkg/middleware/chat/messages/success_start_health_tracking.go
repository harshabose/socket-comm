package messages

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

const SuccessTrackHealthInRoomProtocol message.Protocol = "chat:success_track_health"

type SuccessTrackHealthInRoom struct {
	interceptor.BaseMessage
	RoomID types.RoomID `json:"room_id"`
}

func (m *SuccessTrackHealthInRoom) GetProtocol() message.Protocol {
	return SuccessTrackHealthInRoomProtocol
}

func NewSuccessTrackHealthInRoomMessage(id types.RoomID) (*SuccessTrackHealthInRoom, error) {
	msg := &SuccessTrackHealthInRoom{
		RoomID: id,
	}
	bmsg, err := interceptor.NewBaseMessage(message.NoneProtocol, nil, msg)
	if err != nil {
		return nil, err
	}

	msg.BaseMessage = bmsg
	return msg, nil
}

func NewSuccessTrackHealthInRoomMessageFactory(id types.RoomID) func() (message.Message, error) {
	return func() (message.Message, error) {
		return NewSuccessTrackHealthInRoomMessage(id)
	}
}

func (m *SuccessTrackHealthInRoom) ReadProcess(_ context.Context, _i interceptor.Interceptor, _ interceptor.Connection) error {
	_, ok := _i.(*chat.ServerInterceptor)
	if !ok {
		return interceptor.ErrInterfaceMisMatch
	}

	// NOTE: INTENTIONALLY EMPTY
	return nil
}
