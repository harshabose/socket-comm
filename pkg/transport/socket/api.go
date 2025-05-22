package socket

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
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

func (a *API) NewSocket(ctx context.Context, options ...Option) (*Socket, error) {
	s := NewSocket(NewDefaultSettings(), a.messagesRegistry)

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
