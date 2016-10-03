package ibclient

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"encoding/json"
)

var _ = Describe("Objects", func() {

	Context("EA Object", func() {

		ea := EA{
			"Cloud API Owned":   Bool(true),
			"Tenant Name":       "Engineering01",
			"Maximum Wait Time": 120,
			"DNS Support":       Bool(false),
		}
		eaJSON := `{"Cloud API Owned":{"value":"True"},` +
			`"Tenant Name":{"value":"Engineering01"},` +
			`"Maximum Wait Time":{"value":120},` +
			`"DNS Support":{"value":"False"}}`

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
			nv := NewNetworkView(NetworkView{Name: name})

			It("should set fields correctly", func() {
				Expect(nv.Name).To(Equal(name))
			})

			It("should set base fields correctly", func() {
				Expect(nv.ObjectType()).To(Equal("networkview"))
				Expect(nv.ReturnFields()).To(ConsistOf("extattrs", "name"))
			})
		})

		Context("Network object", func() {
			cidr := "123.0.0.0/24"
			netviewName := "localview"
			nw := NewNetwork(Network{Cidr: cidr, NetviewName: netviewName})
			searchEAs := EA{"Network Name": "shared-net"}
			nw.eaSearch = searchEAs

			It("should set fields correctly", func() {
				Expect(nw.Cidr).To(Equal(cidr))
				Expect(nw.NetviewName).To(Equal(netviewName))
			})

			It("should set base fields correctly", func() {
				Expect(nw.ObjectType()).To(Equal("network"))
				Expect(nw.ReturnFields()).To(ConsistOf("extattrs", "network", "network_view"))
				Expect(nw.EaSearch()).To(Equal(searchEAs))
			})
		})

		Context("NetworkContainer object", func() {
			cidr := "74.0.8.0/24"
			netviewName := "globalview"
			nwc := NewNetworkContainer(NetworkContainer{Cidr: cidr, NetviewName: netviewName})

			It("should set fields correctly", func() {
				Expect(nwc.Cidr).To(Equal(cidr))
				Expect(nwc.NetviewName).To(Equal(netviewName))
			})

			It("should set base fields correctly", func() {
				Expect(nwc.ObjectType()).To(Equal("networkcontainer"))
				Expect(nwc.ReturnFields()).To(ConsistOf("extattrs", "network", "network_view"))
			})
		})

		Context("FixedAddress object", func() {
			netviewName := "globalview"
			cidr := "25.0.7.0/24"
			ipAddress := "25.0.7.59/24"
			mac := "11:22:33:44:55:66"
			fixedAddr := NewFixedAddress(FixedAddress{
				NetviewName: netviewName,
				Cidr:        cidr,
				IPAddress:   ipAddress,
				Mac:         mac})

			It("should set fields correctly", func() {
				Expect(fixedAddr.NetviewName).To(Equal(netviewName))
				Expect(fixedAddr.Cidr).To(Equal(cidr))
				Expect(fixedAddr.IPAddress).To(Equal(ipAddress))
				Expect(fixedAddr.Mac).To(Equal(mac))
			})

			It("should set base fields correctly", func() {
				Expect(fixedAddr.ObjectType()).To(Equal("fixedaddress"))
				Expect(fixedAddr.ReturnFields()).To(ConsistOf("network_view", "network", "ipv4addr", "mac", "extattrs"))
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

	})

	Context("Umnarshalling malformed JSON", func() {
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
