package model

import "gorm.io/gorm"

type OperationLog struct {
	gorm.Model
	Username string `gorm:"type:varchar(255);not null;index:username_idx" json:"username"`
	Role     string `gorm:"type:varchar(255);not null" json:"role"`
	Ip       string `gorm:"type:varchar(255);not null" json:"ip"`
	Location string `gorm:"type:varchar(255);not null" json:"location"`
	Agent    string `gorm:"type:varchar(255);not null" json:"agent"`
	Table    string `gorm:"type:varchar(255);not null" json:"table"`
	Action   string `gorm:"type:varchar(255);not null" json:"action"`
	Content  string `gorm:"type:text" json:"content"`
	Err      string `gorm:"type:text" json:"err"`
}
