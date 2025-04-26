package interceptor

import "github.com/harshabose/socket-comm/internal/util"

type Chain struct {
	interceptors []Interceptor
}

func CreateChain(interceptors []Interceptor) *Chain {
	return &Chain{interceptors: interceptors}
}

func (chain *Chain) BindSocketConnection(connection Connection, writer Writer, reader Reader) (Writer, Reader, error) {
	var (
		w   Writer
		r   Reader
		err error
	)

	for _, interceptor := range chain.interceptors {
		if w, r, err = interceptor.BindSocketConnection(connection, chain.InterceptSocketWriter(writer), chain.InterceptSocketReader(reader)); err != nil {
			return nil, nil, err
		}
	}

	return w, r, nil
}

func (chain *Chain) Init(connection Connection) error {
	for _, interceptor := range chain.interceptors {
		if err := interceptor.Init(connection); err != nil {
			return err
		}
	}

	return nil
}

func (chain *Chain) InterceptSocketWriter(writer Writer) Writer {
	for _, interceptor := range chain.interceptors {
		writer = interceptor.InterceptSocketWriter(writer)
	}

	return writer
}

func (chain *Chain) InterceptSocketReader(reader Reader) Reader {
	for _, interceptor := range chain.interceptors {
		reader = interceptor.InterceptSocketReader(reader)
	}

	return reader
}

func (chain *Chain) UnBindSocketConnection(connection Connection) {
	for _, interceptor := range chain.interceptors {
		interceptor.UnBindSocketConnection(connection)
	}
}

func (chain *Chain) Close() error {
	var merr util.MultiError

	for _, interceptor := range chain.interceptors {
		merr.Add(interceptor.Close())
	}

	return merr.ErrorOrNil()
}
