package service

import (
	"context"
	"time"
)

// RechargeProductRepository 充值套餐数据访问接口
type RechargeProductRepository interface {
	Create(ctx context.Context, product *RechargeProduct) error
	GetByID(ctx context.Context, id int64) (*RechargeProduct, error)
	Update(ctx context.Context, product *RechargeProduct) error
	Delete(ctx context.Context, id int64) error
	ListActive(ctx context.Context) ([]*RechargeProduct, error)
	ListAll(ctx context.Context) ([]*RechargeProduct, error)
}

// RechargeProduct 充值套餐领域模型
type RechargeProduct struct {
	ID            int64
	Name          string  // 套餐名称（如"10元充值"）
	Amount        float64 // 人民币金额
	Balance       float64 // 美元余额（基础金额）
	BonusBalance  float64 // 赠送余额（美元）
	Description   string  // 套餐描述
	SortOrder     int     // 排序权重（数字越小越靠前）
	IsHot         bool    // 是否为热门套餐
	DiscountLabel string  // 折扣标签（如"限时优惠"）
	Status        string  // active/inactive
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// 套餐状态常量
const (
	ProductStatusActive   = "active"   // 启用
	ProductStatusInactive = "inactive" // 禁用
)

// IsActive 是否启用
func (p *RechargeProduct) IsActive() bool {
	return p.Status == ProductStatusActive
}

// GetTotalBalance 获取总到账金额（基础+赠送）
func (p *RechargeProduct) GetTotalBalance() float64 {
	return p.Balance + p.BonusBalance
}
