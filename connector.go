package ibclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
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
	SslCertificate      string
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

func (c *Connector) makeWapiUrl(fmt_str string, args ...interface{}) string {
	wapiBase := fmt.Sprintf("https://%s:%s/wapi/v%s", c.Host, c.WapiPort, c.WapiVersion)

	return fmt.Sprintf("%s/%s", wapiBase, fmt.Sprintf(fmt_str, args...))
}

func (c *Connector) buildUrl(t RequestType, objType string, payload Payload, ref string) url.URL {
	path := []string{"wapi", "v" + c.WapiVersion}
	if len(ref) > 0 {
		path = append(path, ref)
	} else {
		path = append(path, objType)
	}

	u := url.URL{
		Scheme: "https",
		Host:   c.Host + ":" + c.WapiPort,
		Path:   strings.Join(path, "/"),
	}

	switch t {
	case GET, DELETE:
		values := make(url.Values)
		for k, v := range payload {
			values.Add(k, v.(string))
		}
		u.RawQuery = values.Encode()
	}

	return u
}

func (c *Connector) buildBody(t RequestType, payload Payload) io.Reader {
	var jsonStr []byte
	var err error

	switch t {
	case CREATE, UPDATE:
		jsonStr, err = json.Marshal(payload)
		if err != nil {
			log.Printf("Cannot unmarshal payload: '%s'", payload)
			return nil
		}
	default:
		return nil
	}

	return bytes.NewBuffer(jsonStr)
}

func (c *Connector) makeRequest(t RequestType, objType string, payload Payload, ref string) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !c.SslVerify},
	}
	client := &http.Client{Transport: tr}

	url := c.buildUrl(t, objType, payload, ref)
	body := c.buildBody(t, payload)
	req, err1 := http.NewRequest(t.toMethod(), url.String(), body)
	if err1 != nil {
		log.Printf("err1: '%s'", err1)
		return []byte(""), err1
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.Username, c.Password)

	resp, err3 := client.Do(req)
	if !(resp.StatusCode == http.StatusOK ||
		(resp.StatusCode == http.StatusCreated && t == CREATE)) {
		logHttpResponse(resp)
		return []byte(""), err3
	}
	if err3 != nil {
		log.Printf("Error : %s", err3)
		logHttpResponse(resp)
		return []byte(""), err3
	} else {
		defer resp.Body.Close()
		contents, err4 := ioutil.ReadAll(resp.Body)
		if err4 != nil {
			log.Printf("Http Reponse ioutil.ReadAll() Error: '%s'", err4)
			return []byte(""), err4
		}

		return contents, err4
	}

	return []byte(""), nil
}

func (c *Connector) CreateObject(objType string, payload Payload) (string, error) {
	resp, err := c.makeRequest(CREATE, objType, payload, "")

	var ref string

	err = json.Unmarshal(resp, &ref)
	if err != nil {
		log.Printf("Cannot unmarshall '%s', err: '%s'\n", string(resp), err)
		return "", err
	}

	return ref, err
}

func (c *Connector) GetObject(objType string, payload Payload, ref string, res interface{}) error {
	resp, err := c.makeRequest(GET, objType, payload, ref)

	if len(resp) == 0 {
		return err
	}

	err = json.Unmarshal(resp, res)

	if err != nil {
		log.Printf("Cannot unmarshall '%s', err: '%s'\n", string(resp), err)
		return err
	}

	return err
}

func (c *Connector) DeleteObject(ref string) (string, error) {
	resp, err := c.makeRequest(DELETE, "", nil, ref)

	var refRes string

	err = json.Unmarshal(resp, &refRes)
	if err != nil {
		log.Printf("Cannot unmarshall '%s', err: '%s'\n", string(resp), err)
		return "", err
	}

	return refRes, err
}

func NewConnector(host string, wapiVersion string, wapiPort string,
	username string, password string, sslVerify bool, sslCertificate string,
	httpRequestTimeout int, httpPoolConnections int, httpPoolMaxSize int) *Connector {

	return &Connector{
		Host:                host,
		WapiVersion:         wapiVersion,
		WapiPort:            wapiPort,
		Username:            username,
		Password:            password,
		SslVerify:           sslVerify,
		SslCertificate:      sslCertificate,
		HttpRequestTimeout:  httpRequestTimeout,
		HttpPoolConnections: httpPoolConnections,
		HttpPoolMaxSize:     httpPoolMaxSize,
	}
}
