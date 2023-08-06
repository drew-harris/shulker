package types

type ResponseMsg struct {
	Target  ResponseTarget
	Message string
}

type ResponseTarget string

var (
	StartupResponse ResponseTarget = "startup"
	ErrorResponse   ResponseTarget = "error"
)

type FinishedSetupCmd struct {
	ImageId     string
	ContainerId string
}
