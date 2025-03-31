# Whalebox

Whalebox 是一个模拟实现Docker的容器，本项目参考了《动手写Docker》/《Docker容器与容器云》两本书。

## 介绍

Whalebox的主要功能：

- 容器：支持创建、启动、停止、删除容器，支持容器的后台启动，资源限制、网络(非完全体)、存储等功能。
- 镜像：支持镜像的创建、导入、导出等操作，支持镜像的分层存储、共享等功能。
- 网络：支持容器间的网络通信，支持容器的网络隔离等功能，但是不支持容器与外部通信
- 存储：支持容器的存储卷管理，支持容器数据的持久化存储。
- 运行时：支持容器的运行时管理，支持容器的监控、日志等管理。

另外，与原书不同的是，Whalebox实现基于Cgroup v2的资源限制功能，并使用了OverlayFS作为容器的存储层。
Cgroup v2与原书的主要不同之处在于，它允许容器独自拥有一组资源限制，
换句话说，Cgroup v2是一个树状结构，这种结构使得容器管理更加灵活，方便。
而OverlayFS也是一种存储技术，表面看区别并不大，主要在于文件结构和实现有一些不同。
更多详情可以看我的Note，关于本项目的实现全过程。
[Whalebox](https://github.com/R-Goys/Notes/blob/main/Note/%E6%89%8B%E6%90%93/%E6%89%8B%E5%86%99Docker.md)

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
│   ├── cmd_func.go
│   ├── commit.go
│   ├── exec.go
│   ├── list.go
│   ├── logs.go
│   ├── main_command.go
│   ├── main.go
│   ├── remove.go
│   ├── run.go
│   └── stop.go
├── common
│   └── path.go
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
│   ├── example3
│   │   ├── busybox.tar
│   │   ├── image123.tar
│   │   ├── image.tar
│   │   ├── mnt
│   │   ├── workDir
│   │   └── writeLayer
│   ├── example4
│   ├── example5
│   │   └── network_test.go
│   └── volume
├── Godeps
│   └── Godeps.json
├── go.mod
├── go.sum
├── network
│   ├── bridge.go
│   ├── ipam
│   ├── ipam.go
│   ├── network
│   └── network.go
├── nsenter
│   └── nsenter.go
├── pkg
│   └── log
│       ├── logger.go
│       └── record.log
└── README.md

```