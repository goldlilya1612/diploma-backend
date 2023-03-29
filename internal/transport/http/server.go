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

type Router interface {
	Route(rg *gin.RouterGroup)
}

type GinServer struct {
	Config *Config
	Server *gin.Engine
	Router *gin.RouterGroup
}

func NewGinServer(config *Config) *GinServer {

	serv := gin.Default()
	serv.Use(cors.New(config.Cors))
	defaultRoute := serv.Group("/api")

	return &GinServer{
		Config: config,
		Server: serv,
		Router: defaultRoute,
	}
}

func (gs *GinServer) prepareHealthchecker() {
	gs.Router.GET("/healthchecker", func(ctx *gin.Context) {
		message := "Connected"
		status := "success"
		ctx.JSON(http.StatusOK, gin.H{"status": status, "message": message})
	})
}

func (gs *GinServer) start() {
	log.Fatal(gs.Server.Run(":" + gs.Config.Port))
}

func (gs *GinServer) StartGinServer() {
	gs.prepareHealthchecker()
	gs.start()
}

func (gs *GinServer) StartAllRoutes(routers ...Router) {
	for _, r := range routers {
		r.Route(gs.Router)
	}
}
