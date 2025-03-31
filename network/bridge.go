package network

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/R-Goys/Whalebox/container"
	"github.com/R-Goys/Whalebox/pkg/log"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

type BridgeNetworkDriver struct {
}

// Delete implements NetworkDriver.
func (d *BridgeNetworkDriver) Delete(network Network) error {
	bridgeName := network.Name

	br, err := netlink.LinkByName(bridgeName)
	if err != nil {
		log.Error("Failed to get bridge: " + bridgeName + " error: " + err.Error())
		return err
	}
	return netlink.LinkDel(br)
}

// Disconnect implements NetworkDriver.
func (d *BridgeNetworkDriver) Disconnect(*Network, *Endpoint) error {
	return nil
}

func (d *BridgeNetworkDriver) Name() string {
	return "bridge"
}

var _ NetworkDriver = &BridgeNetworkDriver{}

func (d *BridgeNetworkDriver) Create(subnet, name string) (*Network, error) {
	ip, iprange, _ := net.ParseCIDR(subnet)
	iprange.IP = ip
	n := &Network{
		Name:    name,
		IpRange: iprange,
	}

	err := d.initBridge(n)
	if err != nil {
		log.Error("Failed to create bridge network" + err.Error())
		return nil, err
	}
	return n, nil
}

func (d *BridgeNetworkDriver) initBridge(n *Network) error {
	bridgeName := n.Name
	if err := createBridgeInterface(bridgeName); err != nil {
		log.Error("Failed to create bridge" + err.Error())
		return err
	}
	gatewayIP := *n.IpRange
	gatewayIP.IP = n.IpRange.IP
	if err := setInterfaceIP(bridgeName, gatewayIP.String()); err != nil {
		log.Error("Failed to assign Address: " + gatewayIP.String() + " to bridge: " + bridgeName + " error: " + err.Error())
		return err
	}
	if err := setInterfaceUp(bridgeName); err != nil {
		log.Error("Failed to set bridge up: " + bridgeName + " error: " + err.Error())
		return err
	}
	if err := setupIPTables(bridgeName, n.IpRange); err != nil {
		log.Error("Failed to setup IPTables for bridge: " + bridgeName + " error: " + err.Error())
		return err
	}
	return nil
}

func setInterfaceIP(name string, rawip string) error {
	iface, err := netlink.LinkByName(name)
	if err != nil {
		log.Error("Failed to get interface: " + name + " error: " + err.Error())
		return err
	}
	ipNet, err := netlink.ParseIPNet(rawip)
	if err != nil {
		log.Error("Failed to parse IPNet: " + rawip + " error: " + err.Error())
		return err
	}
	addr := &netlink.Addr{IPNet: ipNet, Peer: ipNet, Label: "", Flags: 0, Scope: 0, Broadcast: nil}
	return netlink.AddrAdd(iface, addr)

}

func setInterfaceUp(interfaceName string) error {
	iface, err := netlink.LinkByName(interfaceName)
	if err != nil {
		log.Error("Failed to get interface: " + interfaceName + " error: " + err.Error())
		return err
	}

	//通过linksetup设置接口状态为up
	if err := netlink.LinkSetUp(iface); err != nil {
		log.Error("Failed to set interface up: " + interfaceName + " error: " + err.Error())
		return err
	}
	return nil
}

func setupIPTables(bridgeName string, iprange *net.IPNet) error {
	iptablecmd := fmt.Sprintf("-t nat -A POSTROUTING -s %s ! -o %s -j MASQUERADE", iprange.String(), bridgeName)
	cmd := exec.Command("iptables", strings.Split(iptablecmd, " ")...)
	output, err := cmd.Output()
	if err != nil {
		log.Error("Failed to setup IPTables for bridge: " + bridgeName + " error: " + err.Error() + " output: " + string(output))
		return err
	}
	log.Info("IPTables setup for bridge: " + bridgeName + " output: " + string(output))
	return nil
}

func createBridgeInterface(bridgeName string) error {
	_, err := net.InterfaceByName(bridgeName)
	if err == nil || !strings.Contains(err.Error(), "no such network interface") {
		log.Error("Failed to create bridge interface: " + bridgeName + " error: " + err.Error())
		return err
	}

	la := netlink.NewLinkAttrs()
	la.Name = bridgeName

	br := &netlink.Bridge{LinkAttrs: la}
	if err := netlink.LinkAdd(br); err != nil {
		log.Error("Failed to create bridge interface: " + bridgeName + " error: " + err.Error())
		return fmt.Errorf("Bridge creation failed for bridge %s: %v", bridgeName, err)
	}
	return nil
}

func (d *BridgeNetworkDriver) Connect(n *Network, endpoint *Endpoint) error {
	bridgeName := n.Name
	br, err := netlink.LinkByName(bridgeName)
	if err != nil {
		log.Error("Failed to get bridge: " + bridgeName + " error: " + err.Error())
		return err
	}

	la := netlink.NewLinkAttrs()
	la.Name = endpoint.ID[:5]
	la.MasterIndex = br.Attrs().Index

	endpoint.Device = netlink.Veth{
		LinkAttrs: la,
		PeerName:  "cif-" + endpoint.ID[:5],
	}

	if err := netlink.LinkAdd(&endpoint.Device); err != nil {
		log.Error("Failed to add endpoint device: " + endpoint.ID + " error: " + err.Error())
		return err
	}
	if err := netlink.LinkSetUp(&endpoint.Device); err != nil {
		log.Error("Failed to set endpoint device up: " + endpoint.ID + " error: " + err.Error())
		return err
	}
	return nil
}

func configEndpointIpAddressAndRoute(endpoint *Endpoint, cinfo *container.Container) error {
	peerLink, err := netlink.LinkByName(endpoint.Device.PeerName)
	if err != nil {
		log.Error("Failed to get peer link: " + endpoint.Device.PeerName + " error: " + err.Error())
		return err
	}

	defer enterContainerNamespace(&peerLink, cinfo)()

	interfaceIP := endpoint.Network.IpRange
	interfaceIP.IP = endpoint.IPAddress

	if err := setInterfaceIP(endpoint.Device.Name, interfaceIP.String()); err != nil {
		log.Error("Failed to assign Address: " + interfaceIP.String() + " to endpoint: " + endpoint.ID + " error: " + err.Error())
		return err
	}

	if err = setInterfaceUp(endpoint.Device.PeerName); err != nil {
		log.Error("Failed to set endpoint up: " + endpoint.ID + " error: " + err.Error())
		return err
	}

	if err = setInterfaceUp("lo"); err != nil {
		log.Error("Failed to set lo up: " + err.Error())
		return err
	}
	_, cidr, _ := net.ParseCIDR("0.0.0.0/0")
	defaultRoute := &netlink.Route{
		LinkIndex: peerLink.Attrs().Index,
		Dst:       cidr,
		Scope:     netlink.SCOPE_LINK,
		Gw:        interfaceIP.IP,
	}
	if err := netlink.RouteAdd(defaultRoute); err != nil {
		log.Error("Failed to add default route to endpoint: " + endpoint.ID + " error: " + err.Error())
		return err
	}
	return nil
}

func enterContainerNamespace(link *netlink.Link, cinfo *container.Container) func() {
	f, err := os.OpenFile(fmt.Sprintf("proc/%d/ns/net", cinfo.Pid), os.O_RDONLY, 0)
	if err != nil {
		log.Error("Failed to open container netns: " + err.Error())
		return func() {}
	}
	nsFD := f.Fd()

	runtime.LockOSThread()

	if err := netlink.LinkSetNsFd(*link, int(nsFD)); err != nil {
		log.Error("Failed to set link netns: " + err.Error())
		return func() {}
	}
	//获取namespace，方便关闭
	orign, err := netns.Get()
	if err != nil {
		log.Error("Failed to get current netns: " + err.Error())
		return func() {}
	}

	if err = netns.Set(netns.NsHandle(nsFD)); err != nil {
		log.Error("Failed to set netns: " + err.Error())
		return func() {}
	}

	return func() {
		netns.Set(orign)
		orign.Close()
		runtime.UnlockOSThread()
		f.Close()
	}
}

func configPortMapping(endpoint *Endpoint, cinfo *container.Container) error {
	for _, mapping := range endpoint.PortMapping {
		port := strings.Split(mapping, ":")
		if len(port) != 2 {
			log.Error("Invalid port mapping: " + mapping)
			continue
		}
		iptablescmd := fmt.Sprintf("-t nat -A PREROUTING -p tcp -m tcp --dport %s -j DNAT --to-destination %s:%s", port[0], endpoint.IPAddress.String(), port[1])
		cmd := exec.Command("iptables", strings.Split(iptablescmd, " ")...)
		output, err := cmd.Output()
		if err != nil {
			log.Error("Failed to setup port mapping for endpoint: " + endpoint.ID + " error: " + err.Error() + " output: " + string(output))
			continue
		}
		log.Info("Port mapping setup for endpoint: " + endpoint.ID + " output: " + string(output))
	}
	return nil
}
