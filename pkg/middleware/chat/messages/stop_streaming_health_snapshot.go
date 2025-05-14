package messages

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/process"
)

var StopStreamingHealthSnapshotProtocol message.Protocol = "chat:stop_streaming_health_snapshot"

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
		return errors.ErrInvalidInterceptor
	}

	s, err := i.GetState(connection)
	if err != nil {
		return err
	}

	if err := i.Health.Process(ctx, m, s); err != nil {
		return err
	}

	if err := process.NewSendMessage(NewSuccessStopHealthStreamingMessageFactory(m.Roomid)).Process(ctx, nil, s); err != nil {
		return err
	}

	return nil

}
