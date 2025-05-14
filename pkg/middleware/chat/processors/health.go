package processors

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/health"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/process"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type healthSession struct {
	health                 *health.Health
	healthTrackingDeleter  interfaces.CanBeProcessedBackground
	healthSnapshotStreamer map[types.ClientID]interfaces.CanBeProcessedBackground // receiver[]streamerProcess
}

type Health struct {
	health map[types.RoomID]*healthSession
	mux    sync.RWMutex
	ctx    context.Context
}

func NewHealthProcessor(ctx context.Context) *Health {
	return &Health{
		health: make(map[types.RoomID]*healthSession),
		ctx:    ctx,
	}
}

// ==============================================================================
// ================================ CORE METHODS ================================
// ==============================================================================

func (p *Health) CreateHealth(roomid types.RoomID, allowed []types.ClientID, ttl time.Duration) (*health.Health, error) {
	if p.exists(roomid) {
		return nil, fmt.Errorf("error while creating h with id %s; err: %s", roomid, errors.ErrRoomAlreadyExists)
	}

	p.health[roomid] = &healthSession{
		health:                 health.NewHealth(p.ctx, roomid, allowed, ttl),
		healthTrackingDeleter:  process.NewDeleteHealthWaiter(p.ctx, roomid, ttl).ProcessBackground(nil, p, nil),
		healthSnapshotStreamer: make(map[types.ClientID]interfaces.CanBeProcessedBackground),
	}

	return p.health[roomid].health, nil
}

func (p *Health) DeleteHealth(roomid types.RoomID) error {
	session, err := p.getSession(roomid)
	if err != nil {
		return err
	}

	if err := session.health.Close(); err != nil {
		return err
	}

	if session.healthTrackingDeleter != nil {
		session.healthTrackingDeleter.Stop()
	}

	delete(p.health, roomid)
	return nil
}

func (p *Health) GetHealth(roomid types.RoomID) (*health.Health, error) {
	if !p.exists(roomid) {
		return nil, fmt.Errorf("error while getting h with id %s; err: %s", roomid, errors.ErrRoomNotFound)
	}

	return p.health[roomid].health, nil
}

func (p *Health) getSession(roomid types.RoomID) (*healthSession, error) {
	if !p.exists(roomid) {
		return nil, fmt.Errorf("error while getting h with id %s; err: %s", roomid, errors.ErrRoomNotFound)
	}

	return p.health[roomid], nil
}

func (p *Health) exists(roomid types.RoomID) bool {
	_, exists := p.health[roomid]
	return exists
}

// ==============================================================================
// ========================== INTERFACE IMPLEMENTATIONS =========================
// ==============================================================================

func (p *Health) Process(ctx context.Context, process interfaces.CanBeProcessed, state *state.State) error {
	return process.Process(ctx, p, state)
}

func (p *Health) ProcessBackground(ctx context.Context, process interfaces.CanBeProcessedBackground, state *state.State) interfaces.CanBeProcessedBackground {
	return process.ProcessBackground(ctx, p, state)
}

// GetHealthSnapshot will retrieve the latest snapshot consisting of health stats of the clients within.
func (p *Health) GetHealthSnapshot(roomid types.RoomID) (health.Snapshot, error) {
	p.mux.RLock()
	defer p.mux.RUnlock()

	h, err := p.GetHealth(roomid)
	if err != nil {
		return health.Snapshot{}, err
	}

	// NOTE: FOLLOWING IS DEEP-COPYING SNAPSHOT
	snapshot := health.Snapshot{
		Roomid:       h.Roomid,
		Participants: make(map[types.ClientID]*health.Stat, len(h.Participants)),
	}

	for id, stat := range h.Participants {
		if stat != nil {
			statCopy := *stat
			snapshot.Participants[id] = &statCopy
		} else {
			snapshot.Participants[id] = nil
		}
	}

	return snapshot, nil
}

func (p *Health) AddHealthSnapshotStreamer(roomid types.RoomID, interval time.Duration, s *state.State) error {
	p.mux.Lock()
	defer p.mux.Unlock()

	id, err := s.GetClientID()
	if err != nil {
		return err
	}

	session, err := p.getSession(roomid)
	if err != nil {
		return fmt.Errorf("error while adding streamer for room %s; err: %s", roomid, errors.ErrRoomNotFound)
	}

	if !session.health.IsParticipant(id) {
		return fmt.Errorf("client %s is not a participant of room %s", id, roomid)
	}

	_, exists := session.healthSnapshotStreamer[id]
	if exists {
		fmt.Println(fmt.Errorf("streaming already exists for client %s", id))
		fmt.Println("restarting...")

		if err := p.RemoveHealthSnapshotStreamer(roomid, s); err != nil {
			return err
		}

		// WARN: MAYBE RISK OF INFINITE RECURSION HERE (NO REASON BUT THE EXISTENCE OF RECURSION ITSELF HAS RISKS
		return p.AddHealthSnapshotStreamer(roomid, interval, s)
	}

	// TODO: SEND SNAP MESSAGE
	session.healthSnapshotStreamer[id] = process.NewSendMessageStream(nil, interval).ProcessBackground(p.ctx, nil, s)
	return nil
}

func (p *Health) RemoveHealthSnapshotStreamer(roomid types.RoomID, s *state.State) error {
	p.mux.Lock()
	defer p.mux.Unlock()

	id, err := s.GetClientID()
	if err != nil {
		return err
	}

	session, err := p.getSession(roomid)
	if err != nil {
		return fmt.Errorf("error while removing streamer for room %s; err: %s", roomid, errors.ErrRoomNotFound)
	}

	_, exists := session.healthSnapshotStreamer[id]
	if !exists {
		return fmt.Errorf("streamer for client %s does not exist", id)
	}

	session.healthSnapshotStreamer[id].Stop()
	session.healthSnapshotStreamer[id] = nil
	delete(session.healthSnapshotStreamer, id)
	return nil
}

// Add adds the given client to the health tracking in the given room.
// Only after calling this method, the stat responses from the clients are updated.
func (p *Health) Add(roomid types.RoomID, id types.ClientID) error {
	p.mux.Lock()
	defer p.mux.Unlock()

	h, err := p.GetHealth(roomid)
	if err != nil {
		return err
	}

	return h.Add(roomid, id)
}

// Remove removes the given client from the health tracking from the given room.
// After removing, the stat responses are not updated any more.
func (p *Health) Remove(roomid types.RoomID, id types.ClientID) error {
	p.mux.Lock()
	defer p.mux.Unlock()

	h, err := p.GetHealth(roomid)
	if err != nil {
		return err
	}

	return h.Remove(roomid, id)
}

// Update updates the stats of the given client in the given room.
// If the client is not already added to the list, the update will fail.
func (p *Health) Update(roomid types.RoomID, id types.ClientID, stat *health.Stat) error {
	p.mux.Lock()
	defer p.mux.Unlock()

	h, err := p.GetHealth(roomid)
	if err != nil {
		return err
	}

	return h.Update(roomid, id, stat)
}
