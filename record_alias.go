package ibclient

type RecordAliasOperations interface {
	CreateAliasRecord(recAlias RecordAlias) (*RecordAlias, error)
	GetALiasRecord(recAlias RecordAlias) (*[]RecordAlias, error)
	DeleteAliasRecord(recAlias RecordAlias) (string, error)
	UpdateAliasRecord(recAlias RecordAlias) (*RecordAlias, error)
}

type RecordAlias struct {
	IBBase       `json:"-"`
	Ref          string `json:"_ref,omitempty"`
	TargetName   string `json:"target_name,omitempty"`
	TargetType   string `json:"target_type,omitempty"`
	Name         string `json:"name,omitempty"`
	View         string `json:"view,omitempty"`
	Zone         string `json:"zone,omitempty"`
	Ea           EA     `json:"extattrs,omitempty"`
	AddEA        EA     `json:"omitempty"`
	RemoveEA     EA     `json:"omitempty"`
	Comment      string `json:"comment,omitempty"`
	Creator      string `json:"creator,omitempty"`
	DnsName      string `json:"dns_name,omitempty"`
	Ttl          uint   `json:"ttl,omitempty"`
	UseTtl       bool   `json:"use_ttl,omitempty"`
}

// NewRecordAlias creates a new Alias Record type with objectType and returnFields
func NewRecordAlias(ra RecordAlias) *RecordAlias {
	res := ra
	res.objectType = "record:alias"

	res.returnFields = []string{"target_name", "target_type", "name", "view", "zone", "extattrs", "comment",
		"dns_name", "ttl", "use_ttl"}
	return &res
}

// CreateAliasRecord takes Name, TargetName, TargetType and View of the record to create Alias Record
func (objMgr *ObjectManager) CreateAliasRecord(recAlias RecordAlias) (*RecordAlias, error) {
	recAlias.Ea = objMgr.extendEA(recAlias.Ea)
	recordAlias := NewRecordAlias(recAlias)
	ref, err := objMgr.connector.CreateObject(recordAlias)
	recordAlias.Ref = ref
	return recordAlias, err
}

// GetAliasRecord by passing Name, reference ID, TargetName, TargetType or DNS View
// If no arguments are passed then, all the records are returned
func (objMgr *ObjectManager) GetAliasRecord(recAlias RecordAlias) (*[]RecordAlias, error) {

	var res []RecordAlias
	recordAlias := NewRecordAlias(recAlias)
	var err error
	if len(recAlias.Ref) > 0 {
		err = objMgr.connector.GetObject(recordAlias, recAlias.Ref, &recordAlias)
		res = append(res, *recordAlias)

	} else {
		recordAlias = NewRecordAlias(recAlias)
		err = objMgr.connector.GetObject(recordAlias, "", &res)
		if err != nil || res == nil || len(res) == 0 {
			return nil, err
		}
	}

	return &res, err
}

// DeleteAliasRecord by passing either Reference or Name or TargetName or TargetType
// If a record with same TargetName or TargetType and different Name exists, then name and TargetName or TargetType has to be passed
// to avoid multiple record deletions
func (objMgr *ObjectManager) DeleteAliasRecord(recAlias RecordAlias) (string, error) {
	var res []RecordAlias
	recordName := NewRecordAlias(recAlias)
	if len(recAlias.Ref) > 0 {
		return objMgr.connector.DeleteObject(recAlias.Ref)

	} else {
		recordName = NewRecordAlias(recAlias)
		err := objMgr.connector.GetObject(recordName, "", &res)
		if err != nil || res == nil || len(res) == 0 {
			return "", err
		}
		return objMgr.connector.DeleteObject(res[0].Ref)
	}

}

// UpdateAliasRecord takes Reference ID of the record as an argument
// to update Name, TargetName and TargetType of the record
// returns updated Refernce ID
func (objMgr *ObjectManager) UpdateAliasRecord(recAlias RecordAlias) (*RecordAlias, error) {
	var res RecordAlias
	recordAlias := RecordAlias{}
	recordAlias.returnFields = []string{"name", "target_name", "target_type", "extattrs"}
	err := objMgr.connector.GetObject(&recordAlias, recAlias.Ref, &res)
	if err != nil {
		return nil, err
	}
	res.Name = recAlias.Name
	res.TargetName = recAlias.TargetName
	res.TargetType = recAlias.TargetType
	reference, err := objMgr.connector.UpdateObject(&res, recAlias.Ref)
	res.Ref = reference
	return &res, err
}
