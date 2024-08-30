package data

import (
	"gorm.io/gorm"
	"time"
)

type Order struct {
	gorm.Model
	UserID uint
	Title  string
	ServID uint
	Price  float64    `gorm:"type:decimal(10,2)"`
	Amount float64    `gorm:"type:decimal(10,2)"` // 原价
	DoneAt *time.Time // 完成时间
}
