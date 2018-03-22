package gosoapclient

import "encoding/xml"

func NewSoap() *Envelope {
	return &Envelope{SoapNamespace: "http://schemas.xmlsoap.org/soap/envelope/"}
}

type Envelope struct {
	XMLName       xml.Name   `xml:"SOAP-ENV:Header"`
	SoapNamespace string     `xml:"xmlns:SOAP-ENV,attr"`
	Namespaces    []xml.Attr `xml:",attr"`
	Header        Header     `xml:"SOAP-ENV:Header"`
	Body          Body       `xml:"SOAP-ENV:Body"`
}

type Header struct {
	Content interface{} `xml:",omitempty"`
}

type Body struct {
	Content interface{} `xml:",omitempty"`
}
