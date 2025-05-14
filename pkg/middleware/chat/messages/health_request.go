package messages

import (
	"context"
	"fmt"
	"time"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/process"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

var RequestHealthProtocol message.Protocol = "room:request_health"

// SendHealthStats is sent by the server to request the health from a client.
// The client then responds with UpdateHealthStat with the health.Stat embedded.
type SendHealthStats struct {
	interceptor.BaseMessage
	RoomID                   types.RoomID `json:"room_id"`
	ConnectionStartTimestamp time.Time    `json:"connection_start_timestamp"`
	Timestamp                time.Time    `json:"timestamp"` // in nanoseconds
}

func NewSendHealthStats(id types.RoomID) (*SendHealthStats, error) {
	msg := &SendHealthStats{
		RoomID:    id,
		Timestamp: time.Now(),
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
		return NewSendHealthStats(id)
	}
}

func (m *SendHealthStats) GetProtocol() message.Protocol {
	return RequestHealthProtocol
}

func (m *SendHealthStats) ReadProcess(ctx context.Context, _i interceptor.Interceptor, connection interceptor.Connection) error {
	s, ok := _i.(interfaces.CanGetState)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	ss, err := s.GetState(connection)
	if err != nil {
		return fmt.Errorf("error while read processing 'SendHealthStats' msg; err: %s", err.Error())
	}

	return process.NewSendMessage(NewUpdateHealthStatFactory(m)).Process(ctx, nil, ss)
}

func (m *SendHealthStats) WriteProcess(ctx context.Context, _i interceptor.Interceptor, connection interceptor.Connection) error {
	s, ok := _i.(interfaces.CanGetState)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	ss, err := s.GetState(connection)
	if err != nil {
		return fmt.Errorf("error while read processing 'SendHealthStats' msg; err: %s", err.Error())
	}

	id, err := ss.GetClientID()
	if err != nil {
		return fmt.Errorf("error while read processing 'SendHealthStats' msg; err: %s", err.Error())
	}

	m.SetSender(message.Sender(_i.ID()))
	m.SetReceiver(message.Receiver(id))

	return nil
}
