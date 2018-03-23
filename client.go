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
	soap.Header = header
	soap.Body.Content = body

	namespaces := make(map[string]string)
	namespaces = mergeNamespaces(namespaces, c.collectNamespaces(soap.Header))
	namespaces = mergeNamespaces(namespaces, c.collectNamespaces(soap.Body))
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

	return readResponse(response)
}

func readResponse(response *http.Response) []byte {
	res := bytes.NewBuffer(make([]byte, 0))
	buf := make([]byte, 1024)
	var err error
	var readed int
	for err != io.EOF {
		readed, err = response.Body.Read(buf)
		res.Write(buf[0:readed])
	}

	return res.Bytes()
}

func (c *Client) collectNamespaces(in interface{}) map[string]string {
	inType := reflect.ValueOf(in)
	res := make(map[string]string)

	inTypeName := inType.Type().Name()
	if ns, ok := c.typesNamespaces[inTypeName]; ok {
		nsAlias := c.namespacesAlias[ns]
		res[nsAlias] = ns
	}

	if inType.Kind() == reflect.Struct {
		fieldsNum := inType.NumField()
		for i := 0; i < fieldsNum; i++ {
			val := inType.Field(i)
			res = mergeNamespaces(res, c.collectNamespaces(val.Interface()))
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