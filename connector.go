package ibclient

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HostConfig struct {
	Host     string
	Version  string
	Port     string
	Username string
	Password string
}

type TransportConfig struct {
	SslVerify           bool
	certPool            *x509.CertPool
	HttpRequestTimeout  int // in seconds
	HttpPoolConnections int
}

func NewTransportConfig(sslVerify string, httpRequestTimeout int, httpPoolConnections int) (cfg TransportConfig) {
	switch {
	case "false" == strings.ToLower(sslVerify):
		cfg.SslVerify = false
	case "true" == strings.ToLower(sslVerify):
		cfg.SslVerify = true
	default:
		caPool := x509.NewCertPool()
		cert, err := ioutil.ReadFile(sslVerify)
		if err != nil {
			log.Printf("Cannot load certificate file '%s'", sslVerify)
			return
		}
		if !caPool.AppendCertsFromPEM(cert) {
			err = errors.New(fmt.Sprintf("Cannot append certificate from file '%s'", sslVerify))
			return
		}
		cfg.certPool = caPool
		cfg.SslVerify = true
	}

	return
}

type HttpRequestBuilder interface {
	Init(HostConfig)
	BuildUrl(r RequestType, objType string, ref string, returnFields []string, eaSearch EA) (urlStr string)
	BuildBody(obj IBObject) (jsonStr []byte)
	BuildRequest(r RequestType, obj IBObject, ref string) (req *http.Request, err error)
}

type HttpRequestor interface {
	Init(TransportConfig)
	SendRequest(*http.Request) ([]byte, error)
}

type WapiRequestBuilder struct {
	HostConfig HostConfig
}

type WapiHttpRequestor struct {
	client http.Client
}

type IBConnector interface {
	CreateObject(obj IBObject) (ref string, err error)
	GetObject(obj IBObject, ref string, res interface{}) error
	DeleteObject(ref string) (refRes string, err error)
}

type Connector struct {
	HostConfig      HostConfig
	TransportConfig TransportConfig
	RequestBuilder  HttpRequestBuilder
	Requestor       HttpRequestor
}

type RequestType int

const (
	CREATE RequestType = iota
	GET
	DELETE
	UPDATE
)

func (r RequestType) toMethod() string {
	switch r {
	case CREATE:
		return "POST"
	case GET:
		return "GET"
	case DELETE:
		return "DELETE"
	case UPDATE:
		return "PUT"
	}

	return ""
}

func logHttpResponse(resp *http.Response) {
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	log.Printf("WAPI request error: %d('%s')\nContents:\n%s\n", resp.StatusCode, resp.Status, content)
}

func (whr *WapiHttpRequestor) Init(cfg TransportConfig) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !cfg.SslVerify,
			RootCAs: cfg.certPool},
		MaxIdleConnsPerHost:   cfg.HttpPoolConnections,
		ResponseHeaderTimeout: time.Duration(cfg.HttpRequestTimeout * 1000000000), // ResponseHeaderTimeout is in nanoseconds
	}

	whr.client = http.Client{Transport: tr}
}

func (whr *WapiHttpRequestor) SendRequest(req *http.Request) (res []byte, err error) {
	var resp *http.Response
	resp, err = whr.client.Do(req)
	if err != nil {
		return
	} else if !(resp.StatusCode == http.StatusOK ||
		(resp.StatusCode == http.StatusCreated &&
			req.Method == RequestType(CREATE).toMethod())) {
		logHttpResponse(resp)
		return
	}
	defer resp.Body.Close()
	res, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Http Reponse ioutil.ReadAll() Error: '%s'", err)
		return
	}

	return
}

func (wrb *WapiRequestBuilder) Init(cfg HostConfig) {
	wrb.HostConfig = cfg
}

func (wrb *WapiRequestBuilder) BuildUrl(t RequestType, objType string, ref string, returnFields []string, eaSearch EA) (urlStr string) {
	path := []string{"wapi", "v" + wrb.HostConfig.Version}
	if len(ref) > 0 {
		path = append(path, ref)
	} else {
		path = append(path, objType)
	}

	qry := ""
	vals := url.Values{}
	if t == GET {
		if len(returnFields) > 0 {
			vals.Set("_return_fields", strings.Join(returnFields, ","))
		}

		if len(eaSearch) > 0 {
			for k, v := range eaSearch {
				str, ok := v.(string)
				if !ok {
					log.Printf("Cannot marshal EA Search attribute for '%s'\n", k)
				} else {
					vals.Set("*"+k, str)
				}
			}
		}

		qry = vals.Encode()
	}

	u := url.URL{
		Scheme:   "https",
		Host:     wrb.HostConfig.Host + ":" + wrb.HostConfig.Port,
		Path:     strings.Join(path, "/"),
		RawQuery: qry,
	}

	return u.String()
}

func (wrb *WapiRequestBuilder) BuildBody(obj IBObject) []byte {
	var jsonStr []byte
	var err error

	jsonStr, err = json.Marshal(obj)
	if err != nil {
		log.Printf("Cannot marshal payload: '%s'", obj)
		return nil
	}

	return jsonStr
}

func (wrb *WapiRequestBuilder) BuildRequest(t RequestType, obj IBObject, ref string) (req *http.Request, err error) {
	var (
		objType      string
		returnFields []string
		eaSearch     EA
	)
	if obj != nil {
		objType = obj.ObjectType()
		returnFields = obj.ReturnFields()
		eaSearch = obj.EaSearch()
	}
	urlStr := wrb.BuildUrl(t, objType, ref, returnFields, eaSearch)

	var bodyStr []byte
	if obj != nil {
		bodyStr = wrb.BuildBody(obj)
	}

	req, err = http.NewRequest(t.toMethod(), urlStr, bytes.NewBuffer(bodyStr))
	if err != nil {
		log.Printf("err1: '%s'", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(wrb.HostConfig.Username, wrb.HostConfig.Password)

	return
}

func (c *Connector) makeRequest(t RequestType, obj IBObject, ref string) (res []byte, err error) {
	var req *http.Request
	req, err = c.RequestBuilder.BuildRequest(t, obj, ref)
	res, err = c.Requestor.SendRequest(req)

	return
}

func (c *Connector) CreateObject(obj IBObject) (ref string, err error) {
	ref = ""

	resp, err := c.makeRequest(CREATE, obj, "")
	if err != nil || len(resp) == 0 {
		log.Printf("CreateObject request error: '%s'\n", err)
		return
	}

	err = json.Unmarshal(resp, &ref)
	if err != nil {
		log.Printf("Cannot unmarshall '%s', err: '%s'\n", string(resp), err)
		return
	}

	return
}

func (c *Connector) GetObject(obj IBObject, ref string, res interface{}) (err error) {
	resp, err := c.makeRequest(GET, obj, ref)
	if err != nil {
		log.Printf("GetObject request error: '%s'\n", err)
	}

	if len(resp) == 0 {
		return
	}

	err = json.Unmarshal(resp, res)

	if err != nil {
		log.Printf("Cannot unmarshall '%s', err: '%s'\n", string(resp), err)
		return
	}

	return
}

func (c *Connector) DeleteObject(ref string) (refRes string, err error) {
	refRes = ""

	resp, err := c.makeRequest(DELETE, nil, ref)
	if err != nil {
		log.Printf("DeleteObject request error: '%s'\n", err)
	}

	err = json.Unmarshal(resp, &refRes)
	if err != nil {
		log.Printf("Cannot unmarshall '%s', err: '%s'\n", string(resp), err)
		return
	}

	return
}

func NewConnector(hostConfig HostConfig, transportConfig TransportConfig,
	requestBuilder HttpRequestBuilder, requestor HttpRequestor) (res *Connector, err error) {
	res = nil

	connector := &Connector{
		HostConfig:      hostConfig,
		TransportConfig: transportConfig,
	}

	//connector.RequestBuilder = WapiRequestBuilder{WaipHostConfig: connector.HostConfig}
	connector.RequestBuilder = requestBuilder
	connector.RequestBuilder.Init(connector.HostConfig)

	connector.Requestor = requestor
	connector.Requestor.Init(connector.TransportConfig)

	res = connector

	return
}
