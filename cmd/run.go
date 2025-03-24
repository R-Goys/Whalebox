package main

import (
	"github/R-Goys/Whalebox/container"
	"github/R-Goys/Whalebox/pkg/log"
	"os"
)

func Run(cmd string, tty bool) {
	parent := container.NewParentProcess(tty, cmd)
	if err := parent.Start(); err != nil {
		log.Error(err.Error())
	}
	parent.Wait()
	os.Exit(-1)
}
