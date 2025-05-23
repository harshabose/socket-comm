package chat

import (
	"context"
	"fmt"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/processors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
)

type Option = func(i interceptor.Interceptor) error

type InterceptorFactory struct {
	options []Option
}

func NewInterceptorFactory(options ...Option) *InterceptorFactory {
	return &InterceptorFactory{
		options: options,
	}
}

func WithServerInterceptor(i interceptor.Interceptor) error {
	c, ok := i.(*commonInterceptor)
	if !ok {
		return fmt.Errorf("can only convert common chat interceptor to client chat interceptor; err: %s", interceptor.ErrInterfaceMisMatch.Error())
	}

	i = &ServerInterceptor{
		commonInterceptor: c,
		Rooms:             processors.NewRoomProcessor(c.Ctx()),
		Health:            processors.NewHealthProcessor(c.Ctx()),
	}

	// TODO: add server messages to the registry

	return nil
}

func WithClientInterceptor(i interceptor.Interceptor) error {
	c, ok := i.(*commonInterceptor)
	if !ok {
		return fmt.Errorf("can only convert common chat interceptor to client chat interceptor; err: %s", interceptor.ErrInterfaceMisMatch.Error())
	}

	i = &ClientInterceptor{
		commonInterceptor: c,
		Health:            processors.NewHealthProcessor(c.Ctx()),
	}

	// TODO: add client messages to the registry

	return nil
}

func (f *InterceptorFactory) NewInterceptor(ctx context.Context, id interceptor.ClientID, registry message.Registry) (interceptor.Interceptor, error) {
	i := &commonInterceptor{
		NoOpInterceptor:      interceptor.NewNoOpInterceptor(ctx, id, registry),
		readProcessMessages:  message.NewDefaultRegistry(),
		writeProcessMessages: message.NewDefaultRegistry(),
		states:               state.NewManager(),
	}

	// TODO: add common messages to the registry

	for _, option := range f.options {
		if err := option(i); err != nil {
			return nil, err
		}
	}

	// TODO: copy readProcessMessages and writeProcessMessages to registry

	return i, nil
}
