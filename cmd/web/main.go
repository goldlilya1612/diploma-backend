package main

import (
	"flag"
	"fmt"
	"github.com/goldlilya1612/diploma-backend/internal/app"
	"github.com/goldlilya1612/diploma-backend/internal/database"
	"github.com/goldlilya1612/diploma-backend/internal/services/auth"
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

	psql, err := app.StartPostgreSQL(dbConfig)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	serverConfig := http.NewConfig()
	gs := http.NewGinServer(serverConfig)

	authConfig := auth.NewConfig()
	authRouter := app.PrepareAuthRoute(authConfig, psql.DB)
	usersRouter := app.PrepareUsersRoute(psql.DB, authRouter.AuthController)
	coursesRouter := app.PrepareCoursesRoute(psql.DB, authRouter.AuthController)

	mainGroup := gs.Server.Group("/api")
	authRouter.AuthRoute(mainGroup)
	usersRouter.UsersRoute(mainGroup)
	coursesRouter.CourseRoute(mainGroup)

	app.StartGinServer(gs)
}
