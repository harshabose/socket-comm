package process

import (
	"context"
	"fmt"
	"sync"

	"github.com/harshabose/socket-comm/pkg/interceptor"
)

// AsyncProcess is intended to be embedded in a process (can be message-tagged) to enable async capabilities
// NOTE: TO EMBED THIS, THE EMBEDDER NEEDS TO IMPLEMENT interfaces.CanBeProcessed. THIS CONTRACT IS LEFT TO THE DEVELOPER TO FULFILL
type AsyncProcess struct {
	interceptor.CanBeProcessed
	err    error
	done   chan struct{}
	ctx    context.Context
	cancel context.CancelFunc
	mux    sync.RWMutex
}

func ManualAsyncProcessInitialisation(ctx context.Context, cancel context.CancelFunc) AsyncProcess {
	return AsyncProcess{
		ctx:    ctx,
		cancel: cancel,
		done:   make(chan struct{}),
	}
}

func (p *AsyncProcess) ProcessBackground(ctx context.Context, _p interceptor.CanProcessBackground, s interceptor.State) interceptor.CanBeProcessedBackground {
	if p.CanBeProcessed == nil {
		fmt.Println("WARNING: AsyncProcess.CanBeProcessed is nil; this is not allowed")
		return nil
	}
	if p.done == nil { // ONLY POSSIBLE WHEN NOT ManualAsyncProcessInitialisation-ed
		p.done = make(chan struct{})
	}

	if p.ctx == nil { // ONLY POSSIBLE WHEN NOT ManualAsyncProcessInitialisation-ed
		if ctx == nil {
			fmt.Println("WARNING: ctx is nil; using background context")
			ctx = context.Background()
		}
		// NOTE: THERE IS A ASSUMPTION HERE; IF p.ctx IS NIL, THEN p.cancel IS ALSO NIL
		p.ctx, p.cancel = context.WithCancel(ctx)
	}

	go func() {
		err := p.Process(p.ctx, _p.(interceptor.CanProcess), s)
		p.mux.Lock()
		defer p.mux.Unlock()
		defer p.cancel()

		p.err = err
		p.done <- struct{}{}

		if err != nil {
			fmt.Println(p.err)
		}
	}()

	return p
}

func (p *AsyncProcess) Wait() error {
	<-p.done
	p.mux.RLock()
	defer p.mux.RUnlock()

	return p.err
}

func (p *AsyncProcess) Stop() {
	p.cancel()
}
