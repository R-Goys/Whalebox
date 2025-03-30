package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"syscall"

	"github.com/R-Goys/Whalebox/container"
	"github.com/R-Goys/Whalebox/pkg/log"
)

func stopContainer(containerName string) {
	pid, err := getPidByContainerName(containerName)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to get PID of container %s: %s", containerName, err))
		return
	}
	Pid, err := strconv.Atoi(pid)
	if err != nil {
		log.Error("Failed to convert PID to int: " + err.Error())
		return
	}
	if err := syscall.Kill(Pid, syscall.SIGTERM); err != nil {
		log.Error(fmt.Sprintf("Failed to stop container %s: %s", containerName, err))
		return
	}
	containerInfo, err := getContainerInfoByName(containerName)
	if err != nil {
		log.Error("Failed to get container info:" + err.Error())
		return
	}
	containerInfo.Status = container.STOPPED
	NewContainerInfo, err := json.Marshal(containerInfo)
	if err != nil {
		log.Error("Failed to marshal container info: " + err.Error())
		return
	}
	dir := fmt.Sprintf(container.DEFAULTINFOLOCATION, containerName)
	fileName := dir + "/" + container.CONFIGNAME
	if err := os.WriteFile(fileName, NewContainerInfo, 0622); err != nil {
		log.Error("Failed to write container info: " + err.Error())
		return
	}
	log.Info(containerName + " Container %s stopped")
}

func getContainerInfoByName(containerName string) (*container.Container, error) {
	dirURL := fmt.Sprintf(container.DEFAULTINFOLOCATION, containerName)
	configDir := dirURL + container.CONFIGNAME
	content, err := os.ReadFile(configDir)
	if err != nil {
		log.Error("Failed to read container config file: " + err.Error())
		return nil, err
	}
	var c container.Container
	if err := json.Unmarshal(content, &c); err != nil {
		log.Error("Failed to unmarshal container config file: " + err.Error())
		return nil, err
	}
	return &c, nil
}
