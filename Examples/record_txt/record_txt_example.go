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

	// Create TXT Record
	fmt.Println(objMgr.CreateTXTRecord(ibclient.RecordTXT{Name: "server.test.com", View: "default.test_netview",
		Text: "This is a host server", TTL: 5}))

	// Get Txt Record by Name
	fmt.Println(objMgr.GetTXTRecord(ibclient.RecordTXT{Name: "server.test.com"}))

	// Get TXT Record by Reference ID
	fmt.Println(objMgr.GetTXTRecord(ibclient.RecordTXT{Ref: "record:txt/ZG5zLmJpbmRfdHh0JC4xOC5jb20udGVzdC5zZXJ2ZXIuIlRoaXMiICJpcyIgImEiICJob3N0IiAic2VydmVyIg:server.test.com/default.test_netview server.test.com"}))

	// Get TXT Record by Text
	fmt.Println(objMgr.GetTXTRecord(ibclient.RecordTXT{Text: "This is a host server"}))

	// Get all TXT Records in a specific view
	fmt.Println(objMgr.GetTXTRecord(ibclient.RecordTXT{View: "default.test_netview"}))

	// Update Name of TXT Record
	fmt.Println(objMgr.UpdateTXTRecord(ibclient.RecordTXT{Ref: "record:txt/ZG5zLmJpbmRfdHh0JC4xOC5jb20udGVzdC5zZXJ2ZXIuIlRoaXMiICJpcyIgImEiICJob3N0IiAic2VydmVyIg:server.test.com/default.test_netview server.test.com",
		Name: "server1.test.com"}))

	// Update Text or TTL of TXT Record
	fmt.Println(objMgr.UpdateTXTRecord(ibclient.RecordTXT{Ref: "record:txt/ZG5zLmJpbmRfdHh0JC4xOC5jb20udGVzdC5zZXJ2ZXIxLiJUaGlzIiAiaXMiICJhIiAiaG9zdCIgInNlcnZlciI:server1.test.com/default.test_netview server1.test.com",
		Text: "Host Server"}))

	// Delete TXT Record by Name or Ref or IPv6Addr
	fmt.Println(objMgr.DeleteTXTRecord(ibclient.RecordTXT{Ref: "record:txt/ZG5zLmJpbmRfdHh0JC4xOC5jb20udGVzdC5zZXJ2ZXIxLiJIb3N0IiAiU2VydmVyIg:server1.test.com/default.test_netview"}))
	fmt.Println(objMgr.DeleteTXTRecord(ibclient.RecordTXT{Name: "server1.test.com"}))
	fmt.Println(objMgr.DeleteTXTRecord(ibclient.RecordTXT{Text: "Host Server"}))
}
