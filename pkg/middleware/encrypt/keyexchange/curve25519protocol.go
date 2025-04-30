package keyexchange

import (
	"crypto/rand"
	"io"

	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/ed25519"

	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/encryptionerr"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/types"
)

type Curve25519Protocol struct {
	privKey      types.PrivateKey // also in protocol curve25519protocol.go
	pubKey       types.PublicKey
	peerPubKey   types.PublicKey
	salt         types.Salt // also in protocol curve25519protocol.go
	sessionID    types.EncryptionSessionID
	sharedSecret []byte
	encKey       types.Key
	decKey       types.Key
	state        types.SessionState
	options      Curve25519Options

	// TODO: mutex is needed here for SessionState and Keys
}

type Curve25519Options struct {
	SigningKey       ed25519.PrivateKey
	VerificationKey  ed25519.PublicKey
	RequireSignature bool
}

func (p *Curve25519Protocol) Init(s interfaces.State) error {
	if _, err := io.ReadFull(rand.Reader, p.privKey[:]); err != nil {
		return err
	}

	curve25519.ScalarBaseMult((*[32]byte)(&p.pubKey), (*[32]byte)(&p.privKey))

	if s.GetConfig().IsServer && p.options.RequireSignature {
		if _, err := io.ReadFull(rand.Reader, p.salt[:]); err != nil {
			p.state = types.SessionStateError
			return err
		}

		if _, err := io.ReadFull(rand.Reader, p.sessionID[:]); err != nil {
			p.state = types.SessionStateError
			return err
		}

		sign := ed25519.Sign(p.options.SigningKey, append(p.pubKey[:], p.salt[:]...))

		msg, err := NewInit(p.pubKey, sign, p.sessionID, p.salt)
		if err != nil {
			return err
		}

		if err := s.WriteMessage(msg); err != nil {
			p.state = types.SessionStateError
			return err
		}
	}

	p.state = types.SessionStateInitial
	return nil
}

func (p *Curve25519Protocol) GetKeys() (encKey types.Key, decKey types.Key, err error) {
	if p.state != types.SessionStateCompleted {
		return types.Key{}, types.Key{}, encryptionerr.ErrExchangeNotComplete
	}

	return p.encKey, p.decKey, nil
}

func (p *Curve25519Protocol) GetState() types.SessionState {
	return p.state
}

func (p *Curve25519Protocol) IsComplete() bool {
	return p.state == types.SessionStateCompleted
}

func (p *Curve25519Protocol) Process(msg interfaces.CanProcess, s interfaces.State) error {
	if err := msg.Process(p, s); err != nil {
		p.state = types.SessionStateError
		return err
	}

	return nil
}
