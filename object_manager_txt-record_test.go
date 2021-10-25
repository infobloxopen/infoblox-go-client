package ibclient

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Object Manager: TXT-record", func() {
	Describe("Allocate TXT Record ", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		text := "test-text"
		dnsView := "default"
		recordName := "test"
		ttl := uint(30)
		fakeRefReturn := fmt.Sprintf("record:txt/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)

		aniFakeConnector := &fakeConnector{
			createObjectObj: NewRecordTXT(RecordTXT{
				Name: recordName,
				Text: text,
				Ttl:  ttl,
				View: dnsView,
			}),
			getObjectRef: fakeRefReturn,
			getObjectObj: NewRecordTXT(RecordTXT{
				Name: recordName,
				Text: text,
				View: dnsView,
				Ref:  fakeRefReturn,
				Ttl:  ttl,
			}),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject: NewRecordTXT(RecordTXT{
				Name: recordName,
				Text: text,
				View: dnsView,
				Ttl:  ttl,
				Ref:  fakeRefReturn,
			}),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualRecord *RecordTXT
		var err error
		It("should pass expected TXT record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateTXTRecord(recordName, text, 30, dnsView)
		})
		It("should return expected TXT record Object", func() {
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete TXT Record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		recordName := "test"
		deleteRef := fmt.Sprintf("record:txt/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)
		fakeRefReturn := deleteRef
		nwFakeConnector := &fakeConnector{
			deleteObjectRef: deleteRef,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(nwFakeConnector, cmpType, tenantID)

		var actualRef string
		var err error
		It("should pass expected TXT record Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteTXTRecord(deleteRef)
		})
		It("should return expected TXT record Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})

	Describe("BuildNetworkViewFromRef", func() {
		netviewName := "default_view"
		netviewRef := fmt.Sprintf("networkview/ZG5zLm5ldHdvcmtfdmlldyQyMw:%s/false", netviewName)

		expectedNetworkView := NetworkView{Ref: netviewRef, Name: netviewName}
		It("should return expected Network View Object", func() {
			Expect(*BuildNetworkViewFromRef(netviewRef)).To(Equal(expectedNetworkView))
		})
		It("should failed if bad Network View Ref is provided", func() {
			Expect(BuildNetworkViewFromRef("bad")).To(BeNil())
		})
	})
})
