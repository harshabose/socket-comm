package process

import (
	"context"
	"fmt"

	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
)

type IdentInit struct {
	ctx context.Context
}

func NewIdentInit(ctx context.Context) *IdentInit {
	return &IdentInit{
		ctx: ctx,
	}
}

func (p *IdentInit) Process(r interfaces.Processor, s *state.State) error {
	// TODO: SEND IDENT MESSAGE
	if err := s.Write(); err != nil {
		return err
	}

	waiter := NewIdentWaiter(p.ctx)
	if err := waiter.Process(r, s); err != nil {
		return fmt.Errorf("error while processing IdentInit process; err: %s", err.Error())
	}

	return nil
}
