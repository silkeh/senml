package senml

import (
	"regexp"
	"time"
)

type TestVector struct {
	SkipEncode bool
	JSON, XML  string
	CBOR       []byte
	Result     []Measurement
}

var testVectors = map[string]TestVector{
	"Empty": {
		JSON: `[]`,
		CBOR: []byte{0x80},
		XML:  `<sensml xmlns="urn:ietf:params:xml:ns:senml"></sensml>`,
	},
	"Numeric types": {
		JSON: `[
                 {"v":5},{"v":44},{"v":11290},{"v":739910425},{"v":3177891079421588000},
                 {"v":-6},{"v":-45},{"v":-11291},{"v":-739910426},{"v":-3177891079421588000},
                 {"v":8.1171875},{"v":2.9390623569488525},{"v":1.7},{"v":271.15}
               ]`,
		CBOR: []byte{
			0x8E,
			0xa1, 0x02, 0x05,
			0xa1, 0x02, 0x18, 0x2C,
			0xa1, 0x02, 0x19, 0x2C, 0x1A,
			0xa1, 0x02, 0x1a, 0x2C, 0x1A, 0x23, 0x19,
			0xa1, 0x02, 0x1b, 0x2C, 0x1A, 0x23, 0x19, 0x7b, 0xce, 0x72, 0xd7,
			0xa1, 0x02, 0x25,
			0xa1, 0x02, 0x38, 0x2C,
			0xa1, 0x02, 0x39, 0x2C, 0x1A,
			0xa1, 0x02, 0x3a, 0x2C, 0x1A, 0x23, 0x19,
			0xa1, 0x02, 0x3b, 0x2C, 0x1A, 0x23, 0x19, 0x7b, 0xce, 0x72, 0xd7,
			0xa1, 0x02, 0xf9, 0x48, 0x0f,
			0xa1, 0x02, 0xfa, 0x40, 0x3c, 0x19, 0x99,
			0xa1, 0x02, 0xfb, 0x3f, 0xfb, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33,
			0xa1, 0x02, 0xc4, 0x82, 0x21, 0x19, 0x69, 0xeb,
		},
		Result: []Measurement{
			NewValue("", uint64(5), None, time.Time{}, 0),
			NewValue("", uint64(44), None, time.Time{}, 0),
			NewValue("", uint64(11290), None, time.Time{}, 0),
			NewValue("", uint64(739910425), None, time.Time{}, 0),
			NewValue("", uint64(3177891079421588183), None, time.Time{}, 0),
			NewValue("", int64(-6), None, time.Time{}, 0),
			NewValue("", int64(-45), None, time.Time{}, 0),
			NewValue("", int64(-11291), None, time.Time{}, 0),
			NewValue("", int64(-739910426), None, time.Time{}, 0),
			NewValue("", int64(-3177891079421588184), None, time.Time{}, 0),
			NewValue("", float64(8.1171875), None, time.Time{}, 0),
			NewValue("", float64(2.9390623569488525), None, time.Time{}, 0),
			NewValue("", float64(1.7), None, time.Time{}, 0),
			NewValue("", NewDecimal(-2, 27115), None, time.Time{}, 0),
		},
	},

	// The rest of the test vectors are based on RFC8428
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
		CBOR: []byte{
			0x87, 0xa7, 0x21, 0x78, 0x1b, 0x75, 0x72, 0x6e, 0x3a, 0x64, 0x65, 0x76, 0x3a, 0x6f, 0x77, 0x3a,
			0x31, 0x30, 0x65, 0x32, 0x30, 0x37, 0x33, 0x61, 0x30, 0x31, 0x30, 0x38, 0x30, 0x30, 0x36, 0x3a,
			0x22, 0xfb, 0x41, 0xd3, 0x03, 0xa1, 0x5b, 0x00, 0x10, 0x62, 0x23, 0x61, 0x41, 0x20, 0x05, 0x00,
			0x67, 0x76, 0x6f, 0x6c, 0x74, 0x61, 0x67, 0x65, 0x01, 0x61, 0x56, 0x02, 0xfb, 0x40, 0x5e, 0x06,
			0x66, 0x66, 0x66, 0x66, 0x66, 0xa3, 0x00, 0x67, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x06,
			0x24, 0x02, 0xfb, 0x3f, 0xf3, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0xa3, 0x00, 0x67, 0x63, 0x75,
			0x72, 0x72, 0x65, 0x6e, 0x74, 0x06, 0x23, 0x02, 0xfb, 0x3f, 0xf4, 0xcc, 0xcc, 0xcc, 0xcc, 0xcc,
			0xcd, 0xa3, 0x00, 0x67, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x06, 0x22, 0x02, 0xfb, 0x3f,
			0xf6, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0xa3, 0x00, 0x67, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e,
			0x74, 0x06, 0x21, 0x02, 0xf9, 0x3e, 0x00, 0xa3, 0x00, 0x67, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e,
			0x74, 0x06, 0x20, 0x02, 0xfb, 0x3f, 0xf9, 0x99, 0x99, 0x99, 0x99, 0x99, 0x9a, 0xa3, 0x00, 0x67,
			0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x06, 0x00, 0x02, 0xfb, 0x3f, 0xfb, 0x33, 0x33, 0x33,
			0x33, 0x33, 0x33,
		},
		XML: `<sensml xmlns="urn:ietf:params:xml:ns:senml">
			    <senml bn="urn:dev:ow:10e2073a0108006:" bt="1.276020076001e+09"
			    bu="A" bver="5" n="voltage" u="V" v="120.1"></senml>
			    <senml n="current" t="-5" v="1.2"></senml>
			    <senml n="current" t="-4" v="1.3"></senml>
			    <senml n="current" t="-3" v="1.4"></senml>
			    <senml n="current" t="-2" v="1.5"></senml>
			    <senml n="current" t="-1" v="1.6"></senml>
			    <senml n="current" v="1.7"></senml>
			  </sensml>`,
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
			     {"bn":"urn:dev:ow:10e2073a01080063","bt":1.3200674641e+09,
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
			NewValue("urn:dev:ow:10e2073a01080063", 21.2, RelativeHumidityPercent, floatToTime(1.3200674641e+09+00), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 21.3, RelativeHumidityPercent, floatToTime(1.3200674641e+09+10), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 21.4, RelativeHumidityPercent, floatToTime(1.3200674641e+09+20), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 21.4, RelativeHumidityPercent, floatToTime(1.3200674641e+09+30), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 21.5, RelativeHumidityPercent, floatToTime(1.3200674641e+09+40), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 21.5, RelativeHumidityPercent, floatToTime(1.3200674641e+09+50), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 21.5, RelativeHumidityPercent, floatToTime(1.3200674641e+09+60), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 21.6, RelativeHumidityPercent, floatToTime(1.3200674641e+09+70), 0),
			NewValue("urn:dev:ow:10e2073a01080063", 21.7, RelativeHumidityPercent, floatToTime(1.3200674641e+09+80), 0),
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
