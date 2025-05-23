package messages

import (
	"context"
	"fmt"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

const FailStartHealthTrackingProtocol message.Protocol = "chat:fail_track_health"

type FailStartHealthTracking struct {
	interceptor.BaseMessage
	RoomID types.RoomID `json:"room_id"`
	Error  string       `json:"error"`
}

func NewFailStartHealthTrackingMessage(id types.RoomID, err error) (*FailStartHealthTracking, error) {
	msg := &FailStartHealthTracking{
		RoomID: id,
		Error:  err.Error(),
	}

	bmsg, err := interceptor.NewBaseMessage(message.NoneProtocol, nil, msg)
	if err != nil {
		return nil, err
	}

	msg.BaseMessage = bmsg
	return msg, nil
}

func NewFailStartHealthTrackingMessageFactory(id types.RoomID, err error) func() (message.Message, error) {
	return func() (message.Message, error) {
		return NewFailStartHealthTrackingMessage(id, err)
	}
}

func (m *FailStartHealthTracking) GetProtocol() message.Protocol {
	return FailStartHealthTrackingProtocol
}

func (m *FailStartHealthTracking) ReadProcess(_ context.Context, _i interceptor.Interceptor, _ interceptor.Connection) error {
	_, ok := _i.(*chat.ClientInterceptor)
	if !ok {
		return interceptor.ErrInvalidInterceptor
	}

	fmt.Println("failed to start health tracking:", m.Error)

	// NOTE: INTENTIONALLY EMPTY
	return nil
}
