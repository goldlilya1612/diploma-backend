package main

import (
	"fmt"
	"github.com/goldlilya1612/diploma-backend/internal/app"
	"github.com/goldlilya1612/diploma-backend/internal/database"
	"github.com/goldlilya1612/diploma-backend/internal/services/auth"
	"github.com/goldlilya1612/diploma-backend/internal/transport/http"
	"os"
)

func main() {

	dbConfig := database.NewConfig()
	psql, err := app.StartPostgreSQL(dbConfig)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	serverConfig := http.NewConfig()
	gs := http.NewGinServer(serverConfig)

	authConfig := auth.NewConfig()
	authRouter := app.PrepareAuthRoute(authConfig, psql.DB)
	userRouter := app.PrepareUserRoute(psql.DB)

	authRouter.AuthRoute(gs.Server.Group("/api"))
	userRouter.AuthRoute(gs.Server.Group("/api"))

	app.StartGinServer(gs)
}
