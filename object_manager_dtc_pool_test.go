package ibclient

import (
	"fmt"
	"github.com/infobloxopen/infoblox-go-client/v2/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Object Manager DTC Pool", func() {
	Describe("Create dtc pool with minimal parameters", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		comment := "test client"
		name := "dtc_pool_20"
		LbPreferredMethod := "ROUND_ROBIN"
		var fakeRefReturn = fmt.Sprintf(
			"dtc:pool/ZG5zLmlkbnNfcG9vbCRkdGNfcG9vbF8x:%s",
			name)
		eas := EA{"Site": "blr"}
		useTtl := true
		ttl := uint32(70)
		autoConsolidatedMonitors := false
		objectAsResult := NewDtcPool(comment, name, LbPreferredMethod, nil, nil, nil, nil, "", nil, nil, eas, autoConsolidatedMonitors, "", nil, ttl, true, false, 0)
		objectAsResult.Ref = fakeRefReturn
		aniFakeConnector := &fakeConnector{
			createObjectObj:      NewDtcPool(comment, name, LbPreferredMethod, nil, nil, nil, nil, "", nil, nil, eas, autoConsolidatedMonitors, "", nil, ttl, true, false, 0),
			getObjectRef:         fakeRefReturn,
			getObjectObj:         NewEmptyDtcPool(),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         objectAsResult,
			fakeRefReturn:        fakeRefReturn,
		}
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var PoolDtc *DtcPool
		var err error
		It("should pass expected DTC pool Object to CreateObject", func() {
			PoolDtc, err = objMgr.CreateDtcPool(comment, name, LbPreferredMethod, nil, nil, nil, nil, "", nil, nil, eas, autoConsolidatedMonitors, "", ttl, useTtl, false, 0)

		})
		It("should return expected DTC pool Object", func() {
			Expect(err).To(BeNil())
			Expect(PoolDtc).To(Equal(aniFakeConnector.resultObject))
		})
	})

	Describe("Create DTC pool with with TOPOLOGY as preferred load balancing method ", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		comment := "test client"
		name := "dtc_pool_20"
		var fakeRefReturn = fmt.Sprintf(
			"dtc:pool/ZG5zLmlkbnNfcG9vbCRkdGNfcG9vbF8x:%s",
			name)
		LbPreferredMethod := "TOPOLOGY"
		serverRef := "dtc:server/ZG5zLmlkbnNfc2VydmVyJGR0Y19zZXJ2ZXIuY29t:test-server"
		topologyRef := "dtc:topology/ZG5zLmhvc3QkLmNvbS5hcGkudjI6dGVzdC1wb29s:test-topo"
		topologyName := "test-topo"
		createObjServers := []*DtcServerLink{{Server: serverRef, Ratio: 3}}
		servers := []*DtcServerLink{{Server: "test-server", Ratio: 3}}
		monitor := []Monitor{{Name: "snmp", Type: "snmp"}}
		monitorRef := "dtc:monitor:snmp/ZG5zLmlkbnNfbW9uaXRvcl9odHRwJGh0dHA:snmp"
		createMonitor := []*DtcMonitorHttp{{
			Ref: monitorRef,
		}}
		eas := EA{"Site": "blr"}
		dynamicRatioPreferred := map[string]interface{}{
			"monitor":                Monitor{Name: "snmp", Type: "snmp"},
			"method":                 "MONITOR",
			"monitor_metric":         ".1.1",
			"monitor_invert_monitor": false,
		}
		lbDynamicRatioPreferred := &SettingDynamicratio{
			Method:              "MONITOR",
			Monitor:             monitorRef,
			MonitorMetric:       ".1.1",
			InvertMonitorMetric: false,
		}
		objAsResult := NewDtcPool(comment, name, LbPreferredMethod, nil, createObjServers, createMonitor, &topologyRef, "DYNAMIC_RATIO", nil, lbDynamicRatioPreferred, eas, false, "", nil, 20, true, false, 2)
		objAsResult.Ref = fakeRefReturn
		conn := &fakeConnector{
			createObjectObj: NewDtcPool(comment, name, LbPreferredMethod, nil, createObjServers, createMonitor, &topologyRef, "DYNAMIC_RATIO", nil, lbDynamicRatioPreferred, eas, false, "", nil, 20, true, false, 2),
			getObjectObj: map[string]interface{}{
				"DtcServer":   &DtcServer{},
				"DtcTopology": &DtcTopology{},
				"DtcMonitor":  &DtcMonitorHttp{},
			},
			getObjectQueryParams: map[string]*QueryParams{
				"DtcServer":   NewQueryParams(false, map[string]string{"name": "test-server"}),
				"DtcTopology": NewQueryParams(false, map[string]string{"name": "test-topo"}),
				"DtcMonitor":  NewQueryParams(false, map[string]string{"name": "snmp"}),
			},
			resultObject: map[string]interface{}{
				"DtcTopology": []DtcTopology{{
					Ref:  topologyRef,
					Name: utils.StringPtr("test-topo"),
				}},
				"DtcMonitor": []DtcMonitorHttp{{
					Ref: monitorRef,
				}},
				"DtcServer": []DtcServer{{
					Ref:  serverRef,
					Name: utils.StringPtr("test-server"),
				}},
				"DtcPool": objAsResult,
			},
			fakeRefReturn: fakeRefReturn,
		}
		objMgr := NewObjectManager(conn, cmpType, tenantID)
		It("should pass expected DtcPool Object to CreateObject", func() {
			actualRecord, err := objMgr.CreateDtcPool(comment, name, LbPreferredMethod, nil, servers, monitor, &topologyName, "DYNAMIC_RATIO", nil, dynamicRatioPreferred, eas, false, "", 20, true, false, 2)
			Expect(actualRecord).To(Equal(conn.resultObject.(map[string]interface{})["DtcPool"]))
			Expect(err).To(BeNil())
		})
	})

	Describe("Create DTC pool with with DYNAMIC_RATIO as preferred load balancing method ", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		comment := "test client"
		name := "dtc_pool_20"
		var fakeRefReturn = fmt.Sprintf(
			"dtc:pool/ZG5zLmlkbnNfcG9vbCRkdGNfcG9vbF8x:%s",
			name)
		LbPreferredMethod := "DYNAMIC_RATIO"
		serverRef := "dtc:server/ZG5zLmlkbnNfc2VydmVyJGR0Y19zZXJ2ZXIuY29t:test-server"
		monitor := []Monitor{{Name: "snmp", Type: "snmp"}}
		monitorRef := "dtc:monitor:snmp/ZG5zLmlkbnNfbW9uaXRvcl9odHRwJGh0dHA:snmp"
		createMonitor := []*DtcMonitorHttp{{
			Ref: monitorRef,
		}}
		dynamicRatioPreferred := map[string]interface{}{
			"monitor":                Monitor{Name: "snmp", Type: "snmp"},
			"method":                 "MONITOR",
			"monitor_metric":         ".1.1",
			"monitor_invert_monitor": false,
		}
		lbDynamicRatioPreferred := &SettingDynamicratio{
			Method:              "MONITOR",
			Monitor:             monitorRef,
			MonitorMetric:       ".1.1",
			InvertMonitorMetric: false,
		}
		createObjServers := []*DtcServerLink{{Server: serverRef, Ratio: 3}}
		servers := []*DtcServerLink{{Server: "test-server", Ratio: 3}}
		objAsResult := NewDtcPool(comment, name, LbPreferredMethod, lbDynamicRatioPreferred, createObjServers, createMonitor, nil, "", nil, nil, nil, false, "", nil, 20, true, false, 2)
		objAsResult.Ref = fakeRefReturn
		conn := &fakeConnector{
			createObjectObj: NewDtcPool(comment, name, LbPreferredMethod, lbDynamicRatioPreferred, createObjServers, createMonitor, nil, "", nil, nil, nil, false, "", nil, 20, true, false, 2),
			getObjectObj: map[string]interface{}{
				"DtcServer":  &DtcServer{},
				"DtcMonitor": &DtcMonitorHttp{},
			},
			getObjectQueryParams: map[string]*QueryParams{
				"DtcServer":  NewQueryParams(false, map[string]string{"name": "test-server"}),
				"DtcMonitor": NewQueryParams(false, map[string]string{"name": "snmp"}),
			},
			resultObject: map[string]interface{}{
				"DtcServer": []DtcServer{{
					Ref:  serverRef,
					Name: utils.StringPtr("test-server"),
				}},
				"DtcMonitor": []DtcMonitorHttp{{
					Ref: monitorRef,
				}},
				"DtcPool": objAsResult,
			},
			fakeRefReturn: fakeRefReturn,
		}
		objMgr := NewObjectManager(conn, cmpType, tenantID)
		It("should pass expected DtcPool Object to CreateObject", func() {
			actualRecord, err := objMgr.CreateDtcPool(comment, name, LbPreferredMethod, dynamicRatioPreferred, servers, monitor, nil, "", nil, nil, nil, false, "", 20, true, false, 2)
			Expect(actualRecord).To(Equal(conn.resultObject.(map[string]interface{})["DtcPool"]))
			Expect(err).To(BeNil())
		})
	})

	Describe("Negative case : Create dtc pool gives error when all the required fields are not passed ", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		comment := "test client"
		name := "dtc_pool_20"
		var fakeRefReturn = fmt.Sprintf(
			"dtc:pool/ZG5zLmlkbnNfcG9vbCRkdGNfcG9vbF8x:%s",
			name)
		eas := EA{"Site": "blr"}
		useTtl := true
		ttl := uint32(70)
		autoconsolidatedMonitors := false
		objectAsResult := NewDtcPool(comment, name, "", nil, nil, nil, nil, "", nil, nil, eas, autoconsolidatedMonitors, "", nil, ttl, true, false, 0)
		objectAsResult.Ref = fakeRefReturn
		aniFakeConnector := &fakeConnector{
			createObjectObj:      NewDtcPool(comment, name, "", nil, nil, nil, nil, "", nil, nil, eas, autoconsolidatedMonitors, "", nil, ttl, true, false, 0),
			getObjectRef:         fakeRefReturn,
			getObjectObj:         NewEmptyDtcPool(),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         objectAsResult,
			fakeRefReturn:        fakeRefReturn,
		}
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		_, err := objMgr.CreateDtcPool(comment, name, "", nil, nil, nil, nil, "", nil, nil, eas, autoconsolidatedMonitors, "", ttl, useTtl, false, 0)

		It("Should throw expected error when all the fields are not provided", func() {
			Expect(err).ToNot(BeNil())
		})
	})

	Describe("Negative case: return an error message when preferred method is DYNAMIC_RATIO and required parameters are not provided ", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		comment := "test client"
		name := "dtc_pool_20"
		eas := EA{"Site": "blr"}
		useTtl := true
		ttl := uint32(70)
		lbPreferredMethod := "DYNAMIC_RATIO"
		autoConsolidatedMonitors := false
		conn := &fakeConnector{
			createObjectObj:   NewDtcPool(comment, name, lbPreferredMethod, nil, nil, nil, nil, "", nil, nil, eas, autoConsolidatedMonitors, "", nil, ttl, true, false, 0),
			createObjectError: fmt.Errorf("LbDynamicRatioPreferred cannot be nil when LbPreferredMethod is set to DYNAMIC_RATIO"),
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)
		var actualRecord, expectedObj *DtcPool
		var err error
		expectedObj = nil
		It("should pass expected DTC pool Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateDtcPool(comment, name, lbPreferredMethod, nil, nil, nil, nil, "", nil, nil, eas, autoConsolidatedMonitors, "", ttl, useTtl, false, 0)
			Expect(actualRecord).To(Equal(expectedObj))
			Expect(err).To(Equal(conn.createObjectError))
		})
	})
	Describe("Update Dtc pool", func() {
		var (
			err       error
			objMgr    IBObjectManager
			conn      *fakeConnector
			ref       string
			actualObj *DtcPool
		)

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		name := "dtc_pool_20"
		autoConsolidatedMonitors := false
		It("Updating dtc_pool_20, ttl, useTtl, comment ,LB preferred method and EAs", func() {
			ref = fmt.Sprintf("dtc:pool/ZG5zLmlkbnNfcG9vbCRkdGNfcG9vbF8x:%s", name)
			initialEas := EA{"Site": "Blr"}
			initialLbPreferredMethod := "ROUND_ROBIN"
			initialTtl := uint32(70)
			initObj := NewDtcPool("old comment", name, initialLbPreferredMethod, nil, nil, nil, nil, "", nil, nil, initialEas, autoConsolidatedMonitors, "", nil, initialTtl, false, false, 0)
			initObj.Ref = ref

			expectedEas := EA{"Site": "Blr"}

			updateName := "dtc_pool_21"
			updateComment := "new comment"
			updateUseTtl := true
			updateTtl := uint32(10)
			updateLbPreferredMethod := "ALL_AVAILABLE"
			updatedRef := fmt.Sprintf("dtc:pool/ZG5zLmlkbnNfcG9vbCRkdGNfcG9vbF8x:%s", name)
			monitor := []Monitor{{Name: "snmp", Type: "snmp"}}
			monitorRef := "dtc:monitor:snmp/ZG5zLmlkbnNfbW9uaXRvcl9odHRwJGh0dHA:snmp"
			createMonitor := []*DtcMonitorHttp{{
				Ref: monitorRef,
			}}
			consolidatedMonitor := []map[string]interface{}{
				{
					"monitor":                   Monitor{Name: "snmp", Type: "snmp"},
					"availability":              "ALL",
					"full_health_communication": false,
					"members":                   []string{"infoblox.localdomain"},
				},
			}
			consolidatedMonitorStruct := []*DtcPoolConsolidatedMonitorHealth{
				{
					Monitor:                 monitorRef,
					Members:                 []string{"infoblox.localdomain"},
					Availability:            "ALL",
					FullHealthCommunication: false,
				},
			}
			updateObjIn := NewDtcPool(updateComment, updateName, updateLbPreferredMethod, nil, nil, createMonitor, nil, "", nil, nil, expectedEas, autoConsolidatedMonitors, "", consolidatedMonitorStruct, updateTtl, updateUseTtl, false, 0)
			updateObjIn.Ref = ref

			expectedObj := NewDtcPool(updateComment, updateName, updateLbPreferredMethod, nil, nil, createMonitor, nil, "", nil, nil, expectedEas, autoConsolidatedMonitors, "", consolidatedMonitorStruct, updateTtl, updateUseTtl, false, 0)
			expectedObj.Ref = updatedRef

			conn = &fakeConnector{
				getObjectObj: map[string]interface{}{
					"DtcMonitor": &DtcMonitorHttp{},
					"DtcPool":    NewEmptyDtcPool(),
				},
				getObjectQueryParams: map[string]*QueryParams{
					"DtcMonitor": NewQueryParams(false, map[string]string{"name": "snmp"}),
					"DtcPool":    NewQueryParams(false, nil),
				},
				getObjectRef:   updatedRef,
				getObjectError: nil,
				resultObject: map[string]interface{}{
					"DtcMonitor": []DtcMonitorHttp{{
						Ref: monitorRef,
					}},
					"DtcPool": expectedObj,
				},
				updateObjectObj:   updateObjIn,
				updateObjectRef:   ref,
				updateObjectError: nil,

				fakeRefReturn: updatedRef,
			}
			objMgr = NewObjectManager(conn, cmpType, tenantID)

			actualObj, err = objMgr.UpdateDtcPool(ref, updateComment, updateName, updateLbPreferredMethod, nil, nil, monitor, nil, "", nil, nil, expectedEas, autoConsolidatedMonitors, "", consolidatedMonitor, updateTtl, updateUseTtl, false, 0)
			Expect(err).To(BeNil())
			Expect(*actualObj).To(BeEquivalentTo(*expectedObj))
		})
	})
	Describe("Get pool", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		name := "dtc_pool_20"
		comment := "get pools"
		autoConsolidatedMonitors := false
		fakeRefReturn := fmt.Sprintf("dtc:pool/ZG5zLmlkbnNfcG9vbCRkdGNfcG9vbF8x:%s", name)

		queryParams := NewQueryParams(
			false,
			map[string]string{
				"name":    name,
				"comment": comment,
			})

		conn := &fakeConnector{
			createObjectObj:      NewDtcPool(comment, name, "ROUND_ROBIN", nil, nil, nil, nil, "", nil, nil, nil, autoConsolidatedMonitors, "", nil, 20, true, false, 0),
			getObjectRef:         "",
			getObjectObj:         NewEmptyDtcPool(),
			resultObject:         []DtcPool{*NewDtcPool(comment, name, "ROUND_ROBIN", nil, nil, nil, nil, "", nil, nil, nil, autoConsolidatedMonitors, "", nil, 20, true, false, 0)},
			fakeRefReturn:        fakeRefReturn,
			getObjectQueryParams: queryParams,
		}
		objMgr := NewObjectManager(conn, cmpType, tenantID)

		conn.resultObject.([]DtcPool)[0].Ref = fakeRefReturn

		var actualRecord *DtcPool
		var err error
		It("should pass expected Dtc Pool Object to GetObject", func() {
			actualRecord, err = objMgr.GetDtcPool(queryParams)
			Expect(err).To(BeNil())
			Expect(*actualRecord).To(Equal(conn.resultObject.([]DtcPool)[0]))
		})
	})

	Describe("Delete DTC pool", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		name := "dtc_pool_20"
		deleteRef := fmt.Sprintf("dtc:pool/ZG5zLmlkbnNfcG9vbCRkdGNfcG9vbF8x:%s", name)
		fakeRefReturn := deleteRef
		nwFakeConnector := &fakeConnector{
			deleteObjectRef: deleteRef,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(nwFakeConnector, cmpType, tenantID)

		var actualRef string
		var err error
		It("should pass expected DTC pool Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteDtcPool(deleteRef)
		})
		It("should return expected DTC pool Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})
})
