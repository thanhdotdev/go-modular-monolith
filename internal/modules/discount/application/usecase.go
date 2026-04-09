package discountapplication

import "context"

type UseCase interface {
	GetDiscount(ctx context.Context, code string) (*DiscountDTO, error)
}
