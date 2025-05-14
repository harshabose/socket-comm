package messages

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/process"
)

var MarkRoomForHealthTrackingProtocol message.Protocol = "chat:track_health"

// StartHealthTracking is sent by an interested client (who wants the stats of the whole room) to start tracking the health status.
// This does not send the stats data to the interested client; this just tells the server to track the data.
// If any client wants the data, they can send a StartHealthTracking to the server.
type StartHealthTracking struct {
	interceptor.BaseMessage
	process.StartHealthTracking
}

func (m *StartHealthTracking) GetProtocol() message.Protocol {
	return MarkRoomForHealthTrackingProtocol
}

func (m *StartHealthTracking) ReadProcess(ctx context.Context, _i interceptor.Interceptor, connection interceptor.Connection) error {
	i, ok := _i.(*chat.ServerInterceptor)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	s, err := i.GetState(connection)
	if err != nil {
		return err
	}

	r, ok := i.Rooms.(interfaces.CanGetRoom)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	room, err := r.GetRoom(m.RoomID)
	if err != nil {
		return err
	}

	if err := i.Rooms.Process(ctx, m, nil); err != nil {
		return err
	}

	if err := process.NewCreateHealthRoom(m.RoomID, room.GetAllowed(), room.TTL()).Process(ctx, i.Health, s); err != nil {
		return err
	}

	if err := process.NewSendMessage(NewSuccessTrackHealthInRoomMessageFactory(m.RoomID)).Process(ctx, nil, s); err != nil {
		return err
	}

	return nil
}
