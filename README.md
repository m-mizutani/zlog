# zlog [![Go Reference](https://pkg.go.dev/badge/github.com/m-mizutani/zlog.svg)](https://pkg.go.dev/github.com/m-mizutani/zlog) [![Vulnerability scan](https://github.com/m-mizutani/zlog/actions/workflows/trivy.yml/badge.svg)](https://github.com/m-mizutani/zlog/actions/workflows/trivy.yml) [![Unit test](https://github.com/m-mizutani/zlog/actions/workflows/test.yml/badge.svg)](https://github.com/m-mizutani/zlog/actions/workflows/test.yml) [![Security Scan](https://github.com/m-mizutani/zlog/actions/workflows/gosec.yml/badge.svg)](https://github.com/m-mizutani/zlog/actions/workflows/gosec.yml)

A main distinct feature of `zlog` is secure logging that avoid to output secret/sensitive values to log. The feature reduce risk to store secret values (API token, password and such things) and sensitive data like PII (Personal Identifiable Information) such as address, phone number, email address and etc into logging storage.

`zlog` also has major logger features: contextual logging, leveled logging, structured message, show stacktrace of error. See following usage for mote detail.

## Usage

- [Basic example](#basic-example)
	- [Contextual logging](#contextual-logging)
- [Filter sensitive data](#filter-sensitive-data)
	- [By specified value](#by-specified-value)
	- [By custom type](#by-custom-type)
	- [By struct tag](#by-struct-tag)
	- [By data pattern (e.g. personal information)](#by-data-pattern-eg-personal-information)
- [Customize Log output format](#customize-log-output-format)
	- [Change io.Writer](#change-iowriter)
	- [Change formatter](#change-formatter)
	- [Use original emitter](#use-original-emitter)
- [Leveled Logging](#leveled-logging)
- [Error handling](#error-handling)
	- [Output stack trace](#output-stack-trace)
	- [Output error related values](#output-error-related-values)

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

![example](https://user-images.githubusercontent.com/605953/139518107-e1b1deb0-f9c8-4439-b527-7e3ae4e575a0.png)

#### Contextual logging

`Logger.With(key string, value interface{})` method allows contextual logging that output not only message but also related variables. The method saves a pair of key and value and output it by pretty printing (powered by [k0kubun/pp](https://github.com/k0kubun/pp)).


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
### Customize Log output format

zlog has `Emitter` that is interface to output log event. A default emitter is `Writer` that has `Formatter` to format log message, values and error information and `io.Writer` to output formatted log data.

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

You can use your original Emitter that has `Emit(*zlog.Event) error` method.

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

### Leveled Logging

zlog allows for logging at the following levels.

- `trace` (`zlog.LevelTrace`)
- `debug` (`zlog.LevelDebug`)
- `info` (`zlog.LevelInfo`)
- `warn` (`zlog.LevelWarn`)
- `error` (`zlog.LevelError`)

Log level can be changed by modifying `Logger.Level` or calling `Logger.SetLogLevel()` method.

Modifying `Logger.Level` directly:
```go
	logger.Level = zlog.LevelInfo
	logger.Debug("debugging")
	logger.Info("information")
	// Output: [info] information
```

Using `SetLogLevel()` method. Log level is case insensitive.
```go
	logger.SetLogLevel("InFo")

	logger.Debug("debugging")
	logger.Info("information")
	// Output: [info] information
```

### Error handling

`Logger.Err(err error)` outputs not only error message but also stack trace and error related values.

#### Output stack trace

Support below error packages.

- [github.com/pkg/errors](https://github.com/pkg/errors)
- [github.com/m-mizutani/goerr](https://github.com/m-mizutani/goerr)

```go
func crash() error {
	return errors.New("oops")
}

func main() {
	logger := zlog.New()
	if err := crash(); err != nil {
		logger.Err(err).Error("failed")
	}
}

// Output:
// [error] failed
//
// ------------------
// *errors.fundamental: oops
//
// [StackTrace]
// github.com/m-mizutani/zlog_test.crash
// 	/Users/mizutani/.ghq/github.com/m-mizutani/zlog_test/main.go:xx
// github.com/m-mizutani/zlog_test.main
// 	/Users/mizutani/.ghq/github.com/m-mizutani/zlog_test/main.go:xx
// runtime.main
//	/usr/local/Cellar/go/1.17/libexec/src/runtime/proc.go:255
// runtime.goexit
//	/usr/local/Cellar/go/1.17/libexec/src/runtime/asm_amd64.s:1581
// ------------------
```

#### Output error related values

```go
func crash(args string) error {
	return goerr.New("oops").With("args", args)
}

func main() {
	logger := zlog.New()
	if err := crash("hello"); err != nil {
		logger.Err(err).Error("failed")
	}
}

// Output:
// [error] failed
//
// ------------------
// *goerr.Error: oops
//
// (snip)
//
// [Values]
// args => "hello"
// ------------------
```


## License

- MIT License
- Author: Masayoshi Mizutani <mizutani@hey.com>
