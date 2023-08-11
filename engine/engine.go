package engine

import "github.com/drewharris/shulker/types"

type Engine interface {
	EnsureSetup(log types.Logger) error
	StartServer(log types.Logger) error
	// TODO: boolean for cache
	RebuildAllPlugins(log types.Logger, disableCache bool) error
	Shutdown() error
	// SendCommandToSpigot(cmd string) error
	CanAttach() bool
}
