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

type CoursesRouteController struct {
	CoursesRouteController *CoursesController
	authController         *auth.AuthController
}

func NewCoursesRouteController(coursesController *CoursesController, authController *auth.AuthController) CoursesRouteController {
	return CoursesRouteController{
		CoursesRouteController: coursesController,
		authController:         authController,
	}
}

func (crc *CoursesRouteController) CoursesRoute(rg *gin.RouterGroup) {

	coursesRouter := rg.Group(coursesRout)

	coursesRouter.POST("/create", crc.authController.DeserializeUser(), crc.authController.CheckAccessRole(auth.LecturerRole), crc.CoursesRouteController.CreateCourse)
	coursesRouter.GET("/get-course", crc.authController.DeserializeUser(), crc.CoursesRouteController.GetCourse)
	coursesRouter.GET("/get-courses", crc.authController.DeserializeUser(), crc.CoursesRouteController.GetCourses)
	coursesRouter.PATCH("/update", crc.authController.DeserializeUser(), crc.authController.CheckAccessRole(auth.LecturerRole), crc.CoursesRouteController.UpdateCourse)
	coursesRouter.DELETE("/delete", crc.authController.DeserializeUser(), crc.authController.CheckAccessRole(auth.LecturerRole), crc.CoursesRouteController.DeleteCourse)
	coursesRouter.StaticFS("/images", http.Dir("resources/images"))

	crc.chaptersRoute(coursesRouter)
	crc.articlesRoute(coursesRouter)
}

func (crc *CoursesRouteController) chaptersRoute(rg *gin.RouterGroup) {

	chaptersRouter := rg.Group(chaptersRout)

	chaptersRouter.POST("/create", crc.authController.DeserializeUser(), crc.authController.CheckAccessRole(auth.LecturerRole), crc.CoursesRouteController.CreateChapter)
	chaptersRouter.PATCH("/update", crc.authController.DeserializeUser(), crc.authController.CheckAccessRole(auth.LecturerRole), crc.CoursesRouteController.UpdateChapter)
	_ = chaptersRouter
}

func (crc *CoursesRouteController) articlesRoute(rg *gin.RouterGroup) {

	articlesRouter := rg.Group(articlesRout)

	//chaptersRouter.POST("/create", crc.authController.DeserializeUser(), crc.authController.CheckAccessRole(auth.LecturerRole), crc.CoursesRouteController.CreateCourse)

	_ = articlesRouter
}
