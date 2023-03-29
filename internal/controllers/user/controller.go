package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goldlilya1612/diploma-backend/internal/controllers/auth"
	"github.com/goldlilya1612/diploma-backend/internal/models"
	serv "github.com/goldlilya1612/diploma-backend/internal/transport/http"
	"gorm.io/gorm"
	"net/http"
)

type Controller struct {
	DB             *gorm.DB
	authController *auth.Controller
}

func NewController(DB *gorm.DB, authController *auth.Controller) *Controller {
	return &Controller{
		DB,
		authController,
	}
}

func (c *Controller) Route(rg *gin.RouterGroup) {

	router := rg.Group("/user")
	router.GET("/me", c.authController.DeserializeUser(), c.GetMe)
}

func (c *Controller) GetMe(ctx *gin.Context) {

	currentUser := ctx.MustGet("currentUser").(models.User)

	var groups []string
	switch currentUser.Role {
	case auth.StudentRole:
		student := &models.Student{}
		c.DB.First(&student, "user_id = ?", fmt.Sprint(currentUser.ID))
		groups = []string{student.Group}
	case auth.LecturerRole:
		lecturer := &models.Lecturer{}
		c.DB.First(&lecturer, "user_id = ?", fmt.Sprint(currentUser.ID))
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

	ctx.JSON(http.StatusOK, models.HTTPResponse{
		Status:     serv.SuccessResponseStatus,
		StatusCode: http.StatusOK,
		Data:       map[string]interface{}{"user": userResponse},
	})
}
