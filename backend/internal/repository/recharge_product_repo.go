package repository

import (
	"context"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"gorm.io/gorm"
)

// rechargeProductModel GORM数据库模型
type rechargeProductModel struct {
	ID            int64     `gorm:"primaryKey"`
	Name          string    `gorm:"size:255;not null"`
	Amount        float64   `gorm:"type:decimal(20,2);not null"`
	Balance       float64   `gorm:"type:decimal(20,8);not null"`
	BonusBalance  float64   `gorm:"type:decimal(20,8);default:0"`
	Description   string    `gorm:"type:text"`
	SortOrder     int       `gorm:"index:idx_products_status_sort;default:0"`
	IsHot         bool      `gorm:"default:false"`
	DiscountLabel string    `gorm:"size:50"`
	Status        string    `gorm:"index:idx_products_status_sort;size:20;default:'active';not null"`
	CreatedAt     time.Time `gorm:"not null"`
	UpdatedAt     time.Time `gorm:"not null"`
}

func (rechargeProductModel) TableName() string {
	return "recharge_products"
}

type rechargeProductRepository struct {
	db *gorm.DB
}

// NewRechargeProductRepository 创建充值套餐Repository
func NewRechargeProductRepository(db *gorm.DB) service.RechargeProductRepository {
	return &rechargeProductRepository{db: db}
}

// Create 创建充值套餐
func (r *rechargeProductRepository) Create(ctx context.Context, product *service.RechargeProduct) error {
	m := rechargeProductModelFromService(product)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	product.ID = m.ID
	return nil
}

// GetByID 根据ID获取套餐
func (r *rechargeProductRepository) GetByID(ctx context.Context, id int64) (*service.RechargeProduct, error) {
	var m rechargeProductModel
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, service.ErrProductNotFound
		}
		return nil, err
	}
	return rechargeProductModelToService(&m), nil
}

// Update 更新充值套餐
func (r *rechargeProductRepository) Update(ctx context.Context, product *service.RechargeProduct) error {
	m := rechargeProductModelFromService(product)
	return r.db.WithContext(ctx).Save(m).Error
}

// Delete 删除充值套餐（硬删除）
func (r *rechargeProductRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Delete(&rechargeProductModel{}, id).Error
}

// ListActive 获取启用的套餐列表（按sort_order排序）
func (r *rechargeProductRepository) ListActive(ctx context.Context) ([]*service.RechargeProduct, error) {
	var models []rechargeProductModel
	if err := r.db.WithContext(ctx).
		Where("status = ?", service.ProductStatusActive).
		Order("sort_order ASC, id ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	products := make([]*service.RechargeProduct, len(models))
	for i, m := range models {
		products[i] = rechargeProductModelToService(&m)
	}

	return products, nil
}

// ListAll 获取所有套餐列表（管理员）
func (r *rechargeProductRepository) ListAll(ctx context.Context) ([]*service.RechargeProduct, error) {
	var models []rechargeProductModel
	if err := r.db.WithContext(ctx).
		Order("sort_order ASC, id ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	products := make([]*service.RechargeProduct, len(models))
	for i, m := range models {
		products[i] = rechargeProductModelToService(&m)
	}

	return products, nil
}

// rechargeProductModelFromService 从Service对象转换为GORM模型
func rechargeProductModelFromService(p *service.RechargeProduct) *rechargeProductModel {
	if p == nil {
		return nil
	}
	return &rechargeProductModel{
		ID:            p.ID,
		Name:          p.Name,
		Amount:        p.Amount,
		Balance:       p.Balance,
		BonusBalance:  p.BonusBalance,
		Description:   p.Description,
		SortOrder:     p.SortOrder,
		IsHot:         p.IsHot,
		DiscountLabel: p.DiscountLabel,
		Status:        p.Status,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

// rechargeProductModelToService 从GORM模型转换为Service对象
func rechargeProductModelToService(m *rechargeProductModel) *service.RechargeProduct {
	if m == nil {
		return nil
	}
	return &service.RechargeProduct{
		ID:            m.ID,
		Name:          m.Name,
		Amount:        m.Amount,
		Balance:       m.Balance,
		BonusBalance:  m.BonusBalance,
		Description:   m.Description,
		SortOrder:     m.SortOrder,
		IsHot:         m.IsHot,
		DiscountLabel: m.DiscountLabel,
		Status:        m.Status,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}
