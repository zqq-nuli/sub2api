package repository

import (
	"context"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	dbrechargeproduct "github.com/Wei-Shaw/sub2api/ent/rechargeproduct"
	apperrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"entgo.io/ent/dialect/sql"
)

type rechargeProductRepository struct {
	client *dbent.Client
}

// NewRechargeProductRepository 创建充值套餐Repository
func NewRechargeProductRepository(client *dbent.Client) service.RechargeProductRepository {
	return &rechargeProductRepository{client: client}
}

// Create 创建充值套餐
func (r *rechargeProductRepository) Create(ctx context.Context, product *service.RechargeProduct) error {
	created, err := r.client.RechargeProduct.Create().
		SetName(product.Name).
		SetAmount(product.Amount).
		SetBalance(product.Balance).
		SetBonusBalance(product.BonusBalance).
		SetDescription(product.Description).
		SetSortOrder(product.SortOrder).
		SetIsHot(product.IsHot).
		SetDiscountLabel(product.DiscountLabel).
		SetStatus(product.Status).
		Save(ctx)
	if err != nil {
		return err
	}
	product.ID = created.ID
	product.CreatedAt = created.CreatedAt
	product.UpdatedAt = created.UpdatedAt
	return nil
}

// GetByID 根据ID获取套餐
func (r *rechargeProductRepository) GetByID(ctx context.Context, id int64) (*service.RechargeProduct, error) {
	m, err := r.client.RechargeProduct.Query().
		Where(dbrechargeproduct.IDEQ(id)).
		Only(ctx)
	if err != nil {
		return nil, translateProductError(err)
	}
	return rechargeProductEntityToService(m), nil
}

// Update 更新充值套餐
func (r *rechargeProductRepository) Update(ctx context.Context, product *service.RechargeProduct) error {
	_, err := r.client.RechargeProduct.UpdateOneID(product.ID).
		SetName(product.Name).
		SetAmount(product.Amount).
		SetBalance(product.Balance).
		SetBonusBalance(product.BonusBalance).
		SetDescription(product.Description).
		SetSortOrder(product.SortOrder).
		SetIsHot(product.IsHot).
		SetDiscountLabel(product.DiscountLabel).
		SetStatus(product.Status).
		Save(ctx)
	return err
}

// Delete 删除充值套餐（硬删除）
func (r *rechargeProductRepository) Delete(ctx context.Context, id int64) error {
	return r.client.RechargeProduct.DeleteOneID(id).Exec(ctx)
}

// ListActive 获取启用的套餐列表（按sort_order排序）
func (r *rechargeProductRepository) ListActive(ctx context.Context) ([]*service.RechargeProduct, error) {
	products, err := r.client.RechargeProduct.Query().
		Where(dbrechargeproduct.StatusEQ(service.ProductStatusActive)).
		Order(dbrechargeproduct.BySortOrder(sql.OrderAsc()), dbrechargeproduct.ByID(sql.OrderAsc())).
		All(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*service.RechargeProduct, len(products))
	for i, m := range products {
		result[i] = rechargeProductEntityToService(m)
	}

	return result, nil
}

// ListAll 获取所有套餐列表（管理员）
func (r *rechargeProductRepository) ListAll(ctx context.Context) ([]*service.RechargeProduct, error) {
	products, err := r.client.RechargeProduct.Query().
		Order(dbrechargeproduct.BySortOrder(sql.OrderAsc()), dbrechargeproduct.ByID(sql.OrderAsc())).
		All(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*service.RechargeProduct, len(products))
	for i, m := range products {
		result[i] = rechargeProductEntityToService(m)
	}

	return result, nil
}

// translateProductError 转换 Ent 错误为应用错误
func translateProductError(err error) error {
	if err == nil {
		return nil
	}
	if dbent.IsNotFound(err) {
		return apperrors.NotFound("PRODUCT_NOT_FOUND", "product not found")
	}
	return err
}

// rechargeProductEntityToService 从Ent实体转换为Service对象
func rechargeProductEntityToService(m *dbent.RechargeProduct) *service.RechargeProduct {
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
