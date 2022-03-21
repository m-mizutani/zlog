package zlog

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/k0kubun/colorstring"
	"github.com/k0kubun/pp"
	"github.com/m-mizutani/goerr"
	"github.com/pkg/errors"
)

type Emitter interface {
	Emit(*Log) error
}

// ConsoleEmitter outputs log to console with rich format.
type ConsoleEmitter struct {
	timeFormat string
	noColor    bool
	writer     io.Writer
}

type ConsoleEmitterOption func(x *ConsoleEmitter)

func ConsoleTimeFormat(format string) ConsoleEmitterOption {
	return func(x *ConsoleEmitter) {
		x.timeFormat = format
	}
}
func ConsoleNoColor() ConsoleEmitterOption {
	return func(x *ConsoleEmitter) {
		x.noColor = true
	}
}
func ConsoleWriter(w io.Writer) ConsoleEmitterOption {
	return func(x *ConsoleEmitter) {
		x.writer = w
	}
}

func NewConsoleEmitter(options ...ConsoleEmitterOption) *ConsoleEmitter {
	emitter := &ConsoleEmitter{
		timeFormat: "15:04:05.000",
		noColor:    false,
		writer:     os.Stdout,
	}
	for _, opt := range options {
		opt(emitter)
	}
	return emitter
}

var colorMap = map[LogLevel]string{
	LevelTrace: "blue",
	LevelDebug: "cyan",
	LevelInfo:  "white",
	LevelWarn:  "yellow",
	LevelError: "red",
}

func (x *ConsoleEmitter) Emit(log *Log) error {
	baseFmt := colorstring.Color("%s [[" + colorMap[log.Level] + "][bold]%s[reset]] %s")
	if x.noColor {
		baseFmt = "%s [%s] %s"
	}

	w := x.writer

	base := fmt.Sprintf(baseFmt,
		log.Timestamp.Format(x.timeFormat),
		log.Level.String(),
		log.Msg)
	if _, err := w.Write([]byte(base)); err != nil {
		return goerr.Wrap(err, "fail to write timestamp")
	}
	if _, err := w.Write([]byte("\n")); err != nil {
		return goerr.Wrap(err, "fail to write console")
	}

	if len(log.Values) > 0 {
		for k, v := range log.Values {
			if _, err := w.Write([]byte(fmt.Sprintf("\"%s\" => ", k))); err != nil {
				return goerr.Wrap(err, "fail to write console")
			}
			pp.ColoringEnabled = !x.noColor
			if _, err := pp.Fprint(w, v); err != nil {
				return goerr.Wrap(err)
			}
			if _, err := w.Write([]byte("\n")); err != nil {
				return goerr.Wrap(err, "fail to write console")
			}
		}
		if _, err := w.Write([]byte("\n")); err != nil {
			return goerr.Wrap(err, "fail to write console")
		}
	}

	if log.Error != nil {
		errType := reflect.TypeOf(log.Error)

		fmt.Fprintf(w, "----------------[StackTrace]----------------\n")
		fmt.Fprintf(w, "%s:  %s\n", errType, log.Error.Cause.Error())

		if len(log.Error.StackTrace) > 0 {
			fmt.Fprintf(w, "\n[StackTrace]\n")
			for _, frame := range log.Error.StackTrace {
				fmt.Fprintf(w, "%s\n\t%s:%d\n", frame.Function, frame.Filename, frame.Lineno)
			}
		}
		if log.Error.Values != nil && len(log.Error.Values) > 0 {
			fmt.Fprintf(w, "\n[Values]\n")
			for key, value := range log.Error.Values {
				fmt.Fprintf(w, "%s => ", key)
				pp.ColoringEnabled = !x.noColor
				pp.Fprint(w, value)
				fmt.Fprintf(w, "\n")
			}
		}
		fmt.Fprintf(w, "--------------------------------------------\n")
	}

	return nil
}

// JsonEmitter outputs log as one line JSON text
type JsonEmitter struct {
	timeFormat  string
	writer      io.Writer
	prettyPrint bool
}

func NewJsonEmitter(options ...JsonEmitterOption) *JsonEmitter {
	emitter := &JsonEmitter{
		timeFormat: "2006-01-02T15:04:05.000000",
		writer:     os.Stdout,
	}

	for _, opt := range options {
		opt(emitter)
	}

	return emitter
}

type JsonEmitterOption func(x *JsonEmitter)

func JsonTimeFormat(format string) JsonEmitterOption {
	return func(x *JsonEmitter) {
		x.timeFormat = format
	}
}
func JsonWriter(w io.Writer) JsonEmitterOption {
	return func(x *JsonEmitter) {
		x.writer = w
	}
}
func JsonPrettyPrint() JsonEmitterOption {
	return func(x *JsonEmitter) {
		x.prettyPrint = true
	}
}

type jsonMsg struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Msg       string                 `json:"msg"`
	Values    map[string]interface{} `json:"values,omitempty"`
	Error     *jsonError             `json:"error,omitempty"`
}

type jsonErrorStack struct {
	Function string `json:"function"`
	File     string `json:"file"`
}

type jsonErrorMsg struct {
	Msg  string `json:"msg"`
	Type string `json:"type"`
}

type jsonError struct {
	jsonErrorMsg
	Causes     []*jsonErrorMsg        `json:"causes,omitempty"`
	StackTrace []*jsonErrorStack      `json:"stacktrace,omitempty"`
	Values     map[string]interface{} `json:"values,omitempty"`
}

func newjsonError(err *Error) *jsonError {
	if err == nil {
		return nil
	}

	jerr := &jsonError{
		jsonErrorMsg: jsonErrorMsg{
			Msg:  err.Cause.Error(),
			Type: reflect.TypeOf(err.Cause).String(),
		},
	}

	cause := err.Cause
	for {
		if unwrapped := errors.Cause(cause); cause != unwrapped {
			cause = unwrapped
		} else {
			break
		}

		jerr.Causes = append(jerr.Causes, &jsonErrorMsg{
			Msg:  cause.Error(),
			Type: reflect.TypeOf(cause).String(),
		})
	}

	if len(err.StackTrace) > 0 {
		for _, frame := range err.StackTrace {
			jerr.StackTrace = append(jerr.StackTrace, &jsonErrorStack{
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

func (x *JsonEmitter) Emit(log *Log) error {
	m := jsonMsg{
		Timestamp: log.Timestamp.Format(x.timeFormat),
		Level:     log.Level.String(),
		Msg:       log.Msg,
		Values:    log.Values,
		Error:     newjsonError(log.Error),
	}

	encoder := json.NewEncoder(x.writer)
	if x.prettyPrint {
		encoder.SetIndent("", "  ")
	}
	if err := encoder.Encode(m); err != nil {
		return goerr.Wrap(err)
	}
	return nil
}
