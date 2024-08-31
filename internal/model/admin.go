package model

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	Name     string `gorm:"type:varchar(255);unique"`
	Password string `gorm:"type:varchar(255);not null"`
	Salt     string `gorm:"type:varchar(255);not null"`
}
