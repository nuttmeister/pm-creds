package server

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/nuttmeister/pm-creds/internal/logging"
	"github.com/nuttmeister/pm-creds/internal/paths"
	"github.com/nuttmeister/pm-creds/internal/providers"
	"github.com/pelletier/go-toml"
)

const listen = "localhost:%d"

// config contains the basic configuration for the http server and it's handler.
type config struct {
	certificate   string
	key           string
	caCertificate string

	Port int `mapstructure:"port"`

	AutoApprove []string `mapstructure:"profiles-approve"`
	Warn        []string `mapstructure:"profiles-warn"`
	Deny        []string `mapstructure:"profiles-deny"`

	providers *providers.Providers
	logger    *logging.Logger
}

// Start will start the http server using config file in cfgDir and providers.
func Start(cfgDir string, providers *providers.Providers, logger *logging.Logger) error {
	cfg, err := loadConfig(cfgDir)
	if err != nil {
		return fmt.Errorf("server: couldn't load config. %w", err)
	}
	cfg.providers = providers
	cfg.logger = logger

	ca, err := caPool(cfg.caCertificate)
	if err != nil {
		return fmt.Errorf("server: couldn't create ca pool. %w", err)
	}

	http.HandleFunc("/", cfg.ServerHTTP)
	server := &http.Server{
		Addr: fmt.Sprintf(listen, cfg.Port),
		TLSConfig: &tls.Config{
			ClientCAs:  ca,
			ClientAuth: tls.RequireAndVerifyClientCert,
		},
	}

	cfg.logger.Print("starting listening on https://%s%s", fmt.Sprintf(listen, cfg.Port), logging.Lb())
	if err := server.ListenAndServeTLS(cfg.certificate, cfg.key); err != nil {
		return fmt.Errorf("server: http server error. %w", err)
	}

	return nil
}

// caPool will return a new cert pool only containing cert from fn.
func caPool(fn string) (*x509.CertPool, error) {
	raw, err := os.ReadFile(fn)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("file %q doesn't exist", fn)
		}
		return nil, fmt.Errorf("couldn't read file %q. %w", fn, err)
	}

	pool := x509.NewCertPool()
	if ok := pool.AppendCertsFromPEM(raw); !ok {
		return nil, fmt.Errorf("couldn't add ca cert %q to pool. %w", fn, err)
	}

	return pool, nil
}

// loadConfig will read from file fn and toml unmarshal it's content into config.
func loadConfig(cfgDir string) (*config, error) {
	fn := paths.ConfigFile(cfgDir)

	file, err := os.ReadFile(fn)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("file %q doesn't exist. use --create-config to create default config", fn)
		}
		return nil, fmt.Errorf("couldn't read file %q. %w", fn, err)
	}

	raw := map[string]interface{}{}
	if err := toml.Unmarshal(file, &raw); err != nil {
		return nil, fmt.Errorf("couldn't toml unmarshal file %q. %w", fn, err)
	}

	cfg := &config{}
	if err := mapstructure.Decode(raw, cfg); err != nil {
		return nil, fmt.Errorf("couldn't decode raw to config for %q. %w", fn, err)
	}

	// Set certificates.
	cfg.caCertificate = paths.CaCertFile(cfgDir)
	cfg.key = paths.ServerKeyFile(cfgDir)
	cfg.certificate = paths.ServerCertFile(cfgDir)

	return cfg, nil
}
