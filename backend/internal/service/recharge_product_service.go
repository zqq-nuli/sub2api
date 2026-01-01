package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// RechargeProductService 充值套餐服务
type RechargeProductService struct {
	productRepo RechargeProductRepository
	redisClient *redis.Client
}

// NewRechargeProductService 创建充值套餐服务
func NewRechargeProductService(
	productRepo RechargeProductRepository,
	redisClient *redis.Client,
) *RechargeProductService {
	return &RechargeProductService{
		productRepo: productRepo,
		redisClient: redisClient,
	}
}

const (
	cacheKeyActiveProducts = "recharge:products:active"
	productCacheTTL        = 5 * time.Minute
)

// GetActiveProducts 获取所有启用的充值套餐（带缓存）
func (s *RechargeProductService) GetActiveProducts(ctx context.Context) ([]*RechargeProduct, error) {
	// 1. 尝试从缓存获取
	products, err := s.getProductsFromCache(ctx)
	if err == nil && products != nil {
		return products, nil
	}

	// 2. 从数据库查询
	products, err = s.productRepo.ListActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("list active products: %w", err)
	}

	// 3. 异步设置缓存（避免阻塞主流程）
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.setProductsToCache(cacheCtx, products); err != nil {
			log.Printf("[RechargeProduct] Failed to cache active products: %v", err)
		}
	}()

	return products, nil
}

// GetProductByID 根据ID获取套餐
func (s *RechargeProductService) GetProductByID(ctx context.Context, id int64) (*RechargeProduct, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get product by id: %w", err)
	}
	return product, nil
}

// CreateProductInput 创建套餐输入
type CreateProductInput struct {
	Name          string
	Amount        float64
	Balance       float64
	BonusBalance  float64
	Description   string
	SortOrder     int
	IsHot         bool
	DiscountLabel string
}

// CreateProduct 创建充值套餐（管理员）
func (s *RechargeProductService) CreateProduct(ctx context.Context, input *CreateProductInput) (*RechargeProduct, error) {
	// 验证输入
	if input.Name == "" {
		return nil, fmt.Errorf("product name is required")
	}
	if input.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than 0")
	}
	if input.Balance <= 0 {
		return nil, fmt.Errorf("balance must be greater than 0")
	}

	product := &RechargeProduct{
		Name:          input.Name,
		Amount:        input.Amount,
		Balance:       input.Balance,
		BonusBalance:  input.BonusBalance,
		Description:   input.Description,
		SortOrder:     input.SortOrder,
		IsHot:         input.IsHot,
		DiscountLabel: input.DiscountLabel,
		Status:        ProductStatusActive,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("create product: %w", err)
	}

	// 失效缓存
	s.invalidateCache(ctx)

	return product, nil
}

// UpdateProductInput 更新套餐输入
type UpdateProductInput struct {
	ID            int64
	Name          string
	Amount        float64
	Balance       float64
	BonusBalance  float64
	Description   string
	SortOrder     int
	IsHot         bool
	DiscountLabel string
	Status        string
}

// UpdateProduct 更新充值套餐（管理员）
func (s *RechargeProductService) UpdateProduct(ctx context.Context, input *UpdateProductInput) (*RechargeProduct, error) {
	// 获取现有套餐
	product, err := s.productRepo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, fmt.Errorf("get product: %w", err)
	}

	// 验证输入
	if input.Name == "" {
		return nil, fmt.Errorf("product name is required")
	}
	if input.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than 0")
	}
	if input.Balance <= 0 {
		return nil, fmt.Errorf("balance must be greater than 0")
	}
	if !IsValidProductStatus(input.Status) {
		return nil, fmt.Errorf("invalid product status")
	}

	// 更新字段
	product.Name = input.Name
	product.Amount = input.Amount
	product.Balance = input.Balance
	product.BonusBalance = input.BonusBalance
	product.Description = input.Description
	product.SortOrder = input.SortOrder
	product.IsHot = input.IsHot
	product.DiscountLabel = input.DiscountLabel
	product.Status = input.Status
	product.UpdatedAt = time.Now()

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("update product: %w", err)
	}

	// 失效缓存
	s.invalidateCache(ctx)

	return product, nil
}

// DeleteProduct 删除充值套餐（管理员）
func (s *RechargeProductService) DeleteProduct(ctx context.Context, id int64) error {
	// 检查套餐是否存在
	_, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get product: %w", err)
	}

	// 删除套餐
	if err := s.productRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete product: %w", err)
	}

	// 失效缓存
	s.invalidateCache(ctx)

	return nil
}

// ListAllProducts 获取所有套餐（管理员）
func (s *RechargeProductService) ListAllProducts(ctx context.Context) ([]*RechargeProduct, error) {
	products, err := s.productRepo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list all products: %w", err)
	}
	return products, nil
}

// getProductsFromCache 从Redis缓存获取套餐列表
func (s *RechargeProductService) getProductsFromCache(ctx context.Context) ([]*RechargeProduct, error) {
	if s.redisClient == nil {
		return nil, fmt.Errorf("cache miss")
	}

	data, err := s.redisClient.Get(ctx, cacheKeyActiveProducts).Result()
	if err != nil {
		return nil, err
	}

	var products []*RechargeProduct
	if err := json.Unmarshal([]byte(data), &products); err != nil {
		return nil, err
	}
	return products, nil
}

// setProductsToCache 将套餐列表存入Redis缓存
func (s *RechargeProductService) setProductsToCache(ctx context.Context, products []*RechargeProduct) error {
	if s.redisClient == nil {
		return nil
	}

	data, err := json.Marshal(products)
	if err != nil {
		return err
	}

	return s.redisClient.Set(ctx, cacheKeyActiveProducts, string(data), productCacheTTL).Err()
}

// invalidateCache 失效缓存
func (s *RechargeProductService) invalidateCache(ctx context.Context) {
	if s.redisClient == nil {
		return
	}

	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := s.redisClient.Del(cacheCtx, cacheKeyActiveProducts).Err(); err != nil {
			log.Printf("[RechargeProduct] Failed to invalidate cache: %v", err)
		}
	}()
}

// IsValidProductStatus 验证套餐状态是否有效
func IsValidProductStatus(status string) bool {
	return status == ProductStatusActive || status == ProductStatusInactive
}
