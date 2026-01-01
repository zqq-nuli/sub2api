package repository

import (
	"log"
	"time"

	"gorm.io/gorm"
)

// MaxExpiresAt is the maximum allowed expiration date for subscriptions (year 2099)
// This prevents time.Time JSON serialization errors (RFC 3339 requires year <= 9999)
var maxExpiresAt = time.Date(2099, 12, 31, 23, 59, 59, 0, time.UTC)

// AutoMigrate runs schema migrations for all repository persistence models.
// Persistence models are defined within individual `*_repo.go` files.
func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&userModel{},
		&apiKeyModel{},
		&groupModel{},
		&accountModel{},
		&accountGroupModel{},
		&proxyModel{},
		&redeemCodeModel{},
		&usageLogModel{},
		&settingModel{},
		&userSubscriptionModel{},
		&orderModel{},
		&rechargeProductModel{},
	)
	if err != nil {
		return err
	}

	// 创建默认分组(简易模式支持)
	if err := ensureDefaultGroups(db); err != nil {
		return err
	}

	// 创建默认充值套餐
	if err := ensureDefaultRechargeProducts(db); err != nil {
		return err
	}

	// 修复无效的过期时间（年份超过 2099 会导致 JSON 序列化失败）
	return fixInvalidExpiresAt(db)
}

// fixInvalidExpiresAt 修复 user_subscriptions 表中无效的过期时间
func fixInvalidExpiresAt(db *gorm.DB) error {
	result := db.Model(&userSubscriptionModel{}).
		Where("expires_at > ?", maxExpiresAt).
		Update("expires_at", maxExpiresAt)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected > 0 {
		log.Printf("[AutoMigrate] Fixed %d subscriptions with invalid expires_at (year > 2099)", result.RowsAffected)
	}
	return nil
}

// ensureDefaultGroups 确保默认分组存在(简易模式支持)
// 为每个平台创建一个默认分组,配置最大权限以确保简易模式下不受限制
func ensureDefaultGroups(db *gorm.DB) error {
	defaultGroups := []struct {
		name        string
		platform    string
		description string
	}{
		{
			name:        "anthropic-default",
			platform:    "anthropic",
			description: "Default group for Anthropic accounts (Simple Mode)",
		},
		{
			name:        "openai-default",
			platform:    "openai",
			description: "Default group for OpenAI accounts (Simple Mode)",
		},
		{
			name:        "gemini-default",
			platform:    "gemini",
			description: "Default group for Gemini accounts (Simple Mode)",
		},
	}

	for _, dg := range defaultGroups {
		var count int64
		if err := db.Model(&groupModel{}).Where("name = ?", dg.name).Count(&count).Error; err != nil {
			return err
		}

		if count == 0 {
			group := &groupModel{
				Name:             dg.name,
				Description:      dg.description,
				Platform:         dg.platform,
				RateMultiplier:   1.0,
				IsExclusive:      false,
				Status:           "active",
				SubscriptionType: "standard",
			}
			if err := db.Create(group).Error; err != nil {
				log.Printf("[AutoMigrate] Failed to create default group %s: %v", dg.name, err)
				return err
			}
			log.Printf("[AutoMigrate] Created default group: %s (platform: %s)", dg.name, dg.platform)
		}
	}

	return nil
}

// ensureDefaultRechargeProducts 确保默认充值套餐存在
func ensureDefaultRechargeProducts(db *gorm.DB) error {
	defaultProducts := []rechargeProductModel{
		{
			Name:          "10元充值",
			Amount:        10.00,
			Balance:       1.00,
			BonusBalance:  0.00,
			SortOrder:     1,
			IsHot:         false,
			DiscountLabel: "",
			Status:        "active",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			Name:          "30元充值",
			Amount:        30.00,
			Balance:       3.00,
			BonusBalance:  0.30,
			SortOrder:     2,
			IsHot:         false,
			DiscountLabel: "赠送10%",
			Status:        "active",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			Name:          "50元充值",
			Amount:        50.00,
			Balance:       5.00,
			BonusBalance:  1.00,
			SortOrder:     3,
			IsHot:         true,
			DiscountLabel: "赠送20%",
			Status:        "active",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			Name:          "100元充值",
			Amount:        100.00,
			Balance:       10.00,
			BonusBalance:  2.50,
			SortOrder:     4,
			IsHot:         true,
			DiscountLabel: "赠送25%",
			Status:        "active",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			Name:          "200元充值",
			Amount:        200.00,
			Balance:       20.00,
			BonusBalance:  6.00,
			SortOrder:     5,
			IsHot:         false,
			DiscountLabel: "赠送30%",
			Status:        "active",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}

	// 检查是否已有套餐
	var count int64
	if err := db.Model(&rechargeProductModel{}).Count(&count).Error; err != nil {
		return err
	}

	// 如果没有套餐，创建默认套餐
	if count == 0 {
		for _, product := range defaultProducts {
			if err := db.Create(&product).Error; err != nil {
				log.Printf("[AutoMigrate] Failed to create default recharge product %s: %v", product.Name, err)
				return err
			}
			log.Printf("[AutoMigrate] Created default recharge product: %s (Amount: %.2f CNY, Balance: %.2f USD)",
				product.Name, product.Amount, product.Balance+product.BonusBalance)
		}
	}

	return nil
}
