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

const MarkRoomForHealthTrackingProtocol message.Protocol = "chat:track_health"

// StartHealthTracking is sent by an interested client (who wants the stats of the whole room) to start tracking the health status.
// This does not send the stats data to the interested client; this just tells the server to track the data.
// If any client wants the data, they can send a _ to the server.
type StartHealthTracking struct {
	interceptor.BaseMessage
	RoomID   types.RoomID  `json:"room_id"`
	Interval time.Duration `json:"interval"`
}

func (m *StartHealthTracking) GetProtocol() message.Protocol {
	return MarkRoomForHealthTrackingProtocol
}

func (m *StartHealthTracking) ReadProcess(ctx context.Context, _i interceptor.Interceptor, connection interceptor.Connection) error {
	i, ok := _i.(*chat.ServerInterceptor)
	if !ok {
		// Can't send fail message here because we don't have a state
		return interceptor.ErrInterfaceMisMatch
	}

	s, err := i.GetState(connection)
	if err != nil {
		// Can't send fail message here because we don't have a state
		return err
	}

	r, ok := i.Rooms.(interfaces.CanGetRoom)
	if !ok {
		_ = process.NewSendMessage(NewFailStartHealthTrackingMessageFactory(m.RoomID, interceptor.ErrInterfaceMisMatch)).Process(ctx, nil, s)
		return interceptor.ErrInterfaceMisMatch
	}

	room, err := r.GetRoom(m.RoomID)
	if err != nil {
		_ = process.NewSendMessage(NewFailStartHealthTrackingMessageFactory(m.RoomID, err)).Process(ctx, nil, s)
		return err
	}

	t, ok := i.Rooms.(interfaces.CanStartHealthTracking)
	if !ok {
		_ = process.NewSendMessage(NewFailStartHealthTrackingMessageFactory(m.RoomID, interceptor.ErrInterfaceMisMatch)).Process(ctx, nil, s)
		return interceptor.ErrInterfaceMisMatch
	}

	if err := t.StartHealthTracking(m.RoomID, m.Interval, process.NewSendMessageStreamToAllParticipants(nil, NewRequestHealthFactory(m.RoomID), m.RoomID, m.Interval, room.TTL())); err != nil {
		_ = process.NewSendMessage(NewFailStartHealthTrackingMessageFactory(m.RoomID, err)).Process(ctx, nil, s)
		return err
	}

	healthRoom := process.NewCreateHealthRoom(m.RoomID, room.GetAllowed(), room.TTL())
	if err := healthRoom.Process(ctx, i.Health, s); err != nil {
		_ = process.NewSendMessage(NewFailStartHealthTrackingMessageFactory(m.RoomID, err)).Process(ctx, nil, s)
		return err
	}

	if err := process.NewSendMessage(NewSuccessTrackHealthInRoomMessageFactory(m.RoomID)).Process(ctx, nil, s); err != nil {
		return err
	}

	return nil
}
