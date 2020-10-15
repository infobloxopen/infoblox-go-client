package ibclient

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Testing RecordAliasOperations", func() {
	Context("RecordAlias object", func() {
		targetName := "arec.test.com"
		targetType := "A"
		name := "bind_a.domain.com"
		view := "default"
		zone := "domain.com"

		ra := NewRecordAlias(RecordAlias{
			TargetName: targetName,
			TargetType: targetType,
			Name:       name,
			View:       view,
			Zone:       zone})

		It("should set fields correctly", func() {
			Expect(ra.TargetName).To(Equal(targetName))
			Expect(ra.TargetType).To(Equal(targetType))
			Expect(ra.Name).To(Equal(name))
			Expect(ra.View).To(Equal(view))
			Expect(ra.Zone).To(Equal(zone))
		})

		It("should set base fields correctly", func() {
			Expect(ra.ObjectType()).To(Equal("record:alias"))
			Expect(ra.ReturnFields()).To(ConsistOf("target_name", "target_type", "name", "view", "zone", "extattrs", "comment",
				"dns_name", "ttl", "use_ttl"))
		})
	})

	Describe("Allocate specific Alias Record ", func() {
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		ea1 := EA{"VM ID": vmID, "VM Name": vmName}

		recAlias := RecordAlias{
			TargetName: "arec.test.com",
			TargetType: "A",
			View:       "default",
			Name:       "test",
			Ea:         ea1,
		}
		fakeRefReturn := fmt.Sprintf("record:alias/ZG5zLmJpbmRfY25h:%s/%20%20", recAlias.Name)
		aniFakeConnector := &fakeConnector{
			createObjectObj: NewRecordAlias(recAlias),
			getObjectRef:    fakeRefReturn,
			getObjectObj: NewRecordAlias(RecordAlias{
				Name:       recAlias.Name,
				View:       recAlias.View,
				TargetName: recAlias.TargetName,
				TargetType: recAlias.TargetType,
				Ref:        fakeRefReturn,
				Ea:         ea1,
			}),
			resultObject: NewRecordAlias(RecordAlias{
				TargetName: "arec.test.com",
				TargetType: "A",
				View:       "default",
				Name:       "test",
				Ea:         ea1,
				Ref:        fakeRefReturn,
			}),
			fakeRefReturn: fakeRefReturn,
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		ea := objMgr.getBasicEA(true)
		aniFakeConnector.createObjectObj.(*RecordAlias).Ea = ea
		aniFakeConnector.resultObject.(*RecordAlias).Ea = ea
		aniFakeConnector.resultObject.(*RecordAlias).Ea["VM ID"] = vmID
		aniFakeConnector.resultObject.(*RecordAlias).Ea["VM Name"] = vmName
		aniFakeConnector.getObjectObj.(*RecordAlias).Ea = ea

		var actualRecord *RecordAlias
		var err error
		It("should pass expected Alias record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateAliasRecord(recAlias)
		})

		It("should return expected Alias record Object", func() {
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get Alias Record by TargetName", func() {
		recAlias := RecordAlias{TargetName: "test.example.com"}
		name := "test"
		fakeRefReturn := fmt.Sprintf("record:alias/ZG5zLmJpbmRfY25h:%s/%20%20", name)
		aniFakeConnector := &fakeConnector{
			getObjectObj: NewRecordAlias(recAlias),
			getObjectRef: "",
			resultObject: []RecordAlias{*NewRecordAlias(RecordAlias{Name: recAlias.Name, Ref: fakeRefReturn, TargetName: recAlias.TargetName})},
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualRecord *[]RecordAlias
		var err error
		It("should pass expected Alias record Object to GetObject", func() {
			actualRecord, err = objMgr.GetAliasRecord(recAlias)

		})

		It("should return expected Alias record Object", func() {
			Expect(*actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get Alias Record by Name", func() {
		recAlias := RecordAlias{Name: "test"}
		fakeRefReturn := fmt.Sprintf("record:mx/ZG5zLmJpbmRfY25h:%s/%20%20", recAlias.Name)
		aniFakeConnector := &fakeConnector{
			getObjectObj: NewRecordAlias(recAlias),
			getObjectRef: "",
			resultObject: []RecordAlias{*NewRecordAlias(RecordAlias{Name: recAlias.Name, Ref: fakeRefReturn})},
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualRecord *[]RecordAlias
		var err error
		It("should pass expected Alias record Object to GetObject", func() {
			actualRecord, err = objMgr.GetAliasRecord(recAlias)

		})

		It("should return expected Alias record Object", func() {
			Expect(*actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get Alias Record by TargetType", func() {
		recAlias := RecordAlias{TargetType: "A"}
		name := "test"
		fakeRefReturn := fmt.Sprintf("record:alias/ZG5zLmJpbmRfY25h:%s/%20%20", name)
		aniFakeConnector := &fakeConnector{
			getObjectObj: NewRecordAlias(recAlias),
			getObjectRef: "",
			resultObject: []RecordAlias{*NewRecordAlias(RecordAlias{Name: recAlias.Name, Ref: fakeRefReturn, TargetType: recAlias.TargetType})},
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualRecord *[]RecordAlias
		var err error
		It("should pass expected Alias record Object to GetObject", func() {
			actualRecord, err = objMgr.GetAliasRecord(recAlias)

		})

		It("should return expected Alias record Object", func() {
			Expect(*actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get Alias Record by Reference", func() {
		name := "test"
		fakeRefReturn := fmt.Sprintf("record:alias/ZG5zLmJpbmRfY25h:%s/%20%20", name)
		recAlias := RecordAlias{Ref: fakeRefReturn}
		aniFakeConnector := &fakeConnector{
			getObjectObj: NewRecordAlias(recAlias),
			getObjectRef: fakeRefReturn,
			resultObject: []RecordAlias{*NewRecordAlias(RecordAlias{Name: recAlias.Name, Ref: fakeRefReturn})},
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualRecord *[]RecordAlias
		var err error
		It("should pass expected Alias record Object to GetObject", func() {
			actualRecord, err = objMgr.GetAliasRecord(recAlias)

		})

		It("should return expected Alias record Object", func() {
			Expect(*actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete Alias Record by Reference", func() {

		name := "delete_test"
		fakeRefReturn := fmt.Sprintf("record:alias/ZG5zLmJpbmRfY25h:%s/%20%20", name)
		recAlias := RecordAlias{Ref: fakeRefReturn}
		aniFakeConnector := &fakeConnector{
			deleteObjectRef: recAlias.Ref,
			fakeRefReturn:   fakeRefReturn,
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)
		var actualRef string
		var err error
		It("should pass expected Alais record Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteAliasRecord(recAlias)
		})
		It("should return expected Alais record Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete Alais Record by Name", func() {
		recAlias := RecordAlias{Name: "delete_test"}
		fakeRefReturn := fmt.Sprintf("record:alias/ZG5zLmJpbmRfY25h:%s/%20%20", recAlias.Name)
		aniFakeConnector := &fakeConnector{
			getObjectObj:    NewRecordAlias(recAlias),
			getObjectRef:    "",
			resultObject:    []RecordAlias{*NewRecordAlias(RecordAlias{Name: recAlias.Name, Ref: fakeRefReturn})},
			deleteObjectRef: fakeRefReturn,
			fakeRefReturn:   fakeRefReturn,
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)
		var actualRef string
		var err error
		It("should pass expected Alias record Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteAliasRecord(recAlias)
		})
		It("should return expected Alias record Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete Alias Record by TargetName", func() {
		recAlias := RecordAlias{TargetName: "test.example.com"}
		name := "delete_test"
		fakeRefReturn := fmt.Sprintf("record:mx/ZG5zLmJpbmRfY25h:%s/%20%20", name)
		aniFakeConnector := &fakeConnector{
			getObjectObj:    NewRecordAlias(recAlias),
			getObjectRef:    "",
			resultObject:    []RecordAlias{*NewRecordAlias(RecordAlias{Name: name, Ref: fakeRefReturn, TargetName: recAlias.TargetName})},
			deleteObjectRef: fakeRefReturn,
			fakeRefReturn:   fakeRefReturn,
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)
		var actualRef string
		var err error
		It("should pass expected Alias record Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteAliasRecord(recAlias)
		})
		It("should return expected Alias record Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})
})
