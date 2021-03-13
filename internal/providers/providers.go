package providers

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/nuttmeister/pm-creds/internal/providers/aws"
	"github.com/nuttmeister/pm-creds/internal/providers/types"
	"github.com/pelletier/go-toml"
)

// Providers is a map of configured providers.
type Providers map[string]types.Provider

// Load will load providers from file fn.
func Load(fn string) (Providers, error) {
	providers := Providers(map[string]types.Provider{})

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

	return providers, nil
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
	config, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("couldn't read config data in provider %q", name)
	}

	raw, ok := config["type"]
	if !ok {
		return nil, fmt.Errorf("provider %q is missing config option type", name)
	}

	provType, ok := raw.(string)
	if !ok {
		return nil, fmt.Errorf("provider %q has wrong data type for value type", name)
	}

	// Add more providers that satisfies the Provider interface here.
	switch strings.ToLower(provType) {
	case "aws":
		return aws.Create(name, config)
	}

	return nil, fmt.Errorf("provider type %q is not valid", provType)
}
