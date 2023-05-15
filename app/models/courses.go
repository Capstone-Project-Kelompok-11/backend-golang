package models

import (
  "github.com/shopspring/decimal"
  "skfw/papaya/pigeon/easy"
)

type Courses struct {
  *easy.Model
  Name        string          `gorm:"type:VARCHAR(52);unique;not null" json:"name"`
  Description string          `gorm:"type:VARCHAR(200)" json:"description"`
  Price       decimal.Decimal `gorm:"default:0" json:"price"`
  Modules     []Modules       `gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"modules"`
}

func (Courses) TableName() string {

  return "courses"
}
