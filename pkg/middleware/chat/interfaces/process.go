package interfaces

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
)

type CanProcess interface {
	Process(context.Context, CanBeProcessed, *state.State) error
}

type CanProcessBackground interface {
	ProcessBackground(context.Context, CanBeProcessedBackground, *state.State) CanBeProcessedBackground
}

type Processor interface {
	CanProcess
	CanProcessBackground
}

type CanBeProcessed interface {
	Process(context.Context, Processor, *state.State) error
}

type CanBeProcessedBackground interface {
	ProcessBackground(context.Context, Processor, *state.State) CanBeProcessedBackground
	Wait() error
	Stop()
}
