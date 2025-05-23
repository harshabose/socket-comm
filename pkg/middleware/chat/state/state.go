package state

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
)

type State struct {
	id         interceptor.ClientID
	connection interceptor.Connection
	writer     interceptor.Writer
	reader     interceptor.Reader
	cancel     context.CancelFunc
	ctx        context.Context
}

func NewState(ctx context.Context, cancel context.CancelFunc, connection interceptor.Connection, writer interceptor.Writer, reader interceptor.Reader) *State {
	return &State{
		id:         interceptor.UnknownClientID,
		connection: connection,
		writer:     writer,
		reader:     reader,
		cancel:     cancel,
		ctx:        ctx,
	}
}

func (s *State) Ctx() context.Context {
	return s.ctx
}

func (s *State) GetClientID() (interceptor.ClientID, error) {
	if s.id == interceptor.UnknownClientID {
		return s.id, errors.ErrUnknownClientIDState
	}

	return s.id, nil
}

func (s *State) Write(ctx context.Context, msg message.Message) error {
	return s.writer.Write(ctx, s.connection, msg)
}

func (s *State) SetClientID(id interceptor.ClientID) error {
	if s.id != interceptor.UnknownClientID {
		return interceptor.ErrClientIDNotConsistent
	}

	s.id = id
	return nil
}
