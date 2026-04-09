package customerdomain

import "context"

type Repository interface {
	FindByID(ctx context.Context, id string) (*Customer, error)
}
