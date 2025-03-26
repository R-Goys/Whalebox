package cgroup

import (
	"fmt"
	"os"
	"path"

	"github.com/R-Goys/Whalebox/pkg/log"
)

const (
	cgroupRoot = "/sys/fs/cgroup"
)

// 通过cgroupPath获取cgroup的路径
func GetCgroupPath(cgroupPath string, autoCreate bool) (string, error) {
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
