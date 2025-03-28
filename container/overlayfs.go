package container

import (
	"os"
	"os/exec"
	"strings"

	"github.com/R-Goys/Whalebox/pkg/log"
)

func CreateReadOnlyLayer(RootURL string) {
	busyboxURL := RootURL + "busybox/"
	busyboxTarURL := RootURL + "busybox.tar"
	exist, err := PathExists(busyboxURL)
	if err != nil {
		log.Error("CreateReadOnlyLayer, PathExists error: " + err.Error())
		return
	}
	if !exist {
		if err := os.Mkdir(busyboxURL, 0777); err != nil {
			log.Error("CreateReadOnlyLayer, Mkdir error: " + err.Error())
			return
		}
		if _, err := exec.Command("tar", "-xvf", busyboxTarURL, "-C", busyboxURL).CombinedOutput(); err != nil {
			log.Error("CreateReadOnlyLayer, tar error: " + err.Error())
		}
	}
}

func CreateWriteLayer(RootURL string) {
	writeURL := RootURL + "writeLayer/"
	if err := os.Mkdir(writeURL, 0777); err != nil {
		log.Debug("CreateWriteLayer, Mkdir error: " + err.Error())
	}
}

func CreateMountPoint(RootURL, mntURL string) {
	if err := os.Mkdir(mntURL, 0777); err != nil {
		log.Debug("CreateMountPoint, Mkdir mntURL error: " + err.Error())
		return
	}

	workdirURL := RootURL + "work"
	if err := os.Mkdir(workdirURL, 0777); err != nil {
		log.Debug("CreateMountPoint, Mkdir Workdir error: " + err.Error())
		return
	}

	builder := strings.Builder{}
	builder.WriteString("lowerdir=")
	builder.WriteString(RootURL + "busybox,")
	builder.WriteString("upperdir=")
	builder.WriteString(RootURL + "writeLayer,")
	builder.WriteString("workdir=")
	builder.WriteString(RootURL + "work")

	cmd := exec.Command("mount", "-t", "overlay", "overlay", "-o", builder.String(), mntURL)
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

func DeleteWorkSpace(rootURL, mntURL, volume string) {
	if volume != "" {
		volumeURLs := volumeUrlExtract(volume)
		if len(volumeURLs) == 2 && volumeURLs[0] != "" && volumeURLs[1] != "" {
			DeleteMountPointWithVolume(rootURL, mntURL, volumeURLs)
		} else {
			DeleteMountPoint(rootURL, mntURL)
		}
	} else {
		DeleteMountPoint(rootURL, mntURL)
	}
	DeleteWriteLayer(rootURL)
	DeleteWorkdir(rootURL)
}

func DeleteMountPoint(rootURL, mntURL string) {
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

func DeleteWriteLayer(rootURL string) {
	writeURL := rootURL + "writeLayer/"
	if err := os.RemoveAll(writeURL); err != nil {
		log.Error("DeleteWriteLayer, RemoveAll writeURL error: " + err.Error())
	}
}

func DeleteWorkdir(rootURL string) {
	workdirURL := rootURL + "work"
	if err := os.RemoveAll(workdirURL); err != nil {
		log.Error("DeleteWorkdir, RemoveAll workdirURL error: " + err.Error())
	}
}
