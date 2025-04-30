package messages

import (
	"fmt"

	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/ed25519"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/encryptionerr"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/keyexchange"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/state"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/types"
)

// Protocol constants
const (
	InitProtocol     message.Protocol = "curve25519.init"
	ResponseProtocol message.Protocol = "curve25519.response"
	ConfirmProtocol  message.Protocol = "curve25519.confirm"
)

type Init struct {
	interceptor.BaseMessage
	PublicKey types.PublicKey `json:"public_key"`
	Signature []byte          `json:"signature"`
	SessionID types.SessionID `json:"session_id"`
	Salt      types.Salt      `json:"salt"`
}

func NewInit(sender message.Sender, receiver message.Receiver, key types.PublicKey, sign []byte, sessionID types.SessionID, salt types.Salt) *Init {
	return &Init{
		BaseMessage: interceptor.NewBaseMessage(InitProtocol, sender, receiver),
		PublicKey:   key,
		Signature:   sign,
		SessionID:   sessionID,
		Salt:        salt,
	}
}

func (m *Init) WriteProcess(_ interceptor.Interceptor, _ interceptor.Connection) error {
	return nil
}

func (m *Init) ReadProcess(_interceptor interceptor.Interceptor, connection interceptor.Connection) error {
	s, err := encrypt.GetState(_interceptor, connection)
	if err != nil {
		return err
	}

	return i.keyExchangeManager.Process(s, m)
}

func (m *Init) Process(protocol keyexchange.Protocol, s *state.State) error {
	p, ok := protocol.(*keyexchange.Curve25519Protocol)
	if !ok {
		return encryptionerr.ErrInvalidMessageType
	}

	if p.GetState() != keyexchange.SessionStateInitial {
		return encryptionerr.ErrInvalidSessionState
	}

	sign := append(m.PublicKey[:], m.Salt[:]...)
	if !ed25519.Verify(p.options.VerificationKey, sign, m.Signature) {
		return encryptionerr.ErrInvalidSignature
	}

	p.salt = m.Salt
	shared, err := curve25519.X25519(p.privKey[:], m.PublicKey[:])
	if err != nil {
		return fmt.Errorf("failed to compute shared secret: %w", err)
	}

	encKey, decKey, err := keyexchange.Derive(shared, p.salt, "") // TODO: ADD INFO STRING
	if err != nil {
		return fmt.Errorf("key derivation failed: %w", err)
	}

	p.encKey = encKey
	p.decKey = decKey
	p.sessionID = m.SessionID

	if err := s.Writer.Write(s.Connection, nil); err != nil {
		return err
	} // TODO: ADD RESPONSE MESSAGE

	p.state = SessionStateCompleted
	return nil
}
