package ibclient

type RecordSRVOperations interface {
	CreateSRVRecord(recSRV RecordSRV) (*RecordSRV, error)
	GetSRVRecord(recsSRV RecordSRV) (*[]RecordSRV, error)
	UpdateSRVRecord(recSRV RecordSRV) (*RecordSRV, error)
	DeleteSRVRecord(recSRV RecordSRV) (string, error)
}

type RecordSRV struct {
	IBBase            `json:"-"`
	Ref               string `json:"_ref,omitempty"`
	Port              uint   `json:"port,omitempty"`
	Priority          uint   `json:"priority,omitempty"`
	Target            string `json:"target,omitempty"`
	Weight            uint   `json:"weight,omitempty"`
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

// NewRecordSRV creates a new SRV Record type with objectType and returnFields
func NewRecordSRV(ra RecordSRV) *RecordSRV {
	res := ra
	res.objectType = "record:srv"

	res.returnFields = []string{"port", "priority", "target", "weight", "name", "view", "zone", "extattrs", "comment", "creation_time",
		"ddns_protected", "dns_name", "forbid_reclamation", "reclaimable", "ttl", "use_ttl", "disable"}
	return &res
}

// CreateSRVRecord takes Name, Port, Priority, Target, Weight and View of the record to create SRV Record
// Optional fields: Ea
func (objMgr *ObjectManager) CreateSRVRecord(recSRV RecordSRV) (*RecordSRV, error) {
	recSRV.Ea = objMgr.extendEA(recSRV.Ea)
	recordSRV := NewRecordSRV(recSRV)
	ref, err := objMgr.connector.CreateObject(recordSRV)
	recordSRV.Ref = ref
	return recordSRV, err
}

// GetSRVRecord by passing Name, reference ID, Port, Priority, Target, Weight or DNS View
// If no arguments are passed then, all the records are returned
func (objMgr *ObjectManager) GetSRVRecord(recSRV RecordSRV) (*[]RecordSRV, error) {

	var res []RecordSRV
	recordSRV := NewRecordSRV(recSRV)
	var err error
	if len(recSRV.Ref) > 0 {
		err = objMgr.connector.GetObject(recordSRV, recSRV.Ref, &recordSRV)
		res = append(res, *recordSRV)

	} else {
		err = objMgr.connector.GetObject(recordSRV, "", &res)
		if err != nil || res == nil || len(res) == 0 {
			return nil, err
		}
	}

	return &res, err
}

// DeleteSRVRecord by passing either Reference or Name or Port or Priority or Target or Weight
// If a record with same Port, Priority, Target, Weight and different name exists, then name and Port, Priority, Target, Weight has to be passed
// to avoid multiple record deletions
func (objMgr *ObjectManager) DeleteSRVRecord(recSRV RecordSRV) (string, error) {
	var res []RecordSRV
	recordName := NewRecordSRV(recSRV)
	if len(recSRV.Ref) > 0 {
		return objMgr.connector.DeleteObject(recSRV.Ref)

	} else {
		err := objMgr.connector.GetObject(recordName, "", &res)
		if err != nil || res == nil || len(res) == 0 {
			return "Record doesn't exist", err
		}
		return objMgr.connector.DeleteObject(res[0].Ref)
	}

}

// UpdateSRVRecord takes Reference ID of the record as an argument
// to update Name, Port, Priority, Target and Weight of the record
// returns updated Refernce ID
func (objMgr *ObjectManager) UpdateSRVRecord(recSRV RecordSRV) (*RecordSRV, error) {
	var res RecordSRV
	recordSRV := RecordSRV{}
	recordSRV.returnFields = []string{"name", "extattrs"}
	err := objMgr.connector.GetObject(&recordSRV, recSRV.Ref, &res)
	if err != nil {
		return nil, err
	}
	res.Name = recSRV.Name
	res.Port = recSRV.Port
	res.Priority = recSRV.Priority
	res.Target = recSRV.Target
	res.Weight = recSRV.Weight
	reference, err := objMgr.connector.UpdateObject(&res, recSRV.Ref)
	res.Ref = reference
	return &res, err
}
