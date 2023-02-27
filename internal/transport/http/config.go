package http

import (
	"github.com/gin-contrib/cors"
	"time"
)

type Config struct {
	Port string      `toml:"port"`
	Cors cors.Config `toml:"cors"`
}

func NewConfig() *Config {
	return &Config{
		Port: "8080",
		Cors: cors.Config{
			AllowOrigins:     []string{"http://localhost:3000"},
			AllowMethods:     []string{"*"},
			AllowHeaders:     []string{"*"},
			ExposeHeaders:    []string{"*"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		},
	}
}
