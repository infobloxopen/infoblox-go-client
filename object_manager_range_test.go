package ibclient

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Object Manager: range create", func() {
	Describe("Create IP Range", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		netviewName := "private"
		comment := "test"

		startAddr := "10.10.0.1"
		endAddr := "10.10.0.100"
		rangeName := "rangeObj"
		disable := false

		fakeRefReturn := fmt.Sprintf("range/ZG5zLmJpbmRfY25h:%s/%s/%s", startAddr, endAddr, netviewName)

		rangeObj := NewEmptyRange()
		rangeObj.StartAddr = &startAddr
		rangeObj.EndAddr = &endAddr
		rangeObj.Name = &rangeName
		rangeObj.Comment = &comment
		rangeObj.NetworkView = &netviewName
		rangeObj.Disable = &disable

		resultObj := NewEmptyRange()
		resultObj.StartAddr = &startAddr
		resultObj.EndAddr = &endAddr
		resultObj.Name = &rangeName
		resultObj.NetworkView = &netviewName
		resultObj.Comment = &comment
		resultObj.Disable = &disable

		getRangeObj := NewEmptyRange()
		getRangeObj.Ref = fakeRefReturn

		conn := &fakeConnector{
			createObjectObj:      rangeObj,
			getObjectObj:         getRangeObj,
			getObjectRef:         fakeRefReturn,
			getObjectQueryParams: NewQueryParams(false, nil),
			resultObject:         []Range{*resultObj},
			fakeRefReturn:        fakeRefReturn,
		}

		objMgr := NewObjectManager(conn, cmpType, tenantID)

		It("should pass expected Range to CreateObject", func() {
			createdRange, err := objMgr.CreateRange(netviewName, rangeName, startAddr, endAddr, comment, disable, nil)
			Expect(err).To(BeNil())
			Expect(createdRange.EndAddr).To(Equal(resultObj.EndAddr))
			Expect(createdRange.NetworkView).To(Equal(resultObj.NetworkView))
		})
	})
})
