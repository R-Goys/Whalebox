package cgroup

import (
	"fmt"
	"os"
	"path"

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
type Cgroup struct {
	path string
}

// 新建一个cgroup，由于whalebox的子cgroup需要从whalebox继承内存管理/cpu管理等，所以需要手动将继承选项添加进去。
func NewCgroup(Root string, pid string) *Cgroup {
	Path := Root + "/" + pid
	if CgroupPath, err := GetCgroupPath(Path, true); err == nil {
		if err := os.WriteFile(path.Join(CgroupPath, "cgroup.procs"), []byte(pid), 0644); err != nil {
			log.Error(fmt.Sprintf("failed to add process to cgroup: %v", err))
			return nil
		}

		if err := os.WriteFile(path.Join(CgroupPath[:len(CgroupPath)-len(pid)-1], "cgroup.subtree_control"), []byte("+memory +cpuset +cpu"), 0644); err != nil {
			log.Error(fmt.Sprintf("failed to set cgroup.subtree_control: %v", err))
			return nil
		}
	}
	return &Cgroup{
		path: Path,
	}
}

func (s *Cgroup) Path() string {
	return s.path
}

func (s *Cgroup) Set(resources *ResourceConfig) error {
	if err := s.SetMemoryLimit(resources); err != nil {
		return err
	}
	if err := s.SetCpuShares(resources); err != nil {
		return err
	}
	if err := s.SetCpuLimit(resources); err != nil {
		return err
	}
	return nil
}

// Remove implements CgroupInterface.
func (s *Cgroup) Remove() error {
	if err := os.RemoveAll(s.Path()); err != nil {
		log.Error(fmt.Sprintf("failed to remove cgroup %s: %v", s.Path(), err))
		return fmt.Errorf("failed to remove cgroup %s: %v", s.Path(), err)
	}
	log.Info(fmt.Sprintf("cgroup %s removed", s.Path()))
	return nil
}

var _ CgroupInterface = (*Cgroup)(nil)

// Set implements Subsystem, 此处为cgroup设置资源限制，也就是内存的限制
func (s *Cgroup) SetMemoryLimit(resources *ResourceConfig) error {
	if SubSystemPath, err := GetCgroupPath(s.Path(), true); err == nil {
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

// Set implements Subsystem.
func (s *Cgroup) SetCpuShares(resources *ResourceConfig) error {
	if subsysCgroupPath, err := GetCgroupPath(s.Path(), true); err == nil {
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

// Set implements Subsystem.
func (s *Cgroup) SetCpuLimit(resources *ResourceConfig) error {
	if subsysCgroupPath, err := GetCgroupPath(s.Path(), true); err == nil {
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
