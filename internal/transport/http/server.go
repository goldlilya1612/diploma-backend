package http

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"net/http"
	"time"
)

const (
	ErrResponseStatus     = "error"
	SuccessResponseStatus = "success"
)

type Router interface {
	Route(rg *gin.RouterGroup)
}

type GinServer struct {
	config       *Config
	server       *http.Server
	initialRoute *gin.RouterGroup
}

func NewGinServer(config *Config) *GinServer {

	router := gin.Default()

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
	router.Use(location.New(locationConf), cors.New(corsConf))

	serv := &http.Server{
		Addr:    net.JoinHostPort(config.Host, config.Port),
		Handler: router,
	}

	return &GinServer{
		config:       config,
		server:       serv,
		initialRoute: router.RouterGroup.Group("/api"),
	}
}

func (gs *GinServer) gracefulPowerOff() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := gs.server.Shutdown(ctx); err != nil {
		log.Printf("rror shutting down server %s", err)
	} else {
		log.Println("Server stopping...")
	}

	log.Println("Timeout of 5 seconds")
	select {
	case <-ctx.Done():
	}
	log.Println("Server successfully stopped")
}

func (gs *GinServer) prepareHealthcheck() {
	gs.initialRoute.GET("/healthcheck", func(ctx *gin.Context) {
		message := "Connected"
		ctx.JSON(http.StatusOK, gin.H{"status": SuccessResponseStatus, "message": message})
	})
}

func (gs *GinServer) serve() {
	if err := gs.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}

func (gs *GinServer) StartGinServer(ctxServ context.Context, cancelServ context.CancelFunc) {
	gs.prepareHealthcheck()
	go gs.serve()

	select {
	case <-ctxServ.Done():
		cancelServ()
		gs.gracefulPowerOff()
	}
}

func (gs *GinServer) StartAllRoutes(routers ...Router) {
	for _, r := range routers {
		r.Route(gs.initialRoute)
	}
}
