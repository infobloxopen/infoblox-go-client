# Infoblox Go Client

## Build Status

| Master                                                                                                                                          | Develop                                                                                                                                                           |
| ----------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [![Build Status](https://travis-ci.org/infobloxopen/infoblox-go-client.svg?branch=master)](https://travis-ci.org/infobloxopen/infoblox-go-client) | [![Build Status](https://travis-ci.org/infobloxopen/infoblox-go-client.svg?branch=develop)](https://travis-ci.org/infobloxopen/infoblox-go-client) |


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
   	    objMgr := ibclient.NewObjectManager(conn, "myclient", "")
   	    //Fetches grid information
   	    fmt.Println(objMgr.GetLicense())
       }

## Supported NIOS operations

   * AllocateNetwork
   * CreateDefaultNetviews
   * CreateEADefinition
   * CreateNetwork
   * CreateNetworkContainer
   * CreateNetworkView
   * DeleteNetwork
   * DeleteNetworkView
   * GetAllMembers
   * GetCapacityReport
   * GetEADefinition
   * GetFixedAddress
   * GetNetwork
   * GetNetworkContainer
   * GetNetworkView
   * GetUpgradeStatus (2.7 or above)
   * ReleaseIP
   * UpdateFixedAddress
   * UpdateNetworkViewEA
