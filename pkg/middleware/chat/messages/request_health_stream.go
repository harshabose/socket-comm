package messages

import (
	"fmt"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

var (
	RequestHealthProtocol  message.Protocol = "room:request_health"
	HealthResponseProtocol message.Protocol = "room:health_response"
)

type RequestHealth struct {
	interceptor.BaseMessage
	RoomID types.RoomID `json:"room_id"`
}

func NewRequestHealth(id types.RoomID) *RequestHealth {
	return &RequestHealth{
		RoomID: id,
	}
}

func NewRequestHealthFactory(id types.RoomID) func() message.Message {
	return func() message.Message {
		return NewRequestHealth(id)
	}
}

func (m *RequestHealth) GetProtocol() message.Protocol {
	return RequestHealthProtocol
}

func (m *RequestHealth) ReadProcess(_i interceptor.Interceptor, connection interceptor.Connection) error {

}

func (m *RequestHealth) WriteProcess(_i interceptor.Interceptor, connection interceptor.Connection) error {
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

type HealthResponse struct {
	interceptor.BaseMessage
	RoomID types.RoomID `json:"room_id"`
}
