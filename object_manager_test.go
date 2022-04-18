package ibclient

import (
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type fakeConnector struct {
	// expected object to be passed to CreateObject()
	createObjectObj interface{}

	// expected object and reference to be passed to GetObject()
	getObjectObj         interface{}
	getObjectQueryParams interface{}
	getObjectRef         string

	// expected object and reference to be passed to UpdateObject()
	updateObjectObj interface{}
	updateObjectRef string

	// expected object's reference to be passed to DeleteObject()
	deleteObjectRef string

	// An object to be returned by GetObject() method.
	resultObject interface{}

	// A reference to be returned by Create/Update/Delete (not Get) methods.
	fakeRefReturn string

	// Error which fake Connector is to return on appropriate method call.
	createObjectError error
	getObjectError    error
	updateObjectError error
	deleteObjectError error
}

func (c *fakeConnector) CreateObject(obj IBObject) (string, error) {
	Expect(obj).To(Equal(c.createObjectObj))

	return c.fakeRefReturn, c.createObjectError
}

func (c *fakeConnector) GetObject(obj IBObject, ref string, qp *QueryParams, res interface{}) (err error) {
	Expect(obj).To(Equal(c.getObjectObj))
	Expect(qp).To(Equal(c.getObjectQueryParams))
	Expect(ref).To(Equal(c.getObjectRef))

	if ref == "" {
		switch obj.(type) {
		case *NetworkView:
			*res.(*[]NetworkView) = c.resultObject.([]NetworkView)
		case *NetworkContainer:
			*res.(*[]NetworkContainer) = c.resultObject.([]NetworkContainer)
		case *Network:
			*res.(*[]Network) = c.resultObject.([]Network)
		case *FixedAddress:
			*res.(*[]FixedAddress) = c.resultObject.([]FixedAddress)
		case *EADefinition:
			*res.(*[]EADefinition) = c.resultObject.([]EADefinition)
		case *CapacityReport:
			*res.(*[]CapacityReport) = c.resultObject.([]CapacityReport)
		case *UpgradeStatus:
			*res.(*[]UpgradeStatus) = c.resultObject.([]UpgradeStatus)
		case *Member:
			*res.(*[]Member) = c.resultObject.([]Member)
		case *Grid:
			*res.(*[]Grid) = c.resultObject.([]Grid)
		case *License:
			*res.(*[]License) = c.resultObject.([]License)
		case *HostRecord:
			*res.(*[]HostRecord) = c.resultObject.([]HostRecord)
		case *RecordAAAA:
			*res.(*[]RecordAAAA) = c.resultObject.([]RecordAAAA)
		case *RecordPTR:
			*res.(*[]RecordPTR) = c.resultObject.([]RecordPTR)
		case *RecordTXT:
			*res.(*[]RecordTXT) = c.resultObject.([]RecordTXT)
		case *ZoneDelegated:
			*res.(*[]ZoneDelegated) = c.resultObject.([]ZoneDelegated)
		case *RecordCNAME:
			*res.(*[]RecordCNAME) = c.resultObject.([]RecordCNAME)
		case *RecordA:
			*res.(*[]RecordA) = c.resultObject.([]RecordA)
		}
	} else {
		switch obj.(type) {
		case *ZoneAuth:
			*res.(*ZoneAuth) = *c.resultObject.(*ZoneAuth)
		case *NetworkView:
			*res.(*NetworkView) = *c.resultObject.(*NetworkView)
		case *NetworkContainer:
			*res.(*NetworkContainer) = *c.resultObject.(*NetworkContainer)
		case *Network:
			*res.(*Network) = *c.resultObject.(*Network)
		case *FixedAddress:
			**res.(**FixedAddress) = *c.resultObject.(*FixedAddress)
		case *HostRecord:
			**res.(**HostRecord) = *c.resultObject.(*HostRecord)
		case *RecordPTR:
			**res.(**RecordPTR) = *c.resultObject.(*RecordPTR)
		case *RecordTXT:
			**res.(**RecordTXT) = *c.resultObject.(*RecordTXT)
		case *RecordCNAME:
			**res.(**RecordCNAME) = *c.resultObject.(*RecordCNAME)
		case *RecordA:
			**res.(**RecordA) = *c.resultObject.(*RecordA)
		case *RecordAAAA:
			**res.(**RecordAAAA) = *c.resultObject.(*RecordAAAA)
		}
	}

	err = c.getObjectError
	return
}

func (c *fakeConnector) DeleteObject(ref string) (string, error) {
	Expect(ref).To(Equal(c.deleteObjectRef))

	return c.fakeRefReturn, c.deleteObjectError
}

func (c *fakeConnector) UpdateObject(obj IBObject, ref string) (string, error) {
	Expect(obj).To(Equal(c.updateObjectObj))
	Expect(ref).To(Equal(c.updateObjectRef))

	return c.fakeRefReturn, c.updateObjectError
}

var _ = Describe("Object Manager", func() {
	Describe("Get Capacity report", func() {
		cmpType := "Heka"
		tenantID := "0123"
		var name string = "Member1"
		fakeRefReturn := fmt.Sprintf("member/ZG5zLmJpbmRfY25h:/%s", name)

		queryParams := NewQueryParams(
			false,
			map[string]string{
				"name": name,
			})

		fakeConnector := &fakeConnector{
			getObjectObj:         NewCapcityReport(CapacityReport{}),
			getObjectQueryParams: queryParams,
			getObjectRef:         "",
			resultObject: []CapacityReport{*NewCapcityReport(CapacityReport{
				Ref:  fakeRefReturn,
				Name: name,
			})},
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(fakeConnector, cmpType, tenantID)

		var actualReport []CapacityReport
		var err error

		It("should pass expected Capacityreport object to GetObject", func() {
			actualReport, err = objMgr.GetCapacityReport(name)
		})
		It("should return expected CapacityReport Object", func() {
			Expect(actualReport[0]).To(Equal(fakeConnector.resultObject.([]CapacityReport)[0]))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get upgrade status", func() {
		cmpType := "Heka"
		tenantID := "0123"
		var StatusType = "GRID"
		fakeRefReturn := fmt.Sprintf("upgradestatus/Li51cGdyYWRlc3RhdHVzJHVwZ3JhZGVfc3RhdHVz:test")

		queryParams := NewQueryParams(
			false,
			map[string]string{
				"type": StatusType,
			})

		USFakeConnector := &fakeConnector{
			getObjectObj:         NewUpgradeStatus(UpgradeStatus{}),
			getObjectQueryParams: queryParams,
			getObjectRef:         "",
			resultObject: []UpgradeStatus{*NewUpgradeStatus(UpgradeStatus{
				Ref:  fakeRefReturn,
				Type: StatusType,
			})},
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(USFakeConnector, cmpType, tenantID)

		var actualStatus []UpgradeStatus
		var err error

		It("should pass expected upgradestatus object to GetObject", func() {
			actualStatus, err = objMgr.GetUpgradeStatus(StatusType)
		})
		It("should return expected upgradestatus Object", func() {
			Expect(actualStatus[0]).To(Equal(USFakeConnector.resultObject.([]UpgradeStatus)[0]))
			Expect(err).To(BeNil())
		})

	})
	Describe("Get upgrade status Error case", func() {
		cmpType := "Heka"
		tenantID := "0123"
		StatusType := ""
		fakeRefReturn := fmt.Sprintf("upgradestatus/Li51cGdyYWRlc3RhdHVzJHVwZ3JhZGVfc3RhdHVz:test")
		expectErr := errors.New("Status type can not be nil")
		USFakeConnector := &fakeConnector{
			getObjectObj:         NewUpgradeStatus(UpgradeStatus{Type: StatusType}),
			getObjectRef:         "",
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject: []UpgradeStatus{*NewUpgradeStatus(UpgradeStatus{
				Ref:  fakeRefReturn,
				Type: StatusType,
			})},
			fakeRefReturn: fakeRefReturn,
		}
		objMgr := NewObjectManager(USFakeConnector, cmpType, tenantID)
		It("upgradestatus object to GetObject", func() {
			_, err := objMgr.GetUpgradeStatus(StatusType)
			Expect(err).To(Equal(expectErr))
		})

	})
	Describe("GetAllMembers", func() {
		cmpType := "Heka"
		tenantID := "0123"
		var err error
		fakeRefReturn := fmt.Sprintf("member/Li51cGdyYWRlc3RhdHVzJHVwZ3JhZGVfc3RhdHVz:test")
		returnFields := []string{"host_name", "node_info", "time_zone"}
		MemFakeConnector := &fakeConnector{
			getObjectObj:         NewMember(Member{}),
			getObjectRef:         "",
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject: []Member{*NewMember(Member{
				Ref: fakeRefReturn,
			})},
			fakeRefReturn: fakeRefReturn,
		}
		objMgr := NewObjectManager(MemFakeConnector, cmpType, tenantID)
		var actualMembers []Member
		It("should return expected member Object", func() {
			actualMembers, err = objMgr.GetAllMembers()
			Expect(actualMembers[0]).To(Equal(MemFakeConnector.resultObject.([]Member)[0]))
			Expect(actualMembers[0].returnFields).To(Equal(returnFields))
			Expect(err).To(BeNil())
		})
	})

	Describe("GetGridInfo", func() {
		cmpType := "Heka"
		tenantID := "0123"
		var err error
		fakeRefReturn := fmt.Sprintf("grid/Li511cGdyYWRlc3RhdHVzJHVwZ3JhZGVfc3RhdHVz:test")
		returnFields := []string{"name", "ntp_setting"}
		GridFakeConnector := &fakeConnector{
			getObjectObj:         NewGrid(Grid{}),
			getObjectRef:         "",
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject: []Grid{*NewGrid(Grid{
				Ref: fakeRefReturn,
			})},
			fakeRefReturn: fakeRefReturn,
		}
		objMgr := NewObjectManager(GridFakeConnector, cmpType, tenantID)
		var actualGridInfo []Grid
		It("should return expected Grid Object", func() {
			actualGridInfo, err = objMgr.GetGridInfo()
			Expect(actualGridInfo[0]).To(Equal(GridFakeConnector.resultObject.([]Grid)[0]))
			Expect(actualGridInfo[0].returnFields).To(Equal(returnFields))
			Expect(err).To(BeNil())
		})
	})

	Describe("GetGridLicense", func() {
		cmpType := "Heka"
		tenantID := "0123"
		var err error
		fakeRefReturn := fmt.Sprintf("license/Li511cGdyYWRlc3RhdHVzJHVwZ3JhZGVfc3RhdHVz:test")
		returnFields := []string{"expiration_status",
			"expiry_date",
			"key",
			"limit",
			"limit_context",
			"type"}
		LicFakeConnector := &fakeConnector{
			getObjectObj:         NewGridLicense(License{}),
			getObjectRef:         "",
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject: []License{*NewGridLicense(License{
				Ref: fakeRefReturn,
			})},
			fakeRefReturn: fakeRefReturn,
		}
		objMgr := NewObjectManager(LicFakeConnector, cmpType, tenantID)
		var actualGridLicense []License
		It("should return expected License Object", func() {
			actualGridLicense, err = objMgr.GetGridLicense()
			Expect(actualGridLicense[0]).To(Equal(LicFakeConnector.resultObject.([]License)[0]))
			Expect(actualGridLicense[0].returnFields).To(Equal(returnFields))
			Expect(err).To(BeNil())
		})
	})

	Describe("Create Zone Auth", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		fqdn := "azone.example.com"
		fakeRefReturn := "zone_auth/ZG5zLnpvbmUkLl9kZWZhdWx0LnphLmNvLmFic2EuY2Fhcy5vaG15Z2xiLmdzbGJpYmNsaWVudA:dzone.example.com/default"
		zaFakeConnector := &fakeConnector{
			createObjectObj: NewZoneAuth(ZoneAuth{Fqdn: fqdn}),
			resultObject:    NewZoneAuth(ZoneAuth{Fqdn: fqdn, Ref: fakeRefReturn}),
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(zaFakeConnector, cmpType, tenantID)

		ea := make(EA)

		zaFakeConnector.createObjectObj.(*ZoneAuth).Ea = ea
		zaFakeConnector.createObjectObj.(*ZoneAuth).Ea["Tenant ID"] = tenantID
		zaFakeConnector.createObjectObj.(*ZoneAuth).Ea["CMP Type"] = cmpType

		zaFakeConnector.resultObject.(*ZoneAuth).Ea = ea
		zaFakeConnector.resultObject.(*ZoneAuth).Ea["Tenant ID"] = tenantID
		zaFakeConnector.resultObject.(*ZoneAuth).Ea["CMP Type"] = cmpType

		var actualZoneAuth *ZoneAuth
		var err error
		It("should pass expected ZoneAuth Object to CreateObject", func() {
			actualZoneAuth, err = objMgr.CreateZoneAuth(fqdn, ea)
		})
		It("should return expected ZoneAuth Object", func() {
			Expect(actualZoneAuth).To(Equal(zaFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get AuthZone by ref", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		fqdn := "azone.example.com"
		fakeRefReturn := "zone_delegated/ZG5zLnpvbmUkLl9kZWZhdWx0LnphLmNvLmFic2EuY2Fhcy5vaG15Z2xiLmdzbGJpYmNsaWVudA:azone.example.com/default"
		zdFakeConnector := &fakeConnector{
			getObjectObj:         NewZoneAuth(ZoneAuth{}),
			getObjectRef:         fakeRefReturn,
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         NewZoneAuth(ZoneAuth{Fqdn: fqdn}),
		}

		objMgr := NewObjectManager(zdFakeConnector, cmpType, tenantID)

		var actualZoneAuth, getNoRef *ZoneAuth
		getNoRef = nil
		var err error
		It("should pass expected ZoneAuth Object to GetObject", func() {
			actualZoneAuth, err = objMgr.GetZoneAuthByRef(fakeRefReturn)
		})
		fmt.Printf("doodo  %v", actualZoneAuth)
		It("should return expected ZoneAuth Object", func() {
			Expect(actualZoneAuth).To(Equal(zdFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
		It("should return empty ZoneAuth and nil error if ref is empty", func() {
			zdFakeConnector.getObjectObj.(*ZoneAuth).IBBase.returnFields = nil
			actualZoneAuth, err = objMgr.GetZoneAuthByRef("")
			Expect(actualZoneAuth).To(Equal(getNoRef))
			Expect(err).ToNot(BeNil())
		})
	})

	Describe("Delete ZoneAuth", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		deleteRef := "zone_auth/ZG5zLnpvbmUkLl9kZWZhdWx0LnphLmNvLmFic2EuY2Fhcy5vaG15Z2xiLmdzbGJpYmNsaWVudA:dzone.example.com/default"
		fakeRefReturn := deleteRef
		zaFakeConnector := &fakeConnector{
			deleteObjectRef: deleteRef,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(zaFakeConnector, cmpType, tenantID)

		var actualRef string
		var err error
		It("should pass expected ZoneAuth Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteZoneAuth(deleteRef)
		})
		It("should return expected ZoneAuth Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get Zone Delegated", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		fqdn := "dzone.example.com"
		fakeRefReturn := "zone_delegated/ZG5zLnpvbmUkLl9kZWZhdWx0LnphLmNvLmFic2EuY2Fhcy5vaG15Z2xiLmdzbGJpYmNsaWVudA:dzone.example.com/default"

		queryParams := NewQueryParams(
			false,
			map[string]string{
				"fqdn": fqdn,
			})

		zdFakeConnector := &fakeConnector{
			getObjectObj:         NewZoneDelegated(ZoneDelegated{}),
			getObjectQueryParams: queryParams,
			getObjectRef:         "",
			resultObject:         []ZoneDelegated{*NewZoneDelegated(ZoneDelegated{Fqdn: fqdn, Ref: fakeRefReturn})},
		}

		objMgr := NewObjectManager(zdFakeConnector, cmpType, tenantID)

		var actualZoneDelegated *ZoneDelegated
		var err error
		It("should pass expected ZoneDelegated Object to GetObject", func() {
			actualZoneDelegated, err = objMgr.GetZoneDelegated(fqdn)
		})
		It("should return expected ZoneDelegated Object", func() {
			Expect(*actualZoneDelegated).To(Equal(zdFakeConnector.resultObject.([]ZoneDelegated)[0]))
			Expect(err).To(BeNil())
		})
		It("should return nil if fqdn is empty", func() {
			zdFakeConnector.getObjectObj.(*ZoneDelegated).Fqdn = ""
			actualZoneDelegated, err = objMgr.GetZoneDelegated("")
			Expect(actualZoneDelegated).To(BeNil())
			Expect(err).To(BeNil())
		})
	})

	Describe("Create Zone Delegated", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		fqdn := "dzone.example.com"
		delegateTo := []NameServer{
			{Address: "10.0.0.1", Name: "test1.dzone.example.com"},
			{Address: "10.0.0.2", Name: "test2.dzone.example.com"}}
		fakeRefReturn := "zone_delegated/ZG5zLnpvbmUkLl9kZWZhdWx0LnphLmNvLmFic2EuY2Fhcy5vaG15Z2xiLmdzbGJpYmNsaWVudA:dzone.example.com/default"
		zdFakeConnector := &fakeConnector{
			createObjectObj: NewZoneDelegated(ZoneDelegated{Fqdn: fqdn, DelegateTo: delegateTo}),
			resultObject:    NewZoneDelegated(ZoneDelegated{Fqdn: fqdn, DelegateTo: delegateTo, Ref: fakeRefReturn}),
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(zdFakeConnector, cmpType, tenantID)

		var actualZoneDelegated *ZoneDelegated
		var err error
		It("should pass expected ZoneDelegated Object to CreateObject", func() {
			actualZoneDelegated, err = objMgr.CreateZoneDelegated(fqdn, delegateTo)
		})
		It("should return expected ZoneDelegated Object", func() {
			Expect(actualZoneDelegated).To(Equal(zdFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Update Zone Delegated", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		fakeRefReturn := "zone_delegated/ZG5zLnpvbmUkLl9kZWZhdWx0LnphLmNvLmFic2EuY2Fhcy5vaG15Z2xiLmdzbGJpYmNsaWVudA:dzone.example.com/default"
		delegateTo := []NameServer{
			{Address: "10.0.0.1", Name: "test1.dzone.example.com"},
			{Address: "10.0.0.2", Name: "test2.dzone.example.com"}}

		receiveUpdateObject := NewZoneDelegated(ZoneDelegated{Ref: fakeRefReturn, DelegateTo: delegateTo})
		returnUpdateObject := NewZoneDelegated(ZoneDelegated{DelegateTo: delegateTo, Ref: fakeRefReturn})
		zdFakeConnector := &fakeConnector{
			fakeRefReturn:   fakeRefReturn,
			resultObject:    returnUpdateObject,
			updateObjectObj: receiveUpdateObject,
			updateObjectRef: fakeRefReturn,
		}

		objMgr := NewObjectManager(zdFakeConnector, cmpType, tenantID)

		var updatedObject *ZoneDelegated
		var err error
		It("should pass expected updated object to UpdateObject", func() {
			updatedObject, err = objMgr.UpdateZoneDelegated(fakeRefReturn, delegateTo)
		})
		It("should update zone with new delegation server list with no error", func() {
			Expect(updatedObject).To(Equal(returnUpdateObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete ZoneDelegated", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		deleteRef := "zone_delegated/ZG5zLnpvbmUkLl9kZWZhdWx0LnphLmNvLmFic2EuY2Fhcy5vaG15Z2xiLmdzbGJpYmNsaWVudA:dzone.example.com/default"
		fakeRefReturn := deleteRef
		zdFakeConnector := &fakeConnector{
			deleteObjectRef: deleteRef,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(zdFakeConnector, cmpType, tenantID)

		var actualRef string
		var err error
		It("should pass expected ZoneDelegated Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteZoneDelegated(deleteRef)
		})
		It("should return expected ZoneDelegated Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})
})
