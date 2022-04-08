package ibclient

func (objMgr *ObjectManager) CreateNSRecord(dnsview string, recordname string, ns string, addr []string) (*RecordNS, error) {

	var addresses []ZoneNameServer
	for _, addr := range addr {
		addresses = append(addresses, ZoneNameServer{Address: addr})
	}
	recordNS := NewRecordNS(RecordNS{
		View:       dnsview,
		Name:       recordname,
		NS:         ns,
		Addresses:  addresses,
	})

	ref, err := objMgr.connector.CreateObject(recordNS)
	recordNS.Ref = ref
	return recordNS, err
}

func (objMgr *ObjectManager) DeleteNSRecord(ref string) (string, error) {
	return objMgr.connector.DeleteObject(ref)
}
