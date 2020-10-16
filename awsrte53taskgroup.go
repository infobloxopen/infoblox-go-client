package ibclient

type AWSRoute53TaskOperations interface {
	CreateAWSRoute53TaskGroup(aws AWSRoute53TaskGroup) (*AWSRoute53TaskGroup, error)
	GetAWSRoute53TaskGroup(aws AWSRoute53TaskGroup) (*[]AWSRoute53TaskGroup, error)
	DeleteAWSRoute53TaskGroup(aws AWSRoute53TaskGroup) (string, error)
	UpdateAWSRoute53TaskGroup(aws AWSRoute53TaskGroup) (*AWSRoute53TaskGroup, error)
}

type AWSUser struct {
	IBBase          `json:"-"`
	AccountID       string `json:"account_id,omitempty"`
	AccessKeyId     string `json:"access_key_id,omitempty"`
	LastUsed        int    `json:"last_used,omitempty"`
	Name            string `json:"name,omitempty"`
	NiosUserName    string `json:"nios_user_name,omitempty"`
	SecretAccessKey string `json:"secret_access_key,omitempty"`
	Status          string `json:"status,omitempty"`
}

// NewAWSUser creates a new AWSUser type with objectType and returnFields
func NewAWSUser(aws AWSUser) *AWSUser {
	res := aws
	res.objectType = "awsuser"
	res.returnFields = []string{"account_id", "last_used", "name", "nios_user_name", "status"}
	return &res
}

// GetAWSUser gets the required details of the existing user
func (objMgr *ObjectManager) GetAWSUser(aws AWSUser) (*AWSUser, error) {
	var res []AWSUser
	awsRte53 := NewAWSUser(aws)
	var err error
	err = objMgr.connector.GetObject(awsRte53, "", &res)
	if err != nil || res == nil || len(res) == 0 {
		return nil, err
	}

	return &res[0], err
}

type AWSRoute53TaskGroup struct {
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

// NewAWSRoute53TaskGroup creates a new AWSRoute53TaskGroup type with objectType and returnFields
func NewAWSRoute53TaskGroup(aws AWSRoute53TaskGroup) *AWSRoute53TaskGroup {
	res := aws
	res.objectType = "awsrte53taskgroup"

	res.returnFields = []string{"account_id", "comment", "consolidate_zones", "consolidated_view", "disabled",
		"grid_member", "name", "network_view", "network_view_mapping_policy", "sync_status"}
	return &res
}

// CreateAWSRoute53Task takes Name, View, GridMember as arguments to create AWSRoute53Task
func (objMgr *ObjectManager) CreateAWSRoute53TaskGroup(aws AWSRoute53TaskGroup) (*AWSRoute53TaskGroup, error) {

	if aws.NetworkView != "" {
		aws.NetworkViewMappingPolicy = "DIRECT"
	}
	awsRte53 := NewAWSRoute53TaskGroup(aws)
	ref, err := objMgr.connector.CreateObject(awsRte53)
	awsRte53.Ref = ref
	return awsRte53, err
}

// GetAWSRoute53TaskGroup by passing Name, reference ID or NetworkView
// If no arguments are passed then, all the tasks are returned
func (objMgr *ObjectManager) GetAWSRoute53TaskGroup(aws AWSRoute53TaskGroup) (*[]AWSRoute53TaskGroup, error) {

	var res []AWSRoute53TaskGroup
	awsRte53 := NewAWSRoute53TaskGroup(aws)
	var err error
	if len(aws.Ref) > 0 {
		err = objMgr.connector.GetObject(awsRte53, aws.Ref, &awsRte53)
		res = append(res, *awsRte53)

	} else {
		awsRte53 = NewAWSRoute53TaskGroup(aws)
		err = objMgr.connector.GetObject(awsRte53, "", &res)
		if err != nil || res == nil || len(res) == 0 {
			return nil, err
		}
	}

	return &res, err
}

// DeleteAWSRoute53TaskGroup by passing either Reference or Name
func (objMgr *ObjectManager) DeleteAWSRoute53TaskGroup(aws AWSRoute53TaskGroup) (string, error) {
	var res []AWSRoute53TaskGroup
	awsRte53 := NewAWSRoute53TaskGroup(aws)
	if len(aws.Ref) > 0 {
		return objMgr.connector.DeleteObject(aws.Ref)

	} else {
		awsRte53 = NewAWSRoute53TaskGroup(aws)
		err := objMgr.connector.GetObject(awsRte53, "", &res)
		if err != nil || res == nil || len(res) == 0 {
			return "Task doesn't exist", err
		}
		return objMgr.connector.DeleteObject(res[0].Ref)
	}

}

// UpdateAWSRoute53GroupTask takes Reference ID of the task as an argument
// to update Name
// returns updated Refernce ID
func (objMgr *ObjectManager) UpdateAWSRoute53TaskGroup(aws AWSRoute53TaskGroup) (*AWSRoute53TaskGroup, error) {
	var res AWSRoute53TaskGroup
	awsRte53 := AWSRoute53TaskGroup{}
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
