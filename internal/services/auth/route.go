package auth

import "github.com/gin-gonic/gin"

type AuthRouteController struct {
	AuthController *AuthController
}

func NewAuthRouteController(authController *AuthController) AuthRouteController {
	return AuthRouteController{
		AuthController: authController,
	}
}

func (arc *AuthRouteController) AuthRoute(rg *gin.RouterGroup) {

	router := rg.Group("/auth")

	router.POST("/signup", arc.AuthController.SignUpUser)
	router.POST("/signin", arc.AuthController.SignInUser)
	router.POST("/refresh", arc.AuthController.RefreshAccessToken)
	router.POST("/logout", arc.AuthController.DeserializeUser(), arc.AuthController.LogoutUser)
}
