package service

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// OrderRepository 订单数据访问接口
type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, id int64) (*Order, error)
	GetByOrderNo(ctx context.Context, orderNo string) (*Order, error)
	Update(ctx context.Context, order *Order) error
	UpdateWithTx(tx *gorm.DB, order *Order) error
	ListByUser(ctx context.Context, userID int64, page, limit int) ([]*Order, int64, error)
	AdminList(ctx context.Context, filter *AdminListOrdersFilter) ([]*Order, int64, error)
	GetStatistics(ctx context.Context) (*OrderStatistics, error)
	BeginTx(ctx context.Context) *gorm.DB
	// CountPendingByUser 统计用户未支付订单数量
	CountPendingByUser(ctx context.Context, userID int64) (int64, error)
	// MarkExpiredOrders 将过期的pending订单标记为expired
	MarkExpiredOrders(ctx context.Context) (int64, error)
}

// AdminListOrdersFilter 管理员查询订单过滤条件
type AdminListOrdersFilter struct {
	UserID        *int64
	Status        *string
	PaymentMethod *string
	OrderNo       *string
	StartDate     *time.Time
	EndDate       *time.Time
	Page          int
	Limit         int
}

// Order 订单领域模型
type Order struct {
	ID             int64
	OrderNo        string // 订单号：ORD{timestamp}{random6}
	UserID         int64
	ProductID      *int64  // 可为空（套餐可能被删除）
	ProductName    string  // 冗余存储，防止套餐删除后丢失信息
	Amount         float64 // 支付金额（人民币）
	BonusAmount    float64 // 赠送金额（美元）
	ActualAmount   float64 // 实际到账金额（美元，含赠送）
	PaymentMethod  string  // alipay/wxpay/epusdt
	PaymentGateway string  // epay
	TradeNo        string  // 第三方支付平台订单号
	Status         string  // pending/paid/failed/expired/refunded
	CreatedAt      time.Time
	PaidAt         *time.Time
	ExpiredAt      time.Time // 订单过期时间（15分钟有效期）
	Notes          string    // 管理员备注
	CallbackData   string    // 易支付回调原始数据（JSON）

	// 关联对象
	User    *User
	Product *RechargeProduct
}

// 订单状态常量
const (
	OrderStatusPending  = "pending"  // 待支付
	OrderStatusPaid     = "paid"     // 已支付
	OrderStatusFailed   = "failed"   // 支付失败
	OrderStatusExpired  = "expired"  // 已过期
	OrderStatusRefunded = "refunded" // 已退款
)

// IsPending 是否待支付
func (o *Order) IsPending() bool {
	return o.Status == OrderStatusPending
}

// IsPaid 是否已支付
func (o *Order) IsPaid() bool {
	return o.Status == OrderStatusPaid
}

// IsFailed 是否支付失败
func (o *Order) IsFailed() bool {
	return o.Status == OrderStatusFailed
}

// IsExpired 是否已过期
func (o *Order) IsExpired() bool {
	return o.Status == OrderStatusExpired || time.Now().After(o.ExpiredAt)
}

// CanRefund 是否可以退款
func (o *Order) CanRefund() bool {
	return o.IsPaid()
}
