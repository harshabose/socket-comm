package messages

import (
	"fmt"
	"math"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"

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

func (m *RequestHealth) ReadProcess(_i interceptor.Interceptor, connection interceptor.Connection) error {
	s, ok := _i.(interfaces.CanGetState)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	ss, err := s.GetState(connection)
	if err != nil {
		return fmt.Errorf("error while read processing 'RequestHealth' msg; err: %s", err.Error())
	}

	msg, err := NewHealthResponse(m)
	if err != nil {
		return fmt.Errorf("error while read processing 'RequestHealth' msg; err: %s", err.Error())
	}

	w, ok := ss.(interfaces.CanWriteMessage)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	if err := w.Write(msg); err != nil {
		return fmt.Errorf("error while read processing 'RequestHealth' msg; err: %s", err.Error())
	}

	return nil
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

// NOTE: BASIC HEALTH RESPONSE FOR ROOM MANAGEMENT, OTHER METRICS WILL BE DEALT WITH LATER

type HealthResponse struct {
	interceptor.BaseMessage
	RequestTimeStamp int64                  `json:"-"` // in nanoseconds
	RoomID           types.RoomID           `json:"room_id"`
	ConnectionStatus types.ConnectionState  `json:"connection_status"`
	ConnectionUptime types.ConnectionUptime `json:"connection_uptime"`
	CPUUsage         types.CPUUsage         `json:"cpu_usage"`
	MemoryUsage      types.MemoryUsage      `json:"memory_usage"`
	NetworkUsage     types.NetworkUsage     `json:"network_usage"`
	Latency          types.LatencyMs        `json:"latency"`
}

func NewHealthResponse(request *RequestHealth) (*HealthResponse, error) {
	response := &HealthResponse{}

	response.RequestTimeStamp = request.Timestamp
	response.RoomID = request.RoomID

	response.setConnectionStatus(request.ConnectionStartTime)

	if err := response.setCPUUsage(); err != nil {
		return nil, err
	}

	if err := response.setMemoryUsage(); err != nil {
		return nil, err
	}

	if err := response.setLatency(); err != nil {
		return nil, err
	}

	bmsg, err := interceptor.NewBaseMessage(message.NoneProtocol, nil, response)
	if err != nil {
		return nil, err
	}
	response.BaseMessage = bmsg

	return response, nil
}

func (m *HealthResponse) GetProtocol() message.Protocol {
	return HealthResponseProtocol
}

func (m *HealthResponse) ReadProcess(_i interceptor.Interceptor, connection interceptor.Connection) error {
	// TODO: IMPLEMENT HEALTH MANAGEMENT

	return nil
}

func (m *HealthResponse) WriteProcess(_i interceptor.Interceptor, connection interceptor.Connection) error {
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

func (m *HealthResponse) setConnectionStatus(startTime int64) {
	m.ConnectionStatus = types.ConnectionStateUp
	m.ConnectionUptime = m.getConnectionUptime(startTime)
}

func (m *HealthResponse) getConnectionUptime(startTime int64) types.ConnectionUptime {
	return types.ConnectionUptime(time.Now().Unix() - startTime)
}

func (m *HealthResponse) setCPUUsage() error {
	perCorePercentages, err := cpu.Percent(100*time.Millisecond, false)
	if err != nil {
		return fmt.Errorf("error while reading CPU usage; err: %s", err.Error())
	}

	numCores := len(perCorePercentages)
	if numCores == 0 {
		return fmt.Errorf("error while reading CPU usage; err: 'no cores detected'")
	}

	if numCores > math.MaxUint8 {
		return fmt.Errorf("error while reading CPU usage; err: 'too many cores detected'")
	}

	m.CPUUsage = types.CPUUsage{
		NumCores: uint8(numCores),
		Percent:  perCorePercentages,
	}

	return nil
}

func (m *HealthResponse) setMemoryUsage() error {
	vmStats, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Errorf("error while reading memory usage; err: %s", err.Error())
	}

	total := float32(vmStats.Total) / (1024.0 * 1024.0 * 1024.0)
	used := float32(vmStats.Used) / (1024.0 * 1024.0 * 1024.0)
	available := float32(vmStats.Available) / (1024.0 * 1024.0 * 1024.0)

	m.MemoryUsage = types.MemoryUsage{
		Total:          total,
		Used:           used,
		UsedRatio:      used / total,
		Available:      available,
		AvailableRatio: available / total,
	}

	return nil
}

func (m *HealthResponse) setLatency() error {
	latencyNano := time.Now().UnixNano() - m.RequestTimeStamp

	latencyMs := latencyNano / int64(time.Millisecond)

	if latencyMs < 0 {
		return fmt.Errorf("error while reading latency; err: 'latency is negative'")
	}

	m.Latency = types.LatencyMs(latencyMs)

	return nil
}
