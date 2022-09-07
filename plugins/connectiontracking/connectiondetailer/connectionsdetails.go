package connectiondetailer

import "encoding/xml"

type ConnectionDetails struct {
	XMLName xml.Name `xml:"conntrack"`
	Flows   []Flow   `xml:"flow"`
}

type Flow struct {
	XMLName xml.Name `xml:"flow"`
	Meta    Meta     `xml:"meta"`
}

type Meta struct {
	XMLName   xml.Name `xml:"meta"`
	Direction string   `xml:"direction,attr"`
	Layer3    Layer3   `xml:"layer3"`
	Layer4    Layer4   `xml:"layer4"`
}

type Layer3 struct {
	XMLName   xml.Name `xml:"layer3"`
	Protonum  string   `xml:"protonum,attr"`
	Protoname string   `xml:"protoname,attr"`
	Src       string   `xml:"src"`
	Dst       string   `xml:"dst"`
}

type Layer4 struct {
	XMLName   xml.Name `xml:"layer4"`
	Protonum  string   `xml:"protonum,attr"`
	Protoname string   `xml:"protoname,attr"`
	Sport     string   `xml:"sport"`
	Dport     string   `xml:"dport"`
}
