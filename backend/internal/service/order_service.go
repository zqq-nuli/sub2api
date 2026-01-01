package service

import (
	"context"
	"fmt"
	"log"
	"time"
)

// OrderService 订单服务
type OrderService struct {
	orderRepo           OrderRepository
	userRepo            UserRepository
	billingCacheService *BillingCacheService
}

// NewOrderService 创建订单服务
func NewOrderService(
	orderRepo OrderRepository,
	userRepo UserRepository,
	billingCacheService *BillingCacheService,
) *OrderService {
	return &OrderService{
		orderRepo:           orderRepo,
		userRepo:            userRepo,
		billingCacheService: billingCacheService,
	}
}

// GetOrderByNo 根据订单号查询订单（用户只能查询自己的）
func (s *OrderService) GetOrderByNo(ctx context.Context, userID int64, orderNo string) (*Order, error) {
	order, err := s.orderRepo.GetByOrderNo(ctx, orderNo)
	if err != nil {
		return nil, fmt.Errorf("get order: %w", err)
	}

	// 验证订单归属
	if order.UserID != userID {
		return nil, ErrOrderNotFound
	}

	return order, nil
}

// GetUserOrdersInput 获取用户订单输入
type GetUserOrdersInput struct {
	UserID int64
	Page   int
	Limit  int
}

// GetUserOrdersOutput 获取用户订单输出
type GetUserOrdersOutput struct {
	Orders []*Order
	Total  int64
}

// GetUserOrders 获取用户订单列表（分页）
func (s *OrderService) GetUserOrders(ctx context.Context, input *GetUserOrdersInput) (*GetUserOrdersOutput, error) {
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.Limit <= 0 {
		input.Limit = 20
	}
	if input.Limit > 100 {
		input.Limit = 100
	}

	orders, total, err := s.orderRepo.ListByUser(ctx, input.UserID, input.Page, input.Limit)
	if err != nil {
		return nil, fmt.Errorf("list user orders: %w", err)
	}

	return &GetUserOrdersOutput{
		Orders: orders,
		Total:  total,
	}, nil
}

// AdminListOrdersInput 管理员查询订单输入
type AdminListOrdersInput struct {
	UserID        *int64  // 可选：按用户ID筛选
	Status        *string // 可选：按状态筛选
	PaymentMethod *string // 可选：按支付方式筛选
	OrderNo       *string // 可选：按订单号搜索
	StartDate     *time.Time
	EndDate       *time.Time
	Page          int
	Limit         int
}

// AdminListOrdersOutput 管理员查询订单输出
type AdminListOrdersOutput struct {
	Orders []*Order
	Total  int64
}

// AdminListOrders 管理员查询所有订单（支持筛选和分页）
func (s *OrderService) AdminListOrders(ctx context.Context, input *AdminListOrdersInput) (*AdminListOrdersOutput, error) {
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.Limit <= 0 {
		input.Limit = 20
	}
	if input.Limit > 100 {
		input.Limit = 100
	}

	orders, total, err := s.orderRepo.AdminList(ctx, &AdminListOrdersFilter{
		UserID:        input.UserID,
		Status:        input.Status,
		PaymentMethod: input.PaymentMethod,
		OrderNo:       input.OrderNo,
		StartDate:     input.StartDate,
		EndDate:       input.EndDate,
		Page:          input.Page,
		Limit:         input.Limit,
	})
	if err != nil {
		return nil, fmt.Errorf("admin list orders: %w", err)
	}

	return &AdminListOrdersOutput{
		Orders: orders,
		Total:  total,
	}, nil
}

// AdminGetOrderByID 管理员根据ID获取订单详情
func (s *OrderService) AdminGetOrderByID(ctx context.Context, orderID int64) (*Order, error) {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("get order by id: %w", err)
	}
	return order, nil
}

// AdminUpdateOrderStatusInput 管理员更新订单状态输入
type AdminUpdateOrderStatusInput struct {
	OrderID int64
	Status  string
	Notes   string
}

// AdminUpdateOrderStatus 管理员手动更新订单状态
func (s *OrderService) AdminUpdateOrderStatus(ctx context.Context, input *AdminUpdateOrderStatusInput) error {
	// 1. 获取订单
	order, err := s.orderRepo.GetByID(ctx, input.OrderID)
	if err != nil {
		return fmt.Errorf("get order: %w", err)
	}

	// 2. 验证状态有效性
	if !IsValidOrderStatus(input.Status) {
		return ErrOrderStatusInvalid
	}

	// 记录旧状态，用于判断是否需要更新余额
	oldStatus := order.Status
	newStatus := input.Status

	// 3. 准备更新订单
	order.Status = newStatus
	if input.Notes != "" {
		if order.Notes != "" {
			order.Notes += "\n[Admin] " + input.Notes
		} else {
			order.Notes = "[Admin] " + input.Notes
		}
	}

	// 4. 检查是否需要更新用户余额
	// 情况1: 从非paid变为paid -> 增加余额
	// 情况2: 从paid变为非paid(如refunded) -> 扣减余额
	needAddBalance := oldStatus != OrderStatusPaid && newStatus == OrderStatusPaid
	needDeductBalance := oldStatus == OrderStatusPaid && newStatus == OrderStatusRefunded

	if needAddBalance || needDeductBalance {
		// 更新订单状态
		if newStatus == OrderStatusPaid {
			now := time.Now()
			order.PaidAt = &now
		}

		// 计算余额变化
		var balanceDelta float64
		if needAddBalance {
			balanceDelta = order.ActualAmount
		} else if needDeductBalance {
			balanceDelta = -order.ActualAmount
		}

		// 原子更新订单和用户余额
		if err := s.orderRepo.UpdateOrderAndUserBalance(ctx, order, order.UserID, balanceDelta); err != nil {
			log.Printf("[Order] Admin update: failed to update order and balance: %v", err)
			return fmt.Errorf("update order and balance: %w", err)
		}

		if needAddBalance {
			log.Printf("[Order] Admin confirmed order %s paid, user %d balance +%.2f USD", order.OrderNo, order.UserID, order.ActualAmount)
		} else if needDeductBalance {
			log.Printf("[Order] Admin refunded order %s, user %d balance -%.2f USD", order.OrderNo, order.UserID, order.ActualAmount)
		}

		// 异步失效余额缓存
		go func() {
			cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := s.billingCacheService.InvalidateUserBalance(cacheCtx, order.UserID); err != nil {
				log.Printf("[Order] Failed to invalidate cache for user %d: %v", order.UserID, err)
			}
		}()

		return nil
	}

	// 不涉及余额变更的状态更新，直接更新
	if err := s.orderRepo.Update(ctx, order); err != nil {
		return fmt.Errorf("update order: %w", err)
	}

	return nil
}

// GetOrderStatistics 获取订单统计信息（管理员）
func (s *OrderService) GetOrderStatistics(ctx context.Context) (*OrderStatistics, error) {
	stats, err := s.orderRepo.GetStatistics(ctx)
	if err != nil {
		return nil, fmt.Errorf("get order statistics: %w", err)
	}
	return stats, nil
}

// IsValidOrderStatus 验证订单状态是否有效
func IsValidOrderStatus(status string) bool {
	validStatuses := map[string]bool{
		OrderStatusPending:  true,
		OrderStatusPaid:     true,
		OrderStatusFailed:   true,
		OrderStatusExpired:  true,
		OrderStatusRefunded: true,
	}
	return validStatuses[status]
}
