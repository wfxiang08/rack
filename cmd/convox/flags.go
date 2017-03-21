package main

import "gopkg.in/urfave/cli.v1"

var appFlag = cli.StringFlag{
	Name:  "app, a",
	Usage: "app name inferred from current directory if not specified",
}

var rackFlag = cli.StringFlag{
	Name:  "rack",
	Usage: "rack name",
}

var notifyFlag = cli.BoolFlag{
	Name:  "notify",
	Usage: "send a desktop notification when operation is complete",
}
