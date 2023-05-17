package models

import (
  "skfw/papaya/pigeon/easy"
)

type ReviewCourses struct {
  *easy.Model
  CourseID    string `gorm:"type:VARCHAR(32);not null" json:"course_id"`
  UserID      string `gorm:"type:VARCHAR(32);not null" json:"user_id"`
  Description string `gorm:"type:VARCHAR(200)" json:"description"`
  Rating      int    `json:"rating"`
}

func (ReviewCourses) TableName() string {

  return "review_courses"
}
