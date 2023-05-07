package database

import (
	"github.com/spf13/viper"
)

type Config struct {
	Host     string `mapstructure:"POSTGRES_HOST"`
	Port     string `mapstructure:"POSTGRES_PORT"`
	Database string `mapstructure:"POSTGRES_DB"`
	User     string `mapstructure:"POSTGRES_USER"`
	Password string `mapstructure:"POSTGRES_PASSWORD"`
	SslMode  string `mapstructure:"POSTGRES_SSL_MODE"`
}

var envs = map[string]string{
	"POSTGRES_HOST":     "localhost",
	"POSTGRES_PORT":     "5432",
	"POSTGRES_DB":       "postgres",
	"POSTGRES_USER":     "postgres",
	"POSTGRES_PASSWORD": "password",
	"POSTGRES_SSL_MODE": "disable",
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
