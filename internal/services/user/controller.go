package user

import (
	"github.com/gin-gonic/gin"
	"github.com/goldlilya1612/diploma-backend/internal/models"
	"gorm.io/gorm"
	"net/http"
)

const (
	errStatus     = "error"
	successStatus = "success"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(DB *gorm.DB) *UserController {
	return &UserController{
		DB: DB,
	}
}

func (uc *UserController) GetMe(ctx *gin.Context) {

	currentUser := ctx.MustGet("currentUser").(models.User)

	userResponse := &models.UserResponse{
		ID:    currentUser.ID,
		Name:  currentUser.Name,
		Email: currentUser.Email,
		Role:  currentUser.Role,

		CreatedAt: currentUser.CreatedAt,
		UpdatedAt: currentUser.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": successStatus, "data": gin.H{"user": userResponse}})
}
