package socket

import (
	"context"
	"errors"
	"fmt"
	"net"
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

// Metrics holds server statistics
type Metrics struct {
	ActiveConnections int
	TotalConnections  int
	FailedConnections int
	mux               sync.RWMutex
}

type Socket struct {
	ID              types.SocketID `json:"id"`
	server          *http.Server
	router          *http.ServeMux
	settings        Settings
	interceptor     interceptor.Interceptor
	connections     map[string]interceptor.Connection
	metrics         *Metrics
	messageRegistry message.Registry
	cancel          context.CancelFunc
	ctx             context.Context
	mux             sync.Mutex
}

func NewSocket(ctx context.Context, settings Settings, registry message.Registry) *Socket {
	ctx2, cancel := context.WithCancel(ctx)
	return &Socket{
		ID:              types.SocketID(uuid.NewString()),
		settings:        settings,
		messageRegistry: registry,
		connections:     make(map[string]interceptor.Connection),
		metrics:         &Metrics{},
		cancel:          cancel,
		ctx:             ctx2,
	}
}

func (s *Socket) GetID() types.SocketID {
	return s.ID
}

func (s *Socket) Ctx(l net.Listener) context.Context {
	return s.ctx
}

func (s *Socket) Init() error {
	if err := s.settings.Validate(); err != nil {
		return err
	}

	tlsConfig, err := GetTLSV1(s.settings.TLSCertFile, s.settings.TLSKeyFile)
	if err != nil {
		return err
	}

	s.router = http.NewServeMux()
	s.server = &http.Server{
		Addr:              fmt.Sprintf("%s:%d", s.settings.Address, s.settings.Port),
		ReadTimeout:       s.settings.PopMessageTimeout,
		WriteTimeout:      s.settings.PushMessageTimout,
		ReadHeaderTimeout: s.settings.ReadHeaderTimeout,
		IdleTimeout:       s.settings.IdleTimout,
		TLSConfig:         tlsConfig,
		Handler:           s.router,
		BaseContext:       s.Ctx,
		// TODO: MAYBE ADD MORE
	}

	// Set up routes
	s.router.HandleFunc("/ws", s.handleWebSocket)
	s.router.HandleFunc("/health", s.handleHealth)
	s.router.HandleFunc("/metrics", s.handleMetrics)
	return nil
}

func (s *Socket) Serve() error {
	for {
		select {
		case <-s.ctx.Done():
			return nil // TODO: add error
		default:
			if s.server.TLSConfig != nil {
				if err := s.server.ListenAndServe(); err != nil {
					if err != nil && !errors.Is(err, http.ErrServerClosed) {
						return err
					}
					fmt.Println("error while serving; err: ", err.Error())
					fmt.Println("retrying...")
					time.Sleep(1 * time.Second)
				}
				continue
			}
			if err := s.server.ListenAndServeTLS(s.settings.TLSCertFile, s.settings.TLSKeyFile); err != nil {
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					return err
				}
				fmt.Println("error while serving; err: ", err.Error())
				fmt.Println("retrying...")
				time.Sleep(1 * time.Second)
			}
		}
	}
}

func (s *Socket) Read(ctx context.Context, connection interceptor.Connection) (message.Message, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.settings.PopMessageTimeout)
	defer cancel()

	data, err := connection.Read(ctx)
	if err != nil {
		return nil, err
	}

	return s.messageRegistry.UnmarshalRaw(data)
}

func (s *Socket) Write(ctx context.Context, connection interceptor.Connection, msg message.Message) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.settings.PushMessageTimout)
	defer cancel()

	data, err := msg.Marshal()
	if err != nil {
		return err
	}

	return connection.Write(ctx, data)
}

// registerConnection adds a new connection
func (s *Socket) registerConnection(id string, conn interceptor.Connection) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.connections[id] = conn

	s.metrics.mux.Lock()
	s.metrics.ActiveConnections++
	s.metrics.TotalConnections++
	s.metrics.mux.Unlock()
}

// unregisterConnection removes a connection
func (s *Socket) unregisterConnection(id string) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if _, exists := s.connections[id]; exists {
		delete(s.connections, id)

		s.metrics.mux.Lock()
		s.metrics.ActiveConnections--
		s.metrics.mux.Unlock()
	}
}

// closeAllConnections closes all active connections
func (s *Socket) closeAllConnections() {
	s.mux.Lock()
	defer s.mux.Unlock()

	for _, conn := range s.connections {
		if err := conn.Close(); err != nil {
			fmt.Println("error while closing a connection; err:", err.Error())
		}
	}
}

func (s *Socket) handleWebSocket(writer http.ResponseWriter, request *http.Request) {
	s.metrics.mux.RLock()
	currentConnections := s.metrics.ActiveConnections
	s.metrics.mux.RUnlock()

	if currentConnections >= s.settings.MaxConnections {
		fmt.Println("Connection limit reached", "active", currentConnections)
		http.Error(writer, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}

	conn, err := websocket.Accept(writer, request, nil)
	if err != nil {
		s.metrics.mux.Lock()
		s.metrics.FailedConnections++
		s.metrics.mux.Unlock()

		fmt.Println("error while handling client; removing client...")
		return
	}

	iD := uuid.NewString()
	connection := newAdaptor(request.Context(), iD, conn, s.settings.PopMessageTimeout, s.settings.PushMessageTimout)

	s.registerConnection(iD, connection)
	defer s.unregisterConnection(iD)

	connection.StartReaderWriter()

	if _, _, err := s.interceptor.BindSocketConnection(connection, s, s); err != nil {
		fmt.Println(fmt.Errorf("error while binding socket to interceptors; err: %s", err.Error()))
		fmt.Println("dropping client...")
		return
	}
	defer s.interceptor.UnBindSocketConnection(connection)

	if err := s.interceptor.Init(connection); err != nil {
		fmt.Println("error while connection init; dropping client")
		fmt.Println("dropping client...")
		return
	}

	connection.WaitUntilClose()
}

func (s *Socket) ShutDown(ctx context.Context) error {
	ctx2, cancel := context.WithTimeout(ctx, s.settings.ShutdownTimout)
	defer cancel()

	s.cancel()
	if err := s.server.Shutdown(ctx2); err != nil {
		return fmt.Errorf("server shutdown error: %w", err)
	}

	s.closeAllConnections()
	return nil
}

// handleHealth provides a health check endpoint
func (s *Socket) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

// handleMetrics exposes server metrics
func (s *Socket) handleMetrics(w http.ResponseWriter, _ *http.Request) {
	s.metrics.mux.RLock()
	defer s.metrics.mux.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	_, _ = fmt.Fprintf(w, `{
        "active_connections": %d,
        "total_connections": %d,
        "failed_connections": %d
    }`, s.metrics.ActiveConnections, s.metrics.TotalConnections,
		s.metrics.FailedConnections)
}
