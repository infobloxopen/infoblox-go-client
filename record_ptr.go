package ibclient

import (
	"fmt"
	_ "time"
)

type RecordPTROperations interface {
	CreatePTRRecord(recPTR RecordPTR) (*RecordPTR, error)
	GetPTRRecord(recPTR RecordPTR) (*[]RecordPTR, error)
	DeletePTRRecord(recPTR RecordPTR) (string, error)
	UpdatePTRRecord(recPTR RecordPTR) (*RecordPTR, error)
}

type RecordPTR struct {
	IBBase   `json:"-"`
	Ref              string `json:"_ref,omitempty"`
	Ipv4Addr         string `json:"ipv4addr,omitempty"`
	Ipv6Addr         string `json:"ipv6addr,omitempty"`
	Name             string `json:"name,omitempty"`
	PtrdName         string `json:"ptrdname,omitempty"`
	View             string `json:"view,omitempty"`
	Zone             string `json:"zone,omitempty"`
	Ea               EA     `json:"extattrs,omitempty"`
	NetView          string `json:"omitempty"`
	Cidr             string `json:"omitempty"`
	AddEA            EA     `json:"omitempty"`
	RemoveEA         EA     `json:"omitempty"`
	Creation_time    int    `json:"creation_time,omitempty"`
	Comment          string `json:"comment,omitempty"`
	Creator          string  `json:"creator,omitempty"`
	Ddns_protected   bool    `json:"ddns_protected,omitempty"`
	Dns_name         string  `json:"dns_name,omitempty"`
	Forbid_reclamation bool  `json:"forbid_reclamation,omitempty"`
	Reclaimable      bool     `json:"reclaimable,omitempty"`
 	Ttl              uint     `json:"ttl,omitempty"`
	Use_ttl          bool     `json:"use_ttl,omitempty"`
}

// Define objectType and returnFields for PTR Record 
func NewRecordPTR(rptr RecordPTR) *RecordPTR {
	res := rptr
	res.objectType = "record:ptr"
	res.returnFields = []string{"ipv4addr", "name", "view", "zone","extattrs","comment","creation_time",
		             "creator","ddns_protected","dns_name","cloud_info","forbid_reclamation","last_queried",
		             "reclaimable","ttl","use_ttl","aws_rte53_record_info","ddns_principal","disable","discovered_data","ms_ad_user_data"}

	return &res
}

// CreatePTRRecord takes Name, PtrdName, Ipv4Addr or Ipv6Addr and View of the record to create A Record
// Optional fields: NetView, Ea, Cidr
func (objMgr *ObjectManager) CreatePTRRecord(recPTR RecordPTR) (*RecordPTR, error) {

	recPTR.Ea = objMgr.extendEA(recPTR.Ea)

	recordPTR := NewRecordPTR(recPTR)
	if len(recPTR.Ipv6Addr) > 0 {
		recordPTR.Ipv6Addr = recPTR.Ipv6Addr
	} else {
		if recPTR.Ipv4Addr == "" {
			recordPTR.Ipv4Addr = fmt.Sprintf("func:nextavailableip:%s,%s", recPTR.Cidr, recPTR.NetView)
		} else {
			recordPTR.Ipv4Addr = recPTR.Ipv4Addr
		}
	}

	ref, err := objMgr.connector.CreateObject(recordPTR)
	recordPTR.Ref = ref
	return recordPTR, err
}

// GetPTRRecord by passing Name, reference ID, IPv4 Address, Ptrdname or DNS View
// If no arguments are passed then, all the records are returned
func (objMgr *ObjectManager) GetPTRRecord(recPTR RecordPTR) (*[]RecordPTR, error) {
		
	var res []RecordPTR
	recordPTR := NewRecordPTR(recPTR)

	var err error
	if len(recPTR.Ref)>0 {
		err = objMgr.connector.GetObject(recordPTR, recPTR.Ref, &recordPTR)
		res = append(res,*recordPTR)

	} else {
		recordPTR = NewRecordPTR(recPTR)
		err = objMgr.connector.GetObject(recordPTR, "", &res)
		if err != nil || res == nil || len(res) == 0 {
			return nil, err
		}
	}

	return &res, err
}

// DeletePTRRecord by passing either Reference or Name or IPv4Addr
// If a record with same Ipv4Addr and different name exists, then name and Ipv4Addr has to be passed 
// to avoid multiple record deletions
func (objMgr *ObjectManager) DeletePTRRecord(recPTR RecordPTR) (string, error) {
	var res []RecordPTR
	recordName := NewRecordPTR(RecordPTR{Name: recPTR.Name})
	if len(recPTR.Ref) > 0 {
		return  objMgr.connector.DeleteObject(recPTR.Ref)
	} else {
		recordName = NewRecordPTR(recPTR)
		err := objMgr.connector.GetObject(recordName, "", &res)
		if err != nil || res == nil || len(res) == 0 {
			return "Record doesn't exist", err
		}
		return objMgr.connector.DeleteObject(res[0].Ref)
	}

}

// UpdatePTRRecord takes Reference ID of the record as an argument
// to update Name, Ptrdname and EAs of the record
func (objMgr *ObjectManager) UpdatePTRRecord(recPTR RecordPTR) (*RecordPTR, error) {
	var res RecordPTR
	recordPTR := RecordPTR{Name: recPTR.Name}
	recordPTR.returnFields = []string{"name","extattrs","ptrdname"}
	err := objMgr.connector.GetObject(&recordPTR, recPTR.Ref, &res)
	if err != nil {
		return nil, err
	}
	res.Name = recPTR.Name
	res.PtrdName = recPTR.PtrdName

	for k, v := range recPTR.AddEA {
		res.Ea[k] = v
	}

	for k := range recPTR.RemoveEA {
		_, ok := res.Ea[k]
		if ok {
			delete(res.Ea, k)
		}
	}
	reference, err := objMgr.connector.UpdateObject(&res, recPTR.Ref)
	res.Ref= reference
	err = objMgr.connector.GetObject(&recordPTR, res.Ref, &res)
	return &res, err
}
