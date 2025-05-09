package messages

import (
	"context"
	"time"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/process"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

var TrackHealthProtocol message.Protocol = "room:track_health"

type TrackHealth struct {
	interceptor.BaseMessage
	RoomID          types.RoomID  `json:"room_id"`
	RequestInterval time.Duration `json:"request_interval"`
}

func (m *TrackHealth) GetProtocol() message.Protocol {
	return TrackHealthProtocol
}

func (m *TrackHealth) ReadProcess(ctx context.Context, _i interceptor.Interceptor, _ interceptor.Connection) error {
	i, ok := _i.(*chat.ServerInterceptor)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	r, ok := i.Rooms.(interfaces.CanGetRoom)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}
	room, err := r.GetRoom(m.RoomID)
	if err != nil {
		return err
	}

	_ = i.Rooms.ProcessBackground(ctx, process.NewSendMessageStreamRoom(ctx, NewRequestHealthFactory(m.RoomID), m.RoomID, m.RequestInterval, room.TTL()), nil)
	return nil
}
