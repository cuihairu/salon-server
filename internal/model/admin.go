package model

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	Name     string  `gorm:"type:varchar(255);unique;index:name_idx"`
	Phone    *string `gorm:"type:varchar(255);unique;index:phone_idx"`
	Password string  `gorm:"type:varchar(255);not null"`
	Salt     string  `gorm:"type:varchar(255);not null"`
}
