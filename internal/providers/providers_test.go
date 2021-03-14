package providers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var tests = []struct {
	cfgDir  string
	load    []string
	err     bool
	loadErr bool
}{
	{
		cfgDir: "./testdata/working",
		load:   []string{"aws-full", "aws-credentials", "aws-configs", "aws-default"},
	},
	{
		cfgDir:  "./testdata/working",
		load:    []string{"no-exists"},
		loadErr: true,
	},
	{
		cfgDir: "./testdata/no-dir",
		err:    true,
	},
	{
		cfgDir: "./testdata/error",
		err:    true,
	},
	{
		cfgDir: "./testdata/error-no-type",
		err:    true,
	},
	{
		cfgDir: "./testdata/error-type",
		err:    true,
	},
	{
		cfgDir: "./testdata/error-wrong-type",
		err:    true,
	},
}

func TestLoad(t *testing.T) {
	for _, test := range tests {
		providers, err := Load(test.cfgDir)
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
