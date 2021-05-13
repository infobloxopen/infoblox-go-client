package ibclient

import (
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type fakeConnector struct {
	// expected object to be passed to CreateObject()
	createObjectObj interface{}

	// expected object and reference to be passed to GetObject()
	getObjectObj         interface{}
	getObjectQueryParams interface{}
	getObjectRef         string

	// expected object and reference to be passed to UpdateObject()
	updateObjectObj interface{}
	updateObjectRef string

	// expected object's reference to be passed to DeleteObject()
	deleteObjectRef string

	// An object to be returned by GetObject() method.
	resultObject interface{}

	// A reference to be returned by Create/Update/Delete (not Get) methods.
	fakeRefReturn string

	// Error which fake Connector is to return on appropriate method call.
	createObjectError error
	getObjectError    error
	updateObjectError error
	deleteObjectError error
}

func (c *fakeConnector) CreateObject(obj IBObject) (string, error) {
	Expect(obj).To(Equal(c.createObjectObj))

	return c.fakeRefReturn, c.createObjectError
}

func (c *fakeConnector) GetObject(obj IBObject, ref string, qp *QueryParams, res interface{}) (err error) {
	Expect(obj).To(Equal(c.getObjectObj))
	Expect(qp).To(Equal(c.getObjectQueryParams))
	Expect(ref).To(Equal(c.getObjectRef))

	if ref == "" {
		switch obj.(type) {
		case *NetworkView:
			*res.(*[]NetworkView) = c.resultObject.([]NetworkView)
		case *NetworkContainer:
			*res.(*[]NetworkContainer) = c.resultObject.([]NetworkContainer)
		case *Network:
			*res.(*[]Network) = c.resultObject.([]Network)
		case *FixedAddress:
			*res.(*[]FixedAddress) = c.resultObject.([]FixedAddress)
		case *EADefinition:
			*res.(*[]EADefinition) = c.resultObject.([]EADefinition)
		case *CapacityReport:
			*res.(*[]CapacityReport) = c.resultObject.([]CapacityReport)
		case *UpgradeStatus:
			*res.(*[]UpgradeStatus) = c.resultObject.([]UpgradeStatus)
		case *Member:
			*res.(*[]Member) = c.resultObject.([]Member)
		case *Grid:
			*res.(*[]Grid) = c.resultObject.([]Grid)
		case *License:
			*res.(*[]License) = c.resultObject.([]License)
		case *HostRecord:
			*res.(*[]HostRecord) = c.resultObject.([]HostRecord)
		case *ZoneDelegated:
			*res.(*[]ZoneDelegated) = c.resultObject.([]ZoneDelegated)
		}
	} else {
		switch obj.(type) {
		case *ZoneAuth:
			*res.(*ZoneAuth) = *c.resultObject.(*ZoneAuth)
		case *NetworkView:
			*res.(*NetworkView) = c.resultObject.(NetworkView)
		case *NetworkContainer:
			*res.(*NetworkContainer) = *c.resultObject.(*NetworkContainer)
		case *Network:
			*res.(*Network) = *c.resultObject.(*Network)
		case *FixedAddress:
			**res.(**FixedAddress) = *c.resultObject.(*FixedAddress)
		case *HostRecord:
			**res.(**HostRecord) = *c.resultObject.(*HostRecord)
		}
	}

	err = c.getObjectError
	return
}

func (c *fakeConnector) DeleteObject(ref string) (string, error) {
	Expect(ref).To(Equal(c.deleteObjectRef))

	return c.fakeRefReturn, c.deleteObjectError
}

func (c *fakeConnector) UpdateObject(obj IBObject, ref string) (string, error) {
	Expect(obj).To(Equal(c.updateObjectObj))
	Expect(ref).To(Equal(c.updateObjectRef))

	return c.fakeRefReturn, c.updateObjectError
}

var _ = Describe("Object Manager", func() {

	Describe("Create Network View", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "Default View"
		fakeRefReturn := "networkview/ZG5zLm5ldHdvcmtfdmlldyQyMw:global_view/false"
		nvFakeConnector := &fakeConnector{
			createObjectObj: NewNetworkView(NetworkView{Name: netviewName}),
			resultObject:    NewNetworkView(NetworkView{Name: netviewName, Ref: fakeRefReturn}),
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(nvFakeConnector, cmpType, tenantID)
		nvFakeConnector.createObjectObj.(*NetworkView).Ea = make(EA)
		nvFakeConnector.resultObject.(*NetworkView).Ea = make(EA)

		var actualNetworkView *NetworkView
		var err error
		It("should pass expected NetworkView Object to CreateObject", func() {
			actualNetworkView, err = objMgr.CreateNetworkView(netviewName)
		})
		It("should return expected NetworkView Object", func() {
			Expect(actualNetworkView).To(Equal(nvFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Update Network View", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "Global View"
		fakeRefReturn := "networkview/ZG5zLm5ldHdvcmtfdmlldyQyMw:global_view/false"

		returnGetObject := NetworkView{
			Name: netviewName,
			Ref:  fakeRefReturn,
			Ea:   EA{"network-name": "net1", "Lock": "Removed"}}
		returnUpdateObject := NetworkView{
			Name: netviewName,
			Ref:  fakeRefReturn,
			Ea:   EA{"network-name": "net2", "New": "Added"}}
		getObjectObj := &NetworkView{}
		getObjectObj.returnFields = []string{"extattrs"}
		nvFakeConnector := &fakeConnector{
			getObjectObj:         getObjectObj,
			getObjectQueryParams: NewQueryParams(false, nil),
			getObjectRef:         fakeRefReturn,
			fakeRefReturn:        fakeRefReturn,
			resultObject:         returnGetObject,
			updateObjectObj:      &returnUpdateObject,
			updateObjectRef:      fakeRefReturn,
		}

		objMgr := NewObjectManager(nvFakeConnector, cmpType, tenantID)

		var err error
		It("should pass expected updated object to UpdateObject", func() {
			setEas := EA{"network-name": "net2", "New": "Added"}
			err = objMgr.UpdateNetworkViewEA(fakeRefReturn, setEas)
		})
		It("should updated the GetObject with new EA and with no error", func() {
			Expect(err).To(BeNil())
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
			getObjIn.returnFields = []string{"extattrs", "network", "comment"}

			comment := "test comment 1"
			updateObjIn := NewNetwork(netviewName, cidr, false, comment, expectedEas)
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
			getObjIn.returnFields = []string{"extattrs", "network", "comment"}

			comment := "test comment 1"
			updateObjIn := NewNetwork(netviewName, cidr, false, comment, expectedEas)
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
			getObjIn.returnFields = []string{"extattrs", "network", "comment"}

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
			getObjIn.returnFields = []string{"extattrs", "network", "comment"}

			comment := "test comment 1"
			updateObjIn := NewNetwork(netviewName, cidr, false, comment, nil)
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
			getObjIn.returnFields = []string{"extattrs", "network", "comment"}

			comment := ""
			updateObjIn := NewNetwork(netviewName, cidr, false, comment, nil)
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
			getObjIn.returnFields = []string{"extattrs", "network", "comment"}

			comment := "test comment 1"
			updateObjIn := NewNetwork(netviewName, cidr, false, comment, expectedEas)
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
			getObjIn.returnFields = []string{"extattrs", "network", "comment"}

			comment := "test comment 1"
			updateObjIn := NewNetwork(netviewName, cidr, true, comment, expectedEas)
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
			updateObjIn := NewNetworkContainer(netviewName, ipv4Cidr, false, comment, expectedEas)
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
			updateObjIn := NewNetworkContainer(netviewName, ipv4Cidr, false, comment, nil)
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
			updateObjIn := NewNetworkContainer(netviewName, ipv6Cidr, true, comment, expectedEas)
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
				"", ea, "", isIPv6, comment),
			getObjectObj:         NewEmptyFixedAddress(isIPv6),
			getObjectRef:         fakeRefReturn,
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject: NewFixedAddress(
				netviewName, name,
				GetIPAddressFromRef(fakeRefReturn), cidr, macAddr,
				"", ea, fakeRefReturn, isIPv6, comment),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)

		var actualIP *FixedAddress
		var err error
		It("should pass expected Fixed Address Object to CreateObject", func() {
			actualIP, err = objMgr.AllocateIP(netviewName, cidr, ipAddr, isIPv6, macAddr, name, comment, ea)
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
				"", ea, "", isIPv6, comment),
			getObjectObj:         NewEmptyFixedAddress(isIPv6),
			getObjectRef:         fakeRefReturn,
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject: NewFixedAddress(
				netviewName, name,
				resultIP, cidr, macAddr,
				"", ea, fakeRefReturn, isIPv6, comment),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)

		var actualIP *FixedAddress
		var err error
		It("should pass expected Fixed Address Object to CreateObject", func() {
			actualIP, err = objMgr.AllocateIP(netviewName, cidr, "", isIPv6, macAddr, name, comment, ea)
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
				"", ea, "", isIPv6, comment),
			getObjectObj:         NewEmptyFixedAddress(isIPv6),
			getObjectRef:         fakeRefReturn,
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject: NewFixedAddress(
				netviewName, name,
				ipAddr, cidr, duid,
				"", ea, fakeRefReturn, isIPv6, comment),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)

		var actualIP *FixedAddress
		var err error
		It("should pass expected Fixed Address Object to CreateObject", func() {
			actualIP, err = objMgr.AllocateIP(netviewName, cidr, ipAddr, isIPv6, duid, name, comment, ea)
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
				"", ea, "", isIPv6, comment),
			getObjectObj:         NewEmptyFixedAddress(isIPv6),
			getObjectRef:         fakeRefReturn,
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject: NewFixedAddress(
				netviewName, name,
				resultIP, cidr, duid,
				"", ea, fakeRefReturn, isIPv6, comment),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)

		var actualIP *FixedAddress
		var err error
		It("should pass expected Fixed Address Object to CreateObject", func() {
			actualIP, err = objMgr.AllocateIP(netviewName, cidr, "", isIPv6, duid, name, comment, ea)
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
				"", ea, "", isIPv6, comment),
			createObjectError: fmt.Errorf("the DUID field cannot be left empty"),
			fakeRefReturn:     fakeRefReturn,
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)

		var actualIP *FixedAddress
		var err error
		It("should pass expected Fixed Address Object to CreateObject", func() {
			actualIP, err = objMgr.AllocateIP(netviewName, cidr, "", isIPv6, duid, name, comment, ea)
		})

		It("should return expected Fixed Address Object", func() {
			Expect(actualIP).To(Equal(expectedObj))
			Expect(err).To(Equal(conn.createObjectError))
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
		updateMacAddr := "02:24:46:69:80:cd"
		duid := "01:23:45:67:80:ab"
		updateDuid := "02:24:46:69:80:cd"

		It("IPv4, updating name, MAC Address, comment and EAs", func() {
			ref = fmt.Sprintf("fixedaddress/%s:%s/%s", refBase, ipv4Addr, netviewName)
			initialEas := EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_old_value",
				"ea3": "ea3_value",
				"ea4": "ea4_value",
				"ea5": "ea5_old_value"}
			initObj := NewFixedAddress(netviewName, name, "", ipv4Cidr, macAddr, "MAC_ADDRESS", initialEas, ref, false, "old comment")
			initObj.Ref = ref

			setEas := EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_new_value",
				"ea2": "ea2_new_value",
				"ea5": "ea5_old_value"}
			expectedEas := setEas

			getObjIn := &FixedAddress{}
			getObjIn.objectType = "fixedaddress"
			getObjIn.returnFields = []string{"extattrs", "ipv4addr", "mac", "name", "network", "network_view", "comment"}

			comment := "test comment 1"
			updateObjIn := NewFixedAddress("", updateName, "", "", updateMacAddr, "MAC_ADDRESS", expectedEas, ref, false, comment)
			updateObjIn.Ref = ref

			expectedObj := NewFixedAddress("", updateName, "", "", updateMacAddr, "MAC_ADDRESS", expectedEas, ref, false, comment)
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

			actualObj, err = objMgr.UpdateFixedAddress(ref, updateName, "MAC_ADDRESS", updateMacAddr, comment, setEas)
			Expect(err).To(BeNil())
			Expect(actualObj).To(BeEquivalentTo(expectedObj))
		})

		It("Negative case: Update fails if a valid match client value is not passed", func() {
			ref = fmt.Sprintf("fixedaddress/%s:%s/%s", refBase, ipv4Addr, netviewName)
			matchClient := "MAC"
			initObj := NewFixedAddress("", name, "", "", macAddr, matchClient, nil, ref, false, "")
			initObj.Ref = ref

			getObjIn := &FixedAddress{}
			getObjIn.objectType = "fixedaddress"
			getObjIn.returnFields = []string{"extattrs", "ipv4addr", "mac", "name", "network", "network_view", "comment"}

			comment := "test comment 1"

			conn = &fakeConnector{
				getObjectObj:         getObjIn,
				getObjectQueryParams: NewQueryParams(false, nil),
				getObjectRef:         ref,
				getObjectError:       fmt.Errorf("test error"),
				resultObject:         initObj,
				updateObjectError:    fmt.Errorf("wrong value for match_client passed %s \n ", matchClient),
				fakeRefReturn:        ref,
			}
			objMgr = NewObjectManager(conn, cmpType, tenantID)

			var expectedObj *FixedAddress
			expectedObj = nil
			actualObj, err = objMgr.UpdateFixedAddress(ref, updateName, matchClient, updateMacAddr, comment, nil)
			Expect(actualObj).To(Equal(expectedObj))
			Expect(err).To(Equal(conn.updateObjectError))
		})

		It("IPv6, updating name, MAC Address, comment and EAs", func() {
			ref = fmt.Sprintf("ipv6fixedaddress/%s:%s/%s", refBase, ipv6CidrRef, netviewName)
			initialEas := EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_old_value",
				"ea3": "ea3_value",
				"ea4": "ea4_value",
				"ea5": "ea5_old_value"}
			initObj := NewFixedAddress(netviewName, name, "", ipv6Cidr, duid, "", initialEas, ref, true, "old comment")
			initObj.Ref = ref

			setEas := EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_new_value",
				"ea2": "ea2_new_value",
				"ea5": "ea5_old_value"}
			expectedEas := setEas

			getObjIn := &FixedAddress{}
			getObjIn.objectType = "ipv6fixedaddress"
			getObjIn.returnFields = []string{"extattrs", "ipv6addr", "duid", "name", "network", "network_view", "comment"}

			comment := "test comment 1"
			updateObjIn := NewFixedAddress("", updateName, "", "", updateDuid, "", expectedEas, ref, true, comment)
			updateObjIn.Ref = ref

			expectedObj := NewFixedAddress("", updateName, "", "", updateDuid, "", expectedEas, ref, true, comment)
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

			actualObj, err = objMgr.UpdateFixedAddress(ref, updateName, "", updateDuid, comment, setEas)
			Expect(err).To(BeNil())
			Expect(actualObj).To(BeEquivalentTo(expectedObj))
		})
	})

	Describe("Allocate next available host Record without dns", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "private"
		ipv4Cidr := "53.0.0.0/24"
		macAddr := "01:23:45:67:80:ab"
		ipv4Addr := fmt.Sprintf("func:nextavailableip:%s,%s", ipv4Cidr, netviewName)
		ipv6Cidr := "2003:db8:abcd:14::/64"
		duid := "02:24:46:68:81:cd"
		ipv6Addr := fmt.Sprintf("func:nextavailableip:%s,%s", ipv6Cidr, netviewName)
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		recordName := "test"
		enabledns := false
		enabledhcp := false
		dnsView := "default"
		fakeRefReturn := fmt.Sprintf("record:host/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)
		resultIPV4Addrs := NewHostRecordIpv4Addr(ipv4Addr, macAddr, &enabledhcp, "")
		resultIPv6Addrs := NewHostRecordIpv6Addr(ipv6Addr, duid, &enabledhcp, "")
		enableDNS := new(bool)
		*enableDNS = enabledns
		comment := "test"
		aliases := []string{"abc.test.com"}

		eas := make(EA)
		eas["VM ID"] = vmID
		eas["VM Name"] = vmName

		aniFakeConnector := &fakeConnector{
			createObjectObj: NewHostRecord(
				netviewName, recordName,
				"", "", []HostRecordIpv4Addr{*resultIPV4Addrs}, []HostRecordIpv6Addr{*resultIPv6Addrs},
				eas, enableDNS, dnsView, "", "", comment, aliases),
			getObjectRef: fakeRefReturn,
			getObjectObj: NewHostRecord(
				netviewName, recordName,
				"", "", []HostRecordIpv4Addr{*resultIPV4Addrs}, []HostRecordIpv6Addr{*resultIPv6Addrs},
				eas, enableDNS, dnsView, "", fakeRefReturn, comment, aliases),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject: NewHostRecord(
				netviewName, recordName,
				"", "", []HostRecordIpv4Addr{*resultIPV4Addrs}, []HostRecordIpv6Addr{*resultIPv6Addrs},
				eas, enableDNS, dnsView, "", fakeRefReturn, comment, aliases),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualRecord *HostRecord
		var err error
		It("should pass expected host record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateHostRecord(
				enabledns, false, recordName,
				netviewName, dnsView,
				ipv4Cidr, ipv6Cidr, "", "", macAddr, duid, comment, eas, aliases)
		})
		It("should return expected host record Object", func() {
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Allocate next available host Record with dns", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "private"
		ipv4Cidr := "53.0.0.0/24"
		macAddr := "01:23:45:67:80:ab"
		ipv4Addr := fmt.Sprintf("func:nextavailableip:%s,%s", ipv4Cidr, netviewName)
		ipv6Cidr := "2003:db8:abcd:14::/64"
		duid := "02:24:46:68:81:cd"
		ipv6Addr := fmt.Sprintf("func:nextavailableip:%s,%s", ipv6Cidr, netviewName)
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		recordName := "test"
		enabledns := true
		enabledhcp := false
		dnsView := "default"
		fakeRefReturn := fmt.Sprintf("record:host/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)
		resultIPV4Addrs := NewHostRecordIpv4Addr(ipv4Addr, macAddr, &enabledhcp, "")
		resultIPV6Addrs := NewHostRecordIpv6Addr(ipv6Addr, duid, &enabledhcp, "")
		enableDNS := new(bool)
		*enableDNS = enabledns
		comment := "test"
		aliases := []string{"abc.test.com"}

		aniFakeConnector := &fakeConnector{
			createObjectObj: NewHostRecord(
				netviewName, recordName,
				"", "", []HostRecordIpv4Addr{*resultIPV4Addrs}, []HostRecordIpv6Addr{*resultIPV6Addrs},
				nil, enableDNS, dnsView, "", "", comment, aliases),
			getObjectRef: fakeRefReturn,
			getObjectObj: NewHostRecord(
				netviewName, recordName,
				"", "", []HostRecordIpv4Addr{*resultIPV4Addrs}, []HostRecordIpv6Addr{*resultIPV6Addrs},
				nil, enableDNS, dnsView, "", fakeRefReturn, comment, aliases),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject: NewHostRecord(
				netviewName, recordName,
				"", "", []HostRecordIpv4Addr{*resultIPV4Addrs}, []HostRecordIpv6Addr{*resultIPV6Addrs},
				nil, enableDNS, dnsView, "", fakeRefReturn, comment, aliases),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		ea := make(EA)
		aniFakeConnector.createObjectObj.(*HostRecord).Ea = ea
		aniFakeConnector.createObjectObj.(*HostRecord).Ea["VM ID"] = vmID
		aniFakeConnector.createObjectObj.(*HostRecord).Ea["VM Name"] = vmName

		aniFakeConnector.resultObject.(*HostRecord).Ea = ea
		aniFakeConnector.resultObject.(*HostRecord).Ea["VM ID"] = vmID
		aniFakeConnector.resultObject.(*HostRecord).Ea["VM Name"] = vmName

		aniFakeConnector.getObjectObj.(*HostRecord).Ea = ea
		aniFakeConnector.getObjectObj.(*HostRecord).Ea["VM ID"] = vmID
		aniFakeConnector.getObjectObj.(*HostRecord).Ea["VM Name"] = vmName

		var actualRecord *HostRecord
		var err error
		It("should pass expected host record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateHostRecord(
				enabledns, false, recordName, netviewName, dnsView, "", "", ipv4Addr, ipv6Addr, macAddr, duid, comment, ea, aliases)
		})
		It("should return expected host record Object", func() {
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Allocate specific host Record without dns", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "private"
		ipv4Cidr := "53.0.0.0/24"
		macAddr := "01:23:45:67:80:ab"
		ipv4Addr := "53.0.0.1"
		ipv6Cidr := "2003:db8:abcd:14::/64"
		duid := "02:24:46:68:81:cd"
		ipv6Addr := "2003:db8:abcd:14::1"
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		enabledns := false
		enabledhcp := false
		dnsView := "default"
		recordName := "test"
		fakeRefReturn := fmt.Sprintf("record:host/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)
		resultIPV4Addrs := NewHostRecordIpv4Addr(ipv4Addr, macAddr, &enabledhcp, "")
		resultIPV6Addrs := NewHostRecordIpv6Addr(ipv6Addr, duid, &enabledhcp, "")
		enableDNS := new(bool)
		*enableDNS = enabledns
		comment := "test"
		aliases := []string{"test1"}

		aniFakeConnector := &fakeConnector{
			createObjectObj: NewHostRecord(
				netviewName, recordName,
				"", "", []HostRecordIpv4Addr{*resultIPV4Addrs}, []HostRecordIpv6Addr{*resultIPV6Addrs},
				nil, enableDNS, dnsView, "", "", comment, aliases),
			getObjectRef: fakeRefReturn,
			getObjectObj: NewHostRecord(
				netviewName, recordName,
				"", "", []HostRecordIpv4Addr{*resultIPV4Addrs}, []HostRecordIpv6Addr{*resultIPV6Addrs},
				nil, enableDNS, dnsView, "", fakeRefReturn, comment, aliases),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject: NewHostRecord(
				netviewName, recordName,
				"", "", []HostRecordIpv4Addr{*resultIPV4Addrs}, []HostRecordIpv6Addr{*resultIPV6Addrs},
				nil, enableDNS, dnsView, "", fakeRefReturn, comment, aliases),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		ea := make(EA)
		aniFakeConnector.createObjectObj.(*HostRecord).Ea = ea
		aniFakeConnector.createObjectObj.(*HostRecord).Ea["VM ID"] = vmID
		aniFakeConnector.createObjectObj.(*HostRecord).Ea["VM Name"] = vmName

		aniFakeConnector.resultObject.(*HostRecord).Ea = ea
		aniFakeConnector.resultObject.(*HostRecord).Ea["VM ID"] = vmID
		aniFakeConnector.resultObject.(*HostRecord).Ea["VM Name"] = vmName

		aniFakeConnector.getObjectObj.(*HostRecord).Ea = ea
		aniFakeConnector.getObjectObj.(*HostRecord).Ea["VM ID"] = vmID
		aniFakeConnector.getObjectObj.(*HostRecord).Ea["VM Name"] = vmName

		var actualRecord *HostRecord
		var err error
		It("should pass expected host record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateHostRecord(
				enabledns, false, recordName, netviewName, dnsView, ipv4Cidr, ipv6Cidr, ipv4Addr, ipv6Addr, macAddr, duid, comment, ea, aliases)
		})

		It("should return expected host record Object", func() {
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Allocate specific host Record with dns", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "private"
		ipv4Cidr := "53.0.0.0/24"
		macAddr := "01:23:45:67:80:ab"
		ipv4Addr := "53.0.0.1"
		ipv6Cidr := "2003:db8:abcd:14::/64"
		duid := "02:24:46:68:81:cd"
		ipv6Addr := "2003:db8:abcd:14::1"
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		enabledns := true
		enabledhcp := false
		dnsView := "default"
		recordName := "test"
		fakeRefReturn := fmt.Sprintf("record:host/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)
		resultIPV4Addrs := NewHostRecordIpv4Addr(ipv4Addr, macAddr, &enabledhcp, "")
		resultIPV6Addrs := NewHostRecordIpv6Addr(ipv6Addr, duid, &enabledhcp, "")
		enableDNS := new(bool)
		*enableDNS = enabledns
		comment := "test"
		aliases := []string{"abc.test.com"}

		aniFakeConnector := &fakeConnector{
			createObjectObj: NewHostRecord(
				netviewName, recordName,
				"", "", []HostRecordIpv4Addr{*resultIPV4Addrs}, []HostRecordIpv6Addr{*resultIPV6Addrs},
				nil, enableDNS, dnsView, "", "", comment, aliases),
			getObjectRef: fakeRefReturn,
			getObjectObj: NewHostRecord(
				netviewName, recordName,
				"", "", []HostRecordIpv4Addr{*resultIPV4Addrs}, []HostRecordIpv6Addr{*resultIPV6Addrs},
				nil, enableDNS, dnsView, "", fakeRefReturn, comment, aliases),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject: NewHostRecord(
				netviewName, recordName,
				"", "", []HostRecordIpv4Addr{*resultIPV4Addrs}, []HostRecordIpv6Addr{*resultIPV6Addrs},
				nil, enableDNS, dnsView, "", fakeRefReturn, comment, aliases),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		ea := make(EA)
		aniFakeConnector.createObjectObj.(*HostRecord).Ea = ea
		aniFakeConnector.createObjectObj.(*HostRecord).Ea["VM ID"] = vmID
		aniFakeConnector.createObjectObj.(*HostRecord).Ea["VM Name"] = vmName

		aniFakeConnector.resultObject.(*HostRecord).Ea = ea
		aniFakeConnector.resultObject.(*HostRecord).Ea["VM ID"] = vmID
		aniFakeConnector.resultObject.(*HostRecord).Ea["VM Name"] = vmName

		aniFakeConnector.getObjectObj.(*HostRecord).Ea = ea
		aniFakeConnector.getObjectObj.(*HostRecord).Ea["VM ID"] = vmID
		aniFakeConnector.getObjectObj.(*HostRecord).Ea["VM Name"] = vmName

		var actualRecord *HostRecord
		var err error
		It("should pass expected host record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateHostRecord(
				enabledns, false, recordName, netviewName, dnsView, ipv4Cidr, ipv6Cidr, ipv4Addr, ipv6Addr, macAddr, duid, comment, ea, aliases)
		})

		It("should return expected host record Object", func() {
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Allocate specific A Record ", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "private"
		cidr := "53.0.0.0/24"
		ipAddr := "53.0.0.1"
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		dnsView := "default"
		recordName := "test"
		fakeRefReturn := fmt.Sprintf("record:a/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)
		aniFakeConnector := &fakeConnector{
			createObjectObj: NewRecordA(
				dnsView, "", recordName, ipAddr, nil, ""),
			getObjectRef: fakeRefReturn,
			getObjectObj: NewRecordA(
				dnsView, "", recordName, ipAddr, nil, fakeRefReturn),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject: NewRecordA(
				dnsView, "", recordName, ipAddr, nil, fakeRefReturn),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		ea := make(EA)
		aniFakeConnector.createObjectObj.(*RecordA).Ea = ea
		aniFakeConnector.createObjectObj.(*RecordA).Ea["VM ID"] = vmID
		aniFakeConnector.createObjectObj.(*RecordA).Ea["VM Name"] = vmName

		aniFakeConnector.resultObject.(*RecordA).Ea = ea
		aniFakeConnector.resultObject.(*RecordA).Ea["VM ID"] = vmID
		aniFakeConnector.resultObject.(*RecordA).Ea["VM Name"] = vmName

		aniFakeConnector.getObjectObj.(*RecordA).Ea = ea
		aniFakeConnector.getObjectObj.(*RecordA).Ea["VM ID"] = vmID
		aniFakeConnector.getObjectObj.(*RecordA).Ea["VM Name"] = vmName

		var actualRecord *RecordA
		var err error
		It("should pass expected A record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateARecord(netviewName, dnsView, recordName, cidr, ipAddr, ea)
		})
		It("should return expected A record Object", func() {
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Allocate next available A Record ", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "private"
		cidr := "53.0.0.0/24"
		ipAddr := fmt.Sprintf("func:nextavailableip:%s,%s", cidr, netviewName)
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		dnsView := "default"
		recordName := "test"
		fakeRefReturn := fmt.Sprintf("record:a/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)

		aniFakeConnector := &fakeConnector{
			createObjectObj: NewRecordA(
				dnsView, "", recordName, ipAddr, nil, ""),
			getObjectRef: fakeRefReturn,
			getObjectObj: NewRecordA(
				dnsView, "", recordName, ipAddr, nil, fakeRefReturn),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject: NewRecordA(
				dnsView, "", recordName, ipAddr, nil, fakeRefReturn),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		ea := make(EA)
		aniFakeConnector.createObjectObj.(*RecordA).Ea = ea
		aniFakeConnector.createObjectObj.(*RecordA).Ea["VM ID"] = vmID
		aniFakeConnector.createObjectObj.(*RecordA).Ea["VM Name"] = vmName

		aniFakeConnector.resultObject.(*RecordA).Ea = ea
		aniFakeConnector.resultObject.(*RecordA).Ea["VM ID"] = vmID
		aniFakeConnector.resultObject.(*RecordA).Ea["VM Name"] = vmName

		aniFakeConnector.getObjectObj.(*RecordA).Ea = ea
		aniFakeConnector.getObjectObj.(*RecordA).Ea["VM ID"] = vmID
		aniFakeConnector.getObjectObj.(*RecordA).Ea["VM Name"] = vmName

		var actualRecord *RecordA
		var err error
		It("should pass expected A record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateARecord(netviewName, dnsView, recordName, cidr, ipAddr, ea)
		})
		It("should return expected A record Object", func() {
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Allocate specific PTR Record ", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "private"
		cidr := "53.0.0.0/24"
		ipAddr := "53.0.0.1"
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		dnsView := "default"
		recordName := "test"
		fakeRefReturn := fmt.Sprintf("record:ptr/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)

		aniFakeConnector := &fakeConnector{
			createObjectObj: NewRecordPTR(RecordPTR{
				PtrdName: recordName,
				View:     dnsView,
				Ipv4Addr: ipAddr,
			}),
			getObjectRef: fakeRefReturn,
			getObjectObj: NewRecordPTR(RecordPTR{
				PtrdName: recordName,
				View:     dnsView,
				Ipv4Addr: ipAddr,
				Ref:      fakeRefReturn,
			}),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject: NewRecordPTR(RecordPTR{
				PtrdName: recordName,
				View:     dnsView,
				Ipv4Addr: ipAddr,
				Ref:      fakeRefReturn,
			}),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		ea := make(EA)
		aniFakeConnector.createObjectObj.(*RecordPTR).Ea = ea
		aniFakeConnector.createObjectObj.(*RecordPTR).Ea["VM ID"] = vmID
		aniFakeConnector.createObjectObj.(*RecordPTR).Ea["VM Name"] = vmName

		aniFakeConnector.resultObject.(*RecordPTR).Ea = ea
		aniFakeConnector.resultObject.(*RecordPTR).Ea["VM ID"] = vmID
		aniFakeConnector.resultObject.(*RecordPTR).Ea["VM Name"] = vmName

		aniFakeConnector.getObjectObj.(*RecordPTR).Ea = ea
		aniFakeConnector.getObjectObj.(*RecordPTR).Ea["VM ID"] = vmID
		aniFakeConnector.getObjectObj.(*RecordPTR).Ea["VM Name"] = vmName

		var actualRecord *RecordPTR
		var err error
		It("should pass expected PTR record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreatePTRRecord(netviewName, dnsView, recordName, cidr, ipAddr, ea)
		})
		It("should return expected PTR record Object", func() {
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Allocate next available PTR Record ", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "private"
		cidr := "53.0.0.0/24"
		ipAddr := fmt.Sprintf("func:nextavailableip:%s,%s", cidr, netviewName)
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		dnsView := "default"
		recordName := "test"
		fakeRefReturn := fmt.Sprintf("record:ptr/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)

		aniFakeConnector := &fakeConnector{
			createObjectObj: NewRecordPTR(RecordPTR{
				PtrdName: recordName,
				View:     dnsView,
				Ipv4Addr: ipAddr,
			}),
			getObjectRef: fakeRefReturn,
			getObjectObj: NewRecordPTR(RecordPTR{
				PtrdName: recordName,
				View:     dnsView,
				Ipv4Addr: ipAddr,
				Ref:      fakeRefReturn,
			}),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject: NewRecordPTR(RecordPTR{
				PtrdName: recordName,
				View:     dnsView,
				Ipv4Addr: ipAddr,
				Ref:      fakeRefReturn,
			}),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		ea := make(EA)
		aniFakeConnector.createObjectObj.(*RecordPTR).Ea = ea
		aniFakeConnector.createObjectObj.(*RecordPTR).Ea["VM ID"] = vmID
		aniFakeConnector.createObjectObj.(*RecordPTR).Ea["VM Name"] = vmName

		aniFakeConnector.resultObject.(*RecordPTR).Ea = ea
		aniFakeConnector.resultObject.(*RecordPTR).Ea["VM ID"] = vmID
		aniFakeConnector.resultObject.(*RecordPTR).Ea["VM Name"] = vmName

		aniFakeConnector.getObjectObj.(*RecordPTR).Ea = ea
		aniFakeConnector.getObjectObj.(*RecordPTR).Ea["VM ID"] = vmID
		aniFakeConnector.getObjectObj.(*RecordPTR).Ea["VM Name"] = vmName

		var actualRecord *RecordPTR
		var err error
		It("should pass expected PTR record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreatePTRRecord(netviewName, dnsView, recordName, cidr, ipAddr, ea)
		})
		It("should return expected PTR record Object", func() {
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Allocate CNAME Record ", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		canonical := "test-canonical"
		dnsView := "default"
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		recordName := "test"
		fakeRefReturn := fmt.Sprintf("record:cname/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)

		aniFakeConnector := &fakeConnector{
			createObjectObj: NewRecordCNAME(RecordCNAME{
				Name:      recordName,
				View:      dnsView,
				Canonical: canonical,
			}),
			getObjectRef: fakeRefReturn,
			getObjectObj: NewRecordCNAME(RecordCNAME{
				Name:      recordName,
				View:      dnsView,
				Canonical: canonical,
				Ref:       fakeRefReturn,
			}),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject: NewRecordCNAME(RecordCNAME{
				Name:      recordName,
				View:      dnsView,
				Canonical: canonical,
				Ref:       fakeRefReturn,
			}),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)
		ea := make(EA)
		aniFakeConnector.createObjectObj.(*RecordCNAME).Ea = ea
		aniFakeConnector.createObjectObj.(*RecordCNAME).Ea["VM ID"] = vmID
		aniFakeConnector.createObjectObj.(*RecordCNAME).Ea["VM Name"] = vmName

		aniFakeConnector.resultObject.(*RecordCNAME).Ea = ea
		aniFakeConnector.resultObject.(*RecordCNAME).Ea["VM ID"] = vmID
		aniFakeConnector.resultObject.(*RecordCNAME).Ea["VM Name"] = vmName

		aniFakeConnector.getObjectObj.(*RecordCNAME).Ea = ea
		aniFakeConnector.getObjectObj.(*RecordCNAME).Ea["VM ID"] = vmID
		aniFakeConnector.getObjectObj.(*RecordCNAME).Ea["VM Name"] = vmName
		var actualRecord *RecordCNAME
		var err error
		It("should pass expected CNAME record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateCNAMERecord(canonical, recordName, dnsView, ea)
		})
		It("should return expected CNAME record Object", func() {
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Allocate TXT Record ", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		text := "test-text"
		dnsView := "default"
		recordName := "test"
		ttl := 30
		fakeRefReturn := fmt.Sprintf("record:txt/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)

		aniFakeConnector := &fakeConnector{
			createObjectObj: NewRecordTXT(RecordTXT{
				Name: recordName,
				Text: text,
				TTL:  ttl,
				View: dnsView,
			}),
			getObjectRef: fakeRefReturn,
			getObjectObj: NewRecordTXT(RecordTXT{
				Name: recordName,
				Text: text,
				View: dnsView,
				Ref:  fakeRefReturn,
			}),
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject: NewRecordTXT(RecordTXT{
				Name: recordName,
				Text: text,
				View: dnsView,
				TTL:  ttl,
				Ref:  fakeRefReturn,
			}),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualRecord *RecordTXT
		var err error
		It("should pass expected TXT record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateTXTRecord(recordName, text, 30, dnsView)
		})
		It("should return expected TXT record Object", func() {
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Create EA Definition", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		comment := "Test Extensible Attribute"
		flags := "CGV"
		listValues := []EADefListValue{"True", "False"}
		name := "TestEA"
		eaType := "string"
		allowedTypes := []string{"arecord", "aaarecord", "ptrrecord"}
		ead := EADefinition{
			Name:               name,
			Comment:            comment,
			Flags:              flags,
			ListValues:         listValues,
			Type:               eaType,
			AllowedObjectTypes: allowedTypes}
		fakeRefReturn := "extensibleattributedef/ZG5zLm5ldHdvcmtfdmlldyQyMw:TestEA"
		eadFakeConnector := &fakeConnector{
			createObjectObj: NewEADefinition(ead),
			resultObject:    NewEADefinition(ead),
			fakeRefReturn:   fakeRefReturn,
		}
		eadFakeConnector.resultObject.(*EADefinition).Ref = fakeRefReturn

		objMgr := NewObjectManager(eadFakeConnector, cmpType, tenantID)

		var actualEADef *EADefinition
		var err error
		It("should pass expected EA Definintion Object to CreateObject", func() {
			actualEADef, err = objMgr.CreateEADefinition(ead)
		})
		It("should return expected EA Definition Object", func() {
			Expect(actualEADef).To(Equal(eadFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get Network View", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "Default View"
		fakeRefReturn := "networkview/ZG5zLm5ldHdvcmtfdmlldyQyMw:global_view/false"

		queryParams := NewQueryParams(
			false,
			map[string]string{
				"name": netviewName,
			})

		nvFakeConnector := &fakeConnector{
			getObjectObj:         NewNetworkView(NetworkView{}),
			getObjectQueryParams: queryParams,
			getObjectRef:         "",
			resultObject:         []NetworkView{*NewNetworkView(NetworkView{Name: netviewName, Ref: fakeRefReturn})},
		}

		objMgr := NewObjectManager(nvFakeConnector, cmpType, tenantID)

		var actualNetworkView *NetworkView
		var err error
		It("should pass expected NetworkView Object to GetObject", func() {
			actualNetworkView, err = objMgr.GetNetworkView(netviewName)
		})
		It("should return expected NetworkView Object", func() {
			Expect(*actualNetworkView).To(Equal(nvFakeConnector.resultObject.([]NetworkView)[0]))
			Expect(err).To(BeNil())
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
				"", nil, fakeRefReturn, isIPv6, comment)},
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
				"", nil, fakeRefReturn, isIPv6, comment)},
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

	Describe("Get Ipv4 and IPv6 Host Record Without DNS", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		hostName := "test"
		ipv4Addr := "10.0.0.1"
		ipv6Addr := "2001:db8:abcd:14::1"
		fakeRefReturn := fmt.Sprintf("record:host/ZG5zLmJpbmRfY25h:%s/%20%20", hostName)
		queryParams := NewQueryParams(
			false,
			map[string]string{
				"name":     hostName,
				"ipv4addr": "10.0.0.1",
				"ipv6addr": "2001:db8:abcd:14::1",
			})
		fipFakeConnector := &fakeConnector{
			getObjectObj:         NewEmptyHostRecord(),
			getObjectQueryParams: queryParams,
			getObjectRef:         "",
			resultObject: []HostRecord{*NewHostRecord(
				"", hostName, ipv4Addr, ipv6Addr, nil, nil,
				nil, nil, "", "", fakeRefReturn, "", []string{})},
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(fipFakeConnector, cmpType, tenantID)

		var actualhostRecord *HostRecord
		var err error
		It("should pass expected Host record Object to GetObject", func() {
			actualhostRecord, err = objMgr.GetHostRecord(hostName, ipv4Addr, ipv6Addr)
		})

		It("should return expected Host record Object", func() {
			Expect(*actualhostRecord).To(Equal(fipFakeConnector.resultObject.([]HostRecord)[0]))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get Host record by reference", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		hostName := "test"
		fakeRefReturn := fmt.Sprintf("record:host/ZG5zLmJpbmRfY25h:%s/%20%20", hostName)
		resObj := NewEmptyHostRecord()
		resObj.Ref = fakeRefReturn
		ncFakeConnector := &fakeConnector{
			getObjectObj:         NewEmptyHostRecord(),
			getObjectRef:         fakeRefReturn,
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         resObj,
			fakeRefReturn:        fakeRefReturn,
		}

		objMgr := NewObjectManager(ncFakeConnector, cmpType, tenantID)

		var actualRec *HostRecord
		var err error
		It("should pass expected Network Object to GetObject", func() {
			actualRec, err = objMgr.GetHostRecordByRef(fakeRefReturn)
		})
		It("should return expected Network Object", func() {
			Expect(err).To(BeNil())
			Expect(*actualRec).To(Equal(*resObj))
		})
	})

	Describe("Update host record", func() {
		var (
			err       error
			objMgr    IBObjectManager
			conn      *fakeConnector
			ref       string
			actualObj *HostRecord
		)

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		hostName := "host.test.com"
		refBase := "ZG5zLm5ldHdvcmtfdmlldyQyMw"
		ipv4Addr := "10.0.0.3"
		ipv6Addr := "2003:db8:abcd:14::/64"

		It("Updating name, comment, aliases and EAs", func() {
			enableDNS := new(bool)
			*enableDNS = true
			ref = fmt.Sprintf("record:host/%s:%s", refBase, hostName)
			initialEas := EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_old_value",
				"ea3": "ea3_value",
				"ea4": "ea4_value",
				"ea5": "ea5_old_value"}
			initialAliases := []string{"abc.test.com", "xyz.test.com"}
			initObj := NewHostRecord("", hostName, "", "", []HostRecordIpv4Addr{},
				[]HostRecordIpv6Addr{}, initialEas, enableDNS, "", "", "", "old comment", initialAliases)
			initObj.Ref = ref

			setEas := EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_new_value",
				"ea2": "ea2_new_value",
				"ea5": "ea5_old_value"}
			expectedEas := setEas
			expectedAliases := []string{"abc.test.com", "trial.test.com"}

			getObjIn := &HostRecord{}
			getObjIn.objectType = "record:host"
			getObjIn.returnFields = []string{"extattrs", "ipv4addrs", "ipv6addrs", "name", "view", "zone", "comment", "network_view", "aliases"}

			comment := "test comment 1"
			updateObjIn := NewHostRecord("", "host1.test.com", "", "", []HostRecordIpv4Addr{},
				[]HostRecordIpv6Addr{}, expectedEas, enableDNS, "", "", "", comment, expectedAliases)
			updateObjIn.Ref = ref

			expectedObj := NewHostRecord("", "host1.test.com", "", "", []HostRecordIpv4Addr{},
				[]HostRecordIpv6Addr{}, expectedEas, enableDNS, "", "", "", comment, expectedAliases)
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

			actualObj, err = objMgr.UpdateHostRecord(ref, true, false, "host1.test.com", "",
				"", "", "", comment, setEas, expectedAliases)
			Expect(err).To(BeNil())
			Expect(*actualObj).To(BeEquivalentTo(*expectedObj))
		})

		It("Updating MAC Address and DUID when IPv4 and Ipv6 addresses are passed", func() {
			enableDNS := new(bool)
			*enableDNS = false
			enableDHCP := new(bool)
			*enableDHCP = false
			macAddr := "01:23:45:67:80:ab"
			duid := "02:24:46:68:81:cd"
			resultIPV4Addrs := NewHostRecordIpv4Addr(ipv4Addr, macAddr, enableDHCP, "")
			resultIPV6Addrs := NewHostRecordIpv6Addr(ipv6Addr, duid, enableDHCP, "")
			ref = fmt.Sprintf("record:host/%s:%s", refBase, hostName)

			initObj := NewHostRecord("", hostName, ipv4Addr, ipv6Addr, []HostRecordIpv4Addr{*resultIPV4Addrs},
				[]HostRecordIpv6Addr{*resultIPV6Addrs}, nil, enableDNS, "", "", ref, "", []string{})

			getObjIn := &HostRecord{}
			getObjIn.objectType = "record:host"
			getObjIn.returnFields = []string{"extattrs", "ipv4addrs", "ipv6addrs", "name", "view", "zone", "comment", "network_view", "aliases"}

			updateObjIn := NewHostRecord("", hostName, "", "", []HostRecordIpv4Addr{*resultIPV4Addrs},
				[]HostRecordIpv6Addr{*resultIPV6Addrs}, nil, enableDNS, "", "", "", "", []string{})
			updateObjIn.Ref = ref

			expectedObj := NewHostRecord("", hostName, "", "", []HostRecordIpv4Addr{*resultIPV4Addrs},
				[]HostRecordIpv6Addr{*resultIPV6Addrs}, nil, enableDNS, "", "", "", "", []string{})
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

			actualObj, err = objMgr.UpdateHostRecord(ref, false, false, hostName, ipv4Addr,
				ipv6Addr, macAddr, duid, "", nil, []string{})
			Expect(err).To(BeNil())
			Expect(*actualObj).To(BeEquivalentTo(*expectedObj))
		})
	})

	Describe("Get EA Definition", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		comment := "Test Extensible Attribute"
		flags := "CGV"
		listValues := []EADefListValue{"True", "False"}
		name := "TestEA"
		eaType := "string"
		allowedTypes := []string{"arecord", "aaarecord", "ptrrecord"}
		ead := EADefinition{
			Name: name,
		}
		fakeRefReturn := "extensibleattributedef/ZG5zLm5ldHdvcmtfdmlldyQyMw:TestEA"
		eadRes := EADefinition{
			Name:               name,
			Comment:            comment,
			Flags:              flags,
			ListValues:         listValues,
			Type:               eaType,
			AllowedObjectTypes: allowedTypes,
			Ref:                fakeRefReturn,
		}

		queryParams := NewQueryParams(
			false,
			map[string]string{
				"name": name,
			})

		eadFakeConnector := &fakeConnector{
			getObjectObj:         NewEADefinition(ead),
			getObjectRef:         "",
			getObjectQueryParams: queryParams,
			resultObject:         []EADefinition{*NewEADefinition(eadRes)},
			fakeRefReturn:        fakeRefReturn,
		}

		objMgr := NewObjectManager(eadFakeConnector, cmpType, tenantID)

		var actualEADef *EADefinition
		var err error
		It("should pass expected EA Definintion Object to GetObject", func() {
			actualEADef, err = objMgr.GetEADefinition(name)
		})
		It("should return expected EA Definition Object", func() {
			Expect(*actualEADef).To(Equal(eadFakeConnector.resultObject.([]EADefinition)[0]))
			Expect(err).To(BeNil())
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

	Describe("Delete Network View", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		fakeRefReturn := "networkview/ZG5zLm5ldHdvcmtfdmlldyQyMw:global_view/false"
		deleteRef := fakeRefReturn
		nwFakeConnector := &fakeConnector{
			deleteObjectRef: deleteRef,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(nwFakeConnector, cmpType, tenantID)

		var actualRef string
		var err error
		It("should pass expected Network View Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteNetworkView(deleteRef)
		})
		It("should return expected Network View Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
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
				"",
				nil,
				fakeRefReturn, isIPv6, comment)},
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

	Describe("Delete Host Record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		hostName := "test"
		deleteRef := fmt.Sprintf("record:host/ZG5zLmJpbmRfY25h:%s/%20%20", hostName)
		fakeRefReturn := deleteRef
		nwFakeConnector := &fakeConnector{
			deleteObjectRef: deleteRef,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(nwFakeConnector, cmpType, tenantID)

		var actualRef string
		var err error
		It("should pass expected Host record Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteHostRecord(deleteRef)
		})
		It("should return expected Host record Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete A Record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		recordName := "test"
		deleteRef := fmt.Sprintf("record:a/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)
		fakeRefReturn := deleteRef
		nwFakeConnector := &fakeConnector{
			deleteObjectRef: deleteRef,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(nwFakeConnector, cmpType, tenantID)

		var actualRef string
		var err error
		It("should pass expected A record Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteARecord(deleteRef)
		})
		It("should return expected A record Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete PTR Record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		recordName := "test"
		deleteRef := fmt.Sprintf("record:ptr/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)
		fakeRefReturn := deleteRef
		nwFakeConnector := &fakeConnector{
			deleteObjectRef: deleteRef,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(nwFakeConnector, cmpType, tenantID)

		var actualRef string
		var err error
		It("should pass expected PTR record Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeletePTRRecord(deleteRef)
		})
		It("should return expected PTR record Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete CNAME Record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		recordName := "test"
		deleteRef := fmt.Sprintf("record:CNAME/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)
		fakeRefReturn := deleteRef
		nwFakeConnector := &fakeConnector{
			deleteObjectRef: deleteRef,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(nwFakeConnector, cmpType, tenantID)

		var actualRef string
		var err error
		It("should pass expected CNAME record Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteCNAMERecord(deleteRef)
		})
		It("should return expected CNAME record Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete TXT Record", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		recordName := "test"
		deleteRef := fmt.Sprintf("record:txt/ZG5zLmJpbmRfY25h:%s/%20%20", recordName)
		fakeRefReturn := deleteRef
		nwFakeConnector := &fakeConnector{
			deleteObjectRef: deleteRef,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(nwFakeConnector, cmpType, tenantID)

		var actualRef string
		var err error
		It("should pass expected TXT record Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteTXTRecord(deleteRef)
		})
		It("should return expected TXT record Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})

	Describe("BuildNetworkViewFromRef", func() {
		netviewName := "default_view"
		netviewRef := fmt.Sprintf("networkview/ZG5zLm5ldHdvcmtfdmlldyQyMw:%s/false", netviewName)

		expectedNetworkView := NetworkView{Ref: netviewRef, Name: netviewName}
		It("should return expected Network View Object", func() {
			Expect(*BuildNetworkViewFromRef(netviewRef)).To(Equal(expectedNetworkView))
		})
		It("should failed if bad Network View Ref is provided", func() {
			Expect(BuildNetworkViewFromRef("bad")).To(BeNil())
		})
	})

	Describe("BuildNetworkFromRef", func() {
		netviewName := "test_view"
		cidr := "23.11.0.0/24"
		networkRef := fmt.Sprintf("network/ZG5zLm5ldHdvcmtfdmlldyQyMw:%s/%s", cidr, netviewName)

		expectedNetwork := Network{Ref: networkRef, NetviewName: netviewName, Cidr: cidr}
		expectedNetwork.objectType = "network"
		expectedNetwork.returnFields = []string{"extattrs", "network", "comment"}
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
		expectedNetwork.returnFields = []string{"extattrs", "network", "comment"}
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

	Describe("Get Capacity report", func() {
		cmpType := "Heka"
		tenantID := "0123"
		var name string = "Member1"
		fakeRefReturn := fmt.Sprintf("member/ZG5zLmJpbmRfY25h:/%s", name)

		queryParams := NewQueryParams(
			false,
			map[string]string{
				"name": name,
			})

		fakeConnector := &fakeConnector{
			getObjectObj:         NewCapcityReport(CapacityReport{}),
			getObjectQueryParams: queryParams,
			getObjectRef:         "",
			resultObject: []CapacityReport{*NewCapcityReport(CapacityReport{
				Ref:  fakeRefReturn,
				Name: name,
			})},
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(fakeConnector, cmpType, tenantID)

		var actualReport []CapacityReport
		var err error

		It("should pass expected Capacityreport object to GetObject", func() {
			actualReport, err = objMgr.GetCapacityReport(name)
		})
		It("should return expected CapacityReport Object", func() {
			Expect(actualReport[0]).To(Equal(fakeConnector.resultObject.([]CapacityReport)[0]))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get upgrade status", func() {
		cmpType := "Heka"
		tenantID := "0123"
		var StatusType = "GRID"
		fakeRefReturn := fmt.Sprintf("upgradestatus/Li51cGdyYWRlc3RhdHVzJHVwZ3JhZGVfc3RhdHVz:test")

		queryParams := NewQueryParams(
			false,
			map[string]string{
				"type": StatusType,
			})

		USFakeConnector := &fakeConnector{
			getObjectObj:         NewUpgradeStatus(UpgradeStatus{}),
			getObjectQueryParams: queryParams,
			getObjectRef:         "",
			resultObject: []UpgradeStatus{*NewUpgradeStatus(UpgradeStatus{
				Ref:  fakeRefReturn,
				Type: StatusType,
			})},
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(USFakeConnector, cmpType, tenantID)

		var actualStatus []UpgradeStatus
		var err error

		It("should pass expected upgradestatus object to GetObject", func() {
			actualStatus, err = objMgr.GetUpgradeStatus(StatusType)
		})
		It("should return expected upgradestatus Object", func() {
			Expect(actualStatus[0]).To(Equal(USFakeConnector.resultObject.([]UpgradeStatus)[0]))
			Expect(err).To(BeNil())
		})

	})
	Describe("Get upgrade status Error case", func() {
		cmpType := "Heka"
		tenantID := "0123"
		StatusType := ""
		fakeRefReturn := fmt.Sprintf("upgradestatus/Li51cGdyYWRlc3RhdHVzJHVwZ3JhZGVfc3RhdHVz:test")
		expectErr := errors.New("Status type can not be nil")
		USFakeConnector := &fakeConnector{
			getObjectObj:         NewUpgradeStatus(UpgradeStatus{Type: StatusType}),
			getObjectRef:         "",
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject: []UpgradeStatus{*NewUpgradeStatus(UpgradeStatus{
				Ref:  fakeRefReturn,
				Type: StatusType,
			})},
			fakeRefReturn: fakeRefReturn,
		}
		objMgr := NewObjectManager(USFakeConnector, cmpType, tenantID)
		It("upgradestatus object to GetObject", func() {
			_, err := objMgr.GetUpgradeStatus(StatusType)
			Expect(err).To(Equal(expectErr))
		})

	})
	Describe("GetAllMembers", func() {
		cmpType := "Heka"
		tenantID := "0123"
		var err error
		fakeRefReturn := fmt.Sprintf("member/Li51cGdyYWRlc3RhdHVzJHVwZ3JhZGVfc3RhdHVz:test")
		returnFields := []string{"host_name", "node_info", "time_zone"}
		MemFakeConnector := &fakeConnector{
			getObjectObj:         NewMember(Member{}),
			getObjectRef:         "",
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject: []Member{*NewMember(Member{
				Ref: fakeRefReturn,
			})},
			fakeRefReturn: fakeRefReturn,
		}
		objMgr := NewObjectManager(MemFakeConnector, cmpType, tenantID)
		var actualMembers []Member
		It("should return expected member Object", func() {
			actualMembers, err = objMgr.GetAllMembers()
			Expect(actualMembers[0]).To(Equal(MemFakeConnector.resultObject.([]Member)[0]))
			Expect(actualMembers[0].returnFields).To(Equal(returnFields))
			Expect(err).To(BeNil())
		})
	})

	Describe("GetGridInfo", func() {
		cmpType := "Heka"
		tenantID := "0123"
		var err error
		fakeRefReturn := fmt.Sprintf("grid/Li511cGdyYWRlc3RhdHVzJHVwZ3JhZGVfc3RhdHVz:test")
		returnFields := []string{"name", "ntp_setting"}
		GridFakeConnector := &fakeConnector{
			getObjectObj:         NewGrid(Grid{}),
			getObjectRef:         "",
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject: []Grid{*NewGrid(Grid{
				Ref: fakeRefReturn,
			})},
			fakeRefReturn: fakeRefReturn,
		}
		objMgr := NewObjectManager(GridFakeConnector, cmpType, tenantID)
		var actualGridInfo []Grid
		It("should return expected Grid Object", func() {
			actualGridInfo, err = objMgr.GetGridInfo()
			Expect(actualGridInfo[0]).To(Equal(GridFakeConnector.resultObject.([]Grid)[0]))
			Expect(actualGridInfo[0].returnFields).To(Equal(returnFields))
			Expect(err).To(BeNil())
		})
	})

	Describe("GetGridLicense", func() {
		cmpType := "Heka"
		tenantID := "0123"
		var err error
		fakeRefReturn := fmt.Sprintf("license/Li511cGdyYWRlc3RhdHVzJHVwZ3JhZGVfc3RhdHVz:test")
		returnFields := []string{"expiration_status",
			"expiry_date",
			"key",
			"limit",
			"limit_context",
			"type"}
		LicFakeConnector := &fakeConnector{
			getObjectObj:         NewGridLicense(License{}),
			getObjectRef:         "",
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject: []License{*NewGridLicense(License{
				Ref: fakeRefReturn,
			})},
			fakeRefReturn: fakeRefReturn,
		}
		objMgr := NewObjectManager(LicFakeConnector, cmpType, tenantID)
		var actualGridLicense []License
		It("should return expected License Object", func() {
			actualGridLicense, err = objMgr.GetGridLicense()
			Expect(actualGridLicense[0]).To(Equal(LicFakeConnector.resultObject.([]License)[0]))
			Expect(actualGridLicense[0].returnFields).To(Equal(returnFields))
			Expect(err).To(BeNil())
		})
	})

	Describe("Create Zone Auth", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		fqdn := "azone.example.com"
		fakeRefReturn := "zone_auth/ZG5zLnpvbmUkLl9kZWZhdWx0LnphLmNvLmFic2EuY2Fhcy5vaG15Z2xiLmdzbGJpYmNsaWVudA:dzone.example.com/default"
		zaFakeConnector := &fakeConnector{
			createObjectObj: NewZoneAuth(ZoneAuth{Fqdn: fqdn}),
			resultObject:    NewZoneAuth(ZoneAuth{Fqdn: fqdn, Ref: fakeRefReturn}),
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(zaFakeConnector, cmpType, tenantID)

		ea := make(EA)

		zaFakeConnector.createObjectObj.(*ZoneAuth).Ea = ea
		zaFakeConnector.createObjectObj.(*ZoneAuth).Ea["Tenant ID"] = tenantID
		zaFakeConnector.createObjectObj.(*ZoneAuth).Ea["CMP Type"] = cmpType

		zaFakeConnector.resultObject.(*ZoneAuth).Ea = ea
		zaFakeConnector.resultObject.(*ZoneAuth).Ea["Tenant ID"] = tenantID
		zaFakeConnector.resultObject.(*ZoneAuth).Ea["CMP Type"] = cmpType

		var actualZoneAuth *ZoneAuth
		var err error
		It("should pass expected ZoneAuth Object to CreateObject", func() {
			actualZoneAuth, err = objMgr.CreateZoneAuth(fqdn, ea)
		})
		It("should return expected ZoneAuth Object", func() {
			Expect(actualZoneAuth).To(Equal(zaFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get AuthZone by ref", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		fqdn := "azone.example.com"
		fakeRefReturn := "zone_delegated/ZG5zLnpvbmUkLl9kZWZhdWx0LnphLmNvLmFic2EuY2Fhcy5vaG15Z2xiLmdzbGJpYmNsaWVudA:azone.example.com/default"
		zdFakeConnector := &fakeConnector{
			getObjectObj:         NewZoneAuth(ZoneAuth{}),
			getObjectRef:         fakeRefReturn,
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         NewZoneAuth(ZoneAuth{Fqdn: fqdn}),
		}

		objMgr := NewObjectManager(zdFakeConnector, cmpType, tenantID)

		var actualZoneAuth, getNoRef *ZoneAuth
		getNoRef = nil
		var err error
		It("should pass expected ZoneAuth Object to GetObject", func() {
			actualZoneAuth, err = objMgr.GetZoneAuthByRef(fakeRefReturn)
		})
		fmt.Printf("doodo  %s", actualZoneAuth)
		It("should return expected ZoneAuth Object", func() {
			Expect(actualZoneAuth).To(Equal(zdFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
		It("should return empty ZoneAuth and nil error if ref is empty", func() {
			zdFakeConnector.getObjectObj.(*ZoneAuth).IBBase.objectType = ""
			zdFakeConnector.getObjectObj.(*ZoneAuth).IBBase.returnFields = nil
			actualZoneAuth, err = objMgr.GetZoneAuthByRef("")
			Expect(actualZoneAuth).To(Equal(getNoRef))
			Expect(err).ToNot(BeNil())
		})
	})

	Describe("Delete ZoneAuth", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		deleteRef := "zone_auth/ZG5zLnpvbmUkLl9kZWZhdWx0LnphLmNvLmFic2EuY2Fhcy5vaG15Z2xiLmdzbGJpYmNsaWVudA:dzone.example.com/default"
		fakeRefReturn := deleteRef
		zaFakeConnector := &fakeConnector{
			deleteObjectRef: deleteRef,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(zaFakeConnector, cmpType, tenantID)

		var actualRef string
		var err error
		It("should pass expected ZoneAuth Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteZoneAuth(deleteRef)
		})
		It("should return expected ZoneAuth Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get Zone Delegated", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		fqdn := "dzone.example.com"
		fakeRefReturn := "zone_delegated/ZG5zLnpvbmUkLl9kZWZhdWx0LnphLmNvLmFic2EuY2Fhcy5vaG15Z2xiLmdzbGJpYmNsaWVudA:dzone.example.com/default"

		queryParams := NewQueryParams(
			false,
			map[string]string{
				"fqdn": fqdn,
			})

		zdFakeConnector := &fakeConnector{
			getObjectObj:         NewZoneDelegated(ZoneDelegated{}),
			getObjectQueryParams: queryParams,
			getObjectRef:         "",
			resultObject:         []ZoneDelegated{*NewZoneDelegated(ZoneDelegated{Fqdn: fqdn, Ref: fakeRefReturn})},
		}

		objMgr := NewObjectManager(zdFakeConnector, cmpType, tenantID)

		var actualZoneDelegated *ZoneDelegated
		var err error
		It("should pass expected ZoneDelegated Object to GetObject", func() {
			actualZoneDelegated, err = objMgr.GetZoneDelegated(fqdn)
		})
		It("should return expected ZoneDelegated Object", func() {
			Expect(*actualZoneDelegated).To(Equal(zdFakeConnector.resultObject.([]ZoneDelegated)[0]))
			Expect(err).To(BeNil())
		})
		It("should return nil if fqdn is empty", func() {
			zdFakeConnector.getObjectObj.(*ZoneDelegated).Fqdn = ""
			actualZoneDelegated, err = objMgr.GetZoneDelegated("")
			Expect(actualZoneDelegated).To(BeNil())
			Expect(err).To(BeNil())
		})
	})

	Describe("Create Zone Delegated", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		fqdn := "dzone.example.com"
		delegateTo := []NameServer{
			{Address: "10.0.0.1", Name: "test1.dzone.example.com"},
			{Address: "10.0.0.2", Name: "test2.dzone.example.com"}}
		fakeRefReturn := "zone_delegated/ZG5zLnpvbmUkLl9kZWZhdWx0LnphLmNvLmFic2EuY2Fhcy5vaG15Z2xiLmdzbGJpYmNsaWVudA:dzone.example.com/default"
		zdFakeConnector := &fakeConnector{
			createObjectObj: NewZoneDelegated(ZoneDelegated{Fqdn: fqdn, DelegateTo: delegateTo}),
			resultObject:    NewZoneDelegated(ZoneDelegated{Fqdn: fqdn, DelegateTo: delegateTo, Ref: fakeRefReturn}),
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(zdFakeConnector, cmpType, tenantID)

		var actualZoneDelegated *ZoneDelegated
		var err error
		It("should pass expected ZoneDelegated Object to CreateObject", func() {
			actualZoneDelegated, err = objMgr.CreateZoneDelegated(fqdn, delegateTo)
		})
		It("should return expected ZoneDelegated Object", func() {
			Expect(actualZoneDelegated).To(Equal(zdFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Update Zone Delegated", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		fakeRefReturn := "zone_delegated/ZG5zLnpvbmUkLl9kZWZhdWx0LnphLmNvLmFic2EuY2Fhcy5vaG15Z2xiLmdzbGJpYmNsaWVudA:dzone.example.com/default"
		delegateTo := []NameServer{
			{Address: "10.0.0.1", Name: "test1.dzone.example.com"},
			{Address: "10.0.0.2", Name: "test2.dzone.example.com"}}

		receiveUpdateObject := NewZoneDelegated(ZoneDelegated{Ref: fakeRefReturn, DelegateTo: delegateTo})
		returnUpdateObject := NewZoneDelegated(ZoneDelegated{DelegateTo: delegateTo, Ref: fakeRefReturn})
		zdFakeConnector := &fakeConnector{
			fakeRefReturn:   fakeRefReturn,
			resultObject:    returnUpdateObject,
			updateObjectObj: receiveUpdateObject,
			updateObjectRef: fakeRefReturn,
		}

		objMgr := NewObjectManager(zdFakeConnector, cmpType, tenantID)

		var updatedObject *ZoneDelegated
		var err error
		It("should pass expected updated object to UpdateObject", func() {
			updatedObject, err = objMgr.UpdateZoneDelegated(fakeRefReturn, delegateTo)
		})
		It("should update zone with new delegation server list with no error", func() {
			Expect(updatedObject).To(Equal(returnUpdateObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete ZoneDelegated", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		deleteRef := "zone_delegated/ZG5zLnpvbmUkLl9kZWZhdWx0LnphLmNvLmFic2EuY2Fhcy5vaG15Z2xiLmdzbGJpYmNsaWVudA:dzone.example.com/default"
		fakeRefReturn := deleteRef
		zdFakeConnector := &fakeConnector{
			deleteObjectRef: deleteRef,
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := NewObjectManager(zdFakeConnector, cmpType, tenantID)

		var actualRef string
		var err error
		It("should pass expected ZoneDelegated Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteZoneDelegated(deleteRef)
		})
		It("should return expected ZoneDelegated Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})
})
