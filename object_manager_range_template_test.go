package ibclient

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Object Manager: Range Template", func() {
	Describe("Create Range Template with minimum params", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		comment := "test range template"
		name := "RangeTemplate1"
		numberOfAddresses := uint32(10)
		offset := uint32(40)
		ea := EA{"Site": "Hokkaido"}
		options := []*Dhcpoption{
			{
				Name:        "routers",
				Num:         3,
				Value:       "12.4.0.23",
				VendorClass: "DHCP",
				UseOption:   true,
			},
		}
		useOption := true
		failOverAssociation := ""
		serverAssociationType := "MEMBER"
		member := &Dhcpmember{
			Ipv4Addr: "10.17.21.10",
			Ipv6Addr: "2403:8600:80cf:e10c:3a00::1192",
			Name:     "infoblox.localdomain",
		}
		fakeRefReturn := fmt.Sprintf("rangetemplate/ZG5zLmhvc3QkLZhd3QuaDE:%s", name)

		conn := &fakeConnector{
			createObjectObj:      NewRangeTemplate("", name, numberOfAddresses, offset, comment, ea, options, useOption, serverAssociationType, failOverAssociation, member),
			getObjectObj:         &Rangetemplate{},
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         NewRangeTemplate("", name, numberOfAddresses, offset, comment, ea, options, useOption, serverAssociationType, failOverAssociation, member),
			fakeRefReturn:        fakeRefReturn,
		}
		conn.resultObject.(*Rangetemplate).Ref = fakeRefReturn
		objMgr := NewObjectManager(conn, cmpType, tenantID)
		It("should pass expected Range Template Object to CreateObject", func() {
			actualRecord, err := objMgr.CreateRangeTemplate(name, numberOfAddresses, offset, comment, ea, options, useOption, serverAssociationType, failOverAssociation, member)
			Expect(actualRecord).To(Equal(conn.resultObject))
			Expect(err).To(BeNil())
		})

		It("should fail to create a Range Template object", func() {
			actualRecord, err := objMgr.CreateRangeTemplate("", numberOfAddresses, offset, comment, ea, options, useOption, serverAssociationType, failOverAssociation, member)
			Expect(actualRecord).To(BeNil())
			Expect(err).ToNot(BeNil())
		})
	})

	Describe("Get Range Template", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		comment := "test range template"
		name := "RangeTemplate1"
		numberOfAddresses := uint32(10)
		offset := uint32(40)
		ea := EA{"Site": "Hokkaido"}
		options := []*Dhcpoption{
			{
				Name:        "routers",
				Num:         3,
				Value:       "12.4.0.23",
				VendorClass: "DHCP",
				UseOption:   true,
			},
		}
		useOption := true
		failOverAssociation := ""
		serverAssociationType := "MEMBER"
		member := &Dhcpmember{
			Ipv4Addr: "10.17.21.10",
			Ipv6Addr: "2403:8600:80cf:e10c:3a00::1192",
			Name:     "infoblox.localdomain",
		}
		fakeRefReturn := fmt.Sprintf("rangetemplate/ZG5zLmhvc3QkLZhd3QuaDE:%s", name)
		queryParams := NewQueryParams(false, map[string]string{"name": name})
		res := NewRangeTemplate("", name, numberOfAddresses, offset, comment, ea, options, useOption, serverAssociationType, failOverAssociation, member)

		conn := &fakeConnector{
			getObjectObj:  NewEmptyRangeTemplate(),
			resultObject:  []Rangetemplate{*res},
			fakeRefReturn: fakeRefReturn,
		}
		objMgr := NewObjectManager(conn, cmpType, tenantID)
		It("should get expected Range Template Object from getObject", func() {
			conn.getObjectQueryParams = queryParams
			actualRecord, err := objMgr.GetAllRangeTemplate(queryParams)
			Expect(actualRecord).To(Equal(conn.resultObject))
			Expect(err).To(BeNil())
		})

		It("should fail to get expected Range Template Object from getObject", func() {
			queryParams1 := NewQueryParams(false, map[string]string{"name": "range-template123"})
			conn.getObjectQueryParams = queryParams1
			conn.resultObject = []Rangetemplate{}
			actualRecord, err := objMgr.GetAllRangeTemplate(queryParams1)
			Expect(actualRecord).To(Equal(conn.resultObject))
			Expect(err).To(BeNil())
		})

	})

	Describe("Get Range Template: Negative scenario", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		queryParams2 := NewQueryParams(false, map[string]string{"number_of_addresses": "24"})
		conn := &fakeConnector{
			getObjectObj:         NewEmptyRangeTemplate(),
			getObjectQueryParams: queryParams2,
			resultObject:         []Rangetemplate{},
		}
		objMgr := NewObjectManager(conn, cmpType, tenantID)

		// negative scenario
		conn.getObjectError = fmt.Errorf("Field is not searchable: number_of_addresses")
		It("should fail to get expected Range Template Object from getObject with non searchable field", func() {
			_, err := objMgr.GetAllRangeTemplate(queryParams2)
			Expect(err).ToNot(BeNil())
		})

	})

	Describe("Delete Range Template", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		recordName := "range-template4556"
		deleteRef := fmt.Sprintf("rangetemplate/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)
		nwFakeConnector := &fakeConnector{
			deleteObjectRef: deleteRef,
			fakeRefReturn:   deleteRef,
		}

		objMgr := NewObjectManager(nwFakeConnector, cmpType, tenantID)

		var actualRef string
		var err error
		It("should pass expected Range Template Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteRangeTemplate(deleteRef)
		})
		It("should return expected Range Template Ref", func() {
			Expect(actualRef).To(Equal(deleteRef))
			Expect(err).To(BeNil())
		})

		It("should pass expected Range Template Ref to DeleteObject", func() {
			deleteRef2 := "test-template777"
			nwFakeConnector.deleteObjectRef = deleteRef2
			nwFakeConnector.fakeRefReturn = ""
			nwFakeConnector.deleteObjectError = fmt.Errorf("not found")
			actualRef, err = objMgr.DeleteRangeTemplate(deleteRef2)
		})

		It("should return an error", func() {
			Expect(err).ToNot(BeNil())
		})

	})

	Describe("Update Range Template", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		name := "range-template-new"
		comment := "test range template updated"
		numberOfAddresses := uint32(45)
		offset := uint32(70)
		ea := EA{"Site": "Fukuoka"}
		options := []*Dhcpoption{
			{
				Name:        "routers",
				Num:         3,
				Value:       "12.4.0.23",
				VendorClass: "DHCP",
				UseOption:   true,
			},
		}
		useOption := true
		failOverAssociation := ""
		serverAssociationType := "MEMBER"
		member := &Dhcpmember{
			Ipv4Addr: "10.17.21.10",
			Ipv6Addr: "2403:8600:80cf:e10c:3a00::1192",
			Name:     "infoblox.localdomain",
		}
		updateRef := fmt.Sprintf("rangetemplate/ZG5zLmhvc3QkLZhd3QuaDE:%s", name)

		conn := &fakeConnector{
			getObjectObj:         NewEmptyRangeTemplate(),
			getObjectRef:         updateRef,
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         NewRangeTemplate("", name, numberOfAddresses, offset, comment, ea, options, useOption, serverAssociationType, failOverAssociation, member),
			fakeRefReturn:        updateRef,
			updateObjectObj:      NewRangeTemplate(updateRef, name, numberOfAddresses, offset, comment, ea, options, useOption, serverAssociationType, failOverAssociation, member),
			updateObjectRef:      updateRef,
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)
		It("should pass expected Range Template Object to UpdateObject", func() {
			actualRecord, err := objMgr.UpdateRangeTemplate(updateRef, name, numberOfAddresses, offset, comment, ea, options, useOption, serverAssociationType, failOverAssociation, member)
			Expect(actualRecord).To(Equal(conn.resultObject))
			Expect(err).To(BeNil())
		})

	})

	Describe("Update Range Template with, negative scenario", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		name := "range-template-new"
		comment := "test range template updated"
		numberOfAddresses := uint32(45)
		offset := uint32(70)
		ea := EA{"Site": "Fukuoka"}
		options := []*Dhcpoption{
			{
				Name:        "routers",
				Num:         3,
				Value:       "12.4.0.23",
				VendorClass: "DHCP",
				UseOption:   true,
			},
		}
		useOption := true
		failOverAssociation := ""
		serverAssociationType := "MEMBER"
		member := &Dhcpmember{
			Ipv4Addr: "10.17.21.10",
			Ipv6Addr: "2403:8600:80cf:e10c:3a00::1192",
			Name:     "infoblox.localdomain",
		}
		oldRef := "rangetemplate/ZG5zLmhvc3QkLZhd3QuaDE:range-template-new"

		conn := &fakeConnector{
			getObjectObj:         NewEmptyRangeTemplate(),
			getObjectRef:         oldRef,
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         NewRangeTemplate(oldRef, name, numberOfAddresses, offset, comment, ea, options, useOption, serverAssociationType, failOverAssociation, member),
			getObjectError:       fmt.Errorf("not found"),
			fakeRefReturn:        oldRef,
			updateObjectObj:      NewRangeTemplate(oldRef, name, numberOfAddresses, offset, comment, ea, options, useOption, serverAssociationType, failOverAssociation, member),
			updateObjectRef:      oldRef,
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)
		// negative scenario

		It("should fail to update Range Template Object", func() {
			actualRecord, err := objMgr.UpdateRangeTemplate(oldRef, name, numberOfAddresses, offset, comment, ea, options, useOption, serverAssociationType, failOverAssociation, member)
			Expect(actualRecord).To(BeNil())
			Expect(err).ToNot(BeNil())
		})

	})

})
