package messages

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/process"
)

const CreateRoomProtocol message.Protocol = "room:create_room"

// CreateRoom is the message sent by the client to the server when the client wants to create a room.
// When received by the server, the server will create a room with the given room id and allowed clients.
// A process to KILL the room after TTL is spawned in the background.
// The server will then send a SuccessCreateRoom message to the client.
// NOTE: THIS DOES NOT ADD THE REQUESTING CLIENT TO THE ROOM; THIS IS MANAGED BY JoinRoom MESSAGE
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
		// NOTE: CANNOT SEND FAIL MESSAGE AS STATE IS NOT DISCOVERED
		return interceptor.ErrInterfaceMisMatch
	}

	s, err := i.GetState(connection)
	if err != nil {
		return err
	}

	if err := i.Rooms.Process(ctx, m, s); err != nil {
		_ = process.NewSendMessage(NewFailCreateRoomMessageFactory(m.RoomID, err)).Process(ctx, nil, s)
		return err
	}

	return process.NewSendMessage(NewSuccessCreateRoomMessageFactory(m.RoomID)).Process(ctx, nil, s)
}
