package courses

import (
	"github.com/gin-gonic/gin"
	"github.com/goldlilya1612/diploma-backend/internal/services/auth"
	"net/http"
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

	router := rg.Group("/courses")

	router.POST("/create", crc.authController.DeserializeUser(), crc.authController.CheckAccessRole(auth.LecturerRole), crc.CoursesRouteController.CreateCourse)
	router.GET("/get-course", crc.authController.DeserializeUser(), crc.CoursesRouteController.GetCourse)
	router.GET("/get-courses", crc.authController.DeserializeUser(), crc.CoursesRouteController.GetCourses)
	router.DELETE("/delete", crc.authController.DeserializeUser(), crc.authController.CheckAccessRole(auth.LecturerRole), crc.CoursesRouteController.DeleteCourse)

	router.StaticFS("/images", http.Dir("resources/images"))
}
