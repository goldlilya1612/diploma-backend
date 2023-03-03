package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

const (
	ErrResponseStatus     = "error"
	SuccessResponseStatus = "success"
)

type GinServer struct {
	config *Config
	Server *gin.Engine
}

func NewGinServer(config *Config) *GinServer {

	serv := gin.Default()
	serv.Use(cors.New(config.Cors))

	return &GinServer{
		config: config,
		Server: serv,
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
