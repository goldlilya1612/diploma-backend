package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goldlilya1612/diploma-backend/internal/models"
	"github.com/goldlilya1612/diploma-backend/internal/services"
	"github.com/goldlilya1612/diploma-backend/internal/utils"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
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
		ctx.JSON(http.StatusBadRequest, gin.H{"status": services.ErrStatus, "message": err.Error()})
		return
	}

	if payload.Password != payload.PasswordConfirm {
		message := "Passwords do not match"
		ctx.JSON(http.StatusBadRequest, gin.H{"status": services.ErrStatus, "message": message})
		return
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": services.ErrStatus, "message": err.Error()})
		return
	}

	now := time.Time{}
	newUser := models.User{
		Name:      payload.Name,
		Email:     payload.Email,
		Password:  hashedPassword,
		Role:      payload.Role,
		Verified:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	res := ac.DB.Create(&newUser)
	if res.Error != nil && res.Error.(*pgconn.PgError).Code == pgerrcode.UniqueViolation {
		message := "Email already used"
		ctx.JSON(http.StatusConflict, gin.H{"status": services.ErrStatus, "message": message})
		return
	} else if res.Error != nil {
		ctx.JSON(http.StatusConflict, gin.H{"status": services.ErrStatus, "message": res.Error.Error()})
		return
	}

	switch newUser.Role {
	case services.StudentRole:
		var group string
		if len(payload.Groups) != 0 { // почему проверка лен а не нул  ???
			group = payload.Groups[0]
		} else if len(payload.Groups) > 0 || len(payload.Groups) == 0 {
			message := "No group specified or more than one group specified"
			ctx.JSON(http.StatusConflict, gin.H{"status": services.ErrStatus, "message": message})
			return
		}

		res := ac.DB.Create(&models.Student{
			Name:  newUser.Name,
			Group: group,
			User:  newUser,
		})
		if res.Error != nil {
			ctx.JSON(http.StatusConflict, gin.H{"status": services.ErrStatus, "message": err})
			return
		}
	case services.LecturerRole:
		if len(payload.Groups) == 0 {
			message := "Group not specified"
			ctx.JSON(http.StatusConflict, gin.H{"status": services.ErrStatus, "message": message})
			return
		}

		res := ac.DB.Create(&models.Lecturer{
			Name:   newUser.Name,
			Groups: payload.Groups,
			User:   newUser,
		})
		if res.Error != nil {
			ctx.JSON(http.StatusConflict, gin.H{"status": services.ErrStatus, "message": err})
			return
		}
	default:
		message := "Invalid role"
		ctx.JSON(http.StatusConflict, gin.H{"status": services.ErrStatus, "message": message})
		return
	}

	userResponse := &models.UserResponse{
		ID:        newUser.ID,
		Name:      newUser.Name,
		Groups:    payload.Groups,
		Email:     newUser.Email,
		Role:      newUser.Role,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
	}
	ctx.JSON(http.StatusCreated, gin.H{"status": services.SuccessStatus, "data": gin.H{"user": userResponse}})
}

func (ac *AuthController) SignInUser(ctx *gin.Context) {

	var payload *models.SignIn

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": services.ErrStatus, "message": err.Error()})
		return
	}

	var user *models.User
	res := ac.DB.First(&user, "email = ?", strings.ToLower(payload.Email))
	if res.Error != nil {
		message := "Invalid email or password"
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": services.ErrStatus, "message": message})
		return
	}

	err = utils.VerifyPassword(user.Password, payload.Password)
	if err != nil {
		message := "Invalid email or password"
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": services.ErrStatus, "message": message})
		return
	}

	accessToken, err := utils.CreateToken(ac.config.AccessTokenExpiresIn, user.ID, ac.config.AccessTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": services.ErrStatus, "message": err.Error()})
		return
	}

	refreshToken, err := utils.CreateToken(ac.config.RefreshTokenExpiresIn, user.ID, ac.config.RefreshTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": services.ErrStatus, "message": err.Error()})
		return
	}

	ctx.SetCookie("token", accessToken, ac.config.AccessTokenMaxAge*60, "/", ac.config.Domain, false, true)
	ctx.SetCookie("refresh_token", refreshToken, ac.config.RefreshTokenMaxAge*60, "/", ac.config.Domain, false, true)
	ctx.SetCookie("logged_in", "true", ac.config.AccessTokenMaxAge*60, "/", ac.config.Domain, false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": services.SuccessStatus, "token": accessToken})
}

func (ac *AuthController) RefreshAccessToken(ctx *gin.Context) {

	cookie, err := ctx.Cookie("refresh_token")
	if err != nil {
		message := "could not refresh access token"
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": services.ErrStatus, "message": message})
		return
	}

	sub, err := utils.ValidateToken(cookie, ac.config.RefreshTokenPublicKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": services.ErrStatus, "message": err.Error()})
		return
	}

	var user *models.User
	res := ac.DB.First(&user, "id = ?", fmt.Sprint(sub))
	if res.Error != nil {
		message := "the user belonging to this token no logger exists"
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": services.ErrStatus, "message": message})
		return
	}

	accessToken, err := utils.CreateToken(ac.config.AccessTokenExpiresIn, user.ID, ac.config.AccessTokenPrivateKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": services.ErrStatus, "message": err.Error()})
		return
	}

	ctx.SetCookie("token", accessToken, ac.config.AccessTokenMaxAge*60, "/", ac.config.Domain, false, true)
	ctx.SetCookie("logged_in", "true", ac.config.AccessTokenMaxAge*60, "/", ac.config.Domain, false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "token": accessToken})

}

func (ac *AuthController) LogoutUser(ctx *gin.Context) {
	ctx.SetCookie("token", "", -1, "/", ac.config.Domain, false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", ac.config.Domain, false, true)
	ctx.SetCookie("logged_in", "", -1, "/", ac.config.Domain, false, false)

	ctx.JSON(http.StatusNoContent, gin.H{"status": services.SuccessStatus})
}

func (ac *AuthController) DeserializeUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var accessToken string
		cookie, err := ctx.Cookie("token")

		authorizationHeader := ctx.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			accessToken = fields[1]
		} else if err == nil {
			accessToken = cookie
		}

		if accessToken == "" {
			message := "You are not logged in"
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": services.ErrStatus, "message": message})
			return
		}

		sub, err := utils.ValidateToken(accessToken, ac.config.AccessTokenPublicKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": services.ErrStatus, "message": err.Error()})
			return
		}

		var user models.User
		res := ac.DB.First(&user, "id = ?", fmt.Sprint(sub))
		if res.Error != nil {
			message := "The user belonging to this token no logger exists"
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": services.ErrStatus, "message": message})
			return
		}

		ctx.Set("currentUser", user)
		ctx.Next()
	}
}

func (ac *AuthController) CheckAccessRole(accessRole string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		currentUser := ctx.MustGet("currentUser").(models.User)

		if currentUser.Role != accessRole {
			message := "Invalid role for this method"
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": services.ErrStatus, "message": message})
			return
		}
	}
}
