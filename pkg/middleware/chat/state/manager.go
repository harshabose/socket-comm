package state

import (
	"sync"

	"github.com/harshabose/socket-comm/internal/util"
	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
)

type Manager struct {
	states map[interceptor.Connection]*State
	mux    sync.RWMutex
}

func (m *Manager) GetState(connection interceptor.Connection) (*State, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	s, exists := m.states[connection]
	if !exists {
		return nil, errors.ErrConnectionNotFound
	}

	return s, nil
}

func (m *Manager) SetState(connection interceptor.Connection, s *State) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	if _, exists := m.states[connection]; exists {
		return errors.ErrConnectionExists
	}

	m.states[connection] = s
	return nil
}

// RemoveState removes a connection's State
func (m *Manager) RemoveState(connection interceptor.Connection) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	_, exists := m.states[connection]
	if !exists {
		return errors.ErrConnectionNotFound
	}

	delete(m.states, connection)
	return nil
}

// ForEach executes the provided function for each State in the manager
func (m *Manager) ForEach(fn func(connection interceptor.Connection, state *State) error) error {
	m.mux.RLock()
	defer m.mux.RUnlock()

	var errs util.MultiError
	for conn, s := range m.states {
		if err := fn(conn, s); err != nil {
			errs.Add(err)
		}
	}

	return errs.ErrorOrNil()
}
