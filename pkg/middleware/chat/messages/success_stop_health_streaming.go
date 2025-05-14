package messages

import (
	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/process"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

var SuccessStopHealthStreamingProtocol message.Protocol = "chat:success_stop_health_streaming"

type SuccessStopHealthStreaming struct {
	interceptor.BaseMessage
	process.DeleteHealthRoom
}

func NewSuccessStopHealthStreamingMessage(id types.RoomID) (*SuccessStopHealthStreaming, error) {
	msg := &SuccessStopHealthStreaming{
		DeleteHealthRoom: process.NewDeleteHealthRoom(id),
	}
	bmsg, err := interceptor.NewBaseMessage(message.NoneProtocol, nil, msg)
	if err != nil {
		return nil, err
	}

	msg.BaseMessage = bmsg
	return msg, nil
}

func NewSuccessStopHealthStreamingMessageFactory(id types.RoomID) func() (message.Message, error) {
	return func() (message.Message, error) {
		return NewSuccessStopHealthStreamingMessage(id)
	}
}

func (m *SuccessStopHealthStreaming) GetProtocol() message.Protocol {
	return SuccessStopHealthStreamingProtocol
}
