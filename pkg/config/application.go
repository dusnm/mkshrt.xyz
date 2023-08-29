package config

type (
	Application struct {
		Host           string   `mapstructure:"host"`
		Port           uint     `mapstructure:"port"`
		IsBehindProxy  bool     `mapstructure:"behind_proxy"`
		TrustedProxies []string `mapstructure:"trusted_proxies"`
	}
)
