package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type profileTest struct {
}

var tests = []struct {
	name          string
	provType      string
	config        map[string]interface{}
	result        *Provider
	profiles      map[string]*Profile
	profilesError []string
	err           bool
}{
	{
		name:     "aws1",
		provType: "aws",
		config: map[string]interface{}{
			"credentials": []string{"./testdata/credentials"},
		},
		result: &Provider{
			name:  "aws1",
			creds: []string{"./testdata/credentials"},
		},
		profiles: map[string]*Profile{
			"default": {
				name:    "default",
				payload: []byte(`{"accessKey":"key","secretKey":"secret"}`),
			},
		},
		profilesError: []string{"notexist"},
		err:           false,
	},
	{
		name:     "aws2",
		provType: "aws",
		config: map[string]interface{}{
			"configs": []string{"./testdata/configs"},
		},
		result: &Provider{
			name:    "aws2",
			configs: []string{"./testdata/configs"},
		},
		err: false,
	},
	{
		name:     "aws3",
		provType: "aws",
		config: map[string]interface{}{
			"credentials": []string{"./testdata/credentials"},
			"configs":     []string{"./testdata/configs"},
		},
		result: &Provider{
			name:    "aws3",
			creds:   []string{"./testdata/credentials"},
			configs: []string{"./testdata/configs"},
		},
		err: false,
	},
	{
		name: "aws",
		config: map[string]interface{}{
			"credentials": "wrong-type",
		},
		err: true,
	},
	{
		name: "aws",
		config: map[string]interface{}{
			"configs": "wrong-type",
		},
		err: true,
	},
}

func TestLoad(t *testing.T) {
	for _, test := range tests {
		provider, err := Create(test.name, test.config)
		switch test.err {
		case true:
			assert.Error(t, err)
			continue
		case false:
			assert.NoError(t, err)
		}

		assert.Equal(t, test.result, provider)
		assert.Equal(t, test.name, provider.Name())
		assert.Equal(t, test.provType, provider.Type())

		for name, res := range test.profiles {
			profile, err := provider.Get(name)
			assert.NoError(t, err)
			assert.Equal(t, res, profile)
			assert.Equal(t, res.name, profile.Name())
			assert.Equal(t, res.payload, profile.Payload())
		}

		for _, name := range test.profilesError {
			_, err := provider.Get(name)
			assert.Error(t, err)
		}
	}
}
