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

var envs = map[string]string{
	"GIN_HOST":   "localhost",
	"GIN_PORT":   "8080",
	"GIN_SCHEME": "http",

	"GIN_ALLOW_ORIGINS":     "http://localhost:3000",
	"GIN_ALLOW_METHODS":     "*",
	"GIN_ALLOW_HEADERS":     "*",
	"GIN_EXPOSE_HEADERS":    "*",
	"GIN_ALLOW_CREDENTIALS": "true",
}

func bindEnvs() {
	for k, v := range envs {
		viper.BindEnv(k)
		viper.SetDefault(k, v)
	}
}

func NewConfigFromEnv(path, name string) (*Config, error) {

	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType("env")

	viper.AllowEmptyEnv(true)

	err := viper.ReadInConfig()
	if err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			bindEnvs()
		default:
			return nil, err
		}
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
