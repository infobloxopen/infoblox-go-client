package ibclient

import (
	"errors"
	"fmt"
	"log"
	_"time"
)

type RecordAAAAOperations interface {
	CreateAAAARecord(recA4 RecordAAAA) (*RecordAAAA, error)
	GetAAAARecord(recA4 RecordAAAA) (*[]RecordAAAA, error)
	DeleteAAAARecord(recA4 RecordAAAA) (string, error)
	UpdateAAAARecord(recA4 RecordAAAA) (*RecordAAAA, error)
}

type RecordAAAA struct {
	IBBase   `json:"-"`
	Ref              string `json:"_ref,omitempty"`
	Ipv6Addr         string `json:"ipv6addr,omitempty"`
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
}

func NewRecordAAAA(ra RecordAAAA) *RecordAAAA {
	res := ra
	res.objectType = "record:aaaa"
	res.returnFields = []string{ "ipv6addr", "name", "view", "zone","extattrs","comment","creation_time",
				     "creator","ddns_protected","dns_name","cloud_info","forbid_reclamation","last_queried",
				      "reclaimable","ttl","use_ttl","aws_rte53_record_info","ddns_principal","disable","discovered_data","ms_ad_user_data"}
	return &res
}

// CreateAAAARecord takes Name, Ipv6Addr and View of the record to create AAAA Record
// Optional fields: NetView, Ea, Cidr
// Before creating, it checks if the Name and IP passed already exists in the network
func (objMgr *ObjectManager) CreateAAAARecord(recA4 RecordAAAA) (*RecordAAAA, error) {

	eas := objMgr.extendEA(recA4.Ea)
	recordA4 := NewRecordAAAA(RecordAAAA{
		View: recA4.View,
		Name: recA4.Name,
		Ea:   eas})

	if recA4.Ipv6Addr == "" {
		recordA4.Ipv6Addr = fmt.Sprintf("func:nextavailableip:%s,%s", recA4.Cidr, recA4.NetView)
	} else {
		recordA4.Ipv6Addr = recA4.Ipv6Addr
	}
	ref, err := objMgr.connector.CreateObject(recordA4)
	recordA4.Ref = ref
	recordA4.returnFields = []string{"name", "extattrs", "view", "zone", "ipv6addr"}
	err = objMgr.connector.GetObject(recordA4, ref, &recordA4)
	return recordA4, err
}

// GetAAAARecord by passing Name, reference ID, IP Address or DNS View
// If no arguments are passed then, all the records are returned
func (objMgr *ObjectManager) GetAAAARecord(recA4 RecordAAAA) (*[]RecordAAAA, error) {
		
	var res []RecordAAAA
	recordA4 := NewRecordAAAA(RecordAAAA{})
	var err error
	if len(recA4.Ref)>0 {
		err = objMgr.connector.GetObject(recordA4, recA4.Ref, &recordA4)
		res = append(res,*recordA4)

	} else {
		recordA4 = NewRecordAAAA(recA4)
		err = objMgr.connector.GetObject(recordA4, "", &res)
		if err != nil || res == nil || len(res) == 0 {
			return nil, err
		}
	}

	return &res, err
}

// DeleteAAAARecord by passing either Reference or Name or IPv6Addr
func (objMgr *ObjectManager) DeleteAAAARecord(recA4 RecordAAAA) (string, error) {
	var res []RecordAAAA
	recordName := NewRecordAAAA(RecordAAAA{})
	if len(recA4.Ref) > 0 {
		return  objMgr.connector.DeleteObject(recA4.Ref)

	} else if len (recA4.Name) > 0 {
		recordName = NewRecordAAAA(RecordAAAA{Name: recA4.Name})

	} else if len (recA4.Ipv6Addr) > 0 {
		recordName = NewRecordAAAA(RecordAAAA{Ipv6Addr: recA4.Ipv6Addr})

	} else {
		log.Printf("ERROR, Enter Name or Reference ID of A record")
		return "", errors.New("INVALID INPUT")
	}
	err := objMgr.connector.GetObject(recordName, "", &res)
	if err != nil || res == nil || len(res) == 0 {
		return "Record doesn't exist", err
	}
	return objMgr.connector.DeleteObject(res[0].Ref)
}

// UpdateAAAARecord takes Reference ID of the record as an argument
// to update Name and EAs of the record
func (objMgr *ObjectManager) UpdateAAAARecord(recA4 RecordAAAA) (*RecordAAAA, error) {
	var res RecordAAAA
	recordA4 := RecordAAAA{Name: recA4.Name}
	recordA4.returnFields = []string {"name","extattrs"}
	err := objMgr.connector.GetObject(&recordA4, recA4.Ref, &res)
	if err != nil {
		return nil, err
	}
	res.Name = recA4.Name

	for k, v := range recA4.AddEA {
		res.Ea[k] = v
	}

	for k := range recA4.RemoveEA {
		_, ok := res.Ea[k]
		if ok {
			delete(res.Ea, k)
		}
	}
	reference, err := objMgr.connector.UpdateObject(&res, recA4.Ref)
	res.Ref= reference
	err = objMgr.connector.GetObject(&recordA4, res.Ref, &res)
	return &res, err
}
