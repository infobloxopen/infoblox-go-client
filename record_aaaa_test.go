package ibclient

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Testing RecordAAAAOperations", func() {
	Context("RecordAAAA object", func() {
		ipv6addr := "fd60:e32:f1b9::4"
		name := "bind_a.domain.com"
		view := "default"
		zone := "domain.com"

		ra := NewRecordAAAA(RecordAAAA{
			Ipv6Addr: ipv6addr,
			Name:     name,
			View:     view,
			Zone:     zone})

		It("should set fields correctly", func() {
			Expect(ra.Ipv6Addr).To(Equal(ipv6addr))
			Expect(ra.Name).To(Equal(name))
			Expect(ra.View).To(Equal(view))
			Expect(ra.Zone).To(Equal(zone))
		})

		It("should set base fields correctly", func() {
			Expect(ra.ObjectType()).To(Equal("record:aaaa"))
			Expect(ra.ReturnFields()).To(ConsistOf("ipv6addr", "name", "view", "zone", "extattrs", "comment", "creation_time",
				"creator", "ddns_protected", "dns_name", "forbid_reclamation", "reclaimable", "ttl", "use_ttl"))
		})
	})

	Describe("Allocate specific AAAA Record ", func() {
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		ea1 := EA{"VM ID": vmID, "VM Name": vmName}

		recA4 := RecordAAAA{NetView: "private",
			Cidr:     "fd60:e32:f1b9::/32",
			Ipv6Addr: "fd60:e32:f1b9::4",
			View:     "default",
			Name:     "test",
			Ea:       ea1,
		}
		fakeRefReturn := fmt.Sprintf("record:aaaa/ZG5zLmJpbmRfY25h:%s/%20%20", recA4.Name)
		aniFakeConnector := &fakeConnector{
			createObjectObj: NewRecordAAAA(recA4),
			getObjectRef:    fakeRefReturn,
			getObjectObj: NewRecordAAAA(RecordAAAA{
				Name:     recA4.Name,
				View:     recA4.View,
				Ipv6Addr: recA4.Ipv6Addr,
				Ref:      fakeRefReturn,
				Ea:       ea1,
			}),
			resultObject: NewRecordAAAA(RecordAAAA{
				NetView:  "private",
				Cidr:     "fd60:e32:f1b9::/32",
				Ipv6Addr: "fd60:e32:f1b9::4",
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
		aniFakeConnector.createObjectObj.(*RecordAAAA).Ea = ea
		aniFakeConnector.resultObject.(*RecordAAAA).Ea = ea
		aniFakeConnector.resultObject.(*RecordAAAA).Ea["VM ID"] = vmID
		aniFakeConnector.resultObject.(*RecordAAAA).Ea["VM Name"] = vmName
		aniFakeConnector.getObjectObj.(*RecordAAAA).Ea = ea

		var actualRecord *RecordAAAA
		var err error
		It("should pass expected AAAA record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateAAAARecord(recA4)
		})

		It("should return expected AAAA record Object", func() {
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Allocate AAAA Record in default network view when Cidr and network view is not passed", func() {
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		ea1 := EA{"VM ID": vmID, "VM Name": vmName}

		recA4 := RecordAAAA{
			View:     "default",
			Name:     "test",
			Ipv6Addr: "fd60:e32:f1b9::4",
			Ea:       ea1,
		}
		fakeRefReturn := fmt.Sprintf("record:aaaa/ZG5zLmJpbmRfY25h:%s/%20%20", recA4.Name)
		aniFakeConnector := &fakeConnector{
			createObjectObj: NewRecordAAAA(recA4),
			getObjectRef:    fakeRefReturn,
			getObjectObj: NewRecordAAAA(RecordAAAA{
				Name:     recA4.Name,
				View:     recA4.View,
				Ipv6Addr: "fd60:e32:f1b9::4",
				Ref:      fakeRefReturn,
				Ea:       ea1,
			}),
			resultObject: NewRecordAAAA(RecordAAAA{
				Ipv6Addr: "fd60:e32:f1b9::4",
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
		aniFakeConnector.createObjectObj.(*RecordAAAA).Ea = ea
		aniFakeConnector.resultObject.(*RecordAAAA).Ea = ea
		aniFakeConnector.resultObject.(*RecordAAAA).Ea["VM ID"] = vmID
		aniFakeConnector.resultObject.(*RecordAAAA).Ea["VM Name"] = vmName
		aniFakeConnector.getObjectObj.(*RecordAAAA).Ea = ea

		var actualRecord *RecordAAAA
		var err error
		It("should pass expected AAAA record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateAAAARecord(recA4)
		})

		It("should return expected AAAA record Object", func() {
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Allocate next available AAAA Record ", func() {
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		ea1 := EA{"VM ID": vmID, "VM Name": vmName}

		recA4 := RecordAAAA{NetView: "private",
			Cidr: "fd60:e32:f1b9::/32",
			View: "default",
			Name: "test",
			Ea:   ea1,
		}
		recA4.Ipv6Addr = fmt.Sprintf("func:nextavailableip:%s,%s", recA4.Cidr, recA4.NetView)
		fakeRefReturn := fmt.Sprintf("record:aaaa/ZG5zLmJpbmRfY25h:%s/%20%20", recA4.Name)
		aniFakeConnector := &fakeConnector{
			createObjectObj: NewRecordAAAA(recA4),
			getObjectRef:    fakeRefReturn,
			getObjectObj: NewRecordAAAA(RecordAAAA{
				Name:     recA4.Name,
				View:     recA4.View,
				Ipv6Addr: fmt.Sprintf("func:nextavailableip:%s,%s", recA4.Cidr, recA4.NetView),
				Ref:      fakeRefReturn,
				Ea:       ea1,
			}),
			resultObject: NewRecordAAAA(RecordAAAA{
				NetView:  "private",
				Cidr:     "fd60:e32:f1b9::/32",
				Ipv6Addr: fmt.Sprintf("func:nextavailableip:%s,%s", recA4.Cidr, recA4.NetView),
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
		aniFakeConnector.createObjectObj.(*RecordAAAA).Ea = ea
		aniFakeConnector.resultObject.(*RecordAAAA).Ea = ea
		aniFakeConnector.resultObject.(*RecordAAAA).Ea["VM ID"] = vmID
		aniFakeConnector.resultObject.(*RecordAAAA).Ea["VM Name"] = vmName
		aniFakeConnector.getObjectObj.(*RecordAAAA).Ea = ea

		var actualRecord *RecordAAAA
		var err error
		It("should pass expected AAAA record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateAAAARecord(recA4)
		})

		It("should return expected AAAA record Object", func() {
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get AAAA Record by IPv6Addr", func() {
		recA4 := RecordAAAA{Ipv6Addr: "fd60:e32:f1b9::4"}
		name := "test"
		fakeRefReturn := fmt.Sprintf("record:aaaa/ZG5zLmJpbmRfY25h:%s/%20%20", name)
		aniFakeConnector := &fakeConnector{
			getObjectObj: NewRecordAAAA(recA4),
			getObjectRef: "",
			resultObject: []RecordAAAA{*NewRecordAAAA(RecordAAAA{Name: recA4.Name, Ref: fakeRefReturn, Ipv6Addr: recA4.Ipv6Addr})},
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualRecord *[]RecordAAAA
		var err error
		It("should pass expected AAAA record Object to GetObject", func() {
			actualRecord, err = objMgr.GetAAAARecord(recA4)

		})

		It("should return expected AAAA record Object", func() {
			Expect(*actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get AAAA Record by Name", func() {
		recA4 := RecordAAAA{Name: "test"}
		fakeRefReturn := fmt.Sprintf("record:aaaa/ZG5zLmJpbmRfY25h:%s/%20%20", recA4.Name)
		aniFakeConnector := &fakeConnector{
			getObjectObj: NewRecordAAAA(recA4),
			getObjectRef: "",
			resultObject: []RecordAAAA{*NewRecordAAAA(RecordAAAA{Name: recA4.Name, Ref: fakeRefReturn})},
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualRecord *[]RecordAAAA
		var err error
		It("should pass expected AAAA record Object to GetObject", func() {
			actualRecord, err = objMgr.GetAAAARecord(recA4)

		})

		It("should return expected AAAA record Object", func() {
			Expect(*actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get AAAA Record by Reference", func() {
		name := "test"
		fakeRefReturn := fmt.Sprintf("record:aaaa/ZG5zLmJpbmRfY25h:%s/%20%20", name)
		recA4 := RecordAAAA{Ref: fakeRefReturn}
		aniFakeConnector := &fakeConnector{
			getObjectObj: NewRecordAAAA(recA4),
			getObjectRef: fakeRefReturn,
			resultObject: []RecordAAAA{*NewRecordAAAA(RecordAAAA{Name: recA4.Name, Ref: fakeRefReturn})},
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualRecord *[]RecordAAAA
		var err error
		It("should pass expected AAAA record Object to GetObject", func() {
			actualRecord, err = objMgr.GetAAAARecord(recA4)

		})

		It("should return expected AAAA record Object", func() {
			Expect(*actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete AAAA Record by Reference", func() {

		name := "delete_test"
		fakeRefReturn := fmt.Sprintf("record:aaaa/ZG5zLmJpbmRfY25h:%s/%20%20", name)
		recA4 := RecordAAAA{Ref: fakeRefReturn}
		aniFakeConnector := &fakeConnector{
			deleteObjectRef: recA4.Ref,
			fakeRefReturn:   fakeRefReturn,
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)
		var actualRef string
		var err error
		It("should pass expected AAAA record Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteAAAARecord(recA4)
		})
		It("should return expected AAAA record Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete AAAA Record by Name", func() {
		recA4 := RecordAAAA{Name: "delete_test"}
		fakeRefReturn := fmt.Sprintf("record:aaaa/ZG5zLmJpbmRfY25h:%s/%20%20", recA4.Name)
		aniFakeConnector := &fakeConnector{
			getObjectObj:    NewRecordAAAA(recA4),
			getObjectRef:    "",
			resultObject:    []RecordAAAA{*NewRecordAAAA(RecordAAAA{Name: recA4.Name, Ref: fakeRefReturn})},
			deleteObjectRef: fakeRefReturn,
			fakeRefReturn:   fakeRefReturn,
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)
		var actualRef string
		var err error
		It("should pass expected AAAA record Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteAAAARecord(recA4)
		})
		It("should return expected AAAA record Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete AAAA Record by IPv6Addr", func() {
		recA4 := RecordAAAA{Ipv6Addr: "fd60:e32:f1b9::4"}
		name := "delete_test"
		fakeRefReturn := fmt.Sprintf("record:aaaa/ZG5zLmJpbmRfY25h:%s/%20%20", name)
		aniFakeConnector := &fakeConnector{
			getObjectObj:    NewRecordAAAA(recA4),
			getObjectRef:    "",
			resultObject:    []RecordAAAA{*NewRecordAAAA(RecordAAAA{Name: name, Ref: fakeRefReturn, Ipv6Addr: recA4.Ipv6Addr})},
			deleteObjectRef: fakeRefReturn,
			fakeRefReturn:   fakeRefReturn,
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)
		var actualRef string
		var err error
		It("should pass expected AAAA record Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteAAAARecord(recA4)
		})
		It("should return expected AAAA record Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})
})
