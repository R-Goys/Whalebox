package container

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	Common "github.com/R-Goys/Whalebox/common"
	"github.com/R-Goys/Whalebox/pkg/log"
)

func NewWorkSpace(imageName, containerName, volume string) {
	CreateReadOnlyLayer(imageName)
	CreateWriteLayer(containerName)
	CreateMountPoint(containerName, imageName)
	if volume != "" {
		volumeURLs := volumeUrlExtract(volume)
		length := len(volumeURLs)
		if length == 2 && volumeURLs[0] != "" && volumeURLs[1] != "" {
			MountVolume(containerName, volumeURLs)
			log.Info(fmt.Sprintf("Mount volume: %v", volumeURLs))
		} else {
			log.Info(fmt.Sprintf("Invalid volume format: %s", volume))
		}
	}
}

func volumeUrlExtract(volume string) []string {
	volumeURLs := strings.Split(volume, ":")
	return volumeURLs
}

func MountVolume(contianerName string, volumeURLs []string) {
	parentURL := volumeURLs[0]
	containerURL := volumeURLs[1]
	if err := os.Mkdir(parentURL, 0777); err != nil {
		log.Info("MountVolume, Mkdir parentURL error: " + err.Error())
	}
	mntURL := fmt.Sprintf(Common.MntPath, contianerName)
	containerVolumeURL := mntURL + "/" + containerURL
	log.Debug(fmt.Sprintf("MountVolume, parentURL: %s, containerURL: %s", parentURL, containerVolumeURL))
	if err := os.Mkdir(containerVolumeURL, 0777); err != nil {
		log.Info("MountVolume, Mkdir containerVolumeURL error: " + err.Error())
	}
	cmd := exec.Command("mount", "--bind", parentURL, containerVolumeURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Error("MountVolume, " + containerVolumeURL + " mount error: " + err.Error())
	}
}

func DeleteMountPointWithVolume(containerName string, volumeURLs []string) {
	mntURL := fmt.Sprintf(Common.MntPath, containerName)
	containerURL := mntURL + "/" + volumeURLs[1]
	cmd := exec.Command("umount", containerURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Error("unmountVolume error: " + err.Error())
	}

	cmd = exec.Command("umount", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Error("umount error: " + err.Error())
	}
	if err := os.RemoveAll(mntURL); err != nil {
		log.Error("DeleteMountPointWithVolume, RemoveAll mntURL error: " + err.Error())
	}
}
