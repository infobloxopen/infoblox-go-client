package main

import (
	"fmt"
	ibclient "github.com/infobloxopen/infoblox-go-client"
)

func main() {
	hostConfig := ibclient.HostConfig{
		Host:     "<NIOS grid IP>",
		Version:  "<WAPI version>",
		Port:     "PORT",
		Username: "username",
		Password: "password",
	}
	transportConfig := ibclient.NewTransportConfig("false", 20, 10)
	requestBuilder := &ibclient.WapiRequestBuilder{}
	requestor := &ibclient.WapiHttpRequestor{}
	conn, err := ibclient.NewConnector(hostConfig, transportConfig, requestBuilder, requestor)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Logout()
	objMgr := ibclient.NewObjectManager(conn, "myclient", "akh")

	// Create Alias Record
	ea := ibclient.EA{"Cloud API Owned": ibclient.Bool(false)} // Optional field
	fmt.Println(objMgr.CreateAliasRecord(ibclient.RecordAlias{Name: "alias.test.com", View: "default.test_netview",
		TargetName: "record3.test.com", TargetType: "A", Ea: ea}))

	// Get Alias Record by Name
	fmt.Println(objMgr.GetAliasRecord(ibclient.RecordAlias{Name: "alias.test.com"}))

	// Get Alias Record by Reference ID
	fmt.Println(objMgr.GetAliasRecord(ibclient.RecordAlias{Ref: "record:alias/ZG5zLmFsaWFzX3JlY29yZCQuMTguY29tLnRlc3QuYWxpYXMuQQ:alias.test.com/default.test_netview"}))

	// Get Alias Record by TargetName or TargetType
	fmt.Println(objMgr.GetAliasRecord(ibclient.RecordAlias{TargetName: "record3.test.com"}))
	fmt.Println(objMgr.GetAliasRecord(ibclient.RecordAlias{TargetType: "A"}))

	// Get all Alias Records in a specific view
	fmt.Println(objMgr.GetAliasRecord(ibclient.RecordAlias{View: "default.test_netview"}))

	// Update Name of Alias Record
	fmt.Println(objMgr.UpdateAliasRecord(ibclient.RecordAlias{Ref: "record:alias/ZG5zLmFsaWFzX3JlY29yZCQuMTguY29tLnRlc3QuYWxpYXMuQQ:alias.test.com/default.test_netview",
		Name: "alias1.test.com"}))

	// Update TargetName or TargetType of Alias Record
	fmt.Println(objMgr.UpdateAliasRecord(ibclient.RecordAlias{Ref: "record:alias/ZG5zLmFsaWFzX3JlY29yZCQuMTguY29tLnRlc3QuYWxpYXMxLkE:alias1.test.com/default.test_netview",
		TargetName: "record4.test.com"}))

	// Delete Alias Record by Name or Ref or TargetName
	fmt.Println(objMgr.DeleteAliasRecord(ibclient.RecordAlias{Ref: "record:alias/ZG5zLmFsaWFzX3JlY29yZCQuMTguY29tLnRlc3QuYWxpYXMxLkE:alias1.test.com/default.test_netview"}))
	fmt.Println(objMgr.DeleteAliasRecord(ibclient.RecordAlias{Name: "alias1.test.com"}))
	fmt.Println(objMgr.DeleteAliasRecord(ibclient.RecordAlias{TargetName: "record4.test.com"}))
}
