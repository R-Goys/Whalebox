package main

import (
	"os"

	"github.com/R-Goys/Whalebox/pkg/log"

	"github.com/urfave/cli"
)

const (
	AppName = "Whalebox"
	Version = "0.1.0"
	Usage   = "A container runtime based on containerd"
)

func main() {
	app := cli.NewApp()
	app.Name = AppName
	app.Version = Version
	app.Usage = Usage

	app.Commands = []cli.Command{
		initCommand,
		runCommand,
	}

	app.Before = func(c *cli.Context) error {
		log.InitLogger()
		log.Info("Starting Whalebox...")
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Error(err.Error())
	}
}
