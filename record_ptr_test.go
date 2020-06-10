package ibclient

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Testing RecordPTROperations", func() {
	Context("RecordPTR object", func() {
		ipv4addr := "1.1.1.1"
		ptrdname := "bind_a.domain.com"
		view := "default"
		zone := "domain.com"

		rptr := NewRecordPTR(RecordPTR{
			Ipv4Addr: ipv4addr,
			PtrdName: ptrdname,
			View:     view,
			Zone:     zone})

		It("should set fields correctly", func() {
			Expect(rptr.Ipv4Addr).To(Equal(ipv4addr))
			Expect(rptr.PtrdName).To(Equal(ptrdname))
			Expect(rptr.View).To(Equal(view))
			Expect(rptr.Zone).To(Equal(zone))
		})

		It("should set base fields correctly", func() {
			Expect(rptr.ObjectType()).To(Equal("record:ptr"))
			Expect(rptr.ReturnFields()).To(ConsistOf("ipv4addr", "name", "view", "zone", "extattrs", "comment", "creation_time",
				"creator", "ddns_protected", "dns_name", "forbid_reclamation", "reclaimable", "ttl", "use_ttl"))
		})
	})

	Describe("Allocate specific PTR Record ", func() {
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		ea1 := EA{"VM ID": vmID, "VM Name": vmName}

		recPTR := RecordPTR{NetView: "private",
			Cidr:     "53.0.0.0/24",
			Ipv4Addr: "53.0.0.1",
			View:     "default",
			Name:     "test",
			Ea:       ea1,
		}
		fakeRefReturn := fmt.Sprintf("record:ptr/ZG5zLmJpbmRfY25h:%s/%20%20", recPTR.Name)
		aniFakeConnector := &fakeConnector{
			createObjectObj: NewRecordPTR(recPTR),
			getObjectRef:    fakeRefReturn,
			getObjectObj: NewRecordPTR(RecordPTR{
				Name:     recPTR.Name,
				View:     recPTR.View,
				Ipv4Addr: recPTR.Ipv4Addr,
				Ref:      fakeRefReturn,
				Ea:       ea1,
			}),
			resultObject: NewRecordPTR(RecordPTR{
				NetView:  "private",
				Cidr:     "53.0.0.0/24",
				Ipv4Addr: "53.0.0.1",
				View:     "default",
				Name:     "test",
				Ref:      fakeRefReturn,
				Ea:       ea1,
			}),
			fakeRefReturn: fakeRefReturn,
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		ea := objMgr.getBasicEA(true)
		aniFakeConnector.createObjectObj.(*RecordPTR).Ea = ea
		aniFakeConnector.resultObject.(*RecordPTR).Ea = ea
		aniFakeConnector.resultObject.(*RecordPTR).Ea["VM ID"] = vmID
		aniFakeConnector.resultObject.(*RecordPTR).Ea["VM Name"] = vmName
		aniFakeConnector.getObjectObj.(*RecordPTR).Ea = ea

		var actualRecord *RecordPTR
		var err error
		It("should pass expected PTR record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreatePTRRecord(recPTR)
		})

		It("should return expected PTR record Object", func() {
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Allocate A Record in default network view when Cidr and network view is not passed", func() {
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		ea1 := EA{"VM ID": vmID, "VM Name": vmName}

		recPTR := RecordPTR{
			View:     "default",
			Name:     "test",
			Ipv4Addr: "1.1.1.1",
			Ea:       ea1,
		}
		fakeRefReturn := fmt.Sprintf("record:ptr/ZG5zLmJpbmRfY25h:%s/%20%20", recPTR.Name)
		aniFakeConnector := &fakeConnector{
			createObjectObj: NewRecordPTR(recPTR),
			getObjectRef:    fakeRefReturn,
			getObjectObj: NewRecordPTR(RecordPTR{
				Name:     recPTR.Name,
				View:     recPTR.View,
				Ipv4Addr: "1.1.1.1",
				Ref:      fakeRefReturn,
				Ea:       ea1,
			}),
			resultObject: NewRecordPTR(RecordPTR{
				Ipv4Addr: "1.1.1.1",
				View:     "default",
				Name:     "test",
				Ref:      fakeRefReturn,
				Ea:       ea1,
			}),
			fakeRefReturn: fakeRefReturn,
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		ea := objMgr.getBasicEA(true)
		aniFakeConnector.createObjectObj.(*RecordPTR).Ea = ea
		aniFakeConnector.resultObject.(*RecordPTR).Ea = ea
		aniFakeConnector.resultObject.(*RecordPTR).Ea["VM ID"] = vmID
		aniFakeConnector.resultObject.(*RecordPTR).Ea["VM Name"] = vmName
		aniFakeConnector.getObjectObj.(*RecordPTR).Ea = ea

		var actualRecord *RecordPTR
		var err error
		It("should pass expected PTR record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreatePTRRecord(recPTR)
		})

		It("should return expected PTR record Object", func() {
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Allocate next available PTR Record ", func() {
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		ea1 := EA{"VM ID": vmID, "VM Name": vmName}

		recPTR := RecordPTR{NetView: "private",
			Cidr: "53.0.0.0/24",
			View: "default",
			Name: "test",
			Ea:   ea1,
		}
		recPTR.Ipv4Addr = fmt.Sprintf("func:nextavailableip:%s,%s", recPTR.Cidr, recPTR.NetView)
		fakeRefReturn := fmt.Sprintf("record:ptr/ZG5zLmJpbmRfY25h:%s/%20%20", recPTR.Name)
		aniFakeConnector := &fakeConnector{
			createObjectObj: NewRecordPTR(recPTR),
			getObjectRef:    fakeRefReturn,
			getObjectObj: NewRecordPTR(RecordPTR{
				Name:     recPTR.Name,
				View:     recPTR.View,
				Ipv4Addr: fmt.Sprintf("func:nextavailableip:%s,%s", recPTR.Cidr, recPTR.NetView),
				Ref:      fakeRefReturn,
				Ea:       ea1,
			}),
			resultObject: NewRecordPTR(RecordPTR{
				NetView:  "private",
				Cidr:     "53.0.0.0/24",
				Ipv4Addr: fmt.Sprintf("func:nextavailableip:%s,%s", recPTR.Cidr, recPTR.NetView),
				View:     "default",
				Name:     "test",
				Ref:      fakeRefReturn,
				Ea:       ea1,
			}),
			fakeRefReturn: fakeRefReturn,
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		ea := objMgr.getBasicEA(true)
		aniFakeConnector.createObjectObj.(*RecordPTR).Ea = ea
		aniFakeConnector.resultObject.(*RecordPTR).Ea = ea
		aniFakeConnector.resultObject.(*RecordPTR).Ea["VM ID"] = vmID
		aniFakeConnector.resultObject.(*RecordPTR).Ea["VM Name"] = vmName
		aniFakeConnector.getObjectObj.(*RecordPTR).Ea = ea

		var actualRecord *RecordPTR
		var err error
		It("should pass expected PTR record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreatePTRRecord(recPTR)
		})

		It("should return expected PTR record Object", func() {
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get PTR Record by IPv4Addr", func() {
		recPTR := RecordPTR{Ipv4Addr: "1.1.1.1"}
		name := "test"
		fakeRefReturn := fmt.Sprintf("record:ptr/ZG5zLmJpbmRfY25h:%s/%20%20", name)
		aniFakeConnector := &fakeConnector{
			getObjectObj: NewRecordPTR(recPTR),
			getObjectRef: "",
			resultObject: []RecordPTR{*NewRecordPTR(RecordPTR{Name: recPTR.Name, Ref: fakeRefReturn, Ipv4Addr: recPTR.Ipv4Addr})},
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualRecord *[]RecordPTR
		var err error
		It("should pass expected PTR record Object to GetObject", func() {
			actualRecord, err = objMgr.GetPTRRecord(recPTR)

		})

		It("should return expected PTR record Object", func() {
			Expect(*actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get PTR Record by PtrdName", func() {
		recPTR := RecordPTR{PtrdName: "test"}
		fakeRefReturn := fmt.Sprintf("record:ptr/ZG5zLmJpbmRfY25h:%s/%20%20", recPTR.PtrdName)
		aniFakeConnector := &fakeConnector{
			getObjectObj: NewRecordPTR(recPTR),
			getObjectRef: "",
			resultObject: []RecordPTR{*NewRecordPTR(RecordPTR{Name: recPTR.Name, Ref: fakeRefReturn})},
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualRecord *[]RecordPTR
		var err error
		It("should pass expected PTR record Object to GetObject", func() {
			actualRecord, err = objMgr.GetPTRRecord(recPTR)

		})

		It("should return expected PTR record Object", func() {
			Expect(*actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get PTR Record by Reference", func() {
		name := "test"
		fakeRefReturn := fmt.Sprintf("record:ptr/ZG5zLmJpbmRfY25h:%s/%20%20", name)
		recPTR := RecordPTR{Ref: fakeRefReturn}
		aniFakeConnector := &fakeConnector{
			getObjectObj: NewRecordPTR(recPTR),
			getObjectRef: fakeRefReturn,
			resultObject: []RecordPTR{*NewRecordPTR(RecordPTR{Name: recPTR.Name, Ref: fakeRefReturn})},
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualRecord *[]RecordPTR
		var err error
		It("should pass expected PTR record Object to GetObject", func() {
			actualRecord, err = objMgr.GetPTRRecord(recPTR)

		})

		It("should return expected PTR record Object", func() {
			Expect(*actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete PTR Record by Reference", func() {

		name := "delete_test"
		fakeRefReturn := fmt.Sprintf("record:ptr/ZG5zLmJpbmRfY25h:%s/%20%20", name)
		recPTR := RecordPTR{Ref: fakeRefReturn}
		aniFakeConnector := &fakeConnector{
			deleteObjectRef: recPTR.Ref,
			fakeRefReturn:   fakeRefReturn,
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)
		var actualRef string
		var err error
		It("should pass expected PTR record Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeletePTRRecord(recPTR)
		})
		It("should return expected PTR record Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete PTR Record by Ptrdame", func() {
		recPTR := RecordPTR{PtrdName: "delete_test"}
		name := "test"
		fakeRefReturn := fmt.Sprintf("record:ptr/ZG5zLmJpbmRfY25h:%s/%20%20", name)
		aniFakeConnector := &fakeConnector{
			getObjectObj:    NewRecordPTR(recPTR),
			getObjectRef:    "",
			resultObject:    []RecordPTR{*NewRecordPTR(RecordPTR{PtrdName: recPTR.PtrdName, Ref: fakeRefReturn})},
			deleteObjectRef: fakeRefReturn,
			fakeRefReturn:   fakeRefReturn,
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)
		var actualRef string
		var err error
		It("should pass expected PTR record Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeletePTRRecord(recPTR)
		})
		It("should return expected PTR record Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete PTR Record by IPv4Addr", func() {
		recPTR := RecordPTR{Ipv4Addr: "1.1.1.1"}
		name := "delete_test"
		fakeRefReturn := fmt.Sprintf("record:ptr/ZG5zLmJpbmRfY25h:%s/%20%20", name)
		aniFakeConnector := &fakeConnector{
			getObjectObj:    NewRecordPTR(recPTR),
			getObjectRef:    "",
			resultObject:    []RecordPTR{*NewRecordPTR(RecordPTR{Name: name, Ref: fakeRefReturn, Ipv4Addr: recPTR.Ipv4Addr})},
			deleteObjectRef: fakeRefReturn,
			fakeRefReturn:   fakeRefReturn,
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)
		var actualRef string
		var err error
		It("should pass expected PTR record Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeletePTRRecord(recPTR)
		})
		It("should return expected PTR record Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})
})
