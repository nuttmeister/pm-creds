// Package aws is a provider that can be used by the providers package to
// retrieve aws credentials from aws cli credentials and config files.
package aws

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/mitchellh/mapstructure"
	"github.com/nuttmeister/pm-creds/internal/providers/types"
)

// data is used to configure the aws provider.
type data struct {
	Credentials []string `mapstructure:"credentials"`
	Configs     []string `mapstructure:"configs"`
}

// Create will create a new provider with name based on config and return it.
func Create(name string, raw map[string]interface{}) (*Provider, error) {
	data := &data{}
	if err := mapstructure.Decode(raw, data); err != nil {
		return nil, fmt.Errorf("aws: couldn't decode raw to data for %q. %w", name, err)
	}

	return &Provider{
		name:    name,
		creds:   data.Credentials,
		configs: data.Configs,
	}, nil
}

// Provider satisfies the types.Provider interface and can be used as
// a provider by the providers package.
type Provider struct {
	name string

	creds   []string
	configs []string
}

// Name returns the provider name.
func (p *Provider) Name() string {
	return p.name
}

// Get will retrieve profile name from provider p. If name is $env the credentials will be
// retrieved from current environmental variables.
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
func (p *Profile) Name() string {
	return p.name
}

// Payload returns the profile json payload.
func (p *Profile) Payload() []byte {
	return p.payload
}
