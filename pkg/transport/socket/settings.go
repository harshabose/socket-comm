package socket

import (
	"time"
)

type connectionSettings struct {
	ReadTimout  time.Duration
	WriteTimout time.Duration
}

type Settings struct {
	connectionSettings
	IdleTimout     time.Duration
	ShutdownTimout time.Duration

	// NOTE: TLS SETTINGS ARE OPTIONAL
	TLSCertFile string
	TLSKeyFile  string

	MaxConnections    int
	ConnectionTimeout time.Duration
}

func NewDefaultSettings() *Settings {
	return &Settings{
		connectionSettings: connectionSettings{
			ReadTimout:  30 * time.Second,
			WriteTimout: 30 * time.Second,
		},
		IdleTimout:        120 * time.Second,
		ShutdownTimout:    30 * time.Second,
		TLSCertFile:       "",
		TLSKeyFile:        "",
		MaxConnections:    1000,
		ConnectionTimeout: 10 * time.Second,
	}
}
