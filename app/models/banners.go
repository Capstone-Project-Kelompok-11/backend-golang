package models

import "skfw/papaya/pigeon/easy"

type Banners struct {
  *easy.Model
  Alt string `json:"alt"`
  Src string `json:"src"`
}

func (Banners) TableName() string {

  return "banners"
}
