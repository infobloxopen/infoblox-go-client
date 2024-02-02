package ibclient_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	ibclient "github.com/infobloxopen/infoblox-go-client/v2"
)

var _ = Describe("Object Manager: AAAA-record", func() {
	Describe("Allocate specific AAAA Record ", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		ipAddr := "2001:db8:abcd:14::1"
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		dnsView := "default"
		recordName := "test.domain.com"
		comment := "Test creation"
		fakeRefReturn := fmt.Sprintf("record:aaaa/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)
		ea := ibclient.EA{"VM Name": vmName, "VM ID": vmID}
		conn := &fakeConnector{
			createObjectObj: ibclient.NewRecordAAAA(
				dnsView, recordName, ipAddr, false, 0, comment, ea, ""),
			getObjectRef:         fakeRefReturn,
			getObjectObj:         ibclient.NewEmptyRecordAAAA(),
			getObjectQueryParams: ibclient.NewQueryParams(false, nil),
			resultObject: ibclient.NewRecordAAAA(
				dnsView, recordName, ipAddr, false, 0, comment, ea, fakeRefReturn),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := ibclient.NewObjectManager(conn, cmpType, tenantID)

		var actualRecord *ibclient.RecordAAAA
		var err error
		It("should pass expected AAAA record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateAAAARecord("", dnsView, recordName, "", ipAddr, false, uint32(0), comment, ea)
		})
		It("should return expected AAAA record Object", func() {
			Expect(actualRecord).To(Equal(conn.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Allocate next available AAAA Record ", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "private"
		cidr := "2001:db8:abcd:14::/64"
		ipAddr := fmt.Sprintf("func:nextavailableip:%s,%s", cidr, netviewName)
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		dnsView := "default"
		recordName := "test.domain.com"
		comment := "Test creation"
		fakeRefReturn := fmt.Sprintf("record:aaaa/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)
		ea := ibclient.EA{"VM Name": vmName, "VM ID": vmID}
		conn := &fakeConnector{
			createObjectObj: ibclient.NewRecordAAAA(
				dnsView, recordName, ipAddr, false, 0, comment, ea, ""),
			getObjectRef:         fakeRefReturn,
			getObjectObj:         ibclient.NewEmptyRecordAAAA(),
			getObjectQueryParams: ibclient.NewQueryParams(false, nil),
			resultObject: ibclient.NewRecordAAAA(
				dnsView, recordName, ipAddr, false, 0, comment, ea, fakeRefReturn),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := ibclient.NewObjectManager(conn, cmpType, tenantID)

		var actualRecord *ibclient.RecordAAAA
		var err error
		It("should pass expected AAAA record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateAAAARecord(netviewName, dnsView, recordName, cidr, "", false, uint32(0), comment, ea)
		})
		It("should return expected AAAA record Object", func() {
			Expect(actualRecord).To(Equal(conn.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Negative case: returns an error message when an IPv4 address is passed", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "private"
		cidr := "10.0.0./24"
		ipAddr := fmt.Sprintf("func:nextavailableip:%s,%s", cidr, netviewName)
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		dnsView := "default"
		recordName := "test.domain.com"
		comment := "Test creation"
		ea := ibclient.EA{"VM Name": vmName, "VM ID": vmID}
		conn := &fakeConnector{
			createObjectObj: ibclient.NewRecordAAAA(
				dnsView, recordName, ipAddr, false, 0, comment, ea, ""),
			createObjectError: fmt.Errorf("cannot parse CIDR value: invalid CIDR address: 10.0.0./24"),
		}

		objMgr := ibclient.NewObjectManager(conn, cmpType, tenantID)

		var actualRecord, expectedObj *ibclient.RecordAAAA
		var err error
		It("should pass expected AAAA record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateAAAARecord(netviewName, dnsView, recordName, cidr, "", false, uint32(0), comment, ea)
		})
		It("should return expected AAAA record Object", func() {
			Expect(actualRecord).To(Equal(expectedObj))
			Expect(err).To(Equal(conn.createObjectError))
		})
	})

	Describe("Get AAAA record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		dnsView := "default"
		recordName := "test.domain.com"
		ipAddr := "2001:db8:abcd:14::1"
		fakeRefReturn := fmt.Sprintf("record:aaaa/ZG5zLmJpbmRfY25h:%s/default", recordName)

		queryParams := ibclient.NewQueryParams(
			false,
			map[string]string{
				"view":     dnsView,
				"name":     recordName,
				"ipv6addr": ipAddr,
			})
		conn := &fakeConnector{
			getObjectRef:         "",
			getObjectObj:         ibclient.NewEmptyRecordAAAA(),
			getObjectQueryParams: queryParams,
			resultObject:         []ibclient.RecordAAAA{*ibclient.NewRecordAAAA(dnsView, recordName, ipAddr, false, 0, "", nil, fakeRefReturn)},
			fakeRefReturn:        fakeRefReturn,
		}

		objMgr := ibclient.NewObjectManager(conn, cmpType, tenantID)
		conn.resultObject.([]ibclient.RecordAAAA)[0].Ipv6Addr = &ipAddr
		var actualRecord *ibclient.RecordAAAA
		var err error
		It("should pass expected AAAA record Object to GetObject", func() {
			actualRecord, err = objMgr.GetAAAARecord(dnsView, recordName, ipAddr)
		})

		It("should return expected AAAA record Object", func() {
			Expect(*actualRecord).To(Equal(conn.resultObject.([]ibclient.RecordAAAA)[0]))
			Expect(err).To(BeNil())
		})
	})

	Describe("Negative case: returns an error when all the required fields are not passed", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		recordName := "test.domain.com"
		ipAddr := "2001:db8:abcd:14::1"
		fakeRefReturn := fmt.Sprintf("record:aaaa/ZG5zLmJpbmRfY25h:%s/default", recordName)

		queryParams := ibclient.NewQueryParams(
			false,
			map[string]string{
				"name":     recordName,
				"ipv6addr": ipAddr,
			})
		conn := &fakeConnector{
			getObjectRef:         "",
			getObjectObj:         ibclient.NewEmptyRecordAAAA(),
			getObjectQueryParams: queryParams,
			fakeRefReturn:        fakeRefReturn,
			getObjectError:       fmt.Errorf("DNS view, IPv6 address and record name of the record are required to retreive a unique AAAA record"),
		}

		objMgr := ibclient.NewObjectManager(conn, cmpType, tenantID)
		var actualRecord, expectedObj *ibclient.RecordAAAA
		var err error
		It("should pass expected AAAA record Object to GetObject", func() {
			actualRecord, err = objMgr.GetAAAARecord("", recordName, ipAddr)
		})

		It("should return expected AAAA record Object", func() {
			Expect(actualRecord).To(Equal(expectedObj))
			Expect(err).To(Equal(conn.getObjectError))
		})
	})

	Describe("Update AAAA-record, literal value", func() {
		var (
			err    error
			objMgr ibclient.IBObjectManager
			conn   *fakeConnector
		)

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		dnsView := "default"
		recordName := "test.domain.com"
		refBase := "ZG5zLm5ldHdvcmtfdmlldyQyMw"
		initIpAddr := "2001:db8:abcd:12::1"
		initUseTtl := true
		newRecordName := "test1.domain.com"
		newIpAddr := "2001:db8:abcd:12::2"
		newTtl := uint32(0)
		newUseTtl := false

		It("Negative case: updating an AAAA-record which does not exist", func() {
			initRef := fmt.Sprintf(
				"record:aaaa/%s:%s/%s",
				refBase, recordName, dnsView)
			getObjIn := ibclient.NewEmptyRecordAAAA()

			conn = &fakeConnector{
				getObjectObj:         getObjIn,
				getObjectQueryParams: ibclient.NewQueryParams(false, nil),
				getObjectRef:         initRef,
				getObjectError:       fmt.Errorf("test error"),
				resultObject:         ibclient.NewEmptyRecordAAAA(),
				fakeRefReturn:        "",
			}
			objMgr = ibclient.NewObjectManager(conn, cmpType, tenantID)

			_, err = objMgr.UpdateAAAARecord(initRef, "", newRecordName, "", newIpAddr, newUseTtl, newTtl, "some comment", nil)
			Expect(err).ToNot(BeNil())
		})

		It("Negative case: updating an AAAA-record with no update access", func() {
			initRef := fmt.Sprintf(
				"record:aaaa/%s:%s/%s",
				refBase, recordName, dnsView)
			initialEas := ibclient.EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_old_value",
				"ea3": "ea3_value",
				"ea4": "ea4_value",
				"ea5": "ea5_old_value"}
			initComment := "initial comment"
			initObj := ibclient.NewRecordAAAA(dnsView, recordName, initIpAddr, initUseTtl, newTtl, initComment, initialEas, initRef)

			getObjIn := ibclient.NewEmptyRecordAAAA()

			newEas := ibclient.EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_new_value",
				"ea2": "ea2_new_value",
				"ea5": "ea5_old_value"}

			newComment := "test comment 1"
			updateObjIn := ibclient.NewRecordAAAA("", newRecordName, newIpAddr, newUseTtl, newTtl, newComment, newEas, initRef)

			conn = &fakeConnector{
				getObjectObj:         getObjIn,
				getObjectQueryParams: ibclient.NewQueryParams(false, nil),
				getObjectRef:         initRef,
				getObjectError:       nil,
				resultObject:         initObj,

				updateObjectObj:   updateObjIn,
				updateObjectRef:   initRef,
				updateObjectError: fmt.Errorf("test error"),
				fakeRefReturn:     "",
			}
			objMgr = ibclient.NewObjectManager(conn, cmpType, tenantID)

			_, err = objMgr.UpdateAAAARecord(initRef, "", newRecordName, newIpAddr, "", newUseTtl, newTtl, newComment, newEas)
			Expect(err).ToNot(BeNil())
		})
	})
})
