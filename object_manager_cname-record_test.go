package ibclient

import (
	"fmt"

	. "github.com/onsi/ginkgo"
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
			createObjectObj: &RecordCNAME{
				View: dnsView, Canonical: canonical, Name: recordName,
				UseTtl: useTtl, Ttl: ttl, Comment: comment, Ea: eas,
			},
			getObjectRef:         fakeRefReturn,
			getObjectObj:         &RecordCNAME{},
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject: &RecordCNAME{
				View: dnsView, Canonical: canonical, Name: recordName, UseTtl: useTtl,
				Ttl: ttl, Comment: comment, Ea: eas, Ref: fakeRefReturn,
			},
			fakeRefReturn: fakeRefReturn,
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
			createObjectObj:   &RecordCNAME{View: dnsView, UseTtl: useTtl, Ttl: ttl, Comment: comment, Ea: eas},
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
			getObjectObj:         &RecordCNAME{},
			getObjectQueryParams: queryParams,
			resultObject: []RecordCNAME{
				{View: dnsView, Canonical: canonical, Name: recordName, UseTtl: useTtl, Ttl: ttl, Ref: fakeRefReturn},
			},
			fakeRefReturn: fakeRefReturn,
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
			getObjectObj:         &RecordCNAME{},
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
