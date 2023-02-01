package ibclient

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Object Manager: SRV-Record", func() {
	Describe("Create SRV Record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		dnsView := "default"
		fqdn := "srv.example.com"
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		priority := 10
		weight := 25
		port := 88
		target := "h1.example.com"
		ttl := uint32(70)
		useTtl := true
		comment := "this is a test comment"

		fakeRefReturn := fmt.Sprintf("record:srv/ZG5zLmhvc3RjkuMC4xLg:%s/%s", fqdn, dnsView)

		eas := make(EA)
		eas["VM ID"] = vmID
		eas["VM Name"] = vmName

		aniFakeConnector := &fakeConnector{
			createObjectObj: NewRecordSRV(RecordSRV{
				View:     dnsView,
				Fqdn:     fqdn,
				Priority: priority,
				Weight:   weight,
				Port:     port,
				Target:   target,
				Ttl:      ttl,
				UseTtl:   useTtl,
			}),
			getObjectRef: fakeRefReturn,
			getObjectObj: NewRecordSRV(RecordSRV{
				View: dnsView,
				Fqdn: fqdn,
				Ref:  fakeRefReturn,
			}),
			resultObject: NewRecordSRV(RecordSRV{
				View:     dnsView,
				Fqdn:     fqdn,
				Priority: priority,
				Weight:   weight,
				Port:     port,
				Target:   target,
				Ttl:      ttl,
				UseTtl:   useTtl,
				Ref:      fakeRefReturn,
			}),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualRecord *RecordSRV
		var err error
		It("should pass expected SRV record object to CreateObject", func() {
			actualRecord, err = objMgr.CreateSRVRecord(dnsView, fqdn, priority, weight, port, target, ttl, useTtl, comment, eas)
		})
		It("should return expected SRV record Object", func() {
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Update SRV Record", func() {
		var (
			err       error
			objMgr    IBObjectManager
			conn      *fakeConnector
			ref       string
			actualObj *RecordSRV
		)

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		fqdn := "srv.example.com"
		dnsView := "default"
		priority := 10
		weight := 25
		port := 80
		target := "h2.example.com"
		ttl := uint32(400)
		useTtl := true

		It("Updating fqdn, priority, weight, port, target, comment and EA's", func() {
			ref = fmt.Sprintf("record:srv/ZG5zLmhvc3RjkuMC4xLg:%s/%s", fqdn, dnsView)
			initialEas := EA{"Country": "old value"}
			initObj := NewRecordSRV(RecordSRV{
				View:     dnsView,
				Fqdn:     fqdn,
				Priority: priority,
				Weight:   weight,
				Port:     port,
				Target:   target,
				Comment:  "test comment",
				Ea:       initialEas,
			})
			initObj.Ref = ref

			expectedEas := EA{"Country": "new value"}

			updateFqdn := "new.example.com"
			updatePriority := 15
			updateWeight := 30
			updatePort := 88
			updateTarget := "h3.example.com"
			updateComment := "test comment"
			updateRef := fmt.Sprintf("record:srv/ZG5zLmhvc3RjkugC4xLg:%s/%s", fqdn, dnsView)
			updateObjIn := NewRecordSRV(RecordSRV{
				Fqdn:     updateFqdn,
				Priority: updatePriority,
				Weight:   updateWeight,
				Port:     updatePort,
				Target:   updateTarget,
				Comment:  updateComment,
				Ea:       expectedEas,
			})
			updateObjIn.Ref = ref

			expectedObj := NewRecordSRV(RecordSRV{
				Fqdn:     updateFqdn,
				Priority: updatePriority,
				Weight:   updateWeight,
				Port:     updatePort,
				Target:   updateTarget,
				Comment:  updateComment,
				Ea:       expectedEas,
			})
			expectedObj.Ref = updateRef

			conn = &fakeConnector{
				getObjectObj: NewRecordSRV(RecordSRV{
					Fqdn:     fqdn,
					Priority: initObj.Priority,
					Weight:   initObj.Weight,
					Port:     initObj.Port,
					Target:   initObj.Target,
					Ref:      initObj.Ref,
				}),
				getObjectQueryParams: NewQueryParams(false, nil),
				getObjectRef:         updateRef,
				getObjectError:       nil,
				resultObject:         expectedObj,

				updateObjectObj:   updateObjIn,
				updateObjectRef:   ref,
				updateObjectError: nil,

				fakeRefReturn: updateRef,
			}
			objMgr = NewObjectManager(conn, cmpType, tenantID)
			It("should pass updated SRV record arguments", func() {
				actualObj, err = objMgr.UpdateSRVRecord(ref, dnsView, updateFqdn, updatePriority, updateWeight, updatePort, updateTarget, ttl, useTtl, updateComment, expectedEas)
			})
			It("should return expected SRV record obj", func() {
				Expect(err).To(BeNil())
				Expect(actualObj).To(BeEquivalentTo(expectedObj))
			})
		})
	})

	Describe("Delete SRV Record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		fqdn := "srv.example.com"
		dnsView := "default"
		deleteRef := fmt.Sprintf("record:srv/ZG5zLmhvc3RjkugC4xLg:%s/%s", fqdn, dnsView)
		fakeRefReturn := deleteRef
		nwFakeConnector := &fakeConnector{
			deleteObjectRef: deleteRef,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(nwFakeConnector, cmpType, tenantID)

		var actualRef string
		var err error
		It("should pass expected SRV Record ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteSRVRecord(deleteRef)
		})
		It("should return expected SRV Record Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})
})
