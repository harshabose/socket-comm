package health

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type Stat struct {
	ConnectionStatus types.ConnectionState  `json:"connection_status"`
	ConnectionUptime types.ConnectionUptime `json:"connection_uptime"`
	CPUUsage         types.CPUUsage         `json:"cpu_usage"`
	MemoryUsage      types.MemoryUsage      `json:"memory_usage"`
	NetworkUsage     types.NetworkUsage     `json:"network_usage"`
	Latency          types.LatencyMs        `json:"latency"`
}

type Health struct {
	Snapshot
	mux    sync.RWMutex
	cancel context.CancelFunc
	ctx    context.Context
}

type Snapshot struct {
	Roomid       types.RoomID                   `json:"roomid"` // NOTE: ADDED WHEN HEALTH IS CREATED.
	Allowed      []interceptor.ClientID         `json:"allowed"`
	TTL          time.Duration                  `json:"ttl"`
	Participants map[interceptor.ClientID]*Stat `json:"participants"` // NOTE: ADDED WHEN CLIENT JOINS. UPDATED WHEN CLIENT SENDS HEALTH RESPONSE.
}

// Marshal marshals the health struct into a json byte array.
// NOTE: THIS IS USED TO SEND HEALTH STATS TO THE CLIENT.
func (h *Snapshot) Marshal() ([]byte, error) {
	return json.Marshal(h)
}

func NewHealth(ctx context.Context, id types.RoomID, allowed []interceptor.ClientID, ttl time.Duration) *Health {
	ctx2, cancel := context.WithTimeout(ctx, ttl)
	h := &Health{
		Snapshot: Snapshot{
			Roomid:       id,
			Allowed:      allowed,
			TTL:          ttl,
			Participants: make(map[interceptor.ClientID]*Stat),
		},
		cancel: cancel,
		ctx:    ctx2,
	}

	for _, allowed := range h.Allowed {
		if err := h.Add(id, allowed); err != nil {
			fmt.Printf("error while adding allowed client: %s; but not returning an error", err)
			continue
		}
	}

	return h
}

func (h *Health) Ctx() context.Context {
	h.mux.RLock()
	defer h.mux.RUnlock()

	return h.ctx
}

func (h *Health) ID() types.RoomID {
	h.mux.RLock()
	defer h.mux.RUnlock()

	return h.Roomid
}

func (h *Health) GetTTL() time.Duration {
	h.mux.RLock()
	defer h.mux.RUnlock()

	return h.TTL
}

func (h *Health) GetAllowed() []interceptor.ClientID {
	h.mux.RLock()
	defer h.mux.RUnlock()

	return h.Allowed
}

func (h *Health) Add(roomid types.RoomID, id interceptor.ClientID) error {
	h.mux.Lock()
	defer h.mux.Unlock()

	if roomid != h.Roomid {
		return errors.ErrWrongRoom
	}

	select {
	case <-h.ctx.Done():
		return fmt.Errorf("error while adding client to health stats. client id: %s; room id: %s; err: %s", id, h.Roomid, interceptor.ErrContextCancelled.Error())
	default:
		if !h.IsAllowed(id) {
			return fmt.Errorf("client with id '%s' is not allowed in the health room with id '%s'; err: %s", id, h.Roomid, errors.ErrClientNotAllowed.Error())
		}

		if h.IsParticipant(id) {
			return fmt.Errorf("client with id '%s' already existing in the health stats with id %s; err: %s", id, h.Roomid, errors.ErrClientIsAlreadyParticipant)
		}

		h.Participants[id] = &Stat{}
		return nil
	}
}

func (h *Health) Remove(roomid types.RoomID, id interceptor.ClientID) error {
	h.mux.Lock()
	defer h.mux.Unlock()

	if roomid != h.Roomid {
		return errors.ErrWrongRoom
	}

	select {
	case <-h.ctx.Done():
		return fmt.Errorf("error while removing client to health stats. client id: %s; room id: %s; err: %s", id, h.Roomid, interceptor.ErrContextCancelled.Error())
	default:
		if !h.IsAllowed(id) {
			return fmt.Errorf("client with id '%s' is not allowed in the health room with id '%s'; err: %s", id, h.Roomid, errors.ErrClientNotAllowed.Error())
		}

		if !h.IsParticipant(id) {
			return fmt.Errorf("client with id '%s' does not exist in the health stats with id %s; err: %s", id, h.Roomid, errors.ErrClientNotAParticipant.Error())
		}

		delete(h.Participants, id)
		return nil
	}
}

func (h *Health) Update(roomid types.RoomID, id interceptor.ClientID, s *Stat) error {
	h.mux.Lock()
	defer h.mux.Unlock()

	if roomid != h.Roomid {
		return errors.ErrWrongRoom
	}

	select {
	case <-h.ctx.Done():
		return fmt.Errorf("error while adding client to health stats. client id: %s; room id: %s; err: %s", id, h.Roomid, interceptor.ErrContextCancelled.Error())
	default:
		if !h.IsAllowed(id) {
			return fmt.Errorf("client with id '%s' is not allowed in the health room with id '%s'; err: %s", id, h.Roomid, errors.ErrClientNotAllowed.Error())
		}

		if !h.IsParticipant(id) {
			return fmt.Errorf("client with id '%s' does not exist in the health stats with id %s; err: %s", id, h.Roomid, errors.ErrClientNotAParticipant.Error())
		}

		h.Participants[id] = s
		return nil
	}
}

// Close will close this health room and does not allow any further updates to the participant stats.
// This will also delete all the health stats data.
func (h *Health) Close() error {
	h.mux.Lock()
	defer h.mux.Unlock()

	select {
	case <-h.ctx.Done():
		return fmt.Errorf("error while closing the health room with id %s; err: %s", h.Roomid, interceptor.ErrContextCancelled.Error())
	default:
		h.Allowed = make([]interceptor.ClientID, 0)
		h.Participants = make(map[interceptor.ClientID]*Stat)
		return nil
	}
}

func (h *Health) IsParticipant(id interceptor.ClientID) bool {
	select {
	case <-h.ctx.Done():
		return false
	default:
		_, exists := h.Participants[id]
		return exists
	}
}

func (h *Health) IsAllowed(id interceptor.ClientID) bool {
	select {
	case <-h.ctx.Done():
		return false
	default:
		if len(h.Allowed) == 0 {
			return true
		}

		for _, allowedID := range h.Allowed {
			if allowedID == id {
				return true
			}
		}

		return false
	}
}
