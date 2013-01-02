package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/bmatsuo/mflag"
)

type Options struct {
	Host string `mflag:"host,a remote host running iTunes"`
}

func init() {
	for _, c := range Commands {
		CommandMap[c.name] = c
	}
}

func main() {
	var opts Options
	fs, err := mflag.
		NewFlagSet("go-itunes").
		Mode(mflag.ModeLinear).
		Parse(os.Args[1:], &opts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	CommandHost = opts.Host
	args := fs.Args()

	if len(args) < 1 {
		fmt.Println("usage: go-itunes COMMAND [ARGUMENT ...]")
		os.Exit(1)
	}

	CommandMap[args[0]].action(args[1:])
}

var Commands = []*Command{
	{"status", func(args []string) {
		p, err := iTunesCommand(nil, "player state as string").Output()
		if err != nil {
			fmt.Println(err)
			return
		}
		status := string(p)

		if status == "stopped" {
			fmt.Println(status)
			iTunesCommand(os.Stdout, "play").Run()
			return
		}

		iTunesCommand(os.Stdout, strings.Join([]string{
			"(player state as string)",
			"\" -- \"",
			"(get artist of current track)",
			"\" - \"",
			"(get name of current track)",
			"\" [\" & (get rating of current track as integer / 20) & \" stars]\"",
		}, " & ")).Run()
	}},
	{"play", func(args []string) {
		iTunesCommand(os.Stdout, "play").Run()
	}},
	{"pause", func(args []string) {
		iTunesCommand(os.Stdout, "pause").Run()
	}},
	{"next", func(args []string) {
		iTunesCommand(os.Stdout, "next track").Run()
	}},
	{"prev", func(args []string) {
		iTunesCommand(os.Stdout, "previous track").Run()
	}},
	{"stop", func(args []string) {
		iTunesCommand(os.Stdout, "stop").Run()
	}},
	{"mute", func(args []string) {
		iTunesCommand(os.Stdout, "set mute to true").Run()
	}},
	{"vol", func(args []string) {
		if len(args) == 0 {
			iTunesCommand(os.Stdout, "sound volume as integer").Run()
			return
		}
		arg := args[0]

		var oldv, newv int64
		var err error
		if arg == "up" || arg == "down" {
			p, _ := iTunesCommand(nil, "sound volume as integer").Output()
			oldv, _ = strconv.ParseInt(string(bytes.TrimSpace(p)), 10, 64)
		}
		switch arg {
		case "up":
			newv = oldv + 10
		case "down":
			newv = oldv - 10
		default:
			newv, err = strconv.ParseInt(arg, 10, 64)
			if err != nil {
				fmt.Printf("invalid argument %q for \"vol\"", arg)
				return
			}
		}
		iTunesCommand(os.Stdout, fmt.Sprintf("set sound volume to %d", newv)).Run()
	}},
	{"rate", func(args []string) {
		if len(args) == 0 {
			p, _ := iTunesCommand(nil, "get rating of current track").Output()
			oldr, _ := strconv.ParseInt(string(bytes.TrimSpace(p)), 10, 64)
			fmt.Println(oldr / 20)
			return
		}

		var oldr, newr int64
		var err error
		if args[0] == "+" || args[0] == "-" {
			p, _ := iTunesCommand(nil, "get rating of current track as integer").Output()
			oldr, _ = strconv.ParseInt(string(bytes.TrimSpace(p)), 10, 64)
			if oldr == 0 {
				oldr = 50
			}
		}
		switch args[0] {
		case "+":
			newr = oldr + 20
		case "-":
			newr = oldr - 20
		default:
			newr, err = strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				fmt.Println("invalid argument %q for \"rate\"", args[0])
			}
		}

		iTunesCommand(os.Stdout, fmt.Sprintf("set rating of current track to %d", newr*20)).Run()
	}},
	{"open", func(args []string) {
		iTunesCommand(os.Stdout, "activate").Run()
	}},
	{"quit", func(args []string) {
		iTunesCommand(os.Stdout, "quit").Run()
	}},
}
var CommandMap = make(map[string]*Command, len(Commands))
