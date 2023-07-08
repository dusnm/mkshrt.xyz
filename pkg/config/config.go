package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
	"log"
)

type (
	Config struct {
		Application Application `mapstructure:"application"`
		Database    Database    `mapstructure:"database"`
		Logging     Logging     `mapstructure:"logging"`
	}
)

func (c *Config) update(in fsnotify.Event) {
	slog.Info(in.String())

	if err := viper.Unmarshal(c); err != nil {
		log.Fatalln("Invalid configuration: " + err.Error())
	}
}

func New() (*Config, error) {
	c := &Config{}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$XDG_CONFIG_HOME/mkshrt")
	viper.AddConfigPath("/etc/mkshrt")
	viper.OnConfigChange(c.update)
	viper.WatchConfig()

	if err := viper.ReadInConfig(); err != nil {
		return c, err
	}

	if err := viper.Unmarshal(c); err != nil {
		return c, err
	}

	return c, nil
}
