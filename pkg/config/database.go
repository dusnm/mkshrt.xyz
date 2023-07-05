package config

type (
	Database struct {
		Host     string `mapstructure:"host"`
		Port     uint   `mapstructure:"port"`
		Name     string `mapstructure:"name"`
		Username string `mapstructure:"user"`
		Password string `mapstructure:"pass"`
		Charset  string `mapstructure:"charset"`
	}
)
