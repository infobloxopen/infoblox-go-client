package ibclient

type RecordMXOperations interface {
	CreateMXRecord(recMX RecordMX) (*RecordMX, error)
	GetMXRecord(recMX RecordMX) (*[]RecordMX, error)
	DeleteMXRecord(recMX RecordMX) (string, error)
	UpdateMXRecord(recMX RecordMX) (*RecordMX, error)
}

type RecordMX struct {
	IBBase            `json:"-"`
	Ref               string `json:"_ref,omitempty"`
	MailExchanger     string `json:"mail_exchanger,omitempty"`
	Preference        uint32 `json:"preference,omitempty"`
	Name              string `json:"name,omitempty"`
	View              string `json:"view,omitempty"`
	Zone              string `json:"zone,omitempty"`
	Ea                EA     `json:"extattrs,omitempty"`
	AddEA             EA     `json:"omitempty"`
	RemoveEA          EA     `json:"omitempty"`
	CreationTime      int    `json:"creation_time,omitempty"`
	Comment           string `json:"comment,omitempty"`
	Creator           string `json:"creator,omitempty"`
	DdnsProtected     bool   `json:"ddns_protected,omitempty"`
	DnsName           string `json:"dns_name,omitempty"`
	ForbidReclamation bool   `json:"forbid_reclamation,omitempty"`
	Reclaimable       bool   `json:"reclaimable,omitempty"`
	Ttl               uint   `json:"ttl,omitempty"`
	UseTtl            bool   `json:"use_ttl,omitempty"`
}

// NewRecordMX creates a new MX Record type with objectType and returnFields
func NewRecordMX(ra RecordMX) *RecordMX {
	res := ra
	res.objectType = "record:mx"

	res.returnFields = []string{"mail_exchanger", "preference", "name", "view", "zone", "extattrs", "comment", "creation_time",
		"ddns_protected", "dns_name", "forbid_reclamation", "reclaimable", "ttl", "use_ttl", "disable"}
	return &res
}

// CreateMXRecord takes Name, MailExchanger, Preference and View of the record to create MX Record
func (objMgr *ObjectManager) CreateMXRecord(recMX RecordMX) (*RecordMX, error) {
	recMX.Ea = objMgr.extendEA(recMX.Ea)
	recordMX := NewRecordMX(recMX)
	ref, err := objMgr.connector.CreateObject(recordMX)
	recordMX.Ref = ref
	return recordMX, err
}

// GetMXRecord by passing Name, Reference ID, MailExchanger, Preference or DNS View
// If no arguments are passed then, all the records are returned
func (objMgr *ObjectManager) GetMXRecord(recMX RecordMX) (*[]RecordMX, error) {

	var res []RecordMX
	recordMX := NewRecordMX(recMX)
	var err error
	if len(recMX.Ref) > 0 {
		err = objMgr.connector.GetObject(recordMX, recMX.Ref, &recordMX)
		res = append(res, *recordMX)

	} else {
		recordMX = NewRecordMX(recMX)
		err = objMgr.connector.GetObject(recordMX, "", &res)
		if err != nil || res == nil || len(res) == 0 {
			return nil, err
		}
	}

	return &res, err
}

// DeleteMXRecord by passing either Reference or Name or MailExchanger or Preference
// If a record with same MailExchanger or Preference and different name exists, then name and MailExchanger or Preference has to be passed
// to avoid multiple record deletions
func (objMgr *ObjectManager) DeleteMXRecord(recMX RecordMX) (string, error) {
	var res []RecordMX
	recordName := NewRecordMX(recMX)
	if len(recMX.Ref) > 0 {
		return objMgr.connector.DeleteObject(recMX.Ref)

	} else {
		recordName = NewRecordMX(recMX)
		err := objMgr.connector.GetObject(recordName, "", &res)
		if err != nil || res == nil || len(res) == 0 {
			return "Record doesn't exist", err
		}
		return objMgr.connector.DeleteObject(res[0].Ref)
	}

}

// UpdateMXRecord takes Reference ID of the record as an argument
// to update Name, MailExchanger and Preference of the record
// returns updated Refernce ID
func (objMgr *ObjectManager) UpdateMXRecord(recMX RecordMX) (*RecordMX, error) {
	var res RecordMX
	recordMX := RecordMX{}
	recordMX.returnFields = []string{"name", "mail_exchanger", "preference", "extattrs"}
	err := objMgr.connector.GetObject(&recordMX, recMX.Ref, &res)
	if err != nil {
		return nil, err
	}
	res.Name = recMX.Name
	res.MailExchanger = recMX.MailExchanger
	res.Preference = recMX.Preference
	reference, err := objMgr.connector.UpdateObject(&res, recMX.Ref)
	res.Ref = reference
	return &res, err
}
