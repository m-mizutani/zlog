# zlog [![Vulnerability scan](https://github.com/m-mizutani/zlog/actions/workflows/trivy.yml/badge.svg)](https://github.com/m-mizutani/zlog/actions/workflows/trivy.yml) [![Unit test](https://github.com/m-mizutani/zlog/actions/workflows/test.yml/badge.svg)](https://github.com/m-mizutani/zlog/actions/workflows/test.yml) [![Security Scan](https://github.com/m-mizutani/zlog/actions/workflows/gosec.yml/badge.svg)](https://github.com/m-mizutani/zlog/actions/workflows/gosec.yml)

Structured logger in Go.

## Usage

<!-- TOC depthfrom:undefined -->

- [Basic example](#basic-example)
- [Customize Log output format](#customize-log-output-format)
    - [Change io.Writer](#change-iowriter)
    - [Change formatter](#change-formatter)
    - [Use original emitter](#use-original-emitter)
- [Filter sensitive data](#filter-sensitive-data)
    - [By specified value](#by-specified-value)
    - [By custom type](#by-custom-type)
    - [By struct tag](#by-struct-tag)
    - [By data pattern (e.g. personal information)](#by-data-pattern-eg-personal-information)

<!-- /TOC -->

### Basic example

```go
import "github.com/m-mizutani/zlog"

type myRecord struct {
	Name  string
	EMail string
}

func main() {
	record := myRecord{
		Name:  "mizutani",
		EMail: "mizutani@hey.com",
	}

	logger := zlog.New()
	logger.With("record", record).Info("hello my logger")
}
```

`zlog.New()` creates a new logger with default settings (console formatter).

![example](https://user-images.githubusercontent.com/605953/135705361-a3edcdb7-58c4-45e7-848c-5086270ad312.png)

### Customize Log output format

#### Change io.Writer

For example, change output to standard error.

```go
logger.Emitter = zlog.NewWriterWith(zlog.NewConsoleFormatter(), os.Stderr)
logger.Info("output to stderr")
```

#### Change formatter

For example, use JsonFormatter to output structured json.

```go
logger.Emitter = zlog.NewWriterWith(zlog.NewJsonFormatter(), os.Stdout)
logger.Info("output as json format")
// Output: {"timestamp":"2021-10-02T14:58:11.791258","level":"info","msg":"output as json format","values":null}
```

#### Use original emitter

`Emitter` is interface to output log event. You can use your original Emitter that has `Emit(*zlog.Event) error` method.

```go

type myEmitter struct {
	seq int
}

func (x *myEmitter) Emit(ev *zlog.Event) error {
	x.seq++
	prefix := []string{"＼(^o^)／", "(´・ω・｀)", "(・∀・)"}
	fmt.Println(prefix[x.seq%3], ev.Msg)
	return nil
}

func ExampleEmitter() {
	logger := zlog.New()
	logger.Emitter = &myEmitter{}

	logger.Info("waiwai")
	logger.Info("heyhey")
	// Output:
	// ＼(^o^)／ waiwai
	// (´・ω・｀) heyhey
}
```

### Filter sensitive data

#### By specified value

```go
	const issuedToken = "abcd1234"
	authHeader := "Authorization: Bearer " + issuedToken

	logger := newExampleLogger()
	logger.Filters = []zlog.Filter{
		filter.Value(issuedToken),
	}
	logger.With("auth", authHeader).Info("send header")
	// Output:  [info] send header
	// "auth" => "Authorization: Bearer [filtered]"
```

#### By custom type

```go
	type password string
	type myRecord struct {
		ID    string
		EMail password
	}
	record := myRecord{
		ID:    "m-mizutani",
		EMail: "abcd1234",
	}

	logger := newExampleLogger()
	logger.Filters = []zlog.Filter{
		filter.Type(password("")),
	}
	logger.With("record", record).Info("Got record")
	// Output:  [info] Got record
	// "record" => zlog_test.myRecord{
	//   ID:    "m-mizutani",
	//   EMail: "[filtered]",
	// }
```

#### By struct tag

```go
	type myRecord struct {
		ID    string
		EMail string `zlog:"secure"`
	}
	record := myRecord{
		ID:    "m-mizutani",
		EMail: "mizutani@hey.com",
	}

	logger.Filters = []zlog.Filter{
		filter.Tag(),
	}
	logger.With("record", record).Info("Got record")
	// Output:  [info] Got record
	// "record" => zlog_test.myRecord{
	//   ID:    "m-mizutani",
	//   EMail: "[filtered]",
	// }
```

#### By data pattern (e.g. personal information)

```go
	type myRecord struct {
		ID    string
		Phone string
	}
	record := myRecord{
		ID:    "m-mizutani",
		Phone: "090-0000-0000",
	}

	logger.Filters = []zlog.Filter{
		filter.PhoneNumber(),
	}
	logger.With("record", record).Info("Got record")
	// Output:  [info] Got record
	// "record" => zlog_test.myRecord{
	//   ID:    "m-mizutani",
	//   Phone: "[filtered]",
	// }
```

## License

- MIT License
- Author: Masayoshi Mizutani <mizutani@hey.com>
