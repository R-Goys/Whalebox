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
		&cli.StringSliceFlag{
			Name:  "e",
			Usage: "Set environment variables for container",
		},
		&cli.StringSliceFlag{
			Name:  "p",
			Usage: "Publish a container's port to the host",
		},
		&cli.StringFlag{
			Name:  "net",
			Usage: "Set network mode for container",
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
var networkCommand = cli.Command{
	Name:  "network",
	Usage: "container network commands",
	Subcommands: []cli.Command{
		NetworkCreateCommand,
		ListNetWorkCommand,
		RemoveNetworkCommand,
	},
}

var NetworkCreateCommand = cli.Command{
	Name:  "create",
	Usage: "create a container network",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "driver",
			Usage: "network driver",
		},
		cli.StringFlag{
			Name:  "subnet",
			Usage: "subnet cidr",
		},
	},
	Action: CreateNetworkAction,
}

var ListNetWorkCommand = cli.Command{
	Name:   "list",
	Usage:  "list container network",
	Action: ListNetworkAction,
}

var RemoveNetworkCommand = cli.Command{
	Name:   "remove",
	Usage:  "remove container network",
	Action: RemoveNetworkAction,
}
