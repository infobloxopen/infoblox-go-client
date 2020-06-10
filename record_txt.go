package ibclient

type RecordTXTOperations interface {
	CreateTXTRecord(recTXT RecordTXT) (*RecordTXT, error)
	GetTXTRecord(recTXT RecordTXT) (*[]RecordTXT, error)
	DeleteTXTRecord(recTXT RecordTXT) (string, error)
	UpdateTXTRecord(recTXT RecordTXT) (*RecordTXT, error)
}

type RecordTXT struct {
	IBBase            `json:"-"`
	Ref               string `json:"_ref,omitempty"`
	Name              string `json:"name,omitempty"`
	Text              string `json:"text,omitempty"`
	TTL               int    `json:"ttl,omitempty"`
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
	UseTtl            bool   `json:"use_ttl,omitempty"`
}

// NewRecordTXT creates a new TXT Record type with objectType and returnFields
func NewRecordTXT(rt RecordTXT) *RecordTXT {
	res := rt
	res.objectType = "record:txt"
	res.returnFields = []string{"name", "text", "view", "zone", "extattrs", "comment", "creation_time",
		"creator", "ddns_protected", "dns_name", "forbid_reclamation", "reclaimable", "ttl", "use_ttl"}

	return &res
}

// CreateTXTRecord takes Name, Text, TTL and View of the record to create TXT Record
// Use TTL of 0 to inherit TTL from the Zone
func (objMgr *ObjectManager) CreateTXTRecord(recTXT RecordTXT) (*RecordTXT, error) {
	recTXT.Ea = objMgr.extendEA(recTXT.Ea)
	recordTXT := NewRecordTXT(recTXT)

	ref, err := objMgr.connector.CreateObject(recordTXT)
	recordTXT.Ref = ref
	return recordTXT, err
}

// GetTXTRecord by passing Name, reference ID, TXT, TTL or DNS View
// If no arguments are passed then, all the TXT records are returned
func (objMgr *ObjectManager) GetTXTRecord(recTXT RecordTXT) (*[]RecordTXT, error) {

	var res []RecordTXT
	recordTXT := NewRecordTXT(recTXT)
	var err error
	if len(recTXT.Ref) > 0 {
		err = objMgr.connector.GetObject(recordTXT, recTXT.Ref, &recordTXT)
		res = append(res, *recordTXT)

	} else {
		recordTXT = NewRecordTXT(recTXT)
		err = objMgr.connector.GetObject(recordTXT, "", &res)
		if err != nil || res == nil || len(res) == 0 {
			return nil, err
		}
	}

	return &res, err
}

// DeleteTXTRecord by passing either Reference or Name or Text
// If a record with same Text and different Name exists, then Name and Text has to be passed
// to avoid multiple record deletions
func (objMgr *ObjectManager) DeleteTXTRecord(recTXT RecordTXT) (string, error) {
	var res []RecordTXT
	recordName := NewRecordTXT(recTXT)
	if len(recTXT.Ref) > 0 {
		return objMgr.connector.DeleteObject(recTXT.Ref)

	} else {
		recordName = NewRecordTXT(recTXT)
		err := objMgr.connector.GetObject(recordName, "", &res)
		if err != nil || res == nil || len(res) == 0 {
			return "Record doesn't exist", err
		}
		return objMgr.connector.DeleteObject(res[0].Ref)
	}
}

// UpdateTXTRecord takes Reference ID of the record as an argument
// to update Name, Text and EAs of the record
// returns updated Refernce ID
func (objMgr *ObjectManager) UpdateTXTRecord(recTXT RecordTXT) (*RecordTXT, error) {
	var res RecordTXT

	recordTXT := RecordTXT{}
	recordTXT.returnFields = []string{"name", "text", "ttl", "extattrs"}
	err := objMgr.connector.GetObject(&recordTXT, recTXT.Ref, &res)
	if err != nil {
		return nil, err
	}
	res.Name = recTXT.Name
	res.Text = recTXT.Text
	res.TTL = recTXT.TTL
	res.Zone = "" //  set the Zone value to "" as its a non writable field
	reference, err := objMgr.connector.UpdateObject(&res, recTXT.Ref)
	res.Ref = reference
	return &res, err

}
