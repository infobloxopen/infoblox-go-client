package ibclient

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Object Manager: HTTPS-record", func() {
	Describe("Create a specific HTTPS Record with minimal parameters", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		name := "a7.test.com"
		priority := uint32(30)
		targetName := "test.com"
		dnsView := "default"
		creator := "STATIC"
		fakeRefReturn := fmt.Sprintf("record:https/ZG5zLmJpbmRfaHR0cHMkLl9kZWZhdWx0LmNvbS50ZXN0LmE3LjIwLnRlc3QuY29t:%v/%v", name, dnsView)
		objectForCreation := NewHttpsRecord(name, priority, targetName, "", creator, "", false, nil, false, nil, false, 0, false, dnsView, "")
		objectForResult := NewHttpsRecord(name, priority, targetName, "", creator, "", false, nil, false, nil, false, 0, false, dnsView, fakeRefReturn)
		objectForResult.Ref = fakeRefReturn
		aniFakeConnector := &fakeConnector{
			createObjectObj:      objectForCreation,
			getObjectRef:         fakeRefReturn,
			getObjectObj:         NewEmptyHttpsRecord(),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         objectForResult,
			fakeRefReturn:        fakeRefReturn,
		}
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)
		var actualRecord *RecordHttps
		var err error
		It("should pass expected HTTPS record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateHTTPSRecord(name, priority, targetName, "", creator, "", false, false, nil, false, nil, 0, false, dnsView)
		})
		It("should return expected HTTPS record Object", func() {
			Expect(err).To(BeNil())
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
		})
	})

	Describe("Create a specific HTTPS Record with all parameters", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		name := "a7.test.com"
		priority := uint32(30)
		targetName := "test.com"
		dnsView := "default"
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		eas := make(EA)
		eas["VM ID"] = vmID
		eas["VM Name"] = vmName
		comment := "test comment"
		ttl := uint32(40)
		useTtl := true
		fakeRefReturn := fmt.Sprintf("record:https/ZG5zLmJpbmRfaHR0cHMkLl9kZWZhdWx0LmNvbS50ZXN0LmE3LjIwLnRlc3QuY29t:%v/%v", name, dnsView)
		svcParams := []SVCParams{
			{
				Mandatory: true,
				SvcKey:    "port",
				SvcValue:  []string{"4454"},
			},
		}
		disable := true
		forbidReclamation := true
		creator := "STATIC"
		objectForCreation := NewHttpsRecord(name, priority, targetName, comment, creator, "", false, svcParams, disable, eas, forbidReclamation, ttl, useTtl, dnsView, "")
		objectForResult := NewHttpsRecord(name, priority, targetName, comment, creator, "", false, svcParams, disable, eas, forbidReclamation, ttl, useTtl, dnsView, fakeRefReturn)
		objectForResult.Ref = fakeRefReturn

		aniFakeConnector := &fakeConnector{
			createObjectObj:      objectForCreation,
			getObjectRef:         fakeRefReturn,
			getObjectObj:         NewEmptyHttpsRecord(),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         objectForResult,
			fakeRefReturn:        fakeRefReturn,
		}
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)
		var actualRecord *RecordHttps
		var err error
		It("should pass expected HTTPS record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateHTTPSRecord(name, priority, targetName, comment, creator, "", false, disable, eas, forbidReclamation, svcParams, ttl, useTtl, dnsView)
		})
		It("should return expected HTTPS record Object", func() {
			Expect(err).To(BeNil())
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
		})
	})

	Describe("Get All HTTPS Record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		comment := "get https record"
		name := "a7.test.com"
		dnsView := "default"
		targetName := "test.com"
		fakeRefReturn := fmt.Sprintf("record:https/ZG5zLmJpbmRfaHR0cHMkLl9kZWZhdWx0LmNvbS50ZXN0LmE3LjIwLnRlc3QuY29t:%v/%v", name, dnsView)

		queryParams := NewQueryParams(
			false,
			map[string]string{
				"name":    name,
				"comment": comment,
			})
		creator := "STATIC"
		conn := &fakeConnector{
			createObjectObj: NewHttpsRecord(name, 20, targetName, comment, creator, "", false, nil, false, nil, false, 20, true, dnsView, ""),
			getObjectRef:         "",
			getObjectObj:         NewEmptyHttpsRecord(),
			resultObject:         []RecordHttps{*NewHttpsRecord(name, 20, targetName, comment, creator, "", false, nil, false, nil, false, 20, true, dnsView, "")},
			fakeRefReturn:        fakeRefReturn,
			getObjectQueryParams: queryParams,
		}
		objMgr := NewObjectManager(conn, cmpType, tenantID)

		conn.resultObject.([]RecordHttps)[0].Ref = fakeRefReturn

		var actualRecord []RecordHttps
		var err error
		It("should pass expected HTTPS record Object to GetObject", func() {
			actualRecord, err = objMgr.GetAllHTTPSRecord(queryParams)
			Expect(err).To(BeNil())
			Expect(actualRecord).To(Equal(conn.resultObject.([]RecordHttps)))
		})
	})
	Describe("Delete Https Record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		name := "a7.test.com"
		dnsView := "default"
		deleteRef := fmt.Sprintf("record:https/ZG5zLmJpbmRfaHR0cHMkLl9kZWZhdWx0LmNvbS50ZXN0LmE3LjIwLnRlc3QuY29t:%s/%s", name, dnsView)
		nwFakeConnector := &fakeConnector{
			deleteObjectRef: deleteRef,
			fakeRefReturn:   deleteRef,
		}

		objMgr := NewObjectManager(nwFakeConnector, cmpType, tenantID)

		var actualRef string
		var err error
		It("should pass expected Https Record Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteHTTPSRecord(deleteRef)
		})
		It("should return expected Https Record Ref", func() {
			Expect(actualRef).To(Equal(deleteRef))
			Expect(err).To(BeNil())
		})

		It("should pass expected Https Record Ref to DeleteObject", func() {
			deleteRef2 := "httpsrecord"
			nwFakeConnector.deleteObjectRef = deleteRef2
			nwFakeConnector.fakeRefReturn = ""
			nwFakeConnector.deleteObjectError = fmt.Errorf("not found")
			actualRef, err = objMgr.DeleteHTTPSRecord(deleteRef2)
		})

		It("should return an error", func() {
			Expect(err).ToNot(BeNil())
		})

	})

	Describe("Negative case : Create Https record without all required parameters", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		name := "a7.test.com"
		dnsView := "default"
		creator := "STATIC"
		fakeRefReturn := fmt.Sprintf("record:https/ZG5zLmJpbmRfaHR0cHMkLl9kZWZhdWx0LmNvbS50ZXN0LmE3LjIwLnRlc3QuY29t:%s/%s", name, dnsView)
		aniFakeConnector := &fakeConnector{
			createObjectObj: NewHttpsRecord(name, 0, "", "", creator, "", false, nil, false, nil, false, 0, false, dnsView, ""),
			getObjectRef:         "",
			getObjectObj:         NewEmptyHttpsRecord(),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         NewHttpsRecord(name, 0, "", "", creator, "", false, nil, false, nil, false, 0, false, dnsView, fakeRefReturn),
			createObjectError:    fmt.Errorf("name and targetName are required to create HTTPS Record"),
		}

		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)
		var err error
		It("should return an error when creating HTTPS record", func() {
			_, err = objMgr.CreateHTTPSRecord(name, 0, "", "", creator, "", false, false, nil, false, nil, 0, false, dnsView)
			Expect(err).ToNot(BeNil())
			Expect(err).To(Equal(aniFakeConnector.createObjectError))
		})
	})

	Describe("Update Https Record", func() {
		var (
			err       error
			objMgr    IBObjectManager
			conn      *fakeConnector
			ref       string
			actualObj *RecordHttps
		)

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		name := "a7.test.com"
		priority := uint32(30)
		targetName := "test.com"
		dnsView := "default"
		ttl := uint32(20)
		useTtl := true
		disable := false
		creator := "STATIC"
		ref = fmt.Sprintf("record:https/ZG5zLmJpbmRfaHR0cHMkLl9kZWZhdWx0LmNvbS50ZXN0LmE3LjIwLnRlc3QuY29t:%s/%s", name, dnsView)
		It("Adding SVC Params to the created HTTPS record during Update ", func() {
			ref = fmt.Sprintf("record:https/ZG5zLmJpbmRfaHR0cHMkLl9kZWZhdWx0LmNvbS50ZXN0LmE3LjIwLnRlc3QuY29t:%s/%s", name, dnsView)
			initialEas := EA{"Site": "Blr"}
			initialComment := "old comment"
			initObj := NewHttpsRecord(name, priority, targetName, initialComment, creator, "", false, nil, false, initialEas, false, ttl, useTtl, dnsView, "")
			initObj.Ref = ref

			expectedEas := EA{"Site": "Blr"}

			updateName := "a8.test.com"
			updateComment := "new comment"
			svcParams := []SVCParams{
				{
					Mandatory: true,
					SvcKey:    "port",
					SvcValue:  []string{"4454"},
				},
			}
			updatedRef := fmt.Sprintf("record:https/ZG5zLmJpbmRfaHR0cHMkLl9kZWZhdWx0LmNvbS50ZXN0LmE3LjIwLnRlc3QuY29t:%s/%s", updateName, dnsView)
			updateObjIn := NewHttpsRecord(updateName, priority, targetName, updateComment, creator, "", false, svcParams, false, expectedEas, false, ttl, useTtl, "", ref)
			updateObjIn.Ref = ref

			expectedObj := NewHttpsRecord(updateName, priority, targetName, updateComment, creator, "", false, svcParams, false, expectedEas, false, ttl, useTtl, "", ref)

			expectedObj.Ref = updatedRef

			getObjIn := NewEmptyHttpsRecord()
			conn = &fakeConnector{
				getObjectObj:         getObjIn,
				getObjectQueryParams: NewQueryParams(false, nil),
				getObjectRef:         updatedRef,
				getObjectError:       nil,
				resultObject:         expectedObj,
				updateObjectObj:      updateObjIn,
				updateObjectRef:      ref,
				updateObjectError:    nil,

				fakeRefReturn: updatedRef,
			}
			objMgr = NewObjectManager(conn, cmpType, tenantID)

			actualObj, err = objMgr.UpdateHTTPSRecord(ref, updateName, priority, targetName, updateComment, creator, "", false, disable, expectedEas, false, svcParams, ttl, useTtl)
			Expect(err).To(BeNil())
			Expect(*actualObj).To(BeEquivalentTo(*expectedObj))
		})
	})
	Describe("Get Https Record: Negative scenario", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		queryParams2 := NewQueryParams(false, map[string]string{"disable": "true"})
		conn := &fakeConnector{
			getObjectObj:         NewEmptyHttpsRecord(),
			getObjectQueryParams: queryParams2,
			resultObject:         []RecordHttps{},
		}
		objMgr := NewObjectManager(conn, cmpType, tenantID)

		// negative scenario
		conn.getObjectError = fmt.Errorf("Field is not searchable: disable")
		It("should fail to get expected HTTPS Record from getObject with non searchable field", func() {
			_, err := objMgr.GetAllHTTPSRecord(queryParams2)
			Expect(err).ToNot(BeNil())
		})
	})
})
