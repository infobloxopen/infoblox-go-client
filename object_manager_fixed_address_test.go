package ibclient

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Object Manager: fixed address", func() {
	Describe("Allocate Specific IP", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "private"
		cidr := "53.0.0.0/24"
		ipAddr := "53.0.0.21"
		macAddr := "01:23:45:67:80:ab"
		comment := "test"
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		name := "testvm"
		fakeRefReturn := fmt.Sprintf("fixedaddress/ZG5zLmJpbmRfY25h:%s/private", ipAddr)
		ea := EA{"VM ID": vmID, "VM Name": vmName}
		isIPv6 := false

		conn := &fakeConnector{
			createObjectObj: NewFixedAddress(
				netviewName, name,
				ipAddr, cidr, macAddr,
				nil, ea, "", isIPv6, comment, nil, nil, nil, nil, false, nil, false),
			getObjectObj:         NewEmptyFixedAddress(isIPv6),
			getObjectRef:         fakeRefReturn,
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject: NewFixedAddress(
				netviewName, name,
				GetIPAddressFromRef(fakeRefReturn), cidr, macAddr,
				nil, ea, fakeRefReturn, isIPv6, comment, nil, nil, nil, nil, false, nil, false),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)

		var actualIP *FixedAddress
		var err error
		It("should pass expected Fixed Address Object to CreateObject", func() {
			actualIP, err = objMgr.AllocateIP(netviewName, cidr, ipAddr, isIPv6, macAddr, name, comment, ea, "", "", "", nil, "", false, nil, false)
		})
		It("should return expected Fixed Address Object", func() {
			Expect(actualIP).To(Equal(conn.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Allocate Next Available IP", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "private"
		cidr := "53.0.0.0/24"
		ipAddr := fmt.Sprintf("func:nextavailableip:%s,%s", cidr, netviewName)
		macAddr := "01:23:45:67:80:ab"
		comment := "test"
		isIPv6 := false
		vmID := "93f9249abc039284"
		name := "testvm"
		vmName := "dummyvm"
		resultIP := "53.0.0.32"
		fakeRefReturn := fmt.Sprintf("fixedaddress/ZG5zLmJpbmRfY25h:%s/private", resultIP)
		ea := EA{"VM ID": vmID, "VM Name": vmName}

		conn := &fakeConnector{
			createObjectObj: NewFixedAddress(
				netviewName, name,
				ipAddr, cidr, macAddr,
				nil, ea, "", isIPv6, comment, nil, nil, nil, nil, false, nil, false),
			getObjectObj:         NewEmptyFixedAddress(isIPv6),
			getObjectRef:         fakeRefReturn,
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject: NewFixedAddress(
				netviewName, name,
				resultIP, cidr, macAddr,
				nil, ea, fakeRefReturn, isIPv6, comment, nil, nil, nil, nil, false, nil, false),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)

		var actualIP *FixedAddress
		var err error
		It("should pass expected Fixed Address Object to CreateObject", func() {
			actualIP, err = objMgr.AllocateIP(netviewName, cidr, "", isIPv6, macAddr, name, comment, ea, "", "", "", nil, "", false, nil, false)
		})

		It("should return expected Fixed Address Object", func() {
			Expect(actualIP).To(Equal(conn.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Allocate Specific IPv6 Address", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "private"
		cidr := "2001:db8:abcd:12::/64"
		ipAddr := "2001:db8:abcd:12::1"
		refIp := "2001%3Adb8%3Aabcd%3A12%3A%3A1"
		duid := "01:23:45:67:80:ab"
		comment := "test"
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		name := "testvm"
		ea := EA{"VM ID": vmID, "VM Name": vmName}
		fakeRefReturn := fmt.Sprintf("ipv6fixedaddress/ZG5zLmJpbmRfY25h:%s/private", refIp)
		isIPv6 := true

		conn := &fakeConnector{
			createObjectObj: NewFixedAddress(
				netviewName, name,
				ipAddr, cidr, duid,
				nil, ea, "", isIPv6, comment, nil, nil, nil, nil, false, nil, false),
			getObjectObj:         NewEmptyFixedAddress(isIPv6),
			getObjectRef:         fakeRefReturn,
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject: NewFixedAddress(
				netviewName, name,
				ipAddr, cidr, duid,
				nil, ea, fakeRefReturn, isIPv6, comment, nil, nil, nil, nil, false, nil, false),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)

		var actualIP *FixedAddress
		var err error
		It("should pass expected Fixed Address Object to CreateObject", func() {
			actualIP, err = objMgr.AllocateIP(netviewName, cidr, ipAddr, isIPv6, duid, name, comment, ea, "", "", "", nil, "", false, nil, false)
		})
		It("should return expected Fixed Address Object", func() {
			Expect(actualIP).To(Equal(conn.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Allocate Next Available IPv6 Address", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "private"
		cidr := "2001:db8:abcd:12::/64"
		ipAddr := fmt.Sprintf("func:nextavailableip:%s,%s", cidr, netviewName)
		duid := "01:23:45:67:80:ab"
		comment := "test"
		isIPv6 := true
		vmID := "93f9249abc039284"
		name := "testvm"
		vmName := "dummyvm"
		resultIP := "2001%3Adb8%3Aabcd%3A12%3A%3A1"
		fakeRefReturn := fmt.Sprintf("ipv6fixedaddress/ZG5zLmJpbmRfY25h:%s/private", resultIP)
		ea := EA{"VM ID": vmID, "VM Name": vmName}

		conn := &fakeConnector{
			createObjectObj: NewFixedAddress(
				netviewName, name,
				ipAddr, cidr, duid,
				nil, ea, "", isIPv6, comment, nil, nil, nil, nil, false, nil, false),
			getObjectObj:         NewEmptyFixedAddress(isIPv6),
			getObjectRef:         fakeRefReturn,
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject: NewFixedAddress(
				netviewName, name,
				resultIP, cidr, duid,
				nil, ea, fakeRefReturn, isIPv6, comment, nil, nil, nil, nil, false, nil, false),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)

		var actualIP *FixedAddress
		var err error
		It("should pass expected Fixed Address Object to CreateObject", func() {
			actualIP, err = objMgr.AllocateIP(netviewName, cidr, "", isIPv6, duid, name, comment, ea, "", "", "", nil, "", false, nil, false)
		})

		It("should return expected Fixed Address Object", func() {
			Expect(actualIP).To(Equal(conn.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Negative case:Does not allocate IPv6 Address when DUID is not passed", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "private"
		cidr := "2001:db8:abcd:12::/64"
		ipAddr := fmt.Sprintf("func:nextavailableip:%s,%s", cidr, netviewName)
		duid := ""
		comment := "test"
		isIPv6 := true
		vmID := "93f9249abc039284"
		name := "testvm"
		vmName := "dummyvm"
		resultIP := "2001%3Adb8%3Aabcd%3A12%3A%3A1"
		fakeRefReturn := fmt.Sprintf("ipv6fixedaddress/ZG5zLmJpbmRfY25h:%s/private", resultIP)
		ea := EA{"VM ID": vmID, "VM Name": vmName}
		var expectedObj *FixedAddress
		expectedObj = nil
		conn := &fakeConnector{
			createObjectObj: NewFixedAddress(
				netviewName, name,
				ipAddr, cidr, duid,
				nil, ea, "", isIPv6, comment, nil, nil, nil, nil, false, nil, false),
			createObjectError: fmt.Errorf("the DUID field cannot be left empty"),
			fakeRefReturn:     fakeRefReturn,
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)

		var actualIP *FixedAddress
		var err error
		It("should pass expected Fixed Address Object to CreateObject", func() {
			actualIP, err = objMgr.AllocateIP(netviewName, cidr, "", isIPv6, duid, name, comment, ea, "", "", "", nil, "", false, nil, false)
		})

		It("should return expected Fixed Address Object", func() {
			Expect(actualIP).To(Equal(expectedObj))
			Expect(err).To(Equal(conn.createObjectError))
		})
	})

	Describe("Get Fixed Address", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "private"
		cidr := "53.0.0.0/24"
		ipAddr := "53.0.0.21"
		macAddr := "01:23:45:67:80:ab"
		isIPv6 := false
		comment := "test"
		fakeRefReturn := fmt.Sprintf("fixedaddress/ZG5zLmJpbmRfY25h:%s/private", ipAddr)

		queryParams := NewQueryParams(
			false,
			map[string]string{
				"network_view": netviewName,
				"network":      cidr,
				"ipv4addr":     ipAddr,
				"mac":          macAddr,
			})

		fipFakeConnector := &fakeConnector{
			getObjectObj:         NewEmptyFixedAddress(isIPv6),
			getObjectQueryParams: queryParams,
			getObjectRef:         "",
			resultObject: []FixedAddress{*NewFixedAddress(
				netviewName, "",
				GetIPAddressFromRef(fakeRefReturn), cidr, macAddr,
				nil, nil, fakeRefReturn, isIPv6, comment, nil, nil, nil, nil, false, nil, false)},
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(fipFakeConnector, cmpType, tenantID)

		var actualIP *FixedAddress
		var err error
		It("should pass expected Fixed Address Object to GetObject", func() {
			actualIP, err = objMgr.GetFixedAddress(netviewName, cidr, ipAddr, isIPv6, macAddr)
		})
		It("should return expected Fixed Address Object", func() {
			Expect(*actualIP).To(Equal(fipFakeConnector.resultObject.([]FixedAddress)[0]))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get IPv6 Fixed Address", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "private"
		cidr := "2001:db8:abcd:0012::0/64"
		ipAddr := "2001:db8:abcd:0012::1"
		refIp := "2001%3Adb8%3Aabcd%3A0012%3A%3A1"
		duid := "01:23:45:67:80:ab"
		isIPv6 := true
		comment := "test"
		fakeRefReturn := fmt.Sprintf("ipv6fixedaddress/ZG5zLmJpbmRfY25h:%s/private", refIp)

		queryParams := NewQueryParams(
			false,
			map[string]string{
				"network_view": netviewName,
				"network":      cidr,
				"ipv6addr":     ipAddr,
				"duid":         duid,
			})

		fipFakeConnector := &fakeConnector{
			getObjectObj:         NewEmptyFixedAddress(isIPv6),
			getObjectQueryParams: queryParams,
			getObjectRef:         "",
			resultObject: []FixedAddress{*NewFixedAddress(
				netviewName, "",
				ipAddr, cidr, duid,
				nil, nil, fakeRefReturn, isIPv6, comment, nil, nil, nil, nil, false, nil, false)},
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(fipFakeConnector, cmpType, tenantID)

		var actualIP *FixedAddress
		var err error
		It("should pass expected Fixed Address Object to GetObject", func() {
			actualIP, err = objMgr.GetFixedAddress(netviewName, cidr, ipAddr, isIPv6, duid)
		})
		It("should return expected Fixed Address Object", func() {
			Expect(*actualIP).To(Equal(fipFakeConnector.resultObject.([]FixedAddress)[0]))
			Expect(err).To(BeNil())
		})
	})

	Describe("Update IPv4 Fixed Address", func() {
		var (
			err       error
			objMgr    IBObjectManager
			conn      *fakeConnector
			ref       string
			actualObj *FixedAddress
		)

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "default"
		refBase := "ZG5zLm5ldHdvcmtfdmlldyQyMw"
		ipv4Cidr := "10.2.1.0/20"
		ipv4Addr := "10.2.1.1"
		ipv6Cidr := "2001:db8:abcd:14::/64"
		ipv6CidrRef := "2003%3Adb8%3AAabcd%3A14%3A%3A1"
		name := "test"
		updateName := "test1"
		macAddr := "01:23:45:67:80:ab"
		client := "MAC_ADDRESS"
		updateMacAddr := "02:24:46:69:80:cd"
		duid := "01:23:45:67:80:ab"
		updateDuid := "02:24:46:69:80:cd"

		It("IPv4, updating name, MAC Address, comment and EAs", func() {
			ref = fmt.Sprintf("fixedaddress/%s:%s/%s", refBase, ipv4Addr, netviewName)
			updateIp := "10.0.0.3"
			initialEas := EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_old_value",
				"ea3": "ea3_value",
				"ea4": "ea4_value",
				"ea5": "ea5_old_value"}
			initObj := NewFixedAddress(netviewName, name, "10.0.0.2", ipv4Cidr, macAddr, &client, initialEas, ref, false, "old comment", nil, nil, nil, nil, false, nil, false)
			initObj.Ref = ref

			setEas := EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_new_value",
				"ea2": "ea2_new_value",
				"ea5": "ea5_old_value"}
			expectedEas := setEas

			comment := "test comment 1"
			updateObjIn := NewFixedAddress("", updateName, updateIp, "", updateMacAddr, &client, expectedEas, ref, false, comment, nil, nil, nil, nil, false, nil, false)
			updateObjIn.Ref = ref

			expectedObj := NewFixedAddress("", updateName, updateIp, "", updateMacAddr, &client, expectedEas, ref, false, comment, nil, nil, nil, nil, false, nil, false)
			expectedObj.Ref = ref

			conn = &fakeConnector{
				getObjectObj:         NewEmptyFixedAddress(false),
				getObjectQueryParams: NewQueryParams(false, nil),
				getObjectRef:         ref,
				getObjectError:       nil,
				resultObject:         expectedObj,

				updateObjectObj:   updateObjIn,
				updateObjectRef:   ref,
				updateObjectError: nil,

				fakeRefReturn: ref,
			}
			objMgr = NewObjectManager(conn, cmpType, tenantID)

			actualObj, err = objMgr.UpdateFixedAddress(ref, "", updateName, "", updateIp, "MAC_ADDRESS", updateMacAddr, comment, setEas, "", "", nil, "", false, nil, false)
			Expect(err).To(BeNil())
			Expect(actualObj).To(BeEquivalentTo(expectedObj))
		})

		It("Negative case: Update fails if a valid match client value is not passed", func() {
			ref = fmt.Sprintf("fixedaddress/%s:%s/%s", refBase, ipv4Addr, netviewName)
			matchClient := "MAC"
			initObj := NewFixedAddress("", name, "", "", macAddr, &matchClient, nil, ref, false, "", nil, nil, nil, nil, false, nil, false)
			initObj.Ref = ref

			comment := "test comment 1"

			conn = &fakeConnector{
				getObjectObj:         NewEmptyFixedAddress(false),
				getObjectQueryParams: NewQueryParams(false, nil),
				getObjectRef:         ref,
				getObjectError:       fmt.Errorf("test error"),
				updateObjectError:    fmt.Errorf("wrong value for match_client passed %s \n ", matchClient),
				fakeRefReturn:        ref,
			}
			objMgr = NewObjectManager(conn, cmpType, tenantID)

			var expectedObj *FixedAddress
			expectedObj = nil
			actualObj, err = objMgr.UpdateFixedAddress(ref, "", updateName, "", "", matchClient, updateMacAddr, comment, nil, "", "", nil, "", false, nil, false)
			Expect(actualObj).To(Equal(expectedObj))
			Expect(err).To(Equal(conn.updateObjectError))
		})

		It("IPv6, updating name, MAC Address, comment and EAs", func() {
			ref = fmt.Sprintf("ipv6fixedaddress/%s:%s/%s", refBase, ipv6CidrRef, netviewName)
			updateIp := "2001:db8:abcd:14::2"
			initialEas := EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_old_value",
				"ea3": "ea3_value",
				"ea4": "ea4_value",
				"ea5": "ea5_old_value"}
			initObj := NewFixedAddress(netviewName, name, "2001:db8:abcd:14::1", ipv6Cidr, duid, nil, initialEas, ref, true, "old comment", nil, nil, nil, nil, false, nil, false)
			initObj.Ref = ref

			setEas := EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_new_value",
				"ea2": "ea2_new_value",
				"ea5": "ea5_old_value"}
			expectedEas := setEas

			comment := "test comment 1"
			updateObjIn := NewFixedAddress("", updateName, updateIp, "", updateDuid, nil, expectedEas, ref, true, comment, nil, nil, nil, nil, false, nil, false)
			updateObjIn.Ref = ref

			expectedObj := NewFixedAddress("", updateName, updateIp, "", updateDuid, nil, expectedEas, ref, true, comment, nil, nil, nil, nil, false, nil, false)
			expectedObj.Ref = ref

			conn = &fakeConnector{
				getObjectObj:         NewEmptyFixedAddress(true),
				getObjectQueryParams: NewQueryParams(false, nil),
				getObjectRef:         ref,
				getObjectError:       nil,
				resultObject:         expectedObj,

				updateObjectObj:   updateObjIn,
				updateObjectRef:   ref,
				updateObjectError: nil,

				fakeRefReturn: ref,
			}
			objMgr = NewObjectManager(conn, cmpType, tenantID)

			actualObj, err = objMgr.UpdateFixedAddress(ref, "", updateName, "", updateIp, "", updateDuid, comment, setEas, "", "", nil, "", false, nil, false)
			Expect(err).To(BeNil())
			Expect(actualObj).To(BeEquivalentTo(expectedObj))
		})
	})

	Describe("Delete Fixed Address", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "private"
		cidr := "83.0.101.0/24"
		ipAddr := "83.0.101.68"
		macAddr := "01:23:45:67:80:ab"
		isIPv6 := false
		comment := "test"
		fakeRefReturn := fmt.Sprintf("fixedaddress/ZG5zLmJpbmRfY25h:%s/private", ipAddr)

		queryParams := NewQueryParams(
			false,
			map[string]string{
				"network_view": netviewName,
				"network":      cidr,
				"ipv4addr":     ipAddr,
				"mac":          macAddr,
			})

		fipFakeConnector := &fakeConnector{
			getObjectObj:         NewEmptyFixedAddress(isIPv6),
			getObjectQueryParams: queryParams,
			getObjectRef:         "",
			resultObject: []FixedAddress{*NewFixedAddress(
				netviewName,
				"",
				GetIPAddressFromRef(fakeRefReturn),
				cidr,
				macAddr,
				nil,
				nil,
				fakeRefReturn, isIPv6, comment, nil, nil, nil, nil, false, nil, false)},
			deleteObjectRef: fakeRefReturn,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(fipFakeConnector, cmpType, tenantID)

		var actualRef string
		var err error
		It("should pass expected Fixed Address Object to GetObject and DeleteObject", func() {
			actualRef, err = objMgr.ReleaseIP(netviewName, cidr, ipAddr, isIPv6, macAddr)
		})
		It("should return expected Fixed Address Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})
	Describe("Create Fixed Address with match client as CIRCUIT", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "private"
		cidr := "20.12.12.0/24"
		ipAddr := "20.12.12.4"
		isIPv6 := false
		comment := "test"
		matchClient := "CLIENT_ID"
		agentCircuitId := "43"
		fakeRefReturn := fmt.Sprintf("fixedaddress/ZG5zLmJpbmRfY25h:%s/private", ipAddr)
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		name := "testvm"
		ea := EA{"VM ID": vmID, "VM Name": vmName}
		conn := &fakeConnector{
			createObjectObj: NewFixedAddress(
				netviewName, name,
				ipAddr, cidr, "", &matchClient, ea, "", isIPv6, comment, &agentCircuitId, nil, nil, nil, false, nil, false),
			getObjectObj:         NewEmptyFixedAddress(isIPv6),
			getObjectRef:         fakeRefReturn,
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject: NewFixedAddress(
				netviewName, name,
				GetIPAddressFromRef(fakeRefReturn), cidr, "",
				&matchClient, ea, fakeRefReturn, isIPv6, comment, &agentCircuitId, nil, nil, nil, false, nil, false),
			fakeRefReturn: fakeRefReturn,
		}
		objMgr := NewObjectManager(conn, cmpType, tenantID)
		var actualIP *FixedAddress
		var err error
		It("should pass expected Fixed Address Object to CreateObject", func() {
			actualIP, err = objMgr.AllocateIP(netviewName, cidr, ipAddr, isIPv6, "", name, comment, ea, matchClient, agentCircuitId, "", nil, "", false, nil, false)
		})
		It("should return expected Fixed Address Object", func() {
			Expect(actualIP).To(Equal(conn.resultObject))
			Expect(err).To(BeNil())
		})
	})
	Describe("Update Fixed address by changing the match client to CLIENT_ID", func() {
		var (
			err       error
			objMgr    IBObjectManager
			conn      *fakeConnector
			ref       string
			actualObj *FixedAddress
		)
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "private"
		refBase := "ZG5zLm5ldHdvcmtfdmlldyQyMw"
		ipv4Cidr := "10.2.1.0/20"
		ipv4Addr := "10.2.1.1"
		name := "test"
		updateName := "test1"
		comment := "test comment 1"
		updateComment := "test comment 2"
		oldClients := "MAC_ADDRESS"
		updateClients := "CLIENT_ID"
		updateDhcpClientIdentifier := "43"
		updateClientIdentifierPrependZero := true
		options := []*Dhcpoption{
			{
				Name:        "routers",
				Num:         3,
				VendorClass: "DHCP",
				Value:       "10.2.1.2",
				UseOption:   true,
			},
		}
		useOptions := true
		UpdateOptions := []*Dhcpoption{
			{
				Name:        "dhcp-lease-time",
				Num:         51,
				VendorClass: "DHCP",
				Value:       "43200",
				UseOption:   false,
			},
		}
		eas := EA{
			"ea0": "ea0_old_value",
		}
		ref = fmt.Sprintf("fixedaddress/%s:%s/%s", refBase, ipv4Addr, netviewName)
		initObj := NewFixedAddress(netviewName, name, "10.0.0.2", ipv4Cidr, "", &oldClients, eas, ref, false, comment, nil, nil, nil, nil, false, options, useOptions)
		initObj.Ref = ref
		setEas := EA{
			"ea0": "ea0_new_value",
		}
		updateIp := "10.2.1.3"
		expectedEas := setEas
		updateObjIn := NewFixedAddress("", updateName, updateIp, "", "", &updateClients, expectedEas, ref, false, updateComment, nil, nil, &updateClientIdentifierPrependZero, &updateDhcpClientIdentifier, false, UpdateOptions, false)
		updateObjIn.Ref = ref

		expectedObj := NewFixedAddress("", updateName, updateIp, "", "", &updateClients, expectedEas, ref, false, updateComment, nil, nil, &updateClientIdentifierPrependZero, &updateDhcpClientIdentifier, false, UpdateOptions, false)
		expectedObj.Ref = ref
		conn = &fakeConnector{
			getObjectObj:         NewEmptyFixedAddress(false),
			getObjectQueryParams: NewQueryParams(false, nil),
			getObjectRef:         ref,
			getObjectError:       nil,
			resultObject:         expectedObj,

			updateObjectObj:   updateObjIn,
			updateObjectRef:   ref,
			updateObjectError: nil,

			fakeRefReturn: ref,
		}
		objMgr = NewObjectManager(conn, cmpType, tenantID)

		actualObj, err = objMgr.UpdateFixedAddress(ref, "", updateName, "", updateIp, updateClients, "", updateComment, setEas, "", "", &updateClientIdentifierPrependZero, updateDhcpClientIdentifier, false, UpdateOptions, false)
		Expect(err).To(BeNil())
		Expect(actualObj).To(BeEquivalentTo(expectedObj))
	})
})
