package gosoapclient

import "encoding/xml"

type Fault struct {
	Name        xml.Name `xml:"Fault"`
	FaultCode   string   `xml:"faultCode"`
	FaultString string   `xml:"faultString"`
	FaultActor  string   `xml:"faultActor"`
	Detail      string   `xml:"detail"`
}

func (f *Fault) Error() string {
	return f.FaultString
}