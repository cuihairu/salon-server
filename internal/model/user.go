package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	WechatId string `json:"wechat_id" gorm:"unique"`
	Name     string `json:"name"`
	PhoneNo  string
}
