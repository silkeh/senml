package senml

import "encoding/xml"

// Object represents a SenML object.
// This type is used as an intermediary between the Measurement values and the actual encoding.
// All SenML attributes are supported in this object.
type Object struct {
	XMLName xml.Name `json:"-" xml:"senml" codec:"-"`
	_struct      bool    `codec:",int,omitempty"`
	BaseName     string  `json:"bn,omitempty" xml:"bn,attr,omitempty" codec:"-2"`
	BaseTime     float64 `json:"bt,omitempty" xml:"bt,attr,omitempty" codec:"-3"`
	BaseUnit     string  `json:"bu,omitempty" xml:"bu,attr,omitempty" codec:"-4"`
	BaseValue    float64 `json:"bv,omitempty" xml:"bv,attr,omitempty" codec:"-5"`
	BaseSum      float64 `json:"bs,omitempty" xml:"bs,attr,omitempty" codec:"-6"`
	BaseVersion  int     `json:"bver,omitempty" xml:"bver,attr,omitempty" codec:"-1"`
	Name         string  `json:"n,omitempty" xml:"n,attr,omitempty" codec:"0"`
	Unit         string  `json:"u,omitempty" xml:"u,attr,omitempty" codec:"1"`
	Value        *float64 `json:"v,omitempty" xml:"v,attr,omitempty" codec:"2"`
	StringValue  string  `json:"vs,omitempty" xml:"vs,attr,omitempty" codec:"3"`
	BooleanValue *bool    `json:"vb,omitempty" xml:"vb,attr,omitempty" codec:"4"`
	DataValue    []byte  `json:"vd,omitempty" xml:"vd,attr,omitempty" codec:"8"`
	Sum          *float64 `json:"s,omitempty" xml:"s,attr,omitempty" codec:"5"`
	Time         float64 `json:"t,omitempty" xml:"t,attr,omitempty" codec:"6"`
	UpdateTime   float64 `json:"ut,omitempty" xml:"ut,attr,omitempty" codec:"7"`
}
