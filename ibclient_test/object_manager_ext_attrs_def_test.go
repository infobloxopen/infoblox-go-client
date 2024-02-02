package ibclient_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	ibclient "github.com/infobloxopen/infoblox-go-client/v2"
)

var _ = Describe("Object Manager: EA definition", func() {
	Describe("Create EA Definition", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		comment := "Test Extensible Attribute"
		flags := "CGV"
		listValues := []*ibclient.EADefListValue{{"True"}, {"False"}}
		name := "TestEA"
		eaType := "string"
		allowedTypes := []string{"arecord", "aaarecord", "ptrrecord"}
		ead := ibclient.EADefinition{
			Name:               &name,
			Comment:            &comment,
			Flags:              &flags,
			ListValues:         listValues,
			Type:               eaType,
			AllowedObjectTypes: allowedTypes}
		fakeRefReturn := "extensibleattributedef/ZG5zLm5ldHdvcmtfdmlldyQyMw:TestEA"
		eadFakeConnector := &fakeConnector{
			createObjectObj: ibclient.NewEADefinition(ead),
			resultObject:    ibclient.NewEADefinition(ead),
			fakeRefReturn:   fakeRefReturn,
		}
		eadFakeConnector.resultObject.(*ibclient.EADefinition).Ref = fakeRefReturn

		objMgr := ibclient.NewObjectManager(eadFakeConnector, cmpType, tenantID)

		var actualEADef *ibclient.EADefinition
		var err error
		It("should pass expected EA Definintion Object to CreateObject", func() {
			actualEADef, err = objMgr.CreateEADefinition(ead)
		})
		It("should return expected EA Definition Object", func() {
			Expect(actualEADef).To(Equal(eadFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get EA Definition", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		comment := "Test Extensible Attribute"
		flags := "CGV"
		listValues := []*ibclient.EADefListValue{{"True"}, {"False"}}
		name := "TestEA"
		eaType := "string"
		allowedTypes := []string{"arecord", "aaarecord", "ptrrecord"}
		ead := ibclient.EADefinition{
			Name: &name,
		}
		fakeRefReturn := "extensibleattributedef/ZG5zLm5ldHdvcmtfdmlldyQyMw:TestEA"
		eadRes := ibclient.EADefinition{
			Name:               &name,
			Comment:            &comment,
			Flags:              &flags,
			ListValues:         listValues,
			Type:               eaType,
			AllowedObjectTypes: allowedTypes,
			Ref:                fakeRefReturn,
		}

		queryParams := ibclient.NewQueryParams(
			false,
			map[string]string{
				"name": name,
			})

		eadFakeConnector := &fakeConnector{
			getObjectObj:         ibclient.NewEADefinition(ead),
			getObjectRef:         "",
			getObjectQueryParams: queryParams,
			resultObject:         []ibclient.EADefinition{*ibclient.NewEADefinition(eadRes)},
			fakeRefReturn:        fakeRefReturn,
		}

		objMgr := ibclient.NewObjectManager(eadFakeConnector, cmpType, tenantID)

		var actualEADef *ibclient.EADefinition
		var err error
		It("should pass expected EA Definintion Object to GetObject", func() {
			actualEADef, err = objMgr.GetEADefinition(name)
		})
		It("should return expected EA Definition Object", func() {
			Expect(*actualEADef).To(Equal(eadFakeConnector.resultObject.([]ibclient.EADefinition)[0]))
			Expect(err).To(BeNil())
		})
	})
})
