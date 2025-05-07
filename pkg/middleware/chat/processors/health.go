package processors

import (
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
)

type Health struct {
}

func (p *Health) Process(process interfaces.CanBeProcessed, state *state.State) error {
	return process.Process(p, state)
}

func (p *Health) ProcessBackground(process interfaces.CanBeProcessedBackground, state *state.State) interfaces.CanBeProcessedBackground {
	return process.ProcessBackground(p, state)
}

func (p *Health) GetSnapshot() {

}
