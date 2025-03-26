package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/R-Goys/Whalebox/cgroups"
	subsystem "github.com/R-Goys/Whalebox/cgroups/subsystems"
	"github.com/R-Goys/Whalebox/container"
	"github.com/R-Goys/Whalebox/pkg/log"
)

func Run(tty bool, cmdArray []string, resource *subsystem.ResourceConfig) {
	parent, pipe := container.NewParentProcess(tty)
	if parent == nil {
		log.Error("Failed to create parent process")
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err.Error())
		return
	}
	cgroupManager := cgroups.NewCgroupManager("whalebox")
	defer cgroupManager.Destroy()
	cgroupManager.Apply(parent.Process.Pid)
	cgroupManager.Set(resource)
	sendInitCommand(cmdArray, pipe)
	parent.Wait()
	os.Exit(0)
}

// 发信号给init进程，告诉它要执行的命令
func sendInitCommand(cmdArray []string, pipe *os.File) {
	commamd := strings.Join(cmdArray, " ")
	log.Info(fmt.Sprintf("Sending command to container: %s", commamd))
	pipe.WriteString(commamd)
	pipe.Close()
}
