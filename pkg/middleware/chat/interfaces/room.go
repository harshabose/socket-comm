package interfaces

import (
	"time"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/room"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type CanAdd interface {
	Add(types.RoomID, interceptor.State) error
}

type CanRemove interface {
	Remove(types.RoomID, interceptor.State) error
}

type CanGetRoom interface {
	GetRoom(id types.RoomID) (*room.Room, error)
}

type CanWriteRoomMessage interface {
	WriteRoomMessage(roomid types.RoomID, msg message.Message, from interceptor.ClientID, tos ...interceptor.ClientID) error
}

type CanCreateRoom interface {
	CreateRoom(types.RoomID, []interceptor.ClientID, time.Duration) (*room.Room, error)
}

type CanDeleteRoom interface {
	DeleteRoom(types.RoomID) error
}

type CanStartHealthTracking interface {
	StartHealthTracking(types.RoomID, time.Duration, interceptor.CanBeProcessedBackground) error
	IsHealthTracked(types.RoomID) (bool, error)
}

type CanStopHealthTracking interface {
	StopHealthTracking(types.RoomID) error
}
