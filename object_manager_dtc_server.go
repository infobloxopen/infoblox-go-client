package ibclient

import "fmt"

func NewEmptyDtcServer() *DtcServer {
	DtcServer := &DtcServer{}
	DtcServer.SetReturnFields(append(DtcServer.ReturnFields(), "extattrs", "auto_create_host_record", "disable", "health", "monitors", "sni_hostname", "use_sni_hostname"))
	return DtcServer
}
func NewDtcServer(comment string,
	name string,
	host string,
	AutoCreateHostRecord bool,
	disable bool,
	ea EA,
	Monitors []*DtcServerMonitor,
	SniHostname string,
	UseSniHostname bool,
) *DtcServer {
	DtcServer := NewEmptyDtcServer()
	DtcServer.Comment = &comment
	DtcServer.Name = &name
	DtcServer.Host = &host
	DtcServer.AutoCreateHostRecord = &AutoCreateHostRecord
	DtcServer.Disable = &disable
	DtcServer.Ea = ea
	DtcServer.Monitors = Monitors
	DtcServer.SniHostname = &SniHostname
	DtcServer.UseSniHostname = &UseSniHostname
	return DtcServer
}

func (objMgr *ObjectManager) CreateDtcServer(
	comment string,
	name string,
	host string,
	AutoCreateHostRecord bool,
	disable bool,
	ea EA,
	Monitors []map[string]interface{},
	SniHostname string,
	UseSniHostname bool,
) (*DtcServer, error) {
	if (UseSniHostname && SniHostname == "") || (!UseSniHostname && SniHostname != "") {
		return nil, fmt.Errorf("If 'use_sni_hostname' is enabled then 'sni_hostname' must be provided or if 'sni_hostname' is provided then 'use_sni_hostname' must be enabled")
	}
	var Servermonitors []*DtcServerMonitor
	for _, userMonitor := range Monitors {
		monitor, okMonitor := userMonitor["monitor"].(Monitor)
		MonitorHost, _ := userMonitor["host"].(string)
		if !okMonitor {
			return nil, fmt.Errorf("\"Required field missing: monitor")
		}
		monitorRef, err := getMonitorReference(monitor.Name, monitor.Type, objMgr)
		if err != nil {
			return nil, err
		}

		ServerMonitor := &DtcServerMonitor{
			Monitor: monitorRef,
			Host:    MonitorHost,
		}

		Servermonitors = append(Servermonitors, ServerMonitor)
	}
	DtcServer := NewDtcServer(comment, name, host, AutoCreateHostRecord, disable, ea, Servermonitors, SniHostname, UseSniHostname)
	ref, err := objMgr.connector.CreateObject(DtcServer)
	if err != nil {
		return nil, err
	}
	DtcServer.Ref = ref
	return DtcServer, nil
}

func (objMgr *ObjectManager) GetDtcServer(serverName string, comment string, host string, sni_hostname string) (*DtcServer, error) {
	var res []DtcServer
	ServerDtc := NewEmptyDtcServer()
	sf := map[string]string{
		"name":         serverName,
		"comment":      comment,
		"host":         host,
		"sni_hostname": sni_hostname,
	}
	queryParams := NewQueryParams(false, sf)
	err := objMgr.connector.GetObject(ServerDtc, "", queryParams, &res)

	if err != nil {
		return nil, err
	} else if res == nil || len(res) == 0 {
		return nil, NewNotFoundError(
			fmt.Sprintf(
				"A Dtc Server with name '%s' is not found", serverName))
	}
	return &res[0], nil
}
func (objMgr *ObjectManager) UpdateDtcServer(
	ref string,
	comment string,
	name string,
	host string,
	AutoCreateHostRecord bool,
	disable bool,
	ea EA,
	Monitors []map[string]interface{},
	SniHostname string,
	UseSniHostname bool) (*DtcServer, error) {
	if (UseSniHostname && SniHostname == "") || (!UseSniHostname && SniHostname != "") {
		return nil, fmt.Errorf("If 'use_sni_hostname' is enabled then 'sni_hostname' must be provided or if 'sni_hostname' is provided then 'use_sni_hostname' must be enabled ")
	}
	var Servermonitors []*DtcServerMonitor
	for _, userMonitor := range Monitors {
		monitor, okMonitor := userMonitor["monitor"].(Monitor)
		MonitorHost, _ := userMonitor["host"].(string)
		if !okMonitor {
			return nil, fmt.Errorf("\"Required field missing: monitor")
		}
		monitorRef, err := getMonitorReference(monitor.Name, monitor.Type, objMgr)
		if err != nil {
			return nil, err
		}

		ServerMonitor := &DtcServerMonitor{
			Monitor: monitorRef,
			Host:    MonitorHost,
		}

		Servermonitors = append(Servermonitors, ServerMonitor)
	}
	DtcServer := NewDtcServer(comment, name, host, AutoCreateHostRecord, disable, ea, Servermonitors, SniHostname, UseSniHostname)
	DtcServer.Ref = ref
	ref, err := objMgr.connector.UpdateObject(DtcServer, ref)
	if err != nil {
		return nil, err
	}
	DtcServer.Ref = ref
	return DtcServer, nil

}
func (objMgr *ObjectManager) GetDtcServerByRef(ref string) (*DtcServer, error) {
	ServerDtc := NewEmptyDtcServer()
	err := objMgr.connector.GetObject(
		ServerDtc, ref, NewQueryParams(false, nil), &ServerDtc)
	return ServerDtc, err
}

func (objMgr *ObjectManager) DeleteDtcServer(ref string) (string, error) {
	return objMgr.connector.DeleteObject(ref)
}
