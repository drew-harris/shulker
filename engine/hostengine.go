package engine

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/drewharris/shulker/commands"
	"github.com/drewharris/shulker/config"
	"github.com/drewharris/shulker/types"
	"github.com/xyproto/unzip"
)

type HostEngine struct {
	pwd    string
	config config.Config
	server *exec.Cmd
}

func NewHostEngine(config config.Config) (*HostEngine, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return &HostEngine{pwd: pwd, config: config}, nil
}

func (h *HostEngine) DownloadShulkerbox() error {
	out, err := os.Create(filepath.FromSlash(h.pwd + "/shulkerbox.zip"))
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(h.config.ShulkerboxUrl)
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
	err = h.RebuildAllPlugins(log, false)
	if err != nil {
		return err
	}

	return nil
}

func (h *HostEngine) RebuildAllPlugins(log types.Logger, disableCache bool) error {
	log("Building all plugins...")
	err := commands.RunExternalCommand(log, commands.Command{
		Name: "mvn",
		Args: []string{"package", "-Dmaven.build.cache.enabled=false"},
	})
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	for _, pluginPath := range h.config.PluginCopyPaths {
		err = copyFileContents(filepath.FromSlash(pluginPath), filepath.FromSlash(h.pwd+"/.shulkerbox/plugins/"+filepath.Base(pluginPath)))
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *HostEngine) StartServer(log types.Logger) error {
	var baseDir = filepath.FromSlash(h.pwd + "/.shulkerbox/")
	cmdtmp := commands.Command{
		Name: "java",
		Dir:  baseDir,
		Args: []string{"-Xms1024M", "-Xmx2048M", "-Dfile.encoding=UTF-8", "-jar", "spigot.jar", "--world-dir", "./worlds", "nogui"},
	}

	cmd := exec.Command(cmdtmp.Name, cmdtmp.Args...)
	cmd.Dir = baseDir

	// Display output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	// Combine stdout and stderr so that both are captured
	cmd.Stderr = cmd.Stdout

	// Start the command
	if err := cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdout) // Scanner doesn't return newline byte
	for scanner.Scan() {
		log(strings.ReplaceAll(scanner.Text(), "\n", ""))
	}

	h.server = cmd
	return nil
}

// Not implemented
func (h *HostEngine) Shutdown() error {
	if h.server != nil {
		if h.server.Process != nil {
			err := h.server.Process.Kill()
			return err
		}
	}
	return nil
}

func (h *HostEngine) CanAttach() bool { return false }

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
