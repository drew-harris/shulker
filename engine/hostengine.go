package engine

import (
	"errors"
	"net/http"
	"os"

	"github.com/drewharris/shulker/config"
	"github.com/drewharris/shulker/types"
)

type HostEngine struct {
}

func DownloadShulkerbox(pwd string) error {
	out, err := os.Create(pwd + "/shulker_data.zip")
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(config.ShulkerboxUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (h *HostEngine) EnsureSetup(sub chan types.OutputMsg) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	_, err = os.Stat(pwd + "/shulker_data/spigot")
	if errors.Is(err, os.ErrNotExist) {
		// Download file
		DownloadShulkerbox(pwd)
	} else if err != nil {
		return err // Worst case option
	}

	return nil
}

func (h *HostEngine) StartServer(sub chan types.OutputMsg) error       { return nil }
func (h *HostEngine) RebuildAllPlugins(sub chan types.OutputMsg) error { return nil }
func (h *HostEngine) Shutdown() error                                  { return nil }

// SendCommandToSpigot(cmd string) error
func (h *HostEngine) CanAttach() bool { return false }
