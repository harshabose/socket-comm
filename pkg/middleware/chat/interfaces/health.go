package interfaces

import (
	"time"

	"github.com/harshabose/socket-comm/pkg/middleware/chat/health"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
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

type CanGetHealth interface {
	GetHealth(roomid types.RoomID) (*health.Health, error)
}

type CanAddHealthSnapshotStreamer interface {
	AddHealthSnapshotStreamer(types.RoomID, time.Duration, *state.State) error
}

type CanRemoveHealthSnapshotStreamer interface {
	RemoveHealthSnapshotStreamer(types.RoomID, *state.State) error
}

type CanCreateHealth interface {
	CreateHealth(types.RoomID, []types.ClientID, time.Duration) (*health.Health, error)
}

type CanDeleteHealth interface {
	DeleteHealth(types.RoomID) error
}

type CanGetHealthSnapshot interface {
	GetHealthSnapshot(types.RoomID) (health.Snapshot, error)
}
