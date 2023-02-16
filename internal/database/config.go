package database

type Config struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
	SslMode  string
}

func NewConfig() *Config {
	return &Config{
		Host:     "localhost",
		Port:     "5432",
		Database: "website_data",
		User:     "postgres",
		Password: "password",
		SslMode:  "disable",
	}
}
