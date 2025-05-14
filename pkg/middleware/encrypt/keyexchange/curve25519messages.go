package keyexchange

import (
	"fmt"
	"time"

	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/ed25519"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/encryptionerr"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/types"
)

type Init struct {
	// TODO: MANAGE STATE USING KEY EXCHANGE SESSION id
	interceptor.BaseMessage
	PublicKey types.PublicKey           `json:"public_key"`
	Signature []byte                    `json:"signature"`
	SessionID types.EncryptionSessionID `json:"session_id"`
	Salt      types.Salt                `json:"salt"`
}

func NewInit(pubKey types.PublicKey, sign []byte, sessionID types.EncryptionSessionID, salt types.Salt) (*Init, error) {
	msg := &Init{
		PublicKey: pubKey,
		Signature: sign,
		SessionID: sessionID,
		Salt:      salt,
	}

	bmsg, err := interceptor.NewBaseMessage(message.NoneProtocol, nil, msg)
	if err != nil {
		return nil, err
	}
	msg.BaseMessage = bmsg

	return msg, nil
}

func (m *Init) WriteProcess(_ interceptor.Interceptor, _ interceptor.Connection) error {
	return nil
}

func (m *Init) ReadProcess(_interceptor interceptor.Interceptor, connection interceptor.Connection) error {
	ss, ok := _interceptor.(interfaces.CanGetState)
	if !ok {
		return encryptionerr.ErrInvalidInterceptor
	}

	s, err := ss.GetState(connection)
	if err != nil {
		return err
	}

	pp, ok := _interceptor.(interfaces.ProtocolProcessor)
	if !ok {
		return encryptionerr.ErrInvalidInterceptor
	}

	return pp.Process(m, s)
}

func (m *Init) Process(protocol interfaces.Protocol, s interfaces.State) error {
	p, ok := protocol.(*Curve25519Protocol)
	if !ok {
		return encryptionerr.ErrInvalidMessageType
	}

	if p.state != types.SessionStateInitial {
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

	encKey, decKey, err := Derive(shared, p.salt, "") // TODO: ADD INFO STRING
	if err != nil {
		return fmt.Errorf("key derivation failed: %w", err)
	}

	p.encKey = encKey
	p.decKey = decKey
	p.sessionID = m.SessionID

	msg, err := NewResponse(p.pubKey)
	if err != nil {
		return err
	}

	if err := s.WriteMessage(msg); err != nil {
		return err
	}

	p.state = types.SessionStateInProgress
	return nil
}

type Response struct {
	interceptor.BaseMessage
	PublicKey types.PublicKey `json:"public_key"`
}

func NewResponse(pubKey types.PublicKey) (*Response, error) {
	msg := &Response{
		PublicKey: pubKey,
	}

	bmsg, err := interceptor.NewBaseMessage(message.NoneProtocol, nil, msg)
	if err != nil {
		return nil, err
	}
	msg.BaseMessage = bmsg

	return msg, nil
}

func (m *Response) WriteProcess(_ interceptor.Interceptor, _ interceptor.Connection) error {
	return nil
}

func (m *Response) ReadProcess(_i interceptor.Interceptor, connection interceptor.Connection) error {
	ss, ok := _i.(interfaces.CanGetState)
	if !ok {
		return encryptionerr.ErrInvalidInterceptor
	}

	s, err := ss.GetState(connection)
	if err != nil {
		return err
	}

	pp, ok := _i.(interfaces.ProtocolProcessor)
	if !ok {
		return encryptionerr.ErrInvalidInterceptor
	}

	return pp.Process(m, s)
}

func (m *Response) Process(protocol interfaces.Protocol, s interfaces.State) error {
	p, ok := protocol.(*Curve25519Protocol)
	if !ok {
		return encryptionerr.ErrInvalidMessageType
	}

	if p.state != types.SessionStateInitial {
		return encryptionerr.ErrInvalidSessionState
	}

	shared, err := curve25519.X25519(p.privKey[:], m.PublicKey[:])
	if err != nil {
		return fmt.Errorf("failed to compute shared secret: %w", err)
	}

	decKey, encKey, err := Derive(shared, p.salt, "") // TODO: ADD INFO STRING
	if err != nil {
		return fmt.Errorf("key derivation failed: %w", err)
	}

	p.encKey = encKey
	p.decKey = decKey

	msg, err := NewDone()
	if err != nil {
		return err
	}

	if err := s.WriteMessage(msg); err != nil {
		return err
	}

	p.state = types.SessionStateInProgress
	return nil
}

type Done struct {
	interceptor.BaseMessage
	Timestamp time.Time `json:"timestamp"`
}

func NewDone() (*Done, error) {
	msg := &Done{
		Timestamp: time.Now(),
	}
	bmsg, err := interceptor.NewBaseMessage(message.NoneProtocol, nil, msg)
	if err != nil {
		return nil, err
	}
	msg.BaseMessage = bmsg

	return msg, nil
}

func (m *Done) WriteProcess(_ interceptor.Interceptor, _ interceptor.Connection) error {
	return nil
}

func (m *Done) ReadProcess(_i interceptor.Interceptor, connection interceptor.Connection) error {
	ss, ok := _i.(interfaces.CanGetState)
	if !ok {
		return encryptionerr.ErrInvalidInterceptor
	}

	s, err := ss.GetState(connection)
	if err != nil {
		return err
	}

	pp, ok := _i.(interfaces.ProtocolProcessor)
	if !ok {
		return encryptionerr.ErrInvalidInterceptor
	}

	return pp.Process(m, s)
}

func (m *Done) Process(protocol interfaces.Protocol, s interfaces.State) error {
	p, ok := protocol.(*Curve25519Protocol)
	if !ok {
		return encryptionerr.ErrInvalidMessageType
	}

	msg, err := NewDoneResponse()
	if err != nil {
		return err
	}

	if err := s.WriteMessage(msg); err != nil {
		return err
	}

	p.state = types.SessionStateCompleted
	return nil
}

type DoneResponse struct {
	Done
}

func NewDoneResponse() (*DoneResponse, error) {
	msg := &DoneResponse{
		Done: Done{
			Timestamp: time.Now(),
		},
	}
	bmsg, err := interceptor.NewBaseMessage(message.NoneProtocol, nil, msg)
	if err != nil {
		return nil, err
	}
	msg.BaseMessage = bmsg

	return msg, nil
}

// TODO: ADD WRITE OR READ PROCESS METHODS
// TODO: ADD PROTOCOLS

func (m *DoneResponse) Process(protocol interfaces.Protocol, _ interfaces.State) error {
	p, ok := protocol.(*Curve25519Protocol)
	if !ok {
		return encryptionerr.ErrInvalidMessageType
	}

	p.state = types.SessionStateCompleted
	return nil
}
