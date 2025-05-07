package messages

import (
	"fmt"
	"math"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/health"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

var HealthResponseProtocol message.Protocol = "room:health_response"

// NOTE: BASIC HEALTH RESPONSE FOR ROOM MANAGEMENT, OTHER METRICS WILL BE DEALT WITH LATER

type HealthResponse struct {
	interceptor.BaseMessage
	health.Stat
	RequestTimeStamp int64         `json:"-"` // in nanoseconds
	RoomID           types.RoomID  `json:"room_id"`
	Validity         time.Duration `json:"validity"`
}

func NewHealthResponse(request *RequestHealth, validity time.Duration) (*HealthResponse, error) {
	response := &HealthResponse{}

	response.RequestTimeStamp = request.Timestamp
	response.RoomID = request.RoomID
	response.Validity = validity

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
	i, ok := _i.(*chat.ServerInterceptor)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	s, err := i.GetState(connection)
	if err != nil {
		return err
	}

	return i.Health.Process(m, s)
}

func (m *HealthResponse) Process(p interfaces.Processor, s *state.State) error {
	id, err := s.GetClientID()
	if err != nil {
		return err
	}

	if id != types.ClientID(m.CurrentHeader.Sender) {
		return fmt.Errorf("error while processing 'HealthResponse' message; err: 'sender id does not match'")
	}

	u, ok := p.(interfaces.CanUpdate)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	// NOTE: BE VERY CAREFUL WITH THIS. STAT IS PASSED AS POINTER. ANY CHANGES LATER TO HealthResponse WILL BE REFLECTED IN CanUpdate
	timer := time.NewTimer(m.Validity)
	defer timer.Stop()

	for {
		// NOTE: THIS IS A BLOCKING CALL. WE NEED TO WAIT FOR THE VALIDITY TO EXPIRE
		// NOTE: THIS ALSO MAKES SURE THAT THE ROOM ACTUALLY EXISTS BEFORE UPDATING THE HEALTH STATS
		select {
		case <-timer.C:
			return fmt.Errorf("error while processing 'HealthResponse' message; err: 'validity expired'")
		default:
			err := u.Update(m.RoomID, id, &m.Stat)
			if err == nil {
				return nil
			}
			fmt.Println("error while reading CPU usage; err: ", err.Error())
		}
	}
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

func (m *HealthResponse) setNetworkUsage() error {
	// TODO: IMPLEMENT THIS LATER
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
