package user

import (
	"github.com/gin-gonic/gin"
	"github.com/goldlilya1612/diploma-backend/internal/services/auth"
)

type UsersRouteController struct {
	UsersRouteController *UsersController
	authController       *auth.AuthController
}

func NewUsersRouteController(usersController *UsersController, authController *auth.AuthController) UsersRouteController {
	return UsersRouteController{
		UsersRouteController: usersController,
		authController:       authController,
	}
}

func (urc *UsersRouteController) UsersRoute(rg *gin.RouterGroup) {

	router := rg.Group("/user")
	router.GET("/me", urc.authController.DeserializeUser(), urc.UsersRouteController.GetMe)
}
