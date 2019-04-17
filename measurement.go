package senml

import (
	"bytes"
	"time"
)

// Measurement represents a single SenML measurement value.
// This interface is meant to represent the various Measurement values,
// see: Value, Sum, String, Boolean and Data.
type Measurement interface {
	// Attrs returns a pointer to the measurement of the Measurement value.
	Attrs() *Attributes

	// Equal returns true if the given Measurement value is equal.
	Equal(Measurement) bool

	// Record returns a SenML record representing the value.
	Record() Record
}

// Attributes contains the properties common to a measurement.
type Attributes struct {
	Name       string
	Unit       Unit
	Time       time.Time
	UpdateTime time.Duration
}

// Attrs returns a pointer to the measurement of the Measurement value.
func (m *Attributes) Attrs() *Attributes {
	return m
}

// Equal returns true if the given attribute values are equal.
func (m *Attributes) Equal(s *Attributes) bool {
	return m.Name == s.Name && m.Unit == s.Unit && m.Time.Equal(s.Time) && m.UpdateTime == s.UpdateTime
}

// Record returns a SenML record representing the value.
func (m *Attributes) Record() Record {
	var t float64
	if m.Time.IsZero() {
		t = 0
	} else {
		t = float64(m.Time.Unix())
	}

	return Record{
		Name:       m.Name,
		Unit:       string(m.Unit),
		Time:       t,
		UpdateTime: m.UpdateTime.Seconds(),
	}
}

// Value represents a floating point measurement value.
// It implements Measurement.
type Value struct {
	Attributes
	Value float64
}

// NewValue returns a new Value with the corresponding value and attributes.
func NewValue(name string, value float64, unit Unit, time time.Time, updateTime time.Duration) *Value {
	return &Value{
		Attributes: Attributes{
			Name:       name,
			Unit:       unit,
			Time:       time,
			UpdateTime: updateTime,
		},
		Value: value,
	}
}

// Equal returns true if the given Measurement value is equal.
func (v *Value) Equal(ml Measurement) bool {
	b, ok := ml.(*Value)
	if !ok {
		return false
	}
	return v.Attributes.Equal(&b.Attributes) && v.Value == b.Value
}

// Record returns a SenML record representing the value.
func (v *Value) Record() Record {
	s := v.Attributes.Record()
	s.Value = &v.Value
	return s
}

// Sum represents an integrated floating point measurement value.
// It implements Measurement.
type Sum struct {
	Attributes
	Value float64
}

// NewSum returns a new Sum value with the corresponding value and attributes.
func NewSum(name string, sum float64, unit Unit, time time.Time, updateTime time.Duration) *Sum {
	return &Sum{
		Attributes: Attributes{
			Name:       name,
			Unit:       unit,
			Time:       time,
			UpdateTime: updateTime,
		},
		Value: sum,
	}
}

// Equal returns true if the given Measurement value is equal.
func (v *Sum) Equal(ml Measurement) bool {
	b, ok := ml.(*Sum)
	if !ok {
		return false
	}
	return v.Attributes.Equal(&b.Attributes) && v.Value == b.Value
}

// Record returns a SenML record representing the value.
func (v *Sum) Record() Record {
	s := v.Attributes.Record()
	s.Sum = &v.Value
	return s
}

// String represents a string measurement value.
// It implements Measurement.
type String struct {
	Attributes
	Value string
}

// NewString returns a new String value with the corresponding value and attributes.
func NewString(name string, value string, unit Unit, time time.Time, updateTime time.Duration) *String {
	return &String{
		Attributes: Attributes{
			Name:       name,
			Unit:       unit,
			Time:       time,
			UpdateTime: updateTime,
		},
		Value: value,
	}
}

// Equal returns true if the given Measurement value is equal.
func (v *String) Equal(ml Measurement) bool {
	b, ok := ml.(*String)
	if !ok {
		return false
	}
	return v.Attributes.Equal(&b.Attributes) && v.Value == b.Value
}

// Record returns a SenML record representing the value.
func (v *String) Record() Record {
	s := v.Attributes.Record()
	s.StringValue = v.Value
	return s
}

// Boolean represents a boolean measurement value.
// It implements Measurement.
type Boolean struct {
	Attributes
	Value bool
}

// NewBoolean returns a new Boolean value with the corresponding value and attributes.
func NewBoolean(name string, value bool, unit Unit, time time.Time, updateTime time.Duration) *Boolean {
	return &Boolean{
		Attributes: Attributes{
			Name:       name,
			Unit:       unit,
			Time:       time,
			UpdateTime: updateTime,
		},
		Value: value,
	}
}

// Equal returns true if the given Measurement value is equal.
func (v *Boolean) Equal(ml Measurement) bool {
	b, ok := ml.(*Boolean)
	if !ok {
		return false
	}
	return v.Attributes.Equal(&b.Attributes) && v.Value == b.Value
}

// Record returns a SenML record representing the value.
func (v *Boolean) Record() Record {
	s := v.Attributes.Record()
	s.BooleanValue = &v.Value
	return s
}

// Data represents a measurement value returning binary data.
// It implements Measurement.
type Data struct {
	Attributes
	Value []byte
}

// NewData returns a new Data value with the corresponding value and attributes.
func NewData(name string, value []byte, unit Unit, time time.Time, updateTime time.Duration) *Data {
	return &Data{
		Attributes: Attributes{
			Name:       name,
			Unit:       unit,
			Time:       time,
			UpdateTime: updateTime,
		},
		Value: value,
	}
}

// Equal returns true if the given Measurement value is equal.
func (v *Data) Equal(ml Measurement) bool {
	b, ok := ml.(*Data)
	if !ok {
		return false
	}
	return v.Attributes.Equal(&b.Attributes) && bytes.Equal(v.Value, b.Value)
}

// Record returns a SenML record representing the value.
func (v *Data) Record() Record {
	s := v.Attributes.Record()
	s.DataValue = v.Value
	return s
}
