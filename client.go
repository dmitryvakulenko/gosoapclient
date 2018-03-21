package gosoapclient

import (
	"io"
	"net/http"
	"encoding/xml"
	"bytes"
)

type Poster interface {
	Post(url string, contentType string, body io.Reader) (resp *http.Response, err error)
}

type Client struct {
	url             string
	typesNamespaces map[string]string
	namespacesAlias map[string]string
	createHeader    func(string) interface{}
	poster          Poster
}

func NewClient(url string, headerCreator func(string) interface{}, typesNs, nsAlias map[string]string, poster Poster) *Client {
	return &Client{
		url:             url,
		typesNamespaces: typesNs,
		namespacesAlias: nsAlias,
		createHeader:    headerCreator,
		poster:          poster}
}

func (c *Client) call(soapAction string, body interface{}) []byte {
	soap := NewSoap()
	soap.Header = c.createHeader(soapAction)
	soap.Body = body

	request, err := xml.MarshalIndent(soap, "", "    ")
	if err != nil {
		panic("Wrong xml")
	}
	c.poster.Post(c.url, "text\\xml", bytes.NewReader(request))

	return []byte{}
}
