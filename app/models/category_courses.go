package models

import (
  "skfw/papaya/pigeon/easy"
)

type CategoryCourses struct {
  *easy.Model
  CategoryID string `gorm:"type:VARCHAR(32);not null" json:"category_id"`
  CourseID   string `gorm:"type:VARCHAR(32);not null" json:"course_id"`
}

func (CategoryCourses) TableName() string {

  return "category_courses"
}
