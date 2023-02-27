package courses

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goldlilya1612/diploma-backend/internal/models"
	"github.com/goldlilya1612/diploma-backend/internal/services"
	"github.com/goldlilya1612/diploma-backend/internal/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"strings"
	"time"
)

type CoursesController struct {
	DB *gorm.DB
}

func NewCoursesController(DB *gorm.DB) *CoursesController {
	return &CoursesController{
		DB: DB,
	}
}

func (cc *CoursesController) CreateCourse(ctx *gin.Context) {
	var payload *models.CreateCourse

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": services.ErrStatus, "message": err})
		return
	}

	currentUser := ctx.MustGet("currentUser").(models.User)
	if payload.CreatorID != currentUser.ID {
		message := "Access denied"
		ctx.JSON(http.StatusForbidden, gin.H{"status": services.ErrStatus, "message": message})
		return
	}

	if payload.CreatorName != currentUser.Name {
		message := "The creator's name does not match the current user's name"
		ctx.JSON(http.StatusBadRequest, gin.H{"status": services.ErrStatus, "message": message})
		return
	}

	now := time.Time{}
	newCourse := models.Course{
		Name:        payload.Name,
		CreatorName: payload.CreatorName,
		CreatorID:   payload.CreatorID,
		Image:       payload.Image,
		Category:    payload.Category,
		Description: payload.Description,

		CreateAt: now,
		UpdateAt: now,
	}

	res := cc.DB.Create(&newCourse)
	if res.Error != nil {
		ctx.JSON(http.StatusConflict, gin.H{"status": services.ErrStatus, "message": res.Error.Error()})
		return
	}

	courseResponse := &models.CourseResponse{
		ID:          newCourse.ID,
		Name:        payload.Name,
		CreatorID:   payload.CreatorID,
		CreatorName: payload.CreatorName,
		Image:       payload.Image,
		Category:    payload.Category,
		Description: payload.Description,
		Route:       utils.Latinizer(payload.Name),

		CreateAt: now,
		UpdateAt: now,
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": services.SuccessStatus, "data": gin.H{"course": courseResponse}})
}

func (cc *CoursesController) GetCourse(ctx *gin.Context) {

	var coursesResponse []models.CourseResponse

	params := ctx.Request.URL.Query()
	ids := params["id"]

	for _, id := range ids {
		course := models.Course{}
		res := cc.DB.First(&course, "id = ?", fmt.Sprint(id))
		if res.Error != nil && strings.Contains(res.Error.Error(), "record not found") {
			message := fmt.Sprintf("Course with id=%s not found", id)
			ctx.JSON(http.StatusBadRequest, gin.H{"status": services.ErrStatus, "message": message})
			return
		} else if res.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": services.ErrStatus, "message": res.Error.Error()})
			return
		}

		courseResponse := models.CourseResponse{
			ID:          course.ID,
			Name:        course.Name,
			CreatorID:   course.CreatorID,
			CreatorName: course.Name,
			Image:       course.Image,
			Category:    course.Category,
			Description: course.Description,
			Route:       utils.Latinizer(course.Name),

			CreateAt: course.CreateAt,
			UpdateAt: course.UpdateAt,
		}
		coursesResponse = append(coursesResponse, courseResponse)
	}

	ctx.JSON(http.StatusOK, gin.H{"status": services.SuccessStatus, "data": gin.H{"courses": coursesResponse}})
}

func (cc *CoursesController) GetCourses(ctx *gin.Context) {

	var coursesResponse []models.CourseResponse

	var courses []models.Course
	res := cc.DB.Find(&courses)
	if res.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": services.ErrStatus, "message": res.Error.Error()})
		return
	}

	for _, c := range courses {
		courseResponse := models.CourseResponse{
			ID:          c.ID,
			Name:        c.Name,
			CreatorID:   c.CreatorID,
			CreatorName: c.CreatorName,
			Image:       c.Image,
			Category:    c.Category,
			Description: c.Description,
			Route:       utils.Latinizer(c.Name),

			CreateAt: c.CreateAt,
			UpdateAt: c.UpdateAt,
		}
		coursesResponse = append(coursesResponse, courseResponse)
	}

	ctx.JSON(http.StatusOK, gin.H{"status": services.SuccessStatus, "data": gin.H{"courses": coursesResponse}})
}

func (cc *CoursesController) DeleteCourse(ctx *gin.Context) {

	var coursesResponse []models.CourseResponse

	params := ctx.Request.URL.Query()
	ids := params["id"]

	for _, id := range ids {

		course := models.Course{}
		res := cc.DB.Clauses(clause.Returning{}).Where("id = ?", fmt.Sprint(id)).Delete(&course)
		if res.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": services.ErrStatus, "message": res.Error.Error()})
			return
		}

		courseResponse := models.CourseResponse{
			ID:          course.ID,
			Name:        course.Name,
			CreatorID:   course.CreatorID,
			CreatorName: course.CreatorName,
			Image:       course.Image,
			Category:    course.Category,
			Description: course.Description,
			Route:       utils.Latinizer(course.Name),

			CreateAt: course.CreateAt,
			UpdateAt: course.UpdateAt,
		}
		coursesResponse = append(coursesResponse, courseResponse)
	}

	ctx.JSON(http.StatusOK, gin.H{"status": services.SuccessStatus, "data": gin.H{"deletedCourses": coursesResponse}})
}
