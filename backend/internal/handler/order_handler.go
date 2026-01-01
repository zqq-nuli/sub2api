package handler

import (
	"net/http"

	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	"github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

// OrderHandler 订单处理器（用户侧）
type OrderHandler struct {
	orderService           *service.OrderService
	paymentService         *service.PaymentService
	rechargeProductService *service.RechargeProductService
}

// NewOrderHandler 创建订单处理器
func NewOrderHandler(
	orderService *service.OrderService,
	paymentService *service.PaymentService,
	rechargeProductService *service.RechargeProductService,
) *OrderHandler {
	return &OrderHandler{
		orderService:           orderService,
		paymentService:         paymentService,
		rechargeProductService: rechargeProductService,
	}
}

// GetRechargeProducts 获取充值套餐列表
// GET /api/v1/recharge/products
func (h *OrderHandler) GetRechargeProducts(c *gin.Context) {
	products, err := h.rechargeProductService.GetActiveProducts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.InternalServer("", "Failed to get recharge products"))
		return
	}

	// 转换为DTO
	productDTOs := make([]*dto.RechargeProductDTO, len(products))
	for i, p := range products {
		productDTOs[i] = rechargeProductToDTO(p)
	}

	response.Success(c, productDTOs)
}

// CreateOrder 创建充值订单
// POST /api/v1/orders
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	subject, ok := middleware.GetAuthSubjectFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, errors.Unauthorized("", "User not authenticated"))
		return
	}

	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.BadRequest("", err.Error()))
		return
	}

	// 如果没有提供 ReturnURL，使用默认值
	if req.ReturnURL == "" {
		req.ReturnURL = c.Request.Header.Get("Origin") + "/recharge/result"
	}

	// 创建订单
	result, err := h.paymentService.CreateOrder(c.Request.Context(), &service.CreateOrderInput{
		UserID:        subject.UserID,
		ProductID:     req.ProductID,
		PaymentMethod: req.PaymentMethod,
		ReturnURL:     req.ReturnURL,
	})
	if err != nil {
		if err == service.ErrProductInactive {
			c.JSON(http.StatusBadRequest, errors.BadRequest("", "Recharge product is inactive"))
			return
		}
		if err == service.ErrInvalidPaymentMethod {
			c.JSON(http.StatusBadRequest, errors.BadRequest("", "Invalid payment method"))
			return
		}
		if err == service.ErrPaymentDisabled {
			c.JSON(http.StatusBadRequest, errors.BadRequest("", "Payment is disabled"))
			return
		}
		if err == service.ErrTooManyPendingOrders {
			c.JSON(http.StatusBadRequest, errors.BadRequest("", "Too many pending orders, please complete or cancel existing orders first"))
			return
		}

		c.JSON(http.StatusInternalServerError, errors.InternalServer("", "Failed to create order"))
		return
	}

	// 返回订单和支付表单数据
	response.Success(c, dto.CreateOrderResponse{
		Order:         orderToDTO(result.Order),
		PaymentURL:    result.PaymentURL,
		PaymentParams: result.PaymentParams,
	})
}

// GetOrderByNo 根据订单号查询订单
// GET /api/v1/orders/:order_no
func (h *OrderHandler) GetOrderByNo(c *gin.Context) {
	subject, ok := middleware.GetAuthSubjectFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, errors.Unauthorized("", "User not authenticated"))
		return
	}

	orderNo := c.Param("order_no")
	if orderNo == "" {
		c.JSON(http.StatusBadRequest, errors.BadRequest("", "Order number is required"))
		return
	}

	order, err := h.orderService.GetOrderByNo(c.Request.Context(), subject.UserID, orderNo)
	if err != nil {
		if err == service.ErrOrderNotFound {
			c.JSON(http.StatusNotFound, errors.NotFound("", "Order not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, errors.InternalServer("", "Failed to get order"))
		return
	}

	response.Success(c, orderToDTO(order))
}

// GetUserOrders 获取当前用户订单列表
// GET /api/v1/orders
func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	subject, ok := middleware.GetAuthSubjectFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, errors.Unauthorized("", "User not authenticated"))
		return
	}

	var req dto.GetUserOrdersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.BadRequest("", err.Error()))
		return
	}

	// 默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}

	// 查询订单
	result, err := h.orderService.GetUserOrders(c.Request.Context(), &service.GetUserOrdersInput{
		UserID: subject.UserID,
		Page:   req.Page,
		Limit:  req.Limit,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.InternalServer("", "Failed to get orders"))
		return
	}

	// 转换为DTO
	orderDTOs := make([]*dto.OrderDTO, len(result.Orders))
	for i, order := range result.Orders {
		orderDTOs[i] = orderToDTO(order)
	}

	response.Success(c, dto.GetUserOrdersResponse{
		Orders: orderDTOs,
		Total:  result.Total,
		Page:   req.Page,
		Limit:  req.Limit,
	})
}

// orderToDTO 将订单领域模型转换为DTO
func orderToDTO(order *service.Order) *dto.OrderDTO {
	if order == nil {
		return nil
	}

	orderDTO := &dto.OrderDTO{
		ID:             order.ID,
		OrderNo:        order.OrderNo,
		UserID:         order.UserID,
		ProductID:      order.ProductID,
		ProductName:    order.ProductName,
		Amount:         order.Amount,
		BonusAmount:    order.BonusAmount,
		ActualAmount:   order.ActualAmount,
		PaymentMethod:  order.PaymentMethod,
		PaymentGateway: order.PaymentGateway,
		TradeNo:        order.TradeNo,
		Status:         order.Status,
		CreatedAt:      order.CreatedAt,
		PaidAt:         order.PaidAt,
		ExpiredAt:      order.ExpiredAt,
		Notes:          order.Notes,
	}

	// 如果有关联用户信息（管理员查询）
	if order.User != nil {
		orderDTO.UserEmail = order.User.Email
	}

	return orderDTO
}

// rechargeProductToDTO 将充值套餐领域模型转换为DTO
func rechargeProductToDTO(product *service.RechargeProduct) *dto.RechargeProductDTO {
	if product == nil {
		return nil
	}

	return &dto.RechargeProductDTO{
		ID:            product.ID,
		Name:          product.Name,
		Amount:        product.Amount,
		Balance:       product.Balance,
		BonusBalance:  product.BonusBalance,
		TotalBalance:  product.GetTotalBalance(),
		Description:   product.Description,
		SortOrder:     product.SortOrder,
		IsHot:         product.IsHot,
		DiscountLabel: product.DiscountLabel,
		Status:        product.Status,
		CreatedAt:     product.CreatedAt,
		UpdatedAt:     product.UpdatedAt,
	}
}
