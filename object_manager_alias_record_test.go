package ibclient

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Object Manager Record Alias", func() {
	Describe("Create Alias Record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		name := "alias777.test.com"
		comment := "some comment"
		disable := true
		ea := EA{"Site": "LA"}
		ttl := uint32(120)
		useTtl := false

		view := "default"
		targetName := "aa.xx.com"
		targetType := "PTR"
		fakeRefReturn := "record:alias/ZG5zLmFsaWFzX3JlY29yZCQuX2RlZmF1bHQuY29tLnRlc3QuYWxpYXM3NzcuUFRS:alias777.test.com/default"

		objectAsResult := NewAliasRecord(name, view, targetName, targetType, comment, disable, ea, ttl, useTtl)
		objectAsResult.Ref = fakeRefReturn
		conn := &fakeConnector{
			createObjectObj:      NewAliasRecord(name, view, targetName, targetType, comment, disable, ea, ttl, useTtl),
			getObjectRef:         fakeRefReturn,
			getObjectObj:         NewEmptyAliasRecord(),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         objectAsResult,
			fakeRefReturn:        fakeRefReturn,
		}
		objMgr := NewObjectManager(conn, cmpType, tenantID)
		var aliasRecord *RecordAlias
		var err error
		It("should pass expected Alias record to CreateObject", func() {
			aliasRecord, err = objMgr.CreateAliasRecord(name, view, targetName, targetType, comment, disable, ea, ttl, useTtl)

		})
		It("should return expected Alias record", func() {
			Expect(err).To(BeNil())
			Expect(aliasRecord).To(Equal(conn.resultObject))
		})
	})

	Describe("Create Alias Record, negative scenario", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		name := ""
		comment := "some comment"
		disable := true
		ea := EA{"Site": "LA"}
		ttl := uint32(120)
		useTtl := false

		view := "default"
		targetName := "aa.xx.com"
		targetType := "PTR"
		fakeRefReturn := "record:alias/ZG5zLmFsaWFzX3JlY29yZCQuX2RlZmF1bHQuY29tLnRlc3QuYWxpYXM3NzcuUFRS:alias777.test.com/default"

		objectAsResult := NewAliasRecord(name, view, targetName, targetType, comment, disable, ea, ttl, useTtl)
		objectAsResult.Ref = fakeRefReturn
		conn := &fakeConnector{
			createObjectObj:      NewAliasRecord(name, view, targetName, targetType, comment, disable, ea, ttl, useTtl),
			getObjectRef:         fakeRefReturn,
			getObjectObj:         NewEmptyAliasRecord(),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         objectAsResult,
			fakeRefReturn:        fakeRefReturn,
			createObjectError:    fmt.Errorf("name, targetName and targetType are required to create an Alias Record"),
		}
		objMgr := NewObjectManager(conn, cmpType, tenantID)
		var aliasRecord *RecordAlias
		var err error
		It("should pass expected Alias record to CreateObject", func() {
			aliasRecord, err = objMgr.CreateAliasRecord(name, view, targetName, targetType, comment, disable, ea, ttl, useTtl)

		})
		It("should return expected Alias record", func() {
			Expect(err).To(Equal(conn.createObjectError))
			Expect(aliasRecord).To(BeNil())
		})
	})

	Describe("Get Alias Record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		name := "alias777.test.com"
		comment := "some comment"
		disable := true
		ea := EA{"Site": "LA"}
		ttl := uint32(120)
		useTtl := false

		view := "default"
		targetName := "aa.xx.com"
		targetType := "PTR"
		fakeRefReturn := "record:alias/ZG5zLmFsaWFzX3JlY29yZCQuX2RlZmF1bHQuY29tLnRlc3QuYWxpYXM3NzcuUFRS:alias777.test.com/default"
		queryParams := NewQueryParams(false, map[string]string{"name": name})

		res := NewAliasRecord(name, view, targetName, targetType, comment, disable, ea, ttl, useTtl)
		res.Ref = fakeRefReturn
		conn := &fakeConnector{
			getObjectObj:  NewEmptyAliasRecord(),
			resultObject:  []RecordAlias{*res},
			fakeRefReturn: fakeRefReturn,
		}
		objMgr := NewObjectManager(conn, cmpType, tenantID)
		//var aliasRecord []RecordAlias
		It("should get expected Alias record from GetObject", func() {
			conn.getObjectQueryParams = queryParams
			actualRecord, err := objMgr.GetAllAliasRecord(queryParams)
			Expect(err).To(BeNil())
			Expect(actualRecord).To(Equal(conn.resultObject))

		})
		It("should fail to get expected Alias record from getObject", func() {
			queryParams1 := NewQueryParams(false, map[string]string{"name": "alias999"})
			conn.getObjectQueryParams = queryParams1
			conn.resultObject = []RecordAlias{}
			actualRecord, err := objMgr.GetAllAliasRecord(queryParams1)
			Expect(actualRecord).To(Equal(conn.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get Alias Record: Negative scenario", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		queryParams2 := NewQueryParams(false, map[string]string{"creator": "STATIC"})
		conn := &fakeConnector{
			getObjectObj:         NewEmptyAliasRecord(),
			getObjectQueryParams: queryParams2,
			resultObject:         []RecordAlias{},
		}
		objMgr := NewObjectManager(conn, cmpType, tenantID)

		// negative scenario
		conn.getObjectError = fmt.Errorf("Field is not searchable: creator")
		It("should fail to get expected Alias Record from getObject with non searchable field", func() {
			_, err := objMgr.GetAllAliasRecord(queryParams2)
			Expect(err).ToNot(BeNil())
		})
	})

	Describe("Delete Alias Record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		recordName := "alias111.test.com"
		deleteRef := fmt.Sprintf("record:alias/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)
		//deleteRef := "record:alias/ZG5zLmFsaWFzX3JlY29yZCQuX2RlZmF1bHQuY29tLnRlc3QuYWxpYXM3NzcuUFRS:alias111.test.com/default"
		nwFakeConnector := &fakeConnector{
			deleteObjectRef: deleteRef,
			fakeRefReturn:   deleteRef,
		}
		objMgr := NewObjectManager(nwFakeConnector, cmpType, tenantID)

		var actualRef string
		var err error
		It("should pass expected Alias Record Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteAliasRecord(deleteRef)
		})
		It("should return expected Alias Record Ref", func() {
			Expect(actualRef).To(Equal(deleteRef))
			Expect(err).To(BeNil())
		})

		// negative scenario
		It("should pass expected Alias Record Ref to DeleteObject", func() {
			deleteRef2 := "record:alias/hjwdjdklwhlwijflwkjf"
			nwFakeConnector.deleteObjectRef = deleteRef2
			nwFakeConnector.fakeRefReturn = ""
			nwFakeConnector.deleteObjectError = fmt.Errorf("not found")
			actualRef, err = objMgr.DeleteAliasRecord(deleteRef2)
		})
		It("should return an error", func() {
			Expect(err).ToNot(BeNil())
		})
	})

	Describe("Update Alias record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		name := "alias777.test.com"
		comment := "some comment"
		disable := true
		ea := EA{"Site": "LA"}
		ttl := uint32(120)
		useTtl := false

		view := "default"
		targetName := "aa.xx.com"
		targetType := "PTR"
		updateRef := "record:alias/ZG5zLmFsaWFzX3JlY29yZCQuX2RlZmF1bHQuY29tLnRlc3QuYWxpYXM3NzcuUFRS:alias111.test.com/default"

		conn := &fakeConnector{
			getObjectObj:         NewEmptyAliasRecord(),
			getObjectRef:         updateRef,
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         NewAliasRecord(name, view, targetName, targetType, comment, disable, ea, ttl, useTtl),
			fakeRefReturn:        updateRef,
			updateObjectObj:      NewAliasRecord(name, view, targetName, targetType, comment, disable, ea, ttl, useTtl),
			updateObjectRef:      updateRef,
		}
		conn.resultObject.(*RecordAlias).Ref = updateRef

		objMgr := NewObjectManager(conn, cmpType, tenantID)
		It("should pass expected Alias Record to UpdateObject", func() {
			actualRecord, err := objMgr.UpdateAliasRecord(updateRef, name, view, targetName, targetType, comment, disable, ea, ttl, useTtl)
			Expect(actualRecord).To(Equal(conn.resultObject))
			Expect(err).To(BeNil())
		})

	})

	Describe("Update Alias record: Negative scenario", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		name := "alias777.test.com"
		comment := "some comment"
		disable := true
		ea := EA{"Site": "LA"}
		ttl := uint32(120)
		useTtl := false

		view := "default"
		targetName := "aa.xx.com"
		targetType := "PTR"

		conn := &fakeConnector{
			getObjectObj:         NewEmptyAliasRecord(),
			getObjectError:       fmt.Errorf("not found"),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         NewAliasRecord(name, view, targetName, targetType, comment, disable, ea, ttl, useTtl),
			updateObjectObj:      NewAliasRecord(name, view, targetName, targetType, comment, disable, ea, ttl, useTtl),
			updateObjectError:    fmt.Errorf("not found"),
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)
		// negative scenario
		It("should pass expected Alias Record to UpdateObject, negative scenario", func() {
			actualRecord, err := objMgr.UpdateAliasRecord("", name, view, targetName, targetType, comment, disable, ea, ttl, useTtl)
			Expect(err).NotTo(BeNil())
			Expect(actualRecord).To(BeNil())
		})
	})

})
