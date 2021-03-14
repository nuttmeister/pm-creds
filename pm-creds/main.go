package main

import (
	"github.com/nuttmeister/pm-creds/internal/logging"
	"github.com/nuttmeister/pm-creds/internal/providers"
	"github.com/nuttmeister/pm-creds/internal/server"
)

func main() {
	logger := logging.New()
	providers, err := providers.Load("../providers.default.toml")
	if err != nil {
		logger.Error(err)
	}

	if err := server.Start("../config.default.toml", providers, logger); err != nil {
		logger.Error(err)
	}
}
