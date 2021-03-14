package main

import (
	"fmt"
	"os"

	"github.com/nuttmeister/pm-creds/internal/file"
	"github.com/nuttmeister/pm-creds/internal/logging"
	"github.com/nuttmeister/pm-creds/internal/paths"
)

var configDefault = []byte(`port             = 9999
profiles-warn    = [ "-prod", "-production", "prod-", "production-" ]
profiles-approve = [ "-dev", "development", "dev-", "development-" ]
`)

var providersDefault = []byte("\n")

func createConfiguration(createConfig bool, overwrite bool) {
	if createConfig {
		if err := os.MkdirAll(cfgDir, 0700); err != nil {
			logger.Error(err)
		}
	}

	cfgFile := paths.ConfigFile(cfgDir)
	providersFile := paths.ProvidersFile(cfgDir)

	cfgFiles := []string{cfgFile, providersFile}
	exists, err := file.CheckFilesExists(cfgFiles)
	if err != nil {
		logger.Error(err)
	}

	if exists && !overwrite {
		logger.Error(fmt.Errorf("config/providers files already exist! use --overwrite or delete them first"))
	}

	if err := os.WriteFile(cfgFile, configDefault, 0600); err != nil {
		logger.Error(fmt.Errorf("couldn't write default config to %q. %w", cfgFile, err))
	}

	if err := os.WriteFile(providersFile, providersDefault, 0600); err != nil {
		logger.Error(fmt.Errorf("couldn't write default providers to %q. %w", providersFile, err))
	}

	logger.Print("wrote default config to %q and default providers to %q%s", cfgFile, providersFile, logging.Lb())
	os.Exit(0)
}
