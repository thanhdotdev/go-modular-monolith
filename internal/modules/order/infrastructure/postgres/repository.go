package orderpostgres

import (
	"context"
	"errors"
	orderdomain "project-example/internal/modules/order/domain"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByID(ctx context.Context, id string) (*orderdomain.Order, error) {
	var model OrderModel
	err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, orderdomain.ErrOrderNotFound
	case err != nil:
		return nil, err
	default:
		return toDomain(model), nil
	}
}
