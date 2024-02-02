package ibclient_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/infobloxopen/infoblox-go-client/v2"
)

var _ = Describe("Object Manager: network view", func() {
	Describe("Create Network View", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "Default View"
		comment := "test client"
		setEas := ibclient.EA{"Cloud API Owned": true}
		fakeRefReturn := "networkview/ZG5zLm5ldHdvcmtfdmlldyQyMw:global_view/false"
		nvFakeConnector := &fakeConnector{
			createObjectObj: ibclient.NewNetworkView(netviewName, comment, setEas, ""),
			resultObject:    ibclient.NewNetworkView(netviewName, comment, setEas, fakeRefReturn),
			fakeRefReturn:   fakeRefReturn,
		}

		objMgr := ibclient.NewObjectManager(nvFakeConnector, cmpType, tenantID)

		var actualNetworkView *ibclient.NetworkView
		var err error
		It("should pass expected NetworkView Object to CreateObject", func() {
			actualNetworkView, err = objMgr.CreateNetworkView(netviewName, comment, setEas)
		})
		It("should return expected NetworkView Object", func() {
			Expect(actualNetworkView).To(Equal(nvFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get Network View", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "Default View"
		fakeRefReturn := "networkview/ZG5zLm5ldHdvcmtfdmlldyQyMw:global_view/false"

		queryParams := ibclient.NewQueryParams(
			false,
			map[string]string{
				"name": netviewName,
			})

		nvFakeConnector := &fakeConnector{
			getObjectObj:         ibclient.NewEmptyNetworkView(),
			getObjectQueryParams: queryParams,
			getObjectRef:         "",
			resultObject:         []ibclient.NetworkView{*ibclient.NewNetworkView(netviewName, "", nil, fakeRefReturn)},
		}

		objMgr := ibclient.NewObjectManager(nvFakeConnector, cmpType, tenantID)

		var actualNetworkView *ibclient.NetworkView
		var err error
		It("should pass expected NetworkView Object to GetObject", func() {
			actualNetworkView, err = objMgr.GetNetworkView(netviewName)
		})
		It("should return expected NetworkView Object", func() {
			Expect(*actualNetworkView).To(Equal(nvFakeConnector.resultObject.([]ibclient.NetworkView)[0]))
			Expect(err).To(BeNil())
		})
	})

	Describe("Update Network View", func() {
		var (
			err       error
			objMgr    ibclient.IBObjectManager
			conn      *fakeConnector
			ref       string
			actualObj *ibclient.NetworkView
		)

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "default"
		refBase := "ZG5zLm5ldHdvcmtfdmlldyQyMw"

		It("Updating comment and EAs", func() {
			ref = fmt.Sprintf("networkview/%s:%s", refBase, netviewName)
			initialEas := ibclient.EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_old_value",
				"ea3": "ea3_value",
				"ea4": "ea4_value",
				"ea5": "ea5_old_value"}
			initObj := ibclient.NewNetworkView(netviewName, "old comment", initialEas, ref)

			setEas := ibclient.EA{
				"ea0": "ea0_old_value",
				"ea1": "ea1_new_value",
				"ea2": "ea2_new_value",
				"ea5": "ea5_old_value"}
			expectedEas := setEas

			getObjIn := ibclient.NewEmptyNetworkView()

			comment := "test comment 1"
			updateNetviewName := "default_view"
			updatedRef := fmt.Sprintf("networkview/%s:%s", refBase, updateNetviewName)
			updateObjIn := ibclient.NewNetworkView(updateNetviewName, comment, expectedEas, ref)

			expectedObj := ibclient.NewNetworkView(updateNetviewName, comment, expectedEas, updatedRef)

			conn = &fakeConnector{
				getObjectObj:         getObjIn,
				getObjectQueryParams: ibclient.NewQueryParams(false, nil),
				getObjectRef:         ref,
				getObjectError:       nil,
				resultObject:         initObj,

				updateObjectObj:   updateObjIn,
				updateObjectRef:   ref,
				updateObjectError: nil,

				fakeRefReturn: updatedRef,
			}
			objMgr = ibclient.NewObjectManager(conn, cmpType, tenantID)

			actualObj, err = objMgr.UpdateNetworkView(ref, updateNetviewName, comment, setEas)
			Expect(err).To(BeNil())
			Expect(actualObj).To(BeEquivalentTo(expectedObj))
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

		objMgr := ibclient.NewObjectManager(nwFakeConnector, cmpType, tenantID)

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
})
