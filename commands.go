package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/ota42y/chatwork-command/command"
)

var GlobalFlags = []cli.Flag{}

var Commands = []cli.Command{

	{
		Name:   "send",
		Usage:  "",
		Action: command.CmdSend,
		Flags:  []cli.Flag{},
	},

	{
		Name:   "show",
		Usage:  "",
		Action: command.CmdShow,
		Flags:  []cli.Flag{},
	},
}

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}
