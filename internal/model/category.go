package model

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name string `gorm:"unique;index:name_idx" json:"name"`
}
