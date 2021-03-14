package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/nuttmeister/pm-creds/internal/logging"
	"github.com/nuttmeister/pm-creds/internal/providers"
	"github.com/nuttmeister/pm-creds/internal/server"
)

var (
	logger  = logging.New()
	home, _ = os.UserHomeDir()
	cfgDir  = filepath.Join(home, ".pm-creds")
)

func main() {
	createConfig, createCerts, overwrite := false, false, false
	flag.StringVar(&cfgDir, "config-dir", cfgDir, "Location of the config files")
	flag.BoolVar(&createConfig, "create-config", createConfig, "If the default config should be created")
	flag.BoolVar(&createCerts, "create-certs", createCerts, "If certificates should be generated")
	flag.BoolVar(&overwrite, "overwrite", overwrite, "If new config/certificates should overwrite old")
	flag.Parse()

	createCertificates(createCerts, overwrite)
	createConfiguration(createConfig, overwrite)
	if createCerts || createConfig {
		os.Exit(0)
	}

	providers, err := providers.Load(cfgDir)
	if err != nil {
		logger.Error(err)
	}

	if err := server.Start(cfgDir, providers, logger); err != nil {
		logger.Error(err)
	}
}
