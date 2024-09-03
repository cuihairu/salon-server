package model

import (
	"gorm.io/gorm"
	"time"
)

type Order struct {
	gorm.Model
	UserID     uint       `json:"user_id"`
	Title      string     `json:"title"`
	ServID     uint       `json:"serv_id"`
	Price      float64    `gorm:"type:decimal(10,2)" json:"price"`      // 原价
	Discounted float64    `gorm:"type:decimal(10,2)" json:"discounted"` // 打折之后的价格
	DoneAt     *time.Time `json:"done_at,omitempty"`                    // 完成时间
}
