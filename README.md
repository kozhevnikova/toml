# TOML Parser and Encoder Library for Go

[![Build Status](https://travis-ci.org/naoina/toml.png?branch=master)](https://travis-ci.org/naoina/toml)

A [TOML](https://github.com/toml-lang/toml) parser and encoder library for [Go](http://golang.org/). 

This library is compatible with TOML version [v0.4.0](https://github.com/toml-lang/toml/blob/master/versions/en/toml-v0.4.0.md).

See [API documentation](http://godoc.org/github.com/naoina/toml).

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Mappings](#mappings)
  - [Rules](#rules)
    - [Exact matching](#exact-matching)
    - [Camelcase matching](#camelcase-matching)
    - [Uppercase matching](#uppercase-matching)      
  - [Value](#value)
    - [String](#string)
    - [Integer](#integer)
    - [Float](#float)
    - [Boolean](#boolean)
    - [Datetime](#datetime)
    - [Array](#array)
    - [Table](#table)
    - [Array of Tables](#array-of-tables)
 - [Using encoding.TextUnmarshaler Interface](#using-encodingtextunmarshaler-interface)
 - [Using toml.UnmarshalerRec Interface](#using-tomlunmarshalerrec-interface)
 - [Using toml.Unmarshaler Interface](#using-tomlunmarshaler-interface)
 - [License](#license)

## Installation

    go get -u github.com/naoina/toml

## Usage

Save the following TOML example as `config.toml`:

```toml
[owner]
name = "Lance Uppercut"
birthday = 1979-05-27

[database]
server = "192.168.1.1"
ports = [ 8001, 8001, 8002 ]
connection_max = 5000
enabled = true

[servers]

    [servers.alpha]
    ip = "10.0.0.1"
    dc = "eqdc10"

    [servers.beta]
    ip = "10.0.0.2"
    dc = "eqdc10"

[clients]
data = [ ["gamma", "delta"], [1, 2] ]
hosts = [
  "alpha",
  "omega"
]
```
Then, use the following Go code to map the preceding TOML to a `config` object.

```go
package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/naoina/toml"
)

type tomlConfig struct {
	Owner struct {
		Name     string
		Birthday time.Time
	}

	Database struct {
		Server        string
		Ports         []int
		ConnectionMax uint
		Enabled       bool
	}

	Servers map[string]ServerInfo

	Clients struct {
		Data  [][]interface{}
		Hosts []string
	}
}

type ServerInfo struct {
	IP net.IP
	DC string
}

func main() {
	file, err := os.Open("config.toml")
	if err != nil {
		log.Println(err)
		return
	}

	defer file.Close()

	var config tomlConfig

	if err := toml.NewDecoder(file).Decode(&config); err != nil {
		log.Println(err)
		return
	}

	fmt.Println("Owner info. Name:", config.Owner.Name, "Birthday:",
		config.Owner.Birthday)

	fmt.Println("Database server:", config.Database.Server)
	fmt.Println("Database port:", config.Database.Ports[0])
	fmt.Println("Database connection_max:", config.Database.ConnectionMax)
	fmt.Println("Database enabled:", config.Database.Enabled)

	fmt.Println("IP of server 'alpha':", config.Servers["alpha"].IP)
	fmt.Println("DC of server 'alpha':", config.Servers["alpha"].DC)

	fmt.Println("IP of server 'beta':", config.Servers["beta"].IP)
	fmt.Println("DC of server 'beta':", config.Servers["beta"].DC)

	fmt.Println("First Host:", config.Clients.Hosts[0])
}

```
Output: 

``` go
Owner info. Name: Lance Uppercut Birthday: 1979-05-27 00:00:00 +0000 UTC
Database server: 192.168.1.1
Database port: 8001 
Database connection_max: 5000
Database enabled: true
IP of server 'alpha': 10.0.0.1
DC of server 'alpha': eqdc10
IP of server 'beta': 10.0.0.2   
DC of server 'beta': eqdc10
First Host: alpha
```

## Mappings

A key value pair of TOML maps to the corresponding field.
The fields of a struct for mapping must be exported.

### Rules 

The rules of the mapping of a key are following:

#### Exact Matching

```toml
timeout_seconds = 256
```

```go
type Config struct {
    Timeout_seconds int
}
```

#### Camelcase Matching

```toml
server_name = "srv1"
```

```go
type Config struct {
    ServerName string
}
```

#### Uppercase Matching

```toml
ip = "10.0.0.1"
```

```go
type Config struct {
    IP string
}
```
### Value 

See the following examples how to map different types of values:

#### String

```toml
val = "string"
```

```go
type Config struct {
    Val string
}
```

#### Integer

```toml
val = 100
```

```go
type Config struct {
    Val int
}
```

Types that can be used:

* int8 (from `-128` to `127`)
* int16 (from `-32768` to `32767`)
* int32 (from `-2147483648` to `2147483647`)
* int64 (from `-9223372036854775808` to `9223372036854775807`)
* int (same as `int32` on 32bit environment, or `int64` on 64bit environment)
* uint8 (from `0` to `255`)
* uint16 (from `0` to `65535`)
* uint32 (from `0` to `4294967295`)
* uint64 (from `0` to `18446744073709551615`)
* uint (same as `uint` on 32bit environment, or `uint64` on 64bit environment)

#### Float

```toml
val = 3.1415
```

```go
type Config struct {
    Val float32
}
```

Types that can be used:

* float32
* float64

#### Boolean

```toml
val = true
```

```go
type Config struct {
    Val bool
}
```

#### Datetime

```toml
val = 2014-09-28T21:27:39Z
```

```go
type Config struct {
    Val time.Time
}
```

#### Array

```toml
val = ["a", "b", "c"]
```

```go
type Config struct {
    Val []string
}
```

The following examples also can be mapped:

```toml
val1 = [1, 2, 3]
val2 = [["a", "b"], ["c", "d"]]
val3 = [[1, 2, 3], ["a", "b", "c"]]
val4 = [[1, 2, 3], [["a", "b"], [true, false]]]
```

```go
type Config struct {
    Val1 []int
    Val2 [][]string
    Val3 [][]interface{}
    Val4 [][]interface{}
}
```

#### Table

```toml
[server]
type = "app"

  [server.development]
  ip = "10.0.0.1"

  [server.production]
  ip = "10.0.0.2"
```

```go
type Config struct {
    Server map[string]Server
}

type Server struct {
    IP string
}
```

You can also use the following struct instead of map of struct.

```go
type Config struct {
    Server struct {
        Development Server
	Production Server
    }
}

type Server struct {
    IP string
}
```

#### Array of Tables

```toml
[[fruit]]
  name = "apple"

  [fruit.physical]
    color = "red"
    shape = "round"

  [[fruit.variety]]
    name = "red delicious"

  [[fruit.variety]]
    name = "granny smith"

[[fruit]]
  name = "banana"

  [[fruit.variety]]
    name = "plantain"
```

```go
type Config struct {
    Fruit []struct {
        Name string
	Physical struct {
	    Color string
	    Shape string
	}
	Variety []struct {
	    Name string
	}
    }
}
```

### Using `encoding.TextUnmarshaler` Interface

The package toml supports `encoding.TextUnmarshaler` (and `encoding.TextMarshaler`). You can
use it to apply custom marshaling rules for certain types. The `UnmarshalText` method is
called with the value text found in the TOML input. TOML strings are passed unquoted.

```toml
duration = "10s"
```

```go
import time

type Duration time.Duration

// UnmarshalText implements encoding.TextUnmarshaler
func (d *Duration) UnmarshalText(data []byte) error {
    duration, err := time.ParseDuration(string(data))
    if err == nil {
        *d = Duration(duration)
    }
    
    return err
}

// MarshalText implements encoding.TextMarshaler
func (d Duration) MarshalText() ([]byte, error) {
    return []byte(time.Duration(d).String()), nil
}

type ConfigWithDuration struct {
    Duration Duration
}
```
### Using `toml.UnmarshalerRec` Interface

You can also override marshaling rules specifically for TOML using the `UnmarshalerRec`
and `MarshalerRec` interfaces. These are useful if you want to control how structs or
arrays are handled. You can apply an additional validation or set unexported struct fields.

Note: `encoding.TextUnmarshaler` and `encoding.TextMarshaler` should be preferred for
simple (scalar) values because they're also compatible with other formats like JSON or
YAML.

[See the UnmarshalerRec example](https://godoc.org/github.com/naoina/toml/#example_UnmarshalerRec).

### Using `toml.Unmarshaler` Interface

If you want to deal with raw TOML syntax, use the `Unmarshaler` and `Marshaler`
interfaces. Their input and output is raw TOML syntax. As such, these interfaces are
useful if you want to handle TOML at the syntax level.

[See the Unmarshaler example](https://godoc.org/github.com/naoina/toml/#example_Unmarshaler).

## License

MIT
