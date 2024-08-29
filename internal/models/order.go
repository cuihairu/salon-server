package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID uint
	Amount float64
}
