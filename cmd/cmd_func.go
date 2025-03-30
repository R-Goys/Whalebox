package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	cgroup "github.com/R-Goys/Whalebox/cgroups"
	Common "github.com/R-Goys/Whalebox/common"
	"github.com/R-Goys/Whalebox/container"
	"github.com/R-Goys/Whalebox/pkg/log"
	"github.com/urfave/cli"
)

func initAction(c *cli.Context) error {
	log.Info("init command begin")
	err := container.RunContainerInitProcess()
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func runAction(c *cli.Context) error {
	if len(c.Args()) < 1 {
		log.Error("Please specify a container image name")
		return errors.New("please specify a container image name")
	}
	var cmdArray []string
	for i := 0; i < len(c.Args()); i++ {
		cmdArray = append(cmdArray, c.Args()[i])
	}
	//此处拿到第一条参数
	//此处是获取-ti的参数
	tty := c.Bool("ti")
	detch := c.Bool("d")
	log.Debug("tty: " + strconv.FormatBool(tty) + " detch: " + strconv.FormatBool(detch))
	if tty && detch {
		fmt.Println("Please specify only one of -ti and -d")
		log.Error("Please specify only one of -ti and -d")
		return errors.New("please specify only one of -ti and -d")
	}
	resource := &cgroup.ResourceConfig{
		MemoryLimit: c.String("m"),
		CpuShares:   c.String("cpushare"),
		CpuSet:      c.String("cpuset"),
	}
	volume := c.String("v")
	containerName := c.String("name")
	envSlice := c.StringSlice("e")
	imageName := cmdArray[0]
	cmdArray = cmdArray[1:]
	re, _ := json.Marshal(resource)
	log.Debug(string(re))
	Run(tty, cmdArray, resource, volume, containerName, imageName, envSlice)
	return nil
}

func commitAction(c *cli.Context) error {
	if len(c.Args()) < 2 {
		log.Error("-Commit: Please specify a container name and image name")
		return errors.New("please specify a container name and image name")
	}
	containerName := c.Args().Get(0)
	imageName := c.Args().Get(1)
	commitContainer(containerName, imageName)
	return nil
}

func listAction(c *cli.Context) error {
	listContainers()
	return nil
}

func logAction(c *cli.Context) error {
	if len(c.Args()) == 0 {
		log.Error("please provide a containerName to log")
		return fmt.Errorf("please provide a containerName to log")
	}
	containerName := c.Args().Get(0)
	logContainer(containerName)
	return nil
}

func execAction(c *cli.Context) error {
	if os.Getenv(Common.ENV_EXEC_PID) != "" {
		log.Info("pid callback pid: " + strconv.Itoa(os.Getegid()))
		return nil
	}
	if len(c.Args()) < 2 {
		log.Error("Please specify a container name and command to execute")
		return errors.New("please specify a container name and command to execute")
	}
	containerName := c.Args().Get(0)
	var cmdArray []string
	for _, arg := range c.Args()[1:] {
		cmdArray = append(cmdArray, arg)
	}
	log.Debug("exec containerName: " + containerName + " cmdArray: " + fmt.Sprintf("%v", cmdArray))
	execContainer(containerName, cmdArray)
	return nil
}

func stopAction(c *cli.Context) error {
	if len(c.Args()) < 1 {
		log.Error("Please specify a container name to stop")
		return errors.New("please specify a container name to stop")
	}
	containerName := c.Args().Get(0)
	stopContainer(containerName)
	return nil
}

func removeAction(c *cli.Context) error {
	if len(c.Args()) < 1 {
		log.Error("Please specify a container name to remove")
		return errors.New("please specify a container name to remove")
	}
	containerName := c.Args().Get(0)
	removeContainer(containerName)
	return nil
}
