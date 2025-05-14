package socket

import (
	"context"
	"errors"

	"github.com/coder/websocket"
)

var (
	ErrNotSupportedMessageType = errors.New("not supported message type")
	ErrConnectionClosed        = errors.New("connection closed")
	ErrServerClosed            = errors.New("server closed")
	ErrReaderWriterAlreadySet  = errors.New("reader and writer already set")
)

type adaptor struct {
	id   string
	conn *websocket.Conn
}

func newAdaptor(id string, conn *websocket.Conn) *adaptor {
	return &adaptor{
		id:   id,
		conn: conn,
	}
}

func (a *adaptor) Write(ctx context.Context, p []byte) error {
	return a.conn.Write(ctx, websocket.MessageText, p)
}

func (a *adaptor) Read(ctx context.Context) ([]byte, error) {
	msgType, p, err := a.conn.Read(ctx)
	if err != nil {
		return nil, err
	}

	if msgType != websocket.MessageText {
		return nil, ErrNotSupportedMessageType
	}

	return p, nil
}
