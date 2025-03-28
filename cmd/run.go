package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	cgroup "github.com/R-Goys/Whalebox/cgroups"
	"github.com/R-Goys/Whalebox/container"
	"github.com/R-Goys/Whalebox/pkg/log"
)

func Run(tty bool, cmdArray []string, resource *cgroup.ResourceConfig, volume string) {
	parent, pipe := container.NewParentProcess(tty, volume)
	if parent == nil {
		log.Error("Failed to create parent process")
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err.Error())
		return
	}
	fmt.Println(parent.Process.Pid)
	cgroupManager := cgroup.NewCgroup("whalebox", strconv.Itoa(parent.Process.Pid))
	cgroupManager.Set(resource)
	sendInitCommand(cmdArray, pipe)
	parent.Wait()
	cgroupManager.Remove()
	mntURL := "/home/rinai/PROJECTS/Whalebox/example/example3/mnt"
	rootURL := "/home/rinai/PROJECTS/Whalebox/example/example3/"
	container.DeleteWorkSpace(rootURL, mntURL, volume)
	os.Exit(0)
}

// 发信号给init进程，告诉它要执行的命令
func sendInitCommand(cmdArray []string, pipe *os.File) {
	defer pipe.Close()
	commamd := strings.Join(cmdArray, " ")
	log.Info(fmt.Sprintf("Sending command to container: %s", commamd))
	if _, err := pipe.WriteString(commamd); err != nil {
		log.Error("SendMsg Error: " + err.Error())
	}
	pipe.Sync()
}
