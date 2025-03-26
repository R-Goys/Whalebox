package cgroups

import (
	"fmt"

	subsystem "github.com/R-Goys/Whalebox/cgroups/subsystems"
	"github.com/R-Goys/Whalebox/pkg/log"
)

type CgroupManager struct {
	Resource   *subsystem.ResourceConfig
	CgroupPath string
}

func NewCgroupManager(cgroupPath string) *CgroupManager {
	return &CgroupManager{
		CgroupPath: cgroupPath,
	}
}

func (c *CgroupManager) Apply(pid int) error {
	for _, subsystem := range subsystem.SubSystemIns {
		if err := subsystem.Apply(c.CgroupPath, pid); err != nil {
			log.Error("Apply cgroup failed: " + err.Error())
			return err
		}
	}
	log.Info(fmt.Sprintf("cgroup %s applied", c.CgroupPath))
	return nil
}

func (c *CgroupManager) Set(resources *subsystem.ResourceConfig) error {
	c.Resource = resources
	for _, subsystem := range subsystem.SubSystemIns {
		if err := subsystem.Set(c.CgroupPath, resources); err != nil {
			log.Error("Set cgroup resource failed: " + err.Error())
			return err
		}
	}
	log.Info(fmt.Sprintf("cgroup %s resource set", c.CgroupPath))
	return nil
}

func (c *CgroupManager) Destroy() error {
	for _, subsystem := range subsystem.SubSystemIns {
		if err := subsystem.Remove(c.CgroupPath); err != nil {
			log.Error("Remove cgroup failed: " + err.Error())
			return err
		}
	}
	log.Info(fmt.Sprintf("cgroup %s destroyed", c.CgroupPath))
	return nil
}
