package model

import "gorm.io/gorm"

type Service struct {
	gorm.Model
	Name        string `gorm:"unique;index:name_idx"`
	CategoryId  uint   `gorm:"index:category_idx"`
	Intro       string
	Cover       string   // 封面图片
	Images      []string `gorm:"type:jsonb"` // 详情图片
	Description string
	Time        int     // 分钟
	Price       float64 `gorm:"type:decimal(10,2)"` // 显示价格
	Amount      float64 `gorm:"type:decimal(10,2)"` // 当前售价
}
