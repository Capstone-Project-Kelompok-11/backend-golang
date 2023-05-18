package models

import (
  "skfw/papaya/pigeon/easy"
)

type Checkout struct {
  *easy.Model
  UserID        string `gorm:"type:VARCHAR(32);not null" json:"user_id"`
  CourseID      string `gorm:"type:VARCHAR(32);not null" json:"course_id"`
  PaymentMethod string `gorm:"type:VARCHAR(52);not null" json:"payment_method"`
  Verify        bool   `gorm:"default:FALSE" json:"verify"`
}

func (Checkout) TableName() string {

  return "checkout"
}
