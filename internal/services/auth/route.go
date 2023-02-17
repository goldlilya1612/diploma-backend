package auth

import "github.com/gin-gonic/gin"

type AuthRouteController struct {
	AuthRouteController AuthController
}

func NewAuthRouteController(authController *AuthController) AuthRouteController {
	return AuthRouteController{
		AuthRouteController: *authController,
	}
}

func (arc *AuthRouteController) AuthRoute(rg *gin.RouterGroup) {

	router := rg.Group("/auth")

	router.POST("/signup", arc.AuthRouteController.SignUpUser)
	router.POST("/signin", arc.AuthRouteController.SignInUser)
	router.POST("/refresh", arc.AuthRouteController.RefreshAccessToken)
	router.POST("/logout", arc.AuthRouteController.LogoutUser)
}
