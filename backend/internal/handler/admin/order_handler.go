package admin

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	"github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

// OrderHandler 订单管理处理器（管理员）
type OrderHandler struct {
	orderService *service.OrderService
}

// NewOrderHandler 创建管理员订单处理器
func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// ListOrders 管理员查询订单列表
// GET /api/v1/admin/orders
func (h *OrderHandler) ListOrders(c *gin.Context) {
	var req dto.AdminListOrdersRequest
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

	// 解析日期
	var startDate, endDate *time.Time
	if req.StartDate != nil && *req.StartDate != "" {
		if t, err := time.Parse("2006-01-02", *req.StartDate); err == nil {
			startDate = &t
		}
	}
	if req.EndDate != nil && *req.EndDate != "" {
		if t, err := time.Parse("2006-01-02", *req.EndDate); err == nil {
			endOfDay := t.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			endDate = &endOfDay
		}
	}

	// 查询订单
	result, err := h.orderService.AdminListOrders(c.Request.Context(), &service.AdminListOrdersInput{
		UserID:        req.UserID,
		Status:        req.Status,
		PaymentMethod: req.PaymentMethod,
		OrderNo:       req.OrderNo,
		StartDate:     startDate,
		EndDate:       endDate,
		Page:          req.Page,
		Limit:         req.Limit,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.InternalServer("", "Failed to list orders"))
		return
	}

	// 转换为DTO
	orderDTOs := make([]*dto.OrderDTO, len(result.Orders))
	for i, order := range result.Orders {
		orderDTOs[i] = orderToDTO(order)
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orderDTOs,
		"total":  result.Total,
		"page":   req.Page,
		"limit":  req.Limit,
	})
}

// GetOrderByID 管理员查询订单详情
// GET /api/v1/admin/orders/:id
func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	idStr := c.Param("id")
	orderID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.BadRequest("", "Invalid order ID"))
		return
	}

	order, err := h.orderService.AdminGetOrderByID(c.Request.Context(), orderID)
	if err != nil {
		if err == service.ErrOrderNotFound {
			c.JSON(http.StatusNotFound, errors.NotFound("", "Order not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, errors.InternalServer("", "Failed to get order"))
		return
	}

	c.JSON(http.StatusOK, orderToDTO(order))
}

// UpdateOrderStatus 管理员更新订单状态
// PUT /api/v1/admin/orders/:id
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	idStr := c.Param("id")
	orderID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.BadRequest("", "Invalid order ID"))
		return
	}

	var req dto.AdminUpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.BadRequest("", err.Error()))
		return
	}

	// 更新订单状态
	if err := h.orderService.AdminUpdateOrderStatus(c.Request.Context(), &service.AdminUpdateOrderStatusInput{
		OrderID: orderID,
		Status:  req.Status,
		Notes:   req.Notes,
	}); err != nil {
		if err == service.ErrOrderNotFound {
			c.JSON(http.StatusNotFound, errors.NotFound("", "Order not found"))
			return
		}
		if err == service.ErrOrderStatusInvalid {
			c.JSON(http.StatusBadRequest, errors.BadRequest("", "Invalid order status"))
			return
		}
		c.JSON(http.StatusInternalServerError, errors.InternalServer("", "Failed to update order"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
}

// GetStatistics 获取订单统计信息
// GET /api/v1/admin/orders/statistics
func (h *OrderHandler) GetStatistics(c *gin.Context) {
	stats, err := h.orderService.GetOrderStatistics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.InternalServer("", "Failed to get statistics"))
		return
	}

	c.JSON(http.StatusOK, dto.OrderStatisticsResponse{
		TotalOrders:   stats.TotalOrders,
		PendingOrders: stats.PendingOrders,
		PaidOrders:    stats.PaidOrders,
		FailedOrders:  stats.FailedOrders,
		ExpiredOrders: stats.ExpiredOrders,
		TotalAmount:   stats.TotalAmount,
		TotalBalance:  stats.TotalBalance,
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

	// 如果有关联用户信息
	if order.User != nil {
		orderDTO.UserEmail = order.User.Email
	}

	return orderDTO
}
