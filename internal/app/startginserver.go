package app

import "github.com/goldlilya1612/diploma-backend/internal/transport/http"

func StartGinServer(gs *http.GinServer) {
	gs.PrepareHealthchecker()
	gs.Start()
}
