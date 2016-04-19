package ibclient

import (
	"encoding/json"
)

type Bool bool

type EA map[string]interface{}

type Payload map[string]interface{}

type IBBase struct {
	objectType string `json:"-"`
}

type IBObject interface {
	ObjectType() string
}

func (obj *IBBase) ObjectType() string {
	return obj.objectType
}

type NetworkView struct {
	IBBase `json:"-"`
	Ref    string `json:"_ref,omitempty"`
	Name   string `json:"name,omitempty"`
	Ea     EA     `json:"extattrs,omitempty"`
}

func NewNetworkView() *NetworkView {
	return &NetworkView{IBBase: IBBase{"networkview"}}
}

type Network struct {
	IBBase      `json:"-"`
	Ref         string `json:"_ref,omitempty"`
	NetviewName string `json:"network_view,omitempty"`
	Cidr        string `json:"network,omitempty"`
	Ea          EA     `json:"extattrs,omitempty"`
}

func NewNetwork() *Network {
	return &Network{IBBase: IBBase{"network"}}
}

type NetworkContainer struct {
	IBBase      `json:"-"`
	Ref         string `json:"_ref,omitempty"`
	NetviewName string `json:"network_view,omitempty"`
	Cidr        string `json:"network,omitempty"`
	Ea          EA     `json:"extattrs,omitempty"`
}

func NewNetworkContainer() *NetworkContainer {
	return &NetworkContainer{IBBase: IBBase{"networkcontainer"}}
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

func NewFixedAddress() *FixedAddress {
	return &FixedAddress{IBBase: IBBase{"fixedaddress"}}
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

func (b Bool) MarshalJSON() ([]byte, error) {
	if b {
		return json.Marshal("True")
	}

	return json.Marshal("False")
}
