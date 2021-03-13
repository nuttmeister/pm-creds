// Package aws is a provider that can be used by the providers package to
// retrieve aws credentials from aws cli credentials and config files.
package aws

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var tests = []struct {
	name          string
	config        map[string]interface{}
	result        *Provider
	profiles      map[string]*Profile
	profilesError []string
	err           bool
}{
	{
		name: "aws1",
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
			"default-with-region": {
				name:    "default-with-region",
				payload: []byte(`{"accessKey":"key","secretKey":"secret","region":"eu-north-1"}`),
			},
			"dev-service": {
				name:    "dev-service",
				payload: []byte(`{"accessKey":"dev-key","secretKey":"dev-secret","sessionToken":"dev-token"}`),
			},
			"service-prod": {
				name:    "service-prod",
				payload: []byte(`{"accessKey":"prod-key","secretKey":"prod-secret","sessionToken":"prod-token"}`),
			},
		},
		profilesError: []string{"notexist"},
		err:           false,
	},
	{
		name: "aws2",
		config: map[string]interface{}{
			"credentials": []string{"./testdata/credentials"},
			"configs":     []string{"./testdata/configs"},
		},
		result: &Provider{
			name:    "aws2",
			creds:   []string{"./testdata/credentials"},
			configs: []string{"./testdata/configs"},
		},
		profiles: map[string]*Profile{
			"default": {
				name:    "default",
				payload: []byte(`{"accessKey":"key","secretKey":"secret"}`),
			},
			"default-with-region": {
				name:    "default-with-region",
				payload: []byte(`{"accessKey":"key","secretKey":"secret","region":"eu-north-1"}`),
			},
			"dev-service": {
				name:    "dev-service",
				payload: []byte(`{"accessKey":"dev-key","secretKey":"dev-secret","sessionToken":"dev-token"}`),
			},
			"service-prod": {
				name:    "service-prod",
				payload: []byte(`{"accessKey":"prod-key","secretKey":"prod-secret","sessionToken":"prod-token","region":"eu-north-1"}`),
			},
		},
		profilesError: []string{"notexist"},
		err:           false,
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

func TestDefaultChain(t *testing.T) {
	os.Setenv("AWS_ACCESS_KEY_ID", "def-key")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "def-secret")
	os.Setenv("AWS_SESSION_TOKEN", "def-token")
	os.Setenv("AWS_DEFAULT_REGION", "eu-north-1")

	provider, err := Create("aws-def", map[string]interface{}{})
	assert.NoError(t, err)
	assert.Equal(t, "aws-def", provider.Name())

	profile, err := provider.Get("$default")
	assert.NoError(t, err)

	assert.Equal(t, "$default", profile.Name())
	assert.Equal(t, []byte(`{"accessKey":"def-key","secretKey":"def-secret","sessionToken":"def-token","region":"eu-north-1"}`), profile.Payload())

	os.Unsetenv("AWS_ACCESS_KEY_ID")
	os.Unsetenv("AWS_SECRET_ACCESS_KEY")
	os.Unsetenv("AWS_SESSION_TOKEN")
	os.Unsetenv("AWS_REGION")
}
