package server

import (
	"fmt"
	"net/http"
	"strings"
)

// ServerHTTP is used to deliver credentials.
func (cfg *config) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	remote := fmt.Sprintf("%q (%s)", r.RemoteAddr, r.UserAgent())

	if r.Method != "POST" {
		write(w, 400, "text/plain", []byte(fmt.Sprintf("method %q not allowed", r.Method)))
		cfg.logger.Print("method %q not allowed for %s\n", r.Method, remote)
		return
	}

	path := strings.SplitN(strings.TrimPrefix(r.URL.Path, "/"), "/", 2)
	if len(path) != 2 {
		write(w, 400, "text/plain", []byte(`path must be in format "/provider/profile"`))
		cfg.logger.Print("path must be in format %q for %s\n", "/provider/profile", remote)
		return
	}
	providerName, profileName := path[0], path[1]

	// Deny profile if it matches deny slice.
	if match(profileName, cfg.Deny) {
		write(w, 400, "text/plain", []byte(fmt.Sprintf("profile %q has been denied", profileName)))
		cfg.logger.Warning("profile %q has been denied for %s\n", profileName, remote)
		return
	}

	provider, err := cfg.providers.Get(providerName)
	if err != nil {
		write(w, 400, "text/plain", []byte(fmt.Sprintf("no provider named %q", providerName)))
		cfg.logger.Print("no provider named %q for %s\n", providerName, remote)
		return
	}

	profile, err := provider.Get(profileName)
	if err != nil {
		write(w, 400, "text/plain", []byte(fmt.Sprintf("no profile %q in provider %q", profileName, providerName)))
		cfg.logger.Print("no profile %q (%s) for %s\n", profileName, providerName, remote)
		return
	}

	// Skip prompt if profile matches auto-approve slice.
	switch match(profileName, cfg.AutoApprove) {
	case false:
		prompt := fmt.Sprintf("authorize credentials for %q (%s) %s? [y/n]: ", profileName, providerName, remote)
		switch match(profileName, cfg.Warn) {
		case true:
			cfg.logger.Alert(prompt)
		case false:
			cfg.logger.Warning(prompt)
		}

		text, _ := console.ReadString('\n')

		if strings.ToLower(strings.Replace(text, "\n", "", -1)) != "y" {
			write(w, 401, "text/plain", []byte(fmt.Sprintf("authorization to use %q (%s) denied", profileName, providerName)))
			cfg.logger.Warning("denied credentials for %q (%s) %s\n", profileName, providerName, remote)
			return
		}
		cfg.logger.Notice("approved credentials for %q (%s) %s\n", profileName, providerName, remote)

	case true:
		cfg.logger.Notice("auto-approved credentials for %q (%s) %s\n", profileName, providerName, remote)
	}

	write(w, 200, "application/json", profile.Payload())
}

// write will write body to w with content-type ct and status code status.
func write(w http.ResponseWriter, status int, ct string, body []byte) {
	w.Header().Add("Content-Type", ct)
	w.WriteHeader(status)
	w.Write(body)
}

// match will return true if str matches any of the patterns as either
// prefix, suffix or whole match.
func match(str string, patterns []string) bool {
	for _, pattern := range patterns {
		if strings.HasPrefix(str, pattern) {
			return true
		}
		if strings.HasSuffix(str, pattern) {
			return true
		}
	}
	return false
}
