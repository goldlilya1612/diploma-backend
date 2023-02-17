package app

import (
	"github.com/goldlilya1612/diploma-backend/internal/services/auth"
	"github.com/goldlilya1612/diploma-backend/internal/services/user"
	"gorm.io/gorm"
)

func PrepareAuthRoute(config *auth.Config, db *gorm.DB) auth.AuthRouteController {

	authController := auth.NewAuthController(config, db)
	return auth.NewAuthRouteController(authController)
}

func PrepareUserRoute(db *gorm.DB) user.UserRouteController {

	userController := user.NewUserController(db)
	return user.NewUserRouteController(userController)
}
