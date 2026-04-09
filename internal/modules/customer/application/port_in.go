package customerapplication

import "context"

type UseCase interface {
	GetCustomer(ctx context.Context, id string) (*CustomerDTO, error)
}
