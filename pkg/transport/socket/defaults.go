package socket

import (
	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
)

func WithDefaultInterceptorRegistry(registry *interceptor.Registry) error {
	registry.Register(nil)
	// TODO: IMPLEMENT THIS
	return nil
}

func WithDefaultMessageRegistry(registry message.Registry) error {
	// TODO: IMPLEMENT THIS
	return nil
}

func WithTLSConfig(TLSCertPath string, TLSKeyPath string) Option {
	return func(s *Socket) error {
		s.settings.TLSCertFile = TLSCertPath
		s.settings.TLSKeyFile = TLSKeyPath

		return nil
	}
}

func WithInterceptorRegistry(registry *interceptor.Registry) APIOption {
	return func(api *API) error {
		api.interceptorRegistry = registry
		return nil
	}
}

func WithMessageRegistry(registry message.Registry) APIOption {
	return func(api *API) error {
		api.messagesRegistry = registry
		return nil
	}
}
