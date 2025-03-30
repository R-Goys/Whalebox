package main

import (
	"github.com/urfave/cli"
)

var initCommand = cli.Command{
	Name:   "init",
	Usage:  "Init container process run user's process in container. Do not call it outside.",
	Action: initAction,
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
		&cli.BoolFlag{
			Name:  "d",
			Usage: "detach container",
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
		&cli.StringFlag{
			Name:  "name",
			Usage: "Set container name",
		},
	},
}

var commitCommand = cli.Command{
	Name:   "commit",
	Usage:  "Commit container changes to image",
	Action: commitAction,
}

var listCommand = cli.Command{
	Name:   "ps",
	Usage:  "List all running containers",
	Action: listAction,
}

var logCommand = cli.Command{
	Name:   "logs",
	Usage:  "Show container logs",
	Action: logAction,
}

var execCommand = cli.Command{
	Name:   "exec",
	Usage:  "Run a command in a running container",
	Action: execAction,
}

var stopCommand = cli.Command{
	Name:   "stop",
	Usage:  "Stop a running container",
	Action: stopAction,
}

var removeCommand = cli.Command{
	Name:   "rm",
	Usage:  "Remove a container",
	Action: removeAction,
}
