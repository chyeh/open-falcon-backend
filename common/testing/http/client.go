package http

import (
	"fmt"
	"flag"
	"net/http"
	"github.com/dghubble/sling"
	"io/ioutil"
	json "github.com/bitly/go-simplejson"
	checker "gopkg.in/check.v1"
)

// Slint with checker
type CheckSlint struct {
	Slint *sling.Sling
	LastResponse *http.Response

	checker *checker.C
}

// Initialize a checker with slint support
func NewCheckSlint(checker *checker.C, sling *sling.Sling) *CheckSlint {
	return &CheckSlint{
		Slint: sling,
		checker: checker,
	}
}
// Gets request of slint
func (self *CheckSlint) Request() *http.Request {
	req, err := self.Slint.Request()
	self.checker.Assert(err, checker.IsNil)

	return req
}

// Gets the response for current request
func (self *CheckSlint) GetResponse() *http.Response {
	if self.LastResponse != nil {
		return self.LastResponse
	}

	c := self.checker
	client := &http.Client{}

	var err error
	self.LastResponse, err = client.Do(self.Request())
	c.Assert(err, checker.IsNil)

	return self.LastResponse
}

// Asserts the existing of paging header
func (self *CheckSlint) AssertHasPaging() {
	c := self.checker
	resp := self.GetResponse()

	c.Assert(resp.Header.Get("page-size"), checker.Matches, "\\d+")
	c.Assert(resp.Header.Get("page-pos"), checker.Matches, "\\d+")
	c.Assert(resp.Header.Get("total-count"), checker.Matches, "\\d+")
}

// Gets body as string
//
// The exepcted status is used to get expected status
func (self *CheckSlint) GetStringBody(expectedStatus int) string {
	return string(self.checkAndGetBody(expectedStatus))
}

func (self *CheckSlint) checkAndGetBody(expectedStatus int) []byte {
	c := self.checker

	resp := self.GetResponse()
	defer resp.Body.Close()

	c.Check(resp.StatusCode, checker.Equals, expectedStatus)
	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if c.Failed() {
		if err != nil {
			c.Fatalf("Read response(ioutil.ReadAll()) has error: %v", err)
		} else {
			c.Fatalf("Status code not match. Response: %s.", bodyBytes)
		}
	}

	return bodyBytes
}

// Gets body as JSON
//
// The exepcted status is used to get expected status
func (self *CheckSlint) GetJsonBody(expectedStatus int) *json.Json {
	c := self.checker

	jsonResult, err := json.NewJson(self.checkAndGetBody(expectedStatus))
	c.Assert(err, checker.IsNil)

	return jsonResult
}

// The configuration of http client
type HttpClientConfig struct {
	Ssl bool
	Host string
	Port uint16
	Resource string

	slingBase *sling.Sling
}

// Initialize a client config by flag
//
// 	http_host - host name of http service
// 	http_port - port of http service
// 	http_ssl - whether or not use SSL to test http service
func NewHttpClientConfigByFlag() *HttpClientConfig {
	var host = flag.String("http.host", "127.0.0.1", "Host of HTTP service to be tested")
	var port = flag.Int("http.port", 80, "Port of HTTP service to be tested")
	var ssl = flag.Bool("http.ssl", false, "Whether or not to use SSL for HTTP service to be tested")
	var resource = flag.String("http.resource", "", "resource for http://<host>:<port/<resource>")

	flag.Parse()

	config := &HttpClientConfig {
		Host: *host,
		Port: uint16(*port),
		Ssl: *ssl,
		Resource: *resource,
	}
	config.slingBase = sling.New().Base(
		config.hostAndPort(),
	)

	if config.Resource != "" {
		config.slingBase.Path(config.Resource + "/")
	}

	logger.Infof("Sling URL for testing: %s", config.String())

	return config
}

// Gets the full URL of tested service
func (self *HttpClientConfig) String() string {
	url := self.hostAndPort()

	if self.Resource != "" {
		url += "/" + self.Resource
	}

	return url
}

func (self *HttpClientConfig) NewSlingByBase() *sling.Sling {
	return self.slingBase.New()
}

func (self *HttpClientConfig) hostAndPort() string {
	schema := "http"
	if self.Ssl {
		schema = "https"
	}

	return fmt.Sprintf("%s://%s:%d", schema, self.Host, self.Port)
}
