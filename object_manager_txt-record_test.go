package ibclient

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Object Manager: TXT-record", func() {
	Describe("Create TXT Record ", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		dnsView := "default"
		text := "test-text"
		recordName := "test"
		useTtl := true
		ttl := uint32(70)
		comment := "creation test"
		eas := EA{"Country": "test"}
		fakeRefReturn := fmt.Sprintf("record:txt/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)

		aniFakeConnector := &fakeConnector{
			createObjectObj:      NewRecordTXT(dnsView, "", recordName, text, ttl, useTtl, comment, eas),
			getObjectRef:         fakeRefReturn,
			getObjectObj:         NewEmptyRecordTXT(),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         NewRecordTXT(dnsView, "", recordName, text, ttl, useTtl, comment, eas),
			fakeRefReturn:        fakeRefReturn,
		}

		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualRecord *RecordTXT
		var err error
		It("should pass expected TXT record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateTXTRecord(dnsView, recordName, text, ttl, useTtl, comment, eas)
		})
		It("should return expected TXT record Object", func() {
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Update TXT record", func() {
		var (
			err       error
			objMgr    IBObjectManager
			conn      *fakeConnector
			ref       string
			actualObj *RecordTXT
		)

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		recordName := "test"

		It("Updating text, ttl, useTtl, comment and EAs", func() {
			ref = fmt.Sprintf("record:txt/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)
			initialEas := EA{"Country": "old value"}
			initObj := NewRecordTXT("", "", recordName, "old-text", uint32(70), true, "old comment", initialEas)
			initObj.Ref = ref

			expectedEas := EA{"Country": "new value"}

			updateText := ""
			updateComment := "new comment"
			updateUseTtl := true
			updateTtl := uint32(10)
			updatedRef := fmt.Sprintf("record:txt/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)
			updateObjIn := NewRecordTXT("", "", recordName, updateText, updateTtl, updateUseTtl, updateComment, expectedEas)
			updateObjIn.Ref = ref

			expectedObj := NewRecordTXT("", "", recordName, updateText, updateTtl, updateUseTtl, updateComment, expectedEas)
			expectedObj.Ref = updatedRef

			conn = &fakeConnector{
				getObjectObj:         NewEmptyRecordTXT(),
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

			actualObj, err = objMgr.UpdateTXTRecord(ref, recordName, updateText, updateTtl, updateUseTtl, updateComment, expectedEas)
			Expect(err).To(BeNil())
			Expect(*actualObj).To(BeEquivalentTo(*expectedObj))
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

		expectedNetworkView := NetworkView{Ref: netviewRef, Name: &netviewName}
		It("should return expected Network View Object", func() {
			Expect(*BuildNetworkViewFromRef(netviewRef)).To(Equal(expectedNetworkView))
		})
		It("should failed if bad Network View Ref is provided", func() {
			Expect(BuildNetworkViewFromRef("bad")).To(BeNil())
		})
	})
})
