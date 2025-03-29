package main

import (
	"fmt"
	"io"
	"os"

	"github.com/R-Goys/Whalebox/container"
	"github.com/R-Goys/Whalebox/pkg/log"
)

func logContainer(containerName string) {
	dirURL := fmt.Sprintf(container.DEFAULTINFOLOCATION, containerName)
	logFileLocation := dirURL + container.CONTAINERLOGFILE
	logFile, err := os.Open(logFileLocation)
	if err != nil {
		log.Error("failed to open log file" + logFileLocation)
		return
	}
	defer logFile.Close()
	content, err := io.ReadAll(logFile)
	if err != nil {
		log.Error("failed to read log file" + logFileLocation)
		return
	}
	log.Debug("log content: " + string(content))
	fmt.Fprint(os.Stdout, string(content))
}
