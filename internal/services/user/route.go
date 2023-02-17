package user

import "github.com/gin-gonic/gin"

type UserRouteController struct {
	UserRouteController UserController
}

func NewUserRouteController(userController *UserController) UserRouteController {
	return UserRouteController{
		UserRouteController: *userController,
	}
}

func (arc *UserRouteController) AuthRoute(rg *gin.RouterGroup) {

	router := rg.Group("/users")

	router.GET("/me", arc.UserRouteController.GetMe)
}
