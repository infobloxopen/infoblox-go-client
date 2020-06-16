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

	// Create AAAA Record
	ea := ibclient.EA{"Cloud API Owned": ibclient.Bool(false)} // Optional field
	fmt.Println(objMgr.CreateAAAARecord(ibclient.RecordAAAA{Name: "record1.test.com", View: "default.test_netview",
		Ipv6Addr: "fd60:e32:f1b9::2", Ea: ea}))

	// Get AAAA Record by Name
	fmt.Println(objMgr.GetAAAARecord(ibclient.RecordAAAA{Name: "record1.test.com"}))

	// Get AAAA Record by Reference ID
	fmt.Println(objMgr.GetAAAARecord(ibclient.RecordAAAA{Ref: "record:aaaa/ZG5zLmJpbmRfYWFhYSQuMTguY29tLnRlc3QscmVjb3JkMSxmZDYwOmUzMjpmMWI5Ojoy:record1.test.com/default.test_netview"}))

	// Get AAAA Record by IPv6Addr
	fmt.Println(objMgr.GetAAAARecord(ibclient.RecordAAAA{Ipv6Addr: "fd60:e32:f1b9::2"}))

	// Get all AAAA Records in a specific view
	fmt.Println(objMgr.GetAAAARecord(ibclient.RecordAAAA{View: "default.test_netview"}))

	// Update Name of AAAA Record
	fmt.Println(objMgr.UpdateAAAARecord(ibclient.RecordAAAA{Ref: "record:aaaa/ZG5zLmJpbmRfYWFhYSQuMTguY29tLnRlc3QscmVjb3JkMSxmZDYwOmUzMjpmMWI5Ojoy:record1.test.com/default.test_netview",
		Name: "record2.test.com"}))

	// Update IPv6Addr of AAAA Record
	fmt.Println(objMgr.UpdateAAAARecord(ibclient.RecordAAAA{Ref: "record:aaaa/ZG5zLmJpbmRfYWFhYSQuMTguY29tLnRlc3QscmVjb3JkMixmZDYwOmUzMjpmMWI5Ojoy:record2.test.com/default.test_netview",
		Ipv6Addr: "fd60:e32:f1b9::3"}))

	//Update EAs of AAAA Record (Add or Remove)
	ea1 := ibclient.EA{"Cloud API Owned": ibclient.Bool(true)}
	fmt.Println(objMgr.UpdateAAAARecord(ibclient.RecordAAAA{Ref: "record:aaaa/ZG5zLmJpbmRfYWFhYSQuMTguY29tLnRlc3QscmVjb3JkMixmZDYwOmUzMjpmMWI5Ojoy:record2.test.com/default.test_netview",
		AddEA: ea1}))
	fmt.Println(objMgr.UpdateAAAARecord(ibclient.RecordAAAA{Ref: "record:aaaa/ZG5zLmJpbmRfYWFhYSQuMTguY29tLnRlc3QscmVjb3JkMixmZDYwOmUzMjpmMWI5Ojoy:record2.test.com/default.test_netview",
		RemoveEA: ea1}))

	// Delete AAAA Record by Name or Ref or IPv6Addr
	fmt.Println(objMgr.DeleteAAAARecord(ibclient.RecordAAAA{Ref: "record:aaaa/ZG5zLmJpbmRfYWFhYSQuMTguY29tLnRlc3QscmVjb3JkMixmZDYwOmUzMjpmMWI5Ojoy:record2.test.com/default.test_netview"}))
	fmt.Println(objMgr.DeleteAAAARecord(ibclient.RecordAAAA{Name: "record2.test.com"}))
	fmt.Println(objMgr.DeleteAAAARecord(ibclient.RecordAAAA{Ipv6Addr: "fd60:e32:f1b9::3"}))
}
