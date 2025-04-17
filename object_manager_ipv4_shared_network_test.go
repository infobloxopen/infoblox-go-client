package ibclient

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Object Manager: SharedNetwork-record", func() {
	Describe("Create SharedNetwork Record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		vmID := "93f9249abc039284"
		vmName := "dummyvm"

		networkView := "default"
		recordName := "shared-network1"
		comment := "Test creation"
		networks := []*Ipv4Network{{Ref: "12.12.23.1/32"}}
		nw := []string{"12.12.23.1/32"}
		options := []*Dhcpoption{
			{
				Name:  "domain-name-servers",
				Value: "11.22.3.1",
			},
			{
				Name:  "domain-name",
				Value: "aa.mm.ee",
			},
		}
		fakeRefReturn := fmt.Sprintf("record:SharedNetwork/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)
		ea := EA{"VM Name": vmName, "VM ID": vmID}
		conn := &fakeConnector{
			createObjectObj:      NewIpv4SharedNetwork("", recordName, networks, ea, comment, false, true, options),
			getObjectObj:         NewEmptyIpv4SharedNetwork(),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         NewIpv4SharedNetwork(fakeRefReturn, recordName, networks, ea, comment, false, true, options),
			fakeRefReturn:        fakeRefReturn,
		}
		conn.createObjectObj.(*SharedNetwork).NetworkView = networkView
		conn.resultObject.(*SharedNetwork).NetworkView = networkView
		conn.resultObject.(*SharedNetwork).Networks[0].NetworkView = networkView

		objMgr := NewObjectManager(conn, cmpType, tenantID)

		var actualRecord *SharedNetwork
		var err error
		It("should pass expected sharedNetwork record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateIpv4SharedNetwork(recordName, nw, networkView, ea, comment, false, true, options)
		})
		It("should return expected sharedNetwork record Object", func() {
			Expect(actualRecord).To(Equal(conn.resultObject))
			Expect(err).To(BeNil())
		})

		// negative test case
		It("should fail to create a sharedNetwork object", func() {
			actualRecord, err := objMgr.CreateIpv4SharedNetwork("", nil, networkView, ea, comment, false, false, nil)
			Expect(actualRecord).To(BeNil())
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError("name and networks are required to create a shared network"))
		})
	})

	Describe("Get SharedNetwork Record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		recordName := "shared-network1"
		comment := "Test creation"
		networks := []*Ipv4Network{{Ref: "12.12.23.1/32"}}
		nw := "12.12.23.1/32"
		fakeRefReturn := fmt.Sprintf("record:SharedNetwork/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)
		ea := EA{"VM Name": vmName, "VM ID": vmID}
		queryParams := NewQueryParams(false, map[string]string{"name": recordName})
		conn := &fakeConnector{
			getObjectObj:         NewEmptyIpv4SharedNetwork(),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         []SharedNetwork{*NewIpv4SharedNetwork(fakeRefReturn, recordName, networks, ea, comment, false, false, nil)},
			fakeRefReturn:        fakeRefReturn,
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)

		It("should get expected sharedNetwork from getObject", func() {
			conn.getObjectQueryParams = queryParams
			actualRecord, err := objMgr.GetAllIpv4SharedNetwork(queryParams)
			Expect(err).To(BeNil())
			Expect(actualRecord).To(Equal(conn.resultObject))
		})

		// negative scenario: should fail to get a non-existent sharedNetwork with name shared-network1
		It("should fail to get expected sharedNetwork from getObject", func() {
			queryParams1 := NewQueryParams(false, map[string]string{"name": "shared-network1"})
			conn.getObjectQueryParams = queryParams1
			conn.resultObject = []SharedNetwork{}
			actualRecord, err := objMgr.GetAllIpv4SharedNetwork(queryParams1)
			Expect(err).To(BeNil())
			Expect(actualRecord).To(Equal(conn.resultObject))
		})

		// negative scenario: failed to get sharedNetwork when non-searchable field is passed
		It("should fail to get expected sharedNetwork from getObject with non-searchable fields", func() {
			queryParams2 := NewQueryParams(false, map[string]string{"networks": nw})
			conn.getObjectQueryParams = queryParams2
			conn.getObjectError = fmt.Errorf("Field is not searchable: networks")
			actualRecord, err := objMgr.GetAllIpv4SharedNetwork(queryParams2)
			Expect(err).ToNot(BeNil())
			Expect(actualRecord).To(BeNil())
		})
	})

	Describe("Update SharedNetwork Record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		vmID := "93f9249abc039284"
		vmName := "dummyvm"

		networkView := "default"
		recordName := "shared-network1"
		comment := "Test creation"
		networks := []*Ipv4Network{{Ref: "12.12.23.1/32"}}
		nw := []string{"12.12.23.1/32"}
		options := []*Dhcpoption{
			{
				Name:  "domain-name-servers",
				Value: "11.22.3.1",
			},
			{
				Name:  "domain-name",
				Value: "aa.mm.ee",
			},
		}
		updatedRef := fmt.Sprintf("record:SharedNetwork/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)
		ea := EA{"VM Name": vmName, "VM ID": vmID}
		conn := &fakeConnector{
			getObjectObj:         NewEmptyIpv4SharedNetwork(),
			getObjectRef:         updatedRef,
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         NewIpv4SharedNetwork(updatedRef, recordName, networks, ea, comment, false, true, options),
			fakeRefReturn:        updatedRef,
			updateObjectObj:      NewIpv4SharedNetwork("", recordName, networks, ea, comment, false, true, options),
			updateObjectRef:      updatedRef,
		}
		conn.updateObjectObj.(*SharedNetwork).NetworkView = networkView
		conn.updateObjectObj.(*SharedNetwork).Ref = updatedRef
		conn.resultObject.(*SharedNetwork).NetworkView = networkView
		conn.resultObject.(*SharedNetwork).Networks[0].NetworkView = networkView

		objMgr := NewObjectManager(conn, cmpType, tenantID)

		It("should pass expected sharedNetwork Object to UpdateObject", func() {
			actualRecord, err := objMgr.UpdateIpv4SharedNetwork(updatedRef, recordName, nw, networkView, comment, ea, false, true, options)
			Expect(actualRecord).To(Equal(conn.resultObject))
			Expect(err).To(BeNil())
		})

	})

	Describe("Negative case: Update SharedNetwork Record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		vmID := "93f9249abc039284"
		vmName := "dummyvm"

		networkView := "default"
		recordName := "shared-network1"
		comment := "Test creation"
		networks := []*Ipv4Network{{Ref: "12.12.23.1/32"}}
		nw := []string{"12.12.23.1/32"}
		options := []*Dhcpoption{
			{
				Name:  "domain-name-servers",
				Value: "11.22.3.1",
			},
			{
				Name:  "domain-name",
				Value: "aa.mm.ee",
			},
		}
		ea := EA{"VM Name": vmName, "VM ID": vmID}
		conn := &fakeConnector{
			getObjectError:       fmt.Errorf("not found"),
			getObjectQueryParams: NewQueryParams(false, nil),
			updateObjectError:    fmt.Errorf("not found"),
			updateObjectObj:      NewIpv4SharedNetwork("", recordName, networks, ea, comment, false, true, options),
		}
		conn.updateObjectObj.(*SharedNetwork).NetworkView = networkView
		conn.updateObjectObj.(*SharedNetwork).Networks[0].NetworkView = networkView

		objMgr := NewObjectManager(conn, cmpType, tenantID)

		// negative scenario: should fail to update a sharedNetwork object
		It("should fail to update sharedNetwork Object", func() {
			actualRecord, err := objMgr.UpdateIpv4SharedNetwork("", recordName, nw, networkView, comment, ea, false, true, options)
			Expect(actualRecord).To(BeNil())
			Expect(err).ToNot(BeNil())
		})

	})

	Describe("Delete SharedNetwork Record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		recordName := "shared-network1"

		fakeRefReturn := fmt.Sprintf("record:SharedNetwork/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)
		conn := &fakeConnector{
			fakeRefReturn:   fakeRefReturn,
			deleteObjectRef: fakeRefReturn,
		}
		objMgr := NewObjectManager(conn, cmpType, tenantID)

		var actualRef string
		var err error
		It("should pass expected sharedNetwork Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteIpv4SharedNetwork(fakeRefReturn)
		})
		It("should return expected sharedNetwork Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})

		It("should pass expected sharedNetwork Ref to DeleteObject", func() {
			deleteRef2 := "sharednetwork"
			conn.deleteObjectRef = deleteRef2
			conn.fakeRefReturn = ""
			conn.deleteObjectError = fmt.Errorf("not found")
			actualRef, err = objMgr.DeleteIpv4SharedNetwork(deleteRef2)
		})

		It("should return an error", func() {
			Expect(err).ToNot(BeNil())
		})
	})

})
