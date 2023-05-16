package models

import (
  "skfw/papaya/pigeon/easy"
)

/**
support multiple choices with multiple valid
*/

type Quizzes struct {
  *easy.Model
  ModuleID string `gorm:"type:VARCHAR(32);not null" json:"module_id"`
  Question string `json:"question"`
  ChoiceA  string `gorm:"type:VARCHAR(200)" json:"choice_a"`
  ValidA   bool   `json:"valid_a"`
  ChoiceB  string `gorm:"type:VARCHAR(200)" json:"choice_b"`
  ValidB   bool   `json:"valid_b"`
  ChoiceC  string `gorm:"type:VARCHAR(200)" json:"choice_c"`
  ValidC   bool   `json:"valid_c"`
  ChoiceD  string `gorm:"type:VARCHAR(200)" json:"choice_d"`
  ValidD   bool   `json:"valid_d"`
}

func (Quizzes) TableName() string {

  return "quizzes"
}
