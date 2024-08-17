//go:build !go1.23

package network

import (
	"net"
	"time"
)

// KeepAliveConfig is a copy of net.KeepAliveConfig backported from go1.23.
type KeepAliveConfig struct {
	// If Enable is true, keep-alive probes are enabled.
	Enable bool

	// Idle is the time that the connection must be idle before
	// the first keep-alive probe is sent.
	// If zero, a default value of 15 seconds is used.
	Idle time.Duration

	// Interval is the time between keep-alive probes.
	// If zero, a default value of 15 seconds is used.
	Interval time.Duration

	// Count is the maximum number of keep-alive probes that
	// can go unanswered before dropping a connection.
	// If zero, a default value of 9 is used.
	Count int
}

func SetDialerTCPKeepAlive(dialer *net.Dialer, config KeepAliveConfig) {
	if config.Enable {
		dialer.KeepAlive = config.Idle
	}
}

func SetListenerTCPKeepAlive(listenConfig *net.ListenConfig, config KeepAliveConfig) {
	if config.Enable {
		listenConfig.KeepAlive = config.Idle
	}
}
