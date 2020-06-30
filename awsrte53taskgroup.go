package ibclient

type AwsRte53TaskOperations interface {
	CreateAwsRte53TaskGroup(aws AwsRte53TaskGroup) (*AwsRte53TaskGroup, error)
	GetAwsRte53TaskGroup(aws AwsRte53TaskGroup) (*[]AwsRte53TaskGroup, error)
	DeleteAwsRte53TaskGroup(aws AwsRte53TaskGroup) (string, error)
	UpdateAwsRte53TaskGroup(aws AwsRte53TaskGroup) (*AwsRte53TaskGroup, error)
}

type AwsUser struct {
	IBBase          `json:"-"`
	AccountID       string `json:"account_id,omitempty"`
	AccessKeyId     string `json:"access_key_id,omitempty"`
	LastUsed        int    `json:"last_used,omitempty"`
	Name            string `json:"name,omitempty"`
	NiosUserName    string `json:"nios_user_name,omitempty"`
	SecretAccessKey string `json:"secret_access_key,omitempty"`
	Status          string `json:"status,omitempty"`
}

// NewAwsUser creates a new AwsUser type with objectType and returnFields
func NewAwsUser(aws AwsUser) *AwsUser {
	res := aws
	res.objectType = "awsuser"
	res.returnFields = []string{"account_id", "last_used", "name", "nios_user_name", "status"}
	return &res
}

// GetAwsUser gets the required details of the existing user
func (objMgr *ObjectManager) GetAwsUser(aws AwsUser) (*AwsUser, error) {
	var res []AwsUser
	awsRte53 := NewAwsUser(aws)
	var err error
	err = objMgr.connector.GetObject(awsRte53, "", &res)
	if err != nil || res == nil || len(res) == 0 {
		return nil, err
	}

	return &res[0], err
}
/*type AwsRte53Task struct {
	IBBase           `json:"-"`
	AwsUser          string `json:"aws_user,omitempty"`
	CredentialsType  string `json:"credentials_type,omitempty"`
	Disabled         Bool   `json:"disabled,omitempty"`
	Filter           string `json:"filter,omitempty"`
	LastRun          int    `json:"last_run,omitempty"`
	Name             string `json:"name,omitempty"`
	ScheduleInterval uint   `json:"schedule_interval,omitempty"`
	ScheduleUnits    string `json:"schedule_units,omitempty"`
	State            string `json:"state,omitempty"`
	StateMsg         string `json:"state_msg,omitempty"`
	StatusTimestamp  int    `json:"status_timestamp,omitempty"`
	SyncPrivateZones Bool   `json:"sync_private_zones,omitempty"`
	SyncPublicZones  Bool   `json:"sync_public_zones,omitempty"`
	ZoneCount        uint   `json:"zone_count,omitempty"`
}*/

type AwsRte53TaskGroup struct {
	IBBase                   `json:"-"`
	Ref                      string `json:"_ref,omitempty"`
	AccountId                string `json:"account_id,omitempty"`
	Comment                  string `json:"comment,omitempty"`
	ConsolidateZones         Bool   `json:"consolidate_zones,omitempty"`
	ConsolidatedView         string `json:"consolidated_view,omitempty"`
	Disabled                 Bool   `json:"disabled,omitempty"`
	GridMember               string `json:"grid_member,omitempty"`
	Name                     string `json:"name,omitempty"`
	NetworkView              string `json:"network_view,omitempty"`
	NetworkViewMappingPolicy string `json:"network_view_mapping_policy,omitempty"`
	SyncStatus               string `json:"sync_status,omitempty"`
}

// NewAwsRte53TaskGroup creates a new AwsRte53TaskGroup type with objectType and returnFields
func NewAwsRte53TaskGroup(aws AwsRte53TaskGroup) *AwsRte53TaskGroup {
	res := aws
	res.objectType = "awsrte53taskgroup"

	res.returnFields = []string{"account_id", "comment", "consolidate_zones", "consolidated_view", "disabled",
		"grid_member", "name", "network_view", "network_view_mapping_policy", "sync_status"}
	return &res
}

// CreateAwsRte53Task takes Name, View, GridMember as arguments to create AwsRte53Task
func (objMgr *ObjectManager) CreateAwsRte53TaskGroup(aws AwsRte53TaskGroup) (*AwsRte53TaskGroup, error) {

	if aws.NetworkView != "" {
		aws.NetworkViewMappingPolicy = "DIRECT"
	}
	awsRte53 := NewAwsRte53TaskGroup(aws)
	ref, err := objMgr.connector.CreateObject(awsRte53)
	awsRte53.Ref = ref
	return awsRte53, err
}

// GetAwsRte53TaskGroup by passing Name, reference ID or NetworkView
// If no arguments are passed then, all the tasks are returned
func (objMgr *ObjectManager) GetAwsRte53TaskGroup(aws AwsRte53TaskGroup) (*[]AwsRte53TaskGroup, error) {

	var res []AwsRte53TaskGroup
	awsRte53 := NewAwsRte53TaskGroup(aws)
	var err error
	if len(aws.Ref) > 0 {
		err = objMgr.connector.GetObject(awsRte53, aws.Ref, &awsRte53)
		res = append(res, *awsRte53)

	} else {
		awsRte53 = NewAwsRte53TaskGroup(aws)
		err = objMgr.connector.GetObject(awsRte53, "", &res)
		if err != nil || res == nil || len(res) == 0 {
			return nil, err
		}
	}

	return &res, err
}

// DeleteAwsRte53TaskGroup by passing either Reference or Name
func (objMgr *ObjectManager) DeleteAwsRte53TaskGroup(aws AwsRte53TaskGroup) (string, error) {
	var res []AwsRte53TaskGroup
	awsRte53 := NewAwsRte53TaskGroup(aws)
	if len(aws.Ref) > 0 {
		return objMgr.connector.DeleteObject(aws.Ref)

	} else {
		awsRte53 = NewAwsRte53TaskGroup(aws)
		err := objMgr.connector.GetObject(awsRte53, "", &res)
		if err != nil || res == nil || len(res) == 0 {
			return "Task doesn't exist", err
		}
		return objMgr.connector.DeleteObject(res[0].Ref)
	}

}

// UpdateAwsRte53GroupTask takes Reference ID of the task as an argument
// to update Name
// returns updated Refernce ID
func (objMgr *ObjectManager) UpdateAwsRte53TaskGroup(aws AwsRte53TaskGroup) (*AwsRte53TaskGroup, error) {
	var res AwsRte53TaskGroup
	awsRte53 := AwsRte53TaskGroup{}
	awsRte53.returnFields = []string{"name"}
	err := objMgr.connector.GetObject(&awsRte53, aws.Ref, &res)
	if err != nil {
		return nil, err
	}
	res.Name = aws.Name
	reference, err := objMgr.connector.UpdateObject(&res, aws.Ref)
	res.Ref = reference
	return &res, err
}
