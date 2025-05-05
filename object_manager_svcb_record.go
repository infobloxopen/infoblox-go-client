package ibclient

import "fmt"

func (objMgr *ObjectManager) CreateSVCBRecord(name string, comment string, disable bool, ea EA,
	priority uint32, svcParams []SVCParams, targetName string, useTtl bool, ttl uint32, view string,
	creator string, ddnsPrincipal string, ddnsProtectd bool) (*RecordSVCB, error) {
	if name == "" || priority == 0 || targetName == "" {
		return nil, fmt.Errorf("name, priority and targetName fields are required to create a SVCB Record")
	}
	recordSVCB := NewSVCBRecord("", name, comment, disable, ea, priority, svcParams, targetName, useTtl, ttl, creator,
		ddnsPrincipal, ddnsProtectd)
	recordSVCB.View = view
	ref, err := objMgr.connector.CreateObject(recordSVCB)
	if err != nil {
		return nil, fmt.Errorf("error creating SVCB Record %s, err: %s", name, err)
	}
	recordSVCB.Ref = ref
	return recordSVCB, nil
}

func (objMgr *ObjectManager) DeleteSVCBRecord(ref string) (string, error) {
	return objMgr.connector.DeleteObject(ref)
}

func (objMgr *ObjectManager) GetAllSVCBRecords(queryParams *QueryParams) ([]RecordSVCB, error) {
	var res []RecordSVCB
	recordSVCB := NewEmptyRecordSVCB()
	err := objMgr.connector.GetObject(recordSVCB, "", queryParams, &res)
	if err != nil {
		return nil, fmt.Errorf("failed getting SVCB Record: %s", err)
	}
	return res, nil
}

func (objMgr *ObjectManager) GetSVCBRecordByRef(ref string) (*RecordSVCB, error) {
	recordSVCB := NewEmptyRecordSVCB()
	err := objMgr.connector.GetObject(recordSVCB, ref, NewQueryParams(false, nil), &recordSVCB)
	if err != nil {
		return nil, err
	}
	return recordSVCB, nil
}

func (objMgr *ObjectManager) UpdateSVCBRecord(ref string, name string, comment string, disable bool, ea EA,
	priority uint32, svcParams []SVCParams, targetName string, useTtl bool, ttl uint32, creator string,
	ddnsPrincipal string, ddnsProtectd bool) (*RecordSVCB, error) {
	recordSVCB := NewSVCBRecord(ref, name, comment, disable, ea, priority, svcParams, targetName, useTtl, ttl, creator,
		ddnsPrincipal, ddnsProtectd)
	newRef, err := objMgr.connector.UpdateObject(recordSVCB, ref)
	if err != nil {
		return nil, fmt.Errorf("error updating SVCB Record %s, err: %s", name, err)
	}
	recordSVCB.Ref = newRef
	recordSVCB, err = objMgr.GetSVCBRecordByRef(newRef)
	if err != nil {
		return nil, fmt.Errorf("error getting updated SVCB Record %s, err: %s", name, err)
	}
	return recordSVCB, nil
}

func NewEmptyRecordSVCB() *RecordSVCB {
	recordSVCB := RecordSVCB{}
	recordSVCB.SetReturnFields(append(recordSVCB.returnFields))
	return &recordSVCB
}

func NewSVCBRecord(ref string, name string, comment string, disable bool, ea EA, priority uint32, svcParams []SVCParams,
	targetName string, useTtl bool, ttl uint32, creator string, ddnsPrincipal string, ddnsProtectd bool) *RecordSVCB {
	recordSVCB := NewEmptyRecordSVCB()
	recordSVCB.Ref = ref
	recordSVCB.Name = name
	recordSVCB.Comment = comment
	recordSVCB.Disable = disable
	recordSVCB.Ea = ea
	recordSVCB.Priority = priority
	recordSVCB.SvcParameters = svcParams
	recordSVCB.TargetName = targetName
	recordSVCB.Creator = creator
	recordSVCB.DdnsProtected = ddnsProtectd
	recordSVCB.DdnsPrincipal = ddnsPrincipal
	recordSVCB.UseTtl = useTtl
	recordSVCB.Ttl = ttl
	return recordSVCB
}
