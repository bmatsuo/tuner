package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

var CommandHost string

type Command struct {
	name   string
	action func(args []string)
}

func iTunesScript(script string) string {
	if CommandHost != "" {
		return fmt.Sprintf(`tell application "iTunes" of machine %q to %s`, CommandHost, script)
	}
	return fmt.Sprintf(`tell application "iTunes" to %s`, script)
}

func iTunesCommand(w io.Writer, script string, args ...string) *exec.Cmd {
	command := exec.Command("osascript", "-e", iTunesScript(script))
	command.Stdout = w
	command.Stderr = os.Stderr
	return command
}
