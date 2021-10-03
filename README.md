# zlog [![Vulnerability scan](https://github.com/m-mizutani/zlog/actions/workflows/trivy.yml/badge.svg)](https://github.com/m-mizutani/zlog/actions/workflows/trivy.yml) [![Unit test](https://github.com/m-mizutani/zlog/actions/workflows/test.yml/badge.svg)](https://github.com/m-mizutani/zlog/actions/workflows/test.yml) [![Security Scan](https://github.com/m-mizutani/zlog/actions/workflows/gosec.yml/badge.svg)](https://github.com/m-mizutani/zlog/actions/workflows/gosec.yml)

Structured logger in Go.

## Usage

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

![example](https://user-images.githubusercontent.com/605953/135705361-a3edcdb7-58c4-45e7-848c-5086270ad312.png)

### Change io.Writer

```go
logger.Writer = os.Stderr
logger.Info("output to stderr")
```


### Change formatter

```go
logger.Formatter = zlog.NewJsonFormatter()
logger.Info("output as json format")
// Output: {"timestamp":"2021-10-02T14:58:11.791258","level":"info","msg":"output as json format","values":null}
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

#### By PII (Personally Identifiable Information) data pattern

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
