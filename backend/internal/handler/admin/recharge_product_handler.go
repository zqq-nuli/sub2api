package admin

import (
	"net/http"
	"strconv"

	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	"github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

// RechargeProductHandler 充值套餐管理处理器（管理员）
type RechargeProductHandler struct {
	productService *service.RechargeProductService
}

// NewRechargeProductHandler 创建充值套餐管理处理器
func NewRechargeProductHandler(productService *service.RechargeProductService) *RechargeProductHandler {
	return &RechargeProductHandler{
		productService: productService,
	}
}

// ListProducts 获取所有充值套餐
// GET /api/v1/admin/recharge-products
func (h *RechargeProductHandler) ListProducts(c *gin.Context) {
	products, err := h.productService.ListAllProducts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.InternalServer("", "Failed to list products"))
		return
	}

	// 转换为DTO
	productDTOs := make([]*dto.RechargeProductDTO, len(products))
	for i, p := range products {
		productDTOs[i] = rechargeProductToDTO(p)
	}

	c.JSON(http.StatusOK, productDTOs)
}

// GetProductByID 获取充值套餐详情
// GET /api/v1/admin/recharge-products/:id
func (h *RechargeProductHandler) GetProductByID(c *gin.Context) {
	idStr := c.Param("id")
	productID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.BadRequest("", "Invalid product ID"))
		return
	}

	product, err := h.productService.GetProductByID(c.Request.Context(), productID)
	if err != nil {
		if err == service.ErrProductNotFound {
			c.JSON(http.StatusNotFound, errors.NotFound("", "Product not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, errors.InternalServer("", "Failed to get product"))
		return
	}

	c.JSON(http.StatusOK, rechargeProductToDTO(product))
}

// CreateProduct 创建充值套餐
// POST /api/v1/admin/recharge-products
func (h *RechargeProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateRechargeProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.BadRequest("", err.Error()))
		return
	}

	product, err := h.productService.CreateProduct(c.Request.Context(), &service.CreateProductInput{
		Name:          req.Name,
		Amount:        req.Amount,
		Balance:       req.Balance,
		BonusBalance:  req.BonusBalance,
		Description:   req.Description,
		SortOrder:     req.SortOrder,
		IsHot:         req.IsHot,
		DiscountLabel: req.DiscountLabel,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.InternalServer("", "Failed to create product"))
		return
	}

	c.JSON(http.StatusCreated, rechargeProductToDTO(product))
}

// UpdateProduct 更新充值套餐
// PUT /api/v1/admin/recharge-products/:id
func (h *RechargeProductHandler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	productID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.BadRequest("", "Invalid product ID"))
		return
	}

	var req dto.UpdateRechargeProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.BadRequest("", err.Error()))
		return
	}

	product, err := h.productService.UpdateProduct(c.Request.Context(), &service.UpdateProductInput{
		ID:            productID,
		Name:          req.Name,
		Amount:        req.Amount,
		Balance:       req.Balance,
		BonusBalance:  req.BonusBalance,
		Description:   req.Description,
		SortOrder:     req.SortOrder,
		IsHot:         req.IsHot,
		DiscountLabel: req.DiscountLabel,
		Status:        req.Status,
	})
	if err != nil {
		if err == service.ErrProductNotFound {
			c.JSON(http.StatusNotFound, errors.NotFound("", "Product not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, errors.InternalServer("", "Failed to update product"))
		return
	}

	c.JSON(http.StatusOK, rechargeProductToDTO(product))
}

// DeleteProduct 删除充值套餐
// DELETE /api/v1/admin/recharge-products/:id
func (h *RechargeProductHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	productID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.BadRequest("", "Invalid product ID"))
		return
	}

	if err := h.productService.DeleteProduct(c.Request.Context(), productID); err != nil {
		if err == service.ErrProductNotFound {
			c.JSON(http.StatusNotFound, errors.NotFound("", "Product not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, errors.InternalServer("", "Failed to delete product"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
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
