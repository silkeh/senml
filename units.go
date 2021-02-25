package senml

// SenML unit definitions
const (
	// RFC8428
	None                    Unit = ""
	Meter                   Unit = "m"
	Kilogram                Unit = "kg"
	Gram                    Unit = "g" // not recommended
	Second                  Unit = "s"
	Ampere                  Unit = "A"
	Kelvin                  Unit = "K"
	Candela                 Unit = "cd"
	Mole                    Unit = "mol"
	Hertz                   Unit = "Hz"
	Radian                  Unit = "rad"
	Steradian               Unit = "sr"
	Newton                  Unit = "N"
	Pascal                  Unit = "Pa"
	Joule                   Unit = "J"
	Watt                    Unit = "W"
	Coulomb                 Unit = "C"
	Volt                    Unit = "V"
	Farad                   Unit = "F"
	Ohm                     Unit = "Ohm"
	Siemens                 Unit = "S"
	Weber                   Unit = "Wb"
	Tesla                   Unit = "T"
	Henry                   Unit = "H"
	Celsius                 Unit = "Cel"
	Lumen                   Unit = "lm"
	Lux                     Unit = "lx"
	Becquerel               Unit = "Bq"
	Gray                    Unit = "Gy"
	Sievert                 Unit = "Sv"
	Katal                   Unit = "kat"
	SquareMeter             Unit = "m2"
	CubicMeter              Unit = "m3"
	Liter                   Unit = "l" // not recommended
	MeterPerSecond          Unit = "m/s"
	MeterPerSquareSecond    Unit = "m/s2"
	CubicMeterPerSecond     Unit = "m3/s"
	LiterPerSecond          Unit = "l/s" // not recommended
	WattPerSquareMeter      Unit = "W/m2"
	CandelaPerSquareMeter   Unit = "cd/m2"
	Bit                     Unit = "bit"
	BitPerSecond            Unit = "bit/s"
	Latitude                Unit = "lat"
	Longitude               Unit = "lon"
	PH                      Unit = "pH"
	Decibel                 Unit = "dB"
	DBW                     Unit = "dBW"
	Bel                     Unit = "Bspl" // not recommended
	Count                   Unit = "count"
	Ratio                   Unit = "/"
	Ratio2                  Unit = "%" // not recommended
	RelativeHumidityPercent Unit = "%RH"
	RemainingBatteryPercent Unit = "%EL"
	RemainingBatterySeconds Unit = "EL"
	Rate                    Unit = "1/s"
	RPM                     Unit = "1/min"    // not recommended
	HeartRate               Unit = "beat/min" // not recommended
	HeartBeats              Unit = "beats"    // not recommended
	Conductivity            Unit = "S/m"

	// RFC8798
	Byte                     Unit = "B"
	VoltAmpere               Unit = "VA"
	VoltAmpereSecond         Unit = "VAs"
	VoltAmpereReactive       Unit = "var"
	VoltAmpereReactiveSecond Unit = "vars"
	JoulePerMeter            Unit = "J/m"
	KilogramPerCubicMeter    Unit = "kg/m3"
	Degree                   Unit = "deg" // not recommended

	// ISO 7027-1:2016
	NephelometricTurbidityUnit Unit = "NTU"

	// Secondary units (RFC8798)
	Millisecond            Unit = "ms"
	Minute                 Unit = "min"
	Hour                   Unit = "h"
	Megahertz              Unit = "MHz"
	Kilowatt               Unit = "kW"
	KilovoltAmpere         Unit = "kVA"
	Kilovar                Unit = "kvar"
	AmpereHour             Unit = "Ah"
	WattHour               Unit = "Wh"
	KilowattHour           Unit = "kWh"
	VarHour                Unit = "varh"
	KilovarHour            Unit = "kvarh"
	KilovoltAmpereHour     Unit = "kVAh"
	WattHourPerKilometer   Unit = "Wh/km"
	Kibibyte               Unit = "KiB"
	Gigabyte               Unit = "GB"
	MegabitPerSecond       Unit = "Mbit/s"
	BytePerSecond          Unit = "B/s"
	MegabytePerSecond      Unit = "MB/s"
	Millivolt              Unit = "mV"
	Milliampere            Unit = "mA"
	DecibelMilliwatt       Unit = "dBm"
	MicrogramPerCubicMeter Unit = "ug/m3"
	MillimeterPerHour      Unit = "mm/h"
	MeterPerHour           Unit = "m/h"
	PartsPerMillion        Unit = "ppm"
	Percent                Unit = "/100"
	Permille               Unit = "/1000"
	Hectopascal            Unit = "hPa"
	Millimeter             Unit = "mm"
	Centimeter             Unit = "cm"
	Kilometer              Unit = "km"
	KilometerPerHour       Unit = "km/h"

	// Secondary units (CoRE-1)
	PartsPerBillion  Unit = "ppb"
	PartsPerTrillion Unit = "ppt"
	VoltAmpereHour   Unit = "VAh"
	Milligram        Unit = "mg/l"
	Microgram        Unit = "ug/l"
	GramPerLiter     Unit = "g/l"
)

// Unit represents a SenML defined unit
type Unit string

// String returns the string value of a Unit
func (u Unit) String() string {
	return string(u)
}
