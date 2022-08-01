// Code generated by "infoblox-go-client-generator"; DO NOT EDIT.

package ibclient

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// This test checks if common methods for each
// generated Infoblox object works as expected
var _ = DescribeTable("Common methods of WAPI Objects",
	func(obj IBObject, expectedReturnFields []string, expectedObjectType string) {
		// Check if default return fields are valid
		Expect(obj.ReturnFields()).To(Equal(expectedReturnFields))
		// Check if ObjectType is valid
		Expect(obj.ObjectType()).To(Equal(expectedObjectType))
	},
	Entry(
		"AdAuthService",
		&AdAuthService{},
		[]string{"name"},
		"ad_auth_service",
	),
	Entry(
		"Admingroup",
		&Admingroup{},
		[]string{"comment", "name"},
		"admingroup",
	),
	Entry(
		"Adminrole",
		&Adminrole{},
		[]string{"comment", "name"},
		"adminrole",
	),
	Entry(
		"Adminuser",
		&Adminuser{},
		[]string{"admin_groups", "comment", "name"},
		"adminuser",
	),
	Entry(
		"Allendpoints",
		&Allendpoints{},
		[]string{},
		"allendpoints",
	),
	Entry(
		"Allnsgroup",
		&Allnsgroup{},
		[]string{"name", "type"},
		"allnsgroup",
	),
	Entry(
		"Allrecords",
		&Allrecords{},
		[]string{"comment", "name", "type", "view", "zone"},
		"allrecords",
	),
	Entry(
		"Allrpzrecords",
		&Allrpzrecords{},
		[]string{"comment", "name", "type", "view", "zone"},
		"allrpzrecords",
	),
	Entry(
		"Approvalworkflow",
		&Approvalworkflow{},
		[]string{"approval_group", "submitter_group"},
		"approvalworkflow",
	),
	Entry(
		"Authpolicy",
		&Authpolicy{},
		[]string{"default_group", "usage_type"},
		"authpolicy",
	),
	Entry(
		"Awsrte53taskgroup",
		&Awsrte53taskgroup{},
		[]string{"account_id", "comment", "disabled", "name", "sync_status"},
		"awsrte53taskgroup",
	),
	Entry(
		"Awsuser",
		&Awsuser{},
		[]string{"access_key_id", "account_id", "name"},
		"awsuser",
	),
	Entry(
		"Bfdtemplate",
		&Bfdtemplate{},
		[]string{"name"},
		"bfdtemplate",
	),
	Entry(
		"Bulkhost",
		&Bulkhost{},
		[]string{"comment", "prefix"},
		"bulkhost",
	),
	Entry(
		"Bulkhostnametemplate",
		&Bulkhostnametemplate{},
		[]string{"is_grid_default", "template_format", "template_name"},
		"bulkhostnametemplate",
	),
	Entry(
		"Cacertificate",
		&Cacertificate{},
		[]string{"distinguished_name", "issuer", "serial", "used_by", "valid_not_after", "valid_not_before"},
		"cacertificate",
	),
	Entry(
		"CapacityReport",
		&CapacityReport{},
		[]string{"name", "percent_used", "role"},
		"capacityreport",
	),
	Entry(
		"Captiveportal",
		&Captiveportal{},
		[]string{"name"},
		"captiveportal",
	),
	Entry(
		"CertificateAuthservice",
		&CertificateAuthservice{},
		[]string{"name"},
		"certificate:authservice",
	),
	Entry(
		"CiscoiseEndpoint",
		&CiscoiseEndpoint{},
		[]string{"address", "disable", "resolved_address", "type", "version"},
		"ciscoise:endpoint",
	),
	Entry(
		"Csvimporttask",
		&Csvimporttask{},
		[]string{"action", "admin_name", "end_time", "file_name", "file_size", "import_id", "lines_failed", "lines_processed", "lines_warning", "on_error", "operation", "separator", "start_time", "status", "update_method"},
		"csvimporttask",
	),
	Entry(
		"DbObjects",
		&DbObjects{},
		[]string{"last_sequence_id", "object", "object_type", "unique_id"},
		"db_objects",
	),
	Entry(
		"Dbsnapshot",
		&Dbsnapshot{},
		[]string{"comment", "timestamp"},
		"dbsnapshot",
	),
	Entry(
		"DdnsPrincipalcluster",
		&DdnsPrincipalcluster{},
		[]string{"comment", "group", "name", "principals"},
		"ddns:principalcluster",
	),
	Entry(
		"DdnsPrincipalclusterGroup",
		&DdnsPrincipalclusterGroup{},
		[]string{"comment", "name"},
		"ddns:principalcluster:group",
	),
	Entry(
		"DeletedObjects",
		&DeletedObjects{},
		[]string{"object_type"},
		"deleted_objects",
	),
	Entry(
		"DhcpStatistics",
		&DhcpStatistics{},
		[]string{"dhcp_utilization", "dhcp_utilization_status", "dynamic_hosts", "static_hosts", "total_hosts"},
		"dhcp:statistics",
	),
	Entry(
		"Dhcpfailover",
		&Dhcpfailover{},
		[]string{"name"},
		"dhcpfailover",
	),
	Entry(
		"Dhcpoptiondefinition",
		&Dhcpoptiondefinition{},
		[]string{"code", "name", "type"},
		"dhcpoptiondefinition",
	),
	Entry(
		"Dhcpoptionspace",
		&Dhcpoptionspace{},
		[]string{"comment", "name"},
		"dhcpoptionspace",
	),
	Entry(
		"Discovery",
		&Discovery{},
		[]string{},
		"discovery",
	),
	Entry(
		"DiscoveryCredentialgroup",
		&DiscoveryCredentialgroup{},
		[]string{"name"},
		"discovery:credentialgroup",
	),
	Entry(
		"DiscoveryDevice",
		&DiscoveryDevice{},
		[]string{"address", "name", "network_view"},
		"discovery:device",
	),
	Entry(
		"DiscoveryDevicecomponent",
		&DiscoveryDevicecomponent{},
		[]string{"component_name", "description", "model", "serial", "type"},
		"discovery:devicecomponent",
	),
	Entry(
		"DiscoveryDeviceinterface",
		&DiscoveryDeviceinterface{},
		[]string{"name", "type"},
		"discovery:deviceinterface",
	),
	Entry(
		"DiscoveryDeviceneighbor",
		&DiscoveryDeviceneighbor{},
		[]string{"address", "address_ref", "mac", "name"},
		"discovery:deviceneighbor",
	),
	Entry(
		"DiscoveryDevicesupportbundle",
		&DiscoveryDevicesupportbundle{},
		[]string{"author", "integrated_ind", "name", "version"},
		"discovery:devicesupportbundle",
	),
	Entry(
		"DiscoveryDiagnostictask",
		&DiscoveryDiagnostictask{},
		[]string{"ip_address", "network_view", "task_id"},
		"discovery:diagnostictask",
	),
	Entry(
		"DiscoveryGridproperties",
		&DiscoveryGridproperties{},
		[]string{"grid_name"},
		"discovery:gridproperties",
	),
	Entry(
		"DiscoveryMemberproperties",
		&DiscoveryMemberproperties{},
		[]string{"discovery_member"},
		"discovery:memberproperties",
	),
	Entry(
		"DiscoverySdnnetwork",
		&DiscoverySdnnetwork{},
		[]string{"name", "network_view", "source_sdn_config"},
		"discovery:sdnnetwork",
	),
	Entry(
		"DiscoveryStatus",
		&DiscoveryStatus{},
		[]string{"address", "name", "network_view", "status"},
		"discovery:status",
	),
	Entry(
		"DiscoveryVrf",
		&DiscoveryVrf{},
		[]string{"device", "name", "network_view", "route_distinguisher"},
		"discovery:vrf",
	),
	Entry(
		"Discoverytask",
		&Discoverytask{},
		[]string{"discovery_task_oid", "member_name"},
		"discoverytask",
	),
	Entry(
		"Distributionschedule",
		&Distributionschedule{},
		[]string{"active", "start_time", "time_zone"},
		"distributionschedule",
	),
	Entry(
		"Dns64group",
		&Dns64group{},
		[]string{"comment", "disable", "name"},
		"dns64group",
	),
	Entry(
		"Dtc",
		&Dtc{},
		[]string{},
		"dtc",
	),
	Entry(
		"DtcAllrecords",
		&DtcAllrecords{},
		[]string{"comment", "dtc_server", "type"},
		"dtc:allrecords",
	),
	Entry(
		"DtcCertificate",
		&DtcCertificate{},
		[]string{},
		"dtc:certificate",
	),
	Entry(
		"DtcLbdn",
		&DtcLbdn{},
		[]string{"comment", "name"},
		"dtc:lbdn",
	),
	Entry(
		"DtcMonitor",
		&DtcMonitor{},
		[]string{"comment", "name", "type"},
		"dtc:monitor",
	),
	Entry(
		"DtcMonitorHttp",
		&DtcMonitorHttp{},
		[]string{"comment", "name"},
		"dtc:monitor:http",
	),
	Entry(
		"DtcMonitorIcmp",
		&DtcMonitorIcmp{},
		[]string{"comment", "name"},
		"dtc:monitor:icmp",
	),
	Entry(
		"DtcMonitorPdp",
		&DtcMonitorPdp{},
		[]string{"comment", "name"},
		"dtc:monitor:pdp",
	),
	Entry(
		"DtcMonitorSip",
		&DtcMonitorSip{},
		[]string{"comment", "name"},
		"dtc:monitor:sip",
	),
	Entry(
		"DtcMonitorSnmp",
		&DtcMonitorSnmp{},
		[]string{"comment", "name"},
		"dtc:monitor:snmp",
	),
	Entry(
		"DtcMonitorTcp",
		&DtcMonitorTcp{},
		[]string{"comment", "name"},
		"dtc:monitor:tcp",
	),
	Entry(
		"DtcObject",
		&DtcObject{},
		[]string{"abstract_type", "comment", "display_type", "name", "status"},
		"dtc:object",
	),
	Entry(
		"DtcPool",
		&DtcPool{},
		[]string{"comment", "name"},
		"dtc:pool",
	),
	Entry(
		"DtcRecordA",
		&DtcRecordA{},
		[]string{"dtc_server", "ipv4addr"},
		"dtc:record:a",
	),
	Entry(
		"DtcRecordAaaa",
		&DtcRecordAaaa{},
		[]string{"dtc_server", "ipv6addr"},
		"dtc:record:aaaa",
	),
	Entry(
		"DtcRecordCname",
		&DtcRecordCname{},
		[]string{"canonical", "dtc_server"},
		"dtc:record:cname",
	),
	Entry(
		"DtcRecordNaptr",
		&DtcRecordNaptr{},
		[]string{"dtc_server", "order", "preference", "regexp", "replacement", "services"},
		"dtc:record:naptr",
	),
	Entry(
		"DtcRecordSrv",
		&DtcRecordSrv{},
		[]string{"dtc_server", "name", "port", "priority", "target", "weight"},
		"dtc:record:srv",
	),
	Entry(
		"DtcServer",
		&DtcServer{},
		[]string{"comment", "host", "name"},
		"dtc:server",
	),
	Entry(
		"DtcTopology",
		&DtcTopology{},
		[]string{"comment", "name"},
		"dtc:topology",
	),
	Entry(
		"DtcTopologyLabel",
		&DtcTopologyLabel{},
		[]string{"field", "label"},
		"dtc:topology:label",
	),
	Entry(
		"DtcTopologyRule",
		&DtcTopologyRule{},
		[]string{},
		"dtc:topology:rule",
	),
	Entry(
		"DxlEndpoint",
		&DxlEndpoint{},
		[]string{"disable", "name", "outbound_member_type"},
		"dxl:endpoint",
	),
	Entry(
		"EADefinition",
		&EADefinition{},
		[]string{"comment", "default_value", "name", "type"},
		"extensibleattributedef",
	),
	Entry(
		"Fileop",
		&Fileop{},
		[]string{},
		"fileop",
	),
	Entry(
		"Filterfingerprint",
		&Filterfingerprint{},
		[]string{"comment", "name"},
		"filterfingerprint",
	),
	Entry(
		"Filtermac",
		&Filtermac{},
		[]string{"comment", "name"},
		"filtermac",
	),
	Entry(
		"Filternac",
		&Filternac{},
		[]string{"comment", "name"},
		"filternac",
	),
	Entry(
		"Filteroption",
		&Filteroption{},
		[]string{"comment", "name"},
		"filteroption",
	),
	Entry(
		"Filterrelayagent",
		&Filterrelayagent{},
		[]string{"comment", "name"},
		"filterrelayagent",
	),
	Entry(
		"Fingerprint",
		&Fingerprint{},
		[]string{"comment", "device_class", "name"},
		"fingerprint",
	),
	Entry(
		"Ipv4FixedAddress",
		&Ipv4FixedAddress{},
		[]string{"ipv4addr", "network_view"},
		"fixedaddress",
	),
	Entry(
		"Fixedaddresstemplate",
		&Fixedaddresstemplate{},
		[]string{"comment", "name"},
		"fixedaddresstemplate",
	),
	Entry(
		"Ftpuser",
		&Ftpuser{},
		[]string{"username"},
		"ftpuser",
	),
	Entry(
		"Grid",
		&Grid{},
		[]string{},
		"grid",
	),
	Entry(
		"GridCloudapi",
		&GridCloudapi{},
		[]string{"allow_api_admins", "allowed_api_admins", "enable_recycle_bin"},
		"grid:cloudapi",
	),
	Entry(
		"GridCloudapiCloudstatistics",
		&GridCloudapiCloudstatistics{},
		[]string{"allocated_available_ratio", "allocated_ip_count", "available_ip_count", "fixed_ip_count", "floating_ip_count", "tenant_count", "tenant_ip_count", "tenant_vm_count"},
		"grid:cloudapi:cloudstatistics",
	),
	Entry(
		"GridCloudapiTenant",
		&GridCloudapiTenant{},
		[]string{"comment", "id", "name"},
		"grid:cloudapi:tenant",
	),
	Entry(
		"GridCloudapiVm",
		&GridCloudapiVm{},
		[]string{"comment", "id", "name"},
		"grid:cloudapi:vm",
	),
	Entry(
		"GridCloudapiVmaddress",
		&GridCloudapiVmaddress{},
		[]string{"address", "is_ipv4", "network_view", "port_id", "vm_name"},
		"grid:cloudapi:vmaddress",
	),
	Entry(
		"GridDashboard",
		&GridDashboard{},
		[]string{"analytics_tunneling_event_critical_threshold", "analytics_tunneling_event_warning_threshold", "atp_critical_event_critical_threshold", "atp_critical_event_warning_threshold", "atp_major_event_critical_threshold", "atp_major_event_warning_threshold", "atp_warning_event_critical_threshold", "atp_warning_event_warning_threshold", "rpz_blocked_hit_critical_threshold", "rpz_blocked_hit_warning_threshold", "rpz_passthru_event_critical_threshold", "rpz_passthru_event_warning_threshold", "rpz_substituted_hit_critical_threshold", "rpz_substituted_hit_warning_threshold"},
		"grid:dashboard",
	),
	Entry(
		"GridDhcpproperties",
		&GridDhcpproperties{},
		[]string{"disable_all_nac_filters", "grid"},
		"grid:dhcpproperties",
	),
	Entry(
		"GridDns",
		&GridDns{},
		[]string{},
		"grid:dns",
	),
	Entry(
		"GridFiledistribution",
		&GridFiledistribution{},
		[]string{"allow_uploads", "current_usage", "global_status", "name", "storage_limit"},
		"grid:filedistribution",
	),
	Entry(
		"GridLicensePool",
		&GridLicensePool{},
		[]string{"type"},
		"grid:license_pool",
	),
	Entry(
		"GridLicensePoolContainer",
		&GridLicensePoolContainer{},
		[]string{},
		"grid:license_pool_container",
	),
	Entry(
		"GridMaxminddbinfo",
		&GridMaxminddbinfo{},
		[]string{"binary_major_version", "binary_minor_version", "build_time", "database_type", "deployment_time", "member", "topology_type"},
		"grid:maxminddbinfo",
	),
	Entry(
		"GridMemberCloudapi",
		&GridMemberCloudapi{},
		[]string{"allow_api_admins", "allowed_api_admins", "enable_service", "member", "status"},
		"grid:member:cloudapi",
	),
	Entry(
		"GridServicerestartGroup",
		&GridServicerestartGroup{},
		[]string{"comment", "name", "service"},
		"grid:servicerestart:group",
	),
	Entry(
		"GridServicerestartGroupOrder",
		&GridServicerestartGroupOrder{},
		[]string{},
		"grid:servicerestart:group:order",
	),
	Entry(
		"GridServicerestartRequest",
		&GridServicerestartRequest{},
		[]string{"error", "group", "result", "state"},
		"grid:servicerestart:request",
	),
	Entry(
		"GridServicerestartRequestChangedobject",
		&GridServicerestartRequestChangedobject{},
		[]string{"action", "changed_properties", "changed_time", "object_name", "object_type", "user_name"},
		"grid:servicerestart:request:changedobject",
	),
	Entry(
		"GridServicerestartStatus",
		&GridServicerestartStatus{},
		[]string{"failures", "finished", "grouped", "needed_restart", "no_restart", "parent", "pending", "pending_restart", "processing", "restarting", "success", "timeouts"},
		"grid:servicerestart:status",
	),
	Entry(
		"GridThreatanalytics",
		&GridThreatanalytics{},
		[]string{"enable_auto_download", "enable_scheduled_download", "module_update_policy", "name"},
		"grid:threatanalytics",
	),
	Entry(
		"GridThreatprotection",
		&GridThreatprotection{},
		[]string{"grid_name"},
		"grid:threatprotection",
	),
	Entry(
		"GridX509certificate",
		&GridX509certificate{},
		[]string{"issuer", "serial", "subject"},
		"grid:x509certificate",
	),
	Entry(
		"Hostnamerewritepolicy",
		&Hostnamerewritepolicy{},
		[]string{"name", "replacement_character", "valid_characters"},
		"hostnamerewritepolicy",
	),
	Entry(
		"HsmAllgroups",
		&HsmAllgroups{},
		[]string{"groups"},
		"hsm:allgroups",
	),
	Entry(
		"HsmSafenetgroup",
		&HsmSafenetgroup{},
		[]string{"comment", "hsm_version", "name"},
		"hsm:safenetgroup",
	),
	Entry(
		"HsmThalesgroup",
		&HsmThalesgroup{},
		[]string{"comment", "key_server_ip", "name"},
		"hsm:thalesgroup",
	),
	Entry(
		"IpamStatistics",
		&IpamStatistics{},
		[]string{"cidr", "network", "network_view"},
		"ipam:statistics",
	),
	Entry(
		"IPv4Address",
		&IPv4Address{},
		[]string{"dhcp_client_identifier", "ip_address", "is_conflict", "lease_state", "mac_address", "names", "network", "network_view", "objects", "status", "types", "usage", "username"},
		"ipv4address",
	),
	Entry(
		"IPv6Address",
		&IPv6Address{},
		[]string{"duid", "ip_address", "is_conflict", "lease_state", "names", "network", "network_view", "objects", "status", "types", "usage"},
		"ipv6address",
	),
	Entry(
		"Ipv6dhcpoptiondefinition",
		&Ipv6dhcpoptiondefinition{},
		[]string{"code", "name", "type"},
		"ipv6dhcpoptiondefinition",
	),
	Entry(
		"Ipv6dhcpoptionspace",
		&Ipv6dhcpoptionspace{},
		[]string{"comment", "enterprise_number", "name"},
		"ipv6dhcpoptionspace",
	),
	Entry(
		"Ipv6filteroption",
		&Ipv6filteroption{},
		[]string{"comment", "name"},
		"ipv6filteroption",
	),
	Entry(
		"Ipv6FixedAddress",
		&Ipv6FixedAddress{},
		[]string{"duid", "ipv6addr", "network_view"},
		"ipv6fixedaddress",
	),
	Entry(
		"Ipv6fixedaddresstemplate",
		&Ipv6fixedaddresstemplate{},
		[]string{"comment", "name"},
		"ipv6fixedaddresstemplate",
	),
	Entry(
		"Ipv6Network",
		&Ipv6Network{},
		[]string{"comment", "network", "network_view"},
		"ipv6network",
	),
	Entry(
		"Ipv6NetworkContainer",
		&Ipv6NetworkContainer{},
		[]string{"comment", "network", "network_view"},
		"ipv6networkcontainer",
	),
	Entry(
		"IPv6NetworkTemplate",
		&IPv6NetworkTemplate{},
		[]string{"comment", "name"},
		"ipv6networktemplate",
	),
	Entry(
		"IPv6Range",
		&IPv6Range{},
		[]string{"comment", "end_addr", "network", "network_view", "start_addr"},
		"ipv6range",
	),
	Entry(
		"Ipv6rangetemplate",
		&Ipv6rangetemplate{},
		[]string{"comment", "name", "number_of_addresses", "offset"},
		"ipv6rangetemplate",
	),
	Entry(
		"IPv6SharedNetwork",
		&IPv6SharedNetwork{},
		[]string{"comment", "name", "network_view", "networks"},
		"ipv6sharednetwork",
	),
	Entry(
		"Kerberoskey",
		&Kerberoskey{},
		[]string{"domain", "enctype", "in_use", "principal", "version"},
		"kerberoskey",
	),
	Entry(
		"LdapAuthService",
		&LdapAuthService{},
		[]string{"comment", "disable", "ldap_user_attribute", "mode", "name"},
		"ldap_auth_service",
	),
	Entry(
		"Lease",
		&Lease{},
		[]string{"address", "network_view"},
		"lease",
	),
	Entry(
		"LicenseGridwide",
		&LicenseGridwide{},
		[]string{"type"},
		"license:gridwide",
	),
	Entry(
		"LocaluserAuthservice",
		&LocaluserAuthservice{},
		[]string{"comment", "disabled", "name"},
		"localuser:authservice",
	),
	Entry(
		"MACFilterAddress",
		&MACFilterAddress{},
		[]string{"authentication_time", "comment", "expiration_time", "filter", "guest_custom_field1", "guest_custom_field2", "guest_custom_field3", "guest_custom_field4", "guest_email", "guest_first_name", "guest_last_name", "guest_middle_name", "guest_phone", "is_registered_user", "mac", "never_expires", "reserved_for_infoblox", "username"},
		"macfilteraddress",
	),
	Entry(
		"Mastergrid",
		&Mastergrid{},
		[]string{"address", "enable", "port"},
		"mastergrid",
	),
	Entry(
		"Member",
		&Member{},
		[]string{"config_addr_type", "host_name", "platform", "service_type_configuration"},
		"member",
	),
	Entry(
		"MemberDHCPProperties",
		&MemberDHCPProperties{},
		[]string{"host_name", "ipv4addr", "ipv6addr"},
		"member:dhcpproperties",
	),
	Entry(
		"MemberDns",
		&MemberDns{},
		[]string{"host_name", "ipv4addr", "ipv6addr"},
		"member:dns",
	),
	Entry(
		"MemberFiledistribution",
		&MemberFiledistribution{},
		[]string{"host_name", "ipv4_address", "ipv6_address", "status"},
		"member:filedistribution",
	),
	Entry(
		"MemberLicense",
		&MemberLicense{},
		[]string{"type"},
		"member:license",
	),
	Entry(
		"MemberParentalcontrol",
		&MemberParentalcontrol{},
		[]string{"enable_service", "name"},
		"member:parentalcontrol",
	),
	Entry(
		"MemberThreatanalytics",
		&MemberThreatanalytics{},
		[]string{"host_name", "ipv4_address", "ipv6_address", "status"},
		"member:threatanalytics",
	),
	Entry(
		"MemberThreatprotection",
		&MemberThreatprotection{},
		[]string{},
		"member:threatprotection",
	),
	Entry(
		"Memberdfp",
		&Memberdfp{},
		[]string{},
		"memberdfp",
	),
	Entry(
		"Msserver",
		&Msserver{},
		[]string{"address"},
		"msserver",
	),
	Entry(
		"MsserverAdsitesDomain",
		&MsserverAdsitesDomain{},
		[]string{"name", "netbios", "network_view"},
		"msserver:adsites:domain",
	),
	Entry(
		"MsserverAdsitesSite",
		&MsserverAdsitesSite{},
		[]string{"domain", "name"},
		"msserver:adsites:site",
	),
	Entry(
		"MsserverDhcp",
		&MsserverDhcp{},
		[]string{"address"},
		"msserver:dhcp",
	),
	Entry(
		"MsserverDns",
		&MsserverDns{},
		[]string{"address"},
		"msserver:dns",
	),
	Entry(
		"Mssuperscope",
		&Mssuperscope{},
		[]string{"disable", "name", "network_view"},
		"mssuperscope",
	),
	Entry(
		"Namedacl",
		&Namedacl{},
		[]string{"comment", "name"},
		"namedacl",
	),
	Entry(
		"Natgroup",
		&Natgroup{},
		[]string{"comment", "name"},
		"natgroup",
	),
	Entry(
		"Ipv4Network",
		&Ipv4Network{},
		[]string{"comment", "network", "network_view"},
		"network",
	),
	Entry(
		"NetworkDiscovery",
		&NetworkDiscovery{},
		[]string{},
		"network_discovery",
	),
	Entry(
		"Ipv4NetworkContainer",
		&Ipv4NetworkContainer{},
		[]string{"comment", "network", "network_view"},
		"networkcontainer",
	),
	Entry(
		"NetworkTemplate",
		&NetworkTemplate{},
		[]string{"comment", "name"},
		"networktemplate",
	),
	Entry(
		"Networkuser",
		&Networkuser{},
		[]string{"address", "domainname", "name", "network_view", "user_status"},
		"networkuser",
	),
	Entry(
		"NetworkView",
		&NetworkView{},
		[]string{"comment", "is_default", "name"},
		"networkview",
	),
	Entry(
		"NotificationRestEndpoint",
		&NotificationRestEndpoint{},
		[]string{"name", "outbound_member_type", "uri"},
		"notification:rest:endpoint",
	),
	Entry(
		"NotificationRestTemplate",
		&NotificationRestTemplate{},
		[]string{"content", "name"},
		"notification:rest:template",
	),
	Entry(
		"NotificationRule",
		&NotificationRule{},
		[]string{"event_type", "name", "notification_action", "notification_target"},
		"notification:rule",
	),
	Entry(
		"Nsgroup",
		&Nsgroup{},
		[]string{"comment", "name"},
		"nsgroup",
	),
	Entry(
		"NsgroupDelegation",
		&NsgroupDelegation{},
		[]string{"delegate_to", "name"},
		"nsgroup:delegation",
	),
	Entry(
		"NsgroupForwardingmember",
		&NsgroupForwardingmember{},
		[]string{"forwarding_servers", "name"},
		"nsgroup:forwardingmember",
	),
	Entry(
		"NsgroupForwardstubserver",
		&NsgroupForwardstubserver{},
		[]string{"external_servers", "name"},
		"nsgroup:forwardstubserver",
	),
	Entry(
		"NsgroupStubmember",
		&NsgroupStubmember{},
		[]string{"name"},
		"nsgroup:stubmember",
	),
	Entry(
		"Orderedranges",
		&Orderedranges{},
		[]string{"network", "ranges"},
		"orderedranges",
	),
	Entry(
		"Orderedresponsepolicyzones",
		&Orderedresponsepolicyzones{},
		[]string{"view"},
		"orderedresponsepolicyzones",
	),
	Entry(
		"OutboundCloudclient",
		&OutboundCloudclient{},
		[]string{"enable", "interval"},
		"outbound:cloudclient",
	),
	Entry(
		"ParentalcontrolAvp",
		&ParentalcontrolAvp{},
		[]string{"name", "type", "value_type"},
		"parentalcontrol:avp",
	),
	Entry(
		"ParentalcontrolBlockingpolicy",
		&ParentalcontrolBlockingpolicy{},
		[]string{"name", "value"},
		"parentalcontrol:blockingpolicy",
	),
	Entry(
		"ParentalcontrolSubscriber",
		&ParentalcontrolSubscriber{},
		[]string{"alt_subscriber_id", "local_id", "subscriber_id"},
		"parentalcontrol:subscriber",
	),
	Entry(
		"ParentalcontrolSubscriberrecord",
		&ParentalcontrolSubscriberrecord{},
		[]string{"accounting_session_id", "ip_addr", "ipsd", "localid", "prefix", "site", "subscriber_id"},
		"parentalcontrol:subscriberrecord",
	),
	Entry(
		"ParentalcontrolSubscribersite",
		&ParentalcontrolSubscribersite{},
		[]string{"block_size", "dca_sub_bw_list", "dca_sub_query_count", "first_port", "name", "stop_anycast", "strict_nat"},
		"parentalcontrol:subscribersite",
	),
	Entry(
		"Permission",
		&Permission{},
		[]string{"group", "permission", "resource_type", "role"},
		"permission",
	),
	Entry(
		"PxgridEndpoint",
		&PxgridEndpoint{},
		[]string{"address", "disable", "name", "outbound_member_type"},
		"pxgrid:endpoint",
	),
	Entry(
		"RadiusAuthservice",
		&RadiusAuthservice{},
		[]string{"comment", "disable", "name"},
		"radius:authservice",
	),
	Entry(
		"Range",
		&Range{},
		[]string{"comment", "end_addr", "network", "network_view", "start_addr"},
		"range",
	),
	Entry(
		"Rangetemplate",
		&Rangetemplate{},
		[]string{"comment", "name", "number_of_addresses", "offset"},
		"rangetemplate",
	),
	Entry(
		"RecordA",
		&RecordA{},
		[]string{"ipv4addr", "name", "view"},
		"record:a",
	),
	Entry(
		"RecordAAAA",
		&RecordAAAA{},
		[]string{"ipv6addr", "name", "view"},
		"record:aaaa",
	),
	Entry(
		"RecordAlias",
		&RecordAlias{},
		[]string{"name", "target_name", "target_type", "view"},
		"record:alias",
	),
	Entry(
		"RecordCaa",
		&RecordCaa{},
		[]string{"name", "view"},
		"record:caa",
	),
	Entry(
		"RecordCNAME",
		&RecordCNAME{},
		[]string{"canonical", "name", "view"},
		"record:cname",
	),
	Entry(
		"RecordDhcid",
		&RecordDhcid{},
		[]string{"name", "view"},
		"record:dhcid",
	),
	Entry(
		"RecordDname",
		&RecordDname{},
		[]string{"name", "target", "view"},
		"record:dname",
	),
	Entry(
		"RecordDnskey",
		&RecordDnskey{},
		[]string{"name", "view"},
		"record:dnskey",
	),
	Entry(
		"RecordDs",
		&RecordDs{},
		[]string{"name", "view"},
		"record:ds",
	),
	Entry(
		"RecordDtclbdn",
		&RecordDtclbdn{},
		[]string{"comment", "name", "view", "zone"},
		"record:dtclbdn",
	),
	Entry(
		"HostRecord",
		&HostRecord{},
		[]string{"ipv4addrs", "ipv6addrs", "name", "view"},
		"record:host",
	),
	Entry(
		"HostRecordIpv4Addr",
		&HostRecordIpv4Addr{},
		[]string{"configure_for_dhcp", "host", "ipv4addr", "mac"},
		"record:host_ipv4addr",
	),
	Entry(
		"HostRecordIpv6Addr",
		&HostRecordIpv6Addr{},
		[]string{"configure_for_dhcp", "duid", "host", "ipv6addr"},
		"record:host_ipv6addr",
	),
	Entry(
		"RecordMX",
		&RecordMX{},
		[]string{"mail_exchanger", "name", "preference", "view"},
		"record:mx",
	),
	Entry(
		"RecordNaptr",
		&RecordNaptr{},
		[]string{"name", "order", "preference", "regexp", "replacement", "services", "view"},
		"record:naptr",
	),
	Entry(
		"RecordNS",
		&RecordNS{},
		[]string{"name", "nameserver", "view"},
		"record:ns",
	),
	Entry(
		"RecordNsec",
		&RecordNsec{},
		[]string{"name", "view"},
		"record:nsec",
	),
	Entry(
		"RecordNsec3",
		&RecordNsec3{},
		[]string{"name", "view"},
		"record:nsec3",
	),
	Entry(
		"RecordNsec3param",
		&RecordNsec3param{},
		[]string{"name", "view"},
		"record:nsec3param",
	),
	Entry(
		"RecordPTR",
		&RecordPTR{},
		[]string{"ptrdname", "view"},
		"record:ptr",
	),
	Entry(
		"RecordRpzA",
		&RecordRpzA{},
		[]string{"ipv4addr", "name", "view"},
		"record:rpz:a",
	),
	Entry(
		"RecordRpzAIpaddress",
		&RecordRpzAIpaddress{},
		[]string{"ipv4addr", "name", "view"},
		"record:rpz:a:ipaddress",
	),
	Entry(
		"RecordRpzAaaa",
		&RecordRpzAaaa{},
		[]string{"ipv6addr", "name", "view"},
		"record:rpz:aaaa",
	),
	Entry(
		"RecordRpzAaaaIpaddress",
		&RecordRpzAaaaIpaddress{},
		[]string{"ipv6addr", "name", "view"},
		"record:rpz:aaaa:ipaddress",
	),
	Entry(
		"RecordRpzCname",
		&RecordRpzCname{},
		[]string{"canonical", "name", "view"},
		"record:rpz:cname",
	),
	Entry(
		"RecordRpzCnameClientipaddress",
		&RecordRpzCnameClientipaddress{},
		[]string{"canonical", "name", "view"},
		"record:rpz:cname:clientipaddress",
	),
	Entry(
		"RecordRpzCnameClientipaddressdn",
		&RecordRpzCnameClientipaddressdn{},
		[]string{"canonical", "name", "view"},
		"record:rpz:cname:clientipaddressdn",
	),
	Entry(
		"RecordRpzCnameIpaddress",
		&RecordRpzCnameIpaddress{},
		[]string{"canonical", "name", "view"},
		"record:rpz:cname:ipaddress",
	),
	Entry(
		"RecordRpzCnameIpaddressdn",
		&RecordRpzCnameIpaddressdn{},
		[]string{"canonical", "name", "view"},
		"record:rpz:cname:ipaddressdn",
	),
	Entry(
		"RecordRpzMx",
		&RecordRpzMx{},
		[]string{"mail_exchanger", "name", "preference", "view"},
		"record:rpz:mx",
	),
	Entry(
		"RecordRpzNaptr",
		&RecordRpzNaptr{},
		[]string{"name", "order", "preference", "regexp", "replacement", "services", "view"},
		"record:rpz:naptr",
	),
	Entry(
		"RecordRpzPtr",
		&RecordRpzPtr{},
		[]string{"ptrdname", "view"},
		"record:rpz:ptr",
	),
	Entry(
		"RecordRpzSrv",
		&RecordRpzSrv{},
		[]string{"name", "port", "priority", "target", "view", "weight"},
		"record:rpz:srv",
	),
	Entry(
		"RecordRpzTxt",
		&RecordRpzTxt{},
		[]string{"name", "text", "view"},
		"record:rpz:txt",
	),
	Entry(
		"RecordRrsig",
		&RecordRrsig{},
		[]string{"name", "view"},
		"record:rrsig",
	),
	Entry(
		"RecordSRV",
		&RecordSRV{},
		[]string{"name", "port", "priority", "target", "view", "weight"},
		"record:srv",
	),
	Entry(
		"RecordTlsa",
		&RecordTlsa{},
		[]string{"name", "view"},
		"record:tlsa",
	),
	Entry(
		"RecordTXT",
		&RecordTXT{},
		[]string{"name", "text", "view"},
		"record:txt",
	),
	Entry(
		"RecordUnknown",
		&RecordUnknown{},
		[]string{"name", "view"},
		"record:unknown",
	),
	Entry(
		"Recordnamepolicy",
		&Recordnamepolicy{},
		[]string{"is_default", "name", "regex"},
		"recordnamepolicy",
	),
	Entry(
		"Restartservicestatus",
		&Restartservicestatus{},
		[]string{"dhcp_status", "dns_status", "member", "reporting_status"},
		"restartservicestatus",
	),
	Entry(
		"Rir",
		&Rir{},
		[]string{"communication_mode", "email", "name", "url"},
		"rir",
	),
	Entry(
		"RirOrganization",
		&RirOrganization{},
		[]string{"id", "maintainer", "name", "rir", "sender_email"},
		"rir:organization",
	),
	Entry(
		"RoamingHost",
		&RoamingHost{},
		[]string{"address_type", "name", "network_view"},
		"roaminghost",
	),
	Entry(
		"Ruleset",
		&Ruleset{},
		[]string{"comment", "disabled", "name", "type"},
		"ruleset",
	),
	Entry(
		"SamlAuthservice",
		&SamlAuthservice{},
		[]string{"name"},
		"saml:authservice",
	),
	Entry(
		"Scavengingtask",
		&Scavengingtask{},
		[]string{"action", "associated_object", "status"},
		"scavengingtask",
	),
	Entry(
		"ScheduledTask",
		&ScheduledTask{},
		[]string{"approval_status", "execution_status", "task_id"},
		"scheduledtask",
	),
	Entry(
		"Search",
		&Search{},
		[]string{},
		"search",
	),
	Entry(
		"SharedNetwork",
		&SharedNetwork{},
		[]string{"comment", "name", "network_view", "networks"},
		"sharednetwork",
	),
	Entry(
		"SharedRecordA",
		&SharedRecordA{},
		[]string{"ipv4addr", "name", "shared_record_group"},
		"sharedrecord:a",
	),
	Entry(
		"SharedRecordAAAA",
		&SharedRecordAAAA{},
		[]string{"ipv6addr", "name", "shared_record_group"},
		"sharedrecord:aaaa",
	),
	Entry(
		"SharedrecordCname",
		&SharedrecordCname{},
		[]string{"canonical", "name", "shared_record_group"},
		"sharedrecord:cname",
	),
	Entry(
		"SharedRecordMX",
		&SharedRecordMX{},
		[]string{"mail_exchanger", "name", "preference", "shared_record_group"},
		"sharedrecord:mx",
	),
	Entry(
		"SharedrecordSrv",
		&SharedrecordSrv{},
		[]string{"name", "port", "priority", "shared_record_group", "target", "weight"},
		"sharedrecord:srv",
	),
	Entry(
		"SharedRecordTXT",
		&SharedRecordTXT{},
		[]string{"name", "shared_record_group", "text"},
		"sharedrecord:txt",
	),
	Entry(
		"Sharedrecordgroup",
		&Sharedrecordgroup{},
		[]string{"comment", "name"},
		"sharedrecordgroup",
	),
	Entry(
		"SmartfolderChildren",
		&SmartfolderChildren{},
		[]string{"resource", "value", "value_type"},
		"smartfolder:children",
	),
	Entry(
		"SmartfolderGlobal",
		&SmartfolderGlobal{},
		[]string{"comment", "name"},
		"smartfolder:global",
	),
	Entry(
		"SmartfolderPersonal",
		&SmartfolderPersonal{},
		[]string{"comment", "is_shortcut", "name"},
		"smartfolder:personal",
	),
	Entry(
		"SNMPUser",
		&SNMPUser{},
		[]string{"comment", "name"},
		"snmpuser",
	),
	Entry(
		"Superhost",
		&Superhost{},
		[]string{"comment", "name"},
		"superhost",
	),
	Entry(
		"Superhostchild",
		&Superhostchild{},
		[]string{"comment", "data", "name", "network_view", "parent", "record_parent", "type", "view"},
		"superhostchild",
	),
	Entry(
		"SyslogEndpoint",
		&SyslogEndpoint{},
		[]string{"name", "outbound_member_type"},
		"syslog:endpoint",
	),
	Entry(
		"TacacsplusAuthservice",
		&TacacsplusAuthservice{},
		[]string{"comment", "disable", "name"},
		"tacacsplus:authservice",
	),
	Entry(
		"Taxii",
		&Taxii{},
		[]string{"ipv4addr", "ipv6addr", "name"},
		"taxii",
	),
	Entry(
		"Tftpfiledir",
		&Tftpfiledir{},
		[]string{"directory", "name", "type"},
		"tftpfiledir",
	),
	Entry(
		"ThreatanalyticsAnalyticsWhitelist",
		&ThreatanalyticsAnalyticsWhitelist{},
		[]string{"version"},
		"threatanalytics:analytics_whitelist",
	),
	Entry(
		"ThreatanalyticsModuleset",
		&ThreatanalyticsModuleset{},
		[]string{"version"},
		"threatanalytics:moduleset",
	),
	Entry(
		"ThreatanalyticsWhitelist",
		&ThreatanalyticsWhitelist{},
		[]string{"comment", "disable", "fqdn"},
		"threatanalytics:whitelist",
	),
	Entry(
		"ThreatinsightCloudclient",
		&ThreatinsightCloudclient{},
		[]string{"enable", "interval"},
		"threatinsight:cloudclient",
	),
	Entry(
		"ThreatprotectionGridRule",
		&ThreatprotectionGridRule{},
		[]string{"name", "ruleset", "sid"},
		"threatprotection:grid:rule",
	),
	Entry(
		"ThreatprotectionProfile",
		&ThreatprotectionProfile{},
		[]string{"comment", "name"},
		"threatprotection:profile",
	),
	Entry(
		"ThreatprotectionProfileRule",
		&ThreatprotectionProfileRule{},
		[]string{"profile", "rule"},
		"threatprotection:profile:rule",
	),
	Entry(
		"ThreatprotectionRule",
		&ThreatprotectionRule{},
		[]string{"member", "rule"},
		"threatprotection:rule",
	),
	Entry(
		"ThreatprotectionRulecategory",
		&ThreatprotectionRulecategory{},
		[]string{"name", "ruleset"},
		"threatprotection:rulecategory",
	),
	Entry(
		"ThreatprotectionRuleset",
		&ThreatprotectionRuleset{},
		[]string{"add_type", "version"},
		"threatprotection:ruleset",
	),
	Entry(
		"ThreatprotectionRuletemplate",
		&ThreatprotectionRuletemplate{},
		[]string{"name", "ruleset", "sid"},
		"threatprotection:ruletemplate",
	),
	Entry(
		"ThreatprotectionStatistics",
		&ThreatprotectionStatistics{},
		[]string{"member", "stat_infos"},
		"threatprotection:statistics",
	),
	Entry(
		"Upgradegroup",
		&Upgradegroup{},
		[]string{"comment", "name"},
		"upgradegroup",
	),
	Entry(
		"Upgradeschedule",
		&Upgradeschedule{},
		[]string{"active", "start_time", "time_zone"},
		"upgradeschedule",
	),
	Entry(
		"UpgradeStatus",
		&UpgradeStatus{},
		[]string{"alternate_version", "comment", "current_version", "distribution_version", "element_status", "grid_state", "group_state", "ha_status", "hotfixes", "ipv4_address", "ipv6_address", "member", "message", "pnode_role", "reverted", "status_value", "status_value_update_time", "steps", "steps_completed", "steps_total", "type", "upgrade_group", "upgrade_state", "upgrade_test_status", "upload_version"},
		"upgradestatus",
	),
	Entry(
		"UserProfile",
		&UserProfile{},
		[]string{"name"},
		"userprofile",
	),
	Entry(
		"Vdiscoverytask",
		&Vdiscoverytask{},
		[]string{"name", "state"},
		"vdiscoverytask",
	),
	Entry(
		"View",
		&View{},
		[]string{"comment", "is_default", "name"},
		"view",
	),
	Entry(
		"Vlan",
		&Vlan{},
		[]string{"id", "name", "parent"},
		"vlan",
	),
	Entry(
		"Vlanrange",
		&Vlanrange{},
		[]string{"end_vlan_id", "name", "start_vlan_id", "vlan_view"},
		"vlanrange",
	),
	Entry(
		"Vlanview",
		&Vlanview{},
		[]string{"end_vlan_id", "name", "start_vlan_id"},
		"vlanview",
	),
	Entry(
		"ZoneAuth",
		&ZoneAuth{},
		[]string{"fqdn", "view"},
		"zone_auth",
	),
	Entry(
		"ZoneAuthDiscrepancy",
		&ZoneAuthDiscrepancy{},
		[]string{"description", "severity", "timestamp", "zone"},
		"zone_auth_discrepancy",
	),
	Entry(
		"ZoneDelegated",
		&ZoneDelegated{},
		[]string{"delegate_to", "fqdn", "view"},
		"zone_delegated",
	),
	Entry(
		"ZoneForward",
		&ZoneForward{},
		[]string{"forward_to", "fqdn", "view"},
		"zone_forward",
	),
	Entry(
		"ZoneRp",
		&ZoneRp{},
		[]string{"fqdn", "view"},
		"zone_rp",
	),
	Entry(
		"ZoneStub",
		&ZoneStub{},
		[]string{"fqdn", "stub_from", "view"},
		"zone_stub",
	),
)

var _ = Describe("WAPI_VERSION metadata", func() {
	It("Should be equal to 2.12.1", func() {
		Expect(WAPI_VERSION).To(Equal("2.12.1"))
	})
})
