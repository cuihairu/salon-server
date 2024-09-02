package model

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	Name     string  `gorm:"type:varchar(255);unique;index:name_idx"`
	Phone    *string `gorm:"type:varchar(255);unique;index:phone_idx"`
	Password []byte  `gorm:"type:bytea;not null"`
	Salt     []byte  `gorm:"type:bytea;not null"`
}
