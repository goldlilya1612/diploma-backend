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

func (c *Controller) articlesRoute(rg *gin.RouterGroup) {

	articlesRouter := rg.Group(articlesRout)

	articlesRouter.POST("/create", c.authController.DeserializeUser(), c.authController.CheckAccessRole(auth.LecturerRole), c.CreateArticle)
	articlesRouter.PATCH("/update", c.authController.DeserializeUser(), c.authController.CheckAccessRole(auth.LecturerRole), c.UpdateArticle)
}

func (c *Controller) CreateArticle(ctx *gin.Context) {
	var payload *models.CreateArticle

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
	res := c.DB.Table("courses").Joins("JOIN Chapters ON Chapters.course_id = Courses.id").
		Select("creator_id").Where("Chapters.id = ?", payload.ChapterID).Scan(&creatorID)
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
	newArticle := models.Article{
		Name:      payload.Name,
		ChapterID: payload.ChapterID,

		CreatedAt: now,
		UpdatedAt: now,
	}

	res = c.DB.Create(&newArticle)
	if res.Error != nil {
		ctx.JSON(http.StatusConflict, models.HTTPResponse{
			Status:     serv.ErrResponseStatus,
			StatusCode: http.StatusConflict,
			Message:    res.Error.Error(),
		})
		return
	}

	articleResponse := &models.ArticleResponse{
		ID:        newArticle.ID,
		Name:      newArticle.Name,
		ChapterID: newArticle.ChapterID,
		Route:     utils.Latinizer(payload.Name),

		CreatedAt: newArticle.CreatedAt,
		UpdatedAt: newArticle.UpdatedAt,
	}

	ctx.JSON(http.StatusCreated, models.HTTPResponse{
		Status:     serv.SuccessResponseStatus,
		StatusCode: http.StatusCreated,
		Data:       map[string]interface{}{"createdArticle": articleResponse},
	})
}

func (c *Controller) UpdateArticle(ctx *gin.Context) {
	var payload *models.UpdateArticle

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
		InnerJoins("JOIN Chapters ON Chapters.course_id = Courses.id").
		InnerJoins("JOIN Articles ON Articles.chapter_id = Chapters.id").
		Select("creator_id").
		Where("Articles.id = ?", payload.ID).
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

	updatedArticle := models.Article{
		ID: payload.ID,
	}
	res = c.DB.Model(&updatedArticle).Clauses(clause.Returning{}).Updates(models.Article{
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

	chapterResponse := &models.ArticleResponse{
		ID:        updatedArticle.ID,
		Name:      updatedArticle.Name,
		ChapterID: updatedArticle.ChapterID,
		Route:     utils.Latinizer(payload.Name),

		CreatedAt: updatedArticle.CreatedAt,
		UpdatedAt: updatedArticle.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, models.HTTPResponse{
		Status:     serv.SuccessResponseStatus,
		StatusCode: http.StatusOK,
		Data:       map[string]interface{}{"updatedChapter": chapterResponse},
	})
}

func (c *Controller) DeleteArticles(ctx *gin.Context) {

}
