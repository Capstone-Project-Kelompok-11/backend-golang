package models

import (
  "skfw/papaya/pigeon/easy"
)

type Transactions struct {
  *easy.Model
  UserID        string  `gorm:"type:VARCHAR(32);not null" json:"user_id"`
  PaymentMethod string  `gorm:"type:VARCHAR(52);not null" json:"payment_method"`
  Verify        bool    `gorm:"default:FALSE" json:"verify"`
  Carts         []Carts `gorm:"foreignKey:TransactionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"carts"`
}

func (Transactions) TableName() string {

  return "transactions"
}
