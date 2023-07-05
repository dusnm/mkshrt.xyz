package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

type (
	Config struct {
		Application Application `mapstructure:"application"`
		Database    Database    `mapstructure:"database"`
	}
)

func (c *Config) update(in fsnotify.Event) {
	// TODO: Log event
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
