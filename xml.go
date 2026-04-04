package senml

import (
	"encoding/xml"
	"fmt"
)

const (
	xmlStart     = `<sensml xmlns="urn:ietf:params:xml:ns:senml">`
	xmlEnd       = `</sensml>`
	xmlNamespace = `urn:ietf:params:xml:ns:senml`
)

// XMLDocument represents an XML senml document.
// It contains a list of senml measurement values in the form or records.
type XMLDocument struct {
	XMLName      xml.Name `xml:"sensml" name:"urn:ietf:params:xml:ns:senml"`
	XMLNamespace string   `xml:"xmlns,attr"`
	Records      []Record `xml:"senml"`
}

// UnmarshalXML implements [xml.Unmarshaler].
func (d *XMLDocument) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var container xmlContainer

	if err := decoder.DecodeElement(&container, &start); err != nil {
		return fmt.Errorf("sensml decode: %w", err)
	}

	d.XMLNamespace = container.XMLNamespace
	d.Records = container.records()

	return nil
}

// EncodeXML encodes a list of measurements into XML.
func EncodeXML(list []Measurement) (b []byte, err error) {
	c := XMLDocument{Records: Encode(list), XMLNamespace: xmlNamespace}
	return xml.Marshal(c)
}

// DecodeXML decodes a list of measurements from XML.
func DecodeXML(x []byte) ([]Measurement, error) {
	c := new(xmlContainer)

	err := xml.Unmarshal(x, c)
	if err != nil {
		return nil, err
	}

	return Decode(c.records())
}

type xmlContainer struct {
	XMLName      xml.Name    `xml:"sensml" name:"urn:ietf:params:xml:ns:senml"`
	XMLNamespace string      `xml:"xmlns,attr"`
	Records      []xmlRecord `xml:"senml"`
}

func (c xmlContainer) records() []Record {
	rs := make([]Record, len(c.Records))

	for i := range c.Records {
		rs[i] = c.Records[i].record()
	}

	return rs
}

type xmlRecord struct {
	Record

	BaseTime   *float64 `xml:"bt,attr,omitempty"`
	BaseValue  *float64 `xml:"bv,attr,omitempty"`
	BaseSum    *float64 `xml:"bs,attr,omitempty"`
	Value      *float64 `xml:"v,attr,omitempty"`
	Sum        *float64 `xml:"s,attr,omitempty"`
	Time       *float64 `xml:"t,attr,omitempty"`
	UpdateTime *float64 `xml:"ut,attr,omitempty"`
}

func (r xmlRecord) record() Record {
	r.Record.BaseTime = xmlFloatValue(r.BaseTime)
	r.Record.BaseValue = xmlFloatValue(r.BaseValue)
	r.Record.BaseSum = xmlFloatValue(r.BaseSum)
	r.Record.Value = xmlFloatValue(r.Value)
	r.Record.Sum = xmlFloatValue(r.Sum)
	r.Record.Time = xmlFloatValue(r.Time)
	r.Record.UpdateTime = xmlFloatValue(r.UpdateTime)

	return r.Record
}

func xmlFloatValue(v *float64) Numeric {
	if v != nil {
		return *v
	}

	return nil
}
