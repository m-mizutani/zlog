package zlog_test

import (
	"bytes"
	"testing"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/zlog"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func crash1() error {
	return errors.New("oops")
}

func crash2() error {
	return goerr.New("oops").With("param", "value")
}

func TestErrWithPkgErrors(t *testing.T) {
	logger := zlog.New()
	buf := bytes.Buffer{}
	logger.Emitter = zlog.NewWriterWith(&zlog.ConsoleFormatter{
		TimeFormat: "",
		NoColor:    true,
	}, &buf)

	logger.Err(crash1()).Error("bomb!")

	output := buf.String()
	assert.Contains(t, output, "[StackTrace]\ngithub.com/m-mizutani/zlog_test.crash1\n")
	assert.Contains(t, output, "github.com/m-mizutani/zlog/error_test.go:14\n")
	assert.NotContains(t, output, "[Values]\nparam => \"value\"\n")

	// Output:
	// [error] bomb!
	//
	// ------------------
	// *errors.fundamental:  oops
	//
	// [StackTrace]
	// github.com/m-mizutani/zlog_test.crash1
	// 	/Users/mizutani/.ghq/github.com/m-mizutani/zlog/error_test.go:14
	// github.com/m-mizutani/zlog_test.ExampleErrWithPkgErrors
	// 	/Users/mizutani/.ghq/github.com/m-mizutani/zlog/error_test.go:23
	// testing.runExample
	// 	/usr/local/Cellar/go/1.17/libexec/src/testing/run_example.go:64
	// testing.runExamples
	//	/usr/local/Cellar/go/1.17/libexec/src/testing/example.go:44
	// testing.(*M).Run
	//	/usr/local/Cellar/go/1.17/libexec/src/testing/testing.go:1505
	// main.main
	// 	_testmain.go:61
	// runtime.main
	//	/usr/local/Cellar/go/1.17/libexec/src/runtime/proc.go:255
	// runtime.goexit
	//	/usr/local/Cellar/go/1.17/libexec/src/runtime/asm_amd64.s:1581
	// ------------------
}

func TestErrWithPkgErrorsWithJSON(t *testing.T) {
	logger := zlog.New()
	buf := bytes.Buffer{}
	logger.Emitter = zlog.NewWriterWith(&zlog.JsonFormatter{}, &buf)

	logger.Err(errors.Wrap(crash1(), "wrapped")).Error("bomb!")
	t.Log(buf.String())

	assert.Contains(t, buf.String(), `"msg":"wrapped: oops"`)
	assert.Contains(t, buf.String(), `"function":"github.com/m-mizutani/zlog_test.TestErrWithPkgErrorsWithJSON"`)
}

func TestErrWithGoErrWithJSON(t *testing.T) {
	logger := zlog.New()
	buf := bytes.Buffer{}
	logger.Emitter = zlog.NewWriterWith(&zlog.JsonFormatter{}, &buf)

	logger.Err(goerr.Wrap(crash2(), "wrapped")).Error("bomb!")
	t.Log(buf.String())

	assert.Contains(t, buf.String(), `"msg":"wrapped: oops"`)
	assert.Contains(t, buf.String(), `"function":"github.com/m-mizutani/zlog_test.TestErrWithGoErrWithJSON"`)
}

func TestErrWithGoErr(t *testing.T) {
	logger := zlog.New()
	buf := bytes.Buffer{}
	logger.Emitter = zlog.NewWriterWith(&zlog.ConsoleFormatter{
		TimeFormat: "",
		NoColor:    true,
	}, &buf)

	logger.Err(crash2()).Error("bomb!")

	output := buf.String()
	assert.Contains(t, output, "[StackTrace]\ngithub.com/m-mizutani/zlog_test.crash2\n")
	assert.Contains(t, output, "github.com/m-mizutani/zlog/error_test.go:18\n")
	assert.Contains(t, output, "[Values]\nparam => \"value\"\n")

	// Example Output: (Do not use below as Example test because path/version should be changed by environment)
	// [error] bomb!
	//
	// ------------------
	// *goerr.Error: oops
	//
	// [StackTrace]
	// github.com/m-mizutani/zlog_test.crash2
	// 	/Users/mizutani/.ghq/github.com/m-mizutani/zlog/error_test.go:13
	// github.com/m-mizutani/zlog_test.ExampleErrWithGoErr
	// 	/Users/mizutani/.ghq/github.com/m-mizutani/zlog/error_test.go:30
	// testing.runExample
	// 	/usr/local/Cellar/go/1.17/libexec/src/testing/run_example.go:64
	// testing.runExamples
	// 	/usr/local/Cellar/go/1.17/libexec/src/testing/example.go:44
	// testing.(*M).Run
	// 	/usr/local/Cellar/go/1.17/libexec/src/testing/testing.go:1505
	// main.main
	// 	_testmain.go:61
	// runtime.main
	// 	/usr/local/Cellar/go/1.17/libexec/src/runtime/proc.go:255
	// runtime.goexit
	// 	/usr/local/Cellar/go/1.17/libexec/src/runtime/asm_amd64.s:1581
	//
	// [Values]
	// param => "value"
	// ------------------
}
