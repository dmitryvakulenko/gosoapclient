package gosoapclient

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"log"
	"net/http"
)

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

// TODO: some kind of connection keeping or reconnect
type Client struct {
	url          string
	lastRequest  string
	lastResponse string
	client       httpClient
}

func NewClient(url string) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	return &Client{
		url:    url,
		client: &http.Client{Transport: tr},
	}
}

func (c *Client) Call(soapAction string, header, body interface{}) (*soapResponse, error) {
	var err error
	soap := newSoapEnvelope()
	soap.Header = header

	soap.Body.Value, err = xml.MarshalIndent(body, "", "    ")
	if err != nil {
		panic("Error marshalling request body " + err.Error())
	}

	requestBody, err := xml.MarshalIndent(soap, "", "    ")
	if err != nil {
		panic("Error marshalling request " + err.Error())
	}

	requestBody = append([]byte(xml.Header), requestBody...)
	c.lastRequest = string(requestBody)
	request, _ := http.NewRequest("POST", c.url, bytes.NewReader(requestBody))
	request.Header.Add("Content-Type", "text/xml; charset=\"utf-8\"")
	request.Header.Add("SOAPAction", soapAction)

	httpResponse, err := c.client.Do(request)
	if err != nil {
		log.Printf("Request error %q", err)
		return nilResponse(), err
	}
	defer httpResponse.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(httpResponse.Body)

	c.lastResponse = buf.String()

	res := &soapResponse{}
	err = xml.Unmarshal(buf.Bytes(), res)
	if err != nil {
		log.Printf("Can't unmarshal soap response %q", err)
		return nilResponse(), err
	}

	fault := &Fault{}
	err = xml.Unmarshal(res.Body.Response, fault)
	if err == nil && fault.FaultCode != "" {
		return nilResponse(), fault
	}

	return res, nil
}

func (c *Client) GetLastCommunications() (string, string) {
	return c.lastRequest, c.lastResponse
}
