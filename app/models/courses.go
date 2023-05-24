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
  Video             string              `json:"video"`
  Document          string              `json:"document"`
  Price             decimal.Decimal     `json:"price"`
  Level             string              `json:"level"`
  Rating5           int                 `gorm:"type:INTEGER;default:0" json:"rating_5"`
  Rating4           int                 `gorm:"type:INTEGER;default:0" json:"rating_4"`
  Rating3           int                 `gorm:"type:INTEGER;default:0" json:"rating_3"`
  Rating2           int                 `gorm:"type:INTEGER;default:0" json:"rating_2"`
  Rating1           int                 `gorm:"type:INTEGER;default:0" json:"rating_1"`
  Finished          int                 `gorm:"type:INTEGER;default:0" json:"finished"` // user has been finished this course
  MemberCount       int                 `gorm:"type:INTEGER;default:0" json:"member_count"`
  Modules           []Modules           `gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"modules"`
  Reviews           []ReviewCourses     `gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"reviews"`
  CompletionCourses []CompletionCourses `gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"completion_courses"`
  CategoryCourses   []CategoryCourses   `gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"category_courses"`
  Assignments       []Assignments       `gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"assignments"`
  Checkout          []Checkout          `gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"checkout"`
}

func (Courses) TableName() string {

  return "courses"
}
