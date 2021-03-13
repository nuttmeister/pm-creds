// Package aws is a provider that can be used by the providers package to
// retrieve aws credentials from aws cli credentials and config files.
package aws

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
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
// will be returned. If value is set but isn't a []string nil and error will be returned.
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

// Provider satisfies the types.Provider interface and can be used as
// a provider by the providers package.
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

// Get will retrieve profile name from provider p. If name is $env the credentials will be
// retrieved from current environmental variables.
// Returns *types.Profile and error.
func (p *Provider) Get(name string) (types.Profile, error) {
	creds, region := aws.Credentials{}, ""
	var err error

	switch name {
	case "$default":
		creds, region, err = p.credsFromDefaultChain()
	default:
		creds, region, err = p.credsFromFiles(name)
	}
	if err != nil {
		return nil, fmt.Errorf("aws: couldn't get credentials for %q from %q. %w", name, p.Name(), err)
	}

	raw := &struct {
		AccessKey    string `json:"accessKey"`
		SecretKey    string `json:"secretKey"`
		SessionToken string `json:"sessionToken,omitempty"`
		Region       string `json:"region,omitempty"`
	}{
		AccessKey:    creds.AccessKeyID,
		SecretKey:    creds.SecretAccessKey,
		SessionToken: creds.SessionToken,
		Region:       region,
	}

	payload, _ := json.Marshal(raw)

	return &Profile{name: name, payload: payload}, nil
}

// credsFromFiles will return credentials and region for name from files.
// Returns aws.Credentials and error.
func (p *Provider) credsFromFiles(name string) (aws.Credentials, string, error) {
	opts := func(opts *config.LoadSharedConfigOptions) {
		opts.CredentialsFiles = p.creds
		opts.ConfigFiles = p.configs
		opts.Logger = nil
	}

	ctx := context.Background()
	shared, err := config.LoadSharedConfigProfile(ctx, name, opts)
	if err != nil {
		return aws.Credentials{}, "", err
	}

	return shared.Credentials, shared.Region, nil
}

// credsFromDefaultChain will return credentials and region using the default aws credentials chain.
// The following order will be used to determine the credentials to use.
// Profiles will be read from the default .aws folder.
// 1. Environmental variables.
// 2. Credentials / Config i .aws config. (default or set by AWS_PROFILE/AWS_DEFAULT_PROFILE).
// 3. ECS Task Definition IAM Role.
// 4. EC2 IAM Role.
// Returns aws.Credentials and error.
func (p *Provider) credsFromDefaultChain() (aws.Credentials, string, error) {
	ctx := context.Background()
	def, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return aws.Credentials{}, "", err
	}
	retrived, err := def.Credentials.Retrieve(ctx)
	if err != nil {
		return aws.Credentials{}, "", err
	}

	return retrived, def.Region, nil
}

// Profile satisfies the types.Profile interface and can be used
// as a profile by the providers package.
type Profile struct {
	name    string
	payload []byte
}

// Name returns the profile name.
// Returns string.
func (p *Profile) Name() string {
	return p.name
}

// Payload returns the profile json payload.
// Returns []byte.
func (p *Profile) Payload() []byte {
	return p.payload
}
