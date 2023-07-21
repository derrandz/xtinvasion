package logger

import (
	"fmt"
	"io"
	"os"
)

// Logger represents the logger component.
type Logger struct {
	output io.Writer
}

// NewLogger creates a new Logger instance.
func NewLogger(output io.Writer) *Logger {
	return &Logger{output: output}
}

// Log prints the provided message to the logger's output.
func (l *Logger) Log(message string) {
	fmt.Fprintln(l.output, message)
}

func (l *Logger) Logf(format string, args ...interface{}) {
	fmt.Fprintf(l.output, format+"\n", args...)
}

func NewFileLogger(filename string) (*Logger, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	return NewLogger(file), nil
}

func NewStdoutLogger() *Logger {
	return NewLogger(os.Stdout)
}
