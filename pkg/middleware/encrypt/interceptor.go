package encrypt

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/config"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/keyexchange"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/keyprovider"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/state"
)

type Interceptor struct {
	interceptor.NoOpInterceptor
	nonceValidator     NonceValidator
	keyExchangeManager keyexchange.Manager
	keyProvider        keyprovider.KeyProvider
	stateManager       state.Manager
	config             config.Config
}

func (i *Interceptor) BindSocketConnection(connection interceptor.Connection, writer interceptor.Writer, reader interceptor.Reader) (interceptor.Writer, interceptor.Reader, error) {
	ctx, cancel := context.WithCancel(i.Ctx)

	newState, err := state.NewState(ctx, cancel, i.config, connection, writer, reader)
	if err != nil {
		return nil, nil, err
	}

	if err := i.stateManager.SetState(connection, newState); err != nil {
		return nil, nil, err
	}

	return writer, reader, nil
}

func (i *Interceptor) Init(connection interceptor.Connection) error {
	_state, err := i.stateManager.GetState(connection)
	if err != nil {
		return err
	}

	if err := i.keyExchangeManager.Init(_state, keyexchange.WithKeySignature(i.keyProvider)); err != nil {
		return err
	}
}

func (i *Interceptor) InterceptSocketWriter(writer interceptor.Writer) interceptor.Writer {

}

func (i *Interceptor) InterceptSocketReader(reader interceptor.Reader) interceptor.Reader {

}

func (i *Interceptor) UnBindSocketConnection(connection interceptor.Connection) {

}

func (i *Interceptor) Close() error {

}
