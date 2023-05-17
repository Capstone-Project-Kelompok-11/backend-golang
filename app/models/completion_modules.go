package models

import (
  "skfw/papaya/pigeon/easy"
)

/**

Pre Cache Rating

*/

type CompletionModules struct {
  *easy.Model
  UserID   string `gorm:"type:VARCHAR(32);not null" json:"user_id"`
  CourseID string `gorm:"type:VARCHAR(32);not null" json:"course_id"`
  ModuleID string `gorm:"type:VARCHAR(32);not null" json:"module_id"`
  Score    int    `json:"score"`
}

func (CompletionModules) TableName() string {

  return "completion_modules"
}
