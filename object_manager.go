package ibclient

import (
	"fmt"
	"regexp"
)

type ObjectManager struct {
	connector *Connector

	localAddressSpace  string
	globalAddressSpace string

	dockerID string
}

func NewObjectManager(connector *Connector, globalAddressSpace string, localAddressSpace string, dockerID string) *ObjectManager {
	objMgr := new(ObjectManager)

	objMgr.connector = connector
	objMgr.localAddressSpace = localAddressSpace
	objMgr.globalAddressSpace = globalAddressSpace
	objMgr.dockerID = dockerID

	return objMgr
}

func (objMgr *ObjectManager) getBasicEA(cloudApiOwned Bool) EA {
	ea := make(EA)
	ea["Cloud API Owned"] = cloudApiOwned
	ea["CMP Type"] = "Docker"
	ea["Tenant ID"] = objMgr.dockerID
	return ea
}

func (objMgr *ObjectManager) CreateNetworkView(name string) (*NetworkView, error) {
	networkView := new(NetworkView)
	networkView.Name = name

	payload := make(Payload)
	payload["name"] = name
	payload["extattrs"] = objMgr.getBasicEA(false)

	ref, err := objMgr.connector.CreateObject("networkview", payload)
	networkView.Ref = ref

	return networkView, err
}

func (objMgr *ObjectManager) CreateNetwork(netview string, cidr string) (*Network, error) {
	network := new(Network)
	network.NetviewName = netview
	network.Cidr = cidr

	payload := make(Payload)
	payload["network_view"] = netview
	payload["network"] = cidr
	payload["extattrs"] = objMgr.getBasicEA(true)

	ref, err := objMgr.connector.CreateObject("network", payload)
	network.Ref = ref

	return network, err
}

func (objMgr *ObjectManager) CreateDefaultNetviews() (globalNetview *NetworkView, localNetview *NetworkView) {
	globalNetview, _ = objMgr.GetNetworkView(objMgr.globalAddressSpace)
	if globalNetview == nil {
		globalNetview, _ = objMgr.CreateNetworkView(objMgr.globalAddressSpace)
	}

	localNetview, _ = objMgr.GetNetworkView(objMgr.localAddressSpace)
	if localNetview == nil {
		localNetview, _ = objMgr.CreateNetworkView(objMgr.localAddressSpace)
	}

	return
}

func (objMgr *ObjectManager) GetNetworkView(name string) (*NetworkView, error) {
	res := make([]NetworkView, 1)

	payload := make(Payload)
	payload["name"] = name

	err := objMgr.connector.GetObject("networkview", payload, "", &res)

	if err != nil || res == nil || len(res) == 0 {
		return nil, err
	}

	return &res[0], nil
}

func BuildNetworkViewFromRef(ref string) *NetworkView {
	// networkview/ZG5zLm5ldHdvcmtfdmlldyQyMw:global_view/false
	r := regexp.MustCompile(`networkview/\w+:([^/]+)/\w+`)
	m := r.FindStringSubmatch(ref)

	if m == nil {
		return nil
	}

	return &NetworkView{
		Ref:  ref,
		Name: m[1],
	}
}

func BuildNetworkFromRef(ref string) *Network {
	// network/ZG5zLm5ldHdvcmskODkuMC4wLjAvMjQvMjU:89.0.0.0/24/global_view
	r := regexp.MustCompile(`network/\w+:(\d+\.\d+\.\d+\.\d+/\d+)/(\w+)`)
	m := r.FindStringSubmatch(ref)

	if m == nil {
		return nil
	}

	return &Network{
		Ref:         ref,
		NetviewName: m[2],
		Cidr:        m[1],
	}
}

func (objMgr *ObjectManager) GetNetwork(netview string, cidr string) (*Network, error) {
	res := make([]Network, 1)

	payload := make(Payload)
	payload["network_view"] = netview
	payload["network"] = cidr

	err := objMgr.connector.GetObject("network", payload, "", &res)

	if err != nil || res == nil || len(res) == 0 {
		return nil, err
	}

	return &res[0], nil
}

func GetIPAddressFromRef(ref string) string {
	// fixedaddress/ZG5zLmJpbmRfY25h:12.0.10.1/external
	r := regexp.MustCompile(`fixedaddress/\w+:(\d+\.\d+\.\d+\.\d+)/\w+`)
	m := r.FindStringSubmatch(ref)

	if m != nil {
		return m[1]
	}
	return ""
}

func (objMgr *ObjectManager) AllocateIP(netview string, cidr string, macAddress string) (*FixedAddress, error) {
	fixedAddr := new(FixedAddress)
	fixedAddr.NetviewName = netview
	fixedAddr.Cidr = cidr

	if len(macAddress) == 0 {
		macAddress = "00:00:00:00:00:00"
	}

	payload := make(Payload)
	payload["network_view"] = netview
	payload["ipv4addr"] = fmt.Sprintf("func:nextavailableip:%s,%s", cidr, netview)
	payload["mac"] = macAddress
	payload["extattrs"] = objMgr.getBasicEA(true)

	ref, err := objMgr.connector.CreateObject("fixedaddress", payload)
	fixedAddr.Ref = ref
	fixedAddr.IPAddress = GetIPAddressFromRef(ref)

	return fixedAddr, err
}

func (objMgr *ObjectManager) GetFixedAddress(netview string, ipAddr string) (*FixedAddress, error) {
	res := make([]FixedAddress, 1)

	payload := make(Payload)
	payload["network_view"] = netview
	payload["ipv4addr"] = ipAddr

	err := objMgr.connector.GetObject("fixedaddress", payload, "", &res)

	if err != nil || res == nil || len(res) == 0 {
		return nil, err
	}

	return &res[0], nil
}

func (objMgr *ObjectManager) ReleaseIP(netview string, ipAddr string) (string, error) {
	fixAddress, _ := objMgr.GetFixedAddress(netview, ipAddr)

	return objMgr.connector.DeleteObject(fixAddress.Ref)
}
