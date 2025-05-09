package messages

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/process"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

var SuccessDeleteRoomProtocol message.Protocol = "room:success_delete_room"

type SuccessDeleteRoom struct {
	interceptor.BaseMessage
	process.DeleteHealth
}

func NewSuccessDeleteRoomMessage(id types.RoomID) (*SuccessDeleteRoom, error) {
	msg := &SuccessDeleteRoom{
		DeleteHealth: process.DeleteHealth{
			RoomID: id,
		},
	}

	bmsg, err := interceptor.NewBaseMessage(message.NoneProtocol, nil, msg)
	if err != nil {
		return nil, err
	}

	msg.BaseMessage = bmsg
	return msg, nil
}

func NewSuccessDeleteRoomMessageFactory(id types.RoomID) func() (message.Message, error) {
	return func() (message.Message, error) {
		return NewSuccessDeleteRoomMessage(id)
	}
}

func (m *SuccessDeleteRoom) GetProtocol() message.Protocol {
	return SuccessDeleteRoomProtocol
}

func (m *SuccessDeleteRoom) ReadProcess(ctx context.Context, _i interceptor.Interceptor, connection interceptor.Connection) error {
	i, ok := _i.(*chat.ClientInterceptor)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	if err := i.Health.Process(ctx, m, nil); err != nil {
		return err
	}

	return nil
}
