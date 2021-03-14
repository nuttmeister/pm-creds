package paths

import "path/filepath"

// ConfigFile returns the absolute path to the config file based on cfgDir.
func ConfigFile(cfgDir string) string {
	return filepath.Join(cfgDir, "config.toml")
}

// ProvidersFile returns the absolute path to the providers file based on cfgDir.
func ProvidersFile(cfgDir string) string {
	return filepath.Join(cfgDir, "providers.toml")
}

// CertsDir returns the certificate directory.
func CertsDir(cfgDir string) string {
	return filepath.Join(cfgDir, "certs")
}

// CaKeyFile returns the absolute path to the CA Key file based on cfgDir.
func CaKeyFile(cfgDir string) string {
	return filepath.Join(CertsDir(cfgDir), "ca-key.pem")
}

// CaCertFile returns the absolute path to the CA Certificate file based on cfgDir.
func CaCertFile(cfgDir string) string {
	return filepath.Join(CertsDir(cfgDir), "ca-cert.pem")
}

// ServerKeyFile returns the absolute path to the Server Key file based on cfgDir.
func ServerKeyFile(cfgDir string) string {
	return filepath.Join(CertsDir(cfgDir), "server-key.pem")
}

// ServerCertFile returns the absolute path to the Server Certificate file based on cfgDir.
func ServerCertFile(cfgDir string) string {
	return filepath.Join(CertsDir(cfgDir), "server-cert.pem")
}
