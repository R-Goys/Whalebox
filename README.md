# Whalebox

Whalebox 是一个模拟实现Docker的容器，本项目参考了《动手写Docker》/《Docker容器与容器云》两本书。

## Getting Started

要开始使用Whalebox，请按照以下步骤操作：

1. 克隆仓库：

```bash
git clone https://github.com/rinai/Whalebox.git
```

2. 安装所需的包：

```bash
go mod tidy
``` 

4. 项目结构：

```
.
├── cgroups
│   ├── cgroup.go
│   ├── def_limit.go
│   └── utils.go
├── cmd
│   ├── cmd
│   ├── main_command.go
│   ├── main.go
│   └── run.go
├── container
│   ├── container_process.go
│   ├── init.go
│   ├── overlayfs.go
│   └── volume.go
├── example
│   ├── example1
│   │   ├── main.go
│   │   └── trace.log
│   ├── example2
│   │   ├── lab
│   │   │   └── aufs
│   │   │       ├── container-layer
│   │   │       │   └── image-layer4.txt
│   │   │       ├── image-layer1
│   │   │       │   └── image-layer1.txt
│   │   │       ├── image-layer2
│   │   │       │   └── image-layer2.txt
│   │   │       ├── image-layer3
│   │   │       │   └── image-layer3.txt
│   │   │       └── image-layer4
│   │   │           └── image-layer4.txt
│   │   ├── main.go
│   │   └── new
│   │       └── dockerfile
│   └── example3
│       └── busybox.tar
├── Godeps
│   └── Godeps.json
├── go.mod
├── go.sum
├── pkg
│   └── log
│       ├── logger.go
│       └── record.log
└── README.md
```