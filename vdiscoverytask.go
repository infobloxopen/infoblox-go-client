package ibclient

type VDiscoveryTaskOperations interface {
	CreateVDiscoveryTask(recVDis VDiscoveryTask) (*VDiscoveryTask, error)
	GetVDiscoveryTask(recVDis VDiscoveryTask) (*[]VDiscoveryTask, error)
	DeleteVDiscoveryTask(recVDis VDiscoveryTask) (string, error)
	UpdateVDiscoveryTask(recVDis VDiscoveryTask) (*VDiscoveryTask, error)
}

type VDiscoveryTask struct {
	IBBase                          `json:"-"`
	Ref                             string `json:"_ref,omitempty"`
	Name                            string `json:"name,omitempty"`
	DriverType                      string `json:"driver_type,omitempty"`
	FqdnOrIp                        string `json:"fqdn_or_ip,omitempty"`
	Username                        string `json:"username,omitempty"`
	Password                        string `json:"password,omitempty"`
	MemberV                         string `json:"member,omitempty"`
	Port                            uint   `json:"port,omitempty"`
	Protocol                        string `json:"protocol,omitempty"`
	AutoConsolidateCloudEa          bool   `json:"auto_consolidate_cloud_ea,omitempty"`
	AutoConsolidateManagedTenant    bool   `json:"auto_consolidate_managed_tenant,omitempty"`
	AutoConsolidateManagedVm        bool   `json:"auto_consolidate_managed_vm,omitempty"`
	MergeData                       bool   `json:"merge_data,omitempty"`
	PrivateNetworkView              string `json:"private_network_view,omitempty"`
	PrivateNetworkViewMappingPolicy string `json:"private_network_view_mapping_policy,omitempty"`
	PublicNetworkView               string `json:"public_network_view,omitempty"`
	PublicNetworkViewMappingPolicy  string `json:"public_network_view_mapping_policy,omitempty"`
	AllowUnsecuredConnection        bool   `json:"allow_unsecured_connection,omitempty"`
	AutoCreateDnsHostnameTemplate   string `json:"auto_create_dns_hostname_template,omitempty"`
	AutoCreateDnsRecord             bool   `json:"auto_create_dns_record,omitempty"`
	AutoCreateDnsRecordType         string `json:"auto_create_dns_record_type,omitempty"`
	Comment                         string `json:"comment,omitempty"`
	CredentialsType                 string `json:"credentials_type,omitempty"`
	DnsViewPrivateIp                string `json:"dns_view_private_ip,omitempty"`
	DnsViewPublicIp                 string `json:"dns_view_public_ip,omitempty"`
	DomainName                      string `json:"domain_name,omitempty"`
	Enabled                         bool   `json:"enabled,omitempty"`
	IdentityVersion                 string `json:"identity_version,omitempty"`
	LastRun                         int    `json:"last_run,omitempty"`
	State                           string `json:"state,omitempty"`
	StateMsg                        string `json:"state_msg,omitempty"`
	UpdateDnsViewPrivateIp          bool   `json:"update_dns_view_private_ip,omitempty"`
	UpdateDnsViewPublicIp           bool   `json:"update_dns_view_public_ip,omitempty"`
	UpdateMetadata                  bool   `json:"update_metadata,omitempty"`
}

// NewVDiscoveryTask creates a new vDiscoveryTask with objectType and returnFields
func NewVDiscoveryTask(vdis VDiscoveryTask) *VDiscoveryTask {
	res := vdis
	res.objectType = "vdiscoverytask"

	res.returnFields = []string{"name", "driver_type", "fqdn_or_ip", "username", "member",
		"port", "protocol", "auto_consolidate_cloud_ea", "auto_consolidate_managed_tenant",
		"auto_consolidate_managed_vm", "merge_data", "private_network_view", "private_network_view_mapping_policy",
		"public_network_view", "public_network_view_mapping_policy", "allow_unsecured_connection",
		"auto_create_dns_hostname_template", "auto_create_dns_record", "auto_create_dns_record_type",
		"comment", "credentials_type", "dns_view_private_ip", "dns_view_public_ip", "domain_name",
		"enabled", "identity_version", "state", "state_msg",
		"update_dns_view_private_ip", "update_dns_view_public_ip", "update_metadata"}
	return &res
}

// CreateVDiscoveryTask creates a vDiscovery Task
func (objMgr *ObjectManager) CreateVDiscoveryTask(vDis VDiscoveryTask) (*VDiscoveryTask, error) {

	vDiscovery := NewVDiscoveryTask(vDis)
	ref, err := objMgr.connector.CreateObject(vDiscovery)
	vDiscovery.Ref = ref
	return vDiscovery, err
}

// GetVDiscoveryTask by passing Name, reference ID or DNS View
// If no arguments are passed then, all the tasks are returned
func (objMgr *ObjectManager) GetVDiscoveryTask(vDis VDiscoveryTask) (*[]VDiscoveryTask, error) {

	var res []VDiscoveryTask
	vDiscovery := NewVDiscoveryTask(vDis)
	var err error
	if len(vDis.Ref) > 0 {
		err = objMgr.connector.GetObject(vDiscovery, vDis.Ref, &vDiscovery)
		res = append(res, *vDiscovery)

	} else {
		vDiscovery = NewVDiscoveryTask(vDis)
		err = objMgr.connector.GetObject(vDiscovery, "", &res)
		if err != nil || res == nil || len(res) == 0 {
			return nil, err
		}
	}

	return &res, err
}

// DeleteVDiscoveryTask by passing either Reference or Name
func (objMgr *ObjectManager) DeleteVDiscoveryTask(vDis VDiscoveryTask) (string, error) {
	var res []VDiscoveryTask
	vDiscovery := NewVDiscoveryTask(vDis)
	if len(vDis.Ref) > 0 {
		return objMgr.connector.DeleteObject(vDis.Ref)

	} else {
		vDiscovery = NewVDiscoveryTask(vDis)
		err := objMgr.connector.GetObject(vDiscovery, "", &res)
		if err != nil || res == nil || len(res) == 0 {
			return "Task doesn't exist", err
		}
		return objMgr.connector.DeleteObject(res[0].Ref)
	}

}

// UpdateVDiscoveryTask takes Reference ID of the task as an argument
// to update Name
// returns updated Refernce ID
func (objMgr *ObjectManager) UpdateVDiscoveryTask(vDis VDiscoveryTask) (*VDiscoveryTask, error) {
	var res VDiscoveryTask
	vDiscovery := VDiscoveryTask{}
	vDiscovery.returnFields = []string{"name"}
	err := objMgr.connector.GetObject(&vDiscovery, vDis.Ref, &res)
	if err != nil {
		return nil, err
	}
	res.Name = vDis.Name
	reference, err := objMgr.connector.UpdateObject(&res, vDis.Ref)
	res.Ref = reference
	return &res, err
}
