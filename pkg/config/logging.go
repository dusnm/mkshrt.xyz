package config

const (
	LoggingFormatText = "text"
	LoggingFormatJson = "json"
)

type (
	Logging struct {
		File   string `mapstructure:"file"`
		Level  string `mapstructure:"level"`
		Format string `mapstructure:"format"`
	}
)
