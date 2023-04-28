package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"log"
	"net"
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
	Config       *Config
	Server       *gin.Engine
	InitialRoute *gin.RouterGroup
}

func NewGinServer(config *Config) *GinServer {

	serv := gin.Default()

	locationConf := location.Config{
		Host:   net.JoinHostPort(config.Host, config.Port),
		Scheme: config.Scheme,
	}
	corsConf := cors.Config{
		AllowOrigins:     config.AllowOrigins,
		AllowMethods:     config.AllowMethods,
		AllowHeaders:     config.AllowHeaders,
		ExposeHeaders:    config.ExposeHeaders,
		AllowCredentials: config.AllowCredentials,
	}
	serv.Use(location.New(locationConf), cors.New(corsConf))

	return &GinServer{
		Config:       config,
		Server:       serv,
		InitialRoute: serv.RouterGroup.Group("/api"),
	}
}

func (gs *GinServer) prepareHealthchecker() {
	gs.InitialRoute.GET("/healthcheck", func(ctx *gin.Context) {
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
		r.Route(gs.InitialRoute)
	}
}
