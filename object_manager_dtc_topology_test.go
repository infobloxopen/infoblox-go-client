package ibclient

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Object Manager DTC Server", func() {
	Describe("Create dtc server", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		name := "dtc_Server"
		host := "2.3.4.5"
		fakeRefReturn := fmt.Sprintf(
			"dtc:server/ZG5zLmlkbnNfc2VydmVyJGR0Y19zZXJ2ZXIuY29t:%s",
			name)
		eas := EA{"Site": "blr"}
		comment := "test client"
		monitorRef := "dtc:monitor:snmp/ZG5zLmlkbnNfbW9uaXRvcl9odHRwJGh0dHA:snmp"
		monitors := []map[string]interface{}{
			{
				"monitor": Monitor{Name: "snmp", Type: "snmp"},
				"host":    "2.3.4.5",
			},
		}
		serverMonitor := []*DtcServerMonitor{
			{
				Monitor: monitorRef,
				Host:    "2.3.4.5",
			},
		}
		sniHost := "sni_name"
		useSniHost := true
		objectAsResult := NewDtcServer(comment, name, host, false, false, eas, serverMonitor, sniHost, useSniHost)
		objectAsResult.Ref = fakeRefReturn
		aniFakeConnector := &fakeConnector{
			createObjectObj: NewDtcServer(comment, name, host, false, false, eas, serverMonitor, sniHost, useSniHost),
			getObjectRef:    fakeRefReturn,
			getObjectObj: map[string]interface{}{
				"DtcMonitor": &DtcMonitorHttp{},
			},
			getObjectQueryParams: map[string]*QueryParams{
				"DtcMonitor": NewQueryParams(false, map[string]string{"name": "snmp"}),
			},
			resultObject: map[string]interface{}{
				"DtcMonitor": []DtcMonitorHttp{{
					Ref: monitorRef,
				}},
				"DtcServer": objectAsResult,
			},
			fakeRefReturn: fakeRefReturn,
		}
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)
		var serverDtc *DtcServer
		var err error
		It("should pass expected DTC server Object to CreateObject", func() {
			serverDtc, err = objMgr.CreateDtcServer(comment, name, host, false, false, eas, monitors, sniHost, useSniHost)

		})
		It("should return expected DTC server Object", func() {
			Expect(err).To(BeNil())
			Expect(serverDtc).To(Equal(aniFakeConnector.resultObject.(map[string]interface{})["DtcServer"]))
		})
	})

	Describe("Update Dtc server", func() {
		var (
			err       error
			objMgr    IBObjectManager
			conn      *fakeConnector
			ref       string
			actualObj *DtcServer
		)

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		name := "dtc_server_20"
		It("Updating dtc_server_20, ttl, useTtl, comment and EA's", func() {
			ref = fmt.Sprintf("dtc:server/ZG5zLmlkbnNfc2VydmVyJGR0Y19zZXJ2ZXIuY29t:%s", name)
			initialEas := EA{"Site": "Blr"}
			initialComment := "test client"
			initialHost := "2.4.5.6"
			initObj := NewDtcServer(initialComment, name, initialHost, false, false, initialEas, nil, "", false)
			initObj.Ref = ref

			expectedEas := EA{"Site": "Blr"}

			updateName := "dtc_server_21"
			updateComment := "new comment"
			updateHost := "admin.com"
			updatedRef := fmt.Sprintf("dtc:server/ZG5zLmlkbnNfc2VydmVyJGR0Y19zZXJ2ZXIuY29t:%s", name)
			updateObjIn := NewDtcServer(updateComment, updateName, updateHost, true, false, expectedEas, nil, "sni_name", true)
			updateObjIn.Ref = ref

			expectedObj := NewDtcServer(updateComment, updateName, updateHost, true, false, expectedEas, nil, "sni_name", true)
			expectedObj.Ref = updatedRef

			conn = &fakeConnector{
				getObjectObj:         NewEmptyDtcServer(),
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

			actualObj, err = objMgr.UpdateDtcServer(ref, updateComment, updateName, updateHost, true, false, expectedEas, nil, "sni_name", true)
			Expect(err).To(BeNil())
			Expect(*actualObj).To(BeEquivalentTo(*expectedObj))
		})
	})

	Describe("Negative case: return an error message when sni_hostname is enabled and sni_name is not given ", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		comment := "test client"
		name := "dtc_server"
		host := "2.3.4.5"
		eas := EA{"Site": "blr"}
		conn := &fakeConnector{
			createObjectObj:   NewDtcServer(comment, name, host, true, false, eas, nil, "", true),
			createObjectError: fmt.Errorf("'sni_hostname' must be provided when 'use_sni_hostname' is enabled, and 'use_sni_hostname' must be enabled if 'sni_hostname' is provided"),
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)
		var actualRecord, expectedObj *DtcServer
		var err error
		expectedObj = nil
		It("should return expected dtc server object ", func() {
			actualRecord, err = objMgr.CreateDtcServer(comment, name, host, true, false, eas, nil, "", true)
			Expect(actualRecord).To(Equal(expectedObj))
			Expect(err).To(Equal(conn.createObjectError))
		})
	})
	Describe("Get server", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		name := "dtc_server"
		comment := "get servers"
		host := "2.3.4.5"
		sniHost := "sni_hostname"
		fakeRefReturn := fmt.Sprintf("dtc:server/ZG5zLmlkbnNfc2VydmVyJGR0Y19zZXJ2ZXIuY29t:%s", name)

		queryParams := NewQueryParams(
			false,
			map[string]string{
				"name": name,
				"host": host,
			})

		conn := &fakeConnector{
			createObjectObj:      NewDtcServer(comment, name, host, false, false, nil, nil, sniHost, true),
			getObjectRef:         "",
			getObjectObj:         NewEmptyDtcServer(),
			resultObject:         []DtcServer{*NewDtcServer(comment, name, host, false, false, nil, nil, sniHost, true)},
			fakeRefReturn:        fakeRefReturn,
			getObjectQueryParams: queryParams,
		}
		objMgr := NewObjectManager(conn, cmpType, tenantID)

		conn.resultObject.([]DtcServer)[0].Ref = fakeRefReturn

		var actualRecord *DtcServer
		var err error
		It("should pass expected Dtc Server Object to GetObject", func() {
			actualRecord, err = objMgr.GetDtcServer(name, host)
			Expect(err).To(BeNil())
			Expect(*actualRecord).To(Equal(conn.resultObject.([]DtcServer)[0]))
		})
	})
	Describe("Negative case : Error when name and host is not provided in Get function", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		name := ""
		comment := "servers"
		host := ""
		sniHost := "sni_hostname"
		fakeRefReturn := fmt.Sprintf("dtc:server/ZG5zLmlkbnNfc2VydmVyJGR0Y19zZXJ2ZXIuY29t:%s", name)

		queryParams := NewQueryParams(
			false,
			map[string]string{
				"name": name,
				"host": host,
			})

		conn := &fakeConnector{
			getObjectRef:         "",
			getObjectObj:         NewEmptyDtcServer(),
			resultObject:         []DtcServer{*NewDtcServer(comment, name, host, false, false, nil, nil, sniHost, true)},
			fakeRefReturn:        fakeRefReturn,
			getObjectQueryParams: queryParams,
			getObjectError:       fmt.Errorf("name and host of the server are required to retreive a unique dtc server"),
		}
		objMgr := NewObjectManager(conn, cmpType, tenantID)

		conn.resultObject.([]DtcServer)[0].Ref = fakeRefReturn

		var actualRecord, expectedObj *DtcServer
		var err error
		expectedObj = nil
		It("should return expected dtc server object", func() {
			actualRecord, err = objMgr.GetDtcServer(name, host)
			Expect(actualRecord).To(Equal(expectedObj))
			Expect(err).To(Equal(conn.getObjectError))
		})
	})
	Describe("Get All server", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		name := "dtc_server"
		comment := "get servers"
		host := "2.3.4.5"
		sniHost := "sni_hostname"
		fakeRefReturn := fmt.Sprintf("dtc:server/ZG5zLmlkbnNfc2VydmVyJGR0Y19zZXJ2ZXIuY29t:%s", name)

		queryParams := NewQueryParams(
			false,
			map[string]string{
				"name":         name,
				"host":         host,
				"comment":      comment,
				"sni_hostname": sniHost,
			})

		conn := &fakeConnector{
			createObjectObj:      NewDtcServer(comment, name, host, false, false, nil, nil, sniHost, true),
			getObjectRef:         "",
			getObjectObj:         NewEmptyDtcServer(),
			resultObject:         []DtcServer{*NewDtcServer(comment, name, host, false, false, nil, nil, sniHost, true)},
			fakeRefReturn:        fakeRefReturn,
			getObjectQueryParams: queryParams,
		}
		objMgr := NewObjectManager(conn, cmpType, tenantID)

		conn.resultObject.([]DtcServer)[0].Ref = fakeRefReturn

		var actualRecord []DtcServer
		var err error
		It("should pass expected Dtc Server Object to GetObject", func() {
			actualRecord, err = objMgr.GetAllDtcServer(queryParams)
			Expect(err).To(BeNil())
			Expect(actualRecord).To(Equal(conn.resultObject.([]DtcServer)))
		})
	})
	Describe("Delete DTC server", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		name := "dtc_server_20"
		deleteRef := fmt.Sprintf("dtc:server/ZG5zLmlkbnNfc2VydmVyJGR0Y19zZXJ2ZXIuY29t:%s", name)
		fakeRefReturn := deleteRef
		nwFakeConnector := &fakeConnector{
			deleteObjectRef: deleteRef,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(nwFakeConnector, cmpType, tenantID)

		var actualRef string
		var err error
		It("should pass expected DTC server Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteDtcServer(deleteRef)
		})
		It("should return expected DTC server Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})

})
