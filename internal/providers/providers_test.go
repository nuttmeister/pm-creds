package providers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var tests = []struct {
	file string
	err  bool
}{
	{
		file: "./testdata/providers.toml",
		err:  false,
	},
	{
		file: "./testdata/no-file.toml",
		err:  true,
	},
	{
		file: "./testdata/error.toml",
		err:  true,
	},
	{
		file: "./testdata/error-no-type.toml",
		err:  true,
	},
	{
		file: "./testdata/error-type.toml",
		err:  true,
	},
	{
		file: "./testdata/error-wrong-type.toml",
		err:  true,
	},
}

func TestLoad(t *testing.T) {
	for _, test := range tests {
		_, err := Load(test.file)
		switch test.err {
		case true:
			assert.Error(t, err)
		case false:
			assert.NoError(t, err)
		}

		// t.Log(providers["aws"].Name())
	}
}
