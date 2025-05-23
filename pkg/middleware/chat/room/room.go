package room

import (
	"context"
	"fmt"
	"time"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type Room struct {
	// NOTE: MAYBE A CONFIG FOR ROOM?
	roomid          types.RoomID
	allowed         []interceptor.ClientID
	participants    map[interceptor.ClientID]interceptor.State
	ttl             time.Duration
	isHealthTracked bool
	cancel          context.CancelFunc
	ctx             context.Context
}

// TODO: DO I NEED MUX HERE?
// UPDATE: YES
// TODO: ADD SOME VALIDATION BEFORE CREATING THE ROOM

func NewRoom(ctx context.Context, id types.RoomID, allowed []interceptor.ClientID, ttl time.Duration) *Room {
	ctx2, cancel := context.WithTimeout(ctx, ttl)
	return &Room{
		ctx:          ctx2,
		cancel:       cancel,
		ttl:          ttl,
		roomid:       id,
		allowed:      allowed,
		participants: make(map[interceptor.ClientID]interceptor.State),
	}
}

func (r *Room) Ctx() context.Context {
	return r.ctx
}

func (r *Room) ID() types.RoomID {
	return r.roomid
}

func (r *Room) TTL() time.Duration {
	return r.ttl
}

func (r *Room) Add(roomid types.RoomID, s interceptor.State) error {
	if roomid != r.roomid {
		return errors.ErrWrongRoom
	}

	id, err := s.GetClientID()
	if err != nil {
		return err
	}

	select {
	case <-r.ctx.Done():
		return fmt.Errorf("error while adding client to room. client id: %s; room id: %s; err: %s", id, r.roomid, interceptor.ErrContextCancelled.Error())
	default:
		if !r.isAllowed(id) {
			return fmt.Errorf("error while adding client to room. client id: %s; room id: %s; err: %s", id, r.roomid, errors.ErrClientNotAllowed.Error())
		}

		if r.isParticipant(id) {
			return fmt.Errorf("client with id '%s' already existing in the room with id %s; err: %s", id, r.roomid, errors.ErrClientIsAlreadyParticipant)
		}

		r.participants[id] = s
		return nil
	}
}

func (r *Room) isAllowed(id interceptor.ClientID) bool {
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

func (r *Room) forEachBoolean(f func(id interceptor.ClientID) bool, ids ...interceptor.ClientID) bool {
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

func (r *Room) Remove(roomid types.RoomID, s interceptor.State) error {
	if roomid != r.roomid {
		return errors.ErrWrongRoom
	}

	id, err := s.GetClientID()
	if err != nil {
		return err
	}

	select {
	case <-r.ctx.Done():
		return fmt.Errorf("error while removing client to room. client id: %s; room id: %s; err: %s", id, r.roomid, interceptor.ErrContextCancelled.Error())
	default:
		if !r.isAllowed(id) {
			return fmt.Errorf("error while removing client to room. client id: %s; room id: %s; err: %s", id, r.roomid, errors.ErrClientNotAllowed.Error())
		}

		if !r.isParticipant(id) {
			return fmt.Errorf("client with id '%s' does not exist in the room with id %s; err: %s", id, r.roomid, errors.ErrClientNotAParticipant.Error())
		}

		delete(r.participants, id)
		return nil
	}
}

func (r *Room) isParticipant(id interceptor.ClientID) bool {
	select {
	case <-r.ctx.Done():
		return false
	default:
		_, exists := r.participants[id]
		return exists
	}
}

func (r *Room) WriteRoomMessage(roomid types.RoomID, msg message.Message, from interceptor.ClientID, tos ...interceptor.ClientID) error {
	select {
	case <-r.ctx.Done():
		return fmt.Errorf("error while sending message to peer in room; err: %s", interceptor.ErrContextCancelled.Error())
	default:
		if roomid != r.roomid {
			return errors.ErrWrongRoom
		}

		if len(tos) == 0 {
			return fmt.Errorf("atleast one receiver is need to use 'WriteRoomMessage' message")
		}

		if !r.forEachBoolean(r.isAllowed, append(tos, from)...) {
			return errors.ErrClientNotAllowed
		}

		if !r.forEachBoolean(r.isParticipant, append(tos, from)...) {
			return errors.ErrClientNotAParticipant
		}

		for _, to := range tos {
			if err := r.participants[to].Write(r.ctx, msg); err != nil {
				return fmt.Errorf("error while sending message to peer in room; err: %s", err.Error())
			}
		}

		return nil
	}
}

func (r *Room) StartHealthTracking(roomid types.RoomID) error {
	select {
	case <-r.ctx.Done():
		return fmt.Errorf("error while marking room as health tracked. room id: %s; err: %s", roomid, interceptor.ErrContextCancelled.Error())
	default:
		if roomid != r.roomid {
			return errors.ErrWrongRoom
		}

		if r.IsRoomMarkedForHealthTracking() {
			return fmt.Errorf("room with id %s is already health tracked", roomid)
		}

		r.isHealthTracked = true
		return nil
	}
}

func (r *Room) IsRoomMarkedForHealthTracking() bool {
	return r.isHealthTracked
}

func (r *Room) UnMarkRoomForHealthTracking() error {
	select {
	case <-r.ctx.Done():
		return fmt.Errorf("error while unmarking room as health tracked. room id: %s; err: %s", r.roomid, interceptor.ErrContextCancelled.Error())
	default:
		if !r.IsRoomMarkedForHealthTracking() {
			return fmt.Errorf("room with id %s is not health tracked", r.roomid)
		}

		r.isHealthTracked = false
		return nil
	}
}

func (r *Room) GetParticipants() []interceptor.ClientID {
	select {
	case <-r.ctx.Done():
		return make([]interceptor.ClientID, 0) // EMPTY LIST
	default:
		clients := make([]interceptor.ClientID, 0)
		for id := range r.participants {
			clients = append(clients, id)
		}
		return clients
	}
}

func (r *Room) GetAllowed() []interceptor.ClientID {
	return r.allowed
}

func (r *Room) Close() error {
	r.cancel()
	r.participants = make(map[interceptor.ClientID]interceptor.State)
	r.allowed = make([]interceptor.ClientID, 0)
	return nil
}
