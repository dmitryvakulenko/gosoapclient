package gosoapclient

import (
    "io"
    "net/http"
    "encoding/xml"
    "bytes"
    "crypto/tls"
    "log"
)

type Poster interface {
    Post(url string, contentType string, body io.Reader) (resp *http.Response, err error)
}

type Client struct {
    url          string
    poster       Poster
    client       *http.Client
    lastRequest  string
    lastResponse string
}

func NewClient(url string, poster Poster) *Client {
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{
            InsecureSkipVerify: true,
        },
    }

    return &Client{
        url:    url,
        poster: poster,
        client: &http.Client{Transport: tr},
    }
}

func (c *Client) Call(soapAction string, header, body interface{}) *soapResponse {
    soap := newSoapEnvelope()

    soap.Header = header

    soap.Body.Value, _ = xml.MarshalIndent(body, "", "    ")
    requestBody, err := xml.MarshalIndent(soap, "", "    ")
    if err != nil {
        return nilResponse()
    }

    requestBody = append([]byte(xml.Header), requestBody...)
    c.lastRequest = string(requestBody)
    request, _ := http.NewRequest("POST", c.url, bytes.NewReader(requestBody))
    request.Header.Add("Content-Type", "text/xml; charset=\"utf-8\"")
    request.Header.Add("SOAPAction", soapAction)

    httpResponse, err := c.client.Do(request)
    if err != nil {
        log.Printf("Request error %q", err)
        return nilResponse()
    }
    defer httpResponse.Body.Close()

    buf := new(bytes.Buffer)
    buf.ReadFrom(httpResponse.Body)

    c.lastResponse = buf.String()

    res := &soapResponse{}
    err = xml.Unmarshal(buf.Bytes(), res)
    if err != nil {
        log.Printf("Can't unmarshal soap response %q", err)
        return nilResponse()
    }

    return res
}

func (c *Client) GetLastCommunications() (string, string) {
    return c.lastRequest, c.lastResponse
}