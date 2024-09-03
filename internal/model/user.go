package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	WechatId *string    `gorm:"uniqueIndex" json:"wechat_id,omitempty"`
	AlipayId *string    `gorm:"uniqueIndex" json:"alipay_id,omitempty"`
	DouyinId *string    `gorm:"uniqueIndex" json:"douyin_id,omitempty"`
	Name     string     `json:"name"`
	Phone    *string    `gorm:"uniqueIndex" json:"phone,omitempty"`
	Birthday *time.Time `json:"birthday,omitempty"`
	Address  *string    `json:"address,omitempty"`
}
