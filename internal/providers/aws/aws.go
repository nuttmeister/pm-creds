// Package aws is an provider that can be used by the provider package to
// retrieve aws credentials from the credentials and config files.
package aws

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/nuttmeister/pm-creds/internal/providers/types"
)

const providerType = "aws"

// Create will create a new provider with name based on config and return it.
// Returns *Provider and error.
func Create(name string, config map[string]interface{}) (*Provider, error) {
	creds, err := read(config, "credentials")
	if err != nil {
		return nil, fmt.Errorf("aws: error creating provider %q. %w", name, err)
	}

	configs, err := read(config, "configs")
	if err != nil {
		return nil, fmt.Errorf("aws: error creating provider %q. %w", name, err)
	}

	return &Provider{
		name:    name,
		creds:   creds,
		configs: configs,
	}, nil
}

// read will read option from config and return the value as []string. If value is not set nil
// will be returned. If value is set but isn't an []string nil and error will be returned.
// Returns []string and error.
func read(config map[string]interface{}, option string) ([]string, error) {
	val, exists := config[option]
	if !exists {
		return nil, nil
	}
	slice, ok := val.([]string)
	if !ok {
		return nil, fmt.Errorf("wrong type for option %q (expected []string but got %T)", option, val)
	}

	return slice, nil
}

// Provider fullfills the types.Provider interface and can be used as
// an provider by the providers package.
type Provider struct {
	name string

	creds   []string
	configs []string
}

// Name returns the provider name.
// Returns string.
func (p *Provider) Name() string {
	return p.name
}

// Type returns the provider type.
// Returns string.
func (p *Provider) Type() string {
	return providerType
}

// Get will retrieve profile name from provider p.
// Returns *types.Profile and error.
func (p *Provider) Get(name string) (types.Profile, error) {
	shared, err := config.LoadSharedConfigProfile(
		context.Background(),
		name,
		func(opts *config.LoadSharedConfigOptions) {
			opts.CredentialsFiles = p.creds
			opts.ConfigFiles = p.configs
			opts.Logger = nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("aws: couldn't get credentials for %q from %q. %w", name, p.Name(), err)
	}

	creds := &struct {
		AccessKey    string `json:"accessKey"`
		SecretKey    string `json:"secretKey"`
		SessionToken string `json:"sessionToken,omitempty"`
		Region       string `json:"region,omitempty"`
	}{
		AccessKey:    shared.Credentials.AccessKeyID,
		SecretKey:    shared.Credentials.SecretAccessKey,
		SessionToken: shared.Credentials.SessionToken,
		Region:       shared.Region,
	}

	payload, _ := json.Marshal(creds)

	return &Profile{name: name, payload: payload}, nil
}

// Profile fullfills the types.Profile interface and can be used
// as an profile by the providers package.
type Profile struct {
	name    string
	payload []byte
}

// Name returns the profile name.
// Returns string.
func (p *Profile) Name() string {
	return p.name
}

// Payload returns the profile payload as JSON.
// Returns []byte.
func (p *Profile) Payload() []byte {
	return p.payload
}
