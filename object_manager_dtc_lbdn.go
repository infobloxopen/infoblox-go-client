package ibclient

import (
	"encoding/json"
	"fmt"
)

func (d *DtcLbdn) MarshalJSON() ([]byte, error) {
	type Alias DtcLbdn
	aux := &struct {
		AuthZones []string       `json:"auth_zones"`
		Pools     []*DtcPoolLink `json:"pools"`
		Patterns  []string       `json:"patterns"`
		*Alias
	}{
		Alias: (*Alias)(d),
	}

	// Convert AuthZones to a slice of strings

	for _, zone := range d.AuthZones {
		if zone != nil {
			aux.AuthZones = append(aux.AuthZones, zone.Ref)
		}
	}

	// Convert Pools to a slice of DtcPoolLink
	for _, pool := range d.Pools {
		if pool != nil {
			aux.Pools = append(aux.Pools, pool)
		}
	}

	// Convert Patterns to a slice of strings
	for _, pattern := range d.Patterns {
		aux.Patterns = append(aux.Patterns, pattern)
	}

	// Ensure AuthZones, Pools, and Types are set to empty slices if nil
	if aux.AuthZones == nil {
		aux.AuthZones = make([]string, 0)
	}
	if aux.Pools == nil {
		aux.Pools = make([]*DtcPoolLink, 0)
	}
	if aux.Patterns == nil {
		aux.Patterns = make([]string, 0)
	}

	return json.Marshal(aux)
}

func (objMgr *ObjectManager) CreateDtcLbdn(name string, authzone []string, comment string, disable bool, autoConsolidatedMonitors bool, ea EA,
	lbMethod string, patterns []string, persistence uint32, pools []*DtcPoolLink, priority uint32, topology string, types []string, ttl uint32, usettl bool) (*DtcLbdn, error) {
	// todo: add health and status_member fields

	if name == "" || pools == nil || lbMethod == "" {
		return nil, fmt.Errorf("name, pools and lbMethod fields are required to create a DtcLbdn object")
	}
	// get ref id of authzones and replace
	var zones []*ZoneAuth
	var err error
	if len(authzone) > 0 {
		zones, err = getAuthZones(authzone, objMgr)
		if err != nil {
			return nil, err
		}
	}

	// get ref id of pools and replace
	dtcPoolLink, err := getPools(pools, objMgr)
	if err != nil {
		return nil, err
	}

	//get ref id of topology and replace
	topologyRef, err := getTopology(lbMethod, topology, objMgr)
	if err != nil {
		return nil, err
	}

	dtcLbdn := NewDtcLbdn("", name, zones, comment, disable, autoConsolidatedMonitors, ea,
		lbMethod, patterns, persistence, dtcPoolLink, priority, topologyRef, types, ttl, usettl)
	ref, err := objMgr.connector.CreateObject(dtcLbdn)
	if err != nil {
		return nil, fmt.Errorf("error creating DtcLbdn object %s, err: %s", name, err)
	}
	dtcLbdn.Ref = ref
	return dtcLbdn, nil
}

func getTopology(lbMethod string, topology string, objMgr *ObjectManager) (string, error) {
	var dtcTopology []DtcTopology
	var topologyRef string
	if lbMethod == "TOPOLOGY" {
		if topology == "" {
			return "", fmt.Errorf("topology field is required when lbMethod is TOPOLOGY")
		}
		sf := map[string]string{
			"name": topology,
		}
		err := objMgr.connector.GetObject(&DtcTopology{}, "", NewQueryParams(false, sf), &dtcTopology)
		if err != nil {
			return "", fmt.Errorf("error getting %s DtcTopology object: %s", topology, err)
		}
		topologyRef = dtcTopology[0].Ref
	}
	return topologyRef, nil
}

func getPools(pools []*DtcPoolLink, objMgr *ObjectManager) ([]*DtcPoolLink, error) {
	var dtcPoolLink []*DtcPoolLink

	for _, pool := range pools {
		sf := map[string]string{"name": pool.Pool}
		var dtcPools []DtcPool

		err := objMgr.connector.GetObject(&DtcPool{}, "", NewQueryParams(false, sf), &dtcPools)
		if err != nil {
			return nil, fmt.Errorf("error getting %s DtcPool object: %s", pool.Pool, err)
		}

		if len(dtcPools) == 0 {
			return nil, fmt.Errorf("no DtcPool object found for %s", pool.Pool)
		}

		dtcPoolLink = append(dtcPoolLink, &DtcPoolLink{Pool: dtcPools[0].Ref, Ratio: pool.Ratio})
	}
	return dtcPoolLink, nil
}

func getAuthZones(authzone []string, objMgr *ObjectManager) ([]*ZoneAuth, error) {
	var zoneAuth []ZoneAuth
	var zones []*ZoneAuth
	for i := 0; i < len(authzone); i++ {
		sf := map[string]string{
			"fqdn": authzone[i],
		}
		err := objMgr.connector.GetObject(&ZoneAuth{}, "", NewQueryParams(false, sf), &zoneAuth)
		if err != nil {
			return nil, fmt.Errorf("error getting %s ZoneAuth object: %s", authzone[i], err)
		}
		zones = append(zones, &ZoneAuth{Ref: zoneAuth[0].Ref})
	}
	return zones, nil
}

func NewDtcLbdn(ref string, name string, authzone []*ZoneAuth, comment string, disable bool, autoConsolidatedMonitors bool, ea EA,
	lbMethod string, patterns []string, persistence uint32, pools []*DtcPoolLink, priority uint32, topology string, types []string, ttl uint32, usettl bool) *DtcLbdn {

	lbdn := NewEmptyDtcLbdn()
	lbdn.Name = &name
	lbdn.Ref = ref
	lbdn.AuthZones = authzone
	// todo: add health and status_member fields
	lbdn.Comment = &comment
	lbdn.Disable = &disable
	lbdn.AutoConsolidatedMonitors = &autoConsolidatedMonitors
	lbdn.Ea = ea
	lbdn.LbMethod = lbMethod
	lbdn.Patterns = patterns
	lbdn.Persistence = &persistence
	lbdn.Pools = pools
	if topology != "" {
		lbdn.Topology = &topology
	}
	lbdn.Priority = &priority

	lbdn.Types = types
	lbdn.Ttl = &ttl
	lbdn.UseTtl = &usettl
	return lbdn
}

func NewEmptyDtcLbdn() *DtcLbdn {
	dtcLbdn := &DtcLbdn{}
	dtcLbdn.SetReturnFields(append(dtcLbdn.ReturnFields(), "extattrs", "disable", "auto_consolidated_monitors", "lb_method", "patterns", "persistence", "pools", "priority", "topology", "types", "health", "ttl", "use_ttl"))
	return dtcLbdn
}

func (objMgr *ObjectManager) DeleteDtcLbdn(ref string) (string, error) {
	return objMgr.connector.DeleteObject(ref)
}

func (objMgr *ObjectManager) GetDtcLbdn(queryParams *QueryParams) ([]DtcLbdn, error) {
	var res []DtcLbdn
	lbdn := NewEmptyDtcLbdn()
	err := objMgr.connector.GetObject(lbdn, "", queryParams, &res)
	if err != nil {
		return nil, fmt.Errorf("error getting DtcLbdn object, err: %s", err)
	}
	return res, nil
}

func (objMgr *ObjectManager) UpdateDtcLbdn(ref string, name string, authzone []string, comment string, disable bool, autoConsolidatedMonitors bool, ea EA,
	lbMethod string, patterns []string, persistence uint32, pools []*DtcPoolLink, priority uint32, topology string, types []string, ttl uint32, usettl bool) (*DtcLbdn, error) {

	// get ref id of authzones and replace
	var zones []*ZoneAuth
	var err error
	if len(authzone) > 0 {
		zones, err = getAuthZones(authzone, objMgr)
		if err != nil {
			return nil, err
		}
	}

	// get ref id of pools and replace
	dtcPoolLink, err := getPools(pools, objMgr)
	if err != nil {
		return nil, err
	}

	//get ref id of topology and replace
	topologyRef, err := getTopology(lbMethod, topology, objMgr)
	if err != nil {
		return nil, err
	}

	dtcLbdn := NewDtcLbdn(ref, name, zones, comment, disable, autoConsolidatedMonitors, ea,
		lbMethod, patterns, persistence, dtcPoolLink, priority, topologyRef, types, ttl, usettl)
	newRef, err := objMgr.connector.UpdateObject(dtcLbdn, ref)
	if err != nil {
		return nil, fmt.Errorf("error updating DtcLbdn object %s, err: %s", name, err)
	}
	dtcLbdn.Ref = newRef
	return dtcLbdn, nil
}
