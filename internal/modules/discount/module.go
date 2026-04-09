package discount

import (
	discountapplication "project-example/internal/modules/discount/application"
	discounthttp "project-example/internal/modules/discount/delivery/http"
	discountdomain "project-example/internal/modules/discount/domain"

	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	Repository discountdomain.Repository
}

type Module struct {
	usecase discountapplication.UseCase
	handler *discounthttp.Handler
}

func NewModule(deps Dependencies) *Module {
	usecase := discountapplication.NewService(deps.Repository)
	handler := discounthttp.NewHandler(usecase)

	return &Module{
		usecase: usecase,
		handler: handler,
	}
}

func (m *Module) UseCase() discountapplication.UseCase {
	return m.usecase
}

func (m *Module) RegisterRoutes(rg *gin.RouterGroup) {
	m.handler.RegisterRoutes(rg)
}
