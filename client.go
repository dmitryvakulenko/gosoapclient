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
	soap.Header.Content = c.createHeader(soapAction)
	soap.Body.Content = body

	namespaces := make(map[string]string)
	namespaces = mergeNamespaces(namespaces, c.collectNamespaces(reflect.TypeOf(soap.Header.Content)))
	namespaces = mergeNamespaces(namespaces, c.collectNamespaces(reflect.TypeOf(soap.Body.Content)))
	for alias, ns := range namespaces {
		soap.Namespaces = append(soap.Namespaces, xml.Attr{
			Name: xml.Name{Local: "xmlns:" + alias},
			Value: ns})
	}

	request, err := xml.MarshalIndent(soap, "", "    ")
	if err != nil {
		panic("Wrong xml")
	}
	c.poster.Post(c.url, "text\\xml", bytes.NewReader(request))

	return []byte{}
}

func (c *Client) collectNamespaces(inType reflect.Type) map[string]string {
	res := make(map[string]string)

	if ns, ok := c.typesNamespaces[inType.Name()]; ok {
		nsAlias := c.namespacesAlias[ns]
		res[nsAlias] = ns
	}

	fieldsNum := inType.NumField()
	for i := 0; i < fieldsNum; i++ {
		field := inType.Field(i)
		fieldType := field.Type
		if fieldType.Kind() == reflect.Struct {
			res = mergeNamespaces(res, c.collectNamespaces(fieldType))
		} else if fieldType.Kind() == reflect.Array {
			res = mergeNamespaces(res, c.collectNamespaces(field.Type.Elem()))
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