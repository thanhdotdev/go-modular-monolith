package discounthttp

import (
	discountapplication "project-example/internal/modules/discount/application"
	discountdomain "project-example/internal/modules/discount/domain"
	"project-example/internal/shared/httpx"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	usecase discountapplication.UseCase
}

func NewHandler(usecase discountapplication.UseCase) *Handler {
	return &Handler{usecase: usecase}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	discounts := rg.Group("/discounts")
	{
		discounts.GET("/:code", h.get)
	}
}

func (h *Handler) get(c *gin.Context) {
	discount, err := h.usecase.GetDiscount(c.Request.Context(), c.Param("code"))
	if err != nil {
		httpx.WriteError(
			c,
			err,
			httpx.ErrorMapping{
				Err:    discountapplication.ErrInvalidDiscountCode,
				Status: 400,
				Code:   "invalid_discount_code",
			},
			httpx.ErrorMapping{
				Err:    discountdomain.ErrDiscountNotFound,
				Status: 404,
				Code:   "discount_not_found",
			},
		)
		return
	}

	httpx.OK(c, discount)
}
