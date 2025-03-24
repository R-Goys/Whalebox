package main

import (
	"github/R-Goys/Whalebox/pkg/log"

	"errors"

	"github.com/urfave/cli/v2"
)

var initCommand = &cli.Command{
	Name:   "init",
	Usage:  "Init container process run user's process in container. Do not call it outside.",
	Action: initAction,
}

func initAction(c *cli.Context) error {
	log.Info("init command")
	cmd := c.Args().Get(0)
	log.Info("cmd is:", cmd)
	err := container.RunContainerInitProcess(cmd)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

var runCommand = &cli.Command{
	Name:    "run",
	Aliases: []string{"r"},
	Usage:   "Run a container",
	Action:  runAction,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "ti",
			Usage: "enable tty",
		},
	},
}

func runAction(c *cli.Context) error {
	if c.Args().Len() < 1 {
		log.Error("Please specify a container image name")
		return errors.New("Please specify a container image name")
	}
	cmd := c.Args().Get(0)
	tty := c.Bool("ti")
	Run(cmd, tty)
	return nil
}
