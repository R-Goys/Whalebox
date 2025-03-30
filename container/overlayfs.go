package container

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	Common "github.com/R-Goys/Whalebox/common"
	"github.com/R-Goys/Whalebox/pkg/log"
)

func CreateReadOnlyLayer(imageName string) {
	unTarFolderURL := Common.RootPath + "/" + imageName + "/"
	imageURL := Common.RootPath + "/" + imageName + ".tar"
	exist, err := PathExists(unTarFolderURL)
	if err != nil {
		log.Error("CreateReadOnlyLayer, PathExists error: " + err.Error())
		return
	}
	if !exist {
		if err := os.MkdirAll(unTarFolderURL, 0777); err != nil {
			log.Error("CreateReadOnlyLayer, Mkdir error: " + err.Error())
			return
		}
		if _, err := exec.Command("tar", "-xvf", imageURL, "-C", unTarFolderURL).CombinedOutput(); err != nil {
			log.Error("CreateReadOnlyLayer, tar error: " + err.Error())
		}
	}
}

func CreateWriteLayer(containerName string) {
	writeURL := fmt.Sprintf(Common.WriteLayerURL, containerName)
	if err := os.MkdirAll(writeURL, 0777); err != nil {
		log.Debug("CreateWriteLayer, Mkdir error: " + err.Error())
	}
}

func CreateMountPoint(containerName string, imageName string) {
	mntURL := fmt.Sprintf(Common.MntPath, containerName)
	log.Debug("CreateMountPoint, mntURL: " + mntURL)
	if err := os.MkdirAll(mntURL, 0777); err != nil {
		log.Debug("CreateMountPoint, Mkdir mntURL error: " + err.Error())
		return
	}
	tmpWriteURL := fmt.Sprintf(Common.WriteLayerURL, containerName)
	tmpImageLocation := Common.RootPath + "/" + imageName

	workdirURL := fmt.Sprintf(Common.WorkDirURL, containerName)
	if err := os.MkdirAll(workdirURL, 0777); err != nil {
		log.Debug("CreateMountPoint, Mkdir Workdir error: " + err.Error())
		return
	}

	builder := strings.Builder{}
	builder.WriteString("lowerdir=")
	builder.WriteString(tmpImageLocation + ",")
	builder.WriteString("upperdir=")
	builder.WriteString(tmpWriteURL + ",")
	builder.WriteString("workdir=")
	builder.WriteString(workdirURL)

	cmd := exec.Command("mount", "-t", "overlay", "overlay", "-o", builder.String(), mntURL)
	log.Debug("CreateMountPoint, mount command: " + cmd.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Error("CreateMountPoint, mount error: " + err.Error())
		return
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func DeleteWorkSpace(containerName, volume string) {
	if volume != "" {
		volumeURLs := volumeUrlExtract(volume)
		if len(volumeURLs) == 2 && volumeURLs[0] != "" && volumeURLs[1] != "" {
			DeleteMountPointWithVolume(containerName, volumeURLs)
		} else {
			DeleteMountPoint(containerName)
		}
	} else {
		DeleteMountPoint(containerName)
	}
	DeleteWriteLayer(containerName)
	DeleteWorkdir(containerName)
}

func DeleteMountPoint(containerName string) {
	mntURL := fmt.Sprintf(Common.MntPath, containerName)
	cmd := exec.Command("umount", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Error("DeleteMountPoint, umount error: " + err.Error())
		return
	}

	if err := os.RemoveAll(mntURL); err != nil {
		log.Error("DeleteMountPoint, RemoveAll mntURL error: " + err.Error())
	}
}

func DeleteWriteLayer(containerName string) {
	writeURL := fmt.Sprintf(Common.WriteLayerURL, containerName)
	if err := os.RemoveAll(writeURL); err != nil {
		log.Error("DeleteWriteLayer, RemoveAll writeURL error: " + err.Error())
	}
}

func DeleteWorkdir(containerName string) {
	workdirURL := fmt.Sprintf(Common.WorkDirURL, containerName)
	if err := os.RemoveAll(workdirURL); err != nil {
		log.Error("DeleteWorkdir, RemoveAll workdirURL error: " + err.Error())
	}
}
