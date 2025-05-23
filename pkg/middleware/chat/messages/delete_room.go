package messages

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/process"
)

const DeleteRoomProtocol message.Protocol = "room:delete_room"

type DeleteRoom struct {
	interceptor.BaseMessage
	process.DeleteRoom
}

func (m *DeleteRoom) GetProtocol() message.Protocol {
	return DeleteRoomProtocol
}

func (m *DeleteRoom) ReadProcess(ctx context.Context, _i interceptor.Interceptor, connection interceptor.Connection) error {
	i, ok := _i.(*chat.ServerInterceptor)
	if !ok {
		return interceptor.ErrInterfaceMisMatch
	}

	s, err := i.GetState(connection)
	if err != nil {
		return err
	}

	if err := i.Rooms.Process(ctx, m, s); err != nil {
		_ = process.NewSendMessage(NewFailDeleteRoomMessageFactory(m.RoomID, err)).Process(ctx, nil, s)
		return err
	}

	return process.NewSendMessage(NewSuccessDeleteRoomMessageFactory(m.RoomID)).Process(ctx, nil, s)
}
