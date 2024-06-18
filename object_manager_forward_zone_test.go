package ibclient

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Object Manager: forward zone", func() {
	Describe("Create Forward Zone", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		comment := "test client"
		disable := false
		//eas := EA{"Cloud API Owned": true}
		forwardTo := []NameServer{{Name: "fz1.test.com", Address: "10.0.0.1"}, {Name: "fz2.test.com", Address: "10.0.0.2"}}
		fqdn := "test.fz.com"
		fakeRefReturn := fmt.Sprintf("record:txt/ZG5zLmJpbmRfY25h:%s/%20%20", fqdn)

		conn := &fakeConnector{
			createObjectObj:      NewZoneForward("comment", false, nil, forwardTo, false, nil, fqdn, "", "", ""),
			getObjectRef:         fakeRefReturn,
			getObjectObj:         NewEmptyZoneForward(),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         NewZoneForward("comment", false, nil, forwardTo, false, nil, fqdn, "", "", ""),
			fakeRefReturn:        fakeRefReturn,
		}
		objMgr := NewObjectManager(conn, cmpType, tenantID)

		var actualRecord *ZoneForward
		var err error
		It("should pass expected Forward Zone Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateZoneForward(comment, disable, nil, forwardTo, false, nil, fqdn, "", "", "")
		})
		It("should return expected Forward Zone Object", func() {
			Expect(actualRecord).To(Equal(conn.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Negative case: return an error message if required fields are not passed to create forward zone", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		comment := "test client"
		disable := false
		//eas := EA{"Cloud API Owned": true}
		forwardTo := []NameServer{{Name: "fz1.test.com", Address: "10.0.0.1"}, {Name: "fz2.test.com", Address: "10.0.0.2"}}
		fqdn := "test.fz.com"
		fakeRefReturn := fmt.Sprintf("record:txt/ZG5zLmJpbmRfY25h:%s/%20%20", fqdn)

		conn := &fakeConnector{
			createObjectObj:      NewZoneForward("comment", false, nil, forwardTo, false, nil, "", "", "", ""),
			getObjectRef:         fakeRefReturn,
			getObjectObj:         NewEmptyZoneForward(),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         NewZoneForward("comment", false, nil, forwardTo, false, nil, fqdn, "", "", ""),
			fakeRefReturn:        fakeRefReturn,
		}
		objMgr := NewObjectManager(conn, cmpType, tenantID)

		var actualRecord *ZoneForward
		var err error
		It("should pass expected Forward Zone Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateZoneForward(comment, disable, nil, forwardTo, false, nil, fqdn, "", "", "")
		})
		It("should return expected Forward Zone Object", func() {
			Expect(actualRecord).To(BeNil())
			Expect(err).To(Equal(conn.createObjectError))
		})
	})

	Describe("Get forward zone", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		comment := "test client"
		view := "default"
		//eas := EA{"Cloud API Owned": true}
		forwardTo := []NameServer{{Name: "fz1.test.com", Address: "10.0.0.1"}, {Name: "fz2.test.com", Address: "10.0.0.2"}}
		fqdn := "test.fz.com"
		fakeRefReturn := fmt.Sprintf("record:txt/ZG5zLmJpbmRfY25h:%s/%20%20", fqdn)

		queryParams := NewQueryParams(
			false,
			map[string]string{
				"fqdn":    fqdn,
				"comment": comment,
				"view":    view,
			})
		conn := &fakeConnector{
			createObjectObj:      NewZoneForward("comment", false, nil, forwardTo, false, nil, "", "", "", ""),
			getObjectRef:         "",
			getObjectObj:         NewEmptyZoneForward(),
			getObjectQueryParams: queryParams,
			resultObject:         []ZoneForward{*NewZoneForward("comment", false, nil, forwardTo, false, nil, fqdn, "", "", "")},
			fakeRefReturn:        fakeRefReturn,
		}
		objMgr := NewObjectManager(conn, cmpType, tenantID)

		var actualRecord []ZoneForward
		var err error
		It("should pass expected Forward Zone Object to GetObject", func() {
			actualRecord, err = objMgr.GetZoneForwardFilters(*queryParams)
		})
		It("should return expected Forward Zone Object", func() {
			Expect(actualRecord).To(Equal(conn.resultObject.([]ZoneForward)))
			Expect(err).To(BeNil())
		})
	})

	//Describe("Update Zone Forward", func() {
	//	var (
	//		err       error
	//		objMgr    IBObjectManager
	//		conn      *fakeConnector
	//		ref       string
	//		actualObj *ZoneForward
	//	)
	//
	//	cmpType := "Docker"
	//	tenantID := "01234567890abcdef01234567890abcdef"
	//	fqdn := "test.fz.com"
	//	refBase := "ZG5zLm5ldHdvcmtfdmlldyQyMw"
	//
	//	It("Updating comment and EAs", func() {
	//		ref = fmt.Sprintf("zone_forward/%s:%s", refBase, fqdn)
	//		initialEas := EA{
	//			"ea0": "ea0_old_value",
	//			"ea1": "ea1_old_value",
	//			"ea3": "ea3_value",
	//			"ea4": "ea4_value",
	//			"ea5": "ea5_old_value"}
	//		initObj := NewZoneForward("comment", false, initialEas,
	//			[]NameServer{{Address: "1.2.3.4", Name: "fz1.test.com"}}, false, nil, fqdn, "", "", "")
	//
	//		setEas := EA{
	//			"ea0": "ea0_old_value",
	//			"ea1": "ea1_new_value",
	//			"ea2": "ea2_new_value",
	//			"ea5": "ea5_old_value"}
	//		expectedEas := setEas
	//
	//		getObjIn := NewEmptyZoneForward()
	//
	//		comment := "test comment 1"
	//		updatedRef := fmt.Sprintf("zone_forward/%s:%s", refBase, fqdn)
	//		updateObjIn := NewNetworkView(updateNetviewName, comment, expectedEas, ref)
	//
	//		expectedObj := NewNetworkView(updateNetviewName, comment, expectedEas, updatedRef)
	//
	//		conn = &fakeConnector{
	//			getObjectObj:         getObjIn,
	//			getObjectQueryParams: NewQueryParams(false, nil),
	//			getObjectRef:         ref,
	//			getObjectError:       nil,
	//			resultObject:         initObj,
	//
	//			updateObjectObj:   updateObjIn,
	//			updateObjectRef:   ref,
	//			updateObjectError: nil,
	//
	//			fakeRefReturn: updatedRef,
	//		}
	//		objMgr = NewObjectManager(conn, cmpType, tenantID)
	//
	//		actualObj, err = objMgr.UpdateNetworkView(ref, updateNetviewName, comment, setEas)
	//		Expect(err).To(BeNil())
	//		Expect(actualObj).To(BeEquivalentTo(expectedObj))
	//	})
	//})

	Describe("Delete Zone Forward", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		fakeRefReturn := "zone_forward/ZG5zLm5ldHdvcmtfdmlldyQyMw:test12.ex.com/default"
		deleteRef := fakeRefReturn
		nwFakeConnector := &fakeConnector{
			deleteObjectRef: deleteRef,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(nwFakeConnector, cmpType, tenantID)

		var actualRef string
		var err error
		It("should pass expected zone forward Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteZoneForward(deleteRef)
		})
		It("should return expected zone forward Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})
})
