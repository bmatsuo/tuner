package main

import (
	"fmt"
	"os"

	"github.com/bmatsuo/mflag"
)

var CommandMap map[string]*Command

type Options struct {
	Host string `mflag:"host,a remote host running iTunes"`
}

func init() {
	CommandMap = make(map[string]*Command, len(Commands))
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

	err = CommandMap[args[0]].action(args[1:])
	if err != nil {
		fmt.Println(err)
	}
}
