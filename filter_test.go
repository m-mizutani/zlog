package zlog_test

import (
	"time"

	"github.com/m-mizutani/zlog"
	"github.com/m-mizutani/zlog/filter"
)

func newExampleLogger(options ...zlog.Option) *zlog.Logger {
	emitter := zlog.NewConsoleEmitter(
		zlog.ConsoleNoColor(),
		zlog.ConsoleTimeFormat(""),
	)
	logger := zlog.New(append(options, zlog.WithEmitter(emitter))...)

	return logger
}

func ExampleTypeFilter() {

	type password string
	type myRecord struct {
		ID    string
		EMail password
	}
	record := myRecord{
		ID:    "m-mizutani",
		EMail: "abcd1234",
	}

	logger := newExampleLogger(
		zlog.WithFilters(filter.Type(password(""))),
	)
	logger.With("record", record).Info("Got record")
	// Output: [info] Got record
	// "record" => zlog_test.myRecord{
	//   ID:    "m-mizutani",
	//   EMail: "[filtered]",
	// }
}

func ExampleValueFilter() {
	const issuedToken = "abcd1234"
	authHeader := "Authorization: Bearer " + issuedToken

	logger := newExampleLogger(zlog.WithFilters(
		filter.Value(issuedToken),
	), zlog.WithClock(func() time.Time {
		return time.Date(2021, 4, 20, 5, 12, 19, 0, time.Local)
	}),
	)

	logger.With("auth", authHeader).Info("send header")
	// Output:  [info] send header
	// "auth" => "Authorization: Bearer [filtered]"
}

func ExampleTagFilter() {
	type myRecord struct {
		ID    string
		EMail string `zlog:"secret"`
	}
	record := myRecord{
		ID:    "m-mizutani",
		EMail: "mizutani@hey.com",
	}

	logger := newExampleLogger(
		zlog.WithFilters(filter.Tag()),
		zlog.WithClock(func() time.Time {
			return time.Date(2021, 4, 20, 5, 12, 19, 0, time.Local)
		}),
	)
	logger.With("record", record).Info("Got record")
	// Output:  [info] Got record
	// "record" => zlog_test.myRecord{
	//   ID:    "m-mizutani",
	//   EMail: "[filtered]",
	// }
}

func ExamplePhoneNumberFilter() {
	type myRecord struct {
		ID    string
		Phone string
	}
	record := myRecord{
		ID:    "m-mizutani",
		Phone: "090-0000-0000",
	}

	logger := newExampleLogger(zlog.WithFilters(filter.PhoneNumber()))
	logger.With("record", record).Info("Got record")
	// Output:  [info] Got record
	// "record" => zlog_test.myRecord{
	//   ID:    "m-mizutani",
	//   Phone: "[filtered]",
	// }
}

func ExampleFieldFilter() {
	type myRecord struct {
		ID    string
		Phone string
	}
	record := myRecord{
		ID:    "m-mizutani",
		Phone: "090-0000-0000",
	}

	logger := newExampleLogger(
		zlog.WithFilters(filter.Field("Phone")),
	)
	logger.With("record", record).Info("Got record")
	// Output:  [info] Got record
	// "record" => zlog_test.myRecord{
	//   ID:    "m-mizutani",
	//   Phone: "[filtered]",
	// }
}

func ExampleFieldPrefixFilter() {
	type myRecord struct {
		ID          string
		SecurePhone string
	}
	record := myRecord{
		ID:          "m-mizutani",
		SecurePhone: "090-0000-0000",
	}

	logger := newExampleLogger(zlog.WithFilters(filter.FieldPrefix("Secure")))

	logger.With("record", record).Info("Got record")
	// Output:  [info] Got record
	// "record" => zlog_test.myRecord{
	//   ID:          "m-mizutani",
	//   SecurePhone: "[filtered]",
	// }
}
