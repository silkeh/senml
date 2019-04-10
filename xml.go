package senml

import (
	"encoding/xml"
)

const (
	xmlStart = `<sensml xmlns="urn:ietf:params:xml:ns:senml">`
	xmlEnd = `</sensml>`
	xmlNamespace = `urn:ietf:params:xml:ns:senml`
)

type xmlContainer struct {
	XMLName   xml.Name `xml:"sensml" name:"urn:ietf:params:xml:ns:senml"`
	XMLNamespace string `xml:"xmlns,attr"`
	Objs []Object
}

// EncodeXML encodes a list of measurements into XML.
func EncodeXML(list []Measurement) (b []byte, err error) {
	c := xmlContainer{Objs: Encode(list), XMLNamespace: xmlNamespace}
	return xml.Marshal(c)
}

// DecodeXML decodes a list of measurements from XML.
func DecodeXML(x []byte) ([]Measurement, error) {
	c := new(xmlContainer)
	err := xml.Unmarshal(x, c)
	if err != nil {
		return nil, err
	}
	return Decode(c.Objs)
}
