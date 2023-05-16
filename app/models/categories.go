package models

import (
  "skfw/papaya/pigeon/easy"
)

type Categories struct {
  *easy.Model
  Name            string            `gorm:"type:VARCHAR(52);unique;not null" json:"name"`
  Description     string            `gorm:"type:VARCHAR(200)" json:"description"`
  Thumbnail       string            `json:"thumbnail"`
  CategoryCourses []CategoryCourses `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"category_courses"`
}

func (Categories) TableName() string {

  return "categories"
}
