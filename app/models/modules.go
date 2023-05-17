package models

import (
  "skfw/papaya/pigeon/easy"
)

type Modules struct {
  *easy.Model
  CourseID          string              `gorm:"type:VARCHAR(32);not null" json:"course_id"`
  Name              string              `gorm:"type:VARCHAR(52);unique;not null" json:"name"`
  Description       string              `json:"description"`
  Thumbnail         string              `json:"thumbnail"`
  Video             string              `json:"video"`
  Document          string              `json:"document"`
  Quizzes           []Quizzes           `gorm:"foreignKey:ModuleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"quizzes"`
  CompletionModules []CompletionModules `gorm:"foreignKey:ModuleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"completion_modules"`
}

func (Modules) TableName() string {

  return "modules"
}
