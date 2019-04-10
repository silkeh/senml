package senml

import (
	"encoding/json"
	"log"
)


// Encode encodes a list of measurements to corresponding Measurement objects.
func Encode(list []Measurement) (objects []Object) {
	objects = make([]Object, len(list))

	// Empty list, empty result
	if len(list) == 0 {
		return
	}

	// Analyze the data
	baseTime := list[0].Attrs().Time
	baseName := list[0].Attrs().Name
	units := make(map[Unit]int)
	for _, v := range list {
		m := v.Attrs()

		// Maximum time
		if m.Time.Before(baseTime) {
			baseTime = m.Time
		}

		// Common baseName
		baseName = lcp([]string{baseName, m.Name})

		// Common unit
		if _, ok := units[m.Unit]; ok {
			units[m.Unit] += len(m.Unit)
		} else {
			units[m.Unit] = len(m.Unit)
		}
	}

	// Check base
	var baseUnit Unit
	if _, ok := units[None]; ok {
		baseUnit = None
	} else {
		baseUnit = maxUnit(units)
	}

	// Clear bases when single object
	if len(list) == 1 {
		baseName = ""
		baseUnit = None
		baseTime = time.Time{}
	}

	// Create objects
	for i, m := range list {
		o := m.Object()

		// Apply base values
		if !baseTime.IsZero() {
			o.Time = m.Attrs().Time.Sub(baseTime).Seconds()
		}
		o.Name = o.Name[len(baseName):]
		if o.Unit == string(baseUnit) {
			o.Unit = ""
		}

		objects[i] = o
	}

	// Set base values in first object
	// TODO: BaseValue, BaseSum, BaseVersion
	o := &objects[0]
	o.BaseTime = timeToFloat(baseTime)
	o.BaseName = baseName
	o.BaseUnit = string(baseUnit)

	return
}

// EncodeCBOR encodes a list of measurements into CBOR.
func EncodeCBOR(list []Measurement) ([]byte, error) {
	panic("not implemented")
}

// EncodeJSON encodes a list of measurements into JSON.
func EncodeJSON(list []Measurement) ([]byte, error) {
	return json.Marshal(Encode(list))
}

// EncodeXML encodes a list of measurements into XML.
func EncodeXML(list []Measurement) ([]byte, error) {
	panic("not implemented")
}
