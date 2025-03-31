package test

import (
	"net"
	"testing"

	"github.com/R-Goys/Whalebox/network"
	"github.com/R-Goys/Whalebox/pkg/log"
)

func TestNetwork(t *testing.T) {
	log.InitLogger()
	_, ipnet, _ := net.ParseCIDR("192.168.1.0/24")
	ip, _ := network.IpAllocator.Allocate(ipnet)
	log.Debug(ip.String())
}

func TestNetwork2(t *testing.T) {
	log.InitLogger()
	ip, ipnet, _ := net.ParseCIDR("192.168.1.1/24")
	network.IpAllocator.Release(ipnet, &ip)
}
