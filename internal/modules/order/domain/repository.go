package orderdomain

import "context"

type Repository interface {
	FindByID(ctx context.Context, id string) (*Order, error)
}
