package ibclient

import (
	"bytes"
	"encoding/json"
	"reflect"
)

const MACADDR_ZERO = "00:00:00:00:00:00"

type Bool bool

type EA map[string]interface{}

type EASearch map[string]interface{}

type EADefListValue string

type IBBase struct {
	objectType   string   `json:"-"`
	returnFields []string `json:"-"`
	eaSearch     EASearch `json:"-"`
}

type IBObject interface {
	ObjectType() string
	ReturnFields() []string
	EaSearch() EASearch
}

func (obj *IBBase) ObjectType() string {
	return obj.objectType
}

func (obj *IBBase) ReturnFields() []string {
	return obj.returnFields
}

func (obj *IBBase) EaSearch() EASearch {
	return obj.eaSearch
}

type NetworkView struct {
	IBBase `json:"-"`
	Ref    string `json:"_ref,omitempty"`
	Name   string `json:"name,omitempty"`
	Ea     EA     `json:"extattrs,omitempty"`
}

func NewNetworkView(nv NetworkView) *NetworkView {
	res := nv
	res.objectType = "networkview"
	res.returnFields = []string{"extattrs", "name"}

	return &res
}

type UpgradeStatus struct {
	IBBase `json:"-"`
	Ref    string `json:"_ref,omitempty"`
	Type   string `json:"type"`
}

// SubElementsStatus object representation
type SubElementsStatus struct {
	Ref            string `json:"_ref,omitempty"`
	CurrentVersion string `json:"current_version"`
	ElementStatus  string `json:"element_status"`
	Ipv4Address    string `json:"ipv4_address"`
	Ipv6Address    string `json:"ipv6_address"`
	StatusValue    string `json:"status_value"`
	StepsTotal     int    `json:"steps_total"`
	StepsCompleted int    `json:"steps_completed"`
	NodeType       string `json:"type"`
	Member         string `json:member`
}

type UpgradeStatusResult struct {
	Ref              string              `json:"_ref,omitempty"`
	SubElementStatus []SubElementsStatus `json:"subelements_status",omitempty`
	NodeType         string              `json:"type"`
	UpgradeGroup     string              `json:"upgrade_group"`
}

func NewUpgradeStatus(upgradeStatus UpgradeStatus, returnFields []string) *UpgradeStatus {
	result := upgradeStatus
	result.objectType = "upgradestatus"
	result.returnFields = returnFields
	return &result
}

type Network struct {
	IBBase
	Ref         string `json:"_ref,omitempty"`
	NetviewName string `json:"network_view,omitempty"`
	Cidr        string `json:"network,omitempty"`
	Ea          EA     `json:"extattrs,omitempty"`
}

type Members []Member

type ServiceStatus struct {
	Desciption string `json:"description,omitempty"`
	Service    string `json:"service,omitempty"`
	Status     string `json:"status,omitempty"`
}

type LanHaPortSetting struct {
	HaPortSetting  map[string]interface{} `json:"ha_port_setting,omitempty"`
	LanPortSetting map[string]bool        `json:"lan_port_setting,omitempty"`
}

type NodeInfo struct {
	HaStatus             string                 `json:"ha_status,omitempty"`
	HwId                 string                 `json:"hwid,omitempty"`
	HwModel              string                 `json:"hwmodel,omitempty"`
	HwPlatform           string                 `json:"hwplatform,omitempty"`
	HwType               string                 `json:"hwtype,omitempty"`
	Lan2PhysicalSetting  map[string]bool        `json:"lan2_physical_setting,omitempty"`
	LanHaPortSetting     LanHaPortSetting       `json:"lan_Ha_Port_Setting,omitempty"`
	MgmtNetworkSetting   map[string]interface{} `json:"mgmt_network_setting,omitempty"`
	MgmtPhysicalSetting  map[string]bool        `json:"mgmt_physical_setting,omitempty"`
	PaidNios             bool                   `json:"paid_nios,omitempty"`
	PhysicalOid          string                 `json:"physical_oid,omitempty"`
	ServiceStatus        []ServiceStatus        `json:"service_status,omitempty"`
	V6MgmtNetworkSetting map[string]interface{} `json:"v6_mgmt_network_setting,omitempty"`
}

type Member struct {
	IBBase                   `json:"-"`
	Ref                      string `json:"_ref,omitempty"`
	HostName                 string `json:"host_name,omitempty"`
	ConfigAddrType           string `json:"config_addr_type,omitempty"`
	PLATFORM                 string `json:"platform,omitempty"`
	ServiceTypeConfiguration string `json:"service_type_configuration,omitempty"`
}

// License represents license wapi object
type License struct {
	IBBase           `json:"-"`
	Ref              string `json:"_ref,omitempty"`
	ExpirationStatus string `json:"expiration_status,omitempty"`
	ExpiryDate       int    `json:"expiry_date,omitempty"`
	HwID             string `json:"hwid,omitempty"`
	Key              string `json:"key,omitempty"`
	Kind             string `json:"kind,omitempty"`
	Limit            string `json:"limit,omitempty"`
	LimitContext     string `json:"limit_context,omitempty"`
	Licensetype      string `json:"type,omitempty"`
}

type MemberResult struct {
	Ref      string     `json:"_ref,omitempty"`
	HostName string     `json:"host_name,omitempty"`
	Nodeinfo []NodeInfo `json:"node_info,omitempty"`
	TimeZone string     `json:"time_zone"`
}

// CapacityReport represents capacityreport object
type CapacityReport struct {
	IBBase `json:"-"`
	Ref    string `json:"_ref,omitempty"`

	Name         string                   `json:"name,omitempty"`
	HardwareType string                   `json:"hardware_type,omitempty"`
	MaxCapacity  int                      `json:"max_capacity,omitempty"`
	ObjectCount  []map[string]interface{} `json:"object_counts,omitempty"`
	PercentUsed  int                      `json:"percent_used,omitempty"`
	Role         string                   `json:"role,omitempty"`
	TotalObjects int                      `json:"total_objects,omitempty"`
}
type NTPSetting struct {
	enable_ntp bool                   `json:"enable_ntp,omitempty"`
	NTPAcl     map[string]interface{} `json:"ntp_acl,omitempty"`
	NTPKeys    []string               `json:"ntp_keys,omitempty"`
	NTPKod     bool                   `json:"ntp_kod,omitempty"`
	NTPServers []string               `json:"ntp_servers,omitempty"`
}

type Grid struct {
	IBBase `json:"-"`
	Ref    string `json:"_ref,omitempty"`
}

type GridResult struct {
	Ref        string     `json:"_ref,omitempty"`
	Name       string     `json:"name,omitempty"`
	NTPSetting NTPSetting `json:"ntp_setting"`
}

func NewGrid(grid Grid, returnFields []string) *Grid {
	result := grid
	result.objectType = "grid"
	result.returnFields = returnFields
	return &result
}

func NewLicense(license License, returnFields []string) *License {
	result := license
	result.objectType = "member:license"
	result.returnFields = returnFields
	return &result
}

func NewCapcityReport(capReport CapacityReport, returnFields []string) *CapacityReport {
	res := capReport
	res.objectType = "capacityreport"
	res.returnFields = returnFields
	return &res
}

func NewMember(member Member, returnFields []string) *Member {
	res := member
	res.objectType = "member"
	res.returnFields = returnFields
	return &res
}

func NewNetwork(nw Network) *Network {
	res := nw
	res.objectType = "network"

	res.returnFields = []string{"extattrs", "network", "network_view"}
	return &res
}

type NetworkContainer struct {
	IBBase      `json:"-"`
	Ref         string `json:"_ref,omitempty"`
	NetviewName string `json:"network_view,omitempty"`
	Cidr        string `json:"network,omitempty"`
	Ea          EA     `json:"extattrs,omitempty"`
}

func NewNetworkContainer(nc NetworkContainer) *NetworkContainer {
	res := nc
	res.objectType = "networkcontainer"
	res.returnFields = []string{"extattrs", "network", "network_view"}

	return &res
}

type FixedAddress struct {
	IBBase      `json:"-"`
	Ref         string `json:"_ref,omitempty"`
	NetviewName string `json:"network_view,omitempty"`
	Cidr        string `json:"network,omitempty"`
	IPAddress   string `json:"ipv4addr,omitempty"`
	Mac         string `json:"mac,omitempty"`
	Ea          EA     `json:"extattrs,omitempty"`
}

func NewFixedAddress(fixedAddr FixedAddress) *FixedAddress {
	res := fixedAddr
	res.objectType = "fixedaddress"
	res.returnFields = []string{"extattrs", "ipv4addr", "mac", "network", "network_view"}

	return &res
}

type EADefinition struct {
	IBBase             `json:"-"`
	Ref                string           `json:"_ref,omitempty"`
	Comment            string           `json:"comment,omitempty"`
	Flags              string           `json:"flags,omitempty"`
	ListValues         []EADefListValue `json:"list_values,omitempty"`
	Name               string           `json:"name,omitempty"`
	Type               string           `json:"type,omitempty"`
	AllowedObjectTypes []string         `json:"allowed_object_types,omitempty"`
}

func NewEADefinition(eadef EADefinition) *EADefinition {
	res := eadef
	res.objectType = "extensibleattributedef"
	res.returnFields = []string{"allowed_object_types", "comment", "flags", "list_values", "name", "type"}

	return &res
}

type UserProfile struct {
	IBBase `json:"-"`
	Ref    string `json:"_ref,omitempty"`
	Name   string `json:"name,omitempty"`
}

func NewUserProfile(userprofile UserProfile) *UserProfile {
	res := userprofile
	res.objectType = "userprofile"
	res.returnFields = []string{"name"}

	return &res
}

func (ea EA) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	for k, v := range ea {
		value := make(map[string]interface{})
		value["value"] = v
		m[k] = value
	}

	return json.Marshal(m)
}

func (eas EASearch) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	for k, v := range eas {
		m["*"+k] = v
	}

	return json.Marshal(m)
}

func (val EADefListValue) MarshalJSON() ([]byte, error) {
	m := make(map[string]string)
	m["value"] = string(val)

	return json.Marshal(m)
}

func (b Bool) MarshalJSON() ([]byte, error) {
	if b {
		return json.Marshal("True")
	}

	return json.Marshal("False")
}

func (ea *EA) UnmarshalJSON(b []byte) (err error) {
	var m map[string]map[string]interface{}

	decoder := json.NewDecoder(bytes.NewBuffer(b))
	decoder.UseNumber()
	err = decoder.Decode(&m)
	if err != nil {
		return
	}

	*ea = make(EA)
	for k, v := range m {
		val := v["value"]
		if reflect.TypeOf(val).String() == "json.Number" {
			var i64 int64
			i64, err = val.(json.Number).Int64()
			val = int(i64)
		} else if val.(string) == "True" {
			val = Bool(true)
		} else if val.(string) == "False" {
			val = Bool(false)
		}

		(*ea)[k] = val
	}

	return
}

func (v *EADefListValue) UnmarshalJSON(b []byte) (err error) {
	var m map[string]string
	err = json.Unmarshal(b, &m)
	if err != nil {
		return
	}

	*v = EADefListValue(m["value"])
	return
}

type RequestBody struct {
	Data               map[string]interface{} `json:"data,omitempty"`
	Args               map[string]string      `json:"args,omitempty"`
	Method             string                 `json:"method"`
	Object             string                 `json:"object,omitempty"`
	EnableSubstitution bool                   `json:"enable_substitution,omitempty"`
	AssignState        map[string]string      `json:"assign_state,omitempty"`
	Discard            bool                   `json:"discard,omitempty"`
}

type SingleRequest struct {
	IBBase `json:"-"`
	Body   *RequestBody
}

type MultiRequest struct {
	IBBase `json:"-"`
	Body   []*RequestBody
}

func (r *MultiRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Body)
}

func NewMultiRequest(body []*RequestBody) *MultiRequest {
	req := &MultiRequest{Body: body}
	req.objectType = "request"
	return req
}

func NewRequest(body *RequestBody) *SingleRequest {
	req := &SingleRequest{Body: body}
	req.objectType = "request"
	return req
}
