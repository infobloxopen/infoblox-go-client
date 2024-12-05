package ibclient

import (
	"fmt"
	"github.com/infobloxopen/infoblox-go-client/v2/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Object Manager: Dtc Lbdn", func() {
	Describe("Create Dtc Lbdn with minimum params", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		comment := "test lbdn"
		disable := false
		autoConsolidatedMonitors := false
		name := "TestLbdn1"
		fakeRefReturn := fmt.Sprintf("dtc:lbdn/ZG5zLmhvc3QkLZhd3QuaDE:%s", name)
		lbMethod := "ROUND_ROBIN"
		patterns := []string{"*info.com"}
		persistence := uint32(60)
		priority := uint32(1)
		topology := ""
		types := []string{"A", "CNAME"}
		ttl := uint32(60)
		useTtl := true

		conn := &fakeConnector{
			createObjectObj:      NewDtcLbdn("", name, nil, comment, disable, autoConsolidatedMonitors, nil, lbMethod, patterns, persistence, nil, priority, topology, types, ttl, useTtl),
			getObjectObj:         &DtcLbdn{},
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         NewDtcLbdn("", name, nil, comment, disable, autoConsolidatedMonitors, nil, lbMethod, patterns, persistence, nil, priority, topology, types, ttl, useTtl),
			fakeRefReturn:        fakeRefReturn,
		}
		conn.resultObject.(*DtcLbdn).Ref = fakeRefReturn
		objMgr := NewObjectManager(conn, cmpType, tenantID)
		It("should pass expected DtcLbdn Object to CreateObject", func() {
			actualRecord, err := objMgr.CreateDtcLbdn(name, nil, comment, disable, autoConsolidatedMonitors, nil, lbMethod, patterns, persistence, nil, priority, topology, types, ttl, useTtl)
			Expect(actualRecord).To(Equal(conn.resultObject))
			Expect(err).To(BeNil())
		})

		It("should fail to create a DTC lbdn object", func() {
			actualRecord, err := objMgr.CreateDtcLbdn("", nil, comment, disable, autoConsolidatedMonitors, nil, "", patterns, persistence, nil, priority, topology, types, ttl, useTtl)
			Expect(actualRecord).To(BeNil())
			Expect(err).ToNot(BeNil())
		})
	})

	Describe("Get Dtc Lbdn", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		comment := "test lbdn"
		disable := false
		autoConsolidatedMonitors := false
		name := "TestLbdn1"
		fakeRefReturn := fmt.Sprintf("dtc:lbdn/ZG5zLmhvc3QkLZhd3QuaDE:%s", name)
		lbMethod := "ROUND_ROBIN"
		patterns := []string{"*info.com"}
		persistence := uint32(60)
		priority := uint32(1)
		topology := ""
		types := []string{"A", "CNAME"}
		ttl := uint32(60)
		useTtl := true
		queryParams := NewQueryParams(false, map[string]string{"name": name})
		res := NewDtcLbdn("", name, nil, comment, disable, autoConsolidatedMonitors, nil, lbMethod, patterns, persistence, nil, priority, topology, types, ttl, useTtl)

		conn := &fakeConnector{
			getObjectObj:  NewEmptyDtcLbdn(),
			resultObject:  []DtcLbdn{*res},
			fakeRefReturn: fakeRefReturn,
		}
		objMgr := NewObjectManager(conn, cmpType, tenantID)
		It("should get expected DtcLbdn Object from getObject", func() {
			conn.getObjectQueryParams = queryParams
			actualRecord, err := objMgr.GetDtcLbdn(queryParams)
			Expect(actualRecord).To(Equal(conn.resultObject))
			Expect(err).To(BeNil())
		})

		It("should fail to get expected DtcLbdn Object from getObject", func() {
			qp1 := NewQueryParams(false, map[string]string{"name": "test-lbdn111"})
			conn.getObjectQueryParams = qp1
			conn.resultObject = []DtcLbdn{}
			actualRecord, err := objMgr.GetDtcLbdn(qp1)
			Expect(actualRecord).To(Equal(conn.resultObject))
			Expect(err).To(BeNil())
		})

	})

	Describe("Get Dtc Lbdn: Negative scenario", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		qp2 := NewQueryParams(false, map[string]string{"lb_method": "ROUND_ROBIN"})
		conn := &fakeConnector{
			getObjectObj:         NewEmptyDtcLbdn(),
			getObjectQueryParams: qp2,
			resultObject:         []DtcLbdn{},
		}
		objMgr := NewObjectManager(conn, cmpType, tenantID)

		// negative scenario
		conn.getObjectError = fmt.Errorf("Field is not searchable: lb_method")
		It("should fail to get expected DtcLbdn Object from getObject with non searchable field", func() {
			_, err := objMgr.GetDtcLbdn(qp2)
			Expect(err).ToNot(BeNil())
		})

	})

	Describe("Create Dtc Lbdn with maximum parameters", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		comment := "test lbdn"
		disable := false
		autoConsolidatedMonitors := false
		name := "TestLbdn1"
		fakeRefReturn := fmt.Sprintf("dtc:lbdn/ZG5zLmhvc3QkLZhd3QuaDE:%s", name)
		lbMethod := "TOPOLOGY"
		patterns := []string{"*info.com"}
		persistence := uint32(60)
		pools := []*DtcPoolLink{{Pool: "test-pool", Ratio: 3}}
		priority := uint32(1)
		topology := "test-topo"
		types := []string{"A", "CNAME"}
		ttl := uint32(60)
		useTtl := true
		poolRef := "dtc:pool/ZG5zLmhvc3QkLmNvbS5hcGkudjI6dGVzdC1wb29s:test-pool"
		topologyRef := "dtc:topology/ZG5zLmhvc3QkLmNvbS5hcGkudjI6dGVzdC1wb29s:test-topo"
		createObjPools := []*DtcPoolLink{{Pool: poolRef, Ratio: 3}}
		zone := "test-zone"
		zones := []string{zone}
		zoneRef := "zone_auth/ZG5zLmhvc3QkLmNvbS5hcGkudjI6dGVzdC1wb29s:test-zone"
		zoneAuth := []*ZoneAuth{{Ref: zoneRef, Fqdn: zone}}

		conn := &fakeConnector{
			createObjectObj: NewDtcLbdn("", name, zoneAuth, comment, disable, autoConsolidatedMonitors, nil, lbMethod, patterns, persistence, createObjPools, priority, topologyRef, types, ttl, useTtl),
			getObjectObj: map[string]interface{}{
				"DtcPool":     &DtcPool{},
				"DtcTopology": &DtcTopology{},
				"ZoneAuth":    &ZoneAuth{},
			},
			getObjectQueryParams: map[string]*QueryParams{
				"DtcPool":     NewQueryParams(false, map[string]string{"name": "test-pool"}),
				"DtcTopology": NewQueryParams(false, map[string]string{"name": "test-topo"}),
				"ZoneAuth":    NewQueryParams(false, map[string]string{"fqdn": "test-zone"}),
			},
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
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)
		It("should pass expected DtcLbdn Object to CreateObject", func() {
			actualRecord, err := objMgr.CreateDtcLbdn(name, zones, comment, disable, autoConsolidatedMonitors, nil, lbMethod, patterns, persistence, pools, priority, topology, types, ttl, useTtl)
			Expect(actualRecord).To(Equal(conn.resultObject.(map[string]interface{})["DtcLbdn"]))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete Dtc Lbdn", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		recordName := "test-lbdn"
		deleteRef := fmt.Sprintf("dtc:lbdn/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)
		nwFakeConnector := &fakeConnector{
			deleteObjectRef: deleteRef,
			fakeRefReturn:   deleteRef,
		}

		objMgr := NewObjectManager(nwFakeConnector, cmpType, tenantID)

		var actualRef string
		var err error
		It("should pass expected DTC Lbdn Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteDtcLbdn(deleteRef)
		})
		It("should return expected DTC Lbdn Ref", func() {
			Expect(actualRef).To(Equal(deleteRef))
			Expect(err).To(BeNil())
		})

		It("should pass expected DTC Lbdn Ref to DeleteObject", func() {
			deleteRef2 := "dtc:lbdn"
			nwFakeConnector.deleteObjectRef = deleteRef2
			nwFakeConnector.fakeRefReturn = ""
			nwFakeConnector.deleteObjectError = fmt.Errorf("not found")
			actualRef, err = objMgr.DeleteDtcLbdn(deleteRef2)
		})

		It("should return an error", func() {
			Expect(err).ToNot(BeNil())
		})

	})

	Describe("Update Dtc Lbdn", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		disable := false
		autoConsolidatedMonitors := false
		lbMethod := "RATIO"
		patterns := []string{"*info.com"}
		persistence := uint32(60)
		priority := uint32(1)
		types := []string{"A", "CNAME"}
		ttl := uint32(60)
		useTtl := true
		name1 := "TestLbdn222"
		comment1 := "test lbdn comment updated"
		updateRef := fmt.Sprintf("dtc:lbdn/ZG5zLmhvc3QkLZhd3QuaDE:%s", name1)

		conn := &fakeConnector{
			getObjectObj:         NewEmptyDtcLbdn(),
			getObjectRef:         updateRef,
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         NewDtcLbdn(updateRef, name1, nil, comment1, disable, autoConsolidatedMonitors, nil, lbMethod, patterns, persistence, nil, priority, "", types, ttl, useTtl),
			fakeRefReturn:        updateRef,
			updateObjectObj:      NewDtcLbdn(updateRef, name1, nil, comment1, disable, autoConsolidatedMonitors, nil, lbMethod, patterns, persistence, nil, priority, "", types, ttl, useTtl),
			updateObjectRef:      updateRef,
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)
		It("should pass expected DtcLbdn Object to UpdateObject", func() {
			actualRecord, err := objMgr.UpdateDtcLbdn(updateRef, name1, nil, comment1, disable, autoConsolidatedMonitors, nil, lbMethod, patterns, persistence, nil, priority, "", types, ttl, useTtl)
			Expect(actualRecord).To(Equal(conn.resultObject))
			Expect(err).To(BeNil())
		})

	})

	Describe("Update Dtc Lbdn with maximum params", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		comment := "test lbdn"
		disable := false
		autoConsolidatedMonitors := false
		name := "TestLbdn1"
		fakeRefReturn := fmt.Sprintf("dtc:lbdn/ZG5zLmhvc3QkLZhd3QuaDE:%s", name)
		initLbMethod := "RATIO"
		lbMethod := "TOPOLOGY"
		patterns := []string{"*info.com"}
		persistence := uint32(60)
		pools := []*DtcPoolLink{{Pool: "test-pool", Ratio: 3}}
		priority := uint32(1)
		topology := "test-topo"
		types := []string{"A", "CNAME"}
		ttl := uint32(60)
		useTtl := true
		poolRef := "dtc:pool/ZG5zLmhvc3QkLmNvbS5hcGkudjI6dGVzdC1wb29s:test-pool"
		topologyRef := "dtc:topology/ZG5zLmhvc3QkLmNvbS5hcGkudjI6dGVzdC1wb29s:test-topo"
		createObjPools := []*DtcPoolLink{{Pool: poolRef, Ratio: 3}}
		zone := "test-zone"
		zones := []string{zone}
		zoneRef := "zone_auth/ZG5zLmhvc3QkLmNvbS5hcGkudjI6dGVzdC1wb29s:test-zone"
		zoneAuth := []*ZoneAuth{{Ref: zoneRef, Fqdn: zone}}

		initObject := NewDtcLbdn("", name, nil, comment, disable, autoConsolidatedMonitors, nil, initLbMethod, patterns, persistence, nil, priority, "", types, ttl, useTtl)
		initObject.Ref = fakeRefReturn
		conn := &fakeConnector{
			getObjectObj: map[string]interface{}{
				"DtcPool":     &DtcPool{},
				"DtcTopology": &DtcTopology{},
				"ZoneAuth":    &ZoneAuth{},
			},
			getObjectQueryParams: map[string]*QueryParams{
				"DtcPool":     NewQueryParams(false, map[string]string{"name": "test-pool"}),
				"DtcTopology": NewQueryParams(false, map[string]string{"name": "test-topo"}),
				"ZoneAuth":    NewQueryParams(false, map[string]string{"fqdn": "test-zone"}),
			},
			updateObjectRef: fakeRefReturn,
			updateObjectObj: NewDtcLbdn(fakeRefReturn, name, zoneAuth, comment, disable, autoConsolidatedMonitors, nil, lbMethod, patterns, persistence, createObjPools, priority, topologyRef, types, ttl, useTtl),
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
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)
		It("should pass expected DtcLbdn Object to CreateObject", func() {

			actualRecord, err := objMgr.UpdateDtcLbdn(fakeRefReturn, name, zones, comment, disable, autoConsolidatedMonitors, nil, lbMethod, patterns, persistence, pools, priority, topology, types, ttl, useTtl)
			Expect(actualRecord).To(Equal(conn.resultObject.(map[string]interface{})["DtcLbdn"]))
			Expect(err).To(BeNil())
		})
	})

	Describe("Update Dtc Lbdn with, negative scenario", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		disable := false
		autoConsolidatedMonitors := false
		lbMethod := "RATIO"
		patterns := []string{"*info.com"}
		persistence := uint32(60)
		priority := uint32(1)
		types := []string{"A", "CNAME"}
		ttl := uint32(60)
		useTtl := true

		name2 := "test-LBDN122"
		comment2 := "comment updated"
		oldRef := "dtc:lbdn/ZG5zLmhvc3QkLZhd3QuaDE:test-lbdn121"

		conn := &fakeConnector{
			getObjectObj:         NewEmptyDtcLbdn(),
			getObjectRef:         oldRef,
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         NewDtcLbdn(oldRef, name2, nil, comment2, disable, autoConsolidatedMonitors, nil, lbMethod, patterns, persistence, nil, priority, "", types, ttl, useTtl),
			getObjectError:       fmt.Errorf("not found"),
			fakeRefReturn:        oldRef,
			updateObjectObj:      NewDtcLbdn(oldRef, name2, nil, comment2, disable, autoConsolidatedMonitors, nil, lbMethod, patterns, persistence, nil, priority, "", types, ttl, useTtl),
			updateObjectRef:      oldRef,
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)
		// negative scenario

		It("should fail to update DtcLbdn Object", func() {
			actualRecord, err := objMgr.UpdateDtcLbdn(oldRef, name2, nil, comment2, disable, autoConsolidatedMonitors, nil, lbMethod, patterns, persistence, nil, priority, "", types, ttl, useTtl)
			Expect(actualRecord).To(BeNil())
			Expect(err).ToNot(BeNil())
		})

	})

})
