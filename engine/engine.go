package engine

import "github.com/drewharris/shulker/types"

type Engine interface {
	EnsureSetup(log types.Logger) error
	StartServer(log types.Logger) error
	RebuildAllPlugins(log types.Logger) error
	Shutdown() error
	// SendCommandToSpigot(cmd string) error
	CanAttach() bool
}
