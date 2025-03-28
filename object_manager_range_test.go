package ibclient

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Object Manager: Network Range", func() {
	Describe("Create a Network Range with minimal parameters", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "default"
		network := "12.4.0.0/24"
		startAddr := "12.4.0.120"
		endAddr := "12.4.0.130"
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		fakeRefReturn := fmt.Sprintf(
			"range/ZG5zLmRoY3BfcmFuZ2UkMTIuNC4wLjEyMC8xMi40LjAuMTMwLy8vMC8:%s/%s/%s",
			startAddr,
			endAddr, netviewName)
		eas := make(EA)
		eas["VM ID"] = vmID
		eas["VM Name"] = vmName
		objectForCreation := NewRange(
			"", "", network, startAddr, eas, false, nil, false, endAddr, "", nil, "NONE")
		objectForCreation.NetworkView = &netviewName
		objectAsResult := NewRange(
			"", "", network, startAddr, eas, false, nil, false, endAddr, "", nil, "NONE")
		objectAsResult.NetworkView = &netviewName
		objectAsResult.Ref = fakeRefReturn
		aniFakeConnector := &fakeConnector{
			createObjectObj:      objectForCreation,
			getObjectRef:         fakeRefReturn,
			getObjectObj:         NewEmptyRange(),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         objectAsResult,
			fakeRefReturn:        fakeRefReturn,
		}
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)
		var actualRecord *Range
		var err error
		It("should pass expected Network Range Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateNetworkRange("", "", network, netviewName, startAddr, endAddr, false, eas, nil, "", nil, false, "NONE")
		})
		It("should return expected Network Range Object", func() {
			Expect(err).To(BeNil())
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
		})
	})
	Describe("Create a Network Range with maximal parameters", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "default"
		network := "12.4.0.0/24"
		startAddr := "12.4.0.120"
		endAddr := "12.4.0.130"
		comment := "create a range"
		name := "range for 12.4.0.0/24"
		options := []*Dhcpoption{
			{
				Name:        "routers",
				Num:         3,
				Value:       "12.4.0.23",
				VendorClass: "DHCP",
				UseOption:   true,
			},
		}
		member := &Dhcpmember{
			Ipv4Addr: "10.197.81.120",
			Ipv6Addr: "2403:8600:80cf:e10c:3a00::1192",
			Name:     "infoblox.localdomain",
		}
		serverAssociationType := "MEMBER"
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		fakeRefReturn := fmt.Sprintf(
			"range/ZG5zLmRoY3BfcmFuZ2UkMTIuNC4wLjEyMC8xMi40LjAuMTMwLy8vMC8:%s/%s/%s",
			startAddr,
			endAddr, netviewName)
		eas := make(EA)
		eas["VM ID"] = vmID
		eas["VM Name"] = vmName
		objectForCreation := NewRange(
			comment, name, network, startAddr, eas, false, options, false, endAddr, "", member, serverAssociationType)
		objectForCreation.NetworkView = &netviewName
		objectAsResult := NewRange(
			comment, name, network, startAddr, eas, false, options, false, endAddr, "", member, serverAssociationType)
		objectAsResult.Ref = fakeRefReturn
		objectAsResult.NetworkView = &netviewName
		aniFakeConnector := &fakeConnector{
			createObjectObj:      objectForCreation,
			getObjectRef:         fakeRefReturn,
			getObjectObj:         NewEmptyRange(),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         objectAsResult,
			fakeRefReturn:        fakeRefReturn,
		}
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)
		var actualRecord *Range
		var err error
		It("should pass expected Network Range Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateNetworkRange(comment, name, network, netviewName, startAddr, endAddr, false, eas, member, "", options, false, serverAssociationType)
		})
		It("should return expected Network Range Object", func() {
			Expect(err).To(BeNil())
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
		})
	})
	Describe("Get Network Range", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "default"
		network := "12.4.0.0/24"
		startAddr := "12.4.0.120"
		endAddr := "12.4.0.130"
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		comment := "create a range"
		name := "range for 12.0.0.0/24"
		fakeRefReturn := fmt.Sprintf(
			"range/ZG5zLmRoY3BfcmFuZ2UkMTIuNC4wLjEyMC8xMi40LjAuMTMwLy8vMC8:%s/%s/%s",
			startAddr,
			endAddr, netviewName)
		eas := make(EA)
		eas["VM ID"] = vmID
		eas["VM Name"] = vmName
		options := []*Dhcpoption{
			{
				Name:        "routers",
				Num:         3,
				Value:       "12.4.0.23",
				VendorClass: "DHCP",
				UseOption:   true,
			},
		}
		failOverAssociation := "failOver"
		queryParams := NewQueryParams(
			false,
			map[string]string{
				"network": network,
			})
		conn := &fakeConnector{
			getObjectRef:         "",
			getObjectObj:         NewEmptyRange(),
			getObjectQueryParams: queryParams,
			resultObject:         []Range{*NewRange(comment, name, network, startAddr, eas, false, options, true, endAddr, failOverAssociation, nil, "FAILOVER")},
			fakeRefReturn:        fakeRefReturn,
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)

		conn.resultObject.([]Range)[0].Ref = fakeRefReturn

		var actualRecord []Range
		var err error
		It("should pass expected Network Range to GetObject", func() {
			actualRecord, err = objMgr.GetNetworkRange(queryParams)
			Expect(err).To(BeNil())
			Expect(actualRecord).To(Equal(conn.resultObject.([]Range)))
		})
	})
	Describe("Delete Network Range", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "default"
		startAddr := "12.4.0.120"
		endAddr := "12.4.0.130"
		deleteRef := fmt.Sprintf("range/ZG5zLmRoY3BfcmFuZ2UkMTIuNC4wLjEyMC8xMi40LjAuMTMwLy8vMC8:%s/%s/%s", startAddr, endAddr, netviewName)
		nwFakeConnector := &fakeConnector{
			deleteObjectRef: deleteRef,
			fakeRefReturn:   deleteRef,
		}

		objMgr := NewObjectManager(nwFakeConnector, cmpType, tenantID)

		var actualRef string
		var err error
		It("should pass expected Network Range Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteNetworkRange(deleteRef)
		})
		It("should return expected Network Range Ref", func() {
			Expect(actualRef).To(Equal(deleteRef))
			Expect(err).To(BeNil())
		})

		It("should pass expected Network Range Ref to DeleteObject", func() {
			deleteRef2 := "range"
			nwFakeConnector.deleteObjectRef = deleteRef2
			nwFakeConnector.fakeRefReturn = ""
			nwFakeConnector.deleteObjectError = fmt.Errorf("not found")
			actualRef, err = objMgr.DeleteNetworkRange(deleteRef2)
		})

		It("should return an error", func() {
			Expect(err).ToNot(BeNil())
		})

	})
	Describe("Update Network range", func() {
		var (
			err       error
			objMgr    IBObjectManager
			conn      *fakeConnector
			ref       string
			actualObj *Range
		)

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		name := "network_range_1"
		netviewName := "default"
		startAddr := "12.4.0.120"
		endAddr := "12.4.0.130"
		network := "12.4.0.0/24"
		It("Updating Member association to failOver association ", func() {
			ref = fmt.Sprintf("range/ZG5zLmRoY3BfcmFuZ2UkMTIuNC4wLjEyMC8xMi40LjAuMTMwLy8vMC8:%s/%s/%s", startAddr, endAddr, netviewName)
			initialEas := EA{"Site": "Blr"}
			initalOptions := []*Dhcpoption{
				{
					Name:        "routers",
					Num:         3,
					Value:       "12.4.0.23",
					VendorClass: "DHCP",
					UseOption:   true,
				},
			}
			initialMember := &Dhcpmember{
				Ipv4Addr: "10.197.81.120",
				Ipv6Addr: "2403:8600:80cf:e10c:3a00::1192",
				Name:     "infoblox.localdomain",
			}
			initialServerAssociationType := "MEMBER"
			initialComment := "old comment"
			initObj := NewRange(initialComment, name, network, startAddr, initialEas, false, initalOptions, true, endAddr, "", initialMember, initialServerAssociationType)
			initObj.Ref = ref

			expectedEas := EA{"Site": "Blr"}

			updateName := "network_range_2"
			updateComment := "new comment"
			updateServerAssociationType := "FAILOVER"
			updateFailOverAssociation := "failOver"
			updatedRef := fmt.Sprintf("range/ZG5zLmRoY3BfcmFuZ2UkMTIuNC4wLjEyMC8xMi40LjAuMTMwLy8vMC8:%s/%s/%s", startAddr, endAddr, netviewName)
			updateObjIn := NewRange(updateComment, updateName, network, startAddr, expectedEas, true, initalOptions, true, endAddr, updateFailOverAssociation, nil, updateServerAssociationType)
			updateObjIn.Ref = ref

			expectedObj := NewRange(updateComment, updateName, network, startAddr, expectedEas, true, initalOptions, true, endAddr, updateFailOverAssociation, nil, updateServerAssociationType)
			expectedObj.Ref = updatedRef

			getObjIn := NewEmptyRange()
			conn = &fakeConnector{
				getObjectObj:         getObjIn,
				getObjectQueryParams: NewQueryParams(false, nil),
				getObjectRef:         updatedRef,
				getObjectError:       nil,
				resultObject:         expectedObj,
				updateObjectObj:      updateObjIn,
				updateObjectRef:      ref,
				updateObjectError:    nil,

				fakeRefReturn: updatedRef,
			}
			objMgr = NewObjectManager(conn, cmpType, tenantID)

			actualObj, err = objMgr.UpdateNetworkRange(ref, updateComment, updateName, network, startAddr, endAddr, true, expectedEas, nil, updateFailOverAssociation, initalOptions, true, updateServerAssociationType)
			Expect(err).To(BeNil())
			Expect(*actualObj).To(BeEquivalentTo(*expectedObj))
		})
	})
	Describe("Update Network range with, negative scenario", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		//netviewName := "default"
		network := "12.4.0.0/24"
		name2 := "range_create_2"
		comment2 := "comment updated"
		oldRef := "range/ZG5zLmRoY3BfcmFuZ2UkMTIuNC4wLjEyMC8xMi40LjAuMTMwLy8vMC8:12.4.0.120/12.4.0.130/default"
		expectedObj := NewRange(comment2, name2, network, "", nil, false, nil, false, "", "", nil, "NONE")
		expectedObj.Ref = oldRef
		conn := &fakeConnector{
			getObjectObj:         NewEmptyRange(),
			getObjectRef:         oldRef,
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         expectedObj,
			getObjectError:       fmt.Errorf("not found"),
			fakeRefReturn:        oldRef,
			updateObjectObj:      expectedObj,
			updateObjectRef:      oldRef,
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)
		// negative scenario

		It("should fail to update Network Range Object", func() {
			actualRecord, err := objMgr.UpdateNetworkRange(oldRef, comment2, name2, network, "", "", false, nil, nil, "", nil, false, "NONE")
			Expect(actualRecord).To(BeNil())
			Expect(err).ToNot(BeNil())
		})

	})
	Describe("Negative case: return an error message on create when start_addr or end_addr is not provided", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "default"
		network := "12.4.0.0/24"
		comment := "range for 12.4.0.0/24"
		name := "range"
		conn := &fakeConnector{
			createObjectObj:   NewRange(comment, name, network, "", nil, false, nil, false, "", "", nil, "NONE"),
			createObjectError: fmt.Errorf("start address and end address fields are required to create a range within a Network"),
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)
		var actualRange, expectedObj *Range
		var err error
		expectedObj = nil
		It("should pass expected Network Range Object to CreateObject", func() {
			actualRange, err = objMgr.CreateNetworkRange(comment, name, network, netviewName, "", "", false, nil, nil, "", nil, false, "NONE")
			Expect(actualRange).To(Equal(expectedObj))
			Expect(err).To(Equal(conn.createObjectError))
		})
	})

})
