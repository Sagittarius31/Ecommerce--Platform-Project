package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourname/ecommerce/order-service/internal/domain"
	"github.com/yourname/ecommerce/order-service/internal/service"
	"go.uber.org/zap"
)

type OrderHandler struct {
	svc    *service.OrderService
	logger *zap.Logger
}

func NewOrderHandler(svc *service.OrderService, logger *zap.Logger) *OrderHandler {
	return &OrderHandler{svc: svc, logger: logger}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var in domain.CreateOrderInput
	if err := c.ShouldBindJSON(&in); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"}); return }
	o, err := h.svc.CreateOrder(c.Request.Context(), in)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusCreated, o)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"}); return }
	o, err := h.svc.GetOrder(c.Request.Context(), id)
	if err != nil { c.JSON(http.StatusNotFound, gin.H{"error": "not found"}); return }
	c.JSON(http.StatusOK, o)
}

func (h *OrderHandler) ListOrders(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "list orders"})
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"}); return }
	if err := h.svc.CancelOrder(c.Request.Context(), id); err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "cancel failed"}); return }
	c.Status(http.StatusNoContent)
}
