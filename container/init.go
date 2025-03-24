package container

import (
	"fmt"
	"github/R-Goys/Whalebox/pkg/log"
	"os"
	"syscall"
)

func RunContainerInitProcess(cmd string, args []string) error {
	log.Info(fmt.Sprintf("RunContainerInitProcess, cmd is: %s", cmd))
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	if err := syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), ""); err != nil {
		log.Error(err.Error())
		return err
	}
	argv := []string{cmd}
	if err := syscall.Exec(cmd, argv, os.Environ()); err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}
