package dto

import "time"

// OrderDTO 订单数据传输对象
type OrderDTO struct {
	ID             int64      `json:"id"`
	OrderNo        string     `json:"order_no"`
	UserID         int64      `json:"user_id"`
	UserEmail      string     `json:"user_email,omitempty"` // 管理员查询时返回
	ProductID      *int64     `json:"product_id"`
	ProductName    string     `json:"product_name"`
	Amount         float64    `json:"amount"`          // 人民币金额
	BonusAmount    float64    `json:"bonus_amount"`    // 赠送美元
	ActualAmount   float64    `json:"actual_amount"`   // 实际到账美元
	PaymentMethod  string     `json:"payment_method"`  // alipay/wxpay/epusdt
	PaymentGateway string     `json:"payment_gateway"` // epay
	TradeNo        string     `json:"trade_no,omitempty"`
	Status         string     `json:"status"` // pending/paid/failed/expired/refunded
	CreatedAt      time.Time  `json:"created_at"`
	PaidAt         *time.Time `json:"paid_at,omitempty"`
	ExpiredAt      time.Time  `json:"expired_at"`
	Notes          string     `json:"notes,omitempty"`
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	ProductID     int64  `json:"product_id" binding:"required"`
	PaymentMethod string `json:"payment_method" binding:"required,oneof=alipay wxpay qqpay epusdt"`
	ReturnURL     string `json:"return_url"` // 可选：同步回调URL
}

// CreateOrderResponse 创建订单响应
type CreateOrderResponse struct {
	Order         *OrderDTO         `json:"order"`
	PaymentURL    string            `json:"payment_url"`              // 表单提交URL
	PaymentParams map[string]string `json:"payment_params,omitempty"` // 表单参数（用于POST提交）
}

// GetUserOrdersRequest 获取用户订单请求（查询参数）
type GetUserOrdersRequest struct {
	Page  int `form:"page" binding:"omitempty,min=1"`
	Limit int `form:"limit" binding:"omitempty,min=1,max=100"`
}

// GetUserOrdersResponse 获取用户订单响应
type GetUserOrdersResponse struct {
	Orders []*OrderDTO `json:"orders"`
	Total  int64       `json:"total"`
	Page   int         `json:"page"`
	Limit  int         `json:"limit"`
}

// AdminListOrdersRequest 管理员查询订单请求
type AdminListOrdersRequest struct {
	UserID        *int64  `form:"user_id"`
	Status        *string `form:"status"`
	PaymentMethod *string `form:"payment_method"`
	OrderNo       *string `form:"order_no"`
	StartDate     *string `form:"start_date"` // YYYY-MM-DD
	EndDate       *string `form:"end_date"`   // YYYY-MM-DD
	Page          int     `form:"page" binding:"omitempty,min=1"`
	Limit         int     `form:"limit" binding:"omitempty,min=1,max=100"`
}

// AdminUpdateOrderStatusRequest 管理员更新订单状态请求
type AdminUpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending paid failed expired refunded"`
	Notes  string `json:"notes"`
}

// OrderStatisticsResponse 订单统计响应
type OrderStatisticsResponse struct {
	TotalOrders   int64   `json:"total_orders"`
	PendingOrders int64   `json:"pending_orders"`
	PaidOrders    int64   `json:"paid_orders"`
	FailedOrders  int64   `json:"failed_orders"`
	ExpiredOrders int64   `json:"expired_orders"`
	TotalAmount   float64 `json:"total_amount"`  // 总成交金额（CNY）
	TotalBalance  float64 `json:"total_balance"` // 总充值余额（USD）
}
