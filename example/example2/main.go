package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"

	"github.com/fatih/color"
)

const CgroupPath = "/sys/fs/cgroup/"

func GetAllChinldpids(pids []string) []string {
	for i := 0; i < len(pids); i++ {
		children, err := exec.Command("pgrep", "-P", pids[i]).Output()
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
				continue
			}
			color.Red("Error: %v", err)
			os.Exit(14)
		}
		pids = append(pids, strings.Fields(string(children))...)
	}
	return pids
}

func main() {
	if os.Args[0] == "/proc/self/exe" {
		fmt.Println("pid:", os.Getpid())

		cmd := exec.Command("sh", "-c", `stress --vm-bytes 2048m --vm-keep -m 1`)

		cmd.SysProcAttr = &syscall.SysProcAttr{}

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
			os.Exit(11)
		}
	}

	cmd := exec.Command("/proc/self/exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		fmt.Println(err)
		os.Exit(13)
	} else {
		pids := GetAllChinldpids([]string{fmt.Sprintf("%d", cmd.Process.Pid)})

		fmt.Println("child pid:", pids)

		fmt.Println("pids:", pids)

		os.Mkdir(path.Join(CgroupPath, "TestMemoryLimit"), 0755)

		for _, pid := range pids {
			os.WriteFile(path.Join(CgroupPath, "TestMemoryLimit", "cgroup.procs"), []byte(pid), 0644)
		}

		os.WriteFile(path.Join(CgroupPath, "TestMemoryLimit", "memory.max"), []byte("100m"), 0644)
	}
	cmd.Process.Wait()
}
