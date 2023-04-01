package course

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goldlilya1612/diploma-backend/internal/controllers/auth"
	"github.com/goldlilya1612/diploma-backend/internal/models"
	serv "github.com/goldlilya1612/diploma-backend/internal/transport/http"
	"github.com/goldlilya1612/diploma-backend/internal/utils"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const (
	coursesRout  = "/course"
	chaptersRout = "/chapter"
	articlesRout = "/article"

	imagesDir = "./resources/images/"
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
	coursesRouter := rg.Group(coursesRout)

	coursesRouter.POST("/create", c.authController.DeserializeUser(), c.authController.CheckAccessRole(auth.LecturerRole), c.CreateCourse)
	coursesRouter.GET("/get-course", c.authController.DeserializeUser(), c.GetCourse)
	coursesRouter.GET("/get-courses", c.authController.DeserializeUser(), c.GetCourses)
	coursesRouter.PATCH("/update", c.authController.DeserializeUser(), c.authController.CheckAccessRole(auth.LecturerRole), c.UpdateCourse)
	coursesRouter.DELETE("/delete", c.authController.DeserializeUser(), c.authController.CheckAccessRole(auth.LecturerRole), c.DeleteCourse)
	coursesRouter.StaticFS("/images", http.Dir("resources/images"))

	c.chaptersRoute(coursesRouter)
}

func (c *Controller) CreateCourse(ctx *gin.Context) {
	var payload *models.CreateCourse

	err := ctx.ShouldBind(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Status:     serv.ErrResponseStatus,
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
		return
	}

	image := models.Image{}
	imageName, imagePath, imageURL, err := imageUpload(ctx)
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		ctx.JSON(http.StatusInternalServerError, models.HTTPResponse{
			Status:     serv.ErrResponseStatus,
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
		return
	} else if err == nil {
		image.Name = imageName
		image.Path = imagePath
		image.URL = imageURL
	}

	currentUser := ctx.MustGet("currentUser").(models.User)

	now := time.Time{}
	newCourse := models.Course{
		Name:        payload.Name,
		CreatorID:   currentUser.ID,
		Category:    payload.Category,
		Description: payload.Description,

		Image: image,

		CreatedAt: now,
		UpdatedAt: now,
	}

	res := c.DB.Create(&newCourse)
	if res.Error != nil {
		deleteImage(imagePath)
		ctx.JSON(http.StatusConflict, models.HTTPResponse{
			Status:     serv.ErrResponseStatus,
			StatusCode: http.StatusConflict,
			Message:    res.Error.Error(),
		})
		return
	}

	courseResponse := &models.CourseResponse{
		ID:          newCourse.ID,
		Name:        newCourse.Name,
		CreatorName: currentUser.Name,
		Image:       newCourse.Image,
		Category:    newCourse.Category,
		Description: newCourse.Description,
		Route:       utils.Latinizer(payload.Name),

		CreatedAt: newCourse.CreatedAt,
		UpdatedAt: newCourse.UpdatedAt,
	}

	ctx.JSON(http.StatusCreated, models.HTTPResponse{
		Status:     serv.SuccessResponseStatus,
		StatusCode: http.StatusCreated,
		Data:       map[string]interface{}{"createdCourse": courseResponse},
	})
}

func (c *Controller) GetCourse(ctx *gin.Context) {

	var coursesResponse []models.CourseResponse

	params := ctx.Request.URL.Query()
	ids := params["id"]

	for _, id := range ids {

		course := models.Course{}
		res := c.DB.Preload("Chapters.Articles").
			Joins("Lecturer").Joins("Image").
			First(&course, "Courses.id = ?", id)
		if res.Error != nil && strings.Contains(res.Error.Error(), "record not found") {
			message := fmt.Sprintf("Course with id=%s not found", id)
			ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
				Status:     serv.ErrResponseStatus,
				StatusCode: http.StatusBadRequest,
				Message:    message,
			})
			return
		} else if res.Error != nil {
			ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
				Status:     serv.ErrResponseStatus,
				StatusCode: http.StatusBadRequest,
				Message:    res.Error.Error(),
			})
			return
		}

		courseResponse := models.CourseResponse{
			ID:          course.ID,
			Name:        course.Name,
			CreatorName: course.Lecturer.Name,
			Image:       course.Image,
			Category:    course.Category,
			Description: course.Description,
			Route:       utils.Latinizer(course.Name),

			Chapters: course.Chapters,

			CreatedAt: course.CreatedAt,
			UpdatedAt: course.UpdatedAt,
		}
		coursesResponse = append(coursesResponse, courseResponse)
	}

	ctx.JSON(http.StatusOK, models.HTTPResponse{
		Status:     serv.SuccessResponseStatus,
		StatusCode: http.StatusOK,
		Data:       map[string]interface{}{"courses": coursesResponse},
	})
}

func (c *Controller) GetCourses(ctx *gin.Context) {

	var coursesResponse []models.CourseResponse

	var courses []models.Course
	res := c.DB.Joins("Lecturer").Joins("Image").Order("id").Find(&courses)
	if res.Error != nil {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Status:     serv.ErrResponseStatus,
			StatusCode: http.StatusBadRequest,
			Message:    res.Error.Error(),
		})
		return
	}

	for _, c := range courses {
		courseResponse := models.CourseResponse{
			ID:          c.ID,
			Name:        c.Name,
			CreatorName: c.Lecturer.Name,
			Image:       c.Image,
			Category:    c.Category,
			Description: c.Description,
			Route:       utils.Latinizer(c.Name),

			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		}
		coursesResponse = append(coursesResponse, courseResponse)
	}

	ctx.JSON(http.StatusOK, models.HTTPResponse{
		Status:     serv.SuccessResponseStatus,
		StatusCode: http.StatusOK,
		Data:       map[string]interface{}{"courses": coursesResponse},
	})
}

func (c *Controller) UpdateCourse(ctx *gin.Context) {
	var payload *models.UpdateCourse

	err := ctx.ShouldBind(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Status:     serv.ErrResponseStatus,
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
		return
	}

	course := models.Course{}
	res := c.DB.Joins("Image").Joins("Lecturer").
		First(&course, "Courses.id = ?", payload.ID)
	if res.Error != nil && strings.Contains(res.Error.Error(), "record not found") {
		message := fmt.Sprintf("Course with id=%d not found", payload.ID)
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Status:     serv.ErrResponseStatus,
			StatusCode: http.StatusBadRequest,
			Message:    message,
		})
		return
	} else if res.Error != nil {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Status:     serv.ErrResponseStatus,
			StatusCode: http.StatusBadRequest,
			Message:    res.Error.Error(),
		})
		return
	}
	creatorName := course.Lecturer.Name
	course.Lecturer = models.Lecturer{}

	currentUser := ctx.MustGet("currentUser").(models.User)
	if currentUser.ID != course.CreatorID {
		message := "Access denied"
		ctx.JSON(http.StatusForbidden, models.HTTPResponse{
			Status:     serv.ErrResponseStatus,
			StatusCode: http.StatusForbidden,
			Message:    message,
		})
		return
	}

	var oldImagePath string
	imageName, imagePath, imageURL, err := imageUpload(ctx)
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		ctx.JSON(http.StatusInternalServerError, models.HTTPResponse{
			Status:     serv.ErrResponseStatus,
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
		return
	} else if err == nil {
		oldImagePath = course.Image.Path
		err = c.DB.Model(&course).Association("Image").Replace(&course.Image, &models.Image{
			ID:   course.ImageID,
			Name: imageName,
			Path: imagePath,
			URL:  imageURL,
		})
		if err != nil {
			deleteImage(imagePath)
			ctx.JSON(http.StatusConflict, models.HTTPResponse{
				Status:     serv.ErrResponseStatus,
				StatusCode: http.StatusConflict,
				Message:    err.Error(),
			})
			return
		}
	}

	res = c.DB.Model(&course).Session(&gorm.Session{FullSaveAssociations: true}).Updates(models.Course{
		Name:        payload.Name,
		Category:    payload.Category,
		Description: payload.Description,
	})
	if res.Error != nil {
		deleteImage(imagePath)
		ctx.JSON(http.StatusConflict, models.HTTPResponse{
			Status:     serv.ErrResponseStatus,
			StatusCode: http.StatusConflict,
			Message:    res.Error.Error(),
		})
		return
	}
	deleteImage(oldImagePath)

	courseResponse := &models.CourseResponse{
		ID:          course.ID,
		Name:        course.Name,
		CreatorName: creatorName,
		Image:       course.Image,
		Category:    course.Category,
		Description: course.Description,
		Route:       utils.Latinizer(payload.Name),

		UpdatedAt: course.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, models.HTTPResponse{
		Status:     serv.SuccessResponseStatus,
		StatusCode: http.StatusOK,
		Data:       map[string]interface{}{"updatedCourse": courseResponse},
	})
}

func (c *Controller) DeleteCourse(ctx *gin.Context) {

	var coursesResponse []models.CourseResponse

	params := ctx.Request.URL.Query()
	ids := params["id"]

	for _, id := range ids {

		course := models.Course{}
		res := c.DB.Joins("Image").First(&course, "Courses.id = ?", id)
		if res.Error != nil && strings.Contains(res.Error.Error(), "record not found") {
			message := fmt.Sprintf("Course with id=%s not found", id)
			ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
				Status:     serv.ErrResponseStatus,
				StatusCode: http.StatusBadRequest,
				Message:    message,
			})
			return
		} else if res.Error != nil {
			ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
				Status:     serv.ErrResponseStatus,
				StatusCode: http.StatusBadRequest,
				Message:    res.Error.Error(),
			})
			return
		}

		currentUser := ctx.MustGet("currentUser").(models.User)
		if currentUser.ID != course.CreatorID {
			message := "Access denied"
			ctx.JSON(http.StatusForbidden, models.HTTPResponse{
				Status:     serv.ErrResponseStatus,
				StatusCode: http.StatusForbidden,
				Message:    message,
			})
			return
		}

		//err := cc.DB.Model(&course).Association("Image").Error
		//if err != nil {
		//	ctx.JSON(http.StatusConflict, models.HTTPResponse{
		//		Status:     serv.ErrResponseStatus,
		//		StatusCode: http.StatusConflict,
		//		Message:    err.Error(),
		//	})
		//	return
		//}

		res = c.DB.Joins("join images on images.id = course.image").Where("images.id = ?", fmt.Sprint(course.Image)).Delete(&course)
		if res.Error != nil {
			ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
				Status:     serv.ErrResponseStatus,
				StatusCode: http.StatusBadRequest,
				Message:    res.Error.Error(),
			})
			return
		}
		deleteImage(course.Image.Path)

		courseResponse := models.CourseResponse{
			ID:          course.ID,
			Name:        course.Name,
			Image:       course.Image,
			Category:    course.Category,
			Description: course.Description,
			Route:       utils.Latinizer(course.Name),

			CreatedAt: course.CreatedAt,
			UpdatedAt: course.UpdatedAt,
		}
		coursesResponse = append(coursesResponse, courseResponse)
	}

	ctx.JSON(http.StatusOK, models.HTTPResponse{
		Status:     serv.SuccessResponseStatus,
		StatusCode: http.StatusOK,
		Data:       map[string]interface{}{"deletedCourses": coursesResponse},
	})
}

func imageUpload(ctx *gin.Context) (imageName string, imagePath string, url string, err error) {

	file, header, err := ctx.Request.FormFile("image")
	if err != nil {
		return "", "", "", err
	}

	fileExt := filepath.Ext(header.Filename)
	originalImageName := strings.TrimSuffix(filepath.Base(header.Filename), filepath.Ext(header.Filename))
	now := time.Now()
	imageName = strings.ReplaceAll(strings.ToLower(originalImageName), " ", "-") + "-" + fmt.Sprintf("%v", now.Unix()) + fileExt
	url = path.Join(ctx.Request.Host, "/api", coursesRout, "/images", imageName)

	if _, err := os.Stat(imagesDir); os.IsNotExist(err) {
		err := os.MkdirAll(imagesDir, os.ModePerm)
		if err != nil {
			return "", "", "", err
		}
	}

	imagePath = imagesDir + imageName
	out, err := os.Create(imagePath)
	if err != nil {
		return "", "", "", err
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		return "", "", "", err
	}

	return imageName, imagePath, url, nil
}

func deleteImage(path string) {

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return
	}
	err := os.Remove(path)
	if err != nil {
		// TODO: интегрировать с логером Gin
		log.Printf("[ERROR] File deletion error: %s", err.Error())
	}
}
