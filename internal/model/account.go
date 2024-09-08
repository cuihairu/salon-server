package model

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	Balance  float64 `gorm:"type:decimal(10,2)" json:"balance"`  // DECIMAL(10,2) 表示总共 10 位，其中 2 位是小数位
	Consumed float64 `gorm:"type:decimal(10,2)" json:"consumed"` //
}
