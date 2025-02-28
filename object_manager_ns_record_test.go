package ibclient

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Object Manager: NS-record", func() {
	Describe("Create a specific NS-Record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		nameserver := "testing.test.com"
		view := "default"
		name := "test.com"
		addresses := []*ZoneNameServer{
			{
				Address:       "3.4.4.5",
				AutoCreatePtr: true,
			},
		}
		fakeRefReturn := fmt.Sprintf(
			"record:ns/ZG5zLmJpbmRfbnMkLjMuY29tLnRlc3QuLm5hbWUudGVzdC5jb20:%s/%s/%s",
			nameserver,
			name,
			view)

		objectAsResult := NewRecordNS(name, nameserver, view, addresses, "")
		objectAsResult.Ref = fakeRefReturn
		aniFakeConnector := &fakeConnector{
			createObjectObj:      NewRecordNS(name, nameserver, view, addresses, ""),
			getObjectRef:         fakeRefReturn,
			getObjectObj:         NewEmptyRecordNS(),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         objectAsResult,
			fakeRefReturn:        fakeRefReturn,
		}
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)
		var recordNS *RecordNS
		var err error
		It("should pass expected NS Record Object to CreateObject", func() {
			recordNS, err = objMgr.CreateNSRecord(name, nameserver, view, addresses, "")

		})
		It("should return expected NS record Object", func() {
			Expect(err).To(BeNil())
			Expect(recordNS).To(Equal(aniFakeConnector.resultObject))
		})
	})

	Describe("Create NS Record : Negative scenario ", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		nameserver := "testing.test.com"
		view := "default"
		name := "test.com"
		conn := &fakeConnector{
			createObjectObj:   NewRecordNS(name, nameserver, view, nil, ""),
			createObjectError: fmt.Errorf("name, nameserver and addresses are required on creation"),
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)
		var actualRecord, expectedObj *RecordNS
		var err error
		expectedObj = nil
		It("should throw error", func() {
			actualRecord, err = objMgr.CreateNSRecord(name, nameserver, view, nil, "")
			Expect(actualRecord).To(Equal(expectedObj))
			Expect(err).To(Equal(conn.createObjectError))
		})
	})
	Describe("Get NS-Record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		nameserver := "testing.test.com"
		view := "default"
		name := "test.com"
		creator := "STATIC"
		addresses := []*ZoneNameServer{
			{
				Address:       "3.4.4.5",
				AutoCreatePtr: true,
			},
		}
		queryParams := NewQueryParams(
			false,
			map[string]string{
				"name":       name,
				"view":       view,
				"nameserver": nameserver,
				"creator":    creator,
			})

		fakeRefReturn := fmt.Sprintf("record:ns/ZG5zLmJpbmRfbnMkLjMuY29tLnRlc3QuLm5hbWUudGVzdC5jb20:%s/%s/%s", nameserver, name, view)
		conn := &fakeConnector{
			createObjectObj:      NewRecordNS(name, nameserver, view, addresses, ""),
			getObjectRef:         "",
			getObjectObj:         NewEmptyRecordNS(),
			resultObject:         []RecordNS{*NewRecordNS(name, nameserver, "default", addresses, "")},
			fakeRefReturn:        fakeRefReturn,
			getObjectQueryParams: queryParams,
		}
		objMgr := NewObjectManager(conn, cmpType, tenantID)

		conn.resultObject.([]RecordNS)[0].Ref = fakeRefReturn

		var actualRecord []RecordNS
		var err error
		It("should pass expected NS Record  Object to GetObject", func() {
			actualRecord, err = objMgr.GetAllRecordNS(queryParams)
			Expect(err).To(BeNil())
			Expect(actualRecord).To(Equal(conn.resultObject))
		})
	})
	Describe("Get NS Record: Negative scenario", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		queryParams2 := NewQueryParams(false, map[string]string{"dns_name": "domain_name"})
		conn := &fakeConnector{
			getObjectObj:         NewEmptyRecordNS(),
			getObjectQueryParams: queryParams2,
			resultObject:         []RecordNS{},
		}
		objMgr := NewObjectManager(conn, cmpType, tenantID)

		// negative scenario
		conn.getObjectError = fmt.Errorf("Field is not searchable: dns_name")
		It("should fail to get expected NS Record Object from getObject with non searchable field", func() {
			_, err := objMgr.GetAllRecordNS(queryParams2)
			Expect(err).ToNot(BeNil())
		})
	})
	Describe("Delete NS record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		nameserver := "testing.test.com"
		view := "default"
		name := "test.com"
		deleteRef := fmt.Sprintf("record:ns/ZG5zLmJpbmRfbnMkLjMuY29tLnRlc3QuLm5hbWUudGVzdC5jb20:%s/%s/%s", nameserver, name, view)
		fakeRefReturn := deleteRef
		nwFakeConnector := &fakeConnector{
			deleteObjectRef: deleteRef,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(nwFakeConnector, cmpType, tenantID)

		var actualRef string
		var err error
		It("should pass expected NS Record Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteNSRecord(deleteRef)
		})
		It("should return expected NS Record Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})
	Describe("Update NS Record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		nameserver := "testing.test.com"
		view := "default"
		name := "test.com"
		addresses := []*ZoneNameServer{
			{
				Address:       "3.4.4.5",
				AutoCreatePtr: true,
			},
		}
		updateRef := fmt.Sprintf("record:ns/ZG5zLmJpbmRfbnMkLjMuY29tLnRlc3QuLm5hbWUudGVzdC5jb20:%s/%s/%s", nameserver, name, view)

		resultObj := NewRecordNS(name, nameserver, view, addresses, "")
		resultObj.Ref = updateRef
		expectedObj := NewRecordNS(name, nameserver, view, addresses, "")
		expectedObj.Ref = updateRef
		conn := &fakeConnector{
			getObjectObj:         NewEmptyRecordNS(),
			getObjectRef:         updateRef,
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         resultObj,
			fakeRefReturn:        updateRef,
			updateObjectObj:      expectedObj,
			updateObjectRef:      updateRef,
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)
		It("should pass expected NS Record Object to UpdateObject", func() {
			actualRecord, err := objMgr.UpdateNSRecord(updateRef, name, nameserver, view, addresses, "")
			Expect(actualRecord).To(Equal(conn.resultObject))
			Expect(err).To(BeNil())
		})
	})

})
