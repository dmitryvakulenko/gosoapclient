package gosoapclient

import (
	"io"
	"net/http"
	"encoding/xml"
	"bytes"
	"crypto/tls"
)

type Poster interface {
	Post(url string, contentType string, body io.Reader) (resp *http.Response, err error)
}

type Client struct {
	url             string
	poster          Poster
	client          *http.Client
}

func NewClient(url string, poster Poster) Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	return Client{
		url:             url,
		poster:          poster,
		client:          &http.Client{Transport: tr},
	}
}

func (c *Client) Call(soapAction string, header, body interface{}) []byte {
	soap := NewSoapEnvelope()
	//soap.Header = header
	//soap.Body.Content = body

	//namespaces := make(map[string]string)
	//namespaces = mergeNamespaces(namespaces, c.collectNamespaces(soap.Header))
	//namespaces = mergeNamespaces(namespaces, c.collectNamespaces(soap.Body.Content))

	soap.Namespaces = append(soap.Namespaces, xml.Attr{Name: xml.Name{Local: "xmlns:addr"}, Value: "http://www.w3.org/2005/08/addressing"})

	soap.Header = header
	soap.Body.Value, _ = xml.MarshalIndent(body, "", "    ")

	requestBody, err := xml.MarshalIndent(soap, "", "    ")
	if err != nil {
		panic(err)
	}

	//for alias, ns := range c.Marshaler.namespaces {
	//	soap.Namespaces = append(soap.Namespaces, xml.Attr{
	//		Name:  xml.Name{Local: "xmlns:" + alias},
	//		Value: ns})
	//}

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

