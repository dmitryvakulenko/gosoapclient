package gosoapclient

import "encoding/xml"

func newSoap() *Envelope {
	return &Envelope{SoapNamespace: "http://schemas.xmlsoap.org/soap/envelope/"}
}

type Envelope struct {
	XMLName       xml.Name    `xml:"SOAP-ENV:Envelope"`
	SoapNamespace string      `xml:"xmlns:SOAP-ENV,attr"`
	Namespaces    []xml.Attr  `xml:",attr"`
	Header        interface{} `xml:"SOAP-ENV:Header"`
	Body          Body `xml:"SOAP-ENV:Body"`
}

type Body struct {
	Content interface{}
}