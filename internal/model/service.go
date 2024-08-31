package model

import "gorm.io/gorm"

type Service struct {
	gorm.Model
	Name        string
	Description string
	Price       float64 `gorm:"type:decimal(10,2)"` // 显示价格
	Amount      float64 `gorm:"type:decimal(10,2)"` // 当前售价
}
