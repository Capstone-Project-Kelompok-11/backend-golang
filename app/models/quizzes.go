package models

import (
  "skfw/papaya/pigeon/easy"
)

/**
support multiple choices with multiple valid
*/

type Quizzes struct {
  *easy.Model
  ModuleID string `gorm:"type:VARCHAR(32);unique;not null" json:"module_id"`
  Data     string `gorm:"type:TEXT;not null" json:"data"` // JSON: [{ "question": "text", "choices": [ { "text": "", "valid": true } ] }]
  //Valid    string `gorm:"type:VARCHAR(64);unique;not null" json:"valid"`
}

func (Quizzes) TableName() string {

  return "quizzes"
}
