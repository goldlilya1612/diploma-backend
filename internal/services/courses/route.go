package courses

import (
	"github.com/gin-gonic/gin"
	"github.com/goldlilya1612/diploma-backend/internal/services/auth"
	"net/http"
)

const (
	coursesRout  = "/course"
	chaptersRout = "/chapter"
	articlesRout = "/article"
)

type CourseRouteController struct {
	CoursesRouteController *CourseController
	authController         *auth.AuthController
}

func NewCourseRouteController(coursesController *CourseController, authController *auth.AuthController) CourseRouteController {
	return CourseRouteController{
		CoursesRouteController: coursesController,
		authController:         authController,
	}
}

func (crc *CourseRouteController) CourseRoute(rg *gin.RouterGroup) {

	coursesRouter := rg.Group(coursesRout)

	coursesRouter.POST("/create", crc.authController.DeserializeUser(), crc.authController.CheckAccessRole(auth.LecturerRole), crc.CoursesRouteController.CreateCourse)
	coursesRouter.GET("/get-course", crc.authController.DeserializeUser(), crc.CoursesRouteController.GetCourse)
	coursesRouter.GET("/get-courses", crc.authController.DeserializeUser(), crc.CoursesRouteController.GetCourses)
	coursesRouter.PATCH("/update", crc.authController.DeserializeUser(), crc.authController.CheckAccessRole(auth.LecturerRole), crc.CoursesRouteController.UpdateCourse)
	coursesRouter.DELETE("/delete", crc.authController.DeserializeUser(), crc.authController.CheckAccessRole(auth.LecturerRole), crc.CoursesRouteController.DeleteCourse)
	coursesRouter.StaticFS("/images", http.Dir("resources/images"))

	crc.chaptersRoute(coursesRouter)
}

func (crc *CourseRouteController) chaptersRoute(rg *gin.RouterGroup) {

	chaptersRouter := rg.Group(chaptersRout)

	chaptersRouter.POST("/create", crc.authController.DeserializeUser(), crc.authController.CheckAccessRole(auth.LecturerRole), crc.CoursesRouteController.CreateChapter)
	chaptersRouter.PATCH("/update", crc.authController.DeserializeUser(), crc.authController.CheckAccessRole(auth.LecturerRole), crc.CoursesRouteController.UpdateChapter)

	crc.articlesRoute(chaptersRouter)
}

func (crc *CourseRouteController) articlesRoute(rg *gin.RouterGroup) {

	articlesRouter := rg.Group(articlesRout)

	articlesRouter.POST("/create", crc.authController.DeserializeUser(), crc.authController.CheckAccessRole(auth.LecturerRole), crc.CoursesRouteController.CreateArticle)
	articlesRouter.PATCH("/update", crc.authController.DeserializeUser(), crc.authController.CheckAccessRole(auth.LecturerRole), crc.CoursesRouteController.UpdateArticle)
}
