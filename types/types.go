package types

type OutputMsg struct {
	Target  OutputTarget
	Message string
}

type OutputTarget string

var (
	StartupOutput OutputTarget = "startup"
	ErrorOutput   OutputTarget = "error"
)

type FinishedSetupCmd struct {
	ImageId     string
	ContainerId string
}
