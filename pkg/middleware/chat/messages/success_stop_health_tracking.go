package messages

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

var SuccessUnmarkRoomForHealthTrackingProtocol message.Protocol = "chat:success_untrack_health"

type SuccessUnmarkRoomForHealthTracking struct {
	interceptor.BaseMessage
	RoomID types.RoomID `json:"room_id"`
}

func (m *SuccessUnmarkRoomForHealthTracking) GetProtocol() message.Protocol {
	return SuccessUnmarkRoomForHealthTrackingProtocol
}

func NewSuccessUntrackHealthInRoomMessage(id types.RoomID) (*SuccessUnmarkRoomForHealthTracking, error) {
	msg := &SuccessUnmarkRoomForHealthTracking{
		RoomID: id,
	}

	bmsg, err := interceptor.NewBaseMessage(message.NoneProtocol, nil, msg)
	if err != nil {
		return nil, err
	}

	msg.BaseMessage = bmsg
	return msg, nil
}

func NewSuccessUntrackHealthInRoomMessageFactory(id types.RoomID) func() (message.Message, error) {
	return func() (message.Message, error) {
		return NewSuccessUntrackHealthInRoomMessage(id)
	}
}

func (m *SuccessUnmarkRoomForHealthTracking) ReadProcess(ctx context.Context, _i interceptor.Interceptor, _ interceptor.Connection) error {
	_, ok := _i.(*chat.ClientInterceptor)
	if !ok {
		return errors.ErrInvalidInterceptor
	}

	// NOTE: INTENTIONALLY EMPTY
	return nil
}
