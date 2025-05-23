package processors

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/process"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/room"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type roomSession struct {
	// NOTE: AS OF NOW, I AM MANAGING THE ROOM-DELETION AND HEALTH-TRACKING LIKE THIS;
	// NOTE: I AM NOT SURE IF THIS IS THE RIGHT WAY OR NOT
	room                                        *room.Room
	deletionWaiter, healthTrackingRequestSender interceptor.CanBeProcessedBackground
}

type RoomManager struct {
	rooms map[types.RoomID]*roomSession
	mux   sync.RWMutex
	ctx   context.Context
}

func NewRoomProcessor(ctx context.Context) *RoomManager {
	return &RoomManager{
		rooms: make(map[types.RoomID]*roomSession),
		ctx:   ctx,
	}
}

func (m *RoomManager) Add(id types.RoomID, s interceptor.State) error {
	r, err := m.GetRoom(id)
	if err != nil {
		return err
	}

	return r.Add(id, s)
}

func (m *RoomManager) Remove(id types.RoomID, s interceptor.State) error {
	r, err := m.GetRoom(id)
	if err != nil {
		return err
	}

	return r.Remove(id, s)
}

// CreateRoom creates a new room with the specified id, allowed a client list, and time-to-live duration.
// The room is tagged for deletion when the time expires.
// It returns an error if a room with the given id already exists.
// Parameters:
//   - id: unique identifier for the room
//   - allowed: a list of client IDs that are allowed to join the room
//   - ttl: time-to-live duration after which the room will be automatically deleted
//
// Returns:
//   - *room.Room: pointer to the newly created room
//   - error: nil if successful, ErrRoomAlreadyExists if room already exists
func (m *RoomManager) CreateRoom(id types.RoomID, allowed []interceptor.ClientID, ttl time.Duration) (*room.Room, error) {
	if m.exists(id) {
		return nil, fmt.Errorf("error while creating r with id %s; err: %s", id, errors.ErrRoomAlreadyExists)
	}

	m.rooms[id] = &roomSession{
		room:                        room.NewRoom(m.ctx, id, allowed, ttl),
		deletionWaiter:              process.NewDeleteRoomWaiter(m.ctx, id, ttl).ProcessBackground(nil, m, nil),
		healthTrackingRequestSender: nil,
	}

	return m.rooms[id].room, nil
}

func (m *RoomManager) GetRoom(id types.RoomID) (*room.Room, error) {
	exists := m.exists(id)
	if !exists {
		return nil, fmt.Errorf("error while getting room with id %s; err: %s", id, errors.ErrRoomNotFound)
	}

	return m.rooms[id].room, nil
}

func (m *RoomManager) exists(id types.RoomID) bool {
	s, exists := m.rooms[id]
	return exists && s != nil && s.room != nil && s.room.ID() == id
}

// StartHealthTracking marks a room for health tracking and initiates a background process.
// The background process pings the room participants at the given interval rate until TTL.
// Parameters:
//   - id: unique identifier of the room to be health tracked
//   - interval: interval in seconds between ping; has to be more than one second and at least 10% of the TTL
//
// Returns:
//   - error: nil if successful, ErrRoomNotFound if room does not exist, or other errors if marking fails
func (m *RoomManager) StartHealthTracking(id types.RoomID, interval time.Duration, _p interceptor.CanBeProcessedBackground) error {
	p, ok := _p.(*process.SendMessageStreamRoomToAllParticipants)
	if !ok {
		return interceptor.ErrInterfaceMisMatch
	}

	if interval <= 0 {
		return fmt.Errorf("health tracking interval must be positive, got: %v", interval)
	}

	if interval < 1*time.Second {
		return fmt.Errorf("health tracking interval too small (minimum 1 second), got: %v", interval)
	}

	r, err := m.GetRoom(id)
	if err != nil {
		return err
	}

	if interval > r.TTL()/10 {
		return fmt.Errorf("health tracking interval too large (maximum %v = 10%% of TTL), got: %v", r.TTL()/10, interval)
	}

	p.SetInterval(interval)

	if err := r.StartHealthTracking(id); err != nil {
		return err
	}

	// m.rooms[id].healthTrackingRequestSender = process.NewSendMessageStreamToAllParticipants(
	// 	m.ctx, messages.NewRequestHealthFactory(id), id, interval, r.TTL()).ProcessBackground(nil, m, nil)

	m.rooms[id].healthTrackingRequestSender = p.ProcessBackground(nil, m, nil)

	return nil
}

func (m *RoomManager) IsHealthTracked(id types.RoomID) (bool, error) {
	r, err := m.GetRoom(id)
	if err != nil {
		return false, err
	}

	return r.IsRoomMarkedForHealthTracking(), nil
}

// StopHealthTracking stops health tracking for a room and cancels the background health request pings process.
// Parameters:
//   - id: unique identifier of the room to stop health tracking
//
// Returns:
//   - error: nil if successful, ErrRoomNotFound if room does not exist, or other errors if unmarking fails
func (m *RoomManager) StopHealthTracking(id types.RoomID) error {
	r, err := m.GetRoom(id)
	if err != nil {
		return err
	}

	if err := r.UnMarkRoomForHealthTracking(); err != nil {
		return err
	}

	m.rooms[id].healthTrackingRequestSender.Stop()
	return nil
}

// DeleteRoom deletes an existing room with the specified id.
// It stops all background processes associated with the room and removes it from the manager.
// Parameters:
//   - id: unique identifier of the room to be deleted
//
// Returns:
//   - error: nil if successful, ErrRoomNotFound if room does not exist, or other errors if closing fails
func (m *RoomManager) DeleteRoom(id types.RoomID) error {
	if m.exists(id) {
		return fmt.Errorf("error while deleting r with id: %s; err: %s", id, errors.ErrRoomNotFound)
	}

	session := m.rooms[id]

	if err := session.room.Close(); err != nil {
		return fmt.Errorf("error while deleting r with id: %s; err: %s", id, err.Error())
	}

	if session.healthTrackingRequestSender != nil {
		session.healthTrackingRequestSender.Stop()
	}

	if session.deletionWaiter != nil {
		session.deletionWaiter.Stop()
	}

	delete(m.rooms, id)
	return nil
}

func (m *RoomManager) WriteRoomMessage(roomid types.RoomID, msg message.Message, from interceptor.ClientID, tos ...interceptor.ClientID) error {
	r, err := m.GetRoom(roomid)
	if err != nil {
		return err
	}

	return r.WriteRoomMessage(roomid, msg, from, tos...)
}

func (m *RoomManager) Process(ctx context.Context, process interceptor.CanBeProcessed, state interceptor.State) error {
	return process.Process(ctx, m, state)
}

func (m *RoomManager) ProcessBackground(ctx context.Context, process interceptor.CanBeProcessedBackground, state interceptor.State) interceptor.CanBeProcessedBackground {
	return process.ProcessBackground(ctx, m, state)
}
