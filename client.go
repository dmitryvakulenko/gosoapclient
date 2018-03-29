package gosoapclient

import (
	"io"
	"net/http"
	"encoding/xml"
	"bytes"
	"reflect"
	"crypto/tls"
	//"time"
	//"net"
	//"golang.org/x/net/context"
)

type Poster interface {
	Post(url string, contentType string, body io.Reader) (resp *http.Response, err error)
}

type Client struct {
	url             string
	typesNamespaces map[string]string
	namespacesAlias map[string]string
	poster          Poster
	client          *http.Client
}

func NewClient(url string, typesNamespaces map[string]string, namespacesAlias map[string]string, poster Poster) Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	return Client{
		url:             url,
		typesNamespaces: typesNamespaces,
		namespacesAlias: namespacesAlias,
		poster:          poster,
		client:          &http.Client{Transport: tr},
	}
}

func (c *Client) Call(soapAction string, header interface{}, body interface{}) []byte {
	soap := newSoap()
	soap.Header = header
	soap.Body.Content = body

	namespaces := make(map[string]string)
	namespaces = mergeNamespaces(namespaces, c.collectNamespaces(soap.Header))
	namespaces = mergeNamespaces(namespaces, c.collectNamespaces(soap.Body.Content))
	for alias, ns := range namespaces {
		soap.Namespaces = append(soap.Namespaces, xml.Attr{
			Name:  xml.Name{Local: "xmlns:" + alias},
			Value: ns})
	}

	requestBody, err := xml.MarshalIndent(soap, "", "  ")
	if err != nil {
		panic(err)
	}

	requestBody = append([]byte(xml.Header), requestBody...)

	request, _ := http.NewRequest("POST", c.url, bytes.NewReader(requestBody))
	request.Header.Add("Content-Type", "text/xml; charset=\"utf-8\"")
	request.Header.Add("SOAPAction", soapAction)

	response, _ := c.client.Do(request)
	defer response.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	return buf.Bytes()
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
		val, ok := in.([]interface{})
		if !ok {
			return res
		}
		for _, v := range val {
			res = mergeNamespaces(res, c.collectNamespaces(v))
		}
	} else if inType.Kind() == reflect.Ptr && !inType.IsNil() {
		ptr := inType.Elem().Interface()
		res = mergeNamespaces(res, c.collectNamespaces(ptr))
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
