//go:build go1.23

package network

import "net"

type KeepAliveConfig = net.KeepAliveConfig

func SetDialerTCPKeepAlive(dialer *net.Dialer, config KeepAliveConfig) {
	if config.Enable {
		dialer.KeepAliveConfig = config
	}
}

func SetListenerTCPKeepAlive(listenConfig *net.ListenConfig, config KeepAliveConfig) {
	if config.Enable {
		listenConfig.KeepAliveConfig = config
	}
}
