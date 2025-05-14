package socket

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/transport/types"
)

type API struct {
	interceptorRegistry *interceptor.Registry
	messagesRegistry    message.Registry
}

type APIOption = func(*API) error

func NewAPI(options ...APIOption) (*API, error) {
	api := &API{}

	for _, option := range options {
		if err := option(api); err != nil {
			return nil, err
		}
	}

	if api.interceptorRegistry == nil {
		api.interceptorRegistry = interceptor.NewRegistry()
		if err := WithDefaultInterceptorRegistry(api.interceptorRegistry); err != nil {
			return nil, err
		}
	}

	if api.messagesRegistry == nil {
		api.messagesRegistry = message.NewDefaultRegistry()
		if err := WithDefaultMessageRegistry(api.messagesRegistry); err != nil {
			return nil, err
		}
	}

	return api, nil
}

// TODO: MAKE REGISTRIES TO NON POINTERS

func (a *API) NewSocket(ctx context.Context, id types.SocketID, options ...Option) (*Socket, error) {
	s := &Socket{
		ID:              id,
		settings:        NewDefaultSettings(),
		messageRegistry: a.messagesRegistry,
		ctx:             ctx,
	}

	interceptors, err := a.interceptorRegistry.Build(s.ctx, string(s.ID))
	if err != nil {
		return nil, err
	}

	s.interceptor = interceptors

	for _, option := range options {
		if err := option(s); err != nil {
			return nil, err
		}
	}

	if err := s.Init(); err != nil {
		return nil, err
	}

	return s, nil
}
