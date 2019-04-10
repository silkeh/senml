package senml

import (
	"encoding/json"
	"testing"
	"time"
)

func TestCompareEncoding(t *testing.T) {
	vals := []Measurement{
		NewValue("test_", 24.27, Celsius, time.Now(), 0),
		NewValue("test_", 99912, Pascal, time.Now(), 0),
		NewValue("test_", 33.55, RelativeHumidityPercent, time.Now(), 0),
		NewValue("test_", 507, Lux, time.Now(), 0),
		NewValue("test_tvoc", 5e-9, None, time.Now(), 0),
		NewValue("test_eco2", 481e-6, None, time.Now(), 0),
		NewValue("test_", 24.27, Celsius, time.Now().Add(30 * time.Second), 0),
		NewValue("test_", 99912, Pascal, time.Now().Add(30 * time.Second), 0),
		NewValue("test_", 33.55, RelativeHumidityPercent,time.Now().Add(30 * time.Second), 0),
		NewValue("test_", 507, Lux, time.Now().Add(30 * time.Second), 0),
		NewValue("test_tvoc", 5e-9, None, time.Now().Add(30 * time.Second), 0),
		NewValue("test_eco2", 481e-6, None, time.Now().Add(30 * time.Second), 0),
	}

	res := Encode(vals)
	b, err := json.Marshal(res)
	if err != nil {
		t.Error(err)
	}

	t.Logf("JSON, %v bytes: %s", len(b), b)
	//_ = codec.NewEncoderBytes(&b, &cborHandle).Encode(&ee)
	//t.Logf("CBOR, %v bytes: %q", len(b), b)
}
