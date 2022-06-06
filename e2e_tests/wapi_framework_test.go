package e2e_tests

import (
	ibclient "github.com/infobloxopen/infoblox-go-client/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Go Client", func() {
	var connector *ConnectorFacadeE2E

	BeforeEach(func() {
		hostConfig := ibclient.HostConfig{
			Host:    os.Getenv("INFOBLOX_SERVER"),
			Version: os.Getenv("WAPI_VERSION"),
			Port:    os.Getenv("PORT"),
		}

		authConfig := ibclient.AuthConfig{
			Username: os.Getenv("INFOBLOX_USERNAME"),
			Password: os.Getenv("INFOBLOX_PASSWORD"),
		}

		transportConfig := ibclient.NewTransportConfig("false", 20, 10)
		requestBuilder := &ibclient.WapiRequestBuilder{}
		requestor := &ibclient.WapiHttpRequestor{}
		ibclientConnector, err := ibclient.NewConnector(hostConfig, authConfig, transportConfig, requestBuilder, requestor)
		Expect(err).To(BeNil())
		connector = &ConnectorFacadeE2E{*ibclientConnector, make([]string, 0)}
	})

	AfterEach(func() {
		err := connector.SweepObjects()
		Expect(err).To(BeNil())
	})

	It("Should get the Grid object", func() {
		var res []ibclient.Grid
		search := &ibclient.Grid{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(BeNil())
		Expect(res[0].Ref).To(Equal("grid/b25lLmNsdXN0ZXIkMA:Infoblox"))
	})

	It("Should get the Member object", func() {
		var res []ibclient.Member
		search := &ibclient.Member{}
		search.SetReturnFields([]string{"host_name"})
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(BeNil())
		Expect(res[0].Ref).To(Equal("member/b25lLnZpcnR1YWxfbm9kZSQw:infoblox.localdomain"))
		Expect(res[0].HostName).To(Equal("infoblox.localdomain"))
	})

	It("Should get the Admin User [admin]", func() {
		var res []ibclient.Adminuser
		search := &ibclient.Adminuser{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(BeNil())
		Expect(res[0].Ref).To(Equal("adminuser/b25lLmFkbWluJGFkbWlu:admin"))
		Expect(res[0].AdminGroups[0]).To(Equal("admin-group"))
	})

	//It("Should get the Admin Group [admin-group]", func() {
	//	var res []ibclient.Admingroup
	//	search := &ibclient.Admingroup{}
	//	err := connector.GetObject(search, "", nil, &res)
	//	Expect(err).To(BeNil())
	//	Expect(res[0].Ref).To(Equal("admingroup/b25lLmFkbWluX2dyb3VwJC5hZG1pbi1ncm91cA:admin-group"))
	//	Expect(res[0].Name).To(Equal("admin-group"))
	//})

	It("Should get the AllRecords without search fields (N)", func() {
		var res []ibclient.Allrecords
		search := &ibclient.Allrecords{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).NotTo(BeNil())
		// TODO Check the error string
	})

	//It("Should get the CSV Import task", func() {
	//	var res []ibclient.Csvimporttask
	//	search := &ibclient.Csvimporttask{}
	//	err := connector.GetObject(search, "", nil, &res)
	//	Expect(err).To(BeNil())
	//})

	It("Should get the Discovery object (N)", func() {
		var res []ibclient.Discovery
		search := &ibclient.Discovery{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).NotTo(BeNil())
		// TODO Check the error string
	})

	It("Should get the Discovery Device object (N)", func() {
		var res []ibclient.DiscoveryDevice
		search := &ibclient.DiscoveryDevice{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).NotTo(BeNil())
		// TODO Check the error string
	})

	It("Should get the Discovery Device Component object (N)", func() {
		var res []ibclient.DiscoveryDevicecomponent
		search := &ibclient.DiscoveryDevicecomponent{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).NotTo(BeNil())
		// TODO Check the error string
	})

	It("Should get the Discovery Device Interface object (N)", func() {
		var res []ibclient.DiscoveryDeviceinterface
		search := &ibclient.DiscoveryDeviceinterface{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).NotTo(BeNil())
		// TODO Check the error string
	})

	//It("Should get the Discovery Device Neighbor object", func() {
	//	var res []ibclient.DiscoveryDeviceneighbor
	//	search := &ibclient.DiscoveryDeviceneighbor{}
	//	err := connector.GetObject(search, "", nil, &res)
	//	Expect(err).To(BeNil())
	//})

	It("Should get the Discovery Status object (N)", func() {
		var res []ibclient.DiscoveryStatus
		search := &ibclient.DiscoveryStatus{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).NotTo(BeNil())
		// TODO Check the error string
	})

	It("Should get the DTC (N)", func() {
		var res []ibclient.Dtc
		search := &ibclient.Dtc{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).NotTo(BeNil())
		// TODO Check the error string
	})

	//It("Should get the DTC Certificate object", func() {
	//	var res []ibclient.DtcCertificate
	//	search := &ibclient.DtcCertificate{}
	//	err := connector.GetObject(search, "", nil, &res)
	//	Expect(err).To(BeNil())
	//	Expect(len(res)).To(Equal(0))
	//})

	//It("Should get the DTC LBDN object", func() {
	//	var res []ibclient.DtcLbdn
	//	search := &ibclient.DtcLbdn{}
	//	err := connector.GetObject(search, "", nil, &res)
	//	Expect(err).To(BeNil())
	//	Expect(len(res)).To(Equal(0))
	//})

	//It("Should get the DTC monitor object", func() {
	//	var res []ibclient.DtcMonitor
	//	search := &ibclient.DtcMonitor{}
	//	err := connector.GetObject(search, "", nil, &res)
	//	Expect(err).To(BeNil())
	//	Expect(res[0].Comment).To(Equal("Default ICMP health monitor"))
	//	Expect(res[0].Name).To(Equal("icmp"))
	//	Expect(res[0].Ref).To(Equal("dtc:monitor/Li5sYl9oZWFsdGhfbW9uaXRvciQw:icmp"))
	//	Expect(res[0].Type).To(Equal("ICMP"))
	//
	//	Expect(res[1].Comment).To(Equal("Default HTTP health monitor"))
	//	Expect(res[1].Name).To(Equal("http"))
	//	Expect(res[1].Ref).To(Equal("dtc:monitor/Li5sYl9oZWFsdGhfbW9uaXRvciQx:http"))
	//	Expect(res[1].Type).To(Equal("HTTP"))
	//
	//	Expect(res[2].Comment).To(Equal("Default HTTPS health monitor"))
	//	Expect(res[2].Name).To(Equal("https"))
	//	Expect(res[2].Ref).To(Equal("dtc:monitor/Li5sYl9oZWFsdGhfbW9uaXRvciQy:https"))
	//	Expect(res[2].Type).To(Equal("HTTP"))
	//
	//	Expect(res[3].Comment).To(Equal("Default SIP health monitor"))
	//	Expect(res[3].Name).To(Equal("sip"))
	//	Expect(res[3].Ref).To(Equal("dtc:monitor/Li5sYl9oZWFsdGhfbW9uaXRvciQz:sip"))
	//	Expect(res[3].Type).To(Equal("SIP"))
	//
	//	Expect(res[4].Comment).To(Equal("Default PDP health monitor"))
	//	Expect(res[4].Name).To(Equal("pdp"))
	//	Expect(res[4].Ref).To(Equal("dtc:monitor/Li5sYl9oZWFsdGhfbW9uaXRvciQ0:pdp"))
	//	Expect(res[4].Type).To(Equal("PDP"))
	//})

	It("Should get the DTC HTTP monitor object", func() {
		var res []ibclient.DtcMonitorHttp
		search := &ibclient.DtcMonitorHttp{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(BeNil())
		Expect(res[0].Comment).To(Equal("Default HTTP health monitor"))
		Expect(res[0].Name).To(Equal("http"))
		Expect(res[0].Ref).To(Equal("dtc:monitor:http/ZG5zLmlkbnNfbW9uaXRvcl9odHRwJGh0dHA:http"))

		Expect(res[1].Comment).To(Equal("Default HTTPS health monitor"))
		Expect(res[1].Name).To(Equal("https"))
		Expect(res[1].Ref).To(Equal("dtc:monitor:http/ZG5zLmlkbnNfbW9uaXRvcl9odHRwJGh0dHBz:https"))
	})

	It("Should get the DTC ICMP monitor object", func() {
		var res []ibclient.DtcMonitorIcmp
		search := &ibclient.DtcMonitorIcmp{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(BeNil())
		Expect(res[0].Comment).To(Equal("Default ICMP health monitor"))
		Expect(res[0].Name).To(Equal("icmp"))
		Expect(res[0].Ref).To(Equal("dtc:monitor:icmp/ZG5zLmlkbnNfbW9uaXRvcl9pY21wJGljbXA:icmp"))
	})

	It("Should get the DTC PDP monitor object", func() {
		var res []ibclient.DtcMonitorPdp
		search := &ibclient.DtcMonitorPdp{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(BeNil())
		Expect(res[0].Comment).To(Equal("Default PDP health monitor"))
		Expect(res[0].Name).To(Equal("pdp"))
		Expect(res[0].Ref).To(Equal("dtc:monitor:pdp/ZG5zLmlkbnNfbW9uaXRvcl9wZHAkcGRw:pdp"))
	})

	It("Should get the DTC SIP monitor object", func() {
		var res []ibclient.DtcMonitorSip
		search := &ibclient.DtcMonitorSip{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(BeNil())
		Expect(res[0].Comment).To(Equal("Default SIP health monitor"))
		Expect(res[0].Name).To(Equal("sip"))
		Expect(res[0].Ref).To(Equal("dtc:monitor:sip/ZG5zLmlkbnNfbW9uaXRvcl9zaXAkc2lw:sip"))
	})

	//It("Should get the DTC TCP monitor object", func() {
	//	var res []ibclient.DtcMonitorTcp
	//	search := &ibclient.DtcMonitorTcp{}
	//	err := connector.GetObject(search, "", nil, &res)
	//	Expect(err).To(BeNil())
	//	Expect(len(res)).To(Equal(0))
	//})

	//It("Should get the DTC object", func() {
	//	var res []ibclient.DtcObject
	//	search := &ibclient.DtcObject{}
	//	err := connector.GetObject(search, "", nil, &res)
	//	Expect(err).To(BeNil())
	//	Expect(len(res)).To(Equal(0))
	//})
	//
	//It("Should get the DTC POOL object", func() {
	//	var res []ibclient.DtcPool
	//	search := &ibclient.DtcPool{}
	//	err := connector.GetObject(search, "", nil, &res)
	//	Expect(err).To(BeNil())
	//	Expect(len(res)).To(Equal(0))
	//})
	//
	//It("Should get the DTC Server object", func() {
	//	var res []ibclient.DtcServer
	//	search := &ibclient.DtcServer{}
	//	err := connector.GetObject(search, "", nil, &res)
	//	Expect(err).To(BeNil())
	//	Expect(len(res)).To(Equal(0))
	//})

	//It("Should get the DTC Topology Rule object", func() {
	//	var res []ibclient.DtcTopology
	//	search := &ibclient.DtcTopology{}
	//	err := connector.GetObject(search, "", nil, &res)
	//	Expect(err).To(BeNil())
	//	Expect(len(res)).To(Equal(0))
	//})

	It("Should get the Extensible Attribute Definition object", func() {
		var res []ibclient.EADefinition
		search := &ibclient.EADefinition{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(BeNil())
		Expect(res[0].Name).To(Equal("Site"))
		Expect(res[0].Ref).To(Equal("extensibleattributedef/b25lLmV4dGVuc2libGVfYXR0cmlidXRlc19kZWYkLlNpdGU:Site"))
		Expect(res[0].Type).To(Equal("STRING"))
	})

})
