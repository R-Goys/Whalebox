package container

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	Common "github.com/R-Goys/Whalebox/common"
	"github.com/R-Goys/Whalebox/pkg/log"
)

type Container struct {
	Pid        string `json:"pid"`
	Id         string `json:"id"`
	Name       string `json:"name"`
	Command    string `json:"command"`
	CreateTime string `json:"createTime"`
	Status     string `json:"status"`
}

var (
	RUNNING             = "running"
	STOPPED             = "stopped"
	EXIT                = "exited"
	DEFAULTINFOLOCATION = "/home/rinai/PROJECTS/Whalebox/example/example4/%s/"
	CONFIGNAME          = "config.json"
	CONTAINERLOGFILE    = "container.log"
)

func NewParentProcess(tty bool, volume, containerName string) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		log.Error("NewParentProcess: Failed to create pipe: " + err.Error())
		return nil, nil
	}
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWIPC | syscall.CLONE_NEWNET,
		// syscall.CLONE_NEWUSER,
		// UidMappings: []syscall.SysProcIDMap{
		// 	{
		// 		ContainerID: 0,
		// 		HostID:      0,
		// 		Size:        1,
		// 	},
		// },
		// GidMappings: []syscall.SysProcIDMap{
		// 	{
		// 		ContainerID: 0,
		// 		HostID:      1000,
		// 		Size:        1,
		// 	},
		// },
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		dirURL := fmt.Sprintf(DEFAULTINFOLOCATION, containerName)
		if err := os.MkdirAll(dirURL, 0622); err != nil {
			log.Error("NewParentProcess: Failed to create directory: " + err.Error())
			return nil, nil
		}
		stdLogFilePath := dirURL + CONTAINERLOGFILE
		stdLogFile, err := os.Create(stdLogFilePath)
		if err != nil {
			log.Error("NewParentProcess: Failed to create log file: " + err.Error())
			return nil, nil
		}
		cmd.Stdout = stdLogFile
	}
	cmd.ExtraFiles = []*os.File{readPipe}

	NewWorkSpace(Common.RootPath, Common.MntPath, volume)
	cmd.Dir = Common.MntPath
	log.Info(fmt.Sprintf("Command: %v", cmd))
	return cmd, writePipe
}

func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		log.Error("NewPipe: Failed to create pipe: " + err.Error())
		return nil, nil, err
	}
	log.Info(fmt.Sprintf("New pipe: read: %d, write: %d", read.Fd(), write.Fd()))
	return read, write, nil
}
