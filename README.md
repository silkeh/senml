Package senml
=============

[![godoc](https://godoc.org/github.com/silkeh/senml?status.svg)](https://godoc.org/github.com/silkeh/senml)
[![cirrus ci](https://api.cirrus-ci.com/github/silkeh/senml.svg)](https://cirrus-ci.com/github/silkeh/senml)
[![goreportcard](https://goreportcard.com/badge/github.com/silkeh/senml)](https://goreportcard.com/report/github.com/silkeh/senml)
[![gocover](http://gocover.io/_badge/github.com/silkeh/senml)](http://gocover.io/github.com/silkeh/senml)

Package `senml` implements the [SenML specification][spec],
used for sending simple sensor measurements encoded in JSON, CBOR or XML.

The goal of this package is to not only support the specification,
but to also make it easy to work with within Go.

[spec]: https://tools.ietf.org/html/rfc8428

Examples
--------

Encoding:

```Go
package main

import (
	"fmt"
	"time"

	"github.com/silkeh/senml"
)

func main() {
	now := time.Now()
	list := []senml.Measurement{
		senml.NewValue("sensor:temperature", 23.5, senml.Celsius, now, 0),
		senml.NewValue("sensor:humidity", 33.7, senml.RelativeHumidityPercent, now, 0),
	}

	data, err := senml.EncodeJSON(list)
	if err != nil {
		fmt.Print("Error encoding to JSON:", err)
	}

	fmt.Printf("%s\n", data)
}
```

Decoding:


```Go
package main

import (
	"fmt"

	"github.com/silkeh/senml"
)

func main() {
	json := []byte(`[{"bn":"sensor:","bt":1555487588,"bu":"Cel","n":"temperature","v":23.5},{"n":"humidity","u":"%RH","v":33.7}]`)

	list, err := senml.DecodeJSON(json)
	if err != nil {
		fmt.Print("Error encoding to JSON:", err)
	}

	for _, m := range list {
		fmt.Printf("%#v\n", m)
	}
}
```
