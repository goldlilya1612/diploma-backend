package http

import (
	"github.com/gin-contrib/cors"
	"time"
)

type Config struct {
	Port string
	Cors cors.Config
}

func NewConfig() *Config {
	return &Config{
		Port: "8080",
		Cors: cors.Config{
			AllowOrigins:     []string{"http://localhost:3000"},
			AllowMethods:     []string{"GET", "POST"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		},
	}
}
