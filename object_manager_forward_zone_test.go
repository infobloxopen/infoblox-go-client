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
		forwardTo := NullableNameServers{
			NameServers: []NameServer{
				{Name: "fz1.test.com", Address: "10.0.0.1"},
				{Name: "fz2.test.com", Address: "10.0.0.2"},
			},
			IsNull: false,
		}
		fqdn := "test.fz.com"
		fakeRefReturn := fmt.Sprintf("zone_forward/ZG5zLmhvc3QkLZhd3QuaDE:%s/%20%20", fqdn)

		conn := &fakeConnector{
			createObjectObj:      NewZoneForward(comment, false, nil, forwardTo, false, nil, fqdn, "", "", "", "", ""),
			getObjectRef:         fakeRefReturn,
			getObjectObj:         NewEmptyZoneForward(),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         NewZoneForward(comment, false, nil, forwardTo, false, nil, fqdn, "", "", "", fakeRefReturn, ""),
			fakeRefReturn:        fakeRefReturn,
		}
		objMgr := NewObjectManager(conn, cmpType, tenantID)

		var actualRecord *ZoneForward
		var err error
		It("should pass expected Forward Zone Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateZoneForward(comment, disable, nil, forwardTo, false, nil, fqdn, "", "", "", "")
			Expect(actualRecord).To(Equal(conn.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Negative case: return an error message if required fields are not passed to create forward zone", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		comment := "test client"
		disable := false
		forwardTo := NullableNameServers{
			NameServers: []NameServer{
				{Name: "fz1.test.com", Address: "10.0.0.1"},
				{Name: "fz2.test.com", Address: "10.0.0.2"},
			},
			IsNull: false,
		}
		fqdn := "" //"test.fz.com"
		fakeRefReturn := fmt.Sprintf("zone_forward/ZG5zLmhvc3QkLZhd3QuaDE:%s/%20%20", fqdn)

		conn := &fakeConnector{
			createObjectObj:      NewZoneForward(comment, false, nil, forwardTo, false, nil, "", "", "", "", "", ""),
			getObjectRef:         fakeRefReturn,
			getObjectObj:         NewEmptyZoneForward(),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         NewZoneForward(comment, false, nil, forwardTo, false, nil, fqdn, "", "", "", fakeRefReturn, ""),
			fakeRefReturn:        fakeRefReturn,
			createObjectError:    fmt.Errorf("FQDN is required to create a forward zone"),
		}
		objMgr := NewObjectManager(conn, cmpType, tenantID)

		var actualRecord *ZoneForward
		var err error
		It("should pass expected Forward Zone Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateZoneForward(comment, disable, nil, forwardTo, false, nil, fqdn, "", "", "", "")
			Expect(actualRecord).To(BeNil())
			Expect(err).To(Equal(conn.createObjectError))
		})
	})

	Describe("Get forward zone test", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		comment := "test client"
		view := "default"

		forwardTo := NullableNameServers{
			NameServers: []NameServer{
				{Name: "fz1.test.com", Address: "10.0.0.1"},
				{Name: "fz2.test.com", Address: "10.0.0.2"},
			},
			IsNull: false,
		}
		fqdn := "test.fz.com"
		fakeRefReturn := fmt.Sprintf("zone_forward/ZG5zLmhvc3QkLZhd3QuaDE:%s/%s", fqdn, view)

		filters := map[string]string{
			"fqdn":    fqdn,
			"view":    view,
			"comment": comment,
		}

		queryParams := NewQueryParams(false, filters)
		conn := &fakeConnector{
			createObjectObj:      NewZoneForward(comment, false, nil, forwardTo, false, nil, fqdn, "", view, "", fakeRefReturn, ""),
			getObjectRef:         "",
			getObjectObj:         NewEmptyZoneForward(),
			resultObject:         []ZoneForward{*NewZoneForward(comment, false, nil, forwardTo, false, nil, fqdn, "", view, "", fakeRefReturn, "")},
			fakeRefReturn:        fakeRefReturn,
			getObjectQueryParams: queryParams,
		}
		objMgr := NewObjectManager(conn, cmpType, tenantID)

		conn.resultObject.([]ZoneForward)[0].Ref = fakeRefReturn

		var actualRecord []ZoneForward
		var err error
		It("should pass expected Forward Zone Object to GetObject", func() {
			actualRecord, err = objMgr.GetZoneForwardFilters(queryParams)
			Expect(err).To(BeNil())
			Expect(actualRecord[0]).To(Equal(conn.resultObject.([]ZoneForward)[0]))
		})
	})

	Describe("Get forward zone by Ref", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		comment := "test client"
		view := "default"
		//eas := EA{"Cloud API Owned": true}
		//forwardTo := []NameServer{{Name: "fz1.test.com", Address: "10.0.0.1"}, {Name: "fz2.test.com", Address: "10.0.0.2"}}
		forwardTo := NullableNameServers{
			NameServers: []NameServer{
				{Name: "fz1.test.com", Address: "10.0.0.1"},
				{Name: "fz2.test.com", Address: "10.0.0.2"},
			},
			IsNull: false,
		}
		fqdn := "test12.ex.com"
		fakeRefReturn := fmt.Sprintf("zone_forward/ZG5zLm5ldHdvcmtfdmlldyQyMw:%s/%s", fqdn, view)
		conn := &fakeConnector{
			getObjectRef:         fakeRefReturn,
			getObjectObj:         NewEmptyZoneForward(),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         NewZoneForward(comment, false, nil, forwardTo, false, nil, fqdn, "", view, "", fakeRefReturn, ""),
			fakeRefReturn:        fakeRefReturn,
		}
		objMgr := NewObjectManager(conn, cmpType, tenantID)

		conn.resultObject.(*ZoneForward).Ref = fakeRefReturn

		var actualRecord *ZoneForward
		var err error
		It("should pass expected Forward Zone Ref to GetObject", func() {
			actualRecord, err = objMgr.GetZoneForwardByRef(fakeRefReturn)
			Expect(err).To(BeNil())
			Expect(actualRecord).To(Equal(conn.resultObject))
		})
	})

	//TODO: Implement update forward zone properly

	//Describe("Update Forward Zone", func() {
	//	cmpType := "Docker"
	//	tenantID := "01234567890abcdef01234567890abcdef"
	//	initComment := "comment"
	//	updatedComment := "updated comment"
	//	disable := true
	//	forwardTo := NullForwardTo{
	//		ForwardTo: []NameServer{
	//			{Name: "fz1.test.com", Address: "10.0.0.1"},
	//			{Name: "fz2.test.com", Address: "10.0.0.2"},
	//		},
	//		IsNull: false,
	//	}
	//	view := "default"
	//	fqdn := "updated.fz.com"
	//	initRef := fmt.Sprintf("zone_forward/ZG5zLmhvc3QkLZhd3QuaDE:%s/%s", fqdn, view)
	//	zoneFormat := "Forward"
	//	updatedRef := initRef
	//
	//	initObj := NewZoneForward(initComment, disable, nil, forwardTo, false, nil, fqdn, "", view, zoneFormat, initRef, "")
	//	updatedObj := NewZoneForward(updatedComment, disable, nil, forwardTo, false, nil, fqdn, "", view, zoneFormat, updatedRef, "")
	//	expectedObj := NewZoneForward(updatedComment, disable, nil, forwardTo, false, nil, fqdn, "", view, zoneFormat, updatedRef, "")
	//
	//	conn := &fakeConnector{
	//		getObjectObj:         NewEmptyZoneForward(),
	//		getObjectQueryParams: NewQueryParams(false, nil),
	//		getObjectRef:         initRef,
	//		getObjectError:       nil,
	//		resultObject:         initObj,
	//
	//		updateObjectObj:   updatedObj,
	//		updateObjectRef:   initRef,
	//		updateObjectError: nil,
	//	}
	//	objMgr := NewObjectManager(conn, cmpType, tenantID)
	//
	//	var actualRecord *ZoneForward
	//	var err error
	//	It("should pass expected Forward Zone Object to UpdateObject", func() {
	//		actualRecord, err = objMgr.UpdateZoneForward(initRef, updatedComment, disable, nil, forwardTo, false, nil, "", "")
	//		Expect(actualRecord).To(BeEquivalentTo(expectedObj))
	//		Expect(err).To(BeNil())
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
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})
})
