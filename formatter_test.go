package zlog_test

func ExampleFormatter() {
	s := "blue"
	type myRecord struct {
		StrPtr    *string
		StrPtrNil *string
	}
	record := myRecord{
		StrPtr:    &s,
		StrPtrNil: nil,
	}

	consoleLogger := newExampleLogger()
	consoleLogger.With("record", record).Info("test")
	// Output:  [info] test
	// "record" => zlog_test.myRecord{
	//   StrPtr:    &"blue",
	//   StrPtrNil: (*string)(nil),
	// }
}
