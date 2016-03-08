package ibclient

import (
	"encoding/json"
)

type Bool bool

type EA map[string]interface{}

type Payload map[string]interface{}

type NetworkView struct {
	Ref  string `json:"_ref"`
	Name string `json:"name"`
}

type Network struct {
	Ref         string `json:"_ref"`
	NetviewName string `json:"network_view"`
	Cidr        string `json:"network"`
}

type NetworkContainer struct {
	Ref         string `json:"_ref"`
	NetviewName string `json:"network_view"`
	Cidr        string `json:"network"`
}

type FixedAddress struct {
	Ref         string `json:"_ref"`
	NetviewName string `json:"network_view"`
	Cidr        string `json:"network"`
	IPAddress   string `json:"ipv4addr`
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
