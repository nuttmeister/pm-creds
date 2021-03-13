package providers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var tests = []struct {
	file    string
	load    []string
	err     bool
	loadErr bool
}{
	{
		file: "./testdata/providers.toml",
		load: []string{"aws-full", "aws-credentials", "aws-configs", "aws-default"},
	},
	{
		file:    "./testdata/providers.toml",
		load:    []string{"no-exists"},
		loadErr: true,
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
		providers, err := Load(test.file)
		switch test.err {
		case true:
			assert.Error(t, err)
		case false:
			assert.NoError(t, err)
		}

		for _, name := range test.load {
			provider, err := providers.Get(name)
			switch test.loadErr {
			case true:
				assert.Error(t, err)
			case false:
				assert.NoError(t, err)
				assert.Equal(t, name, provider.Name())
			}
		}
	}
}
