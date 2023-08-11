package engine

import (
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/drewharris/shulker/config"
	"github.com/drewharris/shulker/types"
	"github.com/xyproto/unzip"
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
	out, err := os.Create(h.pwd + "/shulkerbox.zip")
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(config.ShulkerboxUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	io.Copy(out, resp.Body)

	return nil
}

func (h *HostEngine) EnsureSetup(log types.Logger) error {
	_, err := os.Stat(h.pwd + "/shulkerbox/spigot")
	if errors.Is(err, os.ErrNotExist) {
		log("Downloading shulkerbox...")
		// Download file
		err := h.DownloadShulkerbox()
		if err != nil {
			return err
		}
		log("Downloaded!")
		log("Unzipping...")
		err = unzip.Extract("shulkerbox.zip", ".shulkerbox")
		if err != nil {
			return err
		}
		log("Extracted.")

		// Delete zip file
		log("Removing archive")
		err = os.Remove("shulkerbox.zip")
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
