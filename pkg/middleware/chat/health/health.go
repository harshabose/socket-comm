package health

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

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
	mux sync.RWMutex
	ctx context.Context
}

type Snapshot struct {
	Roomid       types.RoomID             `json:"roomid"`       // NOTE: ADDED WHEN HEALTH IS CREATED.
	Allowed      []types.ClientID         `json:"allowed"`      // NOTE: ADDED WHEN HEALTH IS CREATED.
	Participants map[types.ClientID]*Stat `json:"participants"` // NOTE: ADDED WHEN CLIENT JOINS. UPDATED WHEN CLIENT SENDS HEALTH RESPONSE.
}

// Marshal marshals the health struct into a json byte array.
// NOTE: THIS IS USED TO SEND HEALTH STATS TO THE CLIENT.
func (h *Snapshot) Marshal() ([]byte, error) {
	return json.Marshal(h)
}

func NewHealth(ctx context.Context, id types.RoomID, allowed []types.ClientID) *Health {
	return &Health{
		Snapshot: Snapshot{
			Roomid:       id,
			Allowed:      allowed,
			Participants: make(map[types.ClientID]*Stat),
		},
		ctx: ctx,
	}
}

func (h *Health) ID() types.RoomID {
	return h.Roomid
}

func (h *Health) Add(roomid types.RoomID, id types.ClientID) error {
	h.mux.Lock()
	defer h.mux.Unlock()

	if roomid != h.Roomid {
		return errors.ErrWrongRoom
	}

	select {
	case <-h.ctx.Done():
		return fmt.Errorf("error while adding client to health stats. client ID: %s; room ID: %s; err: %s", id, h.Roomid, errors.ErrContextCancelled.Error())
	default:
		if !h.isAllowed(id) {
			return fmt.Errorf("error while adding client to health stats. client ID: %s; room ID: %s; err: %s", id, h.Roomid, errors.ErrClientNotAllowedInRoom.Error())
		}

		if h.isParticipant(id) {
			return fmt.Errorf("client with id '%s' already existing in the health stats with id %s; err: %s", id, h.Roomid, errors.ErrClientIsAlreadyParticipant)
		}

		h.Participants[id] = &Stat{}
		return nil
	}
}

func (h *Health) Remove(roomid types.RoomID, id types.ClientID) error {
	h.mux.Lock()
	defer h.mux.Unlock()

	if roomid != h.Roomid {
		return errors.ErrWrongRoom
	}

	select {
	case <-h.ctx.Done():
		return fmt.Errorf("error while removing client to health stats. client ID: %s; room ID: %s; err: %s", id, h.Roomid, errors.ErrContextCancelled.Error())
	default:
		if !h.isAllowed(id) {
			return fmt.Errorf("error while removing client to health stats. client ID: %s; room ID: %s; err: %s", id, h.Roomid, errors.ErrClientNotAllowedInRoom.Error())
		}

		if !h.isParticipant(id) {
			return fmt.Errorf("client with id '%s' does not exist in the health stats with id %s; err: %s", id, h.Roomid, errors.ErrClientNotAParticipant.Error())
		}

		delete(h.Participants, id)
		return nil
	}
}

func (h *Health) Update(roomid types.RoomID, id types.ClientID, s *Stat) error {
	h.mux.Lock()
	defer h.mux.Unlock()

	if roomid != h.Roomid {
		return errors.ErrWrongRoom
	}

	select {
	case <-h.ctx.Done():
		return fmt.Errorf("error while adding client to health stats. client ID: %s; room ID: %s; err: %s", id, h.Roomid, errors.ErrContextCancelled.Error())
	default:
		if !h.isAllowed(id) {
			return fmt.Errorf("error while adding client to health stats. client ID: %s; room ID: %s; err: %s", id, h.Roomid, errors.ErrClientNotAllowedInRoom.Error())
		}

		if !h.isParticipant(id) {
			return fmt.Errorf("client with id '%s' does not exist in the health stats with id %s; err: %s", id, h.Roomid, errors.ErrClientNotAParticipant.Error())
		}

		h.Participants[id] = s
		return nil
	}
}

func (h *Health) isAllowed(id types.ClientID) bool {
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

func (h *Health) isParticipant(id types.ClientID) bool {
	select {
	case <-h.ctx.Done():
		return false
	default:
		_, exists := h.Participants[id]
		return exists
	}
}
