package messages

import (
	"context"
	"fmt"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

const ForwardMessageProtocol message.Protocol = "room:forward_message"

type ToForward struct {
	interceptor.BaseMessage
	RoomID types.RoomID           `json:"room_id"`
	To     []interceptor.ClientID `json:"to"`
}

func (m *ToForward) GetProtocol() message.Protocol {
	return ForwardMessageProtocol
}

func (m *ToForward) ReadProcess(ctx context.Context, _i interceptor.Interceptor, connection interceptor.Connection) error {
	i, ok := _i.(*chat.ServerInterceptor)
	if !ok {
		return interceptor.ErrInterfaceMisMatch
	}
	s, err := i.GetState(connection)
	if err != nil {
		return err
	}

	return i.Rooms.Process(ctx, m, s)
}

func (m *ToForward) Process(_ context.Context, p interceptor.CanProcess, s interceptor.State) error {
	msg, err := newForwardedMessage(m)
	if err != nil {
		return err
	}

	clientID, err := s.GetClientID()
	if err != nil {
		return err
	}

	if interceptor.ClientID(m.CurrentHeader.Sender) != clientID {
		return fmt.Errorf("error while read processing 'ToForward'; From and ClientID did not match")
	}

	w, ok := p.(interfaces.CanWriteRoomMessage)
	if !ok {
		return interceptor.ErrInterfaceMisMatch
	}

	if err := w.WriteRoomMessage(m.RoomID, msg, clientID, m.To...); err != nil {
		return err
	}

	return nil
}
