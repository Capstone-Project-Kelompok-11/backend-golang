package models

import (
	"gorm.io/gorm"
)

type Modules struct {
	*gorm.Model
	ID          string `gorm:"type:VARCHAR(32);primary" json:"id"`
	CourseID    string `gorm:"type:VARCHAR(32);not null" json:"course_id"`
	Name        string `gorm:"type:VARCHAR(52);unique;not null" json:"name"`
	Description string `gorm:"type:VARCHAR(200)" json:"description"`
}

func (Modules) TableName() string {

	return "modules"
}
