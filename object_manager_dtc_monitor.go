package ibclient

import (
	"encoding/json"
	"fmt"
)

// omit empty value for "client_cert" from marshaling because it holds a ref which would be invalid
func (d *DtcMonitorHttp) MarshalJSON() ([]byte, error) {
	type Alias DtcMonitorHttp
	aux := Alias(*d)
	if aux.ClientCert != nil && *aux.ClientCert == "" {
		aux.ClientCert = nil
	}
	return json.Marshal(aux)
}

func NewEmptyDtcMonitorHttp() *DtcMonitorHttp {
	dtcMonitorHttp := DtcMonitorHttp{}
	dtcMonitorHttp.SetReturnFields(append(dtcMonitorHttp.ReturnFields(),
		"extattrs",
		"interval",
		"retry_down",
		"retry_up",
		"timeout",
		"ciphers",
		"client_cert",
		"content_check",
		"content_check_input",
		"content_check_op",
		"content_check_regex",
		"content_extract_group",
		"content_extract_type",
		"content_extract_value",
		"enable_sni",
		"port",
		"request",
		"result",
		"result_code",
		"secure",
		"validate_cert",
	))
	return &dtcMonitorHttp
}

func NewEmptyDtcMonitorIcmp() *DtcMonitorIcmp {
	dtcMonitorIcmp := DtcMonitorIcmp{}
	dtcMonitorIcmp.SetReturnFields(append(dtcMonitorIcmp.ReturnFields(),
		"extattrs",
		"interval",
		"retry_down",
		"retry_up",
		"timeout",
	))
	return &dtcMonitorIcmp
}

func NewDtcMonitorHttp(
	comment, name string, eas EA, port, interval, retry_down, retry_up, timeout uint32,
	ciphers, client_cert, content_check, content_check_input, content_check_op, content_check_regex string,
	content_extract_group uint32, content_extract_type, content_extract_value string,
	enable_sni bool, request, result string, result_code uint32, secure, validate_cert bool,
) *DtcMonitorHttp {
	dtcMonitorHttp := NewEmptyDtcMonitorHttp()
	dtcMonitorHttp.Comment = &comment
	dtcMonitorHttp.Name = &name
	dtcMonitorHttp.Ea = eas
	dtcMonitorHttp.Port = &port
	dtcMonitorHttp.Interval = &interval
	dtcMonitorHttp.RetryDown = &retry_down
	dtcMonitorHttp.RetryUp = &retry_up
	dtcMonitorHttp.Timeout = &timeout
	dtcMonitorHttp.Ciphers = &ciphers
	dtcMonitorHttp.ClientCert = &client_cert
	dtcMonitorHttp.ContentCheck = content_check
	dtcMonitorHttp.ContentCheckInput = content_check_input
	dtcMonitorHttp.ContentCheckOp = content_check_op
	dtcMonitorHttp.ContentCheckRegex = &content_check_regex
	dtcMonitorHttp.ContentExtractGroup = &content_extract_group
	dtcMonitorHttp.ContentExtractType = content_extract_type
	dtcMonitorHttp.ContentExtractValue = &content_extract_value
	dtcMonitorHttp.EnableSni = &enable_sni
	dtcMonitorHttp.Request = &request
	dtcMonitorHttp.Result = result
	dtcMonitorHttp.ResultCode = &result_code
	dtcMonitorHttp.Secure = &secure
	dtcMonitorHttp.ValidateCert = &validate_cert
	return dtcMonitorHttp
}

func NewDtcMonitorIcmp(
	comment, name string, eas EA,
	interval, retry_down, retry_up, timeout uint32,
) *DtcMonitorIcmp {
	dtcMonitorIcmp := NewEmptyDtcMonitorIcmp()
	dtcMonitorIcmp.Comment = &comment
	dtcMonitorIcmp.Name = &name
	dtcMonitorIcmp.Ea = eas
	dtcMonitorIcmp.Interval = &interval
	dtcMonitorIcmp.RetryDown = &retry_down
	dtcMonitorIcmp.RetryUp = &retry_up
	dtcMonitorIcmp.Timeout = &timeout
	return dtcMonitorIcmp
}

func (objMgr *ObjectManager) CreateDtcMonitorHttp(
	comment, name string, eas EA, port, interval, retry_down, retry_up, timeout uint32,
	ciphers, client_cert, content_check, content_check_input, content_check_op, content_check_regex string,
	content_extract_group uint32, content_extract_type, content_extract_value string,
	enable_sni bool, request, result string, result_code uint32, secure, validate_cert bool,
) (*DtcMonitorHttp, error) {
	if name == "" {
		return nil, fmt.Errorf("name field is required to create a Dtc Monitor object")
	}
	dtcMonitorHttp := NewDtcMonitorHttp(
		comment,
		name,
		eas,
		port,
		interval,
		retry_down,
		retry_up,
		timeout,
		ciphers,
		client_cert,
		content_check,
		content_check_input,
		content_check_op,
		content_check_regex,
		content_extract_group,
		content_extract_type,
		content_extract_value,
		enable_sni,
		request,
		result,
		result_code,
		secure,
		validate_cert,
	)
	ref, err := objMgr.connector.CreateObject(dtcMonitorHttp)
	if err != nil {
		return nil, err
	}
	dtcMonitorHttp.Ref = ref
	return dtcMonitorHttp, nil
}

func (objMgr *ObjectManager) CreateDtcMonitorIcmp(
	comment, name string, eas EA,
	interval, retry_down, retry_up, timeout uint32,
) (*DtcMonitorIcmp, error) {
	if name == "" {
		return nil, fmt.Errorf("name field is required to create a Dtc Monitor object")
	}
	dtcMonitorIcmp := NewDtcMonitorIcmp(
		comment,
		name,
		eas,
		interval,
		retry_down,
		retry_up,
		timeout,
	)
	ref, err := objMgr.connector.CreateObject(dtcMonitorIcmp)
	if err != nil {
		return nil, err
	}
	dtcMonitorIcmp.Ref = ref
	return dtcMonitorIcmp, nil
}

func (objMgr *ObjectManager) GetDtcMonitorHttpByRef(ref string) (*DtcMonitorHttp, error) {
	dtcMonitorHttp := NewEmptyDtcMonitorHttp()
	err := objMgr.connector.GetObject(
		dtcMonitorHttp, ref, NewQueryParams(false, nil), &dtcMonitorHttp)
	return dtcMonitorHttp, err
}

func (objMgr *ObjectManager) GetDtcMonitorIcmpByRef(ref string) (*DtcMonitorIcmp, error) {
	dtcMonitorIcmp := NewEmptyDtcMonitorIcmp()
	err := objMgr.connector.GetObject(
		dtcMonitorIcmp, ref, NewQueryParams(false, nil), &dtcMonitorIcmp)
	return dtcMonitorIcmp, err
}

func (objMgr *ObjectManager) UpdateDtcMonitorHttp(ref string,
	comment, name string, eas EA, port, interval, retry_down, retry_up, timeout uint32,
	ciphers, client_cert, content_check, content_check_input, content_check_op, content_check_regex string,
	content_extract_group uint32, content_extract_type, content_extract_value string,
	enable_sni bool, request, result string, result_code uint32, secure, validate_cert bool,
) (*DtcMonitorHttp, error) {
	dtcMonitorHttp := NewDtcMonitorHttp(
		comment,
		name,
		eas,
		port,
		interval,
		retry_down,
		retry_up,
		timeout,
		ciphers,
		client_cert,
		content_check,
		content_check_input,
		content_check_op,
		content_check_regex,
		content_extract_group,
		content_extract_type,
		content_extract_value,
		enable_sni,
		request,
		result,
		result_code,
		secure,
		validate_cert,
	)
	dtcMonitorHttp.Ref = ref

	ref, err := objMgr.connector.UpdateObject(dtcMonitorHttp, ref)
	if err != nil {
		return nil, err
	}
	dtcMonitorHttp.Ref = ref
	return dtcMonitorHttp, nil
}

func (objMgr *ObjectManager) UpdateDtcMonitorIcmp(ref string,
	comment, name string, eas EA, interval, retry_down, retry_up, timeout uint32,
) (*DtcMonitorIcmp, error) {
	dtcMonitorIcmp := NewDtcMonitorIcmp(
		comment,
		name,
		eas,
		interval,
		retry_down,
		retry_up,
		timeout,
	)
	dtcMonitorIcmp.Ref = ref

	ref, err := objMgr.connector.UpdateObject(dtcMonitorIcmp, ref)
	if err != nil {
		return nil, err
	}
	dtcMonitorIcmp.Ref = ref
	return dtcMonitorIcmp, nil
}

func (objMgr *ObjectManager) DeleteDtcMonitor(ref string) (string, error) {
	return objMgr.connector.DeleteObject(ref)
}
