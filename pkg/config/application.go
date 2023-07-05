package config

type (
	Application struct {
		Host           string   `mapstructure:"host"`
		Port           uint     `mapstructure:"port"`
		Domain         string   `mapstructure:"domain"`
		IsBehindProxy  bool     `mapstructure:"behind_proxy"`
		TrustedProxies []string `mapstructure:"trusted_proxies"`
	}
)
