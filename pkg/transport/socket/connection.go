package socket

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/coder/websocket"
)

var (
	ErrNotSupportedMessageType = errors.New("not supported message type")
	ErrConnectionClosed        = errors.New("connection closed")
)

type adaptor struct {
	connectionSettings
	id         string
	conn       *websocket.Conn
	readQ      Buffer[[]byte]
	writeQ     Buffer[[]byte]
	ctx        context.Context
	cancel     context.CancelFunc
	closeOnce  sync.Once
	wg         sync.WaitGroup
	closeErr   error
	closeErrMu sync.Mutex
}

func newAdaptor(ctx context.Context, id string, conn *websocket.Conn, readTimeout time.Duration, writeTimeout time.Duration) *adaptor {
	// Create a child context with cancellation
	childCtx, cancel := context.WithCancel(ctx)

	return &adaptor{
		connectionSettings: connectionSettings{
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},
		ctx:    childCtx,
		cancel: cancel,
		id:     id,
		conn:   conn,
	}
}

// Write pushes the message of the type '[]byte' to the WriteQ, which will be later sent through the socket
func (a *adaptor) Write(ctx context.Context, p []byte) error {
	if ctx == nil {
		ctx = context.Background()
	}

	// Create a copy of the message to prevent data races
	msgCopy := make([]byte, len(p))
	copy(msgCopy, p)

	return a.writeQ.Push(ctx, msgCopy)
}

// Read reads a message of the type '[]byte' from the ReadQ, which was read from the websocket.
func (a *adaptor) Read(ctx context.Context) ([]byte, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	data, err := a.readQ.Pop(ctx)
	if err != nil {
		return nil, err
	}

	// Return a copy to prevent data races
	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)

	return dataCopy, nil
}

func (a *adaptor) StartReaderWriter() {
	a.wg.Add(2) // One for the reader, one for the writer
	go a.Writer()
	go a.Reader()
}

func (a *adaptor) Writer() {
	defer a.wg.Done()
	defer a.closeWithError(ErrConnectionClosed)

	for {
		select {
		case <-a.ctx.Done():
			return
		default:
			// Use a timeout context for the write operation
			writeCtx, cancel := context.WithTimeout(a.ctx, a.WriteTimeout)
			p, err := a.writeQ.Pop(writeCtx)
			cancel()

			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
					// Check if context was cancelled due to connection close
					select {
					case <-a.ctx.Done():
						return
					default:
						// Just a timeout, continue
						continue
					}
				}

				fmt.Printf("Error while popping message from WriteQ; err: %s\n", err.Error())
				continue
			}

			// Use a timeout context for the websocket write
			writeCtx, cancel = context.WithTimeout(a.ctx, a.WriteTimeout)
			err = a.conn.Write(writeCtx, websocket.MessageText, p)
			cancel()

			if err != nil {
				fmt.Printf("Error while writing message to socket; err: %s\n", err.Error())
				return
			}
		}
	}
}

func (a *adaptor) Reader() {
	defer a.wg.Done()
	defer a.closeWithError(ErrConnectionClosed)

	for {
		select {
		case <-a.ctx.Done():
			return
		default:
			// Use a timeout context for the read operation
			readCtx, cancel := context.WithTimeout(a.ctx, a.ReadTimeout)
			msgType, p, err := a.conn.Read(readCtx)
			cancel()

			if err != nil {
				fmt.Printf("Error while reading message from socket; err: %s\n", err.Error())
				return
			}

			if msgType != websocket.MessageText {
				fmt.Printf("Error while reading message from socket; err: %s\n", ErrNotSupportedMessageType.Error())
				continue
			}

			// Use a background context since we want to buffer the message even if the operation takes time
			err = a.readQ.Push(a.ctx, p)
			if err != nil {
				fmt.Printf("Error while pushing message to ReadQ; err: %s\n", err.Error())
				continue
			}
		}
	}
}

// Close initiates a graceful shutdown of the connection
func (a *adaptor) Close() error {
	var err error
	a.closeOnce.Do(func() {
		// Cancel the context to signal all goroutines to stop
		a.cancel()

		// Try to send a close message to the peer
		closeErr := a.conn.Close(websocket.StatusNormalClosure, "connection closed")
		if closeErr != nil {
			err = closeErr
		}

		// Close the buffers
		a.readQ.Close()
		a.writeQ.Close()
	})
	return err
}

// closeWithError stores the error that caused the connection to close
func (a *adaptor) closeWithError(err error) {
	a.closeErrMu.Lock()
	if a.closeErr == nil {
		a.closeErr = err
	}
	a.closeErrMu.Unlock()
	a.Close()
}

// GetCloseError returns the error that caused the connection to close
func (a *adaptor) GetCloseError() error {
	a.closeErrMu.Lock()
	defer a.closeErrMu.Unlock()
	return a.closeErr
}

// WaitUntilClose blocks until all reader and writer goroutines exit
func (a *adaptor) WaitUntilClose() {
	a.wg.Wait()
}
