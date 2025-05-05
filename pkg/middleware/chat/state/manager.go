package state

import (
	"sync"

	"github.com/harshabose/socket-comm/internal/util"
	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
)

type Manager struct {
	states map[interceptor.Connection]interfaces.State
	mux    sync.RWMutex
}

func (m *Manager) GetState(connection interceptor.Connection) (interfaces.State, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	state, exists := m.states[connection]
	if !exists {
		return nil, errors.ErrConnectionNotFound
	}

	return state, nil
}

func (m *Manager) SetState(connection interceptor.Connection, s interfaces.State) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	if _, exists := m.states[connection]; exists {
		return errors.ErrConnectionExists
	}

	m.states[connection] = s
	return nil
}

// RemoveState removes a connection's state
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

// ForEach executes the provided function for each state in the manager
func (m *Manager) ForEach(fn func(connection interceptor.Connection, state interfaces.State) error) error {
	m.mux.RLock()
	defer m.mux.RUnlock()

	var errs util.MultiError
	for conn, state := range m.states {
		if err := fn(conn, state); err != nil {
			errs.Add(err)
		}
	}

	return errs.ErrorOrNil()
}
