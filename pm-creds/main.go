package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nuttmeister/pm-creds/internal/certs"
	"github.com/nuttmeister/pm-creds/internal/file"
	"github.com/nuttmeister/pm-creds/internal/logging"
	"github.com/nuttmeister/pm-creds/internal/providers"
	"github.com/nuttmeister/pm-creds/internal/server"
)

var logger = logging.New()

func main() {
	dir, _ := os.UserHomeDir()
	providersFile := filepath.Join(dir, "/.pm-creds/providers.toml")
	configFile := filepath.Join(dir, "/.pm-creds/config.toml")
	createCerts, overwrite := false, false

	flag.StringVar(&providersFile, "providers", providersFile, "Location of the providers file")
	flag.StringVar(&configFile, "config", configFile, "Location of the config file")
	flag.BoolVar(&createCerts, "certs", createCerts, "If certificates should be generated")
	flag.BoolVar(&overwrite, "overwrite", overwrite, "If new certificates should overwrite old (if they exist)")
	flag.Parse()

	createCertificates(createCerts, overwrite)

	providers, err := providers.Load(providersFile)
	if err != nil {
		logger.Error(err)
	}

	if err := server.Start(configFile, providers, logger); err != nil {
		logger.Error(err)
	}
}

// createCertificates will create certificates if createCerts is true.
// If certificate files already exists and overwrite is true existing
// files will be overwritten.
func createCertificates(createCerts bool, overwrite bool) {
	if createCerts {
		certFiles := []string{"./ca-key.pem", "./ca-cert.pem", "./server-key.pem", "./server-cert.pem"}
		exists, err := file.CheckFilesExists(certFiles)
		if err != nil {
			logger.Error(err)
		}

		if exists && !overwrite {
			logger.Error(fmt.Errorf("certificate files already exist! use --overwrite or delete them first"))
		}

		if err := certs.Create(certFiles[0], certFiles[1], certFiles[2], certFiles[3]); err != nil {
			logger.Error(err)
		}

		logger.Print(
			"generated %q, %q, %q, %q. please setup pm-creds and postman with these%s",
			certFiles[0], certFiles[1], certFiles[2], certFiles[3], logging.Lb(),
		)
		os.Exit(0)
	}
}
