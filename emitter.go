package zlog

import (
	"io"
	"os"
)

type Emitter interface {
	Emit(*Event) error
}

type Writer struct {
	Formatter Formatter
	Output    io.Writer
}

func NewWriter() *Writer {
	return NewWriterWith(NewConsoleFormatter(), os.Stdout)
}

func NewWriterWith(formatter Formatter, output io.Writer) *Writer {
	return &Writer{
		Formatter: formatter,
		Output:    output,
	}
}

func (x *Writer) Emit(ev *Event) error {
	if x.Formatter == nil {
		panic("Writer.Formatter is not set")
	}
	if x.Output == nil {
		panic("Writer.Output is not set")
	}

	return x.Formatter.Write(ev, x.Output)
}
