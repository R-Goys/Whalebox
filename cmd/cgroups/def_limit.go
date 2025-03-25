package cgroups

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/R-Goys/Whalebox/pkg/log"
)

//关于资源限制的详细实现

type MemorySubSystem struct{}

// Apply implements Subsystem, 此处将进程添加到指定的cgroup
func (m *MemorySubSystem) Apply(CgroupPath string, pid int) error {
	if subsysCgroupPath, err := GetCgroupPath(m.Name(), CgroupPath, true); err == nil {
		if err := os.WriteFile(path.Join(subsysCgroupPath, "cgroup.procs"), []byte(strconv.Itoa(pid)), 0644); err != nil {
			log.Error("failed to add process to cgroup: %v" + err.Error())
			return fmt.Errorf("failed to add process to cgroup: %v", err)
		}
		return nil
	} else {
		log.Error("failed to get cgroup path: %v" + err.Error())
		return fmt.Errorf("failed to get cgroup path: %v", err)
	}
}

// Name implements Subsystem.
func (m *MemorySubSystem) Name() string {
	return "memory"
}

// Remove implements Subsystem, 此处将指定的cgroup删除
func (m *MemorySubSystem) Remove(path string) error {
	if subsysGroupPath, err := GetCgroupPath(m.Name(), path, false); err == nil {
		if err := os.Remove(subsysGroupPath); err != nil {
			return fmt.Errorf("failed to remove cgroup path %s: %v", subsysGroupPath, err)
		}
		return nil
	} else {
		return fmt.Errorf("failed to get cgroup path: %v", err)
	}
}

// Set implements Subsystem, 此处为cgroup设置资源限制，也就是内存的限制
func (s *MemorySubSystem) Set(cgroupPath string, resources *ResourceConfig) error {
	if SubSystemPath, err := GetCgroupPath(s.Name(), cgroupPath, true); err == nil {
		//设置内存限制
		if resources.MemoryLimit != "" {
			if err := os.WriteFile(path.Join(SubSystemPath, "memory.max"), []byte(resources.MemoryLimit), 0644); err != nil {
				return fmt.Errorf("failed to set memory limit: %v", err)
			}
			return nil
		}
		return fmt.Errorf("failed to set memory limit: memory limit is empty")
	} else {
		return fmt.Errorf("failed to get cgroup path: %v", err)
	}
}

var _ Subsystem = (*MemorySubSystem)(nil)

type CpuSubSystem struct{}

// Apply implements Subsystem.
func (c *CpuSubSystem) Apply(CgroupPath string, pid int) error {
	panic("unimplemented")
}

// Name implements Subsystem.
func (c *CpuSubSystem) Name() string {
	return "cpu"
}

// Remove implements Subsystem.
func (c *CpuSubSystem) Remove(path string) error {
	panic("unimplemented")
}

// Set implements Subsystem.
func (c *CpuSubSystem) Set(path string, resources *ResourceConfig) error {
	panic("unimplemented")
}

var _ Subsystem = (*CpuSubSystem)(nil)

type CpusetSubSystem struct{}

// Apply implements Subsystem.
func (c *CpusetSubSystem) Apply(CgroupPath string, pid int) error {
	panic("unimplemented")
}

// Name implements Subsystem.
func (c *CpusetSubSystem) Name() string {
	return "cpuset"
}

// Remove implements Subsystem.
func (c *CpusetSubSystem) Remove(path string) error {
	panic("unimplemented")
}

// Set implements Subsystem.
func (c *CpusetSubSystem) Set(path string, resources *ResourceConfig) error {
	panic("unimplemented")
}

var _ Subsystem = (*CpusetSubSystem)(nil)
