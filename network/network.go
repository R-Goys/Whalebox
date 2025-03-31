package network

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path"
	"path/filepath"
	"text/tabwriter"

	"github.com/R-Goys/Whalebox/container"
	"github.com/R-Goys/Whalebox/pkg/log"
	"github.com/vishvananda/netlink"
)

var (
	defaultNetworkPath = "/home/rinai/PROJECTS/Whalebox/network/network/"
	drivers            = map[string]NetworkDriver{}
	networks           = map[string]*Network{}
)

type Network struct {
	Name    string
	IpRange *net.IPNet
	Driver  string
}

type Endpoint struct {
	ID          string           `json:"id"`
	Device      netlink.Veth     `json:"device"`
	IPAddress   net.IP           `json:"ipAddress"`
	MacAdress   net.HardwareAddr `json:"macAdress"`
	PortMapping []string         `json:"portMapping"`
	Network     *Network         `json:"network"`
}

type NetworkDriver interface {
	Name() string
	Create(subnet string, name string) (*Network, error)
	Delete(network Network) error
	Connect(*Network, *Endpoint) error
	Disconnect(*Network, *Endpoint) error
}

func CreateNetwork(driver, subnet, name string) error {
	_, cidr, _ := net.ParseCIDR(subnet)
	gatewayIp, err := IpAllocator.Allocate(cidr)
	if err != nil {
		log.Error("Failed to allocate IP address for network " + name)
		return err
	}
	cidr.IP = gatewayIp

	nw, err := drivers[driver].Create(cidr.String(), name)
	if err != nil {
		log.Error("Failed to create network " + name)
		return err
	}
	return nw.dump(defaultNetworkPath)
}

func (nw *Network) dump(dumpPath string) error {
	if _, err := os.Stat(dumpPath); err != nil {
		if os.IsNotExist(err) {
			log.Debug("Creating network directory " + dumpPath)
			os.MkdirAll(dumpPath, 0755)
		} else {
			log.Error("Failed to check network directory " + dumpPath)
			return err
		}
	}
	nwPath := path.Join(dumpPath, nw.Name)
	nwFile, err := os.OpenFile(nwPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Error("Failed to open network file " + nwPath)
		return err
	}
	defer nwFile.Close()

	nwJson, err := json.Marshal(nw)
	if err != nil {
		log.Error("Failed to marshal network " + nw.Name)
		return err
	}
	_, err = nwFile.Write(nwJson)
	if err != nil {
		log.Error("Failed to write network file " + nwPath)
		return err
	}
	return nil
}

func (nw *Network) load(dumpPath string) error {
	nwConfigFile, err := os.Open(dumpPath)
	if err != nil {
		log.Error("Failed to open network file " + dumpPath)
		return err
	}
	defer nwConfigFile.Close()
	nwJson := make([]byte, 1024)
	n, err := nwConfigFile.Read(nwJson)
	if err != nil {
		log.Error("Failed to read network file " + dumpPath)
		return err
	}
	err = json.Unmarshal(nwJson[:n], nw)
	if err != nil {
		log.Error("Failed to unmarshal network file " + dumpPath)
		return err
	}
	return nil
}

func Connect(networkName string, cInfo *container.Container) error {
	network, ok := networks[networkName]
	if !ok {
		log.Error("Network " + networkName + " not found")
		return fmt.Errorf("Network %s not found", networkName)
	}
	ip, err := IpAllocator.Allocate(network.IpRange)
	if err != nil {
		log.Error("Failed to allocate IP address for container " + cInfo.Name)
		return err
	}
	endpoint := Endpoint{
		ID:          fmt.Sprintf("%s-%s", cInfo.Id, cInfo.Name),
		IPAddress:   ip,
		Network:     network,
		PortMapping: cInfo.PortMapping,
	}

	if err = drivers[network.Driver].Connect(network, &endpoint); err != nil {
		log.Error("Failed to connect container " + cInfo.Name + " to network " + networkName)
		return err
	}
	if err = configEndpointIpAddressAndRoute(&endpoint, cInfo); err != nil {
		log.Error("Failed to configure container " + cInfo.Name + " IP address and route")
		return err
	}
	return configPortMapping(&endpoint, cInfo)
}

func Init() error {
	var bridgeDriver = &BridgeNetworkDriver{}
	drivers[bridgeDriver.Name()] = bridgeDriver
	if _, err := os.Stat(defaultNetworkPath); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(defaultNetworkPath, 0644)
		} else {
			log.Error("Failed to check network directory " + defaultNetworkPath)
			return err
		}
	}
	filepath.Walk(defaultNetworkPath, func(nwPath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		_, nwName := path.Split(nwPath)
		nw := &Network{
			Name: nwName,
		}

		if err = nw.load(nwPath); err != nil {
			log.Error("Failed to load network " + nwName)
			return err
		}
		networks[nw.Name] = nw
		return nil
	})
	return nil
}

func ListNetworks() {
	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	fmt.Fprintf(w, "Name\tIpRange\tDriver\n")

	for _, nw := range networks {
		fmt.Fprintf(w, "%s\t%s\t%s\n", nw.Name, nw.IpRange.String(), nw.Driver)
	}
	if err := w.Flush(); err != nil {
		log.Error("Failed to flush network list")
	}
}

// 删除网络
func DeleteNetwork(networkName string) error {
	network, ok := networks[networkName]
	if !ok {
		log.Error("Network " + networkName + " not found")
		return fmt.Errorf("Network %s not found", networkName)
	}
	if err := IpAllocator.Release(network.IpRange, &network.IpRange.IP); err != nil {
		log.Error("Failed to release IP address for network " + networkName)
		return err
	}
	if err := drivers[network.Driver].Delete(*network); err != nil {
		log.Error("Failed to delete network " + networkName)
		return err
	}
	return network.remove(defaultNetworkPath)
}

func (nw *Network) remove(dumpPath string) error {
	nwPath := path.Join(dumpPath, nw.Name)
	if _, err := os.Stat(nwPath); err != nil {
		if os.IsNotExist(err) {
			log.Debug("Network file " + nwPath + " not found")
			return nil
		} else {
			log.Error("Failed to check network file " + nwPath)
			return err
		}
	}
	if err := os.Remove(nwPath); err != nil {
		log.Error("Failed to remove network file " + nwPath)
		return err
	}
	return nil
}
