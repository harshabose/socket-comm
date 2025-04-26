package keyexchange

import (
	"crypto/rand"
	"io"

	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/ed25519"

	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/encryptionerr"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/messages"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/state"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/types"
)

type Curve25519Protocol struct {
	privKey      types.PrivateKey // also in protocol curve25519protocol.go
	pubKey       types.PublicKey
	peerPubKey   types.PublicKey
	salt         types.Salt // also in protocol curve25519protocol.go
	sessionID    types.SessionID
	sharedSecret []byte
	encKey       types.Key
	decKey       types.Key
	state        SessionState
	options      Curve25519Options

	// TODO: mutex is needed here for SessionState and Keys
}

type Curve25519Options struct {
	SigningKey       ed25519.PrivateKey
	VerificationKey  ed25519.PublicKey
	RequireSignature bool
}

func (p *Curve25519Protocol) Init(s *state.State) error {
	if _, err := io.ReadFull(rand.Reader, p.privKey[:]); err != nil {
		return err
	}

	curve25519.ScalarBaseMult((*[32]byte)(&p.pubKey), (*[32]byte)(&p.privKey))

	if s.InterceptorConfig.IsServer && p.options.RequireSignature {
		if _, err := io.ReadFull(rand.Reader, p.salt[:]); err != nil {
			p.state = SessionStateError
			return err
		}

		if _, err := io.ReadFull(rand.Reader, p.sessionID[:]); err != nil {
			p.state = SessionStateError
			return err
		}

		sign := ed25519.Sign(p.options.SigningKey, append(p.pubKey[:], p.salt[:]...))

		if err := s.Writer.Write(s.Connection, messages.NewInit("", s.PeerID, p.pubKey, sign, p.sessionID, p.salt)); err != nil {
			p.state = SessionStateError
			return err
		}
	}

	p.state = SessionStateInProgress
	return nil
}

func (p *Curve25519Protocol) GetKeys() (encKey types.Key, decKey types.Key, err error) {
	if p.state != SessionStateCompleted {
		return types.Key{}, types.Key{}, encryptionerr.ErrExchangeNotComplete
	}

	return p.encKey, p.decKey, nil
}

func (p *Curve25519Protocol) GetState() SessionState {
	return p.state
}

func (p *Curve25519Protocol) IsComplete() bool {
	return p.state == SessionStateCompleted
}

func (p *Curve25519Protocol) Process(msg MessageProcessor, s *state.State) error {
	if err := msg.Process(p, s); err != nil {
		p.state = SessionStateError
		return err
	}

	return nil
}
