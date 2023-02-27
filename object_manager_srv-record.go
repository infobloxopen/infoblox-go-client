package ibclient

import (
	"fmt"
	"regexp"
)

func (objMgr *ObjectManager) CreateSRVRecord(
	dnsView string,
	name string,
	priority int,
	weight int,
	port uint32,
	target string,
	ttl uint32,
	useTtl bool,
	comment string,
	eas EA) (*RecordSRV, error) {

	nameRegex := `^_[a-z]+\._[a-z]+\.[a-z0-9-]+\.[a-z]+$`
	targetRegex := `^[a-z]+\.[a-z0-9-]+\.[a-z]+$`

	valid, _ := regexp.MatchString(nameRegex, name)
	valid_tg, _ := regexp.MatchString(targetRegex, name)

	if !valid {
		return nil, fmt.Errorf("'name' format is not valid")
	}

	if !valid_tg {
		return nil, fmt.Errorf("'target' is not in valid format")
	}

	if dnsView == "" {
		dnsView = "default"
	}

	if name == "" {
		return nil, fmt.Errorf("'name' must not be empty")
	}

	if priority < 0 || weight < 0 {
		return nil, fmt.Errorf("'priority' and 'weight' can't be a negative number")
	}

	if port > 65535 {
		return nil, fmt.Errorf("'port' value should between 0 to 65535")
	}

	if target == "" {
		return nil, fmt.Errorf("'target' must not be empty")
	}

	recordSRV := NewRecordSRV(RecordSRV{
		View:     dnsView,
		Name:     name,
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
	return recordSRV, nil

}

func (objMgr *ObjectManager) GetSRVRecord(dnsView string, name string) (*[]RecordSRV, error) {
	if dnsView == "" || name == "" {
		return nil, fmt.Errorf("'DNS view' and 'name' are required to retrieve a unique srv record")
	}
	var res []RecordSRV

	recordSRV := NewRecordSRV(RecordSRV{})

	sf := map[string]string{
		"view": dnsView,
		"name": name,
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
				name, dnsView))
	}

	return &res, nil
}

func (objMgr *ObjectManager) GetSRVRecordByRef(ref string) (*RecordSRV, error) {
	recordSRV := NewRecordSRV(RecordSRV{})
	err := objMgr.connector.GetObject(recordSRV, ref, NewQueryParams(false, nil), &recordSRV)

	return recordSRV, err
}

func (objMgr *ObjectManager) UpdateSRVRecord(
	ref string,
	name string,
	priority int,
	weight int,
	port uint32,
	target string,
	ttl uint32,
	useTtl bool,
	comment string,
	eas EA) (*RecordSRV, error) {

	_, err := objMgr.GetSRVRecordByRef(ref)
	nameRegex := `^_[a-z]+\._[a-z]+\.[a-z0-9-]+\.[a-z]+$`
	targetRegex := `^[a-z]+\.[a-z0-9-]+\.[a-z]+$`

	valid, _ := regexp.MatchString(nameRegex, name)
	valid_tg, _ := regexp.MatchString(targetRegex, name)

	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, fmt.Errorf("'name' format is not valid")
	}

	if !valid_tg {
		return nil, fmt.Errorf("'target' is not in valid format")
	}

	if priority < 0 {
		return nil, fmt.Errorf("'priority' field must not be a negative number")
	}

	if port > 65535 {
		return nil, fmt.Errorf("'port' value should between 0 to 65535")
	}

	if weight < 0 {
		return nil, fmt.Errorf("'weight' must not be a negative number")
	}

	recordSRV := NewRecordSRV(RecordSRV{
		Name:     name,
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

	return recordSRV, nil
}

func (objMgr *ObjectManager) DeleteSRVRecord(ref string) (string, error) {
	return objMgr.connector.DeleteObject(ref)
}
