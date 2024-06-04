package ibclient

import (
	"fmt"
)

func (objMgr *ObjectManager) CreateZoneForward(
	comment string,
	disable bool,
	eas EA,
	forwardTo []NameServer,
	forwardersOnly bool,
	forwardingServers []*Forwardingmemberserver,
	fqdn string,
	nsGroup string,
	view string,
	zoneFormat string) (*ZoneForward, error) {
	if fqdn == "" && forwardTo == nil {
		return nil, fmt.Errorf("FQDN and forwardTo fields are required to create a forward zone")
	}

	zoneForward := NewZoneForward(comment, disable, eas, forwardTo, forwardersOnly, forwardingServers, fqdn, nsGroup, view, zoneFormat)
	ref, err := objMgr.connector.CreateObject(zoneForward)
	if err != nil {
		return nil, err
	}
	zoneForward.Ref = ref
	return zoneForward, nil
}

func (objMgr *ObjectManager) DeleteZoneForward(ref string) (string, error) {
	return objMgr.connector.DeleteObject(ref)

}

func (objMgr *ObjectManager) GetZoneForwardByRef(ref string) (*ZoneForward, error) {
	zoneForward := NewEmptyZoneForward()
	err := objMgr.connector.GetObject(zoneForward, ref, NewQueryParams(false, nil), &zoneForward)
	if err != nil {
		return nil, err
	}
	return zoneForward, nil
}

func (objMgr *ObjectManager) GetZoneForwardFilters(filters map[string]string) ([]ZoneForward, error) {

	var res []ZoneForward
	zoneForward := NewEmptyZoneForward()

	err := objMgr.connector.GetObject(
		zoneForward, "", NewQueryParams(false, filters), &res)
	if err != nil {
		return nil, err
	}

	return res, err
}

func (objMgr *ObjectManager) UpdateZoneForward(
	ref string,
	comment string,
	disable bool,
	eas EA,
	forwardTo []NameServer,
	forwardersOnly bool,
	forwardingServers []*Forwardingmemberserver,
	nsGroup string) (*ZoneForward, error) {

	zoneForward := NewEmptyZoneForward()

	zoneForward.Comment = &comment
	zoneForward.Disable = &disable
	zoneForward.Ea = eas
	zoneForward.ForwardTo = forwardTo
	zoneForward.ForwardersOnly = &forwardersOnly
	zoneForward.ForwardingServers = forwardingServers
	if nsGroup != "" {
		zoneForward.NsGroup = &nsGroup
	} else {
		zoneForward.NsGroup = nil
	}

	new_ref, err := objMgr.connector.UpdateObject(zoneForward, ref)
	if err != nil {
		return nil, err
	}
	zoneForward.Ref = new_ref
	return zoneForward, nil

}

func NewEmptyZoneForward() *ZoneForward {
	zoneForward := &ZoneForward{}
	zoneForward.SetReturnFields(append(zoneForward.ReturnFields(), "zone_format", "ns_group", "comment", "disable", "extattrs", "forwarders_only", "forwarding_servers"))
	return zoneForward
}

func NewZoneForward(comment string,
	disable bool,
	eas EA,
	forwardTo []NameServer,
	forwardersOnly bool,
	forwardingServers []*Forwardingmemberserver,
	fqdn string,
	nsGroup string,
	view string,
	zoneFormat string) *ZoneForward {

	zoneForward := NewEmptyZoneForward()

	zoneForward.Comment = &comment
	zoneForward.Disable = &disable
	zoneForward.Ea = eas
	zoneForward.ForwardTo = forwardTo
	zoneForward.ForwardersOnly = &forwardersOnly
	zoneForward.ForwardingServers = forwardingServers
	zoneForward.Fqdn = fqdn
	if nsGroup == "" {
		zoneForward.NsGroup = nil
	} else {
		zoneForward.NsGroup = &nsGroup
	}
	if view == "" {
		view = "default"
	}
	zoneForward.View = &view
	if zoneFormat == "" {
		zoneFormat = "FORWARD"
	}
	zoneForward.ZoneFormat = zoneFormat

	return zoneForward
}
