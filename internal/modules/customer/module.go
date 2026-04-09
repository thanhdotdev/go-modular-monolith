package customer

import (
	customerapplication "project-example/internal/modules/customer/application"
	customerhttp "project-example/internal/modules/customer/delivery/http"
	customerdomain "project-example/internal/modules/customer/domain"

	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	Repository customerdomain.Repository
}

type Module struct {
	handler *customerhttp.Handler
}

func NewModule(deps Dependencies) *Module {
	usecase := customerapplication.NewService(deps.Repository)
	handler := customerhttp.NewHandler(usecase)

	return &Module{handler: handler}
}

func (m *Module) RegisterRoutes(rg *gin.RouterGroup) {
	m.handler.RegisterRoutes(rg)
}
