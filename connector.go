package ibclient

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Connector struct {
	Host                string
	WapiVersion         string
	WapiPort            string
	Username            string
	Password            string
	SslVerify           bool
	certPool            *x509.CertPool
	HttpRequestTimeout  int
	HttpPoolConnections int
	HttpPoolMaxSize     int
	url                 url.URL
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

func (c *Connector) buildUrl(t RequestType, objType string, ref string, returnFields []string) url.URL {
	path := []string{"wapi", "v" + c.WapiVersion}
	if len(ref) > 0 {
		path = append(path, ref)
	} else {
		path = append(path, objType)
	}

	qry := ""
	if t == GET && len(returnFields) > 0 {
		v := url.Values{}
		v.Set("_return_fields", strings.Join(returnFields, ","))
		qry = v.Encode()
	}

	u := url.URL{
		Scheme:   "https",
		Host:     c.Host + ":" + c.WapiPort,
		Path:     strings.Join(path, "/"),
		RawQuery: qry,
	}

	return u
}

func (c *Connector) buildBody(t RequestType, obj IBObject) io.Reader {
	var jsonStr []byte
	var err error

	jsonStr, err = json.Marshal(obj)
	if err != nil {
		log.Printf("Cannot marshal payload: '%s'", obj)
		return nil
	}

	return bytes.NewBuffer(jsonStr)
}

func (c *Connector) makeRequest(t RequestType, obj IBObject, ref string) (res []byte, err error) {
	res = []byte("")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !c.SslVerify, RootCAs: c.certPool},
	}
	client := &http.Client{Transport: tr}

	var objType string = ""
	if obj != nil {
		objType = obj.ObjectType()
	}
	url := c.buildUrl(t, objType, ref, obj.ReturnFields())

	var body io.Reader = nil
	if obj != nil {
		body = c.buildBody(t, obj)
	}

	req, err := http.NewRequest(t.toMethod(), url.String(), body)
	if err != nil {
		log.Printf("err1: '%s'", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.Username, c.Password)

	resp, err := client.Do(req)
	if err != nil {
		return
	} else if !(resp.StatusCode == http.StatusOK ||
		(resp.StatusCode == http.StatusCreated && t == CREATE)) {
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

func (c *Connector) CreateObject(obj IBObject) (ref string, err error) {
	ref = ""

	resp, err := c.makeRequest(CREATE, obj, "")
	if err != nil || len(resp) == 0 {
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

	err = json.Unmarshal(resp, &refRes)
	if err != nil {
		log.Printf("Cannot unmarshall '%s', err: '%s'\n", string(resp), err)
		return
	}

	return
}

func NewConnector(host string, wapiVersion string, wapiPort string,
	username string, password string, sslVerify string, httpRequestTimeout int,
	httpPoolConnections int, httpPoolMaxSize int) (res *Connector, err error) {
	res = nil

	connector := &Connector{
		Host:                host,
		WapiVersion:         wapiVersion,
		WapiPort:            wapiPort,
		Username:            username,
		Password:            password,
		SslVerify:           false,
		certPool:            nil,
		HttpRequestTimeout:  httpRequestTimeout,
		HttpPoolConnections: httpPoolConnections,
		HttpPoolMaxSize:     httpPoolMaxSize,
	}

	switch {
	case "false" == strings.ToLower(sslVerify):
		connector.SslVerify = false
	case "true" == strings.ToLower(sslVerify):
		connector.SslVerify = true
	default:
		var cert []byte
		caPool := x509.NewCertPool()
		cert, err = ioutil.ReadFile(sslVerify)
		if err != nil {
			log.Printf("Cannot load certificate file '%s'", sslVerify)
			return
		}
		if !caPool.AppendCertsFromPEM(cert) {
			err = errors.New(fmt.Sprintf("Cannot append certificate from file '%s'", sslVerify))
			return
		}
		connector.certPool = caPool
		connector.SslVerify = true
	}
	res = connector

	return
}
