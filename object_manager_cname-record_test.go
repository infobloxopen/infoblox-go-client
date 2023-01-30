package ibclient

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Object Manager: CNAME-record", func() {
	Describe("Allocate CNAME Record ", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		canonical := "test-canonical.domain.com"
		dnsView := "default"
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		recordName := "test.domain.com"
		useTtl := false
		ttl := uint32(0)
		comment := "test CNAME record creation"
		fakeRefReturn := fmt.Sprintf("record:cname/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)
		eas := EA{"VM Name": vmName, "VM ID": vmID}

		conn := &fakeConnector{
			createObjectObj:      NewRecordCNAME(dnsView, canonical, recordName, useTtl, ttl, comment, eas, ""),
			getObjectRef:         fakeRefReturn,
			getObjectObj:         NewEmptyRecordCNAME(),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         NewRecordCNAME(dnsView, canonical, recordName, useTtl, ttl, comment, eas, fakeRefReturn),
			fakeRefReturn:        fakeRefReturn,
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)
		var actualRecord *RecordCNAME
		var err error
		It("should pass expected CNAME record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateCNAMERecord(dnsView, canonical, recordName, useTtl, ttl, comment, eas)
		})
		It("should return expected CNAME record Object", func() {
			Expect(actualRecord).To(Equal(conn.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Negative case: returns an error message if required fields are not passed", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		dnsView := "default"
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		useTtl := false
		ttl := uint32(0)
		comment := "test CNAME record creation"
		eas := EA{"VM Name": vmName, "VM ID": vmID}

		conn := &fakeConnector{
			createObjectObj:   NewRecordCNAME(dnsView, "", "", useTtl, ttl, comment, eas, ""),
			createObjectError: fmt.Errorf("canonical name and record name fields are required to create a CNAME record"),
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)
		var actualRecord, expectedObj *RecordCNAME
		var err error
		It("should pass expected CNAME record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateCNAMERecord(dnsView, "", "", useTtl, ttl, comment, eas)
		})
		It("should return expected CNAME record Object", func() {
			Expect(actualRecord).To(Equal(expectedObj))
			Expect(err).To(Equal(conn.createObjectError))
		})
	})

	Describe("Get CNAME Record ", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		canonical := "test-canonical.domain.com"
		dnsView := "default"
		recordName := "test.domain.com"
		useTtl := false
		ttl := uint32(0)
		fakeRefReturn := fmt.Sprintf("record:cname/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)

		queryParams := NewQueryParams(
			false,
			map[string]string{
				"view":      dnsView,
				"canonical": canonical,
				"name":      recordName,
			})
		conn := &fakeConnector{
			getObjectRef:         "",
			getObjectObj:         NewEmptyRecordCNAME(),
			getObjectQueryParams: queryParams,
			resultObject:         []RecordCNAME{*NewRecordCNAME(dnsView, canonical, recordName, useTtl, ttl, "", nil, fakeRefReturn)},
			fakeRefReturn:        fakeRefReturn,
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)
		var actualRecord *RecordCNAME
		var err error
		It("should pass expected CNAME record Object to CreateObject", func() {
			actualRecord, err = objMgr.GetCNAMERecord(dnsView, canonical, recordName)
		})
		It("should return expected CNAME record Object", func() {
			Expect(*actualRecord).To(Equal(conn.resultObject.([]RecordCNAME)[0]))
			Expect(err).To(BeNil())
		})
	})

	Describe("Negative case: return an error mesage when the required fields to get a unique record are not passed", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		canonical := "test-canonical.domain.com"
		recordName := "test.domain.com"

		queryParams := NewQueryParams(
			false,
			map[string]string{
				"canonical": canonical,
				"name":      recordName,
			})
		conn := &fakeConnector{
			getObjectRef:         "",
			getObjectObj:         NewEmptyRecordCNAME(),
			getObjectQueryParams: queryParams,
			getObjectError:       fmt.Errorf("DNS view, canonical name and record name of the record are required to retreive a unique CNAME record"),
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)
		var actualRecord, expectedObj *RecordCNAME
		var err error
		expectedObj = nil
		It("should pass expected CNAME record Object to CreateObject", func() {
			actualRecord, err = objMgr.GetCNAMERecord("", canonical, recordName)
		})
		It("should return expected CNAME record Object", func() {
			Expect(actualRecord).To(Equal(expectedObj))
			Expect(err).To(Equal(conn.getObjectError))
		})
	})

	Describe("Update CNAME record", func() {
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
		dnsView := "default"

		It("IPv4, updating ptrdname, IPv4 address, comment and EAs", func() {
			ref = fmt.Sprintf("record:cname/%s:%s", refBase, recordName)
			initialEas := EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_old_value",
				"ea3": "ea3_value",
				"ea4": "ea4_value",
				"ea5": "ea5_old_value"}
			initObj := NewRecordCNAME(dnsView, canonical, recordName, useTtl, ttl, "old comment", initialEas, ref)
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
			updateObjIn := NewRecordCNAME(dnsView, newCanonical, newRecordName, updateUseTtl, updateTtl, comment, expectedEas, ref)

			expectedObj := NewRecordCNAME(dnsView, newCanonical, newRecordName, updateUseTtl, updateTtl, comment, expectedEas, updatedRef)

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

			actualObj, err = objMgr.UpdateCNAMERecord(ref, dnsView, newCanonical, newRecordName, updateUseTtl, updateTtl, comment, setEas)
			Expect(err).To(BeNil())
			Expect(*actualObj).To(BeEquivalentTo(*expectedObj))
		})
	})

	Describe("Delete CNAME Record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		recordName := "test"
		deleteRef := fmt.Sprintf("record:CNAME/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)
		fakeRefReturn := deleteRef
		nwFakeConnector := &fakeConnector{
			deleteObjectRef: deleteRef,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(nwFakeConnector, cmpType, tenantID)

		var actualRef string
		var err error
		It("should pass expected CNAME record Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteCNAMERecord(deleteRef)
		})
		It("should return expected CNAME record Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})
})
