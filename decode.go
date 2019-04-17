package senml

import (
	"fmt"
	"time"
)

// AutoTime toggles the automatic setting of zero timestamps to now.
// Disabling this option results in timestamps relative to zero time when no exact time is given.
var AutoTime = true

// Decode decodes a list of Measurement records into measurement values.
func Decode(records []Record) (list []Measurement, err error) {
	list = make([]Measurement, len(records))

	var now time.Time
	if AutoTime {
		now = time.Now()
	}

	var baseName string
	var baseTime float64
	var baseUnit Unit
	var baseValue float64
	var baseSum float64

	for i, o := range records {
		if o.BaseName != "" {
			baseName = o.BaseName
		}
		if o.BaseTime != 0 {
			baseTime = o.BaseTime
		}
		if o.BaseUnit != "" {
			baseUnit = Unit(o.BaseUnit)
		}
		if o.BaseValue != 0 {
			baseValue = o.BaseValue
		}
		if o.BaseSum != 0 {
			baseSum = o.BaseSum
		}

		var unit Unit
		if o.Unit != "" {
			unit = Unit(o.Unit)
		} else {
			unit = baseUnit
		}

		m := Attributes{
			Name:       baseName + o.Name,
			Unit:       unit,
			Time:       parseTime(baseTime, o.Time, now),
			UpdateTime: floatToDuration(o.UpdateTime),
		}

		switch {
		case o.Value != nil:
			list[i] = &Value{Attributes: m, Value: baseValue + *o.Value}
		case o.Sum != nil:
			list[i] = &Sum{Attributes: m, Value: baseSum + *o.Sum}
		case o.StringValue != "":
			list[i] = &String{Attributes: m, Value: o.StringValue}
		case len(o.DataValue) > 0:
			list[i] = &Data{Attributes: m, Value: o.DataValue}
		case o.BooleanValue != nil:
			list[i] = &Boolean{Attributes: m, Value: *o.BooleanValue}
		default:
			return nil, fmt.Errorf("record has no value attribute set: %#v", o)
		}
	}

	return list, nil
}
