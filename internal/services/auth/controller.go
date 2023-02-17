package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goldlilya1612/diploma-backend/internal/models"
	"github.com/goldlilya1612/diploma-backend/internal/utils"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

const (
	errStatus     = "error"
	successStatus = "success"

	userRole = "user"
)

type AuthController struct {
	config *Config
	DB     *gorm.DB
}

func NewAuthController(config *Config, DB *gorm.DB) *AuthController {
	return &AuthController{
		config: config,
		DB:     DB,
	}
}

func (ac *AuthController) SignUpUser(ctx *gin.Context) {

	var payload *models.SignUp

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": errStatus, "message": err.Error()})
		return
	}

	if payload.Password != payload.PasswordConfirm {
		message := "Passwords do not match"
		ctx.JSON(http.StatusBadRequest, gin.H{"status": errStatus, "message": message})
		return
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": errStatus, "message": err.Error()})
		return
	}

	now := time.Time{}
	newUser := &models.User{
		Name:      payload.Name,
		Email:     payload.Email,
		Password:  hashedPassword,
		Role:      userRole,
		Verified:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	res := ac.DB.Create(newUser)

	if res.Error != nil && strings.Contains(res.Error.Error(), "duplicate key value violates unique") {
		message := "Email already used"
		ctx.JSON(http.StatusConflict, gin.H{"status": errStatus, "message": message})
		return
	}

	userResponse := &models.UserResponse{
		ID:        newUser.ID,
		Name:      newUser.Name,
		Email:     newUser.Email,
		Role:      newUser.Role,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
	}
	ctx.JSON(http.StatusCreated, gin.H{"status": successStatus, "data": gin.H{"user": userResponse}})
}

func (ac *AuthController) SignInUser(ctx *gin.Context) {

	var payload *models.SignIn

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": errStatus, "message": err.Error()})
		return
	}

	var user *models.User
	res := ac.DB.First(&user, "email = ?", strings.ToLower(payload.Email))
	if res.Error != nil {
		message := "Invalid email or password"
		ctx.JSON(http.StatusBadRequest, gin.H{"status": errStatus, "message": message})
		return
	}

	err = utils.VerifyPassword(user.Password, payload.Password)
	if err != nil {
		message := "Invalid email or password"
		ctx.JSON(http.StatusBadRequest, gin.H{"status": errStatus, "message": message})
		return
	}

	accessToken, err := utils.CreateToken(ac.config.AccessTokenExpiresIn, user.ID, ac.config.AccessTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": errStatus, "message": err.Error()})
		return
	}

	refreshToken, err := utils.CreateToken(ac.config.RefreshTokenExpiresIn, user.ID, ac.config.RefreshTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": errStatus, "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", accessToken, ac.config.AccessTokenMaxAge*60, "/", ac.config.Domain, false, true)
	ctx.SetCookie("refresh_token", refreshToken, ac.config.RefreshTokenMaxAge*60, "/", ac.config.Domain, false, true)
	ctx.SetCookie("logged_in", "true", ac.config.AccessTokenMaxAge*60, "/", ac.config.Domain, false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": successStatus, "access_token": accessToken})
}

func (ac *AuthController) RefreshAccessToken(ctx *gin.Context) {

	cookie, err := ctx.Cookie("refresh_token")
	if err != nil {
		message := "could not refresh access token"
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": errStatus, "message": message})
		return
	}

	sub, err := utils.ValidateToken(cookie, ac.config.RefreshTokenPublicKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": errStatus, "message": err.Error()})
		return
	}

	var user *models.User
	res := ac.DB.First(&user, "id = ?", fmt.Sprint(sub))
	if res.Error != nil {
		message := "the user belonging to this token no logger exists"
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": errStatus, "message": message})
		return
	}

	accessToken, err := utils.CreateToken(ac.config.AccessTokenExpiresIn, user.ID, ac.config.AccessTokenPrivateKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": errStatus, "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", accessToken, ac.config.AccessTokenMaxAge*60, "/", ac.config.Domain, false, true)
	ctx.SetCookie("logged_in", "true", ac.config.AccessTokenMaxAge*60, "/", ac.config.Domain, false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})

}

func (ac *AuthController) LogoutUser(ctx *gin.Context) {
	ctx.SetCookie("access_token", "", -1, "/", ac.config.Domain, false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", ac.config.Domain, false, true)
	ctx.SetCookie("logged_in", "", -1, "/", ac.config.Domain, false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (ac *AuthController) DeserializeUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var accessToken string
		cookie, err := ctx.Cookie("access_token")

		authorizationHeader := ctx.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			accessToken = fields[1]
		} else if err == nil {
			accessToken = cookie
		}

		if accessToken == "" {
			message := "You are not logged in"
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": errStatus, "message": message})
			return
		}

		sub, err := utils.ValidateToken(accessToken, ac.config.AccessTokenPublicKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": errStatus, "message": err.Error()})
			return
		}

		var user models.User
		res := ac.DB.First(&user, "id = ?", fmt.Sprint(sub))
		if res.Error != nil {
			message := "the user belonging to this token no logger exists"
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": errStatus, "message": message})
			return
		}

		ctx.Set("currentUser", user)
		ctx.Next()
	}
}
