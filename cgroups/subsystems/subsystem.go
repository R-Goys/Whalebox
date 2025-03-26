package subsystem

type ResourceConfig struct {
	MemoryLimit string //内存限制
	CpuShares   string
	CpuSet      string
}

type Subsystem interface {
	//返回subsystem的名称
	Name() string
	//为某个cgroup设置资源限制
	Set(path string, resources *ResourceConfig) error
	//添加进程到cgroup
	Apply(CgroupPath string, pid int) error
	//删除指定的cgroup
	Remove(path string) error
}

var (
	SubSystemIns = []Subsystem{
		&CpuSubSystem{},
		&MemorySubSystem{},
		&CpusetSubSystem{},
	}
)
