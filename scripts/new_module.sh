#!/usr/bin/env bash

set -euo pipefail

if [ "${1-}" = "" ]; then
  echo "usage: $0 <module_name>"
  exit 1
fi

MODULE_NAME="$1"

if ! [[ "$MODULE_NAME" =~ ^[a-z][a-z0-9_]*$ ]]; then
  echo "module name must match: ^[a-z][a-z0-9_]*$"
  exit 1
fi

MODULE_DIR="internal/modules/${MODULE_NAME}"

if [ -d "$MODULE_DIR" ]; then
  echo "module already exists: $MODULE_DIR"
  exit 1
fi

PACKAGE_PREFIX="${MODULE_NAME}"
TYPE_NAME="$(tr '[:lower:]' '[:upper:]' <<< "${MODULE_NAME:0:1}")${MODULE_NAME:1}"

mkdir -p \
  "$MODULE_DIR/domain" \
  "$MODULE_DIR/application" \
  "$MODULE_DIR/delivery/http" \
  "$MODULE_DIR/infrastructure/memory"

cat <<EOF > "$MODULE_DIR/module.go"
package ${PACKAGE_PREFIX}

import (
	${PACKAGE_PREFIX}application "project-example/internal/modules/${MODULE_NAME}/application"
	${PACKAGE_PREFIX}http "project-example/internal/modules/${MODULE_NAME}/delivery/http"
	${PACKAGE_PREFIX}domain "project-example/internal/modules/${MODULE_NAME}/domain"

	"github.com/gin-gonic/gin"
)

type Module struct {
	handler *${PACKAGE_PREFIX}http.Handler
}

type Dependencies struct {
	Repository ${PACKAGE_PREFIX}domain.Repository
}

func NewModule(deps Dependencies) *Module {
	usecase := ${PACKAGE_PREFIX}application.NewService(deps.Repository)
	handler := ${PACKAGE_PREFIX}http.NewHandler(usecase)

	return &Module{handler: handler}
}

func (m *Module) RegisterRoutes(rg *gin.RouterGroup) {
	m.handler.RegisterRoutes(rg)
}
EOF

cat <<EOF > "$MODULE_DIR/domain/${MODULE_NAME}.go"
package ${PACKAGE_PREFIX}domain

type ${TYPE_NAME} struct {
	ID string
}
EOF

cat <<EOF > "$MODULE_DIR/domain/errors.go"
package ${PACKAGE_PREFIX}domain

import "errors"

var Err${TYPE_NAME}NotFound = errors.New("${MODULE_NAME} not found")
EOF

cat <<EOF > "$MODULE_DIR/domain/repository.go"
package ${PACKAGE_PREFIX}domain

import "context"

type Repository interface {
	FindByID(ctx context.Context, id string) (*${TYPE_NAME}, error)
}
EOF

cat <<EOF > "$MODULE_DIR/application/dto.go"
package ${PACKAGE_PREFIX}application

type ${TYPE_NAME}DTO struct {
	ID string \`json:"id"\`
}
EOF

cat <<EOF > "$MODULE_DIR/application/usecase.go"
package ${PACKAGE_PREFIX}application

import "context"

type UseCase interface {
	Get${TYPE_NAME}(ctx context.Context, id string) (*${TYPE_NAME}DTO, error)
}
EOF

cat <<EOF > "$MODULE_DIR/application/service.go"
package ${PACKAGE_PREFIX}application

import (
	"context"
	"errors"
	${PACKAGE_PREFIX}domain "project-example/internal/modules/${MODULE_NAME}/domain"
	"strings"
)

var ErrInvalid${TYPE_NAME}ID = errors.New("${MODULE_NAME} id is required")

type Service struct {
	repo ${PACKAGE_PREFIX}domain.Repository
}

func NewService(repo ${PACKAGE_PREFIX}domain.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Get${TYPE_NAME}(ctx context.Context, id string) (*${TYPE_NAME}DTO, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, ErrInvalid${TYPE_NAME}ID
	}

	entity, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &${TYPE_NAME}DTO{
		ID: entity.ID,
	}, nil
}
EOF

cat <<EOF > "$MODULE_DIR/application/service_test.go"
package ${PACKAGE_PREFIX}application

import (
	"context"
	"errors"
	${PACKAGE_PREFIX}domain "project-example/internal/modules/${MODULE_NAME}/domain"
	"testing"
)

type fakeRepository struct {
	entity *${PACKAGE_PREFIX}domain.${TYPE_NAME}
	err    error
}

func (f fakeRepository) FindByID(_ context.Context, _ string) (*${PACKAGE_PREFIX}domain.${TYPE_NAME}, error) {
	if f.err != nil {
		return nil, f.err
	}

	return f.entity, nil
}

func TestGet${TYPE_NAME}(t *testing.T) {
	service := NewService(fakeRepository{
		entity: &${PACKAGE_PREFIX}domain.${TYPE_NAME}{
			ID: "${MODULE_NAME}-001",
		},
	})

	got, err := service.Get${TYPE_NAME}(context.Background(), "${MODULE_NAME}-001")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if got.ID != "${MODULE_NAME}-001" {
		t.Fatalf("expected ${MODULE_NAME} id ${MODULE_NAME}-001, got %s", got.ID)
	}
}

func TestGet${TYPE_NAME}RequiresID(t *testing.T) {
	service := NewService(fakeRepository{})

	_, err := service.Get${TYPE_NAME}(context.Background(), " ")
	if !errors.Is(err, ErrInvalid${TYPE_NAME}ID) {
		t.Fatalf("expected ErrInvalid${TYPE_NAME}ID, got %v", err)
	}
}
EOF

cat <<EOF > "$MODULE_DIR/delivery/http/handler.go"
package ${PACKAGE_PREFIX}http

import (
	${PACKAGE_PREFIX}application "project-example/internal/modules/${MODULE_NAME}/application"
	${PACKAGE_PREFIX}domain "project-example/internal/modules/${MODULE_NAME}/domain"
	"project-example/internal/shared/httpx"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	usecase ${PACKAGE_PREFIX}application.UseCase
}

func NewHandler(usecase ${PACKAGE_PREFIX}application.UseCase) *Handler {
	return &Handler{usecase: usecase}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/${MODULE_NAME}s")
	{
		group.GET("/:id", h.get)
	}
}

func (h *Handler) get(c *gin.Context) {
	entity, err := h.usecase.Get${TYPE_NAME}(c.Request.Context(), c.Param("id"))
	if err != nil {
		httpx.WriteError(
			c,
			err,
			httpx.ErrorMapping{
				Err:    ${PACKAGE_PREFIX}application.ErrInvalid${TYPE_NAME}ID,
				Status: 400,
				Code:   "invalid_${MODULE_NAME}_id",
			},
			httpx.ErrorMapping{
				Err:    ${PACKAGE_PREFIX}domain.Err${TYPE_NAME}NotFound,
				Status: 404,
				Code:   "${MODULE_NAME}_not_found",
			},
		)
		return
	}

	httpx.OK(c, entity)
}
EOF

cat <<EOF > "$MODULE_DIR/delivery/http/handler_test.go"
package ${PACKAGE_PREFIX}http

import (
	"context"
	"net/http"
	"net/http/httptest"
	${PACKAGE_PREFIX}application "project-example/internal/modules/${MODULE_NAME}/application"
	${PACKAGE_PREFIX}domain "project-example/internal/modules/${MODULE_NAME}/domain"
	"testing"

	"github.com/gin-gonic/gin"
)

type stubUseCase struct {
	entity *${PACKAGE_PREFIX}application.${TYPE_NAME}DTO
	err    error
}

func (s stubUseCase) Get${TYPE_NAME}(_ context.Context, _ string) (*${PACKAGE_PREFIX}application.${TYPE_NAME}DTO, error) {
	if s.err != nil {
		return nil, s.err
	}

	return s.entity, nil
}

func TestGet${TYPE_NAME}ReturnsNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	handler := NewHandler(stubUseCase{err: ${PACKAGE_PREFIX}domain.Err${TYPE_NAME}NotFound})
	handler.RegisterRoutes(router.Group("/api/v1"))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/${MODULE_NAME}s/${MODULE_NAME}-404", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, recorder.Code)
	}
}
EOF

cat <<EOF > "$MODULE_DIR/infrastructure/memory/repository.go"
package ${MODULE_NAME}memory

import (
	"context"
	${PACKAGE_PREFIX}domain "project-example/internal/modules/${MODULE_NAME}/domain"
	"project-example/internal/shared/collection"
	"project-example/internal/shared/ptr"
)

type Repository struct {
	entities map[string]${PACKAGE_PREFIX}domain.${TYPE_NAME}
}

func NewRepository(seed []${PACKAGE_PREFIX}domain.${TYPE_NAME}) *Repository {
	return &Repository{
		entities: collection.IndexBy(seed, func(entity ${PACKAGE_PREFIX}domain.${TYPE_NAME}) string {
			return entity.ID
		}),
	}
}

func Seed${TYPE_NAME}s() []${PACKAGE_PREFIX}domain.${TYPE_NAME} {
	return []${PACKAGE_PREFIX}domain.${TYPE_NAME}{
		{
			ID: "${MODULE_NAME}-001",
		},
	}
}

func (r *Repository) FindByID(_ context.Context, id string) (*${PACKAGE_PREFIX}domain.${TYPE_NAME}, error) {
	entity, ok := r.entities[id]
	if !ok {
		return nil, ${PACKAGE_PREFIX}domain.Err${TYPE_NAME}NotFound
	}

	return ptr.Of(entity), nil
}
EOF

gofmt -w "$MODULE_DIR"

echo "created module skeleton at $MODULE_DIR"
echo "next steps:"
echo "  1. register the module in internal/app/app.go"
echo "  2. fill in entity fields and use case logic"
echo "  3. replace memory repository when you need real persistence"
