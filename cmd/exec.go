package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	Common "github.com/R-Goys/Whalebox/common"
	"github.com/R-Goys/Whalebox/container"
	_ "github.com/R-Goys/Whalebox/nsenter"
	"github.com/R-Goys/Whalebox/pkg/log"
)

func execContainer(containerName string, cmdArray []string) {
	pid, err := getPidByContainerName(containerName)
	if err != nil {
		log.Error("Failed to get pid by container name" + err.Error())
		return
	}
	cmdStr := strings.Join(cmdArray, " ")
	log.Info("Executing command in container " + containerName + " : " + cmdStr)
	cmd := exec.Command("/proc/self/exe", "exec")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	os.Setenv(Common.ENV_EXEC_PID, pid)
	os.Setenv(Common.ENV_EXEC_CMD, cmdStr)
	if err := cmd.Run(); err != nil {
		log.Error("Failed to execute command in container " + containerName + " : " + err.Error())
	}
}

func getPidByContainerName(containerName string) (string, error) {
	dirURL := fmt.Sprintf(container.DEFAULTINFOLOCATION, containerName)
	configFilePath := dirURL + container.CONFIGNAME
	configBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		return "", err
	}
	var containerInfo container.Container
	if err := json.Unmarshal(configBytes, &containerInfo); err != nil {
		return "", err
	}
	return containerInfo.Pid, nil
}
