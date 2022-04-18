package ibclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

const MACADDR_ZERO = "00:00:00:00:00:00"

type Bool bool

func (b Bool) MarshalJSON() ([]byte, error) {
	if b {
		return json.Marshal("True")
	}

	return json.Marshal("False")
}

type EA map[string]interface{}

func (ea EA) Count() int {
	return len(ea)
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
		switch valType := reflect.TypeOf(val).String(); valType {
		case "json.Number":
			var i64 int64
			i64, err = val.(json.Number).Int64()
			val = int(i64)
		case "string":
			if val.(string) == "True" {
				val = Bool(true)
			} else if val.(string) == "False" {
				val = Bool(false)
			}
		case "[]interface {}":
			nval := val.([]interface{})
			nVals := make([]string, len(nval))
			for i, v := range nval {
				nVals[i] = fmt.Sprintf("%v", v)
			}
			val = nVals
		default:
			val = fmt.Sprintf("%v", val)
		}

		(*ea)[k] = val
	}

	return
}

type EASearch map[string]interface{}

func (eas EASearch) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	for k, v := range eas {
		m["*"+k] = v
	}

	return json.Marshal(m)
}

type IBBase struct {
	returnFields []string
	eaSearch     EASearch
}

type IBObject interface {
	ObjectType() string
	ReturnFields() []string
	EaSearch() EASearch
	SetReturnFields([]string)
}

func (obj *IBBase) ReturnFields() []string {
	return obj.returnFields
}

func (obj *IBBase) SetReturnFields(rf []string) {
	obj.returnFields = rf
}

func (obj *IBBase) EaSearch() EASearch {
	return obj.eaSearch
}

// QueryParams is a general struct to add query params used in makeRequest
type QueryParams struct {
	forceProxy bool

	searchFields map[string]string
}

func NewQueryParams(forceProxy bool, searchFields map[string]string) *QueryParams {
	qp := QueryParams{forceProxy: forceProxy}
	if searchFields != nil {
		qp.searchFields = searchFields
	} else {
		qp.searchFields = make(map[string]string)
	}

	return &qp
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
	return req
}

func (MultiRequest) ObjectType() string {
	return "request"
}

func NewRequest(body *RequestBody) *SingleRequest {
	req := &SingleRequest{Body: body}
	return req
}

func (SingleRequest) ObjectType() string {
	return "request"
}

type FixedAddress struct {
	IBBase      `json:"-"`
	objectType  string
	Ref         string `json:"_ref,omitempty"`
	NetviewName string `json:"network_view,omitempty"`
	Cidr        string `json:"network,omitempty"`
	Comment     string `json:"comment"`
	IPv4Address string `json:"ipv4addr,omitempty"`
	IPv6Address string `json:"ipv6addr,omitempty"`
	Duid        string `json:"duid,omitempty"`
	Mac         string `json:"mac,omitempty"`
	Name        string `json:"name,omitempty"`
	MatchClient string `json:"match_client,omitempty"`
	Ea          EA     `json:"extattrs"`
}

func (fa FixedAddress) ObjectType() string {
	return fa.objectType
}

func NewEmptyFixedAddress(isIPv6 bool) *FixedAddress {
	res := &FixedAddress{}
	res.Ea = make(EA)
	if isIPv6 {
		res.objectType = "ipv6fixedaddress"
		res.returnFields = []string{"extattrs", "ipv6addr", "duid", "name", "network", "network_view", "comment"}
	} else {
		res.objectType = "fixedaddress"
		res.returnFields = []string{"extattrs", "ipv4addr", "mac", "name", "network", "network_view", "comment"}
	}
	return res
}

func NewFixedAddress(
	netView string,
	name string,
	ipAddr string,
	cidr string,
	macOrDuid string,
	clients string,
	eas EA,
	ref string,
	isIPv6 bool,
	comment string) *FixedAddress {

	res := NewEmptyFixedAddress(isIPv6)
	res.NetviewName = netView
	res.Name = name
	res.Cidr = cidr
	res.MatchClient = clients
	res.Ea = eas
	res.Ref = ref
	res.Comment = comment
	if isIPv6 {
		res.IPv6Address = ipAddr
		res.Duid = macOrDuid
	} else {
		res.IPv4Address = ipAddr
		res.Mac = macOrDuid
	}
	return res
}

type Network struct {
	IBBase      `json:"-"`
	objectType  string
	Ref         string `json:"_ref,omitempty"`
	NetviewName string `json:"network_view,omitempty"`
	Cidr        string `json:"network,omitempty"`
	Ea          EA     `json:"extattrs"`
	Comment     string `json:"comment"`
}

func (n Network) ObjectType() string {
	return n.objectType
}

func NewNetwork(netview string, cidr string, isIPv6 bool, comment string, ea EA) *Network {
	var res Network
	res.NetviewName = netview
	res.Cidr = cidr
	res.Ea = ea
	res.Comment = comment
	if isIPv6 {
		res.objectType = "ipv6network"
	} else {
		res.objectType = "network"
	}
	res.returnFields = []string{"extattrs", "network", "comment"}

	return &res
}

type NetworkContainer struct {
	IBBase      `json:"-"`
	objectType  string
	Ref         string `json:"_ref,omitempty"`
	NetviewName string `json:"network_view,omitempty"`
	Cidr        string `json:"network,omitempty"`
	Comment     string `json:"comment"`
	Ea          EA     `json:"extattrs"`
}

func (nc NetworkContainer) ObjectType() string {
	return nc.objectType
}

func NewNetworkContainer(netview, cidr string, isIPv6 bool, comment string, ea EA) *NetworkContainer {
	nc := NetworkContainer{
		NetviewName: netview,
		Cidr:        cidr,
		Ea:          ea,
		Comment:     comment,
	}

	if isIPv6 {
		nc.objectType = "ipv6networkcontainer"
	} else {
		nc.objectType = "networkcontainer"
	}
	nc.returnFields = []string{"extattrs", "network", "network_view", "comment"}

	return &nc
}

// License represents license wapi object
type License struct {
	IBBase           `json:"-"`
	objectType       string
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

func (l License) ObjectType() string {
	return l.objectType
}

func NewGridLicense(license License) *License {
	result := license
	result.objectType = "license:gridwide"
	returnFields := []string{"expiration_status",
		"expiry_date",
		"key",
		"limit",
		"limit_context",
		"type"}
	result.returnFields = returnFields
	return &result
}

func NewLicense(license License) *License {
	result := license
	returnFields := []string{"expiration_status",
		"expiry_date",
		"hwid",
		"key",
		"kind",
		"limit",
		"limit_context",
		"type"}
	result.objectType = "member:license"
	result.returnFields = returnFields
	return &result
}

//   AUTOGENERATED CODE BELOW

// Allrpzrecords represents Infoblox object allrpzrecords
type Allrpzrecords struct {
	IBBase         `json:"-"`
	Ref            string `json:"_ref,omitempty"`
	AlertType      string `json:"alert_type,omitempty"`
	Comment        string `json:"comment,omitempty"`
	Disable        bool   `json:"disable,omitempty"`
	ExpirationTime uint32 `json:"expiration_time,omitempty"`
	LastUpdated    uint32 `json:"last_updated,omitempty"`
	Name           string `json:"name,omitempty"`
	Record         string `json:"record,omitempty"`
	RpzRule        string `json:"rpz_rule,omitempty"`
	Ttl            uint32 `json:"ttl,omitempty"`
	Type           string `json:"type,omitempty"`
	View           string `json:"view,omitempty"`
	Zone           string `json:"zone,omitempty"`
}

func (Allrpzrecords) ObjectType() string {
	return "allrpzrecords"
}

func (obj Allrpzrecords) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name", "type", "view", "zone"}
	}
	return obj.returnFields
}

// Allnsgroup represents Infoblox object allnsgroup
type Allnsgroup struct {
	IBBase  `json:"-"`
	Ref     string `json:"_ref,omitempty"`
	Comment string `json:"comment,omitempty"`
	Name    string `json:"name,omitempty"`
	Type    string `json:"type,omitempty"`
}

func (Allnsgroup) ObjectType() string {
	return "allnsgroup"
}

func (obj Allnsgroup) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name", "type"}
	}
	return obj.returnFields
}

// Adminrole represents Infoblox object adminrole
type Adminrole struct {
	IBBase  `json:"-"`
	Ref     string `json:"_ref,omitempty"`
	Comment string `json:"comment,omitempty"`
	Disable bool   `json:"disable,omitempty"`
	Ea      EA     `json:"extattrs,omitempty"`
	Name    string `json:"name,omitempty"`
}

func (Adminrole) ObjectType() string {
	return "adminrole"
}

func (obj Adminrole) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name"}
	}
	return obj.returnFields
}

// Authpolicy represents Infoblox object authpolicy
type Authpolicy struct {
	IBBase       `json:"-"`
	Ref          string                  `json:"_ref,omitempty"`
	AdminGroups  []string                `json:"admin_groups,omitempty"`
	AuthServices []*LocaluserAuthservice `json:"auth_services,omitempty"`
	DefaultGroup string                  `json:"default_group,omitempty"`
	UsageType    string                  `json:"usage_type,omitempty"`
}

func (Authpolicy) ObjectType() string {
	return "authpolicy"
}

func (obj Authpolicy) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"default_group", "usage_type"}
	}
	return obj.returnFields
}

// Allendpoints represents Infoblox object allendpoints
type Allendpoints struct {
	IBBase            `json:"-"`
	Ref               string `json:"_ref,omitempty"`
	Address           string `json:"address,omitempty"`
	Comment           string `json:"comment,omitempty"`
	Disable           bool   `json:"disable,omitempty"`
	SubscribingMember string `json:"subscribing_member,omitempty"`
	Type              string `json:"type,omitempty"`
	Version           string `json:"version,omitempty"`
}

func (Allendpoints) ObjectType() string {
	return "allendpoints"
}

func (obj Allendpoints) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{}
	}
	return obj.returnFields
}

// Allrecords represents Infoblox object allrecords
type Allrecords struct {
	IBBase        `json:"-"`
	Ref           string `json:"_ref,omitempty"`
	Address       string `json:"address,omitempty"`
	Comment       string `json:"comment,omitempty"`
	Creator       string `json:"creator,omitempty"`
	DdnsPrincipal string `json:"ddns_principal,omitempty"`
	DdnsProtected bool   `json:"ddns_protected,omitempty"`
	Disable       bool   `json:"disable,omitempty"`
	DtcObscured   string `json:"dtc_obscured,omitempty"`
	Name          string `json:"name,omitempty"`
	Reclaimable   bool   `json:"reclaimable,omitempty"`
	Record        string `json:"record,omitempty"`
	Ttl           uint32 `json:"ttl,omitempty"`
	Type          string `json:"type,omitempty"`
	View          string `json:"view,omitempty"`
	Zone          string `json:"zone,omitempty"`
}

func (Allrecords) ObjectType() string {
	return "allrecords"
}

func (obj Allrecords) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name", "type", "view", "zone"}
	}
	return obj.returnFields
}

// Adminuser represents Infoblox object adminuser
type Adminuser struct {
	IBBase                          `json:"-"`
	Ref                             string   `json:"_ref,omitempty"`
	AdminGroups                     []string `json:"admin_groups,omitempty"`
	AuthType                        string   `json:"auth_type,omitempty"`
	CaCertificateIssuer             string   `json:"ca_certificate_issuer,omitempty"`
	ClientCertificateSerialNumber   string   `json:"client_certificate_serial_number,omitempty"`
	Comment                         string   `json:"comment,omitempty"`
	Disable                         bool     `json:"disable,omitempty"`
	Email                           string   `json:"email,omitempty"`
	EnableCertificateAuthentication bool     `json:"enable_certificate_authentication,omitempty"`
	Ea                              EA       `json:"extattrs,omitempty"`
	Name                            string   `json:"name,omitempty"`
	Password                        string   `json:"password,omitempty"`
	Status                          string   `json:"status,omitempty"`
	TimeZone                        string   `json:"time_zone,omitempty"`
	UseTimeZone                     bool     `json:"use_time_zone,omitempty"`
}

func (Adminuser) ObjectType() string {
	return "adminuser"
}

func (obj Adminuser) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"admin_groups", "comment", "name"}
	}
	return obj.returnFields
}

// AdAuthService represents Infoblox object ad_auth_service
type AdAuthService struct {
	IBBase                   `json:"-"`
	Ref                      string          `json:"_ref,omitempty"`
	AdDomain                 string          `json:"ad_domain,omitempty"`
	AdditionalSearchPaths    []string        `json:"additional_search_paths,omitempty"`
	Comment                  string          `json:"comment,omitempty"`
	DisableDefaultSearchPath bool            `json:"disable_default_search_path,omitempty"`
	Disabled                 bool            `json:"disabled,omitempty"`
	DomainControllers        []*AdAuthServer `json:"domain_controllers,omitempty"`
	Name                     string          `json:"name,omitempty"`
	NestedGroupQuerying      bool            `json:"nested_group_querying,omitempty"`
	Timeout                  uint32          `json:"timeout,omitempty"`
}

func (AdAuthService) ObjectType() string {
	return "ad_auth_service"
}

func (obj AdAuthService) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name"}
	}
	return obj.returnFields
}

// Approvalworkflow represents Infoblox object approvalworkflow
type Approvalworkflow struct {
	IBBase                  `json:"-"`
	Ref                     string `json:"_ref,omitempty"`
	ApprovalGroup           string `json:"approval_group,omitempty"`
	ApprovalNotifyTo        string `json:"approval_notify_to,omitempty"`
	ApprovedNotifyTo        string `json:"approved_notify_to,omitempty"`
	ApproverComment         string `json:"approver_comment,omitempty"`
	EnableApprovalNotify    bool   `json:"enable_approval_notify,omitempty"`
	EnableApprovedNotify    bool   `json:"enable_approved_notify,omitempty"`
	EnableFailedNotify      bool   `json:"enable_failed_notify,omitempty"`
	EnableNotifyGroup       bool   `json:"enable_notify_group,omitempty"`
	EnableNotifyUser        bool   `json:"enable_notify_user,omitempty"`
	EnableRejectedNotify    bool   `json:"enable_rejected_notify,omitempty"`
	EnableRescheduledNotify bool   `json:"enable_rescheduled_notify,omitempty"`
	EnableSucceededNotify   bool   `json:"enable_succeeded_notify,omitempty"`
	Ea                      EA     `json:"extattrs,omitempty"`
	FailedNotifyTo          string `json:"failed_notify_to,omitempty"`
	RejectedNotifyTo        string `json:"rejected_notify_to,omitempty"`
	RescheduledNotifyTo     string `json:"rescheduled_notify_to,omitempty"`
	SubmitterComment        string `json:"submitter_comment,omitempty"`
	SubmitterGroup          string `json:"submitter_group,omitempty"`
	SucceededNotifyTo       string `json:"succeeded_notify_to,omitempty"`
	TicketNumber            string `json:"ticket_number,omitempty"`
}

func (Approvalworkflow) ObjectType() string {
	return "approvalworkflow"
}

func (obj Approvalworkflow) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"approval_group", "submitter_group"}
	}
	return obj.returnFields
}

// Awsuser represents Infoblox object awsuser
type Awsuser struct {
	IBBase          `json:"-"`
	Ref             string     `json:"_ref,omitempty"`
	AccessKeyId     string     `json:"access_key_id,omitempty"`
	AccountId       string     `json:"account_id,omitempty"`
	LastUsed        *time.Time `json:"last_used,omitempty"`
	Name            string     `json:"name,omitempty"`
	NiosUserName    string     `json:"nios_user_name,omitempty"`
	SecretAccessKey string     `json:"secret_access_key,omitempty"`
	Status          string     `json:"status,omitempty"`
}

func (Awsuser) ObjectType() string {
	return "awsuser"
}

func (obj Awsuser) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"access_key_id", "account_id", "name"}
	}
	return obj.returnFields
}

// Bulkhost represents Infoblox object bulkhost
type Bulkhost struct {
	IBBase          `json:"-"`
	Ref             string            `json:"_ref,omitempty"`
	CloudInfo       *GridCloudapiInfo `json:"cloud_info,omitempty"`
	Comment         string            `json:"comment,omitempty"`
	Disable         bool              `json:"disable,omitempty"`
	DnsPrefix       string            `json:"dns_prefix,omitempty"`
	EndAddr         string            `json:"end_addr,omitempty"`
	Ea              EA                `json:"extattrs,omitempty"`
	LastQueried     *time.Time        `json:"last_queried,omitempty"`
	NameTemplate    string            `json:"name_template,omitempty"`
	NetworkView     string            `json:"network_view,omitempty"`
	Policy          string            `json:"policy,omitempty"`
	Prefix          string            `json:"prefix,omitempty"`
	Reverse         bool              `json:"reverse,omitempty"`
	StartAddr       string            `json:"start_addr,omitempty"`
	TemplateFormat  string            `json:"template_format,omitempty"`
	Ttl             uint32            `json:"ttl,omitempty"`
	UseNameTemplate bool              `json:"use_name_template,omitempty"`
	UseTtl          bool              `json:"use_ttl,omitempty"`
	View            string            `json:"view,omitempty"`
	Zone            string            `json:"zone,omitempty"`
}

func (Bulkhost) ObjectType() string {
	return "bulkhost"
}

func (obj Bulkhost) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "prefix"}
	}
	return obj.returnFields
}

// Bulkhostnametemplate represents Infoblox object bulkhostnametemplate
type Bulkhostnametemplate struct {
	IBBase         `json:"-"`
	Ref            string `json:"_ref,omitempty"`
	IsGridDefault  bool   `json:"is_grid_default,omitempty"`
	PreDefined     bool   `json:"pre_defined,omitempty"`
	TemplateFormat string `json:"template_format,omitempty"`
	TemplateName   string `json:"template_name,omitempty"`
}

func (Bulkhostnametemplate) ObjectType() string {
	return "bulkhostnametemplate"
}

func (obj Bulkhostnametemplate) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"is_grid_default", "template_format", "template_name"}
	}
	return obj.returnFields
}

// Bfdtemplate represents Infoblox object bfdtemplate
type Bfdtemplate struct {
	IBBase              `json:"-"`
	Ref                 string `json:"_ref,omitempty"`
	AuthenticationKey   string `json:"authentication_key,omitempty"`
	AuthenticationKeyId uint32 `json:"authentication_key_id,omitempty"`
	AuthenticationType  string `json:"authentication_type,omitempty"`
	DetectionMultiplier uint32 `json:"detection_multiplier,omitempty"`
	MinRxInterval       uint32 `json:"min_rx_interval,omitempty"`
	MinTxInterval       uint32 `json:"min_tx_interval,omitempty"`
	Name                string `json:"name,omitempty"`
}

func (Bfdtemplate) ObjectType() string {
	return "bfdtemplate"
}

func (obj Bfdtemplate) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name"}
	}
	return obj.returnFields
}

// Cacertificate represents Infoblox object cacertificate
type Cacertificate struct {
	IBBase            `json:"-"`
	Ref               string     `json:"_ref,omitempty"`
	DistinguishedName string     `json:"distinguished_name,omitempty"`
	Issuer            string     `json:"issuer,omitempty"`
	Serial            string     `json:"serial,omitempty"`
	UsedBy            string     `json:"used_by,omitempty"`
	ValidNotAfter     *time.Time `json:"valid_not_after,omitempty"`
	ValidNotBefore    *time.Time `json:"valid_not_before,omitempty"`
}

func (Cacertificate) ObjectType() string {
	return "cacertificate"
}

func (obj Cacertificate) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"distinguished_name", "issuer", "serial", "used_by", "valid_not_after", "valid_not_before"}
	}
	return obj.returnFields
}

// Awsrte53taskgroup represents Infoblox object awsrte53taskgroup
type Awsrte53taskgroup struct {
	IBBase                   `json:"-"`
	Ref                      string          `json:"_ref,omitempty"`
	AccountId                string          `json:"account_id,omitempty"`
	Comment                  string          `json:"comment,omitempty"`
	ConsolidateZones         bool            `json:"consolidate_zones,omitempty"`
	ConsolidatedView         string          `json:"consolidated_view,omitempty"`
	Disabled                 bool            `json:"disabled,omitempty"`
	GridMember               string          `json:"grid_member,omitempty"`
	Name                     string          `json:"name,omitempty"`
	NetworkView              string          `json:"network_view,omitempty"`
	NetworkViewMappingPolicy string          `json:"network_view_mapping_policy,omitempty"`
	SyncStatus               string          `json:"sync_status,omitempty"`
	TaskControl              *Taskcontrol    `json:"task_control,omitempty"`
	TaskList                 []*Awsrte53task `json:"task_list,omitempty"`
}

func (Awsrte53taskgroup) ObjectType() string {
	return "awsrte53taskgroup"
}

func (obj Awsrte53taskgroup) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"account_id", "comment", "disabled", "name", "sync_status"}
	}
	return obj.returnFields
}

// CapacityReport represents Infoblox object capacityreport
type CapacityReport struct {
	IBBase       `json:"-"`
	Ref          string                       `json:"_ref,omitempty"`
	HardwareType string                       `json:"hardware_type,omitempty"`
	MaxCapacity  uint32                       `json:"max_capacity,omitempty"`
	Name         string                       `json:"name,omitempty"`
	ObjectCounts []*CapacityreportObjectcount `json:"object_counts,omitempty"`
	PercentUsed  uint32                       `json:"percent_used,omitempty"`
	Role         string                       `json:"role,omitempty"`
	TotalObjects uint32                       `json:"total_objects,omitempty"`
}

func (CapacityReport) ObjectType() string {
	return "capacityreport"
}

func (obj CapacityReport) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name", "percent_used", "role"}
	}
	return obj.returnFields
}

func NewCapcityReport(capReport CapacityReport) *CapacityReport {
	res := capReport
	returnFields := []string{"name", "hardware_type", "max_capacity", "object_counts", "percent_used", "role", "total_objects"}
	res.returnFields = returnFields
	return &res
}

// Captiveportal represents Infoblox object captiveportal
type Captiveportal struct {
	IBBase                    `json:"-"`
	Ref                       string               `json:"_ref,omitempty"`
	AuthnServerGroup          string               `json:"authn_server_group,omitempty"`
	CompanyName               string               `json:"company_name,omitempty"`
	EnableSyslogAuthFailure   bool                 `json:"enable_syslog_auth_failure,omitempty"`
	EnableSyslogAuthSuccess   bool                 `json:"enable_syslog_auth_success,omitempty"`
	EnableUserType            string               `json:"enable_user_type,omitempty"`
	Encryption                string               `json:"encryption,omitempty"`
	Files                     []*CaptiveportalFile `json:"files,omitempty"`
	GuestCustomField1Name     string               `json:"guest_custom_field1_name,omitempty"`
	GuestCustomField1Required bool                 `json:"guest_custom_field1_required,omitempty"`
	GuestCustomField2Name     string               `json:"guest_custom_field2_name,omitempty"`
	GuestCustomField2Required bool                 `json:"guest_custom_field2_required,omitempty"`
	GuestCustomField3Name     string               `json:"guest_custom_field3_name,omitempty"`
	GuestCustomField3Required bool                 `json:"guest_custom_field3_required,omitempty"`
	GuestCustomField4Name     string               `json:"guest_custom_field4_name,omitempty"`
	GuestCustomField4Required bool                 `json:"guest_custom_field4_required,omitempty"`
	GuestEmailRequired        bool                 `json:"guest_email_required,omitempty"`
	GuestFirstNameRequired    bool                 `json:"guest_first_name_required,omitempty"`
	GuestLastNameRequired     bool                 `json:"guest_last_name_required,omitempty"`
	GuestMiddleNameRequired   bool                 `json:"guest_middle_name_required,omitempty"`
	GuestPhoneRequired        bool                 `json:"guest_phone_required,omitempty"`
	HelpdeskMessage           string               `json:"helpdesk_message,omitempty"`
	ListenAddressIp           string               `json:"listen_address_ip,omitempty"`
	ListenAddressType         string               `json:"listen_address_type,omitempty"`
	Name                      string               `json:"name,omitempty"`
	NetworkView               string               `json:"network_view,omitempty"`
	Port                      uint32               `json:"port,omitempty"`
	ServiceEnabled            bool                 `json:"service_enabled,omitempty"`
	SyslogAuthFailureLevel    string               `json:"syslog_auth_failure_level,omitempty"`
	SyslogAuthSuccessLevel    string               `json:"syslog_auth_success_level,omitempty"`
	WelcomeMessage            string               `json:"welcome_message,omitempty"`
}

func (Captiveportal) ObjectType() string {
	return "captiveportal"
}

func (obj Captiveportal) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name"}
	}
	return obj.returnFields
}

// CertificateAuthservice represents Infoblox object certificate:authservice
type CertificateAuthservice struct {
	IBBase                    `json:"-"`
	Ref                       string                     `json:"_ref,omitempty"`
	AutoPopulateLogin         string                     `json:"auto_populate_login,omitempty"`
	CaCertificates            []*Cacertificate           `json:"ca_certificates,omitempty"`
	Comment                   string                     `json:"comment,omitempty"`
	Disabled                  bool                       `json:"disabled,omitempty"`
	EnablePasswordRequest     bool                       `json:"enable_password_request,omitempty"`
	EnableRemoteLookup        bool                       `json:"enable_remote_lookup,omitempty"`
	MaxRetries                uint32                     `json:"max_retries,omitempty"`
	Name                      string                     `json:"name,omitempty"`
	OcspCheck                 string                     `json:"ocsp_check,omitempty"`
	OcspResponders            []*OcspResponder           `json:"ocsp_responders,omitempty"`
	RecoveryInterval          uint32                     `json:"recovery_interval,omitempty"`
	RemoteLookupPassword      string                     `json:"remote_lookup_password,omitempty"`
	RemoteLookupService       string                     `json:"remote_lookup_service,omitempty"`
	RemoteLookupUsername      string                     `json:"remote_lookup_username,omitempty"`
	ResponseTimeout           uint32                     `json:"response_timeout,omitempty"`
	TestOcspResponderSettings *Testocsprespondersettings `json:"test_ocsp_responder_settings,omitempty"`
	TrustModel                string                     `json:"trust_model,omitempty"`
	UserMatchType             string                     `json:"user_match_type,omitempty"`
}

func (CertificateAuthservice) ObjectType() string {
	return "certificate:authservice"
}

func (obj CertificateAuthservice) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name"}
	}
	return obj.returnFields
}

// Admingroup represents Infoblox object admingroup
type Admingroup struct {
	IBBase                            `json:"-"`
	Ref                               string                                     `json:"_ref,omitempty"`
	AccessMethod                      []string                                   `json:"access_method,omitempty"`
	AdminSetCommands                  *AdmingroupAdminsetcommands                `json:"admin_set_commands,omitempty"`
	AdminShowCommands                 *AdmingroupAdminshowcommands               `json:"admin_show_commands,omitempty"`
	AdminToplevelCommands             *AdmingroupAdmintoplevelcommands           `json:"admin_toplevel_commands,omitempty"`
	CloudSetCommands                  *AdmingroupCloudsetcommands                `json:"cloud_set_commands,omitempty"`
	Comment                           string                                     `json:"comment,omitempty"`
	DatabaseSetCommands               *AdmingroupDatabasesetcommands             `json:"database_set_commands,omitempty"`
	DatabaseShowCommands              *AdmingroupDatabaseshowcommands            `json:"database_show_commands,omitempty"`
	DhcpSetCommands                   *AdmingroupDhcpsetcommands                 `json:"dhcp_set_commands,omitempty"`
	DhcpShowCommands                  *AdmingroupDhcpshowcommands                `json:"dhcp_show_commands,omitempty"`
	Disable                           bool                                       `json:"disable,omitempty"`
	DisableConcurrentLogin            bool                                       `json:"disable_concurrent_login,omitempty"`
	DnsSetCommands                    *AdmingroupDnssetcommands                  `json:"dns_set_commands,omitempty"`
	DnsShowCommands                   *AdmingroupDnsshowcommands                 `json:"dns_show_commands,omitempty"`
	DnsToplevelCommands               *AdmingroupDnstoplevelcommands             `json:"dns_toplevel_commands,omitempty"`
	DockerSetCommands                 *AdmingroupDockersetcommands               `json:"docker_set_commands,omitempty"`
	DockerShowCommands                *AdmingroupDockershowcommands              `json:"docker_show_commands,omitempty"`
	EmailAddresses                    []string                                   `json:"email_addresses,omitempty"`
	EnableRestrictedUserAccess        bool                                       `json:"enable_restricted_user_access,omitempty"`
	Ea                                EA                                         `json:"extattrs,omitempty"`
	GridSetCommands                   *AdmingroupGridsetcommands                 `json:"grid_set_commands,omitempty"`
	GridShowCommands                  *AdmingroupGridshowcommands                `json:"grid_show_commands,omitempty"`
	InactivityLockoutSetting          *SettingInactivelockout                    `json:"inactivity_lockout_setting,omitempty"`
	LicensingSetCommands              *AdmingroupLicensingsetcommands            `json:"licensing_set_commands,omitempty"`
	LicensingShowCommands             *AdmingroupLicensingshowcommands           `json:"licensing_show_commands,omitempty"`
	LockoutSetting                    *AdmingroupLockoutsetting                  `json:"lockout_setting,omitempty"`
	MachineControlToplevelCommands    *AdmingroupMachinecontroltoplevelcommands  `json:"machine_control_toplevel_commands,omitempty"`
	Name                              string                                     `json:"name,omitempty"`
	NetworkingSetCommands             *AdmingroupNetworkingsetcommands           `json:"networking_set_commands,omitempty"`
	NetworkingShowCommands            *AdmingroupNetworkingshowcommands          `json:"networking_show_commands,omitempty"`
	PasswordSetting                   *AdmingroupPasswordsetting                 `json:"password_setting,omitempty"`
	Roles                             []string                                   `json:"roles,omitempty"`
	SamlSetting                       *AdmingroupSamlsetting                     `json:"saml_setting,omitempty"`
	SecuritySetCommands               *AdmingroupSecuritysetcommands             `json:"security_set_commands,omitempty"`
	SecurityShowCommands              *AdmingroupSecurityshowcommands            `json:"security_show_commands,omitempty"`
	Superuser                         bool                                       `json:"superuser,omitempty"`
	TroubleShootingToplevelCommands   *AdmingroupTroubleshootingtoplevelcommands `json:"trouble_shooting_toplevel_commands,omitempty"`
	UseAccountInactivityLockoutEnable bool                                       `json:"use_account_inactivity_lockout_enable,omitempty"`
	UseDisableConcurrentLogin         bool                                       `json:"use_disable_concurrent_login,omitempty"`
	UseLockoutSetting                 bool                                       `json:"use_lockout_setting,omitempty"`
	UsePasswordSetting                bool                                       `json:"use_password_setting,omitempty"`
	UserAccess                        []*Addressac                               `json:"user_access,omitempty"`
}

func (Admingroup) ObjectType() string {
	return "admingroup"
}

func (obj Admingroup) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name"}
	}
	return obj.returnFields
}

// DdnsPrincipalclusterGroup represents Infoblox object ddns:principalcluster:group
type DdnsPrincipalclusterGroup struct {
	IBBase   `json:"-"`
	Ref      string                  `json:"_ref,omitempty"`
	Clusters []*DdnsPrincipalcluster `json:"clusters,omitempty"`
	Comment  string                  `json:"comment,omitempty"`
	Name     string                  `json:"name,omitempty"`
}

func (DdnsPrincipalclusterGroup) ObjectType() string {
	return "ddns:principalcluster:group"
}

func (obj DdnsPrincipalclusterGroup) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name"}
	}
	return obj.returnFields
}

// DdnsPrincipalcluster represents Infoblox object ddns:principalcluster
type DdnsPrincipalcluster struct {
	IBBase     `json:"-"`
	Ref        string   `json:"_ref,omitempty"`
	Comment    string   `json:"comment,omitempty"`
	Group      string   `json:"group,omitempty"`
	Name       string   `json:"name,omitempty"`
	Principals []string `json:"principals,omitempty"`
}

func (DdnsPrincipalcluster) ObjectType() string {
	return "ddns:principalcluster"
}

func (obj DdnsPrincipalcluster) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "group", "name", "principals"}
	}
	return obj.returnFields
}

// Csvimporttask represents Infoblox object csvimporttask
type Csvimporttask struct {
	IBBase         `json:"-"`
	Ref            string     `json:"_ref,omitempty"`
	Action         string     `json:"action,omitempty"`
	AdminName      string     `json:"admin_name,omitempty"`
	EndTime        *time.Time `json:"end_time,omitempty"`
	FileName       string     `json:"file_name,omitempty"`
	FileSize       uint32     `json:"file_size,omitempty"`
	ImportId       uint32     `json:"import_id,omitempty"`
	LinesFailed    uint32     `json:"lines_failed,omitempty"`
	LinesProcessed uint32     `json:"lines_processed,omitempty"`
	LinesWarning   uint32     `json:"lines_warning,omitempty"`
	OnError        string     `json:"on_error,omitempty"`
	Operation      string     `json:"operation,omitempty"`
	Separator      string     `json:"separator,omitempty"`
	StartTime      *time.Time `json:"start_time,omitempty"`
	Status         string     `json:"status,omitempty"`
	Stop           *Stopcsv   `json:"stop,omitempty"`
	UpdateMethod   string     `json:"update_method,omitempty"`
}

func (Csvimporttask) ObjectType() string {
	return "csvimporttask"
}

func (obj Csvimporttask) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"action", "admin_name", "end_time", "file_name", "file_size", "import_id", "lines_failed", "lines_processed", "lines_warning", "on_error", "operation", "separator", "start_time", "status", "update_method"}
	}
	return obj.returnFields
}

// CiscoiseEndpoint represents Infoblox object ciscoise:endpoint
type CiscoiseEndpoint struct {
	IBBase                           `json:"-"`
	Ref                              string                    `json:"_ref,omitempty"`
	Address                          string                    `json:"address,omitempty"`
	BulkDownloadCertificateSubject   string                    `json:"bulk_download_certificate_subject,omitempty"`
	BulkDownloadCertificateToken     string                    `json:"bulk_download_certificate_token,omitempty"`
	BulkDownloadCertificateValidFrom *time.Time                `json:"bulk_download_certificate_valid_from,omitempty"`
	BulkDownloadCertificateValidTo   *time.Time                `json:"bulk_download_certificate_valid_to,omitempty"`
	ClientCertificateSubject         string                    `json:"client_certificate_subject,omitempty"`
	ClientCertificateToken           string                    `json:"client_certificate_token,omitempty"`
	ClientCertificateValidFrom       *time.Time                `json:"client_certificate_valid_from,omitempty"`
	ClientCertificateValidTo         *time.Time                `json:"client_certificate_valid_to,omitempty"`
	Comment                          string                    `json:"comment,omitempty"`
	ConnectionStatus                 string                    `json:"connection_status,omitempty"`
	ConnectionTimeout                uint32                    `json:"connection_timeout,omitempty"`
	Disable                          bool                      `json:"disable,omitempty"`
	Ea                               EA                        `json:"extattrs,omitempty"`
	NetworkView                      string                    `json:"network_view,omitempty"`
	PublishSettings                  *CiscoisePublishsetting   `json:"publish_settings,omitempty"`
	ResolvedAddress                  string                    `json:"resolved_address,omitempty"`
	ResolvedSecondaryAddress         string                    `json:"resolved_secondary_address,omitempty"`
	SecondaryAddress                 string                    `json:"secondary_address,omitempty"`
	SubscribeSettings                *CiscoiseSubscribesetting `json:"subscribe_settings,omitempty"`
	SubscribingMember                string                    `json:"subscribing_member,omitempty"`
	TestConnection                   *Testendpointconnection   `json:"test_connection,omitempty"`
	Type                             string                    `json:"type,omitempty"`
	Version                          string                    `json:"version,omitempty"`
}

func (CiscoiseEndpoint) ObjectType() string {
	return "ciscoise:endpoint"
}

func (obj CiscoiseEndpoint) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"address", "disable", "resolved_address", "type", "version"}
	}
	return obj.returnFields
}

// DbObjects represents Infoblox object db_objects
type DbObjects struct {
	IBBase          `json:"-"`
	Ref             string `json:"_ref,omitempty"`
	LastSequenceId  string `json:"last_sequence_id,omitempty"`
	Object          string `json:"object,omitempty"`
	ObjectTypeField string `json:"object_type,omitempty"`
	UniqueId        string `json:"unique_id,omitempty"`
}

func (DbObjects) ObjectType() string {
	return "db_objects"
}

func (obj DbObjects) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"last_sequence_id", "object", "object_type", "unique_id"}
	}
	return obj.returnFields
}

// Dbsnapshot represents Infoblox object dbsnapshot
type Dbsnapshot struct {
	IBBase             `json:"-"`
	Ref                string          `json:"_ref,omitempty"`
	Comment            string          `json:"comment,omitempty"`
	RollbackDbSnapshot *Emptyparams    `json:"rollback_db_snapshot,omitempty"`
	SaveDbSnapshot     *Savedbsnapshot `json:"save_db_snapshot,omitempty"`
	Timestamp          *time.Time      `json:"timestamp,omitempty"`
}

func (Dbsnapshot) ObjectType() string {
	return "dbsnapshot"
}

func (obj Dbsnapshot) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "timestamp"}
	}
	return obj.returnFields
}

// DeletedObjects represents Infoblox object deleted_objects
type DeletedObjects struct {
	IBBase          `json:"-"`
	Ref             string `json:"_ref,omitempty"`
	ObjectTypeField string `json:"object_type,omitempty"`
}

func (DeletedObjects) ObjectType() string {
	return "deleted_objects"
}

func (obj DeletedObjects) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"object_type"}
	}
	return obj.returnFields
}

// DhcpStatistics represents Infoblox object dhcp:statistics
type DhcpStatistics struct {
	IBBase                `json:"-"`
	Ref                   string `json:"_ref,omitempty"`
	DhcpUtilization       uint32 `json:"dhcp_utilization,omitempty"`
	DhcpUtilizationStatus string `json:"dhcp_utilization_status,omitempty"`
	DynamicHosts          uint32 `json:"dynamic_hosts,omitempty"`
	StaticHosts           uint32 `json:"static_hosts,omitempty"`
	TotalHosts            uint32 `json:"total_hosts,omitempty"`
}

func (DhcpStatistics) ObjectType() string {
	return "dhcp:statistics"
}

func (obj DhcpStatistics) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"dhcp_utilization", "dhcp_utilization_status", "dynamic_hosts", "static_hosts", "total_hosts"}
	}
	return obj.returnFields
}

// Dhcpfailover represents Infoblox object dhcpfailover
type Dhcpfailover struct {
	IBBase                           `json:"-"`
	Ref                              string                      `json:"_ref,omitempty"`
	AssociationType                  string                      `json:"association_type,omitempty"`
	Comment                          string                      `json:"comment,omitempty"`
	Ea                               EA                          `json:"extattrs,omitempty"`
	FailoverPort                     uint32                      `json:"failover_port,omitempty"`
	LoadBalanceSplit                 uint32                      `json:"load_balance_split,omitempty"`
	MaxClientLeadTime                uint32                      `json:"max_client_lead_time,omitempty"`
	MaxLoadBalanceDelay              uint32                      `json:"max_load_balance_delay,omitempty"`
	MaxResponseDelay                 uint32                      `json:"max_response_delay,omitempty"`
	MaxUnackedUpdates                uint32                      `json:"max_unacked_updates,omitempty"`
	MsAssociationMode                string                      `json:"ms_association_mode,omitempty"`
	MsEnableAuthentication           bool                        `json:"ms_enable_authentication,omitempty"`
	MsEnableSwitchoverInterval       bool                        `json:"ms_enable_switchover_interval,omitempty"`
	MsFailoverMode                   string                      `json:"ms_failover_mode,omitempty"`
	MsFailoverPartner                string                      `json:"ms_failover_partner,omitempty"`
	MsHotstandbyPartnerRole          string                      `json:"ms_hotstandby_partner_role,omitempty"`
	MsIsConflict                     bool                        `json:"ms_is_conflict,omitempty"`
	MsPreviousState                  string                      `json:"ms_previous_state,omitempty"`
	MsServer                         string                      `json:"ms_server,omitempty"`
	MsSharedSecret                   string                      `json:"ms_shared_secret,omitempty"`
	MsState                          string                      `json:"ms_state,omitempty"`
	MsSwitchoverInterval             uint32                      `json:"ms_switchover_interval,omitempty"`
	Name                             string                      `json:"name,omitempty"`
	Primary                          string                      `json:"primary,omitempty"`
	PrimaryServerType                string                      `json:"primary_server_type,omitempty"`
	PrimaryState                     string                      `json:"primary_state,omitempty"`
	RecycleLeases                    bool                        `json:"recycle_leases,omitempty"`
	Secondary                        string                      `json:"secondary,omitempty"`
	SecondaryServerType              string                      `json:"secondary_server_type,omitempty"`
	SecondaryState                   string                      `json:"secondary_state,omitempty"`
	SetDhcpFailoverPartnerDown       *Setdhcpfailoverpartnerdown `json:"set_dhcp_failover_partner_down,omitempty"`
	SetDhcpFailoverSecondaryRecovery *Emptyparams                `json:"set_dhcp_failover_secondary_recovery,omitempty"`
	UseFailoverPort                  bool                        `json:"use_failover_port,omitempty"`
	UseMsSwitchoverInterval          bool                        `json:"use_ms_switchover_interval,omitempty"`
	UseRecycleLeases                 bool                        `json:"use_recycle_leases,omitempty"`
}

func (Dhcpfailover) ObjectType() string {
	return "dhcpfailover"
}

func (obj Dhcpfailover) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name"}
	}
	return obj.returnFields
}

// Dhcpoptionspace represents Infoblox object dhcpoptionspace
type Dhcpoptionspace struct {
	IBBase            `json:"-"`
	Ref               string   `json:"_ref,omitempty"`
	Comment           string   `json:"comment,omitempty"`
	Name              string   `json:"name,omitempty"`
	OptionDefinitions []string `json:"option_definitions,omitempty"`
	SpaceType         string   `json:"space_type,omitempty"`
}

func (Dhcpoptionspace) ObjectType() string {
	return "dhcpoptionspace"
}

func (obj Dhcpoptionspace) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name"}
	}
	return obj.returnFields
}

// DiscoveryDevicecomponent represents Infoblox object discovery:devicecomponent
type DiscoveryDevicecomponent struct {
	IBBase        `json:"-"`
	Ref           string `json:"_ref,omitempty"`
	ComponentName string `json:"component_name,omitempty"`
	Description   string `json:"description,omitempty"`
	Device        string `json:"device,omitempty"`
	Model         string `json:"model,omitempty"`
	Serial        string `json:"serial,omitempty"`
	Type          string `json:"type,omitempty"`
}

func (DiscoveryDevicecomponent) ObjectType() string {
	return "discovery:devicecomponent"
}

func (obj DiscoveryDevicecomponent) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"component_name", "description", "model", "serial", "type"}
	}
	return obj.returnFields
}

// DiscoveryCredentialgroup represents Infoblox object discovery:credentialgroup
type DiscoveryCredentialgroup struct {
	IBBase `json:"-"`
	Ref    string `json:"_ref,omitempty"`
	Name   string `json:"name,omitempty"`
}

func (DiscoveryCredentialgroup) ObjectType() string {
	return "discovery:credentialgroup"
}

func (obj DiscoveryCredentialgroup) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name"}
	}
	return obj.returnFields
}

// DiscoveryDevice represents Infoblox object discovery:device
type DiscoveryDevice struct {
	IBBase                         `json:"-"`
	Ref                            string                         `json:"_ref,omitempty"`
	Address                        string                         `json:"address,omitempty"`
	AddressRef                     string                         `json:"address_ref,omitempty"`
	AvailableMgmtIps               []string                       `json:"available_mgmt_ips,omitempty"`
	CapAdminStatusInd              bool                           `json:"cap_admin_status_ind,omitempty"`
	CapAdminStatusNaReason         string                         `json:"cap_admin_status_na_reason,omitempty"`
	CapDescriptionInd              bool                           `json:"cap_description_ind,omitempty"`
	CapDescriptionNaReason         string                         `json:"cap_description_na_reason,omitempty"`
	CapNetDeprovisioningInd        bool                           `json:"cap_net_deprovisioning_ind,omitempty"`
	CapNetDeprovisioningNaReason   string                         `json:"cap_net_deprovisioning_na_reason,omitempty"`
	CapNetProvisioningInd          bool                           `json:"cap_net_provisioning_ind,omitempty"`
	CapNetProvisioningNaReason     string                         `json:"cap_net_provisioning_na_reason,omitempty"`
	CapNetVlanProvisioningInd      bool                           `json:"cap_net_vlan_provisioning_ind,omitempty"`
	CapNetVlanProvisioningNaReason string                         `json:"cap_net_vlan_provisioning_na_reason,omitempty"`
	CapVlanAssignmentInd           bool                           `json:"cap_vlan_assignment_ind,omitempty"`
	CapVlanAssignmentNaReason      string                         `json:"cap_vlan_assignment_na_reason,omitempty"`
	CapVoiceVlanInd                bool                           `json:"cap_voice_vlan_ind,omitempty"`
	CapVoiceVlanNaReason           string                         `json:"cap_voice_vlan_na_reason,omitempty"`
	ChassisSerialNumber            string                         `json:"chassis_serial_number,omitempty"`
	Description                    string                         `json:"description,omitempty"`
	Ea                             EA                             `json:"extattrs,omitempty"`
	Interfaces                     []*DiscoveryDeviceinterface    `json:"interfaces,omitempty"`
	Location                       string                         `json:"location,omitempty"`
	Model                          string                         `json:"model,omitempty"`
	MsAdUserData                   *MsserverAduserData            `json:"ms_ad_user_data,omitempty"`
	Name                           string                         `json:"name,omitempty"`
	Neighbors                      []*DiscoveryDeviceneighbor     `json:"neighbors,omitempty"`
	Network                        string                         `json:"network,omitempty"`
	NetworkInfos                   []*DiscoveryNetworkinfo        `json:"network_infos,omitempty"`
	NetworkView                    string                         `json:"network_view,omitempty"`
	Networks                       []*Ipv4Network                 `json:"networks,omitempty"`
	OsVersion                      string                         `json:"os_version,omitempty"`
	PortStats                      *DiscoveryDevicePortstatistics `json:"port_stats,omitempty"`
	PrivilegedPolling              bool                           `json:"privileged_polling,omitempty"`
	Type                           string                         `json:"type,omitempty"`
	UserDefinedMgmtIp              string                         `json:"user_defined_mgmt_ip,omitempty"`
	Vendor                         string                         `json:"vendor,omitempty"`
	VlanInfos                      []*DiscoveryVlaninfo           `json:"vlan_infos,omitempty"`
}

func (DiscoveryDevice) ObjectType() string {
	return "discovery:device"
}

func (obj DiscoveryDevice) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"address", "name", "network_view"}
	}
	return obj.returnFields
}

// Dhcpoptiondefinition represents Infoblox object dhcpoptiondefinition
type Dhcpoptiondefinition struct {
	IBBase `json:"-"`
	Ref    string `json:"_ref,omitempty"`
	Code   uint32 `json:"code,omitempty"`
	Name   string `json:"name,omitempty"`
	Space  string `json:"space,omitempty"`
	Type   string `json:"type,omitempty"`
}

func (Dhcpoptiondefinition) ObjectType() string {
	return "dhcpoptiondefinition"
}

func (obj Dhcpoptiondefinition) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"code", "name", "type"}
	}
	return obj.returnFields
}

// DiscoveryDeviceneighbor represents Infoblox object discovery:deviceneighbor
type DiscoveryDeviceneighbor struct {
	IBBase     `json:"-"`
	Ref        string               `json:"_ref,omitempty"`
	Address    string               `json:"address,omitempty"`
	AddressRef string               `json:"address_ref,omitempty"`
	Device     string               `json:"device,omitempty"`
	Interface  string               `json:"interface,omitempty"`
	Mac        string               `json:"mac,omitempty"`
	Name       string               `json:"name,omitempty"`
	VlanInfos  []*DiscoveryVlaninfo `json:"vlan_infos,omitempty"`
}

func (DiscoveryDeviceneighbor) ObjectType() string {
	return "discovery:deviceneighbor"
}

func (obj DiscoveryDeviceneighbor) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"address", "address_ref", "mac", "name"}
	}
	return obj.returnFields
}

// DiscoveryDeviceinterface represents Infoblox object discovery:deviceinterface
type DiscoveryDeviceinterface struct {
	IBBase                             `json:"-"`
	Ref                                string                          `json:"_ref,omitempty"`
	AdminStatus                        string                          `json:"admin_status,omitempty"`
	AdminStatusTaskInfo                *DiscoveryPortConfigAdminstatus `json:"admin_status_task_info,omitempty"`
	AggrInterfaceName                  string                          `json:"aggr_interface_name,omitempty"`
	CapIfAdminStatusInd                bool                            `json:"cap_if_admin_status_ind,omitempty"`
	CapIfAdminStatusNaReason           string                          `json:"cap_if_admin_status_na_reason,omitempty"`
	CapIfDescriptionInd                bool                            `json:"cap_if_description_ind,omitempty"`
	CapIfDescriptionNaReason           string                          `json:"cap_if_description_na_reason,omitempty"`
	CapIfNetDeprovisioningIpv4Ind      bool                            `json:"cap_if_net_deprovisioning_ipv4_ind,omitempty"`
	CapIfNetDeprovisioningIpv4NaReason string                          `json:"cap_if_net_deprovisioning_ipv4_na_reason,omitempty"`
	CapIfNetDeprovisioningIpv6Ind      bool                            `json:"cap_if_net_deprovisioning_ipv6_ind,omitempty"`
	CapIfNetDeprovisioningIpv6NaReason string                          `json:"cap_if_net_deprovisioning_ipv6_na_reason,omitempty"`
	CapIfNetProvisioningIpv4Ind        bool                            `json:"cap_if_net_provisioning_ipv4_ind,omitempty"`
	CapIfNetProvisioningIpv4NaReason   string                          `json:"cap_if_net_provisioning_ipv4_na_reason,omitempty"`
	CapIfNetProvisioningIpv6Ind        bool                            `json:"cap_if_net_provisioning_ipv6_ind,omitempty"`
	CapIfNetProvisioningIpv6NaReason   string                          `json:"cap_if_net_provisioning_ipv6_na_reason,omitempty"`
	CapIfVlanAssignmentInd             bool                            `json:"cap_if_vlan_assignment_ind,omitempty"`
	CapIfVlanAssignmentNaReason        string                          `json:"cap_if_vlan_assignment_na_reason,omitempty"`
	CapIfVoiceVlanInd                  bool                            `json:"cap_if_voice_vlan_ind,omitempty"`
	CapIfVoiceVlanNaReason             string                          `json:"cap_if_voice_vlan_na_reason,omitempty"`
	Description                        string                          `json:"description,omitempty"`
	DescriptionTaskInfo                *DiscoveryPortConfigDescription `json:"description_task_info,omitempty"`
	Device                             string                          `json:"device,omitempty"`
	Duplex                             string                          `json:"duplex,omitempty"`
	Ea                                 EA                              `json:"extattrs,omitempty"`
	IfaddrInfos                        []*DiscoveryIfaddrinfo          `json:"ifaddr_infos,omitempty"`
	Index                              int                             `json:"index,omitempty"`
	LastChange                         *time.Time                      `json:"last_change,omitempty"`
	LinkAggregation                    bool                            `json:"link_aggregation,omitempty"`
	Mac                                string                          `json:"mac,omitempty"`
	MsAdUserData                       *MsserverAduserData             `json:"ms_ad_user_data,omitempty"`
	Name                               string                          `json:"name,omitempty"`
	NetworkView                        string                          `json:"network_view,omitempty"`
	OperStatus                         string                          `json:"oper_status,omitempty"`
	PortFast                           string                          `json:"port_fast,omitempty"`
	ReservedObject                     string                          `json:"reserved_object,omitempty"`
	Speed                              uint32                          `json:"speed,omitempty"`
	TrunkStatus                        string                          `json:"trunk_status,omitempty"`
	Type                               string                          `json:"type,omitempty"`
	VlanInfoTaskInfo                   *DiscoveryPortConfigVlaninfo    `json:"vlan_info_task_info,omitempty"`
	VlanInfos                          []*DiscoveryVlaninfo            `json:"vlan_infos,omitempty"`
	VpcPeer                            string                          `json:"vpc_peer,omitempty"`
	VpcPeerDevice                      string                          `json:"vpc_peer_device,omitempty"`
	VrfDescription                     string                          `json:"vrf_description,omitempty"`
	VrfName                            string                          `json:"vrf_name,omitempty"`
	VrfRd                              string                          `json:"vrf_rd,omitempty"`
}

func (DiscoveryDeviceinterface) ObjectType() string {
	return "discovery:deviceinterface"
}

func (obj DiscoveryDeviceinterface) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name", "type"}
	}
	return obj.returnFields
}

// DiscoveryDevicesupportbundle represents Infoblox object discovery:devicesupportbundle
type DiscoveryDevicesupportbundle struct {
	IBBase        `json:"-"`
	Ref           string `json:"_ref,omitempty"`
	Author        string `json:"author,omitempty"`
	IntegratedInd bool   `json:"integrated_ind,omitempty"`
	Name          string `json:"name,omitempty"`
	Version       string `json:"version,omitempty"`
}

func (DiscoveryDevicesupportbundle) ObjectType() string {
	return "discovery:devicesupportbundle"
}

func (obj DiscoveryDevicesupportbundle) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"author", "integrated_ind", "name", "version"}
	}
	return obj.returnFields
}

// Discovery represents Infoblox object discovery
type Discovery struct {
	IBBase                     `json:"-"`
	Ref                        string                      `json:"_ref,omitempty"`
	ClearNetworkPortAssignment *Clearnetworkportassignment `json:"clear_network_port_assignment,omitempty"`
	ControlSwitchPort          *Controlswitchport          `json:"control_switch_port,omitempty"`
	DiscoveryDataConversion    *Discoverydataconversion    `json:"discovery_data_conversion,omitempty"`
	GetDeviceSupportInfo       *Getdevicesupportinfo       `json:"get_device_support_info,omitempty"`
	GetJobDevices              *Getjobdevices              `json:"get_job_devices,omitempty"`
	GetJobProcessDetails       *Getjobprocessdetails       `json:"get_job_process_details,omitempty"`
	ImportDeviceSupportBundle  *Importdevicesupportbundle  `json:"import_device_support_bundle,omitempty"`
	ModifySdnAssignment        *Modifysdnnetworkassignment `json:"modify_sdn_assignment,omitempty"`
	ModifyVrfAssignment        *Modifyvrfassignment        `json:"modify_vrf_assignment,omitempty"`
	ProvisionNetworkDhcpRelay  *Provisionnetworkdhcprelay  `json:"provision_network_dhcp_relay,omitempty"`
	ProvisionNetworkPort       *Provisionnetworkport       `json:"provision_network_port,omitempty"`
}

func (Discovery) ObjectType() string {
	return "discovery"
}

func (obj Discovery) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{}
	}
	return obj.returnFields
}

// DiscoveryDiagnostictask represents Infoblox object discovery:diagnostictask
type DiscoveryDiagnostictask struct {
	IBBase          `json:"-"`
	Ref             string     `json:"_ref,omitempty"`
	CommunityString string     `json:"community_string,omitempty"`
	DebugSnmp       bool       `json:"debug_snmp,omitempty"`
	ForceTest       bool       `json:"force_test,omitempty"`
	IpAddress       string     `json:"ip_address,omitempty"`
	NetworkView     string     `json:"network_view,omitempty"`
	StartTime       *time.Time `json:"start_time,omitempty"`
	TaskId          string     `json:"task_id,omitempty"`
}

func (DiscoveryDiagnostictask) ObjectType() string {
	return "discovery:diagnostictask"
}

func (obj DiscoveryDiagnostictask) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"ip_address", "network_view", "task_id"}
	}
	return obj.returnFields
}

// DiscoveryStatus represents Infoblox object discovery:status
type DiscoveryStatus struct {
	IBBase                `json:"-"`
	Ref                   string               `json:"_ref,omitempty"`
	Address               string               `json:"address,omitempty"`
	CliCollectionEnabled  bool                 `json:"cli_collection_enabled,omitempty"`
	CliCredentialInfo     *DiscoveryStatusinfo `json:"cli_credential_info,omitempty"`
	ExistenceInfo         *DiscoveryStatusinfo `json:"existence_info,omitempty"`
	FingerprintEnabled    bool                 `json:"fingerprint_enabled,omitempty"`
	FingerprintInfo       *DiscoveryStatusinfo `json:"fingerprint_info,omitempty"`
	FirstSeen             *time.Time           `json:"first_seen,omitempty"`
	LastAction            string               `json:"last_action,omitempty"`
	LastSeen              *time.Time           `json:"last_seen,omitempty"`
	LastTimestamp         *time.Time           `json:"last_timestamp,omitempty"`
	Name                  string               `json:"name,omitempty"`
	NetworkView           string               `json:"network_view,omitempty"`
	ReachableInfo         *DiscoveryStatusinfo `json:"reachable_info,omitempty"`
	SdnCollectionEnabled  bool                 `json:"sdn_collection_enabled,omitempty"`
	SdnCollectionInfo     *DiscoveryStatusinfo `json:"sdn_collection_info,omitempty"`
	SnmpCollectionEnabled bool                 `json:"snmp_collection_enabled,omitempty"`
	SnmpCollectionInfo    *DiscoveryStatusinfo `json:"snmp_collection_info,omitempty"`
	SnmpCredentialInfo    *DiscoveryStatusinfo `json:"snmp_credential_info,omitempty"`
	Status                string               `json:"status,omitempty"`
	Type                  string               `json:"type,omitempty"`
}

func (DiscoveryStatus) ObjectType() string {
	return "discovery:status"
}

func (obj DiscoveryStatus) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"address", "name", "network_view", "status"}
	}
	return obj.returnFields
}

// DiscoverySdnnetwork represents Infoblox object discovery:sdnnetwork
type DiscoverySdnnetwork struct {
	IBBase          `json:"-"`
	Ref             string     `json:"_ref,omitempty"`
	FirstSeen       *time.Time `json:"first_seen,omitempty"`
	Name            string     `json:"name,omitempty"`
	NetworkView     string     `json:"network_view,omitempty"`
	SourceSdnConfig string     `json:"source_sdn_config,omitempty"`
}

func (DiscoverySdnnetwork) ObjectType() string {
	return "discovery:sdnnetwork"
}

func (obj DiscoverySdnnetwork) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name", "network_view", "source_sdn_config"}
	}
	return obj.returnFields
}

// Discoverytask represents Infoblox object discoverytask
type Discoverytask struct {
	IBBase                  `json:"-"`
	Ref                     string                         `json:"_ref,omitempty"`
	CsvFileName             string                         `json:"csv_file_name,omitempty"`
	DisableIpScanning       bool                           `json:"disable_ip_scanning,omitempty"`
	DisableVmwareScanning   bool                           `json:"disable_vmware_scanning,omitempty"`
	DiscoveryTaskOid        string                         `json:"discovery_task_oid,omitempty"`
	MemberName              string                         `json:"member_name,omitempty"`
	MergeData               bool                           `json:"merge_data,omitempty"`
	Mode                    string                         `json:"mode,omitempty"`
	NetworkDiscoveryControl *Networkdiscoverycontrolparams `json:"network_discovery_control,omitempty"`
	NetworkView             string                         `json:"network_view,omitempty"`
	Networks                []*Ipv4Network                 `json:"networks,omitempty"`
	PingRetries             uint32                         `json:"ping_retries,omitempty"`
	PingTimeout             uint32                         `json:"ping_timeout,omitempty"`
	ScheduledRun            *SettingSchedule               `json:"scheduled_run,omitempty"`
	State                   string                         `json:"state,omitempty"`
	StateTime               *time.Time                     `json:"state_time,omitempty"`
	Status                  string                         `json:"status,omitempty"`
	StatusTime              *time.Time                     `json:"status_time,omitempty"`
	TcpPorts                []*Discoverytaskport           `json:"tcp_ports,omitempty"`
	TcpScanTechnique        string                         `json:"tcp_scan_technique,omitempty"`
	VNetworkView            string                         `json:"v_network_view,omitempty"`
	Vservers                []*Discoverytaskvserver        `json:"vservers,omitempty"`
	Warning                 string                         `json:"warning,omitempty"`
}

func (Discoverytask) ObjectType() string {
	return "discoverytask"
}

func (obj Discoverytask) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"discovery_task_oid", "member_name"}
	}
	return obj.returnFields
}

// Distributionschedule represents Infoblox object distributionschedule
type Distributionschedule struct {
	IBBase        `json:"-"`
	Ref           string                  `json:"_ref,omitempty"`
	Active        bool                    `json:"active,omitempty"`
	StartTime     *time.Time              `json:"start_time,omitempty"`
	TimeZone      string                  `json:"time_zone,omitempty"`
	UpgradeGroups []*UpgradegroupSchedule `json:"upgrade_groups,omitempty"`
}

func (Distributionschedule) ObjectType() string {
	return "distributionschedule"
}

func (obj Distributionschedule) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"active", "start_time", "time_zone"}
	}
	return obj.returnFields
}

// DiscoveryVrf represents Infoblox object discovery:vrf
type DiscoveryVrf struct {
	IBBase             `json:"-"`
	Ref                string `json:"_ref,omitempty"`
	Description        string `json:"description,omitempty"`
	Device             string `json:"device,omitempty"`
	Name               string `json:"name,omitempty"`
	NetworkView        string `json:"network_view,omitempty"`
	RouteDistinguisher string `json:"route_distinguisher,omitempty"`
}

func (DiscoveryVrf) ObjectType() string {
	return "discovery:vrf"
}

func (obj DiscoveryVrf) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"device", "name", "network_view", "route_distinguisher"}
	}
	return obj.returnFields
}

// DiscoveryMemberproperties represents Infoblox object discovery:memberproperties
type DiscoveryMemberproperties struct {
	IBBase                 `json:"-"`
	Ref                    string                      `json:"_ref,omitempty"`
	Address                string                      `json:"address,omitempty"`
	CliCredentials         []*DiscoveryClicredential   `json:"cli_credentials,omitempty"`
	DefaultSeedRouters     []*DiscoverySeedrouter      `json:"default_seed_routers,omitempty"`
	DiscoveryMember        string                      `json:"discovery_member,omitempty"`
	EnableService          bool                        `json:"enable_service,omitempty"`
	GatewaySeedRouters     []*DiscoverySeedrouter      `json:"gateway_seed_routers,omitempty"`
	IsSa                   bool                        `json:"is_sa,omitempty"`
	Role                   string                      `json:"role,omitempty"`
	ScanInterfaces         []*DiscoveryScaninterface   `json:"scan_interfaces,omitempty"`
	SdnConfigs             []*DiscoverySdnconfig       `json:"sdn_configs,omitempty"`
	SeedRouters            []*DiscoverySeedrouter      `json:"seed_routers,omitempty"`
	Snmpv1v2Credentials    []*DiscoverySnmpcredential  `json:"snmpv1v2_credentials,omitempty"`
	Snmpv3Credentials      []*DiscoverySnmp3credential `json:"snmpv3_credentials,omitempty"`
	UseCliCredentials      bool                        `json:"use_cli_credentials,omitempty"`
	UseSnmpv1v2Credentials bool                        `json:"use_snmpv1v2_credentials,omitempty"`
	UseSnmpv3Credentials   bool                        `json:"use_snmpv3_credentials,omitempty"`
}

func (DiscoveryMemberproperties) ObjectType() string {
	return "discovery:memberproperties"
}

func (obj DiscoveryMemberproperties) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"discovery_member"}
	}
	return obj.returnFields
}

// Dns64group represents Infoblox object dns64group
type Dns64group struct {
	IBBase            `json:"-"`
	Ref               string       `json:"_ref,omitempty"`
	Clients           []*Addressac `json:"clients,omitempty"`
	Comment           string       `json:"comment,omitempty"`
	Disable           bool         `json:"disable,omitempty"`
	EnableDnssecDns64 bool         `json:"enable_dnssec_dns64,omitempty"`
	Exclude           []*Addressac `json:"exclude,omitempty"`
	Ea                EA           `json:"extattrs,omitempty"`
	Mapped            []*Addressac `json:"mapped,omitempty"`
	Name              string       `json:"name,omitempty"`
	Prefix            string       `json:"prefix,omitempty"`
}

func (Dns64group) ObjectType() string {
	return "dns64group"
}

func (obj Dns64group) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "disable", "name"}
	}
	return obj.returnFields
}

// Dtc represents Infoblox object dtc
type Dtc struct {
	IBBase               `json:"-"`
	Ref                  string           `json:"_ref,omitempty"`
	AddCertificate       *Addcertificate  `json:"add_certificate,omitempty"`
	GenerateEaTopologyDb *Emptyparams     `json:"generate_ea_topology_db,omitempty"`
	ImportMaxminddb      *Importmaxminddb `json:"import_maxminddb,omitempty"`
	Query                *Query           `json:"query,omitempty"`
}

func (Dtc) ObjectType() string {
	return "dtc"
}

func (obj Dtc) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{}
	}
	return obj.returnFields
}

// DtcAllrecords represents Infoblox object dtc:allrecords
type DtcAllrecords struct {
	IBBase    `json:"-"`
	Ref       string `json:"_ref,omitempty"`
	Comment   string `json:"comment,omitempty"`
	Disable   bool   `json:"disable,omitempty"`
	DtcServer string `json:"dtc_server,omitempty"`
	Record    string `json:"record,omitempty"`
	Ttl       uint32 `json:"ttl,omitempty"`
	Type      string `json:"type,omitempty"`
}

func (DtcAllrecords) ObjectType() string {
	return "dtc:allrecords"
}

func (obj DtcAllrecords) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "dtc_server", "type"}
	}
	return obj.returnFields
}

// DiscoveryGridproperties represents Infoblox object discovery:gridproperties
type DiscoveryGridproperties struct {
	IBBase                           `json:"-"`
	Ref                              string                            `json:"_ref,omitempty"`
	AdvancedPollingSettings          *DiscoveryAdvancedpollsetting     `json:"advanced_polling_settings,omitempty"`
	AdvancedSdnPollingSettings       *DiscoveryAdvancedsdnpollsettings `json:"advanced_sdn_polling_settings,omitempty"`
	AdvisorRunNow                    *Runnow                           `json:"advisor_run_now,omitempty"`
	AdvisorSettings                  *DiscoveryAdvisorsetting          `json:"advisor_settings,omitempty"`
	AdvisorTestConnection            *Advisortestconnection            `json:"advisor_test_connection,omitempty"`
	AutoConversionSettings           []*DiscoveryAutoconversionsetting `json:"auto_conversion_settings,omitempty"`
	BasicPollingSettings             *DiscoveryBasicpollsettings       `json:"basic_polling_settings,omitempty"`
	BasicSdnPollingSettings          *DiscoveryBasicsdnpollsettings    `json:"basic_sdn_polling_settings,omitempty"`
	CliCredentials                   []*DiscoveryClicredential         `json:"cli_credentials,omitempty"`
	Diagnostic                       *Discoverydiagnostic              `json:"diagnostic,omitempty"`
	DiagnosticStatus                 *Discoverydiagnosticstatus        `json:"diagnostic_status,omitempty"`
	DiscoveryBlackoutSetting         *PropertiesBlackoutsetting        `json:"discovery_blackout_setting,omitempty"`
	DnsLookupOption                  string                            `json:"dns_lookup_option,omitempty"`
	DnsLookupThrottle                uint32                            `json:"dns_lookup_throttle,omitempty"`
	EnableAdvisor                    bool                              `json:"enable_advisor,omitempty"`
	EnableAutoConversion             bool                              `json:"enable_auto_conversion,omitempty"`
	EnableAutoUpdates                bool                              `json:"enable_auto_updates,omitempty"`
	GridName                         string                            `json:"grid_name,omitempty"`
	IgnoreConflictDuration           uint32                            `json:"ignore_conflict_duration,omitempty"`
	PortControlBlackoutSetting       *PropertiesBlackoutsetting        `json:"port_control_blackout_setting,omitempty"`
	Ports                            []*DiscoveryPort                  `json:"ports,omitempty"`
	SamePortControlDiscoveryBlackout bool                              `json:"same_port_control_discovery_blackout,omitempty"`
	Snmpv1v2Credentials              []*DiscoverySnmpcredential        `json:"snmpv1v2_credentials,omitempty"`
	Snmpv3Credentials                []*DiscoverySnmp3credential       `json:"snmpv3_credentials,omitempty"`
	UnmanagedIpsLimit                uint32                            `json:"unmanaged_ips_limit,omitempty"`
	UnmanagedIpsTimeout              uint32                            `json:"unmanaged_ips_timeout,omitempty"`
	VrfMappingPolicy                 string                            `json:"vrf_mapping_policy,omitempty"`
	VrfMappingRules                  []*DiscoveryVrfmappingrule        `json:"vrf_mapping_rules,omitempty"`
}

func (DiscoveryGridproperties) ObjectType() string {
	return "discovery:gridproperties"
}

func (obj DiscoveryGridproperties) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"grid_name"}
	}
	return obj.returnFields
}

// DtcCertificate represents Infoblox object dtc:certificate
type DtcCertificate struct {
	IBBase      `json:"-"`
	Ref         string `json:"_ref,omitempty"`
	Certificate string `json:"certificate,omitempty"`
	InUse       bool   `json:"in_use,omitempty"`
}

func (DtcCertificate) ObjectType() string {
	return "dtc:certificate"
}

func (obj DtcCertificate) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{}
	}
	return obj.returnFields
}

// DtcLbdn represents Infoblox object dtc:lbdn
type DtcLbdn struct {
	IBBase                   `json:"-"`
	Ref                      string         `json:"_ref,omitempty"`
	AuthZones                []*ZoneAuth    `json:"auth_zones,omitempty"`
	AutoConsolidatedMonitors bool           `json:"auto_consolidated_monitors,omitempty"`
	Comment                  string         `json:"comment,omitempty"`
	Disable                  bool           `json:"disable,omitempty"`
	Ea                       EA             `json:"extattrs,omitempty"`
	Health                   *DtcHealth     `json:"health,omitempty"`
	LbMethod                 string         `json:"lb_method,omitempty"`
	Name                     string         `json:"name,omitempty"`
	Patterns                 []string       `json:"patterns,omitempty"`
	Persistence              uint32         `json:"persistence,omitempty"`
	Pools                    []*DtcPoolLink `json:"pools,omitempty"`
	Priority                 uint32         `json:"priority,omitempty"`
	Topology                 string         `json:"topology,omitempty"`
	Ttl                      uint32         `json:"ttl,omitempty"`
	Types                    []string       `json:"types,omitempty"`
	UseTtl                   bool           `json:"use_ttl,omitempty"`
}

func (DtcLbdn) ObjectType() string {
	return "dtc:lbdn"
}

func (obj DtcLbdn) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name"}
	}
	return obj.returnFields
}

// DtcMonitor represents Infoblox object dtc:monitor
type DtcMonitor struct {
	IBBase    `json:"-"`
	Ref       string `json:"_ref,omitempty"`
	Comment   string `json:"comment,omitempty"`
	Ea        EA     `json:"extattrs,omitempty"`
	Interval  uint32 `json:"interval,omitempty"`
	Monitor   string `json:"monitor,omitempty"`
	Name      string `json:"name,omitempty"`
	Port      uint32 `json:"port,omitempty"`
	RetryDown uint32 `json:"retry_down,omitempty"`
	RetryUp   uint32 `json:"retry_up,omitempty"`
	Timeout   uint32 `json:"timeout,omitempty"`
	Type      string `json:"type,omitempty"`
}

func (DtcMonitor) ObjectType() string {
	return "dtc:monitor"
}

func (obj DtcMonitor) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name", "type"}
	}
	return obj.returnFields
}

// DtcMonitorIcmp represents Infoblox object dtc:monitor:icmp
type DtcMonitorIcmp struct {
	IBBase    `json:"-"`
	Ref       string `json:"_ref,omitempty"`
	Comment   string `json:"comment,omitempty"`
	Ea        EA     `json:"extattrs,omitempty"`
	Interval  uint32 `json:"interval,omitempty"`
	Name      string `json:"name,omitempty"`
	RetryDown uint32 `json:"retry_down,omitempty"`
	RetryUp   uint32 `json:"retry_up,omitempty"`
	Timeout   uint32 `json:"timeout,omitempty"`
}

func (DtcMonitorIcmp) ObjectType() string {
	return "dtc:monitor:icmp"
}

func (obj DtcMonitorIcmp) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name"}
	}
	return obj.returnFields
}

// DtcMonitorHttp represents Infoblox object dtc:monitor:http
type DtcMonitorHttp struct {
	IBBase              `json:"-"`
	Ref                 string `json:"_ref,omitempty"`
	Ciphers             string `json:"ciphers,omitempty"`
	ClientCert          string `json:"client_cert,omitempty"`
	Comment             string `json:"comment,omitempty"`
	ContentCheck        string `json:"content_check,omitempty"`
	ContentCheckInput   string `json:"content_check_input,omitempty"`
	ContentCheckOp      string `json:"content_check_op,omitempty"`
	ContentCheckRegex   string `json:"content_check_regex,omitempty"`
	ContentExtractGroup uint32 `json:"content_extract_group,omitempty"`
	ContentExtractType  string `json:"content_extract_type,omitempty"`
	ContentExtractValue string `json:"content_extract_value,omitempty"`
	EnableSni           bool   `json:"enable_sni,omitempty"`
	Ea                  EA     `json:"extattrs,omitempty"`
	Interval            uint32 `json:"interval,omitempty"`
	Name                string `json:"name,omitempty"`
	Port                uint32 `json:"port,omitempty"`
	Request             string `json:"request,omitempty"`
	Result              string `json:"result,omitempty"`
	ResultCode          uint32 `json:"result_code,omitempty"`
	RetryDown           uint32 `json:"retry_down,omitempty"`
	RetryUp             uint32 `json:"retry_up,omitempty"`
	Secure              bool   `json:"secure,omitempty"`
	Timeout             uint32 `json:"timeout,omitempty"`
	ValidateCert        bool   `json:"validate_cert,omitempty"`
}

func (DtcMonitorHttp) ObjectType() string {
	return "dtc:monitor:http"
}

func (obj DtcMonitorHttp) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name"}
	}
	return obj.returnFields
}

// DtcMonitorTcp represents Infoblox object dtc:monitor:tcp
type DtcMonitorTcp struct {
	IBBase    `json:"-"`
	Ref       string `json:"_ref,omitempty"`
	Comment   string `json:"comment,omitempty"`
	Ea        EA     `json:"extattrs,omitempty"`
	Interval  uint32 `json:"interval,omitempty"`
	Name      string `json:"name,omitempty"`
	Port      uint32 `json:"port,omitempty"`
	RetryDown uint32 `json:"retry_down,omitempty"`
	RetryUp   uint32 `json:"retry_up,omitempty"`
	Timeout   uint32 `json:"timeout,omitempty"`
}

func (DtcMonitorTcp) ObjectType() string {
	return "dtc:monitor:tcp"
}

func (obj DtcMonitorTcp) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name"}
	}
	return obj.returnFields
}

// DtcMonitorSnmp represents Infoblox object dtc:monitor:snmp
type DtcMonitorSnmp struct {
	IBBase    `json:"-"`
	Ref       string               `json:"_ref,omitempty"`
	Comment   string               `json:"comment,omitempty"`
	Community string               `json:"community,omitempty"`
	Context   string               `json:"context,omitempty"`
	EngineId  string               `json:"engine_id,omitempty"`
	Ea        EA                   `json:"extattrs,omitempty"`
	Interval  uint32               `json:"interval,omitempty"`
	Name      string               `json:"name,omitempty"`
	Oids      []*DtcMonitorSnmpOid `json:"oids,omitempty"`
	Port      uint32               `json:"port,omitempty"`
	RetryDown uint32               `json:"retry_down,omitempty"`
	RetryUp   uint32               `json:"retry_up,omitempty"`
	Timeout   uint32               `json:"timeout,omitempty"`
	User      string               `json:"user,omitempty"`
	Version   string               `json:"version,omitempty"`
}

func (DtcMonitorSnmp) ObjectType() string {
	return "dtc:monitor:snmp"
}

func (obj DtcMonitorSnmp) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name"}
	}
	return obj.returnFields
}

// DtcObject represents Infoblox object dtc:object
type DtcObject struct {
	IBBase          `json:"-"`
	Ref             string     `json:"_ref,omitempty"`
	AbstractType    string     `json:"abstract_type,omitempty"`
	Comment         string     `json:"comment,omitempty"`
	DisplayType     string     `json:"display_type,omitempty"`
	Ea              EA         `json:"extattrs,omitempty"`
	Ipv4AddressList []string   `json:"ipv4_address_list,omitempty"`
	Ipv6AddressList []string   `json:"ipv6_address_list,omitempty"`
	Name            string     `json:"name,omitempty"`
	Object          string     `json:"object,omitempty"`
	Status          string     `json:"status,omitempty"`
	StatusTime      *time.Time `json:"status_time,omitempty"`
}

func (DtcObject) ObjectType() string {
	return "dtc:object"
}

func (obj DtcObject) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"abstract_type", "comment", "display_type", "name", "status"}
	}
	return obj.returnFields
}

// DtcPool represents Infoblox object dtc:pool
type DtcPool struct {
	IBBase                   `json:"-"`
	Ref                      string                              `json:"_ref,omitempty"`
	AutoConsolidatedMonitors bool                                `json:"auto_consolidated_monitors,omitempty"`
	Availability             string                              `json:"availability,omitempty"`
	Comment                  string                              `json:"comment,omitempty"`
	ConsolidatedMonitors     []*DtcPoolConsolidatedMonitorHealth `json:"consolidated_monitors,omitempty"`
	Disable                  bool                                `json:"disable,omitempty"`
	Ea                       EA                                  `json:"extattrs,omitempty"`
	Health                   *DtcHealth                          `json:"health,omitempty"`
	LbAlternateMethod        string                              `json:"lb_alternate_method,omitempty"`
	LbAlternateTopology      string                              `json:"lb_alternate_topology,omitempty"`
	LbDynamicRatioAlternate  *SettingDynamicratio                `json:"lb_dynamic_ratio_alternate,omitempty"`
	LbDynamicRatioPreferred  *SettingDynamicratio                `json:"lb_dynamic_ratio_preferred,omitempty"`
	LbPreferredMethod        string                              `json:"lb_preferred_method,omitempty"`
	LbPreferredTopology      string                              `json:"lb_preferred_topology,omitempty"`
	Monitors                 []*DtcMonitorHttp                   `json:"monitors,omitempty"`
	Name                     string                              `json:"name,omitempty"`
	Quorum                   uint32                              `json:"quorum,omitempty"`
	Servers                  []*DtcServerLink                    `json:"servers,omitempty"`
	Ttl                      uint32                              `json:"ttl,omitempty"`
	UseTtl                   bool                                `json:"use_ttl,omitempty"`
}

func (DtcPool) ObjectType() string {
	return "dtc:pool"
}

func (obj DtcPool) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name"}
	}
	return obj.returnFields
}

// DtcRecordA represents Infoblox object dtc:record:a
type DtcRecordA struct {
	IBBase      `json:"-"`
	Ref         string `json:"_ref,omitempty"`
	AutoCreated string `json:"auto_created,omitempty"`
	Comment     string `json:"comment,omitempty"`
	Disable     bool   `json:"disable,omitempty"`
	DtcServer   string `json:"dtc_server,omitempty"`
	Ipv4Addr    string `json:"ipv4addr,omitempty"`
	Ttl         uint32 `json:"ttl,omitempty"`
	UseTtl      bool   `json:"use_ttl,omitempty"`
}

func (DtcRecordA) ObjectType() string {
	return "dtc:record:a"
}

func (obj DtcRecordA) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"dtc_server", "ipv4addr"}
	}
	return obj.returnFields
}

// DtcRecordAaaa represents Infoblox object dtc:record:aaaa
type DtcRecordAaaa struct {
	IBBase      `json:"-"`
	Ref         string `json:"_ref,omitempty"`
	AutoCreated string `json:"auto_created,omitempty"`
	Comment     string `json:"comment,omitempty"`
	Disable     bool   `json:"disable,omitempty"`
	DtcServer   string `json:"dtc_server,omitempty"`
	Ipv6Addr    string `json:"ipv6addr,omitempty"`
	Ttl         uint32 `json:"ttl,omitempty"`
	UseTtl      bool   `json:"use_ttl,omitempty"`
}

func (DtcRecordAaaa) ObjectType() string {
	return "dtc:record:aaaa"
}

func (obj DtcRecordAaaa) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"dtc_server", "ipv6addr"}
	}
	return obj.returnFields
}

// DtcMonitorPdp represents Infoblox object dtc:monitor:pdp
type DtcMonitorPdp struct {
	IBBase    `json:"-"`
	Ref       string `json:"_ref,omitempty"`
	Comment   string `json:"comment,omitempty"`
	Ea        EA     `json:"extattrs,omitempty"`
	Interval  uint32 `json:"interval,omitempty"`
	Name      string `json:"name,omitempty"`
	Port      uint32 `json:"port,omitempty"`
	RetryDown uint32 `json:"retry_down,omitempty"`
	RetryUp   uint32 `json:"retry_up,omitempty"`
	Timeout   uint32 `json:"timeout,omitempty"`
}

func (DtcMonitorPdp) ObjectType() string {
	return "dtc:monitor:pdp"
}

func (obj DtcMonitorPdp) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name"}
	}
	return obj.returnFields
}

// DtcRecordCname represents Infoblox object dtc:record:cname
type DtcRecordCname struct {
	IBBase       `json:"-"`
	Ref          string `json:"_ref,omitempty"`
	AutoCreated  string `json:"auto_created,omitempty"`
	Canonical    string `json:"canonical,omitempty"`
	Comment      string `json:"comment,omitempty"`
	Disable      bool   `json:"disable,omitempty"`
	DnsCanonical string `json:"dns_canonical,omitempty"`
	DtcServer    string `json:"dtc_server,omitempty"`
	Ttl          uint32 `json:"ttl,omitempty"`
	UseTtl       bool   `json:"use_ttl,omitempty"`
}

func (DtcRecordCname) ObjectType() string {
	return "dtc:record:cname"
}

func (obj DtcRecordCname) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"canonical", "dtc_server"}
	}
	return obj.returnFields
}

// DtcServer represents Infoblox object dtc:server
type DtcServer struct {
	IBBase               `json:"-"`
	Ref                  string              `json:"_ref,omitempty"`
	AutoCreateHostRecord bool                `json:"auto_create_host_record,omitempty"`
	Comment              string              `json:"comment,omitempty"`
	Disable              bool                `json:"disable,omitempty"`
	Ea                   EA                  `json:"extattrs,omitempty"`
	Health               *DtcHealth          `json:"health,omitempty"`
	Host                 string              `json:"host,omitempty"`
	Monitors             []*DtcServerMonitor `json:"monitors,omitempty"`
	Name                 string              `json:"name,omitempty"`
	SniHostname          string              `json:"sni_hostname,omitempty"`
	UseSniHostname       bool                `json:"use_sni_hostname,omitempty"`
}

func (DtcServer) ObjectType() string {
	return "dtc:server"
}

func (obj DtcServer) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "host", "name"}
	}
	return obj.returnFields
}

// DtcMonitorSip represents Infoblox object dtc:monitor:sip
type DtcMonitorSip struct {
	IBBase       `json:"-"`
	Ref          string `json:"_ref,omitempty"`
	Ciphers      string `json:"ciphers,omitempty"`
	ClientCert   string `json:"client_cert,omitempty"`
	Comment      string `json:"comment,omitempty"`
	Ea           EA     `json:"extattrs,omitempty"`
	Interval     uint32 `json:"interval,omitempty"`
	Name         string `json:"name,omitempty"`
	Port         uint32 `json:"port,omitempty"`
	Request      string `json:"request,omitempty"`
	Result       string `json:"result,omitempty"`
	ResultCode   uint32 `json:"result_code,omitempty"`
	RetryDown    uint32 `json:"retry_down,omitempty"`
	RetryUp      uint32 `json:"retry_up,omitempty"`
	Timeout      uint32 `json:"timeout,omitempty"`
	Transport    string `json:"transport,omitempty"`
	ValidateCert bool   `json:"validate_cert,omitempty"`
}

func (DtcMonitorSip) ObjectType() string {
	return "dtc:monitor:sip"
}

func (obj DtcMonitorSip) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name"}
	}
	return obj.returnFields
}

// DtcRecordSrv represents Infoblox object dtc:record:srv
type DtcRecordSrv struct {
	IBBase    `json:"-"`
	Ref       string `json:"_ref,omitempty"`
	Comment   string `json:"comment,omitempty"`
	Disable   bool   `json:"disable,omitempty"`
	DtcServer string `json:"dtc_server,omitempty"`
	Name      string `json:"name,omitempty"`
	Port      uint32 `json:"port,omitempty"`
	Priority  uint32 `json:"priority,omitempty"`
	Target    string `json:"target,omitempty"`
	Ttl       uint32 `json:"ttl,omitempty"`
	UseTtl    bool   `json:"use_ttl,omitempty"`
	Weight    uint32 `json:"weight,omitempty"`
}

func (DtcRecordSrv) ObjectType() string {
	return "dtc:record:srv"
}

func (obj DtcRecordSrv) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"dtc_server", "name", "port", "priority", "target", "weight"}
	}
	return obj.returnFields
}

// DtcRecordNaptr represents Infoblox object dtc:record:naptr
type DtcRecordNaptr struct {
	IBBase      `json:"-"`
	Ref         string `json:"_ref,omitempty"`
	Comment     string `json:"comment,omitempty"`
	Disable     bool   `json:"disable,omitempty"`
	DtcServer   string `json:"dtc_server,omitempty"`
	Flags       string `json:"flags,omitempty"`
	Order       uint32 `json:"order,omitempty"`
	Preference  uint32 `json:"preference,omitempty"`
	Regexp      string `json:"regexp,omitempty"`
	Replacement string `json:"replacement,omitempty"`
	Services    string `json:"services,omitempty"`
	Ttl         uint32 `json:"ttl,omitempty"`
	UseTtl      bool   `json:"use_ttl,omitempty"`
}

func (DtcRecordNaptr) ObjectType() string {
	return "dtc:record:naptr"
}

func (obj DtcRecordNaptr) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"dtc_server", "order", "preference", "regexp", "replacement", "services"}
	}
	return obj.returnFields
}

// DtcTopology represents Infoblox object dtc:topology
type DtcTopology struct {
	IBBase  `json:"-"`
	Ref     string             `json:"_ref,omitempty"`
	Comment string             `json:"comment,omitempty"`
	Ea      EA                 `json:"extattrs,omitempty"`
	Name    string             `json:"name,omitempty"`
	Rules   []*DtcTopologyRule `json:"rules,omitempty"`
}

func (DtcTopology) ObjectType() string {
	return "dtc:topology"
}

func (obj DtcTopology) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name"}
	}
	return obj.returnFields
}

// EADefinition represents Infoblox object extensibleattributedef
type EADefinition struct {
	IBBase             `json:"-"`
	Ref                string                             `json:"_ref,omitempty"`
	AllowedObjectTypes []string                           `json:"allowed_object_types,omitempty"`
	Comment            string                             `json:"comment,omitempty"`
	DefaultValue       string                             `json:"default_value,omitempty"`
	DescendantsAction  *ExtensibleattributedefDescendants `json:"descendants_action,omitempty"`
	Flags              string                             `json:"flags,omitempty"`
	ListValues         []*EADefListValue                  `json:"list_values,omitempty"`
	Max                uint32                             `json:"max,omitempty"`
	Min                uint32                             `json:"min,omitempty"`
	Name               string                             `json:"name,omitempty"`
	Namespace          string                             `json:"namespace,omitempty"`
	Type               string                             `json:"type,omitempty"`
}

func (EADefinition) ObjectType() string {
	return "extensibleattributedef"
}

func (obj EADefinition) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "default_value", "name", "type"}
	}
	return obj.returnFields
}

func NewEADefinition(eadef EADefinition) *EADefinition {
	res := eadef
	res.returnFields = []string{"allowed_object_types", "comment", "flags", "list_values", "name", "type"}

	return &res
}

// Filterfingerprint represents Infoblox object filterfingerprint
type Filterfingerprint struct {
	IBBase      `json:"-"`
	Ref         string   `json:"_ref,omitempty"`
	Comment     string   `json:"comment,omitempty"`
	Ea          EA       `json:"extattrs,omitempty"`
	Fingerprint []string `json:"fingerprint,omitempty"`
	Name        string   `json:"name,omitempty"`
}

func (Filterfingerprint) ObjectType() string {
	return "filterfingerprint"
}

func (obj Filterfingerprint) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name"}
	}
	return obj.returnFields
}

// Filternac represents Infoblox object filternac
type Filternac struct {
	IBBase     `json:"-"`
	Ref        string        `json:"_ref,omitempty"`
	Comment    string        `json:"comment,omitempty"`
	Expression string        `json:"expression,omitempty"`
	Ea         EA            `json:"extattrs,omitempty"`
	LeaseTime  uint32        `json:"lease_time,omitempty"`
	Name       string        `json:"name,omitempty"`
	Options    []*Dhcpoption `json:"options,omitempty"`
}

func (Filternac) ObjectType() string {
	return "filternac"
}

func (obj Filternac) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name"}
	}
	return obj.returnFields
}

// Filtermac represents Infoblox object filtermac
type Filtermac struct {
	IBBase                      `json:"-"`
	Ref                         string        `json:"_ref,omitempty"`
	Comment                     string        `json:"comment,omitempty"`
	DefaultMacAddressExpiration uint32        `json:"default_mac_address_expiration,omitempty"`
	Disable                     bool          `json:"disable,omitempty"`
	EnforceExpirationTimes      bool          `json:"enforce_expiration_times,omitempty"`
	Ea                          EA            `json:"extattrs,omitempty"`
	LeaseTime                   uint32        `json:"lease_time,omitempty"`
	Name                        string        `json:"name,omitempty"`
	NeverExpires                bool          `json:"never_expires,omitempty"`
	Options                     []*Dhcpoption `json:"options,omitempty"`
	ReservedForInfoblox         string        `json:"reserved_for_infoblox,omitempty"`
}

func (Filtermac) ObjectType() string {
	return "filtermac"
}

func (obj Filtermac) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name"}
	}
	return obj.returnFields
}

// DtcTopologyLabel represents Infoblox object dtc:topology:label
type DtcTopologyLabel struct {
	IBBase `json:"-"`
	Ref    string `json:"_ref,omitempty"`
	Field  string `json:"field,omitempty"`
	Label  string `json:"label,omitempty"`
}

func (DtcTopologyLabel) ObjectType() string {
	return "dtc:topology:label"
}

func (obj DtcTopologyLabel) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"field", "label"}
	}
	return obj.returnFields
}

// Filteroption represents Infoblox object filteroption
type Filteroption struct {
	IBBase       `json:"-"`
	Ref          string        `json:"_ref,omitempty"`
	ApplyAsClass bool          `json:"apply_as_class,omitempty"`
	Bootfile     string        `json:"bootfile,omitempty"`
	Bootserver   string        `json:"bootserver,omitempty"`
	Comment      string        `json:"comment,omitempty"`
	Expression   string        `json:"expression,omitempty"`
	Ea           EA            `json:"extattrs,omitempty"`
	LeaseTime    uint32        `json:"lease_time,omitempty"`
	Name         string        `json:"name,omitempty"`
	NextServer   string        `json:"next_server,omitempty"`
	OptionList   []*Dhcpoption `json:"option_list,omitempty"`
	OptionSpace  string        `json:"option_space,omitempty"`
	PxeLeaseTime uint32        `json:"pxe_lease_time,omitempty"`
}

func (Filteroption) ObjectType() string {
	return "filteroption"
}

func (obj Filteroption) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name"}
	}
	return obj.returnFields
}

// Filterrelayagent represents Infoblox object filterrelayagent
type Filterrelayagent struct {
	IBBase                   `json:"-"`
	Ref                      string `json:"_ref,omitempty"`
	CircuitIdName            string `json:"circuit_id_name,omitempty"`
	CircuitIdSubstringLength uint32 `json:"circuit_id_substring_length,omitempty"`
	CircuitIdSubstringOffset uint32 `json:"circuit_id_substring_offset,omitempty"`
	Comment                  string `json:"comment,omitempty"`
	Ea                       EA     `json:"extattrs,omitempty"`
	IsCircuitId              string `json:"is_circuit_id,omitempty"`
	IsCircuitIdSubstring     bool   `json:"is_circuit_id_substring,omitempty"`
	IsRemoteId               string `json:"is_remote_id,omitempty"`
	IsRemoteIdSubstring      bool   `json:"is_remote_id_substring,omitempty"`
	Name                     string `json:"name,omitempty"`
	RemoteIdName             string `json:"remote_id_name,omitempty"`
	RemoteIdSubstringLength  uint32 `json:"remote_id_substring_length,omitempty"`
	RemoteIdSubstringOffset  uint32 `json:"remote_id_substring_offset,omitempty"`
}

func (Filterrelayagent) ObjectType() string {
	return "filterrelayagent"
}

func (obj Filterrelayagent) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name"}
	}
	return obj.returnFields
}

// DxlEndpoint represents Infoblox object dxl:endpoint
type DxlEndpoint struct {
	IBBase                     `json:"-"`
	Ref                        string                            `json:"_ref,omitempty"`
	Brokers                    []*DxlEndpointBroker              `json:"brokers,omitempty"`
	BrokersImportToken         string                            `json:"brokers_import_token,omitempty"`
	ClearOutboundWorkerLog     *Clearworkerlog                   `json:"clear_outbound_worker_log,omitempty"`
	ClientCertificateSubject   string                            `json:"client_certificate_subject,omitempty"`
	ClientCertificateToken     string                            `json:"client_certificate_token,omitempty"`
	ClientCertificateValidFrom *time.Time                        `json:"client_certificate_valid_from,omitempty"`
	ClientCertificateValidTo   *time.Time                        `json:"client_certificate_valid_to,omitempty"`
	Comment                    string                            `json:"comment,omitempty"`
	Disable                    bool                              `json:"disable,omitempty"`
	Ea                         EA                                `json:"extattrs,omitempty"`
	LogLevel                   string                            `json:"log_level,omitempty"`
	Name                       string                            `json:"name,omitempty"`
	OutboundMemberType         string                            `json:"outbound_member_type,omitempty"`
	OutboundMembers            []string                          `json:"outbound_members,omitempty"`
	TemplateInstance           *NotificationRestTemplateinstance `json:"template_instance,omitempty"`
	TestBrokerConnectivity     *Testdxlbrokerconnectivity        `json:"test_broker_connectivity,omitempty"`
	Timeout                    uint32                            `json:"timeout,omitempty"`
	Topics                     []string                          `json:"topics,omitempty"`
	VendorIdentifier           string                            `json:"vendor_identifier,omitempty"`
	WapiUserName               string                            `json:"wapi_user_name,omitempty"`
	WapiUserPassword           string                            `json:"wapi_user_password,omitempty"`
}

func (DxlEndpoint) ObjectType() string {
	return "dxl:endpoint"
}

func (obj DxlEndpoint) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"disable", "name", "outbound_member_type"}
	}
	return obj.returnFields
}

// DtcTopologyRule represents Infoblox object dtc:topology:rule
type DtcTopologyRule struct {
	IBBase          `json:"-"`
	Ref             string                   `json:"_ref,omitempty"`
	DestType        string                   `json:"dest_type,omitempty"`
	DestinationLink string                   `json:"destination_link,omitempty"`
	ReturnType      string                   `json:"return_type,omitempty"`
	Sources         []*DtcTopologyRuleSource `json:"sources,omitempty"`
	Topology        string                   `json:"topology,omitempty"`
	Valid           bool                     `json:"valid,omitempty"`
}

func (DtcTopologyRule) ObjectType() string {
	return "dtc:topology:rule"
}

func (obj DtcTopologyRule) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{}
	}
	return obj.returnFields
}

// Fingerprint represents Infoblox object fingerprint
type Fingerprint struct {
	IBBase             `json:"-"`
	Ref                string   `json:"_ref,omitempty"`
	Comment            string   `json:"comment,omitempty"`
	DeviceClass        string   `json:"device_class,omitempty"`
	Disable            bool     `json:"disable,omitempty"`
	Ea                 EA       `json:"extattrs,omitempty"`
	Ipv6OptionSequence []string `json:"ipv6_option_sequence,omitempty"`
	Name               string   `json:"name,omitempty"`
	OptionSequence     []string `json:"option_sequence,omitempty"`
	Type               string   `json:"type,omitempty"`
	VendorId           []string `json:"vendor_id,omitempty"`
}

func (Fingerprint) ObjectType() string {
	return "fingerprint"
}

func (obj Fingerprint) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "device_class", "name"}
	}
	return obj.returnFields
}

// Fileop represents Infoblox object fileop
type Fileop struct {
	IBBase                      `json:"-"`
	Ref                         string                             `json:"_ref,omitempty"`
	CsvErrorLog                 *Csverrorlog                       `json:"csv_error_log,omitempty"`
	CsvExport                   *Csvexport                         `json:"csv_export,omitempty"`
	CsvImport                   *Csvimport                         `json:"csv_import,omitempty"`
	CsvSnapshotFile             *Csvsnapshot                       `json:"csv_snapshot_file,omitempty"`
	CsvUploadedFile             *Csvuploaded                       `json:"csv_uploaded_file,omitempty"`
	DownloadAtpRuleUpdate       *Emptyparams                       `json:"download_atp_rule_update,omitempty"`
	DownloadPoolStatus          *Downloadpoolstatus                `json:"download_pool_status,omitempty"`
	Downloadcertificate         *Downloadcertificate               `json:"downloadcertificate,omitempty"`
	Downloadcomplete            *Datagetcomplete                   `json:"downloadcomplete,omitempty"`
	Generatecsr                 *Generatecsr                       `json:"generatecsr,omitempty"`
	Generatedxlendpointcerts    *Generatedxlendpointcerts          `json:"generatedxlendpointcerts,omitempty"`
	Generatesafenetclientcert   *Generatesafenetclientcert         `json:"generatesafenetclientcert,omitempty"`
	Generateselfsignedcert      *Generateselfsignedcert            `json:"generateselfsignedcert,omitempty"`
	GetFileUrl                  *Getfileurl                        `json:"get_file_url,omitempty"`
	GetLastUploadedAtpRuleset   *Getlastuploadedruleset            `json:"get_last_uploaded_atp_ruleset,omitempty"`
	GetLogFiles                 *Getlogfiles                       `json:"get_log_files,omitempty"`
	GetSupportBundle            *Supportbundle                     `json:"get_support_bundle,omitempty"`
	Getgriddata                 *Getgriddata                       `json:"getgriddata,omitempty"`
	Getleasehistoryfiles        *Getleasehistoryfiles              `json:"getleasehistoryfiles,omitempty"`
	Getmemberdata               *Getmemberdata                     `json:"getmemberdata,omitempty"`
	Getsafenetclientcert        *Getsafenetclientcert              `json:"getsafenetclientcert,omitempty"`
	Read                        *Read                              `json:"read,omitempty"`
	RestapiTemplateExport       *Restapitemplateexportparams       `json:"restapi_template_export,omitempty"`
	RestapiTemplateExportSchema *Restapitemplateexportschemaparams `json:"restapi_template_export_schema,omitempty"`
	RestapiTemplateImport       *Restapitemplateimportparams       `json:"restapi_template_import,omitempty"`
	Restoredatabase             *Restoredatabase                   `json:"restoredatabase,omitempty"`
	Restoredtcconfig            *Restoredtcconfig                  `json:"restoredtcconfig,omitempty"`
	SetCaptivePortalFile        *Setcaptiveportalfile              `json:"set_captive_portal_file,omitempty"`
	SetDhcpLeases               *Setdhcpleases                     `json:"set_dhcp_leases,omitempty"`
	SetDowngradeFile            *Downgrade                         `json:"set_downgrade_file,omitempty"`
	SetLastUploadedAtpRuleset   *Setlastuploadedruleset            `json:"set_last_uploaded_atp_ruleset,omitempty"`
	SetUpgradeFile              *Upgrade                           `json:"set_upgrade_file,omitempty"`
	Setdiscoverycsv             *Setdiscoverycsv                   `json:"setdiscoverycsv,omitempty"`
	Setfiledest                 *Setdatafiledest                   `json:"setfiledest,omitempty"`
	Setleasehistoryfiles        *Setleasehistoryfiles              `json:"setleasehistoryfiles,omitempty"`
	Setmemberdata               *Setmemberdata                     `json:"setmemberdata,omitempty"`
	UpdateAtpRuleset            *Updateatprulesetparams            `json:"update_atp_ruleset,omitempty"`
	UpdateLicenses              *Updatelicenses                    `json:"update_licenses,omitempty"`
	Uploadcertificate           *Uploadcertificate                 `json:"uploadcertificate,omitempty"`
	Uploadinit                  *Datauploadinit                    `json:"uploadinit,omitempty"`
	Uploadserviceaccount        *Uploadserviceaccount              `json:"uploadserviceaccount,omitempty"`
}

func (Fileop) ObjectType() string {
	return "fileop"
}

func (obj Fileop) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{}
	}
	return obj.returnFields
}

// Fixedaddresstemplate represents Infoblox object fixedaddresstemplate
type Fixedaddresstemplate struct {
	IBBase                         `json:"-"`
	Ref                            string             `json:"_ref,omitempty"`
	Bootfile                       string             `json:"bootfile,omitempty"`
	Bootserver                     string             `json:"bootserver,omitempty"`
	Comment                        string             `json:"comment,omitempty"`
	DdnsDomainname                 string             `json:"ddns_domainname,omitempty"`
	DdnsHostname                   string             `json:"ddns_hostname,omitempty"`
	DenyBootp                      bool               `json:"deny_bootp,omitempty"`
	EnableDdns                     bool               `json:"enable_ddns,omitempty"`
	EnablePxeLeaseTime             bool               `json:"enable_pxe_lease_time,omitempty"`
	Ea                             EA                 `json:"extattrs,omitempty"`
	IgnoreDhcpOptionListRequest    bool               `json:"ignore_dhcp_option_list_request,omitempty"`
	LogicFilterRules               []*Logicfilterrule `json:"logic_filter_rules,omitempty"`
	Name                           string             `json:"name,omitempty"`
	Nextserver                     string             `json:"nextserver,omitempty"`
	NumberOfAddresses              uint32             `json:"number_of_addresses,omitempty"`
	Offset                         uint32             `json:"offset,omitempty"`
	Options                        []*Dhcpoption      `json:"options,omitempty"`
	PxeLeaseTime                   uint32             `json:"pxe_lease_time,omitempty"`
	UseBootfile                    bool               `json:"use_bootfile,omitempty"`
	UseBootserver                  bool               `json:"use_bootserver,omitempty"`
	UseDdnsDomainname              bool               `json:"use_ddns_domainname,omitempty"`
	UseDenyBootp                   bool               `json:"use_deny_bootp,omitempty"`
	UseEnableDdns                  bool               `json:"use_enable_ddns,omitempty"`
	UseIgnoreDhcpOptionListRequest bool               `json:"use_ignore_dhcp_option_list_request,omitempty"`
	UseLogicFilterRules            bool               `json:"use_logic_filter_rules,omitempty"`
	UseNextserver                  bool               `json:"use_nextserver,omitempty"`
	UseOptions                     bool               `json:"use_options,omitempty"`
	UsePxeLeaseTime                bool               `json:"use_pxe_lease_time,omitempty"`
}

func (Fixedaddresstemplate) ObjectType() string {
	return "fixedaddresstemplate"
}

func (obj Fixedaddresstemplate) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name"}
	}
	return obj.returnFields
}

// Ftpuser represents Infoblox object ftpuser
type Ftpuser struct {
	IBBase        `json:"-"`
	Ref           string `json:"_ref,omitempty"`
	CreateHomeDir bool   `json:"create_home_dir,omitempty"`
	Ea            EA     `json:"extattrs,omitempty"`
	HomeDir       string `json:"home_dir,omitempty"`
	Password      string `json:"password,omitempty"`
	Permission    string `json:"permission,omitempty"`
	Username      string `json:"username,omitempty"`
}

func (Ftpuser) ObjectType() string {
	return "ftpuser"
}

func (obj Ftpuser) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"username"}
	}
	return obj.returnFields
}

// GridCloudapiCloudstatistics represents Infoblox object grid:cloudapi:cloudstatistics
type GridCloudapiCloudstatistics struct {
	IBBase                  `json:"-"`
	Ref                     string `json:"_ref,omitempty"`
	AllocatedAvailableRatio uint32 `json:"allocated_available_ratio,omitempty"`
	AllocatedIpCount        uint32 `json:"allocated_ip_count,omitempty"`
	AvailableIpCount        string `json:"available_ip_count,omitempty"`
	FixedIpCount            uint32 `json:"fixed_ip_count,omitempty"`
	FloatingIpCount         uint32 `json:"floating_ip_count,omitempty"`
	TenantCount             uint32 `json:"tenant_count,omitempty"`
	TenantIpCount           uint32 `json:"tenant_ip_count,omitempty"`
	TenantVmCount           uint32 `json:"tenant_vm_count,omitempty"`
}

func (GridCloudapiCloudstatistics) ObjectType() string {
	return "grid:cloudapi:cloudstatistics"
}

func (obj GridCloudapiCloudstatistics) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"allocated_available_ratio", "allocated_ip_count", "available_ip_count", "fixed_ip_count", "floating_ip_count", "tenant_count", "tenant_ip_count", "tenant_vm_count"}
	}
	return obj.returnFields
}

// GridCloudapi represents Infoblox object grid:cloudapi
type GridCloudapi struct {
	IBBase           `json:"-"`
	Ref              string                     `json:"_ref,omitempty"`
	AllowApiAdmins   string                     `json:"allow_api_admins,omitempty"`
	AllowedApiAdmins []*GridCloudapiUser        `json:"allowed_api_admins,omitempty"`
	EnableRecycleBin bool                       `json:"enable_recycle_bin,omitempty"`
	GatewayConfig    *GridCloudapiGatewayConfig `json:"gateway_config,omitempty"`
}

func (GridCloudapi) ObjectType() string {
	return "grid:cloudapi"
}

func (obj GridCloudapi) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"allow_api_admins", "allowed_api_admins", "enable_recycle_bin"}
	}
	return obj.returnFields
}

// GridCloudapiTenant represents Infoblox object grid:cloudapi:tenant
type GridCloudapiTenant struct {
	IBBase       `json:"-"`
	Ref          string            `json:"_ref,omitempty"`
	CloudInfo    *GridCloudapiInfo `json:"cloud_info,omitempty"`
	Comment      string            `json:"comment,omitempty"`
	CreatedTs    *time.Time        `json:"created_ts,omitempty"`
	Id           string            `json:"id,omitempty"`
	LastEventTs  *time.Time        `json:"last_event_ts,omitempty"`
	Name         string            `json:"name,omitempty"`
	NetworkCount uint32            `json:"network_count,omitempty"`
	VmCount      uint32            `json:"vm_count,omitempty"`
}

func (GridCloudapiTenant) ObjectType() string {
	return "grid:cloudapi:tenant"
}

func (obj GridCloudapiTenant) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "id", "name"}
	}
	return obj.returnFields
}

// GridCloudapiVm represents Infoblox object grid:cloudapi:vm
type GridCloudapiVm struct {
	IBBase            `json:"-"`
	Ref               string            `json:"_ref,omitempty"`
	AvailabilityZone  string            `json:"availability_zone,omitempty"`
	CloudInfo         *GridCloudapiInfo `json:"cloud_info,omitempty"`
	Comment           string            `json:"comment,omitempty"`
	ElasticIpAddress  string            `json:"elastic_ip_address,omitempty"`
	Ea                EA                `json:"extattrs,omitempty"`
	FirstSeen         *time.Time        `json:"first_seen,omitempty"`
	Hostname          string            `json:"hostname,omitempty"`
	Id                string            `json:"id,omitempty"`
	KernelId          string            `json:"kernel_id,omitempty"`
	LastSeen          *time.Time        `json:"last_seen,omitempty"`
	Name              string            `json:"name,omitempty"`
	NetworkCount      uint32            `json:"network_count,omitempty"`
	OperatingSystem   string            `json:"operating_system,omitempty"`
	PrimaryMacAddress string            `json:"primary_mac_address,omitempty"`
	SubnetAddress     string            `json:"subnet_address,omitempty"`
	SubnetCidr        uint32            `json:"subnet_cidr,omitempty"`
	SubnetId          string            `json:"subnet_id,omitempty"`
	TenantName        string            `json:"tenant_name,omitempty"`
	VmType            string            `json:"vm_type,omitempty"`
	VpcAddress        string            `json:"vpc_address,omitempty"`
	VpcCidr           uint32            `json:"vpc_cidr,omitempty"`
	VpcId             string            `json:"vpc_id,omitempty"`
	VpcName           string            `json:"vpc_name,omitempty"`
}

func (GridCloudapiVm) ObjectType() string {
	return "grid:cloudapi:vm"
}

func (obj GridCloudapiVm) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "id", "name"}
	}
	return obj.returnFields
}

// GridDashboard represents Infoblox object grid:dashboard
type GridDashboard struct {
	IBBase                                   `json:"-"`
	Ref                                      string `json:"_ref,omitempty"`
	AnalyticsTunnelingEventCriticalThreshold uint32 `json:"analytics_tunneling_event_critical_threshold,omitempty"`
	AnalyticsTunnelingEventWarningThreshold  uint32 `json:"analytics_tunneling_event_warning_threshold,omitempty"`
	AtpCriticalEventCriticalThreshold        uint32 `json:"atp_critical_event_critical_threshold,omitempty"`
	AtpCriticalEventWarningThreshold         uint32 `json:"atp_critical_event_warning_threshold,omitempty"`
	AtpMajorEventCriticalThreshold           uint32 `json:"atp_major_event_critical_threshold,omitempty"`
	AtpMajorEventWarningThreshold            uint32 `json:"atp_major_event_warning_threshold,omitempty"`
	AtpWarningEventCriticalThreshold         uint32 `json:"atp_warning_event_critical_threshold,omitempty"`
	AtpWarningEventWarningThreshold          uint32 `json:"atp_warning_event_warning_threshold,omitempty"`
	RpzBlockedHitCriticalThreshold           uint32 `json:"rpz_blocked_hit_critical_threshold,omitempty"`
	RpzBlockedHitWarningThreshold            uint32 `json:"rpz_blocked_hit_warning_threshold,omitempty"`
	RpzPassthruEventCriticalThreshold        uint32 `json:"rpz_passthru_event_critical_threshold,omitempty"`
	RpzPassthruEventWarningThreshold         uint32 `json:"rpz_passthru_event_warning_threshold,omitempty"`
	RpzSubstitutedHitCriticalThreshold       uint32 `json:"rpz_substituted_hit_critical_threshold,omitempty"`
	RpzSubstitutedHitWarningThreshold        uint32 `json:"rpz_substituted_hit_warning_threshold,omitempty"`
}

func (GridDashboard) ObjectType() string {
	return "grid:dashboard"
}

func (obj GridDashboard) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"analytics_tunneling_event_critical_threshold", "analytics_tunneling_event_warning_threshold", "atp_critical_event_critical_threshold", "atp_critical_event_warning_threshold", "atp_major_event_critical_threshold", "atp_major_event_warning_threshold", "atp_warning_event_critical_threshold", "atp_warning_event_warning_threshold", "rpz_blocked_hit_critical_threshold", "rpz_blocked_hit_warning_threshold", "rpz_passthru_event_critical_threshold", "rpz_passthru_event_warning_threshold", "rpz_substituted_hit_critical_threshold", "rpz_substituted_hit_warning_threshold"}
	}
	return obj.returnFields
}

// GridCloudapiVmaddress represents Infoblox object grid:cloudapi:vmaddress
type GridCloudapiVmaddress struct {
	IBBase                `json:"-"`
	Ref                   string              `json:"_ref,omitempty"`
	Address               string              `json:"address,omitempty"`
	AddressType           string              `json:"address_type,omitempty"`
	AssociatedIp          string              `json:"associated_ip,omitempty"`
	AssociatedObjectTypes []string            `json:"associated_object_types,omitempty"`
	AssociatedObjects     []*Ipv4FixedAddress `json:"associated_objects,omitempty"`
	CloudInfo             *GridCloudapiInfo   `json:"cloud_info,omitempty"`
	DnsNames              []string            `json:"dns_names,omitempty"`
	ElasticAddress        string              `json:"elastic_address,omitempty"`
	InterfaceName         string              `json:"interface_name,omitempty"`
	IsIpv4                bool                `json:"is_ipv4,omitempty"`
	MacAddress            string              `json:"mac_address,omitempty"`
	MsAdUserData          *MsserverAduserData `json:"ms_ad_user_data,omitempty"`
	Network               string              `json:"network,omitempty"`
	NetworkView           string              `json:"network_view,omitempty"`
	PortId                uint32              `json:"port_id,omitempty"`
	PrivateAddress        string              `json:"private_address,omitempty"`
	PrivateHostname       string              `json:"private_hostname,omitempty"`
	PublicAddress         string              `json:"public_address,omitempty"`
	PublicHostname        string              `json:"public_hostname,omitempty"`
	SubnetAddress         string              `json:"subnet_address,omitempty"`
	SubnetCidr            uint32              `json:"subnet_cidr,omitempty"`
	SubnetId              string              `json:"subnet_id,omitempty"`
	Tenant                string              `json:"tenant,omitempty"`
	VmAvailabilityZone    string              `json:"vm_availability_zone,omitempty"`
	VmComment             string              `json:"vm_comment,omitempty"`
	VmCreationTime        *time.Time          `json:"vm_creation_time,omitempty"`
	VmHostname            string              `json:"vm_hostname,omitempty"`
	VmId                  string              `json:"vm_id,omitempty"`
	VmKernelId            string              `json:"vm_kernel_id,omitempty"`
	VmLastUpdateTime      *time.Time          `json:"vm_last_update_time,omitempty"`
	VmName                string              `json:"vm_name,omitempty"`
	VmNetworkCount        uint32              `json:"vm_network_count,omitempty"`
	VmOperatingSystem     string              `json:"vm_operating_system,omitempty"`
	VmType                string              `json:"vm_type,omitempty"`
	VmVpcAddress          string              `json:"vm_vpc_address,omitempty"`
	VmVpcCidr             uint32              `json:"vm_vpc_cidr,omitempty"`
	VmVpcId               string              `json:"vm_vpc_id,omitempty"`
	VmVpcName             string              `json:"vm_vpc_name,omitempty"`
	VmVpcRef              string              `json:"vm_vpc_ref,omitempty"`
}

func (GridCloudapiVmaddress) ObjectType() string {
	return "grid:cloudapi:vmaddress"
}

func (obj GridCloudapiVmaddress) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"address", "is_ipv4", "network_view", "port_id", "vm_name"}
	}
	return obj.returnFields
}

// Ipv4FixedAddress represents Infoblox object fixedaddress
type Ipv4FixedAddress struct {
	IBBase                         `json:"-"`
	Ref                            string                    `json:"_ref,omitempty"`
	AgentCircuitId                 string                    `json:"agent_circuit_id,omitempty"`
	AgentRemoteId                  string                    `json:"agent_remote_id,omitempty"`
	AllowTelnet                    bool                      `json:"allow_telnet,omitempty"`
	AlwaysUpdateDns                bool                      `json:"always_update_dns,omitempty"`
	Bootfile                       string                    `json:"bootfile,omitempty"`
	Bootserver                     string                    `json:"bootserver,omitempty"`
	CliCredentials                 []*DiscoveryClicredential `json:"cli_credentials,omitempty"`
	ClientIdentifierPrependZero    bool                      `json:"client_identifier_prepend_zero,omitempty"`
	CloudInfo                      *GridCloudapiInfo         `json:"cloud_info,omitempty"`
	Comment                        string                    `json:"comment,omitempty"`
	DdnsDomainname                 string                    `json:"ddns_domainname,omitempty"`
	DdnsHostname                   string                    `json:"ddns_hostname,omitempty"`
	DenyBootp                      bool                      `json:"deny_bootp,omitempty"`
	DeviceDescription              string                    `json:"device_description,omitempty"`
	DeviceLocation                 string                    `json:"device_location,omitempty"`
	DeviceType                     string                    `json:"device_type,omitempty"`
	DeviceVendor                   string                    `json:"device_vendor,omitempty"`
	DhcpClientIdentifier           string                    `json:"dhcp_client_identifier,omitempty"`
	Disable                        bool                      `json:"disable,omitempty"`
	DisableDiscovery               bool                      `json:"disable_discovery,omitempty"`
	DiscoverNowStatus              string                    `json:"discover_now_status,omitempty"`
	DiscoveredData                 *Discoverydata            `json:"discovered_data,omitempty"`
	EnableDdns                     bool                      `json:"enable_ddns,omitempty"`
	EnableImmediateDiscovery       bool                      `json:"enable_immediate_discovery,omitempty"`
	EnablePxeLeaseTime             bool                      `json:"enable_pxe_lease_time,omitempty"`
	Ea                             EA                        `json:"extattrs,omitempty"`
	IgnoreDhcpOptionListRequest    bool                      `json:"ignore_dhcp_option_list_request,omitempty"`
	Ipv4Addr                       string                    `json:"ipv4addr,omitempty"`
	IsInvalidMac                   bool                      `json:"is_invalid_mac,omitempty"`
	LogicFilterRules               []*Logicfilterrule        `json:"logic_filter_rules,omitempty"`
	Mac                            string                    `json:"mac,omitempty"`
	MatchClient                    string                    `json:"match_client,omitempty"`
	MsAdUserData                   *MsserverAduserData       `json:"ms_ad_user_data,omitempty"`
	MsOptions                      []*Msdhcpoption           `json:"ms_options,omitempty"`
	MsServer                       *Msdhcpserver             `json:"ms_server,omitempty"`
	Name                           string                    `json:"name,omitempty"`
	Network                        string                    `json:"network,omitempty"`
	NetworkView                    string                    `json:"network_view,omitempty"`
	Nextserver                     string                    `json:"nextserver,omitempty"`
	Options                        []*Dhcpoption             `json:"options,omitempty"`
	PxeLeaseTime                   uint32                    `json:"pxe_lease_time,omitempty"`
	ReservedInterface              string                    `json:"reserved_interface,omitempty"`
	RestartIfNeeded                bool                      `json:"restart_if_needed,omitempty"`
	Snmp3Credential                *DiscoverySnmp3credential `json:"snmp3_credential,omitempty"`
	SnmpCredential                 *DiscoverySnmpcredential  `json:"snmp_credential,omitempty"`
	Template                       string                    `json:"template,omitempty"`
	UseBootfile                    bool                      `json:"use_bootfile,omitempty"`
	UseBootserver                  bool                      `json:"use_bootserver,omitempty"`
	UseCliCredentials              bool                      `json:"use_cli_credentials,omitempty"`
	UseDdnsDomainname              bool                      `json:"use_ddns_domainname,omitempty"`
	UseDenyBootp                   bool                      `json:"use_deny_bootp,omitempty"`
	UseEnableDdns                  bool                      `json:"use_enable_ddns,omitempty"`
	UseIgnoreDhcpOptionListRequest bool                      `json:"use_ignore_dhcp_option_list_request,omitempty"`
	UseLogicFilterRules            bool                      `json:"use_logic_filter_rules,omitempty"`
	UseMsOptions                   bool                      `json:"use_ms_options,omitempty"`
	UseNextserver                  bool                      `json:"use_nextserver,omitempty"`
	UseOptions                     bool                      `json:"use_options,omitempty"`
	UsePxeLeaseTime                bool                      `json:"use_pxe_lease_time,omitempty"`
	UseSnmp3Credential             bool                      `json:"use_snmp3_credential,omitempty"`
	UseSnmpCredential              bool                      `json:"use_snmp_credential,omitempty"`
}

func (Ipv4FixedAddress) ObjectType() string {
	return "fixedaddress"
}

func (obj Ipv4FixedAddress) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"ipv4addr", "network_view"}
	}
	return obj.returnFields
}

// GridFiledistribution represents Infoblox object grid:filedistribution
type GridFiledistribution struct {
	IBBase             `json:"-"`
	Ref                string `json:"_ref,omitempty"`
	AllowUploads       bool   `json:"allow_uploads,omitempty"`
	BackupStorage      bool   `json:"backup_storage,omitempty"`
	CurrentUsage       uint32 `json:"current_usage,omitempty"`
	EnableAnonymousFtp bool   `json:"enable_anonymous_ftp,omitempty"`
	GlobalStatus       string `json:"global_status,omitempty"`
	Name               string `json:"name,omitempty"`
	StorageLimit       uint32 `json:"storage_limit,omitempty"`
}

func (GridFiledistribution) ObjectType() string {
	return "grid:filedistribution"
}

func (obj GridFiledistribution) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"allow_uploads", "current_usage", "global_status", "name", "storage_limit"}
	}
	return obj.returnFields
}

// GridLicensePool represents Infoblox object grid:license_pool
type GridLicensePool struct {
	IBBase           `json:"-"`
	Ref              string                `json:"_ref,omitempty"`
	Assigned         uint32                `json:"assigned,omitempty"`
	ExpirationStatus string                `json:"expiration_status,omitempty"`
	ExpiryDate       *time.Time            `json:"expiry_date,omitempty"`
	Installed        uint32                `json:"installed,omitempty"`
	Key              string                `json:"key,omitempty"`
	Limit            string                `json:"limit,omitempty"`
	LimitContext     string                `json:"limit_context,omitempty"`
	Model            string                `json:"model,omitempty"`
	Subpools         []*GridLicensesubpool `json:"subpools,omitempty"`
	TempAssigned     uint32                `json:"temp_assigned,omitempty"`
	Type             string                `json:"type,omitempty"`
}

func (GridLicensePool) ObjectType() string {
	return "grid:license_pool"
}

func (obj GridLicensePool) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"type"}
	}
	return obj.returnFields
}

// GridLicensePoolContainer represents Infoblox object grid:license_pool_container
type GridLicensePoolContainer struct {
	IBBase                `json:"-"`
	Ref                   string            `json:"_ref,omitempty"`
	AllocateLicenses      *Allocatelicenses `json:"allocate_licenses,omitempty"`
	LastEntitlementUpdate *time.Time        `json:"last_entitlement_update,omitempty"`
	LpcUid                string            `json:"lpc_uid,omitempty"`
}

func (GridLicensePoolContainer) ObjectType() string {
	return "grid:license_pool_container"
}

func (obj GridLicensePoolContainer) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{}
	}
	return obj.returnFields
}

// GridMaxminddbinfo represents Infoblox object grid:maxminddbinfo
type GridMaxminddbinfo struct {
	IBBase             `json:"-"`
	Ref                string     `json:"_ref,omitempty"`
	BinaryMajorVersion uint32     `json:"binary_major_version,omitempty"`
	BinaryMinorVersion uint32     `json:"binary_minor_version,omitempty"`
	BuildTime          *time.Time `json:"build_time,omitempty"`
	DatabaseType       string     `json:"database_type,omitempty"`
	DeploymentTime     *time.Time `json:"deployment_time,omitempty"`
	Member             string     `json:"member,omitempty"`
	TopologyType       string     `json:"topology_type,omitempty"`
}

func (GridMaxminddbinfo) ObjectType() string {
	return "grid:maxminddbinfo"
}

func (obj GridMaxminddbinfo) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"binary_major_version", "binary_minor_version", "build_time", "database_type", "deployment_time", "member", "topology_type"}
	}
	return obj.returnFields
}

// GridDhcpproperties represents Infoblox object grid:dhcpproperties
type GridDhcpproperties struct {
	IBBase                               `json:"-"`
	Ref                                  string                   `json:"_ref,omitempty"`
	Authority                            bool                     `json:"authority,omitempty"`
	Bootfile                             string                   `json:"bootfile,omitempty"`
	Bootserver                           string                   `json:"bootserver,omitempty"`
	CaptureHostname                      bool                     `json:"capture_hostname,omitempty"`
	DdnsDomainname                       string                   `json:"ddns_domainname,omitempty"`
	DdnsGenerateHostname                 bool                     `json:"ddns_generate_hostname,omitempty"`
	DdnsRetryInterval                    uint32                   `json:"ddns_retry_interval,omitempty"`
	DdnsServerAlwaysUpdates              bool                     `json:"ddns_server_always_updates,omitempty"`
	DdnsTtl                              uint32                   `json:"ddns_ttl,omitempty"`
	DdnsUpdateFixedAddresses             bool                     `json:"ddns_update_fixed_addresses,omitempty"`
	DdnsUseOption81                      bool                     `json:"ddns_use_option81,omitempty"`
	DenyBootp                            bool                     `json:"deny_bootp,omitempty"`
	DisableAllNacFilters                 bool                     `json:"disable_all_nac_filters,omitempty"`
	DnsUpdateStyle                       string                   `json:"dns_update_style,omitempty"`
	EmailList                            []string                 `json:"email_list,omitempty"`
	EnableDdns                           bool                     `json:"enable_ddns,omitempty"`
	EnableDhcpThresholds                 bool                     `json:"enable_dhcp_thresholds,omitempty"`
	EnableEmailWarnings                  bool                     `json:"enable_email_warnings,omitempty"`
	EnableFingerprint                    bool                     `json:"enable_fingerprint,omitempty"`
	EnableGssTsig                        bool                     `json:"enable_gss_tsig,omitempty"`
	EnableHostnameRewrite                bool                     `json:"enable_hostname_rewrite,omitempty"`
	EnableLeasequery                     bool                     `json:"enable_leasequery,omitempty"`
	EnableRoamingHosts                   bool                     `json:"enable_roaming_hosts,omitempty"`
	EnableSnmpWarnings                   bool                     `json:"enable_snmp_warnings,omitempty"`
	FormatLogOption82                    string                   `json:"format_log_option_82,omitempty"`
	Grid                                 string                   `json:"grid,omitempty"`
	GssTsigKeys                          []*Kerberoskey           `json:"gss_tsig_keys,omitempty"`
	HighWaterMark                        uint32                   `json:"high_water_mark,omitempty"`
	HighWaterMarkReset                   uint32                   `json:"high_water_mark_reset,omitempty"`
	HostnameRewritePolicy                string                   `json:"hostname_rewrite_policy,omitempty"`
	IgnoreDhcpOptionListRequest          bool                     `json:"ignore_dhcp_option_list_request,omitempty"`
	IgnoreId                             string                   `json:"ignore_id,omitempty"`
	IgnoreMacAddresses                   []string                 `json:"ignore_mac_addresses,omitempty"`
	ImmediateFaConfiguration             bool                     `json:"immediate_fa_configuration,omitempty"`
	Ipv6CaptureHostname                  bool                     `json:"ipv6_capture_hostname,omitempty"`
	Ipv6DdnsDomainname                   string                   `json:"ipv6_ddns_domainname,omitempty"`
	Ipv6DdnsEnableOptionFqdn             bool                     `json:"ipv6_ddns_enable_option_fqdn,omitempty"`
	Ipv6DdnsServerAlwaysUpdates          bool                     `json:"ipv6_ddns_server_always_updates,omitempty"`
	Ipv6DdnsTtl                          uint32                   `json:"ipv6_ddns_ttl,omitempty"`
	Ipv6DefaultPrefix                    string                   `json:"ipv6_default_prefix,omitempty"`
	Ipv6DnsUpdateStyle                   string                   `json:"ipv6_dns_update_style,omitempty"`
	Ipv6DomainName                       string                   `json:"ipv6_domain_name,omitempty"`
	Ipv6DomainNameServers                []string                 `json:"ipv6_domain_name_servers,omitempty"`
	Ipv6EnableDdns                       bool                     `json:"ipv6_enable_ddns,omitempty"`
	Ipv6EnableGssTsig                    bool                     `json:"ipv6_enable_gss_tsig,omitempty"`
	Ipv6EnableLeaseScavenging            bool                     `json:"ipv6_enable_lease_scavenging,omitempty"`
	Ipv6EnableRetryUpdates               bool                     `json:"ipv6_enable_retry_updates,omitempty"`
	Ipv6GenerateHostname                 bool                     `json:"ipv6_generate_hostname,omitempty"`
	Ipv6GssTsigKeys                      []*Kerberoskey           `json:"ipv6_gss_tsig_keys,omitempty"`
	Ipv6KdcServer                        string                   `json:"ipv6_kdc_server,omitempty"`
	Ipv6LeaseScavengingTime              uint32                   `json:"ipv6_lease_scavenging_time,omitempty"`
	Ipv6MicrosoftCodePage                string                   `json:"ipv6_microsoft_code_page,omitempty"`
	Ipv6Options                          []*Dhcpoption            `json:"ipv6_options,omitempty"`
	Ipv6Prefixes                         []string                 `json:"ipv6_prefixes,omitempty"`
	Ipv6RecycleLeases                    bool                     `json:"ipv6_recycle_leases,omitempty"`
	Ipv6RememberExpiredClientAssociation bool                     `json:"ipv6_remember_expired_client_association,omitempty"`
	Ipv6RetryUpdatesInterval             uint32                   `json:"ipv6_retry_updates_interval,omitempty"`
	Ipv6TxtRecordHandling                string                   `json:"ipv6_txt_record_handling,omitempty"`
	Ipv6UpdateDnsOnLeaseRenewal          bool                     `json:"ipv6_update_dns_on_lease_renewal,omitempty"`
	KdcServer                            string                   `json:"kdc_server,omitempty"`
	LeaseLoggingMember                   string                   `json:"lease_logging_member,omitempty"`
	LeasePerClientSettings               string                   `json:"lease_per_client_settings,omitempty"`
	LeaseScavengeTime                    int                      `json:"lease_scavenge_time,omitempty"`
	LogLeaseEvents                       bool                     `json:"log_lease_events,omitempty"`
	LogicFilterRules                     []*Logicfilterrule       `json:"logic_filter_rules,omitempty"`
	LowWaterMark                         uint32                   `json:"low_water_mark,omitempty"`
	LowWaterMarkReset                    uint32                   `json:"low_water_mark_reset,omitempty"`
	MicrosoftCodePage                    string                   `json:"microsoft_code_page,omitempty"`
	Nextserver                           string                   `json:"nextserver,omitempty"`
	Option60MatchRules                   []*Option60matchrule     `json:"option60_match_rules,omitempty"`
	Options                              []*Dhcpoption            `json:"options,omitempty"`
	PingCount                            uint32                   `json:"ping_count,omitempty"`
	PingTimeout                          uint32                   `json:"ping_timeout,omitempty"`
	PreferredLifetime                    uint32                   `json:"preferred_lifetime,omitempty"`
	PrefixLengthMode                     string                   `json:"prefix_length_mode,omitempty"`
	ProtocolHostnameRewritePolicies      []*Hostnamerewritepolicy `json:"protocol_hostname_rewrite_policies,omitempty"`
	PxeLeaseTime                         uint32                   `json:"pxe_lease_time,omitempty"`
	RecycleLeases                        bool                     `json:"recycle_leases,omitempty"`
	RestartSetting                       *GridServicerestart      `json:"restart_setting,omitempty"`
	RetryDdnsUpdates                     bool                     `json:"retry_ddns_updates,omitempty"`
	SyslogFacility                       string                   `json:"syslog_facility,omitempty"`
	TxtRecordHandling                    string                   `json:"txt_record_handling,omitempty"`
	UpdateDnsOnLeaseRenewal              bool                     `json:"update_dns_on_lease_renewal,omitempty"`
	ValidLifetime                        uint32                   `json:"valid_lifetime,omitempty"`
}

func (GridDhcpproperties) ObjectType() string {
	return "grid:dhcpproperties"
}

func (obj GridDhcpproperties) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"disable_all_nac_filters", "grid"}
	}
	return obj.returnFields
}

// GridServicerestartGroupOrder represents Infoblox object grid:servicerestart:group:order
type GridServicerestartGroupOrder struct {
	IBBase `json:"-"`
	Ref    string   `json:"_ref,omitempty"`
	Groups []string `json:"groups,omitempty"`
}

func (GridServicerestartGroupOrder) ObjectType() string {
	return "grid:servicerestart:group:order"
}

func (obj GridServicerestartGroupOrder) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{}
	}
	return obj.returnFields
}

// GridServicerestartGroup represents Infoblox object grid:servicerestart:group
type GridServicerestartGroup struct {
	IBBase            `json:"-"`
	Ref               string                           `json:"_ref,omitempty"`
	Comment           string                           `json:"comment,omitempty"`
	Ea                EA                               `json:"extattrs,omitempty"`
	IsDefault         bool                             `json:"is_default,omitempty"`
	LastUpdatedTime   *time.Time                       `json:"last_updated_time,omitempty"`
	Members           []string                         `json:"members,omitempty"`
	Mode              string                           `json:"mode,omitempty"`
	Name              string                           `json:"name,omitempty"`
	Position          uint32                           `json:"position,omitempty"`
	RecurringSchedule *GridServicerestartGroupSchedule `json:"recurring_schedule,omitempty"`
	Requests          []*GridServicerestartRequest     `json:"requests,omitempty"`
	Service           string                           `json:"service,omitempty"`
	Status            string                           `json:"status,omitempty"`
}

func (GridServicerestartGroup) ObjectType() string {
	return "grid:servicerestart:group"
}

func (obj GridServicerestartGroup) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name", "service"}
	}
	return obj.returnFields
}

// GridMemberCloudapi represents Infoblox object grid:member:cloudapi
type GridMemberCloudapi struct {
	IBBase           `json:"-"`
	Ref              string                     `json:"_ref,omitempty"`
	AllowApiAdmins   string                     `json:"allow_api_admins,omitempty"`
	AllowedApiAdmins []*GridCloudapiUser        `json:"allowed_api_admins,omitempty"`
	EnableService    bool                       `json:"enable_service,omitempty"`
	Ea               EA                         `json:"extattrs,omitempty"`
	GatewayConfig    *GridCloudapiGatewayConfig `json:"gateway_config,omitempty"`
	Member           *Dhcpmember                `json:"member,omitempty"`
	Status           string                     `json:"status,omitempty"`
}

func (GridMemberCloudapi) ObjectType() string {
	return "grid:member:cloudapi"
}

func (obj GridMemberCloudapi) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"allow_api_admins", "allowed_api_admins", "enable_service", "member", "status"}
	}
	return obj.returnFields
}

// GridServicerestartRequestChangedobject represents Infoblox object grid:servicerestart:request:changedobject
type GridServicerestartRequestChangedobject struct {
	IBBase            `json:"-"`
	Ref               string     `json:"_ref,omitempty"`
	Action            string     `json:"action,omitempty"`
	ChangedProperties []string   `json:"changed_properties,omitempty"`
	ChangedTime       *time.Time `json:"changed_time,omitempty"`
	ObjectName        string     `json:"object_name,omitempty"`
	ObjectTypeField   string     `json:"object_type,omitempty"`
	UserName          string     `json:"user_name,omitempty"`
}

func (GridServicerestartRequestChangedobject) ObjectType() string {
	return "grid:servicerestart:request:changedobject"
}

func (obj GridServicerestartRequestChangedobject) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"action", "changed_properties", "changed_time", "object_name", "object_type", "user_name"}
	}
	return obj.returnFields
}

// GridServicerestartStatus represents Infoblox object grid:servicerestart:status
type GridServicerestartStatus struct {
	IBBase         `json:"-"`
	Ref            string `json:"_ref,omitempty"`
	Failures       uint32 `json:"failures,omitempty"`
	Finished       uint32 `json:"finished,omitempty"`
	Grouped        string `json:"grouped,omitempty"`
	NeededRestart  uint32 `json:"needed_restart,omitempty"`
	NoRestart      uint32 `json:"no_restart,omitempty"`
	Parent         string `json:"parent,omitempty"`
	Pending        uint32 `json:"pending,omitempty"`
	PendingRestart uint32 `json:"pending_restart,omitempty"`
	Processing     uint32 `json:"processing,omitempty"`
	Restarting     uint32 `json:"restarting,omitempty"`
	Success        uint32 `json:"success,omitempty"`
	Timeouts       uint32 `json:"timeouts,omitempty"`
}

func (GridServicerestartStatus) ObjectType() string {
	return "grid:servicerestart:status"
}

func (obj GridServicerestartStatus) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"failures", "finished", "grouped", "needed_restart", "no_restart", "parent", "pending", "pending_restart", "processing", "restarting", "success", "timeouts"}
	}
	return obj.returnFields
}

// GridServicerestartRequest represents Infoblox object grid:servicerestart:request
type GridServicerestartRequest struct {
	IBBase          `json:"-"`
	Ref             string     `json:"_ref,omitempty"`
	Error           string     `json:"error,omitempty"`
	Forced          bool       `json:"forced,omitempty"`
	Group           string     `json:"group,omitempty"`
	LastUpdatedTime *time.Time `json:"last_updated_time,omitempty"`
	Member          string     `json:"member,omitempty"`
	Needed          string     `json:"needed,omitempty"`
	Order           int        `json:"order,omitempty"`
	Result          string     `json:"result,omitempty"`
	Service         string     `json:"service,omitempty"`
	State           string     `json:"state,omitempty"`
}

func (GridServicerestartRequest) ObjectType() string {
	return "grid:servicerestart:request"
}

func (obj GridServicerestartRequest) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"error", "group", "result", "state"}
	}
	return obj.returnFields
}

// Grid represents Infoblox object grid
type Grid struct {
	IBBase                           `json:"-"`
	Ref                              string                             `json:"_ref,omitempty"`
	AllowRecursiveDeletion           string                             `json:"allow_recursive_deletion,omitempty"`
	AuditLogFormat                   string                             `json:"audit_log_format,omitempty"`
	AuditToSyslogEnable              bool                               `json:"audit_to_syslog_enable,omitempty"`
	AutomatedTrafficCaptureSetting   *SettingAutomatedtrafficcapture    `json:"automated_traffic_capture_setting,omitempty"`
	ConsentBannerSetting             *GridConsentbannersetting          `json:"consent_banner_setting,omitempty"`
	ControlIpAddress                 *Controlipaddress                  `json:"control_ip_address,omitempty"`
	CspApiConfig                     *GridCspapiconfig                  `json:"csp_api_config,omitempty"`
	CspGridSetting                   *GridCspgridsetting                `json:"csp_grid_setting,omitempty"`
	DenyMgmSnapshots                 bool                               `json:"deny_mgm_snapshots,omitempty"`
	DescendantsAction                *ExtensibleattributedefDescendants `json:"descendants_action,omitempty"`
	DnsResolverSetting               *SettingDnsresolver                `json:"dns_resolver_setting,omitempty"`
	Dscp                             uint32                             `json:"dscp,omitempty"`
	EmailSetting                     *SettingEmail                      `json:"email_setting,omitempty"`
	EmptyRecycleBin                  *Emptyparams                       `json:"empty_recycle_bin,omitempty"`
	EnableGuiApiForLanVip            bool                               `json:"enable_gui_api_for_lan_vip,omitempty"`
	EnableLom                        bool                               `json:"enable_lom,omitempty"`
	EnableMemberRedirect             bool                               `json:"enable_member_redirect,omitempty"`
	EnableRecycleBin                 bool                               `json:"enable_recycle_bin,omitempty"`
	EnableRirSwip                    bool                               `json:"enable_rir_swip,omitempty"`
	ExternalSyslogBackupServers      []*Extsyslogbackupserver           `json:"external_syslog_backup_servers,omitempty"`
	ExternalSyslogServerEnable       bool                               `json:"external_syslog_server_enable,omitempty"`
	GenerateTsigKey                  *Generatetsigkeyparams             `json:"generate_tsig_key,omitempty"`
	GetAllTemplateVendorId           *Getvendoridentifiersparams        `json:"get_all_template_vendor_id,omitempty"`
	GetGridRevertStatus              *Getgridrevertstatusparams         `json:"get_grid_revert_status,omitempty"`
	GetRpzThreatDetails              *Threatdetails                     `json:"get_rpz_threat_details,omitempty"`
	GetTemplateSchemaVersions        *Gettemplateschemaversionsparams   `json:"get_template_schema_versions,omitempty"`
	HttpProxyServerSetting           *SettingHttpproxyserver            `json:"http_proxy_server_setting,omitempty"`
	InformationalBannerSetting       *GridInformationalbannersetting    `json:"informational_banner_setting,omitempty"`
	IsGridVisualizationVisible       bool                               `json:"is_grid_visualization_visible,omitempty"`
	Join                             *Gridjoin                          `json:"join,omitempty"`
	JoinMgm                          *Joinmgm                           `json:"join_mgm,omitempty"`
	LeaveMgm                         *Emptyparams                       `json:"leave_mgm,omitempty"`
	LockoutSetting                   *GridLockoutsetting                `json:"lockout_setting,omitempty"`
	LomUsers                         []*Lomuser                         `json:"lom_users,omitempty"`
	MemberUpgrade                    *Memberupgrade                     `json:"member_upgrade,omitempty"`
	MgmStrictDelegateMode            bool                               `json:"mgm_strict_delegate_mode,omitempty"`
	MsSetting                        *SettingMsserver                   `json:"ms_setting,omitempty"`
	Name                             string                             `json:"name,omitempty"`
	NatGroups                        []string                           `json:"nat_groups,omitempty"`
	NTPSetting                       *NTPSetting                        `json:"ntp_setting,omitempty"`
	ObjectsChangesTrackingSetting    *Objectschangestrackingsetting     `json:"objects_changes_tracking_setting,omitempty"`
	PasswordSetting                  *SettingPassword                   `json:"password_setting,omitempty"`
	PublishChanges                   *Publishchanges                    `json:"publish_changes,omitempty"`
	QueryFqdnOnMember                *Queryfqdnonmemberparams           `json:"query_fqdn_on_member,omitempty"`
	Requestrestartservicestatus      *Requestgridservicestatus          `json:"requestrestartservicestatus,omitempty"`
	RestartBannerSetting             *GridRestartbannersetting          `json:"restart_banner_setting,omitempty"`
	RestartStatus                    string                             `json:"restart_status,omitempty"`
	Restartservices                  *Gridrestartservices               `json:"restartservices,omitempty"`
	RpzHitRateInterval               uint32                             `json:"rpz_hit_rate_interval,omitempty"`
	RpzHitRateMaxQuery               uint32                             `json:"rpz_hit_rate_max_query,omitempty"`
	RpzHitRateMinQuery               uint32                             `json:"rpz_hit_rate_min_query,omitempty"`
	ScheduledBackup                  *Scheduledbackup                   `json:"scheduled_backup,omitempty"`
	Secret                           string                             `json:"secret,omitempty"`
	SecurityBannerSetting            *SettingSecuritybanner             `json:"security_banner_setting,omitempty"`
	SecuritySetting                  *SettingSecurity                   `json:"security_setting,omitempty"`
	ServiceStatus                    string                             `json:"service_status,omitempty"`
	SkipMemberUpgrade                *Skipmemberupgrade                 `json:"skip_member_upgrade,omitempty"`
	SnmpSetting                      *SettingSnmp                       `json:"snmp_setting,omitempty"`
	StartDiscovery                   *Startdiscovery                    `json:"start_discovery,omitempty"`
	SupportBundleDownloadTimeout     uint32                             `json:"support_bundle_download_timeout,omitempty"`
	SyslogFacility                   string                             `json:"syslog_facility,omitempty"`
	SyslogServers                    []*Syslogserver                    `json:"syslog_servers,omitempty"`
	SyslogSize                       uint32                             `json:"syslog_size,omitempty"`
	TestSyslogBackupServerConnection *Testsyslogbackup                  `json:"test_syslog_backup_server_connection,omitempty"`
	TestSyslogConnection             *Testsyslog                        `json:"test_syslog_connection,omitempty"`
	ThresholdTraps                   []*Thresholdtrap                   `json:"threshold_traps,omitempty"`
	TimeZone                         string                             `json:"time_zone,omitempty"`
	TokenUsageDelay                  uint32                             `json:"token_usage_delay,omitempty"`
	TrafficCaptureAuthDnsSetting     *SettingTriggeruthdnslatency       `json:"traffic_capture_auth_dns_setting,omitempty"`
	TrafficCaptureChrSetting         *SettingTrafficcapturechr          `json:"traffic_capture_chr_setting,omitempty"`
	TrafficCaptureQpsSetting         *SettingTrafficcaptureqps          `json:"traffic_capture_qps_setting,omitempty"`
	TrafficCaptureRecDnsSetting      *SettingTriggerrecdnslatency       `json:"traffic_capture_rec_dns_setting,omitempty"`
	TrafficCaptureRecQueriesSetting  *SettingTriggerrecqueries          `json:"traffic_capture_rec_queries_setting,omitempty"`
	TrapNotifications                []*Trapnotification                `json:"trap_notifications,omitempty"`
	UpdatesDownloadMemberConfig      []*Updatesdownloadmemberconfig     `json:"updates_download_member_config,omitempty"`
	Upgrade                          *Gridupgrade                       `json:"upgrade,omitempty"`
	UpgradeGroupNow                  *Upgradegroupnowparams             `json:"upgrade_group_now,omitempty"`
	UploadKeytab                     *Uploadkeytab                      `json:"upload_keytab,omitempty"`
	VpnPort                          uint32                             `json:"vpn_port,omitempty"`
}

func (Grid) ObjectType() string {
	return "grid"
}

func (obj Grid) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{}
	}
	return obj.returnFields
}

func NewGrid(grid Grid) *Grid {
	result := grid
	returnFields := []string{"name", "ntp_setting"}
	result.returnFields = returnFields
	return &result
}

// GridThreatanalytics represents Infoblox object grid:threatanalytics
type GridThreatanalytics struct {
	IBBase                                  `json:"-"`
	Ref                                     string                                  `json:"_ref,omitempty"`
	ConfigureDomainCollapsing               bool                                    `json:"configure_domain_collapsing,omitempty"`
	CurrentModuleset                        string                                  `json:"current_moduleset,omitempty"`
	CurrentWhitelist                        string                                  `json:"current_whitelist,omitempty"`
	DnsTunnelBlackListRpzZones              []*ZoneRp                               `json:"dns_tunnel_black_list_rpz_zones,omitempty"`
	DomainCollapsingLevel                   uint32                                  `json:"domain_collapsing_level,omitempty"`
	DownloadThreatAnalyticsModulesetUpdate  *Emptyparams                            `json:"download_threat_analytics_moduleset_update,omitempty"`
	DownloadThreatAnalyticsWhitelistUpdate  *Downloadthreatanalyticswhitelistparams `json:"download_threat_analytics_whitelist_update,omitempty"`
	EnableAutoDownload                      bool                                    `json:"enable_auto_download,omitempty"`
	EnableScheduledDownload                 bool                                    `json:"enable_scheduled_download,omitempty"`
	EnableWhitelistAutoDownload             bool                                    `json:"enable_whitelist_auto_download,omitempty"`
	EnableWhitelistScheduledDownload        bool                                    `json:"enable_whitelist_scheduled_download,omitempty"`
	LastCheckedForUpdate                    *time.Time                              `json:"last_checked_for_update,omitempty"`
	LastCheckedForWhitelistUpdate           *time.Time                              `json:"last_checked_for_whitelist_update,omitempty"`
	LastModuleUpdateTime                    *time.Time                              `json:"last_module_update_time,omitempty"`
	LastModuleUpdateVersion                 string                                  `json:"last_module_update_version,omitempty"`
	LastWhitelistUpdateTime                 *time.Time                              `json:"last_whitelist_update_time,omitempty"`
	LastWhitelistUpdateVersion              string                                  `json:"last_whitelist_update_version,omitempty"`
	ModuleUpdatePolicy                      string                                  `json:"module_update_policy,omitempty"`
	MoveBlacklistRpzToWhiteList             *Moveblacklistrpztowhitelistparams      `json:"move_blacklist_rpz_to_white_list,omitempty"`
	Name                                    string                                  `json:"name,omitempty"`
	ScheduledDownload                       *SettingSchedule                        `json:"scheduled_download,omitempty"`
	ScheduledWhitelistDownload              *SettingSchedule                        `json:"scheduled_whitelist_download,omitempty"`
	SetLastUploadedThreatAnalyticsModuleset *Setlastuploadedmodulesetparams         `json:"set_last_uploaded_threat_analytics_moduleset,omitempty"`
	TestThreatAnalyticsServerConnectivity   *Testanalyticsserverconnectivityparams  `json:"test_threat_analytics_server_connectivity,omitempty"`
	UpdateThreatAnalyticsModuleset          *Emptyparams                            `json:"update_threat_analytics_moduleset,omitempty"`
	WhitelistUpdatePolicy                   string                                  `json:"whitelist_update_policy,omitempty"`
}

func (GridThreatanalytics) ObjectType() string {
	return "grid:threatanalytics"
}

func (obj GridThreatanalytics) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"enable_auto_download", "enable_scheduled_download", "module_update_policy", "name"}
	}
	return obj.returnFields
}

// GridDns represents Infoblox object grid:dns
type GridDns struct {
	IBBase                              `json:"-"`
	Ref                                 string                        `json:"_ref,omitempty"`
	AddClientIpMacOptions               bool                          `json:"add_client_ip_mac_options,omitempty"`
	AllowBulkhostDdns                   string                        `json:"allow_bulkhost_ddns,omitempty"`
	AllowGssTsigZoneUpdates             bool                          `json:"allow_gss_tsig_zone_updates,omitempty"`
	AllowQuery                          []*Addressac                  `json:"allow_query,omitempty"`
	AllowRecursiveQuery                 bool                          `json:"allow_recursive_query,omitempty"`
	AllowTransfer                       []*Addressac                  `json:"allow_transfer,omitempty"`
	AllowUpdate                         []*Addressac                  `json:"allow_update,omitempty"`
	AnonymizeResponseLogging            bool                          `json:"anonymize_response_logging,omitempty"`
	AttackMitigation                    *GridAttackmitigation         `json:"attack_mitigation,omitempty"`
	AutoBlackhole                       *GridAutoblackhole            `json:"auto_blackhole,omitempty"`
	BindCheckNamesPolicy                string                        `json:"bind_check_names_policy,omitempty"`
	BindHostnameDirective               string                        `json:"bind_hostname_directive,omitempty"`
	BlackholeList                       []*Addressac                  `json:"blackhole_list,omitempty"`
	BlacklistAction                     string                        `json:"blacklist_action,omitempty"`
	BlacklistLogQuery                   bool                          `json:"blacklist_log_query,omitempty"`
	BlacklistRedirectAddresses          []string                      `json:"blacklist_redirect_addresses,omitempty"`
	BlacklistRedirectTtl                uint32                        `json:"blacklist_redirect_ttl,omitempty"`
	BlacklistRulesets                   []string                      `json:"blacklist_rulesets,omitempty"`
	BulkHostNameTemplates               []*Bulkhostnametemplate       `json:"bulk_host_name_templates,omitempty"`
	CaptureDnsQueriesOnAllDomains       bool                          `json:"capture_dns_queries_on_all_domains,omitempty"`
	CheckNamesForDdnsAndZoneTransfer    bool                          `json:"check_names_for_ddns_and_zone_transfer,omitempty"`
	ClientSubnetDomains                 []*Clientsubnetdomain         `json:"client_subnet_domains,omitempty"`
	ClientSubnetIpv4PrefixLength        uint32                        `json:"client_subnet_ipv4_prefix_length,omitempty"`
	ClientSubnetIpv6PrefixLength        uint32                        `json:"client_subnet_ipv6_prefix_length,omitempty"`
	CopyClientIpMacOptions              bool                          `json:"copy_client_ip_mac_options,omitempty"`
	CopyXferToNotify                    bool                          `json:"copy_xfer_to_notify,omitempty"`
	CustomRootNameServers               []NameServer                  `json:"custom_root_name_servers,omitempty"`
	DdnsForceCreationTimestampUpdate    bool                          `json:"ddns_force_creation_timestamp_update,omitempty"`
	DdnsPrincipalGroup                  string                        `json:"ddns_principal_group,omitempty"`
	DdnsPrincipalTracking               bool                          `json:"ddns_principal_tracking,omitempty"`
	DdnsRestrictPatterns                bool                          `json:"ddns_restrict_patterns,omitempty"`
	DdnsRestrictPatternsList            []string                      `json:"ddns_restrict_patterns_list,omitempty"`
	DdnsRestrictProtected               bool                          `json:"ddns_restrict_protected,omitempty"`
	DdnsRestrictSecure                  bool                          `json:"ddns_restrict_secure,omitempty"`
	DdnsRestrictStatic                  bool                          `json:"ddns_restrict_static,omitempty"`
	DefaultBulkHostNameTemplate         string                        `json:"default_bulk_host_name_template,omitempty"`
	DefaultTtl                          uint32                        `json:"default_ttl,omitempty"`
	DisableEdns                         bool                          `json:"disable_edns,omitempty"`
	Dns64Groups                         []string                      `json:"dns64_groups,omitempty"`
	DnsCacheAccelerationTtl             uint32                        `json:"dns_cache_acceleration_ttl,omitempty"`
	DnsHealthCheckAnycastControl        bool                          `json:"dns_health_check_anycast_control,omitempty"`
	DnsHealthCheckDomainList            []string                      `json:"dns_health_check_domain_list,omitempty"`
	DnsHealthCheckInterval              uint32                        `json:"dns_health_check_interval,omitempty"`
	DnsHealthCheckRecursionFlag         bool                          `json:"dns_health_check_recursion_flag,omitempty"`
	DnsHealthCheckRetries               uint32                        `json:"dns_health_check_retries,omitempty"`
	DnsHealthCheckTimeout               uint32                        `json:"dns_health_check_timeout,omitempty"`
	DnsQueryCaptureFileTimeLimit        uint32                        `json:"dns_query_capture_file_time_limit,omitempty"`
	DnssecBlacklistEnabled              bool                          `json:"dnssec_blacklist_enabled,omitempty"`
	DnssecDns64Enabled                  bool                          `json:"dnssec_dns64_enabled,omitempty"`
	DnssecEnabled                       bool                          `json:"dnssec_enabled,omitempty"`
	DnssecExpiredSignaturesEnabled      bool                          `json:"dnssec_expired_signatures_enabled,omitempty"`
	DnssecKeyParams                     *Dnsseckeyparams              `json:"dnssec_key_params,omitempty"`
	DnssecNegativeTrustAnchors          []string                      `json:"dnssec_negative_trust_anchors,omitempty"`
	DnssecNxdomainEnabled               bool                          `json:"dnssec_nxdomain_enabled,omitempty"`
	DnssecRpzEnabled                    bool                          `json:"dnssec_rpz_enabled,omitempty"`
	DnssecTrustedKeys                   []*Dnssectrustedkey           `json:"dnssec_trusted_keys,omitempty"`
	DnssecValidationEnabled             bool                          `json:"dnssec_validation_enabled,omitempty"`
	DnstapSetting                       *Dnstapsetting                `json:"dnstap_setting,omitempty"`
	DomainsToCaptureDnsQueries          []string                      `json:"domains_to_capture_dns_queries,omitempty"`
	DtcDnsQueriesSpecificBehavior       string                        `json:"dtc_dns_queries_specific_behavior,omitempty"`
	DtcDnssecMode                       string                        `json:"dtc_dnssec_mode,omitempty"`
	DtcEdnsPreferClientSubnet           bool                          `json:"dtc_edns_prefer_client_subnet,omitempty"`
	DtcScheduledBackup                  *Scheduledbackup              `json:"dtc_scheduled_backup,omitempty"`
	DtcTopologyEaList                   []string                      `json:"dtc_topology_ea_list,omitempty"`
	EdnsUdpSize                         uint32                        `json:"edns_udp_size,omitempty"`
	Email                               string                        `json:"email,omitempty"`
	EnableBlackhole                     bool                          `json:"enable_blackhole,omitempty"`
	EnableBlacklist                     bool                          `json:"enable_blacklist,omitempty"`
	EnableCaptureDnsQueries             bool                          `json:"enable_capture_dns_queries,omitempty"`
	EnableCaptureDnsResponses           bool                          `json:"enable_capture_dns_responses,omitempty"`
	EnableClientSubnetForwarding        bool                          `json:"enable_client_subnet_forwarding,omitempty"`
	EnableClientSubnetRecursive         bool                          `json:"enable_client_subnet_recursive,omitempty"`
	EnableDeleteAssociatedPtr           bool                          `json:"enable_delete_associated_ptr,omitempty"`
	EnableDns64                         bool                          `json:"enable_dns64,omitempty"`
	EnableDnsHealthCheck                bool                          `json:"enable_dns_health_check,omitempty"`
	EnableDnstapQueries                 bool                          `json:"enable_dnstap_queries,omitempty"`
	EnableDnstapResponses               bool                          `json:"enable_dnstap_responses,omitempty"`
	EnableExcludedDomainNames           bool                          `json:"enable_excluded_domain_names,omitempty"`
	EnableFixedRrsetOrderFqdns          bool                          `json:"enable_fixed_rrset_order_fqdns,omitempty"`
	EnableFtc                           bool                          `json:"enable_ftc,omitempty"`
	EnableGssTsig                       bool                          `json:"enable_gss_tsig,omitempty"`
	EnableHostRrsetOrder                bool                          `json:"enable_host_rrset_order,omitempty"`
	EnableHsmSigning                    bool                          `json:"enable_hsm_signing,omitempty"`
	EnableNotifySourcePort              bool                          `json:"enable_notify_source_port,omitempty"`
	EnableQueryRewrite                  bool                          `json:"enable_query_rewrite,omitempty"`
	EnableQuerySourcePort               bool                          `json:"enable_query_source_port,omitempty"`
	ExcludedDomainNames                 []string                      `json:"excluded_domain_names,omitempty"`
	ExpireAfter                         uint32                        `json:"expire_after,omitempty"`
	FileTransferSetting                 *Filetransfersetting          `json:"file_transfer_setting,omitempty"`
	FilterAaaa                          string                        `json:"filter_aaaa,omitempty"`
	FilterAaaaList                      []*Addressac                  `json:"filter_aaaa_list,omitempty"`
	FixedRrsetOrderFqdns                []*GridDnsFixedrrsetorderfqdn `json:"fixed_rrset_order_fqdns,omitempty"`
	ForwardOnly                         bool                          `json:"forward_only,omitempty"`
	ForwardUpdates                      bool                          `json:"forward_updates,omitempty"`
	Forwarders                          []string                      `json:"forwarders,omitempty"`
	FtcExpiredRecordTimeout             uint32                        `json:"ftc_expired_record_timeout,omitempty"`
	FtcExpiredRecordTtl                 uint32                        `json:"ftc_expired_record_ttl,omitempty"`
	GssTsigKeys                         []*Kerberoskey                `json:"gss_tsig_keys,omitempty"`
	LameTtl                             uint32                        `json:"lame_ttl,omitempty"`
	LastQueriedAcl                      []*Addressac                  `json:"last_queried_acl,omitempty"`
	LoggingCategories                   *GridLoggingcategories        `json:"logging_categories,omitempty"`
	MaxCacheTtl                         uint32                        `json:"max_cache_ttl,omitempty"`
	MaxCachedLifetime                   uint32                        `json:"max_cached_lifetime,omitempty"`
	MaxNcacheTtl                        uint32                        `json:"max_ncache_ttl,omitempty"`
	MaxUdpSize                          uint32                        `json:"max_udp_size,omitempty"`
	MemberSecondaryNotify               bool                          `json:"member_secondary_notify,omitempty"`
	NegativeTtl                         uint32                        `json:"negative_ttl,omitempty"`
	NotifyDelay                         uint32                        `json:"notify_delay,omitempty"`
	NotifySourcePort                    uint32                        `json:"notify_source_port,omitempty"`
	NsgroupDefault                      string                        `json:"nsgroup_default,omitempty"`
	Nsgroups                            []string                      `json:"nsgroups,omitempty"`
	NxdomainLogQuery                    bool                          `json:"nxdomain_log_query,omitempty"`
	NxdomainRedirect                    bool                          `json:"nxdomain_redirect,omitempty"`
	NxdomainRedirectAddresses           []string                      `json:"nxdomain_redirect_addresses,omitempty"`
	NxdomainRedirectAddressesV6         []string                      `json:"nxdomain_redirect_addresses_v6,omitempty"`
	NxdomainRedirectTtl                 uint32                        `json:"nxdomain_redirect_ttl,omitempty"`
	NxdomainRulesets                    []string                      `json:"nxdomain_rulesets,omitempty"`
	PreserveHostRrsetOrderOnSecondaries bool                          `json:"preserve_host_rrset_order_on_secondaries,omitempty"`
	ProtocolRecordNamePolicies          []*Recordnamepolicy           `json:"protocol_record_name_policies,omitempty"`
	QueryRewriteDomainNames             []string                      `json:"query_rewrite_domain_names,omitempty"`
	QueryRewritePrefix                  string                        `json:"query_rewrite_prefix,omitempty"`
	QuerySourcePort                     uint32                        `json:"query_source_port,omitempty"`
	RecursiveQueryList                  []*Addressac                  `json:"recursive_query_list,omitempty"`
	RefreshTimer                        uint32                        `json:"refresh_timer,omitempty"`
	ResolverQueryTimeout                uint32                        `json:"resolver_query_timeout,omitempty"`
	ResponseRateLimiting                *GridResponseratelimiting     `json:"response_rate_limiting,omitempty"`
	RestartSetting                      *GridServicerestart           `json:"restart_setting,omitempty"`
	RetryTimer                          uint32                        `json:"retry_timer,omitempty"`
	RootNameServerType                  string                        `json:"root_name_server_type,omitempty"`
	RpzDisableNsdnameNsip               bool                          `json:"rpz_disable_nsdname_nsip,omitempty"`
	RpzDropIpRuleEnabled                bool                          `json:"rpz_drop_ip_rule_enabled,omitempty"`
	RpzDropIpRuleMinPrefixLengthIpv4    uint32                        `json:"rpz_drop_ip_rule_min_prefix_length_ipv4,omitempty"`
	RpzDropIpRuleMinPrefixLengthIpv6    uint32                        `json:"rpz_drop_ip_rule_min_prefix_length_ipv6,omitempty"`
	RpzQnameWaitRecurse                 bool                          `json:"rpz_qname_wait_recurse,omitempty"`
	RunScavenging                       *Runscavenging                `json:"run_scavenging,omitempty"`
	ScavengingSettings                  *SettingScavenging            `json:"scavenging_settings,omitempty"`
	SerialQueryRate                     uint32                        `json:"serial_query_rate,omitempty"`
	ServerIdDirective                   string                        `json:"server_id_directive,omitempty"`
	Sortlist                            []*Sortlist                   `json:"sortlist,omitempty"`
	StoreLocally                        bool                          `json:"store_locally,omitempty"`
	SyslogFacility                      string                        `json:"syslog_facility,omitempty"`
	TransferExcludedServers             []string                      `json:"transfer_excluded_servers,omitempty"`
	TransferFormat                      string                        `json:"transfer_format,omitempty"`
	TransfersIn                         uint32                        `json:"transfers_in,omitempty"`
	TransfersOut                        uint32                        `json:"transfers_out,omitempty"`
	TransfersPerNs                      uint32                        `json:"transfers_per_ns,omitempty"`
	ZoneDeletionDoubleConfirm           bool                          `json:"zone_deletion_double_confirm,omitempty"`
}

func (GridDns) ObjectType() string {
	return "grid:dns"
}

func (obj GridDns) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{}
	}
	return obj.returnFields
}

// GridThreatprotection represents Infoblox object grid:threatprotection
type GridThreatprotection struct {
	IBBase                                `json:"-"`
	Ref                                   string                     `json:"_ref,omitempty"`
	AtpObjectReset                        *Atpobjectreset            `json:"atp_object_reset,omitempty"`
	CurrentRuleset                        string                     `json:"current_ruleset,omitempty"`
	DisableMultipleDnsTcpRequest          bool                       `json:"disable_multiple_dns_tcp_request,omitempty"`
	EnableAccelRespBeforeThreatProtection bool                       `json:"enable_accel_resp_before_threat_protection,omitempty"`
	EnableAutoDownload                    bool                       `json:"enable_auto_download,omitempty"`
	EnableNatRules                        bool                       `json:"enable_nat_rules,omitempty"`
	EnableScheduledDownload               bool                       `json:"enable_scheduled_download,omitempty"`
	EventsPerSecondPerRule                uint32                     `json:"events_per_second_per_rule,omitempty"`
	GridName                              string                     `json:"grid_name,omitempty"`
	LastCheckedForUpdate                  *time.Time                 `json:"last_checked_for_update,omitempty"`
	LastRuleUpdateTimestamp               *time.Time                 `json:"last_rule_update_timestamp,omitempty"`
	LastRuleUpdateVersion                 string                     `json:"last_rule_update_version,omitempty"`
	NatRules                              []*ThreatprotectionNatrule `json:"nat_rules,omitempty"`
	OutboundSettings                      *SettingAtpoutbound        `json:"outbound_settings,omitempty"`
	RuleUpdatePolicy                      string                     `json:"rule_update_policy,omitempty"`
	ScheduledDownload                     *SettingSchedule           `json:"scheduled_download,omitempty"`
	TestAtpServerConnectivity             *Testatpserverconnectivity `json:"test_atp_server_connectivity,omitempty"`
}

func (GridThreatprotection) ObjectType() string {
	return "grid:threatprotection"
}

func (obj GridThreatprotection) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"grid_name"}
	}
	return obj.returnFields
}

// GridX509certificate represents Infoblox object grid:x509certificate
type GridX509certificate struct {
	IBBase         `json:"-"`
	Ref            string     `json:"_ref,omitempty"`
	Issuer         string     `json:"issuer,omitempty"`
	Serial         string     `json:"serial,omitempty"`
	Subject        string     `json:"subject,omitempty"`
	ValidNotAfter  *time.Time `json:"valid_not_after,omitempty"`
	ValidNotBefore *time.Time `json:"valid_not_before,omitempty"`
}

func (GridX509certificate) ObjectType() string {
	return "grid:x509certificate"
}

func (obj GridX509certificate) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"issuer", "serial", "subject"}
	}
	return obj.returnFields
}

// Hostnamerewritepolicy represents Infoblox object hostnamerewritepolicy
type Hostnamerewritepolicy struct {
	IBBase               `json:"-"`
	Ref                  string `json:"_ref,omitempty"`
	IsDefault            bool   `json:"is_default,omitempty"`
	Name                 string `json:"name,omitempty"`
	PreDefined           bool   `json:"pre_defined,omitempty"`
	ReplacementCharacter string `json:"replacement_character,omitempty"`
	ValidCharacters      string `json:"valid_characters,omitempty"`
}

func (Hostnamerewritepolicy) ObjectType() string {
	return "hostnamerewritepolicy"
}

func (obj Hostnamerewritepolicy) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name", "replacement_character", "valid_characters"}
	}
	return obj.returnFields
}

// HsmAllgroups represents Infoblox object hsm:allgroups
type HsmAllgroups struct {
	IBBase `json:"-"`
	Ref    string             `json:"_ref,omitempty"`
	Groups []*HsmSafenetgroup `json:"groups,omitempty"`
}

func (HsmAllgroups) ObjectType() string {
	return "hsm:allgroups"
}

func (obj HsmAllgroups) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"groups"}
	}
	return obj.returnFields
}

// HsmThalesgroup represents Infoblox object hsm:thalesgroup
type HsmThalesgroup struct {
	IBBase        `json:"-"`
	Ref           string               `json:"_ref,omitempty"`
	CardName      string               `json:"card_name,omitempty"`
	Comment       string               `json:"comment,omitempty"`
	KeyServerIp   string               `json:"key_server_ip,omitempty"`
	KeyServerPort uint32               `json:"key_server_port,omitempty"`
	Name          string               `json:"name,omitempty"`
	PassPhrase    string               `json:"pass_phrase,omitempty"`
	Protection    string               `json:"protection,omitempty"`
	RefreshHsm    *Hsmrefreshparams    `json:"refresh_hsm,omitempty"`
	Status        string               `json:"status,omitempty"`
	TestHsmStatus *Hsmteststatusparams `json:"test_hsm_status,omitempty"`
	ThalesHsm     []*HsmThales         `json:"thales_hsm,omitempty"`
}

func (HsmThalesgroup) ObjectType() string {
	return "hsm:thalesgroup"
}

func (obj HsmThalesgroup) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "key_server_ip", "name"}
	}
	return obj.returnFields
}

// HsmSafenetgroup represents Infoblox object hsm:safenetgroup
type HsmSafenetgroup struct {
	IBBase        `json:"-"`
	Ref           string               `json:"_ref,omitempty"`
	Comment       string               `json:"comment,omitempty"`
	GroupSn       string               `json:"group_sn,omitempty"`
	HsmSafenet    []*HsmSafenet        `json:"hsm_safenet,omitempty"`
	HsmVersion    string               `json:"hsm_version,omitempty"`
	Name          string               `json:"name,omitempty"`
	PassPhrase    string               `json:"pass_phrase,omitempty"`
	RefreshHsm    *Hsmrefreshparams    `json:"refresh_hsm,omitempty"`
	Status        string               `json:"status,omitempty"`
	TestHsmStatus *Hsmteststatusparams `json:"test_hsm_status,omitempty"`
}

func (HsmSafenetgroup) ObjectType() string {
	return "hsm:safenetgroup"
}

func (obj HsmSafenetgroup) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "hsm_version", "name"}
	}
	return obj.returnFields
}

// IpamStatistics represents Infoblox object ipam:statistics
type IpamStatistics struct {
	IBBase            `json:"-"`
	Ref               string              `json:"_ref,omitempty"`
	Cidr              uint32              `json:"cidr,omitempty"`
	ConflictCount     uint32              `json:"conflict_count,omitempty"`
	MsAdUserData      *MsserverAduserData `json:"ms_ad_user_data,omitempty"`
	Network           string              `json:"network,omitempty"`
	NetworkView       string              `json:"network_view,omitempty"`
	UnmanagedCount    uint32              `json:"unmanaged_count,omitempty"`
	Utilization       uint32              `json:"utilization,omitempty"`
	UtilizationUpdate *time.Time          `json:"utilization_update,omitempty"`
}

func (IpamStatistics) ObjectType() string {
	return "ipam:statistics"
}

func (obj IpamStatistics) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"cidr", "network", "network_view"}
	}
	return obj.returnFields
}

// Ipv4address represents Infoblox object ipv4address
type Ipv4address struct {
	IBBase               `json:"-"`
	Ref                  string              `json:"_ref,omitempty"`
	Comment              string              `json:"comment,omitempty"`
	ConflictTypes        []string            `json:"conflict_types,omitempty"`
	DhcpClientIdentifier string              `json:"dhcp_client_identifier,omitempty"`
	DiscoverNowStatus    string              `json:"discover_now_status,omitempty"`
	DiscoveredData       *Discoverydata      `json:"discovered_data,omitempty"`
	Ea                   EA                  `json:"extattrs,omitempty"`
	Fingerprint          string              `json:"fingerprint,omitempty"`
	IpAddress            string              `json:"ip_address,omitempty"`
	IsConflict           bool                `json:"is_conflict,omitempty"`
	IsInvalidMac         bool                `json:"is_invalid_mac,omitempty"`
	LeaseState           string              `json:"lease_state,omitempty"`
	MacAddress           string              `json:"mac_address,omitempty"`
	MsAdUserData         *MsserverAduserData `json:"ms_ad_user_data,omitempty"`
	Names                []string            `json:"names,omitempty"`
	Network              string              `json:"network,omitempty"`
	NetworkView          string              `json:"network_view,omitempty"`
	Objects              string              `json:"objects,omitempty"`
	ReservedPort         string              `json:"reserved_port,omitempty"`
	Status               string              `json:"status,omitempty"`
	Types                []string            `json:"types,omitempty"`
	Usage                []string            `json:"usage,omitempty"`
	Username             string              `json:"username,omitempty"`
}

func (Ipv4address) ObjectType() string {
	return "ipv4address"
}

func (obj Ipv4address) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"dhcp_client_identifier", "ip_address", "is_conflict", "lease_state", "mac_address", "names", "network", "network_view", "objects", "status", "types", "usage", "username"}
	}
	return obj.returnFields
}

// Ipv6dhcpoptiondefinition represents Infoblox object ipv6dhcpoptiondefinition
type Ipv6dhcpoptiondefinition struct {
	IBBase `json:"-"`
	Ref    string `json:"_ref,omitempty"`
	Code   uint32 `json:"code,omitempty"`
	Name   string `json:"name,omitempty"`
	Space  string `json:"space,omitempty"`
	Type   string `json:"type,omitempty"`
}

func (Ipv6dhcpoptiondefinition) ObjectType() string {
	return "ipv6dhcpoptiondefinition"
}

func (obj Ipv6dhcpoptiondefinition) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"code", "name", "type"}
	}
	return obj.returnFields
}

// Ipv6dhcpoptionspace represents Infoblox object ipv6dhcpoptionspace
type Ipv6dhcpoptionspace struct {
	IBBase            `json:"-"`
	Ref               string   `json:"_ref,omitempty"`
	Comment           string   `json:"comment,omitempty"`
	EnterpriseNumber  uint32   `json:"enterprise_number,omitempty"`
	Name              string   `json:"name,omitempty"`
	OptionDefinitions []string `json:"option_definitions,omitempty"`
}

func (Ipv6dhcpoptionspace) ObjectType() string {
	return "ipv6dhcpoptionspace"
}

func (obj Ipv6dhcpoptionspace) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "enterprise_number", "name"}
	}
	return obj.returnFields
}

// Ipv6fixedaddresstemplate represents Infoblox object ipv6fixedaddresstemplate
type Ipv6fixedaddresstemplate struct {
	IBBase               `json:"-"`
	Ref                  string        `json:"_ref,omitempty"`
	Comment              string        `json:"comment,omitempty"`
	DomainName           string        `json:"domain_name,omitempty"`
	DomainNameServers    []string      `json:"domain_name_servers,omitempty"`
	Ea                   EA            `json:"extattrs,omitempty"`
	Name                 string        `json:"name,omitempty"`
	NumberOfAddresses    uint32        `json:"number_of_addresses,omitempty"`
	Offset               uint32        `json:"offset,omitempty"`
	Options              []*Dhcpoption `json:"options,omitempty"`
	PreferredLifetime    uint32        `json:"preferred_lifetime,omitempty"`
	UseDomainName        bool          `json:"use_domain_name,omitempty"`
	UseDomainNameServers bool          `json:"use_domain_name_servers,omitempty"`
	UseOptions           bool          `json:"use_options,omitempty"`
	UsePreferredLifetime bool          `json:"use_preferred_lifetime,omitempty"`
	UseValidLifetime     bool          `json:"use_valid_lifetime,omitempty"`
	ValidLifetime        uint32        `json:"valid_lifetime,omitempty"`
}

func (Ipv6fixedaddresstemplate) ObjectType() string {
	return "ipv6fixedaddresstemplate"
}

func (obj Ipv6fixedaddresstemplate) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name"}
	}
	return obj.returnFields
}

// Ipv6address represents Infoblox object ipv6address
type Ipv6address struct {
	IBBase            `json:"-"`
	Ref               string              `json:"_ref,omitempty"`
	Comment           string              `json:"comment,omitempty"`
	ConflictTypes     []string            `json:"conflict_types,omitempty"`
	DiscoverNowStatus string              `json:"discover_now_status,omitempty"`
	DiscoveredData    *Discoverydata      `json:"discovered_data,omitempty"`
	Duid              string              `json:"duid,omitempty"`
	Ea                EA                  `json:"extattrs,omitempty"`
	Fingerprint       string              `json:"fingerprint,omitempty"`
	IpAddress         string              `json:"ip_address,omitempty"`
	IsConflict        bool                `json:"is_conflict,omitempty"`
	LeaseState        string              `json:"lease_state,omitempty"`
	MsAdUserData      *MsserverAduserData `json:"ms_ad_user_data,omitempty"`
	Names             []string            `json:"names,omitempty"`
	Network           string              `json:"network,omitempty"`
	NetworkView       string              `json:"network_view,omitempty"`
	Objects           string              `json:"objects,omitempty"`
	ReservedPort      string              `json:"reserved_port,omitempty"`
	Status            string              `json:"status,omitempty"`
	Types             []string            `json:"types,omitempty"`
	Usage             []string            `json:"usage,omitempty"`
}

func (Ipv6address) ObjectType() string {
	return "ipv6address"
}

func (obj Ipv6address) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"duid", "ip_address", "is_conflict", "lease_state", "names", "network", "network_view", "objects", "status", "types", "usage"}
	}
	return obj.returnFields
}

// Ipv6FixedAddress represents Infoblox object ipv6fixedaddress
type Ipv6FixedAddress struct {
	IBBase                   `json:"-"`
	Ref                      string                    `json:"_ref,omitempty"`
	AddressType              string                    `json:"address_type,omitempty"`
	AllowTelnet              bool                      `json:"allow_telnet,omitempty"`
	CliCredentials           []*DiscoveryClicredential `json:"cli_credentials,omitempty"`
	CloudInfo                *GridCloudapiInfo         `json:"cloud_info,omitempty"`
	Comment                  string                    `json:"comment,omitempty"`
	DeviceDescription        string                    `json:"device_description,omitempty"`
	DeviceLocation           string                    `json:"device_location,omitempty"`
	DeviceType               string                    `json:"device_type,omitempty"`
	DeviceVendor             string                    `json:"device_vendor,omitempty"`
	Disable                  bool                      `json:"disable,omitempty"`
	DisableDiscovery         bool                      `json:"disable_discovery,omitempty"`
	DiscoverNowStatus        string                    `json:"discover_now_status,omitempty"`
	DiscoveredData           *Discoverydata            `json:"discovered_data,omitempty"`
	DomainName               string                    `json:"domain_name,omitempty"`
	DomainNameServers        []string                  `json:"domain_name_servers,omitempty"`
	Duid                     string                    `json:"duid,omitempty"`
	EnableImmediateDiscovery bool                      `json:"enable_immediate_discovery,omitempty"`
	Ea                       EA                        `json:"extattrs,omitempty"`
	Ipv6Addr                 string                    `json:"ipv6addr,omitempty"`
	Ipv6prefix               string                    `json:"ipv6prefix,omitempty"`
	Ipv6prefixBits           uint32                    `json:"ipv6prefix_bits,omitempty"`
	MsAdUserData             *MsserverAduserData       `json:"ms_ad_user_data,omitempty"`
	Name                     string                    `json:"name,omitempty"`
	Network                  string                    `json:"network,omitempty"`
	NetworkView              string                    `json:"network_view,omitempty"`
	Options                  []*Dhcpoption             `json:"options,omitempty"`
	PreferredLifetime        uint32                    `json:"preferred_lifetime,omitempty"`
	ReservedInterface        string                    `json:"reserved_interface,omitempty"`
	RestartIfNeeded          bool                      `json:"restart_if_needed,omitempty"`
	Snmp3Credential          *DiscoverySnmp3credential `json:"snmp3_credential,omitempty"`
	SnmpCredential           *DiscoverySnmpcredential  `json:"snmp_credential,omitempty"`
	Template                 string                    `json:"template,omitempty"`
	UseCliCredentials        bool                      `json:"use_cli_credentials,omitempty"`
	UseDomainName            bool                      `json:"use_domain_name,omitempty"`
	UseDomainNameServers     bool                      `json:"use_domain_name_servers,omitempty"`
	UseOptions               bool                      `json:"use_options,omitempty"`
	UsePreferredLifetime     bool                      `json:"use_preferred_lifetime,omitempty"`
	UseSnmp3Credential       bool                      `json:"use_snmp3_credential,omitempty"`
	UseSnmpCredential        bool                      `json:"use_snmp_credential,omitempty"`
	UseValidLifetime         bool                      `json:"use_valid_lifetime,omitempty"`
	ValidLifetime            uint32                    `json:"valid_lifetime,omitempty"`
}

func (Ipv6FixedAddress) ObjectType() string {
	return "ipv6fixedaddress"
}

func (obj Ipv6FixedAddress) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"duid", "ipv6addr", "network_view"}
	}
	return obj.returnFields
}

// Ipv6networktemplate represents Infoblox object ipv6networktemplate
type Ipv6networktemplate struct {
	IBBase                     `json:"-"`
	Ref                        string        `json:"_ref,omitempty"`
	AllowAnyNetmask            bool          `json:"allow_any_netmask,omitempty"`
	AutoCreateReversezone      bool          `json:"auto_create_reversezone,omitempty"`
	Cidr                       uint32        `json:"cidr,omitempty"`
	CloudApiCompatible         bool          `json:"cloud_api_compatible,omitempty"`
	Comment                    string        `json:"comment,omitempty"`
	DdnsDomainname             string        `json:"ddns_domainname,omitempty"`
	DdnsEnableOptionFqdn       bool          `json:"ddns_enable_option_fqdn,omitempty"`
	DdnsGenerateHostname       bool          `json:"ddns_generate_hostname,omitempty"`
	DdnsServerAlwaysUpdates    bool          `json:"ddns_server_always_updates,omitempty"`
	DdnsTtl                    uint32        `json:"ddns_ttl,omitempty"`
	DelegatedMember            *Dhcpmember   `json:"delegated_member,omitempty"`
	DomainName                 string        `json:"domain_name,omitempty"`
	DomainNameServers          []string      `json:"domain_name_servers,omitempty"`
	EnableDdns                 bool          `json:"enable_ddns,omitempty"`
	Ea                         EA            `json:"extattrs,omitempty"`
	FixedAddressTemplates      []string      `json:"fixed_address_templates,omitempty"`
	Ipv6prefix                 string        `json:"ipv6prefix,omitempty"`
	Members                    []*Dhcpmember `json:"members,omitempty"`
	Name                       string        `json:"name,omitempty"`
	Options                    []*Dhcpoption `json:"options,omitempty"`
	PreferredLifetime          uint32        `json:"preferred_lifetime,omitempty"`
	RangeTemplates             []string      `json:"range_templates,omitempty"`
	RecycleLeases              bool          `json:"recycle_leases,omitempty"`
	Rir                        string        `json:"rir,omitempty"`
	RirOrganization            string        `json:"rir_organization,omitempty"`
	RirRegistrationAction      string        `json:"rir_registration_action,omitempty"`
	RirRegistrationStatus      string        `json:"rir_registration_status,omitempty"`
	SendRirRequest             bool          `json:"send_rir_request,omitempty"`
	UpdateDnsOnLeaseRenewal    bool          `json:"update_dns_on_lease_renewal,omitempty"`
	UseDdnsDomainname          bool          `json:"use_ddns_domainname,omitempty"`
	UseDdnsEnableOptionFqdn    bool          `json:"use_ddns_enable_option_fqdn,omitempty"`
	UseDdnsGenerateHostname    bool          `json:"use_ddns_generate_hostname,omitempty"`
	UseDdnsTtl                 bool          `json:"use_ddns_ttl,omitempty"`
	UseDomainName              bool          `json:"use_domain_name,omitempty"`
	UseDomainNameServers       bool          `json:"use_domain_name_servers,omitempty"`
	UseEnableDdns              bool          `json:"use_enable_ddns,omitempty"`
	UseOptions                 bool          `json:"use_options,omitempty"`
	UsePreferredLifetime       bool          `json:"use_preferred_lifetime,omitempty"`
	UseRecycleLeases           bool          `json:"use_recycle_leases,omitempty"`
	UseUpdateDnsOnLeaseRenewal bool          `json:"use_update_dns_on_lease_renewal,omitempty"`
	UseValidLifetime           bool          `json:"use_valid_lifetime,omitempty"`
	ValidLifetime              uint32        `json:"valid_lifetime,omitempty"`
}

func (Ipv6networktemplate) ObjectType() string {
	return "ipv6networktemplate"
}

func (obj Ipv6networktemplate) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name"}
	}
	return obj.returnFields
}

// Ipv6rangetemplate represents Infoblox object ipv6rangetemplate
type Ipv6rangetemplate struct {
	IBBase                `json:"-"`
	Ref                   string                    `json:"_ref,omitempty"`
	CloudApiCompatible    bool                      `json:"cloud_api_compatible,omitempty"`
	Comment               string                    `json:"comment,omitempty"`
	DelegatedMember       *Dhcpmember               `json:"delegated_member,omitempty"`
	Exclude               []*Exclusionrangetemplate `json:"exclude,omitempty"`
	Member                *Dhcpmember               `json:"member,omitempty"`
	Name                  string                    `json:"name,omitempty"`
	NumberOfAddresses     uint32                    `json:"number_of_addresses,omitempty"`
	Offset                uint32                    `json:"offset,omitempty"`
	RecycleLeases         bool                      `json:"recycle_leases,omitempty"`
	ServerAssociationType string                    `json:"server_association_type,omitempty"`
	UseRecycleLeases      bool                      `json:"use_recycle_leases,omitempty"`
}

func (Ipv6rangetemplate) ObjectType() string {
	return "ipv6rangetemplate"
}

func (obj Ipv6rangetemplate) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name", "number_of_addresses", "offset"}
	}
	return obj.returnFields
}

// Kerberoskey represents Infoblox object kerberoskey
type Kerberoskey struct {
	IBBase          `json:"-"`
	Ref             string     `json:"_ref,omitempty"`
	Domain          string     `json:"domain,omitempty"`
	Enctype         string     `json:"enctype,omitempty"`
	InUse           bool       `json:"in_use,omitempty"`
	Members         []string   `json:"members,omitempty"`
	Principal       string     `json:"principal,omitempty"`
	UploadTimestamp *time.Time `json:"upload_timestamp,omitempty"`
	Version         uint32     `json:"version,omitempty"`
}

func (Kerberoskey) ObjectType() string {
	return "kerberoskey"
}

func (obj Kerberoskey) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"domain", "enctype", "in_use", "principal", "version"}
	}
	return obj.returnFields
}

// Ipv6sharednetwork represents Infoblox object ipv6sharednetwork
type Ipv6sharednetwork struct {
	IBBase                     `json:"-"`
	Ref                        string         `json:"_ref,omitempty"`
	Comment                    string         `json:"comment,omitempty"`
	DdnsDomainname             string         `json:"ddns_domainname,omitempty"`
	DdnsGenerateHostname       bool           `json:"ddns_generate_hostname,omitempty"`
	DdnsServerAlwaysUpdates    bool           `json:"ddns_server_always_updates,omitempty"`
	DdnsTtl                    uint32         `json:"ddns_ttl,omitempty"`
	DdnsUseOption81            bool           `json:"ddns_use_option81,omitempty"`
	Disable                    bool           `json:"disable,omitempty"`
	DomainName                 string         `json:"domain_name,omitempty"`
	DomainNameServers          []string       `json:"domain_name_servers,omitempty"`
	EnableDdns                 bool           `json:"enable_ddns,omitempty"`
	Ea                         EA             `json:"extattrs,omitempty"`
	Name                       string         `json:"name,omitempty"`
	NetworkView                string         `json:"network_view,omitempty"`
	Networks                   []*Ipv6Network `json:"networks,omitempty"`
	Options                    []*Dhcpoption  `json:"options,omitempty"`
	PreferredLifetime          uint32         `json:"preferred_lifetime,omitempty"`
	UpdateDnsOnLeaseRenewal    bool           `json:"update_dns_on_lease_renewal,omitempty"`
	UseDdnsDomainname          bool           `json:"use_ddns_domainname,omitempty"`
	UseDdnsGenerateHostname    bool           `json:"use_ddns_generate_hostname,omitempty"`
	UseDdnsTtl                 bool           `json:"use_ddns_ttl,omitempty"`
	UseDdnsUseOption81         bool           `json:"use_ddns_use_option81,omitempty"`
	UseDomainName              bool           `json:"use_domain_name,omitempty"`
	UseDomainNameServers       bool           `json:"use_domain_name_servers,omitempty"`
	UseEnableDdns              bool           `json:"use_enable_ddns,omitempty"`
	UseOptions                 bool           `json:"use_options,omitempty"`
	UsePreferredLifetime       bool           `json:"use_preferred_lifetime,omitempty"`
	UseUpdateDnsOnLeaseRenewal bool           `json:"use_update_dns_on_lease_renewal,omitempty"`
	UseValidLifetime           bool           `json:"use_valid_lifetime,omitempty"`
	ValidLifetime              uint32         `json:"valid_lifetime,omitempty"`
}

func (Ipv6sharednetwork) ObjectType() string {
	return "ipv6sharednetwork"
}

func (obj Ipv6sharednetwork) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name", "network_view", "networks"}
	}
	return obj.returnFields
}

// Ipv6Network represents Infoblox object ipv6network
type Ipv6Network struct {
	IBBase                           `json:"-"`
	Ref                              string                      `json:"_ref,omitempty"`
	AutoCreateReversezone            bool                        `json:"auto_create_reversezone,omitempty"`
	CloudInfo                        *GridCloudapiInfo           `json:"cloud_info,omitempty"`
	Comment                          string                      `json:"comment,omitempty"`
	DdnsDomainname                   string                      `json:"ddns_domainname,omitempty"`
	DdnsEnableOptionFqdn             bool                        `json:"ddns_enable_option_fqdn,omitempty"`
	DdnsGenerateHostname             bool                        `json:"ddns_generate_hostname,omitempty"`
	DdnsServerAlwaysUpdates          bool                        `json:"ddns_server_always_updates,omitempty"`
	DdnsTtl                          uint32                      `json:"ddns_ttl,omitempty"`
	DeleteReason                     string                      `json:"delete_reason,omitempty"`
	Disable                          bool                        `json:"disable,omitempty"`
	DiscoverNowStatus                string                      `json:"discover_now_status,omitempty"`
	DiscoveredBgpAs                  string                      `json:"discovered_bgp_as,omitempty"`
	DiscoveredBridgeDomain           string                      `json:"discovered_bridge_domain,omitempty"`
	DiscoveredTenant                 string                      `json:"discovered_tenant,omitempty"`
	DiscoveredVlanId                 string                      `json:"discovered_vlan_id,omitempty"`
	DiscoveredVlanName               string                      `json:"discovered_vlan_name,omitempty"`
	DiscoveredVrfDescription         string                      `json:"discovered_vrf_description,omitempty"`
	DiscoveredVrfName                string                      `json:"discovered_vrf_name,omitempty"`
	DiscoveredVrfRd                  string                      `json:"discovered_vrf_rd,omitempty"`
	DiscoveryBasicPollSettings       *DiscoveryBasicpollsettings `json:"discovery_basic_poll_settings,omitempty"`
	DiscoveryBlackoutSetting         *PropertiesBlackoutsetting  `json:"discovery_blackout_setting,omitempty"`
	DiscoveryEngineType              string                      `json:"discovery_engine_type,omitempty"`
	DiscoveryMember                  string                      `json:"discovery_member,omitempty"`
	DomainName                       string                      `json:"domain_name,omitempty"`
	DomainNameServers                []string                    `json:"domain_name_servers,omitempty"`
	EnableDdns                       bool                        `json:"enable_ddns,omitempty"`
	EnableDiscovery                  bool                        `json:"enable_discovery,omitempty"`
	EnableIfmapPublishing            bool                        `json:"enable_ifmap_publishing,omitempty"`
	EnableImmediateDiscovery         bool                        `json:"enable_immediate_discovery,omitempty"`
	EndpointSources                  []*CiscoiseEndpoint         `json:"endpoint_sources,omitempty"`
	ExpandNetwork                    *Expandnetwork              `json:"expand_network,omitempty"`
	Ea                               EA                          `json:"extattrs,omitempty"`
	LastRirRegistrationUpdateSent    *time.Time                  `json:"last_rir_registration_update_sent,omitempty"`
	LastRirRegistrationUpdateStatus  string                      `json:"last_rir_registration_update_status,omitempty"`
	Members                          []*Dhcpmember               `json:"members,omitempty"`
	MgmPrivate                       bool                        `json:"mgm_private,omitempty"`
	MgmPrivateOverridable            bool                        `json:"mgm_private_overridable,omitempty"`
	MsAdUserData                     *MsserverAduserData         `json:"ms_ad_user_data,omitempty"`
	Network                          string                      `json:"network,omitempty"`
	NetworkContainer                 string                      `json:"network_container,omitempty"`
	NetworkView                      string                      `json:"network_view,omitempty"`
	NextAvailableIp                  *Nextavailableip6           `json:"next_available_ip,omitempty"`
	NextAvailableNetwork             *Nextavailablenet6          `json:"next_available_network,omitempty"`
	NextAvailableVlan                *Nextavailablevlan          `json:"next_available_vlan,omitempty"`
	Options                          []*Dhcpoption               `json:"options,omitempty"`
	PortControlBlackoutSetting       *PropertiesBlackoutsetting  `json:"port_control_blackout_setting,omitempty"`
	PreferredLifetime                uint32                      `json:"preferred_lifetime,omitempty"`
	RecycleLeases                    bool                        `json:"recycle_leases,omitempty"`
	RestartIfNeeded                  bool                        `json:"restart_if_needed,omitempty"`
	Rir                              string                      `json:"rir,omitempty"`
	RirOrganization                  string                      `json:"rir_organization,omitempty"`
	RirRegistrationAction            string                      `json:"rir_registration_action,omitempty"`
	RirRegistrationStatus            string                      `json:"rir_registration_status,omitempty"`
	SamePortControlDiscoveryBlackout bool                        `json:"same_port_control_discovery_blackout,omitempty"`
	SendRirRequest                   bool                        `json:"send_rir_request,omitempty"`
	SplitNetwork                     *Splitipv6network           `json:"split_network,omitempty"`
	SubscribeSettings                *CiscoiseSubscribesetting   `json:"subscribe_settings,omitempty"`
	Template                         string                      `json:"template,omitempty"`
	Unmanaged                        bool                        `json:"unmanaged,omitempty"`
	UnmanagedCount                   uint32                      `json:"unmanaged_count,omitempty"`
	UpdateDnsOnLeaseRenewal          bool                        `json:"update_dns_on_lease_renewal,omitempty"`
	UseBlackoutSetting               bool                        `json:"use_blackout_setting,omitempty"`
	UseDdnsDomainname                bool                        `json:"use_ddns_domainname,omitempty"`
	UseDdnsEnableOptionFqdn          bool                        `json:"use_ddns_enable_option_fqdn,omitempty"`
	UseDdnsGenerateHostname          bool                        `json:"use_ddns_generate_hostname,omitempty"`
	UseDdnsTtl                       bool                        `json:"use_ddns_ttl,omitempty"`
	UseDiscoveryBasicPollingSettings bool                        `json:"use_discovery_basic_polling_settings,omitempty"`
	UseDomainName                    bool                        `json:"use_domain_name,omitempty"`
	UseDomainNameServers             bool                        `json:"use_domain_name_servers,omitempty"`
	UseEnableDdns                    bool                        `json:"use_enable_ddns,omitempty"`
	UseEnableDiscovery               bool                        `json:"use_enable_discovery,omitempty"`
	UseEnableIfmapPublishing         bool                        `json:"use_enable_ifmap_publishing,omitempty"`
	UseMgmPrivate                    bool                        `json:"use_mgm_private,omitempty"`
	UseOptions                       bool                        `json:"use_options,omitempty"`
	UsePreferredLifetime             bool                        `json:"use_preferred_lifetime,omitempty"`
	UseRecycleLeases                 bool                        `json:"use_recycle_leases,omitempty"`
	UseSubscribeSettings             bool                        `json:"use_subscribe_settings,omitempty"`
	UseUpdateDnsOnLeaseRenewal       bool                        `json:"use_update_dns_on_lease_renewal,omitempty"`
	UseValidLifetime                 bool                        `json:"use_valid_lifetime,omitempty"`
	UseZoneAssociations              bool                        `json:"use_zone_associations,omitempty"`
	ValidLifetime                    uint32                      `json:"valid_lifetime,omitempty"`
	Vlans                            []*Vlanlink                 `json:"vlans,omitempty"`
	ZoneAssociations                 []*Zoneassociation          `json:"zone_associations,omitempty"`
}

func (Ipv6Network) ObjectType() string {
	return "ipv6network"
}

func (obj Ipv6Network) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "network", "network_view"}
	}
	return obj.returnFields
}

// Ipv6range represents Infoblox object ipv6range
type Ipv6range struct {
	IBBase                           `json:"-"`
	Ref                              string                      `json:"_ref,omitempty"`
	AddressType                      string                      `json:"address_type,omitempty"`
	CloudInfo                        *GridCloudapiInfo           `json:"cloud_info,omitempty"`
	Comment                          string                      `json:"comment,omitempty"`
	Disable                          bool                        `json:"disable,omitempty"`
	DiscoverNowStatus                string                      `json:"discover_now_status,omitempty"`
	DiscoveryBasicPollSettings       *DiscoveryBasicpollsettings `json:"discovery_basic_poll_settings,omitempty"`
	DiscoveryBlackoutSetting         *PropertiesBlackoutsetting  `json:"discovery_blackout_setting,omitempty"`
	DiscoveryMember                  string                      `json:"discovery_member,omitempty"`
	EnableDiscovery                  bool                        `json:"enable_discovery,omitempty"`
	EnableImmediateDiscovery         bool                        `json:"enable_immediate_discovery,omitempty"`
	EndAddr                          string                      `json:"end_addr,omitempty"`
	EndpointSources                  []*CiscoiseEndpoint         `json:"endpoint_sources,omitempty"`
	Exclude                          []*Exclusionrange           `json:"exclude,omitempty"`
	Ea                               EA                          `json:"extattrs,omitempty"`
	Ipv6EndPrefix                    string                      `json:"ipv6_end_prefix,omitempty"`
	Ipv6PrefixBits                   uint32                      `json:"ipv6_prefix_bits,omitempty"`
	Ipv6StartPrefix                  string                      `json:"ipv6_start_prefix,omitempty"`
	Member                           *Dhcpmember                 `json:"member,omitempty"`
	Name                             string                      `json:"name,omitempty"`
	Network                          string                      `json:"network,omitempty"`
	NetworkView                      string                      `json:"network_view,omitempty"`
	NextAvailableIp                  *Nextavailableip6           `json:"next_available_ip,omitempty"`
	PortControlBlackoutSetting       *PropertiesBlackoutsetting  `json:"port_control_blackout_setting,omitempty"`
	RecycleLeases                    bool                        `json:"recycle_leases,omitempty"`
	RestartIfNeeded                  bool                        `json:"restart_if_needed,omitempty"`
	SamePortControlDiscoveryBlackout bool                        `json:"same_port_control_discovery_blackout,omitempty"`
	ServerAssociationType            string                      `json:"server_association_type,omitempty"`
	StartAddr                        string                      `json:"start_addr,omitempty"`
	SubscribeSettings                *CiscoiseSubscribesetting   `json:"subscribe_settings,omitempty"`
	Template                         string                      `json:"template,omitempty"`
	UseBlackoutSetting               bool                        `json:"use_blackout_setting,omitempty"`
	UseDiscoveryBasicPollingSettings bool                        `json:"use_discovery_basic_polling_settings,omitempty"`
	UseEnableDiscovery               bool                        `json:"use_enable_discovery,omitempty"`
	UseRecycleLeases                 bool                        `json:"use_recycle_leases,omitempty"`
	UseSubscribeSettings             bool                        `json:"use_subscribe_settings,omitempty"`
}

func (Ipv6range) ObjectType() string {
	return "ipv6range"
}

func (obj Ipv6range) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "end_addr", "network", "network_view", "start_addr"}
	}
	return obj.returnFields
}

// Ipv6NetworkContainer represents Infoblox object ipv6networkcontainer
type Ipv6NetworkContainer struct {
	IBBase                           `json:"-"`
	Ref                              string                      `json:"_ref,omitempty"`
	AutoCreateReversezone            bool                        `json:"auto_create_reversezone,omitempty"`
	CloudInfo                        *GridCloudapiInfo           `json:"cloud_info,omitempty"`
	Comment                          string                      `json:"comment,omitempty"`
	DdnsDomainname                   string                      `json:"ddns_domainname,omitempty"`
	DdnsEnableOptionFqdn             bool                        `json:"ddns_enable_option_fqdn,omitempty"`
	DdnsGenerateHostname             bool                        `json:"ddns_generate_hostname,omitempty"`
	DdnsServerAlwaysUpdates          bool                        `json:"ddns_server_always_updates,omitempty"`
	DdnsTtl                          uint32                      `json:"ddns_ttl,omitempty"`
	DeleteReason                     string                      `json:"delete_reason,omitempty"`
	DiscoverNowStatus                string                      `json:"discover_now_status,omitempty"`
	DiscoveryBasicPollSettings       *DiscoveryBasicpollsettings `json:"discovery_basic_poll_settings,omitempty"`
	DiscoveryBlackoutSetting         *PropertiesBlackoutsetting  `json:"discovery_blackout_setting,omitempty"`
	DiscoveryEngineType              string                      `json:"discovery_engine_type,omitempty"`
	DiscoveryMember                  string                      `json:"discovery_member,omitempty"`
	DomainNameServers                []string                    `json:"domain_name_servers,omitempty"`
	EnableDdns                       bool                        `json:"enable_ddns,omitempty"`
	EnableDiscovery                  bool                        `json:"enable_discovery,omitempty"`
	EnableImmediateDiscovery         bool                        `json:"enable_immediate_discovery,omitempty"`
	EndpointSources                  []*CiscoiseEndpoint         `json:"endpoint_sources,omitempty"`
	Ea                               EA                          `json:"extattrs,omitempty"`
	LastRirRegistrationUpdateSent    *time.Time                  `json:"last_rir_registration_update_sent,omitempty"`
	LastRirRegistrationUpdateStatus  string                      `json:"last_rir_registration_update_status,omitempty"`
	MgmPrivate                       bool                        `json:"mgm_private,omitempty"`
	MgmPrivateOverridable            bool                        `json:"mgm_private_overridable,omitempty"`
	MsAdUserData                     *MsserverAduserData         `json:"ms_ad_user_data,omitempty"`
	Network                          string                      `json:"network,omitempty"`
	NetworkContainer                 string                      `json:"network_container,omitempty"`
	NetworkView                      string                      `json:"network_view,omitempty"`
	NextAvailableNetwork             *Nextavailablenet6          `json:"next_available_network,omitempty"`
	Options                          []*Dhcpoption               `json:"options,omitempty"`
	PortControlBlackoutSetting       *PropertiesBlackoutsetting  `json:"port_control_blackout_setting,omitempty"`
	PreferredLifetime                uint32                      `json:"preferred_lifetime,omitempty"`
	RemoveSubnets                    bool                        `json:"remove_subnets,omitempty"`
	RestartIfNeeded                  bool                        `json:"restart_if_needed,omitempty"`
	Rir                              string                      `json:"rir,omitempty"`
	RirOrganization                  string                      `json:"rir_organization,omitempty"`
	RirRegistrationAction            string                      `json:"rir_registration_action,omitempty"`
	RirRegistrationStatus            string                      `json:"rir_registration_status,omitempty"`
	SamePortControlDiscoveryBlackout bool                        `json:"same_port_control_discovery_blackout,omitempty"`
	SendRirRequest                   bool                        `json:"send_rir_request,omitempty"`
	SubscribeSettings                *CiscoiseSubscribesetting   `json:"subscribe_settings,omitempty"`
	Unmanaged                        bool                        `json:"unmanaged,omitempty"`
	UpdateDnsOnLeaseRenewal          bool                        `json:"update_dns_on_lease_renewal,omitempty"`
	UseBlackoutSetting               bool                        `json:"use_blackout_setting,omitempty"`
	UseDdnsDomainname                bool                        `json:"use_ddns_domainname,omitempty"`
	UseDdnsEnableOptionFqdn          bool                        `json:"use_ddns_enable_option_fqdn,omitempty"`
	UseDdnsGenerateHostname          bool                        `json:"use_ddns_generate_hostname,omitempty"`
	UseDdnsTtl                       bool                        `json:"use_ddns_ttl,omitempty"`
	UseDiscoveryBasicPollingSettings bool                        `json:"use_discovery_basic_polling_settings,omitempty"`
	UseDomainNameServers             bool                        `json:"use_domain_name_servers,omitempty"`
	UseEnableDdns                    bool                        `json:"use_enable_ddns,omitempty"`
	UseEnableDiscovery               bool                        `json:"use_enable_discovery,omitempty"`
	UseMgmPrivate                    bool                        `json:"use_mgm_private,omitempty"`
	UseOptions                       bool                        `json:"use_options,omitempty"`
	UsePreferredLifetime             bool                        `json:"use_preferred_lifetime,omitempty"`
	UseSubscribeSettings             bool                        `json:"use_subscribe_settings,omitempty"`
	UseUpdateDnsOnLeaseRenewal       bool                        `json:"use_update_dns_on_lease_renewal,omitempty"`
	UseValidLifetime                 bool                        `json:"use_valid_lifetime,omitempty"`
	UseZoneAssociations              bool                        `json:"use_zone_associations,omitempty"`
	Utilization                      uint32                      `json:"utilization,omitempty"`
	ValidLifetime                    uint32                      `json:"valid_lifetime,omitempty"`
	ZoneAssociations                 []*Zoneassociation          `json:"zone_associations,omitempty"`
}

func (Ipv6NetworkContainer) ObjectType() string {
	return "ipv6networkcontainer"
}

func (obj Ipv6NetworkContainer) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "network", "network_view"}
	}
	return obj.returnFields
}

// LdapAuthService represents Infoblox object ldap_auth_service
type LdapAuthService struct {
	IBBase                      `json:"-"`
	Ref                         string                   `json:"_ref,omitempty"`
	CheckLdapServerSettings     *Checkldapserversettings `json:"check_ldap_server_settings,omitempty"`
	Comment                     string                   `json:"comment,omitempty"`
	Disable                     bool                     `json:"disable,omitempty"`
	EaMapping                   []*LdapEamapping         `json:"ea_mapping,omitempty"`
	LdapGroupAttribute          string                   `json:"ldap_group_attribute,omitempty"`
	LdapGroupAuthenticationType string                   `json:"ldap_group_authentication_type,omitempty"`
	LdapUserAttribute           string                   `json:"ldap_user_attribute,omitempty"`
	Mode                        string                   `json:"mode,omitempty"`
	Name                        string                   `json:"name,omitempty"`
	RecoveryInterval            uint32                   `json:"recovery_interval,omitempty"`
	Retries                     uint32                   `json:"retries,omitempty"`
	SearchScope                 string                   `json:"search_scope,omitempty"`
	Servers                     []*LdapServer            `json:"servers,omitempty"`
	Timeout                     uint32                   `json:"timeout,omitempty"`
}

func (LdapAuthService) ObjectType() string {
	return "ldap_auth_service"
}

func (obj LdapAuthService) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "disable", "ldap_user_attribute", "mode", "name"}
	}
	return obj.returnFields
}

// Lease represents Infoblox object lease
type Lease struct {
	IBBase                `json:"-"`
	Ref                   string              `json:"_ref,omitempty"`
	Address               string              `json:"address,omitempty"`
	BillingClass          string              `json:"billing_class,omitempty"`
	BindingState          string              `json:"binding_state,omitempty"`
	ClientHostname        string              `json:"client_hostname,omitempty"`
	Cltt                  *time.Time          `json:"cltt,omitempty"`
	DiscoveredData        *Discoverydata      `json:"discovered_data,omitempty"`
	Ends                  *time.Time          `json:"ends,omitempty"`
	Fingerprint           string              `json:"fingerprint,omitempty"`
	Hardware              string              `json:"hardware,omitempty"`
	Ipv6Duid              string              `json:"ipv6_duid,omitempty"`
	Ipv6Iaid              string              `json:"ipv6_iaid,omitempty"`
	Ipv6PreferredLifetime uint32              `json:"ipv6_preferred_lifetime,omitempty"`
	Ipv6PrefixBits        uint32              `json:"ipv6_prefix_bits,omitempty"`
	IsInvalidMac          bool                `json:"is_invalid_mac,omitempty"`
	MsAdUserData          *MsserverAduserData `json:"ms_ad_user_data,omitempty"`
	Network               string              `json:"network,omitempty"`
	NetworkView           string              `json:"network_view,omitempty"`
	NeverEnds             bool                `json:"never_ends,omitempty"`
	NeverStarts           bool                `json:"never_starts,omitempty"`
	NextBindingState      string              `json:"next_binding_state,omitempty"`
	OnCommit              string              `json:"on_commit,omitempty"`
	OnExpiry              string              `json:"on_expiry,omitempty"`
	OnRelease             string              `json:"on_release,omitempty"`
	Option                string              `json:"option,omitempty"`
	Protocol              string              `json:"protocol,omitempty"`
	RemoteId              string              `json:"remote_id,omitempty"`
	ServedBy              string              `json:"served_by,omitempty"`
	ServerHostName        string              `json:"server_host_name,omitempty"`
	Starts                *time.Time          `json:"starts,omitempty"`
	Tsfp                  *time.Time          `json:"tsfp,omitempty"`
	Tstp                  *time.Time          `json:"tstp,omitempty"`
	Uid                   string              `json:"uid,omitempty"`
	Username              string              `json:"username,omitempty"`
	Variable              string              `json:"variable,omitempty"`
}

func (Lease) ObjectType() string {
	return "lease"
}

func (obj Lease) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"address", "network_view"}
	}
	return obj.returnFields
}

// LicenseGridwide represents Infoblox object license:gridwide
type LicenseGridwide struct {
	IBBase           `json:"-"`
	Ref              string     `json:"_ref,omitempty"`
	ExpirationStatus string     `json:"expiration_status,omitempty"`
	ExpiryDate       *time.Time `json:"expiry_date,omitempty"`
	Key              string     `json:"key,omitempty"`
	Limit            string     `json:"limit,omitempty"`
	LimitContext     string     `json:"limit_context,omitempty"`
	Type             string     `json:"type,omitempty"`
}

func (LicenseGridwide) ObjectType() string {
	return "license:gridwide"
}

func (obj LicenseGridwide) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"type"}
	}
	return obj.returnFields
}

// LocaluserAuthservice represents Infoblox object localuser:authservice
type LocaluserAuthservice struct {
	IBBase   `json:"-"`
	Ref      string `json:"_ref,omitempty"`
	Comment  string `json:"comment,omitempty"`
	Disabled bool   `json:"disabled,omitempty"`
	Name     string `json:"name,omitempty"`
}

func (LocaluserAuthservice) ObjectType() string {
	return "localuser:authservice"
}

func (obj LocaluserAuthservice) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "disabled", "name"}
	}
	return obj.returnFields
}

// Mastergrid represents Infoblox object mastergrid
type Mastergrid struct {
	IBBase              `json:"-"`
	Ref                 string     `json:"_ref,omitempty"`
	Address             string     `json:"address,omitempty"`
	ConnectionDisabled  bool       `json:"connection_disabled,omitempty"`
	ConnectionTimestamp *time.Time `json:"connection_timestamp,omitempty"`
	Detached            bool       `json:"detached,omitempty"`
	Enable              bool       `json:"enable,omitempty"`
	Joined              bool       `json:"joined,omitempty"`
	LastEvent           string     `json:"last_event,omitempty"`
	LastEventDetails    string     `json:"last_event_details,omitempty"`
	LastSyncTimestamp   *time.Time `json:"last_sync_timestamp,omitempty"`
	Port                uint32     `json:"port,omitempty"`
	Status              string     `json:"status,omitempty"`
	UseMgmtPort         bool       `json:"use_mgmt_port,omitempty"`
}

func (Mastergrid) ObjectType() string {
	return "mastergrid"
}

func (obj Mastergrid) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"address", "enable", "port"}
	}
	return obj.returnFields
}

// Macfilteraddress represents Infoblox object macfilteraddress
type Macfilteraddress struct {
	IBBase              `json:"-"`
	Ref                 string     `json:"_ref,omitempty"`
	AuthenticationTime  *time.Time `json:"authentication_time,omitempty"`
	Comment             string     `json:"comment,omitempty"`
	ExpirationTime      *time.Time `json:"expiration_time,omitempty"`
	Ea                  EA         `json:"extattrs,omitempty"`
	Filter              string     `json:"filter,omitempty"`
	Fingerprint         string     `json:"fingerprint,omitempty"`
	GuestCustomField1   string     `json:"guest_custom_field1,omitempty"`
	GuestCustomField2   string     `json:"guest_custom_field2,omitempty"`
	GuestCustomField3   string     `json:"guest_custom_field3,omitempty"`
	GuestCustomField4   string     `json:"guest_custom_field4,omitempty"`
	GuestEmail          string     `json:"guest_email,omitempty"`
	GuestFirstName      string     `json:"guest_first_name,omitempty"`
	GuestLastName       string     `json:"guest_last_name,omitempty"`
	GuestMiddleName     string     `json:"guest_middle_name,omitempty"`
	GuestPhone          string     `json:"guest_phone,omitempty"`
	IsRegisteredUser    bool       `json:"is_registered_user,omitempty"`
	Mac                 string     `json:"mac,omitempty"`
	NeverExpires        bool       `json:"never_expires,omitempty"`
	ReservedForInfoblox string     `json:"reserved_for_infoblox,omitempty"`
	Username            string     `json:"username,omitempty"`
}

func (Macfilteraddress) ObjectType() string {
	return "macfilteraddress"
}

func (obj Macfilteraddress) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"authentication_time", "comment", "expiration_time", "filter", "guest_custom_field1", "guest_custom_field2", "guest_custom_field3", "guest_custom_field4", "guest_email", "guest_first_name", "guest_last_name", "guest_middle_name", "guest_phone", "is_registered_user", "mac", "never_expires", "reserved_for_infoblox", "username"}
	}
	return obj.returnFields
}

// MemberFiledistribution represents Infoblox object member:filedistribution
type MemberFiledistribution struct {
	IBBase            `json:"-"`
	Ref               string       `json:"_ref,omitempty"`
	AllowUploads      bool         `json:"allow_uploads,omitempty"`
	Comment           string       `json:"comment,omitempty"`
	EnableFtp         bool         `json:"enable_ftp,omitempty"`
	EnableFtpFilelist bool         `json:"enable_ftp_filelist,omitempty"`
	EnableFtpPassive  bool         `json:"enable_ftp_passive,omitempty"`
	EnableHttp        bool         `json:"enable_http,omitempty"`
	EnableHttpAcl     bool         `json:"enable_http_acl,omitempty"`
	EnableTftp        bool         `json:"enable_tftp,omitempty"`
	FtpAcls           []*Addressac `json:"ftp_acls,omitempty"`
	FtpPort           uint32       `json:"ftp_port,omitempty"`
	FtpStatus         string       `json:"ftp_status,omitempty"`
	HostName          string       `json:"host_name,omitempty"`
	HttpAcls          []*Addressac `json:"http_acls,omitempty"`
	HttpStatus        string       `json:"http_status,omitempty"`
	Ipv4Address       string       `json:"ipv4_address,omitempty"`
	Ipv6Address       string       `json:"ipv6_address,omitempty"`
	Status            string       `json:"status,omitempty"`
	TftpAcls          []*Addressac `json:"tftp_acls,omitempty"`
	TftpPort          uint32       `json:"tftp_port,omitempty"`
	TftpStatus        string       `json:"tftp_status,omitempty"`
	UseAllowUploads   bool         `json:"use_allow_uploads,omitempty"`
}

func (MemberFiledistribution) ObjectType() string {
	return "member:filedistribution"
}

func (obj MemberFiledistribution) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"host_name", "ipv4_address", "ipv6_address", "status"}
	}
	return obj.returnFields
}

// MemberLicense represents Infoblox object member:license
type MemberLicense struct {
	IBBase           `json:"-"`
	Ref              string     `json:"_ref,omitempty"`
	ExpirationStatus string     `json:"expiration_status,omitempty"`
	ExpiryDate       *time.Time `json:"expiry_date,omitempty"`
	Hwid             string     `json:"hwid,omitempty"`
	Key              string     `json:"key,omitempty"`
	Kind             string     `json:"kind,omitempty"`
	Limit            string     `json:"limit,omitempty"`
	LimitContext     string     `json:"limit_context,omitempty"`
	Type             string     `json:"type,omitempty"`
}

func (MemberLicense) ObjectType() string {
	return "member:license"
}

func (obj MemberLicense) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"type"}
	}
	return obj.returnFields
}

// MemberThreatanalytics represents Infoblox object member:threatanalytics
type MemberThreatanalytics struct {
	IBBase        `json:"-"`
	Ref           string `json:"_ref,omitempty"`
	Comment       string `json:"comment,omitempty"`
	EnableService bool   `json:"enable_service,omitempty"`
	HostName      string `json:"host_name,omitempty"`
	Ipv4Address   string `json:"ipv4_address,omitempty"`
	Ipv6Address   string `json:"ipv6_address,omitempty"`
	Status        string `json:"status,omitempty"`
}

func (MemberThreatanalytics) ObjectType() string {
	return "member:threatanalytics"
}

func (obj MemberThreatanalytics) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"host_name", "ipv4_address", "ipv6_address", "status"}
	}
	return obj.returnFields
}

// MemberParentalcontrol represents Infoblox object member:parentalcontrol
type MemberParentalcontrol struct {
	IBBase        `json:"-"`
	Ref           string `json:"_ref,omitempty"`
	EnableService bool   `json:"enable_service,omitempty"`
	Name          string `json:"name,omitempty"`
}

func (MemberParentalcontrol) ObjectType() string {
	return "member:parentalcontrol"
}

func (obj MemberParentalcontrol) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"enable_service", "name"}
	}
	return obj.returnFields
}

// MemberThreatprotection represents Infoblox object member:threatprotection
type MemberThreatprotection struct {
	IBBase                                   `json:"-"`
	Ref                                      string                     `json:"_ref,omitempty"`
	Comment                                  string                     `json:"comment,omitempty"`
	CurrentRuleset                           string                     `json:"current_ruleset,omitempty"`
	DisableMultipleDnsTcpRequest             bool                       `json:"disable_multiple_dns_tcp_request,omitempty"`
	EnableAccelRespBeforeThreatProtection    bool                       `json:"enable_accel_resp_before_threat_protection,omitempty"`
	EnableNatRules                           bool                       `json:"enable_nat_rules,omitempty"`
	EnableService                            bool                       `json:"enable_service,omitempty"`
	EventsPerSecondPerRule                   uint32                     `json:"events_per_second_per_rule,omitempty"`
	HardwareModel                            string                     `json:"hardware_model,omitempty"`
	HardwareType                             string                     `json:"hardware_type,omitempty"`
	HostName                                 string                     `json:"host_name,omitempty"`
	Ipv4address                              string                     `json:"ipv4address,omitempty"`
	Ipv6address                              string                     `json:"ipv6address,omitempty"`
	NatRules                                 []*ThreatprotectionNatrule `json:"nat_rules,omitempty"`
	OutboundSettings                         *SettingAtpoutbound        `json:"outbound_settings,omitempty"`
	Profile                                  string                     `json:"profile,omitempty"`
	UseCurrentRuleset                        bool                       `json:"use_current_ruleset,omitempty"`
	UseDisableMultipleDnsTcpRequest          bool                       `json:"use_disable_multiple_dns_tcp_request,omitempty"`
	UseEnableAccelRespBeforeThreatProtection bool                       `json:"use_enable_accel_resp_before_threat_protection,omitempty"`
	UseEnableNatRules                        bool                       `json:"use_enable_nat_rules,omitempty"`
	UseEventsPerSecondPerRule                bool                       `json:"use_events_per_second_per_rule,omitempty"`
	UseOutboundSettings                      bool                       `json:"use_outbound_settings,omitempty"`
}

func (MemberThreatprotection) ObjectType() string {
	return "member:threatprotection"
}

func (obj MemberThreatprotection) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{}
	}
	return obj.returnFields
}

// MemberDhcpproperties represents Infoblox object member:dhcpproperties
type MemberDhcpproperties struct {
	IBBase                                `json:"-"`
	Ref                                   string                   `json:"_ref,omitempty"`
	AuthServerGroup                       string                   `json:"auth_server_group,omitempty"`
	AuthnCaptivePortal                    string                   `json:"authn_captive_portal,omitempty"`
	AuthnCaptivePortalAuthenticatedFilter string                   `json:"authn_captive_portal_authenticated_filter,omitempty"`
	AuthnCaptivePortalEnabled             bool                     `json:"authn_captive_portal_enabled,omitempty"`
	AuthnCaptivePortalGuestFilter         string                   `json:"authn_captive_portal_guest_filter,omitempty"`
	AuthnServerGroupEnabled               bool                     `json:"authn_server_group_enabled,omitempty"`
	Authority                             bool                     `json:"authority,omitempty"`
	Bootfile                              string                   `json:"bootfile,omitempty"`
	Bootserver                            string                   `json:"bootserver,omitempty"`
	ClearNacAuthCache                     *Clearnacauthcacheparams `json:"clear_nac_auth_cache,omitempty"`
	DdnsDomainname                        string                   `json:"ddns_domainname,omitempty"`
	DdnsGenerateHostname                  bool                     `json:"ddns_generate_hostname,omitempty"`
	DdnsRetryInterval                     uint32                   `json:"ddns_retry_interval,omitempty"`
	DdnsServerAlwaysUpdates               bool                     `json:"ddns_server_always_updates,omitempty"`
	DdnsTtl                               uint32                   `json:"ddns_ttl,omitempty"`
	DdnsUpdateFixedAddresses              bool                     `json:"ddns_update_fixed_addresses,omitempty"`
	DdnsUseOption81                       bool                     `json:"ddns_use_option81,omitempty"`
	DdnsZonePrimaries                     []*Dhcpddns              `json:"ddns_zone_primaries,omitempty"`
	DenyBootp                             bool                     `json:"deny_bootp,omitempty"`
	DhcpUtilization                       uint32                   `json:"dhcp_utilization,omitempty"`
	DhcpUtilizationStatus                 string                   `json:"dhcp_utilization_status,omitempty"`
	DnsUpdateStyle                        string                   `json:"dns_update_style,omitempty"`
	DynamicHosts                          uint32                   `json:"dynamic_hosts,omitempty"`
	EmailList                             []string                 `json:"email_list,omitempty"`
	EnableDdns                            bool                     `json:"enable_ddns,omitempty"`
	EnableDhcp                            bool                     `json:"enable_dhcp,omitempty"`
	EnableDhcpOnIpv6Lan2                  bool                     `json:"enable_dhcp_on_ipv6_lan2,omitempty"`
	EnableDhcpOnLan2                      bool                     `json:"enable_dhcp_on_lan2,omitempty"`
	EnableDhcpThresholds                  bool                     `json:"enable_dhcp_thresholds,omitempty"`
	EnableDhcpv6Service                   bool                     `json:"enable_dhcpv6_service,omitempty"`
	EnableEmailWarnings                   bool                     `json:"enable_email_warnings,omitempty"`
	EnableFingerprint                     bool                     `json:"enable_fingerprint,omitempty"`
	EnableGssTsig                         bool                     `json:"enable_gss_tsig,omitempty"`
	EnableHostnameRewrite                 bool                     `json:"enable_hostname_rewrite,omitempty"`
	EnableLeasequery                      bool                     `json:"enable_leasequery,omitempty"`
	EnableSnmpWarnings                    bool                     `json:"enable_snmp_warnings,omitempty"`
	Ea                                    EA                       `json:"extattrs,omitempty"`
	GssTsigKeys                           []*Kerberoskey           `json:"gss_tsig_keys,omitempty"`
	HighWaterMark                         uint32                   `json:"high_water_mark,omitempty"`
	HighWaterMarkReset                    uint32                   `json:"high_water_mark_reset,omitempty"`
	HostName                              string                   `json:"host_name,omitempty"`
	HostnameRewritePolicy                 string                   `json:"hostname_rewrite_policy,omitempty"`
	IgnoreDhcpOptionListRequest           bool                     `json:"ignore_dhcp_option_list_request,omitempty"`
	IgnoreId                              string                   `json:"ignore_id,omitempty"`
	IgnoreMacAddresses                    []string                 `json:"ignore_mac_addresses,omitempty"`
	ImmediateFaConfiguration              bool                     `json:"immediate_fa_configuration,omitempty"`
	Ipv4Addr                              string                   `json:"ipv4addr,omitempty"`
	Ipv6DdnsDomainname                    string                   `json:"ipv6_ddns_domainname,omitempty"`
	Ipv6DdnsEnableOptionFqdn              bool                     `json:"ipv6_ddns_enable_option_fqdn,omitempty"`
	Ipv6DdnsHostname                      string                   `json:"ipv6_ddns_hostname,omitempty"`
	Ipv6DdnsServerAlwaysUpdates           bool                     `json:"ipv6_ddns_server_always_updates,omitempty"`
	Ipv6DdnsTtl                           uint32                   `json:"ipv6_ddns_ttl,omitempty"`
	Ipv6DnsUpdateStyle                    string                   `json:"ipv6_dns_update_style,omitempty"`
	Ipv6DomainName                        string                   `json:"ipv6_domain_name,omitempty"`
	Ipv6DomainNameServers                 []string                 `json:"ipv6_domain_name_servers,omitempty"`
	Ipv6EnableDdns                        bool                     `json:"ipv6_enable_ddns,omitempty"`
	Ipv6EnableGssTsig                     bool                     `json:"ipv6_enable_gss_tsig,omitempty"`
	Ipv6EnableLeaseScavenging             bool                     `json:"ipv6_enable_lease_scavenging,omitempty"`
	Ipv6EnableRetryUpdates                bool                     `json:"ipv6_enable_retry_updates,omitempty"`
	Ipv6GenerateHostname                  bool                     `json:"ipv6_generate_hostname,omitempty"`
	Ipv6GssTsigKeys                       []*Kerberoskey           `json:"ipv6_gss_tsig_keys,omitempty"`
	Ipv6KdcServer                         string                   `json:"ipv6_kdc_server,omitempty"`
	Ipv6LeaseScavengingTime               uint32                   `json:"ipv6_lease_scavenging_time,omitempty"`
	Ipv6MicrosoftCodePage                 string                   `json:"ipv6_microsoft_code_page,omitempty"`
	Ipv6Options                           []*Dhcpoption            `json:"ipv6_options,omitempty"`
	Ipv6RecycleLeases                     bool                     `json:"ipv6_recycle_leases,omitempty"`
	Ipv6RememberExpiredClientAssociation  bool                     `json:"ipv6_remember_expired_client_association,omitempty"`
	Ipv6RetryUpdatesInterval              uint32                   `json:"ipv6_retry_updates_interval,omitempty"`
	Ipv6ServerDuid                        string                   `json:"ipv6_server_duid,omitempty"`
	Ipv6UpdateDnsOnLeaseRenewal           bool                     `json:"ipv6_update_dns_on_lease_renewal,omitempty"`
	Ipv6Addr                              string                   `json:"ipv6addr,omitempty"`
	KdcServer                             string                   `json:"kdc_server,omitempty"`
	LeasePerClientSettings                string                   `json:"lease_per_client_settings,omitempty"`
	LeaseScavengeTime                     int                      `json:"lease_scavenge_time,omitempty"`
	LogLeaseEvents                        bool                     `json:"log_lease_events,omitempty"`
	LogicFilterRules                      []*Logicfilterrule       `json:"logic_filter_rules,omitempty"`
	LowWaterMark                          uint32                   `json:"low_water_mark,omitempty"`
	LowWaterMarkReset                     uint32                   `json:"low_water_mark_reset,omitempty"`
	MicrosoftCodePage                     string                   `json:"microsoft_code_page,omitempty"`
	Nextserver                            string                   `json:"nextserver,omitempty"`
	Option60MatchRules                    []*Option60matchrule     `json:"option60_match_rules,omitempty"`
	Options                               []*Dhcpoption            `json:"options,omitempty"`
	PingCount                             uint32                   `json:"ping_count,omitempty"`
	PingTimeout                           uint32                   `json:"ping_timeout,omitempty"`
	PreferredLifetime                     uint32                   `json:"preferred_lifetime,omitempty"`
	PrefixLengthMode                      string                   `json:"prefix_length_mode,omitempty"`
	PurgeIfmapData                        *Emptyparams             `json:"purge_ifmap_data,omitempty"`
	PxeLeaseTime                          uint32                   `json:"pxe_lease_time,omitempty"`
	RecycleLeases                         bool                     `json:"recycle_leases,omitempty"`
	RetryDdnsUpdates                      bool                     `json:"retry_ddns_updates,omitempty"`
	StaticHosts                           uint32                   `json:"static_hosts,omitempty"`
	SyslogFacility                        string                   `json:"syslog_facility,omitempty"`
	TotalHosts                            uint32                   `json:"total_hosts,omitempty"`
	UpdateDnsOnLeaseRenewal               bool                     `json:"update_dns_on_lease_renewal,omitempty"`
	UseAuthority                          bool                     `json:"use_authority,omitempty"`
	UseBootfile                           bool                     `json:"use_bootfile,omitempty"`
	UseBootserver                         bool                     `json:"use_bootserver,omitempty"`
	UseDdnsDomainname                     bool                     `json:"use_ddns_domainname,omitempty"`
	UseDdnsGenerateHostname               bool                     `json:"use_ddns_generate_hostname,omitempty"`
	UseDdnsTtl                            bool                     `json:"use_ddns_ttl,omitempty"`
	UseDdnsUpdateFixedAddresses           bool                     `json:"use_ddns_update_fixed_addresses,omitempty"`
	UseDdnsUseOption81                    bool                     `json:"use_ddns_use_option81,omitempty"`
	UseDenyBootp                          bool                     `json:"use_deny_bootp,omitempty"`
	UseDnsUpdateStyle                     bool                     `json:"use_dns_update_style,omitempty"`
	UseEmailList                          bool                     `json:"use_email_list,omitempty"`
	UseEnableDdns                         bool                     `json:"use_enable_ddns,omitempty"`
	UseEnableDhcpThresholds               bool                     `json:"use_enable_dhcp_thresholds,omitempty"`
	UseEnableFingerprint                  bool                     `json:"use_enable_fingerprint,omitempty"`
	UseEnableGssTsig                      bool                     `json:"use_enable_gss_tsig,omitempty"`
	UseEnableHostnameRewrite              bool                     `json:"use_enable_hostname_rewrite,omitempty"`
	UseEnableLeasequery                   bool                     `json:"use_enable_leasequery,omitempty"`
	UseEnableOneLeasePerClient            bool                     `json:"use_enable_one_lease_per_client,omitempty"`
	UseGssTsigKeys                        bool                     `json:"use_gss_tsig_keys,omitempty"`
	UseIgnoreDhcpOptionListRequest        bool                     `json:"use_ignore_dhcp_option_list_request,omitempty"`
	UseIgnoreId                           bool                     `json:"use_ignore_id,omitempty"`
	UseImmediateFaConfiguration           bool                     `json:"use_immediate_fa_configuration,omitempty"`
	UseIpv6DdnsDomainname                 bool                     `json:"use_ipv6_ddns_domainname,omitempty"`
	UseIpv6DdnsEnableOptionFqdn           bool                     `json:"use_ipv6_ddns_enable_option_fqdn,omitempty"`
	UseIpv6DdnsHostname                   bool                     `json:"use_ipv6_ddns_hostname,omitempty"`
	UseIpv6DdnsTtl                        bool                     `json:"use_ipv6_ddns_ttl,omitempty"`
	UseIpv6DnsUpdateStyle                 bool                     `json:"use_ipv6_dns_update_style,omitempty"`
	UseIpv6DomainName                     bool                     `json:"use_ipv6_domain_name,omitempty"`
	UseIpv6DomainNameServers              bool                     `json:"use_ipv6_domain_name_servers,omitempty"`
	UseIpv6EnableDdns                     bool                     `json:"use_ipv6_enable_ddns,omitempty"`
	UseIpv6EnableGssTsig                  bool                     `json:"use_ipv6_enable_gss_tsig,omitempty"`
	UseIpv6EnableRetryUpdates             bool                     `json:"use_ipv6_enable_retry_updates,omitempty"`
	UseIpv6GenerateHostname               bool                     `json:"use_ipv6_generate_hostname,omitempty"`
	UseIpv6GssTsigKeys                    bool                     `json:"use_ipv6_gss_tsig_keys,omitempty"`
	UseIpv6LeaseScavenging                bool                     `json:"use_ipv6_lease_scavenging,omitempty"`
	UseIpv6MicrosoftCodePage              bool                     `json:"use_ipv6_microsoft_code_page,omitempty"`
	UseIpv6Options                        bool                     `json:"use_ipv6_options,omitempty"`
	UseIpv6RecycleLeases                  bool                     `json:"use_ipv6_recycle_leases,omitempty"`
	UseIpv6UpdateDnsOnLeaseRenewal        bool                     `json:"use_ipv6_update_dns_on_lease_renewal,omitempty"`
	UseLeasePerClientSettings             bool                     `json:"use_lease_per_client_settings,omitempty"`
	UseLeaseScavengeTime                  bool                     `json:"use_lease_scavenge_time,omitempty"`
	UseLogLeaseEvents                     bool                     `json:"use_log_lease_events,omitempty"`
	UseLogicFilterRules                   bool                     `json:"use_logic_filter_rules,omitempty"`
	UseMicrosoftCodePage                  bool                     `json:"use_microsoft_code_page,omitempty"`
	UseNextserver                         bool                     `json:"use_nextserver,omitempty"`
	UseOptions                            bool                     `json:"use_options,omitempty"`
	UsePingCount                          bool                     `json:"use_ping_count,omitempty"`
	UsePingTimeout                        bool                     `json:"use_ping_timeout,omitempty"`
	UsePreferredLifetime                  bool                     `json:"use_preferred_lifetime,omitempty"`
	UsePrefixLengthMode                   bool                     `json:"use_prefix_length_mode,omitempty"`
	UsePxeLeaseTime                       bool                     `json:"use_pxe_lease_time,omitempty"`
	UseRecycleLeases                      bool                     `json:"use_recycle_leases,omitempty"`
	UseRetryDdnsUpdates                   bool                     `json:"use_retry_ddns_updates,omitempty"`
	UseSyslogFacility                     bool                     `json:"use_syslog_facility,omitempty"`
	UseUpdateDnsOnLeaseRenewal            bool                     `json:"use_update_dns_on_lease_renewal,omitempty"`
	UseValidLifetime                      bool                     `json:"use_valid_lifetime,omitempty"`
	ValidLifetime                         uint32                   `json:"valid_lifetime,omitempty"`
}

func (MemberDhcpproperties) ObjectType() string {
	return "member:dhcpproperties"
}

func (obj MemberDhcpproperties) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"host_name", "ipv4addr", "ipv6addr"}
	}
	return obj.returnFields
}

// Memberdfp represents Infoblox object memberdfp
type Memberdfp struct {
	IBBase          `json:"-"`
	Ref             string `json:"_ref,omitempty"`
	DfpForwardFirst bool   `json:"dfp_forward_first,omitempty"`
	Ea              EA     `json:"extattrs,omitempty"`
	HostName        string `json:"host_name,omitempty"`
	IsDfpOverride   bool   `json:"is_dfp_override,omitempty"`
}

func (Memberdfp) ObjectType() string {
	return "memberdfp"
}

func (obj Memberdfp) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{}
	}
	return obj.returnFields
}

// MsserverAdsitesDomain represents Infoblox object msserver:adsites:domain
type MsserverAdsitesDomain struct {
	IBBase           `json:"-"`
	Ref              string `json:"_ref,omitempty"`
	EaDefinition     string `json:"ea_definition,omitempty"`
	MsSyncMasterName string `json:"ms_sync_master_name,omitempty"`
	Name             string `json:"name,omitempty"`
	Netbios          string `json:"netbios,omitempty"`
	NetworkView      string `json:"network_view,omitempty"`
	ReadOnly         bool   `json:"read_only,omitempty"`
}

func (MsserverAdsitesDomain) ObjectType() string {
	return "msserver:adsites:domain"
}

func (obj MsserverAdsitesDomain) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name", "netbios", "network_view"}
	}
	return obj.returnFields
}

// MsserverAdsitesSite represents Infoblox object msserver:adsites:site
type MsserverAdsitesSite struct {
	IBBase      `json:"-"`
	Ref         string         `json:"_ref,omitempty"`
	Domain      string         `json:"domain,omitempty"`
	MoveSubnets *Movesubnets   `json:"move_subnets,omitempty"`
	Name        string         `json:"name,omitempty"`
	Networks    []*Ipv4Network `json:"networks,omitempty"`
}

func (MsserverAdsitesSite) ObjectType() string {
	return "msserver:adsites:site"
}

func (obj MsserverAdsitesSite) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"domain", "name"}
	}
	return obj.returnFields
}

// MemberDns represents Infoblox object member:dns
type MemberDns struct {
	IBBase                           `json:"-"`
	Ref                              string                        `json:"_ref,omitempty"`
	AddClientIpMacOptions            bool                          `json:"add_client_ip_mac_options,omitempty"`
	AdditionalIpList                 []string                      `json:"additional_ip_list,omitempty"`
	AdditionalIpListStruct           []*MemberDnsip                `json:"additional_ip_list_struct,omitempty"`
	AllowGssTsigZoneUpdates          bool                          `json:"allow_gss_tsig_zone_updates,omitempty"`
	AllowQuery                       []*Addressac                  `json:"allow_query,omitempty"`
	AllowRecursiveQuery              bool                          `json:"allow_recursive_query,omitempty"`
	AllowTransfer                    []*Addressac                  `json:"allow_transfer,omitempty"`
	AllowUpdate                      []*Addressac                  `json:"allow_update,omitempty"`
	AnonymizeResponseLogging         bool                          `json:"anonymize_response_logging,omitempty"`
	AtcFwdEnable                     bool                          `json:"atc_fwd_enable,omitempty"`
	AttackMitigation                 *GridAttackmitigation         `json:"attack_mitigation,omitempty"`
	AutoBlackhole                    *GridAutoblackhole            `json:"auto_blackhole,omitempty"`
	AutoCreateAAndPtrForLan2         bool                          `json:"auto_create_a_and_ptr_for_lan2,omitempty"`
	AutoCreateAaaaAndIpv6ptrForLan2  bool                          `json:"auto_create_aaaa_and_ipv6ptr_for_lan2,omitempty"`
	AutoSortViews                    bool                          `json:"auto_sort_views,omitempty"`
	BindCheckNamesPolicy             string                        `json:"bind_check_names_policy,omitempty"`
	BindHostnameDirective            string                        `json:"bind_hostname_directive,omitempty"`
	BindHostnameDirectiveFqdn        string                        `json:"bind_hostname_directive_fqdn,omitempty"`
	BlackholeList                    []*Addressac                  `json:"blackhole_list,omitempty"`
	BlacklistAction                  string                        `json:"blacklist_action,omitempty"`
	BlacklistLogQuery                bool                          `json:"blacklist_log_query,omitempty"`
	BlacklistRedirectAddresses       []string                      `json:"blacklist_redirect_addresses,omitempty"`
	BlacklistRedirectTtl             uint32                        `json:"blacklist_redirect_ttl,omitempty"`
	BlacklistRulesets                []string                      `json:"blacklist_rulesets,omitempty"`
	CaptureDnsQueriesOnAllDomains    bool                          `json:"capture_dns_queries_on_all_domains,omitempty"`
	CheckNamesForDdnsAndZoneTransfer bool                          `json:"check_names_for_ddns_and_zone_transfer,omitempty"`
	ClearDnsCache                    *Cleardnscache                `json:"clear_dns_cache,omitempty"`
	CopyClientIpMacOptions           bool                          `json:"copy_client_ip_mac_options,omitempty"`
	CopyXferToNotify                 bool                          `json:"copy_xfer_to_notify,omitempty"`
	CustomRootNameServers            []NameServer                  `json:"custom_root_name_servers,omitempty"`
	DisableEdns                      bool                          `json:"disable_edns,omitempty"`
	Dns64Groups                      []string                      `json:"dns64_groups,omitempty"`
	DnsCacheAccelerationStatus       string                        `json:"dns_cache_acceleration_status,omitempty"`
	DnsCacheAccelerationTtl          uint32                        `json:"dns_cache_acceleration_ttl,omitempty"`
	DnsHealthCheckAnycastControl     bool                          `json:"dns_health_check_anycast_control,omitempty"`
	DnsHealthCheckDomainList         []string                      `json:"dns_health_check_domain_list,omitempty"`
	DnsHealthCheckInterval           uint32                        `json:"dns_health_check_interval,omitempty"`
	DnsHealthCheckRecursionFlag      bool                          `json:"dns_health_check_recursion_flag,omitempty"`
	DnsHealthCheckRetries            uint32                        `json:"dns_health_check_retries,omitempty"`
	DnsHealthCheckTimeout            uint32                        `json:"dns_health_check_timeout,omitempty"`
	DnsNotifyTransferSource          string                        `json:"dns_notify_transfer_source,omitempty"`
	DnsNotifyTransferSourceAddress   string                        `json:"dns_notify_transfer_source_address,omitempty"`
	DnsQueryCaptureFileTimeLimit     uint32                        `json:"dns_query_capture_file_time_limit,omitempty"`
	DnsQuerySourceAddress            string                        `json:"dns_query_source_address,omitempty"`
	DnsQuerySourceInterface          string                        `json:"dns_query_source_interface,omitempty"`
	DnsViewAddressSettings           []*SettingViewaddress         `json:"dns_view_address_settings,omitempty"`
	DnssecBlacklistEnabled           bool                          `json:"dnssec_blacklist_enabled,omitempty"`
	DnssecDns64Enabled               bool                          `json:"dnssec_dns64_enabled,omitempty"`
	DnssecEnabled                    bool                          `json:"dnssec_enabled,omitempty"`
	DnssecExpiredSignaturesEnabled   bool                          `json:"dnssec_expired_signatures_enabled,omitempty"`
	DnssecNegativeTrustAnchors       []string                      `json:"dnssec_negative_trust_anchors,omitempty"`
	DnssecNxdomainEnabled            bool                          `json:"dnssec_nxdomain_enabled,omitempty"`
	DnssecRpzEnabled                 bool                          `json:"dnssec_rpz_enabled,omitempty"`
	DnssecTrustedKeys                []*Dnssectrustedkey           `json:"dnssec_trusted_keys,omitempty"`
	DnssecValidationEnabled          bool                          `json:"dnssec_validation_enabled,omitempty"`
	DnstapSetting                    *Dnstapsetting                `json:"dnstap_setting,omitempty"`
	DomainsToCaptureDnsQueries       []string                      `json:"domains_to_capture_dns_queries,omitempty"`
	DtcDnsQueriesSpecificBehavior    string                        `json:"dtc_dns_queries_specific_behavior,omitempty"`
	DtcEdnsPreferClientSubnet        bool                          `json:"dtc_edns_prefer_client_subnet,omitempty"`
	DtcHealthSource                  string                        `json:"dtc_health_source,omitempty"`
	DtcHealthSourceAddress           string                        `json:"dtc_health_source_address,omitempty"`
	EdnsUdpSize                      uint32                        `json:"edns_udp_size,omitempty"`
	EnableBlackhole                  bool                          `json:"enable_blackhole,omitempty"`
	EnableBlacklist                  bool                          `json:"enable_blacklist,omitempty"`
	EnableCaptureDnsQueries          bool                          `json:"enable_capture_dns_queries,omitempty"`
	EnableCaptureDnsResponses        bool                          `json:"enable_capture_dns_responses,omitempty"`
	EnableDns                        bool                          `json:"enable_dns,omitempty"`
	EnableDns64                      bool                          `json:"enable_dns64,omitempty"`
	EnableDnsCacheAcceleration       bool                          `json:"enable_dns_cache_acceleration,omitempty"`
	EnableDnsHealthCheck             bool                          `json:"enable_dns_health_check,omitempty"`
	EnableDnstapQueries              bool                          `json:"enable_dnstap_queries,omitempty"`
	EnableDnstapResponses            bool                          `json:"enable_dnstap_responses,omitempty"`
	EnableExcludedDomainNames        bool                          `json:"enable_excluded_domain_names,omitempty"`
	EnableFixedRrsetOrderFqdns       bool                          `json:"enable_fixed_rrset_order_fqdns,omitempty"`
	EnableFtc                        bool                          `json:"enable_ftc,omitempty"`
	EnableGssTsig                    bool                          `json:"enable_gss_tsig,omitempty"`
	EnableNotifySourcePort           bool                          `json:"enable_notify_source_port,omitempty"`
	EnableQueryRewrite               bool                          `json:"enable_query_rewrite,omitempty"`
	EnableQuerySourcePort            bool                          `json:"enable_query_source_port,omitempty"`
	ExcludedDomainNames              []string                      `json:"excluded_domain_names,omitempty"`
	Ea                               EA                            `json:"extattrs,omitempty"`
	FileTransferSetting              *Filetransfersetting          `json:"file_transfer_setting,omitempty"`
	FilterAaaa                       string                        `json:"filter_aaaa,omitempty"`
	FilterAaaaList                   []*Addressac                  `json:"filter_aaaa_list,omitempty"`
	FixedRrsetOrderFqdns             []*GridDnsFixedrrsetorderfqdn `json:"fixed_rrset_order_fqdns,omitempty"`
	ForwardOnly                      bool                          `json:"forward_only,omitempty"`
	ForwardUpdates                   bool                          `json:"forward_updates,omitempty"`
	Forwarders                       []string                      `json:"forwarders,omitempty"`
	FtcExpiredRecordTimeout          uint32                        `json:"ftc_expired_record_timeout,omitempty"`
	FtcExpiredRecordTtl              uint32                        `json:"ftc_expired_record_ttl,omitempty"`
	GlueRecordAddresses              []*MemberDnsgluerecordaddr    `json:"glue_record_addresses,omitempty"`
	GssTsigKeys                      []*Kerberoskey                `json:"gss_tsig_keys,omitempty"`
	HostName                         string                        `json:"host_name,omitempty"`
	Ipv4Addr                         string                        `json:"ipv4addr,omitempty"`
	Ipv6GlueRecordAddresses          []*MemberDnsgluerecordaddr    `json:"ipv6_glue_record_addresses,omitempty"`
	Ipv6Addr                         string                        `json:"ipv6addr,omitempty"`
	IsUnboundCapable                 bool                          `json:"is_unbound_capable,omitempty"`
	LameTtl                          uint32                        `json:"lame_ttl,omitempty"`
	Lan1Ipsd                         string                        `json:"lan1_ipsd,omitempty"`
	Lan1Ipv6Ipsd                     string                        `json:"lan1_ipv6_ipsd,omitempty"`
	Lan2Ipsd                         string                        `json:"lan2_ipsd,omitempty"`
	Lan2Ipv6Ipsd                     string                        `json:"lan2_ipv6_ipsd,omitempty"`
	LoggingCategories                *GridLoggingcategories        `json:"logging_categories,omitempty"`
	MaxCacheTtl                      uint32                        `json:"max_cache_ttl,omitempty"`
	MaxCachedLifetime                uint32                        `json:"max_cached_lifetime,omitempty"`
	MaxNcacheTtl                     uint32                        `json:"max_ncache_ttl,omitempty"`
	MaxUdpSize                       uint32                        `json:"max_udp_size,omitempty"`
	MgmtIpsd                         string                        `json:"mgmt_ipsd,omitempty"`
	MgmtIpv6Ipsd                     string                        `json:"mgmt_ipv6_ipsd,omitempty"`
	MinimalResp                      bool                          `json:"minimal_resp,omitempty"`
	NotifyDelay                      uint32                        `json:"notify_delay,omitempty"`
	NotifySourcePort                 uint32                        `json:"notify_source_port,omitempty"`
	NxdomainLogQuery                 bool                          `json:"nxdomain_log_query,omitempty"`
	NxdomainRedirect                 bool                          `json:"nxdomain_redirect,omitempty"`
	NxdomainRedirectAddresses        []string                      `json:"nxdomain_redirect_addresses,omitempty"`
	NxdomainRedirectAddressesV6      []string                      `json:"nxdomain_redirect_addresses_v6,omitempty"`
	NxdomainRedirectTtl              uint32                        `json:"nxdomain_redirect_ttl,omitempty"`
	NxdomainRulesets                 []string                      `json:"nxdomain_rulesets,omitempty"`
	QuerySourcePort                  uint32                        `json:"query_source_port,omitempty"`
	RecordNamePolicy                 string                        `json:"record_name_policy,omitempty"`
	RecursiveClientLimit             uint32                        `json:"recursive_client_limit,omitempty"`
	RecursiveQueryList               []*Addressac                  `json:"recursive_query_list,omitempty"`
	RecursiveResolver                string                        `json:"recursive_resolver,omitempty"`
	ResolverQueryTimeout             uint32                        `json:"resolver_query_timeout,omitempty"`
	ResponseRateLimiting             *GridResponseratelimiting     `json:"response_rate_limiting,omitempty"`
	RootNameServerType               string                        `json:"root_name_server_type,omitempty"`
	RpzDisableNsdnameNsip            bool                          `json:"rpz_disable_nsdname_nsip,omitempty"`
	RpzDropIpRuleEnabled             bool                          `json:"rpz_drop_ip_rule_enabled,omitempty"`
	RpzDropIpRuleMinPrefixLengthIpv4 uint32                        `json:"rpz_drop_ip_rule_min_prefix_length_ipv4,omitempty"`
	RpzDropIpRuleMinPrefixLengthIpv6 uint32                        `json:"rpz_drop_ip_rule_min_prefix_length_ipv6,omitempty"`
	RpzQnameWaitRecurse              bool                          `json:"rpz_qname_wait_recurse,omitempty"`
	SerialQueryRate                  uint32                        `json:"serial_query_rate,omitempty"`
	ServerIdDirective                string                        `json:"server_id_directive,omitempty"`
	ServerIdDirectiveString          string                        `json:"server_id_directive_string,omitempty"`
	SkipInGridRpzQueries             bool                          `json:"skip_in_grid_rpz_queries,omitempty"`
	Sortlist                         []*Sortlist                   `json:"sortlist,omitempty"`
	StoreLocally                     bool                          `json:"store_locally,omitempty"`
	SyslogFacility                   string                        `json:"syslog_facility,omitempty"`
	TransferExcludedServers          []string                      `json:"transfer_excluded_servers,omitempty"`
	TransferFormat                   string                        `json:"transfer_format,omitempty"`
	TransfersIn                      uint32                        `json:"transfers_in,omitempty"`
	TransfersOut                     uint32                        `json:"transfers_out,omitempty"`
	TransfersPerNs                   uint32                        `json:"transfers_per_ns,omitempty"`
	UnboundLoggingLevel              string                        `json:"unbound_logging_level,omitempty"`
	UseAddClientIpMacOptions         bool                          `json:"use_add_client_ip_mac_options,omitempty"`
	UseAllowQuery                    bool                          `json:"use_allow_query,omitempty"`
	UseAllowTransfer                 bool                          `json:"use_allow_transfer,omitempty"`
	UseAttackMitigation              bool                          `json:"use_attack_mitigation,omitempty"`
	UseAutoBlackhole                 bool                          `json:"use_auto_blackhole,omitempty"`
	UseBindHostnameDirective         bool                          `json:"use_bind_hostname_directive,omitempty"`
	UseBlackhole                     bool                          `json:"use_blackhole,omitempty"`
	UseBlacklist                     bool                          `json:"use_blacklist,omitempty"`
	UseCaptureDnsQueriesOnAllDomains bool                          `json:"use_capture_dns_queries_on_all_domains,omitempty"`
	UseCopyClientIpMacOptions        bool                          `json:"use_copy_client_ip_mac_options,omitempty"`
	UseCopyXferToNotify              bool                          `json:"use_copy_xfer_to_notify,omitempty"`
	UseDisableEdns                   bool                          `json:"use_disable_edns,omitempty"`
	UseDns64                         bool                          `json:"use_dns64,omitempty"`
	UseDnsCacheAccelerationTtl       bool                          `json:"use_dns_cache_acceleration_ttl,omitempty"`
	UseDnsHealthCheck                bool                          `json:"use_dns_health_check,omitempty"`
	UseDnssec                        bool                          `json:"use_dnssec,omitempty"`
	UseDnstapSetting                 bool                          `json:"use_dnstap_setting,omitempty"`
	UseDtcDnsQueriesSpecificBehavior bool                          `json:"use_dtc_dns_queries_specific_behavior,omitempty"`
	UseDtcEdnsPreferClientSubnet     bool                          `json:"use_dtc_edns_prefer_client_subnet,omitempty"`
	UseEdnsUdpSize                   bool                          `json:"use_edns_udp_size,omitempty"`
	UseEnableCaptureDns              bool                          `json:"use_enable_capture_dns,omitempty"`
	UseEnableExcludedDomainNames     bool                          `json:"use_enable_excluded_domain_names,omitempty"`
	UseEnableGssTsig                 bool                          `json:"use_enable_gss_tsig,omitempty"`
	UseEnableQueryRewrite            bool                          `json:"use_enable_query_rewrite,omitempty"`
	UseFilterAaaa                    bool                          `json:"use_filter_aaaa,omitempty"`
	UseFixedRrsetOrderFqdns          bool                          `json:"use_fixed_rrset_order_fqdns,omitempty"`
	UseForwardUpdates                bool                          `json:"use_forward_updates,omitempty"`
	UseForwarders                    bool                          `json:"use_forwarders,omitempty"`
	UseFtc                           bool                          `json:"use_ftc,omitempty"`
	UseGssTsigKeys                   bool                          `json:"use_gss_tsig_keys,omitempty"`
	UseLameTtl                       bool                          `json:"use_lame_ttl,omitempty"`
	UseLan2Ipv6Port                  bool                          `json:"use_lan2_ipv6_port,omitempty"`
	UseLan2Port                      bool                          `json:"use_lan2_port,omitempty"`
	UseLanIpv6Port                   bool                          `json:"use_lan_ipv6_port,omitempty"`
	UseLanPort                       bool                          `json:"use_lan_port,omitempty"`
	UseLoggingCategories             bool                          `json:"use_logging_categories,omitempty"`
	UseMaxCacheTtl                   bool                          `json:"use_max_cache_ttl,omitempty"`
	UseMaxCachedLifetime             bool                          `json:"use_max_cached_lifetime,omitempty"`
	UseMaxNcacheTtl                  bool                          `json:"use_max_ncache_ttl,omitempty"`
	UseMaxUdpSize                    bool                          `json:"use_max_udp_size,omitempty"`
	UseMgmtIpv6Port                  bool                          `json:"use_mgmt_ipv6_port,omitempty"`
	UseMgmtPort                      bool                          `json:"use_mgmt_port,omitempty"`
	UseNotifyDelay                   bool                          `json:"use_notify_delay,omitempty"`
	UseNxdomainRedirect              bool                          `json:"use_nxdomain_redirect,omitempty"`
	UseRecordNamePolicy              bool                          `json:"use_record_name_policy,omitempty"`
	UseRecursiveClientLimit          bool                          `json:"use_recursive_client_limit,omitempty"`
	UseRecursiveQuerySetting         bool                          `json:"use_recursive_query_setting,omitempty"`
	UseResolverQueryTimeout          bool                          `json:"use_resolver_query_timeout,omitempty"`
	UseResponseRateLimiting          bool                          `json:"use_response_rate_limiting,omitempty"`
	UseRootNameServer                bool                          `json:"use_root_name_server,omitempty"`
	UseRootServerForAllViews         bool                          `json:"use_root_server_for_all_views,omitempty"`
	UseRpzDisableNsdnameNsip         bool                          `json:"use_rpz_disable_nsdname_nsip,omitempty"`
	UseRpzDropIpRule                 bool                          `json:"use_rpz_drop_ip_rule,omitempty"`
	UseRpzQnameWaitRecurse           bool                          `json:"use_rpz_qname_wait_recurse,omitempty"`
	UseSerialQueryRate               bool                          `json:"use_serial_query_rate,omitempty"`
	UseServerIdDirective             bool                          `json:"use_server_id_directive,omitempty"`
	UseSortlist                      bool                          `json:"use_sortlist,omitempty"`
	UseSourcePorts                   bool                          `json:"use_source_ports,omitempty"`
	UseSyslogFacility                bool                          `json:"use_syslog_facility,omitempty"`
	UseTransfersIn                   bool                          `json:"use_transfers_in,omitempty"`
	UseTransfersOut                  bool                          `json:"use_transfers_out,omitempty"`
	UseTransfersPerNs                bool                          `json:"use_transfers_per_ns,omitempty"`
	UseUpdateSetting                 bool                          `json:"use_update_setting,omitempty"`
	UseZoneTransferFormat            bool                          `json:"use_zone_transfer_format,omitempty"`
	Views                            []string                      `json:"views,omitempty"`
}

func (MemberDns) ObjectType() string {
	return "member:dns"
}

func (obj MemberDns) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"host_name", "ipv4addr", "ipv6addr"}
	}
	return obj.returnFields
}

// Msserver represents Infoblox object msserver
type Msserver struct {
	IBBase                      `json:"-"`
	Ref                         string          `json:"_ref,omitempty"`
	AdDomain                    string          `json:"ad_domain,omitempty"`
	AdSites                     *Adsites        `json:"ad_sites,omitempty"`
	AdUser                      *MsserverAduser `json:"ad_user,omitempty"`
	Address                     string          `json:"address,omitempty"`
	Comment                     string          `json:"comment,omitempty"`
	ConnectionStatus            string          `json:"connection_status,omitempty"`
	ConnectionStatusDetail      string          `json:"connection_status_detail,omitempty"`
	DhcpServer                  *Dhcpserver     `json:"dhcp_server,omitempty"`
	Disabled                    bool            `json:"disabled,omitempty"`
	DnsServer                   *Dnsserver      `json:"dns_server,omitempty"`
	DnsView                     string          `json:"dns_view,omitempty"`
	Ea                          EA              `json:"extattrs,omitempty"`
	GridMember                  string          `json:"grid_member,omitempty"`
	LastSeen                    *time.Time      `json:"last_seen,omitempty"`
	LogDestination              string          `json:"log_destination,omitempty"`
	LogLevel                    string          `json:"log_level,omitempty"`
	LoginName                   string          `json:"login_name,omitempty"`
	LoginPassword               string          `json:"login_password,omitempty"`
	ManagingMember              string          `json:"managing_member,omitempty"`
	MsMaxConnection             uint32          `json:"ms_max_connection,omitempty"`
	MsRpcTimeoutInSeconds       uint32          `json:"ms_rpc_timeout_in_seconds,omitempty"`
	NetworkView                 string          `json:"network_view,omitempty"`
	ReadOnly                    bool            `json:"read_only,omitempty"`
	RootAdDomain                string          `json:"root_ad_domain,omitempty"`
	ServerName                  string          `json:"server_name,omitempty"`
	SynchronizationMinDelay     uint32          `json:"synchronization_min_delay,omitempty"`
	SynchronizationStatus       string          `json:"synchronization_status,omitempty"`
	SynchronizationStatusDetail string          `json:"synchronization_status_detail,omitempty"`
	UseLogDestination           bool            `json:"use_log_destination,omitempty"`
	UseMsMaxConnection          bool            `json:"use_ms_max_connection,omitempty"`
	UseMsRpcTimeoutInSeconds    bool            `json:"use_ms_rpc_timeout_in_seconds,omitempty"`
	Version                     string          `json:"version,omitempty"`
}

func (Msserver) ObjectType() string {
	return "msserver"
}

func (obj Msserver) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"address"}
	}
	return obj.returnFields
}

// MsserverDhcp represents Infoblox object msserver:dhcp
type MsserverDhcp struct {
	IBBase                     `json:"-"`
	Ref                        string     `json:"_ref,omitempty"`
	Address                    string     `json:"address,omitempty"`
	Comment                    string     `json:"comment,omitempty"`
	DhcpUtilization            uint32     `json:"dhcp_utilization,omitempty"`
	DhcpUtilizationStatus      string     `json:"dhcp_utilization_status,omitempty"`
	DynamicHosts               uint32     `json:"dynamic_hosts,omitempty"`
	LastSyncTs                 *time.Time `json:"last_sync_ts,omitempty"`
	LoginName                  string     `json:"login_name,omitempty"`
	LoginPassword              string     `json:"login_password,omitempty"`
	NetworkView                string     `json:"network_view,omitempty"`
	NextSyncControl            string     `json:"next_sync_control,omitempty"`
	ReadOnly                   bool       `json:"read_only,omitempty"`
	ServerName                 string     `json:"server_name,omitempty"`
	StaticHosts                uint32     `json:"static_hosts,omitempty"`
	Status                     string     `json:"status,omitempty"`
	StatusDetail               string     `json:"status_detail,omitempty"`
	StatusLastUpdated          *time.Time `json:"status_last_updated,omitempty"`
	SupportsFailover           bool       `json:"supports_failover,omitempty"`
	SynchronizationInterval    uint32     `json:"synchronization_interval,omitempty"`
	TotalHosts                 uint32     `json:"total_hosts,omitempty"`
	UseLogin                   bool       `json:"use_login,omitempty"`
	UseSynchronizationInterval bool       `json:"use_synchronization_interval,omitempty"`
}

func (MsserverDhcp) ObjectType() string {
	return "msserver:dhcp"
}

func (obj MsserverDhcp) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"address"}
	}
	return obj.returnFields
}

// MsserverDns represents Infoblox object msserver:dns
type MsserverDns struct {
	IBBase                     `json:"-"`
	Ref                        string `json:"_ref,omitempty"`
	Address                    string `json:"address,omitempty"`
	EnableDnsReportsSync       bool   `json:"enable_dns_reports_sync,omitempty"`
	LoginName                  string `json:"login_name,omitempty"`
	LoginPassword              string `json:"login_password,omitempty"`
	SynchronizationInterval    uint32 `json:"synchronization_interval,omitempty"`
	UseEnableDnsReportsSync    bool   `json:"use_enable_dns_reports_sync,omitempty"`
	UseLogin                   bool   `json:"use_login,omitempty"`
	UseSynchronizationInterval bool   `json:"use_synchronization_interval,omitempty"`
}

func (MsserverDns) ObjectType() string {
	return "msserver:dns"
}

func (obj MsserverDns) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"address"}
	}
	return obj.returnFields
}

// Mssuperscope represents Infoblox object mssuperscope
type Mssuperscope struct {
	IBBase                `json:"-"`
	Ref                   string   `json:"_ref,omitempty"`
	Comment               string   `json:"comment,omitempty"`
	DhcpUtilization       uint32   `json:"dhcp_utilization,omitempty"`
	DhcpUtilizationStatus string   `json:"dhcp_utilization_status,omitempty"`
	Disable               bool     `json:"disable,omitempty"`
	DynamicHosts          uint32   `json:"dynamic_hosts,omitempty"`
	Ea                    EA       `json:"extattrs,omitempty"`
	HighWaterMark         uint32   `json:"high_water_mark,omitempty"`
	HighWaterMarkReset    uint32   `json:"high_water_mark_reset,omitempty"`
	LowWaterMark          uint32   `json:"low_water_mark,omitempty"`
	LowWaterMarkReset     uint32   `json:"low_water_mark_reset,omitempty"`
	Name                  string   `json:"name,omitempty"`
	NetworkView           string   `json:"network_view,omitempty"`
	Ranges                []*Range `json:"ranges,omitempty"`
	StaticHosts           uint32   `json:"static_hosts,omitempty"`
	TotalHosts            uint32   `json:"total_hosts,omitempty"`
}

func (Mssuperscope) ObjectType() string {
	return "mssuperscope"
}

func (obj Mssuperscope) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"disable", "name", "network_view"}
	}
	return obj.returnFields
}

// Namedacl represents Infoblox object namedacl
type Namedacl struct {
	IBBase             `json:"-"`
	Ref                string            `json:"_ref,omitempty"`
	AccessList         []*Addressac      `json:"access_list,omitempty"`
	Comment            string            `json:"comment,omitempty"`
	ExplodedAccessList []*Addressac      `json:"exploded_access_list,omitempty"`
	Ea                 EA                `json:"extattrs,omitempty"`
	Name               string            `json:"name,omitempty"`
	ValidateAclItems   *Validateaclitems `json:"validate_acl_items,omitempty"`
}

func (Namedacl) ObjectType() string {
	return "namedacl"
}

func (obj Namedacl) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name"}
	}
	return obj.returnFields
}

// Natgroup represents Infoblox object natgroup
type Natgroup struct {
	IBBase  `json:"-"`
	Ref     string `json:"_ref,omitempty"`
	Comment string `json:"comment,omitempty"`
	Name    string `json:"name,omitempty"`
}

func (Natgroup) ObjectType() string {
	return "natgroup"
}

func (obj Natgroup) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name"}
	}
	return obj.returnFields
}

// NetworkDiscovery represents Infoblox object network_discovery
type NetworkDiscovery struct {
	IBBase             `json:"-"`
	Ref                string              `json:"_ref,omitempty"`
	ClearDiscoveryData *Cleardiscoverydata `json:"clear_discovery_data,omitempty"`
}

func (NetworkDiscovery) ObjectType() string {
	return "network_discovery"
}

func (obj NetworkDiscovery) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{}
	}
	return obj.returnFields
}

// Member represents Infoblox object member
type Member struct {
	IBBase                          `json:"-"`
	Ref                             string                          `json:"_ref,omitempty"`
	ActivePosition                  string                          `json:"active_position,omitempty"`
	AdditionalIpList                []*Interface                    `json:"additional_ip_list,omitempty"`
	AutomatedTrafficCaptureSetting  *SettingAutomatedtrafficcapture `json:"automated_traffic_capture_setting,omitempty"`
	BgpAs                           []*Bgpas                        `json:"bgp_as,omitempty"`
	CaptureTrafficControl           *Membercapturecontrolparams     `json:"capture_traffic_control,omitempty"`
	CaptureTrafficStatus            *Membercapturestatusparams      `json:"capture_traffic_status,omitempty"`
	Comment                         string                          `json:"comment,omitempty"`
	ConfigAddrType                  string                          `json:"config_addr_type,omitempty"`
	CreateToken                     *Pnodetokenoperation            `json:"create_token,omitempty"`
	CspAccessKey                    []string                        `json:"csp_access_key,omitempty"`
	CspMemberSetting                *MemberCspmembersetting         `json:"csp_member_setting,omitempty"`
	DnsResolverSetting              *SettingDnsresolver             `json:"dns_resolver_setting,omitempty"`
	Dscp                            uint32                          `json:"dscp,omitempty"`
	EmailSetting                    *SettingEmail                   `json:"email_setting,omitempty"`
	EnableHa                        bool                            `json:"enable_ha,omitempty"`
	EnableLom                       bool                            `json:"enable_lom,omitempty"`
	EnableMemberRedirect            bool                            `json:"enable_member_redirect,omitempty"`
	EnableRoApiAccess               bool                            `json:"enable_ro_api_access,omitempty"`
	Ea                              EA                              `json:"extattrs,omitempty"`
	ExternalSyslogBackupServers     []*Extsyslogbackupserver        `json:"external_syslog_backup_servers,omitempty"`
	ExternalSyslogServerEnable      bool                            `json:"external_syslog_server_enable,omitempty"`
	HostName                        string                          `json:"host_name,omitempty"`
	Ipv6Setting                     *Ipv6setting                    `json:"ipv6_setting,omitempty"`
	Ipv6StaticRoutes                []*Ipv6networksetting           `json:"ipv6_static_routes,omitempty"`
	IsDscpCapable                   bool                            `json:"is_dscp_capable,omitempty"`
	Lan2Enabled                     bool                            `json:"lan2_enabled,omitempty"`
	Lan2PortSetting                 *Lan2portsetting                `json:"lan2_port_setting,omitempty"`
	LcdInput                        bool                            `json:"lcd_input,omitempty"`
	LomNetworkConfig                []*Lomnetworkconfig             `json:"lom_network_config,omitempty"`
	LomUsers                        []*Lomuser                      `json:"lom_users,omitempty"`
	MasterCandidate                 bool                            `json:"master_candidate,omitempty"`
	MemberAdminOperation            *Memberadminoperationparams     `json:"member_admin_operation,omitempty"`
	MemberServiceCommunication      []*Memberservicecommunication   `json:"member_service_communication,omitempty"`
	MgmtPortSetting                 *Mgmtportsetting                `json:"mgmt_port_setting,omitempty"`
	MmdbEaBuildTime                 *time.Time                      `json:"mmdb_ea_build_time,omitempty"`
	MmdbGeoipBuildTime              *time.Time                      `json:"mmdb_geoip_build_time,omitempty"`
	NatSetting                      *Natsetting                     `json:"nat_setting,omitempty"`
	NodeInfo                        []*Nodeinfo                     `json:"node_info,omitempty"`
	NTPSetting                      *MemberNtp                      `json:"ntp_setting,omitempty"`
	OspfList                        []*Ospf                         `json:"ospf_list,omitempty"`
	PassiveHaArpEnabled             bool                            `json:"passive_ha_arp_enabled,omitempty"`
	Platform                        string                          `json:"platform,omitempty"`
	PreProvisioning                 *Preprovision                   `json:"pre_provisioning,omitempty"`
	PreserveIfOwnsDelegation        bool                            `json:"preserve_if_owns_delegation,omitempty"`
	ReadToken                       *Pnodetokenoperation            `json:"read_token,omitempty"`
	RemoteConsoleAccessEnable       bool                            `json:"remote_console_access_enable,omitempty"`
	Requestrestartservicestatus     *Requestmemberservicestatus     `json:"requestrestartservicestatus,omitempty"`
	Restartservices                 *Memberrestartservices          `json:"restartservices,omitempty"`
	RouterId                        uint32                          `json:"router_id,omitempty"`
	ServiceStatus                   []*Memberservicestatus          `json:"service_status,omitempty"`
	ServiceTypeConfiguration        string                          `json:"service_type_configuration,omitempty"`
	SnmpSetting                     *SettingSnmp                    `json:"snmp_setting,omitempty"`
	StaticRoutes                    []*SettingNetwork               `json:"static_routes,omitempty"`
	SupportAccessEnable             bool                            `json:"support_access_enable,omitempty"`
	SupportAccessInfo               string                          `json:"support_access_info,omitempty"`
	SyslogProxySetting              *SettingSyslogproxy             `json:"syslog_proxy_setting,omitempty"`
	SyslogServers                   []*Syslogserver                 `json:"syslog_servers,omitempty"`
	SyslogSize                      uint32                          `json:"syslog_size,omitempty"`
	ThresholdTraps                  []*Thresholdtrap                `json:"threshold_traps,omitempty"`
	TimeZone                        string                          `json:"time_zone,omitempty"`
	TrafficCaptureAuthDnsSetting    *SettingTriggeruthdnslatency    `json:"traffic_capture_auth_dns_setting,omitempty"`
	TrafficCaptureChrSetting        *SettingTrafficcapturechr       `json:"traffic_capture_chr_setting,omitempty"`
	TrafficCaptureQpsSetting        *SettingTrafficcaptureqps       `json:"traffic_capture_qps_setting,omitempty"`
	TrafficCaptureRecDnsSetting     *SettingTriggerrecdnslatency    `json:"traffic_capture_rec_dns_setting,omitempty"`
	TrafficCaptureRecQueriesSetting *SettingTriggerrecqueries       `json:"traffic_capture_rec_queries_setting,omitempty"`
	TrapNotifications               []*Trapnotification             `json:"trap_notifications,omitempty"`
	UpgradeGroup                    string                          `json:"upgrade_group,omitempty"`
	UseAutomatedTrafficCapture      bool                            `json:"use_automated_traffic_capture,omitempty"`
	UseDnsResolverSetting           bool                            `json:"use_dns_resolver_setting,omitempty"`
	UseDscp                         bool                            `json:"use_dscp,omitempty"`
	UseEmailSetting                 bool                            `json:"use_email_setting,omitempty"`
	UseEnableLom                    bool                            `json:"use_enable_lom,omitempty"`
	UseEnableMemberRedirect         bool                            `json:"use_enable_member_redirect,omitempty"`
	UseExternalSyslogBackupServers  bool                            `json:"use_external_syslog_backup_servers,omitempty"`
	UseLcdInput                     bool                            `json:"use_lcd_input,omitempty"`
	UseRemoteConsoleAccessEnable    bool                            `json:"use_remote_console_access_enable,omitempty"`
	UseSnmpSetting                  bool                            `json:"use_snmp_setting,omitempty"`
	UseSupportAccessEnable          bool                            `json:"use_support_access_enable,omitempty"`
	UseSyslogProxySetting           bool                            `json:"use_syslog_proxy_setting,omitempty"`
	UseThresholdTraps               bool                            `json:"use_threshold_traps,omitempty"`
	UseTimeZone                     bool                            `json:"use_time_zone,omitempty"`
	UseTrafficCaptureAuthDns        bool                            `json:"use_traffic_capture_auth_dns,omitempty"`
	UseTrafficCaptureChr            bool                            `json:"use_traffic_capture_chr,omitempty"`
	UseTrafficCaptureQps            bool                            `json:"use_traffic_capture_qps,omitempty"`
	UseTrafficCaptureRecDns         bool                            `json:"use_traffic_capture_rec_dns,omitempty"`
	UseTrafficCaptureRecQueries     bool                            `json:"use_traffic_capture_rec_queries,omitempty"`
	UseTrapNotifications            bool                            `json:"use_trap_notifications,omitempty"`
	UseV4Vrrp                       bool                            `json:"use_v4_vrrp,omitempty"`
	VipSetting                      *SettingNetwork                 `json:"vip_setting,omitempty"`
	VpnMtu                          uint32                          `json:"vpn_mtu,omitempty"`
}

func (Member) ObjectType() string {
	return "member"
}

func (obj Member) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"config_addr_type", "host_name", "platform", "service_type_configuration"}
	}
	return obj.returnFields
}

func NewMember(member Member) *Member {
	res := member
	returnFields := []string{"host_name", "node_info", "time_zone"}
	res.returnFields = returnFields
	return &res
}

// Networktemplate represents Infoblox object networktemplate
type Networktemplate struct {
	IBBase                         `json:"-"`
	Ref                            string                `json:"_ref,omitempty"`
	AllowAnyNetmask                bool                  `json:"allow_any_netmask,omitempty"`
	Authority                      bool                  `json:"authority,omitempty"`
	AutoCreateReversezone          bool                  `json:"auto_create_reversezone,omitempty"`
	Bootfile                       string                `json:"bootfile,omitempty"`
	Bootserver                     string                `json:"bootserver,omitempty"`
	CloudApiCompatible             bool                  `json:"cloud_api_compatible,omitempty"`
	Comment                        string                `json:"comment,omitempty"`
	DdnsDomainname                 string                `json:"ddns_domainname,omitempty"`
	DdnsGenerateHostname           bool                  `json:"ddns_generate_hostname,omitempty"`
	DdnsServerAlwaysUpdates        bool                  `json:"ddns_server_always_updates,omitempty"`
	DdnsTtl                        uint32                `json:"ddns_ttl,omitempty"`
	DdnsUpdateFixedAddresses       bool                  `json:"ddns_update_fixed_addresses,omitempty"`
	DdnsUseOption81                bool                  `json:"ddns_use_option81,omitempty"`
	DelegatedMember                *Dhcpmember           `json:"delegated_member,omitempty"`
	DenyBootp                      bool                  `json:"deny_bootp,omitempty"`
	EmailList                      []string              `json:"email_list,omitempty"`
	EnableDdns                     bool                  `json:"enable_ddns,omitempty"`
	EnableDhcpThresholds           bool                  `json:"enable_dhcp_thresholds,omitempty"`
	EnableEmailWarnings            bool                  `json:"enable_email_warnings,omitempty"`
	EnablePxeLeaseTime             bool                  `json:"enable_pxe_lease_time,omitempty"`
	EnableSnmpWarnings             bool                  `json:"enable_snmp_warnings,omitempty"`
	Ea                             EA                    `json:"extattrs,omitempty"`
	FixedAddressTemplates          []string              `json:"fixed_address_templates,omitempty"`
	HighWaterMark                  uint32                `json:"high_water_mark,omitempty"`
	HighWaterMarkReset             uint32                `json:"high_water_mark_reset,omitempty"`
	IgnoreDhcpOptionListRequest    bool                  `json:"ignore_dhcp_option_list_request,omitempty"`
	IpamEmailAddresses             []string              `json:"ipam_email_addresses,omitempty"`
	IpamThresholdSettings          *SettingIpamThreshold `json:"ipam_threshold_settings,omitempty"`
	IpamTrapSettings               *SettingIpamTrap      `json:"ipam_trap_settings,omitempty"`
	LeaseScavengeTime              int                   `json:"lease_scavenge_time,omitempty"`
	LogicFilterRules               []*Logicfilterrule    `json:"logic_filter_rules,omitempty"`
	LowWaterMark                   uint32                `json:"low_water_mark,omitempty"`
	LowWaterMarkReset              uint32                `json:"low_water_mark_reset,omitempty"`
	Members                        []*Msdhcpserver       `json:"members,omitempty"`
	Name                           string                `json:"name,omitempty"`
	Netmask                        uint32                `json:"netmask,omitempty"`
	Nextserver                     string                `json:"nextserver,omitempty"`
	Options                        []*Dhcpoption         `json:"options,omitempty"`
	PxeLeaseTime                   uint32                `json:"pxe_lease_time,omitempty"`
	RangeTemplates                 []string              `json:"range_templates,omitempty"`
	RecycleLeases                  bool                  `json:"recycle_leases,omitempty"`
	Rir                            string                `json:"rir,omitempty"`
	RirOrganization                string                `json:"rir_organization,omitempty"`
	RirRegistrationAction          string                `json:"rir_registration_action,omitempty"`
	RirRegistrationStatus          string                `json:"rir_registration_status,omitempty"`
	SendRirRequest                 bool                  `json:"send_rir_request,omitempty"`
	UpdateDnsOnLeaseRenewal        bool                  `json:"update_dns_on_lease_renewal,omitempty"`
	UseAuthority                   bool                  `json:"use_authority,omitempty"`
	UseBootfile                    bool                  `json:"use_bootfile,omitempty"`
	UseBootserver                  bool                  `json:"use_bootserver,omitempty"`
	UseDdnsDomainname              bool                  `json:"use_ddns_domainname,omitempty"`
	UseDdnsGenerateHostname        bool                  `json:"use_ddns_generate_hostname,omitempty"`
	UseDdnsTtl                     bool                  `json:"use_ddns_ttl,omitempty"`
	UseDdnsUpdateFixedAddresses    bool                  `json:"use_ddns_update_fixed_addresses,omitempty"`
	UseDdnsUseOption81             bool                  `json:"use_ddns_use_option81,omitempty"`
	UseDenyBootp                   bool                  `json:"use_deny_bootp,omitempty"`
	UseEmailList                   bool                  `json:"use_email_list,omitempty"`
	UseEnableDdns                  bool                  `json:"use_enable_ddns,omitempty"`
	UseEnableDhcpThresholds        bool                  `json:"use_enable_dhcp_thresholds,omitempty"`
	UseIgnoreDhcpOptionListRequest bool                  `json:"use_ignore_dhcp_option_list_request,omitempty"`
	UseIpamEmailAddresses          bool                  `json:"use_ipam_email_addresses,omitempty"`
	UseIpamThresholdSettings       bool                  `json:"use_ipam_threshold_settings,omitempty"`
	UseIpamTrapSettings            bool                  `json:"use_ipam_trap_settings,omitempty"`
	UseLeaseScavengeTime           bool                  `json:"use_lease_scavenge_time,omitempty"`
	UseLogicFilterRules            bool                  `json:"use_logic_filter_rules,omitempty"`
	UseNextserver                  bool                  `json:"use_nextserver,omitempty"`
	UseOptions                     bool                  `json:"use_options,omitempty"`
	UsePxeLeaseTime                bool                  `json:"use_pxe_lease_time,omitempty"`
	UseRecycleLeases               bool                  `json:"use_recycle_leases,omitempty"`
	UseUpdateDnsOnLeaseRenewal     bool                  `json:"use_update_dns_on_lease_renewal,omitempty"`
}

func (Networktemplate) ObjectType() string {
	return "networktemplate"
}

func (obj Networktemplate) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name"}
	}
	return obj.returnFields
}

// Networkuser represents Infoblox object networkuser
type Networkuser struct {
	IBBase          `json:"-"`
	Ref             string     `json:"_ref,omitempty"`
	Address         string     `json:"address,omitempty"`
	AddressObject   string     `json:"address_object,omitempty"`
	DataSource      string     `json:"data_source,omitempty"`
	DataSourceIp    string     `json:"data_source_ip,omitempty"`
	Domainname      string     `json:"domainname,omitempty"`
	FirstSeenTime   *time.Time `json:"first_seen_time,omitempty"`
	Guid            string     `json:"guid,omitempty"`
	LastSeenTime    *time.Time `json:"last_seen_time,omitempty"`
	LastUpdatedTime *time.Time `json:"last_updated_time,omitempty"`
	LogonId         string     `json:"logon_id,omitempty"`
	LogoutTime      *time.Time `json:"logout_time,omitempty"`
	Name            string     `json:"name,omitempty"`
	Network         string     `json:"network,omitempty"`
	NetworkView     string     `json:"network_view,omitempty"`
	UserStatus      string     `json:"user_status,omitempty"`
}

func (Networkuser) ObjectType() string {
	return "networkuser"
}

func (obj Networkuser) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"address", "domainname", "name", "network_view", "user_status"}
	}
	return obj.returnFields
}

// NotificationRestEndpoint represents Infoblox object notification:rest:endpoint
type NotificationRestEndpoint struct {
	IBBase                     `json:"-"`
	Ref                        string                            `json:"_ref,omitempty"`
	ClearOutboundWorkerLog     *Clearworkerlog                   `json:"clear_outbound_worker_log,omitempty"`
	ClientCertificateSubject   string                            `json:"client_certificate_subject,omitempty"`
	ClientCertificateToken     string                            `json:"client_certificate_token,omitempty"`
	ClientCertificateValidFrom *time.Time                        `json:"client_certificate_valid_from,omitempty"`
	ClientCertificateValidTo   *time.Time                        `json:"client_certificate_valid_to,omitempty"`
	Comment                    string                            `json:"comment,omitempty"`
	Ea                         EA                                `json:"extattrs,omitempty"`
	LogLevel                   string                            `json:"log_level,omitempty"`
	Name                       string                            `json:"name,omitempty"`
	OutboundMemberType         string                            `json:"outbound_member_type,omitempty"`
	OutboundMembers            []string                          `json:"outbound_members,omitempty"`
	Password                   string                            `json:"password,omitempty"`
	ServerCertValidation       string                            `json:"server_cert_validation,omitempty"`
	SyncDisabled               bool                              `json:"sync_disabled,omitempty"`
	TemplateInstance           *NotificationRestTemplateinstance `json:"template_instance,omitempty"`
	TestConnection             *Testconnectivityparams           `json:"test_connection,omitempty"`
	Timeout                    uint32                            `json:"timeout,omitempty"`
	Uri                        string                            `json:"uri,omitempty"`
	Username                   string                            `json:"username,omitempty"`
	VendorIdentifier           string                            `json:"vendor_identifier,omitempty"`
	WapiUserName               string                            `json:"wapi_user_name,omitempty"`
	WapiUserPassword           string                            `json:"wapi_user_password,omitempty"`
}

func (NotificationRestEndpoint) ObjectType() string {
	return "notification:rest:endpoint"
}

func (obj NotificationRestEndpoint) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name", "outbound_member_type", "uri"}
	}
	return obj.returnFields
}

// NotificationRestTemplate represents Infoblox object notification:rest:template
type NotificationRestTemplate struct {
	IBBase           `json:"-"`
	Ref              string                               `json:"_ref,omitempty"`
	ActionName       string                               `json:"action_name,omitempty"`
	AddedOn          *time.Time                           `json:"added_on,omitempty"`
	Comment          string                               `json:"comment,omitempty"`
	Content          string                               `json:"content,omitempty"`
	EventType        []string                             `json:"event_type,omitempty"`
	Name             string                               `json:"name,omitempty"`
	OutboundType     string                               `json:"outbound_type,omitempty"`
	Parameters       []*NotificationRestTemplateparameter `json:"parameters,omitempty"`
	TemplateType     string                               `json:"template_type,omitempty"`
	VendorIdentifier string                               `json:"vendor_identifier,omitempty"`
}

func (NotificationRestTemplate) ObjectType() string {
	return "notification:rest:template"
}

func (obj NotificationRestTemplate) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"content", "name"}
	}
	return obj.returnFields
}

// NetworkView represents Infoblox object networkview
type NetworkView struct {
	IBBase               `json:"-"`
	Ref                  string                    `json:"_ref,omitempty"`
	AssociatedDnsViews   []string                  `json:"associated_dns_views,omitempty"`
	AssociatedMembers    []*NetworkviewAssocmember `json:"associated_members,omitempty"`
	CloudInfo            *GridCloudapiInfo         `json:"cloud_info,omitempty"`
	Comment              string                    `json:"comment,omitempty"`
	DdnsDnsView          string                    `json:"ddns_dns_view,omitempty"`
	DdnsZonePrimaries    []*Dhcpddns               `json:"ddns_zone_primaries,omitempty"`
	Ea                   EA                        `json:"extattrs,omitempty"`
	InternalForwardZones []*ZoneAuth               `json:"internal_forward_zones,omitempty"`
	IsDefault            bool                      `json:"is_default,omitempty"`
	MgmPrivate           bool                      `json:"mgm_private,omitempty"`
	MsAdUserData         *MsserverAduserData       `json:"ms_ad_user_data,omitempty"`
	Name                 string                    `json:"name,omitempty"`
	RemoteForwardZones   []*Remoteddnszone         `json:"remote_forward_zones,omitempty"`
	RemoteReverseZones   []*Remoteddnszone         `json:"remote_reverse_zones,omitempty"`
}

func (NetworkView) ObjectType() string {
	return "networkview"
}

func (obj NetworkView) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "is_default", "name"}
	}
	return obj.returnFields
}

func NewEmptyNetworkView() *NetworkView {
	res := &NetworkView{}
	res.returnFields = []string{"extattrs", "name", "comment"}
	return res
}

func NewNetworkView(name string, comment string, eas EA, ref string) *NetworkView {
	res := NewEmptyNetworkView()
	res.Name = name
	res.Comment = comment
	res.Ea = eas
	res.Ref = ref
	return res
}

// NotificationRule represents Infoblox object notification:rule
type NotificationRule struct {
	IBBase                           `json:"-"`
	Ref                              string                            `json:"_ref,omitempty"`
	AllMembers                       bool                              `json:"all_members,omitempty"`
	Comment                          string                            `json:"comment,omitempty"`
	Disable                          bool                              `json:"disable,omitempty"`
	EnableEventDeduplication         bool                              `json:"enable_event_deduplication,omitempty"`
	EnableEventDeduplicationLog      bool                              `json:"enable_event_deduplication_log,omitempty"`
	EventDeduplicationFields         []string                          `json:"event_deduplication_fields,omitempty"`
	EventDeduplicationLookbackPeriod uint32                            `json:"event_deduplication_lookback_period,omitempty"`
	EventPriority                    string                            `json:"event_priority,omitempty"`
	EventType                        string                            `json:"event_type,omitempty"`
	ExpressionList                   []*NotificationRuleexpressionop   `json:"expression_list,omitempty"`
	Name                             string                            `json:"name,omitempty"`
	NotificationAction               string                            `json:"notification_action,omitempty"`
	NotificationTarget               string                            `json:"notification_target,omitempty"`
	PublishSettings                  *CiscoisePublishsetting           `json:"publish_settings,omitempty"`
	ScheduledEvent                   *SettingSchedule                  `json:"scheduled_event,omitempty"`
	SelectedMembers                  []string                          `json:"selected_members,omitempty"`
	TemplateInstance                 *NotificationRestTemplateinstance `json:"template_instance,omitempty"`
	TriggerOutbound                  *Triggeroutboundparams            `json:"trigger_outbound,omitempty"`
	UsePublishSettings               bool                              `json:"use_publish_settings,omitempty"`
}

func (NotificationRule) ObjectType() string {
	return "notification:rule"
}

func (obj NotificationRule) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"event_type", "name", "notification_action", "notification_target"}
	}
	return obj.returnFields
}

// NsgroupDelegation represents Infoblox object nsgroup:delegation
type NsgroupDelegation struct {
	IBBase     `json:"-"`
	Ref        string       `json:"_ref,omitempty"`
	Comment    string       `json:"comment,omitempty"`
	DelegateTo []NameServer `json:"delegate_to,omitempty"`
	Ea         EA           `json:"extattrs,omitempty"`
	Name       string       `json:"name,omitempty"`
}

func (NsgroupDelegation) ObjectType() string {
	return "nsgroup:delegation"
}

func (obj NsgroupDelegation) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"delegate_to", "name"}
	}
	return obj.returnFields
}

// Ipv4NetworkContainer represents Infoblox object networkcontainer
type Ipv4NetworkContainer struct {
	IBBase                           `json:"-"`
	Ref                              string                      `json:"_ref,omitempty"`
	Authority                        bool                        `json:"authority,omitempty"`
	AutoCreateReversezone            bool                        `json:"auto_create_reversezone,omitempty"`
	Bootfile                         string                      `json:"bootfile,omitempty"`
	Bootserver                       string                      `json:"bootserver,omitempty"`
	CloudInfo                        *GridCloudapiInfo           `json:"cloud_info,omitempty"`
	Comment                          string                      `json:"comment,omitempty"`
	DdnsDomainname                   string                      `json:"ddns_domainname,omitempty"`
	DdnsGenerateHostname             bool                        `json:"ddns_generate_hostname,omitempty"`
	DdnsServerAlwaysUpdates          bool                        `json:"ddns_server_always_updates,omitempty"`
	DdnsTtl                          uint32                      `json:"ddns_ttl,omitempty"`
	DdnsUpdateFixedAddresses         bool                        `json:"ddns_update_fixed_addresses,omitempty"`
	DdnsUseOption81                  bool                        `json:"ddns_use_option81,omitempty"`
	DeleteReason                     string                      `json:"delete_reason,omitempty"`
	DenyBootp                        bool                        `json:"deny_bootp,omitempty"`
	DiscoverNowStatus                string                      `json:"discover_now_status,omitempty"`
	DiscoveryBasicPollSettings       *DiscoveryBasicpollsettings `json:"discovery_basic_poll_settings,omitempty"`
	DiscoveryBlackoutSetting         *PropertiesBlackoutsetting  `json:"discovery_blackout_setting,omitempty"`
	DiscoveryEngineType              string                      `json:"discovery_engine_type,omitempty"`
	DiscoveryMember                  string                      `json:"discovery_member,omitempty"`
	EmailList                        []string                    `json:"email_list,omitempty"`
	EnableDdns                       bool                        `json:"enable_ddns,omitempty"`
	EnableDhcpThresholds             bool                        `json:"enable_dhcp_thresholds,omitempty"`
	EnableDiscovery                  bool                        `json:"enable_discovery,omitempty"`
	EnableEmailWarnings              bool                        `json:"enable_email_warnings,omitempty"`
	EnableImmediateDiscovery         bool                        `json:"enable_immediate_discovery,omitempty"`
	EnablePxeLeaseTime               bool                        `json:"enable_pxe_lease_time,omitempty"`
	EnableSnmpWarnings               bool                        `json:"enable_snmp_warnings,omitempty"`
	EndpointSources                  []*CiscoiseEndpoint         `json:"endpoint_sources,omitempty"`
	Ea                               EA                          `json:"extattrs,omitempty"`
	HighWaterMark                    uint32                      `json:"high_water_mark,omitempty"`
	HighWaterMarkReset               uint32                      `json:"high_water_mark_reset,omitempty"`
	IgnoreDhcpOptionListRequest      bool                        `json:"ignore_dhcp_option_list_request,omitempty"`
	IgnoreId                         string                      `json:"ignore_id,omitempty"`
	IgnoreMacAddresses               []string                    `json:"ignore_mac_addresses,omitempty"`
	IpamEmailAddresses               []string                    `json:"ipam_email_addresses,omitempty"`
	IpamThresholdSettings            *SettingIpamThreshold       `json:"ipam_threshold_settings,omitempty"`
	IpamTrapSettings                 *SettingIpamTrap            `json:"ipam_trap_settings,omitempty"`
	LastRirRegistrationUpdateSent    *time.Time                  `json:"last_rir_registration_update_sent,omitempty"`
	LastRirRegistrationUpdateStatus  string                      `json:"last_rir_registration_update_status,omitempty"`
	LeaseScavengeTime                int                         `json:"lease_scavenge_time,omitempty"`
	LogicFilterRules                 []*Logicfilterrule          `json:"logic_filter_rules,omitempty"`
	LowWaterMark                     uint32                      `json:"low_water_mark,omitempty"`
	LowWaterMarkReset                uint32                      `json:"low_water_mark_reset,omitempty"`
	MgmPrivate                       bool                        `json:"mgm_private,omitempty"`
	MgmPrivateOverridable            bool                        `json:"mgm_private_overridable,omitempty"`
	MsAdUserData                     *MsserverAduserData         `json:"ms_ad_user_data,omitempty"`
	Network                          string                      `json:"network,omitempty"`
	NetworkContainer                 string                      `json:"network_container,omitempty"`
	NetworkView                      string                      `json:"network_view,omitempty"`
	NextAvailableNetwork             *Nextavailablenet           `json:"next_available_network,omitempty"`
	Nextserver                       string                      `json:"nextserver,omitempty"`
	Options                          []*Dhcpoption               `json:"options,omitempty"`
	PortControlBlackoutSetting       *PropertiesBlackoutsetting  `json:"port_control_blackout_setting,omitempty"`
	PxeLeaseTime                     uint32                      `json:"pxe_lease_time,omitempty"`
	RecycleLeases                    bool                        `json:"recycle_leases,omitempty"`
	RemoveSubnets                    bool                        `json:"remove_subnets,omitempty"`
	Resize                           *Resizenetwork              `json:"resize,omitempty"`
	RestartIfNeeded                  bool                        `json:"restart_if_needed,omitempty"`
	Rir                              string                      `json:"rir,omitempty"`
	RirOrganization                  string                      `json:"rir_organization,omitempty"`
	RirRegistrationAction            string                      `json:"rir_registration_action,omitempty"`
	RirRegistrationStatus            string                      `json:"rir_registration_status,omitempty"`
	SamePortControlDiscoveryBlackout bool                        `json:"same_port_control_discovery_blackout,omitempty"`
	SendRirRequest                   bool                        `json:"send_rir_request,omitempty"`
	SubscribeSettings                *CiscoiseSubscribesetting   `json:"subscribe_settings,omitempty"`
	Unmanaged                        bool                        `json:"unmanaged,omitempty"`
	UpdateDnsOnLeaseRenewal          bool                        `json:"update_dns_on_lease_renewal,omitempty"`
	UseAuthority                     bool                        `json:"use_authority,omitempty"`
	UseBlackoutSetting               bool                        `json:"use_blackout_setting,omitempty"`
	UseBootfile                      bool                        `json:"use_bootfile,omitempty"`
	UseBootserver                    bool                        `json:"use_bootserver,omitempty"`
	UseDdnsDomainname                bool                        `json:"use_ddns_domainname,omitempty"`
	UseDdnsGenerateHostname          bool                        `json:"use_ddns_generate_hostname,omitempty"`
	UseDdnsTtl                       bool                        `json:"use_ddns_ttl,omitempty"`
	UseDdnsUpdateFixedAddresses      bool                        `json:"use_ddns_update_fixed_addresses,omitempty"`
	UseDdnsUseOption81               bool                        `json:"use_ddns_use_option81,omitempty"`
	UseDenyBootp                     bool                        `json:"use_deny_bootp,omitempty"`
	UseDiscoveryBasicPollingSettings bool                        `json:"use_discovery_basic_polling_settings,omitempty"`
	UseEmailList                     bool                        `json:"use_email_list,omitempty"`
	UseEnableDdns                    bool                        `json:"use_enable_ddns,omitempty"`
	UseEnableDhcpThresholds          bool                        `json:"use_enable_dhcp_thresholds,omitempty"`
	UseEnableDiscovery               bool                        `json:"use_enable_discovery,omitempty"`
	UseIgnoreDhcpOptionListRequest   bool                        `json:"use_ignore_dhcp_option_list_request,omitempty"`
	UseIgnoreId                      bool                        `json:"use_ignore_id,omitempty"`
	UseIpamEmailAddresses            bool                        `json:"use_ipam_email_addresses,omitempty"`
	UseIpamThresholdSettings         bool                        `json:"use_ipam_threshold_settings,omitempty"`
	UseIpamTrapSettings              bool                        `json:"use_ipam_trap_settings,omitempty"`
	UseLeaseScavengeTime             bool                        `json:"use_lease_scavenge_time,omitempty"`
	UseLogicFilterRules              bool                        `json:"use_logic_filter_rules,omitempty"`
	UseMgmPrivate                    bool                        `json:"use_mgm_private,omitempty"`
	UseNextserver                    bool                        `json:"use_nextserver,omitempty"`
	UseOptions                       bool                        `json:"use_options,omitempty"`
	UsePxeLeaseTime                  bool                        `json:"use_pxe_lease_time,omitempty"`
	UseRecycleLeases                 bool                        `json:"use_recycle_leases,omitempty"`
	UseSubscribeSettings             bool                        `json:"use_subscribe_settings,omitempty"`
	UseUpdateDnsOnLeaseRenewal       bool                        `json:"use_update_dns_on_lease_renewal,omitempty"`
	UseZoneAssociations              bool                        `json:"use_zone_associations,omitempty"`
	Utilization                      uint32                      `json:"utilization,omitempty"`
	ZoneAssociations                 []*Zoneassociation          `json:"zone_associations,omitempty"`
}

func (Ipv4NetworkContainer) ObjectType() string {
	return "networkcontainer"
}

func (obj Ipv4NetworkContainer) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "network", "network_view"}
	}
	return obj.returnFields
}

// Ipv4Network represents Infoblox object network
type Ipv4Network struct {
	IBBase                           `json:"-"`
	Ref                              string                      `json:"_ref,omitempty"`
	Authority                        bool                        `json:"authority,omitempty"`
	AutoCreateReversezone            bool                        `json:"auto_create_reversezone,omitempty"`
	Bootfile                         string                      `json:"bootfile,omitempty"`
	Bootserver                       string                      `json:"bootserver,omitempty"`
	CloudInfo                        *GridCloudapiInfo           `json:"cloud_info,omitempty"`
	Comment                          string                      `json:"comment,omitempty"`
	ConflictCount                    uint32                      `json:"conflict_count,omitempty"`
	DdnsDomainname                   string                      `json:"ddns_domainname,omitempty"`
	DdnsGenerateHostname             bool                        `json:"ddns_generate_hostname,omitempty"`
	DdnsServerAlwaysUpdates          bool                        `json:"ddns_server_always_updates,omitempty"`
	DdnsTtl                          uint32                      `json:"ddns_ttl,omitempty"`
	DdnsUpdateFixedAddresses         bool                        `json:"ddns_update_fixed_addresses,omitempty"`
	DdnsUseOption81                  bool                        `json:"ddns_use_option81,omitempty"`
	DeleteReason                     string                      `json:"delete_reason,omitempty"`
	DenyBootp                        bool                        `json:"deny_bootp,omitempty"`
	DhcpUtilization                  uint32                      `json:"dhcp_utilization,omitempty"`
	DhcpUtilizationStatus            string                      `json:"dhcp_utilization_status,omitempty"`
	Disable                          bool                        `json:"disable,omitempty"`
	DiscoverNowStatus                string                      `json:"discover_now_status,omitempty"`
	DiscoveredBgpAs                  string                      `json:"discovered_bgp_as,omitempty"`
	DiscoveredBridgeDomain           string                      `json:"discovered_bridge_domain,omitempty"`
	DiscoveredTenant                 string                      `json:"discovered_tenant,omitempty"`
	DiscoveredVlanId                 string                      `json:"discovered_vlan_id,omitempty"`
	DiscoveredVlanName               string                      `json:"discovered_vlan_name,omitempty"`
	DiscoveredVrfDescription         string                      `json:"discovered_vrf_description,omitempty"`
	DiscoveredVrfName                string                      `json:"discovered_vrf_name,omitempty"`
	DiscoveredVrfRd                  string                      `json:"discovered_vrf_rd,omitempty"`
	DiscoveryBasicPollSettings       *DiscoveryBasicpollsettings `json:"discovery_basic_poll_settings,omitempty"`
	DiscoveryBlackoutSetting         *PropertiesBlackoutsetting  `json:"discovery_blackout_setting,omitempty"`
	DiscoveryEngineType              string                      `json:"discovery_engine_type,omitempty"`
	DiscoveryMember                  string                      `json:"discovery_member,omitempty"`
	DynamicHosts                     uint32                      `json:"dynamic_hosts,omitempty"`
	EmailList                        []string                    `json:"email_list,omitempty"`
	EnableDdns                       bool                        `json:"enable_ddns,omitempty"`
	EnableDhcpThresholds             bool                        `json:"enable_dhcp_thresholds,omitempty"`
	EnableDiscovery                  bool                        `json:"enable_discovery,omitempty"`
	EnableEmailWarnings              bool                        `json:"enable_email_warnings,omitempty"`
	EnableIfmapPublishing            bool                        `json:"enable_ifmap_publishing,omitempty"`
	EnableImmediateDiscovery         bool                        `json:"enable_immediate_discovery,omitempty"`
	EnablePxeLeaseTime               bool                        `json:"enable_pxe_lease_time,omitempty"`
	EnableSnmpWarnings               bool                        `json:"enable_snmp_warnings,omitempty"`
	EndpointSources                  []*CiscoiseEndpoint         `json:"endpoint_sources,omitempty"`
	ExpandNetwork                    *Expandnetwork              `json:"expand_network,omitempty"`
	Ea                               EA                          `json:"extattrs,omitempty"`
	HighWaterMark                    uint32                      `json:"high_water_mark,omitempty"`
	HighWaterMarkReset               uint32                      `json:"high_water_mark_reset,omitempty"`
	IgnoreDhcpOptionListRequest      bool                        `json:"ignore_dhcp_option_list_request,omitempty"`
	IgnoreId                         string                      `json:"ignore_id,omitempty"`
	IgnoreMacAddresses               []string                    `json:"ignore_mac_addresses,omitempty"`
	IpamEmailAddresses               []string                    `json:"ipam_email_addresses,omitempty"`
	IpamThresholdSettings            *SettingIpamThreshold       `json:"ipam_threshold_settings,omitempty"`
	IpamTrapSettings                 *SettingIpamTrap            `json:"ipam_trap_settings,omitempty"`
	Ipv4Addr                         string                      `json:"ipv4addr,omitempty"`
	LastRirRegistrationUpdateSent    *time.Time                  `json:"last_rir_registration_update_sent,omitempty"`
	LastRirRegistrationUpdateStatus  string                      `json:"last_rir_registration_update_status,omitempty"`
	LeaseScavengeTime                int                         `json:"lease_scavenge_time,omitempty"`
	LogicFilterRules                 []*Logicfilterrule          `json:"logic_filter_rules,omitempty"`
	LowWaterMark                     uint32                      `json:"low_water_mark,omitempty"`
	LowWaterMarkReset                uint32                      `json:"low_water_mark_reset,omitempty"`
	Members                          []*Msdhcpserver             `json:"members,omitempty"`
	MgmPrivate                       bool                        `json:"mgm_private,omitempty"`
	MgmPrivateOverridable            bool                        `json:"mgm_private_overridable,omitempty"`
	MsAdUserData                     *MsserverAduserData         `json:"ms_ad_user_data,omitempty"`
	Netmask                          uint32                      `json:"netmask,omitempty"`
	Network                          string                      `json:"network,omitempty"`
	NetworkContainer                 string                      `json:"network_container,omitempty"`
	NetworkView                      string                      `json:"network_view,omitempty"`
	NextAvailableIp                  *Nextavailableip            `json:"next_available_ip,omitempty"`
	NextAvailableNetwork             *Nextavailablenet           `json:"next_available_network,omitempty"`
	NextAvailableVlan                *Nextavailablevlan          `json:"next_available_vlan,omitempty"`
	Nextserver                       string                      `json:"nextserver,omitempty"`
	Options                          []*Dhcpoption               `json:"options,omitempty"`
	PortControlBlackoutSetting       *PropertiesBlackoutsetting  `json:"port_control_blackout_setting,omitempty"`
	PxeLeaseTime                     uint32                      `json:"pxe_lease_time,omitempty"`
	RecycleLeases                    bool                        `json:"recycle_leases,omitempty"`
	Resize                           *Resizenetwork              `json:"resize,omitempty"`
	RestartIfNeeded                  bool                        `json:"restart_if_needed,omitempty"`
	Rir                              string                      `json:"rir,omitempty"`
	RirOrganization                  string                      `json:"rir_organization,omitempty"`
	RirRegistrationAction            string                      `json:"rir_registration_action,omitempty"`
	RirRegistrationStatus            string                      `json:"rir_registration_status,omitempty"`
	SamePortControlDiscoveryBlackout bool                        `json:"same_port_control_discovery_blackout,omitempty"`
	SendRirRequest                   bool                        `json:"send_rir_request,omitempty"`
	SplitNetwork                     *Splitnetwork               `json:"split_network,omitempty"`
	StaticHosts                      uint32                      `json:"static_hosts,omitempty"`
	SubscribeSettings                *CiscoiseSubscribesetting   `json:"subscribe_settings,omitempty"`
	Template                         string                      `json:"template,omitempty"`
	TotalHosts                       uint32                      `json:"total_hosts,omitempty"`
	Unmanaged                        bool                        `json:"unmanaged,omitempty"`
	UnmanagedCount                   uint32                      `json:"unmanaged_count,omitempty"`
	UpdateDnsOnLeaseRenewal          bool                        `json:"update_dns_on_lease_renewal,omitempty"`
	UseAuthority                     bool                        `json:"use_authority,omitempty"`
	UseBlackoutSetting               bool                        `json:"use_blackout_setting,omitempty"`
	UseBootfile                      bool                        `json:"use_bootfile,omitempty"`
	UseBootserver                    bool                        `json:"use_bootserver,omitempty"`
	UseDdnsDomainname                bool                        `json:"use_ddns_domainname,omitempty"`
	UseDdnsGenerateHostname          bool                        `json:"use_ddns_generate_hostname,omitempty"`
	UseDdnsTtl                       bool                        `json:"use_ddns_ttl,omitempty"`
	UseDdnsUpdateFixedAddresses      bool                        `json:"use_ddns_update_fixed_addresses,omitempty"`
	UseDdnsUseOption81               bool                        `json:"use_ddns_use_option81,omitempty"`
	UseDenyBootp                     bool                        `json:"use_deny_bootp,omitempty"`
	UseDiscoveryBasicPollingSettings bool                        `json:"use_discovery_basic_polling_settings,omitempty"`
	UseEmailList                     bool                        `json:"use_email_list,omitempty"`
	UseEnableDdns                    bool                        `json:"use_enable_ddns,omitempty"`
	UseEnableDhcpThresholds          bool                        `json:"use_enable_dhcp_thresholds,omitempty"`
	UseEnableDiscovery               bool                        `json:"use_enable_discovery,omitempty"`
	UseEnableIfmapPublishing         bool                        `json:"use_enable_ifmap_publishing,omitempty"`
	UseIgnoreDhcpOptionListRequest   bool                        `json:"use_ignore_dhcp_option_list_request,omitempty"`
	UseIgnoreId                      bool                        `json:"use_ignore_id,omitempty"`
	UseIpamEmailAddresses            bool                        `json:"use_ipam_email_addresses,omitempty"`
	UseIpamThresholdSettings         bool                        `json:"use_ipam_threshold_settings,omitempty"`
	UseIpamTrapSettings              bool                        `json:"use_ipam_trap_settings,omitempty"`
	UseLeaseScavengeTime             bool                        `json:"use_lease_scavenge_time,omitempty"`
	UseLogicFilterRules              bool                        `json:"use_logic_filter_rules,omitempty"`
	UseMgmPrivate                    bool                        `json:"use_mgm_private,omitempty"`
	UseNextserver                    bool                        `json:"use_nextserver,omitempty"`
	UseOptions                       bool                        `json:"use_options,omitempty"`
	UsePxeLeaseTime                  bool                        `json:"use_pxe_lease_time,omitempty"`
	UseRecycleLeases                 bool                        `json:"use_recycle_leases,omitempty"`
	UseSubscribeSettings             bool                        `json:"use_subscribe_settings,omitempty"`
	UseUpdateDnsOnLeaseRenewal       bool                        `json:"use_update_dns_on_lease_renewal,omitempty"`
	UseZoneAssociations              bool                        `json:"use_zone_associations,omitempty"`
	Utilization                      uint32                      `json:"utilization,omitempty"`
	UtilizationUpdate                *time.Time                  `json:"utilization_update,omitempty"`
	Vlans                            []*Vlanlink                 `json:"vlans,omitempty"`
	ZoneAssociations                 []*Zoneassociation          `json:"zone_associations,omitempty"`
}

func (Ipv4Network) ObjectType() string {
	return "network"
}

func (obj Ipv4Network) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "network", "network_view"}
	}
	return obj.returnFields
}

// Nsgroup represents Infoblox object nsgroup
type Nsgroup struct {
	IBBase              `json:"-"`
	Ref                 string          `json:"_ref,omitempty"`
	Comment             string          `json:"comment,omitempty"`
	Ea                  EA              `json:"extattrs,omitempty"`
	ExternalPrimaries   []NameServer    `json:"external_primaries,omitempty"`
	ExternalSecondaries []NameServer    `json:"external_secondaries,omitempty"`
	GridPrimary         []*Memberserver `json:"grid_primary,omitempty"`
	GridSecondaries     []*Memberserver `json:"grid_secondaries,omitempty"`
	IsGridDefault       bool            `json:"is_grid_default,omitempty"`
	IsMultimaster       bool            `json:"is_multimaster,omitempty"`
	Name                string          `json:"name,omitempty"`
	UseExternalPrimary  bool            `json:"use_external_primary,omitempty"`
}

func (Nsgroup) ObjectType() string {
	return "nsgroup"
}

func (obj Nsgroup) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name"}
	}
	return obj.returnFields
}

// NsgroupForwardstubserver represents Infoblox object nsgroup:forwardstubserver
type NsgroupForwardstubserver struct {
	IBBase          `json:"-"`
	Ref             string       `json:"_ref,omitempty"`
	Comment         string       `json:"comment,omitempty"`
	Ea              EA           `json:"extattrs,omitempty"`
	ExternalServers []NameServer `json:"external_servers,omitempty"`
	Name            string       `json:"name,omitempty"`
}

func (NsgroupForwardstubserver) ObjectType() string {
	return "nsgroup:forwardstubserver"
}

func (obj NsgroupForwardstubserver) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"external_servers", "name"}
	}
	return obj.returnFields
}

// NsgroupStubmember represents Infoblox object nsgroup:stubmember
type NsgroupStubmember struct {
	IBBase      `json:"-"`
	Ref         string          `json:"_ref,omitempty"`
	Comment     string          `json:"comment,omitempty"`
	Ea          EA              `json:"extattrs,omitempty"`
	Name        string          `json:"name,omitempty"`
	StubMembers []*Memberserver `json:"stub_members,omitempty"`
}

func (NsgroupStubmember) ObjectType() string {
	return "nsgroup:stubmember"
}

func (obj NsgroupStubmember) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name"}
	}
	return obj.returnFields
}

// NsgroupForwardingmember represents Infoblox object nsgroup:forwardingmember
type NsgroupForwardingmember struct {
	IBBase            `json:"-"`
	Ref               string                    `json:"_ref,omitempty"`
	Comment           string                    `json:"comment,omitempty"`
	Ea                EA                        `json:"extattrs,omitempty"`
	ForwardingServers []*Forwardingmemberserver `json:"forwarding_servers,omitempty"`
	Name              string                    `json:"name,omitempty"`
}

func (NsgroupForwardingmember) ObjectType() string {
	return "nsgroup:forwardingmember"
}

func (obj NsgroupForwardingmember) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"forwarding_servers", "name"}
	}
	return obj.returnFields
}

// Orderedranges represents Infoblox object orderedranges
type Orderedranges struct {
	IBBase  `json:"-"`
	Ref     string   `json:"_ref,omitempty"`
	Network string   `json:"network,omitempty"`
	Ranges  []*Range `json:"ranges,omitempty"`
}

func (Orderedranges) ObjectType() string {
	return "orderedranges"
}

func (obj Orderedranges) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"network", "ranges"}
	}
	return obj.returnFields
}

// Orderedresponsepolicyzones represents Infoblox object orderedresponsepolicyzones
type Orderedresponsepolicyzones struct {
	IBBase  `json:"-"`
	Ref     string   `json:"_ref,omitempty"`
	RpZones []string `json:"rp_zones,omitempty"`
	View    string   `json:"view,omitempty"`
}

func (Orderedresponsepolicyzones) ObjectType() string {
	return "orderedresponsepolicyzones"
}

func (obj Orderedresponsepolicyzones) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"view"}
	}
	return obj.returnFields
}

// ParentalcontrolAvp represents Infoblox object parentalcontrol:avp
type ParentalcontrolAvp struct {
	IBBase       `json:"-"`
	Ref          string   `json:"_ref,omitempty"`
	Comment      string   `json:"comment,omitempty"`
	DomainTypes  []string `json:"domain_types,omitempty"`
	IsRestricted bool     `json:"is_restricted,omitempty"`
	Name         string   `json:"name,omitempty"`
	Type         uint32   `json:"type,omitempty"`
	UserDefined  bool     `json:"user_defined,omitempty"`
	ValueType    string   `json:"value_type,omitempty"`
	VendorId     uint32   `json:"vendor_id,omitempty"`
	VendorType   uint32   `json:"vendor_type,omitempty"`
}

func (ParentalcontrolAvp) ObjectType() string {
	return "parentalcontrol:avp"
}

func (obj ParentalcontrolAvp) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name", "type", "value_type"}
	}
	return obj.returnFields
}

// OutboundCloudclient represents Infoblox object outbound:cloudclient
type OutboundCloudclient struct {
	IBBase                    `json:"-"`
	Ref                       string                      `json:"_ref,omitempty"`
	Enable                    bool                        `json:"enable,omitempty"`
	GridMember                string                      `json:"grid_member,omitempty"`
	Interval                  uint32                      `json:"interval,omitempty"`
	OutboundCloudClientEvents []*OutboundCloudclientEvent `json:"outbound_cloud_client_events,omitempty"`
}

func (OutboundCloudclient) ObjectType() string {
	return "outbound:cloudclient"
}

func (obj OutboundCloudclient) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"enable", "interval"}
	}
	return obj.returnFields
}

// ParentalcontrolIpspacediscriminator represents Infoblox object parentalcontrol:ipspacediscriminator
type ParentalcontrolIpspacediscriminator struct {
	IBBase `json:"-"`
	Ref    string `json:"_ref,omitempty"`
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
}

func (ParentalcontrolIpspacediscriminator) ObjectType() string {
	return "parentalcontrol:ipspacediscriminator"
}

func (obj ParentalcontrolIpspacediscriminator) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name", "value"}
	}
	return obj.returnFields
}

// ParentalcontrolBlockingpolicy represents Infoblox object parentalcontrol:blockingpolicy
type ParentalcontrolBlockingpolicy struct {
	IBBase `json:"-"`
	Ref    string `json:"_ref,omitempty"`
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
}

func (ParentalcontrolBlockingpolicy) ObjectType() string {
	return "parentalcontrol:blockingpolicy"
}

func (obj ParentalcontrolBlockingpolicy) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name", "value"}
	}
	return obj.returnFields
}

// PxgridEndpoint represents Infoblox object pxgrid:endpoint
type PxgridEndpoint struct {
	IBBase                     `json:"-"`
	Ref                        string                            `json:"_ref,omitempty"`
	Address                    string                            `json:"address,omitempty"`
	ClientCertificateSubject   string                            `json:"client_certificate_subject,omitempty"`
	ClientCertificateToken     string                            `json:"client_certificate_token,omitempty"`
	ClientCertificateValidFrom *time.Time                        `json:"client_certificate_valid_from,omitempty"`
	ClientCertificateValidTo   *time.Time                        `json:"client_certificate_valid_to,omitempty"`
	Comment                    string                            `json:"comment,omitempty"`
	Disable                    bool                              `json:"disable,omitempty"`
	Ea                         EA                                `json:"extattrs,omitempty"`
	LogLevel                   string                            `json:"log_level,omitempty"`
	Name                       string                            `json:"name,omitempty"`
	NetworkView                string                            `json:"network_view,omitempty"`
	OutboundMemberType         string                            `json:"outbound_member_type,omitempty"`
	OutboundMembers            []string                          `json:"outbound_members,omitempty"`
	PublishSettings            *CiscoisePublishsetting           `json:"publish_settings,omitempty"`
	SubscribeSettings          *CiscoiseSubscribesetting         `json:"subscribe_settings,omitempty"`
	TemplateInstance           *NotificationRestTemplateinstance `json:"template_instance,omitempty"`
	TestConnection             *Testendpointconnection           `json:"test_connection,omitempty"`
	Timeout                    uint32                            `json:"timeout,omitempty"`
	VendorIdentifier           string                            `json:"vendor_identifier,omitempty"`
	WapiUserName               string                            `json:"wapi_user_name,omitempty"`
	WapiUserPassword           string                            `json:"wapi_user_password,omitempty"`
}

func (PxgridEndpoint) ObjectType() string {
	return "pxgrid:endpoint"
}

func (obj PxgridEndpoint) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"address", "disable", "name", "outbound_member_type"}
	}
	return obj.returnFields
}

// ParentalcontrolSubscriber represents Infoblox object parentalcontrol:subscriber
type ParentalcontrolSubscriber struct {
	IBBase                       `json:"-"`
	Ref                          string   `json:"_ref,omitempty"`
	AltSubscriberId              string   `json:"alt_subscriber_id,omitempty"`
	AltSubscriberIdRegexp        string   `json:"alt_subscriber_id_regexp,omitempty"`
	AltSubscriberIdSubexpression uint32   `json:"alt_subscriber_id_subexpression,omitempty"`
	Ancillaries                  []string `json:"ancillaries,omitempty"`
	CatAcctname                  string   `json:"cat_acctname,omitempty"`
	CatPassword                  string   `json:"cat_password,omitempty"`
	CatUpdateFrequency           uint32   `json:"cat_update_frequency,omitempty"`
	CategoryUrl                  string   `json:"category_url,omitempty"`
	EnableMgmtOnlyNas            bool     `json:"enable_mgmt_only_nas,omitempty"`
	EnableParentalControl        bool     `json:"enable_parental_control,omitempty"`
	Ident                        string   `json:"ident,omitempty"`
	InterimAccountingInterval    uint32   `json:"interim_accounting_interval,omitempty"`
	IpAnchors                    []string `json:"ip_anchors,omitempty"`
	IpSpaceDiscRegexp            string   `json:"ip_space_disc_regexp,omitempty"`
	IpSpaceDiscSubexpression     uint32   `json:"ip_space_disc_subexpression,omitempty"`
	IpSpaceDiscriminator         string   `json:"ip_space_discriminator,omitempty"`
	LocalId                      string   `json:"local_id,omitempty"`
	LocalIdRegexp                string   `json:"local_id_regexp,omitempty"`
	LocalIdSubexpression         uint32   `json:"local_id_subexpression,omitempty"`
	LogGuestLookups              bool     `json:"log_guest_lookups,omitempty"`
	NasContextInfo               string   `json:"nas_context_info,omitempty"`
	PcZoneName                   string   `json:"pc_zone_name,omitempty"`
	ProxyPassword                string   `json:"proxy_password,omitempty"`
	ProxyUrl                     string   `json:"proxy_url,omitempty"`
	ProxyUsername                string   `json:"proxy_username,omitempty"`
	SubscriberId                 string   `json:"subscriber_id,omitempty"`
	SubscriberIdRegexp           string   `json:"subscriber_id_regexp,omitempty"`
	SubscriberIdSubexpression    uint32   `json:"subscriber_id_subexpression,omitempty"`
}

func (ParentalcontrolSubscriber) ObjectType() string {
	return "parentalcontrol:subscriber"
}

func (obj ParentalcontrolSubscriber) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"alt_subscriber_id", "local_id", "subscriber_id"}
	}
	return obj.returnFields
}

// ParentalcontrolSubscribersite represents Infoblox object parentalcontrol:subscribersite
type ParentalcontrolSubscribersite struct {
	IBBase             `json:"-"`
	Ref                string                       `json:"_ref,omitempty"`
	Abss               []*ParentalcontrolAbs        `json:"abss,omitempty"`
	BlockSize          uint32                       `json:"block_size,omitempty"`
	BlockingIpv4Vip1   string                       `json:"blocking_ipv4_vip1,omitempty"`
	BlockingIpv4Vip2   string                       `json:"blocking_ipv4_vip2,omitempty"`
	BlockingIpv6Vip1   string                       `json:"blocking_ipv6_vip1,omitempty"`
	BlockingIpv6Vip2   string                       `json:"blocking_ipv6_vip2,omitempty"`
	Comment            string                       `json:"comment,omitempty"`
	Ea                 EA                           `json:"extattrs,omitempty"`
	FirstPort          uint32                       `json:"first_port,omitempty"`
	MaximumSubscribers uint32                       `json:"maximum_subscribers,omitempty"`
	Members            []*ParentalcontrolSitemember `json:"members,omitempty"`
	Msps               []*ParentalcontrolMsp        `json:"msps,omitempty"`
	Name               string                       `json:"name,omitempty"`
	NasGateways        []*ParentalcontrolNasgateway `json:"nas_gateways,omitempty"`
	NasPort            uint32                       `json:"nas_port,omitempty"`
	Spms               []*ParentalcontrolSpm        `json:"spms,omitempty"`
	StrictNat          bool                         `json:"strict_nat,omitempty"`
}

func (ParentalcontrolSubscribersite) ObjectType() string {
	return "parentalcontrol:subscribersite"
}

func (obj ParentalcontrolSubscribersite) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"block_size", "first_port", "name", "strict_nat"}
	}
	return obj.returnFields
}

// ParentalcontrolSubscriberrecord represents Infoblox object parentalcontrol:subscriberrecord
type ParentalcontrolSubscriberrecord struct {
	IBBase                 `json:"-"`
	Ref                    string `json:"_ref,omitempty"`
	AccountingSessionId    string `json:"accounting_session_id,omitempty"`
	AltIpAddr              string `json:"alt_ip_addr,omitempty"`
	Ans0                   string `json:"ans0,omitempty"`
	Ans1                   string `json:"ans1,omitempty"`
	Ans2                   string `json:"ans2,omitempty"`
	Ans3                   string `json:"ans3,omitempty"`
	Ans4                   string `json:"ans4,omitempty"`
	BlackList              string `json:"black_list,omitempty"`
	Bwflag                 bool   `json:"bwflag,omitempty"`
	DynamicCategoryPolicy  bool   `json:"dynamic_category_policy,omitempty"`
	Flags                  string `json:"flags,omitempty"`
	IpAddr                 string `json:"ip_addr,omitempty"`
	Ipsd                   string `json:"ipsd,omitempty"`
	Localid                string `json:"localid,omitempty"`
	NasContextual          string `json:"nas_contextual,omitempty"`
	ParentalControlPolicy  string `json:"parental_control_policy,omitempty"`
	Prefix                 uint32 `json:"prefix,omitempty"`
	ProxyAll               bool   `json:"proxy_all,omitempty"`
	Site                   string `json:"site,omitempty"`
	SubscriberId           string `json:"subscriber_id,omitempty"`
	SubscriberSecurePolicy string `json:"subscriber_secure_policy,omitempty"`
	UnknownCategoryPolicy  bool   `json:"unknown_category_policy,omitempty"`
	WhiteList              string `json:"white_list,omitempty"`
	WpcCategoryPolicy      string `json:"wpc_category_policy,omitempty"`
}

func (ParentalcontrolSubscriberrecord) ObjectType() string {
	return "parentalcontrol:subscriberrecord"
}

func (obj ParentalcontrolSubscriberrecord) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"accounting_session_id", "ip_addr", "ipsd", "localid", "prefix", "site", "subscriber_id"}
	}
	return obj.returnFields
}

// Permission represents Infoblox object permission
type Permission struct {
	IBBase       `json:"-"`
	Ref          string `json:"_ref,omitempty"`
	Group        string `json:"group,omitempty"`
	Object       string `json:"object,omitempty"`
	Permission   string `json:"permission,omitempty"`
	ResourceType string `json:"resource_type,omitempty"`
	Role         string `json:"role,omitempty"`
}

func (Permission) ObjectType() string {
	return "permission"
}

func (obj Permission) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"group", "permission", "resource_type", "role"}
	}
	return obj.returnFields
}

// RadiusAuthservice represents Infoblox object radius:authservice
type RadiusAuthservice struct {
	IBBase                    `json:"-"`
	Ref                       string                     `json:"_ref,omitempty"`
	AcctRetries               uint32                     `json:"acct_retries,omitempty"`
	AcctTimeout               uint32                     `json:"acct_timeout,omitempty"`
	AuthRetries               uint32                     `json:"auth_retries,omitempty"`
	AuthTimeout               uint32                     `json:"auth_timeout,omitempty"`
	CacheTtl                  uint32                     `json:"cache_ttl,omitempty"`
	CheckRadiusServerSettings *Checkradiusserversettings `json:"check_radius_server_settings,omitempty"`
	Comment                   string                     `json:"comment,omitempty"`
	Disable                   bool                       `json:"disable,omitempty"`
	EnableCache               bool                       `json:"enable_cache,omitempty"`
	Mode                      string                     `json:"mode,omitempty"`
	Name                      string                     `json:"name,omitempty"`
	RecoveryInterval          uint32                     `json:"recovery_interval,omitempty"`
	Servers                   []*RadiusServer            `json:"servers,omitempty"`
}

func (RadiusAuthservice) ObjectType() string {
	return "radius:authservice"
}

func (obj RadiusAuthservice) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "disable", "name"}
	}
	return obj.returnFields
}

// Rangetemplate represents Infoblox object rangetemplate
type Rangetemplate struct {
	IBBase                         `json:"-"`
	Ref                            string                    `json:"_ref,omitempty"`
	Bootfile                       string                    `json:"bootfile,omitempty"`
	Bootserver                     string                    `json:"bootserver,omitempty"`
	CloudApiCompatible             bool                      `json:"cloud_api_compatible,omitempty"`
	Comment                        string                    `json:"comment,omitempty"`
	DdnsDomainname                 string                    `json:"ddns_domainname,omitempty"`
	DdnsGenerateHostname           bool                      `json:"ddns_generate_hostname,omitempty"`
	DelegatedMember                *Dhcpmember               `json:"delegated_member,omitempty"`
	DenyAllClients                 bool                      `json:"deny_all_clients,omitempty"`
	DenyBootp                      bool                      `json:"deny_bootp,omitempty"`
	EmailList                      []string                  `json:"email_list,omitempty"`
	EnableDdns                     bool                      `json:"enable_ddns,omitempty"`
	EnableDhcpThresholds           bool                      `json:"enable_dhcp_thresholds,omitempty"`
	EnableEmailWarnings            bool                      `json:"enable_email_warnings,omitempty"`
	EnablePxeLeaseTime             bool                      `json:"enable_pxe_lease_time,omitempty"`
	EnableSnmpWarnings             bool                      `json:"enable_snmp_warnings,omitempty"`
	Exclude                        []*Exclusionrangetemplate `json:"exclude,omitempty"`
	Ea                             EA                        `json:"extattrs,omitempty"`
	FailoverAssociation            string                    `json:"failover_association,omitempty"`
	FingerprintFilterRules         []*Filterrule             `json:"fingerprint_filter_rules,omitempty"`
	HighWaterMark                  uint32                    `json:"high_water_mark,omitempty"`
	HighWaterMarkReset             uint32                    `json:"high_water_mark_reset,omitempty"`
	IgnoreDhcpOptionListRequest    bool                      `json:"ignore_dhcp_option_list_request,omitempty"`
	KnownClients                   string                    `json:"known_clients,omitempty"`
	LeaseScavengeTime              int                       `json:"lease_scavenge_time,omitempty"`
	LogicFilterRules               []*Logicfilterrule        `json:"logic_filter_rules,omitempty"`
	LowWaterMark                   uint32                    `json:"low_water_mark,omitempty"`
	LowWaterMarkReset              uint32                    `json:"low_water_mark_reset,omitempty"`
	MacFilterRules                 []*Filterrule             `json:"mac_filter_rules,omitempty"`
	Member                         *Dhcpmember               `json:"member,omitempty"`
	MsOptions                      []*Msdhcpoption           `json:"ms_options,omitempty"`
	MsServer                       *Msdhcpserver             `json:"ms_server,omitempty"`
	NacFilterRules                 []*Filterrule             `json:"nac_filter_rules,omitempty"`
	Name                           string                    `json:"name,omitempty"`
	Nextserver                     string                    `json:"nextserver,omitempty"`
	NumberOfAddresses              uint32                    `json:"number_of_addresses,omitempty"`
	Offset                         uint32                    `json:"offset,omitempty"`
	OptionFilterRules              []*Filterrule             `json:"option_filter_rules,omitempty"`
	Options                        []*Dhcpoption             `json:"options,omitempty"`
	PxeLeaseTime                   uint32                    `json:"pxe_lease_time,omitempty"`
	RecycleLeases                  bool                      `json:"recycle_leases,omitempty"`
	RelayAgentFilterRules          []*Filterrule             `json:"relay_agent_filter_rules,omitempty"`
	ServerAssociationType          string                    `json:"server_association_type,omitempty"`
	UnknownClients                 string                    `json:"unknown_clients,omitempty"`
	UpdateDnsOnLeaseRenewal        bool                      `json:"update_dns_on_lease_renewal,omitempty"`
	UseBootfile                    bool                      `json:"use_bootfile,omitempty"`
	UseBootserver                  bool                      `json:"use_bootserver,omitempty"`
	UseDdnsDomainname              bool                      `json:"use_ddns_domainname,omitempty"`
	UseDdnsGenerateHostname        bool                      `json:"use_ddns_generate_hostname,omitempty"`
	UseDenyBootp                   bool                      `json:"use_deny_bootp,omitempty"`
	UseEmailList                   bool                      `json:"use_email_list,omitempty"`
	UseEnableDdns                  bool                      `json:"use_enable_ddns,omitempty"`
	UseEnableDhcpThresholds        bool                      `json:"use_enable_dhcp_thresholds,omitempty"`
	UseIgnoreDhcpOptionListRequest bool                      `json:"use_ignore_dhcp_option_list_request,omitempty"`
	UseKnownClients                bool                      `json:"use_known_clients,omitempty"`
	UseLeaseScavengeTime           bool                      `json:"use_lease_scavenge_time,omitempty"`
	UseLogicFilterRules            bool                      `json:"use_logic_filter_rules,omitempty"`
	UseMsOptions                   bool                      `json:"use_ms_options,omitempty"`
	UseNextserver                  bool                      `json:"use_nextserver,omitempty"`
	UseOptions                     bool                      `json:"use_options,omitempty"`
	UsePxeLeaseTime                bool                      `json:"use_pxe_lease_time,omitempty"`
	UseRecycleLeases               bool                      `json:"use_recycle_leases,omitempty"`
	UseUnknownClients              bool                      `json:"use_unknown_clients,omitempty"`
	UseUpdateDnsOnLeaseRenewal     bool                      `json:"use_update_dns_on_lease_renewal,omitempty"`
}

func (Rangetemplate) ObjectType() string {
	return "rangetemplate"
}

func (obj Rangetemplate) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name", "number_of_addresses", "offset"}
	}
	return obj.returnFields
}

// RecordCaa represents Infoblox object record:caa
type RecordCaa struct {
	IBBase            `json:"-"`
	Ref               string            `json:"_ref,omitempty"`
	CaFlag            uint32            `json:"ca_flag,omitempty"`
	CaTag             string            `json:"ca_tag,omitempty"`
	CaValue           string            `json:"ca_value,omitempty"`
	CloudInfo         *GridCloudapiInfo `json:"cloud_info,omitempty"`
	Comment           string            `json:"comment,omitempty"`
	CreationTime      *time.Time        `json:"creation_time,omitempty"`
	Creator           string            `json:"creator,omitempty"`
	DdnsPrincipal     string            `json:"ddns_principal,omitempty"`
	DdnsProtected     bool              `json:"ddns_protected,omitempty"`
	Disable           bool              `json:"disable,omitempty"`
	DnsName           string            `json:"dns_name,omitempty"`
	Ea                EA                `json:"extattrs,omitempty"`
	ForbidReclamation bool              `json:"forbid_reclamation,omitempty"`
	LastQueried       *time.Time        `json:"last_queried,omitempty"`
	Name              string            `json:"name,omitempty"`
	Reclaimable       bool              `json:"reclaimable,omitempty"`
	Ttl               uint32            `json:"ttl,omitempty"`
	UseTtl            bool              `json:"use_ttl,omitempty"`
	View              string            `json:"view,omitempty"`
	Zone              string            `json:"zone,omitempty"`
}

func (RecordCaa) ObjectType() string {
	return "record:caa"
}

func (obj RecordCaa) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name", "view"}
	}
	return obj.returnFields
}

// RecordAlias represents Infoblox object record:alias
type RecordAlias struct {
	IBBase             `json:"-"`
	Ref                string              `json:"_ref,omitempty"`
	AwsRte53RecordInfo *Awsrte53recordinfo `json:"aws_rte53_record_info,omitempty"`
	CloudInfo          *GridCloudapiInfo   `json:"cloud_info,omitempty"`
	Comment            string              `json:"comment,omitempty"`
	Creator            string              `json:"creator,omitempty"`
	Disable            bool                `json:"disable,omitempty"`
	DnsName            string              `json:"dns_name,omitempty"`
	DnsTargetName      string              `json:"dns_target_name,omitempty"`
	Ea                 EA                  `json:"extattrs,omitempty"`
	LastQueried        *time.Time          `json:"last_queried,omitempty"`
	Name               string              `json:"name,omitempty"`
	TargetName         string              `json:"target_name,omitempty"`
	TargetType         string              `json:"target_type,omitempty"`
	Ttl                uint32              `json:"ttl,omitempty"`
	UseTtl             bool                `json:"use_ttl,omitempty"`
	View               string              `json:"view,omitempty"`
	Zone               string              `json:"zone,omitempty"`
}

func (RecordAlias) ObjectType() string {
	return "record:alias"
}

func (obj RecordAlias) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name", "target_name", "target_type", "view"}
	}
	return obj.returnFields
}

// RecordDhcid represents Infoblox object record:dhcid
type RecordDhcid struct {
	IBBase       `json:"-"`
	Ref          string     `json:"_ref,omitempty"`
	CreationTime *time.Time `json:"creation_time,omitempty"`
	Creator      string     `json:"creator,omitempty"`
	Dhcid        string     `json:"dhcid,omitempty"`
	DnsName      string     `json:"dns_name,omitempty"`
	Name         string     `json:"name,omitempty"`
	Ttl          uint32     `json:"ttl,omitempty"`
	UseTtl       bool       `json:"use_ttl,omitempty"`
	View         string     `json:"view,omitempty"`
	Zone         string     `json:"zone,omitempty"`
}

func (RecordDhcid) ObjectType() string {
	return "record:dhcid"
}

func (obj RecordDhcid) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name", "view"}
	}
	return obj.returnFields
}

// RecordDname represents Infoblox object record:dname
type RecordDname struct {
	IBBase            `json:"-"`
	Ref               string            `json:"_ref,omitempty"`
	CloudInfo         *GridCloudapiInfo `json:"cloud_info,omitempty"`
	Comment           string            `json:"comment,omitempty"`
	CreationTime      *time.Time        `json:"creation_time,omitempty"`
	Creator           string            `json:"creator,omitempty"`
	DdnsPrincipal     string            `json:"ddns_principal,omitempty"`
	DdnsProtected     bool              `json:"ddns_protected,omitempty"`
	Disable           bool              `json:"disable,omitempty"`
	DnsName           string            `json:"dns_name,omitempty"`
	DnsTarget         string            `json:"dns_target,omitempty"`
	Ea                EA                `json:"extattrs,omitempty"`
	ForbidReclamation bool              `json:"forbid_reclamation,omitempty"`
	LastQueried       *time.Time        `json:"last_queried,omitempty"`
	Name              string            `json:"name,omitempty"`
	Reclaimable       bool              `json:"reclaimable,omitempty"`
	SharedRecordGroup string            `json:"shared_record_group,omitempty"`
	Target            string            `json:"target,omitempty"`
	Ttl               uint32            `json:"ttl,omitempty"`
	UseTtl            bool              `json:"use_ttl,omitempty"`
	View              string            `json:"view,omitempty"`
	Zone              string            `json:"zone,omitempty"`
}

func (RecordDname) ObjectType() string {
	return "record:dname"
}

func (obj RecordDname) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name", "target", "view"}
	}
	return obj.returnFields
}

// RecordDnskey represents Infoblox object record:dnskey
type RecordDnskey struct {
	IBBase       `json:"-"`
	Ref          string     `json:"_ref,omitempty"`
	Algorithm    string     `json:"algorithm,omitempty"`
	Comment      string     `json:"comment,omitempty"`
	CreationTime *time.Time `json:"creation_time,omitempty"`
	Creator      string     `json:"creator,omitempty"`
	DnsName      string     `json:"dns_name,omitempty"`
	Flags        int        `json:"flags,omitempty"`
	KeyTag       uint32     `json:"key_tag,omitempty"`
	LastQueried  *time.Time `json:"last_queried,omitempty"`
	Name         string     `json:"name,omitempty"`
	PublicKey    string     `json:"public_key,omitempty"`
	Ttl          uint32     `json:"ttl,omitempty"`
	UseTtl       bool       `json:"use_ttl,omitempty"`
	View         string     `json:"view,omitempty"`
	Zone         string     `json:"zone,omitempty"`
}

func (RecordDnskey) ObjectType() string {
	return "record:dnskey"
}

func (obj RecordDnskey) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name", "view"}
	}
	return obj.returnFields
}

// RecordCNAME represents Infoblox object record:cname
type RecordCNAME struct {
	IBBase             `json:"-"`
	Ref                string              `json:"_ref,omitempty"`
	AwsRte53RecordInfo *Awsrte53recordinfo `json:"aws_rte53_record_info,omitempty"`
	Canonical          string              `json:"canonical,omitempty"`
	CloudInfo          *GridCloudapiInfo   `json:"cloud_info,omitempty"`
	Comment            string              `json:"comment,omitempty"`
	CreationTime       *time.Time          `json:"creation_time,omitempty"`
	Creator            string              `json:"creator,omitempty"`
	DdnsPrincipal      string              `json:"ddns_principal,omitempty"`
	DdnsProtected      bool                `json:"ddns_protected,omitempty"`
	Disable            bool                `json:"disable,omitempty"`
	DnsCanonical       string              `json:"dns_canonical,omitempty"`
	DnsName            string              `json:"dns_name,omitempty"`
	Ea                 EA                  `json:"extattrs,omitempty"`
	ForbidReclamation  bool                `json:"forbid_reclamation,omitempty"`
	LastQueried        *time.Time          `json:"last_queried,omitempty"`
	Name               string              `json:"name,omitempty"`
	Reclaimable        bool                `json:"reclaimable,omitempty"`
	SharedRecordGroup  string              `json:"shared_record_group,omitempty"`
	Ttl                uint32              `json:"ttl,omitempty"`
	UseTtl             bool                `json:"use_ttl,omitempty"`
	View               string              `json:"view,omitempty"`
	Zone               string              `json:"zone,omitempty"`
}

func (RecordCNAME) ObjectType() string {
	return "record:cname"
}

func (obj RecordCNAME) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"canonical", "name", "view"}
	}
	return obj.returnFields
}

func NewEmptyRecordCNAME() *RecordCNAME {
	res := &RecordCNAME{}
	res.returnFields = []string{"extattrs", "canonical", "name", "view", "zone", "comment", "ttl", "use_ttl"}

	return res
}

func NewRecordCNAME(dnsView string,
	canonical string,
	recordName string,
	useTtl bool,
	ttl uint32,
	comment string,
	ea EA,
	ref string) *RecordCNAME {

	res := NewEmptyRecordCNAME()
	res.View = dnsView
	res.Canonical = canonical
	res.Name = recordName
	res.UseTtl = useTtl
	res.Ttl = ttl
	res.Comment = comment
	res.Ea = ea
	res.Ref = ref

	return res
}

// Range represents Infoblox object range
type Range struct {
	IBBase                           `json:"-"`
	Ref                              string                      `json:"_ref,omitempty"`
	AlwaysUpdateDns                  bool                        `json:"always_update_dns,omitempty"`
	Bootfile                         string                      `json:"bootfile,omitempty"`
	Bootserver                       string                      `json:"bootserver,omitempty"`
	CloudInfo                        *GridCloudapiInfo           `json:"cloud_info,omitempty"`
	Comment                          string                      `json:"comment,omitempty"`
	DdnsDomainname                   string                      `json:"ddns_domainname,omitempty"`
	DdnsGenerateHostname             bool                        `json:"ddns_generate_hostname,omitempty"`
	DenyAllClients                   bool                        `json:"deny_all_clients,omitempty"`
	DenyBootp                        bool                        `json:"deny_bootp,omitempty"`
	DhcpUtilization                  uint32                      `json:"dhcp_utilization,omitempty"`
	DhcpUtilizationStatus            string                      `json:"dhcp_utilization_status,omitempty"`
	Disable                          bool                        `json:"disable,omitempty"`
	DiscoverNowStatus                string                      `json:"discover_now_status,omitempty"`
	DiscoveryBasicPollSettings       *DiscoveryBasicpollsettings `json:"discovery_basic_poll_settings,omitempty"`
	DiscoveryBlackoutSetting         *PropertiesBlackoutsetting  `json:"discovery_blackout_setting,omitempty"`
	DiscoveryMember                  string                      `json:"discovery_member,omitempty"`
	DynamicHosts                     uint32                      `json:"dynamic_hosts,omitempty"`
	EmailList                        []string                    `json:"email_list,omitempty"`
	EnableDdns                       bool                        `json:"enable_ddns,omitempty"`
	EnableDhcpThresholds             bool                        `json:"enable_dhcp_thresholds,omitempty"`
	EnableDiscovery                  bool                        `json:"enable_discovery,omitempty"`
	EnableEmailWarnings              bool                        `json:"enable_email_warnings,omitempty"`
	EnableIfmapPublishing            bool                        `json:"enable_ifmap_publishing,omitempty"`
	EnableImmediateDiscovery         bool                        `json:"enable_immediate_discovery,omitempty"`
	EnablePxeLeaseTime               bool                        `json:"enable_pxe_lease_time,omitempty"`
	EnableSnmpWarnings               bool                        `json:"enable_snmp_warnings,omitempty"`
	EndAddr                          string                      `json:"end_addr,omitempty"`
	EndpointSources                  []*CiscoiseEndpoint         `json:"endpoint_sources,omitempty"`
	Exclude                          []*Exclusionrange           `json:"exclude,omitempty"`
	Ea                               EA                          `json:"extattrs,omitempty"`
	FailoverAssociation              string                      `json:"failover_association,omitempty"`
	FingerprintFilterRules           []*Filterrule               `json:"fingerprint_filter_rules,omitempty"`
	HighWaterMark                    uint32                      `json:"high_water_mark,omitempty"`
	HighWaterMarkReset               uint32                      `json:"high_water_mark_reset,omitempty"`
	IgnoreDhcpOptionListRequest      bool                        `json:"ignore_dhcp_option_list_request,omitempty"`
	IgnoreId                         string                      `json:"ignore_id,omitempty"`
	IgnoreMacAddresses               []string                    `json:"ignore_mac_addresses,omitempty"`
	IsSplitScope                     bool                        `json:"is_split_scope,omitempty"`
	KnownClients                     string                      `json:"known_clients,omitempty"`
	LeaseScavengeTime                int                         `json:"lease_scavenge_time,omitempty"`
	LogicFilterRules                 []*Logicfilterrule          `json:"logic_filter_rules,omitempty"`
	LowWaterMark                     uint32                      `json:"low_water_mark,omitempty"`
	LowWaterMarkReset                uint32                      `json:"low_water_mark_reset,omitempty"`
	MacFilterRules                   []*Filterrule               `json:"mac_filter_rules,omitempty"`
	Member                           *Dhcpmember                 `json:"member,omitempty"`
	MsAdUserData                     *MsserverAduserData         `json:"ms_ad_user_data,omitempty"`
	MsOptions                        []*Msdhcpoption             `json:"ms_options,omitempty"`
	MsServer                         *Msdhcpserver               `json:"ms_server,omitempty"`
	NacFilterRules                   []*Filterrule               `json:"nac_filter_rules,omitempty"`
	Name                             string                      `json:"name,omitempty"`
	Network                          string                      `json:"network,omitempty"`
	NetworkView                      string                      `json:"network_view,omitempty"`
	NextAvailableIp                  *Nextavailableip            `json:"next_available_ip,omitempty"`
	Nextserver                       string                      `json:"nextserver,omitempty"`
	OptionFilterRules                []*Filterrule               `json:"option_filter_rules,omitempty"`
	Options                          []*Dhcpoption               `json:"options,omitempty"`
	PortControlBlackoutSetting       *PropertiesBlackoutsetting  `json:"port_control_blackout_setting,omitempty"`
	PxeLeaseTime                     uint32                      `json:"pxe_lease_time,omitempty"`
	RecycleLeases                    bool                        `json:"recycle_leases,omitempty"`
	RelayAgentFilterRules            []*Filterrule               `json:"relay_agent_filter_rules,omitempty"`
	RestartIfNeeded                  bool                        `json:"restart_if_needed,omitempty"`
	SamePortControlDiscoveryBlackout bool                        `json:"same_port_control_discovery_blackout,omitempty"`
	ServerAssociationType            string                      `json:"server_association_type,omitempty"`
	SplitMember                      *Msdhcpserver               `json:"split_member,omitempty"`
	SplitScopeExclusionPercent       uint32                      `json:"split_scope_exclusion_percent,omitempty"`
	StartAddr                        string                      `json:"start_addr,omitempty"`
	StaticHosts                      uint32                      `json:"static_hosts,omitempty"`
	SubscribeSettings                *CiscoiseSubscribesetting   `json:"subscribe_settings,omitempty"`
	Template                         string                      `json:"template,omitempty"`
	TotalHosts                       uint32                      `json:"total_hosts,omitempty"`
	UnknownClients                   string                      `json:"unknown_clients,omitempty"`
	UpdateDnsOnLeaseRenewal          bool                        `json:"update_dns_on_lease_renewal,omitempty"`
	UseBlackoutSetting               bool                        `json:"use_blackout_setting,omitempty"`
	UseBootfile                      bool                        `json:"use_bootfile,omitempty"`
	UseBootserver                    bool                        `json:"use_bootserver,omitempty"`
	UseDdnsDomainname                bool                        `json:"use_ddns_domainname,omitempty"`
	UseDdnsGenerateHostname          bool                        `json:"use_ddns_generate_hostname,omitempty"`
	UseDenyBootp                     bool                        `json:"use_deny_bootp,omitempty"`
	UseDiscoveryBasicPollingSettings bool                        `json:"use_discovery_basic_polling_settings,omitempty"`
	UseEmailList                     bool                        `json:"use_email_list,omitempty"`
	UseEnableDdns                    bool                        `json:"use_enable_ddns,omitempty"`
	UseEnableDhcpThresholds          bool                        `json:"use_enable_dhcp_thresholds,omitempty"`
	UseEnableDiscovery               bool                        `json:"use_enable_discovery,omitempty"`
	UseEnableIfmapPublishing         bool                        `json:"use_enable_ifmap_publishing,omitempty"`
	UseIgnoreDhcpOptionListRequest   bool                        `json:"use_ignore_dhcp_option_list_request,omitempty"`
	UseIgnoreId                      bool                        `json:"use_ignore_id,omitempty"`
	UseKnownClients                  bool                        `json:"use_known_clients,omitempty"`
	UseLeaseScavengeTime             bool                        `json:"use_lease_scavenge_time,omitempty"`
	UseLogicFilterRules              bool                        `json:"use_logic_filter_rules,omitempty"`
	UseMsOptions                     bool                        `json:"use_ms_options,omitempty"`
	UseNextserver                    bool                        `json:"use_nextserver,omitempty"`
	UseOptions                       bool                        `json:"use_options,omitempty"`
	UsePxeLeaseTime                  bool                        `json:"use_pxe_lease_time,omitempty"`
	UseRecycleLeases                 bool                        `json:"use_recycle_leases,omitempty"`
	UseSubscribeSettings             bool                        `json:"use_subscribe_settings,omitempty"`
	UseUnknownClients                bool                        `json:"use_unknown_clients,omitempty"`
	UseUpdateDnsOnLeaseRenewal       bool                        `json:"use_update_dns_on_lease_renewal,omitempty"`
}

func (Range) ObjectType() string {
	return "range"
}

func (obj Range) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "end_addr", "network", "network_view", "start_addr"}
	}
	return obj.returnFields
}

// RecordDs represents Infoblox object record:ds
type RecordDs struct {
	IBBase       `json:"-"`
	Ref          string            `json:"_ref,omitempty"`
	Algorithm    string            `json:"algorithm,omitempty"`
	CloudInfo    *GridCloudapiInfo `json:"cloud_info,omitempty"`
	Comment      string            `json:"comment,omitempty"`
	CreationTime *time.Time        `json:"creation_time,omitempty"`
	Creator      string            `json:"creator,omitempty"`
	Digest       string            `json:"digest,omitempty"`
	DigestType   string            `json:"digest_type,omitempty"`
	DnsName      string            `json:"dns_name,omitempty"`
	KeyTag       uint32            `json:"key_tag,omitempty"`
	LastQueried  *time.Time        `json:"last_queried,omitempty"`
	Name         string            `json:"name,omitempty"`
	Ttl          uint32            `json:"ttl,omitempty"`
	UseTtl       bool              `json:"use_ttl,omitempty"`
	View         string            `json:"view,omitempty"`
	Zone         string            `json:"zone,omitempty"`
}

func (RecordDs) ObjectType() string {
	return "record:ds"
}

func (obj RecordDs) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name", "view"}
	}
	return obj.returnFields
}

// RecordDtclbdn represents Infoblox object record:dtclbdn
type RecordDtclbdn struct {
	IBBase      `json:"-"`
	Ref         string     `json:"_ref,omitempty"`
	Comment     string     `json:"comment,omitempty"`
	Disable     bool       `json:"disable,omitempty"`
	Ea          EA         `json:"extattrs,omitempty"`
	LastQueried *time.Time `json:"last_queried,omitempty"`
	Lbdn        string     `json:"lbdn,omitempty"`
	Name        string     `json:"name,omitempty"`
	Pattern     string     `json:"pattern,omitempty"`
	View        string     `json:"view,omitempty"`
	Zone        string     `json:"zone,omitempty"`
}

func (RecordDtclbdn) ObjectType() string {
	return "record:dtclbdn"
}

func (obj RecordDtclbdn) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name", "view", "zone"}
	}
	return obj.returnFields
}

// RecordNaptr represents Infoblox object record:naptr
type RecordNaptr struct {
	IBBase            `json:"-"`
	Ref               string            `json:"_ref,omitempty"`
	CloudInfo         *GridCloudapiInfo `json:"cloud_info,omitempty"`
	Comment           string            `json:"comment,omitempty"`
	CreationTime      *time.Time        `json:"creation_time,omitempty"`
	Creator           string            `json:"creator,omitempty"`
	DdnsPrincipal     string            `json:"ddns_principal,omitempty"`
	DdnsProtected     bool              `json:"ddns_protected,omitempty"`
	Disable           bool              `json:"disable,omitempty"`
	DnsName           string            `json:"dns_name,omitempty"`
	DnsReplacement    string            `json:"dns_replacement,omitempty"`
	Ea                EA                `json:"extattrs,omitempty"`
	Flags             string            `json:"flags,omitempty"`
	ForbidReclamation bool              `json:"forbid_reclamation,omitempty"`
	LastQueried       *time.Time        `json:"last_queried,omitempty"`
	Name              string            `json:"name,omitempty"`
	Order             uint32            `json:"order,omitempty"`
	Preference        uint32            `json:"preference,omitempty"`
	Reclaimable       bool              `json:"reclaimable,omitempty"`
	Regexp            string            `json:"regexp,omitempty"`
	Replacement       string            `json:"replacement,omitempty"`
	Services          string            `json:"services,omitempty"`
	Ttl               uint32            `json:"ttl,omitempty"`
	UseTtl            bool              `json:"use_ttl,omitempty"`
	View              string            `json:"view,omitempty"`
	Zone              string            `json:"zone,omitempty"`
}

func (RecordNaptr) ObjectType() string {
	return "record:naptr"
}

func (obj RecordNaptr) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name", "order", "preference", "regexp", "replacement", "services", "view"}
	}
	return obj.returnFields
}

// HostRecord represents Infoblox object record:host
type HostRecord struct {
	IBBase                   `json:"-"`
	Ref                      string                    `json:"_ref,omitempty"`
	Aliases                  []string                  `json:"aliases,omitempty"`
	AllowTelnet              bool                      `json:"allow_telnet,omitempty"`
	CliCredentials           []*DiscoveryClicredential `json:"cli_credentials,omitempty"`
	CloudInfo                *GridCloudapiInfo         `json:"cloud_info,omitempty"`
	Comment                  string                    `json:"comment,omitempty"`
	EnableDns                bool                      `json:"configure_for_dns,omitempty"`
	DdnsProtected            bool                      `json:"ddns_protected,omitempty"`
	DeviceDescription        string                    `json:"device_description,omitempty"`
	DeviceLocation           string                    `json:"device_location,omitempty"`
	DeviceType               string                    `json:"device_type,omitempty"`
	DeviceVendor             string                    `json:"device_vendor,omitempty"`
	Disable                  bool                      `json:"disable,omitempty"`
	DisableDiscovery         bool                      `json:"disable_discovery,omitempty"`
	DnsAliases               []string                  `json:"dns_aliases,omitempty"`
	DnsName                  string                    `json:"dns_name,omitempty"`
	EnableImmediateDiscovery bool                      `json:"enable_immediate_discovery,omitempty"`
	Ea                       EA                        `json:"extattrs,omitempty"`
	Ipv4Addrs                []HostRecordIpv4Addr      `json:"ipv4addrs,omitempty"`
	Ipv6Addrs                []HostRecordIpv6Addr      `json:"ipv6addrs,omitempty"`
	LastQueried              *time.Time                `json:"last_queried,omitempty"`
	MsAdUserData             *MsserverAduserData       `json:"ms_ad_user_data,omitempty"`
	Name                     string                    `json:"name,omitempty"`
	NetworkView              string                    `json:"network_view,omitempty"`
	RestartIfNeeded          bool                      `json:"restart_if_needed,omitempty"`
	RrsetOrder               string                    `json:"rrset_order,omitempty"`
	Snmp3Credential          *DiscoverySnmp3credential `json:"snmp3_credential,omitempty"`
	SnmpCredential           *DiscoverySnmpcredential  `json:"snmp_credential,omitempty"`
	Ttl                      uint32                    `json:"ttl,omitempty"`
	UseCliCredentials        bool                      `json:"use_cli_credentials,omitempty"`
	UseSnmp3Credential       bool                      `json:"use_snmp3_credential,omitempty"`
	UseSnmpCredential        bool                      `json:"use_snmp_credential,omitempty"`
	UseTtl                   bool                      `json:"use_ttl,omitempty"`
	View                     string                    `json:"view,omitempty"`
	Zone                     string                    `json:"zone,omitempty"`
}

func (HostRecord) ObjectType() string {
	return "record:host"
}

func (obj HostRecord) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"ipv4addrs", "ipv6addrs", "name", "view"}
	}
	return obj.returnFields
}

func NewEmptyHostRecord() *HostRecord {
	res := &HostRecord{}
	res.returnFields = []string{"extattrs", "ipv4addrs", "ipv6addrs", "name", "view", "zone", "comment", "network_view", "aliases", "use_ttl", "ttl", "configure_for_dns"}
	return res
}

func NewHostRecord(
	netView string,
	name string,
	ipv4Addr string,
	ipv6Addr string,
	ipv4AddrList []HostRecordIpv4Addr,
	ipv6AddrList []HostRecordIpv6Addr,
	eas EA,
	enableDNS bool,
	dnsView string,
	zone string,
	ref string,
	useTtl bool,
	ttl uint32,
	comment string,
	aliases []string) *HostRecord {

	res := NewEmptyHostRecord()
	res.NetworkView = netView
	res.Name = name
	res.Ea = eas
	res.View = dnsView
	res.Zone = zone
	res.Ref = ref
	res.Comment = comment
	//res.Ipv4Addr = ipv4Addr
	//res.Ipv6Addr = ipv6Addr
	res.Ipv4Addrs = ipv4AddrList
	res.Ipv6Addrs = ipv6AddrList
	res.UseTtl = useTtl
	res.Ttl = ttl
	res.Aliases = aliases
	res.EnableDns = enableDNS

	return res
}

// RecordNS represents Infoblox object record:ns
type RecordNS struct {
	IBBase           `json:"-"`
	Ref              string            `json:"_ref,omitempty"`
	Addresses        []*ZoneNameServer `json:"addresses,omitempty"`
	CloudInfo        *GridCloudapiInfo `json:"cloud_info,omitempty"`
	Creator          string            `json:"creator,omitempty"`
	DnsName          string            `json:"dns_name,omitempty"`
	LastQueried      *time.Time        `json:"last_queried,omitempty"`
	MsDelegationName string            `json:"ms_delegation_name,omitempty"`
	Name             string            `json:"name,omitempty"`
	Nameserver       string            `json:"nameserver,omitempty"`
	Policy           string            `json:"policy,omitempty"`
	View             string            `json:"view,omitempty"`
	Zone             string            `json:"zone,omitempty"`
}

func (RecordNS) ObjectType() string {
	return "record:ns"
}

func (obj RecordNS) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name", "nameserver", "view"}
	}
	return obj.returnFields
}

// RecordAAAA represents Infoblox object record:aaaa
type RecordAAAA struct {
	IBBase              `json:"-"`
	Ref                 string              `json:"_ref,omitempty"`
	AwsRte53RecordInfo  *Awsrte53recordinfo `json:"aws_rte53_record_info,omitempty"`
	CloudInfo           *GridCloudapiInfo   `json:"cloud_info,omitempty"`
	Comment             string              `json:"comment,omitempty"`
	CreationTime        *time.Time          `json:"creation_time,omitempty"`
	Creator             string              `json:"creator,omitempty"`
	DdnsPrincipal       string              `json:"ddns_principal,omitempty"`
	DdnsProtected       bool                `json:"ddns_protected,omitempty"`
	Disable             bool                `json:"disable,omitempty"`
	DiscoveredData      *Discoverydata      `json:"discovered_data,omitempty"`
	DnsName             string              `json:"dns_name,omitempty"`
	Ea                  EA                  `json:"extattrs,omitempty"`
	ForbidReclamation   bool                `json:"forbid_reclamation,omitempty"`
	Ipv6Addr            string              `json:"ipv6addr,omitempty"`
	LastQueried         *time.Time          `json:"last_queried,omitempty"`
	MsAdUserData        *MsserverAduserData `json:"ms_ad_user_data,omitempty"`
	Name                string              `json:"name,omitempty"`
	Reclaimable         bool                `json:"reclaimable,omitempty"`
	RemoveAssociatedPtr bool                `json:"remove_associated_ptr,omitempty"`
	SharedRecordGroup   string              `json:"shared_record_group,omitempty"`
	Ttl                 uint32              `json:"ttl,omitempty"`
	UseTtl              bool                `json:"use_ttl,omitempty"`
	View                string              `json:"view,omitempty"`
	Zone                string              `json:"zone,omitempty"`
}

func (RecordAAAA) ObjectType() string {
	return "record:aaaa"
}

func (obj RecordAAAA) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"ipv6addr", "name", "view"}
	}
	return obj.returnFields
}

func NewEmptyRecordAAAA() *RecordAAAA {
	res := &RecordAAAA{}
	res.returnFields = []string{"extattrs", "ipv6addr", "name", "view", "zone", "use_ttl", "ttl", "comment"}

	return res
}

func NewRecordAAAA(
	view string,
	name string,
	ipAddr string,
	useTtl bool,
	ttl uint32,
	comment string,
	eas EA,
	ref string) *RecordAAAA {

	res := NewEmptyRecordAAAA()
	res.View = view
	res.Name = name
	res.Ipv6Addr = ipAddr
	res.UseTtl = useTtl
	res.Ttl = ttl
	res.Comment = comment
	res.Ea = eas
	res.Ref = ref

	return res
}

// RecordMx represents Infoblox object record:mx
type RecordMx struct {
	IBBase             `json:"-"`
	Ref                string              `json:"_ref,omitempty"`
	AwsRte53RecordInfo *Awsrte53recordinfo `json:"aws_rte53_record_info,omitempty"`
	CloudInfo          *GridCloudapiInfo   `json:"cloud_info,omitempty"`
	Comment            string              `json:"comment,omitempty"`
	CreationTime       *time.Time          `json:"creation_time,omitempty"`
	Creator            string              `json:"creator,omitempty"`
	DdnsPrincipal      string              `json:"ddns_principal,omitempty"`
	DdnsProtected      bool                `json:"ddns_protected,omitempty"`
	Disable            bool                `json:"disable,omitempty"`
	DnsMailExchanger   string              `json:"dns_mail_exchanger,omitempty"`
	DnsName            string              `json:"dns_name,omitempty"`
	Ea                 EA                  `json:"extattrs,omitempty"`
	ForbidReclamation  bool                `json:"forbid_reclamation,omitempty"`
	LastQueried        *time.Time          `json:"last_queried,omitempty"`
	MailExchanger      string              `json:"mail_exchanger,omitempty"`
	Name               string              `json:"name,omitempty"`
	Preference         uint32              `json:"preference,omitempty"`
	Reclaimable        bool                `json:"reclaimable,omitempty"`
	SharedRecordGroup  string              `json:"shared_record_group,omitempty"`
	Ttl                uint32              `json:"ttl,omitempty"`
	UseTtl             bool                `json:"use_ttl,omitempty"`
	View               string              `json:"view,omitempty"`
	Zone               string              `json:"zone,omitempty"`
}

func (RecordMx) ObjectType() string {
	return "record:mx"
}

func (obj RecordMx) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"mail_exchanger", "name", "preference", "view"}
	}
	return obj.returnFields
}

// HostRecordIpv6Addr represents Infoblox object record:host_ipv6addr
type HostRecordIpv6Addr struct {
	IBBase               `json:"-"`
	Ref                  string              `json:"_ref,omitempty"`
	AddressType          string              `json:"address_type,omitempty"`
	EnableDhcp           bool                `json:"configure_for_dhcp,omitempty"`
	DiscoverNowStatus    string              `json:"discover_now_status,omitempty"`
	DiscoveredData       *Discoverydata      `json:"discovered_data,omitempty"`
	DomainName           string              `json:"domain_name,omitempty"`
	DomainNameServers    []string            `json:"domain_name_servers,omitempty"`
	Duid                 string              `json:"duid,omitempty"`
	Host                 string              `json:"host,omitempty"`
	Ipv6Addr             string              `json:"ipv6addr,omitempty"`
	Ipv6prefix           string              `json:"ipv6prefix,omitempty"`
	Ipv6prefixBits       uint32              `json:"ipv6prefix_bits,omitempty"`
	LastQueried          *time.Time          `json:"last_queried,omitempty"`
	MatchClient          string              `json:"match_client,omitempty"`
	MsAdUserData         *MsserverAduserData `json:"ms_ad_user_data,omitempty"`
	Network              string              `json:"network,omitempty"`
	NetworkView          string              `json:"network_view,omitempty"`
	Options              []*Dhcpoption       `json:"options,omitempty"`
	PreferredLifetime    uint32              `json:"preferred_lifetime,omitempty"`
	ReservedInterface    string              `json:"reserved_interface,omitempty"`
	UseDomainName        bool                `json:"use_domain_name,omitempty"`
	UseDomainNameServers bool                `json:"use_domain_name_servers,omitempty"`
	UseForEaInheritance  bool                `json:"use_for_ea_inheritance,omitempty"`
	UseOptions           bool                `json:"use_options,omitempty"`
	UsePreferredLifetime bool                `json:"use_preferred_lifetime,omitempty"`
	UseValidLifetime     bool                `json:"use_valid_lifetime,omitempty"`
	ValidLifetime        uint32              `json:"valid_lifetime,omitempty"`
}

func (HostRecordIpv6Addr) ObjectType() string {
	return "record:host_ipv6addr"
}

func (obj HostRecordIpv6Addr) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"configure_for_dhcp", "duid", "host", "ipv6addr"}
	}
	return obj.returnFields
}

func NewEmptyHostRecordIpv6Addr() *HostRecordIpv6Addr {
	return &HostRecordIpv6Addr{}
}

func NewHostRecordIpv6Addr(
	ipAddr string,
	duid string,
	enableDhcp bool,
	ref string) *HostRecordIpv6Addr {

	res := NewEmptyHostRecordIpv6Addr()
	res.Ipv6Addr = ipAddr
	res.Duid = duid
	res.Ref = ref
	res.EnableDhcp = enableDhcp

	return res
}

// HostRecordIpv4Addr represents Infoblox object record:host_ipv4addr
type HostRecordIpv4Addr struct {
	IBBase                          `json:"-"`
	Ref                             string              `json:"_ref,omitempty"`
	Bootfile                        string              `json:"bootfile,omitempty"`
	Bootserver                      string              `json:"bootserver,omitempty"`
	EnableDhcp                      bool                `json:"configure_for_dhcp,omitempty"`
	DenyBootp                       bool                `json:"deny_bootp,omitempty"`
	DiscoverNowStatus               string              `json:"discover_now_status,omitempty"`
	DiscoveredData                  *Discoverydata      `json:"discovered_data,omitempty"`
	EnablePxeLeaseTime              bool                `json:"enable_pxe_lease_time,omitempty"`
	Host                            string              `json:"host,omitempty"`
	IgnoreClientRequestedOptions    bool                `json:"ignore_client_requested_options,omitempty"`
	Ipv4Addr                        string              `json:"ipv4addr,omitempty"`
	IsInvalidMac                    bool                `json:"is_invalid_mac,omitempty"`
	LastQueried                     *time.Time          `json:"last_queried,omitempty"`
	LogicFilterRules                []*Logicfilterrule  `json:"logic_filter_rules,omitempty"`
	Mac                             string              `json:"mac,omitempty"`
	MatchClient                     string              `json:"match_client,omitempty"`
	MsAdUserData                    *MsserverAduserData `json:"ms_ad_user_data,omitempty"`
	Network                         string              `json:"network,omitempty"`
	NetworkView                     string              `json:"network_view,omitempty"`
	Nextserver                      string              `json:"nextserver,omitempty"`
	Options                         []*Dhcpoption       `json:"options,omitempty"`
	PxeLeaseTime                    uint32              `json:"pxe_lease_time,omitempty"`
	ReservedInterface               string              `json:"reserved_interface,omitempty"`
	UseBootfile                     bool                `json:"use_bootfile,omitempty"`
	UseBootserver                   bool                `json:"use_bootserver,omitempty"`
	UseDenyBootp                    bool                `json:"use_deny_bootp,omitempty"`
	UseForEaInheritance             bool                `json:"use_for_ea_inheritance,omitempty"`
	UseIgnoreClientRequestedOptions bool                `json:"use_ignore_client_requested_options,omitempty"`
	UseLogicFilterRules             bool                `json:"use_logic_filter_rules,omitempty"`
	UseNextserver                   bool                `json:"use_nextserver,omitempty"`
	UseOptions                      bool                `json:"use_options,omitempty"`
	UsePxeLeaseTime                 bool                `json:"use_pxe_lease_time,omitempty"`
}

func (HostRecordIpv4Addr) ObjectType() string {
	return "record:host_ipv4addr"
}

func (obj HostRecordIpv4Addr) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"configure_for_dhcp", "host", "ipv4addr", "mac"}
	}
	return obj.returnFields
}

func NewEmptyHostRecordIpv4Addr() *HostRecordIpv4Addr {
	return &HostRecordIpv4Addr{}
}

func NewHostRecordIpv4Addr(
	ipAddr string,
	macAddr string,
	enableDhcp bool,
	ref string) *HostRecordIpv4Addr {

	res := NewEmptyHostRecordIpv4Addr()
	res.Ipv4Addr = ipAddr
	res.Mac = macAddr
	res.Ref = ref
	res.EnableDhcp = enableDhcp

	return res
}

// RecordNsec represents Infoblox object record:nsec
type RecordNsec struct {
	IBBase           `json:"-"`
	Ref              string            `json:"_ref,omitempty"`
	CloudInfo        *GridCloudapiInfo `json:"cloud_info,omitempty"`
	CreationTime     *time.Time        `json:"creation_time,omitempty"`
	Creator          string            `json:"creator,omitempty"`
	DnsName          string            `json:"dns_name,omitempty"`
	DnsNextOwnerName string            `json:"dns_next_owner_name,omitempty"`
	LastQueried      *time.Time        `json:"last_queried,omitempty"`
	Name             string            `json:"name,omitempty"`
	NextOwnerName    string            `json:"next_owner_name,omitempty"`
	RrsetTypes       []string          `json:"rrset_types,omitempty"`
	Ttl              uint32            `json:"ttl,omitempty"`
	UseTtl           bool              `json:"use_ttl,omitempty"`
	View             string            `json:"view,omitempty"`
	Zone             string            `json:"zone,omitempty"`
}

func (RecordNsec) ObjectType() string {
	return "record:nsec"
}

func (obj RecordNsec) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name", "view"}
	}
	return obj.returnFields
}

// RecordNsec3 represents Infoblox object record:nsec3
type RecordNsec3 struct {
	IBBase        `json:"-"`
	Ref           string            `json:"_ref,omitempty"`
	Algorithm     string            `json:"algorithm,omitempty"`
	CloudInfo     *GridCloudapiInfo `json:"cloud_info,omitempty"`
	CreationTime  *time.Time        `json:"creation_time,omitempty"`
	Creator       string            `json:"creator,omitempty"`
	DnsName       string            `json:"dns_name,omitempty"`
	Flags         uint32            `json:"flags,omitempty"`
	Iterations    uint32            `json:"iterations,omitempty"`
	LastQueried   *time.Time        `json:"last_queried,omitempty"`
	Name          string            `json:"name,omitempty"`
	NextOwnerName string            `json:"next_owner_name,omitempty"`
	RrsetTypes    []string          `json:"rrset_types,omitempty"`
	Salt          string            `json:"salt,omitempty"`
	Ttl           uint32            `json:"ttl,omitempty"`
	UseTtl        bool              `json:"use_ttl,omitempty"`
	View          string            `json:"view,omitempty"`
	Zone          string            `json:"zone,omitempty"`
}

func (RecordNsec3) ObjectType() string {
	return "record:nsec3"
}

func (obj RecordNsec3) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name", "view"}
	}
	return obj.returnFields
}

// RecordA represents Infoblox object record:a
type RecordA struct {
	IBBase              `json:"-"`
	Ref                 string              `json:"_ref,omitempty"`
	AwsRte53RecordInfo  *Awsrte53recordinfo `json:"aws_rte53_record_info,omitempty"`
	CloudInfo           *GridCloudapiInfo   `json:"cloud_info,omitempty"`
	Comment             string              `json:"comment,omitempty"`
	CreationTime        *time.Time          `json:"creation_time,omitempty"`
	Creator             string              `json:"creator,omitempty"`
	DdnsPrincipal       string              `json:"ddns_principal,omitempty"`
	DdnsProtected       bool                `json:"ddns_protected,omitempty"`
	Disable             bool                `json:"disable,omitempty"`
	DiscoveredData      *Discoverydata      `json:"discovered_data,omitempty"`
	DnsName             string              `json:"dns_name,omitempty"`
	Ea                  EA                  `json:"extattrs,omitempty"`
	ForbidReclamation   bool                `json:"forbid_reclamation,omitempty"`
	Ipv4Addr            string              `json:"ipv4addr,omitempty"`
	LastQueried         *time.Time          `json:"last_queried,omitempty"`
	MsAdUserData        *MsserverAduserData `json:"ms_ad_user_data,omitempty"`
	Name                string              `json:"name,omitempty"`
	Reclaimable         bool                `json:"reclaimable,omitempty"`
	RemoveAssociatedPtr bool                `json:"remove_associated_ptr,omitempty"`
	SharedRecordGroup   string              `json:"shared_record_group,omitempty"`
	Ttl                 uint32              `json:"ttl,omitempty"`
	UseTtl              bool                `json:"use_ttl,omitempty"`
	View                string              `json:"view,omitempty"`
	Zone                string              `json:"zone,omitempty"`
}

func (RecordA) ObjectType() string {
	return "record:a"
}

func (obj RecordA) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"ipv4addr", "name", "view"}
	}
	return obj.returnFields
}

func NewEmptyRecordA() *RecordA {
	res := &RecordA{}
	res.returnFields = []string{
		"extattrs", "ipv4addr", "name", "view", "zone", "comment", "ttl", "use_ttl"}

	return res
}

func NewRecordA(
	view string,
	zone string,
	name string,
	ipAddr string,
	ttl uint32,
	useTTL bool,
	comment string,
	eas EA,
	ref string) *RecordA {

	res := NewEmptyRecordA()
	res.View = view
	res.Zone = zone
	res.Name = name
	res.Ipv4Addr = ipAddr
	res.Ttl = ttl
	res.UseTtl = useTTL
	res.Comment = comment
	res.Ea = eas
	res.Ref = ref

	return res
}

// RecordNsec3param represents Infoblox object record:nsec3param
type RecordNsec3param struct {
	IBBase       `json:"-"`
	Ref          string            `json:"_ref,omitempty"`
	Algorithm    string            `json:"algorithm,omitempty"`
	CloudInfo    *GridCloudapiInfo `json:"cloud_info,omitempty"`
	CreationTime *time.Time        `json:"creation_time,omitempty"`
	Creator      string            `json:"creator,omitempty"`
	DnsName      string            `json:"dns_name,omitempty"`
	Flags        uint32            `json:"flags,omitempty"`
	Iterations   uint32            `json:"iterations,omitempty"`
	LastQueried  *time.Time        `json:"last_queried,omitempty"`
	Name         string            `json:"name,omitempty"`
	Salt         string            `json:"salt,omitempty"`
	Ttl          uint32            `json:"ttl,omitempty"`
	UseTtl       bool              `json:"use_ttl,omitempty"`
	View         string            `json:"view,omitempty"`
	Zone         string            `json:"zone,omitempty"`
}

func (RecordNsec3param) ObjectType() string {
	return "record:nsec3param"
}

func (obj RecordNsec3param) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name", "view"}
	}
	return obj.returnFields
}

// RecordRpzAIpaddress represents Infoblox object record:rpz:a:ipaddress
type RecordRpzAIpaddress struct {
	IBBase   `json:"-"`
	Ref      string `json:"_ref,omitempty"`
	Comment  string `json:"comment,omitempty"`
	Disable  bool   `json:"disable,omitempty"`
	Ea       EA     `json:"extattrs,omitempty"`
	Ipv4Addr string `json:"ipv4addr,omitempty"`
	Name     string `json:"name,omitempty"`
	RpZone   string `json:"rp_zone,omitempty"`
	Ttl      uint32 `json:"ttl,omitempty"`
	UseTtl   bool   `json:"use_ttl,omitempty"`
	View     string `json:"view,omitempty"`
	Zone     string `json:"zone,omitempty"`
}

func (RecordRpzAIpaddress) ObjectType() string {
	return "record:rpz:a:ipaddress"
}

func (obj RecordRpzAIpaddress) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"ipv4addr", "name", "view"}
	}
	return obj.returnFields
}

// RecordRpzA represents Infoblox object record:rpz:a
type RecordRpzA struct {
	IBBase   `json:"-"`
	Ref      string `json:"_ref,omitempty"`
	Comment  string `json:"comment,omitempty"`
	Disable  bool   `json:"disable,omitempty"`
	Ea       EA     `json:"extattrs,omitempty"`
	Ipv4Addr string `json:"ipv4addr,omitempty"`
	Name     string `json:"name,omitempty"`
	RpZone   string `json:"rp_zone,omitempty"`
	Ttl      uint32 `json:"ttl,omitempty"`
	UseTtl   bool   `json:"use_ttl,omitempty"`
	View     string `json:"view,omitempty"`
	Zone     string `json:"zone,omitempty"`
}

func (RecordRpzA) ObjectType() string {
	return "record:rpz:a"
}

func (obj RecordRpzA) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"ipv4addr", "name", "view"}
	}
	return obj.returnFields
}

// RecordRpzAaaa represents Infoblox object record:rpz:aaaa
type RecordRpzAaaa struct {
	IBBase   `json:"-"`
	Ref      string `json:"_ref,omitempty"`
	Comment  string `json:"comment,omitempty"`
	Disable  bool   `json:"disable,omitempty"`
	Ea       EA     `json:"extattrs,omitempty"`
	Ipv6Addr string `json:"ipv6addr,omitempty"`
	Name     string `json:"name,omitempty"`
	RpZone   string `json:"rp_zone,omitempty"`
	Ttl      uint32 `json:"ttl,omitempty"`
	UseTtl   bool   `json:"use_ttl,omitempty"`
	View     string `json:"view,omitempty"`
	Zone     string `json:"zone,omitempty"`
}

func (RecordRpzAaaa) ObjectType() string {
	return "record:rpz:aaaa"
}

func (obj RecordRpzAaaa) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"ipv6addr", "name", "view"}
	}
	return obj.returnFields
}

// RecordRpzAaaaIpaddress represents Infoblox object record:rpz:aaaa:ipaddress
type RecordRpzAaaaIpaddress struct {
	IBBase   `json:"-"`
	Ref      string `json:"_ref,omitempty"`
	Comment  string `json:"comment,omitempty"`
	Disable  bool   `json:"disable,omitempty"`
	Ea       EA     `json:"extattrs,omitempty"`
	Ipv6Addr string `json:"ipv6addr,omitempty"`
	Name     string `json:"name,omitempty"`
	RpZone   string `json:"rp_zone,omitempty"`
	Ttl      uint32 `json:"ttl,omitempty"`
	UseTtl   bool   `json:"use_ttl,omitempty"`
	View     string `json:"view,omitempty"`
	Zone     string `json:"zone,omitempty"`
}

func (RecordRpzAaaaIpaddress) ObjectType() string {
	return "record:rpz:aaaa:ipaddress"
}

func (obj RecordRpzAaaaIpaddress) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"ipv6addr", "name", "view"}
	}
	return obj.returnFields
}

// RecordRpzCname represents Infoblox object record:rpz:cname
type RecordRpzCname struct {
	IBBase    `json:"-"`
	Ref       string `json:"_ref,omitempty"`
	Canonical string `json:"canonical,omitempty"`
	Comment   string `json:"comment,omitempty"`
	Disable   bool   `json:"disable,omitempty"`
	Ea        EA     `json:"extattrs,omitempty"`
	Name      string `json:"name,omitempty"`
	RpZone    string `json:"rp_zone,omitempty"`
	Ttl       uint32 `json:"ttl,omitempty"`
	UseTtl    bool   `json:"use_ttl,omitempty"`
	View      string `json:"view,omitempty"`
	Zone      string `json:"zone,omitempty"`
}

func (RecordRpzCname) ObjectType() string {
	return "record:rpz:cname"
}

func (obj RecordRpzCname) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"canonical", "name", "view"}
	}
	return obj.returnFields
}

// RecordRpzCnameClientipaddress represents Infoblox object record:rpz:cname:clientipaddress
type RecordRpzCnameClientipaddress struct {
	IBBase    `json:"-"`
	Ref       string `json:"_ref,omitempty"`
	Canonical string `json:"canonical,omitempty"`
	Comment   string `json:"comment,omitempty"`
	Disable   bool   `json:"disable,omitempty"`
	Ea        EA     `json:"extattrs,omitempty"`
	IsIpv4    bool   `json:"is_ipv4,omitempty"`
	Name      string `json:"name,omitempty"`
	RpZone    string `json:"rp_zone,omitempty"`
	Ttl       uint32 `json:"ttl,omitempty"`
	UseTtl    bool   `json:"use_ttl,omitempty"`
	View      string `json:"view,omitempty"`
	Zone      string `json:"zone,omitempty"`
}

func (RecordRpzCnameClientipaddress) ObjectType() string {
	return "record:rpz:cname:clientipaddress"
}

func (obj RecordRpzCnameClientipaddress) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"canonical", "name", "view"}
	}
	return obj.returnFields
}

// RecordPTR represents Infoblox object record:ptr
type RecordPTR struct {
	IBBase             `json:"-"`
	Ref                string              `json:"_ref,omitempty"`
	AwsRte53RecordInfo *Awsrte53recordinfo `json:"aws_rte53_record_info,omitempty"`
	CloudInfo          *GridCloudapiInfo   `json:"cloud_info,omitempty"`
	Comment            string              `json:"comment,omitempty"`
	CreationTime       *time.Time          `json:"creation_time,omitempty"`
	Creator            string              `json:"creator,omitempty"`
	DdnsPrincipal      string              `json:"ddns_principal,omitempty"`
	DdnsProtected      bool                `json:"ddns_protected,omitempty"`
	Disable            bool                `json:"disable,omitempty"`
	DiscoveredData     *Discoverydata      `json:"discovered_data,omitempty"`
	DnsName            string              `json:"dns_name,omitempty"`
	DnsPtrdname        string              `json:"dns_ptrdname,omitempty"`
	Ea                 EA                  `json:"extattrs,omitempty"`
	ForbidReclamation  bool                `json:"forbid_reclamation,omitempty"`
	Ipv4Addr           string              `json:"ipv4addr,omitempty"`
	Ipv6Addr           string              `json:"ipv6addr,omitempty"`
	LastQueried        *time.Time          `json:"last_queried,omitempty"`
	MsAdUserData       *MsserverAduserData `json:"ms_ad_user_data,omitempty"`
	Name               string              `json:"name,omitempty"`
	PtrdName           string              `json:"ptrdname,omitempty"`
	Reclaimable        bool                `json:"reclaimable,omitempty"`
	SharedRecordGroup  string              `json:"shared_record_group,omitempty"`
	Ttl                uint32              `json:"ttl,omitempty"`
	UseTtl             bool                `json:"use_ttl,omitempty"`
	View               string              `json:"view,omitempty"`
	Zone               string              `json:"zone,omitempty"`
}

func (RecordPTR) ObjectType() string {
	return "record:ptr"
}

func (obj RecordPTR) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"ptrdname", "view"}
	}
	return obj.returnFields
}

func NewEmptyRecordPTR() *RecordPTR {
	res := RecordPTR{}
	res.returnFields = []string{"extattrs", "ipv4addr", "ipv6addr", "name", "ptrdname", "view", "zone", "comment", "use_ttl", "ttl"}

	return &res
}

func NewRecordPTR(dnsView string, ptrdname string, useTtl bool, ttl uint32, comment string, ea EA) *RecordPTR {
	res := NewEmptyRecordPTR()
	res.View = dnsView
	res.PtrdName = ptrdname
	res.UseTtl = useTtl
	res.Ttl = ttl
	res.Comment = comment
	res.Ea = ea

	return res
}

// RecordRpzCnameClientipaddressdn represents Infoblox object record:rpz:cname:clientipaddressdn
type RecordRpzCnameClientipaddressdn struct {
	IBBase    `json:"-"`
	Ref       string `json:"_ref,omitempty"`
	Canonical string `json:"canonical,omitempty"`
	Comment   string `json:"comment,omitempty"`
	Disable   bool   `json:"disable,omitempty"`
	Ea        EA     `json:"extattrs,omitempty"`
	IsIpv4    bool   `json:"is_ipv4,omitempty"`
	Name      string `json:"name,omitempty"`
	RpZone    string `json:"rp_zone,omitempty"`
	Ttl       uint32 `json:"ttl,omitempty"`
	UseTtl    bool   `json:"use_ttl,omitempty"`
	View      string `json:"view,omitempty"`
	Zone      string `json:"zone,omitempty"`
}

func (RecordRpzCnameClientipaddressdn) ObjectType() string {
	return "record:rpz:cname:clientipaddressdn"
}

func (obj RecordRpzCnameClientipaddressdn) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"canonical", "name", "view"}
	}
	return obj.returnFields
}

// RecordRpzCnameIpaddress represents Infoblox object record:rpz:cname:ipaddress
type RecordRpzCnameIpaddress struct {
	IBBase    `json:"-"`
	Ref       string `json:"_ref,omitempty"`
	Canonical string `json:"canonical,omitempty"`
	Comment   string `json:"comment,omitempty"`
	Disable   bool   `json:"disable,omitempty"`
	Ea        EA     `json:"extattrs,omitempty"`
	IsIpv4    bool   `json:"is_ipv4,omitempty"`
	Name      string `json:"name,omitempty"`
	RpZone    string `json:"rp_zone,omitempty"`
	Ttl       uint32 `json:"ttl,omitempty"`
	UseTtl    bool   `json:"use_ttl,omitempty"`
	View      string `json:"view,omitempty"`
	Zone      string `json:"zone,omitempty"`
}

func (RecordRpzCnameIpaddress) ObjectType() string {
	return "record:rpz:cname:ipaddress"
}

func (obj RecordRpzCnameIpaddress) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"canonical", "name", "view"}
	}
	return obj.returnFields
}

// RecordRpzCnameIpaddressdn represents Infoblox object record:rpz:cname:ipaddressdn
type RecordRpzCnameIpaddressdn struct {
	IBBase    `json:"-"`
	Ref       string `json:"_ref,omitempty"`
	Canonical string `json:"canonical,omitempty"`
	Comment   string `json:"comment,omitempty"`
	Disable   bool   `json:"disable,omitempty"`
	Ea        EA     `json:"extattrs,omitempty"`
	IsIpv4    bool   `json:"is_ipv4,omitempty"`
	Name      string `json:"name,omitempty"`
	RpZone    string `json:"rp_zone,omitempty"`
	Ttl       uint32 `json:"ttl,omitempty"`
	UseTtl    bool   `json:"use_ttl,omitempty"`
	View      string `json:"view,omitempty"`
	Zone      string `json:"zone,omitempty"`
}

func (RecordRpzCnameIpaddressdn) ObjectType() string {
	return "record:rpz:cname:ipaddressdn"
}

func (obj RecordRpzCnameIpaddressdn) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"canonical", "name", "view"}
	}
	return obj.returnFields
}

// RecordRpzMx represents Infoblox object record:rpz:mx
type RecordRpzMx struct {
	IBBase        `json:"-"`
	Ref           string `json:"_ref,omitempty"`
	Comment       string `json:"comment,omitempty"`
	Disable       bool   `json:"disable,omitempty"`
	Ea            EA     `json:"extattrs,omitempty"`
	MailExchanger string `json:"mail_exchanger,omitempty"`
	Name          string `json:"name,omitempty"`
	Preference    uint32 `json:"preference,omitempty"`
	RpZone        string `json:"rp_zone,omitempty"`
	Ttl           uint32 `json:"ttl,omitempty"`
	UseTtl        bool   `json:"use_ttl,omitempty"`
	View          string `json:"view,omitempty"`
	Zone          string `json:"zone,omitempty"`
}

func (RecordRpzMx) ObjectType() string {
	return "record:rpz:mx"
}

func (obj RecordRpzMx) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"mail_exchanger", "name", "preference", "view"}
	}
	return obj.returnFields
}

// RecordRpzPtr represents Infoblox object record:rpz:ptr
type RecordRpzPtr struct {
	IBBase   `json:"-"`
	Ref      string `json:"_ref,omitempty"`
	Comment  string `json:"comment,omitempty"`
	Disable  bool   `json:"disable,omitempty"`
	Ea       EA     `json:"extattrs,omitempty"`
	Ipv4Addr string `json:"ipv4addr,omitempty"`
	Ipv6Addr string `json:"ipv6addr,omitempty"`
	Name     string `json:"name,omitempty"`
	PtrdName string `json:"ptrdname,omitempty"`
	RpZone   string `json:"rp_zone,omitempty"`
	Ttl      uint32 `json:"ttl,omitempty"`
	UseTtl   bool   `json:"use_ttl,omitempty"`
	View     string `json:"view,omitempty"`
	Zone     string `json:"zone,omitempty"`
}

func (RecordRpzPtr) ObjectType() string {
	return "record:rpz:ptr"
}

func (obj RecordRpzPtr) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"ptrdname", "view"}
	}
	return obj.returnFields
}

// RecordRpzNaptr represents Infoblox object record:rpz:naptr
type RecordRpzNaptr struct {
	IBBase      `json:"-"`
	Ref         string     `json:"_ref,omitempty"`
	Comment     string     `json:"comment,omitempty"`
	Disable     bool       `json:"disable,omitempty"`
	Ea          EA         `json:"extattrs,omitempty"`
	Flags       string     `json:"flags,omitempty"`
	LastQueried *time.Time `json:"last_queried,omitempty"`
	Name        string     `json:"name,omitempty"`
	Order       uint32     `json:"order,omitempty"`
	Preference  uint32     `json:"preference,omitempty"`
	Regexp      string     `json:"regexp,omitempty"`
	Replacement string     `json:"replacement,omitempty"`
	RpZone      string     `json:"rp_zone,omitempty"`
	Services    string     `json:"services,omitempty"`
	Ttl         uint32     `json:"ttl,omitempty"`
	UseTtl      bool       `json:"use_ttl,omitempty"`
	View        string     `json:"view,omitempty"`
	Zone        string     `json:"zone,omitempty"`
}

func (RecordRpzNaptr) ObjectType() string {
	return "record:rpz:naptr"
}

func (obj RecordRpzNaptr) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name", "order", "preference", "regexp", "replacement", "services", "view"}
	}
	return obj.returnFields
}

// RecordRpzSrv represents Infoblox object record:rpz:srv
type RecordRpzSrv struct {
	IBBase   `json:"-"`
	Ref      string `json:"_ref,omitempty"`
	Comment  string `json:"comment,omitempty"`
	Disable  bool   `json:"disable,omitempty"`
	Ea       EA     `json:"extattrs,omitempty"`
	Name     string `json:"name,omitempty"`
	Port     uint32 `json:"port,omitempty"`
	Priority uint32 `json:"priority,omitempty"`
	RpZone   string `json:"rp_zone,omitempty"`
	Target   string `json:"target,omitempty"`
	Ttl      uint32 `json:"ttl,omitempty"`
	UseTtl   bool   `json:"use_ttl,omitempty"`
	View     string `json:"view,omitempty"`
	Weight   uint32 `json:"weight,omitempty"`
	Zone     string `json:"zone,omitempty"`
}

func (RecordRpzSrv) ObjectType() string {
	return "record:rpz:srv"
}

func (obj RecordRpzSrv) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name", "port", "priority", "target", "view", "weight"}
	}
	return obj.returnFields
}

// RecordRpzTxt represents Infoblox object record:rpz:txt
type RecordRpzTxt struct {
	IBBase  `json:"-"`
	Ref     string `json:"_ref,omitempty"`
	Comment string `json:"comment,omitempty"`
	Disable bool   `json:"disable,omitempty"`
	Ea      EA     `json:"extattrs,omitempty"`
	Name    string `json:"name,omitempty"`
	RpZone  string `json:"rp_zone,omitempty"`
	Text    string `json:"text,omitempty"`
	Ttl     uint32 `json:"ttl,omitempty"`
	UseTtl  bool   `json:"use_ttl,omitempty"`
	View    string `json:"view,omitempty"`
	Zone    string `json:"zone,omitempty"`
}

func (RecordRpzTxt) ObjectType() string {
	return "record:rpz:txt"
}

func (obj RecordRpzTxt) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name", "text", "view"}
	}
	return obj.returnFields
}

// RecordRrsig represents Infoblox object record:rrsig
type RecordRrsig struct {
	IBBase         `json:"-"`
	Ref            string            `json:"_ref,omitempty"`
	Algorithm      string            `json:"algorithm,omitempty"`
	CloudInfo      *GridCloudapiInfo `json:"cloud_info,omitempty"`
	CreationTime   *time.Time        `json:"creation_time,omitempty"`
	Creator        string            `json:"creator,omitempty"`
	DnsName        string            `json:"dns_name,omitempty"`
	DnsSignerName  string            `json:"dns_signer_name,omitempty"`
	ExpirationTime *time.Time        `json:"expiration_time,omitempty"`
	InceptionTime  *time.Time        `json:"inception_time,omitempty"`
	KeyTag         uint32            `json:"key_tag,omitempty"`
	Labels         uint32            `json:"labels,omitempty"`
	LastQueried    *time.Time        `json:"last_queried,omitempty"`
	Name           string            `json:"name,omitempty"`
	OriginalTtl    uint32            `json:"original_ttl,omitempty"`
	Signature      string            `json:"signature,omitempty"`
	SignerName     string            `json:"signer_name,omitempty"`
	Ttl            uint32            `json:"ttl,omitempty"`
	TypeCovered    string            `json:"type_covered,omitempty"`
	UseTtl         bool              `json:"use_ttl,omitempty"`
	View           string            `json:"view,omitempty"`
	Zone           string            `json:"zone,omitempty"`
}

func (RecordRrsig) ObjectType() string {
	return "record:rrsig"
}

func (obj RecordRrsig) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name", "view"}
	}
	return obj.returnFields
}

// RecordSrv represents Infoblox object record:srv
type RecordSrv struct {
	IBBase             `json:"-"`
	Ref                string              `json:"_ref,omitempty"`
	AwsRte53RecordInfo *Awsrte53recordinfo `json:"aws_rte53_record_info,omitempty"`
	CloudInfo          *GridCloudapiInfo   `json:"cloud_info,omitempty"`
	Comment            string              `json:"comment,omitempty"`
	CreationTime       *time.Time          `json:"creation_time,omitempty"`
	Creator            string              `json:"creator,omitempty"`
	DdnsPrincipal      string              `json:"ddns_principal,omitempty"`
	DdnsProtected      bool                `json:"ddns_protected,omitempty"`
	Disable            bool                `json:"disable,omitempty"`
	DnsName            string              `json:"dns_name,omitempty"`
	DnsTarget          string              `json:"dns_target,omitempty"`
	Ea                 EA                  `json:"extattrs,omitempty"`
	ForbidReclamation  bool                `json:"forbid_reclamation,omitempty"`
	LastQueried        *time.Time          `json:"last_queried,omitempty"`
	Name               string              `json:"name,omitempty"`
	Port               uint32              `json:"port,omitempty"`
	Priority           uint32              `json:"priority,omitempty"`
	Reclaimable        bool                `json:"reclaimable,omitempty"`
	SharedRecordGroup  string              `json:"shared_record_group,omitempty"`
	Target             string              `json:"target,omitempty"`
	Ttl                uint32              `json:"ttl,omitempty"`
	UseTtl             bool                `json:"use_ttl,omitempty"`
	View               string              `json:"view,omitempty"`
	Weight             uint32              `json:"weight,omitempty"`
	Zone               string              `json:"zone,omitempty"`
}

func (RecordSrv) ObjectType() string {
	return "record:srv"
}

func (obj RecordSrv) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name", "port", "priority", "target", "view", "weight"}
	}
	return obj.returnFields
}

// RecordUnknown represents Infoblox object record:unknown
type RecordUnknown struct {
	IBBase               `json:"-"`
	Ref                  string            `json:"_ref,omitempty"`
	CloudInfo            *GridCloudapiInfo `json:"cloud_info,omitempty"`
	Comment              string            `json:"comment,omitempty"`
	Creator              string            `json:"creator,omitempty"`
	Disable              bool              `json:"disable,omitempty"`
	DisplayRdata         string            `json:"display_rdata,omitempty"`
	DnsName              string            `json:"dns_name,omitempty"`
	EnableHostNamePolicy bool              `json:"enable_host_name_policy,omitempty"`
	Ea                   EA                `json:"extattrs,omitempty"`
	LastQueried          *time.Time        `json:"last_queried,omitempty"`
	Name                 string            `json:"name,omitempty"`
	Policy               string            `json:"policy,omitempty"`
	RecordType           string            `json:"record_type,omitempty"`
	SubfieldValues       []*Rdatasubfield  `json:"subfield_values,omitempty"`
	Ttl                  uint32            `json:"ttl,omitempty"`
	UseTtl               bool              `json:"use_ttl,omitempty"`
	View                 string            `json:"view,omitempty"`
	Zone                 string            `json:"zone,omitempty"`
}

func (RecordUnknown) ObjectType() string {
	return "record:unknown"
}

func (obj RecordUnknown) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name", "view"}
	}
	return obj.returnFields
}

// RecordTXT represents Infoblox object record:txt
type RecordTXT struct {
	IBBase             `json:"-"`
	Ref                string              `json:"_ref,omitempty"`
	AwsRte53RecordInfo *Awsrte53recordinfo `json:"aws_rte53_record_info,omitempty"`
	CloudInfo          *GridCloudapiInfo   `json:"cloud_info,omitempty"`
	Comment            string              `json:"comment,omitempty"`
	CreationTime       *time.Time          `json:"creation_time,omitempty"`
	Creator            string              `json:"creator,omitempty"`
	DdnsPrincipal      string              `json:"ddns_principal,omitempty"`
	DdnsProtected      bool                `json:"ddns_protected,omitempty"`
	Disable            bool                `json:"disable,omitempty"`
	DnsName            string              `json:"dns_name,omitempty"`
	Ea                 EA                  `json:"extattrs,omitempty"`
	ForbidReclamation  bool                `json:"forbid_reclamation,omitempty"`
	LastQueried        *time.Time          `json:"last_queried,omitempty"`
	Name               string              `json:"name,omitempty"`
	Reclaimable        bool                `json:"reclaimable,omitempty"`
	SharedRecordGroup  string              `json:"shared_record_group,omitempty"`
	Text               string              `json:"text,omitempty"`
	Ttl                uint32              `json:"ttl,omitempty"`
	UseTtl             bool                `json:"use_ttl,omitempty"`
	View               string              `json:"view,omitempty"`
	Zone               string              `json:"zone,omitempty"`
}

func (RecordTXT) ObjectType() string {
	return "record:txt"
}

func (obj RecordTXT) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name", "text", "view"}
	}
	return obj.returnFields
}

func NewEmptyRecordTXT() *RecordTXT {
	res := RecordTXT{}
	res.returnFields = []string{"view", "zone", "name", "text", "ttl", "use_ttl", "comment", "extattrs"}

	return &res
}

func NewRecordTXT(
	dnsview string,
	zone string,
	recordname string,
	text string,
	ttl uint32,
	useTtl bool,
	comment string,
	eas EA) *RecordTXT {

	res := NewEmptyRecordTXT()
	res.View = dnsview
	res.Zone = zone
	res.Name = recordname
	res.Text = text
	res.Ttl = ttl
	res.UseTtl = useTtl
	res.Comment = comment
	res.Ea = eas

	return res
}

// RecordTlsa represents Infoblox object record:tlsa
type RecordTlsa struct {
	IBBase           `json:"-"`
	Ref              string            `json:"_ref,omitempty"`
	CertificateData  string            `json:"certificate_data,omitempty"`
	CertificateUsage uint32            `json:"certificate_usage,omitempty"`
	CloudInfo        *GridCloudapiInfo `json:"cloud_info,omitempty"`
	Comment          string            `json:"comment,omitempty"`
	Creator          string            `json:"creator,omitempty"`
	Disable          bool              `json:"disable,omitempty"`
	DnsName          string            `json:"dns_name,omitempty"`
	Ea               EA                `json:"extattrs,omitempty"`
	LastQueried      *time.Time        `json:"last_queried,omitempty"`
	MatchedType      uint32            `json:"matched_type,omitempty"`
	Name             string            `json:"name,omitempty"`
	Selector         uint32            `json:"selector,omitempty"`
	Ttl              uint32            `json:"ttl,omitempty"`
	UseTtl           bool              `json:"use_ttl,omitempty"`
	View             string            `json:"view,omitempty"`
	Zone             string            `json:"zone,omitempty"`
}

func (RecordTlsa) ObjectType() string {
	return "record:tlsa"
}

func (obj RecordTlsa) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name", "view"}
	}
	return obj.returnFields
}

// Restartservicestatus represents Infoblox object restartservicestatus
type Restartservicestatus struct {
	IBBase          `json:"-"`
	Ref             string `json:"_ref,omitempty"`
	DhcpStatus      string `json:"dhcp_status,omitempty"`
	DnsStatus       string `json:"dns_status,omitempty"`
	Member          string `json:"member,omitempty"`
	ReportingStatus string `json:"reporting_status,omitempty"`
}

func (Restartservicestatus) ObjectType() string {
	return "restartservicestatus"
}

func (obj Restartservicestatus) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"dhcp_status", "dns_status", "member", "reporting_status"}
	}
	return obj.returnFields
}

// Recordnamepolicy represents Infoblox object recordnamepolicy
type Recordnamepolicy struct {
	IBBase     `json:"-"`
	Ref        string `json:"_ref,omitempty"`
	IsDefault  bool   `json:"is_default,omitempty"`
	Name       string `json:"name,omitempty"`
	PreDefined bool   `json:"pre_defined,omitempty"`
	Regex      string `json:"regex,omitempty"`
}

func (Recordnamepolicy) ObjectType() string {
	return "recordnamepolicy"
}

func (obj Recordnamepolicy) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"is_default", "name", "regex"}
	}
	return obj.returnFields
}

// RirOrganization represents Infoblox object rir:organization
type RirOrganization struct {
	IBBase      `json:"-"`
	Ref         string `json:"_ref,omitempty"`
	Ea          EA     `json:"extattrs,omitempty"`
	Id          string `json:"id,omitempty"`
	Maintainer  string `json:"maintainer,omitempty"`
	Name        string `json:"name,omitempty"`
	Password    string `json:"password,omitempty"`
	Rir         string `json:"rir,omitempty"`
	SenderEmail string `json:"sender_email,omitempty"`
}

func (RirOrganization) ObjectType() string {
	return "rir:organization"
}

func (obj RirOrganization) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"id", "maintainer", "name", "rir", "sender_email"}
	}
	return obj.returnFields
}

// Rir represents Infoblox object rir
type Rir struct {
	IBBase            `json:"-"`
	Ref               string `json:"_ref,omitempty"`
	CommunicationMode string `json:"communication_mode,omitempty"`
	Email             string `json:"email,omitempty"`
	Name              string `json:"name,omitempty"`
	Url               string `json:"url,omitempty"`
	UseEmail          bool   `json:"use_email,omitempty"`
	UseUrl            bool   `json:"use_url,omitempty"`
}

func (Rir) ObjectType() string {
	return "rir"
}

func (obj Rir) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"communication_mode", "email", "name", "url"}
	}
	return obj.returnFields
}

// Roaminghost represents Infoblox object roaminghost
type Roaminghost struct {
	IBBase                         `json:"-"`
	Ref                            string        `json:"_ref,omitempty"`
	AddressType                    string        `json:"address_type,omitempty"`
	Bootfile                       string        `json:"bootfile,omitempty"`
	Bootserver                     string        `json:"bootserver,omitempty"`
	ClientIdentifierPrependZero    bool          `json:"client_identifier_prepend_zero,omitempty"`
	Comment                        string        `json:"comment,omitempty"`
	DdnsDomainname                 string        `json:"ddns_domainname,omitempty"`
	DdnsHostname                   string        `json:"ddns_hostname,omitempty"`
	DenyBootp                      bool          `json:"deny_bootp,omitempty"`
	DhcpClientIdentifier           string        `json:"dhcp_client_identifier,omitempty"`
	Disable                        bool          `json:"disable,omitempty"`
	EnableDdns                     bool          `json:"enable_ddns,omitempty"`
	EnablePxeLeaseTime             bool          `json:"enable_pxe_lease_time,omitempty"`
	Ea                             EA            `json:"extattrs,omitempty"`
	ForceRoamingHostname           bool          `json:"force_roaming_hostname,omitempty"`
	IgnoreDhcpOptionListRequest    bool          `json:"ignore_dhcp_option_list_request,omitempty"`
	Ipv6ClientHostname             string        `json:"ipv6_client_hostname,omitempty"`
	Ipv6DdnsDomainname             string        `json:"ipv6_ddns_domainname,omitempty"`
	Ipv6DdnsHostname               string        `json:"ipv6_ddns_hostname,omitempty"`
	Ipv6DomainName                 string        `json:"ipv6_domain_name,omitempty"`
	Ipv6DomainNameServers          []string      `json:"ipv6_domain_name_servers,omitempty"`
	Ipv6Duid                       string        `json:"ipv6_duid,omitempty"`
	Ipv6EnableDdns                 bool          `json:"ipv6_enable_ddns,omitempty"`
	Ipv6ForceRoamingHostname       bool          `json:"ipv6_force_roaming_hostname,omitempty"`
	Ipv6MatchOption                string        `json:"ipv6_match_option,omitempty"`
	Ipv6Options                    []*Dhcpoption `json:"ipv6_options,omitempty"`
	Ipv6Template                   string        `json:"ipv6_template,omitempty"`
	Mac                            string        `json:"mac,omitempty"`
	MatchClient                    string        `json:"match_client,omitempty"`
	Name                           string        `json:"name,omitempty"`
	NetworkView                    string        `json:"network_view,omitempty"`
	Nextserver                     string        `json:"nextserver,omitempty"`
	Options                        []*Dhcpoption `json:"options,omitempty"`
	PreferredLifetime              uint32        `json:"preferred_lifetime,omitempty"`
	PxeLeaseTime                   uint32        `json:"pxe_lease_time,omitempty"`
	Template                       string        `json:"template,omitempty"`
	UseBootfile                    bool          `json:"use_bootfile,omitempty"`
	UseBootserver                  bool          `json:"use_bootserver,omitempty"`
	UseDdnsDomainname              bool          `json:"use_ddns_domainname,omitempty"`
	UseDenyBootp                   bool          `json:"use_deny_bootp,omitempty"`
	UseEnableDdns                  bool          `json:"use_enable_ddns,omitempty"`
	UseIgnoreDhcpOptionListRequest bool          `json:"use_ignore_dhcp_option_list_request,omitempty"`
	UseIpv6DdnsDomainname          bool          `json:"use_ipv6_ddns_domainname,omitempty"`
	UseIpv6DomainName              bool          `json:"use_ipv6_domain_name,omitempty"`
	UseIpv6DomainNameServers       bool          `json:"use_ipv6_domain_name_servers,omitempty"`
	UseIpv6EnableDdns              bool          `json:"use_ipv6_enable_ddns,omitempty"`
	UseIpv6Options                 bool          `json:"use_ipv6_options,omitempty"`
	UseNextserver                  bool          `json:"use_nextserver,omitempty"`
	UseOptions                     bool          `json:"use_options,omitempty"`
	UsePreferredLifetime           bool          `json:"use_preferred_lifetime,omitempty"`
	UsePxeLeaseTime                bool          `json:"use_pxe_lease_time,omitempty"`
	UseValidLifetime               bool          `json:"use_valid_lifetime,omitempty"`
	ValidLifetime                  uint32        `json:"valid_lifetime,omitempty"`
}

func (Roaminghost) ObjectType() string {
	return "roaminghost"
}

func (obj Roaminghost) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"address_type", "name", "network_view"}
	}
	return obj.returnFields
}

// SamlAuthservice represents Infoblox object saml:authservice
type SamlAuthservice struct {
	IBBase         `json:"-"`
	Ref            string   `json:"_ref,omitempty"`
	Comment        string   `json:"comment,omitempty"`
	Idp            *SamlIdp `json:"idp,omitempty"`
	Name           string   `json:"name,omitempty"`
	SessionTimeout uint32   `json:"session_timeout,omitempty"`
}

func (SamlAuthservice) ObjectType() string {
	return "saml:authservice"
}

func (obj SamlAuthservice) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name"}
	}
	return obj.returnFields
}

// Ruleset represents Infoblox object ruleset
type Ruleset struct {
	IBBase        `json:"-"`
	Ref           string          `json:"_ref,omitempty"`
	Comment       string          `json:"comment,omitempty"`
	Disabled      bool            `json:"disabled,omitempty"`
	Name          string          `json:"name,omitempty"`
	NxdomainRules []*Nxdomainrule `json:"nxdomain_rules,omitempty"`
	Type          string          `json:"type,omitempty"`
}

func (Ruleset) ObjectType() string {
	return "ruleset"
}

func (obj Ruleset) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "disabled", "name", "type"}
	}
	return obj.returnFields
}

// Scavengingtask represents Infoblox object scavengingtask
type Scavengingtask struct {
	IBBase             `json:"-"`
	Ref                string     `json:"_ref,omitempty"`
	Action             string     `json:"action,omitempty"`
	AssociatedObject   string     `json:"associated_object,omitempty"`
	EndTime            *time.Time `json:"end_time,omitempty"`
	ProcessedRecords   uint32     `json:"processed_records,omitempty"`
	ReclaimableRecords uint32     `json:"reclaimable_records,omitempty"`
	ReclaimedRecords   uint32     `json:"reclaimed_records,omitempty"`
	StartTime          *time.Time `json:"start_time,omitempty"`
	Status             string     `json:"status,omitempty"`
}

func (Scavengingtask) ObjectType() string {
	return "scavengingtask"
}

func (obj Scavengingtask) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"action", "associated_object", "status"}
	}
	return obj.returnFields
}

// Scheduledtask represents Infoblox object scheduledtask
type Scheduledtask struct {
	IBBase               `json:"-"`
	Ref                  string           `json:"_ref,omitempty"`
	ApprovalStatus       string           `json:"approval_status,omitempty"`
	Approver             string           `json:"approver,omitempty"`
	ApproverComment      string           `json:"approver_comment,omitempty"`
	AutomaticRestart     bool             `json:"automatic_restart,omitempty"`
	ChangedObjects       []*Changedobject `json:"changed_objects,omitempty"`
	DependentTasks       []*Scheduledtask `json:"dependent_tasks,omitempty"`
	ExecuteNow           bool             `json:"execute_now,omitempty"`
	ExecutionDetails     []string         `json:"execution_details,omitempty"`
	ExecutionDetailsType string           `json:"execution_details_type,omitempty"`
	ExecutionStatus      string           `json:"execution_status,omitempty"`
	ExecutionTime        *time.Time       `json:"execution_time,omitempty"`
	IsNetworkInsightTask bool             `json:"is_network_insight_task,omitempty"`
	Member               string           `json:"member,omitempty"`
	PredecessorTask      string           `json:"predecessor_task,omitempty"`
	ReExecuteTask        bool             `json:"re_execute_task,omitempty"`
	ScheduledTime        *time.Time       `json:"scheduled_time,omitempty"`
	SubmitTime           *time.Time       `json:"submit_time,omitempty"`
	Submitter            string           `json:"submitter,omitempty"`
	SubmitterComment     string           `json:"submitter_comment,omitempty"`
	TaskId               uint32           `json:"task_id,omitempty"`
	TaskType             string           `json:"task_type,omitempty"`
	TicketNumber         string           `json:"ticket_number,omitempty"`
}

func (Scheduledtask) ObjectType() string {
	return "scheduledtask"
}

func (obj Scheduledtask) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"approval_status", "execution_status", "task_id"}
	}
	return obj.returnFields
}

// Search represents Infoblox object search
type Search struct {
	IBBase `json:"-"`
	Ref    string `json:"_ref,omitempty"`
}

func (Search) ObjectType() string {
	return "search"
}

func (obj Search) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{}
	}
	return obj.returnFields
}

// SharedrecordA represents Infoblox object sharedrecord:a
type SharedrecordA struct {
	IBBase            `json:"-"`
	Ref               string `json:"_ref,omitempty"`
	Comment           string `json:"comment,omitempty"`
	Disable           bool   `json:"disable,omitempty"`
	DnsName           string `json:"dns_name,omitempty"`
	Ea                EA     `json:"extattrs,omitempty"`
	Ipv4Addr          string `json:"ipv4addr,omitempty"`
	Name              string `json:"name,omitempty"`
	SharedRecordGroup string `json:"shared_record_group,omitempty"`
	Ttl               uint32 `json:"ttl,omitempty"`
	UseTtl            bool   `json:"use_ttl,omitempty"`
}

func (SharedrecordA) ObjectType() string {
	return "sharedrecord:a"
}

func (obj SharedrecordA) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"ipv4addr", "name", "shared_record_group"}
	}
	return obj.returnFields
}

// Sharednetwork represents Infoblox object sharednetwork
type Sharednetwork struct {
	IBBase                         `json:"-"`
	Ref                            string              `json:"_ref,omitempty"`
	Authority                      bool                `json:"authority,omitempty"`
	Bootfile                       string              `json:"bootfile,omitempty"`
	Bootserver                     string              `json:"bootserver,omitempty"`
	Comment                        string              `json:"comment,omitempty"`
	DdnsGenerateHostname           bool                `json:"ddns_generate_hostname,omitempty"`
	DdnsServerAlwaysUpdates        bool                `json:"ddns_server_always_updates,omitempty"`
	DdnsTtl                        uint32              `json:"ddns_ttl,omitempty"`
	DdnsUpdateFixedAddresses       bool                `json:"ddns_update_fixed_addresses,omitempty"`
	DdnsUseOption81                bool                `json:"ddns_use_option81,omitempty"`
	DenyBootp                      bool                `json:"deny_bootp,omitempty"`
	DhcpUtilization                uint32              `json:"dhcp_utilization,omitempty"`
	DhcpUtilizationStatus          string              `json:"dhcp_utilization_status,omitempty"`
	Disable                        bool                `json:"disable,omitempty"`
	DynamicHosts                   uint32              `json:"dynamic_hosts,omitempty"`
	EnableDdns                     bool                `json:"enable_ddns,omitempty"`
	EnablePxeLeaseTime             bool                `json:"enable_pxe_lease_time,omitempty"`
	Ea                             EA                  `json:"extattrs,omitempty"`
	IgnoreClientIdentifier         bool                `json:"ignore_client_identifier,omitempty"`
	IgnoreDhcpOptionListRequest    bool                `json:"ignore_dhcp_option_list_request,omitempty"`
	IgnoreId                       string              `json:"ignore_id,omitempty"`
	IgnoreMacAddresses             []string            `json:"ignore_mac_addresses,omitempty"`
	LeaseScavengeTime              int                 `json:"lease_scavenge_time,omitempty"`
	LogicFilterRules               []*Logicfilterrule  `json:"logic_filter_rules,omitempty"`
	MsAdUserData                   *MsserverAduserData `json:"ms_ad_user_data,omitempty"`
	Name                           string              `json:"name,omitempty"`
	NetworkView                    string              `json:"network_view,omitempty"`
	Networks                       []*Ipv4Network      `json:"networks,omitempty"`
	Nextserver                     string              `json:"nextserver,omitempty"`
	Options                        []*Dhcpoption       `json:"options,omitempty"`
	PxeLeaseTime                   uint32              `json:"pxe_lease_time,omitempty"`
	StaticHosts                    uint32              `json:"static_hosts,omitempty"`
	TotalHosts                     uint32              `json:"total_hosts,omitempty"`
	UpdateDnsOnLeaseRenewal        bool                `json:"update_dns_on_lease_renewal,omitempty"`
	UseAuthority                   bool                `json:"use_authority,omitempty"`
	UseBootfile                    bool                `json:"use_bootfile,omitempty"`
	UseBootserver                  bool                `json:"use_bootserver,omitempty"`
	UseDdnsGenerateHostname        bool                `json:"use_ddns_generate_hostname,omitempty"`
	UseDdnsTtl                     bool                `json:"use_ddns_ttl,omitempty"`
	UseDdnsUpdateFixedAddresses    bool                `json:"use_ddns_update_fixed_addresses,omitempty"`
	UseDdnsUseOption81             bool                `json:"use_ddns_use_option81,omitempty"`
	UseDenyBootp                   bool                `json:"use_deny_bootp,omitempty"`
	UseEnableDdns                  bool                `json:"use_enable_ddns,omitempty"`
	UseIgnoreClientIdentifier      bool                `json:"use_ignore_client_identifier,omitempty"`
	UseIgnoreDhcpOptionListRequest bool                `json:"use_ignore_dhcp_option_list_request,omitempty"`
	UseIgnoreId                    bool                `json:"use_ignore_id,omitempty"`
	UseLeaseScavengeTime           bool                `json:"use_lease_scavenge_time,omitempty"`
	UseLogicFilterRules            bool                `json:"use_logic_filter_rules,omitempty"`
	UseNextserver                  bool                `json:"use_nextserver,omitempty"`
	UseOptions                     bool                `json:"use_options,omitempty"`
	UsePxeLeaseTime                bool                `json:"use_pxe_lease_time,omitempty"`
	UseUpdateDnsOnLeaseRenewal     bool                `json:"use_update_dns_on_lease_renewal,omitempty"`
}

func (Sharednetwork) ObjectType() string {
	return "sharednetwork"
}

func (obj Sharednetwork) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name", "network_view", "networks"}
	}
	return obj.returnFields
}

// SharedrecordAaaa represents Infoblox object sharedrecord:aaaa
type SharedrecordAaaa struct {
	IBBase            `json:"-"`
	Ref               string `json:"_ref,omitempty"`
	Comment           string `json:"comment,omitempty"`
	Disable           bool   `json:"disable,omitempty"`
	DnsName           string `json:"dns_name,omitempty"`
	Ea                EA     `json:"extattrs,omitempty"`
	Ipv6Addr          string `json:"ipv6addr,omitempty"`
	Name              string `json:"name,omitempty"`
	SharedRecordGroup string `json:"shared_record_group,omitempty"`
	Ttl               uint32 `json:"ttl,omitempty"`
	UseTtl            bool   `json:"use_ttl,omitempty"`
}

func (SharedrecordAaaa) ObjectType() string {
	return "sharedrecord:aaaa"
}

func (obj SharedrecordAaaa) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"ipv6addr", "name", "shared_record_group"}
	}
	return obj.returnFields
}

// SharedrecordCname represents Infoblox object sharedrecord:cname
type SharedrecordCname struct {
	IBBase            `json:"-"`
	Ref               string `json:"_ref,omitempty"`
	Canonical         string `json:"canonical,omitempty"`
	Comment           string `json:"comment,omitempty"`
	Disable           bool   `json:"disable,omitempty"`
	DnsCanonical      string `json:"dns_canonical,omitempty"`
	DnsName           string `json:"dns_name,omitempty"`
	Ea                EA     `json:"extattrs,omitempty"`
	Name              string `json:"name,omitempty"`
	SharedRecordGroup string `json:"shared_record_group,omitempty"`
	Ttl               uint32 `json:"ttl,omitempty"`
	UseTtl            bool   `json:"use_ttl,omitempty"`
}

func (SharedrecordCname) ObjectType() string {
	return "sharedrecord:cname"
}

func (obj SharedrecordCname) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"canonical", "name", "shared_record_group"}
	}
	return obj.returnFields
}

// SharedrecordMx represents Infoblox object sharedrecord:mx
type SharedrecordMx struct {
	IBBase            `json:"-"`
	Ref               string `json:"_ref,omitempty"`
	Comment           string `json:"comment,omitempty"`
	Disable           bool   `json:"disable,omitempty"`
	DnsMailExchanger  string `json:"dns_mail_exchanger,omitempty"`
	DnsName           string `json:"dns_name,omitempty"`
	Ea                EA     `json:"extattrs,omitempty"`
	MailExchanger     string `json:"mail_exchanger,omitempty"`
	Name              string `json:"name,omitempty"`
	Preference        uint32 `json:"preference,omitempty"`
	SharedRecordGroup string `json:"shared_record_group,omitempty"`
	Ttl               uint32 `json:"ttl,omitempty"`
	UseTtl            bool   `json:"use_ttl,omitempty"`
}

func (SharedrecordMx) ObjectType() string {
	return "sharedrecord:mx"
}

func (obj SharedrecordMx) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"mail_exchanger", "name", "preference", "shared_record_group"}
	}
	return obj.returnFields
}

// SharedrecordTxt represents Infoblox object sharedrecord:txt
type SharedrecordTxt struct {
	IBBase            `json:"-"`
	Ref               string `json:"_ref,omitempty"`
	Comment           string `json:"comment,omitempty"`
	Disable           bool   `json:"disable,omitempty"`
	DnsName           string `json:"dns_name,omitempty"`
	Ea                EA     `json:"extattrs,omitempty"`
	Name              string `json:"name,omitempty"`
	SharedRecordGroup string `json:"shared_record_group,omitempty"`
	Text              string `json:"text,omitempty"`
	Ttl               uint32 `json:"ttl,omitempty"`
	UseTtl            bool   `json:"use_ttl,omitempty"`
}

func (SharedrecordTxt) ObjectType() string {
	return "sharedrecord:txt"
}

func (obj SharedrecordTxt) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name", "shared_record_group", "text"}
	}
	return obj.returnFields
}

// Sharedrecordgroup represents Infoblox object sharedrecordgroup
type Sharedrecordgroup struct {
	IBBase              `json:"-"`
	Ref                 string   `json:"_ref,omitempty"`
	Comment             string   `json:"comment,omitempty"`
	Ea                  EA       `json:"extattrs,omitempty"`
	Name                string   `json:"name,omitempty"`
	RecordNamePolicy    string   `json:"record_name_policy,omitempty"`
	UseRecordNamePolicy bool     `json:"use_record_name_policy,omitempty"`
	ZoneAssociations    []string `json:"zone_associations,omitempty"`
}

func (Sharedrecordgroup) ObjectType() string {
	return "sharedrecordgroup"
}

func (obj Sharedrecordgroup) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name"}
	}
	return obj.returnFields
}

// SharedrecordSrv represents Infoblox object sharedrecord:srv
type SharedrecordSrv struct {
	IBBase            `json:"-"`
	Ref               string `json:"_ref,omitempty"`
	Comment           string `json:"comment,omitempty"`
	Disable           bool   `json:"disable,omitempty"`
	DnsName           string `json:"dns_name,omitempty"`
	DnsTarget         string `json:"dns_target,omitempty"`
	Ea                EA     `json:"extattrs,omitempty"`
	Name              string `json:"name,omitempty"`
	Port              uint32 `json:"port,omitempty"`
	Priority          uint32 `json:"priority,omitempty"`
	SharedRecordGroup string `json:"shared_record_group,omitempty"`
	Target            string `json:"target,omitempty"`
	Ttl               uint32 `json:"ttl,omitempty"`
	UseTtl            bool   `json:"use_ttl,omitempty"`
	Weight            uint32 `json:"weight,omitempty"`
}

func (SharedrecordSrv) ObjectType() string {
	return "sharedrecord:srv"
}

func (obj SharedrecordSrv) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name", "port", "priority", "shared_record_group", "target", "weight"}
	}
	return obj.returnFields
}

// SmartfolderChildren represents Infoblox object smartfolder:children
type SmartfolderChildren struct {
	IBBase    `json:"-"`
	Ref       string                     `json:"_ref,omitempty"`
	Resource  string                     `json:"resource,omitempty"`
	Value     *SmartfolderQueryitemvalue `json:"value,omitempty"`
	ValueType string                     `json:"value_type,omitempty"`
}

func (SmartfolderChildren) ObjectType() string {
	return "smartfolder:children"
}

func (obj SmartfolderChildren) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"resource", "value", "value_type"}
	}
	return obj.returnFields
}

// SmartfolderGlobal represents Infoblox object smartfolder:global
type SmartfolderGlobal struct {
	IBBase     `json:"-"`
	Ref        string                   `json:"_ref,omitempty"`
	Comment    string                   `json:"comment,omitempty"`
	GroupBys   []*SmartfolderGroupby    `json:"group_bys,omitempty"`
	Name       string                   `json:"name,omitempty"`
	QueryItems []*SmartfolderQueryitem  `json:"query_items,omitempty"`
	SaveAs     *Smartfoldersaveasparams `json:"save_as,omitempty"`
}

func (SmartfolderGlobal) ObjectType() string {
	return "smartfolder:global"
}

func (obj SmartfolderGlobal) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name"}
	}
	return obj.returnFields
}

// SmartfolderPersonal represents Infoblox object smartfolder:personal
type SmartfolderPersonal struct {
	IBBase     `json:"-"`
	Ref        string                   `json:"_ref,omitempty"`
	Comment    string                   `json:"comment,omitempty"`
	GroupBys   []*SmartfolderGroupby    `json:"group_bys,omitempty"`
	IsShortcut bool                     `json:"is_shortcut,omitempty"`
	Name       string                   `json:"name,omitempty"`
	QueryItems []*SmartfolderQueryitem  `json:"query_items,omitempty"`
	SaveAs     *Smartfoldersaveasparams `json:"save_as,omitempty"`
}

func (SmartfolderPersonal) ObjectType() string {
	return "smartfolder:personal"
}

func (obj SmartfolderPersonal) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "is_shortcut", "name"}
	}
	return obj.returnFields
}

// Snmpuser represents Infoblox object snmpuser
type Snmpuser struct {
	IBBase                 `json:"-"`
	Ref                    string `json:"_ref,omitempty"`
	AuthenticationPassword string `json:"authentication_password,omitempty"`
	AuthenticationProtocol string `json:"authentication_protocol,omitempty"`
	Comment                string `json:"comment,omitempty"`
	Disable                bool   `json:"disable,omitempty"`
	Ea                     EA     `json:"extattrs,omitempty"`
	Name                   string `json:"name,omitempty"`
	PrivacyPassword        string `json:"privacy_password,omitempty"`
	PrivacyProtocol        string `json:"privacy_protocol,omitempty"`
}

func (Snmpuser) ObjectType() string {
	return "snmpuser"
}

func (obj Snmpuser) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name"}
	}
	return obj.returnFields
}

// Superhostchild represents Infoblox object superhostchild
type Superhostchild struct {
	IBBase            `json:"-"`
	Ref               string     `json:"_ref,omitempty"`
	AssociatedObject  string     `json:"associated_object,omitempty"`
	Comment           string     `json:"comment,omitempty"`
	CreationTimestamp *time.Time `json:"creation_timestamp,omitempty"`
	Data              string     `json:"data,omitempty"`
	Disabled          bool       `json:"disabled,omitempty"`
	Name              string     `json:"name,omitempty"`
	NetworkView       string     `json:"network_view,omitempty"`
	Parent            string     `json:"parent,omitempty"`
	RecordParent      string     `json:"record_parent,omitempty"`
	Type              string     `json:"type,omitempty"`
	View              string     `json:"view,omitempty"`
}

func (Superhostchild) ObjectType() string {
	return "superhostchild"
}

func (obj Superhostchild) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "data", "name", "network_view", "parent", "record_parent", "type", "view"}
	}
	return obj.returnFields
}

// Superhost represents Infoblox object superhost
type Superhost struct {
	IBBase                  `json:"-"`
	Ref                     string              `json:"_ref,omitempty"`
	Comment                 string              `json:"comment,omitempty"`
	DeleteAssociatedObjects bool                `json:"delete_associated_objects,omitempty"`
	DhcpAssociatedObjects   []*Ipv4FixedAddress `json:"dhcp_associated_objects,omitempty"`
	Disabled                bool                `json:"disabled,omitempty"`
	DnsAssociatedObjects    []*RecordA          `json:"dns_associated_objects,omitempty"`
	Ea                      EA                  `json:"extattrs,omitempty"`
	Name                    string              `json:"name,omitempty"`
}

func (Superhost) ObjectType() string {
	return "superhost"
}

func (obj Superhost) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name"}
	}
	return obj.returnFields
}

// Taxii represents Infoblox object taxii
type Taxii struct {
	IBBase         `json:"-"`
	Ref            string            `json:"_ref,omitempty"`
	EnableService  bool              `json:"enable_service,omitempty"`
	Ipv4Addr       string            `json:"ipv4addr,omitempty"`
	Ipv6Addr       string            `json:"ipv6addr,omitempty"`
	Name           string            `json:"name,omitempty"`
	TaxiiRpzConfig []*TaxiiRpzconfig `json:"taxii_rpz_config,omitempty"`
}

func (Taxii) ObjectType() string {
	return "taxii"
}

func (obj Taxii) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"ipv4addr", "ipv6addr", "name"}
	}
	return obj.returnFields
}

// TacacsplusAuthservice represents Infoblox object tacacsplus:authservice
type TacacsplusAuthservice struct {
	IBBase                        `json:"-"`
	Ref                           string                         `json:"_ref,omitempty"`
	AcctRetries                   uint32                         `json:"acct_retries,omitempty"`
	AcctTimeout                   uint32                         `json:"acct_timeout,omitempty"`
	AuthRetries                   uint32                         `json:"auth_retries,omitempty"`
	AuthTimeout                   uint32                         `json:"auth_timeout,omitempty"`
	CheckTacacsplusServerSettings *Checktacacsplusserversettings `json:"check_tacacsplus_server_settings,omitempty"`
	Comment                       string                         `json:"comment,omitempty"`
	Disable                       bool                           `json:"disable,omitempty"`
	Name                          string                         `json:"name,omitempty"`
	Servers                       []*TacacsplusServer            `json:"servers,omitempty"`
}

func (TacacsplusAuthservice) ObjectType() string {
	return "tacacsplus:authservice"
}

func (obj TacacsplusAuthservice) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "disable", "name"}
	}
	return obj.returnFields
}

// Tftpfiledir represents Infoblox object tftpfiledir
type Tftpfiledir struct {
	IBBase          `json:"-"`
	Ref             string            `json:"_ref,omitempty"`
	Directory       string            `json:"directory,omitempty"`
	IsSyncedToGm    bool              `json:"is_synced_to_gm,omitempty"`
	LastModify      *time.Time        `json:"last_modify,omitempty"`
	Name            string            `json:"name,omitempty"`
	Type            string            `json:"type,omitempty"`
	VtftpDirMembers []*Vtftpdirmember `json:"vtftp_dir_members,omitempty"`
}

func (Tftpfiledir) ObjectType() string {
	return "tftpfiledir"
}

func (obj Tftpfiledir) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"directory", "name", "type"}
	}
	return obj.returnFields
}

// ThreatanalyticsAnalyticsWhitelist represents Infoblox object threatanalytics:analytics_whitelist
type ThreatanalyticsAnalyticsWhitelist struct {
	IBBase  `json:"-"`
	Ref     string `json:"_ref,omitempty"`
	Version string `json:"version,omitempty"`
}

func (ThreatanalyticsAnalyticsWhitelist) ObjectType() string {
	return "threatanalytics:analytics_whitelist"
}

func (obj ThreatanalyticsAnalyticsWhitelist) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"version"}
	}
	return obj.returnFields
}

// SyslogEndpoint represents Infoblox object syslog:endpoint
type SyslogEndpoint struct {
	IBBase               `json:"-"`
	Ref                  string                            `json:"_ref,omitempty"`
	Ea                   EA                                `json:"extattrs,omitempty"`
	LogLevel             string                            `json:"log_level,omitempty"`
	Name                 string                            `json:"name,omitempty"`
	OutboundMemberType   string                            `json:"outbound_member_type,omitempty"`
	OutboundMembers      []string                          `json:"outbound_members,omitempty"`
	SyslogServers        []*SyslogEndpointServers          `json:"syslog_servers,omitempty"`
	TemplateInstance     *NotificationRestTemplateinstance `json:"template_instance,omitempty"`
	TestSyslogConnection *Testsyslog                       `json:"test_syslog_connection,omitempty"`
	Timeout              uint32                            `json:"timeout,omitempty"`
	VendorIdentifier     string                            `json:"vendor_identifier,omitempty"`
	WapiUserName         string                            `json:"wapi_user_name,omitempty"`
	WapiUserPassword     string                            `json:"wapi_user_password,omitempty"`
}

func (SyslogEndpoint) ObjectType() string {
	return "syslog:endpoint"
}

func (obj SyslogEndpoint) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name", "outbound_member_type"}
	}
	return obj.returnFields
}

// ThreatanalyticsModuleset represents Infoblox object threatanalytics:moduleset
type ThreatanalyticsModuleset struct {
	IBBase  `json:"-"`
	Ref     string `json:"_ref,omitempty"`
	Version string `json:"version,omitempty"`
}

func (ThreatanalyticsModuleset) ObjectType() string {
	return "threatanalytics:moduleset"
}

func (obj ThreatanalyticsModuleset) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"version"}
	}
	return obj.returnFields
}

// ThreatanalyticsWhitelist represents Infoblox object threatanalytics:whitelist
type ThreatanalyticsWhitelist struct {
	IBBase  `json:"-"`
	Ref     string `json:"_ref,omitempty"`
	Comment string `json:"comment,omitempty"`
	Disable bool   `json:"disable,omitempty"`
	Fqdn    string `json:"fqdn,omitempty"`
	Type    string `json:"type,omitempty"`
}

func (ThreatanalyticsWhitelist) ObjectType() string {
	return "threatanalytics:whitelist"
}

func (obj ThreatanalyticsWhitelist) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "disable", "fqdn"}
	}
	return obj.returnFields
}

// ThreatinsightCloudclient represents Infoblox object threatinsight:cloudclient
type ThreatinsightCloudclient struct {
	IBBase           `json:"-"`
	Ref              string    `json:"_ref,omitempty"`
	BlacklistRpzList []*ZoneRp `json:"blacklist_rpz_list,omitempty"`
	Enable           bool      `json:"enable,omitempty"`
	ForceRefresh     bool      `json:"force_refresh,omitempty"`
	Interval         uint32    `json:"interval,omitempty"`
}

func (ThreatinsightCloudclient) ObjectType() string {
	return "threatinsight:cloudclient"
}

func (obj ThreatinsightCloudclient) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"enable", "interval"}
	}
	return obj.returnFields
}

// ThreatprotectionProfile represents Infoblox object threatprotection:profile
type ThreatprotectionProfile struct {
	IBBase                          `json:"-"`
	Ref                             string   `json:"_ref,omitempty"`
	Comment                         string   `json:"comment,omitempty"`
	CurrentRuleset                  string   `json:"current_ruleset,omitempty"`
	DisableMultipleDnsTcpRequest    bool     `json:"disable_multiple_dns_tcp_request,omitempty"`
	EventsPerSecondPerRule          uint32   `json:"events_per_second_per_rule,omitempty"`
	Ea                              EA       `json:"extattrs,omitempty"`
	Members                         []string `json:"members,omitempty"`
	Name                            string   `json:"name,omitempty"`
	SourceMember                    string   `json:"source_member,omitempty"`
	SourceProfile                   string   `json:"source_profile,omitempty"`
	UseCurrentRuleset               bool     `json:"use_current_ruleset,omitempty"`
	UseDisableMultipleDnsTcpRequest bool     `json:"use_disable_multiple_dns_tcp_request,omitempty"`
	UseEventsPerSecondPerRule       bool     `json:"use_events_per_second_per_rule,omitempty"`
}

func (ThreatprotectionProfile) ObjectType() string {
	return "threatprotection:profile"
}

func (obj ThreatprotectionProfile) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name"}
	}
	return obj.returnFields
}

// ThreatprotectionGridRule represents Infoblox object threatprotection:grid:rule
type ThreatprotectionGridRule struct {
	IBBase                `json:"-"`
	Ref                   string                      `json:"_ref,omitempty"`
	AllowedActions        []string                    `json:"allowed_actions,omitempty"`
	Category              string                      `json:"category,omitempty"`
	Comment               string                      `json:"comment,omitempty"`
	Config                *ThreatprotectionRuleconfig `json:"config,omitempty"`
	Description           string                      `json:"description,omitempty"`
	Disabled              bool                        `json:"disabled,omitempty"`
	IsFactoryResetEnabled bool                        `json:"is_factory_reset_enabled,omitempty"`
	Name                  string                      `json:"name,omitempty"`
	Ruleset               string                      `json:"ruleset,omitempty"`
	Sid                   uint32                      `json:"sid,omitempty"`
	Template              string                      `json:"template,omitempty"`
	Type                  string                      `json:"type,omitempty"`
}

func (ThreatprotectionGridRule) ObjectType() string {
	return "threatprotection:grid:rule"
}

func (obj ThreatprotectionGridRule) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name", "ruleset", "sid"}
	}
	return obj.returnFields
}

// ThreatprotectionRule represents Infoblox object threatprotection:rule
type ThreatprotectionRule struct {
	IBBase     `json:"-"`
	Ref        string                      `json:"_ref,omitempty"`
	Config     *ThreatprotectionRuleconfig `json:"config,omitempty"`
	Disable    bool                        `json:"disable,omitempty"`
	Member     string                      `json:"member,omitempty"`
	Rule       string                      `json:"rule,omitempty"`
	Sid        uint32                      `json:"sid,omitempty"`
	UseConfig  bool                        `json:"use_config,omitempty"`
	UseDisable bool                        `json:"use_disable,omitempty"`
}

func (ThreatprotectionRule) ObjectType() string {
	return "threatprotection:rule"
}

func (obj ThreatprotectionRule) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"member", "rule"}
	}
	return obj.returnFields
}

// ThreatprotectionProfileRule represents Infoblox object threatprotection:profile:rule
type ThreatprotectionProfileRule struct {
	IBBase     `json:"-"`
	Ref        string                      `json:"_ref,omitempty"`
	Config     *ThreatprotectionRuleconfig `json:"config,omitempty"`
	Disable    bool                        `json:"disable,omitempty"`
	Profile    string                      `json:"profile,omitempty"`
	Rule       string                      `json:"rule,omitempty"`
	Sid        uint32                      `json:"sid,omitempty"`
	UseConfig  bool                        `json:"use_config,omitempty"`
	UseDisable bool                        `json:"use_disable,omitempty"`
}

func (ThreatprotectionProfileRule) ObjectType() string {
	return "threatprotection:profile:rule"
}

func (obj ThreatprotectionProfileRule) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"profile", "rule"}
	}
	return obj.returnFields
}

// ThreatprotectionRulecategory represents Infoblox object threatprotection:rulecategory
type ThreatprotectionRulecategory struct {
	IBBase                `json:"-"`
	Ref                   string `json:"_ref,omitempty"`
	IsFactoryResetEnabled bool   `json:"is_factory_reset_enabled,omitempty"`
	Name                  string `json:"name,omitempty"`
	Ruleset               string `json:"ruleset,omitempty"`
}

func (ThreatprotectionRulecategory) ObjectType() string {
	return "threatprotection:rulecategory"
}

func (obj ThreatprotectionRulecategory) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name", "ruleset"}
	}
	return obj.returnFields
}

// ThreatprotectionRuleset represents Infoblox object threatprotection:ruleset
type ThreatprotectionRuleset struct {
	IBBase                `json:"-"`
	Ref                   string     `json:"_ref,omitempty"`
	AddType               string     `json:"add_type,omitempty"`
	AddedTime             *time.Time `json:"added_time,omitempty"`
	Comment               string     `json:"comment,omitempty"`
	DoNotDelete           bool       `json:"do_not_delete,omitempty"`
	IsFactoryResetEnabled bool       `json:"is_factory_reset_enabled,omitempty"`
	UsedBy                []string   `json:"used_by,omitempty"`
	Version               string     `json:"version,omitempty"`
}

func (ThreatprotectionRuleset) ObjectType() string {
	return "threatprotection:ruleset"
}

func (obj ThreatprotectionRuleset) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"add_type", "version"}
	}
	return obj.returnFields
}

// ThreatprotectionStatistics represents Infoblox object threatprotection:statistics
type ThreatprotectionStatistics struct {
	IBBase    `json:"-"`
	Ref       string                      `json:"_ref,omitempty"`
	Member    string                      `json:"member,omitempty"`
	StatInfos []*ThreatprotectionStatinfo `json:"stat_infos,omitempty"`
}

func (ThreatprotectionStatistics) ObjectType() string {
	return "threatprotection:statistics"
}

func (obj ThreatprotectionStatistics) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"member", "stat_infos"}
	}
	return obj.returnFields
}

// ThreatprotectionRuletemplate represents Infoblox object threatprotection:ruletemplate
type ThreatprotectionRuletemplate struct {
	IBBase         `json:"-"`
	Ref            string                      `json:"_ref,omitempty"`
	AllowedActions []string                    `json:"allowed_actions,omitempty"`
	Category       string                      `json:"category,omitempty"`
	DefaultConfig  *ThreatprotectionRuleconfig `json:"default_config,omitempty"`
	Description    string                      `json:"description,omitempty"`
	Name           string                      `json:"name,omitempty"`
	Ruleset        string                      `json:"ruleset,omitempty"`
	Sid            uint32                      `json:"sid,omitempty"`
}

func (ThreatprotectionRuletemplate) ObjectType() string {
	return "threatprotection:ruletemplate"
}

func (obj ThreatprotectionRuletemplate) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name", "ruleset", "sid"}
	}
	return obj.returnFields
}

// Upgradegroup represents Infoblox object upgradegroup
type Upgradegroup struct {
	IBBase                     `json:"-"`
	Ref                        string                `json:"_ref,omitempty"`
	Comment                    string                `json:"comment,omitempty"`
	DistributionDependentGroup string                `json:"distribution_dependent_group,omitempty"`
	DistributionPolicy         string                `json:"distribution_policy,omitempty"`
	DistributionTime           *time.Time            `json:"distribution_time,omitempty"`
	Members                    []*UpgradegroupMember `json:"members,omitempty"`
	Name                       string                `json:"name,omitempty"`
	TimeZone                   string                `json:"time_zone,omitempty"`
	UpgradeDependentGroup      string                `json:"upgrade_dependent_group,omitempty"`
	UpgradePolicy              string                `json:"upgrade_policy,omitempty"`
	UpgradeTime                *time.Time            `json:"upgrade_time,omitempty"`
}

func (Upgradegroup) ObjectType() string {
	return "upgradegroup"
}

func (obj Upgradegroup) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "name"}
	}
	return obj.returnFields
}

// Upgradeschedule represents Infoblox object upgradeschedule
type Upgradeschedule struct {
	IBBase        `json:"-"`
	Ref           string                  `json:"_ref,omitempty"`
	Active        bool                    `json:"active,omitempty"`
	StartTime     *time.Time              `json:"start_time,omitempty"`
	TimeZone      string                  `json:"time_zone,omitempty"`
	UpgradeGroups []*UpgradegroupSchedule `json:"upgrade_groups,omitempty"`
}

func (Upgradeschedule) ObjectType() string {
	return "upgradeschedule"
}

func (obj Upgradeschedule) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"active", "start_time", "time_zone"}
	}
	return obj.returnFields
}

// UserProfile represents Infoblox object userprofile
type UserProfile struct {
	IBBase                 `json:"-"`
	Ref                    string     `json:"_ref,omitempty"`
	ActiveDashboardType    string     `json:"active_dashboard_type,omitempty"`
	AdminGroup             string     `json:"admin_group,omitempty"`
	DaysToExpire           int        `json:"days_to_expire,omitempty"`
	Email                  string     `json:"email,omitempty"`
	GlobalSearchOnEa       bool       `json:"global_search_on_ea,omitempty"`
	GlobalSearchOnNiData   bool       `json:"global_search_on_ni_data,omitempty"`
	GridAdminGroups        []string   `json:"grid_admin_groups,omitempty"`
	LastLogin              *time.Time `json:"last_login,omitempty"`
	LbTreeNodesAtGenLevel  uint32     `json:"lb_tree_nodes_at_gen_level,omitempty"`
	LbTreeNodesAtLastLevel uint32     `json:"lb_tree_nodes_at_last_level,omitempty"`
	MaxCountWidgets        uint32     `json:"max_count_widgets,omitempty"`
	Name                   string     `json:"name,omitempty"`
	OldPassword            string     `json:"old_password,omitempty"`
	Password               string     `json:"password,omitempty"`
	TableSize              uint32     `json:"table_size,omitempty"`
	TimeZone               string     `json:"time_zone,omitempty"`
	UseTimeZone            bool       `json:"use_time_zone,omitempty"`
	UserType               string     `json:"user_type,omitempty"`
}

func (UserProfile) ObjectType() string {
	return "userprofile"
}

func (obj UserProfile) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name"}
	}
	return obj.returnFields
}

func NewUserProfile(userprofile UserProfile) *UserProfile {
	res := userprofile
	res.returnFields = []string{"name"}
	return &res
}

// UpgradeStatus represents Infoblox object upgradestatus
type UpgradeStatus struct {
	IBBase                      `json:"-"`
	Ref                         string           `json:"_ref,omitempty"`
	AllowDistribution           bool             `json:"allow_distribution,omitempty"`
	AllowDistributionScheduling bool             `json:"allow_distribution_scheduling,omitempty"`
	AllowUpgrade                bool             `json:"allow_upgrade,omitempty"`
	AllowUpgradeCancel          bool             `json:"allow_upgrade_cancel,omitempty"`
	AllowUpgradePause           bool             `json:"allow_upgrade_pause,omitempty"`
	AllowUpgradeResume          bool             `json:"allow_upgrade_resume,omitempty"`
	AllowUpgradeScheduling      bool             `json:"allow_upgrade_scheduling,omitempty"`
	AllowUpgradeTest            bool             `json:"allow_upgrade_test,omitempty"`
	AllowUpload                 bool             `json:"allow_upload,omitempty"`
	AlternateVersion            string           `json:"alternate_version,omitempty"`
	Comment                     string           `json:"comment,omitempty"`
	CurrentVersion              string           `json:"current_version,omitempty"`
	CurrentVersionSummary       string           `json:"current_version_summary,omitempty"`
	DistributionScheduleActive  bool             `json:"distribution_schedule_active,omitempty"`
	DistributionScheduleTime    *time.Time       `json:"distribution_schedule_time,omitempty"`
	DistributionState           string           `json:"distribution_state,omitempty"`
	DistributionVersion         string           `json:"distribution_version,omitempty"`
	DistributionVersionSummary  string           `json:"distribution_version_summary,omitempty"`
	ElementStatus               string           `json:"element_status,omitempty"`
	GridState                   string           `json:"grid_state,omitempty"`
	GroupState                  string           `json:"group_state,omitempty"`
	HaStatus                    string           `json:"ha_status,omitempty"`
	Hotfixes                    []*Hotfix        `json:"hotfixes,omitempty"`
	Ipv4Address                 string           `json:"ipv4_address,omitempty"`
	Ipv6Address                 string           `json:"ipv6_address,omitempty"`
	Member                      string           `json:"member,omitempty"`
	Message                     string           `json:"message,omitempty"`
	PnodeRole                   string           `json:"pnode_role,omitempty"`
	Reverted                    bool             `json:"reverted,omitempty"`
	StatusTime                  *time.Time       `json:"status_time,omitempty"`
	StatusValue                 string           `json:"status_value,omitempty"`
	StatusValueUpdateTime       *time.Time       `json:"status_value_update_time,omitempty"`
	Steps                       []*Upgradestep   `json:"steps,omitempty"`
	StepsCompleted              int              `json:"steps_completed,omitempty"`
	StepsTotal                  int              `json:"steps_total,omitempty"`
	SubelementType              string           `json:"subelement_type,omitempty"`
	SubelementsCompleted        int              `json:"subelements_completed,omitempty"`
	SubelementsStatus           []*UpgradeStatus `json:"subelements_status,omitempty"`
	SubelementsTotal            int              `json:"subelements_total,omitempty"`
	Type                        string           `json:"type,omitempty"`
	UpgradeGroup                string           `json:"upgrade_group,omitempty"`
	UpgradeScheduleActive       bool             `json:"upgrade_schedule_active,omitempty"`
	UpgradeState                string           `json:"upgrade_state,omitempty"`
	UpgradeTestStatus           string           `json:"upgrade_test_status,omitempty"`
	UploadVersion               string           `json:"upload_version,omitempty"`
	UploadVersionSummary        string           `json:"upload_version_summary,omitempty"`
}

func (UpgradeStatus) ObjectType() string {
	return "upgradestatus"
}

func (obj UpgradeStatus) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"alternate_version", "comment", "current_version", "distribution_version", "element_status", "grid_state", "group_state", "ha_status", "hotfixes", "ipv4_address", "ipv6_address", "member", "message", "pnode_role", "reverted", "status_value", "status_value_update_time", "steps", "steps_completed", "steps_total", "type", "upgrade_group", "upgrade_state", "upgrade_test_status", "upload_version"}
	}
	return obj.returnFields
}

func NewUpgradeStatus(upgradeStatus UpgradeStatus) *UpgradeStatus {
	result := upgradeStatus
	returnFields := []string{"subelements_status", "type"}
	result.returnFields = returnFields
	return &result
}

// Vlanrange represents Infoblox object vlanrange
type Vlanrange struct {
	IBBase              `json:"-"`
	Ref                 string               `json:"_ref,omitempty"`
	Comment             string               `json:"comment,omitempty"`
	DeleteVlans         bool                 `json:"delete_vlans,omitempty"`
	EndVlanId           uint32               `json:"end_vlan_id,omitempty"`
	Ea                  EA                   `json:"extattrs,omitempty"`
	Name                string               `json:"name,omitempty"`
	NextAvailableVlanId *Nextavailablevlanid `json:"next_available_vlan_id,omitempty"`
	PreCreateVlan       bool                 `json:"pre_create_vlan,omitempty"`
	StartVlanId         uint32               `json:"start_vlan_id,omitempty"`
	VlanNamePrefix      string               `json:"vlan_name_prefix,omitempty"`
	VlanView            string               `json:"vlan_view,omitempty"`
}

func (Vlanrange) ObjectType() string {
	return "vlanrange"
}

func (obj Vlanrange) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"end_vlan_id", "name", "start_vlan_id", "vlan_view"}
	}
	return obj.returnFields
}

// Vlan represents Infoblox object vlan
type Vlan struct {
	IBBase      `json:"-"`
	Ref         string         `json:"_ref,omitempty"`
	AssignedTo  []*Ipv4Network `json:"assigned_to,omitempty"`
	Comment     string         `json:"comment,omitempty"`
	Contact     string         `json:"contact,omitempty"`
	Department  string         `json:"department,omitempty"`
	Description string         `json:"description,omitempty"`
	Ea          EA             `json:"extattrs,omitempty"`
	Id          uint32         `json:"id,omitempty"`
	Name        string         `json:"name,omitempty"`
	Parent      string         `json:"parent,omitempty"`
	Reserved    bool           `json:"reserved,omitempty"`
	Status      string         `json:"status,omitempty"`
}

func (Vlan) ObjectType() string {
	return "vlan"
}

func (obj Vlan) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"id", "name", "parent"}
	}
	return obj.returnFields
}

// Vlanview represents Infoblox object vlanview
type Vlanview struct {
	IBBase                `json:"-"`
	Ref                   string               `json:"_ref,omitempty"`
	AllowRangeOverlapping bool                 `json:"allow_range_overlapping,omitempty"`
	Comment               string               `json:"comment,omitempty"`
	EndVlanId             uint32               `json:"end_vlan_id,omitempty"`
	Ea                    EA                   `json:"extattrs,omitempty"`
	Name                  string               `json:"name,omitempty"`
	NextAvailableVlanId   *Nextavailablevlanid `json:"next_available_vlan_id,omitempty"`
	PreCreateVlan         bool                 `json:"pre_create_vlan,omitempty"`
	StartVlanId           uint32               `json:"start_vlan_id,omitempty"`
	VlanNamePrefix        string               `json:"vlan_name_prefix,omitempty"`
}

func (Vlanview) ObjectType() string {
	return "vlanview"
}

func (obj Vlanview) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"end_vlan_id", "name", "start_vlan_id"}
	}
	return obj.returnFields
}

// ZoneAuthDiscrepancy represents Infoblox object zone_auth_discrepancy
type ZoneAuthDiscrepancy struct {
	IBBase      `json:"-"`
	Ref         string     `json:"_ref,omitempty"`
	Description string     `json:"description,omitempty"`
	Severity    string     `json:"severity,omitempty"`
	Timestamp   *time.Time `json:"timestamp,omitempty"`
	Zone        string     `json:"zone,omitempty"`
}

func (ZoneAuthDiscrepancy) ObjectType() string {
	return "zone_auth_discrepancy"
}

func (obj ZoneAuthDiscrepancy) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"description", "severity", "timestamp", "zone"}
	}
	return obj.returnFields
}

// Vdiscoverytask represents Infoblox object vdiscoverytask
type Vdiscoverytask struct {
	IBBase                          `json:"-"`
	Ref                             string             `json:"_ref,omitempty"`
	AllowUnsecuredConnection        bool               `json:"allow_unsecured_connection,omitempty"`
	AutoConsolidateCloudEa          bool               `json:"auto_consolidate_cloud_ea,omitempty"`
	AutoConsolidateManagedTenant    bool               `json:"auto_consolidate_managed_tenant,omitempty"`
	AutoConsolidateManagedVm        bool               `json:"auto_consolidate_managed_vm,omitempty"`
	AutoCreateDnsHostnameTemplate   string             `json:"auto_create_dns_hostname_template,omitempty"`
	AutoCreateDnsRecord             bool               `json:"auto_create_dns_record,omitempty"`
	AutoCreateDnsRecordType         string             `json:"auto_create_dns_record_type,omitempty"`
	Comment                         string             `json:"comment,omitempty"`
	CredentialsType                 string             `json:"credentials_type,omitempty"`
	DnsViewPrivateIp                string             `json:"dns_view_private_ip,omitempty"`
	DnsViewPublicIp                 string             `json:"dns_view_public_ip,omitempty"`
	DomainName                      string             `json:"domain_name,omitempty"`
	DriverType                      string             `json:"driver_type,omitempty"`
	Enabled                         bool               `json:"enabled,omitempty"`
	FqdnOrIp                        string             `json:"fqdn_or_ip,omitempty"`
	IdentityVersion                 string             `json:"identity_version,omitempty"`
	LastRun                         *time.Time         `json:"last_run,omitempty"`
	Member                          string             `json:"member,omitempty"`
	MergeData                       bool               `json:"merge_data,omitempty"`
	Name                            string             `json:"name,omitempty"`
	Password                        string             `json:"password,omitempty"`
	Port                            uint32             `json:"port,omitempty"`
	PrivateNetworkView              string             `json:"private_network_view,omitempty"`
	PrivateNetworkViewMappingPolicy string             `json:"private_network_view_mapping_policy,omitempty"`
	Protocol                        string             `json:"protocol,omitempty"`
	PublicNetworkView               string             `json:"public_network_view,omitempty"`
	PublicNetworkViewMappingPolicy  string             `json:"public_network_view_mapping_policy,omitempty"`
	ScheduledRun                    *SettingSchedule   `json:"scheduled_run,omitempty"`
	ServiceAccountFile              string             `json:"service_account_file,omitempty"`
	State                           string             `json:"state,omitempty"`
	StateMsg                        string             `json:"state_msg,omitempty"`
	UpdateDnsViewPrivateIp          bool               `json:"update_dns_view_private_ip,omitempty"`
	UpdateDnsViewPublicIp           bool               `json:"update_dns_view_public_ip,omitempty"`
	UpdateMetadata                  bool               `json:"update_metadata,omitempty"`
	UseIdentity                     bool               `json:"use_identity,omitempty"`
	Username                        string             `json:"username,omitempty"`
	VdiscoveryControl               *Vdiscoverycontrol `json:"vdiscovery_control,omitempty"`
}

func (Vdiscoverytask) ObjectType() string {
	return "vdiscoverytask"
}

func (obj Vdiscoverytask) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"name", "state"}
	}
	return obj.returnFields
}

// ZoneDelegated represents Infoblox object zone_delegated
type ZoneDelegated struct {
	IBBase                 `json:"-"`
	Ref                    string          `json:"_ref,omitempty"`
	Address                string          `json:"address,omitempty"`
	Comment                string          `json:"comment,omitempty"`
	DelegateTo             []NameServer    `json:"delegate_to,omitempty"`
	DelegatedTtl           uint32          `json:"delegated_ttl,omitempty"`
	Disable                bool            `json:"disable,omitempty"`
	DisplayDomain          string          `json:"display_domain,omitempty"`
	DnsFqdn                string          `json:"dns_fqdn,omitempty"`
	EnableRfc2317Exclusion bool            `json:"enable_rfc2317_exclusion,omitempty"`
	Ea                     EA              `json:"extattrs,omitempty"`
	Fqdn                   string          `json:"fqdn,omitempty"`
	LockUnlockZone         *Lockunlockzone `json:"lock_unlock_zone,omitempty"`
	Locked                 bool            `json:"locked,omitempty"`
	LockedBy               string          `json:"locked_by,omitempty"`
	MaskPrefix             string          `json:"mask_prefix,omitempty"`
	MsAdIntegrated         bool            `json:"ms_ad_integrated,omitempty"`
	MsDdnsMode             string          `json:"ms_ddns_mode,omitempty"`
	MsManaged              string          `json:"ms_managed,omitempty"`
	MsReadOnly             bool            `json:"ms_read_only,omitempty"`
	MsSyncMasterName       string          `json:"ms_sync_master_name,omitempty"`
	NsGroup                string          `json:"ns_group,omitempty"`
	Parent                 string          `json:"parent,omitempty"`
	Prefix                 string          `json:"prefix,omitempty"`
	UseDelegatedTtl        bool            `json:"use_delegated_ttl,omitempty"`
	UsingSrgAssociations   bool            `json:"using_srg_associations,omitempty"`
	View                   string          `json:"view,omitempty"`
	ZoneFormat             string          `json:"zone_format,omitempty"`
}

func (ZoneDelegated) ObjectType() string {
	return "zone_delegated"
}

func (obj ZoneDelegated) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"delegate_to", "fqdn", "view"}
	}
	return obj.returnFields
}

func NewZoneDelegated(za ZoneDelegated) *ZoneDelegated {
	res := za
	res.returnFields = []string{"extattrs", "fqdn", "view", "delegate_to"}

	return &res
}

// View represents Infoblox object view
type View struct {
	IBBase                              `json:"-"`
	Ref                                 string                        `json:"_ref,omitempty"`
	BlacklistAction                     string                        `json:"blacklist_action,omitempty"`
	BlacklistLogQuery                   bool                          `json:"blacklist_log_query,omitempty"`
	BlacklistRedirectAddresses          []string                      `json:"blacklist_redirect_addresses,omitempty"`
	BlacklistRedirectTtl                uint32                        `json:"blacklist_redirect_ttl,omitempty"`
	BlacklistRulesets                   []string                      `json:"blacklist_rulesets,omitempty"`
	CloudInfo                           *GridCloudapiInfo             `json:"cloud_info,omitempty"`
	Comment                             string                        `json:"comment,omitempty"`
	CustomRootNameServers               []NameServer                  `json:"custom_root_name_servers,omitempty"`
	DdnsForceCreationTimestampUpdate    bool                          `json:"ddns_force_creation_timestamp_update,omitempty"`
	DdnsPrincipalGroup                  string                        `json:"ddns_principal_group,omitempty"`
	DdnsPrincipalTracking               bool                          `json:"ddns_principal_tracking,omitempty"`
	DdnsRestrictPatterns                bool                          `json:"ddns_restrict_patterns,omitempty"`
	DdnsRestrictPatternsList            []string                      `json:"ddns_restrict_patterns_list,omitempty"`
	DdnsRestrictProtected               bool                          `json:"ddns_restrict_protected,omitempty"`
	DdnsRestrictSecure                  bool                          `json:"ddns_restrict_secure,omitempty"`
	DdnsRestrictStatic                  bool                          `json:"ddns_restrict_static,omitempty"`
	Disable                             bool                          `json:"disable,omitempty"`
	Dns64Enabled                        bool                          `json:"dns64_enabled,omitempty"`
	Dns64Groups                         []string                      `json:"dns64_groups,omitempty"`
	DnssecEnabled                       bool                          `json:"dnssec_enabled,omitempty"`
	DnssecExpiredSignaturesEnabled      bool                          `json:"dnssec_expired_signatures_enabled,omitempty"`
	DnssecNegativeTrustAnchors          []string                      `json:"dnssec_negative_trust_anchors,omitempty"`
	DnssecTrustedKeys                   []*Dnssectrustedkey           `json:"dnssec_trusted_keys,omitempty"`
	DnssecValidationEnabled             bool                          `json:"dnssec_validation_enabled,omitempty"`
	EdnsUdpSize                         uint32                        `json:"edns_udp_size,omitempty"`
	EnableBlacklist                     bool                          `json:"enable_blacklist,omitempty"`
	EnableFixedRrsetOrderFqdns          bool                          `json:"enable_fixed_rrset_order_fqdns,omitempty"`
	EnableMatchRecursiveOnly            bool                          `json:"enable_match_recursive_only,omitempty"`
	Ea                                  EA                            `json:"extattrs,omitempty"`
	FilterAaaa                          string                        `json:"filter_aaaa,omitempty"`
	FilterAaaaList                      []*Addressac                  `json:"filter_aaaa_list,omitempty"`
	FixedRrsetOrderFqdns                []*GridDnsFixedrrsetorderfqdn `json:"fixed_rrset_order_fqdns,omitempty"`
	ForwardOnly                         bool                          `json:"forward_only,omitempty"`
	Forwarders                          []string                      `json:"forwarders,omitempty"`
	IsDefault                           bool                          `json:"is_default,omitempty"`
	LameTtl                             uint32                        `json:"lame_ttl,omitempty"`
	LastQueriedAcl                      []*Addressac                  `json:"last_queried_acl,omitempty"`
	MatchClients                        []*Addressac                  `json:"match_clients,omitempty"`
	MatchDestinations                   []*Addressac                  `json:"match_destinations,omitempty"`
	MaxCacheTtl                         uint32                        `json:"max_cache_ttl,omitempty"`
	MaxNcacheTtl                        uint32                        `json:"max_ncache_ttl,omitempty"`
	MaxUdpSize                          uint32                        `json:"max_udp_size,omitempty"`
	Name                                string                        `json:"name,omitempty"`
	NetworkView                         string                        `json:"network_view,omitempty"`
	NotifyDelay                         uint32                        `json:"notify_delay,omitempty"`
	NxdomainLogQuery                    bool                          `json:"nxdomain_log_query,omitempty"`
	NxdomainRedirect                    bool                          `json:"nxdomain_redirect,omitempty"`
	NxdomainRedirectAddresses           []string                      `json:"nxdomain_redirect_addresses,omitempty"`
	NxdomainRedirectAddressesV6         []string                      `json:"nxdomain_redirect_addresses_v6,omitempty"`
	NxdomainRedirectTtl                 uint32                        `json:"nxdomain_redirect_ttl,omitempty"`
	NxdomainRulesets                    []string                      `json:"nxdomain_rulesets,omitempty"`
	Recursion                           bool                          `json:"recursion,omitempty"`
	ResponseRateLimiting                *GridResponseratelimiting     `json:"response_rate_limiting,omitempty"`
	RootNameServerType                  string                        `json:"root_name_server_type,omitempty"`
	RpzDropIpRuleEnabled                bool                          `json:"rpz_drop_ip_rule_enabled,omitempty"`
	RpzDropIpRuleMinPrefixLengthIpv4    uint32                        `json:"rpz_drop_ip_rule_min_prefix_length_ipv4,omitempty"`
	RpzDropIpRuleMinPrefixLengthIpv6    uint32                        `json:"rpz_drop_ip_rule_min_prefix_length_ipv6,omitempty"`
	RpzQnameWaitRecurse                 bool                          `json:"rpz_qname_wait_recurse,omitempty"`
	RunScavenging                       *Runscavenging                `json:"run_scavenging,omitempty"`
	ScavengingSettings                  *SettingScavenging            `json:"scavenging_settings,omitempty"`
	Sortlist                            []*Sortlist                   `json:"sortlist,omitempty"`
	UseBlacklist                        bool                          `json:"use_blacklist,omitempty"`
	UseDdnsForceCreationTimestampUpdate bool                          `json:"use_ddns_force_creation_timestamp_update,omitempty"`
	UseDdnsPatternsRestriction          bool                          `json:"use_ddns_patterns_restriction,omitempty"`
	UseDdnsPrincipalSecurity            bool                          `json:"use_ddns_principal_security,omitempty"`
	UseDdnsRestrictProtected            bool                          `json:"use_ddns_restrict_protected,omitempty"`
	UseDdnsRestrictStatic               bool                          `json:"use_ddns_restrict_static,omitempty"`
	UseDns64                            bool                          `json:"use_dns64,omitempty"`
	UseDnssec                           bool                          `json:"use_dnssec,omitempty"`
	UseEdnsUdpSize                      bool                          `json:"use_edns_udp_size,omitempty"`
	UseFilterAaaa                       bool                          `json:"use_filter_aaaa,omitempty"`
	UseFixedRrsetOrderFqdns             bool                          `json:"use_fixed_rrset_order_fqdns,omitempty"`
	UseForwarders                       bool                          `json:"use_forwarders,omitempty"`
	UseLameTtl                          bool                          `json:"use_lame_ttl,omitempty"`
	UseMaxCacheTtl                      bool                          `json:"use_max_cache_ttl,omitempty"`
	UseMaxNcacheTtl                     bool                          `json:"use_max_ncache_ttl,omitempty"`
	UseMaxUdpSize                       bool                          `json:"use_max_udp_size,omitempty"`
	UseNxdomainRedirect                 bool                          `json:"use_nxdomain_redirect,omitempty"`
	UseRecursion                        bool                          `json:"use_recursion,omitempty"`
	UseResponseRateLimiting             bool                          `json:"use_response_rate_limiting,omitempty"`
	UseRootNameServer                   bool                          `json:"use_root_name_server,omitempty"`
	UseRpzDropIpRule                    bool                          `json:"use_rpz_drop_ip_rule,omitempty"`
	UseRpzQnameWaitRecurse              bool                          `json:"use_rpz_qname_wait_recurse,omitempty"`
	UseScavengingSettings               bool                          `json:"use_scavenging_settings,omitempty"`
	UseSortlist                         bool                          `json:"use_sortlist,omitempty"`
}

func (View) ObjectType() string {
	return "view"
}

func (obj View) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"comment", "is_default", "name"}
	}
	return obj.returnFields
}

// ZoneStub represents Infoblox object zone_stub
type ZoneStub struct {
	IBBase               `json:"-"`
	Ref                  string          `json:"_ref,omitempty"`
	Address              string          `json:"address,omitempty"`
	Comment              string          `json:"comment,omitempty"`
	Disable              bool            `json:"disable,omitempty"`
	DisableForwarding    bool            `json:"disable_forwarding,omitempty"`
	DisplayDomain        string          `json:"display_domain,omitempty"`
	DnsFqdn              string          `json:"dns_fqdn,omitempty"`
	Ea                   EA              `json:"extattrs,omitempty"`
	ExternalNsGroup      string          `json:"external_ns_group,omitempty"`
	Fqdn                 string          `json:"fqdn,omitempty"`
	LockUnlockZone       *Lockunlockzone `json:"lock_unlock_zone,omitempty"`
	Locked               bool            `json:"locked,omitempty"`
	LockedBy             string          `json:"locked_by,omitempty"`
	MaskPrefix           string          `json:"mask_prefix,omitempty"`
	MsAdIntegrated       bool            `json:"ms_ad_integrated,omitempty"`
	MsDdnsMode           string          `json:"ms_ddns_mode,omitempty"`
	MsManaged            string          `json:"ms_managed,omitempty"`
	MsReadOnly           bool            `json:"ms_read_only,omitempty"`
	MsSyncMasterName     string          `json:"ms_sync_master_name,omitempty"`
	NsGroup              string          `json:"ns_group,omitempty"`
	Parent               string          `json:"parent,omitempty"`
	Prefix               string          `json:"prefix,omitempty"`
	SoaEmail             string          `json:"soa_email,omitempty"`
	SoaExpire            uint32          `json:"soa_expire,omitempty"`
	SoaMname             string          `json:"soa_mname,omitempty"`
	SoaNegativeTtl       uint32          `json:"soa_negative_ttl,omitempty"`
	SoaRefresh           uint32          `json:"soa_refresh,omitempty"`
	SoaRetry             uint32          `json:"soa_retry,omitempty"`
	SoaSerialNumber      uint32          `json:"soa_serial_number,omitempty"`
	StubFrom             []NameServer    `json:"stub_from,omitempty"`
	StubMembers          []*Memberserver `json:"stub_members,omitempty"`
	StubMsservers        []*Msdnsserver  `json:"stub_msservers,omitempty"`
	UsingSrgAssociations bool            `json:"using_srg_associations,omitempty"`
	View                 string          `json:"view,omitempty"`
	ZoneFormat           string          `json:"zone_format,omitempty"`
}

func (ZoneStub) ObjectType() string {
	return "zone_stub"
}

func (obj ZoneStub) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"fqdn", "stub_from", "view"}
	}
	return obj.returnFields
}

// ZoneForward represents Infoblox object zone_forward
type ZoneForward struct {
	IBBase               `json:"-"`
	Ref                  string                    `json:"_ref,omitempty"`
	Address              string                    `json:"address,omitempty"`
	Comment              string                    `json:"comment,omitempty"`
	Disable              bool                      `json:"disable,omitempty"`
	DisableNsGeneration  bool                      `json:"disable_ns_generation,omitempty"`
	DisplayDomain        string                    `json:"display_domain,omitempty"`
	DnsFqdn              string                    `json:"dns_fqdn,omitempty"`
	Ea                   EA                        `json:"extattrs,omitempty"`
	ExternalNsGroup      string                    `json:"external_ns_group,omitempty"`
	ForwardTo            []NameServer              `json:"forward_to,omitempty"`
	ForwardersOnly       bool                      `json:"forwarders_only,omitempty"`
	ForwardingServers    []*Forwardingmemberserver `json:"forwarding_servers,omitempty"`
	Fqdn                 string                    `json:"fqdn,omitempty"`
	LockUnlockZone       *Lockunlockzone           `json:"lock_unlock_zone,omitempty"`
	Locked               bool                      `json:"locked,omitempty"`
	LockedBy             string                    `json:"locked_by,omitempty"`
	MaskPrefix           string                    `json:"mask_prefix,omitempty"`
	MsAdIntegrated       bool                      `json:"ms_ad_integrated,omitempty"`
	MsDdnsMode           string                    `json:"ms_ddns_mode,omitempty"`
	MsManaged            string                    `json:"ms_managed,omitempty"`
	MsReadOnly           bool                      `json:"ms_read_only,omitempty"`
	MsSyncMasterName     string                    `json:"ms_sync_master_name,omitempty"`
	NsGroup              string                    `json:"ns_group,omitempty"`
	Parent               string                    `json:"parent,omitempty"`
	Prefix               string                    `json:"prefix,omitempty"`
	UsingSrgAssociations bool                      `json:"using_srg_associations,omitempty"`
	View                 string                    `json:"view,omitempty"`
	ZoneFormat           string                    `json:"zone_format,omitempty"`
}

func (ZoneForward) ObjectType() string {
	return "zone_forward"
}

func (obj ZoneForward) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"forward_to", "fqdn", "view"}
	}
	return obj.returnFields
}

// ZoneAuth represents Infoblox object zone_auth
type ZoneAuth struct {
	IBBase                                  `json:"-"`
	Ref                                     string                        `json:"_ref,omitempty"`
	Address                                 string                        `json:"address,omitempty"`
	AllowActiveDir                          []*Addressac                  `json:"allow_active_dir,omitempty"`
	AllowFixedRrsetOrder                    bool                          `json:"allow_fixed_rrset_order,omitempty"`
	AllowGssTsigForUnderscoreZone           bool                          `json:"allow_gss_tsig_for_underscore_zone,omitempty"`
	AllowGssTsigZoneUpdates                 bool                          `json:"allow_gss_tsig_zone_updates,omitempty"`
	AllowQuery                              []*Addressac                  `json:"allow_query,omitempty"`
	AllowTransfer                           []*Addressac                  `json:"allow_transfer,omitempty"`
	AllowUpdate                             []*Addressac                  `json:"allow_update,omitempty"`
	AllowUpdateForwarding                   bool                          `json:"allow_update_forwarding,omitempty"`
	AwsRte53ZoneInfo                        *Awsrte53zoneinfo             `json:"aws_rte53_zone_info,omitempty"`
	CloudInfo                               *GridCloudapiInfo             `json:"cloud_info,omitempty"`
	Comment                                 string                        `json:"comment,omitempty"`
	CopyXferToNotify                        bool                          `json:"copy_xfer_to_notify,omitempty"`
	Copyzonerecords                         *Copyzonerecords              `json:"copyzonerecords,omitempty"`
	CreatePtrForBulkHosts                   bool                          `json:"create_ptr_for_bulk_hosts,omitempty"`
	CreatePtrForHosts                       bool                          `json:"create_ptr_for_hosts,omitempty"`
	CreateUnderscoreZones                   bool                          `json:"create_underscore_zones,omitempty"`
	DdnsForceCreationTimestampUpdate        bool                          `json:"ddns_force_creation_timestamp_update,omitempty"`
	DdnsPrincipalGroup                      string                        `json:"ddns_principal_group,omitempty"`
	DdnsPrincipalTracking                   bool                          `json:"ddns_principal_tracking,omitempty"`
	DdnsRestrictPatterns                    bool                          `json:"ddns_restrict_patterns,omitempty"`
	DdnsRestrictPatternsList                []string                      `json:"ddns_restrict_patterns_list,omitempty"`
	DdnsRestrictProtected                   bool                          `json:"ddns_restrict_protected,omitempty"`
	DdnsRestrictSecure                      bool                          `json:"ddns_restrict_secure,omitempty"`
	DdnsRestrictStatic                      bool                          `json:"ddns_restrict_static,omitempty"`
	Disable                                 bool                          `json:"disable,omitempty"`
	DisableForwarding                       bool                          `json:"disable_forwarding,omitempty"`
	DisplayDomain                           string                        `json:"display_domain,omitempty"`
	DnsFqdn                                 string                        `json:"dns_fqdn,omitempty"`
	DnsIntegrityEnable                      bool                          `json:"dns_integrity_enable,omitempty"`
	DnsIntegrityFrequency                   uint32                        `json:"dns_integrity_frequency,omitempty"`
	DnsIntegrityMember                      string                        `json:"dns_integrity_member,omitempty"`
	DnsIntegrityVerboseLogging              bool                          `json:"dns_integrity_verbose_logging,omitempty"`
	DnsSoaEmail                             string                        `json:"dns_soa_email,omitempty"`
	DnssecExport                            *Dnssecexport                 `json:"dnssec_export,omitempty"`
	DnssecGetZoneKeys                       *Dnssecgetzonekeys            `json:"dnssec_get_zone_keys,omitempty"`
	DnssecKeyParams                         *Dnsseckeyparams              `json:"dnssec_key_params,omitempty"`
	DnssecKeys                              []*Dnsseckey                  `json:"dnssec_keys,omitempty"`
	DnssecKskRolloverDate                   *time.Time                    `json:"dnssec_ksk_rollover_date,omitempty"`
	DnssecOperation                         *Dnssecoperation              `json:"dnssec_operation,omitempty"`
	DnssecSetZoneKeys                       *Dnssecsetzonekeys            `json:"dnssec_set_zone_keys,omitempty"`
	DnssecZskRolloverDate                   *time.Time                    `json:"dnssec_zsk_rollover_date,omitempty"`
	Dnssecgetkskrollover                    *Dnssecgetkskrollover         `json:"dnssecgetkskrollover,omitempty"`
	DoHostAbstraction                       bool                          `json:"do_host_abstraction,omitempty"`
	EffectiveCheckNamesPolicy               string                        `json:"effective_check_names_policy,omitempty"`
	EffectiveRecordNamePolicy               string                        `json:"effective_record_name_policy,omitempty"`
	ExecuteDnsParentCheck                   *Parentcheck                  `json:"execute_dns_parent_check,omitempty"`
	Ea                                      EA                            `json:"extattrs,omitempty"`
	ExternalPrimaries                       []NameServer                  `json:"external_primaries,omitempty"`
	ExternalSecondaries                     []NameServer                  `json:"external_secondaries,omitempty"`
	Fqdn                                    string                        `json:"fqdn,omitempty"`
	GridPrimary                             []*Memberserver               `json:"grid_primary,omitempty"`
	GridPrimarySharedWithMsParentDelegation bool                          `json:"grid_primary_shared_with_ms_parent_delegation,omitempty"`
	GridSecondaries                         []*Memberserver               `json:"grid_secondaries,omitempty"`
	ImportFrom                              string                        `json:"import_from,omitempty"`
	IsDnssecEnabled                         bool                          `json:"is_dnssec_enabled,omitempty"`
	IsDnssecSigned                          bool                          `json:"is_dnssec_signed,omitempty"`
	IsMultimaster                           bool                          `json:"is_multimaster,omitempty"`
	LastQueried                             *time.Time                    `json:"last_queried,omitempty"`
	LastQueriedAcl                          []*Addressac                  `json:"last_queried_acl,omitempty"`
	LockUnlockZone                          *Lockunlockzone               `json:"lock_unlock_zone,omitempty"`
	Locked                                  bool                          `json:"locked,omitempty"`
	LockedBy                                string                        `json:"locked_by,omitempty"`
	MaskPrefix                              string                        `json:"mask_prefix,omitempty"`
	MemberSoaMnames                         []*GridmemberSoamname         `json:"member_soa_mnames,omitempty"`
	MemberSoaSerials                        []*GridmemberSoaserial        `json:"member_soa_serials,omitempty"`
	MsAdIntegrated                          bool                          `json:"ms_ad_integrated,omitempty"`
	MsAllowTransfer                         []*Addressac                  `json:"ms_allow_transfer,omitempty"`
	MsAllowTransferMode                     string                        `json:"ms_allow_transfer_mode,omitempty"`
	MsDcNsRecordCreation                    []*MsserverDcnsrecordcreation `json:"ms_dc_ns_record_creation,omitempty"`
	MsDdnsMode                              string                        `json:"ms_ddns_mode,omitempty"`
	MsManaged                               string                        `json:"ms_managed,omitempty"`
	MsPrimaries                             []*Msdnsserver                `json:"ms_primaries,omitempty"`
	MsReadOnly                              bool                          `json:"ms_read_only,omitempty"`
	MsSecondaries                           []*Msdnsserver                `json:"ms_secondaries,omitempty"`
	MsSyncDisabled                          bool                          `json:"ms_sync_disabled,omitempty"`
	MsSyncMasterName                        string                        `json:"ms_sync_master_name,omitempty"`
	NetworkAssociations                     []*Ipv4Network                `json:"network_associations,omitempty"`
	NetworkView                             string                        `json:"network_view,omitempty"`
	NotifyDelay                             uint32                        `json:"notify_delay,omitempty"`
	NsGroup                                 string                        `json:"ns_group,omitempty"`
	Parent                                  string                        `json:"parent,omitempty"`
	Prefix                                  string                        `json:"prefix,omitempty"`
	PrimaryType                             string                        `json:"primary_type,omitempty"`
	RecordNamePolicy                        string                        `json:"record_name_policy,omitempty"`
	RecordsMonitored                        bool                          `json:"records_monitored,omitempty"`
	RestartIfNeeded                         bool                          `json:"restart_if_needed,omitempty"`
	RrNotQueriedEnabledTime                 *time.Time                    `json:"rr_not_queried_enabled_time,omitempty"`
	RunScavenging                           *Runscavenging                `json:"run_scavenging,omitempty"`
	ScavengingSettings                      *SettingScavenging            `json:"scavenging_settings,omitempty"`
	SetSoaSerialNumber                      bool                          `json:"set_soa_serial_number,omitempty"`
	SoaDefaultTtl                           uint32                        `json:"soa_default_ttl,omitempty"`
	SoaEmail                                string                        `json:"soa_email,omitempty"`
	SoaExpire                               uint32                        `json:"soa_expire,omitempty"`
	SoaNegativeTtl                          uint32                        `json:"soa_negative_ttl,omitempty"`
	SoaRefresh                              uint32                        `json:"soa_refresh,omitempty"`
	SoaRetry                                uint32                        `json:"soa_retry,omitempty"`
	SoaSerialNumber                         uint32                        `json:"soa_serial_number,omitempty"`
	Srgs                                    []string                      `json:"srgs,omitempty"`
	UpdateForwarding                        []*Addressac                  `json:"update_forwarding,omitempty"`
	UseAllowActiveDir                       bool                          `json:"use_allow_active_dir,omitempty"`
	UseAllowQuery                           bool                          `json:"use_allow_query,omitempty"`
	UseAllowTransfer                        bool                          `json:"use_allow_transfer,omitempty"`
	UseAllowUpdate                          bool                          `json:"use_allow_update,omitempty"`
	UseAllowUpdateForwarding                bool                          `json:"use_allow_update_forwarding,omitempty"`
	UseCheckNamesPolicy                     bool                          `json:"use_check_names_policy,omitempty"`
	UseCopyXferToNotify                     bool                          `json:"use_copy_xfer_to_notify,omitempty"`
	UseDdnsForceCreationTimestampUpdate     bool                          `json:"use_ddns_force_creation_timestamp_update,omitempty"`
	UseDdnsPatternsRestriction              bool                          `json:"use_ddns_patterns_restriction,omitempty"`
	UseDdnsPrincipalSecurity                bool                          `json:"use_ddns_principal_security,omitempty"`
	UseDdnsRestrictProtected                bool                          `json:"use_ddns_restrict_protected,omitempty"`
	UseDdnsRestrictStatic                   bool                          `json:"use_ddns_restrict_static,omitempty"`
	UseDnssecKeyParams                      bool                          `json:"use_dnssec_key_params,omitempty"`
	UseExternalPrimary                      bool                          `json:"use_external_primary,omitempty"`
	UseGridZoneTimer                        bool                          `json:"use_grid_zone_timer,omitempty"`
	UseImportFrom                           bool                          `json:"use_import_from,omitempty"`
	UseNotifyDelay                          bool                          `json:"use_notify_delay,omitempty"`
	UseRecordNamePolicy                     bool                          `json:"use_record_name_policy,omitempty"`
	UseScavengingSettings                   bool                          `json:"use_scavenging_settings,omitempty"`
	UseSoaEmail                             bool                          `json:"use_soa_email,omitempty"`
	UsingSrgAssociations                    bool                          `json:"using_srg_associations,omitempty"`
	View                                    string                        `json:"view,omitempty"`
	ZoneFormat                              string                        `json:"zone_format,omitempty"`
	ZoneNotQueriedEnabledTime               *time.Time                    `json:"zone_not_queried_enabled_time,omitempty"`
}

func (ZoneAuth) ObjectType() string {
	return "zone_auth"
}

func (obj ZoneAuth) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"fqdn", "view"}
	}
	return obj.returnFields
}

func NewZoneAuth(za ZoneAuth) *ZoneAuth {
	res := za
	res.returnFields = []string{"extattrs", "fqdn", "view"}
	return &res
}

// ZoneRp represents Infoblox object zone_rp
type ZoneRp struct {
	IBBase                           `json:"-"`
	Ref                              string                 `json:"_ref,omitempty"`
	Address                          string                 `json:"address,omitempty"`
	Comment                          string                 `json:"comment,omitempty"`
	CopyRpzRecords                   *Copyrpzrecords        `json:"copy_rpz_records,omitempty"`
	Disable                          bool                   `json:"disable,omitempty"`
	DisplayDomain                    string                 `json:"display_domain,omitempty"`
	DnsSoaEmail                      string                 `json:"dns_soa_email,omitempty"`
	Ea                               EA                     `json:"extattrs,omitempty"`
	ExternalPrimaries                []NameServer           `json:"external_primaries,omitempty"`
	ExternalSecondaries              []NameServer           `json:"external_secondaries,omitempty"`
	FireeyeRuleMapping               *FireeyeRulemapping    `json:"fireeye_rule_mapping,omitempty"`
	Fqdn                             string                 `json:"fqdn,omitempty"`
	GridPrimary                      []*Memberserver        `json:"grid_primary,omitempty"`
	GridSecondaries                  []*Memberserver        `json:"grid_secondaries,omitempty"`
	LockUnlockZone                   *Lockunlockzone        `json:"lock_unlock_zone,omitempty"`
	Locked                           bool                   `json:"locked,omitempty"`
	LockedBy                         string                 `json:"locked_by,omitempty"`
	LogRpz                           bool                   `json:"log_rpz,omitempty"`
	MaskPrefix                       string                 `json:"mask_prefix,omitempty"`
	MemberSoaMnames                  []*GridmemberSoamname  `json:"member_soa_mnames,omitempty"`
	MemberSoaSerials                 []*GridmemberSoaserial `json:"member_soa_serials,omitempty"`
	NetworkView                      string                 `json:"network_view,omitempty"`
	NsGroup                          string                 `json:"ns_group,omitempty"`
	Parent                           string                 `json:"parent,omitempty"`
	Prefix                           string                 `json:"prefix,omitempty"`
	PrimaryType                      string                 `json:"primary_type,omitempty"`
	RecordNamePolicy                 string                 `json:"record_name_policy,omitempty"`
	RpzDropIpRuleEnabled             bool                   `json:"rpz_drop_ip_rule_enabled,omitempty"`
	RpzDropIpRuleMinPrefixLengthIpv4 uint32                 `json:"rpz_drop_ip_rule_min_prefix_length_ipv4,omitempty"`
	RpzDropIpRuleMinPrefixLengthIpv6 uint32                 `json:"rpz_drop_ip_rule_min_prefix_length_ipv6,omitempty"`
	RpzLastUpdatedTime               *time.Time             `json:"rpz_last_updated_time,omitempty"`
	RpzPolicy                        string                 `json:"rpz_policy,omitempty"`
	RpzPriority                      uint32                 `json:"rpz_priority,omitempty"`
	RpzPriorityEnd                   uint32                 `json:"rpz_priority_end,omitempty"`
	RpzSeverity                      string                 `json:"rpz_severity,omitempty"`
	RpzType                          string                 `json:"rpz_type,omitempty"`
	SetSoaSerialNumber               bool                   `json:"set_soa_serial_number,omitempty"`
	SoaDefaultTtl                    uint32                 `json:"soa_default_ttl,omitempty"`
	SoaEmail                         string                 `json:"soa_email,omitempty"`
	SoaExpire                        uint32                 `json:"soa_expire,omitempty"`
	SoaNegativeTtl                   uint32                 `json:"soa_negative_ttl,omitempty"`
	SoaRefresh                       uint32                 `json:"soa_refresh,omitempty"`
	SoaRetry                         uint32                 `json:"soa_retry,omitempty"`
	SoaSerialNumber                  uint32                 `json:"soa_serial_number,omitempty"`
	SubstituteName                   string                 `json:"substitute_name,omitempty"`
	UseExternalPrimary               bool                   `json:"use_external_primary,omitempty"`
	UseGridZoneTimer                 bool                   `json:"use_grid_zone_timer,omitempty"`
	UseLogRpz                        bool                   `json:"use_log_rpz,omitempty"`
	UseRecordNamePolicy              bool                   `json:"use_record_name_policy,omitempty"`
	UseRpzDropIpRule                 bool                   `json:"use_rpz_drop_ip_rule,omitempty"`
	UseSoaEmail                      bool                   `json:"use_soa_email,omitempty"`
	View                             string                 `json:"view,omitempty"`
}

func (ZoneRp) ObjectType() string {
	return "zone_rp"
}

func (obj ZoneRp) ReturnFields() []string {
	if obj.returnFields == nil {
		obj.returnFields = []string{"fqdn", "view"}
	}
	return obj.returnFields
}

// AdAuthServer represents Infoblox struct ad_auth_server
type AdAuthServer struct {
	FqdnOrIp    string `json:"fqdn_or_ip,omitempty"`
	AuthPort    uint32 `json:"auth_port,omitempty"`
	Comment     string `json:"comment,omitempty"`
	Disabled    bool   `json:"disabled,omitempty"`
	Encryption  string `json:"encryption,omitempty"`
	MgmtPort    bool   `json:"mgmt_port,omitempty"`
	UseMgmtPort bool   `json:"use_mgmt_port,omitempty"`
}

// Addcertificate represents Infoblox struct addcertificate
type Addcertificate struct{}

// Addressac represents Infoblox struct addressac
type Addressac struct {
	Address        string `json:"address,omitempty"`
	Permission     string `json:"permission,omitempty"`
	TsigKey        string `json:"tsig_key,omitempty"`
	TsigKeyAlg     string `json:"tsig_key_alg,omitempty"`
	TsigKeyName    string `json:"tsig_key_name,omitempty"`
	UseTsigKeyName bool   `json:"use_tsig_key_name,omitempty"`
}

// AdmingroupAdminsetcommands represents Infoblox struct admingroup:adminsetcommands
type AdmingroupAdminsetcommands struct {
	SetAdminGroupAcl             bool `json:"set_admin_group_acl,omitempty"`
	EtBfd                        bool `json:"et_bfd,omitempty"`
	SetBfd                       bool `json:"set_bfd,omitempty"`
	SetBgp                       bool `json:"set_bgp,omitempty"`
	SetBloxtools                 bool `json:"set_bloxtools,omitempty"`
	SetCleanMscache              bool `json:"set_clean_mscache,omitempty"`
	SetDebug                     bool `json:"set_debug,omitempty"`
	SetDebugAnalytics            bool `json:"set_debug_analytics,omitempty"`
	SetDeleteTasksInterval       bool `json:"set_delete_tasks_interval,omitempty"`
	SetDisableGuiOneClickSupport bool `json:"set_disable_gui_one_click_support,omitempty"`
	SetHardwareType              bool `json:"set_hardware_type,omitempty"`
	SetIbtrap                    bool `json:"set_ibtrap,omitempty"`
	SetLcd                       bool `json:"set_lcd,omitempty"`
	SetLcdSettings               bool `json:"set_lcd_settings,omitempty"`
	SetLines                     bool `json:"set_lines,omitempty"`
	SetMsMaxConnection           bool `json:"set_ms_max_connection,omitempty"`
	SetNosafemode                bool `json:"set_nosafemode,omitempty"`
	SetOcsp                      bool `json:"set_ocsp,omitempty"`
	SetPurgeRestartObjects       bool `json:"set_purge_restart_objects,omitempty"`
	SetReportingUserCapabilities bool `json:"set_reporting_user_capabilities,omitempty"`
	SetRpzRecursiveOnly          bool `json:"set_rpz_recursive_only,omitempty"`
	SetSafemode                  bool `json:"set_safemode,omitempty"`
	SetScheduled                 bool `json:"set_scheduled,omitempty"`
	SetSnmptrap                  bool `json:"set_snmptrap,omitempty"`
	SetSysname                   bool `json:"set_sysname,omitempty"`
	SetTerm                      bool `json:"set_term,omitempty"`
	SetThresholdtrap             bool `json:"set_thresholdtrap,omitempty"`
	SetExpertmode                bool `json:"set_expertmode,omitempty"`
	SetMaintenancemode           bool `json:"set_maintenancemode,omitempty"`
	SetTransferReportingData     bool `json:"set_transfer_reporting_data,omitempty"`
	SetTransferSupportbundle     bool `json:"set_transfer_supportbundle,omitempty"`
	EnableAll                    bool `json:"enable_all,omitempty"`
	DisableAll                   bool `json:"disable_all,omitempty"`
}

// AdmingroupAdminshowcommands represents Infoblox struct admingroup:adminshowcommands
type AdmingroupAdminshowcommands struct {
	ShowAdminGroupAcl             bool `json:"show_admin_group_acl,omitempty"`
	ShowAnalyticsParameter        bool `json:"show_analytics_parameter,omitempty"`
	ShowArp                       bool `json:"show_arp,omitempty"`
	ShowBfd                       bool `json:"show_bfd,omitempty"`
	ShowBgp                       bool `json:"show_bgp,omitempty"`
	ShowBloxtools                 bool `json:"show_bloxtools,omitempty"`
	ShowCapacity                  bool `json:"show_capacity,omitempty"`
	ShowClusterdInfo              bool `json:"show_clusterd_info,omitempty"`
	ShowConfig                    bool `json:"show_config,omitempty"`
	ShowCpu                       bool `json:"show_cpu,omitempty"`
	ShowDate                      bool `json:"show_date,omitempty"`
	ShowDebug                     bool `json:"show_debug,omitempty"`
	ShowDebugAnalytics            bool `json:"show_debug_analytics,omitempty"`
	ShowDeleteTasksInterval       bool `json:"show_delete_tasks_interval,omitempty"`
	ShowDisk                      bool `json:"show_disk,omitempty"`
	ShowFile                      bool `json:"show_file,omitempty"`
	ShowHardwareType              bool `json:"show_hardware_type,omitempty"`
	ShowHardwareStatus            bool `json:"show_hardware_status,omitempty"`
	ShowHwid                      bool `json:"show_hwid,omitempty"`
	ShowIbtrap                    bool `json:"show_ibtrap,omitempty"`
	ShowLcd                       bool `json:"show_lcd,omitempty"`
	ShowLcdInfo                   bool `json:"show_lcd_info,omitempty"`
	ShowLcdSettings               bool `json:"show_lcd_settings,omitempty"`
	ShowLog                       bool `json:"show_log,omitempty"`
	ShowLogfiles                  bool `json:"show_logfiles,omitempty"`
	ShowMemory                    bool `json:"show_memory,omitempty"`
	ShowNtp                       bool `json:"show_ntp,omitempty"`
	ShowReportingUserCapabilities bool `json:"show_reporting_user_capabilities,omitempty"`
	ShowRpzRecursiveOnly          bool `json:"show_rpz_recursive_only,omitempty"`
	ShowScheduled                 bool `json:"show_scheduled,omitempty"`
	ShowSnmp                      bool `json:"show_snmp,omitempty"`
	ShowStatus                    bool `json:"show_status,omitempty"`
	ShowTechSupport               bool `json:"show_tech_support,omitempty"`
	ShowTemperature               bool `json:"show_temperature,omitempty"`
	ShowThresholdtrap             bool `json:"show_thresholdtrap,omitempty"`
	ShowUpgradeHistory            bool `json:"show_upgrade_history,omitempty"`
	ShowUptime                    bool `json:"show_uptime,omitempty"`
	ShowVersion                   bool `json:"show_version,omitempty"`
	EnableAll                     bool `json:"enable_all,omitempty"`
	DisableAll                    bool `json:"disable_all,omitempty"`
}

// AdmingroupAdmintoplevelcommands represents Infoblox struct admingroup:admintoplevelcommands
type AdmingroupAdmintoplevelcommands struct {
	Ps             bool `json:"ps,omitempty"`
	Iostat         bool `json:"iostat,omitempty"`
	Netstat        bool `json:"netstat,omitempty"`
	Vmstat         bool `json:"vmstat,omitempty"`
	Tcpdump        bool `json:"tcpdump,omitempty"`
	Rndc           bool `json:"rndc,omitempty"`
	Sar            bool `json:"sar,omitempty"`
	Resilver       bool `json:"resilver,omitempty"`
	RestartProduct bool `json:"restart_product,omitempty"`
	Scrape         bool `json:"scrape,omitempty"`
	SamlRestart    bool `json:"saml_restart,omitempty"`
	EnableAll      bool `json:"enable_all,omitempty"`
	DisableAll     bool `json:"disable_all,omitempty"`
}

// AdmingroupCloudsetcommands represents Infoblox struct admingroup:cloudsetcommands
type AdmingroupCloudsetcommands struct {
	SetCloudServicesPortalForceRefresh bool `json:"set_cloud_services_portal_force_refresh,omitempty"`
	EnableAll                          bool `json:"enable_all,omitempty"`
	DisableAll                         bool `json:"disable_all,omitempty"`
}

// AdmingroupDatabasesetcommands represents Infoblox struct admingroup:databasesetcommands
type AdmingroupDatabasesetcommands struct {
	SetNamedMaxJournalSize bool `json:"set_named_max_journal_size,omitempty"`
	SetTxnTrace            bool `json:"set_txn_trace,omitempty"`
	SetDatabaseTransfer    bool `json:"set_database_transfer,omitempty"`
	EnableAll              bool `json:"enable_all,omitempty"`
	DisableAll             bool `json:"disable_all,omitempty"`
}

// AdmingroupDatabaseshowcommands represents Infoblox struct admingroup:databaseshowcommands
type AdmingroupDatabaseshowcommands struct {
	ShowNamedMaxJournalSize    bool `json:"show_named_max_journal_size,omitempty"`
	ShowTxnTrace               bool `json:"show_txn_trace,omitempty"`
	ShowDatabaseTransferStatus bool `json:"show_database_transfer_status,omitempty"`
	EnableAll                  bool `json:"enable_all,omitempty"`
	DisableAll                 bool `json:"disable_all,omitempty"`
}

// AdmingroupDhcpsetcommands represents Infoblox struct admingroup:dhcpsetcommands
type AdmingroupDhcpsetcommands struct {
	SetDhcpdRecvSockBufSize bool `json:"set_dhcpd_recv_sock_buf_size,omitempty"`
	SetLogTxnId             bool `json:"set_log_txn_id,omitempty"`
	SetOverloadBootp        bool `json:"set_overload_bootp,omitempty"`
	EnableAll               bool `json:"enable_all,omitempty"`
	DisableAll              bool `json:"disable_all,omitempty"`
}

// AdmingroupDhcpshowcommands represents Infoblox struct admingroup:dhcpshowcommands
type AdmingroupDhcpshowcommands struct {
	ShowDhcpGssTsig          bool `json:"show_dhcp_gss_tsig,omitempty"`
	ShowDhcpv6GssTsig        bool `json:"show_dhcpv6_gss_tsig,omitempty"`
	ShowDhcpdRecvSockBufSize bool `json:"show_dhcpd_recv_sock_buf_size,omitempty"`
	ShowOverloadBootp        bool `json:"show_overload_bootp,omitempty"`
	ShowLogTxnId             bool `json:"show_log_txn_id,omitempty"`
	EnableAll                bool `json:"enable_all,omitempty"`
	DisableAll               bool `json:"disable_all,omitempty"`
}

// AdmingroupDnssetcommands represents Infoblox struct admingroup:dnssetcommands
type AdmingroupDnssetcommands struct {
	SetDns                          bool `json:"set_dns,omitempty"`
	SetDnsRrl                       bool `json:"set_dns_rrl,omitempty"`
	SetEnableDnstap                 bool `json:"set_enable_dnstap,omitempty"`
	SetEnableMatchRecursiveOnly     bool `json:"set_enable_match_recursive_only,omitempty"`
	SetExtraDnsNameValidations      bool `json:"set_extra_dns_name_validations,omitempty"`
	SetLogGuestLookups              bool `json:"set_log_guest_lookups,omitempty"`
	SetMaxRecursionDepth            bool `json:"set_max_recursion_depth,omitempty"`
	SetMaxRecursionQueries          bool `json:"set_max_recursion_queries,omitempty"`
	SetMonitor                      bool `json:"set_monitor,omitempty"`
	SetMsDnsReportsSyncInterval     bool `json:"set_ms_dns_reports_sync_interval,omitempty"`
	SetMsStickyIp                   bool `json:"set_ms_sticky_ip,omitempty"`
	SetRestartAnycastWithDnsRestart bool `json:"set_restart_anycast_with_dns_restart,omitempty"`
	EnableAll                       bool `json:"enable_all,omitempty"`
	DisableAll                      bool `json:"disable_all,omitempty"`
}

// AdmingroupDnsshowcommands represents Infoblox struct admingroup:dnsshowcommands
type AdmingroupDnsshowcommands struct {
	ShowLogGuestLookups              bool `json:"show_log_guest_lookups,omitempty"`
	ShowDnsGssTsig                   bool `json:"show_dns_gss_tsig,omitempty"`
	ShowDns                          bool `json:"show_dns,omitempty"`
	ShowDnstapStats                  bool `json:"show_dnstap_stats,omitempty"`
	ShowDnstapStatus                 bool `json:"show_dnstap_status,omitempty"`
	ShowExtraDnsNameValidations      bool `json:"show_extra_dns_name_validations,omitempty"`
	ShowMsStickyIp                   bool `json:"show_ms_sticky_ip,omitempty"`
	ShowDnsRrl                       bool `json:"show_dns_rrl,omitempty"`
	ShowEnableMatchRecursiveOnly     bool `json:"show_enable_match_recursive_only,omitempty"`
	ShowMaxRecursionDepth            bool `json:"show_max_recursion_depth,omitempty"`
	ShowMaxRecursionQueries          bool `json:"show_max_recursion_queries,omitempty"`
	ShowMonitor                      bool `json:"show_monitor,omitempty"`
	ShowQueryCapture                 bool `json:"show_query_capture,omitempty"`
	ShowDtcEa                        bool `json:"show_dtc_ea,omitempty"`
	ShowDtcGeoip                     bool `json:"show_dtc_geoip,omitempty"`
	ShowRestartAnycastWithDnsRestart bool `json:"show_restart_anycast_with_dns_restart,omitempty"`
	EnableAll                        bool `json:"enable_all,omitempty"`
	DisableAll                       bool `json:"disable_all,omitempty"`
}

// AdmingroupDnstoplevelcommands represents Infoblox struct admingroup:dnstoplevelcommands
type AdmingroupDnstoplevelcommands struct {
	DdnsAdd          bool `json:"ddns_add,omitempty"`
	DdnsDelete       bool `json:"ddns_delete,omitempty"`
	Delete           bool `json:"delete,omitempty"`
	DnsARecordDelete bool `json:"dns_a_record_delete,omitempty"`
	EnableAll        bool `json:"enable_all,omitempty"`
	DisableAll       bool `json:"disable_all,omitempty"`
}

// AdmingroupDockersetcommands represents Infoblox struct admingroup:dockersetcommands
type AdmingroupDockersetcommands struct {
	SetDockerBridge bool `json:"set_docker_bridge,omitempty"`
	EnableAll       bool `json:"enable_all,omitempty"`
	DisableAll      bool `json:"disable_all,omitempty"`
}

// AdmingroupDockershowcommands represents Infoblox struct admingroup:dockershowcommands
type AdmingroupDockershowcommands struct {
	ShowDockerBridge bool `json:"show_docker_bridge,omitempty"`
	EnableAll        bool `json:"enable_all,omitempty"`
	DisableAll       bool `json:"disable_all,omitempty"`
}

// AdmingroupGridsetcommands represents Infoblox struct admingroup:gridsetcommands
type AdmingroupGridsetcommands struct {
	SetDefaultRevertWindow bool `json:"set_default_revert_window,omitempty"`
	SetDscp                bool `json:"set_dscp,omitempty"`
	SetMembership          bool `json:"set_membership,omitempty"`
	SetNogrid              bool `json:"set_nogrid,omitempty"`
	SetNomastergrid        bool `json:"set_nomastergrid,omitempty"`
	SetPromoteMaster       bool `json:"set_promote_master,omitempty"`
	SetRevertGrid          bool `json:"set_revert_grid,omitempty"`
	SetToken               bool `json:"set_token,omitempty"`
	SetTestPromoteMaster   bool `json:"set_test_promote_master,omitempty"`
	EnableAll              bool `json:"enable_all,omitempty"`
	DisableAll             bool `json:"disable_all,omitempty"`
}

// AdmingroupGridshowcommands represents Infoblox struct admingroup:gridshowcommands
type AdmingroupGridshowcommands struct {
	ShowTestPromoteMaster bool `json:"show_test_promote_master,omitempty"`
	ShowToken             bool `json:"show_token,omitempty"`
	EnableAll             bool `json:"enable_all,omitempty"`
	DisableAll            bool `json:"disable_all,omitempty"`
	ShowDscp              bool `json:"show_dscp,omitempty"`
}

// AdmingroupLicensingsetcommands represents Infoblox struct admingroup:licensingsetcommands
type AdmingroupLicensingsetcommands struct {
	SetLicense               bool `json:"set_license,omitempty"`
	SetReportingResetLicense bool `json:"set_reporting_reset_license,omitempty"`
	SetTempLicense           bool `json:"set_temp_license,omitempty"`
	EnableAll                bool `json:"enable_all,omitempty"`
	DisableAll               bool `json:"disable_all,omitempty"`
}

// AdmingroupLicensingshowcommands represents Infoblox struct admingroup:licensingshowcommands
type AdmingroupLicensingshowcommands struct {
	ShowLicense              bool `json:"show_license,omitempty"`
	ShowLicensePoolContainer bool `json:"show_license_pool_container,omitempty"`
	ShowLicenseUid           bool `json:"show_license_uid,omitempty"`
	EnableAll                bool `json:"enable_all,omitempty"`
	DisableAll               bool `json:"disable_all,omitempty"`
}

// AdmingroupLockoutsetting represents Infoblox struct admingroup:lockoutsetting
type AdmingroupLockoutsetting struct {
	EnableSequentialFailedLoginAttemptsLockout bool   `json:"enable_sequential_failed_login_attempts_lockout,omitempty"`
	SequentialAttempts                         uint32 `json:"sequential_attempts,omitempty"`
	FailedLockoutDuration                      uint32 `json:"failed_lockout_duration,omitempty"`
	NeverUnlockUser                            bool   `json:"never_unlock_user,omitempty"`
}

// AdmingroupMachinecontroltoplevelcommands represents Infoblox struct admingroup:machinecontroltoplevelcommands
type AdmingroupMachinecontroltoplevelcommands struct {
	Reboot     bool `json:"reboot,omitempty"`
	Reset      bool `json:"reset,omitempty"`
	Shutdown   bool `json:"shutdown,omitempty"`
	Restart    bool `json:"restart,omitempty"`
	EnableAll  bool `json:"enable_all,omitempty"`
	DisableAll bool `json:"disable_all,omitempty"`
}

// AdmingroupNetworkingsetcommands represents Infoblox struct admingroup:networkingsetcommands
type AdmingroupNetworkingsetcommands struct {
	SetConnectionLimit      bool `json:"set_connection_limit,omitempty"`
	SetDefaultRoute         bool `json:"set_default_route,omitempty"`
	SetInterface            bool `json:"set_interface,omitempty"`
	SetIpRateLimit          bool `json:"set_ip_rate_limit,omitempty"`
	SetIpv6DisableOnDad     bool `json:"set_ipv6_disable_on_dad,omitempty"`
	SetIpv6Neighbor         bool `json:"set_ipv6_neighbor,omitempty"`
	SetIpv6Ospf             bool `json:"set_ipv6_ospf,omitempty"`
	SetIpv6Status           bool `json:"set_ipv6_status,omitempty"`
	SetLom                  bool `json:"set_lom,omitempty"`
	SetMldVersion1          bool `json:"set_mld_version_1,omitempty"`
	SetNamedRecvSockBufSize bool `json:"set_named_recv_sock_buf_size,omitempty"`
	SetNamedTcpClientsLimit bool `json:"set_named_tcp_clients_limit,omitempty"`
	SetNetwork              bool `json:"set_network,omitempty"`
	SetOspf                 bool `json:"set_ospf,omitempty"`
	SetPrompt               bool `json:"set_prompt,omitempty"`
	SetRemoteConsole        bool `json:"set_remote_console,omitempty"`
	SetStaticRoute          bool `json:"set_static_route,omitempty"`
	SetTcpTimestamps        bool `json:"set_tcp_timestamps,omitempty"`
	SetTrafficCapture       bool `json:"set_traffic_capture,omitempty"`
	SetWinsForwarding       bool `json:"set_wins_forwarding,omitempty"`
	EnableAll               bool `json:"enable_all,omitempty"`
	DisableAll              bool `json:"disable_all,omitempty"`
}

// AdmingroupNetworkingshowcommands represents Infoblox struct admingroup:networkingshowcommands
type AdmingroupNetworkingshowcommands struct {
	ShowConnectionLimit      bool `json:"show_connection_limit,omitempty"`
	ShowConnections          bool `json:"show_connections,omitempty"`
	ShowInterface            bool `json:"show_interface,omitempty"`
	ShowIpRateLimit          bool `json:"show_ip_rate_limit,omitempty"`
	ShowIpv6Bgp              bool `json:"show_ipv6_bgp,omitempty"`
	ShowIpv6DisableOnDad     bool `json:"show_ipv6_disable_on_dad,omitempty"`
	ShowIpv6Neighbor         bool `json:"show_ipv6_neighbor,omitempty"`
	ShowIpv6Ospf             bool `json:"show_ipv6_ospf,omitempty"`
	ShowLom                  bool `json:"show_lom,omitempty"`
	ShowMldVersion           bool `json:"show_mld_version,omitempty"`
	ShowNamedRecvSockBufSize bool `json:"show_named_recv_sock_buf_size,omitempty"`
	ShowNamedTcpClientsLimit bool `json:"show_named_tcp_clients_limit,omitempty"`
	ShowNetwork              bool `json:"show_network,omitempty"`
	ShowOspf                 bool `json:"show_ospf,omitempty"`
	ShowRemoteConsole        bool `json:"show_remote_console,omitempty"`
	ShowRoutes               bool `json:"show_routes,omitempty"`
	ShowStaticRoutes         bool `json:"show_static_routes,omitempty"`
	ShowTcpTimestamps        bool `json:"show_tcp_timestamps,omitempty"`
	ShowTrafficCaptureStatus bool `json:"show_traffic_capture_status,omitempty"`
	ShowWinsForwarding       bool `json:"show_wins_forwarding,omitempty"`
	EnableAll                bool `json:"enable_all,omitempty"`
	DisableAll               bool `json:"disable_all,omitempty"`
}

// AdmingroupPasswordsetting represents Infoblox struct admingroup:passwordsetting
type AdmingroupPasswordsetting struct {
	ExpireEnable bool   `json:"expire_enable,omitempty"`
	ExpireDays   uint32 `json:"expire_days,omitempty"`
	ReminderDays uint32 `json:"reminder_days,omitempty"`
}

// AdmingroupSamlsetting represents Infoblox struct admingroup:samlsetting
type AdmingroupSamlsetting struct {
	AutoCreateUser         bool `json:"auto_create_user,omitempty"`
	PersistAutoCreatedUser bool `json:"persist_auto_created_user,omitempty"`
}

// AdmingroupSecuritysetcommands represents Infoblox struct admingroup:securitysetcommands
type AdmingroupSecuritysetcommands struct {
	SetAdp                          bool `json:"set_adp,omitempty"`
	SetApacheHttpsCert              bool `json:"set_apache_https_cert,omitempty"`
	SetCcMode                       bool `json:"set_cc_mode,omitempty"`
	SetCertificateAuthAdmins        bool `json:"set_certificate_auth_admins,omitempty"`
	SetCertificateAuthServices      bool `json:"set_certificate_auth_services,omitempty"`
	SetCheckAuthNs                  bool `json:"set_check_auth_ns,omitempty"`
	SetCheckSslCertificate          bool `json:"set_check_ssl_certificate,omitempty"`
	SetDisableHttpsCertRegeneration bool `json:"set_disable_https_cert_regeneration,omitempty"`
	SetFipsMode                     bool `json:"set_fips_mode,omitempty"`
	SetReportingCert                bool `json:"set_reporting_cert,omitempty"`
	SetSecurity                     bool `json:"set_security,omitempty"`
	SetSessionTimeout               bool `json:"set_session_timeout,omitempty"`
	SetSubscriberSecureData         bool `json:"set_subscriber_secure_data,omitempty"`
	SetSupportAccess                bool `json:"set_support_access,omitempty"`
	SetSupportInstall               bool `json:"set_support_install,omitempty"`
	SetAdpDebug                     bool `json:"set_adp_debug,omitempty"`
	EnableAll                       bool `json:"enable_all,omitempty"`
	DisableAll                      bool `json:"disable_all,omitempty"`
}

// AdmingroupSecurityshowcommands represents Infoblox struct admingroup:securityshowcommands
type AdmingroupSecurityshowcommands struct {
	ShowFipsMode                bool `json:"show_fips_mode,omitempty"`
	ShowCcMode                  bool `json:"show_cc_mode,omitempty"`
	ShowCertificateAuthAdmins   bool `json:"show_certificate_auth_admins,omitempty"`
	ShowCertificateAuthServices bool `json:"show_certificate_auth_services,omitempty"`
	ShowCheckAuthNs             bool `json:"show_check_auth_ns,omitempty"`
	ShowCheckSslCertificate     bool `json:"show_check_ssl_certificate,omitempty"`
	ShowSecurity                bool `json:"show_security,omitempty"`
	ShowSessionTimeout          bool `json:"show_session_timeout,omitempty"`
	ShowSubscriberSecureData    bool `json:"show_subscriber_secure_data,omitempty"`
	ShowSupportAccess           bool `json:"show_support_access,omitempty"`
	ShowVpnCertDates            bool `json:"show_vpn_cert_dates,omitempty"`
	ShowAdp                     bool `json:"show_adp,omitempty"`
	ShowAdpDebug                bool `json:"show_adp_debug,omitempty"`
	EnableAll                   bool `json:"enable_all,omitempty"`
	DisableAll                  bool `json:"disable_all,omitempty"`
}

// AdmingroupTroubleshootingtoplevelcommands represents Infoblox struct admingroup:troubleshootingtoplevelcommands
type AdmingroupTroubleshootingtoplevelcommands struct {
	Ping           bool `json:"ping,omitempty"`
	Ping6          bool `json:"ping6,omitempty"`
	Strace         bool `json:"strace,omitempty"`
	Traceroute     bool `json:"traceroute,omitempty"`
	TrafficCapture bool `json:"traffic_capture,omitempty"`
	Dig            bool `json:"dig,omitempty"`
	Rotate         bool `json:"rotate,omitempty"`
	Snmpwalk       bool `json:"snmpwalk,omitempty"`
	Snmpget        bool `json:"snmpget,omitempty"`
	Console        bool `json:"console,omitempty"`
	Tracepath      bool `json:"tracepath,omitempty"`
	EnableAll      bool `json:"enable_all,omitempty"`
	DisableAll     bool `json:"disable_all,omitempty"`
}

// Adsites represents Infoblox struct adsites
type Adsites struct {
	UseDefaultIpSiteLink       bool       `json:"use_default_ip_site_link,omitempty"`
	DefaultIpSiteLink          string     `json:"default_ip_site_link,omitempty"`
	UseLogin                   bool       `json:"use_login,omitempty"`
	LoginName                  string     `json:"login_name,omitempty"`
	LoginPassword              string     `json:"login_password,omitempty"`
	UseSynchronizationMinDelay bool       `json:"use_synchronization_min_delay,omitempty"`
	SynchronizationMinDelay    uint32     `json:"synchronization_min_delay,omitempty"`
	UseLdapTimeout             bool       `json:"use_ldap_timeout,omitempty"`
	LdapTimeout                uint32     `json:"ldap_timeout,omitempty"`
	LdapAuthPort               uint32     `json:"ldap_auth_port,omitempty"`
	LdapEncryption             string     `json:"ldap_encryption,omitempty"`
	Managed                    bool       `json:"managed,omitempty"`
	ReadOnly                   bool       `json:"read_only,omitempty"`
	LastSyncTs                 *time.Time `json:"last_sync_ts,omitempty"`
	LastSyncStatus             string     `json:"last_sync_status,omitempty"`
	LastSyncDetail             string     `json:"last_sync_detail,omitempty"`
	SupportsIpv6               bool       `json:"supports_ipv6,omitempty"`
}

// Advisortestconnection represents Infoblox struct advisortestconnection
type Advisortestconnection struct{}

// Allocatelicenses represents Infoblox struct allocatelicenses
type Allocatelicenses struct{}

// Atpobjectreset represents Infoblox struct atpobjectreset
type Atpobjectreset struct{}

// Awsrte53recordinfo represents Infoblox struct awsrte53recordinfo
type Awsrte53recordinfo struct {
	AliasTargetDnsName              string `json:"alias_target_dns_name,omitempty"`
	AliasTargetHostedZoneId         string `json:"alias_target_hosted_zone_id,omitempty"`
	AliasTargetEvaluateTargetHealth bool   `json:"alias_target_evaluate_target_health,omitempty"`
	Failover                        string `json:"failover,omitempty"`
	GeolocationContinentCode        string `json:"geolocation_continent_code,omitempty"`
	GeolocationCountryCode          string `json:"geolocation_country_code,omitempty"`
	GeolocationSubdivisionCode      string `json:"geolocation_subdivision_code,omitempty"`
	HealthCheckId                   string `json:"health_check_id,omitempty"`
	Region                          string `json:"region,omitempty"`
	SetIdentifier                   string `json:"set_identifier,omitempty"`
	Type                            string `json:"type,omitempty"`
	Weight                          uint32 `json:"weight,omitempty"`
}

// Awsrte53task represents Infoblox struct awsrte53task
type Awsrte53task struct {
	Name             string     `json:"name,omitempty"`
	Disabled         bool       `json:"disabled,omitempty"`
	State            string     `json:"state,omitempty"`
	StateMsg         string     `json:"state_msg,omitempty"`
	Filter           string     `json:"filter,omitempty"`
	ScheduleInterval uint32     `json:"schedule_interval,omitempty"`
	ScheduleUnits    string     `json:"schedule_units,omitempty"`
	AwsUser          string     `json:"aws_user,omitempty"`
	StatusTimestamp  *time.Time `json:"status_timestamp,omitempty"`
	LastRun          *time.Time `json:"last_run,omitempty"`
	SyncPublicZones  bool       `json:"sync_public_zones,omitempty"`
	SyncPrivateZones bool       `json:"sync_private_zones,omitempty"`
	ZoneCount        uint32     `json:"zone_count,omitempty"`
	CredentialsType  string     `json:"credentials_type,omitempty"`
}

// Awsrte53zoneinfo represents Infoblox struct awsrte53zoneinfo
type Awsrte53zoneinfo struct {
	AssociatedVpcs  []string `json:"associated_vpcs,omitempty"`
	CallerReference string   `json:"caller_reference,omitempty"`
	DelegationSetId string   `json:"delegation_set_id,omitempty"`
	HostedZoneId    string   `json:"hosted_zone_id,omitempty"`
	NameServers     []string `json:"name_servers,omitempty"`
	RecordSetCount  uint32   `json:"record_set_count,omitempty"`
	Type            string   `json:"type,omitempty"`
}

// Bgpas represents Infoblox struct bgpas
type Bgpas struct {
	As         uint32         `json:"as,omitempty"`
	Keepalive  uint32         `json:"keepalive,omitempty"`
	Holddown   uint32         `json:"holddown,omitempty"`
	Neighbors  []*Bgpneighbor `json:"neighbors,omitempty"`
	LinkDetect bool           `json:"link_detect,omitempty"`
}

// Bgpneighbor represents Infoblox struct bgpneighbor
type Bgpneighbor struct {
	Interface          string `json:"interface,omitempty"`
	NeighborIp         string `json:"neighbor_ip,omitempty"`
	RemoteAs           uint32 `json:"remote_as,omitempty"`
	AuthenticationMode string `json:"authentication_mode,omitempty"`
	BgpNeighborPass    string `json:"bgp_neighbor_pass,omitempty"`
	Comment            string `json:"comment,omitempty"`
	Multihop           bool   `json:"multihop,omitempty"`
	MultihopTtl        uint32 `json:"multihop_ttl,omitempty"`
	BfdTemplate        string `json:"bfd_template,omitempty"`
	EnableBfd          bool   `json:"enable_bfd,omitempty"`
}

// CapacityreportObjectcount represents Infoblox struct capacityreport:objectcount
type CapacityreportObjectcount struct {
	TypeName string `json:"type_name,omitempty"`
	Count    uint32 `json:"count,omitempty"`
}

// CaptiveportalFile represents Infoblox struct captiveportal:file
type CaptiveportalFile struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

// Changedobject represents Infoblox struct changedobject
type Changedobject struct {
	Action     string   `json:"action,omitempty"`
	Name       string   `json:"name,omitempty"`
	Type       string   `json:"type,omitempty"`
	ObjectType string   `json:"object_type,omitempty"`
	Properties []string `json:"properties,omitempty"`
}

// Checkldapserversettings represents Infoblox struct checkldapserversettings
type Checkldapserversettings struct{}

// Checkradiusserversettings represents Infoblox struct checkradiusserversettings
type Checkradiusserversettings struct{}

// Checktacacsplusserversettings represents Infoblox struct checktacacsplusserversettings
type Checktacacsplusserversettings struct{}

// CiscoiseEaassociation represents Infoblox struct ciscoise:eaassociation
type CiscoiseEaassociation struct {
	Name     string `json:"name,omitempty"`
	MappedEa string `json:"mapped_ea,omitempty"`
}

// CiscoisePublishsetting represents Infoblox struct ciscoise:publishsetting
type CiscoisePublishsetting struct {
	EnabledAttributes []string `json:"enabled_attributes,omitempty"`
}

// CiscoiseSubscribesetting represents Infoblox struct ciscoise:subscribesetting
type CiscoiseSubscribesetting struct {
	EnabledAttributes  []string                 `json:"enabled_attributes,omitempty"`
	MappedEaAttributes []*CiscoiseEaassociation `json:"mapped_ea_attributes,omitempty"`
}

// Cleardiscoverydata represents Infoblox struct cleardiscoverydata
type Cleardiscoverydata struct{}

// Cleardnscache represents Infoblox struct cleardnscache
type Cleardnscache struct{}

// Clearnacauthcacheparams represents Infoblox struct clearnacauthcacheparams
type Clearnacauthcacheparams struct{}

// Clearnetworkportassignment represents Infoblox struct clearnetworkportassignment
type Clearnetworkportassignment struct{}

// Clearworkerlog represents Infoblox struct clearworkerlog
type Clearworkerlog struct{}

// Clientsubnetdomain represents Infoblox struct clientsubnetdomain
type Clientsubnetdomain struct {
	Domain     string `json:"domain,omitempty"`
	Permission string `json:"permission,omitempty"`
}

// Controlipaddress represents Infoblox struct controlipaddress
type Controlipaddress struct{}

// Controlswitchport represents Infoblox struct controlswitchport
type Controlswitchport struct{}

// Copyrpzrecords represents Infoblox struct copyrpzrecords
type Copyrpzrecords struct{}

// Copyzonerecords represents Infoblox struct copyzonerecords
type Copyzonerecords struct{}

// Csverrorlog represents Infoblox struct csverrorlog
type Csverrorlog struct{}

// Csvexport represents Infoblox struct csvexport
type Csvexport struct{}

// Csvimport represents Infoblox struct csvimport
type Csvimport struct{}

// Csvsnapshot represents Infoblox struct csvsnapshot
type Csvsnapshot struct{}

// Csvuploaded represents Infoblox struct csvuploaded
type Csvuploaded struct{}

// Datagetcomplete represents Infoblox struct datagetcomplete
type Datagetcomplete struct{}

// Datauploadinit represents Infoblox struct datauploadinit
type Datauploadinit struct{}

// Dhcpddns represents Infoblox struct dhcpddns
type Dhcpddns struct {
	ZoneMatch      string `json:"zone_match,omitempty"`
	DnsGridZone    string `json:"dns_grid_zone,omitempty"`
	DnsGridPrimary string `json:"dns_grid_primary,omitempty"`
	DnsExtZone     string `json:"dns_ext_zone,omitempty"`
	DnsExtPrimary  string `json:"dns_ext_primary,omitempty"`
}

// Dhcpmember represents Infoblox struct dhcpmember
type Dhcpmember struct {
	Ipv4Addr string `json:"ipv4addr,omitempty"`
	Ipv6Addr string `json:"ipv6addr,omitempty"`
	Name     string `json:"name,omitempty"`
}

// Dhcpoption represents Infoblox struct dhcpoption
type Dhcpoption struct {
	Name        string `json:"name,omitempty"`
	Num         uint32 `json:"num,omitempty"`
	VendorClass string `json:"vendor_class,omitempty"`
	Value       string `json:"value,omitempty"`
	UseOption   bool   `json:"use_option,omitempty"`
}

// Dhcpserver represents Infoblox struct dhcpserver
type Dhcpserver struct {
	UseLogin                   bool       `json:"use_login,omitempty"`
	LoginName                  string     `json:"login_name,omitempty"`
	LoginPassword              string     `json:"login_password,omitempty"`
	Managed                    bool       `json:"managed,omitempty"`
	NextSyncControl            string     `json:"next_sync_control,omitempty"`
	Status                     string     `json:"status,omitempty"`
	StatusLastUpdated          *time.Time `json:"status_last_updated,omitempty"`
	UseEnableMonitoring        bool       `json:"use_enable_monitoring,omitempty"`
	EnableMonitoring           bool       `json:"enable_monitoring,omitempty"`
	UseEnableInvalidMac        bool       `json:"use_enable_invalid_mac,omitempty"`
	EnableInvalidMac           bool       `json:"enable_invalid_mac,omitempty"`
	SupportsFailover           bool       `json:"supports_failover,omitempty"`
	UseSynchronizationMinDelay bool       `json:"use_synchronization_min_delay,omitempty"`
	SynchronizationMinDelay    uint32     `json:"synchronization_min_delay,omitempty"`
}

// DiscoveryAdvancedpollsetting represents Infoblox struct discovery:advancedpollsetting
type DiscoveryAdvancedpollsetting struct {
	TcpScanTechnique                      string `json:"tcp_scan_technique,omitempty"`
	PingTimeout                           uint32 `json:"ping_timeout,omitempty"`
	PingRetries                           uint32 `json:"ping_retries,omitempty"`
	PurgeExpiredDeviceData                uint32 `json:"purge_expired_device_data,omitempty"`
	EnablePurgeExpiredEndhostData         bool   `json:"enable_purge_expired_endhost_data,omitempty"`
	PurgeExpiredEndhostData               uint32 `json:"purge_expired_endhost_data,omitempty"`
	ArpAggregateLimit                     uint32 `json:"arp_aggregate_limit,omitempty"`
	RouteLimit                            uint32 `json:"route_limit,omitempty"`
	PingSweepInterval                     uint32 `json:"ping_sweep_interval,omitempty"`
	ArpCacheRefreshInterval               uint32 `json:"arp_cache_refresh_interval,omitempty"`
	PollingAuthenticateSnmpv2cOrLaterOnly bool   `json:"polling_authenticate_snmpv2c_or_later_only,omitempty"`
	DisableDiscoveryOutsideIpam           bool   `json:"disable_discovery_outside_ipam,omitempty"`
	DhcpRouterAsSeed                      bool   `json:"dhcp_router_as_seed,omitempty"`
	SyslogIpamEvents                      bool   `json:"syslog_ipam_events,omitempty"`
	SyslogNetworkEvents                   bool   `json:"syslog_network_events,omitempty"`
}

// DiscoveryAdvancedsdnpollsettings represents Infoblox struct discovery:advancedsdnpollsettings
type DiscoveryAdvancedsdnpollsettings struct {
	NetworksMappingPolicy          string `json:"networks_mapping_policy,omitempty"`
	DisableSdnDiscoveryOutsideIpam bool   `json:"disable_sdn_discovery_outside_ipam,omitempty"`
}

// DiscoveryAdvisorsetting represents Infoblox struct discovery:advisorsetting
type DiscoveryAdvisorsetting struct {
	EnableProxy               bool       `json:"enable_proxy,omitempty"`
	ProxyAddress              string     `json:"proxy_address,omitempty"`
	ProxyPort                 uint32     `json:"proxy_port,omitempty"`
	UseProxyUsernamePasswd    bool       `json:"use_proxy_username_passwd,omitempty"`
	ProxyUsername             string     `json:"proxy_username,omitempty"`
	ProxyPassword             string     `json:"proxy_password,omitempty"`
	ExecutionInterval         uint32     `json:"execution_interval,omitempty"`
	ExecutionHour             uint32     `json:"execution_hour,omitempty"`
	NetworkInterfaceType      string     `json:"network_interface_type,omitempty"`
	NetworkInterfaceVirtualIp string     `json:"network_interface_virtual_ip,omitempty"`
	Address                   string     `json:"address,omitempty"`
	Port                      uint32     `json:"port,omitempty"`
	AuthType                  string     `json:"auth_type,omitempty"`
	AuthToken                 string     `json:"auth_token,omitempty"`
	Username                  string     `json:"username,omitempty"`
	Password                  string     `json:"password,omitempty"`
	MinSeverity               string     `json:"min_severity,omitempty"`
	LastExecTime              *time.Time `json:"last_exec_time,omitempty"`
	LastExecStatus            string     `json:"last_exec_status,omitempty"`
	LastExecDetails           string     `json:"last_exec_details,omitempty"`
	LastRunNowTime            *time.Time `json:"last_run_now_time,omitempty"`
	LastRunNowStatus          string     `json:"last_run_now_status,omitempty"`
	LastRunNowDetails         string     `json:"last_run_now_details,omitempty"`
}

// DiscoveryAutoconversionsetting represents Infoblox struct discovery:autoconversionsetting
type DiscoveryAutoconversionsetting struct {
	NetworkView string `json:"network_view,omitempty"`
	Type        string `json:"type,omitempty"`
	Format      string `json:"format,omitempty"`
	Condition   string `json:"condition,omitempty"`
	Comment     string `json:"comment,omitempty"`
}

// DiscoveryBasicpollsettings represents Infoblox struct discovery:basicpollsettings
type DiscoveryBasicpollsettings struct {
	PortScanning                            bool             `json:"port_scanning,omitempty"`
	DeviceProfile                           bool             `json:"device_profile,omitempty"`
	SnmpCollection                          bool             `json:"snmp_collection,omitempty"`
	CliCollection                           bool             `json:"cli_collection,omitempty"`
	NetbiosScanning                         bool             `json:"netbios_scanning,omitempty"`
	CompletePingSweep                       bool             `json:"complete_ping_sweep,omitempty"`
	SmartSubnetPingSweep                    bool             `json:"smart_subnet_ping_sweep,omitempty"`
	AutoArpRefreshBeforeSwitchPortPolling   bool             `json:"auto_arp_refresh_before_switch_port_polling,omitempty"`
	SwitchPortDataCollectionPolling         string           `json:"switch_port_data_collection_polling,omitempty"`
	SwitchPortDataCollectionPollingSchedule *SettingSchedule `json:"switch_port_data_collection_polling_schedule,omitempty"`
	SwitchPortDataCollectionPollingInterval uint32           `json:"switch_port_data_collection_polling_interval,omitempty"`
	CredentialGroup                         string           `json:"credential_group,omitempty"`
}

// DiscoveryBasicsdnpollsettings represents Infoblox struct discovery:basicsdnpollsettings
type DiscoveryBasicsdnpollsettings struct {
	SdnDiscovery           bool             `json:"sdn_discovery,omitempty"`
	DefaultNetworkView     string           `json:"default_network_view,omitempty"`
	EndHostPolling         string           `json:"end_host_polling,omitempty"`
	EndHostPollingInterval uint32           `json:"end_host_polling_interval,omitempty"`
	EndHostPollingSchedule *SettingSchedule `json:"end_host_polling_schedule,omitempty"`
}

// DiscoveryClicredential represents Infoblox struct discovery:clicredential
type DiscoveryClicredential struct {
	User            string `json:"user,omitempty"`
	Password        string `json:"password,omitempty"`
	CredentialType  string `json:"credential_type,omitempty"`
	Comment         string `json:"comment,omitempty"`
	Id              uint32 `json:"id,omitempty"`
	CredentialGroup string `json:"credential_group,omitempty"`
}

// DiscoveryDevicePortstatistics represents Infoblox struct discovery:device:portstatistics
type DiscoveryDevicePortstatistics struct {
	InterfacesCount        uint32 `json:"interfaces_count,omitempty"`
	AdminUpOperUpCount     uint32 `json:"admin_up_oper_up_count,omitempty"`
	AdminUpOperDownCount   uint32 `json:"admin_up_oper_down_count,omitempty"`
	AdminDownOperDownCount uint32 `json:"admin_down_oper_down_count,omitempty"`
}

// DiscoveryIfaddrinfo represents Infoblox struct discovery:ifaddrinfo
type DiscoveryIfaddrinfo struct {
	Address       string `json:"address,omitempty"`
	AddressObject string `json:"address_object,omitempty"`
	Network       string `json:"network,omitempty"`
}

// DiscoveryNetworkinfo represents Infoblox struct discovery:networkinfo
type DiscoveryNetworkinfo struct {
	Network    string `json:"network,omitempty"`
	NetworkStr string `json:"network_str,omitempty"`
}

// DiscoveryPort represents Infoblox struct discovery:port
type DiscoveryPort struct {
	Port    uint32 `json:"port,omitempty"`
	Type    string `json:"type,omitempty"`
	Comment string `json:"comment,omitempty"`
}

// DiscoveryPortConfigAdminstatus represents Infoblox struct discovery:port:config:adminstatus
type DiscoveryPortConfigAdminstatus struct {
	Status  string                           `json:"status,omitempty"`
	Details *DiscoveryPortControlTaskdetails `json:"details,omitempty"`
}

// DiscoveryPortConfigDescription represents Infoblox struct discovery:port:config:description
type DiscoveryPortConfigDescription struct {
	Description string                           `json:"description,omitempty"`
	Details     *DiscoveryPortControlTaskdetails `json:"details,omitempty"`
}

// DiscoveryPortConfigVlaninfo represents Infoblox struct discovery:port:config:vlaninfo
type DiscoveryPortConfigVlaninfo struct {
	DataVlanInfo  *DiscoveryVlaninfo               `json:"data_vlan_info,omitempty"`
	VoiceVlanInfo *DiscoveryVlaninfo               `json:"voice_vlan_info,omitempty"`
	Details       *DiscoveryPortControlTaskdetails `json:"details,omitempty"`
}

// DiscoveryPortControlTaskdetails represents Infoblox struct discovery:port:control:taskdetails
type DiscoveryPortControlTaskdetails struct {
	Id             uint32 `json:"id,omitempty"`
	Status         string `json:"status,omitempty"`
	IsSynchronized bool   `json:"is_synchronized,omitempty"`
}

// DiscoveryScaninterface represents Infoblox struct discovery:scaninterface
type DiscoveryScaninterface struct {
	NetworkView   string `json:"network_view,omitempty"`
	Type          string `json:"type,omitempty"`
	ScanVirtualIp string `json:"scan_virtual_ip,omitempty"`
}

// DiscoverySdnconfig represents Infoblox struct discovery:sdnconfig
type DiscoverySdnconfig struct {
	SdnType                   string   `json:"sdn_type,omitempty"`
	Addresses                 []string `json:"addresses,omitempty"`
	NetworkView               string   `json:"network_view,omitempty"`
	Protocol                  string   `json:"protocol,omitempty"`
	Handle                    string   `json:"handle,omitempty"`
	Password                  string   `json:"password,omitempty"`
	Username                  string   `json:"username,omitempty"`
	ApiKey                    string   `json:"api_key,omitempty"`
	OnPrem                    bool     `json:"on_prem,omitempty"`
	UseGlobalProxy            bool     `json:"use_global_proxy,omitempty"`
	Comment                   string   `json:"comment,omitempty"`
	NetworkInterfaceType      string   `json:"network_interface_type,omitempty"`
	NetworkInterfaceVirtualIp string   `json:"network_interface_virtual_ip,omitempty"`
	Uuid                      string   `json:"uuid,omitempty"`
}

// DiscoverySeedrouter represents Infoblox struct discovery:seedrouter
type DiscoverySeedrouter struct {
	Address     string `json:"address,omitempty"`
	NetworkView string `json:"network_view,omitempty"`
	Comment     string `json:"comment,omitempty"`
}

// DiscoverySnmp3credential represents Infoblox struct discovery:snmp3credential
type DiscoverySnmp3credential struct {
	User                   string `json:"user,omitempty"`
	AuthenticationProtocol string `json:"authentication_protocol,omitempty"`
	AuthenticationPassword string `json:"authentication_password,omitempty"`
	PrivacyProtocol        string `json:"privacy_protocol,omitempty"`
	PrivacyPassword        string `json:"privacy_password,omitempty"`
	Comment                string `json:"comment,omitempty"`
	CredentialGroup        string `json:"credential_group,omitempty"`
}

// DiscoverySnmpcredential represents Infoblox struct discovery:snmpcredential
type DiscoverySnmpcredential struct {
	CommunityString string `json:"community_string,omitempty"`
	Comment         string `json:"comment,omitempty"`
	CredentialGroup string `json:"credential_group,omitempty"`
}

// DiscoveryStatusinfo represents Infoblox struct discovery:statusinfo
type DiscoveryStatusinfo struct {
	Status    string     `json:"status,omitempty"`
	Message   string     `json:"message,omitempty"`
	Timestamp *time.Time `json:"timestamp,omitempty"`
}

// DiscoveryVlaninfo represents Infoblox struct discovery:vlaninfo
type DiscoveryVlaninfo struct {
	Id   uint32 `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// DiscoveryVrfmappingrule represents Infoblox struct discovery:vrfmappingrule
type DiscoveryVrfmappingrule struct {
	NetworkView string `json:"network_view,omitempty"`
	Criteria    string `json:"criteria,omitempty"`
	Comment     string `json:"comment,omitempty"`
}

// Discoverydata represents Infoblox struct discoverydata
type Discoverydata struct {
	DeviceModel                     string     `json:"device_model,omitempty"`
	DevicePortName                  string     `json:"device_port_name,omitempty"`
	DevicePortType                  string     `json:"device_port_type,omitempty"`
	DeviceType                      string     `json:"device_type,omitempty"`
	DeviceVendor                    string     `json:"device_vendor,omitempty"`
	DiscoveredName                  string     `json:"discovered_name,omitempty"`
	Discoverer                      string     `json:"discoverer,omitempty"`
	Duid                            string     `json:"duid,omitempty"`
	FirstDiscovered                 *time.Time `json:"first_discovered,omitempty"`
	IprgNo                          uint32     `json:"iprg_no,omitempty"`
	IprgState                       string     `json:"iprg_state,omitempty"`
	IprgType                        string     `json:"iprg_type,omitempty"`
	LastDiscovered                  *time.Time `json:"last_discovered,omitempty"`
	MacAddress                      string     `json:"mac_address,omitempty"`
	MgmtIpAddress                   string     `json:"mgmt_ip_address,omitempty"`
	NetbiosName                     string     `json:"netbios_name,omitempty"`
	NetworkComponentDescription     string     `json:"network_component_description,omitempty"`
	NetworkComponentIp              string     `json:"network_component_ip,omitempty"`
	NetworkComponentModel           string     `json:"network_component_model,omitempty"`
	NetworkComponentName            string     `json:"network_component_name,omitempty"`
	NetworkComponentPortDescription string     `json:"network_component_port_description,omitempty"`
	NetworkComponentPortName        string     `json:"network_component_port_name,omitempty"`
	NetworkComponentPortNumber      string     `json:"network_component_port_number,omitempty"`
	NetworkComponentType            string     `json:"network_component_type,omitempty"`
	NetworkComponentVendor          string     `json:"network_component_vendor,omitempty"`
	OpenPorts                       string     `json:"open_ports,omitempty"`
	Os                              string     `json:"os,omitempty"`
	PortDuplex                      string     `json:"port_duplex,omitempty"`
	PortLinkStatus                  string     `json:"port_link_status,omitempty"`
	PortSpeed                       string     `json:"port_speed,omitempty"`
	PortStatus                      string     `json:"port_status,omitempty"`
	PortType                        string     `json:"port_type,omitempty"`
	PortVlanDescription             string     `json:"port_vlan_description,omitempty"`
	PortVlanName                    string     `json:"port_vlan_name,omitempty"`
	PortVlanNumber                  string     `json:"port_vlan_number,omitempty"`
	VAdapter                        string     `json:"v_adapter,omitempty"`
	VCluster                        string     `json:"v_cluster,omitempty"`
	VDatacenter                     string     `json:"v_datacenter,omitempty"`
	VEntityName                     string     `json:"v_entity_name,omitempty"`
	VEntityType                     string     `json:"v_entity_type,omitempty"`
	VHost                           string     `json:"v_host,omitempty"`
	VSwitch                         string     `json:"v_switch,omitempty"`
	VmiName                         string     `json:"vmi_name,omitempty"`
	VmiId                           string     `json:"vmi_id,omitempty"`
	VlanPortGroup                   string     `json:"vlan_port_group,omitempty"`
	VswitchName                     string     `json:"vswitch_name,omitempty"`
	VswitchId                       string     `json:"vswitch_id,omitempty"`
	VswitchType                     string     `json:"vswitch_type,omitempty"`
	VswitchIpv6Enabled              bool       `json:"vswitch_ipv6_enabled,omitempty"`
	VportName                       string     `json:"vport_name,omitempty"`
	VportMacAddress                 string     `json:"vport_mac_address,omitempty"`
	VportLinkStatus                 string     `json:"vport_link_status,omitempty"`
	VportConfSpeed                  string     `json:"vport_conf_speed,omitempty"`
	VportConfMode                   string     `json:"vport_conf_mode,omitempty"`
	VportSpeed                      string     `json:"vport_speed,omitempty"`
	VportMode                       string     `json:"vport_mode,omitempty"`
	VswitchSegmentType              string     `json:"vswitch_segment_type,omitempty"`
	VswitchSegmentName              string     `json:"vswitch_segment_name,omitempty"`
	VswitchSegmentId                string     `json:"vswitch_segment_id,omitempty"`
	VswitchSegmentPortGroup         string     `json:"vswitch_segment_port_group,omitempty"`
	VswitchAvailablePortsCount      uint32     `json:"vswitch_available_ports_count,omitempty"`
	VswitchTepType                  string     `json:"vswitch_tep_type,omitempty"`
	VswitchTepIp                    string     `json:"vswitch_tep_ip,omitempty"`
	VswitchTepPortGroup             string     `json:"vswitch_tep_port_group,omitempty"`
	VswitchTepVlan                  string     `json:"vswitch_tep_vlan,omitempty"`
	VswitchTepDhcpServer            string     `json:"vswitch_tep_dhcp_server,omitempty"`
	VswitchTepMulticast             string     `json:"vswitch_tep_multicast,omitempty"`
	VmhostIpAddress                 string     `json:"vmhost_ip_address,omitempty"`
	VmhostName                      string     `json:"vmhost_name,omitempty"`
	VmhostMacAddress                string     `json:"vmhost_mac_address,omitempty"`
	VmhostSubnetCidr                uint32     `json:"vmhost_subnet_cidr,omitempty"`
	VmhostNicNames                  string     `json:"vmhost_nic_names,omitempty"`
	VmiTenantId                     string     `json:"vmi_tenant_id,omitempty"`
	CmpType                         string     `json:"cmp_type,omitempty"`
	VmiIpType                       string     `json:"vmi_ip_type,omitempty"`
	VmiPrivateAddress               string     `json:"vmi_private_address,omitempty"`
	VmiIsPublicAddress              bool       `json:"vmi_is_public_address,omitempty"`
	CiscoIseSsid                    string     `json:"cisco_ise_ssid,omitempty"`
	CiscoIseEndpointProfile         string     `json:"cisco_ise_endpoint_profile,omitempty"`
	CiscoIseSessionState            string     `json:"cisco_ise_session_state,omitempty"`
	CiscoIseSecurityGroup           string     `json:"cisco_ise_security_group,omitempty"`
	TaskName                        string     `json:"task_name,omitempty"`
	NetworkComponentLocation        string     `json:"network_component_location,omitempty"`
	NetworkComponentContact         string     `json:"network_component_contact,omitempty"`
	DeviceLocation                  string     `json:"device_location,omitempty"`
	DeviceContact                   string     `json:"device_contact,omitempty"`
	ApName                          string     `json:"ap_name,omitempty"`
	ApIpAddress                     string     `json:"ap_ip_address,omitempty"`
	ApSsid                          string     `json:"ap_ssid,omitempty"`
	BridgeDomain                    string     `json:"bridge_domain,omitempty"`
	EndpointGroups                  string     `json:"endpoint_groups,omitempty"`
	Tenant                          string     `json:"tenant,omitempty"`
	VrfName                         string     `json:"vrf_name,omitempty"`
	VrfDescription                  string     `json:"vrf_description,omitempty"`
	VrfRd                           string     `json:"vrf_rd,omitempty"`
	BgpAs                           uint32     `json:"bgp_as,omitempty"`
}

// Discoverydataconversion represents Infoblox struct discoverydataconversion
type Discoverydataconversion struct{}

// Discoverydiagnostic represents Infoblox struct discoverydiagnostic
type Discoverydiagnostic struct{}

// Discoverydiagnosticstatus represents Infoblox struct discoverydiagnosticstatus
type Discoverydiagnosticstatus struct{}

// Discoverytaskport represents Infoblox struct discoverytaskport
type Discoverytaskport struct {
	Number  uint32 `json:"number,omitempty"`
	Comment string `json:"comment,omitempty"`
}

// Discoverytaskvserver represents Infoblox struct discoverytaskvserver
type Discoverytaskvserver struct {
	Disable            bool   `json:"disable,omitempty"`
	ConnectionProtocol string `json:"connection_protocol,omitempty"`
	FqdnOrIp           string `json:"fqdn_or_ip,omitempty"`
	Password           string `json:"password,omitempty"`
	Port               uint32 `json:"port,omitempty"`
	Username           string `json:"username,omitempty"`
}

// Dnssecexport represents Infoblox struct dnssecexport
type Dnssecexport struct{}

// Dnssecgetkskrollover represents Infoblox struct dnssecgetkskrollover
type Dnssecgetkskrollover struct{}

// Dnssecgetzonekeys represents Infoblox struct dnssecgetzonekeys
type Dnssecgetzonekeys struct{}

// Dnsseckey represents Infoblox struct dnsseckey
type Dnsseckey struct {
	Tag           uint32     `json:"tag,omitempty"`
	Status        string     `json:"status,omitempty"`
	NextEventDate *time.Time `json:"next_event_date,omitempty"`
	Type          string     `json:"type,omitempty"`
	Algorithm     string     `json:"algorithm,omitempty"`
	PublicKey     string     `json:"public_key,omitempty"`
}

// Dnsseckeyalgorithm represents Infoblox struct dnsseckeyalgorithm
type Dnsseckeyalgorithm struct {
	Algorithm string `json:"algorithm,omitempty"`
	Size      uint32 `json:"size,omitempty"`
}

// Dnsseckeyparams represents Infoblox struct dnsseckeyparams
type Dnsseckeyparams struct {
	EnableKskAutoRollover         bool                  `json:"enable_ksk_auto_rollover,omitempty"`
	KskAlgorithm                  string                `json:"ksk_algorithm,omitempty"`
	KskAlgorithms                 []*Dnsseckeyalgorithm `json:"ksk_algorithms,omitempty"`
	KskRollover                   uint32                `json:"ksk_rollover,omitempty"`
	KskSize                       uint32                `json:"ksk_size,omitempty"`
	NextSecureType                string                `json:"next_secure_type,omitempty"`
	KskRolloverNotificationConfig string                `json:"ksk_rollover_notification_config,omitempty"`
	KskSnmpNotificationEnabled    bool                  `json:"ksk_snmp_notification_enabled,omitempty"`
	KskEmailNotificationEnabled   bool                  `json:"ksk_email_notification_enabled,omitempty"`
	Nsec3SaltMinLength            uint32                `json:"nsec3_salt_min_length,omitempty"`
	Nsec3SaltMaxLength            uint32                `json:"nsec3_salt_max_length,omitempty"`
	Nsec3Iterations               uint32                `json:"nsec3_iterations,omitempty"`
	SignatureExpiration           uint32                `json:"signature_expiration,omitempty"`
	ZskAlgorithm                  string                `json:"zsk_algorithm,omitempty"`
	ZskAlgorithms                 []*Dnsseckeyalgorithm `json:"zsk_algorithms,omitempty"`
	ZskRollover                   uint32                `json:"zsk_rollover,omitempty"`
	ZskRolloverMechanism          string                `json:"zsk_rollover_mechanism,omitempty"`
	ZskSize                       uint32                `json:"zsk_size,omitempty"`
}

// Dnssecoperation represents Infoblox struct dnssecoperation
type Dnssecoperation struct{}

// Dnssecsetzonekeys represents Infoblox struct dnssecsetzonekeys
type Dnssecsetzonekeys struct{}

// Dnssectrustedkey represents Infoblox struct dnssectrustedkey
type Dnssectrustedkey struct {
	Fqdn               string `json:"fqdn,omitempty"`
	Algorithm          string `json:"algorithm,omitempty"`
	Key                string `json:"key,omitempty"`
	SecureEntryPoint   bool   `json:"secure_entry_point,omitempty"`
	DnssecMustBeSecure bool   `json:"dnssec_must_be_secure,omitempty"`
}

// Dnsserver represents Infoblox struct dnsserver
type Dnsserver struct {
	UseLogin                   bool       `json:"use_login,omitempty"`
	LoginName                  string     `json:"login_name,omitempty"`
	LoginPassword              string     `json:"login_password,omitempty"`
	Managed                    bool       `json:"managed,omitempty"`
	NextSyncControl            string     `json:"next_sync_control,omitempty"`
	Status                     string     `json:"status,omitempty"`
	StatusDetail               string     `json:"status_detail,omitempty"`
	StatusLastUpdated          *time.Time `json:"status_last_updated,omitempty"`
	LastSyncTs                 *time.Time `json:"last_sync_ts,omitempty"`
	LastSyncStatus             string     `json:"last_sync_status,omitempty"`
	LastSyncDetail             string     `json:"last_sync_detail,omitempty"`
	Forwarders                 string     `json:"forwarders,omitempty"`
	SupportsIpv6               bool       `json:"supports_ipv6,omitempty"`
	SupportsIpv6Reverse        bool       `json:"supports_ipv6_reverse,omitempty"`
	SupportsRrDname            bool       `json:"supports_rr_dname,omitempty"`
	SupportsDnssec             bool       `json:"supports_dnssec,omitempty"`
	SupportsActiveDirectory    bool       `json:"supports_active_directory,omitempty"`
	Address                    string     `json:"address,omitempty"`
	SupportsRrNaptr            bool       `json:"supports_rr_naptr,omitempty"`
	UseEnableMonitoring        bool       `json:"use_enable_monitoring,omitempty"`
	EnableMonitoring           bool       `json:"enable_monitoring,omitempty"`
	UseSynchronizationMinDelay bool       `json:"use_synchronization_min_delay,omitempty"`
	SynchronizationMinDelay    uint32     `json:"synchronization_min_delay,omitempty"`
	UseEnableDnsReportsSync    bool       `json:"use_enable_dns_reports_sync,omitempty"`
	EnableDnsReportsSync       bool       `json:"enable_dns_reports_sync,omitempty"`
}

// Dnstapsetting represents Infoblox struct dnstapsetting
type Dnstapsetting struct {
	DnstapReceiverAddress string `json:"dnstap_receiver_address,omitempty"`
	DnstapReceiverPort    uint32 `json:"dnstap_receiver_port,omitempty"`
	DnstapIdentity        string `json:"dnstap_identity,omitempty"`
	DnstapVersion         string `json:"dnstap_version,omitempty"`
}

// Downgrade represents Infoblox struct downgrade
type Downgrade struct{}

// Downloadcertificate represents Infoblox struct downloadcertificate
type Downloadcertificate struct{}

// Downloadpoolstatus represents Infoblox struct downloadpoolstatus
type Downloadpoolstatus struct{}

// Downloadthreatanalyticswhitelistparams represents Infoblox struct downloadthreatanalyticswhitelistparams
type Downloadthreatanalyticswhitelistparams struct{}

// DtcHealth represents Infoblox struct dtc:health
type DtcHealth struct {
	Availability string `json:"availability,omitempty"`
	EnabledState string `json:"enabled_state,omitempty"`
	Description  string `json:"description,omitempty"`
}

// DtcMonitorSnmpOid represents Infoblox struct dtc:monitor:snmp:oid
type DtcMonitorSnmpOid struct {
	Oid       string `json:"oid,omitempty"`
	Comment   string `json:"comment,omitempty"`
	Type      string `json:"type,omitempty"`
	Condition string `json:"condition,omitempty"`
	First     string `json:"first,omitempty"`
	Last      string `json:"last,omitempty"`
}

// DtcPoolConsolidatedMonitorHealth represents Infoblox struct dtc:pool:consolidated_monitor_health
type DtcPoolConsolidatedMonitorHealth struct {
	Members                 []string `json:"members,omitempty"`
	Monitor                 string   `json:"monitor,omitempty"`
	Availability            string   `json:"availability,omitempty"`
	FullHealthCommunication bool     `json:"full_health_communication,omitempty"`
}

// DtcPoolLink represents Infoblox struct dtc:pool:link
type DtcPoolLink struct {
	Pool  string `json:"pool,omitempty"`
	Ratio uint32 `json:"ratio,omitempty"`
}

// DtcServerLink represents Infoblox struct dtc:server:link
type DtcServerLink struct {
	Server string `json:"server,omitempty"`
	Ratio  uint32 `json:"ratio,omitempty"`
}

// DtcServerMonitor represents Infoblox struct dtc:server:monitor
type DtcServerMonitor struct {
	Monitor string `json:"monitor,omitempty"`
	Host    string `json:"host,omitempty"`
}

// DtcTopologyRuleSource represents Infoblox struct dtc:topology:rule:source
type DtcTopologyRuleSource struct {
	SourceType  string `json:"source_type,omitempty"`
	SourceOp    string `json:"source_op,omitempty"`
	SourceValue string `json:"source_value,omitempty"`
}

// DxlEndpointBroker represents Infoblox struct dxl:endpoint:broker
type DxlEndpointBroker struct {
	HostName string `json:"host_name,omitempty"`
	Address  string `json:"address,omitempty"`
	Port     uint32 `json:"port,omitempty"`
	UniqueId string `json:"unique_id,omitempty"`
}

// Eaexpressionop represents Infoblox struct eaexpressionop
type Eaexpressionop struct {
	Op      string `json:"op,omitempty"`
	Op1     string `json:"op1,omitempty"`
	Op1Type string `json:"op1_type,omitempty"`
	Op2     string `json:"op2,omitempty"`
	Op2Type string `json:"op2_type,omitempty"`
}

// Emptyparams represents Infoblox struct emptyparams
type Emptyparams struct{}

// Exclusionrange represents Infoblox struct exclusionrange
type Exclusionrange struct {
	StartAddress string `json:"start_address,omitempty"`
	EndAddress   string `json:"end_address,omitempty"`
	Comment      string `json:"comment,omitempty"`
}

// Exclusionrangetemplate represents Infoblox struct exclusionrangetemplate
type Exclusionrangetemplate struct {
	Offset            uint32 `json:"offset,omitempty"`
	NumberOfAddresses uint32 `json:"number_of_addresses,omitempty"`
	Comment           string `json:"comment,omitempty"`
}

// Expandnetwork represents Infoblox struct expandnetwork
type Expandnetwork struct{}

// Expressionop represents Infoblox struct expressionop
type Expressionop struct {
	Op      string `json:"op,omitempty"`
	Op1     string `json:"op1,omitempty"`
	Op1Type string `json:"op1_type,omitempty"`
	Op2     string `json:"op2,omitempty"`
	Op2Type string `json:"op2_type,omitempty"`
}

// ExtensibleattributedefDescendants represents Infoblox struct extensibleattributedef:descendants
type ExtensibleattributedefDescendants struct {
	OptionWithEa    string `json:"option_with_ea,omitempty"`
	OptionWithoutEa string `json:"option_without_ea,omitempty"`
	OptionDeleteEa  string `json:"option_delete_ea,omitempty"`
}

// EADefListValue represents Infoblox struct extensibleattributedef:listvalues
type EADefListValue struct {
	Value string `json:"value,omitempty"`
}

// NameServer represents Infoblox struct extserver
type NameServer struct {
	Address                      string `json:"address,omitempty"`
	Name                         string `json:"name,omitempty"`
	SharedWithMsParentDelegation bool   `json:"shared_with_ms_parent_delegation,omitempty"`
	Stealth                      bool   `json:"stealth,omitempty"`
	TsigKey                      string `json:"tsig_key,omitempty"`
	TsigKeyAlg                   string `json:"tsig_key_alg,omitempty"`
	TsigKeyName                  string `json:"tsig_key_name,omitempty"`
	UseTsigKeyName               bool   `json:"use_tsig_key_name,omitempty"`
}

// Extsyslogbackupserver represents Infoblox struct extsyslogbackupserver
type Extsyslogbackupserver struct {
	Address       string `json:"address,omitempty"`
	DirectoryPath string `json:"directory_path,omitempty"`
	Enable        bool   `json:"enable,omitempty"`
	Password      string `json:"password,omitempty"`
	Port          uint32 `json:"port,omitempty"`
	Protocol      string `json:"protocol,omitempty"`
	Username      string `json:"username,omitempty"`
}

// Filetransfersetting represents Infoblox struct filetransfersetting
type Filetransfersetting struct {
	Directory string `json:"directory,omitempty"`
	Host      string `json:"host,omitempty"`
	Password  string `json:"password,omitempty"`
	Type      string `json:"type,omitempty"`
	Username  string `json:"username,omitempty"`
	Port      uint32 `json:"port,omitempty"`
}

// Filterrule represents Infoblox struct filterrule
type Filterrule struct {
	Filter     string `json:"filter,omitempty"`
	Permission string `json:"permission,omitempty"`
}

// FireeyeAlertmap represents Infoblox struct fireeye:alertmap
type FireeyeAlertmap struct {
	AlertType string `json:"alert_type,omitempty"`
	RpzRule   string `json:"rpz_rule,omitempty"`
	Lifetime  uint32 `json:"lifetime,omitempty"`
}

// FireeyeRulemapping represents Infoblox struct fireeye:rulemapping
type FireeyeRulemapping struct {
	AptOverride           string             `json:"apt_override,omitempty"`
	FireeyeAlertMapping   []*FireeyeAlertmap `json:"fireeye_alert_mapping,omitempty"`
	SubstitutedDomainName string             `json:"substituted_domain_name,omitempty"`
}

// Forwardingmemberserver represents Infoblox struct forwardingmemberserver
type Forwardingmemberserver struct {
	Name                  string       `json:"name,omitempty"`
	ForwardersOnly        bool         `json:"forwarders_only,omitempty"`
	ForwardTo             []NameServer `json:"forward_to,omitempty"`
	UseOverrideForwarders bool         `json:"use_override_forwarders,omitempty"`
}

// Generatecsr represents Infoblox struct generatecsr
type Generatecsr struct{}

// Generatedxlendpointcerts represents Infoblox struct generatedxlendpointcerts
type Generatedxlendpointcerts struct{}

// Generatesafenetclientcert represents Infoblox struct generatesafenetclientcert
type Generatesafenetclientcert struct{}

// Generateselfsignedcert represents Infoblox struct generateselfsignedcert
type Generateselfsignedcert struct{}

// Generatetsigkeyparams represents Infoblox struct generatetsigkeyparams
type Generatetsigkeyparams struct{}

// Getdevicesupportinfo represents Infoblox struct getdevicesupportinfo
type Getdevicesupportinfo struct{}

// Getfileurl represents Infoblox struct getfileurl
type Getfileurl struct{}

// Getgriddata represents Infoblox struct getgriddata
type Getgriddata struct{}

// Getgridrevertstatusparams represents Infoblox struct getgridrevertstatusparams
type Getgridrevertstatusparams struct{}

// Getjobdevices represents Infoblox struct getjobdevices
type Getjobdevices struct{}

// Getjobprocessdetails represents Infoblox struct getjobprocessdetails
type Getjobprocessdetails struct{}

// Getlastuploadedruleset represents Infoblox struct getlastuploadedruleset
type Getlastuploadedruleset struct{}

// Getleasehistoryfiles represents Infoblox struct getleasehistoryfiles
type Getleasehistoryfiles struct{}

// Getlogfiles represents Infoblox struct getlogfiles
type Getlogfiles struct{}

// Getmemberdata represents Infoblox struct getmemberdata
type Getmemberdata struct{}

// Getsafenetclientcert represents Infoblox struct getsafenetclientcert
type Getsafenetclientcert struct{}

// Gettemplateschemaversionsparams represents Infoblox struct gettemplateschemaversionsparams
type Gettemplateschemaversionsparams struct{}

// Getvendoridentifiersparams represents Infoblox struct getvendoridentifiersparams
type Getvendoridentifiersparams struct{}

// GridAttackdetect represents Infoblox struct grid:attackdetect
type GridAttackdetect struct {
	Enable       bool   `json:"enable,omitempty"`
	High         uint32 `json:"high,omitempty"`
	IntervalMax  uint32 `json:"interval_max,omitempty"`
	IntervalMin  uint32 `json:"interval_min,omitempty"`
	IntervalTime uint32 `json:"interval_time,omitempty"`
	Low          uint32 `json:"low,omitempty"`
}

// GridAttackmitigation represents Infoblox struct grid:attackmitigation
type GridAttackmitigation struct {
	DetectChr               *GridAttackdetect `json:"detect_chr,omitempty"`
	DetectChrGrace          uint32            `json:"detect_chr_grace,omitempty"`
	DetectNxdomainResponses *GridAttackdetect `json:"detect_nxdomain_responses,omitempty"`
	DetectUdpDrop           *GridAttackdetect `json:"detect_udp_drop,omitempty"`
	Interval                uint32            `json:"interval,omitempty"`
	MitigateNxdomainLru     bool              `json:"mitigate_nxdomain_lru,omitempty"`
}

// GridAutoblackhole represents Infoblox struct grid:autoblackhole
type GridAutoblackhole struct {
	EnableFetchesPerServer bool   `json:"enable_fetches_per_server,omitempty"`
	EnableFetchesPerZone   bool   `json:"enable_fetches_per_zone,omitempty"`
	EnableHolddown         bool   `json:"enable_holddown,omitempty"`
	FetchesPerServer       uint32 `json:"fetches_per_server,omitempty"`
	FetchesPerZone         uint32 `json:"fetches_per_zone,omitempty"`
	FpsFreq                uint32 `json:"fps_freq,omitempty"`
	Holddown               uint32 `json:"holddown,omitempty"`
	HolddownThreshold      uint32 `json:"holddown_threshold,omitempty"`
	HolddownTimeout        uint32 `json:"holddown_timeout,omitempty"`
}

// GridCloudapiGatewayConfig represents Infoblox struct grid:cloudapi:gateway:config
type GridCloudapiGatewayConfig struct {
	EnableProxyService bool                                  `json:"enable_proxy_service,omitempty"`
	Port               uint32                                `json:"port,omitempty"`
	EndpointMapping    []*GridCloudapiGatewayEndpointmapping `json:"endpoint_mapping,omitempty"`
}

// GridCloudapiGatewayEndpointmapping represents Infoblox struct grid:cloudapi:gateway:endpointmapping
type GridCloudapiGatewayEndpointmapping struct {
	GatewayFqdn  string `json:"gateway_fqdn,omitempty"`
	EndpointFqdn string `json:"endpoint_fqdn,omitempty"`
}

// GridCloudapiInfo represents Infoblox struct grid:cloudapi:info
type GridCloudapiInfo struct {
	DelegatedMember *Dhcpmember `json:"delegated_member,omitempty"`
	DelegatedScope  string      `json:"delegated_scope,omitempty"`
	DelegatedRoot   string      `json:"delegated_root,omitempty"`
	OwnedByAdaptor  bool        `json:"owned_by_adaptor,omitempty"`
	Usage           string      `json:"usage,omitempty"`
	Tenant          string      `json:"tenant,omitempty"`
	MgmtPlatform    string      `json:"mgmt_platform,omitempty"`
	AuthorityType   string      `json:"authority_type,omitempty"`
}

// GridCloudapiUser represents Infoblox struct grid:cloudapi:user
type GridCloudapiUser struct {
	IsRemote    bool   `json:"is_remote,omitempty"`
	RemoteAdmin string `json:"remote_admin,omitempty"`
	LocalAdmin  string `json:"local_admin,omitempty"`
}

// GridConsentbannersetting represents Infoblox struct grid:consentbannersetting
type GridConsentbannersetting struct {
	Enable  bool   `json:"enable,omitempty"`
	Message string `json:"message,omitempty"`
}

// GridCspapiconfig represents Infoblox struct grid:cspapiconfig
type GridCspapiconfig struct {
	Url      string `json:"url,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// GridCspgridsetting represents Infoblox struct grid:cspgridsetting
type GridCspgridsetting struct {
	CspJoinToken   string `json:"csp_join_token,omitempty"`
	CspDnsResolver string `json:"csp_dns_resolver,omitempty"`
	CspHttpsProxy  string `json:"csp_https_proxy,omitempty"`
}

// GridDnsFixedrrsetorderfqdn represents Infoblox struct grid:dns:fixedrrsetorderfqdn
type GridDnsFixedrrsetorderfqdn struct {
	Fqdn       string `json:"fqdn,omitempty"`
	RecordType string `json:"record_type,omitempty"`
}

// GridInformationalbannersetting represents Infoblox struct grid:informationalbannersetting
type GridInformationalbannersetting struct {
	Enable  bool   `json:"enable,omitempty"`
	Message string `json:"message,omitempty"`
	Color   string `json:"color,omitempty"`
}

// GridLicensesubpool represents Infoblox struct grid:licensesubpool
type GridLicensesubpool struct {
	Key        string     `json:"key,omitempty"`
	Installed  uint32     `json:"installed,omitempty"`
	ExpiryDate *time.Time `json:"expiry_date,omitempty"`
}

// GridLockoutsetting represents Infoblox struct grid:lockoutsetting
type GridLockoutsetting struct {
	EnableSequentialFailedLoginAttemptsLockout bool   `json:"enable_sequential_failed_login_attempts_lockout,omitempty"`
	SequentialAttempts                         uint32 `json:"sequential_attempts,omitempty"`
	FailedLockoutDuration                      uint32 `json:"failed_lockout_duration,omitempty"`
	NeverUnlockUser                            bool   `json:"never_unlock_user,omitempty"`
}

// GridLoggingcategories represents Infoblox struct grid:loggingcategories
type GridLoggingcategories struct {
	LogDtcGslb        bool `json:"log_dtc_gslb,omitempty"`
	LogDtcHealth      bool `json:"log_dtc_health,omitempty"`
	LogGeneral        bool `json:"log_general,omitempty"`
	LogClient         bool `json:"log_client,omitempty"`
	LogConfig         bool `json:"log_config,omitempty"`
	LogDatabase       bool `json:"log_database,omitempty"`
	LogDnssec         bool `json:"log_dnssec,omitempty"`
	LogLameServers    bool `json:"log_lame_servers,omitempty"`
	LogNetwork        bool `json:"log_network,omitempty"`
	LogNotify         bool `json:"log_notify,omitempty"`
	LogQueries        bool `json:"log_queries,omitempty"`
	LogQueryRewrite   bool `json:"log_query_rewrite,omitempty"`
	LogResponses      bool `json:"log_responses,omitempty"`
	LogResolver       bool `json:"log_resolver,omitempty"`
	LogSecurity       bool `json:"log_security,omitempty"`
	LogUpdate         bool `json:"log_update,omitempty"`
	LogXferIn         bool `json:"log_xfer_in,omitempty"`
	LogXferOut        bool `json:"log_xfer_out,omitempty"`
	LogUpdateSecurity bool `json:"log_update_security,omitempty"`
	LogRateLimit      bool `json:"log_rate_limit,omitempty"`
	LogRpz            bool `json:"log_rpz,omitempty"`
}

// NTPSetting represents Infoblox struct grid:ntp
type NTPSetting struct {
	EnableNTP  bool         `json:"enable_ntp,omitempty"`
	NTPServers []*NTPserver `json:"ntp_servers,omitempty"`
	NTPKeys    []*Ntpkey    `json:"ntp_keys,omitempty"`
	NTPAcl     *Ntpaccess   `json:"ntp_acl,omitempty"`
	NTPKod     bool         `json:"ntp_kod,omitempty"`
}

// GridResponseratelimiting represents Infoblox struct grid:responseratelimiting
type GridResponseratelimiting struct {
	EnableRrl          bool   `json:"enable_rrl,omitempty"`
	LogOnly            bool   `json:"log_only,omitempty"`
	ResponsesPerSecond uint32 `json:"responses_per_second,omitempty"`
	Window             uint32 `json:"window,omitempty"`
	Slip               uint32 `json:"slip,omitempty"`
}

// GridRestartbannersetting represents Infoblox struct grid:restartbannersetting
type GridRestartbannersetting struct {
	Enabled                  bool `json:"enabled,omitempty"`
	EnableDoubleConfirmation bool `json:"enable_double_confirmation,omitempty"`
}

// GridServicerestart represents Infoblox struct grid:servicerestart
type GridServicerestart struct {
	Delay          uint32 `json:"delay,omitempty"`
	Timeout        int    `json:"timeout,omitempty"`
	RestartOffline bool   `json:"restart_offline,omitempty"`
}

// GridServicerestartGroupSchedule represents Infoblox struct grid:servicerestart:group:schedule
type GridServicerestartGroupSchedule struct {
	Services []string         `json:"services,omitempty"`
	Mode     string           `json:"mode,omitempty"`
	Schedule *SettingSchedule `json:"schedule,omitempty"`
	Force    bool             `json:"force,omitempty"`
}

// Gridjoin represents Infoblox struct gridjoin
type Gridjoin struct{}

// GridmemberSoamname represents Infoblox struct gridmember_soamname
type GridmemberSoamname struct {
	GridPrimary     string `json:"grid_primary,omitempty"`
	MsServerPrimary string `json:"ms_server_primary,omitempty"`
	Mname           string `json:"mname,omitempty"`
	DnsMname        string `json:"dns_mname,omitempty"`
}

// GridmemberSoaserial represents Infoblox struct gridmember_soaserial
type GridmemberSoaserial struct {
	GridPrimary     string `json:"grid_primary,omitempty"`
	MsServerPrimary string `json:"ms_server_primary,omitempty"`
	Serial          uint32 `json:"serial,omitempty"`
}

// Gridrestartservices represents Infoblox struct gridrestartservices
type Gridrestartservices struct{}

// Gridupgrade represents Infoblox struct gridupgrade
type Gridupgrade struct{}

// Hotfix represents Infoblox struct hotfix
type Hotfix struct {
	StatusText string `json:"status_text,omitempty"`
	UniqueId   string `json:"unique_id,omitempty"`
}

// HsmSafenet represents Infoblox struct hsm:safenet
type HsmSafenet struct {
	Name                  string `json:"name,omitempty"`
	PartitionSerialNumber string `json:"partition_serial_number,omitempty"`
	Disable               bool   `json:"disable,omitempty"`
	PartitionId           string `json:"partition_id,omitempty"`
	IsFipsCompliant       bool   `json:"is_fips_compliant,omitempty"`
	ServerCert            string `json:"server_cert,omitempty"`
	PartitionCapacity     uint32 `json:"partition_capacity,omitempty"`
	Status                string `json:"status,omitempty"`
}

// HsmThales represents Infoblox struct hsm:thales
type HsmThales struct {
	RemoteIp   string `json:"remote_ip,omitempty"`
	RemotePort uint32 `json:"remote_port,omitempty"`
	Status     string `json:"status,omitempty"`
	RemoteEsn  string `json:"remote_esn,omitempty"`
	Keyhash    string `json:"keyhash,omitempty"`
	Disable    bool   `json:"disable,omitempty"`
}

// Hsmrefreshparams represents Infoblox struct hsmrefreshparams
type Hsmrefreshparams struct{}

// Hsmteststatusparams represents Infoblox struct hsmteststatusparams
type Hsmteststatusparams struct{}

// Importdevicesupportbundle represents Infoblox struct importdevicesupportbundle
type Importdevicesupportbundle struct{}

// Importmaxminddb represents Infoblox struct importmaxminddb
type Importmaxminddb struct{}

// Interface represents Infoblox struct interface
type Interface struct {
	Anycast            bool            `json:"anycast,omitempty"`
	Ipv4NetworkSetting *SettingNetwork `json:"ipv4_network_setting,omitempty"`
	Ipv6NetworkSetting *Ipv6setting    `json:"ipv6_network_setting,omitempty"`
	Comment            string          `json:"comment,omitempty"`
	EnableBgp          bool            `json:"enable_bgp,omitempty"`
	EnableOspf         bool            `json:"enable_ospf,omitempty"`
	Interface          string          `json:"interface,omitempty"`
}

// Ipv6networksetting represents Infoblox struct ipv6networksetting
type Ipv6networksetting struct {
	Address string `json:"address,omitempty"`
	Cidr    uint32 `json:"cidr,omitempty"`
	Gateway string `json:"gateway,omitempty"`
}

// Ipv6setting represents Infoblox struct ipv6setting
type Ipv6setting struct {
	Enabled                 bool   `json:"enabled,omitempty"`
	VirtualIp               string `json:"virtual_ip,omitempty"`
	CidrPrefix              uint32 `json:"cidr_prefix,omitempty"`
	Gateway                 string `json:"gateway,omitempty"`
	AutoRouterConfigEnabled bool   `json:"auto_router_config_enabled,omitempty"`
	VlanId                  uint32 `json:"vlan_id,omitempty"`
	Primary                 bool   `json:"primary,omitempty"`
	Dscp                    uint32 `json:"dscp,omitempty"`
	UseDscp                 bool   `json:"use_dscp,omitempty"`
}

// Joinmgm represents Infoblox struct joinmgm
type Joinmgm struct{}

// Lan2portsetting represents Infoblox struct lan2portsetting
type Lan2portsetting struct {
	VirtualRouterId             uint32          `json:"virtual_router_id,omitempty"`
	Enabled                     bool            `json:"enabled,omitempty"`
	NetworkSetting              *SettingNetwork `json:"network_setting,omitempty"`
	V6NetworkSetting            *Ipv6setting    `json:"v6_network_setting,omitempty"`
	NicFailoverEnabled          bool            `json:"nic_failover_enabled,omitempty"`
	NicFailoverEnablePrimary    bool            `json:"nic_failover_enable_primary,omitempty"`
	DefaultRouteFailoverEnabled bool            `json:"default_route_failover_enabled,omitempty"`
}

// Lanhaportsetting represents Infoblox struct lanhaportsetting
type Lanhaportsetting struct {
	MgmtLan        string               `json:"mgmt_lan,omitempty"`
	MgmtIpv6addr   string               `json:"mgmt_ipv6addr,omitempty"`
	HaIpAddress    string               `json:"ha_ip_address,omitempty"`
	LanPortSetting *Physicalportsetting `json:"lan_port_setting,omitempty"`
	HaPortSetting  *Physicalportsetting `json:"ha_port_setting,omitempty"`
}

// LdapEamapping represents Infoblox struct ldap_eamapping
type LdapEamapping struct {
	Name     string `json:"name,omitempty"`
	MappedEa string `json:"mapped_ea,omitempty"`
}

// LdapServer represents Infoblox struct ldap_server
type LdapServer struct {
	Address            string `json:"address,omitempty"`
	AuthenticationType string `json:"authentication_type,omitempty"`
	BaseDn             string `json:"base_dn,omitempty"`
	BindPassword       string `json:"bind_password,omitempty"`
	BindUserDn         string `json:"bind_user_dn,omitempty"`
	Comment            string `json:"comment,omitempty"`
	Disable            bool   `json:"disable,omitempty"`
	Encryption         string `json:"encryption,omitempty"`
	Port               uint32 `json:"port,omitempty"`
	UseMgmtPort        bool   `json:"use_mgmt_port,omitempty"`
	Version            string `json:"version,omitempty"`
}

// Lockunlockzone represents Infoblox struct lockunlockzone
type Lockunlockzone struct{}

// Logicfilterrule represents Infoblox struct logicfilterrule
type Logicfilterrule struct {
	Filter string `json:"filter,omitempty"`
	Type   string `json:"type,omitempty"`
}

// Lomnetworkconfig represents Infoblox struct lomnetworkconfig
type Lomnetworkconfig struct {
	Address      string `json:"address,omitempty"`
	Gateway      string `json:"gateway,omitempty"`
	SubnetMask   string `json:"subnet_mask,omitempty"`
	IsLomCapable bool   `json:"is_lom_capable,omitempty"`
}

// Lomuser represents Infoblox struct lomuser
type Lomuser struct {
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role,omitempty"`
	Disable  bool   `json:"disable,omitempty"`
	Comment  string `json:"comment,omitempty"`
}

// MemberCspmembersetting represents Infoblox struct member:cspmembersetting
type MemberCspmembersetting struct {
	UseCspJoinToken   bool   `json:"use_csp_join_token,omitempty"`
	UseCspDnsResolver bool   `json:"use_csp_dns_resolver,omitempty"`
	UseCspHttpsProxy  bool   `json:"use_csp_https_proxy,omitempty"`
	CspJoinToken      string `json:"csp_join_token,omitempty"`
	CspDnsResolver    string `json:"csp_dns_resolver,omitempty"`
	CspHttpsProxy     string `json:"csp_https_proxy,omitempty"`
}

// MemberDnsgluerecordaddr represents Infoblox struct member:dnsgluerecordaddr
type MemberDnsgluerecordaddr struct {
	AttachEmptyRecursiveView bool   `json:"attach_empty_recursive_view,omitempty"`
	GlueRecordAddress        string `json:"glue_record_address,omitempty"`
	View                     string `json:"view,omitempty"`
	GlueAddressChoice        string `json:"glue_address_choice,omitempty"`
}

// MemberDnsip represents Infoblox struct member:dnsip
type MemberDnsip struct {
	IpAddress string `json:"ip_address,omitempty"`
	Ipsd      string `json:"ipsd,omitempty"`
}

// MemberNtp represents Infoblox struct member:ntp
type MemberNtp struct {
	EnableNTP                  bool         `json:"enable_ntp,omitempty"`
	NTPServers                 []*NTPserver `json:"ntp_servers,omitempty"`
	NTPKeys                    []*Ntpkey    `json:"ntp_keys,omitempty"`
	NTPAcl                     *Ntpaccess   `json:"ntp_acl,omitempty"`
	NTPKod                     bool         `json:"ntp_kod,omitempty"`
	EnableExternalNtpServers   bool         `json:"enable_external_ntp_servers,omitempty"`
	ExcludeGridMasterNtpServer bool         `json:"exclude_grid_master_ntp_server,omitempty"`
	UseNtpServers              bool         `json:"use_ntp_servers,omitempty"`
	UseNtpKeys                 bool         `json:"use_ntp_keys,omitempty"`
	UseNtpAcl                  bool         `json:"use_ntp_acl,omitempty"`
	UseNtpKod                  bool         `json:"use_ntp_kod,omitempty"`
}

// Memberadminoperationparams represents Infoblox struct memberadminoperationparams
type Memberadminoperationparams struct{}

// Membercapturecontrolparams represents Infoblox struct membercapturecontrolparams
type Membercapturecontrolparams struct{}

// Membercapturestatusparams represents Infoblox struct membercapturestatusparams
type Membercapturestatusparams struct{}

// Memberrestartservices represents Infoblox struct memberrestartservices
type Memberrestartservices struct{}

// Memberserver represents Infoblox struct memberserver
type Memberserver struct {
	Name                     string       `json:"name,omitempty"`
	Stealth                  bool         `json:"stealth,omitempty"`
	GridReplicate            bool         `json:"grid_replicate,omitempty"`
	Lead                     bool         `json:"lead,omitempty"`
	PreferredPrimaries       []NameServer `json:"preferred_primaries,omitempty"`
	EnablePreferredPrimaries bool         `json:"enable_preferred_primaries,omitempty"`
}

// Memberservicecommunication represents Infoblox struct memberservicecommunication
type Memberservicecommunication struct {
	Service string `json:"service,omitempty"`
	Type    string `json:"type,omitempty"`
	Option  string `json:"option,omitempty"`
}

// Memberservicestatus represents Infoblox struct memberservicestatus
type Memberservicestatus struct {
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
	Service     string `json:"service,omitempty"`
}

// Memberupgrade represents Infoblox struct memberupgrade
type Memberupgrade struct{}

// Mgmtportsetting represents Infoblox struct mgmtportsetting
type Mgmtportsetting struct {
	Enabled               bool `json:"enabled,omitempty"`
	VpnEnabled            bool `json:"vpn_enabled,omitempty"`
	SecurityAccessEnabled bool `json:"security_access_enabled,omitempty"`
}

// Modifysdnnetworkassignment represents Infoblox struct modifysdnnetworkassignment
type Modifysdnnetworkassignment struct{}

// Modifyvrfassignment represents Infoblox struct modifyvrfassignment
type Modifyvrfassignment struct{}

// Monitoreddomains represents Infoblox struct monitoreddomains
type Monitoreddomains struct {
	DomainName string `json:"domain_name,omitempty"`
	RecordType string `json:"record_type,omitempty"`
}

// Moveblacklistrpztowhitelistparams represents Infoblox struct moveblacklistrpztowhitelistparams
type Moveblacklistrpztowhitelistparams struct{}

// Movesubnets represents Infoblox struct movesubnets
type Movesubnets struct{}

// Msdhcpoption represents Infoblox struct msdhcpoption
type Msdhcpoption struct {
	Num         uint32 `json:"num,omitempty"`
	Value       string `json:"value,omitempty"`
	Name        string `json:"name,omitempty"`
	VendorClass string `json:"vendor_class,omitempty"`
	UserClass   string `json:"user_class,omitempty"`
	Type        string `json:"type,omitempty"`
}

// Msdhcpserver represents Infoblox struct msdhcpserver
type Msdhcpserver struct {
	Ipv4Addr string `json:"ipv4addr,omitempty"`
}

// Msdnsserver represents Infoblox struct msdnsserver
type Msdnsserver struct {
	Address                      string `json:"address,omitempty"`
	IsMaster                     bool   `json:"is_master,omitempty"`
	NsIp                         string `json:"ns_ip,omitempty"`
	NsName                       string `json:"ns_name,omitempty"`
	Stealth                      bool   `json:"stealth,omitempty"`
	SharedWithMsParentDelegation bool   `json:"shared_with_ms_parent_delegation,omitempty"`
}

// MsserverAduser represents Infoblox struct msserver:aduser
type MsserverAduser struct {
	LoginName                  string     `json:"login_name,omitempty"`
	LoginPassword              string     `json:"login_password,omitempty"`
	EnableUserSync             bool       `json:"enable_user_sync,omitempty"`
	SynchronizationInterval    uint32     `json:"synchronization_interval,omitempty"`
	LastSyncTime               *time.Time `json:"last_sync_time,omitempty"`
	LastSyncStatus             string     `json:"last_sync_status,omitempty"`
	LastSyncDetail             string     `json:"last_sync_detail,omitempty"`
	LastSuccessSyncTime        *time.Time `json:"last_success_sync_time,omitempty"`
	UseLogin                   bool       `json:"use_login,omitempty"`
	UseEnableAdUserSync        bool       `json:"use_enable_ad_user_sync,omitempty"`
	UseSynchronizationMinDelay bool       `json:"use_synchronization_min_delay,omitempty"`
	UseEnableUserSync          bool       `json:"use_enable_user_sync,omitempty"`
	UseSynchronizationInterval bool       `json:"use_synchronization_interval,omitempty"`
}

// MsserverAduserData represents Infoblox struct msserver:aduser:data
type MsserverAduserData struct {
	ActiveUsersCount uint32 `json:"active_users_count,omitempty"`
}

// MsserverDcnsrecordcreation represents Infoblox struct msserver:dcnsrecordcreation
type MsserverDcnsrecordcreation struct {
	Address string `json:"address,omitempty"`
	Comment string `json:"comment,omitempty"`
}

// Natsetting represents Infoblox struct natsetting
type Natsetting struct {
	Enabled           bool   `json:"enabled,omitempty"`
	ExternalVirtualIp string `json:"external_virtual_ip,omitempty"`
	Group             string `json:"group,omitempty"`
}

// Networkdiscoverycontrolparams represents Infoblox struct networkdiscoverycontrolparams
type Networkdiscoverycontrolparams struct{}

// NetworkviewAssocmember represents Infoblox struct networkview:assocmember
type NetworkviewAssocmember struct {
	Member    string   `json:"member,omitempty"`
	Failovers []string `json:"failovers,omitempty"`
}

// Nextavailableip represents Infoblox struct nextavailableip
type Nextavailableip struct{}

// Nextavailableip6 represents Infoblox struct nextavailableip6
type Nextavailableip6 struct{}

// Nextavailablenet represents Infoblox struct nextavailablenet
type Nextavailablenet struct{}

// Nextavailablenet6 represents Infoblox struct nextavailablenet6
type Nextavailablenet6 struct{}

// Nextavailablevlan represents Infoblox struct nextavailablevlan
type Nextavailablevlan struct{}

// Nextavailablevlanid represents Infoblox struct nextavailablevlanid
type Nextavailablevlanid struct{}

// Nodeinfo represents Infoblox struct nodeinfo
type Nodeinfo struct {
	ServiceStatus        []*Servicestatus     `json:"service_status,omitempty"`
	PhysicalOid          string               `json:"physical_oid,omitempty"`
	HaStatus             string               `json:"ha_status,omitempty"`
	Hwplatform           string               `json:"hwplatform,omitempty"`
	Hwid                 string               `json:"hwid,omitempty"`
	Hwmodel              string               `json:"hwmodel,omitempty"`
	Hwtype               string               `json:"hwtype,omitempty"`
	PaidNios             bool                 `json:"paid_nios,omitempty"`
	MgmtNetworkSetting   *SettingNetwork      `json:"mgmt_network_setting,omitempty"`
	LanHaPortSetting     *Lanhaportsetting    `json:"lan_ha_port_setting,omitempty"`
	MgmtPhysicalSetting  *Physicalportsetting `json:"mgmt_physical_setting,omitempty"`
	Lan2PhysicalSetting  *Physicalportsetting `json:"lan2_physical_setting,omitempty"`
	NatExternalIp        string               `json:"nat_external_ip,omitempty"`
	V6MgmtNetworkSetting *Ipv6setting         `json:"v6_mgmt_network_setting,omitempty"`
}

// NotificationRestTemplateinstance represents Infoblox struct notification:rest:templateinstance
type NotificationRestTemplateinstance struct {
	Template   string                               `json:"template,omitempty"`
	Parameters []*NotificationRestTemplateparameter `json:"parameters,omitempty"`
}

// NotificationRestTemplateparameter represents Infoblox struct notification:rest:templateparameter
type NotificationRestTemplateparameter struct {
	Name         string `json:"name,omitempty"`
	Value        string `json:"value,omitempty"`
	DefaultValue string `json:"default_value,omitempty"`
	Syntax       string `json:"syntax,omitempty"`
}

// NotificationRuleexpressionop represents Infoblox struct notification:ruleexpressionop
type NotificationRuleexpressionop struct {
	Op      string `json:"op,omitempty"`
	Op1     string `json:"op1,omitempty"`
	Op1Type string `json:"op1_type,omitempty"`
	Op2     string `json:"op2,omitempty"`
	Op2Type string `json:"op2_type,omitempty"`
}

// Ntpac represents Infoblox struct ntpac
type Ntpac struct {
	AddressAc *Addressac `json:"address_ac,omitempty"`
	Service   string     `json:"service,omitempty"`
}

// Ntpaccess represents Infoblox struct ntpaccess
type Ntpaccess struct {
	AclType  string   `json:"acl_type,omitempty"`
	AcList   []*Ntpac `json:"ac_list,omitempty"`
	NamedAcl string   `json:"named_acl,omitempty"`
	Service  string   `json:"service,omitempty"`
}

// Ntpkey represents Infoblox struct ntpkey
type Ntpkey struct {
	Number uint32 `json:"number,omitempty"`
	String string `json:"string,omitempty"`
	Type   string `json:"type,omitempty"`
}

// NTPserver represents Infoblox struct ntpserver
type NTPserver struct {
	Address              string `json:"address,omitempty"`
	EnableAuthentication bool   `json:"enable_authentication,omitempty"`
	NtpKeyNumber         uint32 `json:"ntp_key_number,omitempty"`
	Preferred            bool   `json:"preferred,omitempty"`
	Burst                bool   `json:"burst,omitempty"`
	IBurst               bool   `json:"iburst,omitempty"`
}

// Nxdomainrule represents Infoblox struct nxdomainrule
type Nxdomainrule struct {
	Action  string `json:"action,omitempty"`
	Pattern string `json:"pattern,omitempty"`
}

// Objectschangestrackingsetting represents Infoblox struct objectschangestrackingsetting
type Objectschangestrackingsetting struct {
	Enable           bool   `json:"enable,omitempty"`
	EnableCompletion uint32 `json:"enable_completion,omitempty"`
	State            string `json:"state,omitempty"`
	MaxTimeToTrack   uint32 `json:"max_time_to_track,omitempty"`
	MaxObjsToTrack   uint32 `json:"max_objs_to_track,omitempty"`
}

// OcspResponder represents Infoblox struct ocsp_responder
type OcspResponder struct {
	FqdnOrIp         string `json:"fqdn_or_ip,omitempty"`
	Port             uint32 `json:"port,omitempty"`
	Comment          string `json:"comment,omitempty"`
	Disabled         bool   `json:"disabled,omitempty"`
	Certificate      string `json:"certificate,omitempty"`
	CertificateToken string `json:"certificate_token,omitempty"`
}

// Option60matchrule represents Infoblox struct option60matchrule
type Option60matchrule struct {
	MatchValue      string `json:"match_value,omitempty"`
	OptionSpace     string `json:"option_space,omitempty"`
	IsSubstring     bool   `json:"is_substring,omitempty"`
	SubstringOffset uint32 `json:"substring_offset,omitempty"`
	SubstringLength uint32 `json:"substring_length,omitempty"`
}

// Ospf represents Infoblox struct ospf
type Ospf struct {
	AreaId                 string `json:"area_id,omitempty"`
	AreaType               string `json:"area_type,omitempty"`
	AuthenticationKey      string `json:"authentication_key,omitempty"`
	AuthenticationType     string `json:"authentication_type,omitempty"`
	AutoCalcCostEnabled    bool   `json:"auto_calc_cost_enabled,omitempty"`
	Comment                string `json:"comment,omitempty"`
	Cost                   uint32 `json:"cost,omitempty"`
	DeadInterval           uint32 `json:"dead_interval,omitempty"`
	HelloInterval          uint32 `json:"hello_interval,omitempty"`
	Interface              string `json:"interface,omitempty"`
	IsIpv4                 bool   `json:"is_ipv4,omitempty"`
	KeyId                  uint32 `json:"key_id,omitempty"`
	RetransmitInterval     uint32 `json:"retransmit_interval,omitempty"`
	TransmitDelay          uint32 `json:"transmit_delay,omitempty"`
	AdvertiseInterfaceVlan string `json:"advertise_interface_vlan,omitempty"`
	BfdTemplate            string `json:"bfd_template,omitempty"`
	EnableBfd              bool   `json:"enable_bfd,omitempty"`
}

// OutboundCloudclientEvent represents Infoblox struct outbound:cloudclient:event
type OutboundCloudclientEvent struct {
	EventType string `json:"event_type,omitempty"`
	Enabled   bool   `json:"enabled,omitempty"`
}

// ParentalcontrolAbs represents Infoblox struct parentalcontrol:abs
type ParentalcontrolAbs struct {
	IpAddress      string `json:"ip_address,omitempty"`
	BlockingPolicy string `json:"blocking_policy,omitempty"`
}

// ParentalcontrolMsp represents Infoblox struct parentalcontrol:msp
type ParentalcontrolMsp struct {
	IpAddress string `json:"ip_address,omitempty"`
}

// ParentalcontrolNasgateway represents Infoblox struct parentalcontrol:nasgateway
type ParentalcontrolNasgateway struct {
	Name         string `json:"name,omitempty"`
	IpAddress    string `json:"ip_address,omitempty"`
	SharedSecret string `json:"shared_secret,omitempty"`
	SendAck      bool   `json:"send_ack,omitempty"`
	MessageRate  uint32 `json:"message_rate,omitempty"`
	Comment      string `json:"comment,omitempty"`
}

// ParentalcontrolSitemember represents Infoblox struct parentalcontrol:sitemember
type ParentalcontrolSitemember struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

// ParentalcontrolSpm represents Infoblox struct parentalcontrol:spm
type ParentalcontrolSpm struct {
	IpAddress string `json:"ip_address,omitempty"`
}

// Parentcheck represents Infoblox struct parentcheck
type Parentcheck struct{}

// Physicalportsetting represents Infoblox struct physicalportsetting
type Physicalportsetting struct {
	AutoPortSettingEnabled bool   `json:"auto_port_setting_enabled,omitempty"`
	Speed                  string `json:"speed,omitempty"`
	Duplex                 string `json:"duplex,omitempty"`
}

// Pnodetokenoperation represents Infoblox struct pnodetokenoperation
type Pnodetokenoperation struct{}

// Preprovision represents Infoblox struct preprovision
type Preprovision struct {
	HardwareInfo []*Preprovisionhardware `json:"hardware_info,omitempty"`
	Licenses     []string                `json:"licenses,omitempty"`
}

// Preprovisionhardware represents Infoblox struct preprovisionhardware
type Preprovisionhardware struct {
	Hwtype  string `json:"hwtype,omitempty"`
	Hwmodel string `json:"hwmodel,omitempty"`
}

// PropertiesBlackoutsetting represents Infoblox struct properties:blackoutsetting
type PropertiesBlackoutsetting struct {
	EnableBlackout   bool             `json:"enable_blackout,omitempty"`
	BlackoutDuration uint32           `json:"blackout_duration,omitempty"`
	BlackoutSchedule *SettingSchedule `json:"blackout_schedule,omitempty"`
}

// Provisionnetworkdhcprelay represents Infoblox struct provisionnetworkdhcprelay
type Provisionnetworkdhcprelay struct{}

// Provisionnetworkport represents Infoblox struct provisionnetworkport
type Provisionnetworkport struct{}

// Publishchanges represents Infoblox struct publishchanges
type Publishchanges struct{}

// Queriesuser represents Infoblox struct queriesuser
type Queriesuser struct {
	User    string `json:"user,omitempty"`
	Comment string `json:"comment,omitempty"`
}

// Query represents Infoblox struct query
type Query struct{}

// Queryfqdnonmemberparams represents Infoblox struct queryfqdnonmemberparams
type Queryfqdnonmemberparams struct{}

// RadiusServer represents Infoblox struct radius:server
type RadiusServer struct {
	AcctPort      uint32 `json:"acct_port,omitempty"`
	AuthPort      uint32 `json:"auth_port,omitempty"`
	AuthType      string `json:"auth_type,omitempty"`
	Comment       string `json:"comment,omitempty"`
	Disable       bool   `json:"disable,omitempty"`
	Address       string `json:"address,omitempty"`
	SharedSecret  string `json:"shared_secret,omitempty"`
	UseAccounting bool   `json:"use_accounting,omitempty"`
	UseMgmtPort   bool   `json:"use_mgmt_port,omitempty"`
}

// Rdatasubfield represents Infoblox struct rdatasubfield
type Rdatasubfield struct {
	FieldValue    string `json:"field_value,omitempty"`
	FieldType     string `json:"field_type,omitempty"`
	IncludeLength string `json:"include_length,omitempty"`
}

// Read represents Infoblox struct read
type Read struct{}

// Remoteddnszone represents Infoblox struct remoteddnszone
type Remoteddnszone struct {
	Fqdn                string `json:"fqdn,omitempty"`
	ServerAddress       string `json:"server_address,omitempty"`
	GssTsigDnsPrincipal string `json:"gss_tsig_dns_principal,omitempty"`
	GssTsigDomain       string `json:"gss_tsig_domain,omitempty"`
	TsigKey             string `json:"tsig_key,omitempty"`
	TsigKeyAlg          string `json:"tsig_key_alg,omitempty"`
	TsigKeyName         string `json:"tsig_key_name,omitempty"`
	KeyType             string `json:"key_type,omitempty"`
}

// Requestgridservicestatus represents Infoblox struct requestgridservicestatus
type Requestgridservicestatus struct{}

// Requestmemberservicestatus represents Infoblox struct requestmemberservicestatus
type Requestmemberservicestatus struct{}

// Resizenetwork represents Infoblox struct resizenetwork
type Resizenetwork struct{}

// Restapitemplateexportparams represents Infoblox struct restapitemplateexportparams
type Restapitemplateexportparams struct{}

// Restapitemplateexportschemaparams represents Infoblox struct restapitemplateexportschemaparams
type Restapitemplateexportschemaparams struct{}

// Restapitemplateimportparams represents Infoblox struct restapitemplateimportparams
type Restapitemplateimportparams struct{}

// Restoredatabase represents Infoblox struct restoredatabase
type Restoredatabase struct{}

// Restoredtcconfig represents Infoblox struct restoredtcconfig
type Restoredtcconfig struct{}

// Runnow represents Infoblox struct runnow
type Runnow struct{}

// Runscavenging represents Infoblox struct runscavenging
type Runscavenging struct{}

// SamlIdp represents Infoblox struct saml:idp
type SamlIdp struct {
	IdpType        string `json:"idp_type,omitempty"`
	Comment        string `json:"comment,omitempty"`
	MetadataUrl    string `json:"metadata_url,omitempty"`
	MetadataToken  string `json:"metadata_token,omitempty"`
	Groupname      string `json:"groupname,omitempty"`
	SsoRedirectUrl string `json:"sso_redirect_url,omitempty"`
}

// Savedbsnapshot represents Infoblox struct savedbsnapshot
type Savedbsnapshot struct{}

// Scheduledbackup represents Infoblox struct scheduledbackup
type Scheduledbackup struct {
	Status          string `json:"status,omitempty"`
	Execute         string `json:"execute,omitempty"`
	Operation       string `json:"operation,omitempty"`
	BackupType      string `json:"backup_type,omitempty"`
	KeepLocalCopy   bool   `json:"keep_local_copy,omitempty"`
	BackupFrequency string `json:"backup_frequency,omitempty"`
	Weekday         string `json:"weekday,omitempty"`
	HourOfDay       uint32 `json:"hour_of_day,omitempty"`
	MinutesPastHour uint32 `json:"minutes_past_hour,omitempty"`
	Username        string `json:"username,omitempty"`
	Password        string `json:"password,omitempty"`
	BackupServer    string `json:"backup_server,omitempty"`
	Path            string `json:"path,omitempty"`
	RestoreType     string `json:"restore_type,omitempty"`
	RestoreServer   string `json:"restore_server,omitempty"`
	RestoreUsername string `json:"restore_username,omitempty"`
	RestorePassword string `json:"restore_password,omitempty"`
	RestorePath     string `json:"restore_path,omitempty"`
	NiosData        bool   `json:"nios_data,omitempty"`
	DiscoveryData   bool   `json:"discovery_data,omitempty"`
	SplunkAppData   bool   `json:"splunk_app_data,omitempty"`
	Enable          bool   `json:"enable,omitempty"`
	UseKeys         bool   `json:"use_keys,omitempty"`
	KeyType         string `json:"key_type,omitempty"`
	UploadKeys      bool   `json:"upload_keys,omitempty"`
	DownloadKeys    bool   `json:"download_keys,omitempty"`
}

// Servicestatus represents Infoblox struct servicestatus
type Servicestatus struct {
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
	Service     string `json:"service,omitempty"`
}

// Setcaptiveportalfile represents Infoblox struct setcaptiveportalfile
type Setcaptiveportalfile struct{}

// Setdatafiledest represents Infoblox struct setdatafiledest
type Setdatafiledest struct{}

// Setdhcpfailoverpartnerdown represents Infoblox struct setdhcpfailoverpartnerdown
type Setdhcpfailoverpartnerdown struct{}

// Setdhcpleases represents Infoblox struct setdhcpleases
type Setdhcpleases struct{}

// Setdiscoverycsv represents Infoblox struct setdiscoverycsv
type Setdiscoverycsv struct{}

// Setlastuploadedmodulesetparams represents Infoblox struct setlastuploadedmodulesetparams
type Setlastuploadedmodulesetparams struct{}

// Setlastuploadedruleset represents Infoblox struct setlastuploadedruleset
type Setlastuploadedruleset struct{}

// Setleasehistoryfiles represents Infoblox struct setleasehistoryfiles
type Setleasehistoryfiles struct{}

// Setmemberdata represents Infoblox struct setmemberdata
type Setmemberdata struct{}

// SettingAtpoutbound represents Infoblox struct setting:atpoutbound
type SettingAtpoutbound struct {
	EnableQueryFqdn bool   `json:"enable_query_fqdn,omitempty"`
	QueryFqdnLimit  uint32 `json:"query_fqdn_limit,omitempty"`
}

// SettingAutomatedtrafficcapture represents Infoblox struct setting:automatedtrafficcapture
type SettingAutomatedtrafficcapture struct {
	TrafficCaptureEnable    bool   `json:"traffic_capture_enable,omitempty"`
	Destination             string `json:"destination,omitempty"`
	Duration                uint32 `json:"duration,omitempty"`
	IncludeSupportBundle    bool   `json:"include_support_bundle,omitempty"`
	KeepLocalCopy           bool   `json:"keep_local_copy,omitempty"`
	DestinationHost         string `json:"destination_host,omitempty"`
	TrafficCaptureDirectory string `json:"traffic_capture_directory,omitempty"`
	SupportBundleDirectory  string `json:"support_bundle_directory,omitempty"`
	Username                string `json:"username,omitempty"`
	Password                string `json:"password,omitempty"`
}

// SettingDnsresolver represents Infoblox struct setting:dnsresolver
type SettingDnsresolver struct {
	Resolvers     []string `json:"resolvers,omitempty"`
	SearchDomains []string `json:"search_domains,omitempty"`
}

// SettingDynamicratio represents Infoblox struct setting:dynamicratio
type SettingDynamicratio struct {
	Method              string `json:"method,omitempty"`
	Monitor             string `json:"monitor,omitempty"`
	MonitorMetric       string `json:"monitor_metric,omitempty"`
	MonitorWeighing     string `json:"monitor_weighing,omitempty"`
	InvertMonitorMetric bool   `json:"invert_monitor_metric,omitempty"`
}

// SettingEmail represents Infoblox struct setting:email
type SettingEmail struct {
	Enabled           bool   `json:"enabled,omitempty"`
	FromAddress       string `json:"from_address,omitempty"`
	Address           string `json:"address,omitempty"`
	RelayEnabled      bool   `json:"relay_enabled,omitempty"`
	Relay             string `json:"relay,omitempty"`
	Password          string `json:"password,omitempty"`
	Smtps             bool   `json:"smtps,omitempty"`
	PortNumber        uint32 `json:"port_number,omitempty"`
	UseAuthentication bool   `json:"use_authentication,omitempty"`
}

// SettingHttpproxyserver represents Infoblox struct setting:httpproxyserver
type SettingHttpproxyserver struct {
	Address                   string `json:"address,omitempty"`
	Port                      uint32 `json:"port,omitempty"`
	EnableProxy               bool   `json:"enable_proxy,omitempty"`
	EnableContentInspection   bool   `json:"enable_content_inspection,omitempty"`
	VerifyCname               bool   `json:"verify_cname,omitempty"`
	Comment                   string `json:"comment,omitempty"`
	Username                  string `json:"username,omitempty"`
	Password                  string `json:"password,omitempty"`
	Certificate               string `json:"certificate,omitempty"`
	EnableUsernameAndPassword bool   `json:"enable_username_and_password,omitempty"`
}

// SettingInactivelockout represents Infoblox struct setting:inactivelockout
type SettingInactivelockout struct {
	AccountInactivityLockoutEnable   bool   `json:"account_inactivity_lockout_enable,omitempty"`
	InactiveDays                     uint32 `json:"inactive_days,omitempty"`
	ReminderDays                     uint32 `json:"reminder_days,omitempty"`
	ReactivateViaSerialConsoleEnable bool   `json:"reactivate_via_serial_console_enable,omitempty"`
	ReactivateViaRemoteConsoleEnable bool   `json:"reactivate_via_remote_console_enable,omitempty"`
}

// SettingIpamThreshold represents Infoblox struct setting:ipam:threshold
type SettingIpamThreshold struct {
	TriggerValue uint32 `json:"trigger_value,omitempty"`
	ResetValue   uint32 `json:"reset_value,omitempty"`
}

// SettingIpamTrap represents Infoblox struct setting:ipam:trap
type SettingIpamTrap struct {
	EnableEmailWarnings bool `json:"enable_email_warnings,omitempty"`
	EnableSnmpWarnings  bool `json:"enable_snmp_warnings,omitempty"`
}

// SettingMsserver represents Infoblox struct setting:msserver
type SettingMsserver struct {
	LogDestination       string `json:"log_destination,omitempty"`
	EnableInvalidMac     bool   `json:"enable_invalid_mac,omitempty"`
	MaxConnection        uint32 `json:"max_connection,omitempty"`
	RpcTimeout           uint32 `json:"rpc_timeout,omitempty"`
	EnableDhcpMonitoring bool   `json:"enable_dhcp_monitoring,omitempty"`
	EnableDnsMonitoring  bool   `json:"enable_dns_monitoring,omitempty"`
	LdapTimeout          uint32 `json:"ldap_timeout,omitempty"`
	DefaultIpSiteLink    string `json:"default_ip_site_link,omitempty"`
	EnableNetworkUsers   bool   `json:"enable_network_users,omitempty"`
	EnableAdUserSync     bool   `json:"enable_ad_user_sync,omitempty"`
	AdUserDefaultTimeout uint32 `json:"ad_user_default_timeout,omitempty"`
	EnableDnsReportsSync bool   `json:"enable_dns_reports_sync,omitempty"`
}

// SettingNetwork represents Infoblox struct setting:network
type SettingNetwork struct {
	Address    string `json:"address,omitempty"`
	Gateway    string `json:"gateway,omitempty"`
	SubnetMask string `json:"subnet_mask,omitempty"`
	VlanId     uint32 `json:"vlan_id,omitempty"`
	Primary    bool   `json:"primary,omitempty"`
	Dscp       uint32 `json:"dscp,omitempty"`
	UseDscp    bool   `json:"use_dscp,omitempty"`
}

// SettingPassword represents Infoblox struct setting:password
type SettingPassword struct {
	PasswordMinLength uint32 `json:"password_min_length,omitempty"`
	NumLowerChar      uint32 `json:"num_lower_char,omitempty"`
	NumUpperChar      uint32 `json:"num_upper_char,omitempty"`
	NumNumericChar    uint32 `json:"num_numeric_char,omitempty"`
	NumSymbolChar     uint32 `json:"num_symbol_char,omitempty"`
	CharsToChange     uint32 `json:"chars_to_change,omitempty"`
	ExpireDays        uint32 `json:"expire_days,omitempty"`
	ReminderDays      uint32 `json:"reminder_days,omitempty"`
	ForceResetEnable  bool   `json:"force_reset_enable,omitempty"`
	ExpireEnable      bool   `json:"expire_enable,omitempty"`
	HistoryEnable     bool   `json:"history_enable,omitempty"`
	NumPasswordsSaved uint32 `json:"num_passwords_saved,omitempty"`
	MinPasswordAge    uint32 `json:"min_password_age,omitempty"`
}

// SettingScavenging represents Infoblox struct setting:scavenging
type SettingScavenging struct {
	EnableScavenging          bool              `json:"enable_scavenging,omitempty"`
	EnableRecurrentScavenging bool              `json:"enable_recurrent_scavenging,omitempty"`
	EnableAutoReclamation     bool              `json:"enable_auto_reclamation,omitempty"`
	EnableRrLastQueried       bool              `json:"enable_rr_last_queried,omitempty"`
	EnableZoneLastQueried     bool              `json:"enable_zone_last_queried,omitempty"`
	ReclaimAssociatedRecords  bool              `json:"reclaim_associated_records,omitempty"`
	ScavengingSchedule        *SettingSchedule  `json:"scavenging_schedule,omitempty"`
	ExpressionList            []*Expressionop   `json:"expression_list,omitempty"`
	EaExpressionList          []*Eaexpressionop `json:"ea_expression_list,omitempty"`
}

// SettingSchedule represents Infoblox struct setting:schedule
type SettingSchedule struct {
	Weekdays        []string   `json:"weekdays,omitempty"`
	TimeZone        string     `json:"time_zone,omitempty"`
	RecurringTime   *time.Time `json:"recurring_time,omitempty"`
	Frequency       string     `json:"frequency,omitempty"`
	Every           uint32     `json:"every,omitempty"`
	MinutesPastHour uint32     `json:"minutes_past_hour,omitempty"`
	HourOfDay       uint32     `json:"hour_of_day,omitempty"`
	Year            uint32     `json:"year,omitempty"`
	Month           uint32     `json:"month,omitempty"`
	DayOfMonth      uint32     `json:"day_of_month,omitempty"`
	Repeat          string     `json:"repeat,omitempty"`
	Disable         bool       `json:"disable,omitempty"`
}

// SettingSecurity represents Infoblox struct setting:security
type SettingSecurity struct {
	AuditLogRollingEnable             bool                    `json:"audit_log_rolling_enable,omitempty"`
	AdminAccessItems                  []*Addressac            `json:"admin_access_items,omitempty"`
	HttpRedirectEnable                bool                    `json:"http_redirect_enable,omitempty"`
	LcdInputEnable                    bool                    `json:"lcd_input_enable,omitempty"`
	LoginBannerEnable                 bool                    `json:"login_banner_enable,omitempty"`
	LoginBannerText                   string                  `json:"login_banner_text,omitempty"`
	RemoteConsoleAccessEnable         bool                    `json:"remote_console_access_enable,omitempty"`
	SecurityAccessEnable              bool                    `json:"security_access_enable,omitempty"`
	SecurityAccessRemoteConsoleEnable bool                    `json:"security_access_remote_console_enable,omitempty"`
	SessionTimeout                    uint32                  `json:"session_timeout,omitempty"`
	SshPermEnable                     bool                    `json:"ssh_perm_enable,omitempty"`
	SupportAccessEnable               bool                    `json:"support_access_enable,omitempty"`
	SupportAccessInfo                 string                  `json:"support_access_info,omitempty"`
	DisableConcurrentLogin            bool                    `json:"disable_concurrent_login,omitempty"`
	InactivityLockoutSetting          *SettingInactivelockout `json:"inactivity_lockout_setting,omitempty"`
}

// SettingSecuritybanner represents Infoblox struct setting:securitybanner
type SettingSecuritybanner struct {
	Color   string `json:"color,omitempty"`
	Level   string `json:"level,omitempty"`
	Message string `json:"message,omitempty"`
	Enable  bool   `json:"enable,omitempty"`
}

// SettingSnmp represents Infoblox struct setting:snmp
type SettingSnmp struct {
	EngineId               []string        `json:"engine_id,omitempty"`
	QueriesCommunityString string          `json:"queries_community_string,omitempty"`
	QueriesEnable          bool            `json:"queries_enable,omitempty"`
	Snmpv3QueriesEnable    bool            `json:"snmpv3_queries_enable,omitempty"`
	Snmpv3QueriesUsers     []*Queriesuser  `json:"snmpv3_queries_users,omitempty"`
	Snmpv3TrapsEnable      bool            `json:"snmpv3_traps_enable,omitempty"`
	Syscontact             []string        `json:"syscontact,omitempty"`
	Sysdescr               []string        `json:"sysdescr,omitempty"`
	Syslocation            []string        `json:"syslocation,omitempty"`
	Sysname                []string        `json:"sysname,omitempty"`
	TrapReceivers          []*Trapreceiver `json:"trap_receivers,omitempty"`
	TrapsCommunityString   string          `json:"traps_community_string,omitempty"`
	TrapsEnable            bool            `json:"traps_enable,omitempty"`
}

// SettingSyslogproxy represents Infoblox struct setting:syslogproxy
type SettingSyslogproxy struct {
	Enable     bool         `json:"enable,omitempty"`
	TcpEnable  bool         `json:"tcp_enable,omitempty"`
	TcpPort    uint32       `json:"tcp_port,omitempty"`
	UdpEnable  bool         `json:"udp_enable,omitempty"`
	UdpPort    uint32       `json:"udp_port,omitempty"`
	ClientAcls []*Addressac `json:"client_acls,omitempty"`
}

// SettingTrafficcapturechr represents Infoblox struct setting:trafficcapturechr
type SettingTrafficcapturechr struct {
	ChrTriggerEnable       bool   `json:"chr_trigger_enable,omitempty"`
	ChrThreshold           uint32 `json:"chr_threshold,omitempty"`
	ChrReset               uint32 `json:"chr_reset,omitempty"`
	ChrMinCacheUtilization uint32 `json:"chr_min_cache_utilization,omitempty"`
}

// SettingTrafficcaptureqps represents Infoblox struct setting:trafficcaptureqps
type SettingTrafficcaptureqps struct {
	QpsTriggerEnable bool   `json:"qps_trigger_enable,omitempty"`
	QpsThreshold     uint32 `json:"qps_threshold,omitempty"`
	QpsReset         uint32 `json:"qps_reset,omitempty"`
}

// SettingTriggerrecdnslatency represents Infoblox struct setting:triggerrecdnslatency
type SettingTriggerrecdnslatency struct {
	RecDnsLatencyTriggerEnable  bool                `json:"rec_dns_latency_trigger_enable,omitempty"`
	RecDnsLatencyThreshold      uint32              `json:"rec_dns_latency_threshold,omitempty"`
	RecDnsLatencyReset          uint32              `json:"rec_dns_latency_reset,omitempty"`
	RecDnsLatencyListenOnSource string              `json:"rec_dns_latency_listen_on_source,omitempty"`
	RecDnsLatencyListenOnIp     string              `json:"rec_dns_latency_listen_on_ip,omitempty"`
	KpiMonitoredDomains         []*Monitoreddomains `json:"kpi_monitored_domains,omitempty"`
}

// SettingTriggerrecqueries represents Infoblox struct setting:triggerrecqueries
type SettingTriggerrecqueries struct {
	RecursiveClientsCountTriggerEnable bool   `json:"recursive_clients_count_trigger_enable,omitempty"`
	RecursiveClientsCountThreshold     uint32 `json:"recursive_clients_count_threshold,omitempty"`
	RecursiveClientsCountReset         uint32 `json:"recursive_clients_count_reset,omitempty"`
}

// SettingTriggeruthdnslatency represents Infoblox struct setting:triggeruthdnslatency
type SettingTriggeruthdnslatency struct {
	AuthDnsLatencyTriggerEnable  bool   `json:"auth_dns_latency_trigger_enable,omitempty"`
	AuthDnsLatencyThreshold      uint32 `json:"auth_dns_latency_threshold,omitempty"`
	AuthDnsLatencyReset          uint32 `json:"auth_dns_latency_reset,omitempty"`
	AuthDnsLatencyListenOnSource string `json:"auth_dns_latency_listen_on_source,omitempty"`
	AuthDnsLatencyListenOnIp     string `json:"auth_dns_latency_listen_on_ip,omitempty"`
}

// SettingViewaddress represents Infoblox struct setting:viewaddress
type SettingViewaddress struct {
	ViewName                       string `json:"view_name,omitempty"`
	DnsNotifyTransferSource        string `json:"dns_notify_transfer_source,omitempty"`
	DnsNotifyTransferSourceAddress string `json:"dns_notify_transfer_source_address,omitempty"`
	DnsQuerySourceInterface        string `json:"dns_query_source_interface,omitempty"`
	DnsQuerySourceAddress          string `json:"dns_query_source_address,omitempty"`
	EnableNotifySourcePort         bool   `json:"enable_notify_source_port,omitempty"`
	NotifySourcePort               uint32 `json:"notify_source_port,omitempty"`
	EnableQuerySourcePort          bool   `json:"enable_query_source_port,omitempty"`
	QuerySourcePort                uint32 `json:"query_source_port,omitempty"`
	NotifyDelay                    uint32 `json:"notify_delay,omitempty"`
	UseSourcePorts                 bool   `json:"use_source_ports,omitempty"`
	UseNotifyDelay                 bool   `json:"use_notify_delay,omitempty"`
}

// Skipmemberupgrade represents Infoblox struct skipmemberupgrade
type Skipmemberupgrade struct{}

// SmartfolderGroupby represents Infoblox struct smartfolder:groupby
type SmartfolderGroupby struct {
	Value          string `json:"value,omitempty"`
	ValueType      string `json:"value_type,omitempty"`
	EnableGrouping bool   `json:"enable_grouping,omitempty"`
}

// SmartfolderQueryitem represents Infoblox struct smartfolder:queryitem
type SmartfolderQueryitem struct {
	Name      string                     `json:"name,omitempty"`
	FieldType string                     `json:"field_type,omitempty"`
	Operator  string                     `json:"operator,omitempty"`
	OpMatch   bool                       `json:"op_match,omitempty"`
	ValueType string                     `json:"value_type,omitempty"`
	Value     *SmartfolderQueryitemvalue `json:"value,omitempty"`
}

// SmartfolderQueryitemvalue represents Infoblox struct smartfolder:queryitemvalue
type SmartfolderQueryitemvalue struct {
	ValueInteger int        `json:"value_integer,omitempty"`
	ValueString  string     `json:"value_string,omitempty"`
	ValueDate    *time.Time `json:"value_date,omitempty"`
	ValueBoolean bool       `json:"value_boolean,omitempty"`
}

// Smartfoldersaveasparams represents Infoblox struct smartfoldersaveasparams
type Smartfoldersaveasparams struct{}

// Sortlist represents Infoblox struct sortlist
type Sortlist struct {
	Address   string   `json:"address,omitempty"`
	MatchList []string `json:"match_list,omitempty"`
}

// Splitipv6network represents Infoblox struct splitipv6network
type Splitipv6network struct{}

// Splitnetwork represents Infoblox struct splitnetwork
type Splitnetwork struct{}

// Startdiscovery represents Infoblox struct startdiscovery
type Startdiscovery struct{}

// Stopcsv represents Infoblox struct stopcsv
type Stopcsv struct{}

// Supportbundle represents Infoblox struct supportbundle
type Supportbundle struct{}

// SyslogEndpointServers represents Infoblox struct syslog:endpoint:servers
type SyslogEndpointServers struct {
	Address          string `json:"address,omitempty"`
	ConnectionType   string `json:"connection_type,omitempty"`
	Port             uint32 `json:"port,omitempty"`
	Hostname         string `json:"hostname,omitempty"`
	Format           string `json:"format,omitempty"`
	Facility         string `json:"facility,omitempty"`
	Severity         string `json:"severity,omitempty"`
	Certificate      string `json:"certificate,omitempty"`
	CertificateToken string `json:"certificate_token,omitempty"`
}

// Syslogserver represents Infoblox struct syslogserver
type Syslogserver struct {
	Address          string   `json:"address,omitempty"`
	Certificate      string   `json:"certificate,omitempty"`
	CertificateToken string   `json:"certificate_token,omitempty"`
	ConnectionType   string   `json:"connection_type,omitempty"`
	Port             uint32   `json:"port,omitempty"`
	LocalInterface   string   `json:"local_interface,omitempty"`
	MessageSource    string   `json:"message_source,omitempty"`
	MessageNodeId    string   `json:"message_node_id,omitempty"`
	Severity         string   `json:"severity,omitempty"`
	CategoryList     []string `json:"category_list,omitempty"`
	OnlyCategoryList bool     `json:"only_category_list,omitempty"`
}

// TacacsplusServer represents Infoblox struct tacacsplus:server
type TacacsplusServer struct {
	Address       string `json:"address,omitempty"`
	Port          uint32 `json:"port,omitempty"`
	SharedSecret  string `json:"shared_secret,omitempty"`
	AuthType      string `json:"auth_type,omitempty"`
	Comment       string `json:"comment,omitempty"`
	Disable       bool   `json:"disable,omitempty"`
	UseMgmtPort   bool   `json:"use_mgmt_port,omitempty"`
	UseAccounting bool   `json:"use_accounting,omitempty"`
}

// Taskcontrol represents Infoblox struct taskcontrol
type Taskcontrol struct{}

// TaxiiRpzconfig represents Infoblox struct taxii:rpzconfig
type TaxiiRpzconfig struct {
	CollectionName string `json:"collection_name,omitempty"`
	Zone           string `json:"zone,omitempty"`
}

// Testanalyticsserverconnectivityparams represents Infoblox struct testanalyticsserverconnectivityparams
type Testanalyticsserverconnectivityparams struct{}

// Testatpserverconnectivity represents Infoblox struct testatpserverconnectivity
type Testatpserverconnectivity struct{}

// Testconnectivityparams represents Infoblox struct testconnectivityparams
type Testconnectivityparams struct{}

// Testdxlbrokerconnectivity represents Infoblox struct testdxlbrokerconnectivity
type Testdxlbrokerconnectivity struct{}

// Testendpointconnection represents Infoblox struct testendpointconnection
type Testendpointconnection struct{}

// Testocsprespondersettings represents Infoblox struct testocsprespondersettings
type Testocsprespondersettings struct{}

// Testsyslog represents Infoblox struct testsyslog
type Testsyslog struct{}

// Testsyslogbackup represents Infoblox struct testsyslogbackup
type Testsyslogbackup struct{}

// Threatdetails represents Infoblox struct threatdetails
type Threatdetails struct{}

// ThreatprotectionNatport represents Infoblox struct threatprotection:natport
type ThreatprotectionNatport struct {
	StartPort uint32 `json:"start_port,omitempty"`
	EndPort   uint32 `json:"end_port,omitempty"`
	BlockSize uint32 `json:"block_size,omitempty"`
}

// ThreatprotectionNatrule represents Infoblox struct threatprotection:natrule
type ThreatprotectionNatrule struct {
	RuleType     string                     `json:"rule_type,omitempty"`
	Address      string                     `json:"address,omitempty"`
	Network      string                     `json:"network,omitempty"`
	Cidr         uint32                     `json:"cidr,omitempty"`
	StartAddress string                     `json:"start_address,omitempty"`
	EndAddress   string                     `json:"end_address,omitempty"`
	NatPorts     []*ThreatprotectionNatport `json:"nat_ports,omitempty"`
}

// ThreatprotectionRuleconfig represents Infoblox struct threatprotection:ruleconfig
type ThreatprotectionRuleconfig struct {
	Action      string                       `json:"action,omitempty"`
	LogSeverity string                       `json:"log_severity,omitempty"`
	Params      []*ThreatprotectionRuleparam `json:"params,omitempty"`
}

// ThreatprotectionRuleparam represents Infoblox struct threatprotection:ruleparam
type ThreatprotectionRuleparam struct {
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Syntax      string   `json:"syntax,omitempty"`
	Value       string   `json:"value,omitempty"`
	Min         uint32   `json:"min,omitempty"`
	Max         uint32   `json:"max,omitempty"`
	ReadOnly    bool     `json:"read_only,omitempty"`
	EnumValues  []string `json:"enum_values,omitempty"`
}

// ThreatprotectionStatinfo represents Infoblox struct threatprotection:statinfo
type ThreatprotectionStatinfo struct {
	Timestamp     *time.Time `json:"timestamp,omitempty"`
	Critical      uint64     `json:"critical,omitempty"`
	Major         uint64     `json:"major,omitempty"`
	Warning       uint64     `json:"warning,omitempty"`
	Informational uint64     `json:"informational,omitempty"`
	Total         uint64     `json:"total,omitempty"`
}

// Thresholdtrap represents Infoblox struct thresholdtrap
type Thresholdtrap struct {
	TrapType    string `json:"trap_type,omitempty"`
	TrapReset   uint32 `json:"trap_reset,omitempty"`
	TrapTrigger uint32 `json:"trap_trigger,omitempty"`
}

// Trapnotification represents Infoblox struct trapnotification
type Trapnotification struct {
	TrapType    string `json:"trap_type,omitempty"`
	EnableEmail bool   `json:"enable_email,omitempty"`
	EnableTrap  bool   `json:"enable_trap,omitempty"`
}

// Trapreceiver represents Infoblox struct trapreceiver
type Trapreceiver struct {
	Address string `json:"address,omitempty"`
	User    string `json:"user,omitempty"`
	Comment string `json:"comment,omitempty"`
}

// Triggeroutboundparams represents Infoblox struct triggeroutboundparams
type Triggeroutboundparams struct{}

// Tsigac represents Infoblox struct tsigac
type Tsigac struct {
	Address        string `json:"address,omitempty"`
	Permission     string `json:"permission,omitempty"`
	TsigKey        string `json:"tsig_key,omitempty"`
	TsigKeyAlg     string `json:"tsig_key_alg,omitempty"`
	TsigKeyName    string `json:"tsig_key_name,omitempty"`
	UseTsigKeyName bool   `json:"use_tsig_key_name,omitempty"`
}

// Updateatprulesetparams represents Infoblox struct updateatprulesetparams
type Updateatprulesetparams struct{}

// Updatelicenses represents Infoblox struct updatelicenses
type Updatelicenses struct{}

// Updatesdownloadmemberconfig represents Infoblox struct updatesdownloadmemberconfig
type Updatesdownloadmemberconfig struct {
	Member    string `json:"member,omitempty"`
	Interface string `json:"interface,omitempty"`
	IsOnline  bool   `json:"is_online,omitempty"`
}

// Upgrade represents Infoblox struct upgrade
type Upgrade struct{}

// UpgradegroupMember represents Infoblox struct upgradegroup:member
type UpgradegroupMember struct {
	Member   string `json:"member,omitempty"`
	TimeZone string `json:"time_zone,omitempty"`
}

// UpgradegroupSchedule represents Infoblox struct upgradegroup:schedule
type UpgradegroupSchedule struct {
	Name                       string     `json:"name,omitempty"`
	TimeZone                   string     `json:"time_zone,omitempty"`
	DistributionDependentGroup string     `json:"distribution_dependent_group,omitempty"`
	UpgradeDependentGroup      string     `json:"upgrade_dependent_group,omitempty"`
	DistributionTime           *time.Time `json:"distribution_time,omitempty"`
	UpgradeTime                *time.Time `json:"upgrade_time,omitempty"`
}

// Upgradegroupnowparams represents Infoblox struct upgradegroupnowparams
type Upgradegroupnowparams struct{}

// Upgradestep represents Infoblox struct upgradestep
type Upgradestep struct {
	StatusValue string `json:"status_value,omitempty"`
	StatusText  string `json:"status_text,omitempty"`
}

// Uploadcertificate represents Infoblox struct uploadcertificate
type Uploadcertificate struct{}

// Uploadkeytab represents Infoblox struct uploadkeytab
type Uploadkeytab struct{}

// Uploadserviceaccount represents Infoblox struct uploadserviceaccount
type Uploadserviceaccount struct{}

// Validateaclitems represents Infoblox struct validateaclitems
type Validateaclitems struct{}

// Vdiscoverycontrol represents Infoblox struct vdiscoverycontrol
type Vdiscoverycontrol struct{}

// Vlanlink represents Infoblox struct vlanlink
type Vlanlink struct {
	//Vlan *Subobj `json:"vlan,omitempty"`
	Id   uint32 `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// Vtftpdirmember represents Infoblox struct vtftpdirmember
type Vtftpdirmember struct {
	Member       string `json:"member,omitempty"`
	IpType       string `json:"ip_type,omitempty"`
	Address      string `json:"address,omitempty"`
	StartAddress string `json:"start_address,omitempty"`
	EndAddress   string `json:"end_address,omitempty"`
	Network      string `json:"network,omitempty"`
	Cidr         uint32 `json:"cidr,omitempty"`
}

// Zoneassociation represents Infoblox struct zoneassociation
type Zoneassociation struct {
	Fqdn      string `json:"fqdn,omitempty"`
	IsDefault bool   `json:"is_default,omitempty"`
	View      string `json:"view,omitempty"`
}

// ZoneNameServer represents Infoblox struct zonenameserver
type ZoneNameServer struct {
	Address       string `json:"address,omitempty"`
	AutoCreatePtr bool   `json:"auto_create_ptr,omitempty"`
}
