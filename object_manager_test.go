package ibclient

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
)

type fakeConnector struct {
	createObjectObj interface{}

	getObjectObj interface{}
	getObjectRef string

	deleteObjectRef string

	updateObjectObj interface{}
	updateObjectRef string

	resultObject interface{}

	fakeRefReturn string
}

func (c *fakeConnector) CreateObject(obj IBObject) (string, error) {
	Expect(obj).To(Equal(c.createObjectObj))

	return c.fakeRefReturn, nil
}

func (c *fakeConnector) GetObject(obj IBObject, ref string, res interface{}) (err error) {
	Expect(obj).To(Equal(c.getObjectObj))
	Expect(ref).To(Equal(c.getObjectRef))

	if ref == "" {
		switch obj.(type) {
		case *NetworkView:
			*res.(*[]NetworkView) = c.resultObject.([]NetworkView)
		case *NetworkContainer:
			*res.(*[]NetworkContainer) = c.resultObject.([]NetworkContainer)
		case *Network:
			*res.(*[]Network) = c.resultObject.([]Network)
		case *FixedAddress:
			*res.(*[]FixedAddress) = c.resultObject.([]FixedAddress)
		case *EADefinition:
			*res.(*[]EADefinition) = c.resultObject.([]EADefinition)
		case *CapacityReport:
			*res.(*[]CapacityReport) = c.resultObject.([]CapacityReport)
		case *UpgradeStatus:
			*res.(*[]UpgradeStatus) = c.resultObject.([]UpgradeStatus)
		case *Member:
			*res.(*[]Member) = c.resultObject.([]Member)
		case *Grid:
			*res.(*[]Grid) = c.resultObject.([]Grid)
		case *License:
			*res.(*[]License) = c.resultObject.([]License)
		case *HostRecord:
			*res.(*[]HostRecord) = c.resultObject.([]HostRecord)
		case *ZoneDelegated:
			*res.(*[]ZoneDelegated) = c.resultObject.([]ZoneDelegated)
     case *RecordA:
			*res.(*[]RecordA) = c.resultObject.([]RecordA)
		case *RecordAAAA:
			*res.(*[]RecordAAAA) = c.resultObject.([]RecordAAAA)
		case *RecordPTR:
			*res.(*[]RecordPTR) = c.resultObject.([]RecordPTR)
		case *RecordTXT:
			*res.(*[]RecordTXT) = c.resultObject.([]RecordTXT)
		case *RecordMX:
			*res.(*[]RecordMX) = c.resultObject.([]RecordMX)
		case *RecordSRV:
			*res.(*[]RecordSRV) = c.resultObject.([]RecordSRV)
		case *RecordAlias:
			*res.(*[]RecordAlias) = c.resultObject.([]RecordAlias)
		}
	} else {
		switch obj.(type) {
		case *NetworkView:
			*res.(*NetworkView) = c.resultObject.(NetworkView)
		}
	}

	err = nil
	return
}

func (c *fakeConnector) DeleteObject(ref string) (string, error) {
	Expect(ref).To(Equal(c.deleteObjectRef))

	return c.fakeRefReturn, nil
}

func (c *fakeConnector) UpdateObject(obj IBObject, ref string) (string, error) {
	Expect(obj).To(Equal(c.updateObjectObj))
	Expect(ref).To(Equal(c.updateObjectRef))

	return c.fakeRefReturn, nil
}

var _ = Describe("Object Manager", func() {

	Describe("Create Network View", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "Default View"
		fakeRefReturn := "networkview/ZG5zLm5ldHdvcmtfdmlldyQyMw:global_view/false"
		nvFakeConnector := &fakeConnector{
			createObjectObj: NewNetworkView(NetworkView{Name: netviewName}),
			resultObject:    NewNetworkView(NetworkView{Name: netviewName, Ref: fakeRefReturn}),
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(nvFakeConnector, cmpType, tenantID)
		nvFakeConnector.createObjectObj.(*NetworkView).Ea = objMgr.getBasicEA(false)
		nvFakeConnector.resultObject.(*NetworkView).Ea = objMgr.getBasicEA(false)

		var actualNetworkView *NetworkView
		var err error
		It("should pass expected NetworkView Object to CreateObject", func() {
			actualNetworkView, err = objMgr.CreateNetworkView(netviewName)
		})
		It("should return expected NetworkView Object", func() {
			Expect(actualNetworkView).To(Equal(nvFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Update Network View", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "Global View"
		fakeRefReturn := "networkview/ZG5zLm5ldHdvcmtfdmlldyQyMw:global_view/false"

		returnGetObject := NetworkView{Name: netviewName, Ref: fakeRefReturn, Ea: EA{"network-name": "net1", "Lock": "Removed"}}
		returnUpdateObject := NetworkView{Name: netviewName, Ref: fakeRefReturn, Ea: EA{"network-name": "net2", "New": "Added"}}
		getObjectObj := &NetworkView{}
		getObjectObj.returnFields = []string{"extattrs"}
		nvFakeConnector := &fakeConnector{
			getObjectObj:    getObjectObj,
			getObjectRef:    fakeRefReturn,
			fakeRefReturn:   fakeRefReturn,
			resultObject:    returnGetObject,
			updateObjectObj: &returnUpdateObject,
			updateObjectRef: fakeRefReturn,
		}

		objMgr := NewObjectManager(nvFakeConnector, cmpType, tenantID)

		var err error
		It("should pass expected updated object to UpdateObject", func() {
			addEA := EA{"network-name": "net2", "New": "Added"}
			delEA := EA{"Lock": "Removed"}
			err = objMgr.UpdateNetworkViewEA(fakeRefReturn, addEA, delEA)
		})
		It("should updated the GetObject with new EA and with no error", func() {
			Expect(returnGetObject).To(Equal(returnUpdateObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Create Network Container", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "Default View"
		cidr := "43.0.11.0/24"
		fakeRefReturn := "networkcontainer/ZG5zLm5ldHdvcmtfdmlldyQyMw:global_view/false"
		ncFakeConnector := &fakeConnector{
			createObjectObj: NewNetworkContainer(NetworkContainer{NetviewName: netviewName, Cidr: cidr}),
			resultObject:    NewNetworkContainer(NetworkContainer{NetviewName: netviewName, Cidr: cidr, Ref: fakeRefReturn}),
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(ncFakeConnector, cmpType, tenantID)
		ncFakeConnector.createObjectObj.(*NetworkContainer).Ea = objMgr.getBasicEA(true)
		ncFakeConnector.resultObject.(*NetworkContainer).Ea = objMgr.getBasicEA(true)

		var actualNetworkContainer *NetworkContainer
		var err error
		It("should pass expected NetworkContainer Object to CreateObject", func() {
			actualNetworkContainer, err = objMgr.CreateNetworkContainer(netviewName, cidr)
		})
		It("should return expected NetworkContainer Object", func() {
			Expect(actualNetworkContainer).To(Equal(ncFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Create Network", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "default_view"
		cidr := "43.0.11.0/24"
		networkName := "private-net"
		fakeRefReturn := "network/ZG5zLm5ldHdvcmskODkuMC4wLjAvMjQvMjU:43.0.11.0/24/default_view"
		nwFakeConnector := &fakeConnector{
			createObjectObj: NewNetwork(Network{NetviewName: netviewName, Cidr: cidr}),
			resultObject:    NewNetwork(Network{NetviewName: netviewName, Cidr: cidr, Ref: fakeRefReturn}),
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(nwFakeConnector, cmpType, tenantID)

		nwFakeConnector.createObjectObj.(*Network).Ea = objMgr.getBasicEA(true)
		nwFakeConnector.createObjectObj.(*Network).Ea["Network Name"] = networkName

		nwFakeConnector.resultObject.(*Network).Ea = objMgr.getBasicEA(true)
		nwFakeConnector.resultObject.(*Network).Ea["Network Name"] = networkName

		var actualNetwork *Network
		var err error
		It("should pass expected Network Object to CreateObject", func() {
			actualNetwork, err = objMgr.CreateNetwork(netviewName, cidr, networkName)
		})
		It("should return expected Network Object", func() {
			Expect(actualNetwork).To(Equal(nwFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Allocate Network", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "default_view"
		cidr := "142.0.22.0/24"
		prefixLen := uint(24)
		networkName := "private-net"
		fakeRefReturn := fmt.Sprintf("network/ZG5zLm5ldHdvcmskODkuMC4wLjAvMjQvMjU:%s/%s", cidr, netviewName)
		anFakeConnector := &fakeConnector{
			createObjectObj: NewNetwork(Network{
				NetviewName: netviewName,
				Cidr:        fmt.Sprintf("func:nextavailablenetwork:%s,%s,%d", cidr, netviewName, prefixLen),
			}),
			resultObject:  BuildNetworkFromRef(fakeRefReturn),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(anFakeConnector, cmpType, tenantID)

		anFakeConnector.createObjectObj.(*Network).Ea = objMgr.getBasicEA(true)
		anFakeConnector.createObjectObj.(*Network).Ea["Network Name"] = networkName

		var actualNetwork *Network
		var err error
		It("should pass expected Network Object to CreateObject", func() {
			actualNetwork, err = objMgr.AllocateNetwork(netviewName, cidr, prefixLen, networkName)
		})
		It("should return expected Network Object", func() {
			Expect(actualNetwork).To(Equal(anFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Allocate Specific IP", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "private"
		cidr := "53.0.0.0/24"
		ipAddr := "53.0.0.21"
		macAddr := "01:23:45:67:80:ab"
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		name := "testvm"
		fakeRefReturn := fmt.Sprintf("fixedaddress/ZG5zLmJpbmRfY25h:%s/private", ipAddr)

		asiFakeConnector := &fakeConnector{
			createObjectObj: NewFixedAddress(FixedAddress{
				NetviewName: netviewName,
				Cidr:        cidr,
				IPAddress:   ipAddr,
				Mac:         macAddr,
				Name:        name,
			}),
			resultObject: NewFixedAddress(FixedAddress{
				NetviewName: netviewName,
				Cidr:        cidr,
				IPAddress:   GetIPAddressFromRef(fakeRefReturn),
				Mac:         macAddr,
				Ref:         fakeRefReturn,
				Name:        name,
			}),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(asiFakeConnector, cmpType, tenantID)

		ea := objMgr.getBasicEA(true)
		asiFakeConnector.createObjectObj.(*FixedAddress).Ea = ea
		asiFakeConnector.createObjectObj.(*FixedAddress).Ea["VM ID"] = vmID
		asiFakeConnector.createObjectObj.(*FixedAddress).Ea["VM Name"] = vmName

		asiFakeConnector.resultObject.(*FixedAddress).Ea = ea
		asiFakeConnector.resultObject.(*FixedAddress).Ea["VM ID"] = vmID
		asiFakeConnector.resultObject.(*FixedAddress).Ea["VM Name"] = vmName

		var actualIP *FixedAddress
		var err error
		It("should pass expected Fixed Address Object to CreateObject", func() {
			actualIP, err = objMgr.AllocateIP(netviewName, cidr, ipAddr, macAddr, name, ea)
		})
		It("should return expected Fixed Address Object", func() {
			Expect(actualIP).To(Equal(asiFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Allocate Next Available IP", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "private"
		cidr := "53.0.0.0/24"
		ipAddr := fmt.Sprintf("func:nextavailableip:%s,%s", cidr, netviewName)
		macAddr := "01:23:45:67:80:ab"
		vmID := "93f9249abc039284"
		name := "testvm"
		vmName := "dummyvm"
		resultIP := "53.0.0.32"
		fakeRefReturn := fmt.Sprintf("fixedaddress/ZG5zLmJpbmRfY25h:%s/private", resultIP)

		aniFakeConnector := &fakeConnector{
			createObjectObj: NewFixedAddress(FixedAddress{
				NetviewName: netviewName,
				Cidr:        cidr,
				IPAddress:   ipAddr,
				Mac:         macAddr,
				Name:        name,
			}),
			resultObject: NewFixedAddress(FixedAddress{
				NetviewName: netviewName,
				Cidr:        cidr,
				IPAddress:   resultIP,
				Mac:         macAddr,
				Ref:         fakeRefReturn,
				Name:        name,
			}),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		ea := objMgr.getBasicEA(true)
		aniFakeConnector.createObjectObj.(*FixedAddress).Ea = ea
		aniFakeConnector.createObjectObj.(*FixedAddress).Ea["VM ID"] = vmID
		aniFakeConnector.createObjectObj.(*FixedAddress).Ea["VM Name"] = vmName

		aniFakeConnector.resultObject.(*FixedAddress).Ea = ea
		aniFakeConnector.resultObject.(*FixedAddress).Ea["VM ID"] = vmID
		aniFakeConnector.resultObject.(*FixedAddress).Ea["VM Name"] = vmName

		var actualIP *FixedAddress
		var err error
		It("should pass expected Fixed Address Object to CreateObject", func() {
			actualIP, err = objMgr.AllocateIP(netviewName, cidr, ipAddr, macAddr, name, ea)
		})

		It("should return expected Fixed Address Object", func() {
			Expect(actualIP).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Allocate next available host Record without dns", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "private"
		cidr := "53.0.0.0/24"
		macAddr := "01:23:45:67:80:ab"
		ipAddr := fmt.Sprintf("func:nextavailableip:%s,%s", cidr, netviewName)
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		recordName := "test"
		enabledns := false
		dnsView := "default"
		fakeRefReturn := fmt.Sprintf("record:host/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)
		resultIPV4Addrs := NewHostRecordIpv4Addr(HostRecordIpv4Addr{Ipv4Addr: ipAddr, Mac: macAddr})
		enableDNS := new(bool)
		*enableDNS = enabledns

		aniFakeConnector := &fakeConnector{
			createObjectObj: NewHostRecord(HostRecord{
				Name:        recordName,
				View:        dnsView,
				EnableDns:   enableDNS,
				NetworkView: netviewName,
				Ipv4Addrs:   []HostRecordIpv4Addr{*resultIPV4Addrs},
			}),
			getObjectRef: fakeRefReturn,
			getObjectObj: NewHostRecord(HostRecord{
				Name:        recordName,
				View:        dnsView,
				EnableDns:   enableDNS,
				NetworkView: netviewName,
				Ipv4Addrs:   []HostRecordIpv4Addr{*resultIPV4Addrs},
				Ref:         fakeRefReturn,
			}),
			resultObject: NewHostRecord(HostRecord{
				Name:        recordName,
				View:        dnsView,
				EnableDns:   enableDNS,
				NetworkView: netviewName,
				Ipv4Addrs:   []HostRecordIpv4Addr{*resultIPV4Addrs},
				Ref:         fakeRefReturn,
			}),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		ea := objMgr.getBasicEA(true)
		aniFakeConnector.createObjectObj.(*HostRecord).Ea = ea
		aniFakeConnector.createObjectObj.(*HostRecord).Ea["VM ID"] = vmID
		aniFakeConnector.createObjectObj.(*HostRecord).Ea["VM Name"] = vmName

		aniFakeConnector.resultObject.(*HostRecord).Ea = ea
		aniFakeConnector.resultObject.(*HostRecord).Ea["VM ID"] = vmID
		aniFakeConnector.resultObject.(*HostRecord).Ea["VM Name"] = vmName

		aniFakeConnector.getObjectObj.(*HostRecord).Ea = ea
		aniFakeConnector.getObjectObj.(*HostRecord).Ea["VM ID"] = vmID
		aniFakeConnector.getObjectObj.(*HostRecord).Ea["VM Name"] = vmName

		var actualRecord *HostRecord
		var err error
		It("should pass expected host record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateHostRecord(enabledns, recordName, netviewName, dnsView, cidr, ipAddr, macAddr, ea)
		})
		It("should return expected host record Object", func() {
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Allocate next available host Record with dns", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "private"
		cidr := "53.0.0.0/24"
		macAddr := "01:23:45:67:80:ab"
		ipAddr := fmt.Sprintf("func:nextavailableip:%s,%s", cidr, netviewName)
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		recordName := "test"
		enabledns := true
		dnsView := "default"
		fakeRefReturn := fmt.Sprintf("record:host/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)
		resultIPV4Addrs := NewHostRecordIpv4Addr(HostRecordIpv4Addr{Ipv4Addr: ipAddr, Mac: macAddr})
		enableDNS := new(bool)
		*enableDNS = enabledns

		aniFakeConnector := &fakeConnector{
			createObjectObj: NewHostRecord(HostRecord{
				Name:        recordName,
				View:        dnsView,
				EnableDns:   enableDNS,
				NetworkView: netviewName,
				Ipv4Addrs:   []HostRecordIpv4Addr{*resultIPV4Addrs},
			}),
			getObjectRef: fakeRefReturn,
			getObjectObj: NewHostRecord(HostRecord{
				Name:        recordName,
				View:        dnsView,
				EnableDns:   enableDNS,
				NetworkView: netviewName,
				Ipv4Addrs:   []HostRecordIpv4Addr{*resultIPV4Addrs},
				Ref:         fakeRefReturn,
			}),
			resultObject: NewHostRecord(HostRecord{
				Name:        recordName,
				View:        dnsView,
				EnableDns:   enableDNS,
				NetworkView: netviewName,
				Ipv4Addrs:   []HostRecordIpv4Addr{*resultIPV4Addrs},
				Ref:         fakeRefReturn,
			}),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		ea := objMgr.getBasicEA(true)
		aniFakeConnector.createObjectObj.(*HostRecord).Ea = ea
		aniFakeConnector.createObjectObj.(*HostRecord).Ea["VM ID"] = vmID
		aniFakeConnector.createObjectObj.(*HostRecord).Ea["VM Name"] = vmName

		aniFakeConnector.resultObject.(*HostRecord).Ea = ea
		aniFakeConnector.resultObject.(*HostRecord).Ea["VM ID"] = vmID
		aniFakeConnector.resultObject.(*HostRecord).Ea["VM Name"] = vmName

		aniFakeConnector.getObjectObj.(*HostRecord).Ea = ea
		aniFakeConnector.getObjectObj.(*HostRecord).Ea["VM ID"] = vmID
		aniFakeConnector.getObjectObj.(*HostRecord).Ea["VM Name"] = vmName

		var actualRecord *HostRecord
		var err error
		It("should pass expected host record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateHostRecord(enabledns, recordName, netviewName, dnsView, cidr, ipAddr, macAddr, ea)
		})
		It("should return expected host record Object", func() {
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Allocate specific host Record without dns", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "private"
		cidr := "53.0.0.0/24"
		macAddr := "01:23:45:67:80:ab"
		ipAddr := "53.0.0.1"
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		enabledns := false
		dnsView := "default"
		recordName := "test"
		fakeRefReturn := fmt.Sprintf("record:host/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)
		resultIPV4Addrs := NewHostRecordIpv4Addr(HostRecordIpv4Addr{Ipv4Addr: ipAddr, Mac: macAddr})
		enableDNS := new(bool)
		*enableDNS = enabledns

		aniFakeConnector := &fakeConnector{
			createObjectObj: NewHostRecord(HostRecord{
				Name:        recordName,
				View:        dnsView,
				EnableDns:   enableDNS,
				NetworkView: netviewName,
				Ipv4Addrs:   []HostRecordIpv4Addr{*resultIPV4Addrs},
			}),
			getObjectRef: fakeRefReturn,
			getObjectObj: NewHostRecord(HostRecord{
				Name:        recordName,
				EnableDns:   enableDNS,
				View:        dnsView,
				NetworkView: netviewName,
				Ipv4Addrs:   []HostRecordIpv4Addr{*resultIPV4Addrs},
				Ref:         fakeRefReturn,
			}),
			resultObject: NewHostRecord(HostRecord{
				Name:        recordName,
				EnableDns:   enableDNS,
				View:        dnsView,
				NetworkView: netviewName,
				Ipv4Addrs:   []HostRecordIpv4Addr{*resultIPV4Addrs},
				Ref:         fakeRefReturn,
			}),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		ea := objMgr.getBasicEA(true)
		aniFakeConnector.createObjectObj.(*HostRecord).Ea = ea
		aniFakeConnector.createObjectObj.(*HostRecord).Ea["VM ID"] = vmID
		aniFakeConnector.createObjectObj.(*HostRecord).Ea["VM Name"] = vmName

		aniFakeConnector.resultObject.(*HostRecord).Ea = ea
		aniFakeConnector.resultObject.(*HostRecord).Ea["VM ID"] = vmID
		aniFakeConnector.resultObject.(*HostRecord).Ea["VM Name"] = vmName

		aniFakeConnector.getObjectObj.(*HostRecord).Ea = ea
		aniFakeConnector.getObjectObj.(*HostRecord).Ea["VM ID"] = vmID
		aniFakeConnector.getObjectObj.(*HostRecord).Ea["VM Name"] = vmName

		var actualRecord *HostRecord
		var err error
		It("should pass expected host record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateHostRecord(enabledns, recordName, netviewName, dnsView, cidr, ipAddr, macAddr, ea)
		})

		It("should return expected host record Object", func() {
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Allocate specific host Record with dns", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "private"
		cidr := "53.0.0.0/24"
		macAddr := "01:23:45:67:80:ab"
		ipAddr := "53.0.0.1"
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		enabledns := true
		dnsView := "default"
		recordName := "test"
		fakeRefReturn := fmt.Sprintf("record:host/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)
		resultIPV4Addrs := NewHostRecordIpv4Addr(HostRecordIpv4Addr{Ipv4Addr: ipAddr, Mac: macAddr})
		enableDNS := new(bool)
		*enableDNS = enabledns

		aniFakeConnector := &fakeConnector{
			createObjectObj: NewHostRecord(HostRecord{
				Name:        recordName,
				View:        dnsView,
				EnableDns:   enableDNS,
				NetworkView: netviewName,
				Ipv4Addrs:   []HostRecordIpv4Addr{*resultIPV4Addrs},
			}),
			getObjectRef: fakeRefReturn,
			getObjectObj: NewHostRecord(HostRecord{
				Name:        recordName,
				EnableDns:   enableDNS,
				View:        dnsView,
				NetworkView: netviewName,
				Ipv4Addrs:   []HostRecordIpv4Addr{*resultIPV4Addrs},
				Ref:         fakeRefReturn,
			}),
			resultObject: NewHostRecord(HostRecord{
				Name:        recordName,
				EnableDns:   enableDNS,
				View:        dnsView,
				NetworkView: netviewName,
				Ipv4Addrs:   []HostRecordIpv4Addr{*resultIPV4Addrs},
				Ref:         fakeRefReturn,
			}),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		ea := objMgr.getBasicEA(true)
		aniFakeConnector.createObjectObj.(*HostRecord).Ea = ea
		aniFakeConnector.createObjectObj.(*HostRecord).Ea["VM ID"] = vmID
		aniFakeConnector.createObjectObj.(*HostRecord).Ea["VM Name"] = vmName

		aniFakeConnector.resultObject.(*HostRecord).Ea = ea
		aniFakeConnector.resultObject.(*HostRecord).Ea["VM ID"] = vmID
		aniFakeConnector.resultObject.(*HostRecord).Ea["VM Name"] = vmName

		aniFakeConnector.getObjectObj.(*HostRecord).Ea = ea
		aniFakeConnector.getObjectObj.(*HostRecord).Ea["VM ID"] = vmID
		aniFakeConnector.getObjectObj.(*HostRecord).Ea["VM Name"] = vmName

		var actualRecord *HostRecord
		var err error
		It("should pass expected host record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateHostRecord(enabledns, recordName, netviewName, dnsView, cidr, ipAddr, macAddr, ea)
		})

		It("should return expected host record Object", func() {
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Allocate CNAME Record ", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		canonical := "test-canonical"
		dnsView := "default"
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		recordName := "test"
		fakeRefReturn := fmt.Sprintf("record:cname/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)

		aniFakeConnector := &fakeConnector{
			createObjectObj: NewRecordCNAME(RecordCNAME{
				Name:      recordName,
				View:      dnsView,
				Canonical: canonical,
			}),
			getObjectRef: fakeRefReturn,
			getObjectObj: NewRecordCNAME(RecordCNAME{
				Name:      recordName,
				View:      dnsView,
				Canonical: canonical,
				Ref:       fakeRefReturn,
			}),
			resultObject: NewRecordCNAME(RecordCNAME{
				Name:      recordName,
				View:      dnsView,
				Canonical: canonical,
				Ref:       fakeRefReturn,
			}),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)
		ea := objMgr.getBasicEA(true)
		aniFakeConnector.createObjectObj.(*RecordCNAME).Ea = ea
		aniFakeConnector.createObjectObj.(*RecordCNAME).Ea["VM ID"] = vmID
		aniFakeConnector.createObjectObj.(*RecordCNAME).Ea["VM Name"] = vmName

		aniFakeConnector.resultObject.(*RecordCNAME).Ea = ea
		aniFakeConnector.resultObject.(*RecordCNAME).Ea["VM ID"] = vmID
		aniFakeConnector.resultObject.(*RecordCNAME).Ea["VM Name"] = vmName

		aniFakeConnector.getObjectObj.(*RecordCNAME).Ea = ea
		aniFakeConnector.getObjectObj.(*RecordCNAME).Ea["VM ID"] = vmID
		aniFakeConnector.getObjectObj.(*RecordCNAME).Ea["VM Name"] = vmName
		var actualRecord *RecordCNAME
		var err error
		It("should pass expected CNAME record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateCNAMERecord(canonical, recordName, dnsView, ea)
		})
		It("should return expected CNAME record Object", func() {
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Create EA Definition", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		comment := "Test Extensible Attribute"
		flags := "CGV"
		listValues := []EADefListValue{"True", "False"}
		name := "TestEA"
		eaType := "string"
		allowedTypes := []string{"arecord", "aaarecord", "ptrrecord"}
		ead := EADefinition{
			Name:               name,
			Comment:            comment,
			Flags:              flags,
			ListValues:         listValues,
			Type:               eaType,
			AllowedObjectTypes: allowedTypes}
		fakeRefReturn := "extensibleattributedef/ZG5zLm5ldHdvcmtfdmlldyQyMw:TestEA"
		eadFakeConnector := &fakeConnector{
			createObjectObj: NewEADefinition(ead),
			resultObject:    NewEADefinition(ead),
			fakeRefReturn:   fakeRefReturn,
		}
		eadFakeConnector.resultObject.(*EADefinition).Ref = fakeRefReturn

		objMgr := NewObjectManager(eadFakeConnector, cmpType, tenantID)

		var actualEADef *EADefinition
		var err error
		It("should pass expected EA Definintion Object to CreateObject", func() {
			actualEADef, err = objMgr.CreateEADefinition(ead)
		})
		It("should return expected EA Definition Object", func() {
			Expect(actualEADef).To(Equal(eadFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get Network View", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "Default View"
		fakeRefReturn := "networkview/ZG5zLm5ldHdvcmtfdmlldyQyMw:global_view/false"
		nvFakeConnector := &fakeConnector{
			getObjectObj: NewNetworkView(NetworkView{Name: netviewName}),
			getObjectRef: "",
			resultObject: []NetworkView{*NewNetworkView(NetworkView{Name: netviewName, Ref: fakeRefReturn})},
		}

		objMgr := NewObjectManager(nvFakeConnector, cmpType, tenantID)

		var actualNetworkView *NetworkView
		var err error
		It("should pass expected NetworkView Object to GetObject", func() {
			actualNetworkView, err = objMgr.GetNetworkView(netviewName)
		})
		It("should return expected NetworkView Object", func() {
			Expect(*actualNetworkView).To(Equal(nvFakeConnector.resultObject.([]NetworkView)[0]))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get Network Container", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "Default View"
		cidr := "43.0.11.0/24"
		fakeRefReturn := "networkcontainer/ZG5zLm5ldHdvcmtfdmlldyQyMw:global_view/false"
		ncFakeConnector := &fakeConnector{
			getObjectObj: NewNetworkContainer(NetworkContainer{NetviewName: netviewName, Cidr: cidr}),
			getObjectRef: "",
			resultObject: []NetworkContainer{*NewNetworkContainer(NetworkContainer{NetviewName: netviewName, Cidr: cidr, Ref: fakeRefReturn})},
		}

		objMgr := NewObjectManager(ncFakeConnector, cmpType, tenantID)

		var actualNetworkContainer *NetworkContainer
		var err error
		It("should pass expected NetworkContainer Object to GetObject", func() {
			actualNetworkContainer, err = objMgr.GetNetworkContainer(netviewName, cidr)
		})
		It("should return expected NetworkContainer Object", func() {
			Expect(*actualNetworkContainer).To(Equal(ncFakeConnector.resultObject.([]NetworkContainer)[0]))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get Network", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "default_view"
		cidr := "28.0.42.0/24"
		networkName := "private-net"
		ea := EA{"Network Name": networkName}
		fakeRefReturn := fmt.Sprintf("network/ZG5zLm5ldHdvcmskODkuMC4wLjAvMjQvMjU:%s/%s", cidr, netviewName)
		nwFakeConnector := &fakeConnector{
			getObjectObj: NewNetwork(Network{NetviewName: netviewName, Cidr: cidr}),
			getObjectRef: "",
			resultObject: []Network{*NewNetwork(Network{NetviewName: netviewName, Cidr: cidr, Ref: fakeRefReturn})},
		}

		nwFakeConnector.getObjectObj.(*Network).eaSearch = EASearch(ea)
		nwFakeConnector.resultObject.([]Network)[0].eaSearch = EASearch(ea)

		objMgr := NewObjectManager(nwFakeConnector, cmpType, tenantID)

		var actualNetwork *Network
		var err error
		It("should pass expected Network Object to GetObject", func() {
			actualNetwork, err = objMgr.GetNetwork(netviewName, cidr, ea)
		})
		It("should return expected Network Object", func() {
			Expect(*actualNetwork).To(Equal(nwFakeConnector.resultObject.([]Network)[0]))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get Network with Reference", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		cidr := "28.0.42.0/24"
		netviewName := "default_view"
		getRef := fmt.Sprintf("network/ZG5zLm5ldHdvcmskODkuMC4wLjAvMjQvMjU:%s/%s", cidr, netviewName)
		fakeRefReturn := getRef
		nwFakeConnector := &fakeConnector{
			getObjectObj:  NewNetwork(Network{}),
			getObjectRef:  getRef,
			resultObject:  []Network{*NewNetwork(Network{})},
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(nwFakeConnector, cmpType, tenantID)

		var actualRef *Network
		var err error
		It("should pass expected Network Ref to getObject", func() {
			actualRef, err = objMgr.GetNetworkwithref(fakeRefReturn)
		})
		It("should return expected Network record Ref", func() {
			Expect(*actualRef).To(Equal(nwFakeConnector.resultObject.([]Network)[0]))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get Fixed Address", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "private"
		cidr := "53.0.0.0/24"
		ipAddr := "53.0.0.21"
		macAddr := "01:23:45:67:80:ab"
		fakeRefReturn := fmt.Sprintf("fixedaddress/ZG5zLmJpbmRfY25h:%s/private", ipAddr)

		fipFakeConnector := &fakeConnector{
			getObjectObj: NewFixedAddress(FixedAddress{
				NetviewName: netviewName,
				Cidr:        cidr,
				IPAddress:   ipAddr,
				Mac:         macAddr,
			}),
			getObjectRef: "",
			resultObject: []FixedAddress{*NewFixedAddress(FixedAddress{
				NetviewName: netviewName,
				Cidr:        cidr,
				IPAddress:   GetIPAddressFromRef(fakeRefReturn),
				Mac:         macAddr,
				Ref:         fakeRefReturn,
			})},
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(fipFakeConnector, cmpType, tenantID)

		var actualIP *FixedAddress
		var err error
		It("should pass expected Fixed Address Object to GetObject", func() {
			actualIP, err = objMgr.GetFixedAddress(netviewName, cidr, ipAddr, macAddr)
		})
		It("should return expected Fixed Address Object", func() {
			Expect(*actualIP).To(Equal(fipFakeConnector.resultObject.([]FixedAddress)[0]))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get Host Record Without DNS", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "private"
		cidr := "53.0.0.0/24"
		ipAddr := "53.0.0.21"
		hostName := "test"
		fakeRefReturn := fmt.Sprintf("record:host/ZG5zLmJpbmRfY25h:%s/%20%20", hostName)
		fipFakeConnector := &fakeConnector{
			getObjectObj: NewHostRecord(HostRecord{
				Name: hostName,
			}),
			getObjectRef: "",
			resultObject: []HostRecord{*NewHostRecord(HostRecord{
				Name: hostName,
				Ref:  fakeRefReturn,
			})},
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(fipFakeConnector, cmpType, tenantID)

		var actualhostRecord *HostRecord
		var err error
		It("should pass expected Host record Object to GetObject", func() {
			actualhostRecord, err = objMgr.GetHostRecord(hostName, netviewName, cidr, ipAddr)
		})
		It("should return expected Host record Object", func() {
			Expect(*actualhostRecord).To(Equal(fipFakeConnector.resultObject.([]HostRecord)[0]))
			Expect(err).To(BeNil())
		})

	})

	Describe("Get EA Definition", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		comment := "Test Extensible Attribute"
		flags := "CGV"
		listValues := []EADefListValue{"True", "False"}
		name := "TestEA"
		eaType := "string"
		allowedTypes := []string{"arecord", "aaarecord", "ptrrecord"}
		ead := EADefinition{
			Name: name,
		}
		fakeRefReturn := "extensibleattributedef/ZG5zLm5ldHdvcmtfdmlldyQyMw:TestEA"
		eadRes := EADefinition{
			Name:               name,
			Comment:            comment,
			Flags:              flags,
			ListValues:         listValues,
			Type:               eaType,
			AllowedObjectTypes: allowedTypes,
			Ref:                fakeRefReturn,
		}

		eadFakeConnector := &fakeConnector{
			getObjectObj:  NewEADefinition(ead),
			getObjectRef:  "",
			resultObject:  []EADefinition{*NewEADefinition(eadRes)},
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(eadFakeConnector, cmpType, tenantID)

		var actualEADef *EADefinition
		var err error
		It("should pass expected EA Definintion Object to GetObject", func() {
			actualEADef, err = objMgr.GetEADefinition(name)
		})
		It("should return expected EA Definition Object", func() {
			Expect(*actualEADef).To(Equal(eadFakeConnector.resultObject.([]EADefinition)[0]))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete Network", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "default_view"
		cidr := "28.0.42.0/24"
		deleteRef := fmt.Sprintf("network/ZG5zLm5ldHdvcmskODkuMC4wLjAvMjQvMjU:%s/%s", cidr, netviewName)
		fakeRefReturn := deleteRef
		nwFakeConnector := &fakeConnector{
			deleteObjectRef: deleteRef,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(nwFakeConnector, cmpType, tenantID)

		var actualRef string
		var err error
		It("should pass expected Network Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteNetwork(deleteRef, netviewName)
		})
		It("should return expected Network Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete Network View", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		fakeRefReturn := "networkview/ZG5zLm5ldHdvcmtfdmlldyQyMw:global_view/false"
		deleteRef := fakeRefReturn
		nwFakeConnector := &fakeConnector{
			deleteObjectRef: deleteRef,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(nwFakeConnector, cmpType, tenantID)

		var actualRef string
		var err error
		It("should pass expected Network View Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteNetworkView(deleteRef)
		})
		It("should return expected Network View Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete Fixed Address", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "private"
		cidr := "83.0.101.0/24"
		ipAddr := "83.0.101.68"
		macAddr := "01:23:45:67:80:ab"
		fakeRefReturn := fmt.Sprintf("fixedaddress/ZG5zLmJpbmRfY25h:%s/private", ipAddr)

		fipFakeConnector := &fakeConnector{
			getObjectObj: NewFixedAddress(FixedAddress{
				NetviewName: netviewName,
				Cidr:        cidr,
				IPAddress:   ipAddr,
				Mac:         macAddr,
			}),
			getObjectRef: "",
			resultObject: []FixedAddress{*NewFixedAddress(FixedAddress{
				NetviewName: netviewName,
				Cidr:        cidr,
				IPAddress:   GetIPAddressFromRef(fakeRefReturn),
				Mac:         macAddr,
				Ref:         fakeRefReturn,
			})},
			deleteObjectRef: fakeRefReturn,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(fipFakeConnector, cmpType, tenantID)

		var actualRef string
		var err error
		It("should pass expected Fixed Address Object to GetObject and DeleteObject", func() {
			actualRef, err = objMgr.ReleaseIP(netviewName, cidr, ipAddr, macAddr)
		})
		It("should return expected Fixed Address Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete Host Record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		hostName := "test"
		deleteRef := fmt.Sprintf("record:host/ZG5zLmJpbmRfY25h:%s/%20%20", hostName)
		fakeRefReturn := deleteRef
		nwFakeConnector := &fakeConnector{
			deleteObjectRef: deleteRef,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(nwFakeConnector, cmpType, tenantID)

		var actualRef string
		var err error
		It("should pass expected Host record Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteHostRecord(deleteRef)
		})
		It("should return expected Host record Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete CNAME Record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		recordName := "test"
		deleteRef := fmt.Sprintf("record:CNAME/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)
		fakeRefReturn := deleteRef
		nwFakeConnector := &fakeConnector{
			deleteObjectRef: deleteRef,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(nwFakeConnector, cmpType, tenantID)

		var actualRef string
		var err error
		It("should pass expected CNAME record Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteCNAMERecord(deleteRef)
		})
		It("should return expected CNAME record Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})

	Describe("BuildNetworkViewFromRef", func() {
		netviewName := "default_view"
		netviewRef := fmt.Sprintf("networkview/ZG5zLm5ldHdvcmtfdmlldyQyMw:%s/false", netviewName)

		expectedNetworkView := NetworkView{Ref: netviewRef, Name: netviewName}
		It("should return expected Network View Object", func() {
			Expect(*BuildNetworkViewFromRef(netviewRef)).To(Equal(expectedNetworkView))
		})
		It("should failed if bad Network View Ref is provided", func() {
			Expect(BuildNetworkViewFromRef("bad")).To(BeNil())
		})
	})

	Describe("BuildNetworkFromRef", func() {
		netviewName := "test_view"
		cidr := "23.11.0.0/24"
		networkRef := fmt.Sprintf("network/ZG5zLm5ldHdvcmtfdmlldyQyMw:%s/%s", cidr, netviewName)

		expectedNetwork := Network{Ref: networkRef, NetviewName: netviewName, Cidr: cidr}
		It("should return expected Network Object", func() {
			Expect(*BuildNetworkFromRef(networkRef)).To(Equal(expectedNetwork))
		})
		It("should failed if bad Network Ref is provided", func() {
			Expect(BuildNetworkFromRef("network/ZG5zLm5ldHdvcmtfdmlldyQyMw")).To(BeNil())
		})
	})

	Describe("Get Capacity report", func() {
		cmpType := "Heka"
		tenantID := "0123"
		var name string = "Member1"
		fakeRefReturn := fmt.Sprintf("member/ZG5zLmJpbmRfY25h:/%s", name)

		fakeConnector := &fakeConnector{
			getObjectObj: NewCapcityReport(CapacityReport{Name: name}),
			getObjectRef: "",
			resultObject: []CapacityReport{*NewCapcityReport(CapacityReport{
				Ref:  fakeRefReturn,
				Name: name,
			})},
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(fakeConnector, cmpType, tenantID)

		var actualReport []CapacityReport
		var err error

		It("should pass expected Capacityreport object to GetObject", func() {
			actualReport, err = objMgr.GetCapacityReport(name)
		})
		It("should return expected CapacityReport Object", func() {
			Expect(actualReport[0]).To(Equal(fakeConnector.resultObject.([]CapacityReport)[0]))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get upgrade status", func() {
		cmpType := "Heka"
		tenantID := "0123"
		var StatusType string = "GRID"
		fakeRefReturn := fmt.Sprintf("upgradestatus/Li51cGdyYWRlc3RhdHVzJHVwZ3JhZGVfc3RhdHVz:test")

		USFakeConnector := &fakeConnector{
			getObjectObj: NewUpgradeStatus(UpgradeStatus{Type: StatusType}),
			getObjectRef: "",
			resultObject: []UpgradeStatus{*NewUpgradeStatus(UpgradeStatus{
				Ref:  fakeRefReturn,
				Type: StatusType,
			})},
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(USFakeConnector, cmpType, tenantID)

		var actualStatus []UpgradeStatus
		var err error

		It("should pass expected upgradestatus object to GetObject", func() {
			actualStatus, err = objMgr.GetUpgradeStatus(StatusType)
		})
		It("should return expected upgradestatus Object", func() {
			Expect(actualStatus[0]).To(Equal(USFakeConnector.resultObject.([]UpgradeStatus)[0]))
			Expect(err).To(BeNil())
		})

	})
	Describe("Get upgrade status Error case", func() {
		cmpType := "Heka"
		tenantID := "0123"
		StatusType := ""
		fakeRefReturn := fmt.Sprintf("upgradestatus/Li51cGdyYWRlc3RhdHVzJHVwZ3JhZGVfc3RhdHVz:test")
		expectErr := errors.New("Status type can not be nil")
		USFakeConnector := &fakeConnector{
			getObjectObj: NewUpgradeStatus(UpgradeStatus{Type: StatusType}),
			getObjectRef: "",
			resultObject: []UpgradeStatus{*NewUpgradeStatus(UpgradeStatus{
				Ref:  fakeRefReturn,
				Type: StatusType,
			})},
			fakeRefReturn: fakeRefReturn,
		}
		objMgr := NewObjectManager(USFakeConnector, cmpType, tenantID)
		It("upgradestatus object to GetObject", func() {
			_, err := objMgr.GetUpgradeStatus(StatusType)
			Expect(err).To(Equal(expectErr))
		})

	})
	Describe("GetAllMembers", func() {
		cmpType := "Heka"
		tenantID := "0123"
		var err error
		fakeRefReturn := fmt.Sprintf("member/Li51cGdyYWRlc3RhdHVzJHVwZ3JhZGVfc3RhdHVz:test")
		returnFields := []string{"host_name", "node_info", "time_zone"}
		MemFakeConnector := &fakeConnector{
			getObjectObj: NewMember(Member{}),
			getObjectRef: "",
			resultObject: []Member{*NewMember(Member{
				Ref: fakeRefReturn,
			})},
			fakeRefReturn: fakeRefReturn,
		}
		objMgr := NewObjectManager(MemFakeConnector, cmpType, tenantID)
		var actualMembers []Member
		It("should return expected member Object", func() {
			actualMembers, err = objMgr.GetAllMembers()
			Expect(actualMembers[0]).To(Equal(MemFakeConnector.resultObject.([]Member)[0]))
			Expect(actualMembers[0].returnFields).To(Equal(returnFields))
			Expect(err).To(BeNil())
		})
	})

	Describe("GetGridInfo", func() {
		cmpType := "Heka"
		tenantID := "0123"
		var err error
		fakeRefReturn := fmt.Sprintf("grid/Li511cGdyYWRlc3RhdHVzJHVwZ3JhZGVfc3RhdHVz:test")
		returnFields := []string{"name", "ntp_setting"}
		GridFakeConnector := &fakeConnector{
			getObjectObj: NewGrid(Grid{}),
			getObjectRef: "",
			resultObject: []Grid{*NewGrid(Grid{
				Ref: fakeRefReturn,
			})},
			fakeRefReturn: fakeRefReturn,
		}
		objMgr := NewObjectManager(GridFakeConnector, cmpType, tenantID)
		var actualGridInfo []Grid
		It("should return expected Grid Object", func() {
			actualGridInfo, err = objMgr.GetGridInfo()
			Expect(actualGridInfo[0]).To(Equal(GridFakeConnector.resultObject.([]Grid)[0]))
			Expect(actualGridInfo[0].returnFields).To(Equal(returnFields))
			Expect(err).To(BeNil())
		})
	})

	Describe("GetGridLicense", func() {
		cmpType := "Heka"
		tenantID := "0123"
		var err error
		fakeRefReturn := fmt.Sprintf("license/Li511cGdyYWRlc3RhdHVzJHVwZ3JhZGVfc3RhdHVz:test")
		returnFields := []string{"expiration_status",
			"expiry_date",
			"key",
			"limit",
			"limit_context",
			"type"}
		LicFakeConnector := &fakeConnector{
			getObjectObj: NewGridLicense(License{}),
			getObjectRef: "",
			resultObject: []License{*NewGridLicense(License{
				Ref: fakeRefReturn,
			})},
			fakeRefReturn: fakeRefReturn,
		}
		objMgr := NewObjectManager(LicFakeConnector, cmpType, tenantID)
		var actualGridLicense []License
		It("should return expected License Object", func() {
			actualGridLicense, err = objMgr.GetGridLicense()
			Expect(actualGridLicense[0]).To(Equal(LicFakeConnector.resultObject.([]License)[0]))
			Expect(actualGridLicense[0].returnFields).To(Equal(returnFields))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get Zone Delegated", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		fqdn := "dzone.example.com"
		fakeRefReturn := "zone_delegated/ZG5zLnpvbmUkLl9kZWZhdWx0LnphLmNvLmFic2EuY2Fhcy5vaG15Z2xiLmdzbGJpYmNsaWVudA:dzone.example.com/default"
		zdFakeConnector := &fakeConnector{
			getObjectObj: NewZoneDelegated(ZoneDelegated{Fqdn: fqdn}),
			getObjectRef: "",
			resultObject: []ZoneDelegated{*NewZoneDelegated(ZoneDelegated{Fqdn: fqdn, Ref: fakeRefReturn})},
		}

		objMgr := NewObjectManager(zdFakeConnector, cmpType, tenantID)

		var actualZoneDelegated *ZoneDelegated
		var err error
		It("should pass expected ZoneDelegated Object to GetObject", func() {
			actualZoneDelegated, err = objMgr.GetZoneDelegated(fqdn)
		})
		It("should return expected ZoneDelegated Object", func() {
			Expect(*actualZoneDelegated).To(Equal(zdFakeConnector.resultObject.([]ZoneDelegated)[0]))
			Expect(err).To(BeNil())
		})
		It("should return nil if fqdn is empty", func() {
			zdFakeConnector.getObjectObj.(*ZoneDelegated).Fqdn = ""
			actualZoneDelegated, err = objMgr.GetZoneDelegated("")
			Expect(actualZoneDelegated).To(BeNil())
			Expect(err).To(BeNil())
		})
	})

	Describe("Create Zone Delegated", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		fqdn := "dzone.example.com"
		delegateTo := []NameServer{
			{Address: "10.0.0.1", Name: "test1.dzone.example.com"},
			{Address: "10.0.0.2", Name: "test2.dzone.example.com"}}
		fakeRefReturn := "zone_delegated/ZG5zLnpvbmUkLl9kZWZhdWx0LnphLmNvLmFic2EuY2Fhcy5vaG15Z2xiLmdzbGJpYmNsaWVudA:dzone.example.com/default"
		zdFakeConnector := &fakeConnector{
			createObjectObj: NewZoneDelegated(ZoneDelegated{Fqdn: fqdn, DelegateTo: delegateTo}),
			resultObject:    NewZoneDelegated(ZoneDelegated{Fqdn: fqdn, DelegateTo: delegateTo, Ref: fakeRefReturn}),
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(zdFakeConnector, cmpType, tenantID)

		var actualZoneDelegated *ZoneDelegated
		var err error
		It("should pass expected ZoneDelegated Object to CreateObject", func() {
			actualZoneDelegated, err = objMgr.CreateZoneDelegated(fqdn, delegateTo)
		})
		It("should return expected ZoneDelegated Object", func() {
			Expect(actualZoneDelegated).To(Equal(zdFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Update Zone Delegated", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		fakeRefReturn := "zone_delegated/ZG5zLnpvbmUkLl9kZWZhdWx0LnphLmNvLmFic2EuY2Fhcy5vaG15Z2xiLmdzbGJpYmNsaWVudA:dzone.example.com/default"
		delegateTo := []NameServer{
			{Address: "10.0.0.1", Name: "test1.dzone.example.com"},
			{Address: "10.0.0.2", Name: "test2.dzone.example.com"}}

		receiveUpdateObject := NewZoneDelegated(ZoneDelegated{Ref: fakeRefReturn, DelegateTo: delegateTo})
		returnUpdateObject := NewZoneDelegated(ZoneDelegated{DelegateTo: delegateTo, Ref: fakeRefReturn})
		zdFakeConnector := &fakeConnector{
			fakeRefReturn:   fakeRefReturn,
			resultObject:    returnUpdateObject,
			updateObjectObj: receiveUpdateObject,
			updateObjectRef: fakeRefReturn,
		}

		objMgr := NewObjectManager(zdFakeConnector, cmpType, tenantID)

		var updatedObject *ZoneDelegated
		var err error
		It("should pass expected updated object to UpdateObject", func() {
			updatedObject, err = objMgr.UpdateZoneDelegated(fakeRefReturn, delegateTo)
		})
		It("should update zone with new delegation server list with no error", func() {
			Expect(updatedObject).To(Equal(returnUpdateObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete ZoneDelegated", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		deleteRef := "zone_delegated/ZG5zLnpvbmUkLl9kZWZhdWx0LnphLmNvLmFic2EuY2Fhcy5vaG15Z2xiLmdzbGJpYmNsaWVudA:dzone.example.com/default"
		fakeRefReturn := deleteRef
		zdFakeConnector := &fakeConnector{
			deleteObjectRef: deleteRef,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(zdFakeConnector, cmpType, tenantID)

		var actualRef string
		var err error
		It("should pass expected ZoneDelegated Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteZoneDelegated(deleteRef)
		})
		It("should return expected ZoneDelegated Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})

})
