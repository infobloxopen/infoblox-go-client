package ibclient

import "fmt"

// Creates TXT Record. Use TTL of 0 to inherit TTL from the Zone
func (objMgr *ObjectManager) CreateTXTRecord(
	recordName string,
	text string,
	dnsView string,
	useTtl bool,
	ttl uint32,
	comment string,
	eas EA) (*RecordTXT, error) {

	recordTXT := NewRecordTXT(recordName, text, dnsView, "", useTtl, ttl, comment, eas)

	ref, err := objMgr.connector.CreateObject(recordTXT)
	if err != nil {
		return nil, err
	}
	recordTXT, err = objMgr.GetTXTRecordByRef(ref)
	return recordTXT, err
}

func (objMgr *ObjectManager) GetTXTRecordByRef(ref string) (*RecordTXT, error) {
	recordTXT := NewEmptyRecordTXT()
	err := objMgr.connector.GetObject(
		recordTXT, ref, NewQueryParams(false, nil), &recordTXT)
	return recordTXT, err
}

func (objMgr *ObjectManager) GetTXTRecord(name string) (*RecordTXT, error) {
	if name == "" {
		return nil, fmt.Errorf("name can not be empty")
	}
	var res []RecordTXT

	recordTXT := NewEmptyRecordTXT()

	sf := map[string]string{
		"name": name,
	}
	queryParams := NewQueryParams(false, sf)
	err := objMgr.connector.GetObject(recordTXT, "", queryParams, &res)

	if err != nil || res == nil || len(res) == 0 {
		return nil, err
	}

	return &res[0], nil
}

func (objMgr *ObjectManager) UpdateTXTRecord(
	ref string,
	recordName string,
	text string,
	useTtl bool,
	ttl uint32,
	comment string,
	eas EA) (*RecordTXT, error) {

	recordTXT := NewRecordTXT(recordName, text, "", "", useTtl, ttl, comment, eas)
	recordTXT.Ref = ref

	reference, err := objMgr.connector.UpdateObject(recordTXT, ref)
	if err != nil {
		return nil, err
	}

	recordTXT, err = objMgr.GetTXTRecordByRef(reference)
	return recordTXT, err
}

func (objMgr *ObjectManager) DeleteTXTRecord(ref string) (string, error) {
	return objMgr.connector.DeleteObject(ref)
}
