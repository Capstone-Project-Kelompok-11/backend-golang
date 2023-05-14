package models

import (
	"database/sql"
	"gorm.io/gorm"
)

type Carts struct {
	*gorm.Model
	ID            string         `gorm:"type:VARCHAR(32);primary" json:"id"`
	UserID        string         `gorm:"type:VARCHAR(32);not null" json:"user_id"`
	TransactionID sql.NullString `gorm:"type:VARCHAR(32)" json:"transaction_id"`
	Qty           int            `gorm:"default:0" json:"qty"`
}

func (Carts) TableName() string {

	return "carts"
}
