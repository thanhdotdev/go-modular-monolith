package orderhttp

import (
	orderapplication "project-example/internal/modules/order/application"
	orderdomain "project-example/internal/modules/order/domain"
	"project-example/internal/shared/httpx"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	usecase orderapplication.UseCase
}

func NewHandler(usecase orderapplication.UseCase) *Handler {
	return &Handler{usecase: usecase}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	orders := rg.Group("/orders")
	{
		orders.GET("/:id", h.get)
	}
}

type GetOrderRequest struct {
	ID string
}

func (h *Handler) get(c *gin.Context) {
	order, err := h.usecase.GetOrder(c.Request.Context(), c.Param("id"))
	if err != nil {
		httpx.WriteError(
			c,
			err,
			httpx.ErrorMapping{
				Err:    orderapplication.ErrInvalidOrderID,
				Status: 400,
				Code:   "invalid_order_id",
			},
			httpx.ErrorMapping{
				Err:    orderdomain.ErrOrderNotFound,
				Status: 404,
				Code:   "order_not_found",
			},
		)
		return
	}

	httpx.OK(c, order)
}
