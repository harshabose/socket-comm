package messages

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
)

const ForwardedMessageProtocol message.Protocol = "room:forwarded_message"

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

func (m *ForwardedMessage) ReadProcess(_ context.Context, _i interceptor.Interceptor, _ interceptor.Connection) error {
	_, ok := _i.(*chat.ClientInterceptor)
	if !ok {
		return interceptor.ErrInterfaceMisMatch
	}

	// NOTE: INTENTIONALLY EMPTY

	return nil
	// NOTE: RETURNING NIL ASSUMING THAT THE NEXT PAYLOAD WILL BE MARSHALLED AFTER PROCESSING
}
