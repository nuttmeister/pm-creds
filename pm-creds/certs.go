package main

import (
	"fmt"
	"os"

	"github.com/nuttmeister/pm-creds/internal/certs"
	"github.com/nuttmeister/pm-creds/internal/file"
	"github.com/nuttmeister/pm-creds/internal/logging"
	"github.com/nuttmeister/pm-creds/internal/paths"
)

// createCertificates will create certificates if createCerts is true.
// If certificate files already exists and overwrite is true existing
// files will be overwritten.
func createCertificates(createCerts bool, overwrite bool) {
	if createCerts {
		caKeyFile := paths.CaKeyFile(cfgDir)
		caCertFile := paths.CaCertFile(cfgDir)
		serverKeyFile := paths.ServerKeyFile(cfgDir)
		serverCertFile := paths.ServerCertFile(cfgDir)

		if err := os.MkdirAll(paths.CertsDir(cfgDir), 0700); err != nil {
			logger.Error(err)
		}

		certFiles := []string{caKeyFile, caCertFile, serverKeyFile, serverCertFile}
		exists, err := file.CheckFilesExists(certFiles)
		if err != nil {
			logger.Error(err)
		}

		if exists && !overwrite {
			logger.Error(fmt.Errorf("certificate files already exist! use --overwrite or delete them first"))
		}

		if err := certs.Create(caKeyFile, caCertFile, serverKeyFile, serverCertFile); err != nil {
			logger.Error(err)
		}

		logger.Print(
			"generated %q, %q, %q, %q. please setup pm-creds and postman with these%s",
			caKeyFile, caCertFile, serverKeyFile, serverCertFile, logging.Lb(),
		)
		os.Exit(0)
	}
}
