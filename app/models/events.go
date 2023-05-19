package models

import (
  "skfw/papaya/pigeon/easy"
)

/**

Pre Cache Rating

*/

type Events struct {
  *easy.Model
  UserID      string `gorm:"type:VARCHAR(32);not null" json:"user_id"`
  Name        string `json:"name"`
  Description string `json:"description"`
}

func (Events) TableName() string {

  return "events"
}
