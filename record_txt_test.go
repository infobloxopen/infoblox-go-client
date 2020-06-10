package ibclient

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Testing RecordTXTOperations", func() {
	Context("RecordTXT object", func() {
		name := "txt.domain.com"
		text := "this is text string"
		view := "default"
		zone := "domain.com"

		rt := NewRecordTXT(RecordTXT{
			Name: name,
			Text: text,
			View: view,
			Zone: zone})

		It("should set fields correctly", func() {
			Expect(rt.Name).To(Equal(name))
			Expect(rt.Text).To(Equal(text))
			Expect(rt.View).To(Equal(view))
			Expect(rt.Zone).To(Equal(zone))
		})

		It("should set base fields correctly", func() {
			Expect(rt.ObjectType()).To(Equal("record:txt"))
			Expect(rt.ReturnFields()).To(ConsistOf("name", "text", "view", "zone", "extattrs", "comment", "creation_time",
				"creator", "ddns_protected", "dns_name", "forbid_reclamation", "reclaimable", "ttl", "use_ttl"))
		})
	})

	Describe("Allocate specific TXT Record ", func() {
		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		ea := EA{"VM ID": vmID, "VM Name": vmName}
		recTXT := RecordTXT{
			Name: "test",
			Text: "test-text",
			TTL:  30,
			View: "default",
			Ea: ea,
		}
		fakeRefReturn := fmt.Sprintf("record:txt/ZG5zLmJpbmRfY25h:%s/%20%20", recTXT.Name)
		aniFakeConnector := &fakeConnector{
			createObjectObj: NewRecordTXT(recTXT),
			getObjectRef:    fakeRefReturn,
			getObjectObj: NewRecordTXT(RecordTXT{
				Name: recTXT.Name,
				View: recTXT.View,
				Text: "test-text",
				TTL:  30,
			}),
			resultObject: NewRecordTXT(RecordTXT{
				Ref:  fakeRefReturn,
				Name: recTXT.Name,
				View: recTXT.View,
				Text: "test-text",
				TTL:  30,
			}),
			fakeRefReturn: fakeRefReturn,
		}

		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)
		ea = objMgr.getBasicEA(true)
		aniFakeConnector.createObjectObj.(*RecordTXT).Ea = ea
		aniFakeConnector.resultObject.(*RecordTXT).Ea = ea
		aniFakeConnector.resultObject.(*RecordTXT).Ea["VM ID"] = vmID
		aniFakeConnector.resultObject.(*RecordTXT).Ea["VM Name"] = vmName
		aniFakeConnector.getObjectObj.(*RecordTXT).Ea = ea
		var actualRecord *RecordTXT
		var err error
		It("should pass expected TXT record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateTXTRecord(recTXT)
		})

		It("should return expected TXT record Object", func() {
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get TXT Record by Text", func() {
		recTXT := RecordTXT{Text: "test-text"}
		name := "test"
		fakeRefReturn := fmt.Sprintf("record:txt/ZG5zLmJpbmRfY25h:%s/%20%20", name)
		aniFakeConnector := &fakeConnector{
			getObjectObj: NewRecordTXT(recTXT),
			getObjectRef: "",
			resultObject: []RecordTXT{*NewRecordTXT(RecordTXT{Name: recTXT.Name, Ref: fakeRefReturn, Text: recTXT.Text})},
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualRecord *[]RecordTXT
		var err error
		It("should pass expected TXT record Object to GetObject", func() {
			actualRecord, err = objMgr.GetTXTRecord(recTXT)

		})

		It("should return expected TXT record Object", func() {
			Expect(*actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get TXT Record by Name", func() {
		recTXT := RecordTXT{Name: "test"}
		fakeRefReturn := fmt.Sprintf("record:txt/ZG5zLmJpbmRfY25h:%s/%20%20", recTXT.Name)
		aniFakeConnector := &fakeConnector{
			getObjectObj: NewRecordTXT(recTXT),
			getObjectRef: "",
			resultObject: []RecordTXT{*NewRecordTXT(RecordTXT{Name: recTXT.Name, Ref: fakeRefReturn})},
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualRecord *[]RecordTXT
		var err error
		It("should pass expected TXT record Object to GetObject", func() {
			actualRecord, err = objMgr.GetTXTRecord(recTXT)

		})

		It("should return expected TXT record Object", func() {
			Expect(*actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get TXT Record by Reference", func() {
		name := "test"
		fakeRefReturn := fmt.Sprintf("record:txt/ZG5zLmJpbmRfY25h:%s/%20%20", name)
		recTXT := RecordTXT{Ref: fakeRefReturn}
		aniFakeConnector := &fakeConnector{
			getObjectObj: NewRecordTXT(recTXT),
			getObjectRef: fakeRefReturn,
			resultObject: []RecordTXT{*NewRecordTXT(RecordTXT{Name: recTXT.Name, Ref: fakeRefReturn})},
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualRecord *[]RecordTXT
		var err error
		It("should pass expected TXT record Object to GetObject", func() {
			actualRecord, err = objMgr.GetTXTRecord(recTXT)

		})

		It("should return expected TXT record Object", func() {
			Expect(*actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete TXT Record by Reference", func() {

		name := "delete_test"
		fakeRefReturn := fmt.Sprintf("record:txt/ZG5zLmJpbmRfY25h:%s/%20%20", name)
		recTXT := RecordTXT{Ref: fakeRefReturn}
		aniFakeConnector := &fakeConnector{
			deleteObjectRef: recTXT.Ref,
			fakeRefReturn:   fakeRefReturn,
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)
		var actualRef string
		var err error
		It("should pass expected TXT record Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteTXTRecord(recTXT)
		})
		It("should return expected TXT record Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete TXT Record by Name", func() {
		recTXT := RecordTXT{Name: "delete_test"}
		fakeRefReturn := fmt.Sprintf("record:txt/ZG5zLmJpbmRfY25h:%s/%20%20", recTXT.Name)
		aniFakeConnector := &fakeConnector{
			getObjectObj:    NewRecordTXT(recTXT),
			getObjectRef:    "",
			resultObject:    []RecordTXT{*NewRecordTXT(RecordTXT{Name: recTXT.Name, Ref: fakeRefReturn})},
			deleteObjectRef: fakeRefReturn,
			fakeRefReturn:   fakeRefReturn,
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)
		var actualRef string
		var err error
		It("should pass expected TXT record Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteTXTRecord(recTXT)
		})
		It("should return expected TXT record Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete TXT Record by Text", func() {
		recTXT := RecordTXT{Text: "test-text"}
		name := "delete_test"
		fakeRefReturn := fmt.Sprintf("record:txt/ZG5zLmJpbmRfY25h:%s/%20%20", name)
		aniFakeConnector := &fakeConnector{
			getObjectObj:    NewRecordTXT(recTXT),
			getObjectRef:    "",
			resultObject:    []RecordTXT{*NewRecordTXT(RecordTXT{Name: name, Ref: fakeRefReturn, Text: recTXT.Text})},
			deleteObjectRef: fakeRefReturn,
			fakeRefReturn:   fakeRefReturn,
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)
		var actualRef string
		var err error
		It("should pass expected TXT record Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteTXTRecord(recTXT)
		})
		It("should return expected TXT record Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})
})
