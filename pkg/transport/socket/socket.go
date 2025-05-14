package socket

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/transport/types"
)

const DEBUG = true

type Option func(*Socket) error

type Socket struct {
	ID              types.SocketID `json:"id"`
	server          *http.Server
	router          *http.ServeMux
	handlerFunc     http.HandlerFunc
	settings        *Settings
	interceptor     interceptor.Interceptor
	messageRegistry message.Registry
	mux             sync.RWMutex
	ctx             context.Context
}

func (s *Socket) GetID() types.SocketID {
	return s.ID
}

func (s *Socket) Init() error {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.router = http.NewServeMux()
	s.server = &http.Server{
		ReadTimeout:  s.settings.ReadTimout,
		WriteTimeout: s.settings.WriteTimout,
		IdleTimeout:  s.settings.IdleTimout,
		// TODO: MAYBE ADD MORE
	}
	s.handlerFunc = s.handler

	return nil
}

func (s *Socket) Serve() error {
	defer s.Close()

	for {
		select {
		case <-s.ctx.Done():
			return nil // TODO: add error

		default:
			if DEBUG {
				if err := s.server.ListenAndServe(); err != nil {
					fmt.Println("error while serving; err: ", err.Error())
					fmt.Println("retrying...")
					time.Sleep(1 * time.Second)
				}
			}
			if err := s.server.ListenAndServeTLS(s.settings.TLSCertFile, s.settings.TLSKeyFile); err != nil {
				fmt.Println("error while serving; err: ", err.Error())
				fmt.Println("retrying...")
				time.Sleep(1 * time.Second)
			}
		}
	}
}

func (s *Socket) Close() {

}

func (s *Socket) Read(ctx context.Context, connection interceptor.Connection) (message.Message, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.settings.ReadTimout)
	defer cancel()

	data, err := connection.Read(ctx)
	if err != nil {
		return nil, err
	}

	return s.messageRegistry.UnmarshalRaw(data)
}

func (s *Socket) Write(ctx context.Context, connection interceptor.Connection, msg message.Message) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.settings.WriteTimout)
	defer cancel()

	data, err := msg.Marshal()
	if err != nil {
		return err
	}

	return connection.Write(ctx, data)
}

func (s *Socket) handler(writer http.ResponseWriter, request *http.Request) {
	conn, err := websocket.Accept(writer, request, nil)
	if err != nil {
		fmt.Println("error while handling client; removing client...")
	}
	connection := newAdaptor(uuid.NewString(), conn)

	w, r, err := s.interceptor.BindSocketConnection(connection, s, s)
	if err != nil {
		return
	}

	for {
		select {
		case <-request.Context().Done():
			return
		case <-s.ctx.Done():
			return
		default:
			msg, err := r.Read(s.ctx, connection)
			if err != nil {
				return
			}

			fmt.Println(msg)
		}
	}

	defer s.interceptor.UnBindSocketConnection(connection)

	if err := s.interceptor.Init(connection); err != nil {
		fmt.Println("error while connection init; dropping client")
		return
	}
}
