package keyexchange

import (
	"crypto/sha256"
	"fmt"
	"io"
	"time"

	"golang.org/x/crypto/hkdf"

	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/encryptionerr"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/types"
)

type Session struct {
	protocol    interfaces.Protocol
	state       interfaces.State
	createdAt   time.Time
	completedAt time.Time
}

type Manager struct {
	registry map[types.KeyExchangeProtocol]ProtocolFactory
	sessions map[types.KeyExchangeSessionID]*Session
}

func (m *Manager) Init(s interfaces.State, options ...interfaces.ProtocolFactoryOption) error {
	sessionID := s.GenerateKeyExchangeSessionID()
	_, exists := m.sessions[sessionID]
	if exists {
		return encryptionerr.ErrExchangeInProgress
	}

	factory, exists := m.registry[s.GetConfig().EncryptionProtocol.KeyExchangeProtocol]
	if !exists {
		return fmt.Errorf("%w: %s", encryptionerr.ErrProtocolNotFound, s.GetConfig().EncryptionProtocol.KeyExchangeProtocol)
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

func (m *Manager) Finalise(s interfaces.State) error {
	sessionID := s.GetKeyExchangeSessionID()
	session, exists := m.sessions[sessionID]
	if !exists {
		return encryptionerr.ErrExchangeNotComplete
	}

	ss, ok := s.(interfaces.KeySetter)
	if !ok {
		return encryptionerr.ErrInvalidInterceptor
	}

	p, ok := (session.protocol).(interfaces.KeyGetter)
	if !ok {
		return encryptionerr.ErrInvalidInterceptor
	}

	encKey, decKey, err := p.GetKeys()
	if err != nil {
		return err
	}

	return ss.SetKeys(encKey, decKey)
}

func (m *Manager) Process(msg interfaces.CanProcess, s interfaces.State) error {
	session, exists := m.sessions[s.GetKeyExchangeSessionID()]
	if !exists {
		return encryptionerr.ErrSessionNotFound
	}

	processor, ok := session.protocol.(interfaces.ProtocolProcessor)
	if !ok {
		return encryptionerr.ErrInvalidMessageType
	}

	return processor.Process(msg, s)
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
