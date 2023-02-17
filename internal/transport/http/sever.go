package http

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type GinServer struct {
	config *Config
	Server *gin.Engine
}

func NewGinServer(config *Config) *GinServer {

	return &GinServer{
		config: config,
		Server: gin.Default(),
	}
}

func (gs *GinServer) PrepareHealthchecker() {
	router := gs.Server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		message := "Connected"
		status := "success"
		ctx.JSON(http.StatusOK, gin.H{"status": status, "message": message})
	})
}

func (gs *GinServer) Start() {
	log.Fatal(gs.Server.Run(":" + gs.config.Port))
}
