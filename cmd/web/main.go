package main

import (
	"fmt"
	"github.com/goldlilya1612/diploma-backend/internal/app"
	"github.com/goldlilya1612/diploma-backend/internal/database"
	"github.com/goldlilya1612/diploma-backend/internal/transport/http"
	"os"
)

func main() {

	dbConfig := database.NewConfig()
	err := app.StartPostgreSQL(dbConfig)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	serverConfig := http.NewConfig()
	app.StartGinServer(serverConfig)
}
