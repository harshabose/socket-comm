package messages

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

const SuccessJoinRoomProtocol message.Protocol = "room:success_join_room"

// SuccessJoinRoom is the message sent by the server to the clients (including the requested client and roommates)
// when the client joins a room successfully.
// NOTE: THIS MESSAGE IS SENT TO ALL CLIENTS IN THE ROOM.
// This marks the end of the JoinRoom topic.
type SuccessJoinRoom struct {
	interceptor.BaseMessage
	RoomID   types.RoomID         `json:"room_id"`
	ClientID interceptor.ClientID `json:"client_id"`
}

func NewSuccessJoinRoomMessage(id types.RoomID, clientID interceptor.ClientID) (*SuccessJoinRoom, error) {
	msg := &SuccessJoinRoom{
		RoomID:   id,
		ClientID: clientID,
	}

	bmsg, err := interceptor.NewBaseMessage(message.NoneProtocol, nil, msg)
	if err != nil {
		return nil, err
	}

	msg.BaseMessage = bmsg
	return msg, nil
}

func NewSuccessJoinRoomMessageFactory(id types.RoomID, clientID interceptor.ClientID) func() (message.Message, error) {
	return func() (message.Message, error) {
		return NewSuccessJoinRoomMessage(id, clientID)
	}
}

func (m *SuccessJoinRoom) GetProtocol() message.Protocol {
	return SuccessJoinRoomProtocol
}

func (m *SuccessJoinRoom) ReadProcess(_ context.Context, _i interceptor.Interceptor, _ interceptor.Connection) error {
	_, ok := _i.(*chat.ClientInterceptor)
	if !ok {
		return interceptor.ErrInvalidInterceptor
	}

	// NOTE: INTENTIONALLY EMPTY
	return nil
}
