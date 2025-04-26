package keyexchange

import (
	"crypto/sha256"
	"fmt"
	"io"
	"time"

	"golang.org/x/crypto/hkdf"

	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/encryptionerr"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/state"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/types"
)

type Session struct {
	protocol    Protocol
	state       *state.State
	createdAt   time.Time
	completedAt time.Time
}

type Manager struct {
	registry map[types.KeyExchangeProtocol]ProtocolFactory
	sessions map[types.KeyExchangeSessionID]*Session
}

func (m *Manager) Init(s *state.State, options ...ProtocolFactoryOption) error {
	sessionID := s.GenerateKeyExchangeSessionID()
	_, exists := m.sessions[sessionID]
	if exists {
		return encryptionerr.ErrExchangeInProgress
	}

	factory, exists := m.registry[s.InterceptorConfig.EncryptionProtocol.KeyExchangeProtocol]
	if !exists {
		return fmt.Errorf("%w: %s", encryptionerr.ErrProtocolNotFound, s.InterceptorConfig.EncryptionProtocol.KeyExchangeProtocol)
	}

	p, err := factory(options...)
	if err != nil {
		return err
	}

	if err := p.Init(s); err != nil {
		return err
	}

	m.sessions[sessionID] = &Session{
		protocol:  p,
		state:     s,
		createdAt: time.Now(),
	}

	return nil
}

func (m *Manager) Process(s *state.State, msg MessageProcessor) error {
	session, exists := m.sessions[s.KeyExchangeSessionID]
	if !exists {
		return encryptionerr.ErrSessionNotFound
	}

	return session.protocol.Process(msg, s)
}

// Derive generates encryption keys from shared secret
func Derive(shared []byte, salt types.Salt, info string) (types.Key, types.Key, error) {
	hkdfReader := hkdf.New(sha256.New, shared, salt[:], []byte(info))

	key1 := types.Key{}
	if _, err := io.ReadFull(hkdfReader, key1[:]); err != nil {
		return types.Key{}, types.Key{}, err
	}

	key2 := types.Key{}
	if _, err := io.ReadFull(hkdfReader, key2[:]); err != nil {
		return types.Key{}, types.Key{}, err
	}

	return key1, key2, nil
}
