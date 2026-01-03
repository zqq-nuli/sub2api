package dto

import "time"

// RechargeProductDTO 充值套餐数据传输对象
type RechargeProductDTO struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Amount        float64   `json:"amount"`        // 人民币金额
	Balance       float64   `json:"balance"`       // 基础美元余额
	BonusBalance  float64   `json:"bonus_balance"` // 赠送美元余额
	TotalBalance  float64   `json:"total_balance"` // 总到账余额
	Description   string    `json:"description"`
	SortOrder     int       `json:"sort_order"`
	IsHot         bool      `json:"is_hot"`
	DiscountLabel string    `json:"discount_label"`
	Status        string    `json:"status"` // active/inactive
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// CreateRechargeProductRequest 创建充值套餐请求
type CreateRechargeProductRequest struct {
	Name          string  `json:"name" binding:"required"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	Balance       float64 `json:"balance" binding:"required,gt=0"`
	BonusBalance  float64 `json:"bonus_balance" binding:"omitempty,gte=0"`
	Description   string  `json:"description"`
	SortOrder     int     `json:"sort_order"`
	IsHot         bool    `json:"is_hot"`
	DiscountLabel string  `json:"discount_label"`
}

// UpdateRechargeProductRequest 更新充值套餐请求
type UpdateRechargeProductRequest struct {
	Name          string  `json:"name" binding:"required"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	Balance       float64 `json:"balance" binding:"required,gt=0"`
	BonusBalance  float64 `json:"bonus_balance" binding:"omitempty,gte=0"`
	Description   string  `json:"description"`
	SortOrder     int     `json:"sort_order"`
	IsHot         bool    `json:"is_hot"`
	DiscountLabel string  `json:"discount_label"`
	Status        string  `json:"status" binding:"required,oneof=active inactive"`
}
