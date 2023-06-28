package ibclient

import (
	"fmt"
	"github.com/infobloxopen/infoblox-go-client/v2/utils"

	. "github.com/onsi/ginkgo/v2"
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
		preference := uint32(10)
		ttl := uint32(70)
		useTtl := true
		comment := "test comment"

		fakeRefReturn := fmt.Sprintf("record:mx/ZG5zLmhvc3RjkuMC4xLg:%s/%s", fqdn, dnsView)

		eas := make(EA)
		eas["VM ID"] = vmID
		eas["VM Name"] = vmName

		aniFakeConnector := &fakeConnector{
			createObjectObj: NewRecordMX(RecordMX{
				View:          &dnsView,
				Name:          &fqdn,
				MailExchanger: &mx,
				Preference:    &preference,
				Ttl:           &ttl,
				UseTtl:        &useTtl,
				Comment:       &comment,
				Ea:            eas,
			}),
			getObjectRef: fakeRefReturn,
			getObjectObj: NewRecordMX(RecordMX{
				Name:          &fqdn,
				MailExchanger: &mx,
				Preference:    &preference,
				Ref:           fakeRefReturn,
			}),
			resultObject: NewRecordMX(RecordMX{
				View:          &dnsView,
				Name:          &fqdn,
				MailExchanger: &mx,
				Preference:    &preference,
				Ttl:           &ttl,
				UseTtl:        &useTtl,
				Ref:           fakeRefReturn,
				Comment:       &comment,
				Ea:            eas,
			}),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualRecord *RecordMX
		var err error
		It("should pass expected MX record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateMXRecord(dnsView, fqdn, mx, preference, ttl, useTtl, comment, eas)
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
		dnsView := "default"

		fqdn := "test.example.com"
		initMx := "mx.test.example.com"
		initComment := "test comment"
		ref = fmt.Sprintf("record:mx/ZG5zLmhvc3RjkuMC4xLg:%s/%s", fqdn, dnsView)
		initialEas := EA{"Country": "old value"}

		initObj := NewRecordMX(RecordMX{
			Ref:           ref,
			View:          &dnsView,
			Name:          &fqdn,
			MailExchanger: &initMx,
			Preference:    utils.Uint32Ptr(10),
			Ttl:           utils.Uint32Ptr(70),
			UseTtl:        utils.BoolPtr(true),
			Comment:       &initComment,
			Ea:            initialEas,
		})

		updatedEAs := EA{"Country": "new value"}
		updatedFqdn := "new.example.com"
		updatedMx := "mx.new.example.com"
		updatedComment := "new comment"
		updatedTtl := uint32(100)
		updatedPreference := uint32(15)
		updatedRef := fmt.Sprintf("record:mx/ZG5zLmhvc3RjkuMC4xLg:%s/%s", fqdn, dnsView)
		updateObjIn := NewRecordMX(RecordMX{
			Ref:           ref,
			View:          &dnsView,
			Name:          &updatedFqdn,
			MailExchanger: &updatedMx,
			Preference:    &updatedPreference,
			Ttl:           &updatedTtl,
			UseTtl:        utils.BoolPtr(true),
			Comment:       &updatedComment,
			Ea:            updatedEAs,
		})

		expectedObj := NewRecordMX(RecordMX{
			Ref:           ref,
			View:          &dnsView,
			Name:          &updatedFqdn,
			MailExchanger: &updatedMx,
			Preference:    &updatedPreference,
			Ttl:           &updatedTtl,
			UseTtl:        utils.BoolPtr(true),
			Comment:       &updatedComment,
			Ea:            updatedEAs,
		})

		conn = &fakeConnector{
			getObjectObj:         NewEmptyRecordMX(),
			getObjectQueryParams: NewQueryParams(false, nil),
			getObjectRef:         updatedRef,
			getObjectError:       nil,
			resultObject:         initObj,

			updateObjectObj:   updateObjIn,
			updateObjectRef:   ref,
			updateObjectError: nil,

			fakeRefReturn: updatedRef,
		}
		objMgr = NewObjectManager(conn, cmpType, tenantID)
		It("should pass updated MX record arguments", func() {
			actualObj, err = objMgr.UpdateMXRecord(ref, dnsView, updatedFqdn, updatedMx, updatedPreference, updatedTtl, true, updatedComment, updatedEAs)
		})
		It("should return expected MX record obj", func() {
			Expect(err).To(BeNil())
			Expect(actualObj).To(BeEquivalentTo(expectedObj))
		})

	})

	Describe("Get MX Record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		dnsView := "default"

		fqdn := "test.example.com"
		mx := "example.com"
		preference := uint32(25)
		ttl := uint32(70)
		comment := "test comment"

		fakeRefReturn := fmt.Sprintf("record:mx/ZG5zLmhvc3RjkuMC4xLg:%s/%s", fqdn, dnsView)

		eas := make(EA)
		eas["VM ID"] = "93f9249abc039284"
		eas["VM Name"] = "dummyvm"

		sf := map[string]string{
			"view":           dnsView,
			"name":           fqdn,
			"mail_exchanger": mx,
			"preference":     fmt.Sprintf("%d", preference),
		}
		nwFakeConnector := &fakeConnector{
			getObjectObj:         NewEmptyRecordMX(),
			getObjectQueryParams: NewQueryParams(false, sf),
			resultObject: []RecordMX{*NewRecordMX(RecordMX{
				View:          &dnsView,
				Name:          &fqdn,
				MailExchanger: &mx,
				Preference:    &preference,
				Ttl:           &ttl,
				UseTtl:        utils.BoolPtr(true),
				Comment:       &comment,
				Ea:            eas,
				Ref:           fakeRefReturn,
			})},
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(nwFakeConnector, cmpType, tenantID)

		var actualRecord *RecordMX
		var err error
		It("should pass expected MX record object to GetObject", func() {
			actualRecord, err = objMgr.GetMXRecord(dnsView, fqdn, mx, preference)
		})
		It("should return expected MX record Object", func() {
			Expect(actualRecord).NotTo(BeNil())
			Expect(*actualRecord).To(Equal(nwFakeConnector.resultObject.([]RecordMX)[0]))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get MX Record By Ref", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		dnsView := "default"

		fqdn := "test.example.com"
		readObjRef := fmt.Sprintf("record:mx/ZG5zLmhvc3RjkuMC4xLg:%s/%s", fqdn, dnsView)
		eas := EA{
			"VM ID":   "93f9249abc039284",
			"VM Name": "dummyvm",
		}
		nwFakeConnector := &fakeConnector{
			getObjectRef:         readObjRef,
			getObjectObj:         NewEmptyRecordMX(),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject: NewRecordMX(RecordMX{
				View:          &dnsView,
				Name:          &fqdn,
				MailExchanger: utils.StringPtr("example.com"),
				Preference:    utils.Uint32Ptr(25),
				Ttl:           utils.Uint32Ptr(70),
				UseTtl:        utils.BoolPtr(true),
				Comment:       utils.StringPtr("test comment"),
				Ea:            eas,
				Ref:           readObjRef,
			}),
			fakeRefReturn: readObjRef,
		}

		objMgr := NewObjectManager(nwFakeConnector, cmpType, tenantID)

		var actualRecord *RecordMX
		var err error
		It("should pass expected MX record ref to GetObject", func() {
			actualRecord, err = objMgr.GetMXRecordByRef(readObjRef)
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
