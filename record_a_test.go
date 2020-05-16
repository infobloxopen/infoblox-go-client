package ibclient

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Testing RecordAOperations", func() {
	Context("RecordA object", func() {
		ipv4addr := "1.1.1.1"
		name := "bind_a.domain.com"
		view := "default"
		zone := "domain.com"

		ra := NewRecordA(RecordA{
			Ipv4Addr: ipv4addr,
			Name:     name,
			View:     view,
			Zone:     zone})

		It("should set fields correctly", func() {
			Expect(ra.Ipv4Addr).To(Equal(ipv4addr))
			Expect(ra.Name).To(Equal(name))
			Expect(ra.View).To(Equal(view))
			Expect(ra.Zone).To(Equal(zone))
		})

		It("should set base fields correctly", func() {
			Expect(ra.ObjectType()).To(Equal("record:a"))
			Expect(ra.ReturnFields()).To(ConsistOf("ipv4addr", "name", "view", "zone","extattrs","comment","creation_time",
				"creator","ddns_protected","dns_name","cloud_info","forbid_reclamation","last_queried",
				"reclaimable","ttl","use_ttl","aws_rte53_record_info","ddns_principal","disable","discovered_data","ms_ad_user_data"))
		})
	})

	Describe("Allocate specific A Record ", func() {
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		ea1 := EA{"VM ID": vmID, "VM Name": vmName}

		recA := RecordA{NetView: "private",
			Cidr:     "53.0.0.0/24",
			Ipv4Addr: "53.0.0.1",
			View:     "default",
			Name:     "test",
			Ea:       ea1,
			}
		fakeRefReturn := fmt.Sprintf("record:a/ZG5zLmJpbmRfY25h:%s/%20%20", recA.Name)
		aniFakeConnector := &fakeConnector{
			createObjectObj: NewRecordA(recA),
			getObjectRef:    fakeRefReturn,
			getObjectObj: NewRecordA(RecordA{
				Name:     recA.Name,
				View:     recA.View,
				Ipv4Addr: recA.Ipv4Addr,
				Ref:      fakeRefReturn,
				Ea:       ea1,
			}),
			resultObject:  NewRecordA(RecordA{
				NetView: "private",
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
		aniFakeConnector.createObjectObj.(*RecordA).Ea = ea
		aniFakeConnector.resultObject.(*RecordA).Ea = ea
		aniFakeConnector.resultObject.(*RecordA).Ea["VM ID"] = vmID
		aniFakeConnector.resultObject.(*RecordA).Ea["VM Name"] = vmName
		aniFakeConnector.getObjectObj.(*RecordA).Ea = ea

		var actualRecord *RecordA
		var err error
		It("should pass expected A record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateARecord(recA)
		})

		It("should return expected A record Object", func() {
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Allocate A Record in default network view when Cidr and network view is not passed", func() {
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		ea1 := EA{"VM ID": vmID, "VM Name": vmName}

		recA := RecordA{
			View:     "default",
			Name:     "test",
			Ipv4Addr: "1.1.1.1",
			Ea:       ea1,
		}
		fakeRefReturn := fmt.Sprintf("record:a/ZG5zLmJpbmRfY25h:%s/%20%20", recA.Name)
		aniFakeConnector := &fakeConnector{
			createObjectObj: NewRecordA(recA),
			getObjectRef:    fakeRefReturn,
			getObjectObj: NewRecordA(RecordA{
				Name:     recA.Name,
				View:     recA.View,
				Ipv4Addr: "1.1.1.1",
				Ref:      fakeRefReturn,
				Ea:       ea1,
			}),
			resultObject:  NewRecordA(RecordA{
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
		aniFakeConnector.createObjectObj.(*RecordA).Ea = ea
		aniFakeConnector.resultObject.(*RecordA).Ea = ea
		aniFakeConnector.resultObject.(*RecordA).Ea["VM ID"] = vmID
		aniFakeConnector.resultObject.(*RecordA).Ea["VM Name"] = vmName
		aniFakeConnector.getObjectObj.(*RecordA).Ea = ea

		var actualRecord *RecordA
		var err error
		It("should pass expected A record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateARecord(recA)
		})

		It("should return expected A record Object", func() {
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Allocate next available A Record ", func() {
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		ea1 := EA{"VM ID": vmID, "VM Name": vmName}

		recA := RecordA{NetView: "private",
			Cidr:     "53.0.0.0/24",
			View:     "default",
			Name:     "test",
			Ea:       ea1,
		}
		recA.Ipv4Addr = fmt.Sprintf("func:nextavailableip:%s,%s", recA.Cidr, recA.NetView)
		fakeRefReturn := fmt.Sprintf("record:a/ZG5zLmJpbmRfY25h:%s/%20%20", recA.Name)
		aniFakeConnector := &fakeConnector{
			createObjectObj: NewRecordA(recA),
			getObjectRef:    fakeRefReturn,
			getObjectObj: NewRecordA(RecordA{
				Name:     recA.Name,
				View:     recA.View,
				Ipv4Addr: fmt.Sprintf("func:nextavailableip:%s,%s", recA.Cidr, recA.NetView),
				Ref:      fakeRefReturn,
				Ea:       ea1,
			}),
			resultObject:  NewRecordA(RecordA{
				NetView: "private",
				Cidr:     "53.0.0.0/24",
				Ipv4Addr: fmt.Sprintf("func:nextavailableip:%s,%s", recA.Cidr, recA.NetView),
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
		aniFakeConnector.createObjectObj.(*RecordA).Ea = ea
		aniFakeConnector.resultObject.(*RecordA).Ea = ea
		aniFakeConnector.resultObject.(*RecordA).Ea["VM ID"] = vmID
		aniFakeConnector.resultObject.(*RecordA).Ea["VM Name"] = vmName
		aniFakeConnector.getObjectObj.(*RecordA).Ea = ea

		var actualRecord *RecordA
		var err error
		It("should pass expected A record Object to CreateObject", func() {
			actualRecord, err = objMgr.CreateARecord(recA)
		})

		It("should return expected A record Object", func() {
			Expect(actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get A Record by IPv4Addr", func() {
		recA := RecordA{Ipv4Addr: "1.1.1.1"}
		name := "test"
		fakeRefReturn := fmt.Sprintf("record:a/ZG5zLmJpbmRfY25h:%s/%20%20", name)
		aniFakeConnector := &fakeConnector{
			getObjectObj: NewRecordA(recA),
			getObjectRef: "",
			resultObject: []RecordA{*NewRecordA(RecordA{Name: recA.Name, Ref: fakeRefReturn, Ipv4Addr: recA.Ipv4Addr})},
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualRecord *[]RecordA
		var err error
		It("should pass expected A record Object to GetObject", func() {
			actualRecord, err = objMgr.GetARecord(recA)

		})

		It("should return expected A record Object", func() {
			Expect(*actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get A Record by Name", func() {
		recA := RecordA{Name: "test"}
		fakeRefReturn := fmt.Sprintf("record:a/ZG5zLmJpbmRfY25h:%s/%20%20", recA.Name)
		aniFakeConnector := &fakeConnector{
			getObjectObj: NewRecordA(recA),
			getObjectRef: "",
			resultObject: []RecordA{*NewRecordA(RecordA{Name: recA.Name, Ref: fakeRefReturn})},
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualRecord *[]RecordA
		var err error
		It("should pass expected A record Object to GetObject", func() {
			actualRecord, err = objMgr.GetARecord(recA)

		})

		It("should return expected A record Object", func() {
			Expect(*actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Get A Record by Reference", func() {
		name := "test"
		fakeRefReturn := fmt.Sprintf("record:a/ZG5zLmJpbmRfY25h:%s/%20%20", name)
		recA := RecordA{Ref: fakeRefReturn}
		aniFakeConnector := &fakeConnector{
			getObjectObj: NewRecordA(recA),
			getObjectRef: fakeRefReturn,
			resultObject: []RecordA{*NewRecordA(RecordA{Name: recA.Name, Ref: fakeRefReturn})},
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)

		var actualRecord *[]RecordA
		var err error
		It("should pass expected A record Object to GetObject", func() {
			actualRecord, err = objMgr.GetARecord(recA)

		})

		It("should return expected A record Object", func() {
			Expect(*actualRecord).To(Equal(aniFakeConnector.resultObject))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete A Record by Reference", func() {

		name := "delete_test"
		fakeRefReturn := fmt.Sprintf("record:a/ZG5zLmJpbmRfY25h:%s/%20%20", name)
		recA := RecordA{Ref: fakeRefReturn}
		aniFakeConnector := &fakeConnector{
			deleteObjectRef: recA.Ref,
			fakeRefReturn:   fakeRefReturn,
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)
		var actualRef string
		var err error
		It("should pass expected A record Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteARecord(recA)
		})
		It("should return expected A record Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete A Record by Name", func() {
		recA := RecordA{Name: "delete_test"}
		fakeRefReturn := fmt.Sprintf("record:a/ZG5zLmJpbmRfY25h:%s/%20%20", recA.Name)
		aniFakeConnector := &fakeConnector{
			getObjectObj: NewRecordA(recA),
			getObjectRef: "",
			resultObject: []RecordA{*NewRecordA(RecordA{Name: recA.Name, Ref: fakeRefReturn})},
			deleteObjectRef: fakeRefReturn,
			fakeRefReturn:   fakeRefReturn,
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)
		var actualRef string
		var err error
		It("should pass expected A record Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteARecord(recA)
		})
		It("should return expected A record Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete A Record by IPv4Addr", func() {
		recA := RecordA{Ipv4Addr: "1.1.1.1"}
		name := "delete_test"
		fakeRefReturn := fmt.Sprintf("record:a/ZG5zLmJpbmRfY25h:%s/%20%20", name)
		aniFakeConnector := &fakeConnector{
			getObjectObj: NewRecordA(recA),
			getObjectRef: "",
			resultObject: []RecordA{*NewRecordA(RecordA{Name: name, Ref: fakeRefReturn, Ipv4Addr: recA.Ipv4Addr})},
			deleteObjectRef: fakeRefReturn,
			fakeRefReturn:   fakeRefReturn,
		}

		cmpType := "Docker"
		tenantID := "01234567890abcdef01234567890abcdef"
		objMgr := NewObjectManager(aniFakeConnector, cmpType, tenantID)
		var actualRef string
		var err error
		It("should pass expected A record Ref to DeleteObject", func() {
			actualRef, err = objMgr.DeleteARecord(recA)
		})
		It("should return expected A record Ref", func() {
			Expect(actualRef).To(Equal(fakeRefReturn))
			Expect(err).To(BeNil())
		})
	})
})