package socket

import (
	"context"
)

// TODO: LIMIT AND SELF KILL BUFFER.
// ELEMENT IS RESPONSIBLE TO KILL AND FREE ITS MEMORY. IT SHOULD ALSO DELETE ITSELF FROM THE PARENT
// HARD LIMIT ON NUMBER OF BUFFERS

type Buffered[T any] struct {
	element T
	ctx     context.Context
}

type Buffer[T any] interface {
	Pop(context.Context) (T, error)
	Push(context.Context, T) error
	Close()
}

type LimitKillBuffer[T any] struct {
	buffer chan Buffered[T]
}

func (b *LimitKillBuffer[T]) Pop(ctx context.Context) (T, error) {

}

func (b *LimitKillBuffer[T]) Push(ctx context.Context, element T) error {

}

func (b *LimitKillBuffer[T]) Close() {

}

func (b *LimitKillBuffer[T]) manager() {

}
