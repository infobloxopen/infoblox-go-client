package ibclient

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Objects", func() {

	Context("Grid object", func() {

		tesNtpserver := NTPserver{
			Address:              "16.4.1.2",
			Burst:                true,
			EnableAuthentication: true,
			IBurst:               true,
			Preffered:            true,
		}
		grid := Grid{Name: "test", NTPSetting: &NTPSetting{EnableNTP: true,
			NTPAcl:     nil,
			NTPKeys:    nil,
			NTPKod:     false,
			NTPServers: []NTPserver{tesNtpserver},
		},
		}
		gridJSON := `{
			"name": "test",
			"ntp_setting": {
				"enable_ntp": true,
				"ntp_servers": [{
					"address": "16.4.1.2",
					"burst": true,
					"enable_authentication": true,
					"iburst": true,
					"preffered": true
					}]
				}
				}`

		Context("Marshalling", func() {
			Context("expected JSON is returned", func() {
				js, err := json.Marshal(grid)

				It("should not error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("should match json expected", func() {
					Expect(js).To(MatchJSON(gridJSON))
				})
			})
		})

		Context("Unmarshalling", func() {
			Context("expected object is returned", func() {
				var actualGrid Grid
				err := json.Unmarshal([]byte(gridJSON), &actualGrid)

				It("should not error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("should match object expected", func() {
					Expect(actualGrid).To(Equal(grid))
				})
			})
		})

	})

	Context("EA Object", func() {

		ea := EA{
			"Cloud API Owned":   Bool(true),
			"Tenant Name":       "Engineering01",
			"Maximum Wait Time": 120,
			"DNS Support":       Bool(false),
			"Routers":           []string{"10.1.2.234", "10.1.2.235"},
		}
		eaJSON := `{"Cloud API Owned":{"value":"True"},` +
			`"Tenant Name":{"value":"Engineering01"},` +
			`"Maximum Wait Time":{"value":120},` +
			`"DNS Support":{"value":"False"},` +
			`"Routers":{"value":["10.1.2.234", "10.1.2.235"]}}`

		Context("Marshalling", func() {
			Context("expected JSON is returned", func() {
				js, err := json.Marshal(ea)

				It("should not error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("should match json expected", func() {
					Expect(js).To(MatchJSON(eaJSON))
				})
			})
		})

		Context("Unmarshalling", func() {
			Context("expected object is returned", func() {
				var actualEA EA
				err := json.Unmarshal([]byte(eaJSON), &actualEA)

				It("should not error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("should match object expected", func() {
					Expect(actualEA).To(Equal(ea))
				})
			})
		})

	})

	Context("EA Search Object", func() {
		eas := EASearch{
			"Network Name": "Shared-Net",
			"Network View": "Global",
		}
		expectedJSON := `{"*Network Name" :"Shared-Net",` +
			`"*Network View" :"Global"}`

		Context("Marshalling", func() {
			Context("expected JSON is returned", func() {
				js, err := json.Marshal(eas)

				It("should not error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("should match json expected", func() {
					Expect(js).To(MatchJSON(expectedJSON))
				})
			})
		})
	})

	Context("EADefListValue Object", func() {
		var eadListVal EADefListValue = "Host Record"

		eadListValJSON := `{"value": "Host Record"}`

		Context("Marshalling", func() {
			Context("expected JSON is returned", func() {
				js, err := json.Marshal(eadListVal)

				It("should not error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("should match json expected", func() {
					Expect(js).To(MatchJSON(eadListValJSON))
				})
			})
		})

		Context("Unmarshalling", func() {
			Context("expected object is returned", func() {
				var actualEadListVal EADefListValue
				err := json.Unmarshal([]byte(eadListValJSON), &actualEadListVal)

				It("should not error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("should match object expected", func() {
					Expect(actualEadListVal).To(Equal(eadListVal))
				})
			})
		})

	})

	Context("Instantiation of", func() {
		Context("NetworkView object", func() {
			name := "myview"
			comment := "test client"
			setEas := EA{"Tenant ID": "client"}
			ref := "networkview/ZG5zLm5ldHdvcmtfdmlldyQyMw:global_view/false"
			nv := NewNetworkView(name, comment, setEas, ref)

			It("should set fields correctly", func() {
				Expect(nv.Name).To(Equal(name))
				Expect(nv.Comment).To(Equal(comment))
				Expect(nv.Ea).To(Equal(setEas))
				Expect(nv.Ref).To(Equal(ref))
			})

			It("should set base fields correctly", func() {
				Expect(nv.ObjectType()).To(Equal("networkview"))
				Expect(nv.ReturnFields()).To(ConsistOf("extattrs", "name", "comment"))
			})
		})

		Context("Network object", func() {
			cidr := "123.0.0.0/24"
			netviewName := "localview"
			comment := "test"
			ea := EA{"Tenant Name": "Engineering"}
			nw := NewNetwork(netviewName, cidr, false, comment, ea)
			searchEAs := EASearch{"Network Name": "shared-net"}
			nw.eaSearch = searchEAs

			It("should set fields correctly", func() {
				Expect(nw.Cidr).To(Equal(cidr))
				Expect(nw.NetviewName).To(Equal(netviewName))
				Expect(nw.Comment).To(Equal(comment))
				Expect(nw.Ea).To(Equal(ea))
			})

			It("should set base fields correctly", func() {
				Expect(nw.ObjectType()).To(Equal("network"))
				Expect(nw.ReturnFields()).To(ConsistOf("extattrs", "network", "network_view", "comment"))
				Expect(nw.EaSearch()).To(Equal(searchEAs))
			})
		})

		Context("IPv6 Network object", func() {
			cidr := "fc00::0100/56"
			netviewName := "localview"
			comment := "test"
			ea := EA{"Tenant Name": "Engineering"}
			nw := NewNetwork(netviewName, cidr, true, comment, ea)
			searchEAs := EASearch{"Network Name": "shared-net"}
			nw.eaSearch = searchEAs

			It("should set fields correctly", func() {
				Expect(nw.Cidr).To(Equal(cidr))
				Expect(nw.NetviewName).To(Equal(netviewName))
				Expect(nw.Comment).To(Equal(comment))
				Expect(nw.Ea).To(Equal(ea))
			})

			It("should set base fields correctly", func() {
				Expect(nw.ObjectType()).To(Equal("ipv6network"))
				Expect(nw.ReturnFields()).To(ConsistOf("extattrs", "network", "network_view", "comment"))
				Expect(nw.EaSearch()).To(Equal(searchEAs))
			})
		})

		Context("IPv4 NetworkContainer object", func() {
			cidr := "74.0.8.0/24"
			netviewName := "globalview"
			comment := "some comment"
			nwc := NewNetworkContainer(netviewName, cidr, false, comment, nil)

			It("should set fields correctly", func() {
				Expect(nwc.Cidr).To(Equal(cidr))
				Expect(nwc.NetviewName).To(Equal(netviewName))
				Expect(nwc.Comment).To(Equal(comment))
			})

			It("should set base fields correctly", func() {
				Expect(nwc.ObjectType()).To(Equal("networkcontainer"))
				Expect(nwc.ReturnFields()).To(ConsistOf("extattrs", "network", "network_view", "comment"))
			})
		})

		Context("IPv6 NetworkContainer object, with EAs", func() {
			cidr := "fc00::0100/56"
			netviewName := "default"
			eas := EA{
				"ea1": "ea1 value",
				"ea2": "ea2 value",
				"ea3 list": []string{
					"ea3 text1",
					"ea3 text2"}}
			comment := "some comment"
			nwc := NewNetworkContainer(netviewName, cidr, true, comment, eas)

			It("should set fields correctly", func() {
				Expect(nwc.Cidr).To(Equal(cidr))
				Expect(nwc.NetviewName).To(Equal(netviewName))
				Expect(nwc.Ea).To(Equal(eas))
				Expect(nwc.Comment).To(Equal(comment))
			})

			It("should set base fields correctly", func() {
				Expect(nwc.ObjectType()).To(Equal("ipv6networkcontainer"))
				Expect(nwc.ReturnFields()).To(ConsistOf("extattrs", "network", "network_view", "comment"))
			})
		})

		Context("FixedAddress object", func() {
			netviewName := "globalview"
			cidr := "25.0.7.0/24"
			ipAddress := "25.0.7.59/24"
			mac := "11:22:33:44:55:66"
			matchClient := "MAC_ADDRESS"
			comment := "test"
			ea := EA{"Tenant Name": "Engineering"}
			fixedAddr := NewFixedAddress(
				netviewName, "",
				ipAddress, cidr, mac,
				matchClient, ea, "", false, comment)

			It("should set fields correctly", func() {
				Expect(fixedAddr.NetviewName).To(Equal(netviewName))
				Expect(fixedAddr.Cidr).To(Equal(cidr))
				Expect(fixedAddr.IPv4Address).To(Equal(ipAddress))
				Expect(fixedAddr.Mac).To(Equal(mac))
				Expect(fixedAddr.MatchClient).To(Equal(matchClient))
				Expect(fixedAddr.Ea).To(Equal(ea))
			})

			It("should set base fields correctly", func() {
				Expect(fixedAddr.ObjectType()).To(Equal("fixedaddress"))
				Expect(fixedAddr.ReturnFields()).To(ConsistOf("extattrs", "ipv4addr", "mac", "name", "network", "network_view", "comment"))
			})
		})

		Context("IPv6 FixedAddress object", func() {
			netviewName := "globalview"
			cidr := "fc00::0100/56"
			ipAddress := "fc00::0100"
			duid := "11:22:33:44:55:66"
			comment := "test"
			ea := EA{"Tenant Name": "Engineering"}
			fixedAddr := NewFixedAddress(
				netviewName, "",
				ipAddress, cidr, duid,
				"", ea, "", true, comment)

			It("should set fields correctly", func() {
				Expect(fixedAddr.NetviewName).To(Equal(netviewName))
				Expect(fixedAddr.Cidr).To(Equal(cidr))
				Expect(fixedAddr.IPv6Address).To(Equal(ipAddress))
				Expect(fixedAddr.Duid).To(Equal(duid))
				Expect(fixedAddr.Ea).To(Equal(ea))
			})

			It("should set base fields correctly", func() {
				Expect(fixedAddr.ObjectType()).To(Equal("ipv6fixedaddress"))
				Expect(fixedAddr.ReturnFields()).To(ConsistOf("extattrs", "ipv6addr", "duid", "name", "network", "network_view", "comment"))
			})
		})

		Context("EADefinition object", func() {
			comment := "Test Extensible Attribute"
			flags := "CGV"
			listValues := []EADefListValue{"True", "False"}
			name := "Test EA"
			eaType := "string"
			allowedTypes := []string{"arecord", "aaarecord", "ptrrecord"}
			eaDef := NewEADefinition(EADefinition{
				Name:               name,
				Comment:            comment,
				Flags:              flags,
				ListValues:         listValues,
				Type:               eaType,
				AllowedObjectTypes: allowedTypes})

			It("should set fields correctly", func() {
				Expect(eaDef.Comment).To(Equal(comment))
				Expect(eaDef.Flags).To(Equal(flags))
				Expect(eaDef.ListValues).To(ConsistOf(listValues))
				Expect(eaDef.Name).To(Equal(name))
				Expect(eaDef.Type).To(Equal(eaType))
				Expect(eaDef.AllowedObjectTypes).To(ConsistOf(allowedTypes))
			})

			It("should set base fields correctly", func() {
				Expect(eaDef.ObjectType()).To(Equal("extensibleattributedef"))
				Expect(eaDef.ReturnFields()).To(ConsistOf("allowed_object_types", "comment", "flags", "list_values", "name", "type"))
			})
		})

		Context("UserProfile object", func() {
			userprofile := NewUserProfile(UserProfile{})

			It("should set base fields correctly", func() {
				Expect(userprofile.ObjectType()).To(Equal("userprofile"))
				Expect(userprofile.ReturnFields()).To(ConsistOf("name"))
			})
		})

		Context("RecordA object", func() {
			ipv4addr := "1.1.1.1"
			name := "bind_a.domain.com"
			view := "default"
			zone := "domain.com"
			ttl := uint32(500)
			useTTL := true
			comment := "testcomment"
			eas := EA{
				"TestEA1":  "testea1 value",
				"Location": "east coast",
			}

			ra := NewRecordA(view, zone, name, ipv4addr, ttl, useTTL, comment, eas, "")

			It("should set fields correctly", func() {
				Expect(ra.Ipv4Addr).To(Equal(ipv4addr))
				Expect(ra.Name).To(Equal(name))
				Expect(ra.View).To(Equal(view))
				Expect(ra.Zone).To(Equal(zone))

				Expect(ra.Ttl).To(Equal(ttl))
				Expect(ra.UseTtl).To(Equal(useTTL))
				Expect(ra.Comment).To(Equal(comment))
				Expect(ra.Ea).To(Equal(eas))
			})

			It("should set base fields correctly", func() {
				Expect(ra.ObjectType()).To(Equal("record:a"))
				Expect(ra.ReturnFields()).To(ConsistOf(
					"extattrs", "ipv4addr", "name", "view", "zone", "comment", "ttl", "use_ttl"))
			})
		})

		Context("RecordAAAA object", func() {
			ipv6addr := "2001:db8:abcd:14::1"
			name := "bind_a.domain.com"
			view := "default"
			useTtl := true
			ttl := uint32(10)
			comment := "test comment"
			ea := EA{"VM Name": "test-vm"}

			ra := NewRecordAAAA(view, name, ipv6addr, useTtl, ttl, comment, ea, "")

			It("should set fields correctly", func() {
				Expect(ra.Ipv6Addr).To(Equal(ipv6addr))
				Expect(ra.Name).To(Equal(name))
				Expect(ra.View).To(Equal(view))
				Expect(ra.UseTtl).To(Equal(useTtl))
				Expect(ra.Ttl).To(Equal(ttl))
				Expect(ra.Comment).To(Equal(comment))
				Expect(ra.Ea).To(Equal(ea))
			})

			It("should set base fields correctly", func() {
				Expect(ra.ObjectType()).To(Equal("record:aaaa"))
				Expect(ra.ReturnFields()).To(ConsistOf("extattrs", "ipv6addr", "name", "view", "zone", "use_ttl", "ttl", "comment"))
			})
		})

		Context("RecordPtr object", func() {
			ipv4addr := "1.1.1.1"
			ipv6addr := "2001::/64"
			ptrdname := "bind_a.domain.com"
			view := "default"
			zone := "domain.com"
			useTtl := true
			ttl := uint32(70)
			comment := "test client"
			eas := EA{"VM Name": "test"}

			rptr := NewRecordPTR(view, ptrdname, useTtl, ttl, comment, eas)
			rptr.Zone = zone
			rptr.Ipv4Addr = ipv4addr
			rptr.Ipv6Addr = ipv6addr

			It("should set fields correctly", func() {
				Expect(rptr.Ipv4Addr).To(Equal(ipv4addr))
				Expect(rptr.Ipv6Addr).To(Equal(ipv6addr))
				Expect(rptr.PtrdName).To(Equal(ptrdname))
				Expect(rptr.View).To(Equal(view))
				Expect(rptr.Zone).To(Equal(zone))
				Expect(rptr.UseTtl).To(Equal(useTtl))
				Expect(rptr.Ttl).To(Equal(ttl))
				Expect(rptr.Comment).To(Equal(comment))
				Expect(rptr.Ea).To(Equal(eas))
			})

			It("should set base fields correctly", func() {
				Expect(rptr.ObjectType()).To(Equal("record:ptr"))
				Expect(rptr.ReturnFields()).To(ConsistOf(
					"extattrs",
					"ipv4addr",
					"ipv6addr",
					"name",
					"ptrdname",
					"view",
					"zone",
					"comment",
					"use_ttl",
					"ttl"))
			})
		})

		Context("RecordCNAME object", func() {
			canonical := "cname.domain.com"
			name := "bind_cname.domain.com"
			useTtl := false
			ttl := uint32(0)
			view := "default"
			comment := "test CNAME"
			eas := EA{"VM Name": "test"}

			rc := NewRecordCNAME(view, canonical, name, useTtl, ttl, comment, eas, "")

			It("should set fields correctly", func() {
				Expect(rc.Canonical).To(Equal(canonical))
				Expect(rc.Name).To(Equal(name))
				Expect(rc.UseTtl).To(Equal(useTtl))
				Expect(rc.Ttl).To(Equal(ttl))
				Expect(rc.View).To(Equal(view))
				Expect(rc.Comment).To(Equal(comment))
				Expect(rc.Ea).To(Equal(eas))
			})

			It("should set base fields correctly", func() {
				Expect(rc.ObjectType()).To(Equal("record:cname"))
				Expect(rc.ReturnFields()).To(ConsistOf("extattrs", "canonical", "name", "view", "zone", "comment", "ttl", "use_ttl"))
			})
		})

		Context("RecordHostIpv4Addr object", func() {
			ipAddress := "25.0.7.59/24"
			mac := "11:22:33:44:55:66"
			enableDHCP := false
			hostAddr := NewHostRecordIpv4Addr(ipAddress, mac, enableDHCP, "")

			It("should set fields correctly", func() {
				Expect(hostAddr.Ipv4Addr).To(Equal(ipAddress))
				Expect(hostAddr.Mac).To(Equal(mac))
			})

			It("should set base fields correctly", func() {
				Expect(hostAddr.ObjectType()).To(Equal("record:host_ipv4addr"))
				//Expect(hostAddr.ReturnFields()).To(ConsistOf("configure_for_dhcp", "host", "ipv4addr", "mac"))
			})
		})

		Context("RecordHostIpv4Addr macaddress empty", func() {
			ipAddress := "25.0.7.59"
			enableDHCP := false
			hostAddr := NewHostRecordIpv4Addr(ipAddress, "", enableDHCP, "")

			It("should set fields correctly", func() {
				Expect(hostAddr.Ipv4Addr).To(Equal(ipAddress))
			})

			It("should set base fields correctly", func() {
				Expect(hostAddr.ObjectType()).To(Equal("record:host_ipv4addr"))
				//Expect(hostAddr.ReturnFields()).To(ConsistOf("configure_for_dhcp", "host", "ipv4addr", "mac"))
			})
		})

		Context("RecordHostIpv6Addr object", func() {
			ipAddress := "fc00::0100"
			duid := "11:22:33:44:55:66"
			enableDHCP := false
			hostAddr := NewHostRecordIpv6Addr(ipAddress, duid, enableDHCP, "")

			It("should set fields correctly", func() {
				Expect(hostAddr.Ipv6Addr).To(Equal(ipAddress))
				Expect(hostAddr.Duid).To(Equal(duid))
			})

			It("should set base fields correctly", func() {
				Expect(hostAddr.ObjectType()).To(Equal("record:host_ipv6addr"))
			})
		})

		Context("RecordHostIpv6Addr duid empty", func() {
			ipAddress := "fc00::0100"
			enableDHCP := false
			hostAddr := NewHostRecordIpv6Addr(ipAddress, "", enableDHCP, "")

			It("should set fields correctly", func() {
				Expect(hostAddr.Ipv6Addr).To(Equal(ipAddress))
			})

			It("should set base fields correctly", func() {
				Expect(hostAddr.ObjectType()).To(Equal("record:host_ipv6addr"))
			})
		})

		Context("RecordHost object", func() {
			ipv4addrs := []HostRecordIpv4Addr{{Ipv4Addr: "1.1.1.1"}, {Ipv4Addr: "2.2.2.2"}}
			ipv6addrs := []HostRecordIpv6Addr{{Ipv6Addr: "fc00::0100"}, {Ipv6Addr: "fc00::0101"}}
			name := "bind_host.domain.com"
			view := "default"
			zone := "domain.com"
			useTtl := true
			ttl := uint32(70)
			comment := "test"
			aliases := []string{"bind_host1.domain.com"}

			rh := NewHostRecord(
				"", name, "", "", ipv4addrs, ipv6addrs,
				nil, true, view, zone, "", useTtl, ttl, comment, aliases)

			It("should set fields correctly", func() {
				Expect(rh.Ipv4Addrs).To(Equal(ipv4addrs))
				Expect(rh.Ipv6Addrs).To(Equal(ipv6addrs))
				Expect(rh.Name).To(Equal(name))
				Expect(rh.View).To(Equal(view))
				Expect(rh.Zone).To(Equal(zone))
				Expect(rh.Comment).To(Equal(comment))
				Expect(rh.Aliases).To(Equal(aliases))
			})

			It("should set base fields correctly", func() {
				Expect(rh.ObjectType()).To(Equal("record:host"))
				Expect(rh.ReturnFields()).To(ConsistOf("extattrs", "ipv4addrs", "ipv6addrs", "name", "view", "zone",
					"comment", "network_view", "aliases", "use_ttl", "ttl", "configure_for_dns"))
			})
		})

		Context("RecordMX object", func() {
			fqdn := "test.example.com"
			mx := "example.com"
			dnsView := "default"
			priority := uint32(10)
			ttl := uint32(70)
			useTtl := true
			comment := "test comment"
			eas := EA{"Country": "test"}

			rm := NewRecordMX(RecordMX{
				Fqdn:       fqdn,
				MX:         mx,
				View:       dnsView,
				Preference: priority,
				Ttl:        ttl,
				UseTtl:     useTtl,
				Comment:    comment,
				Ea:         eas,
			})

			It("should set fields correctly", func() {
				Expect(rm.Fqdn).To(Equal(fqdn))
				Expect(rm.MX).To(Equal(mx))
				Expect(rm.View).To(Equal(dnsView))
				Expect(rm.Preference).To(Equal(priority))
				Expect(rm.Ttl).To(Equal(ttl))
				Expect(rm.UseTtl).To(Equal(useTtl))
				Expect(rm.Comment).To(Equal(comment))
				Expect(rm.Ea).To(Equal(eas))
			})

			It("should set base fields correctly", func() {
				Expect(rm.ObjectType()).To(Equal("record:mx"))
				Expect(rm.ReturnFields()).To(ConsistOf("mail_exchanger", "view", "name", "preference", "ttl", "use_ttl", "comment", "extattrs", "zone"))
			})
		})

		Context("RecordSRV object", func() {
			name := "srv.sample.com"
			dnsView := "default"
			priority := uint32(10)
			weight := uint32(24)
			port := uint32(88)
			target := "h1.sample.com"
			ttl := uint32(300)
			useTtl := true
			comment := "test comment"
			eas := EA{"Country": "test"}

			rv := NewRecordSRV(RecordSRV{
				View:     dnsView,
				Name:     name,
				Priority: priority,
				Weight:   weight,
				Port:     port,
				Target:   target,
				Ttl:      ttl,
				UseTtl:   useTtl,
				Comment:  comment,
				Ea:       eas,
			})

			It("should set field correctly", func() {
				Expect(rv.View).To(Equal(dnsView))
				Expect(rv.Name).To(Equal(name))
				Expect(rv.Priority).To(Equal(priority))
				Expect(rv.Weight).To(Equal(weight))
				Expect(rv.Port).To(Equal(port))
				Expect(rv.Target).To(Equal(target))
				Expect(rv.Ttl).To(Equal(ttl))
				Expect(rv.UseTtl).To(Equal(useTtl))
				Expect(rv.Comment).To(Equal(comment))
				Expect(rv.Ea).To(Equal(eas))
			})
			It("should set base fields correctly", func() {
				Expect(rv.ObjectType()).To(Equal("record:srv"))
				Expect(rv.ReturnFields()).To(ConsistOf("name", "view", "priority", "weight", "port", "target", "ttl", "use_ttl", "comment", "extattrs", "zone"))
			})

		})

		Context("RecordTXT object", func() {
			view := "default"
			name := "txt.domain.com"
			text := "this is text string"
			ttl := uint32(70)
			useTtl := true
			comment := "test client"
			eas := EA{"Country": "test"}

			rt := NewRecordTXT(view, "", name, text, ttl, useTtl, comment, eas)

			It("should set fields correctly", func() {
				Expect(rt.View).To(Equal(view))
				Expect(rt.Name).To(Equal(name))
				Expect(rt.Text).To(Equal(text))
			})

			It("should set base fields correctly", func() {
				Expect(rt.ObjectType()).To(Equal("record:txt"))
				Expect(rt.ReturnFields()).To(ConsistOf("view", "zone", "name", "text", "ttl", "use_ttl", "comment", "extattrs"))
			})
		})

		Context("ZoneAuth object", func() {
			fqdn := "domain.com"
			view := "default"

			za := NewZoneAuth(ZoneAuth{
				Fqdn: fqdn,
				View: view})

			It("should set fields correctly", func() {
				Expect(za.Fqdn).To(Equal(fqdn))
				Expect(za.View).To(Equal(view))
			})

			It("should set base fields correctly", func() {
				Expect(za.ObjectType()).To(Equal("zone_auth"))
				Expect(za.ReturnFields()).To(ConsistOf("extattrs", "fqdn", "view"))
			})
		})

		Context("ZoneDelegated object", func() {
			fqdn := "delegated_zone.domain.com"
			view := "default"

			za := NewZoneDelegated(ZoneDelegated{
				Fqdn: fqdn,
				View: view})

			It("should set fields correctly", func() {
				Expect(za.Fqdn).To(Equal(fqdn))
				Expect(za.View).To(Equal(view))
			})

			It("should set base fields correctly", func() {
				Expect(za.ObjectType()).To(Equal("zone_delegated"))
				Expect(za.ReturnFields()).To(ConsistOf("extattrs", "fqdn", "view", "delegate_to"))
			})
		})

	})

	Context("Unmarshalling malformed JSON", func() {
		Context("for EA", func() {
			badJSON := `""`
			var ea EA
			err := json.Unmarshal([]byte(badJSON), &ea)

			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})

		Context("for EADefListValue", func() {
			badJSON := `""`
			var ead EADefListValue
			err := json.Unmarshal([]byte(badJSON), &ead)

			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})

	})

})
