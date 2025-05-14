package messages

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/process"
)

const UpdateHealthSnapshotProtocol message.Protocol = "chat:update_health_snapshot"

type UpdateHealthSnapshot struct {
	interceptor.BaseMessage
	process.UpdateHealthSnapshot
}

func (m *UpdateHealthSnapshot) GetProtocol() message.Protocol {
	return UpdateHealthSnapshotProtocol
}

func (m *UpdateHealthSnapshot) ReadProcess(ctx context.Context, _i interceptor.Interceptor, _ interceptor.Connection) error {
	i, ok := _i.(*chat.ClientInterceptor)
	if !ok {
		return errors.ErrInvalidInterceptor
	}

	return i.Health.Process(ctx, m, nil)
}
