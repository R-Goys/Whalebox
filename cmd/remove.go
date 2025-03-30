package main

import (
	"fmt"

	cgroup "github.com/R-Goys/Whalebox/cgroups"
	"github.com/R-Goys/Whalebox/container"
	"github.com/R-Goys/Whalebox/pkg/log"
)

func removeContainer(containerName string) error {
	containerInfo, err := getContainerInfoByName(containerName)
	if err != nil {
		log.Error("Error getting container info: " + err.Error())
		return err
	}
	if containerInfo.Status != container.STOPPED {
		log.Error("Container is not stopped, cannot remove it")
		return fmt.Errorf("container is not stopped, cannot remove it")
	}
	cgroupManager := cgroup.GetCgroup("whalebox", containerInfo.Pid)
	volume := containerInfo.Volume
	deleteContainerInfo(containerName)
	container.DeleteWorkSpace(containerName, volume)
	cgroupManager.Remove()
	log.Info("Container removed successfully")
	return nil
}
