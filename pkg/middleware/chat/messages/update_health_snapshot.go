package messages

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/process"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

const UpdateHealthSnapshotProtocol message.Protocol = "chat:update_health_snapshot"

type UpdateHealthSnapshot struct {
	interceptor.BaseMessage
	*process.UpdateHealthSnapshot
}

func NewUpdateHealthSnapshotMessageFactory(roomID types.RoomID, getter interfaces.CanGetHealthSnapshot) func() (message.Message, error) {
	return func() (message.Message, error) {
		return NewUpdateHealthSnapshotMessage(roomID, getter)
	}
}

func NewUpdateHealthSnapshotMessage(roomID types.RoomID, getter interfaces.CanGetHealthSnapshot) (*UpdateHealthSnapshot, error) {
	p, err := process.NewUpdateHealthSnapshot(roomID, getter)
	if err != nil {
		return nil, err
	}
	msg := &UpdateHealthSnapshot{
		UpdateHealthSnapshot: p,
	}

	bmsg, err := interceptor.NewBaseMessage(message.NoneProtocol, nil, msg)
	if err != nil {
		return nil, err
	}

	msg.BaseMessage = bmsg
	return msg, nil
}

func (m *UpdateHealthSnapshot) GetProtocol() message.Protocol {
	return UpdateHealthSnapshotProtocol
}

func (m *UpdateHealthSnapshot) ReadProcess(ctx context.Context, _i interceptor.Interceptor, _ interceptor.Connection) error {
	i, ok := _i.(*chat.ClientInterceptor)
	if !ok {
		return interceptor.ErrInvalidInterceptor
	}

	return i.Health.Process(ctx, m, nil)
}
