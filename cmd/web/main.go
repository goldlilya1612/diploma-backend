package main

import (
	"flag"
	"fmt"
	"github.com/goldlilya1612/diploma-backend/internal/controllers/auth"
	"github.com/goldlilya1612/diploma-backend/internal/controllers/course"
	"github.com/goldlilya1612/diploma-backend/internal/controllers/user"
	"github.com/goldlilya1612/diploma-backend/internal/database"
	"github.com/goldlilya1612/diploma-backend/internal/transport/http"
	"os"
)

func main() {

	var psqlConfigPath string
	var psqlConfigName string
	flag.StringVar(&psqlConfigPath, "psql_conf_path", "configs/", "path to PostgreSQL config file")
	flag.StringVar(&psqlConfigName, "psql_conf_name", "default-psql-conf", "name PostgreSQL config file (without extension)")
	flag.Parse()

	dbConfig, err := database.NewConfigFromEnv(psqlConfigPath, psqlConfigName)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	//viper.OnConfigChange(func(e fsnotify.Event) {
	//	fmt.Println("Config file changed:", e.Name)
	//})
	//viper.WatchConfig()

	psql := database.NewPostgreSQL(dbConfig)
	err = psql.StartPostgreSQL()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	serverConfig := http.NewConfig()
	server := http.NewGinServer(serverConfig)

	authConfig := auth.NewConfig()
	authController := auth.NewController(authConfig, psql.DB)
	userController := user.NewController(psql.DB, authController)
	courseController := course.NewController(psql.DB, authController)

	server.StartAllRoutes(authController, userController, courseController)
	server.StartGinServer()
}
