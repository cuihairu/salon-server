package data

import "gorm.io/gorm"

type Member struct {
	gorm.Model
	Level  uint
	Points uint64
}
