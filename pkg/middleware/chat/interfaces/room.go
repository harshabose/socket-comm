package interfaces

import (
	"io"
	"time"

	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type CanAdd interface {
	Add(types.RoomID, State) error
}

type CanRemove interface {
	Remove(types.RoomID, State) error
}

type Room interface {
	CanAdd
	CanRemove
	ID() types.RoomID
	GetParticipants() []types.ClientID
	io.Closer
}

type CanGetRoom interface {
	GetRoom(id types.RoomID) (Room, error)
}

type CanWriteRoomMessage interface {
	WriteRoomMessage(roomid types.RoomID, msg message.Message, from types.ClientID, tos ...types.ClientID) error
}

type CanCreateRoom interface {
	CreateRoom(types.RoomID, []types.ClientID, time.Duration) (Room, error)
}

type CanDeleteRoom interface {
	DeleteRoom(types.RoomID) error
}

type RoomManager interface {
	CanCreateRoom
	CanDeleteRoom
	CanGetRoom
}

type RoomProcessor interface {
	Process(CanProcess, State) error
}

type RoomProcessorBackground interface {
	ProcessBackground(CanProcessBackground, State) CanProcessBackground
}

type CanProcess interface {
	Process(CanGetRoom, State) error
}

type CanProcessBackground interface {
	ProcessBackground(room CanGetRoom, state State) CanProcessBackground
	Wait() error
	Stop()
}
