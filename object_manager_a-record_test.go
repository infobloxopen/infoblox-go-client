package ibclient

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Object Manager: A-record", func() {
	Describe("Create a specific A-Record ", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "private"
		cidr := "53.0.0.0/24"
		ipAddr := "53.0.0.1"
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		dnsView := "default"
		recordName := "test"
		zone := "example.com"
		comment := "test comment"
		fakeRefReturn := fmt.Sprintf(
			"record:a/ZG5zLmJpbmRfY25h:%s/%s",
			recordName,
			netviewName)

		eas := make(EA)
		eas["VM ID"] = vmID
		eas["VM Name"] = vmName
		objectForCreation := NewRecordA(
			dnsView, "", recordName, ipAddr, 5, true, comment, eas, "")
		objectAsResult := NewRecordA(
			dnsView, zone, recordName, ipAddr, 5, true, comment, eas, fakeRefReturn)

		aniFakeConnector := &fakeConnector{
			createObjectObj:      objectForCreation,
			getObjectRef:         fakeRefReturn,
			getObjectObj:         NewEmptyRecordA(),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         objectAsResult,
			fakeRefReturn:        fakeRefReturn,
		}
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualRecord *RecordA
		var err error
		It("should pass expected A record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateARecord(netviewName, dnsView, recordName, cidr, ipAddr, 5, true, comment, eas)
		})
		It("should return expected A record Object", func() {
			Expect(err).To(BeNil())
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
		})
	})
	Describe("Create an A-Record using AllocateNextAvailableIp", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "private"
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		dnsView := "default"
		recordName := "test"
		comment := "test comment"
		ipAddrType := "IPV4" // Specify IP address type as IPv4
		objectType := "record:a"
		eas := make(EA)
		eas["VM ID"] = vmID
		eas["VM Name"] = vmName
		params := make(map[string][]string)

		fakeRefReturn := fmt.Sprintf(
			"record:a/ZG5zLmJpbmRfY25h:%s/%s",
			recordName,
			netviewName)

		objectForCreation := NewIpNextAvailable(recordName, objectType, nil, params, false, eas, comment, false, nil, ipAddrType,
			false, false, "", "", netviewName, "", false, 0, nil)
		objectAsResult := NewRecordA(
			dnsView, "", recordName, "53.0.0.1", 5, true, comment, eas, fakeRefReturn)

		fakeConnector := &fakeConnector{
			createObjectObj:      objectForCreation,
			getObjectRef:         fakeRefReturn,
			getObjectObj:         NewEmptyRecordA(),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         objectAsResult,
			fakeRefReturn:        fakeRefReturn,
		}

		objMgr := NewObjectManager(fakeConnector, cmpType, tenantID)

		It("should create A record with next available IP", func() {
			// Call the function to allocate next available IP and create A record
			result, err := objMgr.AllocateNextAvailableIp(recordName, objectType, nil, params, false, eas, comment, false, nil, ipAddrType,
				false, false, "", "", netviewName, "", false, 0, nil)

			// Assert no error occurred
			Expect(err).To(BeNil())

			// Assert the returned result is a valid A record
			actualRecord, ok := result.(*RecordA)
			Expect(ok).To(BeTrue())

			// Assert the returned record matches the expected result
			Expect(actualRecord).To(Equal(fakeConnector.resultObject))
		})
	})

	Describe("Get A record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		dnsView := "default"
		recordName := "test.domain.com"
		ipAddr := "10.0.0.2"
		fakeRefReturn := fmt.Sprintf("record:a/ZG5zLmJpbmRfY25h:%s/default", recordName)

		queryParams := NewQueryParams(
			false,
			map[string]string{
				"view":     dnsView,
				"name":     recordName,
				"ipv4addr": ipAddr,
			})
		conn := &fakeConnector{
			getObjectRef:         "",
			getObjectObj:         NewEmptyRecordA(),
			getObjectQueryParams: queryParams,
			resultObject:         []RecordA{*NewRecordA(dnsView, "", recordName, ipAddr, 0, false, "", nil, fakeRefReturn)},
			fakeRefReturn:        fakeRefReturn,
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)
		conn.resultObject.([]RecordA)[0].Ipv4Addr = &ipAddr
		var actualRecord *RecordA
		var err error
		It("should pass expected A record Object to GetObject", func() {
			actualRecord, err = objMgr.GetARecord(dnsView, recordName, ipAddr)
		})

		It("should return expected A record Object", func() {
			Expect(*actualRecord).To(Equal(conn.resultObject.([]RecordA)[0]))
			Expect(err).To(BeNil())
		})
	})

	Describe("Negative case: returns an error when all the required fields are not passed", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		recordName := "test.domain.com"
		ipAddr := "10.0.0.2"
		fakeRefReturn := fmt.Sprintf("record:a/ZG5zLmJpbmRfY25h:%s/default", recordName)

		queryParams := NewQueryParams(
			false,
			map[string]string{
				"name":     recordName,
				"ipv4addr": ipAddr,
			})
		conn := &fakeConnector{
			getObjectRef:         "",
			getObjectObj:         NewEmptyRecordA(),
			getObjectQueryParams: queryParams,
			fakeRefReturn:        fakeRefReturn,
			getObjectError:       fmt.Errorf("DNS view, IPv4 address and record name of the record are required to retreive a unique A record"),
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)
		var actualRecord, expectedObj *RecordA
		var err error
		It("should pass expected A record Object to GetObject", func() {
			actualRecord, err = objMgr.GetARecord("", recordName, ipAddr)
		})

		It("should return expected A record Object", func() {
			Expect(actualRecord).To(Equal(expectedObj))
			Expect(err).To(Equal(conn.getObjectError))
		})
	})

	Describe("Create an A-record by allocating next available IP address", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "private"
		cidr := "53.0.0.0/24"
		ipAddrReq := ""
		ipAddrRes := "53.0.0.1"
		ipAddrFunc := fmt.Sprintf("func:nextavailableip:%s,%s", cidr, netviewName)
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		dnsView := "default"
		recordName := "test"
		fakeRefReturn := fmt.Sprintf(
			"record:a/ZG5zLmJpbmRfY25h:%s/%s/%s",
			recordName,
			ipAddrRes,
			netviewName)

		aniFakeConnector := &fakeConnector{
			createObjectObj: NewRecordA(
				dnsView, "", recordName, ipAddrFunc, 0, false, "", nil, ""),
			getObjectRef:         fakeRefReturn,
			getObjectObj:         NewEmptyRecordA(),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject: NewRecordA(
				dnsView, "", recordName, ipAddrRes, 0, false, "", nil, fakeRefReturn),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		ea := make(EA)
		aniFakeConnector.createObjectObj.(*RecordA).Ea = ea
		aniFakeConnector.createObjectObj.(*RecordA).Ea["VM ID"] = vmID
		aniFakeConnector.createObjectObj.(*RecordA).Ea["VM Name"] = vmName

		aniFakeConnector.resultObject.(*RecordA).Ea = ea
		aniFakeConnector.resultObject.(*RecordA).Ea["VM ID"] = vmID
		aniFakeConnector.resultObject.(*RecordA).Ea["VM Name"] = vmName

		var actualRecord *RecordA
		var err error
		It("should pass expected A record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateARecord(netviewName, dnsView, recordName, cidr, ipAddrReq, 0, false, "", ea)
		})
		It("should return expected A record Object", func() {
			Expect(err).To(BeNil())
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
		})
	})

	Describe("Update A-record, literal value", func() {
		var (
			err    error
			objMgr IBObjectManager
			conn   *fakeConnector
			//actualObj       *RecordA
		)

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"

		//netView := "default"
		//netView2 := "notdefault"
		dnsView := "default"
		dnsZone := "test.loc"
		dnsName := "arec1.test.loc"
		refBase := "ZG5zLm5ldHdvcmtfdmlldyQyMw"
		initIPAddr := "10.2.1.56"
		initTTL := uint32(7)
		initUseTTL := true
		newIPAddr := "10.2.1.57"
		newTTL := uint32(70)
		newUseTTL := false
		//cidr := "10.2.1.0/24"
		//nextAvailableIPRequest := fmt.Sprintf(
		//	"func:nextavailableip:%s,%s", cidr, netView)
		//nextAvailableIPRequest2 := fmt.Sprintf(
		//	"func:nextavailableip:%s,%s", cidr, netView2)

		//It("updating IP address (with a literal value), comment, TTL, EAs", func() {
		//	initRef := fmt.Sprintf(
		//		"record:a/%s:%s/%s/%s",
		//		refBase, initIPAddr, dnsName, dnsView)
		//	newRef := fmt.Sprintf(
		//		"record:a/%s:%s/%s/%s",
		//		refBase, newIPAddr, dnsName, dnsView)
		//	initialEas := EA{
		//		"ea0": "ea0_old_value",
		//		"ea1": "ea1_old_value",
		//		"ea3": "ea3_value",
		//		"ea4": "ea4_value",
		//		"ea5": "ea5_old_value"}
		//	initComment := "initial comment"
		//	initObj := NewRecordA(dnsView, dnsZone, dnsName, initIPAddr, initTTL, initUseTTL, initComment, initialEas, initRef)
		//
		//	newEas := EA{
		//		"ea0": "ea0_old_value",
		//		"ea1": "ea1_new_value",
		//		"ea2": "ea2_new_value",
		//		"ea5": "ea5_old_value"}
		//
		//	getObjIn := NewEmptyRecordA()
		//
		//	newComment := "test comment 1"
		//	updateObjIn := NewRecordA(dnsView, dnsZone, dnsName, newIPAddr, newTTL, newUseTTL, newComment, newEas, initRef)
		//	expectedObj := NewRecordA(dnsView, dnsZone, dnsName, newIPAddr, newTTL, newUseTTL, newComment, newEas, newRef)
		//
		//	conn = &fakeConnector{
		//		getObjectObj:         getObjIn,
		//		getObjectQueryParams: NewQueryParams(false, nil),
		//		getObjectRef:         initRef,
		//		getObjectError:       nil,
		//		resultObject:         initObj,
		//
		//		updateObjectObj:   updateObjIn,
		//		updateObjectRef:   initRef,
		//		updateObjectError: nil,
		//
		//		fakeRefReturn: newRef,
		//	}
		//	objMgr = NewObjectManager(conn, cmpType, tenantID)
		//
		//	actualObj, err = objMgr.UpdateARecord(initRef, newIPAddr, "", "", newTTL, newUseTTL, newComment, newEas)
		//	Expect(err).To(BeNil())
		//	Expect(actualObj).To(BeEquivalentTo(expectedObj))
		//})

		//It("updating IP address (with 'nextavailableip' function), comment, TTL, EAs", func() {
		//	initRef := fmt.Sprintf(
		//		"record:a/%s:%s/%s/%s",
		//		refBase, initIPAddr, dnsName, dnsView)
		//	newRef := fmt.Sprintf(
		//		"record:a/%s:%s/%s/%s",
		//		refBase, newIPAddr, dnsName, dnsView)
		//	initialEas := EA{
		//		"ea0": "ea0_old_value",
		//		"ea1": "ea1_old_value",
		//		"ea3": "ea3_value",
		//		"ea4": "ea4_value",
		//		"ea5": "ea5_old_value"}
		//	initComment := "initial comment"
		//	initObj := NewRecordA(dnsView, dnsZone, dnsName, initIPAddr, initTTL, initUseTTL, initComment, initialEas, initRef)
		//
		//	newEas := EA{
		//		"ea0": "ea0_old_value",
		//		"ea1": "ea1_new_value",
		//		"ea2": "ea2_new_value",
		//		"ea5": "ea5_old_value"}
		//
		//	getObjIn := NewEmptyRecordA()
		//
		//	newComment := "test comment 1"
		//	updateObjIn := NewRecordA(dnsView, dnsZone, dnsName, nextAvailableIPRequest, newTTL, newUseTTL, newComment, newEas, initRef)
		//	expectedObj := NewRecordA(dnsView, dnsZone, dnsName, newIPAddr, newTTL, newUseTTL, newComment, newEas, newRef)
		//
		//	conn = &fakeConnector{
		//		getObjectObj:         getObjIn,
		//		getObjectQueryParams: NewQueryParams(false, nil),
		//		getObjectRef:         initRef,
		//		getObjectError:       nil,
		//		resultObject:         initObj,
		//
		//		updateObjectObj:   updateObjIn,
		//		updateObjectRef:   initRef,
		//		updateObjectError: nil,
		//
		//		fakeRefReturn: newRef,
		//	}
		//	objMgr = NewObjectManager(conn, cmpType, tenantID)
		//
		//	actualObj, err = objMgr.UpdateARecord(initRef, "", cidr, "", newTTL, newUseTTL, newComment, newEas)
		//	Expect(err).To(BeNil())
		//	Expect(actualObj).To(BeEquivalentTo(expectedObj))
		//})

		//It("updating IP address (with 'nextavailableip' function, non-default netview), comment, TTL, EAs", func() {
		//	initRef := fmt.Sprintf(
		//		"record:a/%s:%s/%s/%s",
		//		refBase, initIPAddr, dnsName, dnsView)
		//	newRef := fmt.Sprintf(
		//		"record:a/%s:%s/%s/%s",
		//		refBase, newIPAddr, dnsName, dnsView)
		//	initialEas := EA{
		//		"ea0": "ea0_old_value",
		//		"ea1": "ea1_old_value",
		//		"ea3": "ea3_value",
		//		"ea4": "ea4_value",
		//		"ea5": "ea5_old_value"}
		//	initComment := "initial comment"
		//	initObj := NewRecordA(dnsView, dnsZone, dnsName, initIPAddr, initTTL, initUseTTL, initComment, initialEas, initRef)
		//
		//	newEas := EA{
		//		"ea0": "ea0_old_value",
		//		"ea1": "ea1_new_value",
		//		"ea2": "ea2_new_value",
		//		"ea5": "ea5_old_value"}
		//
		//	getObjIn := NewEmptyRecordA()
		//
		//	newComment := "test comment 1"
		//	updateObjIn := NewRecordA(dnsView, dnsZone, dnsName, nextAvailableIPRequest2, newTTL, newUseTTL, newComment, newEas, initRef)
		//	expectedObj := NewRecordA(dnsView, dnsZone, dnsName, newIPAddr, newTTL, newUseTTL, newComment, newEas, newRef)
		//
		//	conn = &fakeConnector{
		//		getObjectObj:         getObjIn,
		//		getObjectQueryParams: NewQueryParams(false, nil),
		//		getObjectRef:         initRef,
		//		getObjectError:       nil,
		//		resultObject:         initObj,
		//
		//		updateObjectObj:   updateObjIn,
		//		updateObjectRef:   initRef,
		//		updateObjectError: nil,
		//
		//		fakeRefReturn: newRef,
		//	}
		//	objMgr = NewObjectManager(conn, cmpType, tenantID)
		//
		//	actualObj, err = objMgr.UpdateARecord(initRef, "", cidr, netView2, newTTL, newUseTTL, newComment, newEas)
		//	Expect(err).To(BeNil())
		//	Expect(actualObj).To(BeEquivalentTo(expectedObj))
		//})

		It("Negative case: updating an A-record which does not exist", func() {
			initRef := fmt.Sprintf(
				"record:a/%s:%s/%s/%s",
				refBase, initIPAddr, dnsName, dnsView)
			getObjIn := NewEmptyRecordA()

			conn = &fakeConnector{
				getObjectObj:         getObjIn,
				getObjectQueryParams: NewQueryParams(false, nil),
				getObjectRef:         initRef,
				getObjectError:       fmt.Errorf("test error"),
				resultObject:         NewEmptyRecordA(),
				fakeRefReturn:        "",
			}
			objMgr = NewObjectManager(conn, cmpType, tenantID)

			_, err = objMgr.UpdateARecord(initRef, dnsName, newIPAddr, "", "", newTTL, newUseTTL, "some comment", nil)
			Expect(err).ToNot(BeNil())
		})

		It("Negative case: updating an A-record with no update access", func() {
			initRef := fmt.Sprintf(
				"record:a/%s:%s/%s/%s",
				refBase, initIPAddr, dnsName, dnsView)
			initialEas := EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_old_value",
				"ea3": "ea3_value",
				"ea4": "ea4_value",
				"ea5": "ea5_old_value"}
			initComment := "initial comment"
			initObj := NewRecordA(dnsView, dnsZone, dnsName, initIPAddr, initTTL, initUseTTL, initComment, initialEas, initRef)

			getObjIn := NewEmptyRecordA()

			newEas := EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_new_value",
				"ea2": "ea2_new_value",
				"ea5": "ea5_old_value"}

			newComment := "test comment 1"
			updateObjIn := NewRecordA("", "", dnsName, newIPAddr, newTTL, newUseTTL, newComment, newEas, initRef)

			conn = &fakeConnector{
				getObjectObj:         getObjIn,
				getObjectQueryParams: NewQueryParams(false, nil),
				getObjectRef:         initRef,
				getObjectError:       nil,
				resultObject:         initObj,

				updateObjectObj:   updateObjIn,
				updateObjectRef:   initRef,
				updateObjectError: fmt.Errorf("test error"),
				fakeRefReturn:     "",
			}
			objMgr = NewObjectManager(conn, cmpType, tenantID)

			_, err = objMgr.UpdateARecord(initRef, dnsName, newIPAddr, "", "", newTTL, newUseTTL, newComment, newEas)
			Expect(err).ToNot(BeNil())
		})
	})

	Describe("Delete A Record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		recordName := "test"
		deleteRef := fmt.Sprintf("record:a/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)
		fakeRefReturn := deleteRef
		nwFakeConnector := &fakeConnector{
			deleteObjectRef: deleteRef,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(nwFakeConnector, cmpType, tenantID)

		var actualRef string
		var err error
		It("should pass expected A record Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteARecord(deleteRef)
		})
		It("should return expected A record Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})
})
