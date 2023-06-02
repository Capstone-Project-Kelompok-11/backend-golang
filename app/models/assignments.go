package models

import (
  "skfw/papaya/pigeon/easy"
)

/**

Pre Cache Rating

*/

type Assignments struct {
  *easy.Model
  UserID   string `gorm:"type:VARCHAR(32);not null" json:"user_id"`
  CourseID string `gorm:"type:VARCHAR(32);not null" json:"course_id"`
  Video    string `json:"video"`
  Document string `json:"document"`
  Score    int    `json:"score"`
}

func (Assignments) TableName() string {

  return "assignments"
}
