package types

type OutputMsg struct {
	Target  OutputTarget
	Message string
}

type OutputTarget string

var (
	StartupOutput OutputTarget = "startup"
	ErrorOutput   OutputTarget = "error"
	InfoOutput    OutputTarget = "info"
	ServerOutput  OutputTarget = "server"
	BuildOutput   OutputTarget = "build"
)

type FinishedSetupCmd struct {
	ImageId     string
	ContainerId string
}

type CallbackFunc func() error
