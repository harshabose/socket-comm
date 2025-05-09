package room

import (
	"context"
	"fmt"
	"time"

	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type Room struct {
	// NOTE: MAYBE A CONFIG FOR ROOM?
	roomid       types.RoomID
	allowed      []types.ClientID
	participants map[types.ClientID]*state.State
	ttl          time.Duration
	cancel       context.CancelFunc
	ctx          context.Context
}

// TODO: DO I NEED MUX HERE?

func NewRoom(ctx context.Context, id types.RoomID, allowed []types.ClientID, ttl time.Duration) *Room {
	ctx2, cancel := context.WithTimeout(ctx, ttl)
	return &Room{
		ctx:          ctx2,
		cancel:       cancel,
		ttl:          ttl,
		roomid:       id,
		allowed:      allowed,
		participants: make(map[types.ClientID]*state.State),
	}
}

func (r *Room) ID() types.RoomID {
	return r.roomid
}

func (r *Room) TTL() time.Duration {
	return r.ttl
}

func (r *Room) Add(roomid types.RoomID, s *state.State) error {
	if roomid != r.roomid {
		return errors.ErrWrongRoom
	}

	id, err := s.GetClientID()
	if err != nil {
		return err
	}

	select {
	case <-r.ctx.Done():
		return fmt.Errorf("error while adding client to room. client ID: %s; room ID: %s; err: %s", id, r.roomid, errors.ErrContextCancelled.Error())
	default:
		if !r.isAllowed(id) {
			return fmt.Errorf("error while adding client to room. client ID: %s; room ID: %s; err: %s", id, r.roomid, errors.ErrClientNotAllowedInRoom.Error())
		}

		if r.isParticipant(id) {
			return fmt.Errorf("client with id '%s' already existing in the room with id %s; err: %s", id, r.roomid, errors.ErrClientIsAlreadyParticipant)
		}

		r.participants[id] = s
		return nil
	}
}

func (r *Room) isAllowed(id types.ClientID) bool {
	select {
	case <-r.ctx.Done():
		return false
	default:
		if len(r.allowed) == 0 {
			return true
		}

		for _, allowedID := range r.allowed {
			if allowedID == id {
				return true
			}
		}

		return false
	}
}

func (r *Room) forEachBoolean(f func(id types.ClientID) bool, ids ...types.ClientID) bool {
	if len(ids) == 0 {
		return false
	}

	for _, id := range ids {
		if !f(id) {
			return false
		}
	}

	return true
}

func (r *Room) Remove(roomid types.RoomID, s *state.State) error {
	if roomid != r.roomid {
		return errors.ErrWrongRoom
	}

	id, err := s.GetClientID()
	if err != nil {
		return err
	}

	select {
	case <-r.ctx.Done():
		return fmt.Errorf("error while removing client to room. client ID: %s; room ID: %s; err: %s", id, r.roomid, errors.ErrContextCancelled.Error())
	default:
		if !r.isAllowed(id) {
			return fmt.Errorf("error while removing client to room. client ID: %s; room ID: %s; err: %s", id, r.roomid, errors.ErrClientNotAllowedInRoom.Error())
		}

		if !r.isParticipant(id) {
			return fmt.Errorf("client with id '%s' does not exist in the room with id %s; err: %s", id, r.roomid, errors.ErrClientNotAParticipant.Error())
		}

		delete(r.participants, id)
		return nil
	}
}

func (r *Room) isParticipant(id types.ClientID) bool {
	select {
	case <-r.ctx.Done():
		return false
	default:
		_, exists := r.participants[id]
		return exists
	}
}

func (r *Room) WriteRoomMessage(roomid types.RoomID, msg message.Message, from types.ClientID, tos ...types.ClientID) error {
	select {
	case <-r.ctx.Done():
		return fmt.Errorf("error while sending message to peer in room; err: %s", errors.ErrContextCancelled.Error())
	default:
		if roomid != r.roomid {
			return errors.ErrWrongRoom
		}

		if len(tos) == 0 {
			return fmt.Errorf("atleast one receiver is need to use 'WriteRoomMessage' message")
		}

		if !r.forEachBoolean(r.isAllowed, append(tos, from)...) {
			return errors.ErrClientNotAllowedInRoom
		}

		if !r.forEachBoolean(r.isParticipant, append(tos, from)...) {
			return errors.ErrClientNotAParticipant
		}

		for _, to := range tos {
			if err := r.participants[to].Write(msg); err != nil {
				return fmt.Errorf("error while sending message to peer in room; err: %s", err.Error())
			}
		}

		return nil
	}
}

func (r *Room) GetParticipants() []types.ClientID {
	select {
	case <-r.ctx.Done():
		return make([]types.ClientID, 0) // EMPTY LIST
	default:
		clients := make([]types.ClientID, 0)
		for id, _ := range r.participants {
			clients = append(clients, id)
		}
		return clients
	}
}

func (r *Room) Close() error {
	r.cancel()
	r.participants = make(map[types.ClientID]*state.State)
	r.allowed = make([]types.ClientID, 0)
	return nil
}
