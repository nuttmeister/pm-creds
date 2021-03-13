// Package providers is used to load all the supported credential providers
// that can be used to get credentials.
package providers

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/nuttmeister/pm-creds/internal/providers/aws"
	"github.com/nuttmeister/pm-creds/internal/providers/types"
	"github.com/pelletier/go-toml"
)

// Providers contains all providers loaded.
type Providers struct {
	providers map[string]types.Provider
}

// Get will return the provider with name or error if it doesn't exists.
func (p *Providers) Get(name string) (types.Provider, error) {
	provider, ok := p.providers[name]
	if !ok {
		return nil, fmt.Errorf("providers: provider %q doesn't exists", name)
	}

	return provider, nil
}

// Load will load and create providers from file fn.
func Load(fn string) (*Providers, error) {
	providers := map[string]types.Provider{}

	rawProviders, err := loadFile(fn)
	if err != nil {
		return nil, err
	}

	for name, data := range rawProviders {
		provider, err := parseProvider(name, data)
		if err != nil {
			return nil, err
		}
		providers[name] = provider
	}

	return &Providers{providers: providers}, nil
}

// loadFile will read from file fn and toml unmarshal it's content.
func loadFile(fn string) (map[string]interface{}, error) {
	file, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, fmt.Errorf("providers: couldn't read providers file %q. %w", fn, err)
	}

	rawProviders := map[string]interface{}{}
	if err := toml.Unmarshal(file, &rawProviders); err != nil {
		return nil, fmt.Errorf("providers: couldn't toml unmarshal file %q. %w", fn, err)
	}

	return rawProviders, nil
}

// parseProvider will parse the provider data to make sure it satisfies the minimum data
// required for it to be created. It will also call the corrept provider package
// depending on what type was set in the data provided.
func parseProvider(name string, data interface{}) (types.Provider, error) {
	raw, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("couldn't read config of provider %q", name)
	}

	cfg := &struct {
		Type string `mapstructure:"type"`
	}{}
	if err := mapstructure.Decode(raw, cfg); err != nil {
		return nil, fmt.Errorf("couldn't decode field %s from data for %q. %w", "type", name, err)
	}

	// Add more providers that satisfies the Provider interface here.
	switch strings.ToLower(cfg.Type) {
	case "aws":
		return aws.Create(name, raw)
	}

	return nil, fmt.Errorf("provider %q has an invalid %q", cfg.Type, "type")
}
