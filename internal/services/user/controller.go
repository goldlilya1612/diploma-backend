package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goldlilya1612/diploma-backend/internal/models"
	"github.com/goldlilya1612/diploma-backend/internal/services"
	"gorm.io/gorm"
	"net/http"
)

type UsersController struct {
	DB *gorm.DB
}

func NewUsersController(DB *gorm.DB) *UsersController {
	return &UsersController{
		DB: DB,
	}
}

func (uc *UsersController) GetMe(ctx *gin.Context) {

	currentUser := ctx.MustGet("currentUser").(models.User)

	var groups []string
	switch currentUser.Role {
	case services.StudentRole:
		student := &models.Student{}
		uc.DB.First(&student, "user_id = ?", fmt.Sprint(currentUser.ID))
		groups = []string{student.Group}
	case services.LecturerRole:
		lecturer := &models.Lecturer{}
		uc.DB.First(&lecturer, "user_id = ?", fmt.Sprint(currentUser.ID))
		groups = lecturer.Groups
	}

	userResponse := &models.UserResponse{
		ID:     currentUser.ID,
		Name:   currentUser.Name,
		Email:  currentUser.Email,
		Role:   currentUser.Role,
		Groups: groups,

		CreatedAt: currentUser.CreatedAt,
		UpdatedAt: currentUser.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": services.SuccessStatus, "data": gin.H{"user": userResponse}})
}
