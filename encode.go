package senml

import "time"

// Encode encodes a list of measurements to corresponding Measurement records.
func Encode(list []Measurement) (records []Record) {
	records = make([]Record, len(list))

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

	// Clear bases when single record
	if len(list) == 1 {
		baseName = ""
		baseUnit = None
		baseTime = time.Time{}
	}

	// Create records
	for i, m := range list {
		o := m.Record()

		// Apply base values
		if !baseTime.IsZero() {
			o.Time = m.Attrs().Time.Sub(baseTime).Seconds()
		}
		o.Name = o.Name[len(baseName):]
		if o.Unit == string(baseUnit) {
			o.Unit = ""
		}

		records[i] = o
	}

	// Set base values in first record
	// TODO: BaseValue, BaseSum, BaseVersion
	o := &records[0]
	o.BaseTime = timeToFloat(baseTime)
	o.BaseName = baseName
	o.BaseUnit = string(baseUnit)

	return
}
