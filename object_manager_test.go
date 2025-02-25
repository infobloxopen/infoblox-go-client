package ibclient

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"reflect"
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

// when an object has internal GET calls to fetch references of dependent objects, mock data of testcases are set as map[string]interface{}
// where key is the object type and value is the object data.
// In the below example, Dtc:Lbdn object is dependent on Dtc:pool, Dtc:topology and Zone_auth objects where in CreateDtcLbdn() and UpdateDtcLdbn() functions only the name of zoneAuth, pools and topology is given,
// the references for all these 3 objects are fetched inside these functions. So the mock data for getObjectObj, resultObject and getObjectQueryParams are given as maps in fakeConnector. Example:
/*
Describe("Create Dtc Lbdn with maximum parameters", func() {
	conn := &fakeConnector{
		// 3 GET calls are made to fetch the references for DtcPool, DtcTopology and ZoneAuth objects, hence the below mock data is set as map of 3 keys, with object type as key and empty object as value
		getObjectObj: map[string]interface{}{
			"DtcPool":     &DtcPool{},
			"DtcTopology": &DtcTopology{},
			"ZoneAuth":    &ZoneAuth{},
		},
		// the query params for the GET calls are set as map with object type as key and object name as value
		getObjectQueryParams: map[string]*QueryParams{
			"DtcPool":     NewQueryParams(false, map[string]string{"name": "test-pool"}),
			"DtcTopology": NewQueryParams(false, map[string]string{"name": "test-topo"}),
			"ZoneAuth":    NewQueryParams(false, map[string]string{"fqdn": "test-zone"}),
		},
		// the result object is set as map with object type as key and object data as value, DtcPool, DctTopology and ZoneAuth objects are returned for the internal GET calls
		// and DtcLbdn object is returned for the CreateObject call
		resultObject: map[string]interface{}{
			"DtcPool": []DtcPool{{
				Ref:  poolRef,
				Name: utils.StringPtr("test-pool"),
			}},
			"DtcTopology": []DtcTopology{{
				Ref:  topologyRef,
				Name: utils.StringPtr("test-topo"),
			}},
			"DtcLbdn": NewDtcLbdn(fakeRefReturn, name, zoneAuth, comment, disable, autoConsolidatedMonitors, nil, lbMethod, patterns, persistence, createObjPools, priority, topologyRef, types, ttl, useTtl),
			"ZoneAuth": []ZoneAuth{{
				Ref:  zoneRef,
				Fqdn: zone,
			}},
		},
	}
})
*/

func (c *fakeConnector) GetObject(obj IBObject, ref string, qp *QueryParams, res interface{}) (err error) {

	if reflect.TypeOf(c.getObjectObj).Kind() == reflect.Map { //&& c.skipInternalGetcalls {
		switch obj.(type) {
		case *DtcPool:
			if ref == "" {
				*res.(*[]DtcPool) = c.resultObject.(map[string]interface{})["DtcPool"].([]DtcPool)
			} else {
				**res.(**DtcPool) = *c.resultObject.(map[string]interface{})["DtcPool"].(*DtcPool)
			}
		case *DtcTopology:
			*res.(*[]DtcTopology) = c.resultObject.(map[string]interface{})["DtcTopology"].([]DtcTopology)
		case *ZoneAuth:
			*res.(*[]ZoneAuth) = c.resultObject.(map[string]interface{})["ZoneAuth"].([]ZoneAuth)
		case *DtcServer:
			*res.(*[]DtcServer) = c.resultObject.(map[string]interface{})["DtcServer"].([]DtcServer)
		case *DtcMonitorHttp:
			*res.(*[]DtcMonitorHttp) = c.resultObject.(map[string]interface{})["DtcMonitor"].([]DtcMonitorHttp)
		case *DtcLbdn:
			**res.(**DtcLbdn) = *c.resultObject.(map[string]interface{})["DtcLbdn"].(*DtcLbdn)
		default:
			return fmt.Errorf("unsupported object type")
		}
	} else {
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
			case *RecordSRV:
				*res.(*[]RecordSRV) = c.resultObject.([]RecordSRV)
			case *RecordTXT:
				*res.(*[]RecordTXT) = c.resultObject.([]RecordTXT)
			case *ZoneDelegated:
				*res.(*[]ZoneDelegated) = c.resultObject.([]ZoneDelegated)
			case *RecordCNAME:
				*res.(*[]RecordCNAME) = c.resultObject.([]RecordCNAME)
			case *RecordA:
				*res.(*[]RecordA) = c.resultObject.([]RecordA)
			case *RecordMX:
				*res.(*[]RecordMX) = c.resultObject.([]RecordMX)
			case *ZoneForward:
				*res.(*[]ZoneForward) = c.resultObject.([]ZoneForward)
			case *DtcLbdn:
				*res.(*[]DtcLbdn) = c.resultObject.([]DtcLbdn)
			case *DtcPool:
				*res.(*[]DtcPool) = c.resultObject.([]DtcPool)
			case *DtcTopology:
				*res.(*[]DtcTopology) = c.resultObject.([]DtcTopology)
			case *DtcServer:
				*res.(*[]DtcServer) = c.resultObject.([]DtcServer)
			case *RecordAlias:
				*res.(*[]RecordAlias) = c.resultObject.([]RecordAlias)
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
			case *RecordSRV:
				**res.(**RecordSRV) = *c.resultObject.(*RecordSRV)
			case *RecordTXT:
				**res.(**RecordTXT) = *c.resultObject.(*RecordTXT)
			case *RecordCNAME:
				**res.(**RecordCNAME) = *c.resultObject.(*RecordCNAME)
			case *RecordA:
				**res.(**RecordA) = *c.resultObject.(*RecordA)
			case *RecordAAAA:
				**res.(**RecordAAAA) = *c.resultObject.(*RecordAAAA)
			case *RecordMX:
				**res.(**RecordMX) = *c.resultObject.(*RecordMX)
			case *Dns:
				*res.(*[]Dns) = c.resultObject.([]Dns)
			case *Dhcp:
				*res.(*[]Dhcp) = c.resultObject.([]Dhcp)
			case *ZoneForward:
				*res.(**ZoneForward) = c.resultObject.(*ZoneForward)
			case *ZoneDelegated:
				*res.(**ZoneDelegated) = c.resultObject.(*ZoneDelegated)
			case *DtcLbdn:
				**res.(**DtcLbdn) = *c.resultObject.(*DtcLbdn)
			case *DtcPool:
				*res.(*DtcPool) = *c.resultObject.(*DtcPool)
			case *DtcTopology:
				*res.(*DtcTopology) = *c.resultObject.(*DtcTopology)
			}
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
		delegateTo := NullableNameServers{IsNull: false, NameServers: []NameServer{
			{Address: "10.0.0.1", Name: "test1.dzone.example.com"},
			{Address: "10.0.0.2", Name: "test2.dzone.example.com"}}}
		comment := "test comment"
		disable := false
		var ea EA
		locked := false
		var delegatedTtl uint32
		useDelegatedTtl := false
		zoneFormat := "FORWARD"
		view := "default"
		fakeRefReturn := "zone_delegated/LmdzbGJpYmNsaWVudA:dzone.example.com/default"
		zdFakeConnector := &fakeConnector{
			createObjectObj: NewZoneDelegated(
				ZoneDelegated{Fqdn: fqdn, DelegateTo: delegateTo, Comment: &comment, Disable: &disable, Ea: ea, Locked: &locked,
					DelegatedTtl: &delegatedTtl, UseDelegatedTtl: &useDelegatedTtl, ZoneFormat: zoneFormat, View: &view}),
			resultObject: NewZoneDelegated(ZoneDelegated{Fqdn: fqdn, DelegateTo: delegateTo, Ref: fakeRefReturn, Comment: &comment,
				Disable: &disable, Ea: ea, Locked: &locked, DelegatedTtl: &delegatedTtl, UseDelegatedTtl: &useDelegatedTtl, ZoneFormat: zoneFormat, View: &view}),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(zdFakeConnector, cmpType, tenantID)

		var actualZoneDelegated *ZoneDelegated
		var err error
		It("should pass expected ZoneDelegated Object to CreateObject", func() {
			actualZoneDelegated, err = objMgr.CreateZoneDelegated(fqdn, delegateTo, "test comment", false, false, "", 0, false, nil, "", "")
		})
		It("should return expected ZoneDelegated Object", func() {
			Expect(actualZoneDelegated).To(Equal(zdFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Update Zone Delegated", func() {
		var (
			err          error
			objMgr       IBObjectManager
			conn         *fakeConnector
			ref          string
			actualRecord *ZoneDelegated
		)

		BeforeEach(func() {
			cmpType := "Docker"
			tenantID := "01234567890abcdef01234567890abcdef"

			ref = "zone_delegated/LmdzbGJpYmNsaWVudA:dzone.example.com/default"
			fqdn := "dzone.example.com"
			delegateTo := NullableNameServers{IsNull: false, NameServers: []NameServer{{Address: "20.20.0.1", Name: "aa.bb.com"}}}
			comment := "test comment"
			nsGroup := "testgroup"
			disable := false
			locked := false
			delegatedTtl := uint32(1800)
			useDelegatedTtl := true
			eas := EA{"Country": "test"}
			view := "default"
			zoneFormat := "FORWARD"

			newEas := EA{"Country": "new value"}
			updatedRef := "zone_delegated/ZG5zLmhvc3RjkugC4xLg:dzone.example.com/default"
			newTtl := uint32(300)

			It("updating Ttl, Extra attributes", func() {

				initObject := NewZoneDelegated(ZoneDelegated{
					Fqdn:            fqdn,
					DelegateTo:      delegateTo,
					Comment:         &comment,
					Disable:         &disable,
					Locked:          &locked,
					NsGroup:         &nsGroup,
					DelegatedTtl:    &delegatedTtl,
					UseDelegatedTtl: &useDelegatedTtl,
					Ea:              eas,
					View:            &view,
					ZoneFormat:      zoneFormat,
				})
				//initObject, _ := objMgr.CreateZoneDelegated(fqdn, delegateTo, comment, disable, locked, nsGroup, delegatedTtl, useDelegatedTtl, eas, view, zoneFormat)
				initObject.Ref = ref

				updatedObjIn := NewZoneDelegated(ZoneDelegated{
					Fqdn:            fqdn,
					DelegateTo:      delegateTo,
					Comment:         &comment,
					Disable:         &disable,
					Locked:          &locked,
					NsGroup:         &nsGroup,
					DelegatedTtl:    &newTtl,
					UseDelegatedTtl: &useDelegatedTtl,
					Ea:              newEas,
					View:            &view,
					ZoneFormat:      zoneFormat,
				})
				updatedObjIn.Ref = ref

				conn = &fakeConnector{
					getObjectObj:         NewEmptyZoneDelegated(),
					getObjectQueryParams: NewQueryParams(false, nil),
					getObjectRef:         updatedRef,
					getObjectError:       nil,

					updateObjectObj:   updatedObjIn,
					updateObjectRef:   ref,
					updateObjectError: nil,

					fakeRefReturn: updatedRef,
				}
				objMgr = NewObjectManager(conn, cmpType, tenantID)

				actualRecord, _ = objMgr.UpdateZoneDelegated(ref, delegateTo, comment, disable, locked, nsGroup, newTtl, useDelegatedTtl, newEas)
			})
			It("should return expected Zone-delegated obj", func() {
				expectedObj := NewZoneDelegated(ZoneDelegated{
					Fqdn:            fqdn,
					DelegateTo:      delegateTo,
					Comment:         &comment,
					Disable:         &disable,
					Locked:          &locked,
					NsGroup:         &nsGroup,
					DelegatedTtl:    &newTtl,
					UseDelegatedTtl: &useDelegatedTtl,
					Ea:              newEas,
					View:            &view,
					ZoneFormat:      zoneFormat,
				})
				expectedObj.Ref = updatedRef

				Expect(err).To(BeNil())
				Expect(actualRecord).NotTo(BeNil())
				Expect(*actualRecord).To(BeEquivalentTo(*expectedObj))
			})
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
