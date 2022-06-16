package e2e_tests

import (
	ibclient "github.com/infobloxopen/infoblox-go-client/v2"
	. "github.com/onsi/ginkgo/v2"
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

	It("Should get the Admin Group [admin-group]", func() {
		var res []ibclient.Admingroup
		search := &ibclient.Admingroup{}
		qp := ibclient.NewQueryParams(false, map[string]string{"name": "admin-group"})
		err := connector.GetObject(search, "", qp, &res)
		Expect(err).To(BeNil())
		Expect(res[0].Ref).To(MatchRegexp("^admingroup.*admin-group$"))
		Expect(res[0].Name).To(Equal("admin-group"))
	})

	It("Should get the AllRecords without search fields (N)", func() {
		var res []ibclient.Allrecords
		search := &ibclient.Allrecords{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).NotTo(BeNil())
		// TODO Check the error string
	})

	PIt("Should get the CSV Import task", func() {
		var res []ibclient.Csvimporttask
		search := &ibclient.Csvimporttask{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(BeNil())
	})

	It("Should get the Discovery object (N)", func() {
		var res []ibclient.Discovery
		search := &ibclient.Discovery{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).NotTo(BeNil())
		// TODO Check the error string
	})

	PIt("Should get the Discovery Device object (N)", func() {
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

	It("Should get the Discovery Device Neighbor object", func() {
		var res []ibclient.DiscoveryDeviceneighbor
		search := &ibclient.DiscoveryDeviceneighbor{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError("not found"))
	})

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

	It("Should get the DTC Certificate object", func() {
		var res []ibclient.DtcCertificate
		search := &ibclient.DtcCertificate{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError("not found"))
	})

	It("Should get the DTC LBDN object", func() {
		var res []ibclient.DtcLbdn
		search := &ibclient.DtcLbdn{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError("not found"))
	})

	It("Should get the DTC monitor object", func() {
		var res []ibclient.DtcMonitor
		search := &ibclient.DtcMonitor{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(BeNil())
		Expect(res[0].Comment).To(Equal("Default ICMP health monitor"))
		Expect(res[0].Name).To(Equal("icmp"))
		Expect(res[0].Ref).To(MatchRegexp("^dtc:monitor.*icmp$"))
		Expect(res[0].Type).To(Equal("ICMP"))

		Expect(res[1].Comment).To(Equal("Default HTTP health monitor"))
		Expect(res[1].Name).To(Equal("http"))
		Expect(res[1].Ref).To(Equal("dtc:monitor/Li5sYl9oZWFsdGhfbW9uaXRvciQx:http"))
		Expect(res[1].Type).To(Equal("HTTP"))

		Expect(res[2].Comment).To(Equal("Default HTTPS health monitor"))
		Expect(res[2].Name).To(Equal("https"))
		Expect(res[2].Ref).To(Equal("dtc:monitor/Li5sYl9oZWFsdGhfbW9uaXRvciQy:https"))
		Expect(res[2].Type).To(Equal("HTTP"))

		Expect(res[3].Comment).To(Equal("Default SIP health monitor"))
		Expect(res[3].Name).To(Equal("sip"))
		Expect(res[3].Ref).To(Equal("dtc:monitor/Li5sYl9oZWFsdGhfbW9uaXRvciQz:sip"))
		Expect(res[3].Type).To(Equal("SIP"))

		Expect(res[4].Comment).To(Equal("Default PDP health monitor"))
		Expect(res[4].Name).To(Equal("pdp"))
		Expect(res[4].Ref).To(Equal("dtc:monitor/Li5sYl9oZWFsdGhfbW9uaXRvciQ0:pdp"))
		Expect(res[4].Type).To(Equal("PDP"))
	})

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

	It("Should get the DTC TCP monitor object", func() {
		var res []ibclient.DtcMonitorTcp
		search := &ibclient.DtcMonitorTcp{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError("not found"))
	})

	It("Should get the DTC object", func() {
		var res []ibclient.DtcObject
		search := &ibclient.DtcObject{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError("not found"))
	})

	It("Should get the DTC POOL object", func() {
		var res []ibclient.DtcPool
		search := &ibclient.DtcPool{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError("not found"))
	})

	It("Should get the DTC Server object", func() {
		var res []ibclient.DtcServer
		search := &ibclient.DtcServer{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError("not found"))
	})

	It("Should get the DTC Topology Rule object", func() {
		var res []ibclient.DtcTopology
		search := &ibclient.DtcTopology{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError("not found"))
	})

	It("Should get the Extensible Attribute Definition object", func() {
		var res []ibclient.EADefinition
		search := &ibclient.EADefinition{}
		qp := ibclient.NewQueryParams(false, map[string]string{"name": "Site"})
		err := connector.GetObject(search, "", qp, &res)
		Expect(err).To(BeNil())
		Expect(res[0].Name).To(Equal("Site"))
		Expect(res[0].Ref).To(Equal("extensibleattributedef/b25lLmV4dGVuc2libGVfYXR0cmlidXRlc19kZWYkLlNpdGU:Site"))
		Expect(res[0].Type).To(Equal("STRING"))
	})

	It("Should not get the File operations object", func() {
		var res []ibclient.Fileop
		search := &ibclient.Fileop{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).NotTo(BeNil())
		// TODO Check the error string
	})

	It("Should get the Grid Cloud API object", func() {
		var res []ibclient.GridCloudapi
		search := &ibclient.GridCloudapi{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(BeNil())
		Expect(res[0].EnableRecycleBin).To(Equal(true))
		Expect(res[0].AllowedApiAdmins).To(Equal([]*ibclient.GridCloudapiUser{}))
		Expect(res[0].AllowApiAdmins).To(Equal("ALL"))
		Expect(res[0].Ref).To(Equal("grid:cloudapi/b25lLnZjb25uZWN0b3JfY2x1c3RlciQw:grid"))
	})

	It("Should get the Grid Cloud Statistics object", func() {
		var res []ibclient.GridCloudapiCloudstatistics
		search := &ibclient.GridCloudapiCloudstatistics{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(BeNil())
		Expect(res[0].TenantCount).To(Equal(uint32(0)))
		Expect(res[0].AllocatedIpCount).To(Equal(uint32(0)))
		Expect(res[0].FloatingIpCount).To(Equal(uint32(0)))
		Expect(res[0].TenantVmCount).To(Equal(uint32(0)))
		Expect(res[0].AvailableIpCount).To(Equal("0"))
		Expect(res[0].TenantIpCount).To(Equal(uint32(0)))
		Expect(res[0].AllocatedAvailableRatio).To(Equal(uint32(0)))
		Expect(res[0].FixedIpCount).To(Equal(uint32(0)))
	})

	It("Should get the Grid Cloud API Tenant object", func() {
		var res []ibclient.GridCloudapiTenant
		search := &ibclient.GridCloudapiTenant{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError("not found"))
	})

	It("Should get the Grid Cloud API VM address object", func() {
		var res []ibclient.GridCloudapiVmaddress
		search := &ibclient.GridCloudapiVmaddress{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError("not found"))
	})

	It("Should get the Grid DHCP properties object", func() {
		var res []ibclient.GridDhcpproperties
		search := &ibclient.GridDhcpproperties{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(BeNil())
		Expect(res[0].DisableAllNacFilters).To(Equal(false))
		Expect(res[0].Ref).To(Equal("grid:dhcpproperties/ZG5zLmNsdXN0ZXJfZGhjcF9wcm9wZXJ0aWVzJDA:Infoblox"))
	})

	It("Should get the Grid Dns object", func() {
		var res []ibclient.GridDns
		search := &ibclient.GridDns{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(BeNil())
		Expect(res[0].Ref).To(Equal("grid:dns/ZG5zLmNsdXN0ZXJfZG5zX3Byb3BlcnRpZXMkMA:Infoblox"))
	})

	It("Should get the Grid Dns object", func() {
		var res []ibclient.GridDns
		search := &ibclient.GridDns{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(BeNil())
		Expect(res[0].Ref).To(Equal("grid:dns/ZG5zLmNsdXN0ZXJfZG5zX3Byb3BlcnRpZXMkMA:Infoblox"))
	})

	PIt("Should get the MaxMind DB Info object", func() {
		var res []ibclient.GridMaxminddbinfo
		search := &ibclient.GridMaxminddbinfo{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(BeNil())
		Expect(res[0].BinaryMajorVersion).To(Equal(uint32(2)))
		Expect(res[0].Ref).To(Equal("grid:maxminddbinfo/ZG5zLm1heG1pbmRfZGJfaW5mbyQw:maxminddbinfo"))
	})

	It("Should get the Member Cloud API object", func() {
		var res []ibclient.GridMemberCloudapi
		search := &ibclient.GridMemberCloudapi{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError("not found"))
	})

	It("Should get the x509certificate object", func() {
		var res []ibclient.GridX509certificate
		search := &ibclient.GridX509certificate{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).NotTo(BeNil())
		Expect(err).To(MatchError("not found"))
	})

	It("Should add Network View [dhcpview]", func() {
		nv := &ibclient.NetworkView{
			Name:    "dhcpview",
			Comment: "wapi added",
		}
		_, err := connector.CreateObject(nv)
		Expect(err).To(BeNil())
	})

	It("Should add FMZ Auth Zone [wapi.com]", func() {
		nv := &ibclient.ZoneAuth{
			View: "default",
			Fqdn: "wapi.com",
		}
		ref, err := connector.CreateObject(nv)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("zone_auth.*wapi.com/default"))
	})

	When("Auth Zone [wapi.com] exists", func() {
		BeforeEach(func() {
			nv := &ibclient.ZoneAuth{
				View: "default",
				Fqdn: "wapi.com",
			}
			_, err := connector.CreateObject(nv)
			Expect(err).To(BeNil())
		})

		It("Should add CNAME Record [cname.wapi.com]", func() {
			r := &ibclient.RecordCNAME{
				View:      "default",
				Name:      "cname.wapi.com",
				Canonical: "qatest1.com",
				Comment:   "verified",
			}
			ref, err := connector.CreateObject(r)
			Expect(err).To(BeNil())
			Expect(ref).To(MatchRegexp("^record:cname.*cname\\.wapi\\.com/default$"))
		})

		It("Should add A Record [a.wapi.com]", func() {
			r := &ibclient.RecordA{
				Name:     "a.wapi.com",
				Ipv4Addr: "9.0.0.1",
				Comment:  "Added A Record",
				Disable:  false,
				Ttl:      uint32(10),
				UseTtl:   true,
				View:     "default",
			}
			ref, err := connector.CreateObject(r)
			Expect(err).To(BeNil())
			Expect(ref).To(MatchRegexp("^record:a.*a\\.wapi\\.com/default$"))
		})

		It("Should add AAAA Record [aaaa.wapi.com]", func() {
			r := &ibclient.RecordAAAA{
				Name:     "aaaa.wapi.com",
				Ipv6Addr: "99::99",
				Comment:  "Added AAAA Record",
				Disable:  false,
				Ttl:      uint32(10),
				UseTtl:   true,
				View:     "default",
			}
			ref, err := connector.CreateObject(r)
			Expect(err).To(BeNil())
			Expect(ref).To(MatchRegexp("^record:aaaa.*aaaa\\.wapi\\.com/default$"))
		})

		It("Should add AAAA Record [aaaa.wapi.com]", func() {
			r := &ibclient.RecordAAAA{
				Name:     "aaaa.wapi.com",
				Ipv6Addr: "99::99",
				Comment:  "Added AAAA Record",
				Disable:  false,
				Ttl:      uint32(10),
				UseTtl:   true,
				View:     "default",
			}
			ref, err := connector.CreateObject(r)
			Expect(err).To(BeNil())
			Expect(ref).To(MatchRegexp("^record:aaaa.*aaaa\\.wapi\\.com/default$"))
		})

		It("Should add Host record [h1.wapi.com] with both ipv4addrs and ipv6addrs fields", func() {
			r := &ibclient.HostRecord{
				Name: "h1.wapi.com",
				View: "default",
				Ipv4Addrs: []ibclient.HostRecordIpv4Addr{
					{EnableDhcp: false, Ipv4Addr: "20.20.20.20"},
					{Ipv4Addr: "20.20.20.30"},
					{EnableDhcp: false, Ipv4Addr: "20.20.20.40"},
				},
				Ipv6Addrs: []ibclient.HostRecordIpv6Addr{
					{EnableDhcp: true, Ipv6Addr: "2000::1", Duid: "11:10"},
					{Ipv6Addr: "2000::2"},
				},
			}
			ref, err := connector.CreateObject(r)
			Expect(err).To(BeNil())
			Expect(ref).To(MatchRegexp("^record:host.*h1\\.wapi\\.com/default$"))
		})

		It("Should add MX Record [mx.wapi.com]", func() {
			r := &ibclient.RecordMx{
				Name:          "mx.wapi.com",
				MailExchanger: "wapi.com",
				Preference:    uint32(10),
				Comment:       "Creating mx record through infoblox-go-client",
				Disable:       false,
				Ttl:           uint32(20),
				UseTtl:        true,
			}
			ref, err := connector.CreateObject(r)
			Expect(err).To(BeNil())
			Expect(ref).To(MatchRegexp("^record:mx.*mx\\.wapi\\.com/default$"))
		})

		It("Should add TXT Record [txt.wapi.com]", func() {
			r := &ibclient.RecordTXT{
				Name:    "txt.wapi.com",
				Text:    "wapi.com",
				Comment: "Creating txt record through infoblox-go-client",
				Disable: false,
				Ttl:     uint32(20),
				UseTtl:  true,
				View:    "default",
			}
			ref, err := connector.CreateObject(r)
			Expect(err).To(BeNil())
			Expect(ref).To(MatchRegexp("^record:txt.*txt\\.wapi\\.com/default$"))
		})

		It("Should add PTR Record [ptr1.wapi.com]", func() {
			r := &ibclient.RecordPTR{
				Name:     "ptr1.wapi.com",
				PtrdName: "ptr.wapi.com",
				View:     "default",
				Comment:  "wapi added",
				Disable:  false,
				Ttl:      uint32(10),
				UseTtl:   true,
			}
			ref, err := connector.CreateObject(r)
			Expect(err).To(BeNil())
			Expect(ref).To(MatchRegexp("^record:ptr.*ptr1\\.wapi\\.com/default$"))
		})

		It("Should add SRV Record with all the attributes, using free format [name = srv.wapi.com]", func() {
			r := &ibclient.RecordSrv{
				Name:     "srv.wapi.com",
				Weight:   uint32(10),
				Priority: uint32(10),
				Port:     uint32(10),
				Target:   "srv.wapi.com",
				Comment:  "wapi added",
				Disable:  false,
				Ttl:      10,
				UseTtl:   true,
				View:     "default",
			}
			ref, err := connector.CreateObject(r)
			Expect(err).To(BeNil())
			Expect(ref).To(MatchRegexp("^record:srv.*srv\\.wapi\\.com/default$"))
		})
	})

	When("Network View [dhcpview] exists", func() {
		BeforeEach(func() {
			nv := &ibclient.NetworkView{
				Name:    "dhcpview",
				Comment: "wapi added",
			}
			_, err := connector.CreateObject(nv)
			Expect(err).To(BeNil())
		})

		It("Should add Network [92.0.0.0/8] in custom network view [dhcpview]", func() {
			r := &ibclient.Ipv4Network{
				Comment:     "Add ipv4network through WAPI",
				Network:     "92.0.0.0/8",
				NetworkView: "dhcpview",
			}
			ref, err := connector.CreateObject(r)
			Expect(err).To(BeNil())
			Expect(ref).To(MatchRegexp("^network.*92\\.0\\.0\\.0/8/dhcpview$"))
		})

		When("Network [92.0.0.0/8] in custom network view [dhcpview] exists", func() {
			BeforeEach(func() {
				r := &ibclient.Ipv4Network{
					Comment:     "Add ipv4network through WAPI",
					Network:     "92.0.0.0/8",
					NetworkView: "dhcpview",
				}
				ref, err := connector.CreateObject(r)
				Expect(err).To(BeNil())
				Expect(ref).To(MatchRegexp("^network.*92\\.0\\.0\\.0/8/dhcpview$"))
			})

			It("Should add IPv4 Range with mandatory fields start_addr [92.0.0.10] end_addr [92.0.0.20]", func() {
				r := &ibclient.Range{
					StartAddr:   "92.0.0.10",
					EndAddr:     "92.0.0.20",
					Comment:     "Add Range through WAPI",
					NetworkView: "dhcpview",
				}
				ref, err := connector.CreateObject(r)
				Expect(err).To(BeNil())
				Expect(ref).To(MatchRegexp("^range.*92\\.0\\.0\\.10/92\\.0\\.0\\.20/dhcpview$"))
			})

			It("Should add IPv4 fixed address [92.0.0.2] and mac [11:11:11:11:11:15]", func() {
				fa := &ibclient.Ipv4FixedAddress{
					Name:        "wapi-fa1",
					Ipv4Addr:    "92.0.0.2",
					NetworkView: "dhcpview",
					Mac:         "11:11:11:11:11:15",
					Comment:     "HellO",
				}
				ref, err := connector.CreateObject(fa)
				Expect(err).To(BeNil())
				Expect(ref).To(MatchRegexp("^fixedaddress.*92\\.0\\.0\\.2/dhcpview"))
			})
		})
	})

	It("Should add IPv6 Network [1::/16]", func() {
		n := &ibclient.Ipv6Network{
			AutoCreateReversezone: false,
			Comment:               "Add ipv6network through WAPI",
			Network:               "1::/16",
			NetworkView:           "default",
		}
		ref, err := connector.CreateObject(n)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^ipv6network.*1%3A%3A/16/default$"))
	})

	When("IPv6 Network [1::/16] exists", func() {
		BeforeEach(func() {
			n := &ibclient.Ipv6Network{
				AutoCreateReversezone: false,
				Comment:               "Add ipv6network through WAPI",
				Network:               "1::/16",
				NetworkView:           "default",
			}
			ref, err := connector.CreateObject(n)
			Expect(err).To(BeNil())
			Expect(ref).To(MatchRegexp("^ipv6network.*1%3A%3A/16/default$"))
		})

		It("Should add IPv6 Range [start_addr = 1::1; end_addr = 1::20]", func() {
			r := &ibclient.Ipv6range{
				StartAddr: "1::1",
				EndAddr:   "1::20",
				Network:   "1::/16",
				Comment:   "Add Range through WAPI",
			}
			ref, err := connector.CreateObject(r)
			Expect(err).To(BeNil())
			Expect(ref).To(MatchRegexp("^ipv6range.*1%3A%3A1/1%3A%3A20/default$"))
		})

		It("Should add ipv6fixedaddress [1::50]", func() {
			fa := &ibclient.Ipv6FixedAddress{
				Duid:        "ab:34:56:78:90",
				NetworkView: "default",
				Ipv6Addr:    "1::50",
			}
			ref, err := connector.CreateObject(fa)
			Expect(err).To(BeNil())
			Expect(ref).To(MatchRegexp("^ipv6fixedaddress.*1%3A%3A50/default"))
		})

		When("IPv6 Range [start_addr = 1::1; end_addr = 1::20] exits", func() {
			BeforeEach(func() {
				r := &ibclient.Ipv6range{
					StartAddr: "1::1",
					EndAddr:   "1::20",
					Network:   "1::/16",
					Comment:   "Add Range through WAPI",
				}
				ref, err := connector.CreateObject(r)
				Expect(err).To(BeNil())
				Expect(ref).To(MatchRegexp("^ipv6range.*1%3A%3A1/1%3A%3A20/default$"))
			})

			PIt("Should get the IPAM IPv6Address object", func() {
				var res []ibclient.Ipv6address
				search := &ibclient.Ipv6address{}
				qp := ibclient.NewQueryParams(false, map[string]string{"ip_address": "1::1"})
				err := connector.GetObject(search, "", qp, &res)
				Expect(err).To(BeNil())
				Expect(res[0].Status).To(Equal("UNUSED"))
				Expect(res[0].Usage).To(HaveLen(0))
				Expect(res[0].Network).To(Equal("1::/16"))
				Expect(res[0].Duid).To(Equal(""))
				Expect(res[0].LeaseState).To(Equal("FREE"))
				Expect(res[0].NetworkView).To(Equal("default"))
				Expect(res[0].Objects).To(HaveLen(0))
				Expect(res[0].IsConflict).To(Equal(false))
				Expect(res[0].Types).To(HaveLen(1))
				Expect(res[0].Types[0]).To(Equal("RESERVED_RANGE"))
				Expect(res[0].Ref).To(Equal("ipv6address/Li5pcHY2X2FkZHJlc3MkMTo6MS8w:1%3A%3A1"))
				Expect(res[0].IpAddress).To(Equal("1::1"))
				Expect(res[0].Names).To(HaveLen(0))
			})
		})

		When("ipv6fixedaddress [1::50] exists", func() {
			BeforeEach(func() {
				fa := &ibclient.Ipv6FixedAddress{
					Duid:        "ab:34:56:78:90",
					NetworkView: "default",
					Ipv6Addr:    "1::50",
				}
				ref, err := connector.CreateObject(fa)
				Expect(err).To(BeNil())
				Expect(ref).To(MatchRegexp("^ipv6fixedaddress.*1%3A%3A50/default"))
			})

			It("Should get ipv6fixedaddress with default return fields", func() {
				var res []ibclient.Ipv6FixedAddress
				search := &ibclient.Ipv6FixedAddress{}
				qp := ibclient.NewQueryParams(false, map[string]string{"ipv6addr": "1::50"})
				err := connector.GetObject(search, "", qp, &res)
				Expect(err).To(BeNil())
				Expect(res[0].Duid).To(Equal("ab:34:56:78:90"))
				Expect(res[0].NetworkView).To(Equal("default"))
				Expect(res[0].Ref).To(MatchRegexp("ipv6fixedaddress.*1%3A%3A50/default$"))
				Expect(res[0].Ipv6Addr).To(Equal("1::50"))
			})
		})
	})

	It("Should add Network Container [78.0.0.0/8]", func() {
		nc := &ibclient.Ipv4NetworkContainer{
			AutoCreateReversezone: false,
			Comment:               "Add networkcontainer through WAPI",
			Network:               "78.0.0.0/8",
			NetworkView:           "default",
		}
		ref, err := connector.CreateObject(nc)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^networkcontainer.*78\\.0\\.0\\.0/8/default$"))
	})

	When("Network Container [78.0.0.0/8] exits", func() {
		BeforeEach(func() {
			nc := &ibclient.Ipv4NetworkContainer{
				AutoCreateReversezone: false,
				Comment:               "Add networkcontainer through WAPI",
				Network:               "78.0.0.0/8",
				NetworkView:           "default",
			}
			ref, err := connector.CreateObject(nc)
			Expect(err).To(BeNil())
			Expect(ref).To(MatchRegexp("^networkcontainer.*78\\.0\\.0\\.0/8/default$"))
		})

		It("Should add Network [78.0.0.0/30] to the Network Container [78.0.0.0/8]", func() {
			n := &ibclient.Ipv4Network{
				Comment:     "Add a network to container 78.0.0.0/8 through WAPI",
				Network:     "78.0.0.0/30",
				NetworkView: "default",
			}
			ref, err := connector.CreateObject(n)
			Expect(err).To(BeNil())
			Expect(ref).To(MatchRegexp("^network.*78\\.0\\.0\\.0/30/default$"))
		})

		When("Network [78.0.0.0/30] exists", func() {
			BeforeEach(func() {
				n := &ibclient.Ipv4Network{
					Comment:     "Add a network to container 78.0.0.0/8 through WAPI",
					Network:     "78.0.0.0/30",
					NetworkView: "default",
				}
				ref, err := connector.CreateObject(n)
				Expect(err).To(BeNil())
				Expect(ref).To(MatchRegexp("^network.*78\\.0\\.0\\.0/30/default$"))
			})

			PIt("Should get the IPAM IPv4Address object", func() {
				var res []ibclient.Ipv4address
				search := &ibclient.Ipv4address{}
				qp := ibclient.NewQueryParams(false, map[string]string{"ip_address": "78.0.0.1"})
				err := connector.GetObject(search, "", qp, &res)
				Expect(err).To(BeNil())
				Expect(res[0].Status).To(Equal("UNUSED"))
				Expect(res[0].Network).To(Equal("78.0.0.0/30"))
				Expect(res[0].NetworkView).To(Equal("default"))
				Expect(res[0].Usage).To(HaveLen(0))
				Expect(res[0].Objects).To(HaveLen(0))
				Expect(res[0].IsConflict).To(Equal(false))
				Expect(res[0].MacAddress).To(Equal(""))
				Expect(res[0].Types).To(HaveLen(0))
				Expect(res[0].Ref).To(Equal("ipv4address/Li5pcHY0X2FkZHJlc3MkOTIuMC4wLjEvMA:78.0.0.1"))
				Expect(res[0].IpAddress).To(Equal("78.0.0.1"))
				Expect(res[0].Names).To(HaveLen(0))
			})
		})
	})

	It("Should add IPv6 Network Container [2000::/64]", func() {
		n := &ibclient.Ipv6Network{
			AutoCreateReversezone: false,
			Comment:               "Add ipv6networkcontainer through WAPI",
			Network:               "2000::/64",
			NetworkView:           "default",
		}
		ref, err := connector.CreateObject(n)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^ipv6network.*2000%3A%3A/64/default"))
	})

	It("Should add DNS view [view1]", func() {
		v := &ibclient.View{
			Name: "view1",
		}
		ref, err := connector.CreateObject(v)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^view.*view1/false$"))
	})

	It("Should add Named ACL without 'access_list' field [wapi-na2]", func() {
		v := &ibclient.Namedacl{
			Name:    "wapi-na2",
			Comment: "No acls present",
		}
		ref, err := connector.CreateObject(v)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^namedacl.*wapi-na2$"))

		var res []ibclient.Namedacl
		search := &ibclient.Namedacl{}
		err = connector.GetObject(search, "", nil, &res)
		Expect(err).To(BeNil())
		Expect(res[0].Comment).To(Equal("No acls present"))
		Expect(res[0].Ref).To(MatchRegexp("wapi-na2"))
		Expect(res[0].Name).To(Equal("wapi-na2"))
	})

	It("Should get the DHCP network template object", func() {
		var res []ibclient.Networktemplate
		search := &ibclient.Networktemplate{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError("not found"))
	})

	It("Should get the DHCP IPv6 network template object", func() {
		var res []ibclient.Ipv6networktemplate
		search := &ibclient.Ipv6networktemplate{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError("not found"))
	})

	It("Should get the DHCP IPv6 Shared Network object", func() {
		var res []ibclient.Ipv6sharednetwork
		search := &ibclient.Ipv6sharednetwork{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError("not found"))
	})

	It("Should get the DHCP Lease object", func() {
		var res []ibclient.Lease
		search := &ibclient.Lease{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError("not found"))
	})

	It("Should get the MAC Filter Address object", func() {
		var res []ibclient.Macfilteraddress
		search := &ibclient.Macfilteraddress{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError("not found"))
	})

	It("Should get the Member DHCP properties object", func() {
		var res []ibclient.MemberDhcpproperties
		search := &ibclient.MemberDhcpproperties{}
		search.SetReturnFields([]string{"host_name"})
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(BeNil())
		Expect(res[0].HostName).To(Equal("infoblox.localdomain"))
		Expect(res[0].Ref).To(MatchRegexp("^member:dhcpproperties.*infoblox.localdomain$"))
	})

	It("Should get the Member DNS object", func() {
		var res []ibclient.MemberDns
		search := &ibclient.MemberDns{}
		search.SetReturnFields([]string{"host_name"})
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(BeNil())
		Expect(res[0].HostName).To(Equal("infoblox.localdomain"))
		Expect(res[0].Ref).To(MatchRegexp("^member:dns.*infoblox.localdomain$"))
	})

	It("Should get the Active Directory Domain object", func() {
		var res []ibclient.MsserverAdsitesDomain
		search := &ibclient.MsserverAdsitesDomain{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError("not found"))
	})

	It("Should get the Active Directory Site object", func() {
		var res []ibclient.MsserverAdsitesSite
		search := &ibclient.MsserverAdsitesSite{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError("not found"))
	})

	It("Should get the Permissions object", func() {
		var res []ibclient.Permission
		search := &ibclient.Permission{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(BeNil())
		Expect(res[0].Permission).To(Equal("WRITE"))
		Expect(res[0].Role).To(Equal("DNS Admin"))
		Expect(res[0].Ref).To(MatchRegexp("^permission.*Admin/WRITE$"))
		Expect(res[0].ResourceType).To(Equal("VIEW"))
	})

})
