package models

import (
	"skfw/papaya/pigeon/easy"
)

type Articles struct {
	*easy.Model
	CourseID           string     `gorm:"type:VARCHAR(32);not null" json:"course_id"`
	Name               string     `gorm:"type:VARCHAR(52);unique;not null" json:"name"`
	Description        string     `json:"description"`
	Document           string     `json:"document"`
	CompletionArticles []Articles `gorm:"foreignKey:ArticleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"completion_modules"`
}

func (Articles) TableName() string {

	return "articles"
}
