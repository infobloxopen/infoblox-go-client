package ibclient

import "fmt"

func (objMgr *ObjectManager) CreateSRVRecord(
	dnsView string,
	fqdn string,
	priority int,
	weight int,
	port int,
	target string,
	ttl uint32,
	useTtl bool,
	comment string,
	eas EA) (*RecordSRV, error) {

	if dnsView == "" {
		dnsView = "default"
	}

	if fqdn == "" {
		return nil, fmt.Errorf("fqdn must not be empty")
	}

	if priority < 0 || weight < 0 {
		return nil, fmt.Errorf("priority and weight can't be a negative number")
	}

	if target == "" {
		return nil, fmt.Errorf("target must not be empty")
	}

	recordSRV := NewRecordSRV(RecordSRV{
		View:     dnsView,
		Fqdn:     fqdn,
		Priority: priority,
		Weight:   weight,
		Port:     port,
		Target:   target,
		Ttl:      ttl,
		UseTtl:   useTtl,
		Comment:  comment,
		Ea:       eas,
	})

	ref, err := objMgr.connector.CreateObject(recordSRV)

	if err != nil {
		return nil, err
	}

	recordSRV.Ref = ref
	return recordSRV, err

}

func (objMgr *ObjectManager) GetSRVRecord(dnsView string, fqdn string) (*RecordSRV, error) {
	if dnsView == "" || fqdn == "" {
		return nil, fmt.Errorf("DNS view and fqdn are required to retrieve a unique srv record")
	}
	var res []RecordSRV

	recordSRV := NewRecordSRV(RecordSRV{})

	sf := map[string]string{
		"view": dnsView,
		"name": fqdn,
	}
	queryParams := NewQueryParams(false, sf)
	err := objMgr.connector.GetObject(recordSRV, "", queryParams, &res)

	if err != nil {
		return nil, err
	}

	if res == nil || len(res) == 0 {
		return nil, NewNotFoundError(
			fmt.Sprintf(
				"SRV record with name '%s' in DNS view '%s' is not found",
				fqdn, dnsView))
	}

	return &res[0], nil
}

func (objMgr *ObjectManager) GetSRVRecordByRef(ref string) (*RecordSRV, error) {
	recordSRV := NewRecordSRV(RecordSRV{})
	err := objMgr.connector.GetObject(recordSRV, ref, NewQueryParams(false, nil), &recordSRV)

	return recordSRV, err
}

func (objMgr *ObjectManager) UpdateSRVRecord(
	ref string,
	dnsView string,
	fqdn string,
	priority int,
	weight int,
	port int,
	target string,
	ttl uint32,
	useTtl bool,
	comment string,
	eas EA) (*RecordSRV, error) {

	res, err := objMgr.GetSRVRecordByRef(ref)

	if err != nil {
		return nil, err
	}

	if dnsView != res.View {
		return nil, fmt.Errorf("changing dns_view after object creation is not allowed")
	}

	if priority < 0 {
		return nil, fmt.Errorf("priority field must not be a negative number")
	}

	if port < 0 || weight < 0 {
		return nil, fmt.Errorf("port or weight must not be a negative number")
	}

	recordSRV := NewRecordSRV(RecordSRV{
		Fqdn:     fqdn,
		Priority: priority,
		Weight:   weight,
		Port:     port,
		Target:   target,
		Ttl:      ttl,
		UseTtl:   useTtl,
		Comment:  comment,
		Ea:       eas,
	})

	recordSRV.Ref = ref

	nw_ref, err := objMgr.connector.UpdateObject(recordSRV, ref)

	if err != nil {
		return nil, err
	}

	recordSRV.Ref = nw_ref

	return recordSRV, err
}

func (objMgr *ObjectManager) DeleteSRVRecord(ref string) (string, error) {
	return objMgr.connector.DeleteObject(ref)
}
