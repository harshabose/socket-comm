package messages

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/process"
)

const StopStreamingHealthSnapshotProtocol message.Protocol = "chat:stop_streaming_health_snapshot"

type StopStreamingHealthSnapshot struct {
	interceptor.BaseMessage
	process.StopStreamingHealthSnapshots
}

func (m *StopStreamingHealthSnapshot) GetProtocol() message.Protocol {
	return StopStreamingHealthSnapshotProtocol
}

func (m *StopStreamingHealthSnapshot) ReadProcess(ctx context.Context, _i interceptor.Interceptor, connection interceptor.Connection) error {
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

	if err := i.Health.Process(ctx, m, s); err != nil {
		_ = process.NewSendMessage(NewFailStopHealthStreamingMessageFactory(m.Roomid, err)).Process(ctx, nil, s)
		return err
	}

	if err := process.NewSendMessage(NewSuccessStopHealthStreamingMessageFactory(m.Roomid)).Process(ctx, nil, s); err != nil {
		// do not send a fail message here as failing to send a success message also means failing to send a failure message
		return err
	}

	return nil

}
