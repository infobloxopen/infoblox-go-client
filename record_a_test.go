package ibclient

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	"reflect"
)

var _ = Describe("Testing with Ginkgo", func() {
	It("new record a", func() {

		type args struct {
			ra RecordA
		}
		recA := NewRecordA(RecordA{})
		tests := []struct {
			name string
			args args
			want *RecordA
		}{
			{"test", args{ra: RecordA{IBBase: IBBase{objectType: "record:a",
				returnFields: []string{"ipv4addr", "name", "view", "zone", "extattrs", "comment", "creation_time",
					"creator", "ddns_protected", "dns_name", "cloud_info", "forbid_reclamation", "last_queried",
					"reclaimable", "ttl", "use_ttl", "aws_rte53_record_info", "ddns_principal", "disable", "discovered_data", "ms_ad_user_data"}}}}, recA},
		}
		for _, tt := range tests {
			GinkgoT().Log(tt.name, func(t GinkgoTInterface) {
				if got := NewRecordA(tt.args.ra); !reflect.DeepEqual(got, tt.want) {
					GinkgoT().Errorf("NewRecordA() = %v, want %v", got, tt.want)
				}
			})
		}
	})
	It("object manager_ create a record", func() {

		type Fields struct {
			connector *fakeConnector
			cmpType   string
			tenantID  string
		}
		vmID := "93f9249abc039284"
		vmName := "dummyvm"
		ea := EA{"VM ID": vmID, "VM Name": vmName}

		recA := RecordA{NetView: "private",
			Cidr:     "53.0.0.0/24",
			Ipv4Addr: "53.0.0.1",
			View:     "default",
			Name:     "test",
			Ea:       ea,
			AddEA:    ea}

		rec1A := RecordA{NetView: "private",
			Cidr:     "53.0.0.0/24",
			Ipv4Addr: "",
			View:     "default",
			Name:     "test1",
			Ea:       ea,
		}
		rec2A := RecordA{View: "default",
			Name:     "test",
			Ipv4Addr: "53.0.0.2",
			Ea:       ea}

		fakeRefReturn := fmt.Sprintf("record:a/ZG5zLmJpbmRfY25h:%s/%20%20", recA.Name)
		recA.Ref = fakeRefReturn
		aniFakeConnector := &fakeConnector{
			createObjectObj: NewRecordA(recA),
			getObjectRef:    fakeRefReturn,
			getObjectObj: NewRecordA(RecordA{
				Name:     recA.Name,
				View:     recA.View,
				Ipv4Addr: recA.Ipv4Addr,
				Ref:      fakeRefReturn,
				Ea:       ea,
			}),
			resultObject:  NewRecordA(recA),
			fakeRefReturn: fakeRefReturn,
		}
		fakeRefReturn1 := fmt.Sprintf("record:a/ZG5zLmJpbmRfY25h:%s/%20%20", rec1A.Name)
		rec1A.Ref = fakeRefReturn1
		//rec1A.Ipv4Addr = fmt.Sprintf("func:nextavailableip:%s,%s", rec1A.Cidr, rec1A.NetView)
		aniFakeConnector1 := &fakeConnector{
			createObjectObj: NewRecordA(rec1A),
			getObjectRef:    fakeRefReturn1,
			getObjectObj: NewRecordA(RecordA{
				Name:     rec1A.Name,
				View:     rec1A.View,
				Ipv4Addr: fmt.Sprintf("func:nextavailableip:%s,%s", rec1A.Cidr, rec1A.NetView),
				Ref:      fakeRefReturn1,
				Ea:       ea,
			}),
			resultObject: NewRecordA(RecordA{NetView: "private",
				Cidr:     "53.0.0.0/24",
				Name:     rec1A.Name,
				View:     rec1A.View,
				Ipv4Addr: fmt.Sprintf("func:nextavailableip:%s,%s", rec1A.Cidr, rec1A.NetView),
				Ref:      fakeRefReturn1,
				Ea:       ea,
			}),
			fakeRefReturn: fakeRefReturn1,
		}
		fakeRefReturn2 := fmt.Sprintf("record:a/ZG5zLmJpbmRfY25h:%s/%20%20", recA.Name)
		rec2A.Ref = fakeRefReturn2
		aniFakeConnector2 := &fakeConnector{
			createObjectObj: NewRecordA(rec2A),
			getObjectRef:    fakeRefReturn2,
			getObjectObj: NewRecordA(RecordA{
				Name:     rec2A.Name,
				View:     rec2A.View,
				Ipv4Addr: rec2A.Ipv4Addr,
				Ref:      fakeRefReturn2,
				Ea:       ea,
			}),
			resultObject:  NewRecordA(rec2A),
			fakeRefReturn: fakeRefReturn2,
		}

		res := aniFakeConnector.resultObject
		res1 := aniFakeConnector1.resultObject
		res2 := aniFakeConnector2.resultObject
		fields_test := Fields{connector: aniFakeConnector, cmpType: "Docker", tenantID: "01234567890abcdef01234567890abcdef"}
		fields_test1 := Fields{connector: aniFakeConnector1, cmpType: "Docker", tenantID: "01234567890abcdef01234567890abcdef"}
		fields_test2 := Fields{connector: aniFakeConnector2, cmpType: "Docker", tenantID: "01234567890abcdef01234567890abcdef"}
		tests := []struct {
			name    string
			fields  Fields
			recA    RecordA
			want    interface{}
			wantErr bool
		}{
			{"test", fields_test, recA, res, false},
			{"test1", fields_test1, rec1A, res1, false},
			{"test2", fields_test2, rec2A, res2, false},
		}
		for _, tt := range tests {
			GinkgoT().Log(tt.name, func(t GinkgoTInterface) {
				objMgr := &ObjectManager{
					connector: tt.fields.connector,
					cmpType:   tt.fields.cmpType,
					tenantID:  tt.fields.tenantID,
				}
				ea := objMgr.getBasicEA(true)

				//aniFakeConnector.createObjectObj.(*RecordA).Ea = ea
				aniFakeConnector.resultObject.(*RecordA).Ea = ea
				aniFakeConnector.resultObject.(*RecordA).Ea["VM ID"] = vmID
				aniFakeConnector.resultObject.(*RecordA).Ea["VM Name"] = vmName
				//aniFakeConnector.getObjectObj.(*RecordA).Ea = ea

				//aniFakeConnector1.createObjectObj.(*RecordA).Ea = ea
				aniFakeConnector1.resultObject.(*RecordA).Ea = ea
				aniFakeConnector1.resultObject.(*RecordA).Ea["VM ID"] = vmID
				aniFakeConnector1.resultObject.(*RecordA).Ea["VM Name"] = vmName
				//aniFakeConnector1.getObjectObj.(*RecordA).Ea = ea

				//aniFakeConnector2.createObjectObj.(*RecordA).Ea = ea
				aniFakeConnector2.resultObject.(*RecordA).Ea = ea
				aniFakeConnector2.resultObject.(*RecordA).Ea["VM ID"] = vmID
				aniFakeConnector2.resultObject.(*RecordA).Ea["VM Name"] = vmName
				//aniFakeConnector2.getObjectObj.(*RecordA).Ea = ea

				got, err := objMgr.CreateARecord(tt.recA)
				if (err != nil) != tt.wantErr {
					GinkgoT().Errorf("ObjectManager.CreateARecord() error = %v, wantErr %v, got=%v,want=%v", err, tt.wantErr, got, tt.want)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					GinkgoT().Errorf("ObjectManager.CreateARecord() = %v, want %v", got, tt.want)
				}
			})
		}
	})
	It("object manager_ get a record", func() {

		type fields struct {
			connector IBConnector
			cmpType   string
			tenantID  string
		}

		// GetARecord by passing Name
		recA := RecordA{Name: "test"}
		fakeRefReturn := fmt.Sprintf("record:a/ZG5zLmJpbmRfY25h:%s/%20%20", recA.Name)
		aniFakeConnector := &fakeConnector{
			getObjectObj: NewRecordA(RecordA{}),
			getObjectRef: "",
			resultObject: []RecordA{*NewRecordA(RecordA{Name: recA.Name, Ref: fakeRefReturn})},
		}
		res := aniFakeConnector.resultObject
		fields_test := fields{connector: aniFakeConnector, cmpType: "Docker", tenantID: "01234567890abcdef01234567890abcdef"}

		// GetARecord by passing Reference
		name := "test1"
		fakeRefReturn1 := fmt.Sprintf("record:a/ZG5zLmJpbmRfY25h:%s/%20%20", name)
		rec1A := RecordA{Ref: fakeRefReturn1}
		aniFakeConnector1 := &fakeConnector{
			getObjectObj: NewRecordA(RecordA{}),
			getObjectRef: fakeRefReturn1,
			resultObject: []RecordA{*NewRecordA(RecordA{Ref: fakeRefReturn1})},
		}
		res1 := aniFakeConnector1.resultObject
		fields_test1 := fields{connector: aniFakeConnector1, cmpType: "Docker", tenantID: "01234567890abcdef01234567890abcdef"}

		// GetARecord by passing IPv4Addr
		ipaddr := "10.0.0.12"
		name2 := "test2"
		fakeRefReturn2 := fmt.Sprintf("record:a/ZG5zLmJpbmRfY25h:%s/%20%20", name2)
		rec2A := RecordA{Ipv4Addr: ipaddr}
		aniFakeConnector2 := &fakeConnector{
			getObjectObj: NewRecordA(RecordA{}),
			getObjectRef: fakeRefReturn2,
			resultObject: []RecordA{*NewRecordA(RecordA{Ref: fakeRefReturn2})},
		}
		res2 := aniFakeConnector2.resultObject
		fields_test2 := fields{connector: aniFakeConnector2, cmpType: "Docker", tenantID: "01234567890abcdef01234567890abcdef"}
		tests := []struct {
			name    string
			fields  fields
			recA    RecordA
			want    interface{}
			wantErr bool
		}{
			{"test", fields_test, recA, res, false},
			{"test1", fields_test1, rec1A, res1, false},
			{"test2", fields_test2, rec2A, res2, false},
		}
		for _, tt := range tests {
			GinkgoT().Log(tt.name, func(t GinkgoTInterface) {
				objMgr := &ObjectManager{
					connector: tt.fields.connector,
					cmpType:   tt.fields.cmpType,
					tenantID:  tt.fields.tenantID,
				}
				got1, err := objMgr.GetARecord(tt.recA)
				got := *got1
				if (err != nil) != tt.wantErr {
					GinkgoT().Errorf("ObjectManager.GetARecord() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					GinkgoT().Errorf("ObjectManager.GetARecord() = %v, want %v", got, tt.want)
				}
			})
		}
	})
	It("object manager_ delete a record", func() {

		type fields struct {
			connector IBConnector
			cmpType   string
			tenantID  string
		}

		// DeleteARecord by passing Reference
		name := "delete_test"
		fakeRefReturn := fmt.Sprintf("record:a/ZG5zLmJpbmRfY25h:%s/%20%20", name)
		recA := RecordA{Ref: fakeRefReturn}
		nwFakeConnector := &fakeConnector{
			deleteObjectRef: recA.Ref,
			fakeRefReturn:   fakeRefReturn,
		}

		// DeleteARecord by passing Name
		rec1A := RecordA{Name: "delete_test1", View: "default"}
		fakeRefReturn1 := fmt.Sprintf("record:a/ZG5zLmJpbmRfY25h:%s/%20%20", rec1A.Name)
		nwFakeConnector1 := &fakeConnector{
			getObjectObj:    NewRecordA(RecordA{}),
			getObjectRef:    fakeRefReturn1,
			resultObject:    []RecordA{*NewRecordA(RecordA{Ref: fakeRefReturn1, Name: "delete_test1", Ipv4Addr: "10.0.0.12"})},
			deleteObjectRef: fakeRefReturn1,
			fakeRefReturn:   fakeRefReturn1,
		}

		// DeleteARecord by passing IPv4Addr
		rec2A := RecordA{Ipv4Addr: "10.0.0.12"}
		fakeRefReturn2 := fmt.Sprintf("record:a/ZG5zLmJpbmRfY25h:%s/%20%20", "delete_test2")
		nwFakeConnector2 := &fakeConnector{
			getObjectObj:    NewRecordA(RecordA{}),
			getObjectRef:    fakeRefReturn2,
			resultObject:    []RecordA{*NewRecordA(RecordA{Ref: fakeRefReturn2, Name: "delete_test2"})},
			deleteObjectRef: fakeRefReturn2,
			fakeRefReturn:   fakeRefReturn2,
		}

		// Returns error message if anything else is passed
		rec3A := RecordA{NetView: "private"}
		nwFakeConnector3 := &fakeConnector{
			getObjectObj:    NewRecordA(RecordA{}),
			resultObject:    []RecordA{*NewRecordA(RecordA{Name: "delete_test3"})},
			deleteObjectRef: "",
		}
		fields_test := fields{connector: nwFakeConnector, cmpType: "Docker", tenantID: "01234567890abcdef01234567890abcdef"}
		fields_test1 := fields{connector: nwFakeConnector1, cmpType: "Docker", tenantID: "01234567890abcdef01234567890abcdef"}
		fields_test2 := fields{connector: nwFakeConnector2, cmpType: "Docker", tenantID: "01234567890abcdef01234567890abcdef"}
		fields_test3 := fields{connector: nwFakeConnector3, cmpType: "Docker", tenantID: "01234567890abcdef01234567890abcdef"}

		tests := []struct {
			name    string
			fields  fields
			recA    RecordA
			want    interface{}
			wantErr bool
		}{
			{"test", fields_test, recA, nwFakeConnector.deleteObjectRef, false},
			{"test1", fields_test1, rec1A, nwFakeConnector1.deleteObjectRef, false},
			{"test2", fields_test2, rec2A, nwFakeConnector2.deleteObjectRef, false},
			{"test3", fields_test3, rec3A, nwFakeConnector3.deleteObjectRef, true},
		}
		for _, tt := range tests {
			GinkgoT().Log(tt.name, func(t GinkgoTInterface) {
				objMgr := &ObjectManager{
					connector: tt.fields.connector,
					cmpType:   tt.fields.cmpType,
					tenantID:  tt.fields.tenantID,
				}
				got, err := objMgr.DeleteARecord(tt.recA)
				if (err != nil) != tt.wantErr {
					GinkgoT().Errorf("ObjectManager.DeleteARecord() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					GinkgoT().Errorf("ObjectManager.DeleteARecord() = %v, want %v", got, tt.want)
				}
			})
		}
	})
})