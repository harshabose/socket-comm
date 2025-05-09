package messages

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/process"
)

var CreateRoomProtocol message.Protocol = "room:create_room"

// CreateRoom is the message sent by the client to the server when the client wants to create a room.
// When received by the server, the server will create a room with the given room id and allowed clients.
// The room is created with the given TTL.
// A process to KILL the room after TTL is created in the background.
// The server will then send a SuccessCreateRoom message to the client.
// NOTE: THIS DOES NOT ADD THE REQUESTING CLIENT TO THE ROOM; THIS IS MANAGED BY ANOTHER MESSAGE
type CreateRoom struct {
	interceptor.BaseMessage
	process.CreateRoom
}

func (m *CreateRoom) GetProtocol() message.Protocol {
	return CreateRoomProtocol
}

func (m *CreateRoom) ReadProcess(ctx context.Context, _i interceptor.Interceptor, connection interceptor.Connection) error {
	i, ok := _i.(*chat.ServerInterceptor)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	s, err := i.GetState(connection)
	if err != nil {
		return err
	}

	if err := i.Rooms.Process(ctx, m, s); err != nil {
		return err
	}

	return process.NewSendMessage(NewSuccessCreateRoomMessageFactory(m.RoomID)).Process(ctx, nil, s)
}
