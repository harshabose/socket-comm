package processors

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/room"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type RoomManager struct {
	rooms map[types.RoomID]*room.Room
	mux   sync.RWMutex
	ctx   context.Context
}

// TODO: DO I NEED MUX HERE?

func (m *RoomManager) Add(id types.RoomID, s *state.State) error {
	r, err := m.GetRoom(id)
	if err != nil {
		return err
	}

	return r.Add(id, s)
}

func (m *RoomManager) Remove(id types.RoomID, s *state.State) error {
	r, err := m.GetRoom(id)
	if err != nil {
		return err
	}

	return r.Remove(id, s)
}

func (m *RoomManager) CreateRoom(id types.RoomID, allowed []types.ClientID, ttl time.Duration) (*room.Room, error) {
	if m.exists(id) {
		return nil, fmt.Errorf("error while creating r with id %s; err: %s", id, errors.ErrRoomAlreadyExists)
	}

	r := room.NewRoom(m.ctx, id, allowed, ttl)
	m.rooms[id] = r

	return r, nil
}

func (m *RoomManager) GetRoom(id types.RoomID) (*room.Room, error) {
	exists := m.exists(id)
	if !exists {
		return nil, fmt.Errorf("error while getting room with id %s; err: %s", id, errors.ErrRoomNotFound)
	}

	return m.rooms[id], nil
}

func (m *RoomManager) exists(id types.RoomID) bool {
	_, exists := m.rooms[id]
	return exists
}

func (m *RoomManager) DeleteRoom(id types.RoomID) error {
	r, err := m.GetRoom(id)
	if err != nil {
		return fmt.Errorf("error while deleting r with id: %s; err: %s", id, err.Error())
	}

	if err := r.Close(); err != nil {
		return fmt.Errorf("error while deleting r with id: %s; err: %s", id, err.Error())
	}

	delete(m.rooms, id)
	return nil
}

func (m *RoomManager) WriteRoomMessage(roomid types.RoomID, msg message.Message, from types.ClientID, tos ...types.ClientID) error {
	r, err := m.GetRoom(roomid)
	if err != nil {
		return err
	}

	return r.WriteRoomMessage(roomid, msg, from, tos...)
}

func (m *RoomManager) Process(ctx context.Context, process interfaces.CanBeProcessed, state *state.State) error {
	return process.Process(ctx, m, state)
}

func (m *RoomManager) ProcessBackground(ctx context.Context, process interfaces.CanBeProcessedBackground, state *state.State) interfaces.CanBeProcessedBackground {
	return process.ProcessBackground(ctx, m, state)
}
