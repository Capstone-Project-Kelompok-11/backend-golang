package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Courses struct {
	*gorm.Model
	ID          string          `gorm:"type:VARCHAR(32);primary" json:"id"`
	Name        string          `gorm:"type:VARCHAR(52);unique;not null" json:"name"`
	Description string          `gorm:"type:VARCHAR(200)" json:"description"`
	Price       decimal.Decimal `gorm:"default:0" json:"price"`
	Modules     []Modules       `gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"modules"`
}

func (Courses) TableName() string {

	return "courses"
}
