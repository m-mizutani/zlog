package zlog

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/k0kubun/colorstring"
	"github.com/k0kubun/pp"
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
		Values:    ev.LogEntity.values,
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

func (x *ConsoleFormatter) Write(ev *Event, w io.Writer) error {
	baseFmt := colorstring.Color("%s [[red]%s[reset]] %s")
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

	if len(ev.Values()) > 0 {
		if _, err := w.Write([]byte("\n")); err != nil {
			return goerr.Wrap(err, "fail to write console")
		}

		for k, v := range ev.Values() {
			if _, err := w.Write([]byte(fmt.Sprintf("\"%s\" => ", k))); err != nil {
				return goerr.Wrap(err, "fail to write console")
			}
			pp.ColoringEnabled = !x.NoColor
			if _, err := pp.Fprint(w, v); err != nil {
				return goerr.Wrap(err)
			}
		}
	}

	return nil
}
