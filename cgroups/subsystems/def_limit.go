package subsystem

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/R-Goys/Whalebox/pkg/log"
)

//关于资源限制的详细实现

/*
+---------------------------------------------------+
|													|
|													|
|													|
|			Memory Subsystem Implement				|
|													|
|													|
|													|
+---------------------------------------------------+
*/

type MemorySubSystem struct{}

// Apply implements Subsystem, 此处将进程添加到指定的cgroup
func (m *MemorySubSystem) Apply(CgroupPath string, pid int) error {
	if subsysCgroupPath, err := GetCgroupPath(m.Name(), CgroupPath, true); err == nil {
		if err := os.WriteFile(path.Join(subsysCgroupPath, "cgroup.procs"), []byte(strconv.Itoa(pid)), 0644); err != nil {
			log.Error("Memory:" + "failed to add process to cgroup: %v" + err.Error())
			return fmt.Errorf("failed to add process to cgroup: %v", err)
		}
		return nil
	} else {
		log.Error("Memory:" + "failed to get cgroup path: %v" + err.Error())
		return fmt.Errorf("failed to get cgroup path: %v", err)
	}
}

// Name implements Subsystem.
func (m *MemorySubSystem) Name() string {
	return "memory"
}

// Remove implements Subsystem, 此处将指定的cgroup删除
func (m *MemorySubSystem) Remove(CgroupPath string) error {
	if subsysGroupPath, err := GetCgroupPath(m.Name(), CgroupPath, false); err == nil {
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
		log.Debug("memory limit not set")
		return nil
	} else {
		return fmt.Errorf("failed to get cgroup path: %v", err)
	}
}

var _ Subsystem = (*MemorySubSystem)(nil)

/*
+---------------------------------------------------+
|													|
|													|
|													|
|			Cpusub Subsystem Implement				|
|													|
|													|
|													|
+---------------------------------------------------+
*/

type CpuSubSystem struct{}

// Apply implements Subsystem.
func (c *CpuSubSystem) Apply(CgroupPath string, pid int) error {
	if subsysCgroupPath, err := GetCgroupPath(c.Name(), CgroupPath, true); err == nil {
		if err := os.WriteFile(path.Join(subsysCgroupPath, "cgroup.procs"), []byte(strconv.Itoa(pid)), 0644); err != nil {
			log.Error("Cpusub:" + "failed to add process to cgroup: %v" + err.Error())
			return fmt.Errorf("failed to add process to cgroup: %v", err)
		}
		return nil
	} else {
		log.Error("Cpusub:" + "failed to get cgroup path: %v" + err.Error())
		return fmt.Errorf("failed to get cgroup path: %v", err)
	}
}

// Name implements Subsystem.
func (c *CpuSubSystem) Name() string {
	return "cpu"
}

// Remove implements Subsystem.
func (c *CpuSubSystem) Remove(CgroupPath string) error {
	if subsysGroupPath, err := GetCgroupPath(c.Name(), CgroupPath, false); err == nil {
		if err := os.Remove(subsysGroupPath); err != nil {
			log.Error("Cpusub:" + "failed to remove cgroup path %s: %v" + subsysGroupPath + err.Error())
			return fmt.Errorf("failed to remove cgroup path %s: %v", subsysGroupPath, err)
		}
		return nil
	} else {
		log.Error("Cpusub:" + "failed to get cgroup path: %v" + err.Error())
		return fmt.Errorf("failed to get cgroup path: %v", err)
	}
}

// Set implements Subsystem.
func (c *CpuSubSystem) Set(CgroupPath string, resources *ResourceConfig) error {
	if subsysCgroupPath, err := GetCgroupPath(c.Name(), CgroupPath, true); err == nil {
		if resources.CpuShares != "" {
			if err := os.WriteFile(path.Join(subsysCgroupPath, "cpu.shares"), []byte(resources.CpuShares), 0644); err != nil {
				log.Error("Cpusub:" + "failed to set cpu shares: %v" + err.Error())
				return fmt.Errorf("failed to set cpu shares: %v", err)
			}
			return nil
		}
		log.Debug("cpu shares not set")
		return nil
	} else {
		log.Error("Cpusub:" + "failed to get cgroup path: %v" + err.Error())
		return fmt.Errorf("failed to get cgroup path: %v", err)
	}
}

var _ Subsystem = (*CpuSubSystem)(nil)

/*
+---------------------------------------------------+
|													|
|													|
|													|
|			CpusetSub Subsystem Implement			|
|													|
|													|
|													|
+---------------------------------------------------+
*/

type CpusetSubSystem struct{}

// Apply implements Subsystem.
func (c *CpusetSubSystem) Apply(CgroupPath string, pid int) error {
	if subsysCgroupPath, err := GetCgroupPath(c.Name(), CgroupPath, true); err == nil {
		if err := os.WriteFile(path.Join(subsysCgroupPath, "cgroup.procs"), []byte(strconv.Itoa(pid)), 0644); err != nil {
			log.Error("CpusetSub:" + "failed to add process to cgroup: %v" + err.Error())
			return fmt.Errorf("failed to add process to cgroup: %v", err)
		}
		return nil
	} else {
		log.Error("CpusetSub:" + "failed to get cgroup path: %v" + err.Error())
		return fmt.Errorf("failed to get cgroup path: %v", err)
	}
}

// Name implements Subsystem.
func (c *CpusetSubSystem) Name() string {
	return "cpuset"
}

// Remove implements Subsystem.
func (c *CpusetSubSystem) Remove(CgroupPath string) error {
	if subsysGroupPath, err := GetCgroupPath(c.Name(), CgroupPath, false); err == nil {
		if err := os.Remove(subsysGroupPath); err != nil {
			log.Error("CpusetSub:" + "failed to remove cgroup path %s: %v" + subsysGroupPath + err.Error())
			return fmt.Errorf("failed to remove cgroup path %s: %v", subsysGroupPath, err)
		}
		return nil
	} else {
		log.Error("CpusetSub:" + "failed to get cgroup path: %v" + err.Error())
		return fmt.Errorf("failed to get cgroup path: %v", err)
	}
}

// Set implements Subsystem.
func (c *CpusetSubSystem) Set(CgroupPath string, resources *ResourceConfig) error {
	if subsysCgroupPath, err := GetCgroupPath(c.Name(), CgroupPath, true); err == nil {
		if resources.CpuSet != "" {
			if err := os.WriteFile(path.Join(subsysCgroupPath, "cpuset.cpus"), []byte(resources.CpuSet), 0644); err != nil {
				log.Error("CpusetSub:" + "failed to set cpuset: %v" + err.Error())
				return fmt.Errorf("failed to set cpuset: %v", err)
			}
			return nil
		}
		log.Debug("cpuset not set")
		return nil
	} else {
		log.Error("CpusetSub:" + "failed to get cgroup path: %v" + err.Error())
		return fmt.Errorf("failed to get cgroup path: %v", err)
	}
}

var _ Subsystem = (*CpusetSubSystem)(nil)
