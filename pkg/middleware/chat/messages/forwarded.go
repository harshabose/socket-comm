package messages

import (
	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
)

var ForwardedMessageProtocol message.Protocol = "room:forwarded_message"

type ForwardedMessage struct {
	interceptor.BaseMessage
}

func (m *ForwardedMessage) GetProtocol() message.Protocol {
	return ForwardedMessageProtocol
}

func newForwardedMessage(forward *ToForward) (*ForwardedMessage, error) {
	msg := &ForwardedMessage{}
	bmsg, err := interceptor.NewBaseMessage(forward.GetNextProtocol(), forward.NextPayload, msg)
	if err != nil {
		return nil, err
	}

	msg.BaseMessage = bmsg

	return msg, nil
}
