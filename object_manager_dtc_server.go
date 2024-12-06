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
		return nil, fmt.Errorf("if 'use_sni_hostname' is enabled then 'sni_hostname' must be provided or if 'sni_hostname' is provided then 'use_sni_hostname' must be enabled")
	}
	var serverMonitors []*DtcServerMonitor
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

		serverMonitors = append(serverMonitors, ServerMonitor)
	}
	DtcServer := NewDtcServer(comment, name, host, AutoCreateHostRecord, disable, ea, serverMonitors, SniHostname, UseSniHostname)
	ref, err := objMgr.connector.CreateObject(DtcServer)
	if err != nil {
		return nil, err
	}
	DtcServer.Ref = ref
	return DtcServer, nil
}

func (objMgr *ObjectManager) GetDtcServer(queryParams *QueryParams) (*DtcServer, error) {
	var res []DtcServer
	server := NewEmptyDtcServer()
	err := objMgr.connector.GetObject(server, "", queryParams, &res)
	if err != nil {
		return nil, fmt.Errorf("error getting DtcServer object, err: %s", err)
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
	var serverMonitors []*DtcServerMonitor
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

		serverMonitors = append(serverMonitors, ServerMonitor)
	}
	DtcServer := NewDtcServer(comment, name, host, AutoCreateHostRecord, disable, ea, serverMonitors, SniHostname, UseSniHostname)
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
