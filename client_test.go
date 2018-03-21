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
	Test string `xml:"testElem"`
	Action string `xml:"soapAction,attr"`
}

func headerCreator(soapAction string) interface{} {
	return FakeHeader{"Hello", soapAction}
}

func (c *Client) DoSomething(in interface{}) []byte {
	return c.call("this is a soap action string", in)
}


type EmptyStruct struct {

}

func TestEmpty(t *testing.T) {
	mock := MockPoster{}
	c := NewClient("", headerCreator, make(map[string]string), make(map[string]string), &mock)

	c.DoSomething(EmptyStruct{})

	actual := mock.request
	golden := filepath.Join("testdata", t.Name() + ".golden")
	if *update {
		ioutil.WriteFile(golden, actual, 0644)
	}
	expected, _ := ioutil.ReadFile(golden)

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
