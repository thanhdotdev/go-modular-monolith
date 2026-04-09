package discountmemory

import (
	"context"

	discountdomain "project-example/internal/modules/discount/domain"
	"project-example/internal/shared/collection"
	"project-example/internal/shared/ptr"
)

type Repository struct {
	discounts map[string]discountdomain.Discount
}

func NewRepository(seed []discountdomain.Discount) *Repository {
	return &Repository{
		discounts: collection.IndexBy(seed, func(discount discountdomain.Discount) string {
			return discount.Code
		}),
	}
}

func SeedDiscounts() []discountdomain.Discount {
	return []discountdomain.Discount{
		{
			Code:   "WELCOME10",
			Type:   "fixed",
			Value:  10000,
			Active: true,
		},
		{
			Code:   "VIP50",
			Type:   "fixed",
			Value:  50000,
			Active: true,
		},
	}
}

func (r *Repository) FindByCode(_ context.Context, code string) (*discountdomain.Discount, error) {
	discount, ok := r.discounts[code]
	if !ok {
		return nil, discountdomain.ErrDiscountNotFound
	}

	return ptr.Of(discount), nil
}
