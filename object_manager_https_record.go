package ibclient

import "fmt"

func NewEmptyHttpsRecord() *RecordHttps {
	newRecordHttps := &RecordHttps{}
	return newRecordHttps
}
func NewRecordHTTPS(
	name string,
	comment string,
	svcParameters []*Svcparams,
	targetName string,
	disable bool,
	extAttrs EA,
	priority uint32,
	forbidReclamation bool,
	ttl uint32,
	useTtl bool,
	view string,
	ref string) *RecordHttps {

	res := NewEmptyHttpsRecord()

	res.Name = name
	res.Ref = ref
	res.Comment = comment
	res.SvcParameters = svcParameters
	res.TargetName = targetName
	res.Disable = &disable
	res.Ea = extAttrs
	res.Priority = priority
	res.ForbidReclamation = &forbidReclamation
	res.Ttl = ttl
	res.UseTtl = &useTtl
	res.View = view
	return res
}

func (obj *ObjectManager) CreateHTTPSRecord(name string, comment string, svcParameters []*Svcparams, targetName string, disable bool, extAttrs EA, priority uint32, forbidReclamation bool, Ttl uint32, UseTtl bool, view string) (*RecordHttps, error) {

	if priority > 65535 {
		return nil, fmt.Errorf("priority must be between 0 and 65535")
	}

	if name == "" || targetName == "" {
		return nil, fmt.Errorf("name and targetName are required to create HTTPS Record")
	}

	recordHttps := NewRecordHTTPS(name, comment, svcParameters, targetName, disable, extAttrs, priority, forbidReclamation, Ttl, UseTtl, view, "")
	ref, err := obj.connector.CreateObject(recordHttps)
	if err != nil {
		return nil, err
	}
	recordHttps.Ref = ref
	return recordHttps, nil
}

func (objMgr *ObjectManager) GetRecordHTTPSByRef(ref string) (*RecordHttps, error) {
	recordHTTPS := NewEmptyHttpsRecord()
	err := objMgr.connector.GetObject(recordHTTPS, ref, NewQueryParams(false, nil), &recordHTTPS)
	if err != nil {
		return nil, err
	}
	return recordHTTPS, nil
}

func (objMgr *ObjectManager) GetAllHTTPSRecord(queryParams *QueryParams) ([]RecordHttps, error) {
	var res []RecordHttps
	recordHttps := NewEmptyHttpsRecord()
	err := objMgr.connector.GetObject(recordHttps, "", queryParams, &res)
	if err != nil {
		return nil, fmt.Errorf("failed getting HTTPS Record: %s", err)
	}
	return res, nil
}

func (objMgr *ObjectManager) UpdateHttpsRecord(ref string, name string, comment string, svcParameters []*Svcparams, targetName string, disable bool, extAttrs EA, priority uint32, forbidReclamation bool, Ttl uint32, UseTtl bool) (*RecordHttps, error) {
	if priority > 65535 {
		return nil, fmt.Errorf("priority must be between 0 and 65535")
	}
	if name == "" || targetName == "" {
		return nil, fmt.Errorf("name and targetName cannot be empty")
	}
	httpsRecord := NewRecordHTTPS(name, comment, svcParameters, targetName, disable, extAttrs, priority, forbidReclamation, Ttl, UseTtl, "", ref)
	updatedRef, err := objMgr.connector.UpdateObject(httpsRecord, ref)
	if err != nil {
		return nil, err
	}
	httpsRecord.Ref = updatedRef
	return httpsRecord, nil
}

func (objMgr *ObjectManager) DeleteHttpsRecord(ref string) (string, error) {
	return objMgr.connector.DeleteObject(ref)
}
