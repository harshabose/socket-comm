package state

import (
	"sync"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/encryptionerr"
)

type Manager struct {
	states map[interceptor.Connection]*State
	mux    sync.RWMutex
}

func (m *Manager) GetState(connection interceptor.Connection) (*State, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	state, exists := m.states[connection]
	if !exists {
		return nil, encryptionerr.ErrConnectionNotFound
	}

	return state, nil
}

func (m *Manager) SetState(connection interceptor.Connection, s *State) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	if _, exists := m.states[connection]; exists {
		return encryptionerr.ErrConnectionExists
	}

	m.states[connection] = s
	return nil
}

// RemoveState removes a Connection's state
func (m *Manager) RemoveState(connection interceptor.Connection) (*State, error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	state, exists := m.states[connection]
	if !exists {
		return nil, encryptionerr.ErrConnectionNotFound
	}

	delete(m.states, connection)
	return state, nil
}

// ForEach executes the provided function for each state in the manager
func (m *Manager) ForEach(fn func(connection interceptor.Connection, state *State) error) []error {
	m.mux.RLock()
	defer m.mux.RUnlock()

	var errs []error
	for conn, state := range m.states {
		if err := fn(conn, state); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}
