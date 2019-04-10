package senml

// Object represents a SenML object.
// This type is used as an intermediary between the Measurement values and the actual encoding.
// All SenML attributes are supported in this object.
type Object struct {
	BaseName     string  `json:"bn,omitempty" cbor:"-2,omitempty"`
	BaseTime     float64 `json:"bt,omitempty" cbor:"-3,omitempty"`
	BaseUnit     string  `json:"bu,omitempty" cbor:"-4,omitempty"`
	BaseValue    float64 `json:"bv,omitempty" cbor:"-5,omitempty"`
	BaseSum      float64 `json:"bs,omitempty" cbor:"-5,omitempty"`
	BaseVersion  int     `json:"bver,omitempty" cbor:"-,omitempty"`
	Name         string  `json:"n,omitempty" cbor:"0,omitempty"`
	Unit         string  `json:"u,omitempty" cbor:"1,omitempty"`
	Value        float64 `json:"v,omitempty" cbor:"2,omitempty"`
	StringValue  string  `json:"vs,omitempty" cbor:"3,omitempty"`
	BooleanValue bool    `json:"vb,omitempty" cbor:"4,omitempty"`
	DataValue    []byte  `json:"vd,omitempty" cbor:"8,omitempty"`
	Sum          float64 `json:"s,omitempty" cbor:"5,omitempty"`
	Time         float64 `json:"t,omitempty" cbor:"6,omitempty"`
	UpdateTime   float64 `json:"ut,omitempty" cbor:"7,omitempty"`
}


