package messages

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
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
		return errors.ErrInterfaceMisMatch
	}

	s, err := i.GetState(connection)
	if err != nil {
		return err
	}

	return i.Rooms.Process(ctx, m, s)
	// TODO: AFTER SUCCESS, SEND ROOM CURRENT STATE TO THE CLIENT
	// TODO: THEN SEND SuccessJoinRoom MESSAGE TO THE CLIENT
}
