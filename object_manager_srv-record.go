package ibclient

import (
	"fmt"
)

func validateSrvRecArgs(
	name string,
	priority uint32,
	weight uint32,
	port uint32,
	target string) error {

	if err := ValidateSrvRecName(name); err != nil {
		return err
	}

	if name == "" {
		return fmt.Errorf("'name' must not be empty")
	}

	if priority < 0 || priority > 65535 {
		return fmt.Errorf("'priority' value must be in range(0-65535)")
	}

	if port > 65535 {
		return fmt.Errorf("'port' value should between 0 to 65535")
	}

	if weight < 0 || weight > 65535 {
		return fmt.Errorf("'weight' value should in range(0-65535)")
	}

	if target == "" {
		return fmt.Errorf("'target' must not be empty")
	}

	return nil
}

func (objMgr *ObjectManager) CreateSRVRecord(
	dnsView string,
	name string,
	priority uint32,
	weight uint32,
	port uint32,
	target string,
	ttl uint32,
	useTtl bool,
	comment string,
	eas EA) (*RecordSRV, error) {

	err := validateSrvRecArgs(
		name,
		priority,
		weight,
		port,
		target)
	if err != nil {
		return nil, err
	}

	if dnsView == "" {
		dnsView = "default"
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

func (objMgr *ObjectManager) GetSRVRecord(dnsView string, name string, priority uint32, weight uint32) (*RecordSRV, error) {
	if dnsView == "" || name == "" {
		return nil, fmt.Errorf("'DNS view' and 'name' are required to retrieve a unique srv record")
	}
	var res []RecordSRV

	recordSRV := NewEmptyRecordSRV()

	sf := map[string]string{
		"view":     dnsView,
		"name":     name,
		"priority": fmt.Sprintf("%d", priority),
		"weight":   fmt.Sprintf("%d", weight),
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

	return &res[0], nil
}

func (objMgr *ObjectManager) GetSRVRecordByRef(ref string) (*RecordSRV, error) {
	recordSRV := NewRecordSRV(RecordSRV{})
	err := objMgr.connector.GetObject(recordSRV, ref, NewQueryParams(false, nil), &recordSRV)

	return recordSRV, err
}

func (objMgr *ObjectManager) UpdateSRVRecord(
	ref string,
	name string,
	priority uint32,
	weight uint32,
	port uint32,
	target string,
	ttl uint32,
	useTtl bool,
	comment string,
	eas EA) (*RecordSRV, error) {

	err := validateSrvRecArgs(
		name,
		priority,
		weight,
		port,
		target)
	if err != nil {
		return nil, err
	}

	recordSRV := NewRecordSRV(RecordSRV{
		Ref:      ref,
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
