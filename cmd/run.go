package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	cgroup "github.com/R-Goys/Whalebox/cgroups"
	"github.com/R-Goys/Whalebox/container"
	"github.com/R-Goys/Whalebox/pkg/log"
)

func Run(tty bool, cmdArray []string, resource *cgroup.ResourceConfig, volume, containerName, imageName string, envSlice []string, port []string, network string) {
	containerID := randStringBytes(12)
	if containerName == "" {
		containerName = containerID
	}
	parent, pipe := container.NewParentProcess(tty, volume, containerName, imageName, envSlice)
	if parent == nil {
		log.Error("Failed to create parent process")
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err.Error())
		return
	}
	fmt.Println("Container started, pid: ", parent.Process.Pid)
	containerName, err := RecordContainerInfo(parent.Process.Pid, cmdArray, containerName, containerID, volume, imageName)
	if err != nil {
		log.Error("Record container info error" + err.Error())
		return
	}
	cgroupManager := cgroup.NewCgroup("whalebox", strconv.Itoa(parent.Process.Pid))
	cgroupManager.Set(resource)
	if network != "" {
		// config container network
		network.Init()
		containerInfo := &container.Container{
			Id:          containerID,
			Pid:         strconv.Itoa(parent.Process.Pid),
			Name:        containerName,
			PortMapping: port,
		}
		if err := network.Connect(network, containerInfo); err != nil {
			log.Error("Error Connect Network %v", err)
			return
		}
	}
	sendInitCommand(cmdArray, pipe)
	if tty {
		parent.Wait()
		deleteContainerInfo(containerName)
		container.DeleteWorkSpace(containerName, volume)
		cgroupManager.Remove()
	}
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

func randStringBytes(n int) string {
	letterBytes := "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func RecordContainerInfo(ContainerPID int, commandArray []string, containerName, ContainerID, volume, imageName string) (string, error) {
	createTime := time.Now().Format("2006-01-02 15:04:05")
	command := strings.Join(commandArray, " ")
	containerInfo := &container.Container{
		Id:         ContainerID,
		Name:       containerName,
		Pid:        strconv.Itoa(ContainerPID),
		Volume:     volume,
		ImageName:  imageName,
		Command:    command,
		CreateTime: createTime,
		Status:     "running",
	}
	jsonBytes, err := json.Marshal(containerInfo)
	if err != nil {
		log.Error("Record container info error" + err.Error())
		return "", err
	}
	jsonStr := string(jsonBytes)
	log.Debug("Record container info: " + jsonStr)

	dir := fmt.Sprintf(container.DEFAULTINFOLOCATION, containerName)
	if err := os.MkdirAll(dir, 0622); err != nil {
		log.Error("Create container info dir error" + err.Error())
		return "", err
	}

	fileName := dir + "/" + container.CONFIGNAME
	file, err := os.Create(fileName)
	if err != nil {
		log.Error("Create container info file error" + err.Error())
		return "", err
	}
	defer file.Close()

	if _, err := file.WriteString(jsonStr); err != nil {
		log.Error("Write container info error" + err.Error())
		return "", err
	}
	return containerName, nil
}

func deleteContainerInfo(containerName string) {
	dirURL := fmt.Sprintf(container.DEFAULTINFOLOCATION, containerName)
	if err := os.RemoveAll(dirURL); err != nil {
		log.Error("Delete container info error" + err.Error())
	}
}
