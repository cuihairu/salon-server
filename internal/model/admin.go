package model

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	Name     string  `gorm:"type:varchar(255);unique;index:name_idx" json:"name"`
	Phone    *string `gorm:"type:varchar(255);unique;index:phone_idx" json:"phone"`
	Role     string  `json:"role"`
	Password []byte  `gorm:"type:bytea;not null"`
	Salt     []byte  `gorm:"type:bytea;not null"`
}
