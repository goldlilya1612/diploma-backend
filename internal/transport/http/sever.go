package http

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type GinServer struct {
	config *Config
	server *gin.Engine
}

func NewGinServer(config *Config) *GinServer {

	return &GinServer{
		config: config,
		server: gin.Default(),
	}
}

func (gs *GinServer) PrepareHealthchecker() {
	router := gs.server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		message := "Connected"
		status := "success"
		ctx.JSON(http.StatusOK, gin.H{"status": status, "message": message})
	})
}

func (gs *GinServer) Start() {
	log.Fatal(gs.server.Run(":" + gs.config.Port))
}
