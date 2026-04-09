package discountdomain

import "context"

type Repository interface {
	FindByCode(ctx context.Context, code string) (*Discount, error)
}
