package main

import (
	"os"

	"github.com/R-Goys/Whalebox/container"
	"github.com/R-Goys/Whalebox/pkg/log"
)

func Run(cmd string, tty bool) {
	//新建一个父进程来管理一个容器
	parent := container.NewParentProcess(tty, cmd)
	if err := parent.Start(); err != nil {
		log.Error(err.Error())
	}
	parent.Wait()
	os.Exit(-1)
}
