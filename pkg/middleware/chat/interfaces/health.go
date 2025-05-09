package interfaces

import (
	"github.com/harshabose/socket-comm/pkg/middleware/chat/health"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type CanAddHealth interface {
	Add(roomid types.RoomID, id types.ClientID) error
}

type CanRemoveHealth interface {
	Remove(roomid types.RoomID, id types.ClientID) error
}

type CanUpdate interface {
	Update(roomid types.RoomID, id types.ClientID, s *health.Stat) error
}

type CanCreateHealth interface {
	CreateHealth(types.RoomID, []types.ClientID) (*health.Health, error)
}

type CanDeleteHealth interface {
	DeleteHealth(types.RoomID) error
}

type CanGetHealthSnapshot interface {
	GetHealthSnapshot(types.RoomID) (health.Snapshot, error)
}
