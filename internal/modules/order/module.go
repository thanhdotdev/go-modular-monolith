package order

import (
	orderapplication "project-example/internal/modules/order/application"
	orderhttp "project-example/internal/modules/order/delivery/http"
	orderdomain "project-example/internal/modules/order/domain"

	"github.com/gin-gonic/gin"
)

type Module struct {
	handler *orderhttp.Handler
}

func NewModule(repo orderdomain.Repository) *Module {
	usecase := orderapplication.NewService(repo)
	handler := orderhttp.NewHandler(usecase)

	return &Module{handler: handler}
}

func (m *Module) RegisterRoutes(rg *gin.RouterGroup) {
	m.handler.RegisterRoutes(rg)
}
