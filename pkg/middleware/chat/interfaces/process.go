package interfaces

import "github.com/harshabose/socket-comm/pkg/middleware/chat/state"

type CanProcess interface {
	Process(CanBeProcessed, *state.State) error
}

type CanProcessBackground interface {
	ProcessBackground(CanBeProcessedBackground, *state.State) CanBeProcessedBackground
}

type Processor interface {
	CanProcess
	CanProcessBackground
}

type CanBeProcessed interface {
	Process(Processor, *state.State) error
}

type CanBeProcessedBackground interface {
	ProcessBackground(Processor, *state.State) CanBeProcessedBackground
	Wait() error
	Stop()
}
