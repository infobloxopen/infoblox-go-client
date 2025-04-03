package ibclient

import "fmt"

func (objMgr *ObjectManager) CreateRangeTemplate(name string, numberOfAdresses uint32, offset uint32, comment string, ea EA,
	options []*Dhcpoption, useOption bool, serverAssociationType string, failOverAssociation string, member *Dhcpmember) (*Rangetemplate, error) {
	if name == "" {
		return nil, fmt.Errorf("name field is required to create a Range Template object")
	}
	rangeTemplate := NewRangeTemplate("", name, numberOfAdresses, offset, comment, ea, options,
		useOption, serverAssociationType, failOverAssociation, member)
	ref, err := objMgr.connector.CreateObject(rangeTemplate)
	if err != nil {
		return nil, fmt.Errorf("error creating Range Template object %s, err: %s", name, err)
	}
	rangeTemplate.Ref = ref
	return rangeTemplate, nil
}

func (objMgr *ObjectManager) DeleteRangeTemplate(ref string) (string, error) {
	return objMgr.connector.DeleteObject(ref)
}

func (objMgr *ObjectManager) GetAllRangeTemplate(queryParams *QueryParams) ([]Rangetemplate, error) {
	var res []Rangetemplate
	rangeTemplate := NewEmptyRangeTemplate()
	err := objMgr.connector.GetObject(rangeTemplate, "", queryParams, &res)
	if err != nil {
		return nil, fmt.Errorf("error getting Range Template object, err: %s", err)
	}
	return res, nil
}

func (objMgr *ObjectManager) GetRangeTemplateByRef(ref string) (*Rangetemplate, error) {
	rangeTemplate := NewEmptyRangeTemplate()
	err := objMgr.connector.GetObject(rangeTemplate, ref, NewQueryParams(false, nil), &rangeTemplate)
	if err != nil {
		return nil, err
	}
	return rangeTemplate, nil
}

func (objMgr *ObjectManager) UpdateRangeTemplate(ref string, name string, numberOfAddresses uint32, offset uint32, comment string, ea EA,
	options []*Dhcpoption, useOption bool, serverAssociationType string, failOverAssociation string, member *Dhcpmember) (*Rangetemplate, error) {
	if name == "" {
		return nil, fmt.Errorf("name field is required to update a Range Template object")
	}
	rangeTemplate := NewRangeTemplate(ref, name, numberOfAddresses, offset, comment, ea, options, useOption,
		serverAssociationType, failOverAssociation, member)
	newRef, err := objMgr.connector.UpdateObject(rangeTemplate, ref)
	if err != nil {
		return nil, fmt.Errorf("error updating Range Template object %s, err: %s", name, err)
	}
	rangeTemplate.Ref = newRef
	rangeTemplate, err = objMgr.GetRangeTemplateByRef(newRef)
	if err != nil {
		return nil, fmt.Errorf("error getting updated Range Template object %s, err: %s", name, err)
	}
	return rangeTemplate, nil
}

func NewRangeTemplate(ref string, name string, numberOfAddresses uint32, offset uint32, comment string, ea EA,
	options []*Dhcpoption, useOption bool, serverAssociationType string, failOverAssociation string, member *Dhcpmember) *Rangetemplate {
	rangeTemplate := NewEmptyRangeTemplate()
	rangeTemplate.Ref = ref
	rangeTemplate.Name = &name
	rangeTemplate.NumberOfAddresses = &numberOfAddresses
	rangeTemplate.Offset = &offset
	rangeTemplate.Comment = &comment
	rangeTemplate.Ea = ea
	rangeTemplate.Options = options
	rangeTemplate.UseOptions = &useOption
	rangeTemplate.ServerAssociationType = serverAssociationType
	rangeTemplate.FailoverAssociation = &failOverAssociation
	rangeTemplate.Member = member
	return rangeTemplate
}

func NewEmptyRangeTemplate() *Rangetemplate {
	rangeTemplate := &Rangetemplate{}
	rangeTemplate.SetReturnFields(append(rangeTemplate.ReturnFields(), "extattrs", "options", "use_options",
		"server_association_type", "failover_association", "member"))
	return rangeTemplate
}
