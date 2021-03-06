package gosoapclient

import "encoding/xml"

func newSoapEnvelope() *envelope {
	return &envelope{
		SoapNamespace: "http://schemas.xmlsoap.org/soap/envelope/",
		Namespaces: []xml.Attr{
			{Name: xml.Name{Local: "xmlns:addr"}, Value: "http://www.w3.org/2005/08/addressing"}}}
}

type envelope struct {
	XMLName       xml.Name    `xml:"SOAP-ENV:Envelope"`
	SoapNamespace string      `xml:"xmlns:SOAP-ENV,attr"`
	Namespaces    []xml.Attr  `xml:",attr"`
	Header        interface{} `xml:"SOAP-ENV:Header"`
	Body          content     `xml:"SOAP-ENV:Body"`
}

type content struct {
	Value []byte `xml:",innerxml"`
}

type soapResponse struct {
	XMLName string      `xml:"Envelope"`
	Header  *headerBody `xml:"Header"`
	Body    *respBody   `xml:"Body"`
}

type headerBody struct {
	Response []byte `xml:",innerxml"`
}

type respBody struct {
	Response []byte `xml:",innerxml"`
}

func nilResponse() *soapResponse {
	return &soapResponse{}
}
