package chat

import (
	"github.com/harshabose/socket-comm/pkg/interceptor"
)

type ClientInterceptor struct {
	commonInterceptor
}

func (i *ClientInterceptor) BindSocketConnection(connection interceptor.Connection, writer interceptor.Writer, reader interceptor.Reader) (interceptor.Writer, interceptor.Reader, error) {
	return i.commonInterceptor.BindSocketConnection(connection, writer, reader)
}
