package gosoapclient

import (
	"encoding/xml"
	"flag"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
)

var update = flag.Bool("update", false, "update .golden files")

const emptySoapResponse = "<?xml version=\"1.0\" encoding=\"UTF-8\"?><Envelope><Header/><Body/></Envelope>"

type (
	FakeHeader struct {
		XMLName xml.Name `xml:"testElem"`
		Action  string   `xml:"soapAction,attr"`
	}

	EmptyBody struct {
		XMLName xml.Name `xml:"emptyBody"`
	}

	TestClient struct {
		soapClient     *Client
		httpClientMock *testHttpClient
	}

	testHttpClient struct {
		request  *http.Request
		response *http.Response
	}
)

func (t *testHttpClient) Do(request *http.Request) (*http.Response, error) {
	t.request = request
	return t.response, nil
}

func NewTestClient(url string, response *http.Response) *TestClient {
	httpClientMock := &testHttpClient{response: response}
	return &TestClient{
		soapClient: &Client{
			url:    url,
			client: httpClientMock,
		},
		httpClientMock: httpClientMock}
}

func (c *TestClient) DoSomething(header, body interface{}) (*soapResponse, error) {
	return c.soapClient.Call("this is a soap action string", header, body)
}

func TestEmpty(t *testing.T) {
	c := NewTestClient(
		"http://server.com",
		&http.Response{
			Body: ioutil.NopCloser(strings.NewReader(emptySoapResponse)),
		})
	_, err := c.DoSomething(FakeHeader{Action: "this is a soap action string"}, EmptyBody{})
	assert.Nil(t, err)

	expected := string(readUpdateGolden(t.Name(), nil))

	assert.Equal(t, expected, c.soapClient.lastRequest)
	assert.Equal(t, "this is a soap action string", c.httpClientMock.request.Header.Get("SOAPAction"))
}

func TestFaultProcessing(t *testing.T) {
	responseBody, _ := ioutil.ReadFile("testdata/FaultResponse.xml")
	c := NewTestClient(
		"http://server.com",
		&http.Response{
			Body: ioutil.NopCloser(strings.NewReader(string(responseBody))),
		})

	_, err := c.DoSomething(FakeHeader{Action: "action"}, EmptyBody{})
	assert.NotNil(t, err)
	assert.Equal(t, " 95|Session|Inactive conversation", err.Error())
}

func readUpdateGolden(name string, actual []byte) []byte {
	golden := filepath.Join("testdata", name+".golden")
	if *update {
		ioutil.WriteFile(golden, actual, 0644)
	}

	expected, err := ioutil.ReadFile(golden)
	if err != nil {
		panic("Can't read golden file " + golden)
	}

	return expected
}
