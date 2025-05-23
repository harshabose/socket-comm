package messages

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/process"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

const HealthResponseProtocol message.Protocol = "room:health_response"

// UpdateHealthStat is sent by a client to the server in response to SendHealthStats.
type UpdateHealthStat struct {
	interceptor.BaseMessage
	process.UpdateHealthStat
}

func NewUpdateHealthStatFactory(request *SendHealthStats) func() (message.Message, error) {
	return func() (message.Message, error) {
		return NewUpdateHealthStat(request)
	}
}

func NewUpdateHealthStat(request *SendHealthStats) (*UpdateHealthStat, error) {
	response := &UpdateHealthStat{}

	response.RoomID = request.RoomID

	response.setConnectionStatus(request.ConnectionStartTimestamp)

	if err := response.setCPUUsage(); err != nil {
		return nil, err
	}

	if err := response.setMemoryUsage(); err != nil {
		return nil, err
	}

	if err := response.setLatency(request.Timestamp); err != nil {
		return nil, err
	}

	bmsg, err := interceptor.NewBaseMessage(message.NoneProtocol, nil, response)
	if err != nil {
		return nil, err
	}
	response.BaseMessage = bmsg

	return response, nil
}

func (m *UpdateHealthStat) GetProtocol() message.Protocol {
	return HealthResponseProtocol
}

func (m *UpdateHealthStat) ReadProcess(ctx context.Context, _i interceptor.Interceptor, connection interceptor.Connection) error {
	i, ok := _i.(*chat.ServerInterceptor)
	if !ok {
		return interceptor.ErrInterfaceMisMatch
	}

	s, err := i.GetState(connection)
	if err != nil {
		return err
	}

	return i.Health.Process(ctx, m, s)
}

func (m *UpdateHealthStat) WriteProcess(_ context.Context, _i interceptor.Interceptor, connection interceptor.Connection) error {
	s, ok := _i.(*chat.ServerInterceptor)
	if !ok {
		return interceptor.ErrInterfaceMisMatch
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

func (m *UpdateHealthStat) setConnectionStatus(startTime time.Time) {
	m.ConnectionStatus = types.ConnectionStateUp
	m.ConnectionUptime = m.getConnectionUptime(startTime)
}

func (m *UpdateHealthStat) getConnectionUptime(startTime time.Time) types.ConnectionUptime {
	return types.ConnectionUptime(time.Since(startTime))
}

func (m *UpdateHealthStat) setCPUUsage() error {
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

func (m *UpdateHealthStat) setMemoryUsage() error {
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

func (m *UpdateHealthStat) setNetworkUsage() error {
	// TODO: IMPLEMENT THIS LATER
	return nil
}

func (m *UpdateHealthStat) setLatency(requestTimestamp time.Time) error {
	latencyNano := time.Since(requestTimestamp).Nanoseconds()

	latencyMs := latencyNano / int64(time.Millisecond)

	if latencyMs < 0 {
		return fmt.Errorf("error while reading latency; err: 'latency is negative'")
	}

	m.Latency = types.LatencyMs(latencyMs)

	return nil
}
