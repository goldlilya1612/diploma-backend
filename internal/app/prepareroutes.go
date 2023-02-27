package app

import (
	"github.com/goldlilya1612/diploma-backend/internal/services/auth"
	"github.com/goldlilya1612/diploma-backend/internal/services/courses"
	"github.com/goldlilya1612/diploma-backend/internal/services/user"
	"gorm.io/gorm"
)

func PrepareAuthRoute(config *auth.Config, db *gorm.DB) auth.AuthRouteController {

	authController := auth.NewAuthController(config, db)
	return auth.NewAuthRouteController(authController)
}

func PrepareUsersRoute(db *gorm.DB, authController *auth.AuthController) user.UsersRouteController {

	usersController := user.NewUsersController(db)
	return user.NewUsersRouteController(usersController, authController)
}

func PrepareCoursesRoute(db *gorm.DB, authController *auth.AuthController) courses.CoursesRouteController {

	coursesController := courses.NewCoursesController(db)
	return courses.NewCoursesRouteController(coursesController, authController)
}
