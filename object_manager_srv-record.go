package ibclient

import (
	"fmt"
	"regexp"
	"strings"
)

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

	targetRegex := `^[a-z]+\.[a-z0-9-]+\.[a-z]+$`
	valid_tg, _ := regexp.MatchString(targetRegex, target)

	nameSplit := strings.SplitN(name, ".", 3)

	if len(nameSplit) < 3 {
		return nil, fmt.Errorf("SRV Record format: _service._proto.domainName")
	} else {
		serviceRegex := `^_[a-z]+$`
		validService, _ := regexp.MatchString(serviceRegex, nameSplit[0])

		protocolRegex := `^_[a-z0-9-]+$`
		validProtocol, _ := regexp.MatchString(protocolRegex, nameSplit[1])

		domainRegexp := regexp.MustCompile(`^(?i)[a-z0-9-]+(\.[a-z0-9-]+)+\.?$`)
		validDomainName := domainRegexp.MatchString(nameSplit[2])

		if !(validService && validProtocol && validDomainName) {
			return nil, fmt.Errorf("name is not in valid format")
		}
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

	if priority < 0 || priority > 65535 {
		return nil, fmt.Errorf("'priority' value must be in range(0-65535)")
	}

	if port > 65535 {
		return nil, fmt.Errorf("'port' value should between 0 to 65535")
	}

	if weight < 0 || weight > 65535 {
		return nil, fmt.Errorf("'weight' value should in range(0-65535)")
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

func (objMgr *ObjectManager) GetSRVRecord(dnsView string, name string) (*RecordSRV, error) {
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

	_, err := objMgr.GetSRVRecordByRef(ref)
	targetRegex := `^[a-z]+\.[a-z0-9-]+\.[a-z]+$`
	valid_tg, _ := regexp.MatchString(targetRegex, target)

	if err != nil {
		return nil, err
	}

	nameSplit := strings.SplitN(name, ".", 3)

	if len(nameSplit) < 3 {
		return nil, fmt.Errorf("SRV Record format: _service._proto.domainName")
	} else {
		serviceRegex := `^_[a-z]+$`
		validService, _ := regexp.MatchString(serviceRegex, nameSplit[0])

		protocolRegex := `^_[a-z0-9-]+$`
		validProtocol, _ := regexp.MatchString(protocolRegex, nameSplit[1])

		domainRegexp := regexp.MustCompile(`^(?i)[a-z0-9-]+(\.[a-z0-9-]+)+\.?$`)
		validDomainName := domainRegexp.MatchString(nameSplit[2])

		if !(validService && validProtocol && validDomainName) {
			return nil, fmt.Errorf("name is not in valid format")
		}
	}

	if !valid_tg {
		return nil, fmt.Errorf("'target' is not in valid format")
	}

	if priority < 0 || priority > 65535 {
		return nil, fmt.Errorf("priority' value must be in range(0-65535)")
	}

	if port > 65535 {
		return nil, fmt.Errorf("'port' value should between 0 to 65535")
	}

	if weight < 0 || weight > 65535 {
		return nil, fmt.Errorf("'weight' value must be in range(0-65535)")
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
