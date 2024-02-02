package ibclient_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	ibclient "github.com/infobloxopen/infoblox-go-client/v2"
	"github.com/infobloxopen/infoblox-go-client/v2/utils"
)

var _ = Describe("Object Manager: SRV-Record", func() {
	Describe("Create SRV Record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		dnsView := "default"

		name := "_srv._proto.example.com"
		priority := uint32(10)
		weight := uint32(25)
		port := uint32(88)
		target := "h1.example.com"
		ttl := uint32(70)
		useTtl := true
		comment := "this is a test comment"

		fakeRefReturn := fmt.Sprintf("record:srv/ZG5zLmhvc3RjkuMC4xLg:%s/%s", name, dnsView)
		eas := ibclient.EA{
			"VM ID":   "93f9249abc039284",
			"VM Name": "dummyvm",
		}
		aniFakeConnector := &fakeConnector{
			createObjectObj: ibclient.NewRecordSRV(ibclient.RecordSRV{
				View:     dnsView,
				Name:     &name,
				Priority: &priority,
				Weight:   &weight,
				Port:     &port,
				Target:   &target,
				Ttl:      &ttl,
				UseTtl:   &useTtl,
				Comment:  &comment,
				Ea:       eas,
			}),
			getObjectRef: fakeRefReturn,
			getObjectObj: ibclient.NewRecordSRV(ibclient.RecordSRV{
				View: dnsView,
				Name: &name,
				Ref:  fakeRefReturn,
			}),
			resultObject: ibclient.NewRecordSRV(ibclient.RecordSRV{
				View:     dnsView,
				Name:     &name,
				Priority: &priority,
				Weight:   &weight,
				Port:     &port,
				Target:   &target,
				Ttl:      &ttl,
				UseTtl:   &useTtl,
				Ref:      fakeRefReturn,
				Comment:  &comment,
				Ea:       eas,
			}),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := ibclient.NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualRecord *ibclient.RecordSRV
		var err error
		It("should pass expected SRV record object to CreateObject", func() {
			actualRecord, err = objMgr.CreateSRVRecord(dnsView, name, priority, weight, port, target, ttl, useTtl, comment, eas)
		})
		It("should return expected SRV record Object", func() {
			Expect(err).To(BeNil())
			Expect(actualRecord).NotTo(BeNil())
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
		})
	})

	Describe("Update SRV Record", func() {
		var (
			err          error
			objMgr       ibclient.IBObjectManager
			conn         *fakeConnector
			ref          string
			actualRecord *ibclient.RecordSRV
		)

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		dnsView := "default"

		name := "_srv._proto.example.com"
		ref = fmt.Sprintf("record:srv/ZG5zLmhvc3RjkuMC4xLg:%s/%s", name, dnsView)

		newEas := ibclient.EA{"Country": "new value"}
		updateName := "_srv2._proto.example.com"
		updateRef := fmt.Sprintf("record:srv/ZG5zLmhvc3RjkugC4xLg:%s/%s", updateName, dnsView)
		updateObjIn := ibclient.NewRecordSRV(ibclient.RecordSRV{
			Ref:      ref,
			View:     "",
			Name:     &updateName,
			Priority: utils.Uint32Ptr(20),
			Weight:   utils.Uint32Ptr(30),
			Port:     utils.Uint32Ptr(88),
			Target:   utils.StringPtr("h3.example.com"),
			Ttl:      utils.Uint32Ptr(100),
			UseTtl:   utils.BoolPtr(true),
			Comment:  utils.StringPtr("new comment"),
			Ea:       newEas,
		})

		expectedObj := ibclient.NewRecordSRV(ibclient.RecordSRV{
			Ref:      updateRef,
			Name:     &updateName,
			Priority: utils.Uint32Ptr(20),
			Weight:   utils.Uint32Ptr(30),
			Port:     utils.Uint32Ptr(88),
			Target:   utils.StringPtr("h3.example.com"),
			Ttl:      utils.Uint32Ptr(100),
			UseTtl:   utils.BoolPtr(true),
			Comment:  utils.StringPtr("new comment"),
			Ea:       newEas,
		})

		conn = &fakeConnector{
			updateObjectObj:   updateObjIn,
			updateObjectRef:   ref,
			updateObjectError: nil,

			fakeRefReturn: updateRef,
		}
		objMgr = ibclient.NewObjectManager(conn, cmpType, tenantID)
		It("should pass updated SRV record arguments", func() {
			actualRecord, err = objMgr.UpdateSRVRecord(
				ref,
				updateName,
				uint32(20),
				uint32(30),
				uint32(88),
				"h3.example.com",
				uint32(100), true,
				"new comment",
				newEas)
		})
		It("should return expected SRV record obj", func() {
			Expect(err).To(BeNil())
			Expect(actualRecord).NotTo(BeNil())
			Expect(actualRecord).To(BeEquivalentTo(expectedObj))
		})
	})

	Describe("Get SRV Record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		dnsView := "default"
		name := "_srv._proto.example.com"
		priority := uint32(10)
		weight := uint32(25)
		port := uint32(88)
		target := "h1.example.com"
		ttl := uint32(70)
		useTtl := true
		comment := "this is a test comment"
		fakeRefReturn := fmt.Sprintf("record:srv/ZG5zLmhvc3RjkugC4xLg:%s/%s", name, dnsView)

		sf := map[string]string{
			"view":   dnsView,
			"name":   name,
			"target": fmt.Sprintf("%s", target),
			"port":   fmt.Sprintf("%d", port),
		}
		queryParams := ibclient.NewQueryParams(false, sf)

		nwFakeConnector := &fakeConnector{
			getObjectObj:         ibclient.NewEmptyRecordSRV(),
			getObjectQueryParams: queryParams,

			resultObject: []ibclient.RecordSRV{*ibclient.NewRecordSRV(ibclient.RecordSRV{
				View:     dnsView,
				Name:     &name,
				Priority: &priority,
				Weight:   &weight,
				Port:     &port,
				Target:   &target,
				Ttl:      &ttl,
				UseTtl:   &useTtl,
				Comment:  &comment,
				Ref:      fakeRefReturn,
			})},
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := ibclient.NewObjectManager(nwFakeConnector, cmpType, tenantID)

		var actualRecord *ibclient.RecordSRV
		var err error
		It("should pass expected dnsview, name to GetObject", func() {
			actualRecord, err = objMgr.GetSRVRecord(dnsView, name, target, port)
		})
		It("should return expected SRV record Object", func() {
			Expect(err).To(BeNil())
			Expect(actualRecord).NotTo(BeNil())
			Expect(*actualRecord).To(Equal(nwFakeConnector.resultObject.([]ibclient.RecordSRV)[0]))
		})
	})

	Describe("Delete SRV Record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		name := "_srv._proto.example.com"
		dnsView := "default"
		deleteRef := fmt.Sprintf("record:srv/ZG5zLmhvc3RjkugC4xLg:%s/%s", name, dnsView)
		fakeRefReturn := deleteRef
		nwFakeConnector := &fakeConnector{
			deleteObjectRef: deleteRef,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := ibclient.NewObjectManager(nwFakeConnector, cmpType, tenantID)

		var actualRef string
		var err error
		It("should pass expected SRV Record ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteSRVRecord(deleteRef)
		})
		It("should return expected SRV Record Ref", func() {
			Expect(err).To(BeNil())
			Expect(actualRef).To(Equal(fakeRefReturn))
		})
	})
})
