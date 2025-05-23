package messages

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/process"
)

var JoinRoomProtocol message.Protocol = "room:join_room"

type JoinRoom struct {
	interceptor.BaseMessage
	process.AddToRoom
}

func (m *JoinRoom) GetProtocol() message.Protocol {
	return JoinRoomProtocol
}

func (m *JoinRoom) ReadProcess(ctx context.Context, _i interceptor.Interceptor, connection interceptor.Connection) error {
	i, ok := _i.(*chat.ServerInterceptor)
	if !ok {
		return interceptor.ErrInterfaceMisMatch
	}

	s, err := i.GetState(connection)
	if err != nil {
		return err
	}

	if err := i.Rooms.Process(ctx, m, s); err != nil {
		_ = process.NewSendMessage(NewFailJoinRoomMessageFactory(m.RoomID, err)).Process(ctx, nil, s)
		return err
	}

	return process.NewSendMessageToAllParticipantsInRoom(m.RoomID, NewSuccessJoinRoomMessageFactory(m.RoomID, interceptor.ClientID(m.GetCurrentHeader().Sender))).Process(ctx, nil, s)
}
