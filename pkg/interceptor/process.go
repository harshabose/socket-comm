package interceptor

import "context"

type CanProcess interface {
	Process(context.Context, CanBeProcessed, State) error
}

type CanProcessBackground interface {
	ProcessBackground(context.Context, CanBeProcessedBackground, State) CanBeProcessedBackground
}

type Processor interface {
	CanProcess
	CanProcessBackground
}

type CanBeProcessed interface {
	Process(context.Context, CanProcess, State) error
}

type CanBeProcessedBackground interface {
	ProcessBackground(context.Context, CanProcessBackground, State) CanBeProcessedBackground
	Wait() error
	Stop()
}
