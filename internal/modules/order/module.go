package order

import (
	discountapplication "project-example/internal/modules/discount/application"
	orderapplication "project-example/internal/modules/order/application"
	orderhttp "project-example/internal/modules/order/delivery/http"
	orderdomain "project-example/internal/modules/order/domain"

	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	Repository orderdomain.Repository
	Discounts  discountapplication.UseCase
}

type Module struct {
	handler *orderhttp.Handler
}

func NewModule(deps Dependencies) *Module {
	usecase := orderapplication.NewService(deps.Repository, deps.Discounts)
	handler := orderhttp.NewHandler(usecase)

	return &Module{handler: handler}
}

func (m *Module) RegisterRoutes(rg *gin.RouterGroup) {
	m.handler.RegisterRoutes(rg)
}
