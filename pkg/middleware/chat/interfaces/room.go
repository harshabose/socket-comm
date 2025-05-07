package interfaces

import (
	"time"

	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/room"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type CanAdd interface {
	Add(types.RoomID, *state.State) error
}

type CanRemove interface {
	Remove(types.RoomID, *state.State) error
}

type CanGetRoom interface {
	GetRoom(id types.RoomID) (*room.Room, error)
}

type CanWriteRoomMessage interface {
	WriteRoomMessage(roomid types.RoomID, msg message.Message, from types.ClientID, tos ...types.ClientID) error
}

type CanCreateRoom interface {
	CreateRoom(types.RoomID, []types.ClientID, time.Duration) (*room.Room, error)
}

type CanDeleteRoom interface {
	DeleteRoom(types.RoomID) error
}
