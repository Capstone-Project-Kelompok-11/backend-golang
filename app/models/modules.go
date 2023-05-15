package models

import (
  "skfw/papaya/pigeon/easy"
)

type Modules struct {
  *easy.Model
  CourseID    string `gorm:"type:VARCHAR(32);not null" json:"course_id"`
  Name        string `gorm:"type:VARCHAR(52);unique;not null" json:"name"`
  Description string `gorm:"type:VARCHAR(200)" json:"description"`
}

func (Modules) TableName() string {

  return "modules"
}
