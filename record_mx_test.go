package ibclient

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Testing RecordMXOperations", func() {
	Context("RecordMX object", func() {
		mailExchanger := "test.example.com"
		preference := uint32(10)
		name := "bind_a.domain.com"
		view := "default"
		zone := "domain.com"

		ra := NewRecordMX(RecordMX{
			MailExchanger: mailExchanger,
			Preference:    preference,
			Name:          name,
			View:          view,
			Zone:          zone})

		It("should set fields correctly", func() {
			Expect(ra.MailExchanger).To(Equal(mailExchanger))
			Expect(ra.Preference).To(Equal(preference))
			Expect(ra.Name).To(Equal(name))
			Expect(ra.View).To(Equal(view))
			Expect(ra.Zone).To(Equal(zone))
		})

		It("should set base fields correctly", func() {
			Expect(ra.ObjectType()).To(Equal("record:mx"))
			Expect(ra.ReturnFields()).To(ConsistOf("mail_exchanger", "preference", "name", "view", "zone", "extattrs", "comment",
				"creation_time", "ddns_protected", "dns_name", "forbid_reclamation", "reclaimable", "ttl", "use_ttl", "disable"))
		})
	})

	Describe("Allocate specific MX Record ", func() {
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		ea1 := EA{"VM ID": vmID, "VM Name": vmName}

		recMX := RecordMX{
			MailExchanger: "test.example.com",
			Preference:    10,
			View:          "default",
			Name:          "test",
			Ea:            ea1,
		}
		fakeRefReturn := fmt.Sprintf("record:mx/ZG5zLmJpbmRfY25h:%s/%20%20", recMX.Name)
		aniFakeConnector := &fakeConnector{
			createObjectObj: NewRecordMX(recMX),
			getObjectRef:    fakeRefReturn,
			getObjectObj: NewRecordMX(RecordMX{
				Name:          recMX.Name,
				View:          recMX.View,
				Preference:    recMX.Preference,
				MailExchanger: recMX.MailExchanger,
				Ref:           fakeRefReturn,
				Ea:            ea1,
			}),
			resultObject: NewRecordMX(RecordMX{
				MailExchanger: "test.example.com",
				Preference:    10,
				View:          "default",
				Name:          "test",
				Ea:            ea1,
				Ref:           fakeRefReturn,
			}),
			fakeRefReturn: fakeRefReturn,
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		ea := objMgr.getBasicEA(true)
		aniFakeConnector.createObjectObj.(*RecordMX).Ea = ea
		aniFakeConnector.resultObject.(*RecordMX).Ea = ea
		aniFakeConnector.resultObject.(*RecordMX).Ea["VM ID"] = vmID
		aniFakeConnector.resultObject.(*RecordMX).Ea["VM Name"] = vmName
		aniFakeConnector.getObjectObj.(*RecordMX).Ea = ea

		var actualRecord *RecordMX
		var err error
		It("should pass expected MX record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateMXRecord(recMX)
		})

		It("should return expected MX record Object", func() {
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get MX Record by MailExchanger", func() {
		recMX := RecordMX{MailExchanger: "test.example.com"}
		name := "test"
		fakeRefReturn := fmt.Sprintf("record:mx/ZG5zLmJpbmRfY25h:%s/%20%20", name)
		aniFakeConnector := &fakeConnector{
			getObjectObj: NewRecordMX(recMX),
			getObjectRef: "",
			resultObject: []RecordMX{*NewRecordMX(RecordMX{Name: recMX.Name, Ref: fakeRefReturn, MailExchanger: recMX.MailExchanger})},
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualRecord *[]RecordMX
		var err error
		It("should pass expected MX record Object to GetObject", func() {
			actualRecord, err = objMgr.GetMXRecord(recMX)

		})

		It("should return expected MX record Object", func() {
			Expect(*actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get MX Record by Name", func() {
		recMX := RecordMX{Name: "test"}
		fakeRefReturn := fmt.Sprintf("record:mx/ZG5zLmJpbmRfY25h:%s/%20%20", recMX.Name)
		aniFakeConnector := &fakeConnector{
			getObjectObj: NewRecordMX(recMX),
			getObjectRef: "",
			resultObject: []RecordMX{*NewRecordMX(RecordMX{Name: recMX.Name, Ref: fakeRefReturn})},
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualRecord *[]RecordMX
		var err error
		It("should pass expected MX record Object to GetObject", func() {
			actualRecord, err = objMgr.GetMXRecord(recMX)

		})

		It("should return expected MX record Object", func() {
			Expect(*actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get MX Record by Preference", func() {
		recMX := RecordMX{Preference: 10}
		name := "test"
		fakeRefReturn := fmt.Sprintf("record:mx/ZG5zLmJpbmRfY25h:%s/%20%20", name)
		aniFakeConnector := &fakeConnector{
			getObjectObj: NewRecordMX(recMX),
			getObjectRef: "",
			resultObject: []RecordMX{*NewRecordMX(RecordMX{Name: recMX.Name, Ref: fakeRefReturn, Preference: recMX.Preference})},
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualRecord *[]RecordMX
		var err error
		It("should pass expected MX record Object to GetObject", func() {
			actualRecord, err = objMgr.GetMXRecord(recMX)

		})

		It("should return expected MX record Object", func() {
			Expect(*actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get MX Record by Reference", func() {
		name := "test"
		fakeRefReturn := fmt.Sprintf("record:mx/ZG5zLmJpbmRfY25h:%s/%20%20", name)
		recMX := RecordMX{Ref: fakeRefReturn}
		aniFakeConnector := &fakeConnector{
			getObjectObj: NewRecordMX(recMX),
			getObjectRef: fakeRefReturn,
			resultObject: []RecordMX{*NewRecordMX(RecordMX{Name: recMX.Name, Ref: fakeRefReturn})},
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualRecord *[]RecordMX
		var err error
		It("should pass expected MX record Object to GetObject", func() {
			actualRecord, err = objMgr.GetMXRecord(recMX)

		})

		It("should return expected MX record Object", func() {
			Expect(*actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete MX Record by Reference", func() {

		name := "delete_test"
		fakeRefReturn := fmt.Sprintf("record:mx/ZG5zLmJpbmRfY25h:%s/%20%20", name)
		recMX := RecordMX{Ref: fakeRefReturn}
		aniFakeConnector := &fakeConnector{
			deleteObjectRef: recMX.Ref,
			fakeRefReturn:   fakeRefReturn,
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)
		var actualRef string
		var err error
		It("should pass expected MX record Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteMXRecord(recMX)
		})
		It("should return expected MX record Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete MX Record by Name", func() {
		recMX := RecordMX{Name: "delete_test"}
		fakeRefReturn := fmt.Sprintf("record:mx/ZG5zLmJpbmRfY25h:%s/%20%20", recMX.Name)
		aniFakeConnector := &fakeConnector{
			getObjectObj:    NewRecordMX(recMX),
			getObjectRef:    "",
			resultObject:    []RecordMX{*NewRecordMX(RecordMX{Name: recMX.Name, Ref: fakeRefReturn})},
			deleteObjectRef: fakeRefReturn,
			fakeRefReturn:   fakeRefReturn,
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)
		var actualRef string
		var err error
		It("should pass expected Mx record Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteMXRecord(recMX)
		})
		It("should return expected MX record Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete MX Record by MailExchanger", func() {
		recMX := RecordMX{MailExchanger: "test.example.com"}
		name := "delete_test"
		fakeRefReturn := fmt.Sprintf("record:mx/ZG5zLmJpbmRfY25h:%s/%20%20", name)
		aniFakeConnector := &fakeConnector{
			getObjectObj:    NewRecordMX(recMX),
			getObjectRef:    "",
			resultObject:    []RecordMX{*NewRecordMX(RecordMX{Name: name, Ref: fakeRefReturn, MailExchanger: recMX.MailExchanger})},
			deleteObjectRef: fakeRefReturn,
			fakeRefReturn:   fakeRefReturn,
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)
		var actualRef string
		var err error
		It("should pass expected MX record Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteMXRecord(recMX)
		})
		It("should return expected MX record Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})
})
