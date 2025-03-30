package main

import (
	"fmt"
	"os/exec"

	Common "github.com/R-Goys/Whalebox/common"
	"github.com/R-Goys/Whalebox/pkg/log"
)

func commitContainer(containerName, imageName string) {
	mntURL := fmt.Sprintf(Common.MntPath, containerName)
	mntURL += "/"
	imageTar := Common.RootPath + "/" + imageName + ".tar"
	log.Info(fmt.Sprintf("Committing container %s to %s", imageName, imageTar))

	if _, err := exec.Command("tar", "-czf", imageTar, "-C", mntURL, ".").CombinedOutput(); err != nil {
		log.Error(fmt.Sprintf("Tar folder error: %s", err))
	}
}
