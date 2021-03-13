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
	creds, err := read(config, "credentials", nil)
	if err != nil {
		return nil, fmt.Errorf("aws: error creating provider %q. %w", name, err)
	}

	configs, err := read(config, "configs", nil)
	if err != nil {
		return nil, fmt.Errorf("aws: error creating provider %q. %w", name, err)
	}

	return &Provider{
		name:    name,
		creds:   creds,
		configs: configs,
	}, nil
}

// read will read option from config and return the value as []string.
// If value is not set def will be returned. If value is set but isn't
// an []string nil and error will be returned.
// Returns []string and error.
func read(config map[string]interface{}, option string, def []string) ([]string, error) {
	val, exists := config[option]
	if !exists {
		return def, nil
	}
	slice, ok := val.([]string)
	if !ok {
		return nil, fmt.Errorf("wrong type for option %q (expected []string but got %T)", option, val)
	}

	return slice, nil
}

type Provider struct {
	name string

	creds   []string
	configs []string
}

func (p *Provider) Name() string {
	return p.name
}

func (p *Provider) Type() string {
	return providerType
}

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

	creds := &credentials{
		AccessKey:    shared.Credentials.AccessKeyID,
		SecretKey:    shared.Credentials.SecretAccessKey,
		SessionToken: shared.Credentials.SessionToken,
		Region:       shared.Region,
	}

	payload, _ := json.Marshal(creds)

	return &Profile{name: name, payload: payload}, nil
}

type Profile struct {
	name    string
	payload []byte
}

func (p *Profile) Name() string {
	return p.name
}

func (p *Profile) Payload() []byte {
	return p.payload
}

type credentials struct {
	AccessKey    string `json:"accessKey"`
	SecretKey    string `json:"secretKey"`
	SessionToken string `json:"sessionToken,omitempty"`
	Region       string `json:"region,omitempty"`
}
