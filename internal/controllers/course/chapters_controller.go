package course

import (
	"github.com/gin-gonic/gin"
	"github.com/goldlilya1612/diploma-backend/internal/controllers/auth"
	"github.com/goldlilya1612/diploma-backend/internal/models"
	serv "github.com/goldlilya1612/diploma-backend/internal/transport/http"
	"github.com/goldlilya1612/diploma-backend/internal/utils"
	"gorm.io/gorm/clause"
	"net/http"
	"time"
)

func (c *Controller) chaptersRoute(rg *gin.RouterGroup) {

	chaptersRouter := rg.Group(chaptersRout)

	chaptersRouter.POST("/create", c.authController.DeserializeUser(), c.authController.CheckAccessRole(auth.LecturerRole), c.CreateChapter)
	chaptersRouter.PATCH("/update", c.authController.DeserializeUser(), c.authController.CheckAccessRole(auth.LecturerRole), c.UpdateChapter)
	chaptersRouter.DELETE("/delete", c.authController.DeserializeUser(), c.authController.CheckAccessRole(auth.LecturerRole), c.DeleteChapters)

	c.articlesRoute(chaptersRouter)
}

func (c *Controller) CreateChapter(ctx *gin.Context) {
	var payload *models.CreateChapter

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Status:     serv.ErrResponseStatus,
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
		return
	}

	var creatorID string
	res := c.DB.Table("courses").Select("creator_id").Where("id = ?", payload.CourseID).
		Scan(&creatorID)
	if res.Error != nil {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Status:     serv.ErrResponseStatus,
			StatusCode: http.StatusBadRequest,
			Message:    res.Error.Error(),
		})
		return
	}

	currentUser := ctx.MustGet("currentUser").(models.User)
	if currentUser.ID.String() != creatorID {
		message := "Access denied"
		ctx.JSON(http.StatusForbidden, models.HTTPResponse{
			Status:     serv.ErrResponseStatus,
			StatusCode: http.StatusForbidden,
			Message:    message,
		})
		return
	}

	now := time.Time{}
	newChapter := models.Chapter{
		Name:     payload.Name,
		CourseID: payload.CourseID,

		CreatedAt: now,
		UpdatedAt: now,
	}

	res = c.DB.Create(&newChapter)
	if res.Error != nil {
		ctx.JSON(http.StatusConflict, models.HTTPResponse{
			Status:     serv.ErrResponseStatus,
			StatusCode: http.StatusConflict,
			Message:    res.Error.Error(),
		})
		return
	}

	chapterResponse := models.ChapterResponse{
		ID:       newChapter.ID,
		Name:     newChapter.Name,
		CourseID: newChapter.CourseID,
		Route:    utils.Latinizer(newChapter.Name),

		CreatedAt: newChapter.CreatedAt,
		UpdatedAt: newChapter.UpdatedAt,
	}

	ctx.JSON(http.StatusCreated, models.HTTPResponse{
		Status:     serv.SuccessResponseStatus,
		StatusCode: http.StatusCreated,
		Data:       map[string]interface{}{"createdChapter": chapterResponse},
	})
}

func (c *Controller) UpdateChapter(ctx *gin.Context) {
	var payload *models.UpdateChapter

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Status:     serv.ErrResponseStatus,
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
		return
	}

	var creatorID string
	res := c.DB.Table("courses").
		Joins("JOIN Chapters ON Chapters.course_id = Courses.id").
		Select("creator_id").Where("Chapters.id = ?", payload.ID).
		Scan(&creatorID)
	if res.Error != nil {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Status:     serv.ErrResponseStatus,
			StatusCode: http.StatusBadRequest,
			Message:    res.Error.Error(),
		})
		return
	}

	currentUser := ctx.MustGet("currentUser").(models.User)
	if currentUser.ID.String() != creatorID {
		message := "Access denied"
		ctx.JSON(http.StatusForbidden, models.HTTPResponse{
			Status:     serv.ErrResponseStatus,
			StatusCode: http.StatusForbidden,
			Message:    message,
		})
		return
	}

	updatedChapter := models.Chapter{
		ID: payload.ID,
	}
	res = c.DB.Model(&updatedChapter).Clauses(clause.Returning{}).Updates(models.Chapter{
		Name: payload.Name,
	})
	if res.Error != nil {
		ctx.JSON(http.StatusConflict, models.HTTPResponse{
			Status:     serv.ErrResponseStatus,
			StatusCode: http.StatusConflict,
			Message:    res.Error.Error(),
		})
		return
	}

	chapterResponse := models.ChapterResponse{
		ID:       updatedChapter.ID,
		Name:     updatedChapter.Name,
		CourseID: updatedChapter.CourseID,
		Route:    utils.Latinizer(updatedChapter.Name),

		CreatedAt: updatedChapter.CreatedAt,
		UpdatedAt: updatedChapter.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, models.HTTPResponse{
		Status:     serv.SuccessResponseStatus,
		StatusCode: http.StatusOK,
		Data:       map[string]interface{}{"updatedChapter": chapterResponse},
	})
}

func (c *Controller) DeleteChapters(ctx *gin.Context) {
	var chaptersResponse []models.ChapterResponse

	params := ctx.Request.URL.Query()
	ids := params["id"]

	for _, id := range ids {

		chapter := models.Chapter{}
		res := c.DB.Preload("Articles").
			Joins("Course").First(&chapter, "Chapters.id = ?", id)
		if res.Error != nil {
			ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
				Status:     serv.ErrResponseStatus,
				StatusCode: http.StatusBadRequest,
				Message:    res.Error.Error(),
			})
			return
		}

		currentUser := ctx.MustGet("currentUser").(models.User)
		if currentUser.ID != chapter.Course.CreatorID {
			message := "Access denied"
			ctx.JSON(http.StatusForbidden, models.HTTPResponse{
				Status:     serv.ErrResponseStatus,
				StatusCode: http.StatusForbidden,
				Message:    message,
			})
			return
		}

		res = c.DB.Clauses(clause.Returning{}).Delete(&chapter)
		if res.Error != nil {
			ctx.JSON(http.StatusConflict, models.HTTPResponse{
				Status:     serv.ErrResponseStatus,
				StatusCode: http.StatusConflict,
				Message:    res.Error.Error(),
			})
			return
		}

		chapterResponse := models.ChapterResponse{
			ID:       chapter.ID,
			Name:     chapter.Name,
			CourseID: chapter.CourseID,
			Route:    utils.Latinizer(chapter.Name),
			Articles: chapter.Articles,
		}
		chaptersResponse = append(chaptersResponse, chapterResponse)
	}

	ctx.JSON(http.StatusOK, models.HTTPResponse{
		Status:     serv.SuccessResponseStatus,
		StatusCode: http.StatusOK,
		Data:       map[string]interface{}{"deletedChapters": chaptersResponse},
	})
}
