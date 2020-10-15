package ibclient

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Testing RecordSRVOperations", func() {
	Context("RecordSRV object", func() {
		port := uint(443)
		priority := uint(1)
		target := "xmpp-server.example.com"
		weight := uint(10)
		name := "bind_a.domain.com"
		view := "default"
		zone := "domain.com"

		ra := NewRecordSRV(RecordSRV{
			Port:     port,
			Priority: priority,
			Target:   target,
			Weight:   weight,
			Name:     name,
			View:     view,
			Zone:     zone})

		It("should set fields correctly", func() {
			Expect(ra.Priority).To(Equal(priority))
			Expect(ra.Port).To(Equal(port))
			Expect(ra.Target).To(Equal(target))
			Expect(ra.Weight).To(Equal(weight))
			Expect(ra.Name).To(Equal(name))
			Expect(ra.View).To(Equal(view))
			Expect(ra.Zone).To(Equal(zone))
		})

		It("should set base fields correctly", func() {
			Expect(ra.ObjectType()).To(Equal("record:srv"))
			Expect(ra.ReturnFields()).To(ConsistOf("port", "priority", "target", "weight", "name", "view", "zone", "extattrs", "comment", "creation_time",
				"ddns_protected", "dns_name", "forbid_reclamation", "reclaimable", "ttl", "use_ttl", "disable"))
		})
	})

	Describe("Allocate specific SRV Record ", func() {
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		ea1 := EA{"VM ID": vmID, "VM Name": vmName}

		recSRV := RecordSRV{
			Port:     443,
			Priority: 1,
			Target:   "xmpp-server.example.com",
			Weight:   10,
			View:     "default",
			Name:     "test",
			Ea:       ea1,
		}
		fakeRefReturn := fmt.Sprintf("record:mx/ZG5zLmJpbmRfY25h:%s/%20%20", recSRV.Name)
		aniFakeConnector := &fakeConnector{
			createObjectObj: NewRecordSRV(recSRV),
			getObjectRef:    fakeRefReturn,
			getObjectObj: NewRecordSRV(RecordSRV{
				Port:     443,
				Priority: 1,
				Target:   "xmpp-server.example.com",
				Weight:   10,
				View:     "default",
				Name:     "test",
				Ref:      fakeRefReturn,
				Ea:       ea1,
			}),
			resultObject: NewRecordSRV(RecordSRV{
				Port:     443,
				Priority: 1,
				Target:   "xmpp-server.example.com",
				Weight:   10,
				View:     "default",
				Name:     "test",
				Ea:       ea1,
				Ref:      fakeRefReturn,
			}),
			fakeRefReturn: fakeRefReturn,
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		ea := objMgr.getBasicEA(true)
		aniFakeConnector.createObjectObj.(*RecordSRV).Ea = ea
		aniFakeConnector.resultObject.(*RecordSRV).Ea = ea
		aniFakeConnector.resultObject.(*RecordSRV).Ea["VM ID"] = vmID
		aniFakeConnector.resultObject.(*RecordSRV).Ea["VM Name"] = vmName
		aniFakeConnector.getObjectObj.(*RecordSRV).Ea = ea

		var actualRecord *RecordSRV
		var err error
		It("should pass expected SRV record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateSRVRecord(recSRV)
		})

		It("should return expected SRV record Object", func() {
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get SRV Record by Target", func() {
		recSRV := RecordSRV{Target: "xmpp-server.example.com"}
		name := "test"
		fakeRefReturn := fmt.Sprintf("record:srv/ZG5zLmJpbmRfY25h:%s/%20%20", name)
		aniFakeConnector := &fakeConnector{
			getObjectObj: NewRecordSRV(recSRV),
			getObjectRef: "",
			resultObject: []RecordSRV{*NewRecordSRV(RecordSRV{Name: recSRV.Name, Ref: fakeRefReturn, Target: recSRV.Target})},
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualRecord *[]RecordSRV
		var err error
		It("should pass expected SRV record Object to GetObject", func() {
			actualRecord, err = objMgr.GetSRVRecord(recSRV)

		})

		It("should return expected SRV record Object", func() {
			Expect(*actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get SRV Record by Name", func() {
		recSRV := RecordSRV{Name: "test"}
		fakeRefReturn := fmt.Sprintf("record:srv/ZG5zLmJpbmRfY25h:%s/%20%20", recSRV.Name)
		aniFakeConnector := &fakeConnector{
			getObjectObj: NewRecordSRV(recSRV),
			getObjectRef: "",
			resultObject: []RecordSRV{*NewRecordSRV(RecordSRV{Name: recSRV.Name, Ref: fakeRefReturn})},
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualRecord *[]RecordSRV
		var err error
		It("should pass expected SRV record Object to GetObject", func() {
			actualRecord, err = objMgr.GetSRVRecord(recSRV)

		})

		It("should return expected SRV record Object", func() {
			Expect(*actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get SRV Record by Reference", func() {
		name := "test"
		fakeRefReturn := fmt.Sprintf("record:srv/ZG5zLmJpbmRfY25h:%s/%20%20", name)
		recSRV := RecordSRV{Ref: fakeRefReturn}
		aniFakeConnector := &fakeConnector{
			getObjectObj: NewRecordSRV(recSRV),
			getObjectRef: fakeRefReturn,
			resultObject: []RecordSRV{*NewRecordSRV(RecordSRV{Name: recSRV.Name, Ref: fakeRefReturn})},
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualRecord *[]RecordSRV
		var err error
		It("should pass expected SRV record Object to GetObject", func() {
			actualRecord, err = objMgr.GetSRVRecord(recSRV)

		})

		It("should return expected SRV record Object", func() {
			Expect(*actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete SRV Record by Reference", func() {

		name := "delete_test"
		fakeRefReturn := fmt.Sprintf("record:srv/ZG5zLmJpbmRfY25h:%s/%20%20", name)
		recSRV := RecordSRV{Ref: fakeRefReturn}
		aniFakeConnector := &fakeConnector{
			deleteObjectRef: recSRV.Ref,
			fakeRefReturn:   fakeRefReturn,
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)
		var actualRef string
		var err error
		It("should pass expected SRV record Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteSRVRecord(recSRV)
		})
		It("should return expected SRV record Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete SRV Record by Name", func() {
		recSRV := RecordSRV{Name: "delete_test"}
		fakeRefReturn := fmt.Sprintf("record:srv/ZG5zLmJpbmRfY25h:%s/%20%20", recSRV.Name)
		aniFakeConnector := &fakeConnector{
			getObjectObj:    NewRecordSRV(recSRV),
			getObjectRef:    "",
			resultObject:    []RecordSRV{*NewRecordSRV(RecordSRV{Name: recSRV.Name, Ref: fakeRefReturn})},
			deleteObjectRef: fakeRefReturn,
			fakeRefReturn:   fakeRefReturn,
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)
		var actualRef string
		var err error
		It("should pass expected SRV record Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteSRVRecord(recSRV)
		})
		It("should return expected SRV record Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete SRV Record by Target", func() {
		recSRV := RecordSRV{Target: "xmpp-server.example.com"}
		name := "delete_test"
		fakeRefReturn := fmt.Sprintf("record:srv/ZG5zLmJpbmRfY25h:%s/%20%20", name)
		aniFakeConnector := &fakeConnector{
			getObjectObj:    NewRecordSRV(recSRV),
			getObjectRef:    "",
			resultObject:    []RecordSRV{*NewRecordSRV(RecordSRV{Name: name, Ref: fakeRefReturn, Target: recSRV.Target})},
			deleteObjectRef: fakeRefReturn,
			fakeRefReturn:   fakeRefReturn,
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)
		var actualRef string
		var err error
		It("should pass expected SRV record Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteSRVRecord(recSRV)
		})
		It("should return expected SRV record Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})
})
