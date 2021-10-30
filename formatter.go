package zlog

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"github.com/k0kubun/colorstring"
	"github.com/k0kubun/pp"
	"github.com/m-mizutani/goerr"
	"github.com/pkg/errors"
)

type Formatter interface {
	Write(ev *Event, w io.Writer) error
}

type JsonFormatter struct {
	TimeFormat string
}

func NewJsonFormatter() *JsonFormatter {
	return &JsonFormatter{
		TimeFormat: "2006-01-02T15:04:05.000000",
	}
}

type JsonMsg struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Msg       string                 `json:"msg"`
	Values    map[string]interface{} `json:"values,omitempty"`
	Error     *JsonError             `json:"error,omitempty"`
}

type JsonErrorStack struct {
	Function string `json:"function"`
	File     string `json:"file"`
}

type JsonErrorMsg struct {
	Msg  string `json:"msg"`
	Type string `json:"type"`
}

type JsonError struct {
	JsonErrorMsg
	Causes     []*JsonErrorMsg        `json:"causes,omitempty"`
	StackTrace []*JsonErrorStack      `json:"stacktrace,omitempty"`
	Values     map[string]interface{} `json:"values,omitempty"`
}

func newJsonError(err *Error) *JsonError {
	if err == nil {
		return nil
	}

	jerr := &JsonError{
		JsonErrorMsg: JsonErrorMsg{
			Msg:  err.Err.Error(),
			Type: reflect.TypeOf(err.Err).String(),
		},
	}

	cause := err.Err
	for {
		if unwrapped := errors.Cause(cause); cause != unwrapped {
			cause = unwrapped
		} else {
			break
		}

		jerr.Causes = append(jerr.Causes, &JsonErrorMsg{
			Msg:  cause.Error(),
			Type: reflect.TypeOf(cause).String(),
		})
	}

	if len(err.StackTrace) > 0 {
		for _, frame := range err.StackTrace {
			jerr.StackTrace = append(jerr.StackTrace, &JsonErrorStack{
				Function: frame.Function,
				File:     fmt.Sprintf("%s:%d", frame.Filename, frame.Lineno),
			})
		}
	}
	if err.Values != nil && len(err.Values) > 0 {
		jerr.Values = make(map[string]interface{})
		for key, value := range err.Values {
			jerr.Values[key] = value
		}
	}

	return jerr
}

func (x *JsonFormatter) Write(ev *Event, w io.Writer) error {
	m := JsonMsg{
		Timestamp: ev.Timestamp.Format(x.TimeFormat),
		Level:     ev.Level.String(),
		Msg:       ev.Msg,
		Values:    ev.LogEntity.values,
		Error:     newJsonError(ev.err),
	}

	if err := json.NewEncoder(w).Encode(m); err != nil {
		return goerr.Wrap(err)
	}
	return nil
}

type ConsoleFormatter struct {
	TimeFormat string
	NoColor    bool
}

func NewConsoleFormatter() *ConsoleFormatter {
	return &ConsoleFormatter{
		TimeFormat: "15:04:05.000",
		NoColor:    false,
	}
}

var colorMap = map[LogLevel]string{
	LevelTrace: "blue",
	LevelDebug: "cyan",
	LevelInfo:  "white",
	LevelWarn:  "yellow",
	LevelError: "red",
}

func (x *ConsoleFormatter) Write(ev *Event, w io.Writer) error {
	baseFmt := colorstring.Color("%s [[" + colorMap[ev.Level] + "][bold]%s[reset]] %s")
	if x.NoColor {
		baseFmt = "%s [%s] %s"
	}

	base := fmt.Sprintf(baseFmt,
		ev.Timestamp.Format(x.TimeFormat),
		ev.Level.String(),
		ev.Msg)
	if _, err := w.Write([]byte(base)); err != nil {
		return goerr.Wrap(err, "fail to write timestamp")
	}
	if _, err := w.Write([]byte("\n")); err != nil {
		return goerr.Wrap(err, "fail to write console")
	}

	if len(ev.Values()) > 0 {
		for k, v := range ev.Values() {
			if _, err := w.Write([]byte(fmt.Sprintf("\"%s\" => ", k))); err != nil {
				return goerr.Wrap(err, "fail to write console")
			}
			pp.ColoringEnabled = !x.NoColor
			if _, err := pp.Fprint(w, v); err != nil {
				return goerr.Wrap(err)
			}
			if _, err := w.Write([]byte("\n")); err != nil {
				return goerr.Wrap(err, "fail to write console")
			}
		}
	}

	if ev.err != nil {
		errType := reflect.TypeOf(ev.err.Err)

		fmt.Fprintf(w, "\n------------------\n")
		fmt.Fprintf(w, "%s:  %s\n", errType, ev.err.Err.Error())

		if len(ev.err.StackTrace) > 0 {
			fmt.Fprintf(w, "\n[StackTrace]\n")
			for _, frame := range ev.err.StackTrace {
				fmt.Fprintf(w, "%s\n\t%s:%d\n", frame.Function, frame.Filename, frame.Lineno)
			}
		}
		if ev.err.Values != nil && len(ev.err.Values) > 0 {
			fmt.Fprintf(w, "\n[Values]\n")
			for key, value := range ev.err.Values {
				fmt.Fprintf(w, "%s => ", key)
				pp.ColoringEnabled = !x.NoColor
				pp.Fprint(w, value)
				fmt.Fprintf(w, "\n")
			}
		}
		fmt.Fprintf(w, "------------------\n\n")
	}

	return nil
}
