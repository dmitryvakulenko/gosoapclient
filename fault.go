package gosoapclient

import "encoding/xml"

type Fault struct {
	XMLName     xml.Name `xml:"Fault"`
	FaultCode   string   `xml:"faultcode"`
	FaultString string   `xml:"faultstring"`
	FaultActor  string   `xml:"faultactor"`
	Detail      string   `xml:"detail"`
}

func (f *Fault) Error() string {
	return f.FaultString
}