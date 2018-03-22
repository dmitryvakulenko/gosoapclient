package gosoapclient

import (
	"testing"
	"io"
	"net/http"
	"flag"
	"path/filepath"
	"io/ioutil"
	"bytes"
)

var update = flag.Bool("update", false, "update .golden files")

type FakeHeader struct {
	XMLName string `xml:"testElem"`
	Action string `xml:"soapAction,attr"`
}

func headerCreator(soapAction string) interface{} {
	return FakeHeader{"Hello", soapAction}
}

func (c *Client) DoSomething(in interface{}) []byte {
	return c.call("this is a soap action string", in)
}


type EmptyStruct struct {}
func TestEmpty(t *testing.T) {
	mock := MockPoster{}
	c := NewClient("", headerCreator, make(map[string]string), make(map[string]string), &mock)

	c.DoSomething(EmptyStruct{})

	actual := mock.request
	expected := readUpdateGolden(t.Name(), actual)

	if !bytes.Equal(actual, expected) {
		t.Fatalf("Wrong")
	}
}

type Session struct {
	XMLName	string `xml:"ns0:Session"`
	SessionId             string `xml:"ns0:SessionId,omitempty"`
	SequenceNumber        string `xml:"ns0:SequenceNumber,omitempty"`
	SecurityToken         []string `xml:"ns0:SecurityToken,omitempty"`
	TransactionStatusCode string `xml:"TransactionStatusCode,attr,omitempty"`
}
func TestNamespaces(t *testing.T) {
	mock := MockPoster{}
	nsAlias := map[string]string{
		"http://aaa.aaa.aaa": "ns0"}
	typesNs := map[string]string{
		"Session": "http://aaa.aaa.aaa"}

	c := NewClient("", headerCreator, typesNs, nsAlias, &mock)
	session := Session{
		SessionId: "a",
		SequenceNumber: "b",
		SecurityToken: []string{"c"},
		TransactionStatusCode: "d"}
	c.DoSomething(session)

	actual := mock.request
	expected := readUpdateGolden(t.Name(), actual)

	if !bytes.Equal(actual, expected) {
		t.Fatalf("Wrong")
	}
}

type MockPoster struct {
	request []byte
}

func (m *MockPoster) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	m.request = make([]byte, 1024)
	size, _ := body.Read(m.request)
	m.request = m.request[0:size]
	return &http.Response{}, nil
}

func readUpdateGolden(name string, actual []byte) []byte {
	golden := filepath.Join("testdata", name + ".golden")
	if *update {
		ioutil.WriteFile(golden, actual, 0644)
	}

	expected, err := ioutil.ReadFile(golden)
	if err != nil {
		panic("Can't read golden file " + golden)
	}

	return expected
}