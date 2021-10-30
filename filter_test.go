package zlog_test

import (
	"os"

	"github.com/m-mizutani/zlog"
	"github.com/m-mizutani/zlog/filter"
)

func newExampleLogger() *zlog.Logger {
	logger := zlog.New()
	logger.Emitter = zlog.NewWriterWith(&zlog.ConsoleFormatter{
		TimeFormat: "",
		NoColor:    true,
	}, os.Stdout)
	return logger
}

func ExampleTypeFilter() {
	logger := newExampleLogger()

	type password string
	type myRecord struct {
		ID    string
		EMail password
	}
	record := myRecord{
		ID:    "m-mizutani",
		EMail: "abcd1234",
	}

	logger.Filters = []zlog.Filter{
		filter.Type(password("")),
	}
	logger.With("record", record).Info("Got record")
	// Output:  [info] Got record
	// "record" => zlog_test.myRecord{
	//   ID:    "m-mizutani",
	//   EMail: "[filtered]",
	// }
}

func ExampleValueFilter() {
	logger := newExampleLogger()

	const issuedToken = "abcd1234"
	authHeader := "Authorization: Bearer " + issuedToken

	logger.Filters = []zlog.Filter{
		filter.Value(issuedToken),
	}
	logger.With("auth", authHeader).Info("send header")
	// Output:  [info] send header
	// "auth" => "Authorization: Bearer [filtered]"
}

func ExampleTagFilter() {
	logger := newExampleLogger()

	type myRecord struct {
		ID    string
		EMail string `zlog:"secret"`
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
}

func ExamplePhoneNumberFilter() {
	logger := newExampleLogger()

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
}

func ExampleFieldFilter() {
	logger := newExampleLogger()

	type myRecord struct {
		ID    string
		Phone string
	}
	record := myRecord{
		ID:    "m-mizutani",
		Phone: "090-0000-0000",
	}

	logger.Filters = []zlog.Filter{
		filter.Field("Phone"),
	}
	logger.With("record", record).Info("Got record")
	// Output:  [info] Got record
	// "record" => zlog_test.myRecord{
	//   ID:    "m-mizutani",
	//   Phone: "[filtered]",
	// }
}

func ExampleFieldPrefixFilter() {
	logger := newExampleLogger()

	type myRecord struct {
		ID          string
		SecurePhone string
	}
	record := myRecord{
		ID:          "m-mizutani",
		SecurePhone: "090-0000-0000",
	}

	logger.Filters = []zlog.Filter{
		filter.FieldPrefix("Secure"),
	}
	logger.With("record", record).Info("Got record")
	// Output:  [info] Got record
	// "record" => zlog_test.myRecord{
	//   ID:          "m-mizutani",
	//   SecurePhone: "[filtered]",
	// }
}
