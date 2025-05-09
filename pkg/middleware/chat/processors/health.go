package processors

import (
	"context"
	"fmt"
	"sync"

	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/health"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type Health struct {
	health map[types.RoomID]*health.Health
	mux    sync.RWMutex
	ctx    context.Context
}

// ==============================================================================
// ================================ CORE METHODS ================================
// ==============================================================================

func (p *Health) CreateHealth(roomid types.RoomID, allowed []types.ClientID) (*health.Health, error) {
	if p.exists(roomid) {
		return nil, fmt.Errorf("error while creating h with id %s; err: %s", roomid, errors.ErrRoomAlreadyExists)
	}

	h := health.NewHealth(p.ctx, roomid, allowed)

	p.health[roomid] = h
	return h, nil
}

func (p *Health) DeleteHealth(roomid types.RoomID) error {
	if !p.exists(roomid) {
		return fmt.Errorf("error while deleting h with id: %s; err: %s", roomid, errors.ErrRoomNotFound)
	}
	// TODO: DO I NEED TO CLOSE THE HEALTH?
	delete(p.health, roomid)
	return nil
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

func (p *Health) GetHealthSnapshot(roomid types.RoomID) (health.Snapshot, error) {
	p.mux.RLock()
	defer p.mux.RUnlock()

	if !p.exists(roomid) {
		return health.Snapshot{}, fmt.Errorf("error while getting snapshot for room with id: %s; err: %s", roomid, errors.ErrRoomNotFound)
	}

	h := p.health[roomid]

	// NOTE: FOLLOWING IS DEEP-COPYING SNAPSHOT
	snapshot := health.Snapshot{
		Roomid:       h.Roomid,
		Allowed:      make([]types.ClientID, len(h.Allowed)),
		Participants: make(map[types.ClientID]*health.Stat, len(h.Participants)),
	}

	copy(snapshot.Allowed, h.Allowed)

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

func (p *Health) Add(roomid types.RoomID, id types.ClientID) error {
	p.mux.Lock()
	defer p.mux.Unlock()

	if !p.exists(roomid) {
		return fmt.Errorf("error while adding participant with id: %s; err: %s", id, errors.ErrRoomNotFound)
	}

	return p.health[roomid].Add(roomid, id)
}

func (p *Health) Remove(roomid types.RoomID, id types.ClientID) error {
	p.mux.Lock()
	defer p.mux.Unlock()

	if !p.exists(roomid) {
		return fmt.Errorf("error while removing participant with id: %s; err: %s", id, errors.ErrRoomNotFound)
	}

	return p.health[roomid].Remove(roomid, id)
}

func (p *Health) Update(roomid types.RoomID, id types.ClientID, stat *health.Stat) error {
	p.mux.Lock()
	defer p.mux.Unlock()

	if !p.exists(roomid) {
		return fmt.Errorf("error while updating participant with id: %s; err: %s", id, errors.ErrRoomNotFound)
	}

	return p.health[roomid].Update(roomid, id, stat)
}
