package commands

import (
	"bufio"
	"os"
	"os/exec"
	"strings"

	"github.com/drewharris/shulker/types"
)

type Command struct {
	Name string
	Args []string
	Dir  string
}

func RunExternalCommand(log types.Logger, command Command) error {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	cmd := exec.Command(command.Name, command.Args...)
	if command.Dir != "" {
		cmd.Dir = command.Dir
	} else {
		cmd.Dir = cwd
	}

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

	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func GetCommandOutput(command Command) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	cmd := exec.Command(command.Name, command.Args...)
	if command.Dir != "" {
		cmd.Dir = command.Dir
	} else {
		cmd.Dir = cwd
	}

	raw, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(raw), nil
}
