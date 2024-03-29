package e2e_tests

import (
	ibclient "github.com/infobloxopen/infoblox-go-client/v2"
	"github.com/infobloxopen/infoblox-go-client/v2/utils"
	. "github.com/onsi/ginkgo/v2"
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
		It("Should properly serialize/deserialize", Label("RW"), func() {
			nv := &ibclient.NetworkView{
				Name:    utils.StringPtr("e2e_test_view"),
				Comment: utils.StringPtr("Network View created by e2e test"),
			}

			ref, err := connector.CreateObject(nv)
			Expect(err).To(BeNil())

			var res ibclient.NetworkView
			err = connector.GetObject(nv, ref, nil, &res)
			Expect(err).To(BeNil())
			Expect(res.Ref).To(Equal(ref))
			Expect(res.Name).To(Equal(nv.Name))
			Expect(res.Comment).To(Equal(nv.Comment))

			nv.Comment = utils.StringPtr("Network View updated by e2e test")
			updRef, err := connector.UpdateObject(nv, ref)
			Expect(err).To(BeNil())

			_, err = connector.DeleteObject(updRef)
			Expect(err).To(BeNil())
		})
	})

	When("Network View exists", Label("RW"), func() {
		BeforeEach(func() {
			nv := &ibclient.NetworkView{
				Name:    utils.StringPtr("e2e_test_view"),
				Comment: utils.StringPtr("Network View created by e2e test"),
			}

			_, err := connector.CreateObject(nv)
			Expect(err).To(BeNil())
		})

		Describe("IPv4 Network", Label("RW"), func() {
			It("Should properly serialize/deserialize", func() {
				nw := &ibclient.Ipv4Network{
					NetworkView: "e2e_test_view",
					Network:     utils.StringPtr("192.168.1.0/24"),
					Comment:     utils.StringPtr("IPv4 Network created by e2e test"),
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
				nw.Comment = utils.StringPtr("IPv4 Network updated by e2e test")
				updRef, err := connector.UpdateObject(nw, ref)
				Expect(err).To(BeNil())

				_, err = connector.DeleteObject(updRef)
				Expect(err).To(BeNil())
			})

			It("Should be able to remove all EAs", Label("RW", "IPv4 Network"), func() {
				nw := &ibclient.Ipv4Network{
					NetworkView: "e2e_test_view",
					Network:     utils.StringPtr("192.168.1.0/24"),
					Ea: ibclient.EA{
						"Country": "Colombia",
					},
				}

				ref, err := connector.CreateObject(nw)
				Expect(err).To(BeNil())
				nw.SetReturnFields([]string{"network_view", "network", "extattrs"})

				var res ibclient.Ipv4Network
				err = connector.GetObject(nw, ref, nil, &res)
				Expect(err).To(BeNil())
				Expect(res.Ref).To(Equal(ref))
				Expect(res.NetworkView).To(Equal(nw.NetworkView))
				Expect(res.Network).To(Equal(nw.Network))
				Expect(res.Ea["Country"]).To(Equal("Colombia"))

				nw.NetworkView = ""
				nw.Ea = ibclient.EA{}
				updRef, err := connector.UpdateObject(nw, ref)
				Expect(err).To(BeNil())

				res = ibclient.Ipv4Network{}
				err = connector.GetObject(nw, ref, nil, &res)
				Expect(err).To(BeNil())
				Expect(res.Ref).To(Equal(ref))
				Expect(res.Ea["Country"]).To(BeNil())

				_, err = connector.DeleteObject(updRef)
				Expect(err).To(BeNil())
			})
		})

		When("IPv4 Network exists", Label("RW"), func() {
			BeforeEach(func() {
				nw := &ibclient.Ipv4Network{
					NetworkView: "e2e_test_view",
					Network:     utils.StringPtr("192.168.1.0/24"),
					Comment:     utils.StringPtr("IPv4 Network created by e2e test"),
				}

				_, err := connector.CreateObject(nw)
				Expect(err).To(BeNil())
			})

			Describe("IPv4 Fixed Address", func() {
				It("Should properly serialize/deserialize", Label("RW"), func() {
					fa := &ibclient.Ipv4FixedAddress{
						NetworkView: utils.StringPtr("e2e_test_view"),
						Name:        utils.StringPtr("e2e_test_ipv4_fixed_address"),
						Ipv4Addr:    utils.StringPtr("192.168.1.45"),
						Mac:         utils.StringPtr("00:00:00:00:00:00"),
						Comment:     utils.StringPtr("IPv4 Fixed Address created by e2e test"),
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

					fa.Comment = utils.StringPtr("IPv4 Fixed Address updated by e2e test")
					updRef, err := connector.UpdateObject(fa, ref)
					Expect(err).To(BeNil())

					_, err = connector.DeleteObject(updRef)
					Expect(err).To(BeNil())
				})
			})

			Describe("IP Range", func() {
				It("Should properly serialize/deserialize", Label("RW"), func() {
					r := &ibclient.Range{
						NetworkView: utils.StringPtr("e2e_test_view"),
						Name:        utils.StringPtr("e2e_test_ip_range"),
						Comment:     utils.StringPtr("IP Range created by e2e test"),
						StartAddr:   utils.StringPtr("192.168.1.10"),
						EndAddr:     utils.StringPtr("192.168.1.20"),
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

					r.Comment = utils.StringPtr("IP Range updated by e2e test")
					updRef, err := connector.UpdateObject(r, ref)
					Expect(err).To(BeNil())

					_, err = connector.DeleteObject(updRef)
					Expect(err).To(BeNil())
				})
			})

			Describe("IPv4 Host Record", Label("record:host"), func() {
				It("Should properly serialize/deserialize", Label("RW"), func() {
					By("Creating DNS view")
					v := &ibclient.View{
						Name:        utils.StringPtr("e2e_test_dns_view"),
						NetworkView: utils.StringPtr("e2e_test_view"),
						Comment:     utils.StringPtr("DNS View created by e2e test"),
					}
					_, err := connector.CreateObject(v)
					Expect(err).To(BeNil())

					By("Creating forwarding-mapping DNS Auth Zone")
					zf := &ibclient.ZoneAuth{
						View:    utils.StringPtr("e2e_test_dns_view"),
						Fqdn:    "e2e-test.com",
						Comment: utils.StringPtr("Forwarding-mapping DNS Auth Zone created by e2e test"),
					}
					_, err = connector.CreateObject(zf)
					Expect(err).To(BeNil())

					By("Creating reverse-mapping DNS Auth Zone")
					zr := &ibclient.ZoneAuth{
						View:       utils.StringPtr("e2e_test_dns_view"),
						Fqdn:       "192.168.1.0/24",
						ZoneFormat: "IPV4",
						Comment:    utils.StringPtr("Reverse-mapping DNS Auth Zone created by e2e test"),
					}
					_, err = connector.CreateObject(zr)
					Expect(err).To(BeNil())

					fa := &ibclient.HostRecord{
						NetworkView: "e2e_test_view",
						View:        utils.StringPtr("e2e_test_dns_view"),
						Name:        utils.StringPtr("e2e_test_host_record.e2e-test.com"),
						Ipv4Addrs: []ibclient.HostRecordIpv4Addr{
							{
								Ipv4Addr: utils.StringPtr("192.168.1.60"),
								Mac:      utils.StringPtr("00:00:00:00:00:00"),
							},
						},
						Comment: utils.StringPtr("IPv4 Host Record created by e2e test"),
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
					fa.Comment = utils.StringPtr("IPv4 Host Record updated by e2e test")
					updRef, err := connector.UpdateObject(fa, ref)
					Expect(err).To(BeNil())

					_, err = connector.DeleteObject(updRef)
					Expect(err).To(BeNil())
				})

				It("View field should be updatable",
					Label("RW"), func() {
						By("Creating DNS view")
						v := &ibclient.View{
							Name:        utils.StringPtr("e2e_test_dns_view"),
							NetworkView: utils.StringPtr("e2e_test_view"),
							Comment:     utils.StringPtr("DNS View created by e2e test"),
						}
						_, err := connector.CreateObject(v)
						Expect(err).To(BeNil())

						By("Creating forwarding-mapping DNS Auth Zone")
						zf := &ibclient.ZoneAuth{
							View:    utils.StringPtr("e2e_test_dns_view"),
							Fqdn:    "e2e-test.com",
							Comment: utils.StringPtr("Forwarding-mapping DNS Auth Zone created by e2e test"),
						}
						_, err = connector.CreateObject(zf)
						Expect(err).To(BeNil())

						By("Creating reverse-mapping DNS Auth Zone")
						zr := &ibclient.ZoneAuth{
							View:       utils.StringPtr("e2e_test_dns_view"),
							Fqdn:       "192.168.1.0/24",
							ZoneFormat: "IPV4",
							Comment:    utils.StringPtr("Reverse-mapping DNS Auth Zone created by e2e test"),
						}
						_, err = connector.CreateObject(zr)
						Expect(err).To(BeNil())

						hr := &ibclient.HostRecord{
							NetworkView: "e2e_test_view",
							View:        utils.StringPtr("e2e_test_dns_view"),
							Name:        utils.StringPtr("e2e_test_host_record.e2e-test.com"),
							Ipv4Addrs: []ibclient.HostRecordIpv4Addr{
								{
									Ipv4Addr: utils.StringPtr("192.168.1.60"),
									Mac:      utils.StringPtr("00:00:00:00:00:00"),
								},
							},
							Comment: utils.StringPtr("IPv4 Host Record created by e2e test"),
						}

						By("Creating Host Record")
						ref, err := connector.CreateObject(hr)
						Expect(err).To(BeNil())

						hr.SetReturnFields([]string{"view"})
						var res ibclient.HostRecord
						err = connector.GetObject(hr, ref, nil, &res)
						Expect(err).To(BeNil())
						Expect(*res.View).To(Equal("e2e_test_dns_view"))

						By("Creating a second DNS view")
						v2 := &ibclient.View{
							Name: utils.StringPtr("e2e_test_dns_view2"),
						}
						_, err = connector.CreateObject(v2)
						Expect(err).To(BeNil())

						By("Creating a second forwarding-mapping DNS Auth Zone")
						zf2 := &ibclient.ZoneAuth{
							View:    utils.StringPtr("e2e_test_dns_view2"),
							Fqdn:    "e2e-test.com",
							Comment: utils.StringPtr("Forwarding-mapping DNS Auth Zone created by e2e test"),
						}
						_, err = connector.CreateObject(zf2)
						Expect(err).To(BeNil())

						By("Creating a second reverse-mapping DNS Auth Zone")
						zr2 := &ibclient.ZoneAuth{
							View:       utils.StringPtr("e2e_test_dns_view2"),
							Fqdn:       "192.168.1.0/24",
							ZoneFormat: "IPV4",
							Comment:    utils.StringPtr("Reverse-mapping DNS Auth Zone created by e2e test"),
						}
						_, err = connector.CreateObject(zr2)
						Expect(err).To(BeNil())

						By("Updating a DNS View value for the Host Record object")
						hr.NetworkView = ""
						hr.View = utils.StringPtr("e2e_test_dns_view2")
						updRef, err := connector.UpdateObject(hr, ref)
						Expect(err).To(BeNil())

						By("Reading Host Record and checking if view field is updated")
						res = ibclient.HostRecord{}
						err = connector.GetObject(hr, updRef, nil, &res)
						Expect(err).To(BeNil())
						Expect(*res.View).To(Equal("e2e_test_dns_view2"))
					})
			})

		})

		Describe("IPv4 Network Container", func() {
			It("Should properly serialize/deserialize", Label("RW"), func() {
				nc := &ibclient.Ipv4NetworkContainer{
					NetworkView: "e2e_test_view",
					Network:     "192.168.1.0/24",
					Comment:     utils.StringPtr("IPv4 Network Container created by e2e test"),
				}

				ref, err := connector.CreateObject(nc)
				Expect(err).To(BeNil())

				var res ibclient.Ipv4FixedAddress
				err = connector.GetObject(nc, ref, nil, &res)
				Expect(err).To(BeNil())
				Expect(res.Ref).To(Equal(ref))
				Expect(*res.NetworkView).To(Equal(nc.NetworkView))
				Expect(*res.Network).To(Equal(nc.Network))
				Expect(*res.Comment).To(Equal(*nc.Comment))

				nc.Network = ""
				nc.NetworkView = ""
				nc.Comment = utils.StringPtr("IPv4 Network Container updated by e2e test")
				updRef, err := connector.UpdateObject(nc, ref)
				Expect(err).To(BeNil())

				_, err = connector.DeleteObject(updRef)
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("DNS View", func() {
		It("Should properly serialize/deserialize", Label("RW"), func() {
			v := &ibclient.View{
				Name:    utils.StringPtr("e2e_test_dns_view"),
				Comment: utils.StringPtr("DNS View created by e2e test"),
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

			v.Comment = utils.StringPtr("DNS View updated by e2e test")
			updRef, err := connector.UpdateObject(v, ref)
			Expect(err).To(BeNil())

			_, err = connector.DeleteObject(updRef)
			Expect(err).To(BeNil())
		})

		It("Should update view comment to empty string", Label("RW", "DNS View"), func() {
			v := &ibclient.View{
				Name:    utils.StringPtr("e2e_test_dns_view"),
				Comment: utils.StringPtr("DNS View created by e2e test"),
			}

			ref, err := connector.CreateObject(v)
			Expect(err).To(BeNil())

			v.SetReturnFields([]string{"name", "comment"})

			var res ibclient.View
			err = connector.GetObject(v, ref, nil, &res)
			Expect(err).To(BeNil())
			Expect(res.Ref).To(Equal(ref))
			Expect(*res.Comment).To(Equal("DNS View created by e2e test"))

			v.Comment = utils.StringPtr("")
			updRef, err := connector.UpdateObject(v, ref)
			Expect(err).To(BeNil())

			v.SetReturnFields([]string{"name", "comment"})

			res = ibclient.View{}
			err = connector.GetObject(v, ref, nil, &res)
			Expect(err).To(BeNil())
			Expect(res.Ref).To(Equal(ref))
			Expect(res.Comment).To(BeNil())

			_, err = connector.DeleteObject(updRef)
			Expect(err).To(BeNil())
		})
	})

	When("DNS View exists", Label("RW"), func() {
		BeforeEach(func() {
			v := &ibclient.View{
				Name:    utils.StringPtr("e2e_test_dns_view"),
				Comment: utils.StringPtr("DNS View created by e2e test"),
			}

			_, err := connector.CreateObject(v)
			Expect(err).To(BeNil())
		})

		Describe("DNS Zone Auth", func() {
			It("Should support CRUD operations of forwarding-mapping zone", Label("RW"), func() {
				z := &ibclient.ZoneAuth{
					View:    utils.StringPtr("e2e_test_dns_view"),
					Fqdn:    "e2e-test.com",
					Comment: utils.StringPtr("DNS Auth Zone created by e2e test"),
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
				z.Comment = utils.StringPtr("DNS Auth Zone updated by e2e test")
				updRef, err := connector.UpdateObject(z, ref)
				Expect(err).To(BeNil())

				_, err = connector.DeleteObject(updRef)
				Expect(err).To(BeNil())
			})

			It("Should support CRUD operations of reverse-mapping zone", Label("RW"), func() {
				z := &ibclient.ZoneAuth{
					View:       utils.StringPtr("e2e_test_dns_view"),
					Fqdn:       "192.168.1.0/24",
					ZoneFormat: "IPV4",
					Comment:    utils.StringPtr("DNS Auth Zone created by e2e test"),
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
				z.Comment = utils.StringPtr("DNS Auth Zone updated by e2e test")
				updRef, err := connector.UpdateObject(z, ref)
				Expect(err).To(BeNil())

				_, err = connector.DeleteObject(updRef)
				Expect(err).To(BeNil())
			})
		})

		When("forwarding-mapping DNS Zone Exists", Label("RW"), func() {
			BeforeEach(func() {
				z := &ibclient.ZoneAuth{
					View:    utils.StringPtr("e2e_test_dns_view"),
					Fqdn:    "e2e-test.com",
					Comment: utils.StringPtr("Forwarding-mapping DNS Auth Zone created by e2e test"),
				}

				_, err := connector.CreateObject(z)
				Expect(err).To(BeNil())
			})

			Describe("A Record", func() {
				It("Should properly serialize/deserialize", Label("RW"), func() {
					a := &ibclient.RecordA{
						View:     "e2e_test_dns_view",
						Name:     utils.StringPtr("e2e_test_a_record.e2e-test.com"),
						Ipv4Addr: utils.StringPtr("192.168.1.45"),
						Ttl:      utils.Uint32Ptr(5),
						UseTtl:   utils.BoolPtr(true),
						Comment:  utils.StringPtr("A Record created by e2e test"),
						Ea:       ibclient.EA{},
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
					a.Comment = utils.StringPtr("A Record updated by e2e test")
					updRef, err := connector.UpdateObject(a, ref)
					Expect(err).To(BeNil())

					_, err = connector.DeleteObject(updRef)
					Expect(err).To(BeNil())
				})

				It("Should support search by zone field", Label("RW"), func() {
					a := &ibclient.RecordA{
						View:     "e2e_test_dns_view",
						Name:     utils.StringPtr("e2e_test_a_record.e2e-test.com"),
						Ipv4Addr: utils.StringPtr("192.168.1.45"),
						Ttl:      utils.Uint32Ptr(5),
						UseTtl:   utils.BoolPtr(true),
						Comment:  utils.StringPtr("A Record created by e2e test"),
						Ea:       ibclient.EA{},
					}

					ref, err := connector.CreateObject(a)
					Expect(err).To(BeNil())

					aSearch := &ibclient.RecordA{}
					aSearch.SetReturnFields([]string{"view", "comment", "creation_time", "zone"})
					var res []ibclient.RecordA
					queryParams := ibclient.NewQueryParams(false, map[string]string{"view": "e2e_test_dns_view", "zone": "e2e-test.com"})
					err = connector.GetObject(aSearch, "", queryParams, &res)
					Expect(err).To(BeNil())
					Expect(res[0].Ref).To(Equal(ref))
					Expect(res[0].View).To(Equal(a.View))
					Expect(res[0].Comment).To(Equal(a.Comment))
					Expect(res[0].Zone).To(Equal("e2e-test.com"))
					Expect(res[0].CreationTime).To(Not(BeNil()))

					_, err = connector.DeleteObject(ref)
					Expect(err).To(BeNil())
				})
			})

			Describe("AAAA Record", func() {
				It("Should properly serialize/deserialize", Label("RW"), func() {
					aaaa := &ibclient.RecordAAAA{
						View:     "e2e_test_dns_view",
						Name:     utils.StringPtr("e2e_test_a_record.e2e-test.com"),
						Ipv6Addr: utils.StringPtr("2001:db8:abcd:14::1"),
						Ttl:      utils.Uint32Ptr(5),
						UseTtl:   utils.BoolPtr(true),
						Comment:  utils.StringPtr("A Record created by e2e test"),
						Ea:       ibclient.EA{},
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
					aaaa.Comment = utils.StringPtr("A Record updated by e2e test")
					updRef, err := connector.UpdateObject(aaaa, ref)
					Expect(err).To(BeNil())

					_, err = connector.DeleteObject(updRef)
					Expect(err).To(BeNil())
				})
			})

			Describe("CNAME Record", Label("record:cname"), func() {
				It("Should properly serialize/deserialize", Label("RW"), func() {
					cname := &ibclient.RecordCNAME{
						View:      utils.StringPtr("e2e_test_dns_view"),
						Canonical: utils.StringPtr("e2e_test_cname_record.e2e-test.com"),
						Name:      utils.StringPtr("e2e_test_cname_record.e2e-test.com"),
						Ttl:       utils.Uint32Ptr(5),
						UseTtl:    utils.BoolPtr(true),
						Comment:   utils.StringPtr("CNAME Record created by e2e test"),
						Ea:        ibclient.EA{},
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

					cname.View = nil
					cname.Comment = utils.StringPtr("CNAME Record updated by e2e test")
					updRef, err := connector.UpdateObject(cname, ref)
					Expect(err).To(BeNil())

					_, err = connector.DeleteObject(updRef)
					Expect(err).To(BeNil())
				})

				It("View field should be updatable",
					Label("RW"), func() {
						By("Creating a second DNS view")
						v := &ibclient.View{
							Name: utils.StringPtr("e2e_test_dns_view2"),
						}
						_, err := connector.CreateObject(v)
						Expect(err).To(BeNil())

						By("Creating a second forward-mapping DNS zone")
						z := &ibclient.ZoneAuth{
							View: utils.StringPtr("e2e_test_dns_view2"),
							Fqdn: "e2e-test.com",
						}
						_, err = connector.CreateObject(z)
						Expect(err).To(BeNil())

						By("Creating CNAME Record in the first DNS View")
						cname := &ibclient.RecordCNAME{
							View:      utils.StringPtr("e2e_test_dns_view"),
							Canonical: utils.StringPtr("e2e_test_cname_record.e2e-test.com"),
							Name:      utils.StringPtr("e2e_test_cname_record.e2e-test.com"),
							Ttl:       utils.Uint32Ptr(5),
							UseTtl:    utils.BoolPtr(true),
							Comment:   utils.StringPtr("CNAME Record created by e2e test"),
							Ea:        ibclient.EA{},
						}

						ref, err := connector.CreateObject(cname)
						Expect(err).To(BeNil())

						cname.SetReturnFields([]string{"view"})
						var res ibclient.RecordCNAME
						err = connector.GetObject(cname, ref, nil, &res)
						Expect(err).To(BeNil())
						Expect(*res.View).To(Equal("e2e_test_dns_view"))

						By("Updating CNAME record's view field with a new view name")
						cname.View = utils.StringPtr("e2e_test_dns_view2")
						updRef, err := connector.UpdateObject(cname, ref)
						Expect(err).To(BeNil())

						By("Reading the same CNAME record and checking if view field is updated")
						res = ibclient.RecordCNAME{}
						err = connector.GetObject(cname, updRef, nil, &res)
						Expect(err).To(BeNil())
						Expect(*res.View).To(Equal("e2e_test_dns_view2"))

						_, err = connector.DeleteObject(updRef)
						Expect(err).To(BeNil())
					})
			})

			Describe("TXT Record", func() {
				It("Should properly serialize/deserialize", Label("RW"), func() {
					txt := &ibclient.RecordTXT{
						View:    utils.StringPtr("e2e_test_dns_view"),
						Name:    utils.StringPtr("e2e_test_txt_record.e2e-test.com"),
						Text:    utils.StringPtr("TXT Record created by e2e test"),
						Ttl:     utils.Uint32Ptr(5),
						UseTtl:  utils.BoolPtr(true),
						Comment: utils.StringPtr("TXT Record created by e2e test"),
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

					txt.Comment = utils.StringPtr("TXT Record updated by e2e test")
					updRef, err := connector.UpdateObject(txt, ref)
					Expect(err).To(BeNil())

					_, err = connector.DeleteObject(updRef)
					Expect(err).To(BeNil())
				})
			})

			When("reverse-mapping DNS Zone Exists", Label("RW"), func() {
				BeforeEach(func() {
					z := &ibclient.ZoneAuth{
						View:       utils.StringPtr("e2e_test_dns_view"),
						Fqdn:       "192.168.1.0/24",
						ZoneFormat: "IPV4",
						Comment:    utils.StringPtr("Reverse-mapping DNS Auth Zone created by e2e test"),
					}

					_, err := connector.CreateObject(z)
					Expect(err).To(BeNil())
				})

				Describe("PTR Record", func() {
					It("Should properly serialize/deserialize", Label("RW"), func() {
						ptr := &ibclient.RecordPTR{
							View:     "e2e_test_dns_view",
							PtrdName: utils.StringPtr("e2e_test_ptr_record.e2e-test.com"),
							Ipv4Addr: utils.StringPtr("192.168.1.45"),
							Ttl:      utils.Uint32Ptr(5),
							UseTtl:   utils.BoolPtr(true),
							Comment:  utils.StringPtr("PTR Record created by e2e test"),
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
						ptr.Comment = utils.StringPtr("PTR Record updated by e2e test")
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
		It("Should properly serialize/deserialize", Label("RW"), func() {
			eadef := &ibclient.EADefinition{
				Name:       utils.StringPtr("E2E Test EA"),
				Comment:    utils.StringPtr("EA Def created by e2e test"),
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

			eadef.Comment = utils.StringPtr("EA Def updated by e2e test")
			updRef, err := connector.UpdateObject(eadef, ref)
			Expect(err).To(BeNil())

			_, err = connector.DeleteObject(updRef)
			Expect(err).To(BeNil())
		})
	})

})
