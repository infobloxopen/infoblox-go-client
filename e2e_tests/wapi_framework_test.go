package e2e_tests

import (
	ibclient "github.com/infobloxopen/infoblox-go-client/v2"
	"github.com/infobloxopen/infoblox-go-client/v2/utils"
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

	It("Should get the Grid object", Label("ID: 1", "RO"), func() {
		var res []ibclient.Grid
		search := &ibclient.Grid{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(BeNil())
		Expect(res[0].Ref).To(Equal("grid/b25lLmNsdXN0ZXIkMA:Infoblox"))
	})

	It("Should get the Member object", Label("ID: 2", "RO"), func() {
		var res []ibclient.Member
		search := &ibclient.Member{}
		search.SetReturnFields([]string{"host_name"})
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(BeNil())
		Expect(res[0].Ref).To(HavePrefix("member/b25lLnZpcnR1YWxfbm9kZSQw:infoblox."))
		Expect(*res[0].HostName).To(HavePrefix("infoblox."))
	})

	It("Should get the Admin User [admin]", Label("ID: 3", "RO"), func() {
		var res []ibclient.Adminuser
		search := &ibclient.Adminuser{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(BeNil())
		Expect(res[0].Ref).To(Equal("adminuser/b25lLmFkbWluJGFkbWlu:admin"))
		Expect(res[0].AdminGroups[0]).To(Equal("admin-group"))
	})

	It("Should get the Admin Group [admin-group]", Label("ID: 4", "RO"), func() {
		var res []ibclient.Admingroup
		search := &ibclient.Admingroup{}
		qp := ibclient.NewQueryParams(false, map[string]string{"name": "admin-group"})
		err := connector.GetObject(search, "", qp, &res)
		Expect(err).To(BeNil())
		Expect(res[0].Ref).To(MatchRegexp("^admingroup.*admin-group$"))
		Expect(*res[0].Name).To(Equal("admin-group"))
	})

	It("Should get the AllRecords without search fields (N)", Label("ID: 5", "RO"), func() {
		var res []ibclient.Allrecords
		search := &ibclient.Allrecords{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).NotTo(BeNil())
		// TODO Check the error string
	})

	It("Should get the CSV Import task", Label("ID: 6", "RO"), func() {
		var res []ibclient.Csvimporttask
		search := &ibclient.Csvimporttask{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).NotTo(BeNil())
	})

	It("Should get the Discovery object (N)", Label("ID: 7", "RO"), func() {
		var res []ibclient.Discovery
		search := &ibclient.Discovery{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).NotTo(BeNil())
		// TODO Check the error string
	})

	It("Should get the Discovery Device object (N)", Label("ID: 8", "RO"), func() {
		var res []ibclient.DiscoveryDevice
		search := &ibclient.DiscoveryDevice{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).NotTo(BeNil())
		// TODO Check the error string
	})

	It("Should get the Discovery Device Component object (N)", Label("ID: 9", "RO"), func() {
		var res []ibclient.DiscoveryDevicecomponent
		search := &ibclient.DiscoveryDevicecomponent{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).NotTo(BeNil())
		// TODO Check the error string
	})

	It("Should get the Discovery Device Interface object (N)", Label("ID: 10", "RO"), func() {
		var res []ibclient.DiscoveryDeviceinterface
		search := &ibclient.DiscoveryDeviceinterface{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).NotTo(BeNil())
		// TODO Check the error string
	})

	It("Should get the Discovery Device Neighbor object", Label("ID: 11", "RO"), func() {
		var res []ibclient.DiscoveryDeviceneighbor
		search := &ibclient.DiscoveryDeviceneighbor{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError(ibclient.NewNotFoundError("requested object not found")))
	})

	It("Should get the Discovery Status object (N)", Label("ID: 12", "RO"), func() {
		var res []ibclient.DiscoveryStatus
		search := &ibclient.DiscoveryStatus{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).NotTo(BeNil())
		// TODO Check the error string
	})

	It("Should get the DTC (N)", Label("ID: 13", "RO"), func() {
		var res []ibclient.Dtc
		search := &ibclient.Dtc{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).NotTo(BeNil())
		// TODO Check the error string
	})

	It("Should get the DTC Certificate object", Label("ID: 14", "RO"), func() {
		var res []ibclient.DtcCertificate
		search := &ibclient.DtcCertificate{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError(ibclient.NewNotFoundError("requested object not found")))
	})

	It("Should get the DTC LBDN object", Label("ID: 15", "RO"), func() {
		var res []ibclient.DtcLbdn
		search := &ibclient.DtcLbdn{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError(ibclient.NewNotFoundError("requested object not found")))
	})

	It("Should get the DTC monitor object", Label("ID: 16", "RO"), func() {
		var res []ibclient.DtcMonitor
		search := &ibclient.DtcMonitor{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(BeNil())
		Expect(*res[0].Comment).To(Equal("Default ICMP health monitor"))
		Expect(*res[0].Name).To(Equal("icmp"))
		Expect(res[0].Ref).To(MatchRegexp("^dtc:monitor.*icmp$"))
		Expect(res[0].Type).To(Equal("ICMP"))

		Expect(*res[1].Comment).To(Equal("Default HTTP health monitor"))
		Expect(*res[1].Name).To(Equal("http"))
		Expect(res[1].Ref).To(MatchRegexp("^dtc:monitor.*http$"))
		Expect(res[1].Type).To(Equal("HTTP"))

		Expect(*res[2].Comment).To(Equal("Default HTTPS health monitor"))
		Expect(*res[2].Name).To(Equal("https"))
		Expect(res[2].Ref).To(MatchRegexp("^dtc:monitor.*https$"))
		Expect(res[2].Type).To(Equal("HTTP"))

		Expect(*res[3].Comment).To(Equal("Default SIP health monitor"))
		Expect(*res[3].Name).To(Equal("sip"))
		Expect(res[3].Ref).To(MatchRegexp("^dtc:monitor.*sip$"))
		Expect(res[3].Type).To(Equal("SIP"))

		Expect(*res[4].Comment).To(Equal("Default PDP health monitor"))
		Expect(*res[4].Name).To(Equal("pdp"))
		Expect(res[4].Ref).To(MatchRegexp("^dtc:monitor.*pdp$"))
		Expect(res[4].Type).To(Equal("PDP"))
	})

	It("Should get the DTC HTTP monitor object", Label("ID: 17", "RO"), func() {
		var res []ibclient.DtcMonitorHttp
		search := &ibclient.DtcMonitorHttp{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(BeNil())
		Expect(*res[0].Comment).To(Equal("Default HTTP health monitor"))
		Expect(*res[0].Name).To(Equal("http"))
		Expect(res[0].Ref).To(Equal("dtc:monitor:http/ZG5zLmlkbnNfbW9uaXRvcl9odHRwJGh0dHA:http"))

		Expect(*res[1].Comment).To(Equal("Default HTTPS health monitor"))
		Expect(*res[1].Name).To(Equal("https"))
		Expect(res[1].Ref).To(Equal("dtc:monitor:http/ZG5zLmlkbnNfbW9uaXRvcl9odHRwJGh0dHBz:https"))
	})

	It("Should get the DTC ICMP monitor object", Label("ID: 18", "RO"), func() {
		var res []ibclient.DtcMonitorIcmp
		search := &ibclient.DtcMonitorIcmp{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(BeNil())
		Expect(*res[0].Comment).To(Equal("Default ICMP health monitor"))
		Expect(*res[0].Name).To(Equal("icmp"))
		Expect(res[0].Ref).To(Equal("dtc:monitor:icmp/ZG5zLmlkbnNfbW9uaXRvcl9pY21wJGljbXA:icmp"))
	})

	It("Should get the DTC PDP monitor object", Label("ID: 19", "RO"), func() {
		var res []ibclient.DtcMonitorPdp
		search := &ibclient.DtcMonitorPdp{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(BeNil())
		Expect(*res[0].Comment).To(Equal("Default PDP health monitor"))
		Expect(*res[0].Name).To(Equal("pdp"))
		Expect(res[0].Ref).To(Equal("dtc:monitor:pdp/ZG5zLmlkbnNfbW9uaXRvcl9wZHAkcGRw:pdp"))
	})

	It("Should get the DTC SIP monitor object", Label("ID: 20", "RO"), func() {
		var res []ibclient.DtcMonitorSip
		search := &ibclient.DtcMonitorSip{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(BeNil())
		Expect(*res[0].Comment).To(Equal("Default SIP health monitor"))
		Expect(*res[0].Name).To(Equal("sip"))
		Expect(res[0].Ref).To(Equal("dtc:monitor:sip/ZG5zLmlkbnNfbW9uaXRvcl9zaXAkc2lw:sip"))
	})

	It("Should get the DTC TCP monitor object", Label("ID: 21", "RO"), func() {
		var res []ibclient.DtcMonitorTcp
		search := &ibclient.DtcMonitorTcp{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError(ibclient.NewNotFoundError("requested object not found")))
	})

	It("Should get the DTC object", Label("ID: 22", "RO"), func() {
		var res []ibclient.DtcObject
		search := &ibclient.DtcObject{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError(ibclient.NewNotFoundError("requested object not found")))
	})

	It("Should get the DTC POOL object", Label("ID: 23", "RO"), func() {
		var res []ibclient.DtcPool
		search := &ibclient.DtcPool{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError(ibclient.NewNotFoundError("requested object not found")))
	})

	It("Should get the DTC Server object", Label("ID: 24", "RO"), func() {
		var res []ibclient.DtcServer
		search := &ibclient.DtcServer{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError(ibclient.NewNotFoundError("requested object not found")))
	})

	It("Should get the DTC Topology object", Label("ID: 25", "RO"), func() {
		var res []ibclient.DtcTopology
		search := &ibclient.DtcTopology{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError(ibclient.NewNotFoundError("requested object not found")))
	})

	It("Should get the DTC Topology Rule object", Label("ID: 26", "RO"), func() {
		var res []ibclient.DtcTopologyRule
		search := &ibclient.DtcTopologyRule{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError(ibclient.NewNotFoundError("requested object not found")))
	})

	It("Should get the Extensible Attribute Definition object", Label("ID: 27", "RO"), func() {
		var res []ibclient.EADefinition
		search := &ibclient.EADefinition{}
		qp := ibclient.NewQueryParams(false, map[string]string{"name": "Site"})
		err := connector.GetObject(search, "", qp, &res)
		Expect(err).To(BeNil())
		Expect(*res[0].Name).To(Equal("Site"))
		Expect(res[0].Ref).To(Equal("extensibleattributedef/b25lLmV4dGVuc2libGVfYXR0cmlidXRlc19kZWYkLlNpdGU:Site"))
		Expect(res[0].Type).To(Equal("STRING"))
	})

	It("Should not get the File operations object", Label("ID: 28", "RO"), func() {
		var res []ibclient.Fileop
		search := &ibclient.Fileop{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).NotTo(BeNil())
		// TODO Check the error string
	})

	It("Should get the Grid Cloud API object", Label("ID: 29", "RO"), func() {
		var res []ibclient.GridCloudapi
		search := &ibclient.GridCloudapi{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(BeNil())
		Expect(*res[0].EnableRecycleBin).To(Equal(true))
		Expect(res[0].AllowedApiAdmins).To(Equal([]*ibclient.GridCloudapiUser{}))
		Expect(res[0].AllowApiAdmins).To(Equal("ALL"))
		Expect(res[0].Ref).To(Equal("grid:cloudapi/b25lLnZjb25uZWN0b3JfY2x1c3RlciQw:grid"))
	})

	It("Should get the Grid Cloud Statistics object", Label("ID: 30", "RO"), func() {
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

	It("Should get the Grid Cloud API Tenant object", Label("ID: 31", "RO"), func() {
		var res []ibclient.GridCloudapiTenant
		search := &ibclient.GridCloudapiTenant{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError(ibclient.NewNotFoundError("requested object not found")))
	})

	It("Should get the Grid Cloud API VM address object", Label("ID: 32", "RO"), func() {
		var res []ibclient.GridCloudapiVmaddress
		search := &ibclient.GridCloudapiVmaddress{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError(ibclient.NewNotFoundError("requested object not found")))
	})

	It("Should get the Grid DHCP properties object", Label("ID: 33", "RO"), func() {
		var res []ibclient.GridDhcpproperties
		search := &ibclient.GridDhcpproperties{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(BeNil())
		Expect(*res[0].DisableAllNacFilters).To(Equal(false))
		Expect(res[0].Ref).To(Equal("grid:dhcpproperties/ZG5zLmNsdXN0ZXJfZGhjcF9wcm9wZXJ0aWVzJDA:Infoblox"))
	})

	It("Should get the Grid Dns object", Label("ID: 34", "RO"), func() {
		var res []ibclient.GridDns
		search := &ibclient.GridDns{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(BeNil())
		Expect(res[0].Ref).To(Equal("grid:dns/ZG5zLmNsdXN0ZXJfZG5zX3Byb3BlcnRpZXMkMA:Infoblox"))
	})

	It("Should get the MaxMind DB Info object", Label("ID: 35", "RO"), func() {
		var res []ibclient.GridMaxminddbinfo
		search := &ibclient.GridMaxminddbinfo{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).NotTo(BeNil())
	})

	It("Should get the Member Cloud API object", Label("ID: 36", "RO"), func() {
		var res []ibclient.GridMemberCloudapi
		search := &ibclient.GridMemberCloudapi{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError(ibclient.NewNotFoundError("requested object not found")))
	})

	It("Should get the x509certificate object", Label("ID: 37", "RO"), func() {
		var res []ibclient.GridX509certificate
		search := &ibclient.GridX509certificate{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).NotTo(BeNil())
		Expect(err).To(MatchError(ibclient.NewNotFoundError("requested object not found")))
	})

	It("Should add Network View [dhcpview]", Label("ID: 38", "RW"), func() {
		nv := &ibclient.NetworkView{
			Name:    utils.StringPtr("dhcpview"),
			Comment: utils.StringPtr("wapi added"),
		}
		_, err := connector.CreateObject(nv)
		Expect(err).To(BeNil())
	})

	It("Should add FMZ Auth Zone [wapi.com]", Label("ID: 39", "RW"), func() {
		nv := &ibclient.ZoneAuth{
			View: utils.StringPtr("default"),
			Fqdn: "wapi.com",
		}
		ref, err := connector.CreateObject(nv)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("zone_auth.*wapi.com/default"))
	})

	When("Auth Zone [wapi.com] exists", Label("RW"), func() {
		var refZoneAuth string
		BeforeEach(func() {
			nv := &ibclient.ZoneAuth{
				View: utils.StringPtr("default"),
				Fqdn: "wapi.com",
			}
			var err error
			refZoneAuth, err = connector.CreateObject(nv)
			Expect(err).To(BeNil())

			ipv4Cidr := "16.12.1.0/24"
			netView := "default"
			ipv4Comment := "ipv4 network"
			ipv4Ea := ibclient.EA{"Site": "Namibia"}
			ipv4Network := ibclient.NewNetwork(netView, ipv4Cidr, false, ipv4Comment, ipv4Ea)
			_, err1 := connector.CreateObject(ipv4Network)
			Expect(err1).To(BeNil())

		})

		It("Should create a Zone delegation", func() {
			// Create a DNS Forward Zone
			zone := &ibclient.ZoneDelegated{
				Fqdn: "example1.wapi.com",
				DelegateTo: ibclient.NullableNameServers{
					NameServers: []ibclient.NameServer{
						{Name: "test", Address: "1.2.3.4"},
						{Name: "test2", Address: "2.3.4.5"},
					},
					IsNull: false,
				},
				DelegatedTtl:    utils.Uint32Ptr(3600),
				UseDelegatedTtl: utils.BoolPtr(true),
			}
			ref, err := connector.CreateObject(zone)
			Expect(err).To(BeNil())
			Expect(ref).To(MatchRegexp("^zone_delegated.*"))
		})
		It("Should get the Zone delegated", func() {
			zone := &ibclient.ZoneDelegated{
				Fqdn: "example2.wapi.com",
				DelegateTo: ibclient.NullableNameServers{
					NameServers: []ibclient.NameServer{
						{Name: "test", Address: "1.2.3.4"},
						{Name: "test2", Address: "2.3.4.5"},
					},
					IsNull: false,
				},
				Comment: utils.StringPtr("wapi added"),
			}
			ref, err := connector.CreateObject(zone)
			Expect(err).To(BeNil())
			Expect(ref).To(MatchRegexp("^zone_delegated.*"))

			var res []ibclient.ZoneDelegated
			search := &ibclient.ZoneDelegated{}
			errCode := connector.GetObject(search, "", nil, &res)
			Expect(errCode).To(BeNil())
			Expect(res[0].Ref).To(MatchRegexp("^zone_delegated.*"))
		})
		It("Should update the Zone delegated", func() {
			// Create a Zone delegated
			zoneCreate := &ibclient.ZoneDelegated{
				Fqdn: "example3.wapi.com",
				DelegateTo: ibclient.NullableNameServers{
					NameServers: []ibclient.NameServer{
						{Name: "test", Address: "1.2.3.4"},
						{Name: "test2", Address: "1.2.3.5"},
					},
					IsNull: false,
				},
				Comment: utils.StringPtr("wapi added"),
			}
			ref, errCode := connector.CreateObject(zoneCreate)
			Expect(errCode).To(BeNil())
			Expect(ref).To(MatchRegexp("^zone_delegated.*"))

			// Update a Zone-delegated
			zone := &ibclient.ZoneDelegated{
				DelegateTo: ibclient.NullableNameServers{
					NameServers: []ibclient.NameServer{
						{Name: "test", Address: "1.2.3.4"},
						{Name: "test2", Address: "1.2.3.6"},
					},
					IsNull: false,
				},
				Comment: utils.StringPtr("wapi added"),
			}

			var res []ibclient.ZoneDelegated
			search := &ibclient.ZoneDelegated{}
			err := connector.GetObject(search, "", nil, &res)
			ref, err = connector.UpdateObject(zone, res[0].Ref)
			Expect(err).To(BeNil())
			Expect(ref).To(MatchRegexp("^zone_delegated.*"))
		})

		It("Should delete the Zone-delegated", func() {
			// Create a DNS Zone-delegated
			zoneCreate := &ibclient.ZoneDelegated{
				Fqdn: "example4.wapi.com",
				DelegateTo: ibclient.NullableNameServers{
					NameServers: []ibclient.NameServer{
						{Name: "test", Address: "1.2.3.4"},
						{Name: "test2", Address: "1.2.3.5"},
					},
					IsNull: false,
				},
				Comment: utils.StringPtr("wapi added"),
			}
			refCreate, errCode := connector.CreateObject(zoneCreate)
			Expect(errCode).To(BeNil())
			Expect(refCreate).To(MatchRegexp("^zone_delegated.*"))

			var res []ibclient.ZoneDelegated
			search := &ibclient.ZoneDelegated{}
			err := connector.GetObject(search, "", nil, &res)
			ref, err := connector.DeleteObject(res[0].Ref)
			Expect(err).To(BeNil())
			Expect(ref).To(MatchRegexp("^zone_delegated.*"))
		})

		It("Should fail to create a DNS Zone-delegated without mandatory parameters", func() {
			zone := &ibclient.ZoneDelegated{
				// Missing mandatory parameters like Fqdn and DelegatedTo oe NsGroup
				Comment:      utils.StringPtr("wapi added"),
				DelegatedTtl: utils.Uint32Ptr(3600),
			}
			_, err := connector.CreateObject(zone)
			Expect(err).NotTo(BeNil())
		})

		It("Should fail to get a non-existent DNS Zone-delegated", func() {
			var res []ibclient.ZoneDelegated
			search := &ibclient.ZoneDelegated{Fqdn: "nonexistent.test_fwzone.com"}
			err := connector.GetObject(search, "", nil, &res)
			Expect(err).NotTo(BeNil())
		})

		It("Should fail to update a non-existent DNS Zone-delegated", func() {
			zone := &ibclient.ZoneDelegated{
				Fqdn: "nonexistent.com",
				DelegateTo: ibclient.NullableNameServers{
					NameServers: []ibclient.NameServer{
						{Name: "test", Address: "1.2.3.4"},
						{Name: "test2", Address: "1.2.3.6"},
					},
					IsNull: false,
				},
				Comment: utils.StringPtr("wapi added"),
			}

			_, err := connector.UpdateObject(zone, "nonexistent_ref")
			Expect(err).NotTo(BeNil())
		})

		It("Should fail to delete a non-existent DNS Zone-delegated", func() {
			_, err := connector.DeleteObject("nonexistent_ref")
			Expect(err).NotTo(BeNil())
		})

		It("Should add CNAME Record [cname.wapi.com]", Label("ID: 40", "RW"), func() {
			r := &ibclient.RecordCNAME{
				View:      utils.StringPtr("default"),
				Name:      utils.StringPtr("cname.wapi.com"),
				Canonical: utils.StringPtr("qatest1.com"),
				Comment:   utils.StringPtr("verified"),
			}
			ref, err := connector.CreateObject(r)
			Expect(err).To(BeNil())
			Expect(ref).To(MatchRegexp("^record:cname.*cname\\.wapi\\.com/default$"))
		})

		When("CNAME Record [cname.wapi.com] exists", Label("RW"), func() {
			var ref string
			BeforeEach(func() {
				r := &ibclient.RecordCNAME{
					View:      utils.StringPtr("default"),
					Name:      utils.StringPtr("cname.wapi.com"),
					Canonical: utils.StringPtr("qatest1.com"),
					Comment:   utils.StringPtr("verified"),
				}
				var err error
				ref, err = connector.CreateObject(r)
				Expect(err).To(BeNil())
				Expect(ref).To(MatchRegexp("^record:cname.*cname\\.wapi\\.com/default$"))
			})

			It("Should get the CNAME Record [cname.wapi.com]", Label("ID: 108", "RO"), func() {
				// Get the CNAME Record [cname.wapi.com] to validate the above addition
				// case by searching with all the attributes
				var res []ibclient.RecordCNAME
				search := &ibclient.RecordCNAME{}
				qp := ibclient.NewQueryParams(false, map[string]string{
					"view":      "default",
					"name":      "cname.wapi.com",
					"canonical": "qatest1.com",
				})
				err := connector.GetObject(search, "", qp, &res)
				Expect(err).To(BeNil())
				Expect(res[0].Ref).To(MatchRegexp("record:cname.*cname\\.wapi\\.com/default$"))
				Expect(*res[0].View).To(Equal("default"))
				Expect(*res[0].Name).To(Equal("cname.wapi.com"))
				Expect(*res[0].Canonical).To(Equal("qatest1.com"))
			})

			It("Should modify the CNAME Record [cname.wapi.com] of the fields [comment, disable, ttl]",
				Label("ID: 110", "ID: 148", "RW"), func() {
					r := &ibclient.RecordCNAME{
						Comment: utils.StringPtr("Modified CNAME Record"),
						Disable: utils.BoolPtr(true),
						Ttl:     utils.Uint32Ptr(20),
					}
					ref, err := connector.UpdateObject(r, ref)
					Expect(err).To(BeNil())
					Expect(ref).To(MatchRegexp("^record:cname.*cname\\.wapi\\.com/default$"))

					// Get the CNAME Record [cname.wapi.com] to validate the above modified fields
					var res []ibclient.RecordCNAME
					search := &ibclient.RecordCNAME{}
					search.SetReturnFields([]string{"comment", "disable", "ttl", "use_ttl"})
					qp := ibclient.NewQueryParams(false, map[string]string{
						"name": "cname.wapi.com",
					})
					err = connector.GetObject(search, "", qp, &res)
					Expect(err).To(BeNil())
					Expect(res[0].Ref).To(MatchRegexp("record:cname.*cname\\.wapi\\.com/default$"))
					Expect(*res[0].Comment).To(Equal("Modified CNAME Record"))
					Expect(*res[0].Disable).To(Equal(true))
					Expect(*res[0].Ttl).To(Equal(uint32(20)))
					Expect(*res[0].UseTtl).To(Equal(true))
				},
			)
		})

		It("Should add A Record [a.wapi.com]", Label("ID: 41", "RW"), func() {
			r := &ibclient.RecordA{
				Name:     utils.StringPtr("a.wapi.com"),
				Ipv4Addr: utils.StringPtr("9.0.0.1"),
				Comment:  utils.StringPtr("Added A Record"),
				Disable:  utils.BoolPtr(false),
				Ttl:      utils.Uint32Ptr(10),
				UseTtl:   utils.BoolPtr(true),
				View:     "default",
			}
			ref, err := connector.CreateObject(r)
			Expect(err).To(BeNil())
			Expect(ref).To(MatchRegexp("^record:a.*a\\.wapi\\.com/default$"))
		})

		When("A Record [a.wapi.com] exists", Label("RW"), func() {
			var ref string
			BeforeEach(func() {
				r := &ibclient.RecordA{
					Name:     utils.StringPtr("a.wapi.com"),
					Ipv4Addr: utils.StringPtr("9.0.0.1"),
					Comment:  utils.StringPtr("Added A Record"),
					Disable:  utils.BoolPtr(false),
					Ttl:      utils.Uint32Ptr(10),
					UseTtl:   utils.BoolPtr(true),
					View:     "default",
				}
				var err error
				ref, err = connector.CreateObject(r)
				Expect(err).To(BeNil())
				Expect(ref).To(MatchRegexp("^record:a.*a\\.wapi\\.com/default$"))
			})

			It("Should get the A Record [a.wapi.com]", Label("ID: 109", "RO"), func() {
				// Get the A Record [a.wapi.com] to validate the above addition case
				var res []ibclient.RecordA
				search := &ibclient.RecordA{}
				search.SetReturnFields([]string{"name", "ipv4addr", "comment", "disable", "ttl", "use_ttl", "view"})
				qp := ibclient.NewQueryParams(false, map[string]string{
					"view": "default",
					"name": "a.wapi.com",
				})
				err := connector.GetObject(search, "", qp, &res)
				Expect(err).To(BeNil())
				Expect(res[0].Ref).To(MatchRegexp("record:a.*a\\.wapi\\.com/default$"))
				Expect(*res[0].Name).To(Equal("a.wapi.com"))
				Expect(*res[0].Ipv4Addr).To(Equal("9.0.0.1"))
				Expect(*res[0].Comment).To(Equal("Added A Record"))
				Expect(*res[0].Disable).To(Equal(false))
				Expect(*res[0].Ttl).To(Equal(uint32(10)))
				Expect(*res[0].UseTtl).To(Equal(true))
				Expect(res[0].View).To(Equal("default"))
			})

			It("Should modify the A Record [a.wapi.com] of the fields [comment, disable, ttl, use_ttl]",
				Label("ID: 111", "ID: 147", "RW"), func() {
					r := &ibclient.RecordA{
						Comment: utils.StringPtr("Modified A Record"),
						Disable: utils.BoolPtr(true),
						Ttl:     utils.Uint32Ptr(20),
						UseTtl:  utils.BoolPtr(false),
					}
					ref, err := connector.UpdateObject(r, ref)
					Expect(err).To(BeNil())
					Expect(ref).To(MatchRegexp("^record:a.*a\\.wapi\\.com/default$"))

					// Get the A Record [a.wapi.com] to validate the above modified fields
					var res []ibclient.RecordA
					search := &ibclient.RecordA{}
					search.SetReturnFields([]string{"comment", "disable", "ttl", "use_ttl"})
					qp := ibclient.NewQueryParams(false, map[string]string{
						"view": "default",
						"name": "a.wapi.com",
					})
					err = connector.GetObject(search, "", qp, &res)
					Expect(err).To(BeNil())
					Expect(res[0].Ref).To(MatchRegexp("record:a.*a\\.wapi\\.com/default$"))
					Expect(*res[0].Comment).To(Equal("Modified A Record"))
					Expect(*res[0].Disable).To(Equal(true))
					Expect(*res[0].Ttl).To(Equal(uint32(20)))
					Expect(*res[0].UseTtl).To(Equal(false))
				},
			)
		})

		It("Should add AAAA Record [aaaa.wapi.com]", Label("ID: 42", "RW"), func() {
			r := &ibclient.RecordAAAA{
				Name:     utils.StringPtr("aaaa.wapi.com"),
				Ipv6Addr: utils.StringPtr("99::99"),
				Comment:  utils.StringPtr("Added AAAA Record"),
				Disable:  utils.BoolPtr(false),
				Ttl:      utils.Uint32Ptr(10),
				UseTtl:   utils.BoolPtr(true),
				View:     "default",
			}
			ref, err := connector.CreateObject(r)
			Expect(err).To(BeNil())
			Expect(ref).To(MatchRegexp("^record:aaaa.*aaaa\\.wapi\\.com/default$"))
		})

		When("AAAA Record [aaaa.wapi.com] exists", Label("RW"), func() {
			var ref string
			BeforeEach(func() {
				r := &ibclient.RecordAAAA{
					Name:     utils.StringPtr("aaaa.wapi.com"),
					Ipv6Addr: utils.StringPtr("99::99"),
					Comment:  utils.StringPtr("Added AAAA Record"),
					Disable:  utils.BoolPtr(false),
					Ttl:      utils.Uint32Ptr(10),
					UseTtl:   utils.BoolPtr(true),
					View:     "default",
				}
				var err error
				ref, err = connector.CreateObject(r)
				Expect(err).To(BeNil())
				Expect(ref).To(MatchRegexp("^record:aaaa.*aaaa\\.wapi\\.com/default$"))
			})

			It("Should get the AAAA Record [aaaa.wapi.com]", Label("ID: 87", "RO"), func() {
				// Get the AAAA Record [aaaa.wapi.com] to validate the above addition case
				var res []ibclient.RecordAAAA
				search := &ibclient.RecordAAAA{}
				search.SetReturnFields([]string{"name", "ipv6addr", "comment", "disable", "ttl", "use_ttl", "view"})
				qp := ibclient.NewQueryParams(false, map[string]string{"ipv6addr": "99::99"})
				err := connector.GetObject(search, "", qp, &res)
				Expect(err).To(BeNil())
				Expect(res[0].Ref).To(MatchRegexp("record:aaaa.*aaaa\\.wapi\\.com/default$"))
				Expect(*res[0].Name).To(Equal("aaaa.wapi.com"))
				Expect(*res[0].Ipv6Addr).To(Equal("99::99"))
				Expect(*res[0].Comment).To(Equal("Added AAAA Record"))
				Expect(*res[0].Disable).To(Equal(false))
				Expect(*res[0].Ttl).To(Equal(uint32(10)))
				Expect(*res[0].UseTtl).To(Equal(true))
				Expect(res[0].View).To(Equal("default"))
			})

			It("Should modify the AAAA Record [aaaa.wapi.com] of the fields [comment, disable, ttl, use_ttl]",
				Label("ID: 112", "ID: 146", "RW"), func() {
					r := &ibclient.RecordAAAA{
						Comment: utils.StringPtr("Modified AAAA Record"),
						Disable: utils.BoolPtr(true),
						Ttl:     utils.Uint32Ptr(20),
						UseTtl:  utils.BoolPtr(false),
					}
					ref, err := connector.UpdateObject(r, ref)
					Expect(err).To(BeNil())
					Expect(ref).To(MatchRegexp("^record:aaaa.*aaaa\\.wapi\\.com/default$"))

					// Get the AAAA Record [aaaa.wapi.com] to validate the above modified fields
					var res []ibclient.RecordAAAA
					search := &ibclient.RecordAAAA{}
					search.SetReturnFields([]string{"comment", "disable", "ttl", "use_ttl"})
					qp := ibclient.NewQueryParams(false, map[string]string{"ipv6addr": "99::99"})
					err = connector.GetObject(search, "", qp, &res)
					Expect(err).To(BeNil())
					Expect(res[0].Ref).To(MatchRegexp("record:aaaa.*aaaa\\.wapi\\.com/default$"))
					Expect(*res[0].Comment).To(Equal("Modified AAAA Record"))
					Expect(*res[0].Disable).To(Equal(true))
					Expect(*res[0].Ttl).To(Equal(uint32(20)))
					Expect(*res[0].UseTtl).To(Equal(false))
				},
			)
		})

		It("Should add Host record [h1.wapi.com] with both ipv4addrs and ipv6addrs fields",
			Label("ID: 43", "RW"), func() {
				r := &ibclient.HostRecord{
					Name: utils.StringPtr("h1.wapi.com"),
					View: utils.StringPtr("default"),
					Ipv4Addrs: []ibclient.HostRecordIpv4Addr{
						{EnableDhcp: utils.BoolPtr(false), Ipv4Addr: utils.StringPtr("20.20.20.20")},
						{Ipv4Addr: utils.StringPtr("20.20.20.30")},
						{EnableDhcp: utils.BoolPtr(false), Ipv4Addr: utils.StringPtr("20.20.20.40")},
					},
					Ipv6Addrs: []ibclient.HostRecordIpv6Addr{
						{EnableDhcp: utils.BoolPtr(true), Ipv6Addr: utils.StringPtr("2000::1"), Duid: utils.StringPtr("11:10")},
						{Ipv6Addr: utils.StringPtr("2000::2")},
					},
				}
				ref, err := connector.CreateObject(r)
				Expect(err).To(BeNil())
				Expect(ref).To(MatchRegexp("^record:host.*h1\\.wapi\\.com/default$"))
			},
		)

		It("Should add Host record [h1.wapi.com] with both ipv4addrs and aliases fields when dns is enabled",
			Label("ID:44", "RW"), func() {
				r := &ibclient.HostRecord{
					Name:      utils.StringPtr("h1.wapi.com"),
					View:      utils.StringPtr("default"),
					EnableDns: utils.BoolPtr(true),
					Ipv4Addrs: []ibclient.HostRecordIpv4Addr{
						{EnableDhcp: utils.BoolPtr(false), Ipv4Addr: utils.StringPtr("20.20.20.20")},
						{Ipv4Addr: utils.StringPtr("20.20.20.30")},
						{EnableDhcp: utils.BoolPtr(false), Ipv4Addr: utils.StringPtr("20.20.20.40")},
					},
					Aliases: []string{
						"alias1.wapi.com",
						"alias2.wapi.com",
					},
				}
				ref, err := connector.CreateObject(r)
				Expect(err).To(BeNil())
				Expect(ref).To(MatchRegexp("^record:host.*h1\\.wapi\\.com/default$"))
			},
		)
		It("Should add Host record [h1.wapi.com] with ipv4addrs  and aliases fields when dns is disabled",
			Label("ID:45", "RW"), func() {
				r := &ibclient.HostRecord{
					Name:      utils.StringPtr("h1.wapi.com"),
					View:      utils.StringPtr("default"),
					EnableDns: utils.BoolPtr(false),
					Ipv4Addrs: []ibclient.HostRecordIpv4Addr{
						{EnableDhcp: utils.BoolPtr(false), Ipv4Addr: utils.StringPtr("16.12.1.20")},
						{Ipv4Addr: utils.StringPtr("16.12.1.30")},
						{EnableDhcp: utils.BoolPtr(false), Ipv4Addr: utils.StringPtr("16.12.1.40")},
					},
					Aliases: []string{
						"alias1",
						"alias2.wapi.com",
					},
				}
				ref, err := connector.CreateObject(r)
				Expect(err).To(BeNil())
				Expect(ref).To(MatchRegexp("^record:host.*h1\\.wapi\\.com(/.*)?$"))
			},
		)
		When("Host record [h1.wapi.com] with both ipv4addrs and ipv6addrs fields exits",
			Label("RW"), func() {
				BeforeEach(func() {
					r := &ibclient.HostRecord{
						Name: utils.StringPtr("h1.wapi.com"),
						View: utils.StringPtr("default"),
						Ipv4Addrs: []ibclient.HostRecordIpv4Addr{
							{EnableDhcp: utils.BoolPtr(false), Ipv4Addr: utils.StringPtr("20.20.20.20")},
							{Ipv4Addr: utils.StringPtr("20.20.20.30")},
							{EnableDhcp: utils.BoolPtr(false), Ipv4Addr: utils.StringPtr("20.20.20.40")},
						},
						Ipv6Addrs: []ibclient.HostRecordIpv6Addr{
							{EnableDhcp: utils.BoolPtr(true), Ipv6Addr: utils.StringPtr("2000::1"), Duid: utils.StringPtr("11:10")},
							{Ipv6Addr: utils.StringPtr("2000::2")},
						},
					}
					ref, err := connector.CreateObject(r)
					Expect(err).To(BeNil())
					Expect(ref).To(MatchRegexp("^record:host.*h1\\.wapi\\.com/default$"))
				})

				It("Should get the DNS Host record object", Label("ID: 91", "RO"), func() {
					var res []ibclient.HostRecord
					search := &ibclient.HostRecord{}
					qp := ibclient.NewQueryParams(false, map[string]string{
						"view": "default",
						"name": "h1.wapi.com",
					})
					err := connector.GetObject(search, "", qp, &res)
					Expect(err).To(BeNil())
					Expect(*res[0].Name).To(Equal("h1.wapi.com"))

					Expect(*res[0].Ipv6Addrs[0].EnableDhcp).To(Equal(true))
					Expect(res[0].Ipv6Addrs[0].Ref).To(MatchRegexp("^record:host_ipv6addr.*2000%3A%3A1/h1\\.wapi\\.com/default$"))
					Expect(res[0].Ipv6Addrs[0].Host).To(Equal("h1.wapi.com"))
					Expect(*res[0].Ipv6Addrs[0].Ipv6Addr).To(Equal("2000::1"))
					Expect(*res[0].Ipv6Addrs[0].Duid).To(Equal("11:10"))
					Expect(*res[0].Ipv6Addrs[1].EnableDhcp).To(Equal(false))
					Expect(res[0].Ipv6Addrs[1].Ref).To(MatchRegexp("^record:host_ipv6addr.*2000%3A%3A2/h1\\.wapi\\.com/default$"))
					Expect(res[0].Ipv6Addrs[1].Host).To(Equal("h1.wapi.com"))
					Expect(*res[0].Ipv6Addrs[1].Ipv6Addr).To(Equal("2000::2"))

					Expect(*res[0].Ipv4Addrs[0].EnableDhcp).To(Equal(false))
					Expect(res[0].Ipv4Addrs[0].Host).To(Equal("h1.wapi.com"))
					Expect(*res[0].Ipv4Addrs[0].Ipv4Addr).To(Equal("20.20.20.20"))
					Expect(res[0].Ipv4Addrs[0].Ref).To(MatchRegexp("^record:host_ipv4addr.*20\\.20\\.20\\.20/h1\\.wapi\\.com/default$"))
					Expect(*res[0].Ipv4Addrs[1].EnableDhcp).To(Equal(false))
					Expect(res[0].Ipv4Addrs[1].Host).To(Equal("h1.wapi.com"))
					Expect(*res[0].Ipv4Addrs[1].Ipv4Addr).To(Equal("20.20.20.30"))
					Expect(res[0].Ipv4Addrs[1].Ref).To(MatchRegexp("^record:host_ipv4addr.*20\\.20\\.20\\.30/h1\\.wapi\\.com/default$"))
					Expect(*res[0].Ipv4Addrs[2].EnableDhcp).To(Equal(false))
					Expect(res[0].Ipv4Addrs[2].Host).To(Equal("h1.wapi.com"))
					Expect(*res[0].Ipv4Addrs[2].Ipv4Addr).To(Equal("20.20.20.40"))
					Expect(res[0].Ipv4Addrs[2].Ref).To(MatchRegexp("^record:host_ipv4addr.*20\\.20\\.20\\.40/h1\\.wapi\\.com/default$"))

					Expect(res[0].Ref).To(MatchRegexp("^record:host.*h1\\.wapi\\.com/default$"))
					Expect(*res[0].View).To(Equal("default"))
				})

				It("Should get the IPv4 Host address object", Label("ID: 92", "RO"), func() {
					var res []ibclient.HostRecordIpv4Addr
					search := &ibclient.HostRecordIpv4Addr{}
					qp := ibclient.NewQueryParams(false, map[string]string{
						"network_view": "default",
						"ipv4addr":     "20.20.20.20",
					})
					err := connector.GetObject(search, "", qp, &res)
					Expect(err).To(BeNil())
					Expect(*res[0].EnableDhcp).To(Equal(false))
					Expect(res[0].Host).To(Equal("h1.wapi.com"))
					Expect(*res[0].Ipv4Addr).To(Equal("20.20.20.20"))
					Expect(res[0].Ref).To(MatchRegexp("^record:host_ipv4addr.*20\\.20\\.20\\.20/h1\\.wapi\\.com/default$"))

					qp = ibclient.NewQueryParams(false, map[string]string{
						"network_view": "default",
						"ipv4addr":     "20.20.20.30",
					})
					err = connector.GetObject(search, "", qp, &res)
					Expect(err).To(BeNil())
					Expect(*res[0].EnableDhcp).To(Equal(false))
					Expect(res[0].Host).To(Equal("h1.wapi.com"))
					Expect(*res[0].Ipv4Addr).To(Equal("20.20.20.30"))
					Expect(res[0].Ref).To(MatchRegexp("^record:host_ipv4addr.*20\\.20\\.20\\.30/h1\\.wapi\\.com/default$"))

					qp = ibclient.NewQueryParams(false, map[string]string{
						"network_view": "default",
						"ipv4addr":     "20.20.20.40",
					})
					err = connector.GetObject(search, "", qp, &res)
					Expect(*res[0].EnableDhcp).To(Equal(false))
					Expect(res[0].Host).To(Equal("h1.wapi.com"))
					Expect(*res[0].Ipv4Addr).To(Equal("20.20.20.40"))
					Expect(res[0].Ref).To(MatchRegexp("^record:host_ipv4addr.*20\\.20\\.20\\.40/h1\\.wapi\\.com/default$"))

				})

				It("Should get the IPv6 Host address object", Label("ID: 93", "RO"), func() {
					var res []ibclient.HostRecordIpv6Addr
					search := &ibclient.HostRecordIpv6Addr{}
					qp := ibclient.NewQueryParams(false, map[string]string{
						"network_view": "default",
						"ipv6addr":     "2000::1",
					})
					err := connector.GetObject(search, "", qp, &res)
					Expect(err).To(BeNil())
					Expect(*res[0].EnableDhcp).To(Equal(true))
					Expect(res[0].Ref).To(MatchRegexp("^record:host_ipv6addr.*2000%3A%3A1/h1\\.wapi\\.com/default$"))
					Expect(res[0].Host).To(Equal("h1.wapi.com"))
					Expect(*res[0].Ipv6Addr).To(Equal("2000::1"))
					Expect(*res[0].Duid).To(Equal("11:10"))

					qp = ibclient.NewQueryParams(false, map[string]string{
						"network_view": "default",
						"ipv6addr":     "2000::2",
					})
					err = connector.GetObject(search, "", qp, &res)
					Expect(err).To(BeNil())
					Expect(*res[0].EnableDhcp).To(Equal(false))
					Expect(res[0].Ref).To(MatchRegexp("^record:host_ipv6addr.*2000%3A%3A2/h1\\.wapi\\.com/default$"))
					Expect(res[0].Host).To(Equal("h1.wapi.com"))
					Expect(*res[0].Ipv6Addr).To(Equal("2000::2"))

				})
			},
		)

		It("Should add MX Record [mx.wapi.com]", Label("ID: 44", "RW"), func() {
			r := &ibclient.RecordMX{
				Name:          utils.StringPtr("mx.wapi.com"),
				MailExchanger: utils.StringPtr("wapi.com"),
				Preference:    utils.Uint32Ptr(10),
				Comment:       utils.StringPtr("Creating mx record through infoblox-go-client"),
				Disable:       utils.BoolPtr(false),
				Ttl:           utils.Uint32Ptr(20),
				UseTtl:        utils.BoolPtr(true),
			}
			ref, err := connector.CreateObject(r)
			Expect(err).To(BeNil())
			Expect(ref).To(MatchRegexp("^record:mx.*mx\\.wapi\\.com/default$"))
		})

		When("MX Record [mx.wapi.com] exists", Label("RW"), func() {
			var ref string
			BeforeEach(func() {
				r := &ibclient.RecordMX{
					Name:          utils.StringPtr("mx.wapi.com"),
					MailExchanger: utils.StringPtr("wapi.com"),
					Preference:    utils.Uint32Ptr(10),
					Comment:       utils.StringPtr("Creating mx record through infoblox-go-client"),
					Disable:       utils.BoolPtr(false),
					Ttl:           utils.Uint32Ptr(20),
					UseTtl:        utils.BoolPtr(true),
				}
				var err error
				ref, err = connector.CreateObject(r)
				Expect(err).To(BeNil())
				Expect(ref).To(MatchRegexp("^record:mx.*mx\\.wapi\\.com/default$"))
			})

			It("Should get the MX Record [mx.wapi.com]", Label("ID: 86", "RO"), func() {
				// Get the MX Record [mx.wapi.com] to validate the above addition case
				var res []ibclient.RecordMX
				search := &ibclient.RecordMX{}
				search.SetReturnFields([]string{"comment", "disable", "mail_exchanger", "name", "preference", "ttl", "use_ttl", "view", "zone"})
				qp := ibclient.NewQueryParams(false, map[string]string{
					"view": "default",
					"name": "mx.wapi.com",
				})
				err := connector.GetObject(search, "", qp, &res)
				Expect(err).To(BeNil())
				Expect(*res[0].Comment).To(Equal("Creating mx record through infoblox-go-client"))
				Expect(*res[0].Name).To(Equal("mx.wapi.com"))
				Expect(res[0].Zone).To(Equal("wapi.com"))
				Expect(*res[0].MailExchanger).To(Equal("wapi.com"))
				Expect(*res[0].UseTtl).To(Equal(true))
				Expect(*res[0].Disable).To(Equal(false))
				Expect(*res[0].Preference).To(Equal(uint32(10)))
				Expect(*res[0].Ttl).To(Equal(uint32(20)))
				Expect(res[0].Ref).To(MatchRegexp("^record:mx.*mx\\.wapi\\.com/default$"))
				Expect(*res[0].View).To(Equal("default"))
			})

			It("Should modify the MX Record [mx.wapi.com] of the fields [comment, disable, ttl, use_ttl]",
				Label("ID: 113", "ID: 145", "RW"), func() {
					r := &ibclient.RecordMX{
						Comment: utils.StringPtr("Modified mx Record"),
						Disable: utils.BoolPtr(true),
						Ttl:     utils.Uint32Ptr(120),
						UseTtl:  utils.BoolPtr(false),
					}
					ref, err := connector.UpdateObject(r, ref)
					Expect(err).To(BeNil())
					Expect(ref).To(MatchRegexp("^record:mx.*mx\\.wapi\\.com/default$"))

					// Get the MX Record [mx.wapi.com] to validate the above modified fields
					var res []ibclient.RecordMX
					search := &ibclient.RecordMX{}
					search.SetReturnFields([]string{"comment", "disable", "ttl", "use_ttl"})
					qp := ibclient.NewQueryParams(false, map[string]string{
						"view": "default",
						"name": "mx.wapi.com",
					})
					err = connector.GetObject(search, "", qp, &res)
					Expect(err).To(BeNil())
					Expect(res[0].Ref).To(MatchRegexp("^record:mx.*mx\\.wapi\\.com/default$"))
					Expect(*res[0].Comment).To(Equal("Modified mx Record"))
					Expect(*res[0].Disable).To(Equal(true))
					Expect(*res[0].Ttl).To(Equal(uint32(120)))
					Expect(*res[0].UseTtl).To(Equal(false))
				},
			)
		})

		It("Should add TXT Record [txt.wapi.com]", Label("ID: 45", "RW"), func() {
			r := &ibclient.RecordTXT{
				Name:    utils.StringPtr("txt.wapi.com"),
				Text:    utils.StringPtr("wapi.com"),
				Comment: utils.StringPtr("Creating txt record through infoblox-go-client"),
				Disable: utils.BoolPtr(false),
				Ttl:     utils.Uint32Ptr(20),
				UseTtl:  utils.BoolPtr(true),
				View:    utils.StringPtr("default"),
			}
			ref, err := connector.CreateObject(r)
			Expect(err).To(BeNil())
			Expect(ref).To(MatchRegexp("^record:txt.*txt\\.wapi\\.com/default$"))
		})

		When("TXT Record [txt.wapi.com] exists", Label("RW"), func() {
			var ref string
			BeforeEach(func() {
				r := &ibclient.RecordTXT{
					Name:    utils.StringPtr("txt.wapi.com"),
					Text:    utils.StringPtr("wapi.com"),
					Comment: utils.StringPtr("Creating txt record through infoblox-go-client"),
					Disable: utils.BoolPtr(false),
					Ttl:     utils.Uint32Ptr(20),
					UseTtl:  utils.BoolPtr(true),
					View:    utils.StringPtr("default"),
				}
				var err error
				ref, err = connector.CreateObject(r)
				Expect(err).To(BeNil())
				Expect(ref).To(MatchRegexp("^record:txt.*txt\\.wapi\\.com/default$"))
			})

			It("Should get the TXT Record with all the fields", Label("ID: 85", "RO"), func() {
				// Get the TXT Record with all the fields
				var res []ibclient.RecordTXT
				search := &ibclient.RecordTXT{}
				search.SetReturnFields([]string{"comment", "disable", "text", "name", "ttl", "use_ttl", "view", "zone"})
				qp := ibclient.NewQueryParams(false, map[string]string{
					"view": "default",
					"name": "txt.wapi.com",
				})
				err := connector.GetObject(search, "", qp, &res)
				Expect(err).To(BeNil())
				Expect(*res[0].Comment).To(Equal("Creating txt record through infoblox-go-client"))
				Expect(*res[0].Name).To(Equal("txt.wapi.com"))
				Expect(res[0].Zone).To(Equal("wapi.com"))
				Expect(*res[0].Text).To(Equal("wapi.com"))
				Expect(*res[0].UseTtl).To(Equal(true))
				Expect(*res[0].Disable).To(Equal(false))
				Expect(*res[0].Ttl).To(Equal(uint32(20)))
				Expect(res[0].Ref).To(MatchRegexp("record:txt.*txt\\.wapi\\.com/default$"))
				Expect(*res[0].View).To(Equal("default"))
			})

			It("Should modify the TXT Record [txt.wapi.com] of the fields [comment, disable, ttl, use_ttl]",
				Label("ID: 114", "ID: 144", "RW"), func() {
					r := &ibclient.RecordTXT{
						Comment: utils.StringPtr("Modified TXT Record"),
						Disable: utils.BoolPtr(true),
						Ttl:     utils.Uint32Ptr(120),
						UseTtl:  utils.BoolPtr(false),
					}
					ref, err := connector.UpdateObject(r, ref)
					Expect(err).To(BeNil())
					Expect(ref).To(MatchRegexp("^record:txt.*txt\\.wapi\\.com/default$"))

					// Get the TXT Record [txt.wapi.com] to validate the above modified fields
					var res []ibclient.RecordTXT
					search := &ibclient.RecordTXT{}
					search.SetReturnFields([]string{"comment", "disable", "ttl", "use_ttl"})
					qp := ibclient.NewQueryParams(false, map[string]string{
						"view": "default",
						"name": "txt.wapi.com",
					})
					err = connector.GetObject(search, "", qp, &res)
					Expect(err).To(BeNil())
					Expect(res[0].Ref).To(MatchRegexp("record:txt.*txt\\.wapi\\.com/default$"))
					Expect(*res[0].Comment).To(Equal("Modified TXT Record"))
					Expect(*res[0].Disable).To(Equal(true))
					Expect(*res[0].Ttl).To(Equal(uint32(120)))
					Expect(*res[0].UseTtl).To(Equal(false))
				},
			)
		})

		It("Should add PTR Record [ptr1.wapi.com]", Label("ID: 46", "RW"), func() {
			r := &ibclient.RecordPTR{
				Name:     utils.StringPtr("ptr1.wapi.com"),
				PtrdName: utils.StringPtr("ptr.wapi.com"),
				View:     "default",
				Comment:  utils.StringPtr("wapi added"),
				Disable:  utils.BoolPtr(false),
				Ttl:      utils.Uint32Ptr(10),
				UseTtl:   utils.BoolPtr(true),
			}
			ref, err := connector.CreateObject(r)
			Expect(err).To(BeNil())
			Expect(ref).To(MatchRegexp("^record:ptr.*ptr1\\.wapi\\.com/default$"))
		})

		When("PTR Record [ptr1.wapi.com] exists", Label("RW"), func() {
			var ref string
			BeforeEach(func() {
				r := &ibclient.RecordPTR{
					Name:     utils.StringPtr("ptr1.wapi.com"),
					PtrdName: utils.StringPtr("ptr.wapi.com"),
					View:     "default",
					Comment:  utils.StringPtr("wapi added"),
					Disable:  utils.BoolPtr(false),
					Ttl:      utils.Uint32Ptr(10),
					UseTtl:   utils.BoolPtr(true),
				}
				var err error
				ref, err = connector.CreateObject(r)
				Expect(err).To(BeNil())
				Expect(ref).To(MatchRegexp("^record:ptr.*ptr1\\.wapi\\.com/default$"))
			})

			It("Should get the PTR Record [ptr1.wapi.com]", Label("ID: 84", "RO"), func() {
				// Get the PTR Record [ptr1.wapi.com] to validate the above addition case
				var res []ibclient.RecordPTR
				search := &ibclient.RecordPTR{}
				search.SetReturnFields([]string{"comment", "ptrdname", "disable", "name", "ttl", "use_ttl", "view", "zone"})
				qp := ibclient.NewQueryParams(false, map[string]string{
					"view": "default",
					"name": "ptr1.wapi.com",
				})
				err := connector.GetObject(search, "", qp, &res)
				Expect(err).To(BeNil())
				Expect(*res[0].Comment).To(Equal("wapi added"))
				Expect(*res[0].PtrdName).To(Equal("ptr.wapi.com"))
				Expect(*res[0].Name).To(Equal("ptr1.wapi.com"))
				Expect(res[0].Zone).To(Equal("wapi.com"))
				Expect(*res[0].UseTtl).To(Equal(true))
				Expect(*res[0].Disable).To(Equal(false))
				Expect(*res[0].Ttl).To(Equal(uint32(10)))
				Expect(res[0].Ref).To(MatchRegexp("^record:ptr.*ptr1\\.wapi\\.com/default$"))
				Expect(res[0].View).To(Equal("default"))
			})

			It("Should modify the PTR Record [ptr1.wapi.com] of the fields [comment, disable, ttl, use_ttl]",
				Label("ID: 115", "ID: 143", "RW"), func() {
					r := &ibclient.RecordTXT{
						Comment: utils.StringPtr("Modified PTR Record"),
						Disable: utils.BoolPtr(true),
						Ttl:     utils.Uint32Ptr(120),
						UseTtl:  utils.BoolPtr(false),
					}
					ref, err := connector.UpdateObject(r, ref)
					Expect(err).To(BeNil())
					Expect(ref).To(MatchRegexp("^record:ptr.*ptr1\\.wapi\\.com/default$"))

					// Get the PTR Record [ptr1.wapi.com] to validate the above modified fields
					var res []ibclient.RecordPTR
					search := &ibclient.RecordPTR{}
					search.SetReturnFields([]string{"comment", "disable", "ttl", "use_ttl"})
					qp := ibclient.NewQueryParams(false, map[string]string{"name": "ptr1.wapi.com"})
					err = connector.GetObject(search, "", qp, &res)
					Expect(err).To(BeNil())
					Expect(res[0].Ref).To(MatchRegexp("^record:ptr.*ptr1\\.wapi\\.com/default$"))
					Expect(*res[0].Comment).To(Equal("Modified PTR Record"))
					Expect(*res[0].Disable).To(Equal(true))
					Expect(*res[0].Ttl).To(Equal(uint32(120)))
					Expect(*res[0].UseTtl).To(Equal(false))
				},
			)
		})

		It("Should add SRV Record with all the attributes, using free format [name = srv.wapi.com]",
			Label("ID: 47", "RW"), func() {
				r := &ibclient.RecordSRV{
					Name:     utils.StringPtr("srv.wapi.com"),
					Weight:   utils.Uint32Ptr(10),
					Priority: utils.Uint32Ptr(10),
					Port:     utils.Uint32Ptr(10),
					Target:   utils.StringPtr("srv.wapi.com"),
					Comment:  utils.StringPtr("wapi added"),
					Disable:  utils.BoolPtr(false),
					Ttl:      utils.Uint32Ptr(10),
					UseTtl:   utils.BoolPtr(true),
					View:     "default",
				}
				ref, err := connector.CreateObject(r)
				Expect(err).To(BeNil())
				Expect(ref).To(MatchRegexp("^record:srv.*srv\\.wapi\\.com/default$"))
			},
		)

		When("SRV Record [name = srv.wapi.com] exists", Label("RW"), func() {
			var ref string
			BeforeEach(func() {
				r := &ibclient.RecordSRV{
					Name:     utils.StringPtr("srv.wapi.com"),
					Weight:   utils.Uint32Ptr(10),
					Priority: utils.Uint32Ptr(10),
					Port:     utils.Uint32Ptr(10),
					Target:   utils.StringPtr("srv.wapi.com"),
					Comment:  utils.StringPtr("wapi added"),
					Disable:  utils.BoolPtr(false),
					Ttl:      utils.Uint32Ptr(10),
					UseTtl:   utils.BoolPtr(true),
					View:     "default",
				}
				var err error
				ref, err = connector.CreateObject(r)
				Expect(err).To(BeNil())
				Expect(ref).To(MatchRegexp("^record:srv.*srv\\.wapi\\.com/default$"))
			})

			It("Should get SRV record [name = srv.wapi.com]", Label("ID: 83", "RO"), func() {

				// Get SRV record to validate above case
				var res []ibclient.RecordSRV
				search := &ibclient.RecordSRV{}
				search.SetReturnFields([]string{
					"name", "weight", "priority", "port", "target",
					"comment", "disable", "ttl", "use_ttl", "zone", "view",
				})
				qp := ibclient.NewQueryParams(false, map[string]string{"name": "srv.wapi.com"})
				err := connector.GetObject(search, "", qp, &res)
				Expect(err).To(BeNil())
				Expect(*res[0].Comment).To(Equal("wapi added"))
				Expect(res[0].View).To(Equal("default"))
				Expect(*res[0].Name).To(Equal("srv.wapi.com"))
				Expect(*res[0].Weight).To(Equal(uint32(10)))
				Expect(*res[0].UseTtl).To(Equal(true))
				Expect(*res[0].Priority).To(Equal(uint32(10)))
				Expect(*res[0].Disable).To(Equal(false))
				Expect(*res[0].Ttl).To(Equal(uint32(10)))
				Expect(res[0].Ref).To(MatchRegexp("record:srv.*srv\\.wapi\\.com/default$"))
				Expect(res[0].Zone).To(Equal("wapi.com"))
				Expect(*res[0].Port).To(Equal(uint32(10)))
				Expect(*res[0].Target).To(Equal("srv.wapi.com"))
			})

			It("Should modify the SRV Record [srv.wapi.com] of the fields [comment, disable, ttl, use_ttl]",
				Label("ID: 116", "ID: 142", "RW"), func() {
					r := &ibclient.RecordSRV{
						Comment: utils.StringPtr("Modified SRV Record"),
						Disable: utils.BoolPtr(true),
						Ttl:     utils.Uint32Ptr(120),
						UseTtl:  utils.BoolPtr(false),
					}
					ref, err := connector.UpdateObject(r, ref)
					Expect(err).To(BeNil())
					Expect(ref).To(MatchRegexp("^record:srv.*srv\\.wapi\\.com/default$"))

					// Get the SRV Record [srv.wapi.com] to validate the above modified fields
					var res []ibclient.RecordSRV
					search := &ibclient.RecordSRV{}
					search.SetReturnFields([]string{"comment", "disable", "ttl", "use_ttl"})
					qp := ibclient.NewQueryParams(false, map[string]string{"name": "srv.wapi.com"})
					err = connector.GetObject(search, "", qp, &res)
					Expect(err).To(BeNil())
					Expect(res[0].Ref).To(MatchRegexp("record:srv.*srv\\.wapi\\.com/default$"))
					Expect(*res[0].Comment).To(Equal("Modified SRV Record"))
					Expect(*res[0].Disable).To(Equal(true))
					Expect(*res[0].Ttl).To(Equal(uint32(120)))
					Expect(*res[0].UseTtl).To(Equal(false))
				},
			)
		})

		It("Should get the FMZ Auth Zone [wapi.com]", Label("ID: 88", "RO"), func() {
			var res []ibclient.ZoneAuth
			search := &ibclient.ZoneAuth{}
			search.SetReturnFields(append(search.ReturnFields(), "comment"))
			qp := ibclient.NewQueryParams(false, map[string]string{"view": "default", "fqdn": "wapi.com"})
			err := connector.GetObject(search, "", qp, &res)
			Expect(err).To(BeNil())
			Expect(res[0].Ref).To(MatchRegexp("zone_auth.*wapi\\.com/default$"))
			Expect(*res[0].View).To(Equal("default"))
			Expect(res[0].Fqdn).To(Equal("wapi.com"))
		})

		It("Should modify the fields [comment,allow_update_forwarding,notify_delay] for zone [wapi.com]",
			Label("ID: 127", "ID: 131", "RW"), func() {
				za := &ibclient.ZoneAuth{
					Comment:               utils.StringPtr("WAPI Modified Comment"),
					AllowUpdateForwarding: utils.BoolPtr(true),
					NotifyDelay:           utils.Uint32Ptr(100),
				}
				ref, err := connector.UpdateObject(za, refZoneAuth)
				Expect(err).To(BeNil())
				Expect(ref).To(MatchRegexp("^zone_auth.*wapi.com/default$"))

				// Using GMC IP, Get the auth zone [wapi.com] to validate the above modified fields
				var res []ibclient.ZoneAuth
				search := &ibclient.ZoneAuth{}
				qp := ibclient.NewQueryParams(false, map[string]string{"fqdn": "wapi.com"})
				search.SetReturnFields([]string{"comment", "allow_update_forwarding", "notify_delay"})
				err = connector.GetObject(search, "", qp, &res)
				Expect(err).To(BeNil())
				Expect(res[0].Ref).To(MatchRegexp("^zone_auth.*wapi.com/default$"))
				Expect(*res[0].Comment).To(Equal("WAPI Modified Comment"))
				Expect(*res[0].AllowUpdateForwarding).To(Equal(true))
				Expect(*res[0].NotifyDelay).To(Equal(uint32(100)))
			},
		)

		It("Should delete the auth zone [wapi.com]", Label("ID: 149", "RW"), func() {
			ref, err := connector.DeleteObject(refZoneAuth)
			Expect(err).To(BeNil())
			Expect(ref).To(MatchRegexp("^zone_auth.*wapi.com/default$"))
		})
	})

	When("Network View [dhcpview] exists", func() {
		var viewRef string
		BeforeEach(func() {
			nv := &ibclient.NetworkView{
				Name:    utils.StringPtr("dhcpview"),
				Comment: utils.StringPtr("wapi added"),
			}
			var err error
			viewRef, err = connector.CreateObject(nv)
			Expect(err).To(BeNil())
		})

		It("Should add Network [92.0.0.0/8] in custom network view [dhcpview]", Label("ID: 48", "RW"), func() {
			r := &ibclient.Ipv4Network{
				Comment:     utils.StringPtr("Add ipv4network through WAPI"),
				Network:     utils.StringPtr("92.0.0.0/8"),
				NetworkView: "dhcpview",
			}
			ref, err := connector.CreateObject(r)
			Expect(err).To(BeNil())
			Expect(ref).To(MatchRegexp("^network.*92\\.0\\.0\\.0/8/dhcpview$"))
		})

		When("Network [92.0.0.0/8] in custom network view [dhcpview] exists", Label("RW"), func() {
			var refNetwork string
			BeforeEach(func() {
				r := &ibclient.Ipv4Network{
					Comment:     utils.StringPtr("Add ipv4network through WAPI"),
					Network:     utils.StringPtr("92.0.0.0/8"),
					NetworkView: "dhcpview",
				}
				var err error
				refNetwork, err = connector.CreateObject(r)
				Expect(err).To(BeNil())
				Expect(refNetwork).To(MatchRegexp("^network.*92\\.0\\.0\\.0/8/dhcpview$"))
			})

			It("Should add IPv4 Range with mandatory fields start_addr [92.0.0.10] end_addr [92.0.0.20]",
				Label("ID: 49", "ID: 80", "ID: 119", "ID: 139", "RW"), func() {
					r := &ibclient.Range{
						StartAddr:   utils.StringPtr("92.0.0.10"),
						EndAddr:     utils.StringPtr("92.0.0.20"),
						Comment:     utils.StringPtr("Add Range through WAPI"),
						NetworkView: utils.StringPtr("dhcpview"),
					}
					ref, err := connector.CreateObject(r)
					Expect(err).To(BeNil())
					Expect(ref).To(MatchRegexp("^range.*92\\.0\\.0\\.10/92\\.0\\.0\\.20/dhcpview$"))

					var res []ibclient.Range
					search := &ibclient.Range{}
					qp := ibclient.NewQueryParams(false, map[string]string{"start_addr": "92.0.0.10", "end_addr": "92.0.0.20"})
					err = connector.GetObject(search, "", qp, &res)
					Expect(err).To(BeNil())
					Expect(*res[0].Network).To(Equal("92.0.0.0/8"))
					Expect(*res[0].Comment).To(Equal("Add Range through WAPI"))
					Expect(*res[0].NetworkView).To(Equal("dhcpview"))
					Expect(*res[0].StartAddr).To(Equal("92.0.0.10"))
					Expect(*res[0].EndAddr).To(Equal("92.0.0.20"))
					Expect(res[0].Ref).To(MatchRegexp("92.0.0.10/92.0.0.20/dhcpview$"))

					// Modify the Range [92.0.0.10/92.0.0.20] of the
					// fields [comment, bootfile, use_bootfile]
					r = &ibclient.Range{
						Comment:     utils.StringPtr("modified comment"),
						UseBootfile: utils.BoolPtr(true),
						Bootfile:    utils.StringPtr("boot_file"),
					}
					ref, err = connector.UpdateObject(r, ref)
					Expect(err).To(BeNil())
					Expect(ref).To(MatchRegexp("^range.*92\\.0\\.0\\.10/92\\.0\\.0\\.20/dhcpview$"))

					// Get the Range [92.0.0.10/92.0.0.20] to validate the above modified fields
					var resUpd []ibclient.Range
					searchUpd := &ibclient.Range{}
					searchUpd.SetReturnFields([]string{"comment", "bootfile", "use_bootfile"})
					qp = ibclient.NewQueryParams(false, map[string]string{"start_addr": "92.0.0.10", "end_addr": "92.0.0.20"})
					err = connector.GetObject(searchUpd, "", qp, &resUpd)
					Expect(err).To(BeNil())
					Expect(resUpd[0].Ref).To(MatchRegexp("92\\.0\\.0\\.10/92\\.0\\.0\\.20/dhcpview$"))
					Expect(*resUpd[0].Comment).To(Equal("modified comment"))
					Expect(*resUpd[0].UseBootfile).To(Equal(true))
					Expect(*resUpd[0].Bootfile).To(Equal("boot_file"))
				},
			)

			It("Should add IPv4 fixed address [92.0.0.2] and mac [11:11:11:11:11:15]",
				Label("ID: 55", "ID: 75", "ID: 124", "ID: 134", "RW"), func() {
					fa := &ibclient.Ipv4FixedAddress{
						Name:        utils.StringPtr("wapi-fa1"),
						Ipv4Addr:    utils.StringPtr("92.0.0.2"),
						NetworkView: utils.StringPtr("dhcpview"),
						Mac:         utils.StringPtr("11:11:11:11:11:15"),
						Comment:     utils.StringPtr("HellO"),
					}
					ref, err := connector.CreateObject(fa)
					Expect(err).To(BeNil())
					Expect(ref).To(MatchRegexp("^fixedaddress.*92\\.0\\.0\\.2/dhcpview"))

					var res []ibclient.Ipv4FixedAddress
					search := &ibclient.Ipv4FixedAddress{}
					qp := ibclient.NewQueryParams(false, map[string]string{"ipv4addr": "92.0.0.2"})
					err = connector.GetObject(search, "", qp, &res)
					Expect(err).To(BeNil())
					Expect(res[0].Ref).To(MatchRegexp("^fixedaddress.*92\\.0\\.0\\.2/dhcpview"))
					Expect(*res[0].Ipv4Addr).To(MatchRegexp("92.0.0.2"))
					Expect(*res[0].NetworkView).To(MatchRegexp("dhcpview"))

					// Modify the fields [mac,comment,deny_bootp, pxe_lease_time]
					// for IPv4 fixed address [92.0.0.2]
					r := &ibclient.Ipv4FixedAddress{
						Mac:          utils.StringPtr("11:11:11:11:11:18"),
						Comment:      utils.StringPtr("changed it"),
						DenyBootp:    utils.BoolPtr(true),
						PxeLeaseTime: utils.Uint32Ptr(10),
					}
					ref, err = connector.UpdateObject(r, ref)
					Expect(err).To(BeNil())
					Expect(ref).To(MatchRegexp("^fixedaddress.*92\\.0\\.0\\.2/dhcpview"))

					// Get the IPv4 fixed address [92.0.0.2] to validate the above modified fields
					var resUpd []ibclient.Ipv4FixedAddress
					searchUpd := &ibclient.Ipv4FixedAddress{}
					searchUpd.SetReturnFields(append(searchUpd.ReturnFields(), "mac", "comment", "deny_bootp", "pxe_lease_time"))
					err = connector.GetObject(searchUpd, "", qp, &resUpd)
					Expect(err).To(BeNil())
					Expect(*resUpd[0].Comment).To(Equal("changed it"))
					Expect(*resUpd[0].NetworkView).To(Equal("dhcpview"))
					Expect(*resUpd[0].Mac).To(Equal("11:11:11:11:11:18"))
					Expect(*resUpd[0].Ipv4Addr).To(Equal("92.0.0.2"))
					Expect(*resUpd[0].DenyBootp).To(Equal(true))
					Expect(*resUpd[0].PxeLeaseTime).To(Equal(uint32(10)))
					Expect(resUpd[0].Ref).To(MatchRegexp("^fixedaddress.*92\\.0\\.0\\.2/dhcpview"))
				},
			)

			It("Should get the Network [92.0.0.0/8]", Label("ID: 81", "RO"), func() {
				var res []ibclient.Ipv4Network
				search := &ibclient.Ipv4Network{}
				qp := ibclient.NewQueryParams(false, map[string]string{"network": "92.0.0.0/8"})
				err := connector.GetObject(search, "", qp, &res)
				Expect(err).To(BeNil())
				Expect(res[0].Ref).To(MatchRegexp("92\\.0\\.0\\.0/8/dhcpview$"))
				Expect(*res[0].Comment).To(Equal("Add ipv4network through WAPI"))
				Expect(*res[0].Network).To(Equal("92.0.0.0/8"))
				Expect(res[0].NetworkView).To(Equal("dhcpview"))
			})

			It("Should modify the Network [92.0.0.0/8] of the fields [comment, bootfile, use_bootfile]",
				Label("ID: 118", "ID: 140", "RW"), func() {
					r := &ibclient.Ipv4Network{
						Comment:     utils.StringPtr("modified comment"),
						UseBootfile: utils.BoolPtr(true),
						Bootfile:    utils.StringPtr("boot_file"),
					}
					ref, err := connector.UpdateObject(r, refNetwork)
					Expect(err).To(BeNil())
					Expect(ref).To(MatchRegexp("^network.*92\\.0\\.0\\.0/8/dhcpview$"))

					// Get the Network [92.0.0.0/8] to validate the above modified fields
					var res []ibclient.Ipv4Network
					search := &ibclient.Ipv4Network{}
					search.SetReturnFields([]string{"comment", "bootfile", "use_bootfile"})
					qp := ibclient.NewQueryParams(false, map[string]string{"network": "92.0.0.0/8"})
					err = connector.GetObject(search, "", qp, &res)
					Expect(err).To(BeNil())
					Expect(res[0].Ref).To(MatchRegexp("92\\.0\\.0\\.0/8/dhcpview$"))
					Expect(*res[0].Comment).To(Equal("modified comment"))
					Expect(*res[0].UseBootfile).To(Equal(true))
					Expect(*res[0].Bootfile).To(Equal("boot_file"))
				},
			)

			It("Should delete the IPv4 Network [92.0.0.0/8]", Label("ID: 151", "RW"), func() {
				ref, err := connector.DeleteObject(refNetwork)
				Expect(err).To(BeNil())
				Expect(ref).To(MatchRegexp("^network.*92\\.0\\.0\\.0/8/dhcpview$"))
			})
		})

		It("Should modify the Network View [dhcpview] of the fields [comment, name]",
			Label("ID: 117", "ID: 141", "RW"), func() {
				r := &ibclient.NetworkView{
					Comment: utils.StringPtr("Modified Network View"),
					Name:    utils.StringPtr("dhcpview"),
				}
				ref, err := connector.UpdateObject(r, viewRef)
				Expect(err).To(BeNil())
				Expect(ref).To(MatchRegexp("^networkview.*dhcpview/false$"))

				// Get the Network View [dhcpview] to validate the above modified fields
				var res []ibclient.NetworkView
				search := &ibclient.NetworkView{}
				search.SetReturnFields([]string{"comment", "name"})
				qp := ibclient.NewQueryParams(false, map[string]string{"name": "dhcpview"})
				err = connector.GetObject(search, "", qp, &res)
				Expect(err).To(BeNil())
				Expect(res[0].Ref).To(MatchRegexp("^networkview.*dhcpview/false$"))
				Expect(*res[0].Comment).To(Equal("Modified Network View"))
				Expect(*res[0].Name).To(Equal("dhcpview"))
			},
		)

		It("Should delete the Custom Network View [dhcpview]", Label("ID: 157", "RW"), func() {
			ref, err := connector.DeleteObject(viewRef)
			Expect(err).To(BeNil())
			Expect(ref).To(MatchRegexp("^networkview.*dhcpview/false$"))
		})
	})

	It("Should add IPv6 Network [1::/16]", Label("ID: 50", "RW"), func() {
		n := &ibclient.Ipv6Network{
			AutoCreateReversezone: false,
			Comment:               utils.StringPtr("Add ipv6network through WAPI"),
			Network:               utils.StringPtr("1::/16"),
			NetworkView:           utils.StringPtr("default"),
		}
		ref, err := connector.CreateObject(n)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^ipv6network.*1%3A%3A/16/default$"))
	})

	When("IPv6 Network [1::/16] exists", Label("RW"), func() {
		var refNetwork string
		BeforeEach(func() {
			n := &ibclient.Ipv6Network{
				AutoCreateReversezone: false,
				Comment:               utils.StringPtr("Add ipv6network through WAPI"),
				Network:               utils.StringPtr("1::/16"),
				NetworkView:           utils.StringPtr("default"),
			}
			var err error
			refNetwork, err = connector.CreateObject(n)
			Expect(err).To(BeNil())
			Expect(refNetwork).To(MatchRegexp("^ipv6network.*1%3A%3A/16/default$"))
		})

		It("Should add IPv6 Range [start_addr = 1::1; end_addr = 1::20]", Label("ID: 51", "RW"), func() {
			r := &ibclient.IPv6Range{
				StartAddr: utils.StringPtr("1::1"),
				EndAddr:   utils.StringPtr("1::20"),
				Network:   utils.StringPtr("1::/16"),
				Comment:   utils.StringPtr("Add Range through WAPI"),
			}
			ref, err := connector.CreateObject(r)
			Expect(err).To(BeNil())
			Expect(ref).To(MatchRegexp("^ipv6range.*1%3A%3A1/1%3A%3A20/default$"))
		})

		It("Should add ipv6fixedaddress [1::50]", Label("ID: 56", "RW"), func() {
			fa := &ibclient.Ipv6FixedAddress{
				Duid:        utils.StringPtr("ab:34:56:78:90"),
				NetworkView: utils.StringPtr("default"),
				Ipv6Addr:    utils.StringPtr("1::50"),
			}
			ref, err := connector.CreateObject(fa)
			Expect(err).To(BeNil())
			Expect(ref).To(MatchRegexp("^ipv6fixedaddress.*1%3A%3A50/default"))
		})

		It("Should get the IPv6 Network [1::/16]", Label("ID: 79", "RO"), func() {
			var res []ibclient.Ipv6Network
			search := &ibclient.Ipv6Network{}
			qp := ibclient.NewQueryParams(false, map[string]string{"network": "1::/16"})
			err := connector.GetObject(search, "", qp, &res)
			Expect(err).To(BeNil())
			Expect(*res[0].Comment).To(Equal("Add ipv6network through WAPI"))
			Expect(*res[0].NetworkView).To(Equal("default"))
			Expect(res[0].Ref).To(MatchRegexp("1%3A%3A/16/default$"))
			Expect(*res[0].Network).To(Equal("1::/16"))
		})

		When("IPv6 Range [start_addr = 1::1; end_addr = 1::20] exits", Label("RW"), func() {
			var refRange string
			BeforeEach(func() {
				r := &ibclient.IPv6Range{
					StartAddr: utils.StringPtr("1::1"),
					EndAddr:   utils.StringPtr("1::20"),
					Network:   utils.StringPtr("1::/16"),
					Comment:   utils.StringPtr("Add Range through WAPI"),
				}
				var err error
				refRange, err = connector.CreateObject(r)
				Expect(err).To(BeNil())
				Expect(refRange).To(MatchRegexp("^ipv6range.*1%3A%3A1/1%3A%3A20/default$"))
			})

			It("Should get the IPAM IPv6Address object", Label("ID: 63", "RO"), func() {
				var res []ibclient.IPv6Address
				search := &ibclient.IPv6Address{}
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

			It("Should get the IPv6 Range [1::1/1::20] using reference with all default fields", Label("ID: 78", "RO"), func() {
				var res []ibclient.IPv6Range
				search := &ibclient.IPv6Range{}
				err := connector.GetObject(search, "", nil, &res)
				Expect(err).To(BeNil())
				Expect(*res[0].Comment).To(Equal("Add Range through WAPI"))
				Expect(*res[0].Network).To(Equal("1::/16"))
				Expect(*res[0].NetworkView).To(Equal("default"))
				Expect(*res[0].StartAddr).To(Equal("1::1"))
				Expect(*res[0].EndAddr).To(Equal("1::20"))
				Expect(res[0].Ref).To(MatchRegexp("1%3A%3A1/1%3A%3A20/default$"))
			})

			It("Should modify the IPv6 Range [1::1/1::20] of the fields [comment, domain_name, use_domain_name]",
				Label("ID: 121", "ID: 137", "RW"), func() {
					r := &ibclient.IPv6Range{
						Comment:          utils.StringPtr("modified comment"),
						UseRecycleLeases: utils.BoolPtr(true),
						RecycleLeases:    utils.BoolPtr(true),
					}
					ref, err := connector.UpdateObject(r, refRange)
					Expect(err).To(BeNil())
					Expect(ref).To(MatchRegexp("^ipv6range.*1%3A%3A1/1%3A%3A20/default$"))

					// Get the IPv6 Range [1::1/1::20] to validate the above modified fields
					var res []ibclient.IPv6Range
					search := &ibclient.IPv6Range{}
					search.SetReturnFields([]string{"comment", "recycle_leases", "use_recycle_leases"})
					qp := ibclient.NewQueryParams(
						false,
						map[string]string{"start_addr": "1::1", "end_addr": "1::20"},
					)
					err = connector.GetObject(search, "", qp, &res)
					Expect(err).To(BeNil())
					Expect(*res[0].Comment).To(Equal("modified comment"))
					Expect(*res[0].UseRecycleLeases).To(Equal(true))
					Expect(*res[0].RecycleLeases).To(Equal(true))
				},
			)
		})

		When("ipv6fixedaddress [1::50] exists", Label("RW"), Label("RW"), func() {
			var refFixedAddress string
			BeforeEach(func() {
				fa := &ibclient.Ipv6FixedAddress{
					Duid:        utils.StringPtr("ab:34:56:78:90"),
					NetworkView: utils.StringPtr("default"),
					Ipv6Addr:    utils.StringPtr("1::50"),
				}
				var err error
				refFixedAddress, err = connector.CreateObject(fa)
				Expect(err).To(BeNil())
				Expect(refFixedAddress).To(MatchRegexp("^ipv6fixedaddress.*1%3A%3A50/default"))
			})

			It("Should get ipv6fixedaddress with default return fields", Label("ID: 74", "RO"), func() {
				var res []ibclient.Ipv6FixedAddress
				search := &ibclient.Ipv6FixedAddress{}
				qp := ibclient.NewQueryParams(false, map[string]string{"ipv6addr": "1::50"})
				err := connector.GetObject(search, "", qp, &res)
				Expect(err).To(BeNil())
				Expect(*res[0].Duid).To(Equal("ab:34:56:78:90"))
				Expect(*res[0].NetworkView).To(Equal("default"))
				Expect(res[0].Ref).To(MatchRegexp("ipv6fixedaddress.*1%3A%3A50/default$"))
				Expect(*res[0].Ipv6Addr).To(Equal("1::50"))
			})

			It("Should modify the fields [comment, use_preferred_lifetime, preferred_lifetime] in ipv6fixedaddress [1::50]",
				Label("ID: 125", "ID: 133", "RW"), func() {
					r := &ibclient.Ipv6FixedAddress{
						Comment:              utils.StringPtr("Modify the ipv6fixedaddress object"),
						UsePreferredLifetime: utils.BoolPtr(false),
						PreferredLifetime:    utils.Uint32Ptr(220),
					}
					ref, err := connector.UpdateObject(r, refFixedAddress)
					Expect(err).To(BeNil())
					Expect(ref).To(MatchRegexp("^ipv6fixedaddress.*1%3A%3A50/default$"))

					// Using GMC IP, Get the ipv6fixedaddress [1::50] to validate
					// the above modified fields with _return_fields+ method
					var res []ibclient.Ipv6FixedAddress
					search := &ibclient.Ipv6FixedAddress{}
					qp := ibclient.NewQueryParams(false, map[string]string{"ipv6addr": "1::50"})
					search.SetReturnFields(append(search.ReturnFields(), "comment", "use_preferred_lifetime", "preferred_lifetime"))
					err = connector.GetObject(search, "", qp, &res)
					Expect(err).To(BeNil())
					Expect(*res[0].Comment).To(Equal("Modify the ipv6fixedaddress object"))
					Expect(*res[0].NetworkView).To(Equal("default"))
					Expect(*res[0].Duid).To(Equal("ab:34:56:78:90"))
					Expect(*res[0].UsePreferredLifetime).To(Equal(false))
					Expect(*res[0].PreferredLifetime).To(Equal(uint32(220)))
					Expect(res[0].Ref).To(MatchRegexp("^ipv6fixedaddress.*1%3A%3A50/default$"))
					Expect(*res[0].Ipv6Addr).To(Equal("1::50"))
				},
			)
		})

		It("Should modify the IPv6 Network [1::/16] of the fields [comment, domain_name, use_domain_name]",
			Label("ID: 120", "ID: 138", "RW"), func() {
				r := &ibclient.Ipv6Network{
					Comment:       utils.StringPtr("modified comment"),
					UseDomainName: utils.BoolPtr(true),
					DomainName:    utils.StringPtr("boot_file"),
				}
				ref, err := connector.UpdateObject(r, refNetwork)
				Expect(err).To(BeNil())
				Expect(ref).To(MatchRegexp("^ipv6network.*1%3A%3A/16/default$"))

				// Get the IPv6 Network [1::/16] to validate the above modified fields
				var res []ibclient.Ipv6Network
				search := &ibclient.Ipv6Network{}
				qp := ibclient.NewQueryParams(false, map[string]string{"network": "1::/16"})
				search.SetReturnFields([]string{"comment", "domain_name", "use_domain_name"})
				err = connector.GetObject(search, "", qp, &res)
				Expect(err).To(BeNil())
				Expect(*res[0].Comment).To(Equal("modified comment"))
				Expect(*res[0].UseDomainName).To(Equal(true))
				Expect(*res[0].DomainName).To(Equal("boot_file"))
			},
		)

		It("Should delete the IPv6 Network [1::/16]", Label("ID: 153", "RW"), func() {
			ref, err := connector.DeleteObject(refNetwork)
			Expect(err).To(BeNil())
			Expect(ref).To(MatchRegexp("^ipv6network.*1%3A%3A/16/default$"))
		})
	})

	It("Should add Network Container [78.0.0.0/8]", Label("ID: 52", "RW"), func() {
		nc := &ibclient.Ipv4NetworkContainer{
			AutoCreateReversezone: false,
			Comment:               utils.StringPtr("Add networkcontainer through WAPI"),
			Network:               "78.0.0.0/8",
			NetworkView:           "default",
		}
		ref, err := connector.CreateObject(nc)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^networkcontainer.*78\\.0\\.0\\.0/8/default$"))
	})

	When("Network Container [78.0.0.0/8] exits", Label("RW"), func() {
		var refNetworkContainer string
		BeforeEach(func() {
			nc := &ibclient.Ipv4NetworkContainer{
				AutoCreateReversezone: false,
				Comment:               utils.StringPtr("Add networkcontainer through WAPI"),
				Network:               "78.0.0.0/8",
				NetworkView:           "default",
			}
			var err error
			refNetworkContainer, err = connector.CreateObject(nc)
			Expect(err).To(BeNil())
			Expect(refNetworkContainer).To(MatchRegexp("^networkcontainer.*78\\.0\\.0\\.0/8/default$"))
		})

		It("Should add Network [78.0.0.0/30] to the Network Container [78.0.0.0/8]", Label("ID: 53", "RW"), func() {
			n := &ibclient.Ipv4Network{
				Comment:     utils.StringPtr("Add a network to container 78.0.0.0/8 through WAPI"),
				Network:     utils.StringPtr("78.0.0.0/30"),
				NetworkView: "default",
			}
			ref, err := connector.CreateObject(n)
			Expect(err).To(BeNil())
			Expect(ref).To(MatchRegexp("^network.*78\\.0\\.0\\.0/30/default$"))
		})

		When("Network [78.0.0.0/30] exists", Label("RW"), func() {
			var networkRef string
			BeforeEach(func() {
				n := &ibclient.Ipv4Network{
					Comment:     utils.StringPtr("Add a network to container 78.0.0.0/8 through WAPI"),
					Network:     utils.StringPtr("78.0.0.0/30"),
					NetworkView: "default",
				}
				var err error
				networkRef, err = connector.CreateObject(n)
				Expect(err).To(BeNil())
				Expect(networkRef).To(MatchRegexp("^network.*78\\.0\\.0\\.0/30/default$"))
			})

			It("Should get the IPAM IPv4Address object", Label("ID: 62", "RO"), func() {
				var res []ibclient.IPv4Address
				search := &ibclient.IPv4Address{}
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
				Expect(res[0].Ref).To(MatchRegexp("ipv4address.*78\\.0\\.0\\.1"))
				Expect(res[0].IpAddress).To(Equal("78.0.0.1"))
				Expect(res[0].Names).To(HaveLen(0))
			})

			It("Should delete the IPv4 Network [78.0.0.0/30]", Label("ID: 152", "RW"), func() {
				ref, err := connector.DeleteObject(networkRef)
				Expect(err).To(BeNil())
				Expect(ref).To(MatchRegexp("^network.*78\\.0\\.0\\.0/30/default$"))
			})
		})

		It("Should get the Network Container [78.0.0.0/8] using reference with all return fields",
			Label("ID: 77", "RO"), func() {
				var res []ibclient.Ipv4NetworkContainer
				search := &ibclient.Ipv4NetworkContainer{}
				search.SetReturnFields([]string{"comment", "network", "network_view", "network_container"})
				err := connector.GetObject(search, "", nil, &res)
				Expect(err).To(BeNil())
				Expect(*res[0].Comment).To(Equal("Add networkcontainer through WAPI"))
				Expect(res[0].NetworkContainer).To(Equal("/"))
				Expect(res[0].Network).To(Equal("78.0.0.0/8"))
				Expect(res[0].NetworkView).To(Equal("default"))
				Expect(res[0].Ref).To(MatchRegexp("78\\.0\\.0\\.0/8/default$"))
			},
		)

		It("Should modify the Network Container [78.0.0.0/8] of the field [comment]",
			Label("ID: 122", "ID: 136", "RW"), func() {
				nc := &ibclient.Ipv4NetworkContainer{
					Comment: utils.StringPtr("modified comment"),
				}
				ref, err := connector.UpdateObject(nc, refNetworkContainer)
				Expect(err).To(BeNil())
				Expect(ref).To(MatchRegexp("^networkcontainer.*78\\.0\\.0\\.0/8/default$"))

				// Get the Network Container [78.0.0.0/8] to validate the above modified field
				var res []ibclient.Ipv4NetworkContainer
				search := &ibclient.Ipv4NetworkContainer{}
				search.SetReturnFields([]string{"comment"})
				qp := ibclient.NewQueryParams(false, map[string]string{"network": "78.0.0.0/8"})
				err = connector.GetObject(search, "", qp, &res)
				Expect(err).To(BeNil())
				Expect(*res[0].Comment).To(Equal("modified comment"))
				Expect(res[0].Ref).To(MatchRegexp("^networkcontainer.*78\\.0\\.0\\.0/8/default$"))
			},
		)

		It("Should delete the IPv4 Network Container [78.0.0.0/8]", Label("ID: 154", "RW"), func() {
			ref, err := connector.DeleteObject(refNetworkContainer)
			Expect(err).To(BeNil())
			Expect(ref).To(MatchRegexp("^networkcontainer.*78\\.0\\.0\\.0/8/default$"))
		})
	})

	It("Should add IPv6 Network Container [2000::/64]", Label("ID: 54", "RW"), func() {
		n := &ibclient.Ipv6NetworkContainer{
			AutoCreateReversezone: false,
			Comment:               utils.StringPtr("Add ipv6networkcontainer through WAPI"),
			Network:               "2000::/64",
			NetworkView:           utils.StringPtr("default"),
		}
		ref, err := connector.CreateObject(n)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^ipv6network.*2000%3A%3A/64/default"))
	})

	When("IPv6 Network Container [2000::/64] exists", Label("RW"), func() {
		var refNetworkContainer string
		BeforeEach(func() {
			n := &ibclient.Ipv6NetworkContainer{
				AutoCreateReversezone: false,
				Comment:               utils.StringPtr("Add ipv6networkcontainer through WAPI"),
				Network:               "2000::/64",
				NetworkView:           utils.StringPtr("default"),
			}
			var err error
			refNetworkContainer, err = connector.CreateObject(n)
			Expect(err).To(BeNil())
			Expect(refNetworkContainer).To(MatchRegexp("^ipv6network.*2000%3A%3A/64/default"))
		})

		It("Should get the IPv6 Network Container [2000::/64] using reference with all return fields",
			Label("ID: 76", "RO"), func() {
				var res []ibclient.Ipv6NetworkContainer
				search := &ibclient.Ipv6NetworkContainer{}
				search.SetReturnFields([]string{"comment", "network", "network_view", "network_container"})
				err := connector.GetObject(search, "", nil, &res)
				Expect(err).To(BeNil())
				Expect(*res[0].Comment).To(Equal("Add ipv6networkcontainer through WAPI"))
				Expect(res[0].NetworkContainer).To(Equal("/"))
				Expect(res[0].Network).To(Equal("2000::/64"))
				Expect(*res[0].NetworkView).To(Equal("default"))
				Expect(res[0].Ref).To(MatchRegexp("2000%3A%3A/64/default$"))
			},
		)

		It("Should modify the IPv6 Network Container [2000::/64] of the field [comment]",
			Label("ID: 123", "ID: 135", "RW"), func() {
				nc := &ibclient.Ipv6NetworkContainer{
					Comment: utils.StringPtr("modified comment"),
				}
				ref, err := connector.UpdateObject(nc, refNetworkContainer)
				Expect(err).To(BeNil())
				Expect(ref).To(MatchRegexp("^ipv6network.*2000%3A%3A/64/default"))

				// Get the IPv6 Network Container [2000::/64] to validate the above modified field
				var res []ibclient.Ipv6NetworkContainer
				search := &ibclient.Ipv6NetworkContainer{}
				search.SetReturnFields([]string{"comment"})
				qp := ibclient.NewQueryParams(false, map[string]string{"network": "2000::/64"})
				err = connector.GetObject(search, "", qp, &res)
				Expect(err).To(BeNil())
				Expect(*res[0].Comment).To(Equal("modified comment"))
				Expect(res[0].Ref).To(MatchRegexp("2000%3A%3A/64/default$"))
			},
		)

		It("Should delete the IPv6 Network Container [2000::/64]", Label("ID: 155", "RW"), func() {
			ref, err := connector.DeleteObject(refNetworkContainer)
			Expect(err).To(BeNil())
			Expect(ref).To(MatchRegexp("^ipv6network.*2000%3A%3A/64/default"))
		})
	})

	It("Should add DNS view [view1]", Label("ID: 57", "ID: 126", "ID: 132", "RW"), func() {
		v := &ibclient.View{
			Name: utils.StringPtr("view1"),
		}
		ref, err := connector.CreateObject(v)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^view.*view1/false$"))

		// Modify the fields [comment, disable, recursion] in DNS view [view1]
		v = &ibclient.View{
			Comment:   utils.StringPtr("Modify the view object"),
			Disable:   utils.BoolPtr(true),
			Recursion: utils.BoolPtr(true),
		}
		ref, err = connector.UpdateObject(v, ref)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^view.*view1/false$"))

		// Get the DNS view [view1] to validate the above modified fields
		var res []ibclient.View
		search := &ibclient.View{}
		qp := ibclient.NewQueryParams(false, map[string]string{"name": "view1"})
		err = connector.GetObject(search, "", qp, &res)
		Expect(err).To(BeNil())
		Expect(*res[0].Comment).To(Equal("Modify the view object"))
		Expect(res[0].IsDefault).To(Equal(false))
		Expect(res[0].Ref).To(MatchRegexp("^view.*view1/false$"))
		Expect(*res[0].Name).To(Equal("view1"))
	})

	It("Should add Named ACL without 'access_list' field [wapi-na2]",
		Label("ID: 58", "ID: 60", "ID: 128", "ID: 130", "ID: 150", "RW"), func() {
			n := &ibclient.Namedacl{
				Name:    utils.StringPtr("wapi-na2"),
				Comment: utils.StringPtr("No acls present"),
			}
			ref, err := connector.CreateObject(n)
			Expect(err).To(BeNil())
			Expect(ref).To(MatchRegexp("^namedacl.*wapi-na2$"))

			var res []ibclient.Namedacl
			search := &ibclient.Namedacl{}
			err = connector.GetObject(search, "", nil, &res)
			Expect(err).To(BeNil())
			Expect(*res[0].Comment).To(Equal("No acls present"))
			Expect(res[0].Ref).To(MatchRegexp("wapi-na2"))
			Expect(*res[0].Name).To(Equal("wapi-na2"))

			// Modify the fields [name,comment] for the Named ACL [wapi-na2]
			n = &ibclient.Namedacl{
				Name:    utils.StringPtr("wapi-mod"),
				Comment: utils.StringPtr("modified"),
			}
			ref, err = connector.UpdateObject(n, ref)
			Expect(err).To(BeNil())
			Expect(ref).To(MatchRegexp("^namedacl.*wapi-mod$"))

			// Using GMC IP, Get the Named ACL [wapi-mod] to validate the modified fields [name,comment]
			var resMod []ibclient.Namedacl
			search = &ibclient.Namedacl{}
			//qp := ibclient.NewQueryParams(false, map[string]string{"name": "wapi-mod"})
			search.SetReturnFields([]string{"name", "comment"})
			err = connector.GetObject(search, "", nil, &resMod)
			Expect(err).To(BeNil())
			Expect(*resMod[0].Comment).To(Equal("modified"))
			Expect(resMod[0].Ref).To(MatchRegexp("^namedacl.*wapi-mod$"))
			Expect(*resMod[0].Name).To(Equal("wapi-mod"))

			// Delete Named ACL [wapi-mod]
			ref, err = connector.DeleteObject(ref)
			Expect(err).To(BeNil())
			Expect(ref).To(MatchRegexp("^namedacl.*wapi-mod$"))
		},
	)

	It("Should get the DHCP network template object", Label("ID: 64", "RO"), func() {
		var res []ibclient.NetworkTemplate
		search := &ibclient.NetworkTemplate{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError(ibclient.NewNotFoundError("requested object not found")))
	})

	It("Should get the DHCP IPv6 network template object", Label("ID: 65", "RO"), func() {
		var res []ibclient.IPv6NetworkTemplate
		search := &ibclient.IPv6NetworkTemplate{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError(ibclient.NewNotFoundError("requested object not found")))
	})

	It("Should get the DHCP IPv6 Shared Network object", Label("ID: 66", "RO"), func() {
		var res []ibclient.IPv6SharedNetwork
		search := &ibclient.IPv6SharedNetwork{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError(ibclient.NewNotFoundError("requested object not found")))
	})

	It("Should get the DHCP Lease object", Label("ID: 67", "RO"), func() {
		var res []ibclient.Lease
		search := &ibclient.Lease{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError(ibclient.NewNotFoundError("requested object not found")))
	})

	It("Should get the MAC Filter Address object", Label("ID: 68", "RO"), func() {
		var res []ibclient.MACFilterAddress
		search := &ibclient.MACFilterAddress{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError(ibclient.NewNotFoundError("requested object not found")))
	})

	It("Should get the Member DHCP properties object", Label("ID: 69", "RO"), func() {
		var res []ibclient.MemberDHCPProperties
		search := &ibclient.MemberDHCPProperties{}
		search.SetReturnFields([]string{"host_name"})
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(BeNil())
		Expect(res[0].HostName).To(HavePrefix("infoblox."))
		Expect(res[0].Ref).To(MatchRegexp("^member:dhcpproperties.*infoblox\\..*"))
	})

	It("Should get the Member DNS object", Label("ID: 70", "RO"), func() {
		var res []ibclient.MemberDns
		search := &ibclient.MemberDns{}
		search.SetReturnFields([]string{"host_name"})
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(BeNil())
		Expect(res[0].HostName).To(HavePrefix("infoblox."))
		Expect(res[0].Ref).To(MatchRegexp("^member:dns.*infoblox\\..*"))
	})

	It("Should get the Active Directory Domain object", Label("ID: 71", "RO"), func() {
		var res []ibclient.MsserverAdsitesDomain
		search := &ibclient.MsserverAdsitesDomain{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError(ibclient.NewNotFoundError("requested object not found")))
	})

	It("Should get the Active Directory Site object", Label("ID: 72", "RO"), func() {
		var res []ibclient.MsserverAdsitesSite
		search := &ibclient.MsserverAdsitesSite{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError(ibclient.NewNotFoundError("requested object not found")))
	})

	It("Should get the Permissions object", Label("ID: 73", "RO"), func() {
		var res []ibclient.Permission
		search := &ibclient.Permission{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(BeNil())
		Expect(res[0].Permission).To(Equal("WRITE"))
		Expect(*res[0].Role).To(Equal("DNS Admin"))
		Expect(res[0].Ref).To(MatchRegexp("^permission.*Admin/WRITE$"))
		Expect(res[0].ResourceType).To(Equal("VIEW"))
	})

	It("Should get the DNS NAPTR record object", Label("ID: 89", "RO"), func() {
		var res []ibclient.RecordNaptr
		search := &ibclient.RecordNaptr{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError(ibclient.NewNotFoundError("requested object not found")))
	})

	It("Should get the DTC LBDN object", Label("ID: 90", "RO"), func() {
		var res []ibclient.RecordDtclbdn
		search := &ibclient.RecordDtclbdn{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError(ibclient.NewNotFoundError("requested object not found")))
	})

	It("Should get the DTC LBDN object", Label("ID: 90", "RO"), func() {
		var res []ibclient.RecordDtclbdn
		search := &ibclient.RecordDtclbdn{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError(ibclient.NewNotFoundError("requested object not found")))
	})

	It("Should get the DHCP Roaming Host object", Label("ID: 94", "RO"), func() {
		var res []ibclient.RoamingHost
		search := &ibclient.RoamingHost{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError(ibclient.NewNotFoundError("requested object not found")))
	})

	It("Should get the Scheduled Task object", Label("ID: 95", "RO"), func() {
		var res []ibclient.ScheduledTask
		search := &ibclient.ScheduledTask{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError(ibclient.NewNotFoundError("requested object not found")))
	})

	It("Should get the Search object (N)", Label("ID: 96", "RO"), func() {
		var res []ibclient.Search
		search := &ibclient.Search{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).NotTo(BeNil())
		// TODO Check the error string
	})

	It("Should get the DHCP Shared Network object", Label("ID: 97", "RO"), func() {
		var res []ibclient.SharedNetwork
		search := &ibclient.SharedNetwork{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError(ibclient.NewNotFoundError("requested object not found")))
	})

	It("Should get the DHCP Shared A record object", Label("ID: 98", "RO"), func() {
		var res []ibclient.SharedRecordA
		search := &ibclient.SharedRecordA{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError(ibclient.NewNotFoundError("requested object not found")))
	})

	It("Should get the DHCP Shared AAAA record object", Label("ID: 99", "RO"), func() {
		var res []ibclient.SharedRecordAAAA
		search := &ibclient.SharedRecordAAAA{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError(ibclient.NewNotFoundError("requested object not found")))
	})

	It("Should get the DHCP Shared MX record object", Label("ID: 100", "RO"), func() {
		var res []ibclient.SharedRecordMX
		search := &ibclient.SharedRecordMX{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError(ibclient.NewNotFoundError("requested object not found")))
	})

	It("Should get the DHCP Shared TXT record object", Label("ID: 101", "RO"), func() {
		var res []ibclient.SharedRecordTXT
		search := &ibclient.SharedRecordTXT{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError(ibclient.NewNotFoundError("requested object not found")))
	})

	It("Should get the DHCP Shared SRV record object", Label("ID: 102", "RO"), func() {
		var res []ibclient.SharedrecordSrv
		search := &ibclient.SharedrecordSrv{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError(ibclient.NewNotFoundError("requested object not found")))
	})

	It("Should get the snmpuser object (N)", Label("ID: 103", "RO"), func() {
		var res []ibclient.SNMPUser
		search := &ibclient.SNMPUser{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).NotTo(BeNil())
		// TODO Check the error string
	})

	It("Should get the Zone discrepancy information object", Label("ID: 104", "RO"), func() {
		var res []ibclient.ZoneAuthDiscrepancy
		search := &ibclient.ZoneAuthDiscrepancy{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError(ibclient.NewNotFoundError("requested object not found")))
	})

	It("Should get the DNS Delegated Zone object", Label("ID: 105", "RO"), func() {
		var res []ibclient.ZoneDelegated
		search := &ibclient.ZoneDelegated{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError(ibclient.NewNotFoundError("requested object not found")))
	})

	It("Should get the DNS Forward Zone object", Label("ID: 106", "RO"), func() {
		var res []ibclient.ZoneForward
		search := &ibclient.ZoneForward{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError(ibclient.NewNotFoundError("requested object not found")))
	})

	It("Should get the DNS Stub Zone object", Label("ID: 107", "RO"), func() {
		var res []ibclient.ZoneStub
		search := &ibclient.ZoneStub{}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).To(MatchError(ibclient.NewNotFoundError("requested object not found")))
	})
})

var _ = Describe("DNS Forward Zone", func() {
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

	It("Should create a DNS Forward Zone", func() {
		// Create a DNS Forward Zone
		zone := &ibclient.ZoneForward{
			Fqdn: "example.com",
			ForwardTo: ibclient.NullableNameServers{
				NameServers: []ibclient.NameServer{
					{Name: "test", Address: "1.2.3.4"},
					{Name: "test2", Address: "2.3.4.5"},
				},
				IsNull: false,
			},
		}
		ref, err := connector.CreateObject(zone)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^zone_forward.*"))
	})

	It("Should create a DNS Forward Zone with all params", func() {
		// Create a DNS Forward Zone
		zone := &ibclient.ZoneForward{
			Fqdn: "example.com",
			ForwardTo: ibclient.NullableNameServers{
				NameServers: []ibclient.NameServer{
					{Name: "test", Address: "1.2.3.4"},
					{Name: "test2", Address: "2.3.4.5"},
				},
				IsNull: false,
			},
			Comment:        utils.StringPtr("wapi added"),
			ForwardersOnly: utils.BoolPtr(true),
			ForwardingServers: &ibclient.NullableForwardingServers{
				[]*ibclient.Forwardingmemberserver{
					{
						Name:           "infoblox.localdomain",
						ForwardersOnly: true,
						ForwardTo: ibclient.NullableNameServers{
							NameServers: []ibclient.NameServer{
								{Name: "test", Address: "1.2.3.4"},
								{Name: "test2", Address: "2.3.4.5"},
							},
							IsNull: false,
						},
						UseOverrideForwarders: false,
					}},
				false,
			},
		}
		ref, err := connector.CreateObject(zone)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^zone_forward.*"))
	})

	It("Should get the DNS Forward Zone", func() {
		zone := &ibclient.ZoneForward{
			Fqdn: "example.com",
			ForwardTo: ibclient.NullableNameServers{
				NameServers: []ibclient.NameServer{
					{Name: "test", Address: "1.2.3.4"},
					{Name: "test2", Address: "2.3.4.5"},
				},
				IsNull: false,
			},
			Comment: utils.StringPtr("wapi added"),
		}
		ref, err := connector.CreateObject(zone)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^zone_forward.*"))

		var res []ibclient.ZoneForward
		search := &ibclient.ZoneForward{}
		errCode := connector.GetObject(search, "", nil, &res)
		Expect(errCode).To(BeNil())
		Expect(res[0].Ref).To(MatchRegexp("^zone_forward.*"))
	})

	It("Should update the DNS Forward Zone", func() {
		// Create a DNS Forward Zone
		zoneCreate := &ibclient.ZoneForward{
			Fqdn: "example.com",
			ForwardTo: ibclient.NullableNameServers{
				NameServers: []ibclient.NameServer{
					{Name: "test", Address: "1.2.3.4"},
					{Name: "test2", Address: "1.2.3.5"},
				},
				IsNull: false,
			},
			Comment: utils.StringPtr("wapi added"),
		}
		ref, errCode := connector.CreateObject(zoneCreate)
		Expect(errCode).To(BeNil())
		Expect(ref).To(MatchRegexp("^zone_forward.*"))

		// Update a DNS Forward Zone
		zone := &ibclient.ZoneForward{
			ForwardTo: ibclient.NullableNameServers{
				NameServers: []ibclient.NameServer{
					{Name: "test", Address: "1.2.3.4"},
					{Name: "test2", Address: "1.2.3.6"},
				},
				IsNull: false,
			},
			Comment: utils.StringPtr("wapi added"),
		}

		var res []ibclient.ZoneForward
		search := &ibclient.ZoneForward{}
		err := connector.GetObject(search, "", nil, &res)
		ref, err = connector.UpdateObject(zone, res[0].Ref)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^zone_forward.*"))
	})

	It("Should delete the DNS Forward Zone", func() {
		// Create a DNS Forward Zone
		zoneCreate := &ibclient.ZoneForward{
			Fqdn: "example.com",
			ForwardTo: ibclient.NullableNameServers{
				NameServers: []ibclient.NameServer{
					{Name: "test", Address: "1.2.3.4"},
					{Name: "test2", Address: "1.2.3.5"},
				},
				IsNull: false,
			},
			Comment: utils.StringPtr("wapi added"),
		}
		refCreate, errCode := connector.CreateObject(zoneCreate)
		Expect(errCode).To(BeNil())
		Expect(refCreate).To(MatchRegexp("^zone_forward.*"))

		var res []ibclient.ZoneForward
		search := &ibclient.ZoneForward{}
		err := connector.GetObject(search, "", nil, &res)
		ref, err := connector.DeleteObject(res[0].Ref)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^zone_forward.*"))
	})

	It("Should fail to create a DNS Forward Zone with invalid data", func() {
		zone := &ibclient.ZoneForward{
			Fqdn: "invalid..com", // Invalid FQDN
			ForwardTo: ibclient.NullableNameServers{
				NameServers: []ibclient.NameServer{
					{Name: "test", Address: "1.2.3.4"},
					{Name: "test2", Address: "1.2.3.5"},
				},
				IsNull: false,
			},
			Comment:        utils.StringPtr("wapi added"),
			ForwardersOnly: utils.BoolPtr(true),
		}
		_, err := connector.CreateObject(zone)
		Expect(err).NotTo(BeNil())
	})

	It("Should fail to get a non-existent DNS Forward Zone", func() {
		var res []ibclient.ZoneForward
		search := &ibclient.ZoneForward{Fqdn: "nonexistent.com"}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).NotTo(BeNil())
	})

	It("Should fail to update a non-existent DNS Forward Zone", func() {
		zone := &ibclient.ZoneForward{
			Fqdn: "nonexistent.com",
			ForwardTo: ibclient.NullableNameServers{
				NameServers: []ibclient.NameServer{
					{Name: "test", Address: "1.2.3.4"},
					{Name: "test2", Address: "1.2.3.6"},
				},
				IsNull: false,
			},
			Comment: utils.StringPtr("wapi added"),
		}

		_, err := connector.UpdateObject(zone, "nonexistent_ref")
		Expect(err).NotTo(BeNil())
	})

	It("Should fail to delete a non-existent DNS Forward Zone", func() {
		_, err := connector.DeleteObject("nonexistent_ref")
		Expect(err).NotTo(BeNil())
	})

	It("Should fail to create a DNS Forward Zone without mandatory parameters", func() {
		zone := &ibclient.ZoneForward{
			// Missing mandatory parameters like Fqdn and ForwardTo
			Comment:        utils.StringPtr("wapi added"),
			ForwardersOnly: utils.BoolPtr(true),
		}
		_, err := connector.CreateObject(zone)
		Expect(err).NotTo(BeNil())
	})

	It("Should fail to create a DNS Forward Zone without fqdn parameter", func() {
		zone := &ibclient.ZoneForward{
			// Missing mandatory parameter Fqdn
			ForwardTo: ibclient.NullableNameServers{
				NameServers: []ibclient.NameServer{
					{Name: "test", Address: "1.2.3.4"},
					{Name: "test2", Address: "1.2.3.5"},
				},
				IsNull: false,
			},
			Comment:        utils.StringPtr("wapi added"),
			ForwardersOnly: utils.BoolPtr(true),
		}
		_, err := connector.CreateObject(zone)
		Expect(err).NotTo(BeNil())
	})

	It("Should fail to create a DNS Forward Zone without forward_to parameter", func() {
		zone := &ibclient.ZoneForward{
			// Missing mandatory parameter ForwardTo
			Fqdn:           "example.com",
			Comment:        utils.StringPtr("wapi added"),
			ForwardersOnly: utils.BoolPtr(true),
		}
		_, err := connector.CreateObject(zone)
		Expect(err).NotTo(BeNil())
	})
})

var _ = Describe("Allocate next available using EA", func() {
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

		var (
			netView = "default"
			comment = "ipv4 network container"
			ea      = ibclient.EA{"Site": "Burma"}
			cidr    = "18.12.1.0/24"
			err1    error

			ipv4Comment = "ipv4 network"
			ipv4Ea      = ibclient.EA{"Site": "Namibia"}
			ipv4Cidr    = "15.12.1.0/24"

			ipv6Comment = "ipv6 network"
			ipv6Ea      = ibclient.EA{"Site": "Norway"}
			ipv6Cidr    = "2002:db8:85a4::/64"

			ipv6NCComment = "ipv6 network container"
			ipv6NCEa      = ibclient.EA{"Site": "Sri Lanka"}
			ipv6NCCidr    = "3002:db8:85a4::/64"

			ipv4Comment1 = "ipv4 network2"
			ipv4Ea1      = ibclient.EA{"Site": "Japan"}
			ipv4Cidr1    = "22.12.1.0/24"

			ipv6Comment1 = "ipv6 network2"
			ipv6Ea1      = ibclient.EA{"Site": "Japan"}
			ipv6Cidr1    = "4002:db8:85a4::/64"
		)

		nv := &ibclient.ZoneAuth{
			View: utils.StringPtr("default"),
			Fqdn: "wapi.com",
		}
		_, err1 = connector.CreateObject(nv)
		Expect(err1).To(BeNil())

		nc := ibclient.NewNetworkContainer(netView, cidr, false, comment, ea)
		_, err1 = connector.CreateObject(nc)
		Expect(err1).To(BeNil())

		ipv6Network := ibclient.NewNetwork(netView, ipv6Cidr, true, ipv6Comment, ipv6Ea)
		_, err1 = connector.CreateObject(ipv6Network)
		Expect(err1).To(BeNil())

		ipv4Network := ibclient.NewNetwork(netView, ipv4Cidr, false, ipv4Comment, ipv4Ea)
		_, err1 = connector.CreateObject(ipv4Network)
		Expect(err1).To(BeNil())

		ipv6NWC := ibclient.NewNetworkContainer(netView, ipv6NCCidr, true, ipv6NCComment, ipv6NCEa)
		_, err1 = connector.CreateObject(ipv6NWC)
		Expect(err1).To(BeNil())

		ipv6Network1 := ibclient.NewNetwork(netView, ipv6Cidr1, true, ipv6Comment1, ipv6Ea1)
		_, err1 = connector.CreateObject(ipv6Network1)
		Expect(err1).To(BeNil())

		ipv4Network1 := ibclient.NewNetwork(netView, ipv4Cidr1, false, ipv4Comment1, ipv4Ea1)
		_, err1 = connector.CreateObject(ipv4Network1)
		Expect(err1).To(BeNil())
	})

	AfterEach(func() {
		err := connector.SweepObjects()
		Expect(err).To(BeNil())
	})

	It("Should create a ipv4 network within a networkcontainer with EA", func() {
		// Create an ipv4 network within a networkcontainer with EA
		ea := ibclient.EA{"Region": "East"}
		comment := "Test ipv4 network creation with next_available_network"
		eaMap := map[string]string{"*Site": "Burma"}
		prefixLen := uint(26)
		netviewName := "default"
		networkinfo := &ibclient.NetworkContainerNextAvailable{
			Network: &ibclient.NetworkContainerNextAvailableInfo{
				Function:     "next_available_network",
				ResultField:  "networks",
				Object:       "networkcontainer",
				ObjectParams: eaMap,
				Params: map[string]uint{
					"cidr": prefixLen,
				},
				NetviewName: "",
			},
			NetviewName: netviewName,
			Comment:     comment,
			Ea:          ea,
		}
		networkinfo.SetObjectType("network")
		ref, err := connector.CreateObject(networkinfo)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^network.*"))
	})

	It("Should create a ipv4 networkconatiner within a networkcontainer with EA", func() {
		// Create an ipv4 networkcontainer within a networkcontainer with EA
		ea := ibclient.EA{"Region": "East"}
		comment := "Test ipv4 network container creation with next_available_network"
		eaMap := map[string]string{"*Site": "Burma"}
		prefixLen := uint(26)
		netviewName := "default"
		networkinfo := &ibclient.NetworkContainerNextAvailable{
			Network: &ibclient.NetworkContainerNextAvailableInfo{
				Function:     "next_available_network",
				ResultField:  "networks",
				Object:       "networkcontainer",
				ObjectParams: eaMap,
				Params: map[string]uint{
					"cidr": prefixLen,
				},
				NetviewName: "",
			},
			NetviewName: netviewName,
			Comment:     comment,
			Ea:          ea,
		}
		networkinfo.SetObjectType("networkcontainer")
		ref, err := connector.CreateObject(networkinfo)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^network.*"))
	})

	It("Should create a ipv6 network within a networkconatiner with EA", func() {
		// Create an ipv6 network within a network with EA
		ea := ibclient.EA{"Region": "East"}
		comment := "Test ipv6 network creation with next_available_network"
		eaMap := map[string]string{"*Site": "Sri Lanka"}
		prefixLen := uint(66)
		netviewName := "default"
		networkinfo := &ibclient.NetworkContainerNextAvailable{
			Network: &ibclient.NetworkContainerNextAvailableInfo{
				Function:     "next_available_network",
				ResultField:  "networks",
				Object:       "ipv6networkcontainer",
				ObjectParams: eaMap,
				Params: map[string]uint{
					"cidr": prefixLen,
				},
				NetviewName: "",
			},
			NetviewName: netviewName,
			Comment:     comment,
			Ea:          ea,
		}
		networkinfo.SetObjectType("ipv6network")
		ref, err := connector.CreateObject(networkinfo)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^ipv6network.*"))
	})

	It("Should create a ipv6 networkcontainer within a networkcontainer with EA", func() {
		// Create an ipv6 networkcontainer within a networkcontainer with EA
		ea := ibclient.EA{"Region": "East"}
		comment := "Test ipv6 network Container creation with next_available_network"
		eaMap := map[string]string{"*Site": "Sri Lanka"}
		prefixLen := uint(67)
		netviewName := "default"
		networkinfo := &ibclient.NetworkContainerNextAvailable{
			//objectType: "ipv6network",
			Network: &ibclient.NetworkContainerNextAvailableInfo{
				Function:     "next_available_network",
				ResultField:  "networks",
				Object:       "ipv6networkcontainer",
				ObjectParams: eaMap,
				Params: map[string]uint{
					"cidr": prefixLen,
				},
				NetviewName: "",
			},
			NetviewName: netviewName,
			Comment:     comment,
			Ea:          ea,
		}

		networkinfo.SetObjectType("ipv6networkcontainer")
		ref, err := connector.CreateObject(networkinfo)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^ipv6networkcontainer.*"))
	})

	It("Should create a record:a within a network with EA", func() {
		// Create Record:A within a network with EA
		ea := ibclient.EA{"Site": "Basavangudi"}
		comment := "Test next_available_ip for record:a"
		eaMap := map[string]string{"*Site": "Namibia"}
		name := "testa.wapi.com"
		recordA := ibclient.IpNextAvailable{
			Name:                  name,
			NextAvailableIPv4Addr: ibclient.NewIpNextAvailableInfo(eaMap, nil, false, "IPV4"),
			Comment:               comment,
			Ea:                    ea,
			Disable:               false,
		}

		recordA.SetObjectType("record:a")
		ref, err := connector.CreateObject(&recordA)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^record:a.*"))
	})

	It("Should create a record:aaaa within a network with EA", func() {
		// Create Record:AAAA within a network with EA
		ea := ibclient.EA{"Site": "Bangalore"}
		comment := "Test next_available_ip for record:aaaa"
		eaMap := map[string]string{"*Site": "Norway"}
		name := "testaaaa.wapi.com"
		recordAAAA := ibclient.IpNextAvailable{
			Name:                  name,
			NextAvailableIPv6Addr: ibclient.NewIpNextAvailableInfo(eaMap, nil, false, "IPV6"),
			Comment:               comment,
			Ea:                    ea,
			Disable:               false,
		}

		recordAAAA.SetObjectType("record:aaaa")
		ref, err := connector.CreateObject(&recordAAAA)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^record:aaaa.*"))
	})

	It("Should create a record:host within a ipv6 network with EA", func() {
		// Create Record:Host within a network with EA
		ea := ibclient.EA{"Site": "Bangalore"}
		comment := "Test next_available_ip for record:host with ipv6"
		eaMap := map[string]string{"*Site": "Norway"}
		name := "testhost1.wapi.com"
		recordHost := ibclient.NewIpNextAvailable(name, "record:host", eaMap, nil, false, ea, comment, false, nil, "IPV6",
			false, false, "", "", "", "", false, 0, nil)

		recordHost.SetObjectType("record:host")
		ref, err := connector.CreateObject(recordHost)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^record:host.*"))
	})

	It("Should create a record:host within a ipv4 network with EA", func() {
		// Create Record:Host within a network with EA
		ea := ibclient.EA{"Site": "Mangalore"}
		comment := "Test next_available_ip for record:host with ipv4"
		eaMap := map[string]string{"*Site": "Namibia"}
		name := "testhost2.wapi.com"
		recordHost := ibclient.NewIpNextAvailable(name, "record:host", eaMap, nil, false, ea, comment, false, nil, "IPV4",
			false, false, "", "", "", "", false, 0, nil)

		recordHost.SetObjectType("record:host")
		ref, err := connector.CreateObject(recordHost)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^record:host.*"))
	})

	It("Should create a record:host within a ipv4 and ipv6 network with EA", func() {
		// Create Record:Host within a network with EA
		ea := ibclient.EA{"Site": "Mangalore"}
		comment := "Test next_available_ip for record:host with ipv4 and ipv6"
		eaMap := map[string]string{"*Site": "Japan"}
		name := "testhost3.wapi.com"
		recordHost := ibclient.NewIpNextAvailable(name, "record:host", eaMap, nil, false, ea, comment, false, nil, "Both",
			false, false, "", "", "", "", false, 0, nil)

		recordHost.SetObjectType("record:host")
		ref, err := connector.CreateObject(recordHost)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^record:host.*"))
	})

	It("Should fail to create a ipv4 network within a networkcontainer with EA", func() {
		// Create an ipv4 network within a networkcontainer with EA
		ea := ibclient.EA{"Region": "East"}
		comment := "Test ipv4 network creation with next_available_network"
		eaMap := map[string]string{"*Site": "Madagascar"}
		prefixLen := uint(26)
		netviewName := "default"
		networkinfo := &ibclient.NetworkContainerNextAvailable{
			Network: &ibclient.NetworkContainerNextAvailableInfo{
				Function:     "next_available_network",
				ResultField:  "networks",
				Object:       "networkcontainer",
				ObjectParams: eaMap,
				Params: map[string]uint{
					"cidr": prefixLen,
				},
				NetviewName: "",
			},
			NetviewName: netviewName,
			Comment:     comment,
			Ea:          ea,
		}
		networkinfo.SetObjectType("network")
		ref, err := connector.CreateObject(networkinfo)
		Expect(err).NotTo(BeNil())
		Expect(ref).To(BeEmpty())
	})

	It("Should fail to create a ipv6 networkcontainer within a networkcontainer with EA", func() {
		// Create an ipv6 networkcontainer within a networkcontainer with EA
		ea := ibclient.EA{"Region": "East"}
		comment := "Test ipv6 network Container creation with next_available_network"
		eaMap := map[string]string{"*Site": "Lakshwadeep"}
		prefixLen := uint(67)
		netviewName := "default"
		networkinfo := &ibclient.NetworkContainerNextAvailable{
			//objectType: "ipv6network",
			Network: &ibclient.NetworkContainerNextAvailableInfo{
				Function:     "next_available_network",
				ResultField:  "networks",
				Object:       "ipv6networkcontainer",
				ObjectParams: eaMap,
				Params: map[string]uint{
					"cidr": prefixLen,
				},
				NetviewName: "",
			},
			NetviewName: netviewName,
			Comment:     comment,
			Ea:          ea,
		}

		networkinfo.SetObjectType("ipv6networkcontainer")
		ref, err := connector.CreateObject(networkinfo)
		Expect(err).NotTo(BeNil())
		Expect(ref).To(BeEmpty())
	})

	It("Should fail to create a record:a within a network with EA", func() {
		// Create Record:A within a network with EA
		ea := ibclient.EA{"Site": "Basavangudi"}
		comment := "Test next_available_ip for record:a"
		eaMap := map[string]string{"*Site": "Mongolia"}
		name := "testa.wapi.com"
		recordA := ibclient.IpNextAvailable{
			Name:                  name,
			NextAvailableIPv4Addr: ibclient.NewIpNextAvailableInfo(eaMap, nil, false, "IPV4"),
			Comment:               comment,
			Ea:                    ea,
			Disable:               false,
		}

		recordA.SetObjectType("record:a")
		ref, err := connector.CreateObject(&recordA)
		Expect(err).NotTo(BeNil())
		Expect(ref).To(BeEmpty())
	})

	It("Should fail to create a record:aaaa within a network with EA", func() {
		// Create Record:AAAA within a network with EA
		ea := ibclient.EA{"Site": "Bangalore"}
		comment := "Test next_available_ip for record:aaaa"
		eaMap := map[string]string{"*Site": "Mongolia"}
		name := "testaaaa.wapi.com"
		recordAAAA := ibclient.IpNextAvailable{
			Name:                  name,
			NextAvailableIPv6Addr: ibclient.NewIpNextAvailableInfo(eaMap, nil, false, "IPV6"),
			Comment:               comment,
			Ea:                    ea,
			Disable:               false,
		}

		recordAAAA.SetObjectType("record:aaaa")
		ref, err := connector.CreateObject(&recordAAAA)
		Expect(err).NotTo(BeNil())
		Expect(ref).To(BeEmpty())
	})

	It("Should fail to create a record:host within a ipv4 and ipv6 network with EA", func() {
		// Create Record:Host within a network with EA
		ea := ibclient.EA{"Site": "Mangalore"}
		comment := "Test next_available_ip for record:host with ipv4 and ipv6"
		eaMap := map[string]string{"*Site": "Mongolia"}
		name := "testhost3.wapi.com"
		recordHost := ibclient.NewIpNextAvailable(name, "record:host", eaMap, nil, false, ea, comment, false, nil, "Both",
			false, false, "", "", "", "", false, 0, nil)

		recordHost.SetObjectType("record:host")
		ref, err := connector.CreateObject(recordHost)
		Expect(err).NotTo(BeNil())
		Expect(ref).To(BeEmpty())
	})

})

var _ = Describe("DTC Object pools", func() {
	var connector *ConnectorFacadeE2E

	var serverRef string
	var topologyRef string
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
		ibClientConnector, err := ibclient.NewConnector(hostConfig, authConfig, transportConfig, requestBuilder, requestor)
		Expect(err).To(BeNil())
		connector = &ConnectorFacadeE2E{*ibClientConnector, make([]string, 0)}

		var (
			serverName = "server.com"
			host       = "3.6.7.8"
			err1       error

			topologyName = "topology_test"
		)
		server := &ibclient.DtcServer{
			Name: &serverName,
			Host: &host,
		}
		serverRef, err1 = connector.CreateObject(server)
		Expect(err1).To(BeNil())

		topology := &ibclient.DtcTopology{
			Name: &topologyName,
			Rules: []*ibclient.DtcTopologyRule{
				{
					DestType:        "SERVER",
					DestinationLink: &serverRef,
				},
			},
		}
		topologyRef, err1 = connector.CreateObject(topology)
		Expect(err1).To(BeNil())
	})

	AfterEach(func() {
		err := connector.SweepObjects()
		Expect(err).To(BeNil())
	})

	It("Should create a dtc pool with ROUND_ROBIN method", func() {
		eaMap := ibclient.EA{"Site": "Burma"}
		dtcPool := ibclient.DtcPool{
			Name:    utils.StringPtr("dtc_pool_1.com"),
			Comment: utils.StringPtr("pool object creation"),
			Ea:      eaMap,
			Servers: []*ibclient.DtcServerLink{
				{
					Server: serverRef,
					Ratio:  4,
				},
			},
			LbPreferredMethod: "ROUND_ROBIN",
		}
		ref, err := connector.CreateObject(&dtcPool)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("dtc:pool/*"))
	})

	It("Should create a dtc pool with DYNAMIC_RATIO method", func() {
		eaMap := ibclient.EA{"Site": "Burma"}
		sf := map[string]string{"name": "http"}
		queryParams := ibclient.NewQueryParams(false, sf)
		var monitor []ibclient.DtcMonitorHttp
		_ = connector.GetObject(&ibclient.DtcMonitorHttp{}, "dtc:monitor:http", queryParams, &monitor)
		monitorRef := monitor[0].Ref
		dtcPool := ibclient.DtcPool{
			Name:    utils.StringPtr("dtc_pool_2.com"),
			Comment: utils.StringPtr("pool object creation"),
			Ea:      eaMap,
			Servers: []*ibclient.DtcServerLink{
				{
					Server: serverRef,
					Ratio:  4,
				},
			},
			LbPreferredMethod: "DYNAMIC_RATIO",
			Monitors: []*ibclient.DtcMonitorHttp{
				{
					Ref: monitorRef,
				},
			},
			LbDynamicRatioPreferred: &ibclient.SettingDynamicratio{
				Method:  "ROUND_TRIP_DELAY",
				Monitor: monitorRef,
			},
		}
		ref, err := connector.CreateObject(&dtcPool)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("dtc:pool/*"))
	})

	It("Should create a dtc pool with TOPOLOGY method", func() {
		eaMap := ibclient.EA{"Site": "Burma"}
		sf := map[string]string{"name": "http"}
		queryParams := ibclient.NewQueryParams(false, sf)
		var monitor []ibclient.DtcMonitorHttp
		_ = connector.GetObject(&ibclient.DtcMonitorHttp{}, "dtc:monitor:http", queryParams, &monitor)
		monitorRef := monitor[0].Ref
		dtcPool := ibclient.DtcPool{
			Name:    utils.StringPtr("dtc_pool_3.com"),
			Comment: utils.StringPtr("pool object creation"),
			Ea:      eaMap,
			Servers: []*ibclient.DtcServerLink{
				{
					Server: serverRef,
					Ratio:  4,
				},
			},
			Monitors: []*ibclient.DtcMonitorHttp{
				{
					Ref: monitorRef,
				},
			},
			LbPreferredMethod:   "TOPOLOGY",
			LbPreferredTopology: &topologyRef,
			LbAlternateMethod:   "DYNAMIC_RATIO",
			LbDynamicRatioAlternate: &ibclient.SettingDynamicratio{
				Method:  "ROUND_TRIP_DELAY",
				Monitor: monitorRef,
			},
		}
		ref, err := connector.CreateObject(&dtcPool)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("dtc:pool/*"))
	})
	It("should update the dtc pool", func() {
		var gridMembers []ibclient.Member
		err := connector.GetObject(&ibclient.Member{}, "", nil, &gridMembers)
		authZoneCreate := ibclient.ZoneAuth{
			Fqdn: "test.com",
			GridPrimary: []*ibclient.Memberserver{
				{
					Name: *gridMembers[0].HostName,
				},
				{
					Name: *gridMembers[1].HostName,
				},
			},
		}
		refAuthZone, err := connector.CreateObject(&authZoneCreate)
		authZoneCreate.Ref = refAuthZone
		Expect(err).To(BeNil())
		Expect(refAuthZone).To(MatchRegexp("zone_auth/*"))
		dtcPool := ibclient.DtcPool{
			Name: utils.StringPtr("dtc_pool_1.com"),
			Servers: []*ibclient.DtcServerLink{
				{
					Server: serverRef,
					Ratio:  4,
				},
			},
			LbPreferredMethod: "ROUND_ROBIN",
		}
		ref, err := connector.CreateObject(&dtcPool)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("dtc:pool/*"))
		dtcLBDN := ibclient.DtcLbdn{
			Name: utils.StringPtr("dtc_lbdn.com"),
			Pools: []*ibclient.DtcPoolLink{
				{
					Pool:  ref,
					Ratio: 4,
				},
			},
			LbMethod:  "ROUND_ROBIN",
			Patterns:  []string{"*test.com"},
			AuthZones: []*ibclient.ZoneAuth{&authZoneCreate},
		}
		lbdnRef, err := connector.CreateObject(&dtcLBDN)
		Expect(err).To(BeNil())
		Expect(lbdnRef).To(MatchRegexp("dtc:lbdn/*"))
		//update dtc pool
		autoConsolidateMonitors := false
		sf := map[string]string{"name": "http"}
		queryParams := ibclient.NewQueryParams(false, sf)
		var monitor []ibclient.DtcMonitorHttp
		_ = connector.GetObject(&ibclient.DtcMonitorHttp{}, "dtc:monitor:http", queryParams, &monitor)
		monitorRef := monitor[0].Ref
		dtcPoolUpdate := &ibclient.DtcPool{
			Name: utils.StringPtr("dtc_pool_2.com"),
			Servers: []*ibclient.DtcServerLink{
				{
					Server: serverRef,
					Ratio:  4,
				},
			},
			Monitors: []*ibclient.DtcMonitorHttp{
				{
					Ref: monitorRef,
				},
			},
			LbPreferredMethod:        "ROUND_ROBIN",
			Comment:                  utils.StringPtr("pool object update"),
			AutoConsolidatedMonitors: &autoConsolidateMonitors,
			ConsolidatedMonitors: []*ibclient.DtcPoolConsolidatedMonitorHealth{
				{
					Members: []string{
						*gridMembers[0].HostName,
						*gridMembers[1].HostName,
					},
					Monitor:                 monitorRef,
					Availability:            "ANY",
					FullHealthCommunication: true,
				},
			},
		}
		var res []ibclient.DtcPool
		search := &ibclient.DtcPool{}
		err = connector.GetObject(search, "", nil, &res)
		ref, err = connector.UpdateObject(dtcPoolUpdate, res[0].Ref)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("dtc:pool/*"))

		delRef, err := connector.DeleteObject(lbdnRef)
		Expect(err).To(BeNil())
		Expect(delRef).To(MatchRegexp("dtc:lbdn/*"))
	})
	It("should delete the dtc pool ", func() {
		dtcPool := ibclient.DtcPool{
			Name: utils.StringPtr("dtc_pool_1.com"),
			Servers: []*ibclient.DtcServerLink{
				{
					Server: serverRef,
					Ratio:  4,
				},
			},
			LbPreferredMethod: "ROUND_ROBIN",
		}
		ref, err := connector.CreateObject(&dtcPool)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("dtc:pool/*"))

		var res []ibclient.DtcPool
		search := &ibclient.DtcPool{}
		err = connector.GetObject(search, "", nil, &res)
		ref, err = connector.DeleteObject(res[0].Ref)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("dtc:pool/*"))
	})
	It("Should get the DTC pool", func() {
		dtcPool := ibclient.DtcPool{
			Name: utils.StringPtr("dtc_pool_1.com"),
			Servers: []*ibclient.DtcServerLink{
				{
					Server: serverRef,
					Ratio:  4,
				},
			},
			Comment:           utils.StringPtr("pool object creation"),
			LbPreferredMethod: "ROUND_ROBIN",
		}
		ref, err := connector.CreateObject(&dtcPool)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("dtc:pool/*"))

		var res []ibclient.DtcPool
		search := &ibclient.DtcPool{}
		search.SetReturnFields(append(search.ReturnFields(), "servers", "lb_preferred_method"))
		qp := ibclient.NewQueryParams(false, map[string]string{
			"name": "dtc_pool_1.com",
		})
		errCode := connector.GetObject(search, "", qp, &res)
		Expect(errCode).To(BeNil())
		Expect(res[0].Name).To(Equal(utils.StringPtr("dtc_pool_1.com")))
		Expect(res[0].Comment).To(Equal(utils.StringPtr("pool object creation")))
		Expect(res[0].Servers[0].Server).To(Equal(serverRef))
		Expect(res[0].LbPreferredMethod).To(Equal("ROUND_ROBIN"))
		Expect(res[0].Ref).To(MatchRegexp("dtc:pool/*"))
	})
	It("Should fail to create a DTC pool without mandatory parameters", func() {
		dtcPool := &ibclient.DtcPool{
			Comment: utils.StringPtr("wapi added"),
			Name:    utils.StringPtr("dtc_pool_1.com"),
		}
		_, err := connector.CreateObject(dtcPool)
		Expect(err).NotTo(BeNil())
	})
	It("Should fail to get a non-existent DTC pool", func() {
		var res []ibclient.DtcPool
		name := "dtc_pool_wapi"
		search := &ibclient.DtcPool{Name: &name}
		err := connector.GetObject(search, "", nil, &res)
		Expect(err).NotTo(BeNil())
	})
	It("Should fail to update a non-existent DTC pool", func() {
		name := "dtc_pool_wapi"
		dtcPool := &ibclient.DtcPool{
			Name: &name,
			Servers: []*ibclient.DtcServerLink{
				{
					Server: serverRef,
					Ratio:  4,
				},
			},
			Comment: utils.StringPtr("wapi added"),
		}

		_, err := connector.UpdateObject(dtcPool, "nonexistent_ref")
		Expect(err).NotTo(BeNil())
	})
	It("Should fail to delete a non-existent Dtc pool", func() {
		_, err := connector.DeleteObject("nonexistent_ref")
		Expect(err).NotTo(BeNil())
	})
})

var _ = Describe("DTC LBDN and server object", func() {
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

	It("Should create a DTC LBDN object with minimum parameters", func() {
		// Create a DTC LBDN object
		lbdn := ibclient.DtcLbdn{
			Name:     utils.StringPtr("test-LBDN1"),
			LbMethod: "ROUND_ROBIN",
		}
		ref, err := connector.CreateObject(&lbdn)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^dtc:lbdn.*"))
	})

	It("Should create a DTC LBDN object with maximum parameters", func() {

		var (
			topologyRef, poolRef string
			err                  error
		)
		// create zoneAuth
		var gridMembers []ibclient.Member
		err = connector.GetObject(&ibclient.Member{}, "", nil, &gridMembers)
		Expect(err).To(BeNil())
		zones := ibclient.ZoneAuth{
			Fqdn:        "wapi.com",
			GridPrimary: []*ibclient.Memberserver{{Name: *gridMembers[0].HostName}},
		}
		zoneRef, err := connector.CreateObject(&zones)
		Expect(err).To(BeNil())
		zones.Ref = zoneRef

		// create server
		server := ibclient.DtcServer{
			Name: utils.StringPtr("TestServer123"),
			Host: utils.StringPtr("12.12.1.1"),
		}
		serverRef, err := connector.CreateObject(&server)
		Expect(err).To(BeNil())
		Expect(serverRef).To(MatchRegexp("^dtc:server.*"))

		// create pools
		pool := ibclient.DtcPool{
			Name:              utils.StringPtr("test-pool1"),
			LbPreferredMethod: "ROUND_ROBIN",
			Servers: []*ibclient.DtcServerLink{{
				Server: serverRef,
				Ratio:  uint32(2),
			}},
		}

		poolRef, err = connector.CreateObject(&pool)
		Expect(err).To(BeNil())
		Expect(poolRef).To(MatchRegexp("^dtc:pool.*"))

		// create topology
		topology := ibclient.DtcTopology{
			Name: utils.StringPtr("test-topology"),
			Rules: []*ibclient.DtcTopologyRule{
				{
					DestType:        "POOL",
					DestinationLink: utils.StringPtr(poolRef),
				},
			},
		}
		topologyRef, err = connector.CreateObject(&topology)
		Expect(err).To(BeNil())

		// Create a DTC LBDN object with maximum parameters
		lbdn := ibclient.DtcLbdn{
			Name:      utils.StringPtr("test-LBDN12"),
			LbMethod:  "TOPOLOGY",
			Topology:  utils.StringPtr(topologyRef),
			AuthZones: []*ibclient.ZoneAuth{&zones},
			Pools:     []*ibclient.DtcPoolLink{{Pool: poolRef, Ratio: uint32(3)}},
			Types:     []string{"A", "CNAME", "AAAA"},
			Patterns:  []string{"*wapi.com"},
		}

		ref, err := connector.CreateObject(&lbdn)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^dtc:lbdn.*"))
	})

	It("Should get the DTC LBDN object TestLBDN222", func() {
		lbdn := ibclient.DtcLbdn{
			Name:     utils.StringPtr("TestLBDN222"),
			LbMethod: "ROUND_ROBIN",
		}
		ref, err := connector.CreateObject(&lbdn)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^dtc:lbdn.*"))

		var res []ibclient.DtcLbdn
		search := &ibclient.DtcLbdn{}
		errCode := connector.GetObject(search, "", nil, &res)
		Expect(errCode).To(BeNil())
		Expect(res[0].Ref).To(MatchRegexp("^dtc:lbdn.*"))
	})

	It("Should update the DTC LBDN object TestLBDN11", func() {
		lbdn := ibclient.DtcLbdn{
			Name:     utils.StringPtr("TestLBDN11"),
			LbMethod: "ROUND_ROBIN",
			Comment:  utils.StringPtr("sample comment"),
			Priority: utils.Uint32Ptr(3),
			Types:    []string{"A"},
		}
		ref, errCode := connector.CreateObject(&lbdn)
		Expect(errCode).To(BeNil())
		Expect(ref).To(MatchRegexp("^dtc:lbdn.*"))

		// Update the DTC Lbdn object
		lbdnUpdated := ibclient.DtcLbdn{
			Name:     utils.StringPtr("TestLBDN1111"),
			Comment:  utils.StringPtr("sample comment updated"),
			LbMethod: "RATIO",
		}

		var res []ibclient.DtcLbdn
		search := &ibclient.DtcLbdn{}
		err := connector.GetObject(search, "", nil, &res)
		ref, err = connector.UpdateObject(&lbdnUpdated, ref)
		Expect(err).To(BeNil())
		Expect(ref).To(HaveSuffix("TestLBDN1111"))
	})

	It("Should delete the DTC LBDN object TestLBDN22", func() {
		lbdn := ibclient.DtcLbdn{
			Name:     utils.StringPtr("TestLBDN22"),
			LbMethod: "ROUND_ROBIN",
			Comment:  utils.StringPtr("sample comment"),
			Priority: utils.Uint32Ptr(3),
			Types:    []string{"A", "CNAME"},
		}
		ref, errCode := connector.CreateObject(&lbdn)
		Expect(errCode).To(BeNil())
		Expect(ref).To(MatchRegexp("^dtc:lbdn.*"))
		ref, err := connector.DeleteObject(ref)
		Expect(err).To(BeNil())
	})

	It("Should fail to create a DTC LBDN object TestLBDN33", func() {
		lbdn := ibclient.DtcLbdn{
			Name:     utils.StringPtr("TestLBDN33"),
			Comment:  utils.StringPtr("sample comment"),
			Priority: utils.Uint32Ptr(3),
			Types:    []string{"A", "CNAME"},
		}
		_, err := connector.CreateObject(&lbdn)
		Expect(err).NotTo(BeNil())
	})

	It("Should fail to get a non-existent DTC LBDN object Testlbdn1010", func() {
		var res []ibclient.DtcLbdn
		sf := map[string]string{"name": "Testlbdn1010"}
		qp := ibclient.NewQueryParams(false, sf)
		err := connector.GetObject(&ibclient.DtcLbdn{}, "", qp, &res)
		Expect(res).To(BeEmpty())
		Expect(err).NotTo(BeNil())
	})

	It("Should fail to update a non-existent DTC LBDN object TestLBDN44", func() {
		lbdn := ibclient.DtcLbdn{
			Name:     utils.StringPtr("TestLBDN44"),
			Comment:  utils.StringPtr("sample comment"),
			Priority: utils.Uint32Ptr(3),
			Types:    []string{"A", "CNAME"},
		}

		_, err := connector.UpdateObject(&lbdn, "nonexistent_ref")
		Expect(err).NotTo(BeNil())
	})

	It("Should fail to delete a non-existent DTC LBDN object", func() {
		_, err := connector.DeleteObject("nonexistent_ref")
		Expect(err).NotTo(BeNil())
	})

	// create DTC server object
	It("Should create DTC server object with minimum params", func() {
		server := ibclient.DtcServer{
			Name: utils.StringPtr("test-server1"),
			Host: utils.StringPtr("12.12.1.1"),
		}
		ref, err := connector.CreateObject(&server)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^dtc:server.*"))
	})

	// create DTC server object with maximum params
	It("Should create DTC server object with maximum params", func() {

		monitor := ibclient.DtcMonitorHttp{
			Name: utils.StringPtr("test-monitor"),
		}
		monitorRef, err := connector.CreateObject(&monitor)
		Expect(err).To(BeNil())
		Expect(monitorRef).To(MatchRegexp("^dtc:monitor:http.*"))
		server := ibclient.DtcServer{
			Name:                 utils.StringPtr("test-server111"),
			Host:                 utils.StringPtr("12.12.1.11"),
			Comment:              utils.StringPtr("test comment"),
			UseSniHostname:       utils.BoolPtr(true),
			SniHostname:          utils.StringPtr("test-sni"),
			AutoCreateHostRecord: utils.BoolPtr(true),
			Disable:              utils.BoolPtr(true),
			Ea:                   ibclient.EA{"Site": "India"},
			Monitors: []*ibclient.DtcServerMonitor{{
				Monitor: monitorRef,
				Host:    "1.2.3.4",
			}},
		}
		ref, err := connector.CreateObject(&server)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^dtc:server.*"))
	})

	// get Dtc server
	It("Should get Dtc server", func() {
		server := ibclient.DtcServer{
			Name: utils.StringPtr("test-server2"),
			Host: utils.StringPtr("12.12.1.1"),
		}
		ref, err := connector.CreateObject(&server)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^dtc:server.*"))
		var res ibclient.DtcServer
		err = connector.GetObject(&ibclient.DtcServer{}, ref, nil, &res)
		Expect(err).To(BeNil())
		Expect(*res.Name).To(MatchRegexp("test-server2"))
	})

	// update Dtc server
	It("Should get Dtc server", func() {
		server := ibclient.DtcServer{
			Name:    utils.StringPtr("test-server3"),
			Host:    utils.StringPtr("12.12.1.2"),
			Comment: utils.StringPtr("test comment"),
		}
		ref, err := connector.CreateObject(&server)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^dtc:server.*"))
		updatedServer := ibclient.DtcServer{
			Name:    utils.StringPtr("test-server22"),
			Host:    utils.StringPtr("11.11.1.1"),
			Comment: utils.StringPtr("test comment updated"),
		}
		ref, err = connector.UpdateObject(&updatedServer, ref)
		Expect(err).To(BeNil())
		Expect(ref).To(HaveSuffix("test-server22"))
	})

	// delete Dtc server
	It("Should delete a Dtc server", func() {
		server := ibclient.DtcServer{
			Name: utils.StringPtr("test-server3"),
			Host: utils.StringPtr("12.12.1.3"),
		}
		ref, err := connector.CreateObject(&server)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^dtc:server.*"))
		delRef, err := connector.DeleteObject(ref)
		Expect(err).To(BeNil())
		Expect(delRef).To(MatchRegexp("^dtc:server.*"))
	})

	// create DTC server object, -ve scenario
	It("Should fail to create DTC server object with minimum params", func() {
		server := ibclient.DtcServer{
			Host: utils.StringPtr("12.12.1.1"),
		}
		_, err := connector.CreateObject(&server)
		Expect(err).NotTo(BeNil())
	})

	// get Dtc server, -ve scenario
	It("Should fail to get a non existent Dtc server", func() {
		var res []ibclient.DtcServer
		err := connector.GetObject(&ibclient.DtcServer{}, "nonexistent_ref", nil, &res)
		Expect(res).To(BeNil())
		Expect(err).NotTo(BeNil())
	})

	// update Dtc server, -ve scenario
	It("Should fail to update a Dtc server", func() {
		server := ibclient.DtcServer{
			Name: utils.StringPtr("test-server4"),
			Host: utils.StringPtr("12.12.1.4"),
		}
		ref, err := connector.CreateObject(&server)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^dtc:server.*"))
		_, err = connector.UpdateObject(&server, "nonexistent_ref")
		Expect(err).NotTo(BeNil())
	})

	// delete Dtc server, -ve scenario
	It("Should fail to delete a non existent Dtc server", func() {
		_, err := connector.DeleteObject("nonexistent_ref")
		Expect(err).NotTo(BeNil())
	})

})

var _ = Describe("Record Alias", func() {
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

		zones := ibclient.ZoneAuth{
			Fqdn: "wapi.com",
		}
		zoneRef, err := connector.CreateObject(&zones)
		Expect(err).To(BeNil())
		zones.Ref = zoneRef

	})

	AfterEach(func() {
		err := connector.SweepObjects()
		Expect(err).To(BeNil())
	})

	It("Should create Alias record with minimum parameters", func() {
		// Create an Alias Record with minimum parameters
		recordAlias := ibclient.RecordAlias{
			Name:       utils.StringPtr("test-alias1.wapi.com"),
			TargetType: "A",
			TargetName: utils.StringPtr("aa.test.com"),
		}
		ref, err := connector.CreateObject(&recordAlias)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^record:alias.*"))

		var aliasRecord *ibclient.RecordAlias
		err = connector.GetObject(&ibclient.RecordAlias{}, ref, nil, &aliasRecord)
		Expect(err).To(BeNil())
		Expect(aliasRecord).NotTo(BeNil())
	})

	It("Should create an Alias Record with maximum parameters", func() {

		// Create an Alias Record object with maximum parameters
		recordAlias := ibclient.RecordAlias{
			Name:       utils.StringPtr("test-alias2.wapi.com"),
			TargetType: "PTR",
			TargetName: utils.StringPtr("aa.bb.com"),
			Comment:    utils.StringPtr("sample comment"),
			Ea:         ibclient.EA{"Site": "India"},
			Disable:    utils.BoolPtr(true),
			Ttl:        utils.Uint32Ptr(100),
			UseTtl:     utils.BoolPtr(true),
		}

		ref, err := connector.CreateObject(&recordAlias)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^record:alias.*"))

		var aliasRecord *ibclient.RecordAlias
		err = connector.GetObject(&ibclient.RecordAlias{}, ref, nil, &aliasRecord)
		Expect(err).To(BeNil())
		Expect(aliasRecord).NotTo(BeNil())
	})

	It("Should get alias record testalias1.wapi.com", func() {
		recordAlias := ibclient.RecordAlias{
			Name:       utils.StringPtr("testalias1.wapi.com"),
			TargetType: "NAPTR",
			TargetName: utils.StringPtr("aab.bb.com"),
		}
		ref, err := connector.CreateObject(&recordAlias)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^record:alias.*"))

		var res ibclient.RecordAlias
		search := &ibclient.RecordAlias{}
		errCode := connector.GetObject(search, ref, nil, &res)
		Expect(errCode).To(BeNil())
		Expect(res.Ref).To(MatchRegexp("^record:alias.*"))

		Expect(*res.Name).To(Equal("testalias1.wapi.com"))
		Expect(res.TargetType).To(Equal("NAPTR"))
		Expect(*res.TargetName).To(Equal("aab.bb.com"))
	})

	It("Should update record Alias testalias2.wapi.com", func() {
		recordAlias := ibclient.RecordAlias{
			Name:       utils.StringPtr("testalias2.wapi.com"),
			TargetType: "TXT",
			TargetName: utils.StringPtr("aaa.bbb.com"),
			Comment:    utils.StringPtr("test comment"),
			Ea:         ibclient.EA{"Site": "Sri Lanka"},
			Disable:    utils.BoolPtr(true),
			Ttl:        utils.Uint32Ptr(500),
			UseTtl:     utils.BoolPtr(true),
		}
		ref, errCode := connector.CreateObject(&recordAlias)
		Expect(errCode).To(BeNil())
		Expect(ref).To(MatchRegexp("^record:alias.*"))

		// Update the Record Alias
		recordAliasUpdated := ibclient.RecordAlias{
			Name:       utils.StringPtr("testalias222.wapi.com"),
			TargetType: "PTR",
			TargetName: utils.StringPtr("xyz.bnm.com"),
			Comment:    utils.StringPtr("test comment updated"),
		}

		var res []ibclient.RecordAlias
		search := &ibclient.RecordAlias{}
		err := connector.GetObject(search, "", nil, &res)
		ref, err = connector.UpdateObject(&recordAliasUpdated, ref)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^record:alias.*"))

		var aliasRecord *ibclient.RecordAlias
		err = connector.GetObject(&ibclient.RecordAlias{}, ref, nil, &aliasRecord)
		Expect(err).To(BeNil())
		Expect(aliasRecord).NotTo(BeNil())
	})

	It("Should delete record alias ", func() {
		recordAlias := ibclient.RecordAlias{
			Name:       utils.StringPtr("Alias1111.wapi.com"),
			TargetType: "NAPTR",
			TargetName: utils.StringPtr("abs.test.com"),
		}
		ref, errCode := connector.CreateObject(&recordAlias)
		Expect(errCode).To(BeNil())
		Expect(ref).To(MatchRegexp("^record:alias.*"))
		ref, err := connector.DeleteObject(ref)
		Expect(err).To(BeNil())
	})

	It("Should fail to create alias record alias222.wapi.com", func() {
		recordAlias := ibclient.RecordAlias{
			Name:       utils.StringPtr("alias222.wapi.com"),
			TargetType: "NAPTR",
		}
		_, err := connector.CreateObject(&recordAlias)
		Expect(err).NotTo(BeNil())
	})

	It("Should fail to get a non-existent alias record test_alias1", func() {
		var res []ibclient.RecordAlias
		sf := map[string]string{"name": "test_alias1"}
		qp := ibclient.NewQueryParams(false, sf)
		err := connector.GetObject(&ibclient.RecordAlias{}, "", qp, &res)
		Expect(res).To(BeEmpty())
		Expect(err).NotTo(BeNil())
	})

	It("Should fail to update a non-existent alias record alias222", func() {
		recordAlias := ibclient.RecordAlias{
			Name:       utils.StringPtr("alias222.wapi.com"),
			TargetType: "NAPTR",
			TargetName: utils.StringPtr("abs.test.com"),
		}
		_, err := connector.UpdateObject(&recordAlias, "nonexistent_ref")
		Expect(err).NotTo(BeNil())
	})

	It("Should fail to delete a non-existent alias record", func() {
		_, err := connector.DeleteObject("nonexistent_ref")
		Expect(err).NotTo(BeNil())
	})

})

var _ = Describe("NS Record", func() {
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
		ibClientConnector, err := ibclient.NewConnector(hostConfig, authConfig, transportConfig, requestBuilder, requestor)
		Expect(err).To(BeNil())
		connector = &ConnectorFacadeE2E{*ibClientConnector, make([]string, 0)}

		zones := ibclient.ZoneAuth{
			Fqdn: "wapi_test.com",
		}
		zoneRef, err := connector.CreateObject(&zones)
		Expect(err).To(BeNil())
		zones.Ref = zoneRef
	})
	AfterEach(func() {
		err := connector.SweepObjects()
		Expect(err).To(BeNil())
	})

	It("Should create a NS Record", func() {
		nsRecord := ibclient.RecordNS{
			Name:       "wapi_test.com",
			Nameserver: utils.StringPtr("ns1.wapi_test.com"),
			Addresses: []*ibclient.ZoneNameServer{
				{
					Address:       "2.3.4.5",
					AutoCreatePtr: true,
				},
			},
		}
		ref, err := connector.CreateObject(&nsRecord)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("record:ns/*"))
	})

	It("Should update the NS Record", func() {
		nsRecord := ibclient.RecordNS{
			Name:       "wapi_test.com",
			Nameserver: utils.StringPtr("ns1.wapi_test.com"),
			Addresses: []*ibclient.ZoneNameServer{
				{
					Address:       "2.3.4.5",
					AutoCreatePtr: true,
				},
			},
		}
		ref, err := connector.CreateObject(&nsRecord)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("record:ns/*"))

		// Update the DTC Lbdn object
		nsRecordUpdate := ibclient.RecordNS{
			Nameserver: utils.StringPtr("ns2.wapi_test.com"),
			Addresses: []*ibclient.ZoneNameServer{
				{
					Address:       "2.3.4.5",
					AutoCreatePtr: true,
				},
			},
		}

		var res []ibclient.RecordNS
		search := &ibclient.RecordNS{}
		err = connector.GetObject(search, "", nil, &res)
		ref, err = connector.UpdateObject(&nsRecordUpdate, ref)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("record:ns/*"))
	})

	It("Should get the NS Record", func() {
		nsRecord := ibclient.RecordNS{
			Name:       "wapi_test.com",
			Nameserver: utils.StringPtr("ns1.wapi_test.com"),
			Addresses: []*ibclient.ZoneNameServer{
				{
					Address:       "2.3.4.5",
					AutoCreatePtr: true,
				},
			},
		}
		ref, err := connector.CreateObject(&nsRecord)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("record:ns/*"))

		var res []ibclient.RecordNS
		search := &ibclient.RecordNS{}
		search.SetReturnFields(append(search.ReturnFields(), "addresses"))
		qp := ibclient.NewQueryParams(false, map[string]string{
			"nameserver": "ns1.wapi_test.com",
		})
		errCode := connector.GetObject(search, "", qp, &res)
		Expect(errCode).To(BeNil())
		Expect(res[0].Nameserver).To(Equal(utils.StringPtr("ns1.wapi_test.com")))
		Expect(res[0].Name).To(Equal("wapi_test.com"))
		Expect(res[0].Addresses[0].Address).To(Equal("2.3.4.5"))
		Expect(res[0].Addresses[0].AutoCreatePtr).To(Equal(true))
		Expect(res[0].Ref).To(MatchRegexp("record:ns/*"))
	})
	It("Should delete a NS Record", func() {
		nsRecord := ibclient.RecordNS{
			Name:       "wapi_test.com",
			Nameserver: utils.StringPtr("ns1.wapi_test.com"),
			Addresses: []*ibclient.ZoneNameServer{
				{
					Address:       "2.3.4.5",
					AutoCreatePtr: true,
				},
			},
		}
		ref, err := connector.CreateObject(&nsRecord)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("record:ns/*"))
		delRef, err := connector.DeleteObject(ref)
		Expect(err).To(BeNil())
		Expect(delRef).To(MatchRegexp("record:ns/*"))
	})

	It("Should fail to create a NS Record object", func() {
		nsRecord := ibclient.RecordNS{
			Name: "wapi_test.com",
			Addresses: []*ibclient.ZoneNameServer{
				{
					Address:       "2.3.4.5",
					AutoCreatePtr: true,
				},
			},
		}
		_, err := connector.CreateObject(&nsRecord)
		Expect(err).NotTo(BeNil())
	})

	// update Dtc server, -ve scenario
	It("Should fail to update a NS Record with wrong reference", func() {
		nsRecord := ibclient.RecordNS{
			Name:       "wapi_test.com",
			Nameserver: utils.StringPtr("ns3.wapi_test.com"),
			Addresses: []*ibclient.ZoneNameServer{
				{
					Address:       "2.3.4.5",
					AutoCreatePtr: true,
				},
			},
		}
		ref, err := connector.CreateObject(&nsRecord)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("record:ns/*"))
		_, err = connector.UpdateObject(&nsRecord, "nonexistent_ref")
		Expect(err).NotTo(BeNil())
	})
	It("Should fail to update a NS record", func() {
		nsRecord := ibclient.RecordNS{
			Name:       "wapi_test.com",
			Nameserver: utils.StringPtr("ns3.wapi_test.com"),
			Addresses: []*ibclient.ZoneNameServer{
				{
					Address:       "2.3.4.5",
					AutoCreatePtr: true,
				},
			},
		}
		ref, err := connector.CreateObject(&nsRecord)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("record:ns/*"))
		nsRecordUpdate := ibclient.RecordNS{
			Name: "wapi_test2.com",
		}
		var res []ibclient.RecordNS
		search := &ibclient.RecordNS{}
		err = connector.GetObject(search, "", nil, &res)
		ref, err = connector.UpdateObject(&nsRecordUpdate, res[0].Ref)
		Expect(err).NotTo(BeNil())
	})
})

var _ = Describe("Record Range Template", func() {
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

	It("Should create Range Template record with minimum parameters", func() {
		// Create a Range Template Record with minimum parameters
		rangeTemplate := ibclient.Rangetemplate{
			Name:              utils.StringPtr("template1"),
			NumberOfAddresses: utils.Uint32Ptr(10),
			Offset:            utils.Uint32Ptr(20),
		}
		ref, err := connector.CreateObject(&rangeTemplate)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^rangetemplate.*"))

		var templateRange *ibclient.Rangetemplate
		err = connector.GetObject(&ibclient.Rangetemplate{}, ref, nil, &templateRange)
		Expect(err).To(BeNil())
		Expect(templateRange).NotTo(BeNil())
	})

	It("Should create a Range Template Record with maximum parameters", func() {

		// Create a Range Template Record object with maximum parameters
		options := []*ibclient.Dhcpoption{
			{
				Name:  "domain-name-servers",
				Value: "11.22.3.1",
			},
			{
				Name:  "domain-name",
				Value: "aa.mm.ee",
			},
		}
		rangeTemplate := ibclient.Rangetemplate{
			Name:              utils.StringPtr("template2"),
			NumberOfAddresses: utils.Uint32Ptr(10),
			Offset:            utils.Uint32Ptr(20),
			Comment:           utils.StringPtr("test comment"),
			Ea:                ibclient.EA{"Site": "Sapporo"},
			Options:           options,
		}
		ref, err := connector.CreateObject(&rangeTemplate)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^rangetemplate.*"))

		var templateRange *ibclient.Rangetemplate
		err = connector.GetObject(&ibclient.Rangetemplate{}, ref, nil, &templateRange)
		Expect(err).To(BeNil())
		Expect(templateRange).NotTo(BeNil())
	})

	It("Should get Range Template record template111", func() {
		rangeTemplate := ibclient.Rangetemplate{
			Name:              utils.StringPtr("template111"),
			NumberOfAddresses: utils.Uint32Ptr(30),
			Offset:            utils.Uint32Ptr(40),
		}
		ref, err := connector.CreateObject(&rangeTemplate)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^rangetemplate.*"))

		var res ibclient.Rangetemplate
		search := &ibclient.Rangetemplate{}
		errCode := connector.GetObject(search, ref, nil, &res)
		Expect(errCode).To(BeNil())
		Expect(res.Ref).To(MatchRegexp("rangetemplate.*"))

		Expect(*res.Name).To(Equal("template111"))
		Expect(int(*res.NumberOfAddresses)).To(Equal(30))
		Expect(int(*res.Offset)).To(Equal(40))
	})

	It("Should update record Range Template template1234", func() {
		options := []*ibclient.Dhcpoption{
			{
				Name:  "domain-name-servers",
				Value: "11.22.3.1",
			},
			{
				Name:  "domain-name",
				Value: "aa.mm.ee",
			},
		}
		rangeTemplate := ibclient.Rangetemplate{
			Name:              utils.StringPtr("template1234"),
			NumberOfAddresses: utils.Uint32Ptr(60),
			Offset:            utils.Uint32Ptr(50),
			Comment:           utils.StringPtr("test comment"),
			Ea:                ibclient.EA{"Site": "Sapporo"},
			Options:           options,
		}
		ref, errCode := connector.CreateObject(&rangeTemplate)
		Expect(errCode).To(BeNil())
		Expect(ref).To(MatchRegexp("^rangetemplate.*"))

		// Update the Record range template
		rangeTemplateUpdated := ibclient.Rangetemplate{
			Name:              utils.StringPtr("template4567"),
			NumberOfAddresses: utils.Uint32Ptr(70),
			Offset:            utils.Uint32Ptr(80),
			Comment:           utils.StringPtr("test comment updated"),
			Ea:                ibclient.EA{"Site": "Sendai"},
			Options:           options,
		}

		var res []ibclient.Rangetemplate
		search := &ibclient.Rangetemplate{}
		err := connector.GetObject(search, "", nil, &res)
		ref, err = connector.UpdateObject(&rangeTemplateUpdated, ref)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^rangetemplate.*"))

		var templateRange *ibclient.Rangetemplate
		err = connector.GetObject(&ibclient.Rangetemplate{}, ref, nil, &templateRange)
		Expect(err).To(BeNil())
		Expect(templateRange).NotTo(BeNil())
	})

	It("Should delete record Range Template ", func() {
		rangeTemplate := ibclient.Rangetemplate{
			Name:              utils.StringPtr("template222"),
			NumberOfAddresses: utils.Uint32Ptr(33),
			Offset:            utils.Uint32Ptr(44),
		}
		ref, errCode := connector.CreateObject(&rangeTemplate)
		Expect(errCode).To(BeNil())
		Expect(ref).To(MatchRegexp("^rangetemplate.*"))
		ref, err := connector.DeleteObject(ref)
		Expect(err).To(BeNil())
	})

	It("Should fail to create Range Template record template333", func() {
		rangeTemplate := ibclient.Rangetemplate{
			Name:              utils.StringPtr("template333"),
			NumberOfAddresses: utils.Uint32Ptr(33),
		}
		_, err := connector.CreateObject(&rangeTemplate)
		Expect(err).NotTo(BeNil())
	})

	It("Should fail to get a non-existent Range Template record range-template100", func() {
		var res []ibclient.Rangetemplate
		sf := map[string]string{"name": "range-template100"}
		qp := ibclient.NewQueryParams(false, sf)
		err := connector.GetObject(&ibclient.Rangetemplate{}, "", qp, &res)
		Expect(res).To(BeEmpty())
		Expect(err).NotTo(BeNil())
	})

	It("Should fail to update a non-existent Range Template record template444", func() {
		rangeTemplate := ibclient.Rangetemplate{
			Name:              utils.StringPtr("template444"),
			NumberOfAddresses: utils.Uint32Ptr(122),
			Offset:            utils.Uint32Ptr(32),
		}
		_, err := connector.UpdateObject(&rangeTemplate, "nonexistent_ref")
		Expect(err).NotTo(BeNil())
	})

	It("Should fail to delete a non-existent Range Template record", func() {
		_, err := connector.DeleteObject("nonexistent_ref")
		Expect(err).NotTo(BeNil())
	})

})

var _ = Describe("IPV4 fixed address", func() {
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
		ibClientConnector, err := ibclient.NewConnector(hostConfig, authConfig, transportConfig, requestBuilder, requestor)
		Expect(err).To(BeNil())
		connector = &ConnectorFacadeE2E{*ibClientConnector, make([]string, 0)}

		var (
			cidr        = "12.0.0.0/24"
			networkView = "default"
			err1        error
		)
		ipv4Network := ibclient.NewNetwork(networkView, cidr, false, "", nil)
		_, err1 = connector.CreateObject(ipv4Network)
		Expect(err1).To(BeNil())
	})
	AfterEach(func() {
		err := connector.SweepObjects()
		Expect(err).To(BeNil())
	})
	It("should create an IPV4 fixed address object with minimal parameters", func() {
		client := "MAC_ADDRESS"
		fixedAddress := ibclient.NewFixedAddress("", "", "12.0.0.1", "", "43:56:98:98:32:21", &client, nil, "", false, "", nil, nil, nil, nil, false, nil, false)
		ref, err := connector.CreateObject(fixedAddress)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("fixedaddress/*"))
	})
	It("should create an IPV4 fixed address with maximal parameters", func() {
		ea := ibclient.EA{"Site": "India"}
		options := []*ibclient.Dhcpoption{
			{
				Name:        "routers",
				Num:         3,
				VendorClass: "DHCP",
				Value:       "12.0.0.34",
				UseOption:   true,
			},
		}
		client := "CLIENT_ID"
		clientIdentifierPrependZero := false
		dhcpClientIdentifier := "34"
		fixedAddress := ibclient.NewFixedAddress("default", "fixedaddress1", "12.0.0.2", "12.0.0.0/24", "43:56:98:98:32:21", &client, ea, "", false, "test comment", nil, nil, &clientIdentifierPrependZero, &dhcpClientIdentifier, true, options, true)
		ref, err := connector.CreateObject(fixedAddress)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("fixedaddress/*"))
	})
	// get IPV4 fixed address
	It("Should get IPV4 fixed address", func() {
		ea := ibclient.EA{"Site": "India"}
		options := []*ibclient.Dhcpoption{
			{
				Name:        "routers",
				Num:         3,
				VendorClass: "DHCP",
				Value:       "12.0.0.34",
				UseOption:   true,
			},
		}
		client := "CLIENT_ID"
		clientIdentifierPrependZero := false
		dhcpClientIdentifier := "34"
		Name := "fixedaddress1"
		fixedAddress := ibclient.NewFixedAddress("default", Name, "12.0.0.2", "12.0.0.0/24", "43:56:98:98:32:21", &client, ea, "", false, "test comment", nil, nil, &clientIdentifierPrependZero, &dhcpClientIdentifier, true, options, true)
		ref, err := connector.CreateObject(fixedAddress)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("fixedaddress/*"))
		var res ibclient.FixedAddress
		search := &ibclient.FixedAddress{}
		search.SetReturnFields(append(search.ReturnFields(), "name", "network_view", "ipv4addr", "network", "options"))
		errCode := connector.GetObject(search, ref, nil, &res)
		Expect(errCode).To(BeNil())
		Expect(res.Name).To(Equal(&Name))
		Expect(res.NetviewName).To(Equal("default"))
		Expect(res.IPv4Address).To(Equal("12.0.0.2"))
		Expect(res.Cidr).To(Equal("12.0.0.0/24"))
		Expect(res.Options[0].Name).To(Equal("dhcp-lease-time"))
		Expect(res.Options[0].Num).To(Equal(uint32(51)))
		Expect(res.Options[0].Value).To(Equal("43200"))
		Expect(res.Options[0].UseOption).To(Equal(false))
		Expect(res.Options[0].VendorClass).To(Equal("DHCP"))
		Expect(res.Options[1].Name).To(Equal("routers"))
		Expect(res.Options[1].Num).To(Equal(uint32(3)))
		Expect(res.Options[1].Value).To(Equal("12.0.0.34"))
		Expect(res.Options[1].UseOption).To(Equal(true))
		Expect(res.Options[1].VendorClass).To(Equal("DHCP"))
	})
	//It should update IPV4 fixed address
	It("Should update IPV4 fixed address", func() {
		client := "MAC_ADDRESS"
		fixedAddress := ibclient.NewFixedAddress("", "", "12.0.0.3", "", "43:56:98:98:32:21", &client, nil, "", false, "", nil, nil, nil, nil, false, nil, false)
		ref, err := connector.CreateObject(fixedAddress)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("fixedaddress/*"))

		updateOptions := []*ibclient.Dhcpoption{
			{
				Name:        "routers",
				Num:         3,
				VendorClass: "DHCP",
				Value:       "12.0.0.34",
				UseOption:   true,
			},
		}
		updateClient := "REMOTE_ID"
		updateAgentRemoteId := "23"
		updatedFixedAddress := ibclient.NewFixedAddress("", "fixedaddress", "12.0.0.4", "", "", &updateClient, nil, "", false, "comment 1", nil, &updateAgentRemoteId, nil, nil, false, updateOptions, true)
		var res []ibclient.FixedAddress
		search := &ibclient.FixedAddress{}
		err = connector.GetObject(search, "", nil, &res)
		ref, err = connector.UpdateObject(updatedFixedAddress, ref)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("fixedaddress/*"))
	})
	It("Should delete a fixed address", func() {
		client := "MAC_ADDRESS"
		fixedAddress := ibclient.NewFixedAddress("", "", "12.0.0.3", "", "43:56:98:98:32:21", &client, nil, "", false, "", nil, nil, nil, nil, false, nil, false)
		ref, err := connector.CreateObject(fixedAddress)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("fixedaddress/*"))
		delRef, err := connector.DeleteObject(ref)
		Expect(err).To(BeNil())
		Expect(delRef).To(MatchRegexp("fixedaddress/*"))
	})
	// get IPV4 fixed address, -ve scenario
	It("Should fail to get a non existent fixed address", func() {
		var res []ibclient.FixedAddress
		err := connector.GetObject(&ibclient.FixedAddress{}, "nonexistent_ref", nil, &res)
		Expect(res).To(BeNil())
		Expect(err).NotTo(BeNil())
	})
	// update IPV4 fixed address, -ve scenario
	It("Should fail to update a Fixed address", func() {
		client := "MAC_ADDRESS"
		fixedAddress := ibclient.NewFixedAddress("", "", "12.0.0.3", "", "43:56:98:98:32:21", &client, nil, "", false, "", nil, nil, nil, nil, false, nil, false)
		ref, err := connector.CreateObject(fixedAddress)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("fixedaddress/*"))
		_, err = connector.UpdateObject(fixedAddress, "nonexistent_ref")
		Expect(err).NotTo(BeNil())
	})
	// delete IPV4 Fixed address, -ve scenario
	It("Should fail to delete a non existent Fixed address", func() {
		_, err := connector.DeleteObject("nonexistent_ref")
		Expect(err).NotTo(BeNil())
	})

})

var _ = Describe("SharedNetwork Record", func() {
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

	It("Should create SharedNetwork record", func() {
		// create a sharedNetwork with maximum parameters
		ipv4Network1 := ibclient.NewNetwork("default", "22.23.24.0/24", false, "ipv4 network", nil)
		ipv4NetworkRef1, err := connector.CreateObject(ipv4Network1)
		Expect(err).To(BeNil())

		options := []*ibclient.Dhcpoption{
			{
				Name:  "domain-name-servers",
				Value: "11.22.3.1",
			},
			{
				Name:  "domain-name",
				Value: "aa.mm.ee",
			},
		}
		sharedNetwork := ibclient.SharedNetwork{
			Name:        utils.StringPtr("shared-network1"),
			Networks:    []*ibclient.Ipv4Network{{Ref: ipv4NetworkRef1}},
			Options:     options,
			Comment:     utils.StringPtr("sharedNetwork"),
			Disable:     utils.BoolPtr(false),
			UseOptions:  utils.BoolPtr(true),
			Ea:          ibclient.EA{"Site": "Baku"},
			NetworkView: "default",
		}
		ref, err := connector.CreateObject(&sharedNetwork)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^sharednetwork.*"))

		var ipv4SharedNetwork *ibclient.SharedNetwork
		err = connector.GetObject(&ibclient.SharedNetwork{}, ref, nil, &ipv4SharedNetwork)
		Expect(err).To(BeNil())
		Expect(ipv4SharedNetwork).NotTo(BeNil())
	})

	It("Should fail to create SharedNetwork record", func() {

		// Create a sharedNetwork object without mandatory fields
		sharedNetwork := ibclient.SharedNetwork{
			Name: utils.StringPtr("shared-network1"),
		}
		_, err := connector.CreateObject(&sharedNetwork)
		Expect(err).ToNot(BeNil())
	})

	It("Should Get SharedNetwork record", func() {
		// Create a sharedNetwork record and compare its fields
		ipv4Cidr := "23.23.24.0/24"
		ipv4Network1 := ibclient.NewNetwork("default", ipv4Cidr, false, "ipv4 network", nil)
		ipv4NetworkRef1, err := connector.CreateObject(ipv4Network1)
		Expect(err).To(BeNil())

		sharedNetwork := ibclient.SharedNetwork{
			Name:     utils.StringPtr("shared-network2"),
			Networks: []*ibclient.Ipv4Network{{Ref: ipv4NetworkRef1}},
			Ea:       ibclient.EA{"Site": "Baku"},
		}
		ref, err := connector.CreateObject(&sharedNetwork)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^sharednetwork.*"))

		var res *ibclient.SharedNetwork
		// Get by Ref
		err = connector.GetObject(&ibclient.SharedNetwork{}, ref, nil, &res)
		Expect(err).To(BeNil())
		Expect(res).NotTo(BeNil())

		Expect(*res.Name).To(Equal("shared-network2"))
		Expect(res.Networks[0].Ref).To(Equal(ipv4NetworkRef1))

		// Get by name
		var result []ibclient.SharedNetwork
		qp := ibclient.NewQueryParams(false, map[string]string{"name": "shared-network2"})
		err = connector.GetObject(&ibclient.SharedNetwork{}, "", qp, &result)
		Expect(err).To(BeNil())
		Expect(result).NotTo(BeNil())

		Expect(*result[0].Name).To(Equal("shared-network2"))
		Expect(result[0].Networks[0].Ref).To(Equal(ipv4NetworkRef1))
	})

	It("Should fail to Get SharedNetwork record", func() {
		// should fail to get a non-existent sharedNetwork record
		var res []ibclient.SharedNetwork
		sf := map[string]string{"name": "sharedNetwork10"}
		qp := ibclient.NewQueryParams(false, sf)
		err := connector.GetObject(&ibclient.SharedNetwork{}, "", qp, &res)
		Expect(res).To(BeEmpty())
		Expect(err).NotTo(BeNil())
	})

	It("Should update SharedNetwork record", func() {
		// Create a sharedNetwork object and update it
		ipv4Network1 := ibclient.NewNetwork("default", "24.23.24.0/24", false, "ipv4 network", nil)
		ipv4NetworkRef1, err := connector.CreateObject(ipv4Network1)
		Expect(err).To(BeNil())

		options := []*ibclient.Dhcpoption{
			{
				Name:  "domain-name-servers",
				Value: "11.22.3.1",
			},
			{
				Name:  "domain-name",
				Value: "aa.mm.ee",
			},
		}
		sharedNetwork := ibclient.SharedNetwork{
			Name:        utils.StringPtr("shared-network3"),
			Networks:    []*ibclient.Ipv4Network{{Ref: ipv4NetworkRef1}},
			Options:     options,
			Comment:     utils.StringPtr("sharedNetwork"),
			Disable:     utils.BoolPtr(false),
			UseOptions:  utils.BoolPtr(true),
			Ea:          ibclient.EA{"Site": "Baku"},
			NetworkView: "default",
		}
		ref, err := connector.CreateObject(&sharedNetwork)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^sharednetwork.*"))

		ipv4Network2 := ibclient.NewNetwork("default", "24.24.24.0/24", false, "ipv4 network", nil)
		ipv4NetworkRef2, err := connector.CreateObject(ipv4Network2)
		Expect(err).To(BeNil())

		sharedNetworkUpdated := ibclient.SharedNetwork{
			Name:       utils.StringPtr("shared-network-updated"),
			Networks:   []*ibclient.Ipv4Network{{Ref: ipv4NetworkRef2}},
			Comment:    utils.StringPtr("sharedNetwork updated"),
			UseOptions: utils.BoolPtr(false),
			Options:    nil,
		}
		updatedRef, err := connector.UpdateObject(&sharedNetworkUpdated, ref)
		Expect(err).To(BeNil())
		Expect(updatedRef).To(MatchRegexp("^sharednetwork.*"))

		Expect(*sharedNetworkUpdated.Name).To(Equal("shared-network-updated"))
		Expect(sharedNetworkUpdated.Networks[0].Ref).To(Equal(ipv4NetworkRef2))
		Expect(*sharedNetworkUpdated.Comment).To(Equal("sharedNetwork updated"))
		Expect(*sharedNetworkUpdated.UseOptions).To(BeFalse())
		Expect(len(sharedNetworkUpdated.Options)).To(Equal(0))
	})

	It("Should fail to update SharedNetwork record", func() {
		// Create a sharedNetwork record and try to update network_view field
		ipv4Network1 := ibclient.NewNetwork("default", "26.23.24.0/24", false, "ipv4 network", nil)
		ipv4NetworkRef1, err := connector.CreateObject(ipv4Network1)
		Expect(err).To(BeNil())

		options := []*ibclient.Dhcpoption{
			{
				Name:  "domain-name-servers",
				Value: "11.22.3.1",
			},
			{
				Name:  "domain-name",
				Value: "aa.mm.ee",
			},
		}
		sharedNetwork := ibclient.SharedNetwork{
			Name:        utils.StringPtr("shared-network4"),
			Networks:    []*ibclient.Ipv4Network{{Ref: ipv4NetworkRef1}},
			Options:     options,
			Comment:     utils.StringPtr("sharedNetwork"),
			Disable:     utils.BoolPtr(false),
			UseOptions:  utils.BoolPtr(true),
			Ea:          ibclient.EA{"Site": "Baku"},
			NetworkView: "default",
		}
		ref, err := connector.CreateObject(&sharedNetwork)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^sharednetwork.*"))

		ipv4Network2 := ibclient.NewNetwork("default", "24.24.24.0/24", false, "ipv4 network", nil)
		ipv4NetworkRef2, err := connector.CreateObject(ipv4Network2)
		Expect(err).To(BeNil())

		sharedNetworkUpdated := ibclient.SharedNetwork{
			Name:        utils.StringPtr("shared-network-updated"),
			Networks:    []*ibclient.Ipv4Network{{Ref: ipv4NetworkRef2}},
			Comment:     utils.StringPtr("sharedNetwork updated"),
			UseOptions:  utils.BoolPtr(false),
			Options:     nil,
			NetworkView: "view2",
		}
		_, err = connector.UpdateObject(&sharedNetworkUpdated, ref)
		Expect(err).ToNot(BeNil())
	})

	It("Should Delete SharedNetwork record", func() {
		// Create a sharedNetwork object and delete it
		ipv4Network1 := ibclient.NewNetwork("default", "25.23.24.0/24", false, "ipv4 network", nil)
		ipv4NetworkRef1, err := connector.CreateObject(ipv4Network1)
		Expect(err).To(BeNil())

		sharedNetwork := ibclient.SharedNetwork{
			Name:     utils.StringPtr("shared-network4"),
			Networks: []*ibclient.Ipv4Network{{Ref: ipv4NetworkRef1}},
			Ea:       ibclient.EA{"Site": "Baku"},
		}
		ref, err := connector.CreateObject(&sharedNetwork)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^sharednetwork.*"))

		deleteRef, err := connector.DeleteObject(ref)
		Expect(err).To(BeNil())
		Expect(deleteRef).To(Equal(ref))
	})

	It("Should fail to Delete SharedNetwork record", func() {
		// Delete a non-existent sharedNetwork object
		_, err := connector.DeleteObject("non-existent-sharedNetwork")
		Expect(err).ToNot(BeNil())
	})

})

var _ = Describe("Network Range Object", func() {
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
		ibClientConnector, err := ibclient.NewConnector(hostConfig, authConfig, transportConfig, requestBuilder, requestor)
		Expect(err).To(BeNil())
		connector = &ConnectorFacadeE2E{*ibClientConnector, make([]string, 0)}
		var (
			networkCidr = "60.0.0.0/24"
			networkView = "default"
		)
		network := ibclient.NewNetwork(networkView, networkCidr, false, "comment string", nil)
		network.Members = []ibclient.NetworkMember{
			{
				DhcpMember: &ibclient.Dhcpmember{Name: "infoblox.localdomain"},
			},
		}
		_, err = connector.CreateObject(network)
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		err := connector.SweepObjects()
		Expect(err).To(BeNil())
	})

	It("should create a Network Range with minimal parameters", func() {
		networkRange := ibclient.Range{
			StartAddr: utils.StringPtr("60.0.0.10"),
			EndAddr:   utils.StringPtr("60.0.0.20"),
		}
		ref, err := connector.CreateObject(&networkRange)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("range/*"))
	})
	It("should create a Network Range with maximal parameters", func() {
		options := []*ibclient.Dhcpoption{
			{
				Name:        "routers",
				Num:         3,
				VendorClass: "DHCP",
				Value:       "60.0.0.34",
				UseOption:   true,
			},
		}
		networkRange := ibclient.Range{
			StartAddr: utils.StringPtr("60.0.0.31"),
			EndAddr:   utils.StringPtr("60.0.0.35"),
			Comment:   utils.StringPtr("new comment"),
			Member: &ibclient.Dhcpmember{
				Name: "infoblox.localdomain",
			},
			ServerAssociationType: "MEMBER",
			Network:               utils.StringPtr("60.0.0.0/24"),
			NetworkView:           utils.StringPtr("default"),
			Disable:               utils.BoolPtr(true),
			Ea: ibclient.EA{
				"Site": "India",
			},
			Options:    options,
			UseOptions: utils.BoolPtr(true),
		}
		ref, err := connector.CreateObject(&networkRange)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("range/*"))
	})
	It("should get the Network Range object", func() {
		options := []*ibclient.Dhcpoption{
			{
				Name:        "routers",
				Num:         3,
				VendorClass: "DHCP",
				Value:       "60.0.0.34",
				UseOption:   true,
			},
		}
		networkRange := ibclient.Range{
			StartAddr: utils.StringPtr("60.0.0.31"),
			EndAddr:   utils.StringPtr("60.0.0.35"),
			Comment:   utils.StringPtr("new comment"),
			Name:      utils.StringPtr("range 1"),
			Member: &ibclient.Dhcpmember{
				Name: "infoblox.localdomain",
			},
			ServerAssociationType: "MEMBER",
			Network:               utils.StringPtr("60.0.0.0/24"),
			NetworkView:           utils.StringPtr("default"),
			Disable:               utils.BoolPtr(true),
			Ea: ibclient.EA{
				"Site": "India",
			},
			Options:    options,
			UseOptions: utils.BoolPtr(true),
		}
		ref, err := connector.CreateObject(&networkRange)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("range/*"))
		var res ibclient.Range
		search := &ibclient.Range{}
		search.SetReturnFields(append(search.ReturnFields(), "member", "options", "disable", "server_association_type", "name", "extattrs", "use_options"))
		errCode := connector.GetObject(search, ref, nil, &res)
		Expect(errCode).To(BeNil())
		Expect(res.Name).To(Equal(utils.StringPtr("range 1")))
		Expect(res.Member.Name).To(Equal("infoblox.localdomain"))
		Expect(res.Disable).To(Equal(utils.BoolPtr(true)))
		Expect(res.ServerAssociationType).To(Equal("MEMBER"))
		Expect(res.Ea["Site"]).To(Equal("India"))
		Expect(res.Options[0].Name).To(Equal("dhcp-lease-time"))
		Expect(res.Options[0].Num).To(Equal(uint32(51)))
		Expect(res.Options[0].Value).To(Equal("43200"))
		Expect(res.Options[0].UseOption).To(Equal(false))
		Expect(res.Options[0].VendorClass).To(Equal("DHCP"))
		Expect(res.Options[1].Name).To(Equal("routers"))
		Expect(res.Options[1].Num).To(Equal(uint32(3)))
		Expect(res.Options[1].Value).To(Equal("60.0.0.34"))
		Expect(res.Options[1].UseOption).To(Equal(true))
		Expect(res.Options[1].VendorClass).To(Equal("DHCP"))
		Expect(res.UseOptions).To(Equal(utils.BoolPtr(true)))
		Expect(res.Ref).To(MatchRegexp("range/*"))
		Expect(res.Comment).To(Equal(utils.StringPtr("new comment")))
	})
	It("should delete the network range ", func() {
		networkRange := ibclient.Range{
			StartAddr: utils.StringPtr("60.0.0.10"),
			EndAddr:   utils.StringPtr("60.0.0.20"),
		}
		ref, err := connector.CreateObject(&networkRange)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("range/*"))

		var res []ibclient.Range
		search := &ibclient.Range{}
		err = connector.GetObject(search, "", nil, &res)
		ref, err = connector.DeleteObject(res[0].Ref)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("range/*"))
	})
	It("Should fail to update a range with wrong reference", func() {
		networkRange := ibclient.Range{
			StartAddr: utils.StringPtr("60.0.0.10"),
			EndAddr:   utils.StringPtr("60.0.0.20"),
		}
		ref, err := connector.CreateObject(&networkRange)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("range/*"))
		_, err = connector.UpdateObject(&networkRange, "nonexistent_ref")
		Expect(err).NotTo(BeNil())
	})
	It("Should fail to update a range", func() {
		networkRange := ibclient.Range{
			StartAddr: utils.StringPtr("60.0.0.10"),
			EndAddr:   utils.StringPtr("60.0.0.20"),
		}
		ref, err := connector.CreateObject(&networkRange)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("range/*"))
		networkRangeUpdate := ibclient.Range{
			NetworkView: utils.StringPtr("network_view"),
		}
		var res []ibclient.Range
		search := &ibclient.Range{}
		err = connector.GetObject(search, "", nil, &res)
		ref, err = connector.UpdateObject(&networkRangeUpdate, res[0].Ref)
		Expect(err).NotTo(BeNil())
	})
	It("Should fail to create a range object", func() {
		networkRange := ibclient.Range{
			StartAddr: utils.StringPtr("60.0.0.10"),
		}
		_, err := connector.CreateObject(&networkRange)
		Expect(err).NotTo(BeNil())
	})
	It("Should fail to get a non-existent range", func() {
		var res []ibclient.Range
		sf := map[string]string{"name": "range"}
		qp := ibclient.NewQueryParams(false, sf)
		err := connector.GetObject(&ibclient.Range{}, "", qp, &res)
		Expect(res).To(BeEmpty())
		Expect(err).NotTo(BeNil())
	})
	It("Should fail to delete a non-existent range", func() {
		_, err := connector.DeleteObject("nonexistent_ref")
		Expect(err).NotTo(BeNil())
	})
})

var _ = Describe("SVCB Record", func() {
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

		zone := &ibclient.ZoneAuth{
			View: utils.StringPtr("default"),
			Fqdn: "test.com",
		}
		ref, err := connector.CreateObject(zone)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("zone_auth.*test.com/default"))

	})

	AfterEach(func() {
		err := connector.SweepObjects()
		Expect(err).To(BeNil())
	})

	It("Should create SVCB record", func() {
		// create a SVCB Record with maximum parameters
		recordSVCB := ibclient.RecordSVCB{
			Name:              "svcb-record-1.test.com",
			Comment:           "SVCB record",
			Creator:           "STATIC",
			Disable:           false,
			Ea:                map[string]interface{}{"Site": "Baku"},
			Priority:          100,
			Reclaimable:       false,
			TargetName:        "test123.com",
			Ttl:               300,
			UseTtl:            false,
			ForbidReclamation: true,
			SvcParameters: []ibclient.SVCParams{
				{Mandatory: false,
					SvcValue: []string{"443"},
					SvcKey:   "port",
				},
			},
			View: "default",
		}
		ref, err := connector.CreateObject(&recordSVCB)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^record:svcb.*"))

		var svcbRecord *ibclient.RecordSVCB
		err = connector.GetObject(&ibclient.RecordSVCB{}, ref, nil, &svcbRecord)
		Expect(err).To(BeNil())
		Expect(svcbRecord).NotTo(BeNil())
	})

	It("Should fail to create SVCB record", func() {

		// Create a SVCB Record without mandatory fields
		recordSvcb := ibclient.RecordSVCB{
			Name: "svcb-record-2.test.com",
		}
		_, err := connector.CreateObject(&recordSvcb)
		Expect(err).ToNot(BeNil())
	})

	It("Should Get SVCB record", func() {
		// Create a SVCB record and compare its fields
		recordSVCB := ibclient.RecordSVCB{
			Name:              "svcb-record-11.test.com",
			Comment:           "SVCB record",
			Creator:           "STATIC",
			Disable:           false,
			Ea:                map[string]interface{}{"Site": "Baku"},
			Priority:          300,
			Reclaimable:       false,
			TargetName:        "test1234.com",
			Ttl:               400,
			UseTtl:            false,
			ForbidReclamation: false,
			SvcParameters: []ibclient.SVCParams{
				{Mandatory: false,
					SvcValue: []string{"443"},
					SvcKey:   "port",
				},
			},
			View: "default",
		}
		ref, err := connector.CreateObject(&recordSVCB)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^record:svcb.*"))

		var res *ibclient.RecordSVCB
		// Get by Ref
		err = connector.GetObject(&ibclient.RecordSVCB{}, ref, nil, &res)
		Expect(err).To(BeNil())
		Expect(res).NotTo(BeNil())

		Expect(res.Name).To(Equal("svcb-record-11.test.com"))
		Expect(res.TargetName).To(Equal("test1234.com"))
		Expect(res.Priority).To(Equal(uint32(300)))

		// Get by name
		var result []ibclient.RecordSVCB
		qp := ibclient.NewQueryParams(false, map[string]string{"name": "svcb-record-11.test.com"})
		err = connector.GetObject(&ibclient.RecordSVCB{}, "", qp, &result)
		Expect(err).To(BeNil())
		Expect(result).NotTo(BeNil())

		Expect(result[0].Name).To(Equal("svcb-record-11.test.com"))
		Expect(res.TargetName).To(Equal("test1234.com"))
		Expect(res.Priority).To(Equal(uint32(300)))
	})

	It("Should fail to Get Record SVCB record", func() {
		// should fail to get a non-existent SVCB record
		var res []ibclient.RecordSVCB
		sf := map[string]string{"name": "record-svcb-1234"}
		qp := ibclient.NewQueryParams(false, sf)
		err := connector.GetObject(&ibclient.RecordSVCB{}, "", qp, &res)
		Expect(res).To(BeEmpty())
		Expect(err).NotTo(BeNil())
	})

	It("Should update SVCB Reecord", func() {
		// Create a SVCB Record and update it
		recordSVCB := ibclient.RecordSVCB{
			Name:        "svcb-record-1122.test.com",
			Comment:     "SVCB record",
			Creator:     "STATIC",
			Disable:     false,
			Ea:          map[string]interface{}{"Site": "Baku"},
			Priority:    300,
			Reclaimable: false,
			TargetName:  "test1234.com",
			Ttl:         400,
			UseTtl:      false,
			SvcParameters: []ibclient.SVCParams{
				{Mandatory: false,
					SvcValue: []string{"443"},
					SvcKey:   "port",
				},
			},
			View: "default",
		}
		ref, err := connector.CreateObject(&recordSVCB)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^record:svcb.*"))

		recordSvcbUpdated := ibclient.RecordSVCB{
			Name:        "svcb-record-1144.test.com",
			Comment:     "SVCB record updated",
			Creator:     "STATIC",
			Disable:     false,
			Ea:          map[string]interface{}{"Site": "Osaka"},
			Priority:    500,
			Reclaimable: false,
			TargetName:  "test1234.com",
			Ttl:         800,
			UseTtl:      true,
			SvcParameters: []ibclient.SVCParams{
				{Mandatory: false,
					SvcValue: []string{"500"},
					SvcKey:   "port",
				},
			},
		}
		updatedRef, err := connector.UpdateObject(&recordSvcbUpdated, ref)
		Expect(err).To(BeNil())
		Expect(updatedRef).To(MatchRegexp("^record:svcb.*"))

		Expect(recordSvcbUpdated.Name).To(Equal("svcb-record-1144.test.com"))
		Expect(recordSvcbUpdated.Comment).To(Equal("SVCB record updated"))
		Expect(recordSvcbUpdated.Priority).To(Equal(uint32(500)))
		Expect(recordSvcbUpdated.TargetName).To(Equal("test1234.com"))
		Expect(recordSvcbUpdated.Ttl).To(Equal(uint32(800)))
		Expect(recordSvcbUpdated.UseTtl).To(BeTrue())
		Expect(recordSvcbUpdated.Ea["Site"]).To(Equal("Osaka"))
		Expect(len(recordSvcbUpdated.SvcParameters)).To(Equal(1))
		Expect(recordSvcbUpdated.SvcParameters[0].SvcKey).To(Equal("port"))
		Expect(recordSvcbUpdated.SvcParameters[0].SvcValue[0]).To(Equal("500"))
		Expect(recordSvcbUpdated.SvcParameters[0].Mandatory).To(BeFalse())
		Expect(recordSvcbUpdated.Reclaimable).To(BeFalse())
	})

	It("Should fail to update SVCB Record", func() {
		// Create a SVCB Record and try to update view field
		recordSvcb := ibclient.RecordSVCB{
			Name:        "svcb-record-1143.test.com",
			Comment:     "SVCB record updated",
			Creator:     "STATIC",
			Disable:     false,
			Ea:          map[string]interface{}{"Site": "Osaka"},
			Priority:    500,
			Reclaimable: false,
			TargetName:  "test1234.com",
			Ttl:         800,
			UseTtl:      true,
			SvcParameters: []ibclient.SVCParams{
				{Mandatory: false,
					SvcValue: []string{"500"},
					SvcKey:   "port",
				},
			},
			View: "default",
		}
		ref, err := connector.CreateObject(&recordSvcb)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^record:svcb.*"))

		recordSvcbUpdated := ibclient.RecordSVCB{
			Name:        "svcb-record-1144.test.com",
			Comment:     "SVCB record updated",
			Creator:     "STATIC",
			Disable:     false,
			Ea:          map[string]interface{}{"Site": "Osaka"},
			Priority:    500,
			Reclaimable: false,
			TargetName:  "test1234.com",
			Ttl:         800,
			UseTtl:      true,
			SvcParameters: []ibclient.SVCParams{
				{
					Mandatory: false,
					SvcValue:  []string{"500"},
					SvcKey:    "port",
				},
			},
			View: "custom",
		}
		_, err = connector.UpdateObject(&recordSvcbUpdated, ref)
		Expect(err).ToNot(BeNil())
	})

	It("Should Delete SVCB Record", func() {
		// Create a SVCB Record and delete it
		recordSVCB := ibclient.RecordSVCB{
			Name:       "svcb-record-1111.test.com",
			Priority:   300,
			TargetName: "test1234.com",
		}
		ref, err := connector.CreateObject(&recordSVCB)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("^record:svcb.*"))

		deleteRef, err := connector.DeleteObject(ref)
		Expect(err).To(BeNil())
		Expect(deleteRef).To(Equal(ref))
	})

	It("Should fail to Delete SVCB Record", func() {
		// Delete a non-existent SVCB Record
		_, err := connector.DeleteObject("non-existent-record-svcb")
		Expect(err).ToNot(BeNil())
	})

})

var _ = Describe("HTTPS Record Object", func() {
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
		ibClientConnector, err := ibclient.NewConnector(hostConfig, authConfig, transportConfig, requestBuilder, requestor)
		Expect(err).To(BeNil())
		connector = &ConnectorFacadeE2E{*ibClientConnector, make([]string, 0)}
		zones := ibclient.ZoneAuth{
			Fqdn: "testing-https.com",
		}
		zoneRef, err := connector.CreateObject(&zones)
		Expect(err).To(BeNil())
		zones.Ref = zoneRef
	})

	AfterEach(func() {
		err := connector.SweepObjects()
		Expect(err).To(BeNil())
	})
	It("should create a Https record  with minimal parameters", func() {
		httpsRecord := ibclient.RecordHttps{
			Name:       "a1.testing-https.com",
			TargetName: "testing-https.com",
			Priority:   30,
		}
		ref, err := connector.CreateObject(&httpsRecord)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("record:https/*"))
	})
	It("should create a Https record with maximal parameters", func() {
		httpsRecord := ibclient.RecordHttps{
			Name:       "a2.testing-https.com",
			TargetName: "testing-https.com",
			Priority:   30,
			Comment:    "test comment",
			Disable:    false,
			View:       "default",
			UseTtl:     true,
			Ttl:        60,
			SvcParameters: []ibclient.SVCParams{
				{
					SvcKey: "port",
					SvcValue: []string{
						"233"},
					Mandatory: false,
				},
			},
			Ea:                ibclient.EA{"Site": "Bangalore"},
			ForbidReclamation: false,
			Creator:           "SYSTEM",
		}
		ref, err := connector.CreateObject(&httpsRecord)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("record:https/*"))
	})
	It("should GET https record ", func() {
		httpsRecord := ibclient.RecordHttps{
			Name:       "a2.testing-https.com",
			TargetName: "testing-https.com",
			Priority:   30,
			Comment:    "test comment",
			Disable:    false,
			View:       "default",
			UseTtl:     true,
			Ttl:        60,
			SvcParameters: []ibclient.SVCParams{
				{
					SvcKey: "port",
					SvcValue: []string{
						"233"},
					Mandatory: false,
				},
			},
			Ea:                ibclient.EA{"Site": "Bangalore"},
			ForbidReclamation: false,
			Creator:           "SYSTEM",
		}
		ref, err := connector.CreateObject(&httpsRecord)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("record:https/*"))
		var res ibclient.RecordHttps
		search := &ibclient.RecordHttps{}
		search.SetReturnFields([]string{"name","comment","disable", "view", "use_ttl", "ttl", "svc_parameters", "extattrs", "forbid_reclamation", "creator","target_name","priority"})
		errCode := connector.GetObject(search, ref, nil, &res)
		Expect(errCode).To(BeNil())
		Expect(res.Name).To(Equal("a2.testing-https.com"))
		Expect(res.TargetName).To(Equal("testing-https.com"))
		Expect(res.Priority).To(Equal(uint32(30)))
		Expect(res.Comment).To(Equal("test comment"))
		Expect(res.Disable).To(Equal(false))
		Expect(res.View).To(Equal("default"))
		Expect(res.UseTtl).To(Equal(true))
		Expect(res.Ttl).To(Equal(uint32(60)))
		Expect(res.SvcParameters[0].SvcKey).To(Equal("port"))
		Expect(res.SvcParameters[0].SvcValue[0]).To(Equal("233"))
		Expect(res.SvcParameters[0].Mandatory).To(Equal(false))
		Expect(res.Ea["Site"]).To(Equal("Bangalore"))
		Expect(res.ForbidReclamation).To(Equal(false))
		Expect(res.Creator).To(Equal("SYSTEM"))
		Expect(res.Ref).To(MatchRegexp("record:https/*"))
	})
	It("should delete the HTTPS Record ", func() {
		httpsRecord := ibclient.RecordHttps{
			Name:       "a1.testing-https.com",
			TargetName: "testing-https.com",
			Priority:   30,
		}
		ref, err := connector.CreateObject(&httpsRecord)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("record:https/*"))

		var res []ibclient.RecordHttps
		search := &ibclient.RecordHttps{}
		err = connector.GetObject(search, "", nil, &res)
		ref, err = connector.DeleteObject(res[0].Ref)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("record:https/*"))
	})
	It("Should fail to update a HTTPS Record with wrong reference", func() {
		httpsRecord := ibclient.RecordHttps{
			Name:       "a1.testing-https.com",
			TargetName: "testing-https.com",
			Priority:   30,
		}
		ref, err := connector.CreateObject(&httpsRecord)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("record:https/*"))
		_, err = connector.UpdateObject(&httpsRecord, "nonexistent_ref")
		Expect(err).NotTo(BeNil())
	})
	It("Should fail to create a https object", func() {
		httpsRecord := ibclient.RecordHttps{
			Name:       "a1.testing-https.com",
			TargetName: "testing-https.com",
		}
		_, err := connector.CreateObject(&httpsRecord)
		Expect(err).NotTo(BeNil())
	})
	It("Should fail to get a non-existent https record", func() {
		var res []ibclient.RecordHttps
		sf := map[string]string{"name": "range"}
		qp := ibclient.NewQueryParams(false, sf)
		err := connector.GetObject(&ibclient.RecordHttps{}, "", qp, &res)
		Expect(res).To(BeEmpty())
		Expect(err).NotTo(BeNil())
	})
	It("Should fail to delete a non-existent range", func() {
		_, err := connector.DeleteObject("nonexistent_ref")
		Expect(err).NotTo(BeNil())
	})
	It("Should update Https record", func() {
		httpsRecord := ibclient.RecordHttps{
			Name:       "a2.testing-https.com",
			TargetName: "testing-https.com",
			Priority:   30,
			Comment:    "test comment",
			Disable:    false,
			View:       "default",
			UseTtl:     true,
			Ttl:        60,
			SvcParameters: []ibclient.SVCParams{
				{
					SvcKey: "port",
					SvcValue: []string{
						"233"},
					Mandatory: false,
				},
			},
			Ea:                ibclient.EA{"Site": "Bangalore"},
			ForbidReclamation: false,
			Creator:           "SYSTEM",
		}
		ref, err := connector.CreateObject(&httpsRecord)
		Expect(err).To(BeNil())
		Expect(ref).To(MatchRegexp("record:https/*"))

		httpsRecord1 := ibclient.RecordHttps{
			Name:              "a3.testing-https.com",
			TargetName:        "testing-https.com",
			Priority:          30,
			Comment:           "test comment",
			Disable:          false,
			UseTtl:            true,
			Ttl:               60,
			SvcParameters:     []ibclient.SVCParams{},
			Ea:                ibclient.EA{"Site": "India"},
			ForbidReclamation: false,
			Creator:           "SYSTEM",
		}
		updatedRef, err := connector.UpdateObject(&httpsRecord1, ref)
		httpsRecord1.Ref = updatedRef
		Expect(err).To(BeNil())
		Expect(updatedRef).To(MatchRegexp("record:https/*"))
		Expect(httpsRecord1.Name).To(Equal("a3.testing-https.com"))
		Expect(httpsRecord1.TargetName).To(Equal("testing-https.com"))
		Expect(httpsRecord1.Priority).To(Equal(uint32(30)))
		Expect(httpsRecord1.Comment).To(Equal("test comment"))
		Expect(httpsRecord1.Disable).To(Equal(false))
		Expect(httpsRecord1.UseTtl).To(Equal(true))
		Expect(httpsRecord1.Ttl).To(Equal(uint32(60)))
		Expect(httpsRecord1.Ea["Site"]).To(Equal("India"))
		Expect(httpsRecord1.ForbidReclamation).To(Equal(false))
		Expect(httpsRecord1.Creator).To(Equal("SYSTEM"))
		Expect(httpsRecord1.Ref).To(MatchRegexp("record:https/*"))
	})
})
