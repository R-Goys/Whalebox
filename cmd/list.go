package main

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/R-Goys/Whalebox/container"
	"github.com/R-Goys/Whalebox/pkg/log"
)

func listContainers() {
	dirURL := fmt.Sprintf(container.DEFAULTINFOLOCATION, "")
	dirURL = dirURL[:len(dirURL)-1]
	files, err := os.ReadDir(dirURL)
	if err != nil {
		log.Error("Error reading directory: " + err.Error())
		return
	}
	var containers []*container.Container
	for _, f := range files {
		tmpContainer, err := GetContainerInfo(f)
		if err != nil {
			log.Error("Error getting container info: " + err.Error())
			continue
		}
		containers = append(containers, tmpContainer)
	}
	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)

	fmt.Fprint(w, "ID\tNAME\tPID\tSTATUS\tCOMMAND\tCREATED\n")
	for _, c := range containers {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			c.Id,
			c.Name,
			c.Pid,
			c.Status,
			c.Command,
			c.CreateTime)
	}
	if err := w.Flush(); err != nil {
		log.Error("Error flushing writer: " + err.Error())
	}
}

func GetContainerInfo(file os.DirEntry) (*container.Container, error) {
	containerName := file.Name()
	configFileDir := fmt.Sprintf(container.DEFAULTINFOLOCATION, containerName)
	configFileDir = configFileDir + container.CONFIGNAME
	content, err := os.ReadFile(configFileDir)
	if err != nil {
		log.Error("Error reading file: " + err.Error())
		return nil, err
	}
	var c container.Container
	if err := json.Unmarshal(content, &c); err != nil {
		log.Error("Error unmarshalling json: " + err.Error())
		return nil, err
	}
	return &c, nil
}
