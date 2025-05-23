package messages

import (
	"context"
	"time"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/process"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

const GetHealthSnapshotProtocol message.Protocol = "chat:get_health_snapshot"

type StartStreamingHealthSnapshots struct {
	interceptor.BaseMessage
	Roomid   types.RoomID  `json:"roomid"`
	Interval time.Duration `json:"interval"`
}

func (m *StartStreamingHealthSnapshots) GetProtocol() message.Protocol {
	return GetHealthSnapshotProtocol
}

func (m *StartStreamingHealthSnapshots) ReadProcess(ctx context.Context, _i interceptor.Interceptor, connection interceptor.Connection) error {
	i, ok := _i.(*chat.ServerInterceptor)
	if !ok {
		return interceptor.ErrInvalidInterceptor
	}

	s, err := i.GetState(connection)
	if err != nil {
		return err
	}

	r, ok := i.Rooms.(interfaces.CanGetRoom)
	if !ok {
		_ = process.NewSendMessage(NewFailStartHealthStreamingMessageFactory(m.Roomid, interceptor.ErrInterfaceMisMatch)).Process(ctx, nil, s)
		return interceptor.ErrInterfaceMisMatch
	}

	room, err := r.GetRoom(m.Roomid)
	if err != nil {
		_ = process.NewSendMessage(NewFailStartHealthStreamingMessageFactory(m.Roomid, err)).Process(ctx, nil, s)
		return err
	}

	if !room.IsRoomMarkedForHealthTracking() {
		err := interceptor.NewError("to get snapshots room must first be marked for health tracking")
		_ = process.NewSendMessage(NewFailStartHealthStreamingMessageFactory(m.Roomid, err)).Process(ctx, nil, s)
		return err
	}

	h, ok := i.Health.(interfaces.CanAddHealthSnapshotStreamer)
	if !ok {
		_ = process.NewSendMessage(NewFailStartHealthStreamingMessageFactory(m.Roomid, interceptor.ErrInterfaceMisMatch)).Process(ctx, nil, s)
		return interceptor.ErrInterfaceMisMatch
	}

	g, ok := i.Health.(interfaces.CanGetHealthSnapshot)
	if !ok {
		_ = process.NewSendMessage(NewFailStartHealthStreamingMessageFactory(m.Roomid, interceptor.ErrInterfaceMisMatch)).Process(ctx, nil, s)
		return interceptor.ErrInterfaceMisMatch
	}

	if err := h.AddHealthSnapshotStreamer(m.Roomid, s, process.NewSendMessageStream(NewUpdateHealthSnapshotMessageFactory(m.Roomid, g), m.Interval)); err != nil {
		_ = process.NewSendMessage(NewFailStartHealthStreamingMessageFactory(m.Roomid, err)).Process(ctx, nil, s)
		return err
	}

	if err := process.NewSendMessage(NewSuccessStartHealthStreamingMessageFactory(m.Roomid, room.GetAllowed(), room.TTL())).Process(ctx, nil, s); err != nil {
		// do not send a fail message here as failing to send a success message also means failing to send a failure message
		return err
	}

	return nil
}
