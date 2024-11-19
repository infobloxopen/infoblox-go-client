package ibclient

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Object Manager: Dtc Lbdn", func() {
	Describe("Create Dtc Lbdn", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		comment := "test lbdn"
		disable := false
		autoConsolidatedMonitors := false
		name := "TestLbdn1"
		fakeRefReturn := fmt.Sprintf("dtc:lbdn/ZG5zLmhvc3QkLZhd3QuaDE:%s", name)
		//authRef := "zone_auth/ZG5zLnpvbmUkLl9kZWZhdWx0LnphLmNvLmFic2EuY2Fhcy5vaG15Z2xiLmdzbGJpYmNsaWVudA:info.com/default"
		//authZones := []*ZoneAuth{{Fqdn: "info.com", Ref: authRef}}
		//authZones := []*ZoneAuth{}
		lbMethod := "TOPOLOGY"
		patterns := []string{"*info.com"}
		persistence := uint32(60)
		pools := []*DtcPoolLink{{Pool: "test-pool", Ratio: 3}}
		priority := uint32(1)
		topology := "test-topo"
		types := []string{"A", "CNAME"}
		ttl := uint32(60)
		useTtl := true
		//poolRef := "dtc:pool/ZG5zLmhvc3QkLmNvbS5hcGkudjI6dGVzdC1wb29s:test-pool"

		conn := &fakeConnector{
			createObjectObj: NewDtcLbdn("", name, nil, comment, disable, autoConsolidatedMonitors, nil, lbMethod, patterns, persistence, pools, priority, topology, types, ttl, useTtl),
			//getObjectRef:         poolRef,
			getObjectObj:         &DtcPool{},
			getObjectQueryParams: NewQueryParams(false, map[string]string{"name": "test-pool"}),
			//resultObject:         NewDtcLbdn("", name, nil, comment, disable, autoConsolidatedMonitors, nil, lbMethod, patterns, persistence, pools, priority, topology, types, ttl, useTtl),
			resultObject:  []DtcPool{},
			fakeRefReturn: fakeRefReturn,
		}
		//conn.createObjectObj.(*DtcPool).Ref = fakeRefReturn
		//conn.getObjectObj.(*DtcPool).Ref = poolRef
		//conn.getObjectObj.(*DtcPool).Name = utils.StringPtr("test-pool")
		objMgr := NewObjectManager(conn, cmpType, tenantID)
		It("should pass expected DtcLbdn Object to CreateObject", func() {
			actualRecord, err := objMgr.CreateDtcLbdn(name, nil, comment, disable, autoConsolidatedMonitors, nil, lbMethod, patterns, persistence, pools, priority, topology, types, ttl, useTtl)
			conn.getObjectObj = NewEmptyDtcLbdn()
			Expect(actualRecord).To(Equal(conn.resultObject))
			Expect(err).To(BeNil())
		})
	})

})
