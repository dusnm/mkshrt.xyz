package config

type (
	TLS struct {
		Enabled  bool   `mapstructure:"enabled"`
		CertPath string `mapstructure:"cert_path"`
		KeyPath  string `mapstructure:"key_path"`
	}
)
