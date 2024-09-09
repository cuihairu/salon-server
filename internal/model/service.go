package model

import (
	"gorm.io/gorm"
)

type Images []string

type Service struct {
	gorm.Model
	Name         string  `gorm:"unique;index:name_idx" json:"name"`     // 服务名
	CategoryId   uint    `gorm:"index:category_idx" json:"category_id"` // 关联 Category
	CategoryName string  `gorm:"-" json:"category_name"`                // 关联 Category
	Intro        string  `json:"intro"`                                 // 简介
	Cover        string  `json:"cover"`                                 // 封面图片
	Content      string  `gorm:"type:text" json:"content"`              // 内容
	Duration     int     `json:"duration"`                              // 分钟
	Price        float64 `gorm:"type:decimal(10,2)"`                    // 显示价格
	Amount       float64 `gorm:"type:decimal(10,2)"`                    // 当前售价
	Recommend    bool    `json:"recommend"`                             // 是否推荐
}
