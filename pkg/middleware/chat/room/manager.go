package room

import (
	"context"
	"fmt"
	"time"

	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type Manager struct {
	rooms map[types.RoomID]interfaces.Room
	ctx   context.Context
}

func (m *Manager) Add(id types.RoomID, s interfaces.State) error {
	room, err := m.GetRoom(id)
	if err != nil {
		return err
	}

	return room.Add(id, s)
}

func (m *Manager) Remove(id types.RoomID, s interfaces.State) error {
	room, err := m.GetRoom(id)
	if err != nil {
		return err
	}

	return room.Remove(id, s)
}

func (m *Manager) CreateRoom(id types.RoomID, allowed []types.ClientID, ttl time.Duration) (interfaces.Room, error) {
	if m.exists(id) {
		return nil, fmt.Errorf("error while creating room with id %s; err: %s", id, errors.ErrRoomAlreadyExists)
	}

	ctx, cancel := context.WithTimeout(m.ctx, ttl)
	room := NewRoom(ctx, cancel, id, allowed)

	m.rooms[id] = room

	go func() {
		if err := m.Process(NewDeleteRoomWaiter(ctx, m, id, ttl), nil); err != nil {
			fmt.Println(err.Error())
		}
	}() // TODO: DO I NEED THIS?

	return room, nil
}

func (m *Manager) GetRoom(id types.RoomID) (interfaces.Room, error) {
	exists := m.exists(id)
	if !exists {
		return nil, fmt.Errorf("error while getting room with id %s; err: %s", id, errors.ErrRoomNotFound)
	}

	return m.rooms[id], nil
}

func (m *Manager) exists(id types.RoomID) bool {
	_, exists := m.rooms[id]
	return exists
}

func (m *Manager) DeleteRoom(id types.RoomID) error {
	room, err := m.GetRoom(id)
	if err != nil {
		return fmt.Errorf("error while deleting room with id: %s; err: %s", id, err.Error())
	}

	if err := room.Close(); err != nil {
		return fmt.Errorf("error while deleting room with id: %s; err: %s", id, err.Error())
	}

	return nil
}

func (m *Manager) WriteMessage(roomid types.RoomID, msg message.Message, from types.ClientID, tos ...types.ClientID) error {
	room, err := m.GetRoom(roomid)
	if err != nil {
		return err
	}

	w, ok := room.(interfaces.CanWriteRoomMessage)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	return w.WriteMessage(roomid, msg, from, tos...)
}

func (m *Manager) Process(process interfaces.CanProcess, state interfaces.State) error {
	return process.Process(m, state)
}
