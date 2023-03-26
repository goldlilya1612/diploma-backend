package services

import "gorm.io/gorm"

type CoursesController struct {
	DB *gorm.DB
}

func NewCoursesController(DB *gorm.DB) *CoursesController {
	return &CoursesController{
		DB: DB,
	}
}
