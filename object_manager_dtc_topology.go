package ibclient

import (
	"encoding/json"
	"fmt"
)

// helpers for DtcTopology

// find all dtc:topology:rule objects belonging to the given dtc:topology by their reference to it
func getDtcRules(topologyRef string, objMgr *ObjectManager) ([]*DtcTopologyRule, error) {
	var dtcRules []*DtcTopologyRule

	searchFields := map[string]string{
		"topology": topologyRef,
	}
	objMgr.connector.GetObject(&DtcTopologyRule{}, "", NewQueryParams(false, searchFields), &dtcRules)

	return dtcRules, nil
}

// DtcTopology implementation

func (d *DtcTopology) MarshalJSON() ([]byte, error) {
	// correctly encode if rules is empty  (NOTE: copied from DtcServer)
	type Alias DtcTopology
	aux := &struct {
		Rules []*DtcTopologyRule `json:"rules"`
		*Alias
	}{
		Alias: (*Alias)(d),
	}

	if len(d.Rules) == 0 {
		aux.Rules = []*DtcTopologyRule{}
	} else {
		aux.Rules = d.Rules
	}

	return json.Marshal(aux)
}

func NewEmptyDtcTopology() *DtcTopology {
	dtcTopology := &DtcTopology{}
	dtcTopology.SetReturnFields(append(dtcTopology.ReturnFields(), "extattrs",
		"rules.dest_type", "rules.destination_link", "rules.return_type", "rules.sources", "rules.valid"))
	return dtcTopology
}

func NewDtcTopology(
	comment string,
	name string,
	rules []*DtcTopologyRule,
	ea EA,
) *DtcTopology {
	DtcTopology := NewEmptyDtcTopology()
	DtcTopology.Comment = &comment
	DtcTopology.Name = &name
	DtcTopology.Rules = rules
	DtcTopology.Ea = ea
	return DtcTopology
}

func (objMgr *ObjectManager) CreateDtcTopology(
	comment string,
	name string,
	dtcRules []*DtcTopologyRule, // rules []TopologyRule,
	ea EA,
) (*DtcTopology, error) {
	if name == "" {
		return nil, fmt.Errorf("name field is required to create a Dtc Topology object")
	}

	dtcTopology := NewDtcTopology(comment, name, dtcRules, ea)
	ref, err := objMgr.connector.CreateObject(dtcTopology)
	if err != nil {
		return nil, err
	}
	dtcTopology.Ref = ref
	return dtcTopology, nil
}

func (objMgr *ObjectManager) GetAllDtcTopology(queryParams *QueryParams) ([]DtcTopology, error) {
	var res []DtcTopology
	topology := NewEmptyDtcTopology()
	err := objMgr.connector.GetObject(topology, "", queryParams, &res)
	if err != nil {
		return nil, fmt.Errorf("error getting Dtc Topology object, err: %s", err)
	}
	return res, nil
}

func (objMgr *ObjectManager) GetDtcTopology(name string) (*DtcTopology, error) {
	var res []DtcTopology
	topology := NewEmptyDtcTopology()
	if name == "" {
		return nil, fmt.Errorf("name of the topology is required to retreive a unique dtc topology")
	}
	sf := map[string]string{
		"name": name,
	}
	queryParams := NewQueryParams(false, sf)
	err := objMgr.connector.GetObject(topology, "", queryParams, &res)
	if err != nil {
		return nil, err
	} else if res == nil || len(res) == 0 {
		return nil, NewNotFoundError(
			fmt.Sprintf("Dtc topology with name '%s' not found", name))
	}
	return &res[0], nil
}

func (objMgr *ObjectManager) UpdateDtcTopology(
	ref string,
	comment string,
	name string,
	ea EA,
	dtcRules []*DtcTopologyRule, // rules []TopologyRule,
) (*DtcTopology, error) {
	dtcTopology := NewDtcTopology(comment, name, dtcRules, ea)
	dtcTopology.Ref = ref
	ref, err := objMgr.connector.UpdateObject(dtcTopology, ref)
	if err != nil {
		return nil, err
	}
	dtcTopology.Ref = ref
	return dtcTopology, nil
}

func (objMgr *ObjectManager) GetDtcTopologyByRef(ref string) (*DtcTopology, error) {
	topologyDtc := NewEmptyDtcTopology()
	err := objMgr.connector.GetObject(
		topologyDtc, ref, NewQueryParams(false, nil), &topologyDtc)
	return topologyDtc, err
}

func (objMgr *ObjectManager) DeleteDtcTopology(ref string) (string, error) {
	return objMgr.connector.DeleteObject(ref)
}
