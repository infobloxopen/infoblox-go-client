package ibclient

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Object Manager: MX-record", func() {
	Describe("Create MX record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		mx := "example.com"
		dnsView := "default"
		fqdn := "test.example.com"
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		priority := uint32(10)
		ttl := uint32(70)
		useTtl := true
		comment := "test comment"

		fakeRefReturn := fmt.Sprintf("record:mx/ZG5zLmhvc3RjkuMC4xLg:%s/%s", fqdn, dnsView)

		eas := make(EA)
		eas["VM ID"] = vmID
		eas["VM Name"] = vmName

		aniFakeConnector := &fakeConnector{
			createObjectObj: NewRecordMX(RecordMX{
				View:     dnsView,
				Fqdn:     fqdn,
				MX:       mx,
				Priority: priority,
				Ttl:      ttl,
				UseTtl:   useTtl,
			}),
			getObjectRef: fakeRefReturn,
			getObjectObj: NewRecordMX(RecordMX{
				Fqdn:     fqdn,
				MX:       mx,
				Priority: priority,
				Ref:      fakeRefReturn,
			}),
			resultObject: NewRecordMX(RecordMX{
				View:     dnsView,
				Fqdn:     fqdn,
				MX:       mx,
				Priority: priority,
				Ttl:      ttl,
				UseTtl:   useTtl,
				Ref:      fakeRefReturn,
			}),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualRecord *RecordMX
		var err error
		It("should pass expected MX record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateMXRecord(dnsView, fqdn, mx, priority, ttl, useTtl, comment, eas)
		})
		It("should return expected MX record Object", func() {
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Update MX Record", func() {
		var (
			err       error
			objMgr    IBObjectManager
			conn      *fakeConnector
			ref       string
			actualObj *RecordMX
		)

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		fqdn := "test.example.com"
		dnsView := "default"
		ttl := uint32(70)
		useTtl := true

		It("Updating fqdn, comment, priority and EAs", func() {
			ref = fmt.Sprintf("record:mx/ZG5zLmhvc3RjkuMC4xLg:%s/%s", fqdn, dnsView)
			initialEas := EA{"Country": "old value"}
			initObj := NewRecordMX(RecordMX{
				View:     dnsView,
				Fqdn:     fqdn,
				Priority: uint32(10),
				Comment:  "test comment",
				Ttl:      ttl,
				UseTtl:   useTtl,
				Ea:       initialEas,
			})
			initObj.Ref = ref

			expectedEas := EA{"Country": "new value"}

			updateFqdn := "new.example.com"
			updateComment := "new comment"
			updateTtl := uint32(100)
			updatePriority := uint32(15)
			updatedRef := fmt.Sprintf("record:mx/ZG5zLmhvc3RjkuMC4xLg:%s/%s", fqdn, dnsView)
			updateObjIn := NewRecordMX(RecordMX{
				Fqdn:     updateFqdn,
				Priority: updatePriority,
				Comment:  updateComment,
				Ttl:      ttl,
				UseTtl:   useTtl,
				Ea:       expectedEas,
			})
			updateObjIn.Ref = ref

			expectedObj := NewRecordMX(RecordMX{
				Fqdn:     updateFqdn,
				Priority: updatePriority,
				Comment:  updateComment,
				Ttl:      ttl,
				UseTtl:   useTtl,
				Ea:       expectedEas,
			})
			expectedObj.Ref = updatedRef

			conn = &fakeConnector{
				getObjectObj: NewRecordMX(RecordMX{
					Fqdn:     fqdn,
					MX:       actualObj.MX,
					Priority: initObj.Priority,
					Ref:      initObj.Ref,
				}),
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
			It("should pass updated MX record arguments", func() {
				actualObj, err = objMgr.UpdateMXRecord(ref, dnsView, updateFqdn, "", updateTtl, useTtl, updateComment, updatePriority, expectedEas)
			})
			It("should return expected MX record obj", func() {
				Expect(err).To(BeNil())
				Expect(actualObj).To(BeEquivalentTo(expectedObj))
			})

		})
	})

	Describe("Get MX Record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		fqdn := "test.example.com"
		dnsView := "default"
		mx := "example.com"
		priority := uint32(25)
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		ttl := uint32(70)
		useTtl := true
		comment := "test comment"

		fakeRefReturn := fmt.Sprintf("record:mx/ZG5zLmhvc3RjkuMC4xLg:%s/%s", fqdn, dnsView)

		eas := make(EA)
		eas["VM ID"] = vmID
		eas["VM Name"] = vmName

		nwFakeConnector := &fakeConnector{
			getObjectObj: NewRecordMX(RecordMX{
				Fqdn: fqdn,
				View: dnsView,
			}),
			resultObject: NewRecordMX(RecordMX{
				View:     dnsView,
				Fqdn:     fqdn,
				MX:       mx,
				Priority: priority,
				Ttl:      ttl,
				UseTtl:   useTtl,
				Comment:  comment,
				Ea:       eas,
				Ref:      fakeRefReturn,
			}),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(nwFakeConnector, cmpType, tenantID)

		var actualRecord *RecordMX
		var err error
		It("should pass expected MX record object to GetObject", func() {
			actualRecord, err = objMgr.GetMXRecord(dnsView, fqdn)
		})
		It("should return expected MX record Object", func() {
			Expect(actualRecord).To(Equal(nwFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get MX Record By Ref", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		fqdn := "test.example.com"
		dnsView := "default"
		mx := "example.com"
		priority := uint32(25)
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		ttl := uint32(70)
		useTtl := true
		comment := "test comment"

		readobjRef := fmt.Sprintf("record:mx/ZG5zLmhvc3RjkuMC4xLg:%s/%s", fqdn, dnsView)

		eas := make(EA)
		eas["VM ID"] = vmID
		eas["VM Name"] = vmName

		nwFakeConnector := &fakeConnector{
			getObjectRef: readobjRef,
			resultObject: NewRecordMX(RecordMX{
				View:     dnsView,
				Fqdn:     fqdn,
				MX:       mx,
				Priority: priority,
				Ttl:      ttl,
				UseTtl:   useTtl,
				Comment:  comment,
				Ea:       eas,
				Ref:      readobjRef,
			}),
			fakeRefReturn: readobjRef,
		}

		objMgr := NewObjectManager(nwFakeConnector, cmpType, tenantID)

		var actualRecord *RecordMX
		var err error
		It("should pass expected MX record ref to GetObject", func() {
			actualRecord, err = objMgr.GetMXRecordByRef(readobjRef)
		})
		It("should return expected MX record Object", func() {
			Expect(actualRecord).To(Equal(nwFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete MX Record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		fqdn := "test.example.com"
		dnsView := "default"
		deleteRef := fmt.Sprintf("record:mx/ZG5zLmhvc3RjkuMC4xLg:%s/%s", fqdn, dnsView)
		fakeRefReturn := deleteRef
		nwFakeConnector := &fakeConnector{
			deleteObjectRef: deleteRef,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(nwFakeConnector, cmpType, tenantID)

		var actualRef string
		var err error
		It("should pass expected MX record Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteMXRecord(deleteRef)
		})
		It("should return expected MX record Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})

	})
})
