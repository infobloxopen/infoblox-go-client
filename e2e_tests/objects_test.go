package e2e_tests

import (
	ibclient "github.com/infobloxopen/infoblox-go-client/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Objects", func() {
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

	Describe("Network View", func() {
		It("Should properly serialize/deserialize", func() {
			nv := &ibclient.NetworkView{
				Name:    "e2e_test_view",
				Comment: "Network View created by e2e test",
			}

			ref, err := connector.CreateObject(nv)
			Expect(err).To(BeNil())

			var res ibclient.NetworkView
			err = connector.GetObject(nv, ref, nil, &res)
			Expect(err).To(BeNil())
			Expect(res.Ref).To(Equal(ref))
			Expect(res.Name).To(Equal(nv.Name))
			Expect(res.Comment).To(Equal(nv.Comment))

			nv.Comment = "Network View updated by e2e test"
			updRef, err := connector.UpdateObject(nv, ref)
			Expect(err).To(BeNil())

			_, err = connector.DeleteObject(updRef)
			Expect(err).To(BeNil())
		})
	})

	When("Network View exists", func() {
		BeforeEach(func() {
			nv := &ibclient.NetworkView{
				Name:    "e2e_test_view",
				Comment: "Network View created by e2e test",
			}

			_, err := connector.CreateObject(nv)
			Expect(err).To(BeNil())
		})

		Describe("IPv4 Network", func() {
			It("Should properly serialize/deserialize", func() {
				nw := &ibclient.Ipv4Network{
					NetworkView: "e2e_test_view",
					Network:     "192.168.1.0/24",
					Comment:     "IPv4 Network created by e2e test",
				}

				ref, err := connector.CreateObject(nw)
				Expect(err).To(BeNil())

				var res ibclient.Ipv4Network
				err = connector.GetObject(nw, ref, nil, &res)
				Expect(err).To(BeNil())
				Expect(res.Ref).To(Equal(ref))
				Expect(res.NetworkView).To(Equal(nw.NetworkView))
				Expect(res.Network).To(Equal(nw.Network))
				Expect(res.Comment).To(Equal(nw.Comment))

				nw.NetworkView = ""
				nw.Comment = "IPv4 Network updated by e2e test"
				updRef, err := connector.UpdateObject(nw, ref)
				Expect(err).To(BeNil())

				_, err = connector.DeleteObject(updRef)
				Expect(err).To(BeNil())
			})
		})

		When("IPv4 Network exists", func() {
			BeforeEach(func() {
				nw := &ibclient.Ipv4Network{
					NetworkView: "e2e_test_view",
					Network:     "192.168.1.0/24",
					Comment:     "IPv4 Network created by e2e test",
				}

				_, err := connector.CreateObject(nw)
				Expect(err).To(BeNil())
			})

			Describe("IPv4 Fixed Address", func() {
				It("Should properly serialize/deserialize", func() {
					fa := &ibclient.Ipv4FixedAddress{
						NetworkView: "e2e_test_view",
						Name:        "e2e_test_ipv4_fixed_address",
						Ipv4Addr:    "192.168.1.45",
						Mac:         "00:00:00:00:00:00",
						Comment:     "IPv4 Fixed Address created by e2e test",
					}

					ref, err := connector.CreateObject(fa)
					Expect(err).To(BeNil())

					fa.SetReturnFields(append(fa.ReturnFields(), "name", "mac", "comment"))
					var res ibclient.Ipv4FixedAddress
					err = connector.GetObject(fa, ref, nil, &res)
					Expect(err).To(BeNil())
					Expect(res.Ref).To(Equal(ref))
					Expect(res.NetworkView).To(Equal(fa.NetworkView))
					Expect(res.Name).To(Equal(fa.Name))
					Expect(res.Ipv4Addr).To(Equal(fa.Ipv4Addr))
					Expect(res.Mac).To(Equal(fa.Mac))
					Expect(res.Comment).To(Equal(fa.Comment))

					fa.Comment = "IPv4 Fixed Address updated by e2e test"
					updRef, err := connector.UpdateObject(fa, ref)
					Expect(err).To(BeNil())

					_, err = connector.DeleteObject(updRef)
					Expect(err).To(BeNil())
				})
			})

			Describe("IP Range", func() {
				It("Should properly serialize/deserialize", func() {
					r := &ibclient.Range{
						NetworkView: "e2e_test_view",
						Name:        "e2e_test_ip_range",
						Comment:     "IP Range created by e2e test",
						StartAddr:   "192.168.1.10",
						EndAddr:     "192.168.1.20",
					}

					ref, err := connector.CreateObject(r)
					Expect(err).To(BeNil())

					r.SetReturnFields(append(r.ReturnFields(), "name"))
					var res ibclient.Range
					err = connector.GetObject(r, ref, nil, &res)
					Expect(err).To(BeNil())
					Expect(res.Ref).To(Equal(ref))
					Expect(res.NetworkView).To(Equal(r.NetworkView))
					Expect(res.Name).To(Equal(r.Name))
					Expect(res.Comment).To(Equal(r.Comment))
					Expect(res.StartAddr).To(Equal(r.StartAddr))
					Expect(res.EndAddr).To(Equal(r.EndAddr))

					r.Comment = "IP Range updated by e2e test"
					updRef, err := connector.UpdateObject(r, ref)
					Expect(err).To(BeNil())

					_, err = connector.DeleteObject(updRef)
					Expect(err).To(BeNil())
				})
			})

			Describe("IPv4 Host Record", func() {
				It("Should properly serialize/deserialize", func() {
					By("Creating DNS view")
					v := &ibclient.View{
						Name:        "e2e_test_dns_view",
						NetworkView: "e2e_test_view",
						Comment:     "DNS View created by e2e test",
					}
					_, err := connector.CreateObject(v)
					Expect(err).To(BeNil())

					By("Creating forwarding-mapping DNS Auth Zone")
					zf := &ibclient.ZoneAuth{
						View:    "e2e_test_dns_view",
						Fqdn:    "e2e-test.com",
						Comment: "Forwarding-mapping DNS Auth Zone created by e2e test",
					}
					_, err = connector.CreateObject(zf)
					Expect(err).To(BeNil())

					By("Creating reverse-mapping DNS Auth Zone")
					zr := &ibclient.ZoneAuth{
						View:       "e2e_test_dns_view",
						Fqdn:       "192.168.1.0/24",
						ZoneFormat: "IPV4",
						Comment:    "Reverse-mapping DNS Auth Zone created by e2e test",
					}
					_, err = connector.CreateObject(zr)
					Expect(err).To(BeNil())

					fa := &ibclient.HostRecord{
						NetworkView: "e2e_test_view",
						View:        "e2e_test_dns_view",
						Name:        "e2e_test_host_record.e2e-test.com",
						Ipv4Addrs: []ibclient.HostRecordIpv4Addr{
							{
								Ipv4Addr: "192.168.1.60",
								Mac:      "00:00:00:00:00:00",
							},
						},
						Comment: "IPv4 Host Record created by e2e test",
					}

					ref, err := connector.CreateObject(fa)
					Expect(err).To(BeNil())

					fa.SetReturnFields(append(fa.ReturnFields(), "network_view", "comment"))
					var res ibclient.HostRecord
					err = connector.GetObject(fa, ref, nil, &res)
					Expect(err).To(BeNil())
					Expect(res.Ref).To(Equal(ref))
					Expect(res.NetworkView).To(Equal(fa.NetworkView))
					Expect(res.Name).To(Equal(fa.Name))
					Expect(res.Comment).To(Equal(fa.Comment))

					fa.NetworkView = ""
					fa.Comment = "IPv4 Host Record updated by e2e test"
					updRef, err := connector.UpdateObject(fa, ref)
					Expect(err).To(BeNil())

					_, err = connector.DeleteObject(updRef)
					Expect(err).To(BeNil())
				})
			})

		})

		Describe("IPv4 Network Container", func() {
			It("Should properly serialize/deserialize", func() {
				nc := &ibclient.Ipv4NetworkContainer{
					NetworkView: "e2e_test_view",
					Network:     "192.168.1.0/24",
					Comment:     "IPv4 Network Container created by e2e test",
				}

				ref, err := connector.CreateObject(nc)
				Expect(err).To(BeNil())

				var res ibclient.Ipv4FixedAddress
				err = connector.GetObject(nc, ref, nil, &res)
				Expect(err).To(BeNil())
				Expect(res.Ref).To(Equal(ref))
				Expect(res.NetworkView).To(Equal(nc.NetworkView))
				Expect(res.Network).To(Equal(nc.Network))
				Expect(res.Comment).To(Equal(nc.Comment))

				nc.Network = ""
				nc.NetworkView = ""
				nc.Comment = "IPv4 Network Container updated by e2e test"
				updRef, err := connector.UpdateObject(nc, ref)
				Expect(err).To(BeNil())

				_, err = connector.DeleteObject(updRef)
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("DNS View", func() {
		It("Should properly serialize/deserialize", func() {
			v := &ibclient.View{
				Name:    "e2e_test_dns_view",
				Comment: "DNS View created by e2e test",
			}

			ref, err := connector.CreateObject(v)
			Expect(err).To(BeNil())

			v.SetReturnFields([]string{"name", "comment"})

			var res ibclient.View
			err = connector.GetObject(v, ref, nil, &res)
			Expect(err).To(BeNil())
			Expect(res.Ref).To(Equal(ref))
			Expect(res.Name).To(Equal(v.Name))
			Expect(res.Comment).To(Equal(v.Comment))

			v.Comment = "DNS View updated by e2e test"
			updRef, err := connector.UpdateObject(v, ref)
			Expect(err).To(BeNil())

			_, err = connector.DeleteObject(updRef)
			Expect(err).To(BeNil())
		})
	})

	When("DNS View exists", func() {
		BeforeEach(func() {
			v := &ibclient.View{
				Name:    "e2e_test_dns_view",
				Comment: "DNS View created by e2e test",
			}

			_, err := connector.CreateObject(v)
			Expect(err).To(BeNil())
		})

		Describe("DNS Zone Auth", func() {
			It("Should support CRUD operations of forwarding-mapping zone", func() {
				z := &ibclient.ZoneAuth{
					View:    "e2e_test_dns_view",
					Fqdn:    "e2e-test.com",
					Comment: "DNS Auth Zone created by e2e test",
				}

				ref, err := connector.CreateObject(z)
				Expect(err).To(BeNil())

				z.SetReturnFields([]string{"view", "fqdn", "comment"})

				var res ibclient.ZoneAuth
				err = connector.GetObject(z, ref, nil, &res)
				Expect(err).To(BeNil())
				Expect(res.Ref).To(Equal(ref))
				Expect(res.Fqdn).To(Equal(z.Fqdn))
				Expect(res.View).To(Equal(z.View))
				Expect(res.Comment).To(Equal(z.Comment))

				z.Fqdn = ""
				z.Comment = "DNS Auth Zone updated by e2e test"
				updRef, err := connector.UpdateObject(z, ref)
				Expect(err).To(BeNil())

				_, err = connector.DeleteObject(updRef)
				Expect(err).To(BeNil())
			})

			It("Should support CRUD operations of reverse-mapping zone", func() {
				z := &ibclient.ZoneAuth{
					View:       "e2e_test_dns_view",
					Fqdn:       "192.168.1.0/24",
					ZoneFormat: "IPV4",
					Comment:    "DNS Auth Zone created by e2e test",
				}

				ref, err := connector.CreateObject(z)
				Expect(err).To(BeNil())

				z.SetReturnFields([]string{"view", "fqdn", "comment"})

				var res ibclient.ZoneAuth
				err = connector.GetObject(z, ref, nil, &res)
				Expect(err).To(BeNil())
				Expect(res.Ref).To(Equal(ref))
				Expect(res.Fqdn).To(Equal(z.Fqdn))
				Expect(res.View).To(Equal(z.View))
				Expect(res.Comment).To(Equal(z.Comment))

				z.Fqdn = ""
				z.ZoneFormat = ""
				z.Comment = "DNS Auth Zone updated by e2e test"
				updRef, err := connector.UpdateObject(z, ref)
				Expect(err).To(BeNil())

				_, err = connector.DeleteObject(updRef)
				Expect(err).To(BeNil())
			})
		})

		When("forwarding-mapping DNS Zone Exists", func() {
			BeforeEach(func() {
				z := &ibclient.ZoneAuth{
					View:    "e2e_test_dns_view",
					Fqdn:    "e2e-test.com",
					Comment: "Forwarding-mapping DNS Auth Zone created by e2e test",
				}

				_, err := connector.CreateObject(z)
				Expect(err).To(BeNil())
			})

			Describe("A Record", func() {
				It("Should properly serialize/deserialize", func() {
					a := &ibclient.RecordA{
						View:     "e2e_test_dns_view",
						Name:     "e2e_test_a_record.e2e-test.com",
						Ipv4Addr: "192.168.1.45",
						Ttl:      5,
						UseTtl:   true,
						Comment:  "A Record created by e2e test",
						Ea:       make(ibclient.EA),
					}

					ref, err := connector.CreateObject(a)
					Expect(err).To(BeNil())

					a.SetReturnFields([]string{"view", "comment", "creation_time"})
					var res ibclient.RecordA
					err = connector.GetObject(a, ref, nil, &res)
					Expect(err).To(BeNil())
					Expect(res.Ref).To(Equal(ref))
					Expect(res.View).To(Equal(a.View))
					Expect(res.Comment).To(Equal(a.Comment))
					Expect(res.CreationTime).To(Not(BeNil()))

					a.View = ""
					a.Comment = "A Record updated by e2e test"
					updRef, err := connector.UpdateObject(a, ref)
					Expect(err).To(BeNil())

					_, err = connector.DeleteObject(updRef)
					Expect(err).To(BeNil())
				})
			})

			Describe("AAAA Record", func() {
				It("Should properly serialize/deserialize", func() {
					aaaa := &ibclient.RecordAAAA{
						View:     "e2e_test_dns_view",
						Name:     "e2e_test_a_record.e2e-test.com",
						Ipv6Addr: "2001:db8:abcd:14::1",
						Ttl:      5,
						UseTtl:   true,
						Comment:  "A Record created by e2e test",
						Ea:       make(ibclient.EA),
					}

					ref, err := connector.CreateObject(aaaa)
					Expect(err).To(BeNil())

					aaaa.SetReturnFields([]string{"ipv6addr", "name", "ttl", "view", "comment"})
					var res ibclient.RecordAAAA
					err = connector.GetObject(aaaa, ref, nil, &res)
					Expect(err).To(BeNil())
					Expect(res.Ref).To(Equal(ref))
					Expect(res.Name).To(Equal(aaaa.Name))
					Expect(res.Ipv6Addr).To(Equal(aaaa.Ipv6Addr))
					Expect(res.Ttl).To(Equal(aaaa.Ttl))
					Expect(res.Comment).To(Equal(aaaa.Comment))

					aaaa.View = ""
					aaaa.Comment = "A Record updated by e2e test"
					updRef, err := connector.UpdateObject(aaaa, ref)
					Expect(err).To(BeNil())

					_, err = connector.DeleteObject(updRef)
					Expect(err).To(BeNil())
				})
			})

			Describe("CNAME Record", func() {
				It("Should properly serialize/deserialize", func() {
					cname := &ibclient.RecordCNAME{
						View:      "e2e_test_dns_view",
						Canonical: "e2e_test_cname_record.e2e-test.com",
						Name:      "e2e_test_cname_record.e2e-test.com",
						Ttl:       5,
						UseTtl:    true,
						Comment:   "CNAME Record created by e2e test",
						Ea:        make(ibclient.EA),
					}

					ref, err := connector.CreateObject(cname)
					Expect(err).To(BeNil())

					cname.SetReturnFields([]string{"name", "ttl", "view", "comment"})
					var res ibclient.RecordCNAME
					err = connector.GetObject(cname, ref, nil, &res)
					Expect(err).To(BeNil())
					Expect(res.Ref).To(Equal(ref))
					Expect(res.Name).To(Equal(cname.Name))
					Expect(res.Ttl).To(Equal(cname.Ttl))
					Expect(res.Comment).To(Equal(cname.Comment))

					cname.View = ""
					cname.Comment = "CNAME Record updated by e2e test"
					updRef, err := connector.UpdateObject(cname, ref)
					Expect(err).To(BeNil())

					_, err = connector.DeleteObject(updRef)
					Expect(err).To(BeNil())
				})
			})

			Describe("TXT Record", func() {
				It("Should properly serialize/deserialize", func() {
					txt := &ibclient.RecordTXT{
						View:    "e2e_test_dns_view",
						Name:    "e2e_test_txt_record.e2e-test.com",
						Text:    "TXT Record created by e2e test",
						Ttl:     5,
						UseTtl:  true,
						Comment: "TXT Record created by e2e test",
					}

					ref, err := connector.CreateObject(txt)
					Expect(err).To(BeNil())

					txt.SetReturnFields([]string{"view", "name", "text", "ttl", "comment"})
					var res ibclient.RecordTXT
					err = connector.GetObject(txt, ref, nil, &res)
					Expect(err).To(BeNil())
					Expect(res.Ref).To(Equal(ref))
					Expect(res.View).To(Equal(txt.View))
					Expect(res.Name).To(Equal(txt.Name))
					Expect(res.Text).To(Equal(txt.Text))
					Expect(res.Ttl).To(Equal(txt.Ttl))
					Expect(res.Comment).To(Equal(txt.Comment))

					txt.Comment = "TXT Record updated by e2e test"
					updRef, err := connector.UpdateObject(txt, ref)
					Expect(err).To(BeNil())

					_, err = connector.DeleteObject(updRef)
					Expect(err).To(BeNil())
				})
			})

			When("reverse-mapping DNS Zone Exists", func() {
				BeforeEach(func() {
					z := &ibclient.ZoneAuth{
						View:       "e2e_test_dns_view",
						Fqdn:       "192.168.1.0/24",
						ZoneFormat: "IPV4",
						Comment:    "Reverse-mapping DNS Auth Zone created by e2e test",
					}

					_, err := connector.CreateObject(z)
					Expect(err).To(BeNil())
				})

				Describe("PTR Record", func() {
					It("Should properly serialize/deserialize", func() {
						ptr := &ibclient.RecordPTR{
							View:     "e2e_test_dns_view",
							PtrdName: "e2e_test_ptr_record.e2e-test.com",
							Ipv4Addr: "192.168.1.45",
							Ttl:      5,
							UseTtl:   true,
							Comment:  "PTR Record created by e2e test",
						}

						ref, err := connector.CreateObject(ptr)
						Expect(err).To(BeNil())

						ptr.SetReturnFields([]string{"ptrdname", "ipv4addr", "ttl", "view", "comment"})
						var res ibclient.RecordPTR
						err = connector.GetObject(ptr, ref, nil, &res)
						Expect(err).To(BeNil())
						Expect(res.Ref).To(Equal(ref))
						Expect(res.View).To(Equal(ptr.View))
						Expect(res.PtrdName).To(Equal(ptr.PtrdName))
						Expect(res.Ipv4Addr).To(Equal(ptr.Ipv4Addr))
						Expect(res.Ttl).To(Equal(ptr.Ttl))
						Expect(res.Comment).To(Equal(ptr.Comment))

						ptr.View = ""
						ptr.Comment = "PTR Record updated by e2e test"
						updRef, err := connector.UpdateObject(ptr, ref)
						Expect(err).To(BeNil())

						_, err = connector.DeleteObject(updRef)
						Expect(err).To(BeNil())
					})

				})
			})
		})
	})

	Describe("EA Definition", func() {
		It("Should properly serialize/deserialize", func() {
			eadef := &ibclient.EADefinition{
				Name:       "E2E Test EA",
				Comment:    "EA Def created by e2e test",
				ListValues: []*ibclient.EADefListValue{{"value1"}, {"value2"}},
				Type:       "STRING",
			}

			ref, err := connector.CreateObject(eadef)
			Expect(err).To(BeNil())

			eadef.SetReturnFields(append(eadef.ReturnFields(), "list_values"))
			var res ibclient.EADefinition
			err = connector.GetObject(eadef, ref, nil, &res)
			Expect(err).To(BeNil())
			Expect(res.Ref).To(Equal(ref))
			Expect(res.Name).To(Equal(eadef.Name))
			Expect(res.Comment).To(Equal(eadef.Comment))
			Expect(res.ListValues).To(Equal(eadef.ListValues))
			Expect(res.Type).To(Equal(eadef.Type))

			eadef.Comment = "EA Def updated by e2e test"
			updRef, err := connector.UpdateObject(eadef, ref)
			Expect(err).To(BeNil())

			_, err = connector.DeleteObject(updRef)
			Expect(err).To(BeNil())
		})
	})

})
