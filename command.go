package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

var CommandHost string

type UsageOptions []struct{ option, desc string }

type Command struct {
	name    string
	desc    string
	usage   string
	options UsageOptions
	action  func(args []string) error
}

var Commands = []*Command{
	{
		"help",
		"show this list of commands",
		"[command]",
		UsageOptions{{"command", "command to show detailed documentation for"}},
		func(args []string) error {
			if len(args) == 0 {
				w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
				for _, v := range CommandMap {
					fmt.Fprintln(w, v.name+"\t"+v.desc)
				}
				w.Flush()
				return nil
			}

			cmd, ok := CommandMap[args[0]]
			if !ok {
				return fmt.Errorf("unknown command %q", args[0])
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 4, 4, ' ', tabwriter.TabIndent)
			fmt.Fprintln(w)
			fmt.Fprintln(w, cmd.desc)
			w.Flush()
			fmt.Fprintln(w)
			fmt.Fprintln(w, "usage:", "go-itunes", cmd.name, cmd.usage)
			w.Flush()
			if len(cmd.options) > 0 {
				fmt.Fprintln(w)
				for _, opt := range cmd.options {
					fmt.Fprintln(w, "\t"+opt.option+"\t"+opt.desc)
				}
				w.Flush()
			}
			return nil
		},
	},
	{
		"status",
		"player status (playing, paused, stopped)",
		"",
		nil,
		func(args []string) error {
			return Script{"player state as string"}.Run()
		},
	},
	{
		"info",
		"current track info",
		"",
		nil,
		func(args []string) error {
			script := Script{
				"if exists (current track) then",
				strings.Join([]string{
					"(get artist of current track)",
					"\" - \"",
					"(get name of current track)",
					"\" [\" & (get rating of current track as integer / 20) & \" stars]\"",
				}, " & "),
				"else",
				"\"no track playing\"",
				"end if",
			}
			return script.Run()
		},
	},
	{
		"play",
		"start playing the current track",
		"",
		nil,
		func(args []string) error {
			return Script{"play"}.Run()
		},
	},
	{
		"pause",
		"pause the current track",
		"",
		nil,
		func(args []string) error {
			return Script{"pause"}.Run()
		},
	},
	{
		"next",
		"skip to the next track",
		"",
		nil,
		func(args []string) error {
			return Script{"next track"}.Run()
		},
	},
	{
		"prev",
		"skip to the previous track",
		"",
		nil,
		func(args []string) error {
			return Script{"previous track"}.Run()
		},
	},
	{
		"stop",
		"stop playback",
		"",
		nil,
		func(args []string) error {
			return Script{"stop"}.Run()
		},
	},
	{
		"mute",
		"mute/unmute playback",
		"",
		nil,
		func(args []string) error {
			script := Script{
				"if mute then",
				"set mute to false",
				"\"unmuted\"",
				"else",
				"set mute to true",
				"\"muted\"",
				"end if",
			}
			return script.Run()
		},
	},
	{
		"vol",
		"adjust playback volume",
		"LEVEL|up|+|down|-",
		UsageOptions{
			{"LEVEL", "a volume level (0-100)"},
			{"up, +", "nudge volume up 10 levels"},
			{"down, -", "nudge volume down 10 levels"},
		},
		func(args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("usage: go-itunes vol +|-|VOLUME")
			}
			arg := args[0]

			var oldv, newv int64
			var err error
			if arg == "up" || arg == "down" || arg == "+" || arg == "-" {
				out, _ := Script{"sound volume as integer"}.OutputString()
				oldv, _ = strconv.ParseInt(strings.TrimSpace(out), 10, 64)
			}
			switch arg {
			case "up", "+":
				newv = oldv + 10
			case "down", "-":
				newv = oldv - 10
			default:
				newv, err = strconv.ParseInt(arg, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid argument %q for \"vol\"", arg)
				}
			}
			return Script{fmt.Sprintf("set sound volume to %d", newv)}.Run()
		},
	},
	{
		"rate",
		"adjust rating of the current track",
		"RATING|up|+|down|-",
		UsageOptions{
			{"RATING", "a rating (in stars)"},
			{"up, +", "increase rating by one star"},
			{"down, -", "decrease rating by one star"},
		},
		func(args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("usage: go-itunes rate RATING")
			}

			var oldr, newr int64
			var err error
			if args[0] == "up" || args[0] == "down" || args[0] == "+" || args[0] == "-" {
				out, _ := Script{"get rating of current track as integer"}.OutputString()
				oldr, _ = strconv.ParseInt(strings.TrimSpace(out), 10, 64)
				if oldr == 0 {
					oldr = 50
				}
			}
			switch args[0] {
			case "up", "+":
				newr = oldr + 20
			case "down", "-":
				newr = oldr - 20
			default:
				newr, err = strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					return fmt.Errorf("invalid argument %q for \"rate\"", args[0])
				}
			}

			return Script{fmt.Sprintf("set rating of current track to %d", newr*20)}.Run()
		},
	},
	{
		"open",
		"open iTunes",
		"",
		nil,
		func(args []string) error {
			return Script{"activate"}.Run()
		},
	},
	{
		"quit",
		"quit iTunes",
		"",
		nil,
		func(args []string) error {
			return Script{"quit"}.Run()
		},
	},
}
