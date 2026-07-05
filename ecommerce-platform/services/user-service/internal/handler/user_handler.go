package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/yourname/ecommerce/user-service/internal/domain"
	"github.com/yourname/ecommerce/user-service/internal/service"
	"go.uber.org/zap"
)

type UserHandler struct {
	svc      domain.UserService
	validate *validator.Validate
	logger   *zap.Logger
}

func NewUserHandler(svc domain.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{svc: svc, validate: validator.New(), logger: logger}
}

func (h *UserHandler) Register(c *gin.Context) {
	var in domain.RegisterInput
	if err := c.ShouldBindJSON(&in); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"}); return }
	if err := h.validate.Struct(in); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed"}); return }
	resp, err := h.svc.Register(in)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyExists) { c.JSON(http.StatusConflict, gin.H{"error": "email already registered"}); return }
		h.logger.Error("register error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"}); return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *UserHandler) Login(c *gin.Context) {
	var in domain.LoginInput
	if err := c.ShouldBindJSON(&in); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"}); return }
	resp, err := h.svc.Login(in)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) { c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"}); return }
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"}); return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	id, _ := uuid.Parse(c.GetString("user_id"))
	u, err := h.svc.GetProfile(id)
	if err != nil { c.JSON(http.StatusNotFound, gin.H{"error": "not found"}); return }
	c.JSON(http.StatusOK, u)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	id, _ := uuid.Parse(c.GetString("user_id"))
	var in domain.UpdateProfileInput
	if err := c.ShouldBindJSON(&in); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"}); return }
	u, err := h.svc.UpdateProfile(id, in)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"}); return }
	c.JSON(http.StatusOK, u)
}

func (h *UserHandler) DeleteAccount(c *gin.Context) {
	id, _ := uuid.Parse(c.GetString("user_id"))
	if err := h.svc.DeleteAccount(id); err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"}); return }
	c.Status(http.StatusNoContent)
}
