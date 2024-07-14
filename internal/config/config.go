package config

import (
	"github.com/spf13/viper"

	"github.com/ardihikaru/example-golang-google-chat-sender/pkg/config"
)

type Reader interface {
	Get(key string) string
}

type ViperConfigReader struct {
	viper *viper.Viper
}

type Config struct {
	General   config.General   `mapstructure:"general"`
	Log       config.Log       `mapstructure:"log"`
	Google    config.Google    `mapstructure:"google"`
	Scheduler config.Scheduler `mapstructure:"scheduler"`
}

// Validate validates any miss configurations or missing configs
func (*Config) Validate() error {
	// TODO: implements this logic

	return nil
}

// Get gets config object
func Get() (*Config, error) {
	cfg, err := Load()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// Load loads config from the config.yaml
func Load() (*Config, error) {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var Cfg Config
	err := v.Unmarshal(&Cfg)
	if err != nil {
		return nil, err
	}

	return &Cfg, err
}
