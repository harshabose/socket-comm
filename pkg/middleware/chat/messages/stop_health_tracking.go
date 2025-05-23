package messages

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/process"
)

const UntrackHealthInRoomProtocol message.Protocol = "chat:untrack_health"

type StopHealthTracking struct {
	interceptor.BaseMessage
	process.StopHealthTracking
}

func (m *StopHealthTracking) GetProtocol() message.Protocol {
	return UntrackHealthInRoomProtocol
}

func (m *StopHealthTracking) ReadProcess(ctx context.Context, _i interceptor.Interceptor, connection interceptor.Connection) error {
	i, ok := _i.(*chat.ServerInterceptor)
	if !ok {
		// Can't send fail message here because we don't have a state
		return interceptor.ErrInvalidInterceptor
	}

	s, err := i.GetState(connection)
	if err != nil {
		// Can't send fail message here because we don't have a state
		return err
	}

	if err := i.Rooms.Process(ctx, m, nil); err != nil {
		_ = process.NewSendMessage(NewFailStopHealthTrackingMessageFactory(m.RoomID, err)).Process(ctx, nil, s)
		return err
	}

	if err := process.NewDeleteHealthRoom(m.RoomID).Process(ctx, i.Health, s); err != nil {
		_ = process.NewSendMessage(NewFailStopHealthTrackingMessageFactory(m.RoomID, err)).Process(ctx, nil, s)
		return err
	}

	if err := process.NewSendMessage(NewSuccessUntrackHealthInRoomMessageFactory(m.RoomID)).Process(ctx, nil, s); err != nil {
		// do not send a fail message here as failing to send a success message also means failing to send a failure message
		return err
	}

	return nil
}
