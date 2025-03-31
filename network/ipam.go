package network

import (
	"encoding/json"
	"fmt"
	"strings"

	"net"
	"os"
	"path"

	"github.com/R-Goys/Whalebox/pkg/log"
)

const ipamDefaultAllocatorPath = "/home/rinai/PROJECTS/Whalebox/network/ipam/subnet.json"

// IPAM是一个存储IP地址分配信息的结构体
type IPAM struct {
	//存放IP地址分配信息的文件路径
	SubnetAllocatorPath string
	Subnets             *map[string]string
}

var IpAllocator = &IPAM{
	SubnetAllocatorPath: ipamDefaultAllocatorPath,
}

func (ipam *IPAM) load() error {
	if _, err := os.Stat(ipam.SubnetAllocatorPath); err != nil {
		if os.IsNotExist(err) {
			_, err = os.Create(ipam.SubnetAllocatorPath)
			if err != nil {
				log.Error("Failed to create subnet allocator file: " + err.Error())
				return err
			}
		} else {
			log.Error("Other error when checking subnet allocator file directory: " + err.Error())
			return err
		}
	}
	subnetConfigFile, err := os.Open(ipam.SubnetAllocatorPath)
	if err != nil {
		log.Error("Failed to open subnet allocator file for reading: " + err.Error())
		return err
	}
	defer subnetConfigFile.Close()

	subnetJson := make([]byte, 2000)
	n, err := subnetConfigFile.Read(subnetJson)
	if err != nil {
		return err
	}
	err = json.Unmarshal(subnetJson[:n], ipam.Subnets)
	if err != nil {
		log.Error("Failed to unmarshal subnets from json: " + err.Error())
		return err
	}
	log.Debug(fmt.Sprintf("Loaded subnets: %v", *ipam.Subnets))
	return nil
}

// 保存IP地址分配信息到文件
func (ipam *IPAM) dump() error {
	ipamConfigFileDir, _ := path.Split(ipam.SubnetAllocatorPath)
	if _, err := os.Stat(ipamConfigFileDir); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(ipamConfigFileDir, 0755)
			if err != nil {
				log.Error("Failed to mkdir for subnet allocator file: " + err.Error())
				return err
			}
		} else {
			log.Error("Other error when checking subnet allocator file directory: " + err.Error())
			return err
		}
	}
	subnetConfigFile, err := os.OpenFile(ipam.SubnetAllocatorPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Error("Failed to open subnet allocator file for writing: " + err.Error())
		return err
	}
	defer subnetConfigFile.Close()

	ipamJson, err := json.Marshal(ipam.Subnets)
	if err != nil {
		log.Error("Failed to marshal subnets to json: " + err.Error())
		return err
	}
	_, err = subnetConfigFile.Write(ipamJson)
	if err != nil {
		log.Error("Failed to write subnets to file: " + err.Error())
		return err
	}
	return nil
}

// 分配IP地址
func (ipam *IPAM) Allocate(subnet *net.IPNet) (ip net.IP, err error) {
	ipam.Subnets = &map[string]string{}
	if err := ipam.load(); err != nil {
		log.Error("Failed to load subnet allocator file: " + err.Error())
	}
	_, subnet, _ = net.ParseCIDR(subnet.String())
	one, size := subnet.Mask.Size()
	if _, exist := (*ipam.Subnets)[subnet.String()]; !exist {
		(*ipam.Subnets)[subnet.String()] = strings.Repeat("0", 1<<uint8(size-one))
	}
	//遍历子网的每个IP地址
	for c := range (*ipam.Subnets)[subnet.String()] {
		//c并不表示ip地址本身，而是一个索引，如果为0，表示该IP地址未分配，可以分配
		//此处使用位图的方式来表示每个IP地址的分配情况。
		if (*ipam.Subnets)[subnet.String()][c] == '0' {
			ipalloc := []byte((*ipam.Subnets)[subnet.String()])
			//为1，表示该IP地址已分配。
			log.Debug(fmt.Sprintf("Allocating IP offset: %v", c))
			ipalloc[c] = '1'
			(*ipam.Subnets)[subnet.String()] = string(ipalloc)
			ip = subnet.IP
			//通过偏移量计算出IP地址
			for t := uint(4); t > 0; t-- {
				[]byte(ip)[4-t] += uint8(c >> ((t - 1) * 8))
			}
			ip[3] += 1
			break
		}
	}
	//保存IP地址分配信息到文件
	ipam.dump()
	return
}

// 释放IP地址
func (ipam *IPAM) Release(subnet *net.IPNet, ipaddr *net.IP) error {
	ipam.Subnets = &map[string]string{}

	_, subnet, _ = net.ParseCIDR(subnet.String())

	err := ipam.load()
	if err != nil {
		log.Error("Failed to load subnet allocator file: " + err.Error())
	}
	c := 0
	//四个字节表示方式
	releaseIP := ipaddr.To4()
	releaseIP[3]--
	//通过偏移量计算出索引
	for t := uint(4); t > 0; t-- {
		c += int(releaseIP[t-1]-subnet.IP[t-1]) << ((4 - t) * 8)
	}
	log.Debug(fmt.Sprintf("map: %v", (*ipam.Subnets)[subnet.String()]))
	log.Debug(fmt.Sprintf("offset: %d", c))
	ipalloc := []byte((*ipam.Subnets)[subnet.String()])
	ipalloc[c] = '0'
	(*ipam.Subnets)[subnet.String()] = string(ipalloc)
	//保存IP地址分配信息到文件
	ipam.dump()
	return nil
}
