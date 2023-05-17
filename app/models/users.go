package models

import (
  "github.com/shopspring/decimal"
  "skfw/papaya/pigeon/templates/basicAuth/models"
)

type Users struct {
  *models.UserModel
  Balance           decimal.Decimal     `gorm:"default:0" json:"balance"`
  Courses           []Courses           `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"courses"`
  Reviews           []ReviewCourses     `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"reviews"`
  CompletionCourses []CompletionCourses `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"completion_courses"`
  CompletionModules []CompletionModules `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"completion_modules"`
  Transactions      []Transactions      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"transactions"`
  Carts             []Carts             `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"carts"`
}

func (Users) TableName() string {

  return "users"
}
