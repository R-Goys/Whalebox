package main

import (
	"encoding/json"
	"fmt"

	cgroup "github.com/R-Goys/Whalebox/cgroups"
	"github.com/R-Goys/Whalebox/container"
	"github.com/R-Goys/Whalebox/pkg/log"

	"errors"

	"github.com/urfave/cli"
)

var initCommand = cli.Command{
	Name:   "init",
	Usage:  "Init container process run user's process in container. Do not call it outside.",
	Action: initAction,
}

func initAction(c *cli.Context) error {
	log.Info("init command begin")
	err := container.RunContainerInitProcess()
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

var runCommand = cli.Command{
	Name: "run",
	Usage: `Run a container With namespace and cgroup limit.
		    ./cmd run -ti [command]`,
	Action: runAction,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "ti",
			Usage: "enable tty",
		},
		&cli.StringFlag{
			Name:  "m",
			Usage: "Set memory limit for container",
		},
		&cli.StringFlag{
			Name:  "cpuset",
			Usage: "Set CPU limit for container",
		},
		&cli.StringFlag{
			Name:  "cpushare",
			Usage: "Set CPU share for container",
		},
		&cli.StringFlag{
			Name:  "v",
			Usage: "Set volume for container",
		},
	},
}

func runAction(c *cli.Context) error {
	if len(c.Args()) < 1 {
		log.Error("Please specify a container image name")
		return errors.New("please specify a container image name")
	}
	var cmdArray []string
	for i := 0; i < len(c.Args()); i++ {
		log.Debug(fmt.Sprintf("Arg[%d]: %s", i, c.Args()[i]))
		cmdArray = append(cmdArray, c.Args()[i])
	}
	//此处拿到第一条参数
	//此处是获取-ti的参数
	tty := c.Bool("ti")
	resource := &cgroup.ResourceConfig{
		MemoryLimit: c.String("m"),
		CpuShares:   c.String("cpushare"),
		CpuSet:      c.String("cpuset"),
	}
	volume := c.String("v")
	re, _ := json.Marshal(resource)
	log.Debug(string(re))
	Run(tty, cmdArray, resource, volume)
	return nil
}
