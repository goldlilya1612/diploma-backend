package database

import "github.com/spf13/viper"

type Config struct {
	Host     string `mapstructure:"POSTGRES_HOST"`
	Port     string `mapstructure:"POSTGRES_PORT"`
	Database string `mapstructure:"POSTGRES_DB"`
	User     string `mapstructure:"POSTGRES_USER"`
	Password string `mapstructure:"POSTGRES_PASSWORD"`
	SslMode  string `mapstructure:"POSTGRES_SSL_MODE"`
}

func NewDefaultConfig() *Config {
	return &Config{
		Host:     "localhost",
		Port:     "5432",
		Database: "postgres",
		User:     "postgres",
		Password: "password",
		SslMode:  "disable",
	}
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
