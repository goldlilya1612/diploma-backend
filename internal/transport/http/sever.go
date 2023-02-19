package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type GinServer struct {
	config *Config
	Server *gin.Engine
}

func NewGinServer(config *Config) *GinServer {

	serv := gin.Default()
	serv.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
