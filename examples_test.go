package senml

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"
)

type TestVector struct {
	SkipEncode  bool
	JSON   string
	CBOR   []byte
	Result []Measurement
}

var examples = map[string]TestVector{
	"Single Data Point": {
		JSON: `[{"n":"urn:dev:ow:10e2073a01080063","u":"Cel","v":23.1}]`,
		Result: []Measurement{
			NewValue("urn:dev:ow:10e2073a01080063", 23.1, Celsius, time.Time{}, 0),
		},
	},
	"Multiple Data Points 1": {
		JSON: `[
                 {"bn":"urn:dev:ow:10e2073a01080063:","n":"voltage","u":"V","v":120.1},
                 {"n":"current","u":"A","v":1.2}
               ]`,
		Result: []Measurement{
			NewValue("urn:dev:ow:10e2073a01080063:voltage", 120.1, Volt, time.Time{}, 0),
			NewValue("urn:dev:ow:10e2073a01080063:current", 1.2, Ampere, time.Time{}, 0),
		},
	},
	"Multiple Data Points 2": {
		SkipEncode: true, // Encode uses positive increments, base version not supported
		JSON: `[
			     {"bn":"urn:dev:ow:10e2073a0108006:","bt":1.276020076001e+09,
			      "bu":"A","bver":5,
			      "n":"voltage","u":"V","v":120.1},
			     {"n":"current","t":-5,"v":1.2},
			     {"n":"current","t":-4,"v":1.3},
			     {"n":"current","t":-3,"v":1.4},
			     {"n":"current","t":-2,"v":1.5},
			     {"n":"current","t":-1,"v":1.6},
			     {"n":"current","v":1.7}
			   ]`,
		Result: []Measurement{
			NewValue("urn:dev:ow:10e2073a0108006:voltage", 120.1, Volt, floatToTime(1.276020076001e+09), 0),
			NewValue("urn:dev:ow:10e2073a0108006:current", 1.2, Ampere, floatToTime(1.276020076001e+09-5), 0),
			NewValue("urn:dev:ow:10e2073a0108006:current", 1.3, Ampere, floatToTime(1.276020076001e+09-4), 0),
			NewValue("urn:dev:ow:10e2073a0108006:current", 1.4, Ampere, floatToTime(1.276020076001e+09-3), 0),
			NewValue("urn:dev:ow:10e2073a0108006:current", 1.5, Ampere, floatToTime(1.276020076001e+09-2), 0),
			NewValue("urn:dev:ow:10e2073a0108006:current", 1.6, Ampere, floatToTime(1.276020076001e+09-1), 0),
			NewValue("urn:dev:ow:10e2073a0108006:current", 1.7, Ampere, floatToTime(1.276020076001e+09-0), 0),
		},
	},
	"Multiple Data Points 3": {
		JSON: `[
			     {"bn":"urn:dev:ow:10e2073a01080063","bt":1.320067464e+09,
			      "bu":"%RH","v":21.2},
			     {"t":10,"v":21.3},
			     {"t":20,"v":21.4},
			     {"t":30,"v":21.4},
			     {"t":40,"v":21.5},
			     {"t":50,"v":21.5},
			     {"t":60,"v":21.5},
			     {"t":70,"v":21.6},
			     {"t":80,"v":21.7}
			]`,
		Result: []Measurement{
			NewValue("urn:dev:ow:10e2073a01080063", 21.2, RelativeHumidityPercent, floatToTime(1.320067464e+09+00), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 21.3, RelativeHumidityPercent, floatToTime(1.320067464e+09+10), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 21.4, RelativeHumidityPercent, floatToTime(1.320067464e+09+20), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 21.4, RelativeHumidityPercent, floatToTime(1.320067464e+09+30), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 21.5, RelativeHumidityPercent, floatToTime(1.320067464e+09+40), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 21.5, RelativeHumidityPercent, floatToTime(1.320067464e+09+50), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 21.5, RelativeHumidityPercent, floatToTime(1.320067464e+09+60), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 21.6, RelativeHumidityPercent, floatToTime(1.320067464e+09+70), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 21.7, RelativeHumidityPercent, floatToTime(1.320067464e+09+80), 0),
		},
	},
	"Multiple Measurements": {
		SkipEncode: true, // Encode chooses 'lat' as base unit, as '%RH', 'lon' and 'lat' save the same amount of bytes
		JSON: `[
			     {"bn":"urn:dev:ow:10e2073a01080063","bt":1.320067464e+09,
			      "bu":"%RH","v":20},
			     {"u":"lon","v":24.30621},
			     {"u":"lat","v":60.07965},
			     {"t":60,"v":20.3},
			     {"u":"lon","t":60,"v":24.30622},
			     {"u":"lat","t":60,"v":60.07965},
			     {"t":120,"v":20.7},
			     {"u":"lon","t":120,"v":24.30623},
			     {"u":"lat","t":120,"v":60.07966},
			     {"u":"%EL","t":150,"v":98},
			     {"t":180,"v":21.2},
			     {"u":"lon","t":180,"v":24.30628},
			     {"u":"lat","t":180,"v":60.07967}
			   ]`,
		Result: []Measurement{
			NewValue("urn:dev:ow:10e2073a01080063", 20, RelativeHumidityPercent, floatToTime(1.320067464e+09), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 24.30621, Longitude, floatToTime(1.320067464e+09), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 60.07965, Latitude, floatToTime(1.320067464e+09), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 20.3, RelativeHumidityPercent, floatToTime(1.320067464e+09+60), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 24.30622, Longitude, floatToTime(1.320067464e+09+60), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 60.07965, Latitude, floatToTime(1.320067464e+09+60), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 20.7, RelativeHumidityPercent, floatToTime(1.320067464e+09+120), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 24.30623, Longitude, floatToTime(1.320067464e+09+120), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 60.07966, Latitude, floatToTime(1.320067464e+09+120), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 98, RemainingBatteryPercent, floatToTime(1.320067464e+09+150), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 21.2, RelativeHumidityPercent, floatToTime(1.320067464e+09+180), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 24.30628, Longitude, floatToTime(1.320067464e+09+180), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 60.07967, Latitude, floatToTime(1.320067464e+09+180), 0),
		},
	},
	"Resolved Data": {
		SkipEncode: true, // Encode chooses 'lat' as base unit, as '%RH', 'lon' and 'lat' save the same amount of bytes
		JSON: `[
			     {"n":"urn:dev:ow:10e2073a01080063","u":"%RH","t":1.320067464e+09,
			      "v":20},
			     {"n":"urn:dev:ow:10e2073a01080063","u":"lon","t":1.320067464e+09,
			      "v":24.30621},
			     {"n":"urn:dev:ow:10e2073a01080063","u":"lat","t":1.320067464e+09,
			      "v":60.07965},
			     {"n":"urn:dev:ow:10e2073a01080063","u":"%RH","t":1.320067524e+09,
			      "v":20.3},
			     {"n":"urn:dev:ow:10e2073a01080063","u":"lon","t":1.320067524e+09,
			      "v":24.30622},
			     {"n":"urn:dev:ow:10e2073a01080063","u":"lat","t":1.320067524e+09,
			      "v":60.07965},
			     {"n":"urn:dev:ow:10e2073a01080063","u":"%RH","t":1.320067584e+09,
			      "v":20.7},
			     {"n":"urn:dev:ow:10e2073a01080063","u":"lon","t":1.320067584e+09,
			      "v":24.30623},
			     {"n":"urn:dev:ow:10e2073a01080063","u":"lat","t":1.320067584e+09,
			      "v":60.07966},
			     {"n":"urn:dev:ow:10e2073a01080063","u":"%EL","t":1.320067614e+09,
			      "v":98},
			     {"n":"urn:dev:ow:10e2073a01080063","u":"%RH","t":1.320067644e+09,
			      "v":21.2},
			     {"n":"urn:dev:ow:10e2073a01080063","u":"lon","t":1.320067644e+09,
			      "v":24.30628},
			     {"n":"urn:dev:ow:10e2073a01080063","u":"lat","t":1.320067644e+09,
			      "v":60.07967}
			   ]`,
		Result: []Measurement{
			NewValue("urn:dev:ow:10e2073a01080063", 20, RelativeHumidityPercent, floatToTime(1.320067464e+09), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 24.30621, Longitude, floatToTime(1.320067464e+09), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 60.07965, Latitude, floatToTime(1.320067464e+09), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 20.3, RelativeHumidityPercent, floatToTime(1.320067464e+09+60), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 24.30622, Longitude, floatToTime(1.320067464e+09+60), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 60.07965, Latitude, floatToTime(1.320067464e+09+60), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 20.7, RelativeHumidityPercent, floatToTime(1.320067464e+09+120), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 24.30623, Longitude, floatToTime(1.320067464e+09+120), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 60.07966, Latitude, floatToTime(1.320067464e+09+120), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 98, RemainingBatteryPercent, floatToTime(1.320067464e+09+150), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 21.2, RelativeHumidityPercent, floatToTime(1.320067464e+09+180), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 24.30628, Longitude, floatToTime(1.320067464e+09+180), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 60.07967, Latitude, floatToTime(1.320067464e+09+180), 0),
		},
	},
	"Multiple Data Types": {
		SkipEncode: true, // Encode chooses 'lat' as base unit, as '%RH', 'lon' and 'lat' save the same amount of bytes
		JSON: `[
		    {"bn":"urn:dev:ow:10e2073a01080063:","n":"temp","u":"Cel","v":23.1},
		    {"n":"label","vs":"Machine Room"},
		    {"n":"open","vb":false},
		    {"n":"nfc-reader","vd":"aGkgCg=="}
		  ]`,
		Result: []Measurement{
			NewValue("urn:dev:ow:10e2073a01080063:temp", 23.1, Celsius, time.Time{}, 0),
			NewString("urn:dev:ow:10e2073a01080063:label", "Machine Room", None, time.Time{}, 0),
			NewBoolean("urn:dev:ow:10e2073a01080063:open", false, None, time.Time{}, 0),
			NewData("urn:dev:ow:10e2073a01080063:nfc-reader", []byte("hi \n"), None, time.Time{}, 0),
		},
	},
}

var regexpWhitespace = regexp.MustCompile(`\s`)

func TestExamplesDecodeJSON(t *testing.T) {
	for n, example := range examples {
		res, err := DecodeJSON([]byte(example.JSON))
		if err != nil {
			t.Errorf("Decode error in example %s: %s", n, err)
			continue
		}

		if !equal(res, example.Result) {
			t.Errorf("Decode for example %s incorrect, got:\n%s\nexpected:\n%s", n, toString(res), toString(example.Result))
		}
	}
}

func TestExamplesEncode(t *testing.T) {
	for n, example := range examples {
		if example.SkipEncode {
			continue
		}

		var exp []Object
		err := json.Unmarshal([]byte(example.JSON), &exp)
		if err != nil {
			t.Errorf("JSON error in example %s: %s", n, err)
			continue
		}

		res := Encode(example.Result)
		if !reflect.DeepEqual(exp, res) {
			t.Errorf("Encode for example %s incorrect, got:\n%#v\nexpected:\n%#v", n, res, exp)
		}
	}
}


func toString(ml []Measurement) string {
	strs := make([]string, len(ml))
	for i, m := range ml {
		strs[i] = fmt.Sprintf("%#v", m)
	}

	return strings.Join(strs, ",\n")
}

func equal(a, b []Measurement) bool {
	for i := range a {
		if !a[i].Equal(b[i]) {
			log.Printf("\n%#v\n%#v", a[i], b[i])
			return false
		}
	}
	return true
}
