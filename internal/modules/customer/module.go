package customer

import (
	customerapplication "project-example/internal/modules/customer/application"
	customerhttp "project-example/internal/modules/customer/delivery/http"
	customerdomain "project-example/internal/modules/customer/domain"

	"github.com/gin-gonic/gin"
)

type Module struct {
	handler *customerhttp.Handler
}

func NewModule(repo customerdomain.Repository) *Module {
	usecase := customerapplication.NewService(repo)
	handler := customerhttp.NewHandler(usecase)

	return &Module{handler: handler}
}

func (m *Module) RegisterRoutes(rg *gin.RouterGroup) {
	m.handler.RegisterRoutes(rg)
}
