package messages

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/process"
)

var SuccessJoinRoomProtocol message.Protocol = "room:success_join_room"

// SuccessJoinRoom is the message sent by the server to the clients (including the requested client and roommates)
// when the client joins a room successfully.
type SuccessJoinRoom struct {
	interceptor.BaseMessage
	process.CreateHealth
}

func (m *SuccessJoinRoom) GetProtocol() message.Protocol {
	return SuccessJoinRoomProtocol
}

func (m *SuccessJoinRoom) ReadProcess(ctx context.Context, _i interceptor.Interceptor, _ interceptor.Connection) error {
	i, ok := _i.(*chat.ClientInterceptor)
	if !ok {
		return errors.ErrInvalidInterceptor
	}

	return i.Health.Process(ctx, m, nil)
}
