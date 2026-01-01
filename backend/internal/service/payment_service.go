package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
)

// PaymentService 支付服务
type PaymentService struct {
	orderRepo           OrderRepository
	userRepo            UserRepository
	productRepo         RechargeProductRepository
	settingService      *SettingService
	billingCacheService *BillingCacheService
	redisClient         *redis.Client
}

// NewPaymentService 创建支付服务
func NewPaymentService(
	orderRepo OrderRepository,
	userRepo UserRepository,
	productRepo RechargeProductRepository,
	settingService *SettingService,
	billingCacheService *BillingCacheService,
	redisClient *redis.Client,
) *PaymentService {
	return &PaymentService{
		orderRepo:           orderRepo,
		userRepo:            userRepo,
		productRepo:         productRepo,
		settingService:      settingService,
		billingCacheService: billingCacheService,
		redisClient:         redisClient,
	}
}

// CreateOrderInput 创建订单输入
type CreateOrderInput struct {
	UserID        int64
	ProductID     int64
	PaymentMethod string
	ReturnURL     string
}

// CreateOrderOutput 创建订单输出
type CreateOrderOutput struct {
	Order         *Order
	PaymentURL    string            // 表单提交URL
	PaymentParams map[string]string // 表单参数
}

// CreateOrder 创建充值订单
func (s *PaymentService) CreateOrder(ctx context.Context, input *CreateOrderInput) (*CreateOrderOutput, error) {
	// 1. 验证用户
	user, err := s.userRepo.GetByID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if user.Status != StatusActive {
		return nil, fmt.Errorf("user is not active")
	}

	// 2. 检查用户未支付订单数量（限制最多3个未支付订单）
	pendingCount, err := s.orderRepo.CountPendingByUser(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("count pending orders: %w", err)
	}
	if pendingCount >= 3 {
		return nil, ErrTooManyPendingOrders
	}

	// 3. 验证套餐
	product, err := s.productRepo.GetByID(ctx, input.ProductID)
	if err != nil {
		return nil, fmt.Errorf("get product: %w", err)
	}
	if !product.IsActive() {
		return nil, ErrProductInactive
	}

	// 4. 验证支付方式
	if !IsValidPaymentMethod(input.PaymentMethod) {
		return nil, ErrInvalidPaymentMethod
	}

	// 5. 生成唯一订单号
	orderNo := GenerateOrderNo()

	// 6. 创建订单
	actualAmount := product.GetTotalBalance()
	order := &Order{
		OrderNo:        orderNo,
		UserID:         input.UserID,
		ProductID:      &product.ID,
		ProductName:    product.Name,
		Amount:         product.Amount,
		BonusAmount:    product.BonusBalance,
		ActualAmount:   actualAmount,
		PaymentMethod:  input.PaymentMethod,
		PaymentGateway: "epay",
		Status:         OrderStatusPending,
		CreatedAt:      time.Now(),
		ExpiredAt:      time.Now().Add(15 * time.Minute),
	}

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("create order: %w", err)
	}

	// 7. 获取易支付配置
	epayConfig, err := s.getEpayConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("get epay config: %w", err)
	}

	// 8. 构建易支付表单数据
	paymentURL, paymentParams, err := s.buildEpayFormData(epayConfig, order, input.ReturnURL)
	if err != nil {
		return nil, fmt.Errorf("build payment form: %w", err)
	}

	return &CreateOrderOutput{
		Order:         order,
		PaymentURL:    paymentURL,
		PaymentParams: paymentParams,
	}, nil
}

// HandlePaymentCallback 处理易支付异步回调
func (s *PaymentService) HandlePaymentCallback(ctx context.Context, params map[string]string) error {
	// 1. 验证签名
	epayConfig, err := s.getEpayConfig(ctx)
	if err != nil {
		return fmt.Errorf("get epay config: %w", err)
	}

	epayClient := NewEpayClient(epayConfig.ApiURL, epayConfig.MerchantID, epayConfig.MerchantKey)
	if !epayClient.VerifySign(params) {
		log.Printf("[Payment] Invalid signature: %+v", params)
		return ErrInvalidSign
	}

	// 2. 提取订单号
	orderNo := params["out_trade_no"]
	if orderNo == "" {
		return ErrMissingOrderNo
	}

	// 3. 获取分布式锁（防止并发回调）
	lockKey := "payment:callback:lock:" + orderNo
	locked, err := s.redisClient.SetNX(ctx, lockKey, "1", 5*time.Minute).Result()
	if err != nil {
		return fmt.Errorf("acquire lock: %w", err)
	}
	if !locked {
		log.Printf("[Payment] Order %s is locked by another process", orderNo)
		return ErrOrderLocked
	}
	defer func() {
		_ = s.redisClient.Del(ctx, lockKey).Err()
	}()

	// 4. 查询订单
	order, err := s.orderRepo.GetByOrderNo(ctx, orderNo)
	if err != nil {
		return fmt.Errorf("get order: %w", err)
	}

	// 5. 检查订单状态（幂等性保证）
	if order.IsPaid() {
		log.Printf("[Payment] Order %s already paid, skip", orderNo)
		return nil // 已支付，直接返回成功
	}

	if !order.IsPending() {
		log.Printf("[Payment] Order %s status is %s, not pending", orderNo, order.Status)
		return ErrOrderStatusInvalid
	}

	// 6. 检查订单是否已过期
	if time.Now().After(order.ExpiredAt) {
		log.Printf("[Payment] Order %s has expired at %v", orderNo, order.ExpiredAt)
		// 更新订单状态为过期
		order.Status = OrderStatusExpired
		callbackJSON, _ := json.Marshal(params)
		order.CallbackData = string(callbackJSON)
		_ = s.orderRepo.Update(ctx, order)
		return ErrOrderExpired
	}

	// 7. 验证金额（使用 decimal 精确比较，避免浮点数精度问题）
	paidAmountDec, err := decimal.NewFromString(params["money"])
	if err != nil {
		return fmt.Errorf("parse paid amount: %w", err)
	}
	expectedAmountDec := decimal.NewFromFloat(order.Amount)

	// 使用 decimal 精确比较，不允许任何误差
	if !paidAmountDec.Equal(expectedAmountDec) {
		log.Printf("[Payment] Amount mismatch: expected %s, got %s", expectedAmountDec.String(), paidAmountDec.String())
		return ErrAmountMismatch
	}

	// 8. 验证支付状态
	tradeStatus := params["trade_status"]
	if tradeStatus != "TRADE_SUCCESS" {
		// 更新订单为失败状态
		order.Status = OrderStatusFailed
		callbackJSON, _ := json.Marshal(params)
		order.CallbackData = string(callbackJSON)
		_ = s.orderRepo.Update(ctx, order)
		log.Printf("[Payment] Payment failed: %s", tradeStatus)
		return ErrPaymentFailed
	}

	// 9. 更新订单状态和用户余额（原子操作）
	now := time.Now()
	order.Status = OrderStatusPaid
	order.TradeNo = params["trade_no"]
	order.PaidAt = &now
	callbackJSON, _ := json.Marshal(params)
	order.CallbackData = string(callbackJSON)

	if err := s.orderRepo.UpdateOrderAndUserBalance(ctx, order, order.UserID, order.ActualAmount); err != nil {
		log.Printf("[Payment] CRITICAL: Order %s update failed: %v", orderNo, err)
		return fmt.Errorf("update order and balance: %w", err)
	}

	log.Printf("[Payment] Order %s paid successfully, user %d balance +%.2f USD", orderNo, order.UserID, order.ActualAmount)

	// 12. 异步失效余额缓存
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.billingCacheService.InvalidateUserBalance(cacheCtx, order.UserID); err != nil {
			log.Printf("[Payment] Failed to invalidate cache for user %d: %v", order.UserID, err)
		}
	}()

	return nil
}

// EpayConfig 易支付配置
type EpayConfig struct {
	Enabled     bool
	ApiURL      string
	MerchantID  string
	MerchantKey string
	NotifyURL   string
	ReturnURL   string
}

// getEpayConfig 获取易支付配置
func (s *PaymentService) getEpayConfig(ctx context.Context) (*EpayConfig, error) {
	settings, err := s.settingService.GetAllSettings(ctx)
	if err != nil {
		return nil, err
	}

	if !settings.EpayEnabled {
		return nil, ErrPaymentDisabled
	}

	if settings.EpayApiURL == "" || settings.EpayMerchantID == "" || settings.EpayMerchantKey == "" {
		return nil, fmt.Errorf("epay config incomplete")
	}

	return &EpayConfig{
		Enabled:     settings.EpayEnabled,
		ApiURL:      settings.EpayApiURL,
		MerchantID:  settings.EpayMerchantID,
		MerchantKey: settings.EpayMerchantKey,
		NotifyURL:   settings.EpayNotifyURL,
		ReturnURL:   settings.EpayReturnURL,
	}, nil
}

// buildEpayFormData 构建易支付表单数据（用于POST提交）
func (s *PaymentService) buildEpayFormData(config *EpayConfig, order *Order, returnURL string) (string, map[string]string, error) {
	epayClient := NewEpayClient(config.ApiURL, config.MerchantID, config.MerchantKey)

	// 根据支付渠道获取对应的 epay_type
	epayType, err := s.settingService.GetEpayTypeByChannel(context.Background(), order.PaymentMethod)
	if err != nil {
		// 回退到默认值 epay
		epayType = "epay"
	}

	params := map[string]string{
		"type":         epayType,
		"out_trade_no": order.OrderNo,
		"notify_url":   config.NotifyURL,
		"return_url":   returnURL,
		"name":         order.ProductName,
		"money":        fmt.Sprintf("%.2f", order.Amount),
	}

	// 获取表单提交URL和带签名的参数
	return epayClient.CreatePaymentForm(params)
}

// SyncOrderFromUpstream 从上游同步订单状态
// 用于手动触发或轮询时主动查询上游支付状态
func (s *PaymentService) SyncOrderFromUpstream(ctx context.Context, orderNo string, tradeNo string) (*Order, error) {
	// 1. 获取订单
	order, err := s.orderRepo.GetByOrderNo(ctx, orderNo)
	if err != nil {
		return nil, fmt.Errorf("get order: %w", err)
	}

	// 如果订单已经完成，直接返回
	if order.IsPaid() || order.Status == OrderStatusFailed || order.Status == OrderStatusExpired {
		return order, nil
	}

	// 检查订单是否已过期
	if time.Now().After(order.ExpiredAt) {
		order.Status = OrderStatusExpired
		_ = s.orderRepo.Update(ctx, order)
		return order, nil
	}

	// 2. 获取易支付配置
	epayConfig, err := s.getEpayConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("get epay config: %w", err)
	}

	// 如果没有提供 tradeNo，使用订单已有的 tradeNo（如果有）
	if tradeNo == "" && order.TradeNo != "" {
		tradeNo = order.TradeNo
	}

	// 3. 查询上游订单状态
	// 注意：易支付查询API允许只传 out_trade_no（商户订单号）来查询
	epayClient := NewEpayClient(epayConfig.ApiURL, epayConfig.MerchantID, epayConfig.MerchantKey)
	result, err := epayClient.QueryOrderStatus(tradeNo, orderNo)
	if err != nil {
		log.Printf("[Payment] Upstream query failed for order %s: %v", orderNo, err)
		// 查询失败不报错，返回当前订单状态
		return order, nil
	}

	log.Printf("[Payment] Upstream query result for order %s: %+v", orderNo, result)

	// 4. 根据上游状态更新本地订单
	if result.Code == 1 && result.Status == 1 {
		// 上游订单已完成，更新本地订单
		return s.completeOrderFromUpstream(ctx, order, result)
	}

	// 上游订单未完成或查询失败
	return order, nil
}

// completeOrderFromUpstream 根据上游查询结果完成订单
func (s *PaymentService) completeOrderFromUpstream(ctx context.Context, order *Order, upstreamResult *EpayOrderQueryResponse) (*Order, error) {
	// 获取分布式锁
	lockKey := "payment:sync:lock:" + order.OrderNo
	locked, err := s.redisClient.SetNX(ctx, lockKey, "1", 5*time.Minute).Result()
	if err != nil {
		return nil, fmt.Errorf("acquire lock: %w", err)
	}
	if !locked {
		return nil, ErrOrderLocked
	}
	defer func() {
		_ = s.redisClient.Del(ctx, lockKey).Err()
	}()

	// 重新获取订单（确保最新状态）
	order, err = s.orderRepo.GetByOrderNo(ctx, order.OrderNo)
	if err != nil {
		return nil, fmt.Errorf("get order: %w", err)
	}

	// 再次检查订单状态
	if order.IsPaid() {
		return order, nil
	}

	// 更新订单状态和用户余额（原子操作）
	now := time.Now()
	order.Status = OrderStatusPaid
	order.TradeNo = upstreamResult.TradeNo
	order.PaidAt = &now
	order.Notes = "synced from upstream query"

	if err := s.orderRepo.UpdateOrderAndUserBalance(ctx, order, order.UserID, order.ActualAmount); err != nil {
		log.Printf("[Payment] CRITICAL: Order %s update failed: %v", order.OrderNo, err)
		return nil, fmt.Errorf("update order and balance: %w", err)
	}

	log.Printf("[Payment] Order %s synced from upstream, user %d balance +%.2f USD", order.OrderNo, order.UserID, order.ActualAmount)

	// 异步失效余额缓存
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.billingCacheService.InvalidateUserBalance(cacheCtx, order.UserID); err != nil {
			log.Printf("[Payment] Failed to invalidate cache for user %d: %v", order.UserID, err)
		}
	}()

	return order, nil
}
