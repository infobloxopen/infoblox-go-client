package ibclient

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Object Manager: apps", func() {

	Describe("Get DNS member by reference", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		fakeRefReturn := "memmber:dns/ZG5zLm5ldHdvcmtfdmlldyQyMw:infoblox.localdomain"
		resObj := []Dns{*NewDns(Dns{
			Ref:      fakeRefReturn,
			HostName: "infoblox.localdomain",
		})}
		ncFakeConnector := &fakeConnector{
			getObjectObj:         NewDns(Dns{}),
			getObjectRef:         fakeRefReturn,
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         resObj,
		}

		objMgr := NewObjectManager(ncFakeConnector, cmpType, tenantID)

		var actualDns []Dns
		var err error
		It("should pass expected Dns Object to GetObject", func() {
			actualDns, err = objMgr.GetDnsMember(fakeRefReturn)
		})
		It("should return expected Dns Object", func() {
			Expect(err).To(BeNil())
			Expect(actualDns[0]).To(Equal(ncFakeConnector.resultObject.([]Dns)[0]))
		})
	})

	Describe("Get DNS member by reference", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		fakeRefReturn := "memmber:dns/ZG5zLm5ldHdvcmtfdmlldyQyMw:infoblox.localdomain"
		resObj := []Dns{*NewDns(Dns{
			Ref:      fakeRefReturn,
			HostName: "infoblox.localdomain",
		})}
		ncFakeConnector := &fakeConnector{
			getObjectObj:         NewDns(Dns{}),
			getObjectRef:         fakeRefReturn,
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         resObj,
		}

		objMgr := NewObjectManager(ncFakeConnector, cmpType, tenantID)

		var actualDns []Dns
		var err error
		It("should pass expected Dns Object to GetObject", func() {
			actualDns, err = objMgr.GetDnsMember(fakeRefReturn)
		})
		It("should return expected Dns Object", func() {
			Expect(err).To(BeNil())
			Expect(actualDns[0]).To(Equal(ncFakeConnector.resultObject.([]Dns)[0]))
		})
	})

	Describe("Get DHCP member by reference", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		fakeRefReturn := "memmber:dhcpproperties/ZG5zLm5ldHdvcmtfdmlldyQyMw:infoblox.localdomain"
		resObj := []Dhcp{*NewDhcp(Dhcp{
			Ref:      fakeRefReturn,
			HostName: "infoblox.localdomain",
		})}
		ncFakeConnector := &fakeConnector{
			getObjectObj:         NewDhcp(Dhcp{}),
			getObjectRef:         fakeRefReturn,
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         resObj,
		}

		objMgr := NewObjectManager(ncFakeConnector, cmpType, tenantID)

		var actualDhcp []Dhcp
		var err error
		It("should pass expected Dhcp Object to GetObject", func() {
			actualDhcp, err = objMgr.GetDhcpMember(fakeRefReturn)
		})
		It("should return expected Dhcp Object", func() {
			Expect(err).To(BeNil())
			Expect(actualDhcp[0]).To(Equal(ncFakeConnector.resultObject.([]Dhcp)[0]))
		})
	})

	Describe("Update DNS member by reference", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		fakeRefReturn := "updateref"
		resObj := Dns{
			IBBase: IBBase{
				returnFields: []string{"enable_dns", "host_name"},
				eaSearch:     nil,
			},
			EnableDns:  false,
			objectType: "member:dns",
		}
		ncFakeConnector := &fakeConnector{
			updateObjectObj:      NewDns(Dns{}),
			updateObjectRef:      fakeRefReturn,
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         resObj,
			updateObjectError:    nil,
		}

		objMgr := NewObjectManager(ncFakeConnector, cmpType, tenantID)

		var actualDns Dns
		var err error
		It("should pass expected Dns Object to GetObject", func() {
			actualDns, err = objMgr.UpdateDnsStatus(fakeRefReturn, false)
		})
		It("should return expected Dns Object", func() {
			Expect(err).To(BeNil())
			Expect(actualDns).To(Equal(resObj))
		})
	})

	Describe("Update DHCP member by reference", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		fakeRefReturn := "updateref"
		resObj := Dhcp{
			IBBase: IBBase{
				returnFields: []string{"enable_dhcp", "host_name"},
				eaSearch:     nil,
			},
			EnableDhcp: false,
			objectType: "member:dhcpproperties",
		}
		ncFakeConnector := &fakeConnector{
			updateObjectObj:      NewDhcp(Dhcp{}),
			updateObjectRef:      fakeRefReturn,
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         resObj,
			updateObjectError:    nil,
		}

		objMgr := NewObjectManager(ncFakeConnector, cmpType, tenantID)

		var actualDhcp Dhcp
		var err error
		It("should pass expected Dhcp Object to GetObject", func() {
			actualDhcp, err = objMgr.UpdateDhcpStatus(fakeRefReturn, false)
		})
		It("should return expected Dhcp Object", func() {
			Expect(err).To(BeNil())
			Expect(actualDhcp).To(Equal(resObj))
		})
	})

})
