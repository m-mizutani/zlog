package zlog

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/m-mizutani/goerr"
	"github.com/pkg/errors"
)

type pkgErrorsStackTracer interface {
	StackTrace() errors.StackTrace
}

type Frame struct {
	Function string
	Filename string
	Lineno   int
}

type Error struct {
	Cause      error
	StackTrace []*Frame
	Values     map[string]any
}

func newError(err error, m *masking) *Error {
	if err == nil {
		return nil
	}

	values := extractErrorValues(err)

	maskedValues := map[string]any{}
	for key, value := range values {
		maskedValues[key] = m.clone(key, reflect.ValueOf(value), "").Interface()
	}

	return &Error{
		Cause:      err,
		StackTrace: extractStackTrace(err),
		Values:     maskedValues,
	}
}

func extractStackTrace(err error) []*Frame {
	switch e := err.(type) {
	case pkgErrorsStackTracer:
		var frames []*Frame
		for _, frame := range e.StackTrace() {
			// Ignore if failed to parse
			l, _ := strconv.ParseInt(fmt.Sprintf("%d", frame), 10, 64)
			f := strings.Split(fmt.Sprintf("%+s", frame), "\n\t")
			frames = append(frames, &Frame{
				Filename: f[1],
				Function: f[0],
				Lineno:   int(l),
			})
		}
		return frames

	case *goerr.Error:
		var frames []*Frame
		for _, stack := range e.Stacks() {
			frames = append(frames, &Frame{
				Function: stack.Func,
				Filename: stack.File,
				Lineno:   stack.Line,
			})
		}
		return frames
	}

	return nil
}

func extractErrorValues(err error) map[string]any {
	var goErr *goerr.Error
	switch {
	case errors.As(err, &goErr):
		values := map[string]any{}
		for k, v := range goErr.Values() {
			values[fmt.Sprintf("%v", k)] = v
		}
		return values
	}
	return nil
}
