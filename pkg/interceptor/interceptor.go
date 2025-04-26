package interceptor

import (
	"context"
	"io"
	"sync"
)

type Registry struct {
	factories []Factory
}

func (registry *Registry) Register(factory Factory) {
	registry.factories = append(registry.factories, factory)
}

func (registry *Registry) Build(ctx context.Context, id string) (Interceptor, error) {
	if len(registry.factories) == 0 {
		return &NoOpInterceptor{}, nil
	}

	interceptors := make([]Interceptor, 0)
	for _, factory := range registry.factories {
		interceptor, err := factory.NewInterceptor(ctx, id)
		if err != nil {
			return nil, err
		}

		interceptors = append(interceptors, interceptor)
	}

	return CreateChain(interceptors), nil
}

type Factory interface {
	NewInterceptor(context.Context, string) (Interceptor, error)
}

type Connection interface {
	Write(ctx context.Context, p []byte) error
	Read(ctx context.Context) ([]byte, error)
}

type Interceptor interface {
	BindSocketConnection(Connection, Writer, Reader) (Writer, Reader, error)

	Init(Connection) error

	InterceptSocketWriter(Writer) Writer

	InterceptSocketReader(Reader) Reader

	UnBindSocketConnection(Connection)

	io.Closer
}

type Writer interface {
	Write(conn Connection, message Message) error
}

type Reader interface {
	Read(conn Connection) (Message, error)
}

type ReaderFunc func(conn Connection) (Message, error)

func (f ReaderFunc) Read(conn Connection) (Message, error) {
	return f(conn)
}

type WriterFunc func(conn Connection, message Message) error

func (f WriterFunc) Write(conn Connection, message Message) error {
	return f(conn, message)
}

type NoOpInterceptor struct {
	ID    string
	Mutex sync.RWMutex
	Ctx   context.Context
}

func (interceptor *NoOpInterceptor) BindSocketConnection(_ Connection, _ Writer, _ Reader) (Writer, Reader, error) {
	return nil, nil, nil
}

func (interceptor *NoOpInterceptor) Init(_ Connection) error {
	return nil
}

func (interceptor *NoOpInterceptor) InterceptSocketWriter(writer Writer) Writer {
	return writer
}

func (interceptor *NoOpInterceptor) InterceptSocketReader(reader Reader) Reader {
	return reader
}

func (interceptor *NoOpInterceptor) UnBindSocketConnection(_ Connection) {}

func (interceptor *NoOpInterceptor) Close() error {
	return nil
}
