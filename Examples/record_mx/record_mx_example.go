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

	// Create MX Record
	ea := ibclient.EA{"Cloud API Owned": ibclient.Bool(false)} // Optional field
	fmt.Println(objMgr.CreateMXRecord(ibclient.RecordMX{Name: "mx.test.com", View: "default.test_netview",
		MailExchanger: "example.test.com", Preference: 10, Ea: ea}))

	// Get MX Record by Name
	fmt.Println(objMgr.GetMXRecord(ibclient.RecordMX{Name: "mx.test.com"}))

	// Get MX Record by Reference ID
	fmt.Println(objMgr.GetMXRecord(ibclient.RecordMX{Ref: "record:mx/ZG5zLmJpbmRfbXgkLjE4LmNvbS50ZXN0Lm14LmV4YW1wbGUudGVzdC5jb20uMTA:mx.test.com/default.test_netview"}))

	// Get MX Record by MailExchanger or Preference
	fmt.Println(objMgr.GetMXRecord(ibclient.RecordMX{MailExchanger: "example.test.com"}))
	fmt.Println(objMgr.GetMXRecord(ibclient.RecordMX{Preference: 10}))

	// Get all MX Records in a specific view
	fmt.Println(objMgr.GetMXRecord(ibclient.RecordMX{View: "default.test_netview"}))

	// Update Name of MX Record
	fmt.Println(objMgr.UpdateMXRecord(ibclient.RecordMX{Ref: "record:mx/ZG5zLmJpbmRfbXgkLjE4LmNvbS50ZXN0Lm14LmV4YW1wbGUudGVzdC5jb20uMTA:mx.test.com/default.test_netview",
		Name: "mx1.test.com"}))

	// Update MailExchanger or Preference of MX Record
	fmt.Println(objMgr.UpdateMXRecord(ibclient.RecordMX{Ref: "record:mx/ZG5zLmJpbmRfbXgkLjE4LmNvbS50ZXN0Lm14MS5leGFtcGxlLnRlc3QuY29tLjEw:mx1.test.com/default.test_netview",
		MailExchanger: "example1.test.com"}))

	// Delete MX Record by Name or Ref or MailExchanger
	fmt.Println(objMgr.DeleteMXRecord(ibclient.RecordMX{Ref: "record:mx/ZG5zLmJpbmRfbXgkLjE4LmNvbS50ZXN0Lm14MS5leGFtcGxlMS50ZXN0LmNvbS4xMA:mx1.test.com/default.test_netview"}))
	fmt.Println(objMgr.DeleteMXRecord(ibclient.RecordMX{Name: "mx1.test.com"}))
	fmt.Println(objMgr.DeleteMXRecord(ibclient.RecordMX{MailExchanger: "example1.test.com"}))
}
