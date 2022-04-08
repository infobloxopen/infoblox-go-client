package ibclient

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Object Manager: NS-record", func() {
Describe("Allocate NS Record ", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		ns := "ns.example.com"
		dnsView := "default"
		recordName := "test.example.com"
		addr := []string{"10.0.0.1"}
		resultAddr := []ZoneNameServer{{Address: addr[0]}}

		fakeRefReturn := fmt.Sprintf("record:ns/ZG5zLmhvc3RjkuMC4xLg:%s/%s/default", ns, recordName)

		aniFakeConnector := &fakeConnector{
			createObjectObj: NewRecordNS(RecordNS{
				Name:       recordName,
				NS:         ns,
				Addresses:  resultAddr,
				View:       dnsView,
			}),
			getObjectRef: fakeRefReturn,
			getObjectObj:      NewRecordNS(RecordNS{
				Name:      recordName,
				NS:        ns,
				Addresses: resultAddr,
				View:      dnsView,
				Ref:       fakeRefReturn,
			}),
			resultObject: NewRecordNS(RecordNS{
				Name:      recordName,
				NS:        ns,
				Addresses: resultAddr,
				View:      dnsView,
				Ref:       fakeRefReturn,
			}),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualRecord *RecordNS
		var err error
		It("should pass expected NS record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateNSRecord(dnsView, recordName, ns, addr)
		})
		It("should return expected NS record Object", func() {
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})
	Describe("Delete NS Record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		recordName := "test.example.com"
		ns := "ns1.example.com"
		deleteRef := fmt.Sprintf("record:ns/ZG5zLmJpbmRfY25h:%s/%s", ns, recordName)
		fakeRefReturn := deleteRef
		nwFakeConnector := &fakeConnector{
			deleteObjectRef: deleteRef,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(nwFakeConnector, cmpType, tenantID)

		var actualRef string
		var err error
		It("should pass expected NS record Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteTXTRecord(deleteRef)
		})
		It("should return expected NS record Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})
})
