package ibclient

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Compile-time interface checks
var _ IBObjectManager = new(ObjectManager)

type IBObjectManager interface {
	GetDNSView(name string) (*DNSView, error)
	AllocateIP(netview string, cidr string, ipAddr string, isIPv6 bool, macOrDuid string, name string, comment string, eas EA) (*FixedAddress, error)
	AllocateNetwork(netview string, cidr string, isIPv6 bool, prefixLen uint, comment string, eas EA) (network *Network, err error)
	AllocateNetworkContainer(netview string, cidr string, isIPv6 bool, prefixLen uint, comment string, eas EA) (netContainer *NetworkContainer, err error)
	CreateAAAARecord(netView string, dnsView string, recordName string, cidr string, ipAddr string, useTtl bool, ttl uint32, comment string, eas EA) (*RecordAAAA, error)
	CreateARecord(netView string, dnsView string, name string, cidr string, ipAddr string, ttl uint32, useTTL bool, comment string, ea EA) (*RecordA, error)
	CreateCNAMERecord(dnsview string, canonical string, recordname string, useTtl bool, ttl uint32, comment string, eas EA) (*RecordCNAME, error)
	CreateDefaultNetviews(globalNetview string, localNetview string) (globalNetviewRef string, localNetviewRef string, err error)
	CreateEADefinition(eadef EADefinition) (*EADefinition, error)
	CreateHostRecord(enabledns bool, enabledhcp bool, recordName string, netview string, dnsview string, ipv4cidr string, ipv6cidr string, ipv4Addr string, ipv6Addr string, macAddr string, duid string, useTtl bool, ttl uint32, comment string, eas EA, aliases []string) (*HostRecord, error)
	CreateMXRecord(dnsView string, fqdn string, mx string, preference uint32, ttl uint32, useTtl bool, comment string, eas EA) (*RecordMX, error)
	CreateNetwork(netview string, cidr string, isIPv6 bool, comment string, eas EA) (*Network, error)
	CreateNetworkContainer(netview string, cidr string, isIPv6 bool, comment string, eas EA) (*NetworkContainer, error)
	CreateNetworkView(name string, comment string, setEas EA) (*NetworkView, error)
	CreateNsGroupDelegation(name string, delegate_to []NameServer) (*NsGroupDelegation, error)
	CreatePTRRecord(networkView string, dnsView string, ptrdname string, recordName string, cidr string, ipAddr string, useTtl bool, ttl uint32, comment string, eas EA) (*RecordPTR, error)
	CreateSRVRecord(dnsView string, name string, priority uint32, weight uint32, port uint32, target string, ttl uint32, useTtl bool, comment string, eas EA) (*RecordSRV, error)
	CreateTXTRecord(dnsView string, recordName string, text string, ttl uint32, useTtl bool, comment string, eas EA) (*RecordTXT, error)
	CreateZoneAuth(fqdn string, ea EA) (*ZoneAuth, error)
	CreateZoneDelegated(fqdn string, delegate_to []NameServer, NsGroup string) (*ZoneDelegated, error)
	DeleteAAAARecord(ref string) (string, error)
	DeleteARecord(ref string) (string, error)
	DeleteCNAMERecord(ref string) (string, error)
	DeleteFixedAddress(ref string) (string, error)
	DeleteHostRecord(ref string) (string, error)
	DeleteMXRecord(ref string) (string, error)
	DeleteNetwork(ref string) (string, error)
	DeleteNetworkContainer(ref string) (string, error)
	DeleteNetworkView(ref string) (string, error)
	DeleteNsGroupDelegation(ref string) (string, error)
	DeletePTRRecord(ref string) (string, error)
	DeleteSRVRecord(ref string) (string, error)
	DeleteTXTRecord(ref string) (string, error)
	DeleteZoneAuth(ref string) (string, error)
	DeleteZoneDelegated(ref string) (string, error)
	GetAAAARecord(dnsview string, recordName string, ipAddr string) (*RecordAAAA, error)
	GetAAAARecordByRef(ref string) (*RecordAAAA, error)
	GetARecord(dnsview string, recordName string, ipAddr string) (*RecordA, error)
	GetARecordByRef(ref string) (*RecordA, error)
	GetAllMembers() ([]Member, error)
	GetCNAMERecord(dnsview string, canonical string, recordName string) (*RecordCNAME, error)
	GetCNAMERecordByRef(ref string) (*RecordCNAME, error)
	GetCapacityReport(name string) ([]CapacityReport, error)
	GetEADefinition(name string) (*EADefinition, error)
	GetFixedAddress(netview string, cidr string, ipAddr string, isIPv6 bool, macOrDuid string) (*FixedAddress, error)
	GetFixedAddressByRef(ref string) (*FixedAddress, error)
	GetGridInfo() ([]Grid, error)
	GetGridLicense() ([]License, error)
	GetHostRecord(netview string, dnsview string, recordName string, ipv4addr string, ipv6addr string) (*HostRecord, error)
	GetHostRecordByRef(ref string) (*HostRecord, error)
	GetIpAddressFromHostRecord(host HostRecord) (string, error)
	GetMXRecord(dnsView string, fqdn string, mx string, preference uint32) (*RecordMX, error)
	GetMXRecordByRef(ref string) (*RecordMX, error)
	GetNetwork(netview string, cidr string, isIPv6 bool, ea EA) (*Network, error)
	GetNetworkByRef(ref string) (*Network, error)
	GetNetworkContainer(netview string, cidr string, isIPv6 bool, eaSearch EA) (*NetworkContainer, error)
	GetNetworkContainerByRef(ref string) (*NetworkContainer, error)
	GetNetworkView(name string) (*NetworkView, error)
	GetNetworkViewByRef(ref string) (*NetworkView, error)
	GetNsGroupDelegation(name string) (*NsGroupDelegation, error)
	GetPTRRecord(dnsview string, ptrdname string, recordName string, ipAddr string) (*RecordPTR, error)
	GetPTRRecordByRef(ref string) (*RecordPTR, error)
	GetSRVRecord(dnsView string, name string, target string, port uint32) (*RecordSRV, error)
	GetSRVRecordByRef(ref string) (*RecordSRV, error)
	GetTXTRecord(dnsview string, name string) (*RecordTXT, error)
	GetTXTRecordByRef(ref string) (*RecordTXT, error)
	GetUpgradeStatus(statusType string) ([]UpgradeStatus, error)
	GetZoneAuthByRef(ref string) (*ZoneAuth, error)
	GetZoneDelegated(fqdn string) (*ZoneDelegated, error)
	ReleaseIP(netview string, cidr string, ipAddr string, isIPv6 bool, macAddr string) (string, error)
	SearchHostRecordByAltId(internalId string, ref string, eaNameForInternalId string) (*HostRecord, error)
	UpdateAAAARecord(ref string, netView string, recordName string, cidr string, ipAddr string, useTtl bool, ttl uint32, comment string, setEas EA) (*RecordAAAA, error)
	UpdateARecord(ref string, name string, ipAddr string, cidr string, netview string, ttl uint32, useTTL bool, comment string, eas EA) (*RecordA, error)
	UpdateCNAMERecord(ref string, canonical string, recordName string, useTtl bool, ttl uint32, comment string, setEas EA) (*RecordCNAME, error)
	UpdateFixedAddress(fixedAddrRef string, netview string, name string, cidr string, ipAddr string, matchclient string, macOrDuid string, comment string, eas EA) (*FixedAddress, error)
	UpdateHostRecord(hostRref string, enabledns bool, enabledhcp bool, name string, netview string, dnsView string, ipv4cidr string, ipv6cidr string, ipv4Addr string, ipv6Addr string, macAddress string, duid string, useTtl bool, ttl uint32, comment string, eas EA, aliases []string) (*HostRecord, error)
	UpdateMXRecord(ref string, dnsView string, fqdn string, mx string, preference uint32, ttl uint32, useTtl bool, comment string, eas EA) (*RecordMX, error)
	UpdateNetwork(ref string, setEas EA, comment string) (*Network, error)
	UpdateNetworkContainer(ref string, setEas EA, comment string) (*NetworkContainer, error)
	UpdateNetworkView(ref string, name string, comment string, setEas EA) (*NetworkView, error)
	UpdateNsGroupDelegation(ref string, delegate_to []NameServer, ns_group string) (*NsGroupDelegation, error)
	UpdatePTRRecord(ref string, netview string, ptrdname string, name string, cidr string, ipAddr string, useTtl bool, ttl uint32, comment string, setEas EA) (*RecordPTR, error)
	UpdateSRVRecord(ref string, name string, priority uint32, weight uint32, port uint32, target string, ttl uint32, useTtl bool, comment string, eas EA) (*RecordSRV, error)
	UpdateTXTRecord(ref string, recordName string, text string, ttl uint32, useTtl bool, comment string, eas EA) (*RecordTXT, error)
	UpdateZoneDelegated(ref string, delegate_to []NameServer, ns_group string) (*ZoneDelegated, error)
}

type ObjectManager struct {
	connector IBConnector
	cmpType   string
	tenantID  string
}

func NewObjectManager(connector IBConnector, cmpType string, tenantID string) IBObjectManager {
	objMgr := &ObjectManager{}

	objMgr.connector = connector
	objMgr.cmpType = cmpType
	objMgr.tenantID = tenantID

	return objMgr
}

// CreateMultiObject unmarshals the result into slice of maps
func (objMgr *ObjectManager) CreateMultiObject(req *MultiRequest) ([]map[string]interface{}, error) {

	conn := objMgr.connector.(*Connector)
	queryParams := NewQueryParams(false, nil)
	res, err := conn.makeRequest(CREATE, req, "", queryParams)

	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	err = json.Unmarshal(res, &result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetUpgradeStatus returns the grid upgrade information
func (objMgr *ObjectManager) GetUpgradeStatus(statusType string) ([]UpgradeStatus, error) {
	var res []UpgradeStatus

	if statusType == "" {
		// TODO option may vary according to the WAPI version, need to
		// throw relevant  error.
		msg := fmt.Sprintf("Status type can not be nil")
		return res, errors.New(msg)
	}
	upgradestatus := NewUpgradeStatus(UpgradeStatus{})

	sf := map[string]string{
		"type": statusType,
	}
	queryParams := NewQueryParams(false, sf)
	err := objMgr.connector.GetObject(upgradestatus, "", queryParams, &res)

	return res, err
}

// GetAllMembers returns all members information
func (objMgr *ObjectManager) GetAllMembers() ([]Member, error) {
	var res []Member

	memberObj := NewMember(Member{})
	err := objMgr.connector.GetObject(
		memberObj, "", NewQueryParams(false, nil), &res)
	return res, err
}

// GetCapacityReport returns all capacity for members
func (objMgr *ObjectManager) GetCapacityReport(name string) ([]CapacityReport, error) {
	var res []CapacityReport

	capacityReport := NewCapcityReport(CapacityReport{})

	sf := map[string]string{
		"name": name,
	}
	queryParams := NewQueryParams(false, sf)
	err := objMgr.connector.GetObject(capacityReport, "", queryParams, &res)
	return res, err
}

// GetLicense returns the license details for member
func (objMgr *ObjectManager) GetLicense() ([]License, error) {
	var res []License

	licenseObj := NewLicense(License{})
	err := objMgr.connector.GetObject(
		licenseObj, "", NewQueryParams(false, nil), &res)
	return res, err
}

// GetLicense returns the license details for grid
func (objMgr *ObjectManager) GetGridLicense() ([]License, error) {
	var res []License

	licenseObj := NewGridLicense(License{})
	err := objMgr.connector.GetObject(
		licenseObj, "", NewQueryParams(false, nil), &res)
	return res, err
}

// GetGridInfo returns the details for grid
func (objMgr *ObjectManager) GetGridInfo() ([]Grid, error) {
	var res []Grid

	gridObj := NewGrid(Grid{})
	err := objMgr.connector.GetObject(
		gridObj, "", NewQueryParams(false, nil), &res)
	return res, err
}

// CreateZoneAuth creates zones and subs by passing fqdn
func (objMgr *ObjectManager) CreateZoneAuth(
	fqdn string,
	eas EA) (*ZoneAuth, error) {

	zoneAuth := NewZoneAuth(ZoneAuth{
		Fqdn: fqdn,
		Ea:   eas})

	ref, err := objMgr.connector.CreateObject(zoneAuth)
	zoneAuth.Ref = ref
	return zoneAuth, err
}

// Retreive a authortative zone by ref
func (objMgr *ObjectManager) GetZoneAuthByRef(ref string) (*ZoneAuth, error) {
	res := NewZoneAuth(ZoneAuth{})

	if ref == "" {
		return nil, fmt.Errorf("empty reference to an object is not allowed")
	}

	err := objMgr.connector.GetObject(
		res, ref, NewQueryParams(false, nil), res)
	return res, err
}

// DeleteZoneAuth deletes an auth zone
func (objMgr *ObjectManager) DeleteZoneAuth(ref string) (string, error) {
	return objMgr.connector.DeleteObject(ref)
}

// GetZoneAuth returns the authoritatives zones
func (objMgr *ObjectManager) GetZoneAuth() ([]ZoneAuth, error) {
	var res []ZoneAuth

	zoneAuth := NewZoneAuth(ZoneAuth{})
	err := objMgr.connector.GetObject(
		zoneAuth, "", NewQueryParams(false, nil), &res)

	return res, err
}

// GetZoneDelegated returns the delegated zone
func (objMgr *ObjectManager) GetZoneDelegated(fqdn string) (*ZoneDelegated, error) {
	if len(fqdn) == 0 {
		return nil, nil
	}
	var res []ZoneDelegated

	zoneDelegated := NewZoneDelegated(ZoneDelegated{})

	sf := map[string]string{
		"fqdn": fqdn,
	}
	queryParams := NewQueryParams(false, sf)
	err := objMgr.connector.GetObject(zoneDelegated, "", queryParams, &res)

	if err != nil || res == nil || len(res) == 0 {
		return nil, err
	}

	return &res[0], nil
}

// CreateZoneDelegated creates delegated zone
func (objMgr *ObjectManager) CreateZoneDelegated(fqdn string, delegate_to []NameServer, ns_group string) (*ZoneDelegated, error) {
	zoneDelegated := NewZoneDelegated(ZoneDelegated{
		Fqdn:       fqdn,
		DelegateTo: delegate_to,
		NsGroup:    ns_group,
	})

	ref, err := objMgr.connector.CreateObject(zoneDelegated)
	zoneDelegated.Ref = ref

	return zoneDelegated, err
}

// UpdateZoneDelegated updates delegated zone
func (objMgr *ObjectManager) UpdateZoneDelegated(ref string, delegate_to []NameServer, ns_group string) (*ZoneDelegated, error) {
	zoneDelegated := NewZoneDelegated(ZoneDelegated{
		Ref:        ref,
		DelegateTo: delegate_to,
		NsGroup:    ns_group,
	})

	refResp, err := objMgr.connector.UpdateObject(zoneDelegated, ref)
	zoneDelegated.Ref = refResp
	return zoneDelegated, err
}

// DeleteZoneDelegated deletes delegated zone
func (objMgr *ObjectManager) DeleteZoneDelegated(ref string) (string, error) {
	return objMgr.connector.DeleteObject(ref)
}

// GetNsGroupDelegation returns the Delegation NSGroup
func (objMgr *ObjectManager) GetNsGroupDelegation(name string) (*NsGroupDelegation, error) {
	if len(name) == 0 {
		return nil, nil
	}
	var res []NsGroupDelegation

	zoneDelegated := NewNsGroupDelegation(NsGroupDelegation{})

	sf := map[string]string{
		"name": name,
	}
	queryParams := NewQueryParams(false, sf)
	err := objMgr.connector.GetObject(zoneDelegated, "", queryParams, &res)

	if err != nil || res == nil || len(res) == 0 {
		return nil, err
	}

	return &res[0], nil
}

// CreateNsGroupDelegation creates delegated zone
func (objMgr *ObjectManager) CreateNsGroupDelegation(name string, delegate_to []NameServer) (*NsGroupDelegation, error) {
	zoneDelegated := NewNsGroupDelegation(NsGroupDelegation{
		Name:       name,
		DelegateTo: delegate_to,
	})

	ref, err := objMgr.connector.CreateObject(zoneDelegated)
	zoneDelegated.Ref = ref

	return zoneDelegated, err
}

// UpdateNsGroupDelegation updates delegated zone
func (objMgr *ObjectManager) UpdateNsGroupDelegation(ref string, delegate_to []NameServer, ns_group string) (*NsGroupDelegation, error) {
	zoneDelegated := NewNsGroupDelegation(NsGroupDelegation{
		Ref:        ref,
		DelegateTo: delegate_to,
	})

	refResp, err := objMgr.connector.UpdateObject(zoneDelegated, ref)
	zoneDelegated.Ref = refResp
	return zoneDelegated, err
}

// DeleteNsGroupDelegation deletes delegated zone
func (objMgr *ObjectManager) DeleteNsGroupDelegation(ref string) (string, error) {
	return objMgr.connector.DeleteObject(ref)
}
