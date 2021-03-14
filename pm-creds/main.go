package main

import (
	"flag"

	"github.com/nuttmeister/pm-creds/internal/logging"
	"github.com/nuttmeister/pm-creds/internal/providers"
	"github.com/nuttmeister/pm-creds/internal/server"
)

func main() {
	providersFile := "~/.pm-creds/providers.toml"
	configFile := "~/.pm-creds/config.toml"
	certs, overwrite := false, false

	flag.StringVar(&providersFile, "providers", providersFile, "Location of the providers file")
	flag.StringVar(&configFile, "config", configFile, "Location of the config file")
	flag.BoolVar(&certs, "certs", certs, "If certificates should be generated")
	flag.BoolVar(&overwrite, "overwrite", overwrite, "If new certificates should overwrite old (if they exist)")
	flag.Parse()

	logger := logging.New()

	providers, err := providers.Load(providersFile)
	if err != nil {
		logger.Error(err)
	}

	if err := server.Start(configFile, providers, logger); err != nil {
		logger.Error(err)
	}
}
