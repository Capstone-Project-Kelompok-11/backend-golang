package models

import (
  "skfw/papaya/pigeon/easy"
)

/**

Pre Cache Rating

*/

type CompletionCourses struct {
  *easy.Model
  UserID   string `gorm:"type:VARCHAR(32);not null" json:"user_id"`
  CourseID string `gorm:"type:VARCHAR(32);not null" json:"course_id"`
  Score    int    `json:"score"`
}

func (CompletionCourses) TableName() string {

  return "completion_courses"
}
