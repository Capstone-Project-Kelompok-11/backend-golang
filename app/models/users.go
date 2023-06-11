package models

import (
  "github.com/shopspring/decimal"
  "skfw/papaya/pigeon/templates/basicAuth/models"
)

type Users struct {
  *models.UserModel
  Image             string              `json:"image"`
  Balance           decimal.Decimal     `json:"balance"` // admin balance
  Courses           []Courses           `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"courses"`
  Reviews           []ReviewCourses     `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"reviews"`
  CompletionCourses []CompletionCourses `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"completion_courses"`
  CompletionModules []CompletionModules `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"completion_modules"`
  Assignments       []Assignments       `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"assignments"`
  Checkout          []Checkout          `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"checkout"`
  Events            []Events            `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"events"`
}

func (Users) TableName() string {

  return "users"
}
