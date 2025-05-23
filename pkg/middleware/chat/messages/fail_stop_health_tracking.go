package messages

import (
	"context"
	"fmt"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

const FailStopHealthTrackingProtocol message.Protocol = "chat:fail_untrack_health"

type FailStopHealthTracking struct {
	interceptor.BaseMessage
	RoomID types.RoomID `json:"room_id"`
	Error  string       `json:"error"`
}

func NewFailStopHealthTrackingMessage(id types.RoomID, err error) (*FailStopHealthTracking, error) {
	msg := &FailStopHealthTracking{
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

func NewFailStopHealthTrackingMessageFactory(id types.RoomID, err error) func() (message.Message, error) {
	return func() (message.Message, error) {
		return NewFailStopHealthTrackingMessage(id, err)
	}
}

func (m *FailStopHealthTracking) GetProtocol() message.Protocol {
	return FailStopHealthTrackingProtocol
}

func (m *FailStopHealthTracking) ReadProcess(_ context.Context, _i interceptor.Interceptor, _ interceptor.Connection) error {
	_, ok := _i.(*chat.ClientInterceptor)
	if !ok {
		return interceptor.ErrInvalidInterceptor
	}

	fmt.Println("failed to stop health tracking:", m.Error)

	// NOTE: INTENTIONALLY EMPTY
	return nil
}
