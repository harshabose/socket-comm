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
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/types"
)

type State struct {
	InterceptorConfig    config.Config // copy of interceptor config
	PeerID               message.Receiver
	privKey              types.PrivateKey // also in protocol curve25519protocol.go
	salt                 types.Salt       // also in protocol curve25519protocol.go
	sessionID            types.SessionID  // encryption sessionID
	encryptor            encryptor.Encryptor
	Connection           interceptor.Connection
	Writer               interceptor.Writer
	Reader               interceptor.Reader
	cancel               context.CancelFunc
	ctx                  context.Context
	KeyExchangeSessionID types.KeyExchangeSessionID // Used for key exchange tracking
	mux                  sync.RWMutex
}

func NewState(ctx context.Context, cancel context.CancelFunc, config config.Config, connection interceptor.Connection, writer interceptor.Writer, reader interceptor.Reader) (*State, error) {
	newEncryptor, err := encryptor.NewEncryptor(config.EncryptionProtocol.CipherSuite)
	if err != nil {
		return nil, err
	}

	return &State{
		InterceptorConfig: config,
		peerID:            message.UnknownReceiver,
		privKey:           types.PrivateKey{},
		salt:              types.Salt{},
		encryptor:         newEncryptor,
		Connection:        connection,
		Writer:            writer,
		Reader:            reader,
		cancel:            cancel,
		ctx:               ctx,
	}, nil
}

func (s *State) GenerateKeyExchangeSessionID() types.KeyExchangeSessionID {
	if s.KeyExchangeSessionID != "" {
		fmt.Println("KeyExchangeSessionID already exists; creating new")
	}
	s.KeyExchangeSessionID = types.KeyExchangeSessionID(uuid.NewString())

	return s.KeyExchangeSessionID
}
