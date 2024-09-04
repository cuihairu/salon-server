package model

import (
	"gorm.io/gorm"
	"time"
)

type OrderService struct {
	ID       uint    `json:"id"`       // 服务 ID
	Name     string  `json:"name"`     // 服务名字
	Duration int     `json:"duration"` // 服务时长
	Cover    string  `json:"cover"`    // 服务封面
	Price    float64 `json:"price"`    // 服务价格
	Amount   float64 `json:"amount"`   // 总价
	Discount float64 `json:"discount"` // 折扣
	Num      int     `json:"num"`      // 购买数量
}

type OrderState int

const (
	OrderCreated OrderState = iota
	OrderDone
	OrderExpired
	OrderCanceled
)

type Order struct {
	gorm.Model
	OrderId       string         `json:"order_id" gorm:"unique;index:order_id_idx"`
	UserID        uint           `json:"user_id"`
	Title         string         `json:"title"`
	OrderServices []OrderService `gorm:"type:jsonb" json:"services"`
	Price         float64        `gorm:"type:decimal(10,2)" json:"price"`      // 原价
	Discounted    float64        `gorm:"type:decimal(10,2)" json:"discounted"` // 打折之后的价格
	DoneAt        *time.Time     `json:"done_at,omitempty"`                    // 完成时间
	ExpiredAt     *time.Time     `json:"expired_at"`                           // 过期时间
	Status        int            `json:"status"`                               // 订单状态 0 created 1 done 2 expired 3 canceled
}
