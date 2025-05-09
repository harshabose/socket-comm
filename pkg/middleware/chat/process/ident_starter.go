package process

import (
	"context"
	"fmt"

	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
)

type IdentInit struct {
	AsyncProcess
}

func NewIdentInit() *IdentInit {
	return &IdentInit{}
}

func (p *IdentInit) Process(ctx context.Context, _ interfaces.Processor, s *state.State) error {
	// TODO: SEND IDENT MESSAGE
	if err := s.Write(nil); err != nil {
		return err
	}

	waiter := NewIdentWaiter()
	if err := waiter.Process(ctx, nil, s); err != nil {
		return fmt.Errorf("error while processing IdentInit process; err: %s", err.Error())
	}

	return nil
}
