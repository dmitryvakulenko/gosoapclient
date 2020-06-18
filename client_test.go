package gosoapclient

import (
	"encoding/xml"
	"flag"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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

func testHttpRequest(method string, url string, body io.Reader) (*http.Request, error) {
	return httptest.NewRequest(method, url, body), nil
}

func TestEmpty(t *testing.T) {
	c := NewTestClient("http://server.com", &http.Response{
		Body: ioutil.NopCloser(strings.NewReader(emptySoapResponse)),
	})
	_, err := c.DoSomething(FakeHeader{Action: "this is a soap action string"}, EmptyBody{})
	assert.Nil(t, err)

	expected := string(readUpdateGolden(t.Name(), nil))

	assert.Equal(t, expected, c.soapClient.lastRequest)
}

type Session struct {
	XMLName               string   `xml:"ns0:Session"`
	SessionId             string   `xml:"ns0:SessionId,omitempty"`
	SequenceNumber        string   `xml:"ns0:SequenceNumber,omitempty"`
	SecurityToken         []string `xml:"ns0:SecurityToken,omitempty"`
	TransactionStatusCode string   `xml:"TransactionStatusCode,attr,omitempty"`
}

// func TestNamespaces(t *testing.T) {
// 	mock := MockPoster{}
// 	nsAlias := map[string]string{
// 		"http://aaa.aaa.aaa": "ns0"}
// 	typesNs := map[string]string{
// 		"Session": "http://aaa.aaa.aaa"}
//
// 	soapClient := NewClient("", headerCreator, typesNs, nsAlias, &mock)
// 	session := Session{
// 		SessionId:             "a",
// 		SequenceNumber:        "b",
// 		SecurityToken:         []string{"soapClient"},
// 		TransactionStatusCode: "d"}
// 	soapClient.DoSomething(session)
//
// 	actual := mock.request
// 	expected := readUpdateGolden(t.Name(), actual)
//
// 	if !bytes.Equal(actual, expected) {
// 		t.Fatalf("Wrong")
// 	}
// }

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
