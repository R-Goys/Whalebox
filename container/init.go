package container

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/R-Goys/Whalebox/pkg/log"
)

func RunContainerInitProcess() error {
	cmdArray := readUserCommand()
	if cmdArray == nil {
		log.Debug("No command received from parent")
		return errors.New("no command received from parent")
	}
	log.Info(fmt.Sprintf("RunContainerInitProcess, cmd is: %s", cmdArray))

	SetupMount()

	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info(fmt.Sprintf("Find path: %s", path))

	if err := syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func readUserCommand() []string {
	pipe := os.NewFile(uintptr(3), "pipe")

	log.Debug("pipe Create")

	msg, err := io.ReadAll(pipe)
	if err != nil {
		log.Error("init read pipe error:" + err.Error())
		return nil
	}
	log.Info(fmt.Sprintf("Received command from parent: %s", msg))
	return strings.Split(string(msg), " ")
}

/*
+----------------------------------------------+
|  											   |
|  				  挂载点的实现					 |
|  											   |
+----------------------------------------------+
*/

func pivotRoot(root string) error {
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		log.Error("pivotRoot: Failed to bind mount root: " + err.Error())
		return err
	}
	pivotDir := filepath.Join(root, ".pivot_root")
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		log.Error("pivotRoot: Failed to create pivot directory: " + err.Error())
		return err
	}
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		log.Error("pivotRoot: Failed to pivot root: " + err.Error())
		return err
	}
	if err := syscall.Chdir("/"); err != nil {
		log.Error("pivotRoot: Failed to change directory to /: " + err.Error())
		return err
	}
	pivotDir = filepath.Join("/", ".pivot_root")
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		log.Error("pivotRoot: Failed to unmount pivot directory: " + err.Error())
		return err
	}
	return os.Remove(pivotDir)
}

func SetupMount() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Error("SetupMount: Failed to get current directory: " + err.Error())
		return
	}
	log.Info("Current directory: " + pwd)
	syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	pivotRoot(pwd)
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	if err := syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), ""); err != nil {
		log.Error("SetupMount: Failed to mount proc: " + err.Error())
		return
	}
	if err := syscall.Mount("tmpfs", "/tmp", "tmpfs", uintptr(defaultMountFlags), ""); err != nil {
		log.Error("SetupMount: Failed to mount tmpfs: " + err.Error())
		return
	}
}
