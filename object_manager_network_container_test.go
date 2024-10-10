package ibclient

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Object Manager: network container", func() {
	Describe("Create Network Container", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "Default View"
		cidr := "43.0.11.0/24"
		fakeRefReturn := "networkcontainer/ZG5zLm5ldHdvcmtfdmlldyQyMw:global_view/false"
		resObj := NewNetworkContainer(netviewName, cidr, false, "", nil)
		resObj.Ref = fakeRefReturn
		ncFakeConnector := &fakeConnector{
			createObjectObj: NewNetworkContainer(netviewName, cidr, false, "", nil),
			resultObject:    resObj,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(ncFakeConnector, cmpType, tenantID)

		var actualNetworkContainer *NetworkContainer
		var err error
		It("should pass expected NetworkContainer Object to CreateObject", func() {
			actualNetworkContainer, err = objMgr.CreateNetworkContainer(
				netviewName, cidr, false, "", nil)
		})
		It("should return expected NetworkContainer Object", func() {
			Expect(err).To(BeNil())
			Expect(actualNetworkContainer).To(Equal(ncFakeConnector.resultObject))
		})
	})

	Describe("Create IPv6 Network Container", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "default"
		cidr := "fc00::0100/56"
		cidrRef := "fc00%3A%3A0100/56"
		fakeRefReturn := fmt.Sprintf(
			"ipv6networkcontainer/ZZl7Lm5ldHdvcmtfdmlldyQyMw:%s/%s",
			cidrRef, netviewName)

		resObj := &NetworkContainer{
			NetviewName: netviewName,
			Cidr:        cidr,
		}
		resObj.objectType = "ipv6networkcontainer"
		resObj.returnFields = []string{"extattrs", "network", "network_view", "comment"}
		resObj.Ref = fakeRefReturn

		ncFakeConnector := &fakeConnector{
			createObjectObj: NewNetworkContainer(netviewName, cidr, true, "", nil),
			resultObject:    resObj,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(ncFakeConnector, cmpType, tenantID)

		var actualNetworkContainer *NetworkContainer
		var err error
		It("should pass expected NetworkContainer Object to CreateObject", func() {
			ncFakeConnector.createObjectError = nil
			actualNetworkContainer, err = objMgr.CreateNetworkContainer(netviewName, cidr, true, "", nil)
		})
		It("should return expected NetworkContainer Object", func() {
			Expect(err).To(BeNil())
			Expect(actualNetworkContainer).To(Equal(ncFakeConnector.resultObject))
		})

		// Negative test case: error may be returned by some reason.
		It("should pass expected NetworkContainer Object to CreateObject", func() {
			ncFakeConnector.createObjectError = NewNotFoundError("test error")
			actualNetworkContainer, err = objMgr.CreateNetworkContainer(netviewName, cidr, true, "", nil)
		})
		It("should return an error", func() {
			Expect(err).ToNot(BeNil())
			_, ok := err.(*NotFoundError)
			Expect(ok).To(BeTrue())
		})
	})

	Describe("Get Network Container by netview/CIDR", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "Default View"
		cidr := "43.0.11.0/24"
		fakeRefReturn := "networkcontainer/ZG5zLm5ldHdvcmtfdmlldyQyMw:global_view/false"
		resObj := NewNetworkContainer(netviewName, cidr, false, "", nil)
		resObj.Ref = fakeRefReturn

		queryParams := NewQueryParams(
			false,
			map[string]string{
				"network_view": netviewName,
				"network":      cidr,
			})

		ncFakeConnector := &fakeConnector{
			getObjectObj:         NewNetworkContainer(netviewName, cidr, false, "", nil),
			getObjectRef:         "",
			getObjectQueryParams: queryParams,
			resultObject:         []NetworkContainer{*resObj},
		}

		objMgr := NewObjectManager(ncFakeConnector, cmpType, tenantID)

		var actualNetworkContainer *NetworkContainer
		var err error
		It("should pass expected NetworkContainer Object to GetObject", func() {
			actualNetworkContainer, err = objMgr.GetNetworkContainer(netviewName, cidr, false, nil)
		})
		It("should return expected NetworkContainer Object", func() {
			Expect(err).To(BeNil())
			Expect(*actualNetworkContainer).To(Equal(ncFakeConnector.resultObject.([]NetworkContainer)[0]))
		})
	})

	Describe("Get Network Container by reference", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "Default View"
		cidr := "43.0.11.0/24"
		fakeRefReturn := "networkcontainer/ZG5zLm5ldHdvcmtfdmlldyQyMw:global_view/false"
		resObj := NewNetworkContainer(netviewName, cidr, false, "", nil)
		resObj.Ref = fakeRefReturn
		ncFakeConnector := &fakeConnector{
			getObjectObj:         NewNetworkContainer("", "", false, "", nil),
			getObjectRef:         fakeRefReturn,
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         resObj,
		}

		objMgr := NewObjectManager(ncFakeConnector, cmpType, tenantID)

		var actualNetworkContainer *NetworkContainer
		var err error
		It("should pass expected NetworkContainer Object to GetObject", func() {
			actualNetworkContainer, err = objMgr.GetNetworkContainerByRef(fakeRefReturn)
		})
		It("should return expected NetworkContainer Object", func() {
			Expect(err).To(BeNil())
			Expect(*actualNetworkContainer).To(Equal(*resObj))
		})
	})

	Describe("Get IPv6 Network Container", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "default"
		cidr := "fc00::0100/56"
		cidrRef := "fc00%3A%3A0100/56"
		fakeRefReturn := fmt.Sprintf(
			"ipv6networkcontainer/ZZl7Lm5ldHdvcmtfdmlldyQyMw:%s/%s",
			cidrRef, netviewName)

		resObj := NetworkContainer{
			NetviewName: netviewName,
			Cidr:        cidr,
		}
		resObj.objectType = "ipv6networkcontainer"
		resObj.returnFields = []string{"extattrs", "network", "network_view"}
		resObj.Ref = fakeRefReturn

		queryParams := NewQueryParams(
			false,
			map[string]string{
				"network_view": netviewName,
				"network":      cidr,
			})

		ncFakeConnector := &fakeConnector{
			getObjectObj: NewNetworkContainer(
				netviewName, cidr, true, "", nil),
			getObjectQueryParams: queryParams,
			fakeRefReturn:        fakeRefReturn,
		}

		objMgr := NewObjectManager(ncFakeConnector, cmpType, tenantID)

		var actualNetworkContainer *NetworkContainer
		var err error
		It("should pass expected NetworkContainer Object to GetObject", func() {
			resObj.Ea = make(EA)
			ncFakeConnector.resultObject = []NetworkContainer{resObj}
			ncFakeConnector.getObjectError = nil
			actualNetworkContainer, err = objMgr.GetNetworkContainer(netviewName, cidr, true, nil)
		})
		It("should return expected NetworkContainer Object", func() {
			Expect(err).To(BeNil())
			Expect(actualNetworkContainer).To(Equal(&resObj))
		})

		// Negative test case: error may be returned by some reason.
		It("should pass expected NetworkContainer Object to GetObject", func() {
			ncFakeConnector.getObjectError = fmt.Errorf("test error")
			actualNetworkContainer, err = objMgr.GetNetworkContainer(netviewName, cidr, true, nil)
		})
		It("should return an error", func() {
			_, ok := err.(*NotFoundError)
			Expect(ok).To(BeFalse())
		})

		// Negative test case: empty result set.
		It("should pass expected NetworkContainer Object to GetObject", func() {
			ncFakeConnector.getObjectError = nil
			ncFakeConnector.resultObject = []NetworkContainer{}
			actualNetworkContainer, err = objMgr.GetNetworkContainer(netviewName, cidr, true, nil)
		})
		It("should return an error", func() {
			_, ok := err.(*NotFoundError)
			Expect(ok).To(BeTrue())
		})
	})

	Describe("Update network container", func() {
		var (
			err       error
			objMgr    IBObjectManager
			conn      *fakeConnector
			ref       string
			actualObj *NetworkContainer
		)

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "default"
		refBase := "ZG5zLm5ldHdvcmtfdmlldyQyMw"
		ipv4Cidr := "10.2.1.0/20"
		ipv6Cidr := "fc00::0100/56"
		ipv6CidrRef := "fc00%3A%3A0100/56"

		It("IPv4, updating comment and EAs", func() {
			ref = fmt.Sprintf("networkcontainer/%s:%s", refBase, netviewName)
			initialEas := EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_old_value",
				"ea3": "ea3_value",
				"ea4": "ea4_value",
				"ea5": "ea5_old_value"}
			initObj := NewNetworkContainer(netviewName, ipv4Cidr, false, "old comment", initialEas)
			initObj.Ref = ref

			setEas := EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_new_value",
				"ea2": "ea2_new_value",
				"ea5": "ea5_old_value"}
			expectedEas := setEas

			getObjIn := &NetworkContainer{}
			getObjIn.returnFields = []string{"extattrs", "comment"}

			comment := "test comment 1"
			updateObjIn := NewNetworkContainer("", ipv4Cidr, false, comment, expectedEas)
			updateObjIn.Ref = ref

			expectedObj := NewNetworkContainer(netviewName, ipv4Cidr, false, comment, expectedEas)
			expectedObj.Ref = ref

			conn = &fakeConnector{
				getObjectObj:         getObjIn,
				getObjectQueryParams: NewQueryParams(false, nil),
				getObjectRef:         ref,
				getObjectError:       nil,
				resultObject:         initObj,

				updateObjectObj:   updateObjIn,
				updateObjectRef:   ref,
				updateObjectError: nil,

				fakeRefReturn: ref,
			}
			objMgr = NewObjectManager(conn, cmpType, tenantID)

			actualObj, err = objMgr.UpdateNetworkContainer(ref, setEas, comment)
			Expect(err).To(BeNil())
			Expect(actualObj).To(BeEquivalentTo(expectedObj))
		})

		It("Negative case: updating a network container which does not exist", func() {
			ref = fmt.Sprintf("networkcontainer/%s:%s", refBase, netviewName)
			initObj := NewNetworkContainer(netviewName, ipv4Cidr, false, "", nil)
			initObj.Ref = ref

			getObjIn := &NetworkContainer{}
			getObjIn.returnFields = []string{"extattrs", "comment"}

			comment := "test comment 1"

			conn = &fakeConnector{
				getObjectObj:         getObjIn,
				getObjectQueryParams: NewQueryParams(false, nil),
				getObjectRef:         ref,
				getObjectError:       fmt.Errorf("test error"),
				resultObject:         initObj,

				fakeRefReturn: ref,
			}
			objMgr = NewObjectManager(conn, cmpType, tenantID)

			actualObj, err = objMgr.UpdateNetworkContainer(ref, nil, comment)
			Expect(err).ToNot(BeNil())
		})

		It("Negative case: updating a network container with no update access", func() {
			ref = fmt.Sprintf("networkcontainer/%s:%s", refBase, netviewName)
			initObj := NewNetworkContainer(netviewName, ipv4Cidr, false, "old comment", nil)
			initObj.Ref = ref

			getObjIn := &NetworkContainer{}
			getObjIn.returnFields = []string{"extattrs", "comment"}

			comment := "test comment 1"
			updateObjIn := NewNetworkContainer("", ipv4Cidr, false, comment, nil)
			updateObjIn.Ref = ref

			conn = &fakeConnector{
				getObjectObj:         getObjIn,
				getObjectQueryParams: NewQueryParams(false, nil),
				getObjectRef:         ref,
				getObjectError:       nil,
				resultObject:         initObj,

				updateObjectObj:   updateObjIn,
				updateObjectRef:   ref,
				updateObjectError: fmt.Errorf("test error"),

				fakeRefReturn: ref,
			}
			objMgr = NewObjectManager(conn, cmpType, tenantID)

			actualObj, err = objMgr.UpdateNetworkContainer(ref, nil, comment)
			Expect(err).ToNot(BeNil())
		})

		It("IPv6, updating comment and EAs", func() {
			ref = fmt.Sprintf(
				"ipv6networkcontainer/%s:%s:%s",
				refBase, ipv6CidrRef, netviewName)
			initialEas := EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_old_value",
				"ea3": "ea3_value",
				"ea4": "ea4_value",
				"ea5": "ea5_old_value"}
			initObj := NewNetworkContainer(netviewName, ipv6Cidr, true, "", initialEas)
			initObj.Ref = ref

			setEas := EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_new_value",
				"ea2": "ea2_new_value",
				"ea5": "ea5_old_value"}
			expectedEas := setEas

			getObjIn := &NetworkContainer{}
			getObjIn.returnFields = []string{"extattrs", "comment"}

			comment := "test comment 1"
			updateObjIn := NewNetworkContainer("", ipv6Cidr, true, comment, expectedEas)
			updateObjIn.Ref = ref

			expectedObj := NewNetworkContainer(netviewName, ipv6Cidr, true, comment, expectedEas)
			expectedObj.Ref = ref

			conn = &fakeConnector{
				getObjectObj:         getObjIn,
				getObjectQueryParams: NewQueryParams(false, nil),
				getObjectRef:         ref,
				getObjectError:       nil,
				resultObject:         initObj,

				updateObjectObj:   updateObjIn,
				updateObjectRef:   ref,
				updateObjectError: nil,

				fakeRefReturn: ref,
			}
			objMgr = NewObjectManager(conn, cmpType, tenantID)

			actualObj, err = objMgr.UpdateNetworkContainer(ref, setEas, comment)
			Expect(err).To(BeNil())
			Expect(actualObj).To(BeEquivalentTo(expectedObj))
		})
	})

	Describe("Delete Network Container", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "default_view"
		cidr := "28.0.42.0/24"
		cidrRefIPv6 := "fc00%3A%3A0100/56"
		deleteRefIPv4 := fmt.Sprintf("networkcontainer/ZG5zLm5ldHdvcmskODkuMC4wLjAvMjQvMjU:%s/%s", cidr, netviewName)
		deleteRefIPv6 := fmt.Sprintf("networkcontainer/ZG5zLm5ldHdvcmskODkuMC4wLjAvMjQvMjU:%s/%s", cidrRefIPv6, netviewName)
		connector := &fakeConnector{}
		objMgr := NewObjectManager(connector, cmpType, tenantID)

		var actualRef string
		var err error

		It("should pass expected Network Ref to DeleteObject", func() {
			connector.deleteObjectRef = deleteRefIPv4
			connector.fakeRefReturn = deleteRefIPv4
			actualRef, err = objMgr.DeleteNetworkContainer(deleteRefIPv4)
		})
		It("should return expected Network container reference", func() {
			Expect(err).To(BeNil())
			Expect(actualRef).To(Equal(deleteRefIPv4))
		})

		// IPv6 case.
		It("should pass expected Network Ref to DeleteObject", func() {
			connector.deleteObjectRef = deleteRefIPv6
			connector.fakeRefReturn = deleteRefIPv6
			actualRef, err = objMgr.DeleteNetworkContainer(deleteRefIPv6)
		})
		It("should return expected Network container reference", func() {
			Expect(err).To(BeNil())
			Expect(actualRef).To(Equal(deleteRefIPv6))
		})

		var delRef string
		// Negative test case.
		It("should pass expected Network Ref to DeleteObject", func() {
			delRef = "networkcontainer"
			connector.deleteObjectRef = delRef
			connector.fakeRefReturn = ""
			connector.deleteObjectError = nil
			actualRef, err = objMgr.DeleteNetworkContainer(delRef)
		})
		It("should return an error", func() {
			Expect(err).ToNot(BeNil())
		})
		// Negative test case.
		It("should pass expected Network Ref to DeleteObject", func() {
			delRef = fmt.Sprintf(
				"network/ZG5zLm5ldHdvcmskODkuMC4wLjAvMjQvMjU:%s/%s",
				cidr, netviewName)
			connector.deleteObjectRef = delRef
			connector.fakeRefReturn = ""
			connector.deleteObjectError = nil
			actualRef, err = objMgr.DeleteNetworkContainer(delRef)
		})
		It("should return an error", func() {
			Expect(err).ToNot(BeNil())
		})
	})

	Describe("Allocate Network Container", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "default_view"
		cidr := "142.0.22.0/24"
		prefixLen := uint(28)
		networkName := "private-net"
		fakeRefReturn := fmt.Sprintf("networkcontainer/ZG5zLm5ldHdvcmskODkuMC4wLjAvMjQvMjU:%s/%s", cidr, netviewName)
		ea := EA{"Site": "test"}
		comment := "Test network container"
		resObj, err := BuildNetworkContainerFromRef(fakeRefReturn)

		containerInfo := NewNetworkContainerNextAvailableInfo(netviewName, cidr, prefixLen, false)
		container := NewNetworkContainerNextAvailable(containerInfo, false, comment, ea)

		connector := &fakeConnector{
			createObjectObj: container,
			resultObject:    resObj,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(connector, cmpType, tenantID)

		connector.createObjectObj.(*NetworkContainerNextAvailable).Ea = ea
		connector.createObjectObj.(*NetworkContainerNextAvailable).Ea["Network Name"] = networkName

		var actualNetwork *NetworkContainer
		It("should pass expected Network Container Object to CreateObject", func() {
			actualNetwork, err = objMgr.AllocateNetworkContainer(
				netviewName, cidr, false, prefixLen, comment, ea)
		})
		It("should return expected Network Object", func() {
			Expect(actualNetwork).To(Equal(connector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Does not allocate Network Container if an invalid cidr is passed", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "default_view"
		cidr := "10.0.1.0./64"
		prefixLen := uint(65)
		networkName := "private-net"
		fakeRefReturn := fmt.Sprintf("networkcontainer/ZG5zLm5ldHdvcmskODkuMC4wLjAvMjQvMjU:%s/%s", cidr, netviewName)
		ea := EA{"Site": "test"}
		comment := "Test network container"
		resObj, err := BuildNetworkContainerFromRef(fakeRefReturn)

		containerInfo := NewNetworkContainerNextAvailableInfo(netviewName, cidr, prefixLen, false)
		container := NewNetworkContainerNextAvailable(containerInfo, false, comment, ea)

		connector := &fakeConnector{
			createObjectObj: container,
			resultObject:    resObj,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(connector, cmpType, tenantID)

		connector.createObjectObj.(*NetworkContainerNextAvailable).Ea = ea
		connector.createObjectObj.(*NetworkContainerNextAvailable).Ea["Network Name"] = networkName

		var actualNetwork *NetworkContainer
		It("should pass expected Network Container Object with invalid Cidr value to CreateObject", func() {
			actualNetwork, err = objMgr.AllocateNetworkContainer(
				netviewName, cidr, false, prefixLen, comment, ea)
		})
		It("should return nil and an error message", func() {
			Expect(actualNetwork).To(Equal(connector.resultObject))
			Expect(err).To(Equal(fmt.Errorf("CIDR format not matched")))
		})
	})

	Describe("Allocate Network Container", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "default_view"
		cidr := "2003:db8:abcd:14::/64"
		prefixLen := uint(28)
		networkName := "private-net"
		fakeRefReturn := fmt.Sprintf("ipv6networkcontainer/ZG5zLm5ldHdvcmskODkuMC4wLjAvMjQvMjU:%s/%s", cidr, netviewName)
		ea := EA{"Site": "test"}
		comment := "Test network container"
		resObj, err := BuildIPv6NetworkContainerFromRef(fakeRefReturn)
		containerInfo := NewNetworkContainerNextAvailableInfo(netviewName, cidr, prefixLen, true)
		container := NewNetworkContainerNextAvailable(containerInfo, true, comment, ea)

		connector := &fakeConnector{
			createObjectObj: container,
			resultObject:    resObj,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(connector, cmpType, tenantID)

		connector.createObjectObj.(*NetworkContainerNextAvailable).Ea = ea
		connector.createObjectObj.(*NetworkContainerNextAvailable).Ea["Network Name"] = networkName

		var actualNetwork *NetworkContainer
		It("should pass expected Network Container Object to CreateObject", func() {
			actualNetwork, err = objMgr.AllocateNetworkContainer(
				netviewName, cidr, true, prefixLen, comment, ea)
		})
		It("should return expected Network Object", func() {
			Expect(actualNetwork).To(Equal(connector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Allocate IPV4 Network Container by EA", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		cidr := "20.20.1.0/24"
		prefixLen := uint(24)
		netviewName := "default"
		fakeRefReturn := fmt.Sprintf("networkcontainer/ZG5zLm5ldHdvcmskODkuMC4wLjAvMjQvMjU:%s/%s", cidr, netviewName)
		ea := EA{"Site": "test"}
		eaMap := map[string]string{"Site": "Turkey"}
		comment := "Test network container"
		resObj, err := BuildNetworkContainerFromRef(fakeRefReturn)
		container := &NetworkContainerNextAvailable{
			objectType: "networkcontainer",
			Network: &NetworkContainerNextAvailableInfo{
				Function:     "next_available_network",
				ResultField:  "networks",
				Object:       "networkcontainer",
				ObjectParams: eaMap,
				Params:       map[string]uint{"cidr": prefixLen},
				NetviewName:  "",
			},
			NetviewName: netviewName,
			Comment:     comment,
			Ea:          ea,
		}

		connector := &fakeConnector{
			createObjectObj: container,
			resultObject:    resObj,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(connector, cmpType, tenantID)

		var actualNetwork *NetworkContainer
		It("should pass expected Network Container Object to CreateObject", func() {
			actualNetwork, err = objMgr.AllocateNetworkContainerByEA(netviewName, false, comment, ea, eaMap, prefixLen)
		})
		It("should return expected Network container Object", func() {
			Expect(actualNetwork).To(Equal(connector.resultObject))
			Expect(err).To(BeNil())
		})
	})
})
