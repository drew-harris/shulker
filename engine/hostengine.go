package engine

import (
	"archive/zip"
	"errors"
	"net/http"
	"os"

	"github.com/drewharris/shulker/config"
	"github.com/drewharris/shulker/types"
)

type HostEngine struct {
	pwd string
}

func NewHostEngine() (*HostEngine, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return &HostEngine{pwd: pwd}, nil
}

func (h *HostEngine) DownloadShulkerbox() error {
	out, err := os.Create(h.pwd + "/shulker_data.zip")
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

func (h *HostEngine) EnsureSetup(log types.Logger) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	_, err = os.Stat(h.pwd + "/shulker_data/spigot")
	if errors.Is(err, os.ErrNotExist) {
		// Download file
		err := h.DownloadShulkerbox()
		if err != nil {
			return err
		}
		reader, err := zip.OpenReader(pwd + "/shulker_data.zip")

		// Extract file
		defer reader.Close()
		for _, file := range reader.File {
			err := os.MkdirAll(pwd+"/shulker_data/"+file.Name, 0755)
			if err != nil {
				return err
			}
		}

		if err != nil {
			return err
		}

	} else if err != nil {
		return err // Worst case option
	}

	return nil
}

// Not implemented
func (h *HostEngine) StartServer(log types.Logger) error       { return nil }
func (h *HostEngine) RebuildAllPlugins(log types.Logger) error { return nil }
func (h *HostEngine) Shutdown() error                          { return nil }

// SendCommandToSpigot(cmd string) error
func (h *HostEngine) CanAttach() bool { return false }
