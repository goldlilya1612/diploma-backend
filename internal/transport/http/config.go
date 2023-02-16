package http

type Config struct {
	Port string
}

func NewConfig() *Config {
	return &Config{
		Port: "8080",
	}
}
