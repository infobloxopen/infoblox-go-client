package ibclient

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Object Manager: host record", func() {
	Describe("Allocate next available host Record without dns", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "private"
		ipv4Cidr := "53.0.0.0/24"
		macAddr := "01:23:45:67:80:ab"
		ipv4Addr := fmt.Sprintf("func:nextavailableip:%s,%s", ipv4Cidr, netviewName)
		ipv6Cidr := "2003:db8:abcd:14::/64"
		duid := "02:24:46:68:81:cd"
		ipv6Addr := fmt.Sprintf("func:nextavailableip:%s,%s", ipv6Cidr, netviewName)
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		recordName := "test"
		enabledns := false
		enabledhcp := false
		dnsView := "default"
		fakeRefReturn := fmt.Sprintf("record:host/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)
		resultIPV4Addrs := NewHostRecordIpv4Addr(ipv4Addr, macAddr, enabledhcp, "")
		resultIPv6Addrs := NewHostRecordIpv6Addr(ipv6Addr, duid, enabledhcp, "")
		useTtl := true
		ttl := uint32(70)
		comment := "test"
		aliases := []string{"abc.test.com"}

		eas := make(EA)
		eas["VM ID"] = vmID
		eas["VM Name"] = vmName

		aniFakeConnector := &fakeConnector{
			createObjectObj: NewHostRecord(
				netviewName, recordName,
				"", "", []HostRecordIpv4Addr{*resultIPV4Addrs}, []HostRecordIpv6Addr{*resultIPv6Addrs},
				eas, enabledns, dnsView, "", "", useTtl, ttl, comment, aliases),
			getObjectRef: fakeRefReturn,
			getObjectObj: NewHostRecord(
				netviewName, recordName,
				"", "", []HostRecordIpv4Addr{*resultIPV4Addrs}, []HostRecordIpv6Addr{*resultIPv6Addrs},
				eas, enabledns, dnsView, "", fakeRefReturn, useTtl, ttl, comment, aliases),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject: NewHostRecord(
				netviewName, recordName,
				"", "", []HostRecordIpv4Addr{*resultIPV4Addrs}, []HostRecordIpv6Addr{*resultIPv6Addrs},
				eas, enabledns, dnsView, "", fakeRefReturn, useTtl, ttl, comment, aliases),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualRecord *HostRecord
		var err error
		It("should pass expected host record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateHostRecord(
				enabledns, false, recordName,
				netviewName, dnsView,
				ipv4Cidr, ipv6Cidr, "", "", macAddr, duid, useTtl, ttl, comment, eas, aliases)
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
		ipv4Cidr := "53.0.0.0/24"
		macAddr := "01:23:45:67:80:ab"
		ipv4Addr := fmt.Sprintf("func:nextavailableip:%s,%s", ipv4Cidr, netviewName)
		ipv6Cidr := "2003:db8:abcd:14::/64"
		duid := "02:24:46:68:81:cd"
		ipv6Addr := fmt.Sprintf("func:nextavailableip:%s,%s", ipv6Cidr, netviewName)
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		recordName := "test"
		enabledns := true
		enabledhcp := false
		dnsView := "default"
		fakeRefReturn := fmt.Sprintf("record:host/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)
		resultIPV4Addrs := NewHostRecordIpv4Addr(ipv4Addr, macAddr, enabledhcp, "")
		resultIPV6Addrs := NewHostRecordIpv6Addr(ipv6Addr, duid, enabledhcp, "")
		enableDNS := true
		useTtl := true
		ttl := uint32(70)
		comment := "test"
		aliases := []string{"abc.test.com"}

		aniFakeConnector := &fakeConnector{
			createObjectObj: NewHostRecord(
				netviewName, recordName,
				"", "", []HostRecordIpv4Addr{*resultIPV4Addrs}, []HostRecordIpv6Addr{*resultIPV6Addrs},
				nil, enableDNS, dnsView, "", "", useTtl, ttl, comment, aliases),
			getObjectRef: fakeRefReturn,
			getObjectObj: NewHostRecord(
				netviewName, recordName,
				"", "", []HostRecordIpv4Addr{*resultIPV4Addrs}, []HostRecordIpv6Addr{*resultIPV6Addrs},
				nil, enableDNS, dnsView, "", fakeRefReturn, useTtl, ttl, comment, aliases),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject: NewHostRecord(
				netviewName, recordName,
				"", "", []HostRecordIpv4Addr{*resultIPV4Addrs}, []HostRecordIpv6Addr{*resultIPV6Addrs},
				nil, enableDNS, dnsView, "", fakeRefReturn, useTtl, ttl, comment, aliases),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		ea := make(EA)
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
			actualRecord, err = objMgr.CreateHostRecord(
				enabledns, false, recordName, netviewName, dnsView, "", "",
				ipv4Addr, ipv6Addr, macAddr, duid, useTtl, ttl, comment, ea, aliases)
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
		ipv4Cidr := "53.0.0.0/24"
		macAddr := "01:23:45:67:80:ab"
		ipv4Addr := "53.0.0.1"
		ipv6Cidr := "2003:db8:abcd:14::/64"
		duid := "02:24:46:68:81:cd"
		ipv6Addr := "2003:db8:abcd:14::1"
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		enabledns := false
		enabledhcp := false
		dnsView := "default"
		recordName := "test"
		fakeRefReturn := fmt.Sprintf("record:host/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)
		resultIPV4Addrs := NewHostRecordIpv4Addr(ipv4Addr, macAddr, enabledhcp, "")
		resultIPV6Addrs := NewHostRecordIpv6Addr(ipv6Addr, duid, enabledhcp, "")
		useTtl := true
		ttl := uint32(70)
		comment := "test"
		aliases := []string{"test1"}

		aniFakeConnector := &fakeConnector{
			createObjectObj: NewHostRecord(
				netviewName, recordName,
				"", "", []HostRecordIpv4Addr{*resultIPV4Addrs}, []HostRecordIpv6Addr{*resultIPV6Addrs},
				nil, enabledns, dnsView, "", "", useTtl, ttl, comment, aliases),
			getObjectRef: fakeRefReturn,
			getObjectObj: NewHostRecord(
				netviewName, recordName,
				"", "", []HostRecordIpv4Addr{*resultIPV4Addrs}, []HostRecordIpv6Addr{*resultIPV6Addrs},
				nil, enabledns, dnsView, "", fakeRefReturn, useTtl, ttl, comment, aliases),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject: NewHostRecord(
				netviewName, recordName,
				"", "", []HostRecordIpv4Addr{*resultIPV4Addrs}, []HostRecordIpv6Addr{*resultIPV6Addrs},
				nil, enabledns, dnsView, "", fakeRefReturn, useTtl, ttl, comment, aliases),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		ea := make(EA)
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
			actualRecord, err = objMgr.CreateHostRecord(
				enabledns, false, recordName, netviewName, dnsView, ipv4Cidr,
				ipv6Cidr, ipv4Addr, ipv6Addr, macAddr, duid, useTtl, ttl, comment, ea, aliases)
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
		ipv4Cidr := "53.0.0.0/24"
		macAddr := "01:23:45:67:80:ab"
		ipv4Addr := "53.0.0.1"
		ipv6Cidr := "2003:db8:abcd:14::/64"
		duid := "02:24:46:68:81:cd"
		ipv6Addr := "2003:db8:abcd:14::1"
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		enabledns := true
		enabledhcp := false
		dnsView := "default"
		recordName := "test"
		fakeRefReturn := fmt.Sprintf("record:host/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)
		resultIPV4Addrs := NewHostRecordIpv4Addr(ipv4Addr, macAddr, enabledhcp, "")
		resultIPV6Addrs := NewHostRecordIpv6Addr(ipv6Addr, duid, enabledhcp, "")
		enableDNS := true
		useTtl := true
		ttl := uint32(70)
		comment := "test"
		aliases := []string{"abc.test.com"}

		aniFakeConnector := &fakeConnector{
			createObjectObj: NewHostRecord(
				netviewName, recordName,
				"", "", []HostRecordIpv4Addr{*resultIPV4Addrs}, []HostRecordIpv6Addr{*resultIPV6Addrs},
				nil, enableDNS, dnsView, "", "", useTtl, ttl, comment, aliases),
			getObjectRef: fakeRefReturn,
			getObjectObj: NewHostRecord(
				netviewName, recordName,
				"", "", []HostRecordIpv4Addr{*resultIPV4Addrs}, []HostRecordIpv6Addr{*resultIPV6Addrs},
				nil, enableDNS, dnsView, "", fakeRefReturn, useTtl, ttl, comment, aliases),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject: NewHostRecord(
				netviewName, recordName,
				"", "", []HostRecordIpv4Addr{*resultIPV4Addrs}, []HostRecordIpv6Addr{*resultIPV6Addrs},
				nil, enableDNS, dnsView, "", fakeRefReturn, useTtl, ttl, comment, aliases),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		ea := make(EA)
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
			actualRecord, err = objMgr.CreateHostRecord(
				enabledns, false, recordName, netviewName, dnsView, ipv4Cidr, ipv6Cidr,
				ipv4Addr, ipv6Addr, macAddr, duid, useTtl, ttl, comment, ea, aliases)
		})

		It("should return expected host record Object", func() {
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get Ipv4 and IPv6 Host Record Without DNS", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netview := "private"
		dnsview := "private"
		hostName := "test"
		ipv4Addr := "10.0.0.1"
		ipv6Addr := "2001:db8:abcd:14::1"
		fakeRefReturn := fmt.Sprintf("record:host/ZG5zLmJpbmRfY25h:%s/%20%20", hostName)
		queryParams := NewQueryParams(
			false,
			map[string]string{
				"name":         hostName,
				"network_view": netview,
				"view":         dnsview,
				"ipv4addr":     "10.0.0.1",
				"ipv6addr":     "2001:db8:abcd:14::1",
			})
		fipFakeConnector := &fakeConnector{
			getObjectObj:         NewEmptyHostRecord(),
			getObjectQueryParams: queryParams,
			getObjectRef:         "",
			resultObject: []HostRecord{*NewHostRecord(
				netview, hostName, ipv4Addr, ipv6Addr, nil, nil,
				nil, true, dnsview, "", fakeRefReturn, false, 0, "", []string{})},
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(fipFakeConnector, cmpType, tenantID)

		var actualhostRecord *HostRecord
		var err error
		It("should pass expected Host record Object to GetObject", func() {
			actualhostRecord, err = objMgr.GetHostRecord(netview, dnsview, hostName, ipv4Addr, ipv6Addr)
		})

		It("should return expected Host record Object", func() {
			Expect(*actualhostRecord).To(Equal(fipFakeConnector.resultObject.([]HostRecord)[0]))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get Host record by reference", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		hostName := "test"
		fakeRefReturn := fmt.Sprintf("record:host/ZG5zLmJpbmRfY25h:%s/%20%20", hostName)
		resObj := NewEmptyHostRecord()
		resObj.Ref = fakeRefReturn
		ncFakeConnector := &fakeConnector{
			getObjectObj:         NewEmptyHostRecord(),
			getObjectRef:         fakeRefReturn,
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         resObj,
			fakeRefReturn:        fakeRefReturn,
		}

		objMgr := NewObjectManager(ncFakeConnector, cmpType, tenantID)

		var actualRec *HostRecord
		var err error
		It("should pass expected host record object to GetObject", func() {
			actualRec, err = objMgr.GetHostRecordByRef(fakeRefReturn)
		})
		It("should return expected host record object", func() {
			Expect(err).To(BeNil())
			Expect(*actualRec).To(Equal(*resObj))
		})
	})

	Describe("Update host record", func() {
		var (
			err       error
			objMgr    IBObjectManager
			conn      *fakeConnector
			ref       string
			actualObj *HostRecord
		)

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		hostName := "host.test.com"
		refBase := "ZG5zLm5ldHdvcmtfdmlldyQyMw"
		ipv4Addr := "10.0.0.3"
		ipv6Addr := "2003:db8:abcd:14::1"
		useTtl := true
		ttl := uint32(70)

		It("Updating name, comment, aliases and EAs", func() {
			enableDNS := true
			ref = fmt.Sprintf("record:host/%s:%s", refBase, hostName)
			initialEas := EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_old_value",
				"ea3": "ea3_value",
				"ea4": "ea4_value",
				"ea5": "ea5_old_value"}
			initialAliases := []string{"abc.test.com", "xyz.test.com"}
			initObj := NewHostRecord("", hostName, "", "", []HostRecordIpv4Addr{},
				[]HostRecordIpv6Addr{}, initialEas, enableDNS, "someDNSview", "", "", useTtl, ttl, "old comment", initialAliases)
			initObj.Ref = ref

			setEas := EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_new_value",
				"ea2": "ea2_new_value",
				"ea5": "ea5_old_value"}
			expectedEas := setEas
			expectedAliases := []string{"abc.test.com", "trial.test.com"}

			comment := "test comment 1"
			updateUseTtl := false
			updateTtl := uint32(0)
			updateObjIn := NewHostRecord("", "host1.test.com", "", "", []HostRecordIpv4Addr{},
				[]HostRecordIpv6Addr{}, expectedEas, enableDNS, "someDNSview", "", "", updateUseTtl, updateTtl, comment, expectedAliases)
			updateObjIn.Ref = ref

			expectedObj := NewHostRecord("", "host1.test.com", "", "", []HostRecordIpv4Addr{},
				[]HostRecordIpv6Addr{}, expectedEas, enableDNS, "someDNSview", "", "", updateUseTtl, updateTtl, comment, expectedAliases)
			expectedObj.Ref = ref

			conn = &fakeConnector{
				getObjectObj:         NewEmptyHostRecord(),
				getObjectQueryParams: NewQueryParams(false, nil),
				getObjectRef:         ref,
				getObjectError:       nil,
				resultObject:         expectedObj,

				updateObjectObj:   updateObjIn,
				updateObjectRef:   ref,
				updateObjectError: nil,

				fakeRefReturn: ref,
			}
			objMgr = NewObjectManager(conn, cmpType, tenantID)

			actualObj, err = objMgr.UpdateHostRecord(ref, true, false, "host1.test.com", "",
				"someDNSview", "", "", "", "", "", "", updateUseTtl, updateTtl, comment, setEas, expectedAliases)
			Expect(err).To(BeNil())
			Expect(*actualObj).To(BeEquivalentTo(*expectedObj))
		})

		It("Updating MAC Address and DUID when IPv4 and Ipv6 addresses are passed", func() {
			enableDNS := true
			enableDHCP := false
			macAddr := "01:23:45:67:80:ab"
			duid := "02:24:46:68:81:cd"
			resultIPV4Addrs := NewHostRecordIpv4Addr(ipv4Addr, macAddr, enableDHCP, "")
			resultIPV6Addrs := NewHostRecordIpv6Addr(ipv6Addr, duid, enableDHCP, "")
			ref = fmt.Sprintf("record:host/%s:%s", refBase, hostName)

			updateObjIn := NewHostRecord("", hostName, "", "", []HostRecordIpv4Addr{*resultIPV4Addrs},
				[]HostRecordIpv6Addr{*resultIPV6Addrs}, nil, enableDNS, "", "", "", useTtl, ttl, "", []string{})
			updateObjIn.Ref = ref

			expectedObj := NewHostRecord("", hostName, "", "", []HostRecordIpv4Addr{*resultIPV4Addrs},
				[]HostRecordIpv6Addr{*resultIPV6Addrs}, nil, enableDNS, "", "", "", useTtl, ttl, "", []string{})
			expectedObj.Ref = ref

			conn = &fakeConnector{
				getObjectObj:         NewEmptyHostRecord(),
				getObjectQueryParams: NewQueryParams(false, nil),
				getObjectRef:         ref,
				getObjectError:       nil,
				resultObject:         expectedObj,

				updateObjectObj:   updateObjIn,
				updateObjectRef:   ref,
				updateObjectError: nil,

				fakeRefReturn: ref,
			}
			objMgr = NewObjectManager(conn, cmpType, tenantID)

			actualObj, err = objMgr.UpdateHostRecord(ref, enableDNS, false, hostName, "", "", "",
				"", ipv4Addr, ipv6Addr, macAddr, duid, useTtl, ttl, "", nil, []string{})
			Expect(err).To(BeNil())
			Expect(*actualObj).To(BeEquivalentTo(*expectedObj))
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
})
