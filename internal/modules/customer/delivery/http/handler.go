package customerhttp

import (
	customerapplication "project-example/internal/modules/customer/application"
	customerdomain "project-example/internal/modules/customer/domain"
	"project-example/internal/shared/httpx"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	usecase customerapplication.UseCase
}

func NewHandler(usecase customerapplication.UseCase) *Handler {
	return &Handler{usecase: usecase}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	customers := rg.Group("/customers")
	{
		customers.GET("/:id", h.get)
	}
}

func (h *Handler) get(c *gin.Context) {
	customer, err := h.usecase.GetCustomer(c.Request.Context(), c.Param("id"))
	if err != nil {
		httpx.WriteError(
			c,
			err,
			httpx.ErrorMapping{
				Err:    customerapplication.ErrInvalidCustomerID,
				Status: 400,
				Code:   "invalid_customer_id",
			},
			httpx.ErrorMapping{
				Err:    customerdomain.ErrCustomerNotFound,
				Status: 404,
				Code:   "customer_not_found",
			},
		)
		return
	}

	httpx.OK(c, customer)
}
