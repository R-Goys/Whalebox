package cgroups

//关于资源限制的详细实现

type MemorySubSystem struct{}

// Apply implements Subsystem.
func (m *MemorySubSystem) Apply(pid int) error {

}

// Name implements Subsystem.
func (m *MemorySubSystem) Name() string {
	return "memory"
}

// Remove implements Subsystem.
func (m *MemorySubSystem) Remove(path string) error {
	panic("unimplemented")
}

// Set implements Subsystem.
func (m *MemorySubSystem) Set(path string, resources *ResourceConfig) error {
	panic("unimplemented")
}

var _ Subsystem = (*MemorySubSystem)(nil)

type CpuSubSystem struct{}

// Apply implements Subsystem.
func (c *CpuSubSystem) Apply(pid int) error {
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
func (c *CpusetSubSystem) Apply(pid int) error {
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
