package model

import "gorm.io/gorm"

type Member struct {
	gorm.Model
	Level  uint   `json:"level"`
	Points uint64 `json:"points"`
}
