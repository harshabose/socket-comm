package types

type (
	ConnectionState  string
	ConnectionUptime int
	CPUUsage         struct {
		NumCores uint8     `json:"num_cores"`
		Percent  []float64 `json:"percent"`
	}
	MemoryUsage struct {
		Total          float32 `json:"total"`
		Used           float32 `json:"used"`
		UsedRatio      float32 `json:"used_ratio"`
		Available      float32 `json:"available"`
		AvailableRatio float32 `json:"available_ratio"`
	}
	NetworkUsage float64
	LatencyMs    int64
)

const (
	ConnectionStateUp   ConnectionState = "up"
	ConnectionStateDown ConnectionState = "down"
)
