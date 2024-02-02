package ibclient_test

import (
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	ibclient "github.com/infobloxopen/infoblox-go-client/v2"
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

func (c *fakeConnector) CreateObject(obj ibclient.IBObject) (string, error) {
	Expect(obj).To(Equal(c.createObjectObj))

	return c.fakeRefReturn, c.createObjectError
}

func (c *fakeConnector) GetObject(obj ibclient.IBObject, ref string, qp *ibclient.QueryParams, res interface{}) (err error) {
	Expect(obj).To(Equal(c.getObjectObj))
	Expect(qp).To(Equal(c.getObjectQueryParams))
	Expect(ref).To(Equal(c.getObjectRef))

	if ref == "" {
		switch obj.(type) {
		case *ibclient.NetworkView:
			*res.(*[]ibclient.NetworkView) = c.resultObject.([]ibclient.NetworkView)
		case *ibclient.NetworkContainer:
			*res.(*[]ibclient.NetworkContainer) = c.resultObject.([]ibclient.NetworkContainer)
		case *ibclient.Network:
			*res.(*[]ibclient.Network) = c.resultObject.([]ibclient.Network)
		case *ibclient.FixedAddress:
			*res.(*[]ibclient.FixedAddress) = c.resultObject.([]ibclient.FixedAddress)
		case *ibclient.EADefinition:
			*res.(*[]ibclient.EADefinition) = c.resultObject.([]ibclient.EADefinition)
		case *ibclient.CapacityReport:
			*res.(*[]ibclient.CapacityReport) = c.resultObject.([]ibclient.CapacityReport)
		case *ibclient.UpgradeStatus:
			*res.(*[]ibclient.UpgradeStatus) = c.resultObject.([]ibclient.UpgradeStatus)
		case *ibclient.Member:
			*res.(*[]ibclient.Member) = c.resultObject.([]ibclient.Member)
		case *ibclient.Grid:
			*res.(*[]ibclient.Grid) = c.resultObject.([]ibclient.Grid)
		case *ibclient.License:
			*res.(*[]ibclient.License) = c.resultObject.([]ibclient.License)
		case *ibclient.HostRecord:
			*res.(*[]ibclient.HostRecord) = c.resultObject.([]ibclient.HostRecord)
		case *ibclient.RecordAAAA:
			*res.(*[]ibclient.RecordAAAA) = c.resultObject.([]ibclient.RecordAAAA)
		case *ibclient.RecordPTR:
			*res.(*[]ibclient.RecordPTR) = c.resultObject.([]ibclient.RecordPTR)
		case *ibclient.RecordSRV:
			*res.(*[]ibclient.RecordSRV) = c.resultObject.([]ibclient.RecordSRV)
		case *ibclient.RecordTXT:
			*res.(*[]ibclient.RecordTXT) = c.resultObject.([]ibclient.RecordTXT)
		case *ibclient.ZoneDelegated:
			*res.(*[]ibclient.ZoneDelegated) = c.resultObject.([]ibclient.ZoneDelegated)
		case *ibclient.RecordCNAME:
			*res.(*[]ibclient.RecordCNAME) = c.resultObject.([]ibclient.RecordCNAME)
		case *ibclient.RecordA:
			*res.(*[]ibclient.RecordA) = c.resultObject.([]ibclient.RecordA)
		case *ibclient.RecordMX:
			*res.(*[]ibclient.RecordMX) = c.resultObject.([]ibclient.RecordMX)
		}
	} else {
		switch obj.(type) {
		case *ibclient.ZoneAuth:
			*res.(*ibclient.ZoneAuth) = *c.resultObject.(*ibclient.ZoneAuth)
		case *ibclient.NetworkView:
			*res.(*ibclient.NetworkView) = *c.resultObject.(*ibclient.NetworkView)
		case *ibclient.NetworkContainer:
			*res.(*ibclient.NetworkContainer) = *c.resultObject.(*ibclient.NetworkContainer)
		case *ibclient.Network:
			*res.(*ibclient.Network) = *c.resultObject.(*ibclient.Network)
		case *ibclient.FixedAddress:
			**res.(**ibclient.FixedAddress) = *c.resultObject.(*ibclient.FixedAddress)
		case *ibclient.HostRecord:
			**res.(**ibclient.HostRecord) = *c.resultObject.(*ibclient.HostRecord)
		case *ibclient.RecordPTR:
			**res.(**ibclient.RecordPTR) = *c.resultObject.(*ibclient.RecordPTR)
		case *ibclient.RecordSRV:
			**res.(**ibclient.RecordSRV) = *c.resultObject.(*ibclient.RecordSRV)
		case *ibclient.RecordTXT:
			**res.(**ibclient.RecordTXT) = *c.resultObject.(*ibclient.RecordTXT)
		case *ibclient.RecordCNAME:
			**res.(**ibclient.RecordCNAME) = *c.resultObject.(*ibclient.RecordCNAME)
		case *ibclient.RecordA:
			**res.(**ibclient.RecordA) = *c.resultObject.(*ibclient.RecordA)
		case *ibclient.RecordAAAA:
			**res.(**ibclient.RecordAAAA) = *c.resultObject.(*ibclient.RecordAAAA)
		case *ibclient.RecordMX:
			**res.(**ibclient.RecordMX) = *c.resultObject.(*ibclient.RecordMX)
		}
	}

	err = c.getObjectError
	return
}

func (c *fakeConnector) DeleteObject(ref string) (string, error) {
	Expect(ref).To(Equal(c.deleteObjectRef))

	return c.fakeRefReturn, c.deleteObjectError
}

func (c *fakeConnector) UpdateObject(obj ibclient.IBObject, ref string) (string, error) {
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

		queryParams := ibclient.NewQueryParams(
			false,
			map[string]string{
				"name": name,
			})

		fakeConnector := &fakeConnector{
			getObjectObj:         ibclient.NewCapcityReport(ibclient.CapacityReport{}),
			getObjectQueryParams: queryParams,
			getObjectRef:         "",
			resultObject: []ibclient.CapacityReport{*ibclient.NewCapcityReport(ibclient.CapacityReport{
				Ref:  fakeRefReturn,
				Name: name,
			})},
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := ibclient.NewObjectManager(fakeConnector, cmpType, tenantID)

		var actualReport []ibclient.CapacityReport
		var err error

		It("should pass expected Capacityreport object to GetObject", func() {
			actualReport, err = objMgr.GetCapacityReport(name)
		})
		It("should return expected CapacityReport Object", func() {
			Expect(actualReport[0]).To(Equal(fakeConnector.resultObject.([]ibclient.CapacityReport)[0]))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get upgrade status", func() {
		cmpType := "Heka"
		tenantID := "0123"
		var StatusType = "GRID"
		fakeRefReturn := fmt.Sprintf("upgradestatus/Li51cGdyYWRlc3RhdHVzJHVwZ3JhZGVfc3RhdHVz:test")

		queryParams := ibclient.NewQueryParams(
			false,
			map[string]string{
				"type": StatusType,
			})

		USFakeConnector := &fakeConnector{
			getObjectObj:         ibclient.NewUpgradeStatus(ibclient.UpgradeStatus{}),
			getObjectQueryParams: queryParams,
			getObjectRef:         "",
			resultObject: []ibclient.UpgradeStatus{*ibclient.NewUpgradeStatus(ibclient.UpgradeStatus{
				Ref:  fakeRefReturn,
				Type: StatusType,
			})},
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := ibclient.NewObjectManager(USFakeConnector, cmpType, tenantID)

		var actualStatus []ibclient.UpgradeStatus
		var err error

		It("should pass expected upgradestatus object to GetObject", func() {
			actualStatus, err = objMgr.GetUpgradeStatus(StatusType)
		})
		It("should return expected upgradestatus Object", func() {
			Expect(actualStatus[0]).To(Equal(USFakeConnector.resultObject.([]ibclient.UpgradeStatus)[0]))
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
			getObjectObj:         ibclient.NewUpgradeStatus(ibclient.UpgradeStatus{Type: StatusType}),
			getObjectRef:         "",
			getObjectQueryParams: ibclient.NewQueryParams(false, nil),
			resultObject: []ibclient.UpgradeStatus{*ibclient.NewUpgradeStatus(ibclient.UpgradeStatus{
				Ref:  fakeRefReturn,
				Type: StatusType,
			})},
			fakeRefReturn: fakeRefReturn,
		}
		objMgr := ibclient.NewObjectManager(USFakeConnector, cmpType, tenantID)
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
			getObjectObj:         ibclient.NewMember(ibclient.Member{}),
			getObjectRef:         "",
			getObjectQueryParams: ibclient.NewQueryParams(false, nil),
			resultObject: []ibclient.Member{*ibclient.NewMember(ibclient.Member{
				Ref: fakeRefReturn,
			})},
			fakeRefReturn: fakeRefReturn,
		}
		objMgr := ibclient.NewObjectManager(MemFakeConnector, cmpType, tenantID)
		var actualMembers []ibclient.Member
		It("should return expected member Object", func() {
			actualMembers, err = objMgr.GetAllMembers()
			Expect(actualMembers[0]).To(Equal(MemFakeConnector.resultObject.([]ibclient.Member)[0]))
			Expect(actualMembers[0].ReturnFields()).To(Equal(returnFields))
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
			getObjectObj:         ibclient.NewGrid(ibclient.Grid{}),
			getObjectRef:         "",
			getObjectQueryParams: ibclient.NewQueryParams(false, nil),
			resultObject: []ibclient.Grid{*ibclient.NewGrid(ibclient.Grid{
				Ref: fakeRefReturn,
			})},
			fakeRefReturn: fakeRefReturn,
		}
		objMgr := ibclient.NewObjectManager(GridFakeConnector, cmpType, tenantID)
		var actualGridInfo []ibclient.Grid
		It("should return expected Grid Object", func() {
			actualGridInfo, err = objMgr.GetGridInfo()
			Expect(actualGridInfo[0]).To(Equal(GridFakeConnector.resultObject.([]ibclient.Grid)[0]))
			Expect(actualGridInfo[0].ReturnFields()).To(Equal(returnFields))
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
			getObjectObj:         ibclient.NewGridLicense(ibclient.License{}),
			getObjectRef:         "",
			getObjectQueryParams: ibclient.NewQueryParams(false, nil),
			resultObject: []ibclient.License{*ibclient.NewGridLicense(ibclient.License{
				Ref: fakeRefReturn,
			})},
			fakeRefReturn: fakeRefReturn,
		}
		objMgr := ibclient.NewObjectManager(LicFakeConnector, cmpType, tenantID)
		var actualGridLicense []ibclient.License
		It("should return expected License Object", func() {
			actualGridLicense, err = objMgr.GetGridLicense()
			Expect(actualGridLicense[0]).To(Equal(LicFakeConnector.resultObject.([]ibclient.License)[0]))
			Expect(actualGridLicense[0].ReturnFields()).To(Equal(returnFields))
			Expect(err).To(BeNil())
		})
	})

	Describe("Create Zone Auth", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		fqdn := "azone.example.com"
		fakeRefReturn := "zone_auth/ZG5zLnpvbmUkLl9kZWZhdWx0LnphLmNvLmFic2EuY2Fhcy5vaG15Z2xiLmdzbGJpYmNsaWVudA:dzone.example.com/default"
		zaFakeConnector := &fakeConnector{
			createObjectObj: ibclient.NewZoneAuth(ibclient.ZoneAuth{Fqdn: fqdn}),
			resultObject:    ibclient.NewZoneAuth(ibclient.ZoneAuth{Fqdn: fqdn, Ref: fakeRefReturn}),
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := ibclient.NewObjectManager(zaFakeConnector, cmpType, tenantID)

		ea := make(ibclient.EA)

		zaFakeConnector.createObjectObj.(*ibclient.ZoneAuth).Ea = ea
		zaFakeConnector.createObjectObj.(*ibclient.ZoneAuth).Ea["Tenant ID"] = tenantID
		zaFakeConnector.createObjectObj.(*ibclient.ZoneAuth).Ea["CMP Type"] = cmpType

		zaFakeConnector.resultObject.(*ibclient.ZoneAuth).Ea = ea
		zaFakeConnector.resultObject.(*ibclient.ZoneAuth).Ea["Tenant ID"] = tenantID
		zaFakeConnector.resultObject.(*ibclient.ZoneAuth).Ea["CMP Type"] = cmpType

		var actualZoneAuth *ibclient.ZoneAuth
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
			getObjectObj:         ibclient.NewZoneAuth(ibclient.ZoneAuth{}),
			getObjectRef:         fakeRefReturn,
			getObjectQueryParams: ibclient.NewQueryParams(false, nil),
			resultObject:         ibclient.NewZoneAuth(ibclient.ZoneAuth{Fqdn: fqdn}),
		}

		objMgr := ibclient.NewObjectManager(zdFakeConnector, cmpType, tenantID)

		var actualZoneAuth, getNoRef *ibclient.ZoneAuth
		getNoRef = nil
		var err error
		It("should pass expected ZoneAuth Object to GetObject", func() {
			actualZoneAuth, err = objMgr.GetZoneAuthByRef(fakeRefReturn)
		})
		It("should return expected ZoneAuth Object", func() {
			Expect(actualZoneAuth).To(Equal(zdFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
		It("should return empty ZoneAuth and nil error if ref is empty", func() {
			zdFakeConnector.getObjectObj.(*ibclient.ZoneAuth).IBBase.SetReturnFields(nil)
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

		objMgr := ibclient.NewObjectManager(zaFakeConnector, cmpType, tenantID)

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

		queryParams := ibclient.NewQueryParams(
			false,
			map[string]string{
				"fqdn": fqdn,
			})

		zdFakeConnector := &fakeConnector{
			getObjectObj:         ibclient.NewZoneDelegated(ibclient.ZoneDelegated{}),
			getObjectQueryParams: queryParams,
			getObjectRef:         "",
			resultObject:         []ibclient.ZoneDelegated{*ibclient.NewZoneDelegated(ibclient.ZoneDelegated{Fqdn: fqdn, Ref: fakeRefReturn})},
		}

		objMgr := ibclient.NewObjectManager(zdFakeConnector, cmpType, tenantID)

		var actualZoneDelegated *ibclient.ZoneDelegated
		var err error
		It("should pass expected ZoneDelegated Object to GetObject", func() {
			actualZoneDelegated, err = objMgr.GetZoneDelegated(fqdn)
		})
		It("should return expected ZoneDelegated Object", func() {
			Expect(*actualZoneDelegated).To(Equal(zdFakeConnector.resultObject.([]ibclient.ZoneDelegated)[0]))
			Expect(err).To(BeNil())
		})
		It("should return nil if fqdn is empty", func() {
			zdFakeConnector.getObjectObj.(*ibclient.ZoneDelegated).Fqdn = ""
			actualZoneDelegated, err = objMgr.GetZoneDelegated("")
			Expect(actualZoneDelegated).To(BeNil())
			Expect(err).To(BeNil())
		})
	})

	Describe("Create Zone Delegated", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		fqdn := "dzone.example.com"
		delegateTo := []ibclient.NameServer{
			{Address: "10.0.0.1", Name: "test1.dzone.example.com"},
			{Address: "10.0.0.2", Name: "test2.dzone.example.com"}}
		fakeRefReturn := "zone_delegated/ZG5zLnpvbmUkLl9kZWZhdWx0LnphLmNvLmFic2EuY2Fhcy5vaG15Z2xiLmdzbGJpYmNsaWVudA:dzone.example.com/default"
		zdFakeConnector := &fakeConnector{
			createObjectObj: ibclient.NewZoneDelegated(ibclient.ZoneDelegated{Fqdn: fqdn, DelegateTo: delegateTo}),
			resultObject:    ibclient.NewZoneDelegated(ibclient.ZoneDelegated{Fqdn: fqdn, DelegateTo: delegateTo, Ref: fakeRefReturn}),
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := ibclient.NewObjectManager(zdFakeConnector, cmpType, tenantID)

		var actualZoneDelegated *ibclient.ZoneDelegated
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
		delegateTo := []ibclient.NameServer{
			{Address: "10.0.0.1", Name: "test1.dzone.example.com"},
			{Address: "10.0.0.2", Name: "test2.dzone.example.com"}}

		receiveUpdateObject := ibclient.NewZoneDelegated(ibclient.ZoneDelegated{Ref: fakeRefReturn, DelegateTo: delegateTo})
		returnUpdateObject := ibclient.NewZoneDelegated(ibclient.ZoneDelegated{DelegateTo: delegateTo, Ref: fakeRefReturn})
		zdFakeConnector := &fakeConnector{
			fakeRefReturn:   fakeRefReturn,
			resultObject:    returnUpdateObject,
			updateObjectObj: receiveUpdateObject,
			updateObjectRef: fakeRefReturn,
		}

		objMgr := ibclient.NewObjectManager(zdFakeConnector, cmpType, tenantID)

		var updatedObject *ibclient.ZoneDelegated
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

		objMgr := ibclient.NewObjectManager(zdFakeConnector, cmpType, tenantID)

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
