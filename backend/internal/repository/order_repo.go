package repository

import (
	"context"
	"errors"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	dborder "github.com/Wei-Shaw/sub2api/ent/order"
	dbuser "github.com/Wei-Shaw/sub2api/ent/user"
	apperrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"entgo.io/ent/dialect/sql"
)

type orderRepository struct {
	client *dbent.Client
}

// NewOrderRepository 创建订单Repository
func NewOrderRepository(client *dbent.Client) service.OrderRepository {
	return &orderRepository{client: client}
}

// Create 创建订单
func (r *orderRepository) Create(ctx context.Context, order *service.Order) error {
	builder := r.client.Order.Create().
		SetOrderNo(order.OrderNo).
		SetUserID(order.UserID).
		SetProductName(order.ProductName).
		SetAmount(order.Amount).
		SetBonusAmount(order.BonusAmount).
		SetActualAmount(order.ActualAmount).
		SetPaymentMethod(order.PaymentMethod).
		SetPaymentGateway(order.PaymentGateway).
		SetStatus(order.Status).
		SetCreatedAt(order.CreatedAt).
		SetExpiredAt(order.ExpiredAt).
		SetNotes(order.Notes)

	if order.ProductID != nil {
		builder.SetProductID(*order.ProductID)
	}
	if order.TradeNo != "" {
		builder.SetTradeNo(order.TradeNo)
	}
	if order.CallbackData != "" {
		builder.SetCallbackData(order.CallbackData)
	}
	if order.PaidAt != nil {
		builder.SetPaidAt(*order.PaidAt)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		return err
	}
	order.ID = created.ID
	return nil
}

// GetByID 根据ID获取订单
func (r *orderRepository) GetByID(ctx context.Context, id int64) (*service.Order, error) {
	m, err := r.client.Order.Query().
		Where(dborder.IDEQ(id)).
		WithUser().
		WithProduct().
		Only(ctx)
	if err != nil {
		return nil, translateOrderError(err)
	}
	return orderEntityToService(m), nil
}

// GetByOrderNo 根据订单号获取订单
func (r *orderRepository) GetByOrderNo(ctx context.Context, orderNo string) (*service.Order, error) {
	m, err := r.client.Order.Query().
		Where(dborder.OrderNoEQ(orderNo)).
		WithUser().
		WithProduct().
		Only(ctx)
	if err != nil {
		return nil, translateOrderError(err)
	}
	return orderEntityToService(m), nil
}

// Update 更新订单
func (r *orderRepository) Update(ctx context.Context, order *service.Order) error {
	builder := r.client.Order.UpdateOneID(order.ID).
		SetOrderNo(order.OrderNo).
		SetUserID(order.UserID).
		SetProductName(order.ProductName).
		SetAmount(order.Amount).
		SetBonusAmount(order.BonusAmount).
		SetActualAmount(order.ActualAmount).
		SetPaymentMethod(order.PaymentMethod).
		SetPaymentGateway(order.PaymentGateway).
		SetStatus(order.Status).
		SetExpiredAt(order.ExpiredAt).
		SetNotes(order.Notes)

	if order.ProductID != nil {
		builder.SetProductID(*order.ProductID)
	} else {
		builder.ClearProductID()
	}
	if order.TradeNo != "" {
		builder.SetTradeNo(order.TradeNo)
	} else {
		builder.ClearTradeNo()
	}
	if order.CallbackData != "" {
		builder.SetCallbackData(order.CallbackData)
	} else {
		builder.ClearCallbackData()
	}
	if order.PaidAt != nil {
		builder.SetPaidAt(*order.PaidAt)
	} else {
		builder.ClearPaidAt()
	}

	_, err := builder.Save(ctx)
	return err
}

// ListByUser 获取用户订单列表
func (r *orderRepository) ListByUser(ctx context.Context, userID int64, page, limit int) ([]*service.Order, int64, error) {
	// 查询总数
	total, err := r.client.Order.Query().
		Where(dborder.UserIDEQ(userID)).
		Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * limit
	orders, err := r.client.Order.Query().
		Where(dborder.UserIDEQ(userID)).
		WithProduct().
		Order(dborder.ByCreatedAt(sql.OrderDesc())).
		Offset(offset).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*service.Order, len(orders))
	for i, m := range orders {
		result[i] = orderEntityToService(m)
	}

	return result, int64(total), nil
}

// AdminList 管理员查询订单列表（支持灵活筛选）
func (r *orderRepository) AdminList(ctx context.Context, filter *service.AdminListOrdersFilter) ([]*service.Order, int64, error) {
	query := r.client.Order.Query().WithUser().WithProduct()

	// 应用过滤条件
	if filter.UserID != nil {
		query = query.Where(dborder.UserIDEQ(*filter.UserID))
	}
	if filter.Status != nil {
		query = query.Where(dborder.StatusEQ(*filter.Status))
	}
	if filter.PaymentMethod != nil {
		query = query.Where(dborder.PaymentMethodEQ(*filter.PaymentMethod))
	}
	if filter.OrderNo != nil {
		query = query.Where(dborder.OrderNoContains(*filter.OrderNo))
	}
	if filter.StartDate != nil {
		query = query.Where(dborder.CreatedAtGTE(*filter.StartDate))
	}
	if filter.EndDate != nil {
		query = query.Where(dborder.CreatedAtLTE(*filter.EndDate))
	}

	// 查询总数
	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (filter.Page - 1) * filter.Limit
	orders, err := query.
		Order(dborder.ByCreatedAt(sql.OrderDesc())).
		Offset(offset).
		Limit(filter.Limit).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*service.Order, len(orders))
	for i, m := range orders {
		result[i] = orderEntityToService(m)
	}

	return result, int64(total), nil
}

// GetStatistics 获取订单统计信息
func (r *orderRepository) GetStatistics(ctx context.Context) (*service.OrderStatistics, error) {
	stats := &service.OrderStatistics{}

	// 总订单数
	total, err := r.client.Order.Query().Count(ctx)
	if err != nil {
		return nil, err
	}
	stats.TotalOrders = int64(total)

	// 各状态订单数
	pendingCount, _ := r.client.Order.Query().Where(dborder.StatusEQ(service.OrderStatusPending)).Count(ctx)
	paidCount, _ := r.client.Order.Query().Where(dborder.StatusEQ(service.OrderStatusPaid)).Count(ctx)
	failedCount, _ := r.client.Order.Query().Where(dborder.StatusEQ(service.OrderStatusFailed)).Count(ctx)
	expiredCount, _ := r.client.Order.Query().Where(dborder.StatusEQ(service.OrderStatusExpired)).Count(ctx)

	stats.PendingOrders = int64(pendingCount)
	stats.PaidOrders = int64(paidCount)
	stats.FailedOrders = int64(failedCount)
	stats.ExpiredOrders = int64(expiredCount)

	// 总成交金额和余额
	paidOrders, err := r.client.Order.Query().
		Where(dborder.StatusEQ(service.OrderStatusPaid)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	for _, o := range paidOrders {
		stats.TotalAmount += o.Amount
		stats.TotalBalance += o.ActualAmount
	}

	return stats, nil
}

// CountPendingByUser 统计用户未支付订单数量
func (r *orderRepository) CountPendingByUser(ctx context.Context, userID int64) (int64, error) {
	count, err := r.client.Order.Query().
		Where(
			dborder.UserIDEQ(userID),
			dborder.StatusEQ(service.OrderStatusPending),
		).
		Count(ctx)
	return int64(count), err
}

// MarkExpiredOrders 将过期的pending订单标记为expired
func (r *orderRepository) MarkExpiredOrders(ctx context.Context) (int64, error) {
	affected, err := r.client.Order.Update().
		Where(
			dborder.StatusEQ(service.OrderStatusPending),
			dborder.ExpiredAtLT(time.Now()),
		).
		SetStatus(service.OrderStatusExpired).
		Save(ctx)
	return int64(affected), err
}

// UpdateOrderAndUserBalance 原子更新订单状态和用户余额
func (r *orderRepository) UpdateOrderAndUserBalance(ctx context.Context, order *service.Order, userID int64, balanceDelta float64) error {
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return err
	}

	// 确保事务正确处理
	defer func() {
		if v := recover(); v != nil {
			_ = tx.Rollback()
			panic(v)
		}
	}()

	// 更新订单
	orderBuilder := tx.Order.UpdateOneID(order.ID).
		SetStatus(order.Status).
		SetNotes(order.Notes)

	if order.TradeNo != "" {
		orderBuilder.SetTradeNo(order.TradeNo)
	}
	if order.CallbackData != "" {
		orderBuilder.SetCallbackData(order.CallbackData)
	}
	if order.PaidAt != nil {
		orderBuilder.SetPaidAt(*order.PaidAt)
	}

	if _, err := orderBuilder.Save(ctx); err != nil {
		_ = tx.Rollback()
		return err
	}

	// 更新用户余额
	if balanceDelta != 0 {
		// 获取当前余额
		user, err := tx.User.Query().Where(dbuser.IDEQ(userID)).Only(ctx)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		newBalance := user.Balance + balanceDelta
		if _, err := tx.User.UpdateOneID(userID).SetBalance(newBalance).Save(ctx); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// translateOrderError 转换 Ent 错误为应用错误
func translateOrderError(err error) error {
	if err == nil {
		return nil
	}
	if dbent.IsNotFound(err) {
		return apperrors.NotFound("ORDER_NOT_FOUND", "order not found")
	}
	return err
}

// orderEntityToService 从Ent实体转换为Service对象
func orderEntityToService(m *dbent.Order) *service.Order {
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
		Status:         m.Status,
		CreatedAt:      m.CreatedAt,
		PaidAt:         m.PaidAt,
		ExpiredAt:      m.ExpiredAt,
		Notes:          m.Notes,
	}

	if m.TradeNo != nil {
		order.TradeNo = *m.TradeNo
	}
	if m.CallbackData != nil {
		order.CallbackData = *m.CallbackData
	}

	// 关联对象
	if m.Edges.User != nil {
		order.User = userEntityToService(m.Edges.User)
	}
	if m.Edges.Product != nil {
		order.Product = rechargeProductEntityToService(m.Edges.Product)
	}

	return order
}

// ErrOrderNotFound 订单未找到错误（供外部使用）
var ErrOrderNotFound = errors.New("order not found")
