package ibclient

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Object Manager: network", func() {
	Describe("Create Network", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "default_view"
		cidr := "43.0.11.0/24"
		networkName := "private-net"
		fakeRefReturn := "network/ZG5zLm5ldHdvcmskODkuMC4wLjAvMjQvMjU:43.0.11.0/24/default_view"
		ea := EA{"Lock": "added", "Region": "East"}
		comment := "Test network view"
		connector := &fakeConnector{
			createObjectObj: NewNetwork(netviewName, cidr, false, comment, ea),
			resultObject:    NewNetwork(netviewName, cidr, false, comment, ea),
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(connector, cmpType, tenantID)

		connector.createObjectObj.(*Network).Ea = ea
		connector.createObjectObj.(*Network).Ea["Network Name"] = networkName

		connector.resultObject.(*Network).Ref = fakeRefReturn
		connector.resultObject.(*Network).Ea = ea
		connector.resultObject.(*Network).Ea["Network Name"] = networkName

		var actualNetwork *Network
		var err error
		It("should pass expected Network Object to CreateObject", func() {
			actualNetwork, err = objMgr.CreateNetwork(
				netviewName, cidr, false, comment, ea)
		})
		It("should return expected Network Object", func() {
			Expect(actualNetwork).To(Equal(connector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Create IPv6 Network", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "default_view"
		cidr := "2001:db8:abcd:14::/64"
		cidrRef := " 2001%3Adb8%3Aabcd%3A14%3A%3A/64"
		networkName := "private-net"
		fakeRefReturn := fmt.Sprintf("ipv6network/ZG5zLm5ldHdvcmskODkuMC4wLjAvMjQvMjU:%s/default_view", cidrRef)
		ea := EA{"Lock": "added", "Region": "East"}
		comment := "Test network view"
		connector := &fakeConnector{
			createObjectObj: NewNetwork(netviewName, cidr, true, comment, ea),
			resultObject:    NewNetwork(netviewName, cidr, true, comment, ea),
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(connector, cmpType, tenantID)

		connector.createObjectObj.(*Network).Ea = ea
		connector.createObjectObj.(*Network).Ea["Network Name"] = networkName

		connector.resultObject.(*Network).Ref = fakeRefReturn
		connector.resultObject.(*Network).Ea = ea
		connector.resultObject.(*Network).Ea["Network Name"] = networkName

		var actualNetwork *Network
		var err error
		It("should pass expected Network Object to CreateObject", func() {
			actualNetwork, err = objMgr.CreateNetwork(
				netviewName, cidr, true, comment, ea)
		})
		It("should return expected Network Object", func() {
			Expect(actualNetwork).To(Equal(connector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Allocate Network", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "default_view"
		cidr := "142.0.22.0/24"
		prefixLen := uint(28)
		networkName := "private-net"
		cidr1 := fmt.Sprintf("func:nextavailablenetwork:%s,%s,%d", cidr, netviewName, prefixLen)
		fakeRefReturn := fmt.Sprintf("network/ZG5zLm5ldHdvcmskODkuMC4wLjAvMjQvMjU:%s/%s", cidr, netviewName)
		ea := EA{"Lock": "added", "Region": "East"}
		comment := "Test network view"
		resObj, err := BuildNetworkFromRef(fakeRefReturn)
		connector := &fakeConnector{
			createObjectObj: NewNetwork(netviewName, cidr1, false, comment, ea),
			resultObject:    resObj,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(connector, cmpType, tenantID)

		connector.createObjectObj.(*Network).Ea = ea
		connector.createObjectObj.(*Network).Ea["Network Name"] = networkName

		var actualNetwork *Network
		It("should pass expected Network Object to CreateObject", func() {
			actualNetwork, err = objMgr.AllocateNetwork(
				netviewName, cidr, false, prefixLen, comment, ea)
		})
		It("should return expected Network Object", func() {
			Expect(actualNetwork).To(Equal(connector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Does not allocate Network if an invalid cidr is passed", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "default_view"
		cidr := "10.0.1.0./64"
		prefixLen := uint(65)
		networkName := "private-net"
		cidr1 := fmt.Sprintf("func:nextavailablenetwork:%s,%s,%d", cidr, netviewName, prefixLen)
		fakeRefReturn := fmt.Sprintf("network/ZG5zLm5ldHdvcmskODkuMC4wLjAvMjQvMjU:%s/%s", cidr, netviewName)
		ea := EA{"Lock": "added", "Region": "East"}
		comment := "Test network view"
		resObj, err := BuildNetworkFromRef(fakeRefReturn)
		connector := &fakeConnector{
			createObjectObj: NewNetwork(netviewName, cidr1, false, comment, ea),
			resultObject:    resObj,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(connector, cmpType, tenantID)

		connector.createObjectObj.(*Network).Ea = ea
		connector.createObjectObj.(*Network).Ea["Network Name"] = networkName

		var actualNetwork *Network
		It("should pass expected Network Object with invalid Cidr value to CreateObject", func() {
			actualNetwork, err = objMgr.AllocateNetwork(
				netviewName, cidr, false, prefixLen, comment, ea)
		})
		It("should return nil and an error message", func() {
			Expect(actualNetwork).To(Equal(connector.resultObject))
			Expect(err).To(Equal(fmt.Errorf("CIDR format not matched")))
		})
	})

	Describe("Allocate IPv6 Network", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "default_view"
		cidr := "2003:db8:abcd:14::/64"
		prefixLen := uint(28)
		networkName := "private-net"
		cidr1 := fmt.Sprintf("func:nextavailablenetwork:%s,%s,%d", cidr, netviewName, prefixLen)
		fakeRefReturn := fmt.Sprintf("ipv6network/ZG5zLm5ldHdvcmskODkuMC4wLjAvMjQvMjU:%s/%s", cidr, netviewName)
		ea := EA{"Lock": "added", "Region": "East", "Network Name": networkName}
		comment := "Test network view"

		resObj, err := BuildIPv6NetworkFromRef(fakeRefReturn)
		connector := &fakeConnector{
			createObjectObj: NewNetwork(netviewName, cidr1, true, comment, ea),
			resultObject:    resObj,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(connector, cmpType, tenantID)

		var actualNetwork *Network
		It("should pass expected Network Object to CreateObject", func() {
			actualNetwork, err = objMgr.AllocateNetwork(
				netviewName, cidr, true, prefixLen, comment, ea)
		})
		It("should return expected Network Object", func() {
			Expect(actualNetwork).To(Equal(connector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get Network", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "default_view"
		cidr := "28.0.42.0/24"
		networkName := "private-net"
		ea := EA{"Network Name": networkName}
		fakeRefReturn := fmt.Sprintf("network/ZG5zLm5ldHdvcmskODkuMC4wLjAvMjQvMjU:%s/%s", cidr, netviewName)

		queryParams := NewQueryParams(
			false,
			map[string]string{
				"network_view": netviewName,
				"network":      cidr,
			})

		connector := &fakeConnector{
			getObjectObj:         NewNetwork(netviewName, cidr, false, "", ea),
			getObjectRef:         "",
			getObjectQueryParams: queryParams,
			resultObject:         []Network{*NewNetwork(netviewName, cidr, false, "", ea)},
		}

		connector.resultObject.([]Network)[0].Ref = fakeRefReturn
		connector.getObjectObj.(*Network).eaSearch = EASearch(ea)
		connector.resultObject.([]Network)[0].eaSearch = EASearch(ea)

		objMgr := NewObjectManager(connector, cmpType, tenantID)

		var actualNetwork *Network
		var err error
		It("should pass expected Network Object to GetObject", func() {
			actualNetwork, err = objMgr.GetNetwork(netviewName, cidr, false, ea)
		})
		It("should return expected Network Object", func() {
			Expect(*actualNetwork).To(Equal(connector.resultObject.([]Network)[0]))
			Expect(err).To(BeNil())
		})
	})

	Describe("Does not fetch the Network if required arguments are not passed", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := ""
		cidr := "10.0.0.0/24"
		networkName := "private-net"
		ea := EA{"Network Name": networkName}
		comment := "Test network view"
		connector := &fakeConnector{
			getObjectObj:         NewNetwork(netviewName, cidr, false, comment, ea),
			getObjectRef:         "",
			getObjectQueryParams: NewQueryParams(false, nil),
		}

		connector.getObjectObj.(*Network).eaSearch = EASearch(ea)
		objMgr := NewObjectManager(connector, cmpType, tenantID)

		var actualNetwork, resultObj *Network
		resultObj = nil
		var err error
		It("should pass expected Network Object to GetObject", func() {
			actualNetwork, err = objMgr.GetNetwork(netviewName, cidr, false, ea)
		})
		It("should return nil and an error message", func() {
			Expect(actualNetwork).To(Equal(resultObj))
			Expect(err).To(Equal(fmt.Errorf("both network view and cidr values are required")))
		})
	})

	Describe("Get IPv6 Network", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "default_view"
		cidr := "2001:db8:abcd:14::/64"
		cidrRef := " 2001%3Adb8%3Aabcd%3A14%3A%3A/64"
		networkName := "private-net"
		ea := EA{"Network Name": networkName}
		fakeRefReturn := fmt.Sprintf("ipv6network/ZG5zLm5ldHdvcmskODkuMC4wLjAvMjQvMjU:%s/%s", cidrRef, netviewName)

		queryParams := NewQueryParams(
			false,
			map[string]string{
				"network_view": netviewName,
				"network":      cidr,
			})

		connector := &fakeConnector{
			getObjectObj:         NewNetwork(netviewName, cidr, true, "", ea),
			getObjectRef:         "",
			getObjectQueryParams: queryParams,
			resultObject:         []Network{*NewNetwork(netviewName, cidr, true, "", ea)},
		}

		connector.resultObject.([]Network)[0].Ref = fakeRefReturn
		connector.getObjectObj.(*Network).eaSearch = EASearch(ea)
		connector.resultObject.([]Network)[0].eaSearch = EASearch(ea)

		objMgr := NewObjectManager(connector, cmpType, tenantID)

		var actualNetwork *Network
		var err error
		It("should pass expected Network Object to GetObject", func() {
			actualNetwork, err = objMgr.GetNetwork(netviewName, cidr, true, ea)
		})
		It("should return expected Network Object", func() {
			Expect(*actualNetwork).To(Equal(connector.resultObject.([]Network)[0]))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get Network by reference", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "Default View"
		cidr := "43.0.11.0/24"
		fakeRefReturn := "network/ZG5zLm5ldHdvcmtfdmlldyQyMw:global_view/false"
		resObj := NewNetwork(netviewName, cidr, false, "", nil)
		resObj.Ref = fakeRefReturn
		ncFakeConnector := &fakeConnector{
			getObjectObj:         NewNetwork("", "", false, "", nil),
			getObjectRef:         fakeRefReturn,
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         resObj,
		}

		objMgr := NewObjectManager(ncFakeConnector, cmpType, tenantID)

		var actualNetwork *Network
		var err error
		It("should pass expected Network Object to GetObject", func() {
			actualNetwork, err = objMgr.GetNetworkByRef(fakeRefReturn)
		})
		It("should return expected Network Object", func() {
			Expect(err).To(BeNil())
			Expect(*actualNetwork).To(Equal(*resObj))
		})
	})

	Describe("Update network", func() {
		var (
			err       error
			objMgr    IBObjectManager
			conn      *fakeConnector
			ref       string
			actualObj *Network
		)

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "default"
		refBase := "ZG5zLm5ldHdvcmtfdmlldyQyMw"
		cidr := "10.2.1.0/20"

		It("Updating comment and EAs", func() {
			ref = fmt.Sprintf("network/%s:%s", refBase, netviewName)
			initialEas := EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_old_value",
				"ea3": "ea3_value",
				"ea4": "ea4_value",
				"ea5": "ea5_old_value"}
			initObj := NewNetwork(netviewName, cidr, false, "old comment", initialEas)
			initObj.Ref = ref

			setEas := EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_new_value",
				"ea2": "ea2_new_value",
				"ea5": "ea5_old_value"}
			expectedEas := setEas

			getObjIn := &Network{}
			getObjIn.objectType = "network"
			getObjIn.returnFields = []string{"extattrs", "network", "network_view", "comment"}

			comment := "test comment 1"
			updateObjIn := NewNetwork("", cidr, false, comment, expectedEas)
			updateObjIn.Ref = ref

			expectedObj := NewNetwork(netviewName, cidr, false, comment, expectedEas)
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

			actualObj, err = objMgr.UpdateNetwork(ref, setEas, comment)
			Expect(err).To(BeNil())
			Expect(actualObj).To(BeEquivalentTo(expectedObj))
		})

		It("removing all EAs", func() {
			ref = fmt.Sprintf("network/%s:%s", refBase, netviewName)
			initialEas := EA{
				"ea3": "ea3_value",
				"ea4": "ea4_value"}
			initObj := NewNetwork(netviewName, cidr, false, "old comment", initialEas)
			initObj.Ref = ref

			setEas := EA{}
			expectedEas := setEas

			getObjIn := &Network{}
			getObjIn.objectType = "network"
			getObjIn.returnFields = []string{"extattrs", "network", "network_view", "comment"}

			comment := "test comment 1"
			updateObjIn := NewNetwork("", cidr, false, comment, expectedEas)
			updateObjIn.Ref = ref

			expectedObj := NewNetwork(netviewName, cidr, false, comment, expectedEas)
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

			actualObj, err = objMgr.UpdateNetwork(ref, setEas, comment)
			Expect(err).To(BeNil())
			Expect(actualObj).To(BeEquivalentTo(expectedObj))
		})

		It("Negative case: updating a IPv4 network which does not exist", func() {
			ref = fmt.Sprintf("network/%s:%s", refBase, netviewName)
			initObj := NewNetwork(netviewName, cidr, false, "", nil)
			initObj.Ref = ref

			getObjIn := &Network{}
			getObjIn.objectType = "network"
			getObjIn.returnFields = []string{"extattrs", "network", "network_view", "comment"}

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

			_, err = objMgr.UpdateNetwork(ref, nil, comment)
			Expect(err).ToNot(BeNil())
		})

		It("Negative case: updating an IPv4 network with no update access", func() {
			ref = fmt.Sprintf("network/%s:%s", refBase, netviewName)
			initObj := NewNetwork(netviewName, cidr, false, "old comment", nil)
			initObj.Ref = ref

			getObjIn := &Network{}
			getObjIn.objectType = "network"
			getObjIn.returnFields = []string{"extattrs", "network", "network_view", "comment"}

			comment := "test comment 1"
			updateObjIn := NewNetwork("", cidr, false, comment, nil)
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

			actualObj, err = objMgr.UpdateNetwork(ref, nil, comment)
			Expect(err).ToNot(BeNil())
		})

		It("Clearing the comment field", func() {
			ref = fmt.Sprintf("network/%s:%s", refBase, netviewName)
			initObj := NewNetwork(netviewName, cidr, false, "old comment", nil)
			initObj.Ref = ref

			getObjIn := &Network{}
			getObjIn.objectType = "network"
			getObjIn.returnFields = []string{"extattrs", "network", "network_view", "comment"}

			comment := ""
			updateObjIn := NewNetwork("", cidr, false, comment, nil)
			updateObjIn.Ref = ref

			expectedObj := NewNetwork(netviewName, cidr, false, comment, nil)
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

			actualObj, err = objMgr.UpdateNetwork(ref, nil, comment)
			Expect(err).To(BeNil())
			Expect(actualObj).To(BeEquivalentTo(expectedObj))
		})

		It("Updating EAs only. Comment field unchanged.", func() {
			ref = fmt.Sprintf("network/%s:%s", refBase, netviewName)
			initialEas := EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_old_value",
				"ea3": "ea3_value",
				"ea4": "ea4_value",
				"ea5": "ea5_old_value"}
			initObj := NewNetwork(netviewName, cidr, false, "old comment", initialEas)
			initObj.Ref = ref

			setEas := EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_new_value",
				"ea2": "ea2_new_value",
				"ea5": "ea5_old_value"}
			expectedEas := setEas

			getObjIn := &Network{}
			getObjIn.objectType = "network"
			getObjIn.returnFields = []string{"extattrs", "network", "network_view", "comment"}

			comment := "test comment 1"
			updateObjIn := NewNetwork("", cidr, false, comment, expectedEas)
			updateObjIn.Ref = ref

			expectedObj := NewNetwork(netviewName, cidr, false, comment, expectedEas)
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

			actualObj, err = objMgr.UpdateNetwork(ref, setEas, comment)
			Expect(err).To(BeNil())
			Expect(actualObj).To(BeEquivalentTo(expectedObj))
		})

		It("IPv6, Updating comment and EAs", func() {
			ref = fmt.Sprintf("ipv6network/%s:%s", refBase, netviewName)
			initialEas := EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_old_value",
				"ea3": "ea3_value",
				"ea4": "ea4_value",
				"ea5": "ea5_old_value"}
			initObj := NewNetwork(netviewName, cidr, true, "old comment", initialEas)
			initObj.Ref = ref

			setEas := EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_new_value",
				"ea2": "ea2_new_value",
				"ea5": "ea5_old_value"}
			expectedEas := setEas

			getObjIn := &Network{}
			getObjIn.objectType = "ipv6network"
			getObjIn.returnFields = []string{"extattrs", "network", "network_view", "comment"}

			comment := "test comment 1"
			updateObjIn := NewNetwork("", cidr, true, comment, expectedEas)
			updateObjIn.Ref = ref

			expectedObj := NewNetwork(netviewName, cidr, true, comment, expectedEas)
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

			actualObj, err = objMgr.UpdateNetwork(ref, setEas, comment)
			Expect(err).To(BeNil())
			Expect(actualObj).To(BeEquivalentTo(expectedObj))
		})
	})

	Describe("Delete Network", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "default_view"
		cidr := "28.0.42.0/24"
		deleteRef := fmt.Sprintf("network/ZG5zLm5ldHdvcmskODkuMC4wLjAvMjQvMjU:%s/%s", cidr, netviewName)
		fakeRefReturn := deleteRef
		connector := &fakeConnector{
			deleteObjectRef: deleteRef,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(connector, cmpType, tenantID)

		var actualRef string
		var err error
		It("should pass expected Network Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteNetwork(deleteRef)
		})
		It("should return expected Network Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete IPv6 Fixed Address", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "default_view"
		cidrRef := "2003%3Adb8%3Aabcd%3A14%3A1"
		deleteRef := fmt.Sprintf("ipv6fixedaddress/ZG5zLm5ldHdvcmskODkuMC4wLjAvMjQvMjU:%s/%s", cidrRef, netviewName)
		fakeRefReturn := deleteRef
		connector := &fakeConnector{
			deleteObjectRef: deleteRef,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(connector, cmpType, tenantID)

		var actualRef string
		var err error
		It("should pass expected IPv6 fixed address Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteFixedAddress(deleteRef)
		})
		It("should return expected Network Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})

	Describe("BuildNetworkFromRef", func() {
		netviewName := "test_view"
		cidr := "23.11.0.0/24"
		networkRef := fmt.Sprintf("network/ZG5zLm5ldHdvcmtfdmlldyQyMw:%s/%s", cidr, netviewName)

		expectedNetwork := Network{Ref: networkRef, NetviewName: netviewName, Cidr: cidr}
		expectedNetwork.objectType = "network"
		expectedNetwork.returnFields = []string{"extattrs", "network", "network_view", "comment"}
		resObj, err := BuildNetworkFromRef(networkRef)
		resObj1, err1 := BuildNetworkFromRef("network/ZG5zLm5ldHdvcmtfdmlldyQyMw")
		It("should return expected Network Object", func() {
			Expect(*resObj).To(Equal(expectedNetwork))
			Expect(err).To(BeNil())
		})
		It("should fail if bad Network Ref is provided", func() {
			Expect(resObj1).To(BeNil())
			Expect(err1).To(Equal(fmt.Errorf("CIDR format not matched")))
		})
	})

	Describe("BuildIPv6NetworkFromRef", func() {
		netviewName := "test_view"
		cidr := "2001:db8:abcd:14::/64"
		cidrRef := "2001%3Adb8%3Aabcd%3A14%3A%3A/64"
		networkRef := fmt.Sprintf("ipv6network/ZG5zLm5ldHdvcmtfdmlldyQyMw:%s/%s", cidrRef, netviewName)

		expectedNetwork := Network{Ref: networkRef, NetviewName: netviewName, Cidr: cidr}
		expectedNetwork.objectType = "ipv6network"
		expectedNetwork.returnFields = []string{"extattrs", "network", "network_view", "comment"}
		resObj, err := BuildIPv6NetworkFromRef(networkRef)
		resObj1, err1 := BuildIPv6NetworkFromRef("ipv6network/ZG5zLm5ldHdvcmtfdmlldyQyMw")
		It("should return expected Network Object", func() {
			Expect(*resObj).To(Equal(expectedNetwork))
			Expect(err).To(BeNil())
		})
		It("should fail if bad Network Ref is provided", func() {
			Expect(resObj1).To(BeNil())
			Expect(err1).To(Equal(fmt.Errorf("CIDR format not matched")))
		})
	})
})
