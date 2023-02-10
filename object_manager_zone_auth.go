package ibclient

import (
	"fmt"
)

func (objMgr *ObjectManager) CreateZoneAuth(
	dnsview string,
	fqdn string,
	nsGroup string,
	restartIfNeeded bool,
	comment string,
	soaDefaultTtl int,
	soaExpire int,
	soaNegativeTtl int,
	soaRefresh int,
	soaRetry int,
	zoneFormat string,
	ea EA) (*ZoneAuth, error) {

	zoneAuth := NewZoneAuth(ZoneAuth{
		View:            dnsview,
		Fqdn:            fqdn,
		NsGroup:         nsGroup,
		RestartIfNeeded: restartIfNeeded,
		Comment:         comment,
		SoaDefaultTtl:   soaDefaultTtl,
		SoaExpire:       soaExpire,
		SoaNegativeTtl:  soaNegativeTtl,
		SoaRefresh:      soaRefresh,
		SoaRetry:        soaRetry,
		ZoneFormat:      zoneFormat,
		Ea:              ea})

	ref, err := objMgr.connector.CreateObject(zoneAuth)
	zoneAuth.Ref = ref
	return zoneAuth, err
}

// Retrieve a authoritative zone by ref
func (objMgr *ObjectManager) GetZoneAuthByRef(ref string) (*ZoneAuth, error) {
	res := NewZoneAuth(ZoneAuth{})

	if ref == "" {
		return nil, fmt.Errorf("empty reference to an object is not allowed")
	}

	err := objMgr.connector.GetObject(
		res, ref, NewQueryParams(false, nil), res)
	return res, err
}

// UpdateZoneAuth updates an auth zone, except fields that can't be updated
func (objMgr *ObjectManager) UpdateZoneAuth(
	ref string,
	dnsview string,
	// fqdn string, CANNOT BE UPDATED IN WAPI
	nsGroup string,
	restartIfNeeded bool,
	comment string,
	soaDefaultTtl int,
	soaExpire int,
	soaNegativeTtl int,
	soaRefresh int,
	soaRetry int,
	// zoneFormat string, CANNOT BE UPDATED IN WAPI
	ea EA) (*ZoneAuth, error) {

	inputZoneAuth := NewZoneAuth(ZoneAuth{
		Ref:             ref,
		View:            dnsview,
		NsGroup:         nsGroup,
		RestartIfNeeded: restartIfNeeded,
		Comment:         comment,
		SoaDefaultTtl:   soaDefaultTtl,
		SoaExpire:       soaExpire,
		SoaNegativeTtl:  soaNegativeTtl,
		SoaRefresh:      soaRefresh,
		SoaRetry:        soaRetry,
		Ea:              ea})

	updatedRef, err := objMgr.connector.UpdateObject(inputZoneAuth, ref)
	if err != nil {
		fmt.Printf("failed to update object with this ref '%s': %s", updatedRef, err)
		return nil, err
	}
	inputZoneAuth.Ref = updatedRef
	return inputZoneAuth, err
}

// DeleteZoneAuth deletes an auth zone
func (objMgr *ObjectManager) DeleteZoneAuth(ref string) (string, error) {
	return objMgr.connector.DeleteObject(ref)
}
