package ibclient

import "fmt"

func (objMgr *ObjectManager) CreateMXRecord(
	dnsview string,
	fqdn string,
	mx string,
	priority int,
	comment string,
	eas EA) (*RecordMX, error) {

	if dnsview == "" {
		dnsview = "default"
	}

	if fqdn == "" || mx == "" {
		return nil, fmt.Errorf("fqdn and mx must not be empty")
	}
	recordMx := NewRecordMX(RecordMX{
		dnsView:  dnsview,
		Fqdn:     fqdn,
		MX:       mx,
		Priority: priority,
		Comment:  comment,
		Ea:       eas,
	})

	ref, err := objMgr.connector.CreateObject(recordMx)
	if err != nil {
		return nil, err
	}
	recordMx.Ref = ref
	return recordMx, err
}

func (objMgr *ObjectManager) GetMXRecordByRef(ref string) (*RecordMX, error) {
	recordMX := NewRecordMX(RecordMX{})
	err := objMgr.connector.GetObject(recordMX, ref, NewQueryParams(false, nil), &recordMX)

	return recordMX, err
}

func (objMgr *ObjectManager) GetMXRecord(dnsview string, fqdn string) (*RecordMX, error) {
	if dnsview == "" || fqdn == "" {
		return nil, fmt.Errorf("DNS view and fqdn are required to retrieve a unique mx record")
	}
	var res []RecordMX

	recordMX := NewRecordMX(RecordMX{})

	sf := map[string]string{
		"view": dnsview,
		"name": fqdn,
	}
	queryParams := NewQueryParams(false, sf)
	err := objMgr.connector.GetObject(recordMX, "", queryParams, &res)

	if err != nil {
		return nil, err
	} else if res == nil || len(res) == 0 {
		return nil, NewNotFoundError(
			fmt.Sprintf(
				"MX record with name '%s' in DNS view '%s' is not found",
				fqdn, dnsview))
	}

	return &res[0], err
}

func (objMgr *ObjectManager) UpdateMXRecord(
	ref string,
	dnsview string,
	fqdn string,
	mx string,
	comment string,
	priority int,
	eas EA) (*RecordMX, error) {

	recordMx := NewRecordMX(RecordMX{
		dnsView:  dnsview,
		Fqdn:     fqdn,
		MX:       mx,
		Priority: priority,
		Comment:  comment,
		Ea:       eas,
	})

	recordMx.Ref = ref

	nw_ref, err := objMgr.connector.UpdateObject(recordMx, ref)

	if err != nil {
		return nil, err
	}

	recordMx.Ref = nw_ref

	return recordMx, err
}

func (objMgr *ObjectManager) DeleteMXRecord(ref string) (string, error) {
	return objMgr.connector.DeleteObject(ref)
}
