package ibclient

import (
	"bytes"
	"encoding/json"
	"reflect"
)

const MACADDR_ZERO = "00:00:00:00:00:00"

type Bool bool

type EA map[string]interface{}
type EAS map[string]interface{}

type EADefListValue string

type IBBase struct {
	objectType   string   `json:"-"`
	returnFields []string `json:"-"`
	eaSearch     EAS      `json:"-"`
}

type IBObject interface {
	ObjectType() string
	ReturnFields() []string
	EaSearch() EAS
}

func (obj *IBBase) ObjectType() string {
	return obj.objectType
}

func (obj *IBBase) ReturnFields() []string {
	return obj.returnFields
}

func (obj *IBBase) EaSearch() EAS {
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

type Network struct {
	IBBase
	Ref         string `json:"_ref,omitempty"`
	NetviewName string `json:"network_view,omitempty"`
	Cidr        string `json:"network,omitempty"`
	Ea          EA     `json:"extattrs,omitempty"`
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

func (ea EA) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	for k, v := range ea {
		value := make(map[string]interface{})
		value["value"] = v
		m[k] = value
	}

	return json.Marshal(m)
}

func (eas EAS) MarshalJSON() ([]byte, error) {
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
