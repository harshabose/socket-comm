package interceptor

import (
	"context"
	"io"

	"github.com/harshabose/socket-comm/pkg/message"
)

type Registry struct {
	factories []Factory
	messages  message.Registry
}

func NewRegistry() *Registry {
	return &Registry{
		factories: make([]Factory, 0),
	}
}

func (registry *Registry) Register(factory Factory) {
	registry.factories = append(registry.factories, factory)
}

func (registry *Registry) Build(ctx context.Context, id ClientID) (Interceptor, error) {
	if len(registry.factories) == 0 {
		return &NoOpInterceptor{}, nil
	}

	interceptors := make([]Interceptor, 0)
	for _, factory := range registry.factories {
		interceptor, err := factory.NewInterceptor(ctx, id, registry.messages)
		if err != nil {
			return nil, err
		}

		interceptors = append(interceptors, interceptor)
	}

	return CreateChain(interceptors), nil
}

type Factory interface {
	NewInterceptor(context.Context, ClientID, message.Registry) (Interceptor, error)
}

type Connection interface {
	Write(ctx context.Context, p []byte) error
	Read(ctx context.Context) ([]byte, error)
	io.Closer
}

type Interceptor interface {
	ID() ClientID

	Ctx() context.Context

	GetMessageRegistry() message.Registry

	BindSocketConnection(Connection, Writer, Reader) (Writer, Reader, error)

	Init(Connection) error

	InterceptSocketWriter(Writer) Writer

	InterceptSocketReader(Reader) Reader

	UnBindSocketConnection(Connection)

	io.Closer
}

type Writer interface {
	Write(ctx context.Context, conn Connection, message message.Message) error
}

type Reader interface {
	Read(ctx context.Context, conn Connection) (message.Message, error)
}

type ReaderFunc func(ctx context.Context, connection Connection) (message.Message, error)

func (f ReaderFunc) Read(ctx context.Context, connection Connection) (message.Message, error) {
	return f(ctx, connection)
}

type WriterFunc func(ctx context.Context, connection Connection, message message.Message) error

func (f WriterFunc) Write(ctx context.Context, connection Connection, message message.Message) error {
	return f(ctx, connection, message)
}

type NoOpInterceptor struct {
	iD              ClientID
	messageRegistry message.Registry
	ctx             context.Context
}

func NewNoOpInterceptor(ctx context.Context, id ClientID, registry message.Registry) NoOpInterceptor {
	return NoOpInterceptor{
		ctx:             ctx,
		iD:              id,
		messageRegistry: registry,
	}
}

func (interceptor *NoOpInterceptor) Ctx() context.Context {
	return interceptor.ctx
}

func (interceptor *NoOpInterceptor) ID() ClientID {
	return interceptor.iD
}

func (interceptor *NoOpInterceptor) GetMessageRegistry() message.Registry {
	return interceptor.messageRegistry
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
