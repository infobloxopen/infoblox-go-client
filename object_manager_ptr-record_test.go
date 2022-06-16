package ibclient

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Object Manager: PTR-record", func() {
	Describe("Allocate specific PTR Record with IPv6 Address", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		ipAddr := "2001:db8:abcd:14::1"
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		dnsView := "default"
		ptrdname := "test"
		useTtl := true
		ttl := uint32(70)
		comment := "creation test"
		eas := EA{"VM Name": vmName, "VM ID": vmID}
		fakeRefReturn := fmt.Sprintf("record:ptr/ZG5zLmJpbmRfY25h:1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.4.1.0.0.d.c.b.a.8.b.d.0.1.0.0.2.ip6.arpa/default")

		conn := &fakeConnector{
			createObjectObj:      NewRecordPTR(dnsView, ptrdname, useTtl, ttl, comment, eas),
			getObjectRef:         fakeRefReturn,
			getObjectObj:         NewEmptyRecordPTR(),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         NewRecordPTR(dnsView, ptrdname, useTtl, ttl, comment, eas),
			fakeRefReturn:        fakeRefReturn,
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)
		conn.createObjectObj.(*RecordPTR).Ipv6Addr = ipAddr

		var actualRecord *RecordPTR
		var err error
		It("should pass expected PTR record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreatePTRRecord("", dnsView, ptrdname, "", "", ipAddr, useTtl, ttl, comment, eas)
		})
		It("should return expected PTR record Object", func() {
			Expect(actualRecord).To(Equal(conn.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Allocate specific PTR Record with IPv4 Address", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		ipAddr := "10.0.0.1"
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		dnsView := "default"
		ptrdname := "test"
		useTtl := true
		ttl := uint32(70)
		comment := "creation test"
		eas := EA{"VM Name": vmName, "VM ID": vmID}
		fakeRefReturn := fmt.Sprintf("record:ptr/ZG5zLmJpbmRfY25h:2.0.0.10.in-addr.arpa/default")

		conn := &fakeConnector{
			createObjectObj:      NewRecordPTR(dnsView, ptrdname, useTtl, ttl, comment, eas),
			getObjectRef:         fakeRefReturn,
			getObjectObj:         NewEmptyRecordPTR(),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         NewRecordPTR(dnsView, ptrdname, useTtl, ttl, comment, eas),
			fakeRefReturn:        fakeRefReturn,
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)
		conn.createObjectObj.(*RecordPTR).Ipv4Addr = ipAddr

		var actualRecord *RecordPTR
		var err error
		It("should pass expected PTR record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreatePTRRecord("", dnsView, ptrdname, "", "", ipAddr, useTtl, ttl, comment, eas)
		})
		It("should return expected PTR record Object", func() {
			Expect(actualRecord).To(Equal(conn.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Allocate next available PTR Record-IPv4", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "private"
		cidr := "10.0.0.0/24"
		ipAddr := fmt.Sprintf("func:nextavailableip:%s,%s", cidr, netviewName)
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		dnsView := "default"
		ptrdname := "test"
		comment := "creation test"
		eas := EA{"VM Name": vmName, "VM ID": vmID}
		useTtl := true
		ttl := uint32(70)
		fakeRefReturn := fmt.Sprintf("record:ptr/ZG5zLmJpbmRfY25h:2.0.0.10.in-addr.arpa/default")

		conn := &fakeConnector{
			createObjectObj:      NewRecordPTR(dnsView, ptrdname, useTtl, ttl, comment, eas),
			getObjectRef:         fakeRefReturn,
			getObjectObj:         NewEmptyRecordPTR(),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         NewRecordPTR(dnsView, ptrdname, useTtl, ttl, comment, eas),
			fakeRefReturn:        fakeRefReturn,
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)
		conn.createObjectObj.(*RecordPTR).Ipv4Addr = ipAddr
		var actualRecord *RecordPTR
		var err error
		It("should pass expected PTR record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreatePTRRecord(netviewName, dnsView, ptrdname, "", cidr, "", useTtl, ttl, comment, eas)
		})
		It("should return expected PTR record Object", func() {
			Expect(actualRecord).To(Equal(conn.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Allocate next available PTR Record-IPv6", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "private"
		cidr := "2001:db8:abcd:14::/64"
		ipAddr := fmt.Sprintf("func:nextavailableip:%s,%s", cidr, netviewName)
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		dnsView := "default"
		ptrdname := "test"
		comment := "creation test"
		eas := EA{"VM Name": vmName, "VM ID": vmID}
		useTtl := true
		ttl := uint32(70)
		fakeRefReturn := fmt.Sprintf("record:ptr/ZG5zLmJpbmRfY25h:2.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.4.1.0.0.d.c.b.a.8.b.d.0.1.0.0.2.ip6.arpa/default")

		conn := &fakeConnector{
			createObjectObj:      NewRecordPTR(dnsView, ptrdname, useTtl, ttl, comment, eas),
			getObjectRef:         fakeRefReturn,
			getObjectObj:         NewEmptyRecordPTR(),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         NewRecordPTR(dnsView, ptrdname, useTtl, ttl, comment, eas),
			fakeRefReturn:        fakeRefReturn,
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)
		conn.createObjectObj.(*RecordPTR).Ipv6Addr = ipAddr
		var actualRecord *RecordPTR
		var err error
		It("should pass expected PTR record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreatePTRRecord(netviewName, dnsView, ptrdname, "", cidr, "", useTtl, ttl, comment, eas)
		})
		It("should return expected PTR record Object", func() {
			Expect(actualRecord).To(Equal(conn.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Allocate a PTR Record in forward mapping zone", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		recordName := "test-ptr-record.test.com"
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		dnsView := "default"
		ptrdname := "test"
		comment := "creation test"
		eas := EA{"VM Name": vmName, "VM ID": vmID}
		useTtl := true
		ttl := uint32(70)
		fakeRefReturn := fmt.Sprintf("record:ptr/ZG5zLmJpbmRfY25h:%s/%s", recordName, dnsView)

		conn := &fakeConnector{
			createObjectObj:      NewRecordPTR(dnsView, ptrdname, useTtl, ttl, comment, eas),
			getObjectRef:         fakeRefReturn,
			getObjectObj:         NewEmptyRecordPTR(),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         NewRecordPTR(dnsView, ptrdname, useTtl, ttl, comment, eas),
			fakeRefReturn:        fakeRefReturn,
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)
		conn.createObjectObj.(*RecordPTR).Name = recordName
		var actualRecord *RecordPTR
		var err error
		It("should pass expected PTR record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreatePTRRecord("", dnsView, ptrdname, recordName, "", "", useTtl, ttl, comment, eas)
		})
		It("should return expected PTR record Object", func() {
			Expect(actualRecord).To(Equal(conn.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Negative case: returns an error message if ptrdname is not entered", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		recordName := "test-ptr-record.test.com"
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		dnsView := "default"
		comment := "creation test"
		eas := EA{"VM Name": vmName, "VM ID": vmID}
		useTtl := true
		ttl := uint32(70)

		conn := &fakeConnector{
			createObjectObj:   NewRecordPTR(dnsView, "", useTtl, ttl, comment, eas),
			createObjectError: fmt.Errorf("ptrdname is a required field to create a PTR record"),
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)
		conn.createObjectObj.(*RecordPTR).Name = recordName
		var actualRecord, expectedObj *RecordPTR
		var err error
		expectedObj = nil
		It("should pass expected PTR record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreatePTRRecord("", dnsView, "", recordName, "", "", useTtl, 0, comment, eas)
			Expect(actualRecord).To(Equal(expectedObj))
			Expect(err).To(Not(BeNil()))
		})
	})

	Describe("Negative case: returns an error message if an invalid IP address is passed", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		ipAddr := "10.0.0.300"
		ptrdname := "ptr-test.infoblox.com"
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		dnsView := "default"
		comment := "creation test"
		eas := EA{"VM Name": vmName, "VM ID": vmID}
		useTtl := true
		ttl := uint32(70)

		conn := &fakeConnector{
			createObjectObj:   NewRecordPTR(dnsView, "", useTtl, ttl, comment, eas),
			createObjectError: fmt.Errorf("%s is an invalid IP address", ipAddr),
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)
		conn.createObjectObj.(*RecordPTR).Ipv4Addr = ipAddr
		var actualRecord, expectedObj *RecordPTR
		var err error
		expectedObj = nil
		It("should pass expected PTR record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreatePTRRecord("", dnsView, ptrdname, "", "", ipAddr, useTtl, 0, comment, eas)
			Expect(actualRecord).To(Equal(expectedObj))
			Expect(err).To(Equal(conn.createObjectError))
		})
	})

	Describe("Negative case: returns an error message if the required fields for creation of record is empty", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		dnsView := "default"
		comment := "creation test"
		ptrdname := "ptr-test.infoblox.com"
		eas := EA{"VM Name": vmName, "VM ID": vmID}
		useTtl := true
		ttl := uint32(70)

		conn := &fakeConnector{
			createObjectObj: NewRecordPTR(dnsView, "", useTtl, ttl, comment, eas),
			createObjectError: fmt.Errorf("CIDR and network view are required to allocate a next available IP address\n" +
				"IP address is required to create PTR record in reverse mapping zone\n" +
				"record name is required to create a record in forwarrd mapping zone"),
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)
		var actualRecord, expectedObj *RecordPTR
		var err error
		expectedObj = nil
		It("should pass expected PTR record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreatePTRRecord("", dnsView, ptrdname, "", "", "", useTtl, 0, comment, eas)
			Expect(actualRecord).To(Equal(expectedObj))
			Expect(err).To(Equal(conn.createObjectError))
		})
	})

	Describe("Get PTR record-IPv4", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		dnsView := "default"
		ptrdname := "test"
		ipAddr := "10.0.0.1"
		useTtl := true
		ttl := uint32(70)
		fakeRefReturn := fmt.Sprintf("record:ptr/ZG5zLmJpbmRfY25h:1.0.0.10.in-addr.arpa/default")

		queryParams := NewQueryParams(
			false,
			map[string]string{
				"view":     dnsView,
				"ptrdname": ptrdname,
				"ipv4addr": ipAddr,
			})
		conn := &fakeConnector{
			getObjectRef:         "",
			getObjectObj:         NewEmptyRecordPTR(),
			getObjectQueryParams: queryParams,
			resultObject:         []RecordPTR{*NewRecordPTR(dnsView, ptrdname, useTtl, ttl, "", nil)},
			fakeRefReturn:        fakeRefReturn,
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)
		conn.resultObject.([]RecordPTR)[0].Ipv4Addr = ipAddr
		var actualRecord *RecordPTR
		var err error
		It("should pass expected PTR record Object to GetObject", func() {
			actualRecord, err = objMgr.GetPTRRecord(dnsView, ptrdname, "", ipAddr)
		})

		It("should return expected PTR record Object", func() {
			Expect(*actualRecord).To(Equal(conn.resultObject.([]RecordPTR)[0]))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get PTR record-IPv6", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		dnsView := "default"
		ptrdname := "test"
		ipAddr := "2001:db8:abcd:14::1"
		useTtl := true
		ttl := uint32(70)
		fakeRefReturn := fmt.Sprintf("record:ptr/ZG5zLmJpbmRfY25h:1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.4.1.0.0.d.c.b.a.8.b.d.0.1.0.0.2.ip6.arpa/default")

		queryParams := NewQueryParams(
			false,
			map[string]string{
				"view":     dnsView,
				"ptrdname": ptrdname,
				"ipv6addr": ipAddr,
			})
		conn := &fakeConnector{
			getObjectRef:         "",
			getObjectObj:         NewEmptyRecordPTR(),
			getObjectQueryParams: queryParams,
			resultObject:         []RecordPTR{*NewRecordPTR(dnsView, ptrdname, useTtl, ttl, "", nil)},
			fakeRefReturn:        fakeRefReturn,
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)

		var actualRecord *RecordPTR
		var err error
		It("should pass expected PTR record Object to GetObject", func() {
			actualRecord, err = objMgr.GetPTRRecord(dnsView, ptrdname, "", ipAddr)
		})

		It("should return expected PTR record Object", func() {
			Expect(*actualRecord).To(Equal(conn.resultObject.([]RecordPTR)[0]))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get PTR record-name(forward mapping zone)", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		dnsView := "default"
		ptrdname := "test"
		recordName := "test-ptr-record.test.com"
		useTtl := true
		ttl := uint32(70)
		fakeRefReturn := fmt.Sprintf("record:ptr/ZG5zLmJpbmRfY25h:%s/%s", recordName, dnsView)

		queryParams := NewQueryParams(
			false,
			map[string]string{
				"view":     dnsView,
				"ptrdname": ptrdname,
				"name":     recordName,
			})
		conn := &fakeConnector{
			getObjectRef:         "",
			getObjectObj:         NewEmptyRecordPTR(),
			getObjectQueryParams: queryParams,
			resultObject:         []RecordPTR{*NewRecordPTR(dnsView, ptrdname, useTtl, ttl, "", nil)},
			fakeRefReturn:        fakeRefReturn,
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)

		var actualRecord *RecordPTR
		var err error
		It("should pass expected PTR record Object to GetObject", func() {
			actualRecord, err = objMgr.GetPTRRecord(dnsView, ptrdname, recordName, "")
		})

		It("should return expected PTR record Object", func() {
			Expect(*actualRecord).To(Equal(conn.resultObject.([]RecordPTR)[0]))
			Expect(err).To(BeNil())
		})
	})

	Describe("Update PTR record", func() {
		var (
			err       error
			objMgr    IBObjectManager
			conn      *fakeConnector
			ref       string
			actualObj *RecordPTR
		)

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		refBase := "ZG5zLm5ldHdvcmtfdmlldyQyMw"
		ptrdname := "test"
		ipv4Addr := "10.0.0.1"
		ipv6Addr := "2001:db8:abcd:14::1"
		recordName := "test-ptr-record.test.com"
		useTtl := false
		ttl := uint32(0)
		netview := "private"
		ipv4cidr := "10.0.0.0/24"
		ipv6cidr := "2001:db8:abcd:14::/64"

		It("IPv4, updating ptrdname, IPv4 address, comment and EAs", func() {
			ref = fmt.Sprintf("record:ptr/%s:1.0.0.10.in-addr.arpa/default", refBase)
			initialEas := EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_old_value",
				"ea3": "ea3_value",
				"ea4": "ea4_value",
				"ea5": "ea5_old_value"}
			initObj := NewRecordPTR("", ptrdname, useTtl, ttl, "old comment", initialEas)
			initObj.Ref = ref
			initObj.Ipv4Addr = ipv4Addr

			setEas := EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_new_value",
				"ea2": "ea2_new_value",
				"ea5": "ea5_old_value"}
			expectedEas := setEas

			comment := "test comment 1"
			updateUseTtl := true
			updateTtl := uint32(10)
			newPtrdname := "test-update-ptr.test.com"
			updateIpAddr := "10.0.0.2"
			updatedRef := fmt.Sprintf("record:ptr/%s:2.0.0.10.in-addr.arpa/default", refBase)
			updateObjIn := NewRecordPTR("", newPtrdname, updateUseTtl, updateTtl, comment, expectedEas)
			updateObjIn.Ref = ref
			updateObjIn.Ipv4Addr = updateIpAddr

			expectedObj := NewRecordPTR("", newPtrdname, updateUseTtl, updateTtl, comment, expectedEas)
			expectedObj.Ref = updatedRef
			expectedObj.Ipv4Addr = updateIpAddr

			conn = &fakeConnector{
				getObjectObj:         NewEmptyRecordPTR(),
				getObjectQueryParams: NewQueryParams(false, nil),
				getObjectRef:         updatedRef,
				getObjectError:       nil,
				resultObject:         expectedObj,

				updateObjectObj:   updateObjIn,
				updateObjectRef:   ref,
				updateObjectError: nil,

				fakeRefReturn: updatedRef,
			}
			objMgr = NewObjectManager(conn, cmpType, tenantID)

			actualObj, err = objMgr.UpdatePTRRecord(ref, "", newPtrdname, "", "", updateIpAddr, updateUseTtl, updateTtl, comment, setEas)
			Expect(err).To(BeNil())
			Expect(*actualObj).To(BeEquivalentTo(*expectedObj))
		})

		It("IPv6: updating ptrdname, TTl fields, IPv6 address, comment and EAs", func() {
			ref = fmt.Sprintf("record:ptr/%s:1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.4.1.0.0.d.c.b.a.8.b.d.0.1.0.0.2.ip6.arpa/default", refBase)
			initialEas := EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_old_value",
				"ea3": "ea3_value",
				"ea4": "ea4_value",
				"ea5": "ea5_old_value"}
			initObj := NewRecordPTR("", ptrdname, useTtl, ttl, "old comment", initialEas)
			initObj.Ref = ref
			initObj.Ipv4Addr = ipv6Addr

			setEas := EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_new_value",
				"ea2": "ea2_new_value",
				"ea5": "ea5_old_value"}
			expectedEas := setEas

			comment := "test comment 1"
			updateUseTtl := true
			updateTtl := uint32(10)
			newPtrdname := "test-update"
			updateIpAddr := "2001:db8:abcd:14::2"
			updatedRef := fmt.Sprintf("record:ptr/%s:2.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.4.1.0.0.d.c.b.a.8.b.d.0.1.0.0.2.ip6.arpa/default", refBase)
			updateObjIn := NewRecordPTR("", newPtrdname, updateUseTtl, updateTtl, comment, expectedEas)
			updateObjIn.Ref = ref
			updateObjIn.Ipv6Addr = updateIpAddr

			expectedObj := NewRecordPTR("", newPtrdname, updateUseTtl, updateTtl, comment, expectedEas)
			expectedObj.Ref = updatedRef
			expectedObj.Ipv6Addr = updateIpAddr

			conn = &fakeConnector{
				getObjectObj:         NewEmptyRecordPTR(),
				getObjectQueryParams: NewQueryParams(false, nil),
				getObjectRef:         updatedRef,
				getObjectError:       nil,
				resultObject:         expectedObj,

				updateObjectObj:   updateObjIn,
				updateObjectRef:   ref,
				updateObjectError: nil,

				fakeRefReturn: updatedRef,
			}
			objMgr = NewObjectManager(conn, cmpType, tenantID)

			actualObj, err = objMgr.UpdatePTRRecord(ref, "", newPtrdname, "", "", updateIpAddr, updateUseTtl, updateTtl, comment, setEas)
			Expect(err).To(BeNil())
			Expect(*actualObj).To(BeEquivalentTo(*expectedObj))
		})

		It("IPv4, updating ptrdname, IPv4 address by passing cidr and network view, comment and EAs", func() {
			ref = fmt.Sprintf("record:ptr/%s:1.0.0.10.in-addr.arpa/default", refBase)
			initialEas := EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_old_value",
				"ea3": "ea3_value",
				"ea4": "ea4_value",
				"ea5": "ea5_old_value"}
			initObj := NewRecordPTR("", ptrdname, useTtl, ttl, "old comment", initialEas)
			initObj.Ref = ref
			initObj.Ipv4Addr = ipv4Addr

			setEas := EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_new_value",
				"ea2": "ea2_new_value",
				"ea5": "ea5_old_value"}
			expectedEas := setEas

			comment := "test comment 1"
			updateUseTtl := true
			updateTtl := uint32(10)
			newPtrdname := "test-update-ptr.test.com"
			updateIpAddr := fmt.Sprintf("func:nextavailableip:%s,%s", ipv4cidr, netview)
			updatedRef := fmt.Sprintf("record:ptr/%s:2.0.0.10.in-addr.arpa/default", refBase)
			updateObjIn := NewRecordPTR("", newPtrdname, updateUseTtl, updateTtl, comment, expectedEas)
			updateObjIn.Ref = ref
			updateObjIn.Ipv4Addr = updateIpAddr

			expectedObj := NewRecordPTR("", newPtrdname, updateUseTtl, updateTtl, comment, expectedEas)
			expectedObj.Ref = updatedRef
			expectedObj.Ipv4Addr = updateIpAddr

			conn = &fakeConnector{
				getObjectObj:         NewEmptyRecordPTR(),
				getObjectQueryParams: NewQueryParams(false, nil),
				getObjectRef:         updatedRef,
				getObjectError:       nil,
				resultObject:         expectedObj,

				updateObjectObj:   updateObjIn,
				updateObjectRef:   ref,
				updateObjectError: nil,

				fakeRefReturn: updatedRef,
			}
			objMgr = NewObjectManager(conn, cmpType, tenantID)

			actualObj, err = objMgr.UpdatePTRRecord(ref, netview, newPtrdname, "", ipv4cidr, "", updateUseTtl, updateTtl, comment, setEas)
			Expect(err).To(BeNil())
			Expect(*actualObj).To(BeEquivalentTo(*expectedObj))
		})

		It("IPv6, updating ptrdname, IPv6 address by passing cidr and network view, comment and EAs", func() {
			ref = fmt.Sprintf("record:ptr/%s:1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.4.1.0.0.d.c.b.a.8.b.d.0.1.0.0.2.ip6.arpa/default", refBase)
			initialEas := EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_old_value",
				"ea3": "ea3_value",
				"ea4": "ea4_value",
				"ea5": "ea5_old_value"}
			initObj := NewRecordPTR("", ptrdname, useTtl, ttl, "old comment", initialEas)
			initObj.Ref = ref
			initObj.Ipv6Addr = ipv6Addr

			setEas := EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_new_value",
				"ea2": "ea2_new_value",
				"ea5": "ea5_old_value"}
			expectedEas := setEas

			comment := "test comment 1"
			updateUseTtl := true
			updateTtl := uint32(10)
			newPtrdname := "test-update-ptr.test.com"
			updateIpAddr := fmt.Sprintf("func:nextavailableip:%s,%s", ipv6cidr, netview)
			updatedRef := fmt.Sprintf("record:ptr/%s:2.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.4.1.0.0.d.c.b.a.8.b.d.0.1.0.0.2.ip6.arpa/default", refBase)
			updateObjIn := NewRecordPTR("", newPtrdname, updateUseTtl, updateTtl, comment, expectedEas)
			updateObjIn.Ref = ref
			updateObjIn.Ipv6Addr = updateIpAddr

			expectedObj := NewRecordPTR("", newPtrdname, updateUseTtl, updateTtl, comment, expectedEas)
			expectedObj.Ref = updatedRef
			expectedObj.Ipv6Addr = updateIpAddr

			conn = &fakeConnector{
				getObjectObj:         NewEmptyRecordPTR(),
				getObjectQueryParams: NewQueryParams(false, nil),
				getObjectRef:         updatedRef,
				getObjectError:       nil,
				resultObject:         expectedObj,

				updateObjectObj:   updateObjIn,
				updateObjectRef:   ref,
				updateObjectError: nil,

				fakeRefReturn: updatedRef,
			}
			objMgr = NewObjectManager(conn, cmpType, tenantID)

			actualObj, err = objMgr.UpdatePTRRecord(ref, netview, newPtrdname, "", ipv6cidr, "", updateUseTtl, updateTtl, comment, setEas)
			Expect(err).To(BeNil())
			Expect(*actualObj).To(BeEquivalentTo(*expectedObj))
		})

		It("Updating ptrdname, TTl fields, record name, comment and EAs", func() {
			ref = fmt.Sprintf("record:ptr/%s:%s/default", refBase, recordName)
			initialEas := EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_old_value",
				"ea3": "ea3_value",
				"ea4": "ea4_value",
				"ea5": "ea5_old_value"}
			initObj := NewRecordPTR("", ptrdname, useTtl, ttl, "old comment", initialEas)
			initObj.Ref = ref
			initObj.Name = recordName

			setEas := EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_new_value",
				"ea2": "ea2_new_value",
				"ea5": "ea5_old_value"}
			expectedEas := setEas

			comment := "test comment 1"
			updateUseTtl := true
			updateTtl := uint32(10)
			newPtrdname := "test-update"
			updateName := "test-ptr-update"
			updatedRef := fmt.Sprintf("record:ptr/%s:%s/20", refBase, newPtrdname)
			updateObjIn := NewRecordPTR("", newPtrdname, updateUseTtl, updateTtl, comment, expectedEas)
			updateObjIn.Ref = ref
			updateObjIn.Name = updateName

			expectedObj := NewRecordPTR("", newPtrdname, updateUseTtl, updateTtl, comment, expectedEas)
			expectedObj.Ref = updatedRef
			expectedObj.Name = updateName

			conn = &fakeConnector{
				getObjectObj:         NewEmptyRecordPTR(),
				getObjectQueryParams: NewQueryParams(false, nil),
				getObjectRef:         updatedRef,
				getObjectError:       nil,
				resultObject:         expectedObj,

				updateObjectObj:   updateObjIn,
				updateObjectRef:   ref,
				updateObjectError: nil,

				fakeRefReturn: updatedRef,
			}
			objMgr = NewObjectManager(conn, cmpType, tenantID)

			actualObj, err = objMgr.UpdatePTRRecord(ref, "", newPtrdname, updateName, "", "", updateUseTtl, updateTtl, comment, setEas)
			Expect(err).To(BeNil())
			Expect(*actualObj).To(BeEquivalentTo(*expectedObj))
		})
	})

	Describe("Update PTR record", func() {
		var (
			err       error
			objMgr    IBObjectManager
			conn      *fakeConnector
			ref       string
			actualObj *RecordCNAME
		)

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		refBase := "ZG5zLm5ldHdvcmtfdmlldyQyMw"
		canonical := "test-canonical.domain.com"
		recordName := "test.domain.com"
		useTtl := false
		ttl := uint32(0)

		It("IPv4, updating ptrdname, IPv4 address, comment and EAs", func() {
			ref = fmt.Sprintf("record:cname/%s:%s", refBase, recordName)
			initialEas := EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_old_value",
				"ea3": "ea3_value",
				"ea4": "ea4_value",
				"ea5": "ea5_old_value"}
			initObj := NewRecordCNAME("", canonical, recordName, useTtl, ttl, "old comment", initialEas, ref)
			initObj.Ref = ref

			setEas := EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_new_value",
				"ea2": "ea2_new_value",
				"ea5": "ea5_old_value"}
			expectedEas := setEas

			comment := "test comment 1"
			updateUseTtl := true
			updateTtl := uint32(10)
			newCanonical := "test-canonical-update.domain.com"
			newRecordName := "test-update.domain.com"
			updatedRef := fmt.Sprintf("record:cname/%s:%s", refBase, newRecordName)
			updateObjIn := NewRecordCNAME("", newCanonical, newRecordName, updateUseTtl, updateTtl, comment, expectedEas, ref)

			expectedObj := NewRecordCNAME("", newCanonical, newRecordName, updateUseTtl, updateTtl, comment, expectedEas, updatedRef)

			conn = &fakeConnector{
				getObjectObj:         NewEmptyRecordCNAME(),
				getObjectQueryParams: NewQueryParams(false, nil),
				getObjectRef:         updatedRef,
				getObjectError:       nil,
				resultObject:         expectedObj,

				updateObjectObj:   updateObjIn,
				updateObjectRef:   ref,
				updateObjectError: nil,

				fakeRefReturn: updatedRef,
			}
			objMgr = NewObjectManager(conn, cmpType, tenantID)

			actualObj, err = objMgr.UpdateCNAMERecord(ref, newCanonical, newRecordName, updateUseTtl, updateTtl, comment, setEas)
			Expect(err).To(BeNil())
			Expect(*actualObj).To(BeEquivalentTo(*expectedObj))
		})
	})

	Describe("Delete PTR Record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		recordName := "test"
		deleteRef := fmt.Sprintf("record:ptr/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)
		fakeRefReturn := deleteRef
		nwFakeConnector := &fakeConnector{
			deleteObjectRef: deleteRef,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(nwFakeConnector, cmpType, tenantID)

		var actualRef string
		var err error
		It("should pass expected PTR record Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeletePTRRecord(deleteRef)
		})
		It("should return expected PTR record Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})
})
