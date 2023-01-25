package ibclient

import "fmt"

func (objMgr *ObjectManager) CreateMXRecord(
	dnsView string,
	fqdn string,
	mx string,
	priority int,
	ttl uint32,
	useTtl bool,
	comment string,
	eas EA) (*RecordMX, error) {

	if dnsView == "" {
		dnsView = "default"
	}

	if fqdn == "" || mx == "" {
		return nil, fmt.Errorf("fqdn and mx must not be empty")
	}
	recordMx := NewRecordMX(RecordMX{
		View:     dnsView,
		Fqdn:     fqdn,
		MX:       mx,
		Priority: priority,
		Ttl:      ttl,
		UseTtl:   useTtl,
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

func (objMgr *ObjectManager) GetMXRecord(dnsView string, fqdn string) (*RecordMX, error) {
	if dnsView == "" || fqdn == "" {
		return nil, fmt.Errorf("DNS view and fqdn are required to retrieve a unique mx record")
	}
	var res []RecordMX

	recordMX := NewRecordMX(RecordMX{})

	sf := map[string]string{
		"view": dnsView,
		"name": fqdn,
	}
	queryParams := NewQueryParams(false, sf)
	err := objMgr.connector.GetObject(recordMX, "", queryParams, &res)

	if err != nil {
		return nil, err
	}

	return &res[0], err
}

func (objMgr *ObjectManager) UpdateMXRecord(
	ref string,
	dnsView string,
	fqdn string,
	mx string,
	ttl uint32,
	useTtl bool,
	comment string,
	priority int,
	eas EA) (*RecordMX, error) {

	res, _ := objMgr.GetMXRecordByRef(ref)

	if dnsView != res.View {
		return nil, fmt.Errorf("changing dns_view after resource creation is not allowed")
	}

	if priority < 0 {
		return nil, fmt.Errorf("priority must be greater than zero")
	}

	if mx == "" {
		return nil, fmt.Errorf("mx must not be empty")
	}
	recordMx := NewRecordMX(RecordMX{
		View:     dnsView,
		Fqdn:     fqdn,
		MX:       mx,
		Priority: priority,
		Ttl:      ttl,
		UseTtl:   useTtl,
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
