package repository

import (
	"context"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"gorm.io/gorm"
)

// orderModel GORM数据库模型
type orderModel struct {
	ID             int64      `gorm:"primaryKey"`
	OrderNo        string     `gorm:"uniqueIndex;size:32;not null"`
	UserID         int64      `gorm:"index:idx_orders_user_id;not null"`
	ProductID      *int64     `gorm:"index"`
	ProductName    string     `gorm:"size:255;not null"`
	Amount         float64    `gorm:"type:decimal(20,2);not null"` // 人民币金额
	BonusAmount    float64    `gorm:"type:decimal(20,8);default:0"`
	ActualAmount   float64    `gorm:"type:decimal(20,8);not null"` // 美元金额
	PaymentMethod  string     `gorm:"size:50;not null"`
	PaymentGateway string     `gorm:"size:50;default:'epay';not null"`
	TradeNo        string     `gorm:"size:255"`
	Status         string     `gorm:"index:idx_orders_status;size:20;default:'pending';not null"`
	CreatedAt      time.Time  `gorm:"index:idx_orders_created_at;not null"`
	PaidAt         *time.Time
	ExpiredAt      time.Time  `gorm:"not null"`
	Notes          string  `gorm:"type:text"`
	CallbackData   *string `gorm:"type:jsonb"`

	// 关联
	User    *userModel            `gorm:"foreignKey:UserID;references:ID"`
	Product *rechargeProductModel `gorm:"foreignKey:ProductID;references:ID"`
}

func (orderModel) TableName() string {
	return "orders"
}

type orderRepository struct {
	db *gorm.DB
}

// NewOrderRepository 创建订单Repository
func NewOrderRepository(db *gorm.DB) service.OrderRepository {
	return &orderRepository{db: db}
}

// Create 创建订单
func (r *orderRepository) Create(ctx context.Context, order *service.Order) error {
	m := orderModelFromService(order)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	order.ID = m.ID
	return nil
}

// GetByID 根据ID获取订单
func (r *orderRepository) GetByID(ctx context.Context, id int64) (*service.Order, error) {
	var m orderModel
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Product").
		First(&m, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, service.ErrOrderNotFound
		}
		return nil, err
	}
	return orderModelToService(&m), nil
}

// GetByOrderNo 根据订单号获取订单
func (r *orderRepository) GetByOrderNo(ctx context.Context, orderNo string) (*service.Order, error) {
	var m orderModel
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Product").
		Where("order_no = ?", orderNo).
		First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, service.ErrOrderNotFound
		}
		return nil, err
	}
	return orderModelToService(&m), nil
}

// Update 更新订单
func (r *orderRepository) Update(ctx context.Context, order *service.Order) error {
	m := orderModelFromService(order)
	return r.db.WithContext(ctx).Save(m).Error
}

// UpdateWithTx 使用事务更新订单
func (r *orderRepository) UpdateWithTx(tx *gorm.DB, order *service.Order) error {
	m := orderModelFromService(order)
	return tx.Save(m).Error
}

// ListByUser 获取用户订单列表
func (r *orderRepository) ListByUser(ctx context.Context, userID int64, page, limit int) ([]*service.Order, int64, error) {
	var models []orderModel
	var total int64

	query := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC")

	// 查询总数
	if err := query.Model(&orderModel{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).
		Preload("Product").
		Find(&models).Error; err != nil {
		return nil, 0, err
	}

	orders := make([]*service.Order, len(models))
	for i, m := range models {
		orders[i] = orderModelToService(&m)
	}

	return orders, total, nil
}

// AdminList 管理员查询订单列表（支持灵活筛选）
func (r *orderRepository) AdminList(ctx context.Context, filter *service.AdminListOrdersFilter) ([]*service.Order, int64, error) {
	var models []orderModel
	var total int64

	query := r.db.WithContext(ctx).Preload("User").Preload("Product")

	// 应用过滤条件
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.PaymentMethod != nil {
		query = query.Where("payment_method = ?", *filter.PaymentMethod)
	}
	if filter.OrderNo != nil {
		query = query.Where("order_no LIKE ?", "%"+*filter.OrderNo+"%")
	}
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", *filter.EndDate)
	}

	// 查询总数
	if err := query.Model(&orderModel{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (filter.Page - 1) * filter.Limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(filter.Limit).
		Find(&models).Error; err != nil {
		return nil, 0, err
	}

	orders := make([]*service.Order, len(models))
	for i, m := range models {
		orders[i] = orderModelToService(&m)
	}

	return orders, total, nil
}

// GetStatistics 获取订单统计信息
func (r *orderRepository) GetStatistics(ctx context.Context) (*service.OrderStatistics, error) {
	stats := &service.OrderStatistics{}

	// 总订单数
	if err := r.db.WithContext(ctx).Model(&orderModel{}).Count(&stats.TotalOrders).Error; err != nil {
		return nil, err
	}

	// 各状态订单数
	statusCounts := []struct {
		Status string
		Count  int64
	}{}
	if err := r.db.WithContext(ctx).Model(&orderModel{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&statusCounts).Error; err != nil {
		return nil, err
	}

	for _, sc := range statusCounts {
		switch sc.Status {
		case service.OrderStatusPending:
			stats.PendingOrders = sc.Count
		case service.OrderStatusPaid:
			stats.PaidOrders = sc.Count
		case service.OrderStatusFailed:
			stats.FailedOrders = sc.Count
		case service.OrderStatusExpired:
			stats.ExpiredOrders = sc.Count
		}
	}

	// 总成交金额（已支付订单的人民币金额）
	if err := r.db.WithContext(ctx).Model(&orderModel{}).
		Where("status = ?", service.OrderStatusPaid).
		Select("COALESCE(SUM(amount), 0) as total").
		Scan(&stats.TotalAmount).Error; err != nil {
		return nil, err
	}

	// 总充值余额（已支付订单的美元余额）
	if err := r.db.WithContext(ctx).Model(&orderModel{}).
		Where("status = ?", service.OrderStatusPaid).
		Select("COALESCE(SUM(actual_amount), 0) as total").
		Scan(&stats.TotalBalance).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

// BeginTx 开始数据库事务
func (r *orderRepository) BeginTx(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).Begin()
}

// CountPendingByUser 统计用户未支付订单数量
func (r *orderRepository) CountPendingByUser(ctx context.Context, userID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&orderModel{}).
		Where("user_id = ? AND status = ?", userID, service.OrderStatusPending).
		Count(&count).Error
	return count, err
}

// MarkExpiredOrders 将过期的pending订单标记为expired
func (r *orderRepository) MarkExpiredOrders(ctx context.Context) (int64, error) {
	result := r.db.WithContext(ctx).Model(&orderModel{}).
		Where("status = ? AND expired_at < ?", service.OrderStatusPending, time.Now()).
		Update("status", service.OrderStatusExpired)
	return result.RowsAffected, result.Error
}

// orderModelFromService 从Service对象转换为GORM模型
func orderModelFromService(o *service.Order) *orderModel {
	if o == nil {
		return nil
	}
	m := &orderModel{
		ID:             o.ID,
		OrderNo:        o.OrderNo,
		UserID:         o.UserID,
		ProductID:      o.ProductID,
		ProductName:    o.ProductName,
		Amount:         o.Amount,
		BonusAmount:    o.BonusAmount,
		ActualAmount:   o.ActualAmount,
		PaymentMethod:  o.PaymentMethod,
		PaymentGateway: o.PaymentGateway,
		TradeNo:        o.TradeNo,
		Status:         o.Status,
		CreatedAt:      o.CreatedAt,
		PaidAt:         o.PaidAt,
		ExpiredAt:      o.ExpiredAt,
		Notes:          o.Notes,
	}
	// CallbackData: 空字符串转为 nil（SQL NULL）
	if o.CallbackData != "" {
		m.CallbackData = &o.CallbackData
	}
	return m
}

// orderModelToService 从GORM模型转换为Service对象
func orderModelToService(m *orderModel) *service.Order {
	if m == nil {
		return nil
	}

	order := &service.Order{
		ID:             m.ID,
		OrderNo:        m.OrderNo,
		UserID:         m.UserID,
		ProductID:      m.ProductID,
		ProductName:    m.ProductName,
		Amount:         m.Amount,
		BonusAmount:    m.BonusAmount,
		ActualAmount:   m.ActualAmount,
		PaymentMethod:  m.PaymentMethod,
		PaymentGateway: m.PaymentGateway,
		TradeNo:        m.TradeNo,
		Status:         m.Status,
		CreatedAt:      m.CreatedAt,
		PaidAt:         m.PaidAt,
		ExpiredAt:      m.ExpiredAt,
		Notes:          m.Notes,
	}
	// CallbackData: nil 转为空字符串
	if m.CallbackData != nil {
		order.CallbackData = *m.CallbackData
	}

	// 关联对象
	if m.User != nil {
		order.User = userModelToService(m.User)
	}
	if m.Product != nil {
		order.Product = rechargeProductModelToService(m.Product)
	}

	return order
}
