package cgroup

type ResourceConfig struct {
	MemoryLimit string //内存限制
	CpuShares   string
	CpuSet      string
}

type CgroupInterface interface {
	Path() string
	//为该cgroup设置资源限制
	Set(resources *ResourceConfig) error
	//删除该Cgroup
	Remove() error
}
