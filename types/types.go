package types

type OutputMsg struct {
	Target  OutputTarget
	Message string
}

type OutputTarget string

var (
	StartupOutput OutputTarget = "startup"
	ErrorOutput   OutputTarget = "error"
	ServerOutput  OutputTarget = "server"
	BuildOutput   OutputTarget = "build"
)

type CallbackFunc func() error

type QuickMsg int

const (
	DoneBuilding QuickMsg = iota
	BuildStarted
	ErrorBuilding
	FinishedSetup
	FinishedServerStart
)

type Logger func(msg string)
