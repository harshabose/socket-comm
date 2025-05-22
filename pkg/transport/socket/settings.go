package socket

import (
	"crypto/tls"
	"errors"
	"fmt"
	"time"
)

var ErrSettingsInvalid = errors.New("server settings invalid")

type connectionSettings struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Settings struct {
	connectionSettings
	Address           string
	Port              uint16
	ReadHeaderTimeout time.Duration
	IdleTimout        time.Duration
	ShutdownTimout    time.Duration

	// NOTE: TLS SETTINGS ARE OPTIONAL
	TLSCertFile string
	TLSKeyFile  string

	MaxConnections    int
	ConnectionTimeout time.Duration

	PopMessageTimeout time.Duration
	PushMessageTimout time.Duration
}

func (s Settings) Validate() error {
	// TODO: IMPLEMENT THIS
	return ErrSettingsInvalid
}

func NewDefaultSettings() Settings {
	return Settings{
		connectionSettings: connectionSettings{
			ReadTimeout:  time.Second,
			WriteTimeout: time.Second,
		},
		ReadHeaderTimeout: 5 * time.Second,
		ConnectionTimeout: 10 * time.Second,
		IdleTimout:        120 * time.Second,
		ShutdownTimout:    30 * time.Second,
		TLSCertFile:       "",
		TLSKeyFile:        "",
		MaxConnections:    1000,
		PopMessageTimeout: 30 * time.Second,
		PushMessageTimout: 30 * time.Second,
	}
}

func GetTLSV1(tlsCertPath, tlsKeyFile string) (*tls.Config, error) {
	var tlsConfig *tls.Config
	if tlsCertPath != "" && tlsKeyFile != "" {
		cert, err := tls.LoadX509KeyPair(tlsCertPath, tlsKeyFile)
		if err != nil {
			return nil, fmt.Errorf("loading TLS certificates: %w", err)
		}

		tlsConfig = &tls.Config{
			Certificates:     []tls.Certificate{cert},
			MinVersion:       tls.VersionTLS12,
			CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
		}

		return tlsConfig, nil
	}

	return nil, errors.New("invalid cert or/and file path")
}
