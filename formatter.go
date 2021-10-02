package zlog

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/k0kubun/colorstring"
	"github.com/m-mizutani/goerr"
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
	Values    map[string]interface{} `json:"values"`
}

func (x *JsonFormatter) Write(ev *Event, w io.Writer) error {
	m := JsonMsg{
		Timestamp: ev.Timestamp.Format(x.TimeFormat),
		Level:     ev.Level.String(),
		Msg:       ev.Msg,
		Values:    ev.Entry.values,
	}
	if err := json.NewEncoder(w).Encode(m); err != nil {
		return goerr.Wrap(err)
	}
	return nil
}

type ConsoleFormatter struct {
	TimeFormat string
}

func NewConsoleFormatter() *ConsoleFormatter {
	return &ConsoleFormatter{
		TimeFormat: "15:04:05.000",
	}
}

func (x *ConsoleFormatter) Write(ev *Event, w io.Writer) error {
	base := fmt.Sprintf(colorstring.Color("%s [[red]%s[reset]] %s"),
		ev.Timestamp.Format(x.TimeFormat),
		ev.Level.String(),
		ev.Msg)
	if _, err := w.Write([]byte(base)); err != nil {
		return goerr.Wrap(err, "fail to write timestamp")
	}

	for k, v := range ev.Entry.values {
		obj, err := json.Marshal(v)
		if err != nil {
			return goerr.Wrap(err, "marshal for console").With(k, v)
		}
		msg := fmt.Sprintf(" %s=%s", k, string(obj))
		if _, err := w.Write([]byte(msg)); err != nil {
			return goerr.Wrap(err, "fail to write timestamp")
		}
	}

	return nil
}
