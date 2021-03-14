package file

import (
	"fmt"
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockStat(fname string) (fs.FileInfo, error) {
	return nil, fmt.Errorf("mock stat error")
}

var tests = []struct {
	statFn func(string) (fs.FileInfo, error)
	files  []string
	exists bool
	err    bool
}{
	{
		statFn: os.Stat,
		files:  []string{"./testdata/file1.txt"},
		exists: true,
	},
	{
		statFn: os.Stat,
		files:  []string{"./testdata/file1.txt", "./testdata/file2.txt"},
		exists: true,
	},
	{
		statFn: os.Stat,
		files:  []string{"./testdata/file1.txt", "./testdata/file3.txt"},
		exists: true,
	},
	{
		statFn: os.Stat,
		files:  []string{"./testdata/file3.txt", "./testdata/file1.txt"},
		exists: true,
	},
	{
		statFn: os.Stat,
		files:  []string{"./testdata/file3.txt"},
		exists: false,
	},
	{
		statFn: os.Stat,
		files:  []string{"./testdata/file3.txt", "./testdata/file4.txt"},
		exists: false,
	},
	{
		statFn: mockStat,
		files:  []string{"./testdata/file10.txt"},
		exists: false,
		err:    true,
	},
}

func TestCheckFileExists(t *testing.T) {
	for _, test := range tests {
		stat = test.statFn

		exists, err := CheckFilesExists(test.files)
		switch test.err {
		case true:
			assert.Error(t, err)
		case false:
			assert.NoError(t, err)
		}
		assert.Equal(t, test.exists, exists)
	}
}
