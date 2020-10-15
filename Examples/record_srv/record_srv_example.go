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

	// Create SRV Record
	fmt.Println(objMgr.CreateSRVRecord(ibclient.RecordSRV{Name: "srv1.test.com", View: "default.test_netview",
		Target: "xmpp-server.example.com", Port: 5269, Weight: 10, Priority: 1}))

	// Get SRV Record by Name
	fmt.Println(objMgr.GetSRVRecord(ibclient.RecordSRV{Name: "srv1.test.com"}))

	// Get SRV Record by Reference ID
	fmt.Println(objMgr.GetSRVRecord(ibclient.RecordSRV{Ref: "record:srv/ZG5zLmJpbmRfc3J2JC4xOC5jb20udGVzdC9zcnYxLzEvMTAvNTI2OS94bXBwLXNlcnZlci5leGFtcGxlLmNvbQ:srv1.test.com/default.test_netview"}))

	// Get SRV Record by MailExchanger or Preference
	fmt.Println(objMgr.GetSRVRecord(ibclient.RecordSRV{Target: "xmpp-server.example.com"}))
	fmt.Println(objMgr.GetSRVRecord(ibclient.RecordSRV{Priority: 1}))

	// Get all SRV Records in a specific view
	fmt.Println(objMgr.GetSRVRecord(ibclient.RecordSRV{View: "default.test_netview"}))

	// Update Name of SRV Record
	fmt.Println(objMgr.UpdateSRVRecord(ibclient.RecordSRV{Ref: "record:srv/ZG5zLmJpbmRfc3J2JC4xOC5jb20udGVzdC9teDEvMS8xMC81MjY5L3htcHAtc2VydmVyLmV4YW1wbGUuY29t:mx1.test.com/default.test_netview",
		Name: "srv2.test.com"}))

	// Update  Name, Port, Priority, Target and Weight of SRV Record
	fmt.Println(objMgr.UpdateSRVRecord(ibclient.RecordSRV{Ref: "record:srv/ZG5zLmJpbmRfc3J2JC4xOC5jb20udGVzdC9zcnYyLzEvMTAvNTI2OS9leGFtcGxlMS50ZXN0LmNvbQ:srv2.test.com/default.test_netview",
		Target: "example1.test.com"}))

	// Delete SRV Record by Name or Ref or Target
	fmt.Println(objMgr.DeleteSRVRecord(ibclient.RecordSRV{Ref: "record:srv/ZG5zLmJpbmRfc3J2JC4xOC5jb20udGVzdC9zcnYyLzEvMTAvNTI2OS9leGFtcGxlMS50ZXN0LmNvbQ:srv2.test.com/default.test_netview"}))
	fmt.Println(objMgr.DeleteSRVRecord(ibclient.RecordSRV{Name: "srv2.test.com"}))
	fmt.Println(objMgr.DeleteSRVRecord(ibclient.RecordSRV{Target: "example1.test.com"}))
}
