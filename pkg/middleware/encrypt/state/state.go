package state

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/config"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/encryptor"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/types"
)

type State struct {
	currentConfig        config.Config // copy of interceptor config
	peerID               message.Receiver
	privKey              types.PrivateKey          // also in protocol curve25519protocol.go
	salt                 types.Salt                // also in protocol curve25519protocol.go
	encryptSessionID     types.EncryptionSessionID // encryption encryptSessionID
	encryptor            interfaces.Encryptor
	connection           interceptor.Connection
	writer               interceptor.Writer
	reader               interceptor.Reader
	cancel               context.CancelFunc
	ctx                  context.Context
	keyExchangeSessionID types.KeyExchangeSessionID // Used for key exchange tracking
	mux                  sync.RWMutex
}

func NewState(ctx context.Context, cancel context.CancelFunc, config config.Config, connection interceptor.Connection, writer interceptor.Writer, reader interceptor.Reader) (*State, error) {
	newEncryptor, err := encryptor.NewEncryptor(config.EncryptionProtocol.CipherSuite)
	if err != nil {
		return nil, err
	}

	return &State{
		currentConfig: config,
		peerID:        message.UnknownReceiver,
		privKey:       types.PrivateKey{},
		salt:          types.Salt{},
		encryptor:     newEncryptor,
		connection:    connection,
		writer:        writer,
		reader:        reader,
		cancel:        cancel,
		ctx:           ctx,
	}, nil
}

func (s *State) GenerateKeyExchangeSessionID() types.KeyExchangeSessionID {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.keyExchangeSessionID != "" {
		fmt.Println("keyExchangeSessionID already exists; creating new")
	}
	s.keyExchangeSessionID = types.KeyExchangeSessionID(uuid.NewString())

	return s.keyExchangeSessionID
}

func (s *State) GetConnection() interceptor.Connection {
	return s.connection
}

func (s *State) WriteMessage(msg interceptor.Message) error {
	s.mux.Lock()
	defer s.mux.Unlock()
	// TODO: MANAGE CLIENT DISCOVERY
	return s.writer.Write(s.connection, msg)
}

func (s *State) ReadMessage(msg interceptor.Message) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.reader.Read()
}

func (s *State) GetKeyExchangeSessionID() types.KeyExchangeSessionID {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.keyExchangeSessionID
}

func (s *State) GetConfig() config.Config {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.currentConfig
}

func (s *State) SetKeys(encKey, decKey types.Key) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.encryptor.SetKeys(encKey, decKey)
}

func (s *State) Decrypt(msg message.Message) (message.Message, error) {
	return s.encryptor.Decrypt(msg)
}
