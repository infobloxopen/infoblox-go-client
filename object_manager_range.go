package ibclient

import (
	"fmt"
)

func (objMgr *ObjectManager) CreateRange(netViewName, name, startAddr, endAddr, comment string, disable bool, exclude []*Exclusionrange) (*Range, error) {
	newRange := NewEmptyRange()
	newRange.Name = &name
	newRange.StartAddr = &startAddr
	newRange.EndAddr = &endAddr
	newRange.Comment = &comment
	newRange.NetworkView = &netViewName
	newRange.Exclude = exclude
	newRange.Disable = &disable

	refCreatedObj, err := objMgr.connector.CreateObject(newRange)
	if err != nil {
		return nil, err
	}

	return objMgr.GetRange(refCreatedObj)
}

func (objMgr *ObjectManager) UpdateRange(ref, name, startAddr, endAddr, comment string, disable bool, exclude []*Exclusionrange) (*Range, error) {
	r := &Range{
		Comment: &comment,
		Disable: &disable,
	}

	if len(exclude) > 0 {
		r.Exclude = exclude
	}

	if len(startAddr) > 0 {
		r.StartAddr = &startAddr
	}

	if len(endAddr) > 0 {
		r.EndAddr = &endAddr
	}

	if len(name) > 0 {
		r.Name = &name
	}

	updateRef, err := objMgr.connector.UpdateObject(r, ref)
	if err != nil {
		return nil, err
	}

	return objMgr.GetRange(updateRef)
}

func (objMgr *ObjectManager) GetRange(ref string) (*Range, error) {
	var res []Range
	search := &Range{Ref: ref}
	qp := NewQueryParams(false, nil)
	err := objMgr.connector.GetObject(search, ref, qp, &res)
	if err != nil {
		return nil, err
	} else if res == nil || len(res) == 0 {
		return nil, NewNotFoundError(
			fmt.Sprintf(
				"could not find range object"))
	}

	return &res[0], nil
}

func (objMgr *ObjectManager) DeleteRange(ref string) error {
	_, err := objMgr.connector.DeleteObject(ref)
	return err
}
