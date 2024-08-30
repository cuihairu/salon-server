package data

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	WechatId *string `gorm:"uniqueIndex"`
	AlipayId *string `gorm:"uniqueIndex"`
	DouYinId *string `gorm:"uniqueIndex"`
	Name     string
	Phone    *string `gorm:"uniqueIndex"`
	Birthday *time.Time
	Address  *string
}
