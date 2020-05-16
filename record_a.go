package ibclient

import (
	"fmt"
	_ "time"
)

type RecordAOperations interface {
	CreateARecord(recA RecordA) (*RecordA, error)
	GetARecord(recA RecordA) (*[]RecordA, error)
	DeleteARecord(recA RecordA) (string, error)
	UpdateARecord(recA RecordA) (*RecordA, error)
}
type Dhcpmember struct {
	Ipv4addr string    `json:"ipv4addr,omitempty"`
	Ipv6addr string    `json:"ipv6addr,omitempty"`
	Name     string    `json:"name,omitempty"`
}
type Cloud_info struct {
	Authority_type      string       `json:"authority_type,omitempty"`
	Delegated_member    Dhcpmember   `json:"delegated_member,omitempty"`
	Delegated_root      string       `json:"delegated_root,omitempty"`
	Delegated_scope     string       `json:"delegated_scope,omitempty"`
	Mgmt_platform       string       `json:"msmt_platform,omitempty"`
	Owned_by_adaptor    bool         `json:"owned_by_adaptor,omitempty"`
	Tenant              string       `json:"tenant,omitempty"`
	Usage               string       `json:"usage,omitempty"`
}

type RecordA struct {
	IBBase   `json:"-"`
	Ref              string `json:"_ref,omitempty"`
	Ipv4Addr         string `json:"ipv4addr,omitempty"`
	Name             string `json:"name,omitempty"`
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
	//CloudInfo       Cloud_info `json:"cloud_info,omitempty"`
}

func NewRecordA(ra RecordA) *RecordA {
	res := ra
	res.objectType = "record:a"

	res.returnFields = []string{ "ipv4addr", "name", "view", "zone","extattrs","comment","creation_time",
				     "creator","ddns_protected","dns_name","cloud_info","forbid_reclamation","last_queried",
				      "reclaimable","ttl","use_ttl","aws_rte53_record_info","ddns_principal","disable","discovered_data","ms_ad_user_data"}
	return &res
}

// CreateARecord takes Name, Ipv4Addr and View of the record to create A Record
// Optional fields: NetView, Ea, Cidr
// Before creating, it checks if the Name and IP passed already exists in the network
func (objMgr *ObjectManager) CreateARecord(recA RecordA) (*RecordA, error) {
	recA.Ea = objMgr.extendEA(recA.Ea)
	recordA := NewRecordA(recA)
	if recA.Ipv4Addr == "" {
		recordA.Ipv4Addr = fmt.Sprintf("func:nextavailableip:%s,%s", recA.Cidr, recA.NetView)
	} else {
		recordA.Ipv4Addr = recA.Ipv4Addr
	}
	ref, err := objMgr.connector.CreateObject(recordA)
	recordA.Ref = ref
	return recordA, err
}

// GetARecord by passing Name, reference ID, IP Address or DNS View
// If no arguments are passed then, all the records are returned
func (objMgr *ObjectManager) GetARecord(recA RecordA) (*[]RecordA, error) {

	var res []RecordA
	recordA := NewRecordA(recA)
	var err error
	if len(recA.Ref)>0 {
		err = objMgr.connector.GetObject(recordA, recA.Ref, &recordA)
		res = append(res,*recordA)

	} else {
		recordA = NewRecordA(recA)
		err = objMgr.connector.GetObject(recordA, "", &res)
		if err != nil || res == nil || len(res) == 0 {
			return nil, err
		}
	}

	return &res, err
}

// DeleteARecord by passing either Reference or Name or IPv4Addr
func (objMgr *ObjectManager) DeleteARecord(recA RecordA) (string, error) {
	var res []RecordA
	recordName := NewRecordA(recA)
	if len(recA.Ref) > 0 {
		return  objMgr.connector.DeleteObject(recA.Ref)

	} else {
		recordName = NewRecordA(recA)
		err := objMgr.connector.GetObject(recordName, "", &res)
		if err != nil || res == nil || len(res) == 0 {
			return "Record doesn't exist", err
		}
		return objMgr.connector.DeleteObject(res[0].Ref)
	}

}

// UpdateARecord takes Reference ID of the record as an argument
// to update Name, IPv4Addr and EAs of the record
// returns updated Refernce ID
func (objMgr *ObjectManager) UpdateARecord(recA RecordA) (*RecordA, error) {
	var res RecordA
	recordA := RecordA{}
	recordA.returnFields = []string{"name","extattrs"}
	err := objMgr.connector.GetObject(&recordA, recA.Ref, &res)
	if err != nil {
		return nil, err
	}
	res.Name = recA.Name
	res.Ipv4Addr = recA.Ipv4Addr
	for k, v := range recA.AddEA {
		res.Ea[k] = v
	}

	for k := range recA.RemoveEA {
		_, ok := res.Ea[k]
		if ok {
			delete(res.Ea, k)
		}
	}
	reference, err := objMgr.connector.UpdateObject(&res, recA.Ref)
	res.Ref= reference
	return &res, err
}
