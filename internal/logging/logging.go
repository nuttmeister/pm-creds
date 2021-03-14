package logging

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/fatih/color"
)

const dateFormat = "2006-01-02 15:04:05"

var (
	red    = color.New(color.FgRed)
	yellow = color.New(color.FgYellow)
	green  = color.New(color.FgGreen)
)

// Will return \n except if runtime.GOOS is
// windows. Then returns \r\n.
func Lb() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}

type Logger struct {
	messages chan string
	stdOut   io.Writer
	stdErr   io.Writer
}

// New creates a new Logger.
func New() *Logger {
	logger := &Logger{
		messages: make(chan string),
		stdOut:   os.Stdout,
		stdErr:   os.Stderr,
	}
	go logger.print()

	return logger
}

// Print will print message with format and prefix of current date and time.
func (l *Logger) Print(format string, a ...interface{}) {
	l.messages <- fmt.Sprintf("%s: %s", now(), fmt.Sprintf(format, a...))
}

// Notice will print message with format and prefix of current date and time in green color.
func (l *Logger) Notice(format string, a ...interface{}) {
	l.messages <- green.Sprintf("%s: %s", now(), fmt.Sprintf(format, a...))
}

// Warning will print message with format and prefix of current date and time in yellow color.
func (l *Logger) Warning(format string, a ...interface{}) {
	l.messages <- yellow.Sprintf("%s: %s", now(), fmt.Sprintf(format, a...))
}

// Alert will print message with format and prefix of current date and time in red color.
func (l *Logger) Alert(format string, a ...interface{}) {
	l.messages <- red.Sprintf("%s: %s", now(), fmt.Sprintf(format, a...))
}

// Error will print message with format and prefix of current date and time and then exit 1.
func (l *Logger) Error(err error) {
	fmt.Fprintf(l.stdErr, "%s: %s", now(), fmt.Sprintf("%s%s", err.Error(), Lb()))
	os.Exit(1)
}

// print will read from l.messages and print to l.stdOut.
func (l *Logger) print() {
	for {
		msg := <-l.messages
		fmt.Fprint(l.stdOut, msg)
	}
}

// now returns current date and time in dateFormat.
func now() string {
	return time.Now().Format(dateFormat)
}
