package logging

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
)

const dateFormat = "2006-01-02 15:04:05"

var (
	red    = color.New(color.FgRed)
	yellow = color.New(color.FgYellow)
	green  = color.New(color.FgGreen)
	rt     = runtime.GOOS
)

// Will return \n except if rt is
// windows. Then returns \r\n.
func Lb() string {
	if rt == "windows" {
		return "\r\n"
	}
	return "\n"
}

type Logger struct {
	messages chan *message
	time     timeInterface
	stdOut   io.Writer
	stdErr   io.Writer
}

// message is used to print text and sends done
// once the text has been printed to l.stdOut.
type message struct {
	text string
	done chan bool
}

// New creates a new Logger.
func New() *Logger {
	logger := &Logger{
		messages: make(chan *message),
		time:     &realTime{},
		stdOut:   colorable.NewColorable(os.Stdout),
		stdErr:   colorable.NewColorable(os.Stderr),
	}
	go logger.print()

	return logger
}

// Print will print message with format and prefix of current date and time.
func (l *Logger) Print(format string, a ...interface{}) {
	msg := &message{
		text: fmt.Sprintf("%s: %s", l.now(), fmt.Sprintf(format, a...)),
		done: make(chan bool),
	}
	l.messages <- msg
	<-msg.done
}

// Notice will print message with format and prefix of current date and time in green color.
func (l *Logger) Notice(format string, a ...interface{}) {
	msg := &message{
		text: green.Sprintf("%s: %s", l.now(), fmt.Sprintf(format, a...)),
		done: make(chan bool),
	}
	l.messages <- msg
	<-msg.done
}

// Warning will print message with format and prefix of current date and time in yellow color.
func (l *Logger) Warning(format string, a ...interface{}) {
	msg := &message{
		text: yellow.Sprintf("%s: %s", l.now(), fmt.Sprintf(format, a...)),
		done: make(chan bool),
	}
	l.messages <- msg
	<-msg.done
}

// Alert will print message with format and prefix of current date and time in red color.
func (l *Logger) Alert(format string, a ...interface{}) {
	msg := &message{
		text: red.Sprintf("%s: %s", l.now(), fmt.Sprintf(format, a...)),
		done: make(chan bool),
	}
	l.messages <- msg
	<-msg.done
}

// Error will print message with format and prefix of current date and time and then exit 1.
func (l *Logger) Error(err error) {
	fmt.Fprintf(l.stdErr, "%s: %s", l.now(), fmt.Sprintf("%s%s", err.Error(), Lb()))
	os.Exit(1)
}

// print will read from l.messages and print to l.stdOut.
func (l *Logger) print() {
	for {
		msg := <-l.messages
		fmt.Fprint(l.stdOut, msg.text)
		msg.done <- true
	}
}

// timeInterface implements Now.
type timeInterface interface {
	Now() time.Time
}

// realTime satisfies the timeInterface.
type realTime struct{}

// Now returns time.Time from time.Now().
func (rt *realTime) Now() time.Time {
	return time.Now()
}

// now returns current date and time in dateFormat.
func (l *Logger) now() string {
	return l.time.Now().Format(dateFormat)
}
