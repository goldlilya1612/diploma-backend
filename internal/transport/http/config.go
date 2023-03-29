package http

import (
	"github.com/spf13/viper"
)

type Config struct {
	Host   string `mapstructure:"GIN_HOST"`
	Port   string `mapstructure:"GIN_PORT"`
	Scheme string `mapstructure:"GIN_SCHEME"`

	AllowOrigins     []string `mapstructure:"GIN_ALLOW_ORIGINS"`
	AllowMethods     []string `mapstructure:"GIN_ALLOW_METHODS"`
	AllowHeaders     []string `mapstructure:"GIN_ALLOW_HEADERS"`
	ExposeHeaders    []string `mapstructure:"GIN_EXPOSE_HEADERS"`
	AllowCredentials bool     `mapstructure:"GIN_ALLOW_CREDENTIALS"`
}

func NewConfigFromEnv(path, name string) (*Config, error) {

	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType("env")

	viper.AllowEmptyEnv(true)

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
