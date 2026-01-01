package service

import (
	"errors"
)

// 订单相关错误
var (
	ErrOrderNotFound       = errors.New("order not found")
	ErrOrderStatusInvalid  = errors.New("invalid order status")
	ErrOrderLocked         = errors.New("order is locked by another process")
	ErrOrderExpired        = errors.New("order has expired")
	ErrAmountMismatch      = errors.New("payment amount mismatch")
	ErrPaymentFailed       = errors.New("payment failed")
	ErrInvalidSign         = errors.New("invalid signature")
	ErrMissingOrderNo      = errors.New("missing order number")
	ErrPaymentDisabled     = errors.New("payment is disabled")
	ErrInvalidPaymentMethod = errors.New("invalid payment method")
	ErrTooManyPendingOrders = errors.New("too many pending orders")
)

// 充值套餐相关错误
var (
	ErrProductNotFound  = errors.New("recharge product not found")
	ErrProductInactive  = errors.New("recharge product is inactive")
)
