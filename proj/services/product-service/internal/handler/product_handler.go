package handler
import ("errors"; "net/http"; "strconv"; "github.com/gin-gonic/gin"; "github.com/google/uuid"; "github.com/yourname/ecommerce/product-service/internal/domain"; "github.com/yourname/ecommerce/product-service/internal/service"; "go.uber.org/zap")
type ProductHandler struct { svc *service.ProductService; logger *zap.Logger }
func NewProductHandler(svc *service.ProductService, logger *zap.Logger) *ProductHandler { return &ProductHandler{svc: svc, logger: logger} }
func (h *ProductHandler) ListProducts(c *gin.Context) {
	f := domain.ProductFilter{Page: 1, PageSize: 20}
	if p := c.Query("page"); p != "" { if v, e := strconv.Atoi(p); e == nil { f.Page = v } }
	products, total, err := h.svc.List(c.Request.Context(), f)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"}); return }
	c.JSON(http.StatusOK, gin.H{"data": products, "total": total})
}
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id")); if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"}); return }
	p, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil { if errors.Is(err, domain.ErrNotFound) { c.JSON(http.StatusNotFound, gin.H{"error": "not found"}); return }; c.JSON(http.StatusInternalServerError, gin.H{"error": "error"}); return }
	c.JSON(http.StatusOK, p)
}
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var in domain.CreateProductInput
	if err := c.ShouldBindJSON(&in); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"}); return }
	p, err := h.svc.Create(c.Request.Context(), in); if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"}); return }
	c.JSON(http.StatusCreated, p)
}
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id")); if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"}); return }
	var in domain.CreateProductInput
	if err := c.ShouldBindJSON(&in); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"}); return }
	p, err := h.svc.Update(c.Request.Context(), id, in); if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"}); return }
	c.JSON(http.StatusOK, p)
}
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id")); if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"}); return }
	if err := h.svc.Delete(c.Request.Context(), id); err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"}); return }
	c.Status(http.StatusNoContent)
}
