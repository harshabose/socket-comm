package process

import (
	"context"
	"fmt"

	"github.com/harshabose/socket-comm/pkg/interceptor"
)

type IdentInit struct {
	AsyncProcess
}

func NewIdentInit() *IdentInit {
	return &IdentInit{}
}

func (p *IdentInit) Process(ctx context.Context, _ interceptor.CanProcess, s interceptor.State) error {
	// TODO: SEND IDENT MESSAGE // PROBLEM HERE AS PROCESS MODULE CANNOT IMPORT MESSAGE
	if err := s.Write(nil); err != nil {
		return err
	}

	waiter := NewIdentWaiter()
	if err := waiter.Process(ctx, nil, s); err != nil {
		return fmt.Errorf("error while processing IdentInit process; err: %s", err.Error())
	}

	return nil
}
