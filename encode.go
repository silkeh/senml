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
	minTime := list[0].Attrs().Time
	baseName := list[0].Attrs().Name
	units := make(map[Unit]int)
	for _, v := range list {
		m := v.Attrs()

		// Maximum time
		if m.Time.Before(minTime) {
			minTime = m.Time
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

	// Calculate base time
	var baseTime float64
	if !minTime.IsZero() {
		baseTime = timeToFloat(minTime)
	} else {
		baseTime = 0
	}
	log.Printf("Min Time: %s => %v", minTime, baseTime)

	// Check base
	var baseUnit string
	if _, ok := units[None]; ok {
		baseUnit = ""
	} else {
		baseUnit = string(maxUnit(units))
	}

	// Clear bases when single object
	if len(list) == 1 {
		baseName = ""
		baseUnit = ""
		baseTime = 0
	}

	// Create objects
	for i, m := range list {
		o := m.Object()

		// Apply base values
		o.Time -= baseTime
		o.Name = o.Name[len(baseName):]
		if o.Unit == baseUnit {
			o.Unit = ""
		}

		objects[i] = o
	}

	// Set base values in first object
	// TODO: BaseValue, BaseSum, BaseVersion
	o := &objects[0]
	o.BaseTime = baseTime
	o.BaseName = baseName
	o.BaseUnit = baseUnit

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
