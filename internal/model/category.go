package model

import (
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name        string `gorm:"unique;index:name_idx" json:"name"`
	Description string `json:"description"`
}

func (c *Category) overwriting(other *Category) {
	if other == nil {
		return
	}
	c.Name = other.Name
	c.Description = other.Description
}
