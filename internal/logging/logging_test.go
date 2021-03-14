package logging

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
)

var tests = []struct {
	time   string
	format string
	data   []interface{}
	output []byte
}{
	{
		time:   "2020-01-15 10:01:59",
		format: "test",
		output: []byte("2020-01-15 10:01:59: test"),
	},
	{
		time:   "2020-01-15 15:01:59",
		format: "test %d and %s and %q",
		data:   []interface{}{11, "example", "quoted"},
		output: []byte(`2020-01-15 15:01:59: test 11 and example and "quoted"`),
	},
}

var testErrors = []struct {
	time   string
	err    error
	output []byte
}{
	{
		time:   "2021-01-01 18:59:05",
		err:    fmt.Errorf("generic test error"),
		output: []byte("2021-01-01 18:59:05: generic test error" + Lb()),
	},
}

type timeMock struct {
	time time.Time
}

func (tm *timeMock) Now() time.Time {
	return tm.time
}

func TestPrint(t *testing.T) {
	color.NoColor = true
	std := &bytes.Buffer{}

	logger := New()
	logger.stdOut = std

	for _, test := range tests {
		testTime, err := time.Parse("2006-01-02 15:04:05", test.time)
		logger.time = &timeMock{time: testTime}
		assert.NoError(t, err)
		logger.Print(test.format, test.data...)
		assert.Equal(t, string(test.output), string(std.Bytes()))
		std.Reset()
	}
}

func TestNotice(t *testing.T) {
	color.NoColor = true
	std := &bytes.Buffer{}

	logger := New()
	logger.stdOut = std

	for _, test := range tests {
		testTime, err := time.Parse("2006-01-02 15:04:05", test.time)
		logger.time = &timeMock{time: testTime}
		assert.NoError(t, err)
		logger.Notice(test.format, test.data...)
		assert.Equal(t, string(test.output), string(std.Bytes()))
		std.Reset()
	}
}

func TestWarning(t *testing.T) {
	color.NoColor = true
	std := &bytes.Buffer{}

	logger := New()
	logger.stdOut = std

	for _, test := range tests {
		testTime, err := time.Parse("2006-01-02 15:04:05", test.time)
		logger.time = &timeMock{time: testTime}
		assert.NoError(t, err)
		logger.Warning(test.format, test.data...)
		assert.Equal(t, string(test.output), string(std.Bytes()))
		std.Reset()
	}
}

func TestAlert(t *testing.T) {
	color.NoColor = true
	std := &bytes.Buffer{}

	logger := New()
	logger.stdOut = std

	for _, test := range tests {
		testTime, err := time.Parse("2006-01-02 15:04:05", test.time)
		logger.time = &timeMock{time: testTime}
		assert.NoError(t, err)
		logger.Alert(test.format, test.data...)
		assert.Equal(t, string(test.output), string(std.Bytes()))
		std.Reset()
	}
}

func TestError(t *testing.T) {
	color.NoColor = true
	logger := New()

	for _, test := range testErrors {
		if os.Getenv("RUN_EXIT") == "1" {
			testTime, err := time.Parse("2006-01-02 15:04:05", test.time)
			logger.time = &timeMock{time: testTime}
			assert.NoError(t, err)
			logger.Error(test.err)
			return
		}

		errOut := &bytes.Buffer{}
		cmd := exec.Command(os.Args[0], "-test.run=TestError")
		cmd.Stderr = errOut
		cmd.Env = append(os.Environ(), "RUN_EXIT=1")
		err := cmd.Run()

		assert.Error(t, err)
		assert.Equal(t, test.output, errOut.Bytes())
		errOut.Reset()
	}
}

func TestLb(t *testing.T) {
	old := rt

	rt = "windows"
	nl := Lb()
	assert.Equal(t, "\r\n", nl)
	rt = "darwin"
	nl = Lb()
	assert.Equal(t, "\n", nl)
	rt = "linux"
	nl = Lb()
	assert.Equal(t, "\n", nl)

	rt = old
}

func TestNow(t *testing.T) {
	std := &bytes.Buffer{}

	testTime, err := time.Parse("2006-01-02 15:04:05", "2020-01-01 01:02:03")
	assert.NoError(t, err)

	logger := New()
	logger.stdOut = std
	logger.time = &timeMock{time: testTime}

	str := logger.now()
	assert.Equal(t, testTime.Format("2006-01-02 15:04:05"), str)
}

func TestNowRealTime(t *testing.T) {
	logger := New()
	str := logger.now()
	_, err := time.Parse("2006-01-02 15:04:05", str)
	assert.NoError(t, err)
}
