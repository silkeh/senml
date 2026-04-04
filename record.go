package senml

import (
	"encoding/xml"
)

// Record represents a SenML record.
// This type is used as an intermediary between the Measurement values and the actual encoding.
// All SenML attributes are supported in this record.
type Record struct {
	XMLName      xml.Name `json:"-" xml:"senml" cbor:"-"`
	BaseName     string   `json:"bn,omitempty" xml:"bn,attr,omitempty" cbor:"-2,keyasint,omitempty"`
	BaseTime     Numeric  `json:"bt,omitempty" xml:"bt,attr,omitempty" cbor:"-3,keyasint,omitempty"`
	BaseUnit     string   `json:"bu,omitempty" xml:"bu,attr,omitempty" cbor:"-4,keyasint,omitempty"`
	BaseValue    Numeric  `json:"bv,omitempty" xml:"bv,attr,omitempty" cbor:"-5,keyasint,omitempty"`
	BaseSum      Numeric  `json:"bs,omitempty" xml:"bs,attr,omitempty" cbor:"-6,keyasint,omitempty"`
	BaseVersion  int      `json:"bver,omitempty" xml:"bver,attr,omitempty" cbor:"-1,keyasint,omitempty"`
	Name         string   `json:"n,omitempty" xml:"n,attr,omitempty" cbor:"0,keyasint,omitempty"`
	Unit         string   `json:"u,omitempty" xml:"u,attr,omitempty" cbor:"1,keyasint,omitempty"`
	Value        Numeric  `json:"v,omitempty" xml:"v,attr,omitempty" cbor:"2,keyasint,omitempty"`
	StringValue  string   `json:"vs,omitempty" xml:"vs,attr,omitempty" cbor:"3,keyasint,omitempty"`
	BooleanValue *bool    `json:"vb,omitempty" xml:"vb,attr,omitempty" cbor:"4,keyasint,omitempty"`
	DataValue    []byte   `json:"vd,omitempty" xml:"vd,attr,omitempty" cbor:"8,keyasint,omitempty"`
	Sum          Numeric  `json:"s,omitempty" xml:"s,attr,omitempty" cbor:"5,keyasint,omitempty"`
	Time         Numeric  `json:"t,omitempty" xml:"t,attr,omitempty" cbor:"6,keyasint,omitempty"`
	UpdateTime   Numeric  `json:"ut,omitempty" xml:"ut,attr,omitempty" cbor:"7,keyasint,omitempty"`
}
