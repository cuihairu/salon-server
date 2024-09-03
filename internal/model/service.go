package model

import "gorm.io/gorm"

type Service struct {
	gorm.Model
	Name         string   `gorm:"unique;index:name_idx" json:"name"`
	CategoryId   uint     `gorm:"index:category_idx" json:"category_id"`
	CategoryName string   `gorm:"-" json:"category_name"` // 关联 Category
	Intro        string   `json:"intro"`
	Cover        string   `json:"cover"`                    // 封面图片
	Images       []string `gorm:"type:jsonb" json:"images"` // 详情图片
	Description  string   `json:"description"`
	Duration     int      `json:"duration"`           // 分钟
	Price        float64  `gorm:"type:decimal(10,2)"` // 显示价格
	Amount       float64  `gorm:"type:decimal(10,2)"` // 当前售价
}
