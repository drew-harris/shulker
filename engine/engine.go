package engine

type Engine interface {
	StartServerCmd() error
}

type DockerEngine struct {
}

func (e *DockerEngine) StartServerCmd() error {
	return nil
}
