package ibclient

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Object Manager: SVCB Record", func() {
	Describe("Create SVCB Record with maximum params", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		name := "svcb1111.info.com"
		comment := "test SVCB Record"
		targetName := "target.info.com"
		priority := uint32(100)
		disable := false
		creator := "DYNAMIC"
		ddnsPrincipal := "test-principal"
		ddndProtected := true
		view := "test-view"
		ea := EA{"Site": "Kawasaki"}
		useTtl := true
		ttl := uint32(120)
		ref := ""
		svcParams := []SVCParams{
			{
				SvcKey:    "ipv4hint",
				SvcValue:  []string{"11.11.23.11", "11.11.23.12"},
				Mandatory: true,
			},
			{
				SvcKey:    "ipv6hint",
				SvcValue:  []string{"2001:db8::1", "2001:db8::2"},
				Mandatory: false,
			},
			{
				SvcKey:    "port",
				SvcValue:  []string{"443"},
				Mandatory: true,
			},
		}
		fakeRefReturn := fmt.Sprintf("record:svcb/ZG5zLmhvc3QkLZhd3QuaDE:%s", name)

		conn := &fakeConnector{
			createObjectObj:      NewSVCBRecord(ref, name, comment, disable, ea, priority, svcParams, targetName, useTtl, ttl, creator, ddnsPrincipal, ddndProtected),
			getObjectObj:         &RecordSVCB{},
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         NewSVCBRecord(ref, name, comment, disable, ea, priority, svcParams, targetName, useTtl, ttl, creator, ddnsPrincipal, ddndProtected),
			fakeRefReturn:        fakeRefReturn,
		}
		conn.resultObject.(*RecordSVCB).Ref = fakeRefReturn
		conn.createObjectObj.(*RecordSVCB).View = view
		conn.resultObject.(*RecordSVCB).View = view
		objMgr := NewObjectManager(conn, cmpType, tenantID)
		It("should pass expected SVCB Record Object to CreateObject", func() {
			actualRecord, err := objMgr.CreateSVCBRecord(name, comment, disable, ea, priority, svcParams, targetName, useTtl, ttl, view, creator, ddnsPrincipal, ddndProtected)
			Expect(actualRecord).To(Equal(conn.resultObject))
			Expect(err).To(BeNil())
		})

		// Negative scenario
		It("should fail to create a SVCB Record object", func() {
			actualRecord, err := objMgr.CreateSVCBRecord("", comment, disable, ea, priority, svcParams, "", useTtl, ttl, view, creator, ddnsPrincipal, ddndProtected)
			Expect(actualRecord).To(BeNil())
			Expect(err).ToNot(BeNil())
		})
	})

	Describe("Create SVCB Record with minimum params", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		name := "svcb1112.info.com"
		targetName := "target.info.com"
		priority := uint32(100)
		view := ""
		ref := ""
		fakeRefReturn := fmt.Sprintf("record:svcb/ZG5zLmhvc3QkLZhd3QuaDE:%s", name)

		conn := &fakeConnector{
			createObjectObj:      NewSVCBRecord(ref, name, "", false, nil, priority, nil, targetName, false, 0, "", "", false),
			getObjectObj:         &RecordSVCB{},
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         NewSVCBRecord(ref, name, "", false, nil, priority, nil, targetName, false, 0, "", "", false),
			fakeRefReturn:        fakeRefReturn,
		}
		conn.resultObject.(*RecordSVCB).Ref = fakeRefReturn
		objMgr := NewObjectManager(conn, cmpType, tenantID)
		It("should pass expected SVCB Record Object to CreateObject", func() {
			actualRecord, err := objMgr.CreateSVCBRecord(name, "", false, nil, priority, nil, targetName, false, 0, view, "", "", false)
			Expect(actualRecord).To(Equal(conn.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get SVCB Record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		name := "svcb1113.info.com"
		targetName := "target.info.com"
		priority := uint32(120)
		ref := ""
		fakeRefReturn := fmt.Sprintf("record:svcb/ZG5zLmhvc3QkLZhd3QuaDE:%s", name)
		res := NewSVCBRecord(ref, name, "", false, nil, priority, nil, targetName, false, 0, "", "", false)

		conn := &fakeConnector{
			getObjectObj:  NewEmptyRecordSVCB(),
			resultObject:  []RecordSVCB{*res},
			fakeRefReturn: fakeRefReturn,
		}
		queryParams := NewQueryParams(false, map[string]string{"name": name})
		objMgr := NewObjectManager(conn, cmpType, tenantID)
		It("should get expected SVCB Record from getObject", func() {
			conn.getObjectQueryParams = queryParams
			actualRecord, err := objMgr.GetAllSVCBRecords(queryParams)
			Expect(actualRecord).To(Equal(conn.resultObject))
			Expect(err).To(BeNil())
		})

		It("should fail to get expected SVCB Record from getObject", func() {
			queryParams1 := NewQueryParams(false, map[string]string{"name": "svcb-record123"})
			conn.getObjectQueryParams = queryParams1
			conn.resultObject = []RecordSVCB{}
			actualRecord, err := objMgr.GetAllSVCBRecords(queryParams1)
			Expect(actualRecord).To(Equal(conn.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete Record SVCB", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		recordName := "record-svcb222"
		deleteRef := fmt.Sprintf("record:svcb/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)
		nwFakeConnector := &fakeConnector{
			deleteObjectRef: deleteRef,
			fakeRefReturn:   deleteRef,
		}

		objMgr := NewObjectManager(nwFakeConnector, cmpType, tenantID)

		var actualRef string
		var err error
		It("should pass expected SVCB Record Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteSVCBRecord(deleteRef)
		})
		It("should return expected SVCB Record Ref", func() {
			Expect(actualRef).To(Equal(deleteRef))
			Expect(err).To(BeNil())
		})

		It("should pass expected SVCB Record Ref to DeleteObject", func() {
			deleteRef2 := "svcb-777"
			nwFakeConnector.deleteObjectRef = deleteRef2
			nwFakeConnector.fakeRefReturn = ""
			nwFakeConnector.deleteObjectError = fmt.Errorf("not found")
			actualRef, err = objMgr.DeleteSVCBRecord(deleteRef2)
		})

		It("should return an error", func() {
			Expect(err).ToNot(BeNil())
		})

	})

	Describe("Update SVCB Record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		name := "svcb7777.info.com"
		comment := "test SVCB Record updated"
		targetName := "target12.info.com"
		priority := uint32(1010)
		disable := true
		creator := "STATIC"
		ddnsPrincipal := "test-principal"
		ddndProtected := true
		ea := EA{"Site": "SPAIN"}
		useTtl := true
		ttl := uint32(200)
		svcParams := []SVCParams{
			{
				SvcKey:    "ipv4hint",
				SvcValue:  []string{"21.11.23.11", "11.12.23.12"},
				Mandatory: true,
			},
			{
				SvcKey:    "ipv6hint",
				SvcValue:  []string{"2001:db8::1", "2001:db8::2"},
				Mandatory: false,
			},
			{
				SvcKey:    "port",
				SvcValue:  []string{"43"},
				Mandatory: false,
			},
		}
		updateRef := fmt.Sprintf("record:svcb/ZG5zLmhvc3QkLZhd3QuaDE:%s", name)

		conn := &fakeConnector{
			getObjectObj:         NewEmptyRecordSVCB(),
			getObjectRef:         updateRef,
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         NewSVCBRecord("", name, comment, disable, ea, priority, svcParams, targetName, useTtl, ttl, creator, ddnsPrincipal, ddndProtected),
			fakeRefReturn:        updateRef,
			updateObjectObj:      NewSVCBRecord(updateRef, name, comment, disable, ea, priority, svcParams, targetName, useTtl, ttl, creator, ddnsPrincipal, ddndProtected),
			updateObjectRef:      updateRef,
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)
		It("should pass expected SVCB Record Object to UpdateObject", func() {
			actualRecord, err := objMgr.UpdateSVCBRecord(updateRef, name, comment, disable, ea, priority, svcParams, targetName, useTtl, ttl, creator, ddnsPrincipal, ddndProtected)
			Expect(actualRecord).To(Equal(conn.resultObject))
			Expect(err).To(BeNil())
		})

	})

	Describe("Update SVCB Record with, negative scenario", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		name := "svcb7778.info.com"
		comment := "test SVCB Record updated"
		targetName := "target128.info.com"
		priority := uint32(110)
		disable := true
		creator := "STATIC"
		ddnsPrincipal := "test-principal"
		ddndProtected := true
		ea := EA{"Site": "SPAIN"}
		useTtl := true
		ttl := uint32(800)
		svcParams := []SVCParams{
			{
				SvcKey:    "ipv4hint",
				SvcValue:  []string{"21.11.23.11", "11.12.23.12"},
				Mandatory: true,
			},
			{
				SvcKey:    "ipv6hint",
				SvcValue:  []string{"2001:db8::1", "2001:db8::2"},
				Mandatory: false,
			},
			{
				SvcKey:    "port",
				SvcValue:  []string{"43"},
				Mandatory: false,
			},
		}
		oldRef := fmt.Sprintf("record:svcb/ZG5zLmhvc3QkLZhd3QuaDE:%s", name)

		conn := &fakeConnector{
			getObjectObj:         NewEmptyRecordSVCB(),
			getObjectRef:         oldRef,
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         NewSVCBRecord(oldRef, name, comment, disable, ea, priority, svcParams, targetName, useTtl, ttl, creator, ddnsPrincipal, ddndProtected),
			getObjectError:       fmt.Errorf("not found"),
			fakeRefReturn:        oldRef,
			updateObjectObj:      NewSVCBRecord(oldRef, name, comment, disable, ea, priority, svcParams, targetName, useTtl, ttl, creator, ddnsPrincipal, ddndProtected),
			updateObjectRef:      oldRef,
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)
		// negative scenario

		It("should fail to update SVCB Record Object", func() {
			actualRecord, err := objMgr.UpdateSVCBRecord(oldRef, name, comment, disable, ea, priority, svcParams, targetName, useTtl, ttl, creator, ddnsPrincipal, ddndProtected)
			Expect(actualRecord).To(BeNil())
			Expect(err).ToNot(BeNil())
		})
	})
})
