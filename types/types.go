package types

type ResponseMsg struct {
	Target  ResponseTarget
	Message string
}

type ResponseTarget string

var (
	StartupResponse ResponseTarget = "startup"
)
