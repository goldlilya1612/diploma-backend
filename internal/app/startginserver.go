package app

import "github.com/goldlilya1612/diploma-backend/internal/transport/http"

func StartGinServer(config *http.Config) {

	gs := http.NewGinServer(config)
	gs.PrepareHealthchecker()
	gs.Start()
}
