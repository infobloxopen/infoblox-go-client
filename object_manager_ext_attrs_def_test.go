package ibclient

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Object Manager: EA definition", func() {
	Describe("Create EA Definition", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		comment := "Test Extensible Attribute"
		flags := "CGV"
		listValues := []*EADefListValue{{"True"}, {"False"}}
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

	Describe("Get EA Definition", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		comment := "Test Extensible Attribute"
		flags := "CGV"
		listValues := []*EADefListValue{{"True"}, {"False"}}
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
})
