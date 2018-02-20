# Infoblox Go Client

An Infoblox Client library for Go.

This library is compatible with Go 1.2+

- [Prerequisites](#Prerequisites)
- [Installation](#Installation)
- [Usage](#Usage)

## Prerequisites
   * Infoblox GRID with 2.5 or above WAPI support
   * Go 1.2 or above

## Installation
   go get github.com/infobloxopen/infoblox-go-client

## Usage

   The following is a very simple example for the client usage:

       package main
       import (
   	    "fmt"
   	    ibclient "github.com/infobloxopen/infoblox-go-client"
       )

       func main() {
   	    hostConfig := ibclient.HostConfig{
   		    Host:     "10.196.107.209",
   		    Version:  "2.8",
   		    Port:     "443",
   		    Username: "admin",
   		    Password: "infoblox",
   	    }
   	    transportConfig := ibclient.NewTransportConfig("false", 20, 10)
   	    requestBuilder := &ibclient.WapiRequestBuilder{}
   	    requestor := &ibclient.WapiHttpRequestor{}
   	    conn, err := ibclient.NewConnector(hostConfig, transportConfig, requestBuilder, requestor)
   	    if err != nil {
   		    fmt.Println(err)
   	    }
   	    objMgr := ibclient.NewObjectManager(conn, "myclient", "")
   	    //Fetches grid information
   	    fmt.Println(objMgr.GetLicense())
       }

## Supported NIOS operations

   * CreateNetworkView
   * CreateDefaultNetviews
   * CreateNetwork
   * CreateNetworkContainer
   * GetNetworkView
   * GetNetwork
   * GetNetworkContainer
   * AllocateNetwork
   * UpdateFixedAddress
   * GetFixedAddress
   * ReleaseIP
   * DeleteNetwork
   * GetEADefinition
   * CreateEADefinition
   * UpdateNetworkViewEA
   * GetCapacityReport
   * GetAllMembers
   * GetUpgradeStatus (2.7 or above)
