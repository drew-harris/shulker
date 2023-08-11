package engine

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/drewharris/shulker/commands"
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
	out, err := os.Create(filepath.FromSlash(h.pwd + "/shulkerbox.zip"))
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
	log(h.pwd)
	_, err := os.Stat(filepath.FromSlash(h.pwd + "/.shulkerbox/spigot.jar"))
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

	// Rebuild all plugins
	err = h.RebuildAllPlugins(log)
	if err != nil {
		return err
	}

	return nil
}

func (h *HostEngine) RebuildAllPlugins(log types.Logger) error {
	err := commands.RunExternalCommand(log, commands.Command{
		Name: "mvn",
		Args: []string{"clean", "package"},
	})
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	// copyFileContents(h.pwd + "/plugins/")

	// Move plugins

	return nil
}

// Not implemented
func (h *HostEngine) StartServer(log types.Logger) error { return nil }
func (h *HostEngine) Shutdown() error                    { return nil }
func (h *HostEngine) CanAttach() bool                    { return false }

func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
