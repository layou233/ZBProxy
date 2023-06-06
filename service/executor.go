package service

import (
	"context"

	"github.com/layou233/ZBProxy/config"
)

func ExecuteServices(ctx context.Context) {
	for _, s := range config.Config.Services {
		go StartNewService(ctx, s)
	}
}

func CleanupServices() {
	// sometimes after the program exits on Windows, the ports are still occupied and "listening".
	// so manually closes these listeners when the program exits.
	for _, listener := range Listeners {
		if listener != nil { // avoid null pointers
			listener.Close()
		}
	}
	Listeners = nil
}
