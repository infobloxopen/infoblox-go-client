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

	// Create PTR Record
	ea := ibclient.EA{"Cloud API Owned": ibclient.Bool(false)} // Optional field
	fmt.Println(objMgr.CreatePTRRecord(ibclient.RecordPTR{View: "default.test_netview", Name: "8.2.168.192.in-addr.arpa",
			PtrdName: "ptr_test1.test.com", Ipv4Addr: "192.168.2.8", Ea: ea}))

	fmt.Println(objMgr.CreatePTRRecord(ibclient.RecordPTR{View: "default.test_netview", Name: "3.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.9.b.1.f.2.3.e.0.0.6.d.f.ip6.arpa",
		PtrdName: "ptr_a41.test.com", Ipv6Addr: "fd60:e32:f1b9::3", Ea: ea}))

	// Get PTR Record by Name
	fmt.Println(objMgr.GetPTRRecord(ibclient.RecordPTR{Name: "8.2.168.192.in-addr.arpa"}))

	// Get PTR Record by Reference ID
	fmt.Println(objMgr.GetPTRRecord(ibclient.RecordPTR{Ref: "record:ptr/ZG5zLmJpbmRfcHRyJC4xOC5hcnBhLmluLWFkZHIuMTkyLjE2OC4yLjgucHRyX3Rlc3QxLnRlc3QuY29t:8.2.168.192.in-addr.arpa/default.test_netview 192.168.2.8"}))

	// Get PTR Record by Ipv4Addr or Ipv6Addr
	fmt.Println(objMgr.GetPTRRecord(ibclient.RecordPTR{Ipv4Addr: "192.168.2.8"}))

	// Get PTR Record by PtrdName
	fmt.Println(objMgr.GetPTRRecord(ibclient.RecordPTR{PtrdName: "ptr_test1.test.com"}))

	// Get all PTR Records in a specific view
	fmt.Println(objMgr.GetPTRRecord(ibclient.RecordPTR{View: "default.test_netview"}))

	// Update Name of PTR Record
	fmt.Println(objMgr.UpdatePTRRecord(ibclient.RecordPTR{Ref: "record:ptr/ZG5zLmJpbmRfcHRyJC4xOC5hcnBhLmluLWFkZHIuMTkyLjE2OC4yLjgucHRyX3Rlc3QxLnRlc3QuY29t:8.2.168.192.in-addr.arpa/default.test_netview",
		Name: "9.2.168.192.in-addr.arpa"}))

	// Update PtrdName of PTR Record
	fmt.Println(objMgr.UpdatePTRRecord(ibclient.RecordPTR{Ref: "record:ptr/ZG5zLmJpbmRfcHRyJC4xOC5hcnBhLmluLWFkZHIuMTkyLjE2OC4yLjkucHRyX3Rlc3QxLnRlc3QuY29t:9.2.168.192.in-addr.arpa/default.test_netview ",
		PtrdName: "ptr_a.test.com"}))

	//Update EAs of PTR Record (Add or Remove)
	ea1 := ibclient.EA{"Cloud API Owned": ibclient.Bool(true)}
	fmt.Println(objMgr.UpdatePTRRecord(ibclient.RecordPTR{Ref: "record:ptr/ZG5zLmJpbmRfcHRyJC4xOC5hcnBhLmluLWFkZHIuMTkyLjE2OC4yLjkucHRyX2EudGVzdC5jb20:9.2.168.192.in-addr.arpa/default.test_netview",
		AddEA: ea1}))
	fmt.Println(objMgr.UpdatePTRRecord(ibclient.RecordPTR{Ref: "record:ptr/ZG5zLmJpbmRfcHRyJC4xOC5hcnBhLmluLWFkZHIuMTkyLjE2OC4yLjkucHRyX2EudGVzdC5jb20:9.2.168.192.in-addr.arpa/default.test_netview",
		RemoveEA: ea1}))

	// Delete PTR Record by Name or Ref or PtrdName
	fmt.Println(objMgr.DeletePTRRecord(ibclient.RecordPTR{Ref: "record:ptr/ZG5zLmJpbmRfcHRyJC4xOC5hcnBhLmluLWFkZHIuMTkyLjE2OC4yLjkucHRyX2EudGVzdC5jb20:9.2.168.192.in-addr.arpa/default.test_netview"}))
	fmt.Println(objMgr.DeletePTRRecord(ibclient.RecordPTR{Name: "9.2.168.192.in-addr.arpa"}))
	fmt.Println(objMgr.DeletePTRRecord(ibclient.RecordPTR{PtrdName: "ptr_a.test.com"}))
}
