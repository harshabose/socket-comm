package state

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type state struct {
	id         types.ClientID
	connection interceptor.Connection
	writer     interceptor.Writer
	reader     interceptor.Reader
	cancel     context.CancelFunc
	ctx        context.Context
}

func NewState(ctx context.Context, cancel context.CancelFunc, connection interceptor.Connection, writer interceptor.Writer, reader interceptor.Reader) interfaces.State {
	return &state{
		id:         types.UnKnownClient,
		connection: connection,
		writer:     writer,
		reader:     reader,
		cancel:     cancel,
		ctx:        ctx,
	}
}

func (s *state) Ctx() context.Context {
	return s.ctx
}

func (s *state) GetClientID() (types.ClientID, error) {
	if s.id == types.UnKnownClient {
		return s.id, errors.ErrUnknownClientIDState
	}

	return s.id, nil
}

func (s *state) Write(msg message.Message) error {
	return s.writer.Write(s.connection, msg)
}

func (s *state) SetClientID(id types.ClientID) error {
	if s.id != types.UnKnownClient {
		return errors.ErrClientIDNotConsistent
	}

	s.id = id
	return nil
}
