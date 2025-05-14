package messages

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/process"
)

var UntrackHealthInRoomProtocol message.Protocol = "chat:untrack_health"

type UnMarkRoomForHealthTracking struct {
	interceptor.BaseMessage
	process.UnMarkRoomForHealthTracking
}

func (m *UnMarkRoomForHealthTracking) GetProtocol() message.Protocol {
	return UntrackHealthInRoomProtocol
}

func (m *UnMarkRoomForHealthTracking) ReadProcess(ctx context.Context, _i interceptor.Interceptor, connection interceptor.Connection) error {
	i, ok := _i.(*chat.ServerInterceptor)
	if !ok {
		return errors.ErrInvalidInterceptor
	}

	s, err := i.GetState(connection)
	if err != nil {
		return err
	}

	if err := i.Rooms.Process(ctx, m, nil); err != nil {
		return err
	}

	if err := process.NewDeleteHealthRoom(m.RoomID).Process(ctx, i.Health, s); err != nil {
		return err
	}

	if err := process.NewSendMessage(nil).Process(ctx, nil, s); err != nil {
		return err
	}

	return nil
}
