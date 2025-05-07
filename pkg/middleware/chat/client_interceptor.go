package chat

import (
	"context"
	"fmt"
	"time"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/process"
)

type ClientInterceptor struct {
	commonInterceptor
	Health interfaces.Processor
}

func (i *ClientInterceptor) BindSocketConnection(connection interceptor.Connection, writer interceptor.Writer, reader interceptor.Reader) (interceptor.Writer, interceptor.Reader, error) {
	return i.commonInterceptor.BindSocketConnection(connection, writer, reader)
}

func (i *ClientInterceptor) Init(connection interceptor.Connection) error {
	s, err := i.GetState(connection)
	if err != nil {
		return fmt.Errorf("error while init; err: %s", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p := process.NewIdentWaiter(ctx)
	if err := p.Process(nil, s); err != nil {
		return fmt.Errorf("error while init; err: %s", err.Error())
	}

	return nil
}
