package gosoapclient

import "encoding/xml"

func NewSoapEnvelope() *Envelope {
	return &Envelope{SoapNamespace: "http://schemas.xmlsoap.org/soap/envelope/"}
}

type Envelope struct {
	XMLName       xml.Name    `xml:"SOAP-ENV:Envelope"`
	SoapNamespace string      `xml:"xmlns:SOAP-ENV,attr"`
	Namespaces    []xml.Attr  `xml:",attr"`
	Header        interface{} `xml:"SOAP-ENV:Header"`
	Body          Content     `xml:"SOAP-ENV:Body"`
}

type Content struct {
	Value []byte `xml:",innerxml"`
}

type Response struct {
	XMLName   string `xml:"Envelope"`
	SessionId string `xml:"Header>Session>SessionId"`
	Body      RespBody `xml:"Body"`
}

type RespBody struct {
	Response []byte `xml:",innerxml"`
}
