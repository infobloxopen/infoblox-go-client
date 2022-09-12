package ibclient

import "fmt"

func (objMgr *ObjectManager) CreateZoneAuth(
	fqdn string,
	nsGroup string,
	restartIfNeeded bool,
	comment string,
	soaDefaultTtl int,
	soaExpire int,
	soaNegativeTtl int,
	soaRefresh int,
	soaRetry int,
	eas EA) (*ZoneAuth, error) {

	zoneAuth := NewZoneAuth(ZoneAuth{
		Fqdn:            fqdn,
		NsGroup:         nsGroup,
		RestartIfNeeded: restartIfNeeded,
		Comment:         comment,
		SoaDefaultTtl:   soaDefaultTtl,
		SoaExpire:       soaExpire,
		SoaNegativeTtl:  soaNegativeTtl,
		SoaRefresh:      soaRefresh,
		SoaRetry:        soaRetry,
		Ea:              eas})

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
