package messages

import (
	"context"
	"fmt"
	"time"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

var RequestHealthProtocol message.Protocol = "room:request_health"

type RequestHealth struct {
	interceptor.BaseMessage
	RoomID              types.RoomID `json:"room_id"`
	Timestamp           int64        `json:"timestamp"` // in nanoseconds
	ConnectionStartTime int64        `json:"connection_start_time"`
}

func NewRequestHealth(id types.RoomID) (*RequestHealth, error) {
	msg := &RequestHealth{
		RoomID:    id,
		Timestamp: time.Now().UnixNano(),
	}

	bmsg, err := interceptor.NewBaseMessage(message.NoneProtocol, nil, msg)
	if err != nil {
		panic(err)
	}
	msg.BaseMessage = bmsg

	return msg, nil
}

func NewRequestHealthFactory(id types.RoomID) func() (message.Message, error) {
	return func() (message.Message, error) {
		return NewRequestHealth(id)
	}
}

func (m *RequestHealth) GetProtocol() message.Protocol {
	return RequestHealthProtocol
}

func (m *RequestHealth) ReadProcess(ctx context.Context, _i interceptor.Interceptor, connection interceptor.Connection) error {
	s, ok := _i.(interfaces.CanGetState)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	ss, err := s.GetState(connection)
	if err != nil {
		return fmt.Errorf("error while read processing 'RequestHealth' msg; err: %s", err.Error())
	}

	msg, err := NewHealthResponse(m, 5*time.Second)
	if err != nil {
		return fmt.Errorf("error while read processing 'RequestHealth' msg; err: %s", err.Error())
	}

	if err := ss.Write(msg); err != nil {
		return fmt.Errorf("error while read processing 'RequestHealth' msg; err: %s", err.Error())
	}

	return nil
}

func (m *RequestHealth) WriteProcess(ctx context.Context, _i interceptor.Interceptor, connection interceptor.Connection) error {
	s, ok := _i.(interfaces.CanGetState)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	ss, err := s.GetState(connection)
	if err != nil {
		return fmt.Errorf("error while read processing 'RequestHealth' msg; err: %s", err.Error())
	}

	id, err := ss.GetClientID()
	if err != nil {
		return fmt.Errorf("error while read processing 'RequestHealth' msg; err: %s", err.Error())
	}

	m.SetSender(message.Sender(_i.ID()))
	m.SetReceiver(message.Receiver(id))

	return nil
}
