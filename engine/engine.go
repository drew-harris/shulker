package engine

import "github.com/drewharris/shulker/types"

type Engine interface {
	EnsureSetup(sub chan types.OutputMsg) error
	StartServer(sub chan types.OutputMsg) error
	RebuildAllPlugins(sub chan types.OutputMsg) error
	Shutdown() error
	// SendCommandToSpigot(cmd string) error
	CanAttach() bool
}
