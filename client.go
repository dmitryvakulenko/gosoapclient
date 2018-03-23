package gosoapclient

import (
	"io"
	"net/http"
	"encoding/xml"
	"bytes"
	"reflect"
)

type Poster interface {
	Post(url string, contentType string, body io.Reader) (resp *http.Response, err error)
}

type Client struct {
	url             string
	typesNamespaces map[string]string
	namespacesAlias map[string]string
	poster          Poster
}

func NewClient(url string, typesNs, nsAlias map[string]string, poster Poster) *Client {
	return &Client{
		url:             url,
		typesNamespaces: typesNs,
		namespacesAlias: nsAlias,
		poster:          poster}
}

func (c *Client) Call(soapAction string, header interface{}, body interface{}) []byte {
	soap := NewSoap()
	soap.Header.Content = header
	soap.Body.Content = body

	namespaces := make(map[string]string)
	namespaces = mergeNamespaces(namespaces, c.collectNamespaces(soap.Header.Content))
	namespaces = mergeNamespaces(namespaces, c.collectNamespaces(soap.Body.Content))
	for alias, ns := range namespaces {
		soap.Namespaces = append(soap.Namespaces, xml.Attr{
			Name: xml.Name{Local: "xmlns:" + alias},
			Value: ns})
	}

	request, err := xml.MarshalIndent(soap, "", "    ")
	if err != nil {
		panic(err)
	}
	response, _ := c.poster.Post(c.url, "text\\xml", bytes.NewReader(request))
	defer response.Body.Close()

	res := make([]byte, response.ContentLength)
	response.Body.Read(res)

	return res
}

func (c *Client) collectNamespaces(in interface{}) map[string]string {
	inType := reflect.TypeOf(in)
	res := make(map[string]string)

	if ns, ok := c.typesNamespaces[inType.Name()]; ok {
		nsAlias := c.namespacesAlias[ns]
		res[nsAlias] = ns
	}

	if inType.Kind() == reflect.Struct {
		fieldsNum := inType.NumField()
		for i := 0; i < fieldsNum; i++ {
			res = mergeNamespaces(res, c.collectNamespaces(inType.Field(i).Type))
		}
	} else if inType.Kind() == reflect.Slice {
		val := in.([]interface{})
		for _, v := range val {
			res = mergeNamespaces(res, c.collectNamespaces(v))
		}
	}

	return res
}

func mergeNamespaces(first map[string]string, second map[string]string) map[string]string {
	res := first
	for k, v := range second {
		res[k] = v
	}

	return res
}