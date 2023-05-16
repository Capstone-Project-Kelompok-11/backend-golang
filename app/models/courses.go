package models

import (
  "github.com/shopspring/decimal"
  "skfw/papaya/pigeon/easy"
)

/**

Pre Cache Rating

*/

type Courses struct {
  *easy.Model
  UserID            string              `gorm:"type:VARCHAR(32);not null" json:"user_id"`
  Name              string              `gorm:"type:VARCHAR(52);unique;not null" json:"name"`
  Description       string              `json:"description"`
  Thumbnail         string              `json:"thumbnail"`
  Price             decimal.Decimal     `json:"price"`
  Level             string              `json:"level"`
  Rating5           int                 `json:"rating_5"`
  Rating4           int                 `json:"rating_4"`
  Rating3           int                 `json:"rating_3"`
  Rating2           int                 `json:"rating_2"`
  Rating1           int                 `json:"rating_1"`
  Modules           []Modules           `gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"modules"`
  Reviews           []Reviews           `gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"reviews"`
  CompletionCourses []CompletionCourses `gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"completion_courses"`
  CategoryCourses   []CategoryCourses   `gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"category_courses"`
}

func (Courses) TableName() string {

  return "courses"
}
