package model

import (
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name        string `gorm:"unique;index:name_idx" json:"name"`
	Description string `json:"description"`
}

func (c *Category) Overwriting(other *Category) {
	if other == nil {
		return
	}
	if other.Name != "" {
		c.Name = other.Name
	}
	if other.Description != "" {
		c.Description = other.Description
	}
}
