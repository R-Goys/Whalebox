package cgroups

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/R-Goys/Whalebox/pkg/log"
)

func GetCgroupPath(subsystem, cgroupPath string, autoCreate bool) (string, error) {
	cgroupRoot := FindCgroupMountpoint(subsystem)
	if _, err := os.Stat(path.Join(cgroupRoot, cgroupPath)); err == nil || (autoCreate && os.IsNotExist(err)) {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(path.Join(cgroupRoot, cgroupPath), 0755); err == nil {
				return path.Join(cgroupRoot, cgroupPath), nil
			} else {
				log.Error(err.Error())
				return "", fmt.Errorf("failed to create cgroup path %s: %v", path.Join(cgroupRoot, cgroupPath), err)
			}
		} else {
			return path.Join(cgroupRoot, cgroupPath), nil
		}
	}
	log.Error(fmt.Sprintf("cgroup path %s not found", path.Join(cgroupRoot, cgroupPath)))
	return "", fmt.Errorf("cgroup path %s not found", path.Join(cgroupRoot, cgroupPath))
}

func FindCgroupMountpoint(subsystem string) string {
	//打开当前进程的挂载信息文件
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return ""
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		for _, opt := range strings.Split(fields[3], ",") {
			//如果匹配到了对应的subsystem，则说明当前的挂载点支持这个subsystem
			if opt == subsystem {
				return fields[4]
			}
		}
	}
	if err = scanner.Err(); err != nil {
		return ""
	}
	return ""
}
